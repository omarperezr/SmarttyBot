package Email

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"gopkg.in/gomail.v2"

	dClient "github.com/omarperezr/SmarttyBot/Platforms/Discord"

	"github.com/emersion/go-imap"
	imapClient "github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

type EmailClient struct {
	Email         string
	Password      string
	SMTPServer    string
	SMTPPort      int
	IMAPServer    string
	IMAPPort      int
	Dialer        *gomail.Dialer
	From_list     *map[string]string
	DiscordClient *dClient.DiscordClient
}

func (client *EmailClient) Init() {
	// Connects to SMTP to send messages
	client.Dialer = gomail.NewDialer(client.SMTPServer, client.SMTPPort, client.Email, client.Password)
}

// ConnectIMAO connects to the specified IMAP server
func (client *EmailClient) ConnectIMAP() *imapClient.Client {
	// Connect to server
	c, err := imapClient.DialTLS(fmt.Sprintf("%s:%d", client.IMAPServer, client.IMAPPort), nil)
	if err != nil {
		log.Println(err, 1)
	}
	return c
}

// SendEmail sends an email
func (client *EmailClient) SendEmail(to, subject, body string) {
	m := gomail.NewMessage()
	m.SetHeader("From", client.Email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	if err := client.Dialer.DialAndSend(m); err != nil {
		log.Panic(err)
	}
}

// CheckSender checks if the email sender is in the from list
func (client *EmailClient) CheckSender(address_list []*mail.Address) (string, string) {
	for _, f := range address_list {
		if val, ok := (*client.From_list)[f.Address]; ok {
			return val, f.Address
		}
	}
	return "", ""
}

// CheckMessagesForDiscord checks if there are any email messages that need to be sent to discord channel
func (client *EmailClient) CheckMessagesForDiscord(messages chan *imap.Message) {
	var channel string
	var subject string
	var sender string

	// Unseen messages
	section := &imap.BodySectionName{}
	//section.Specifier = imap.TextSpecifier
	for msg := range messages {
		r := msg.GetBody(section)
		if r == nil {
			log.Println("Server didn't return message body")
		}
		// Create a new mail reader
		mr, err := mail.CreateReader(r)
		if err != nil {
			log.Panic(err)
		}
		header := mr.Header
		if from, err := header.AddressList("From"); err == nil {
			if channel, sender = client.CheckSender(from); channel == "" {
				log.Println("SENDER NOT REGISTERED")
				return
			}
		}
		if subject, err = header.Subject(); err != nil {
			log.Panic("Error getting subject", err)
		}
		// Process each message's part
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Panic(err)
			}

			switch p.Header.(type) {
			case *mail.InlineHeader:
				if strings.Split(p.Header.Get("Content-Type"), ";")[0] == "text/plain" {
					// This is the message's text (can be plain-text or HTML)
					b, _ := ioutil.ReadAll(p.Body)
					mailToDiscord := fmt.Sprintf("Subject: %s\nMessage: %s", subject, string(b))

					channel_id := client.DiscordClient.GetDiscordChannelID(channel)
					if channel_id != "" {
						discord_message := fmt.Sprintf("%s says:\n%s", sender, mailToDiscord)
						client.DiscordClient.Session.ChannelMessageSend(channel_id, discord_message)
					} else {
						log.Println("No text channel created for:", channel)
					}
				}

				// case *mail.AttachmentHeader:
				// // This is an attachment
				// 	filename, _ := h.Filename()
				// 	log.Println("Got attachment: %v", filename)
			}
		}
	}
}

// ReadUnseen reads all messages from INBOX and checks if there are any unread messages
func (client *EmailClient) ReadUnseen() {
	c := client.ConnectIMAP()
	defer c.Logout()

	// Login
	if err := c.Login(client.Email, client.Password); err != nil {
		log.Println(err)
	}

	// Select INBOX
	_, err := c.Select("INBOX", false)
	if err != nil {
		log.Println(err)
	}

	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{"\\Seen"}
	uids, err := c.Search(criteria)
	if err != nil {
		log.Println(err)
	}

	// If there are unread messages
	if len(uids) > 0 {
		client.RetrieveMessagesData(uids, c)
	}
}

// RetrieveMessagesData retrieves data from email messages
func (client *EmailClient) RetrieveMessagesData(uids []uint32, c *imapClient.Client) {
	seqset := new(imap.SeqSet)
	seqset.AddNum(uids...)
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchInternalDate, section.FetchItem()}
	messages := make(chan *imap.Message)
	go func() {
		if err := c.Fetch(seqset, items, messages); err != nil {
			log.Println("Error fetching messages: ", err)
		}
	}()

	client.CheckMessagesForDiscord(messages)
}

func (client *EmailClient) StartListening() {
	for {
		client.ReadUnseen()
		time.Sleep(200)
	}
}
