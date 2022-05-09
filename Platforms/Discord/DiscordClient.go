package Discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/omarperezr/SmarttyBot/Utils"

	"github.com/bwmarrin/discordgo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type DiscordClient struct {
	ApiKey          string
	Session         *discordgo.Session
	TelegramIDs     *map[string]int64
	DiscordMentions *map[string]string
	TelegramClient  *tgbotapi.BotAPI
	Email_from_list *map[string]string
	SendEmail       func(string, string, string)
}

func (client *DiscordClient) Init() {
	// Create a new Discord session using the provided bot token.
	var discord_err error
	client.Session, discord_err = discordgo.New("Bot " + client.ApiKey)
	if discord_err != nil {
		log.Panic("Error creating session,", discord_err)
	}

	client.Session.AddHandler(client.messageHandler)

	// In this example, we only care about receiving message events.
	client.Session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
}

func (client *DiscordClient) StartListening() {
	// Open a websocket connection to Discord and begin listening.
	discord_err := client.Session.Open()
	if discord_err != nil {
		log.Panic("Error opening connection, ", discord_err)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Client is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc

	// Cleanly close down the Discord session.
	client.Session.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func (client *DiscordClient) messageHandler(discordSession *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == discordSession.State.User.ID {
		return
	}

	// If the message is "-a message" send message to discord
	message_parameters := strings.Fields(m.Content)
	amount_parameters := len(message_parameters)
	if amount_parameters == 1 && message_parameters[0] == "-habla" {
		client.SendChannelMessage(discordSession, m.ChannelID, "Omar es cul")

	} else if amount_parameters > 1 && (message_parameters[0] == "-a" || message_parameters[0] == "-answer") {
		discord_channel, err := discordSession.Channel(m.ChannelID)
		if err != nil {
			log.Panic("Error getting discord channel:", err)
		}

		telegram_message := fmt.Sprintf("%s says: %s", strings.ToLower(m.Author.Username), strings.Join(message_parameters[1:], " "))
		telegram_id := (*client.TelegramIDs)[discord_channel.Name]
		msg := tgbotapi.NewMessage(telegram_id, telegram_message)
		// msg.ReplyToMessageID = update.Message.MessageID
		if _, err := client.TelegramClient.Send(msg); err != nil {
			log.Panic("Error sending message to telegram:", err)
		}

	} else if amount_parameters == 2 && message_parameters[0] == "-r" {
		discord_user_mention := fmt.Sprintf("<@%s>", m.Author.ID)
		(*client.DiscordMentions)[strings.ToUpper(message_parameters[1])] = discord_user_mention
		Utils.SerializeObject(client.DiscordMentions, "data/discord_mentions.gob")

		message := fmt.Sprintf("%s was registered correctly", message_parameters[1])
		client.SendChannelMessage(discordSession, m.ChannelID, message)

	} else if amount_parameters == 2 && message_parameters[0] == "-rm" {
		discord_channel, _ := discordSession.Channel(m.ChannelID)
		(*client.Email_from_list)[message_parameters[1]] = discord_channel.Name
		Utils.SerializeObject(*client.Email_from_list, "data/from_list.gob")
		message := fmt.Sprintf("%s was registered correctly for this channel (%s)", message_parameters[1], discord_channel.Name)
		client.SendChannelMessage(discordSession, m.ChannelID, message)

	} else if amount_parameters > 2 && (message_parameters[0] == "-m" || message_parameters[0] == "-send-email") {
		to := message_parameters[1]
		subject := message_parameters[2]
		body := strings.Join(message_parameters[3:], " ")
		client.SendEmail(to, subject, body)
		message := fmt.Sprintf("Email sent to %s!", to)
		client.SendChannelMessage(discordSession, m.ChannelID, message)
	}
}

// SendChannelMessage sends a message to a channel using the channel ID
func (client *DiscordClient) SendChannelMessage(discordSession *discordgo.Session, id_for_channel, message string) {
	discord_channel, err := discordSession.Channel(id_for_channel)
	if err != nil {
		log.Panic("Error getting discord channel:", err)
	}
	channel_id := client.GetDiscordChannelID(discord_channel.Name)
	client.Session.ChannelMessageSend(channel_id, message)
}

// GetDiscordChannelID gets the channel ID based on the channel name
func (client *DiscordClient) GetDiscordChannelID(channel_name string) string {
	for _, guild := range client.Session.State.Guilds {

		// Get channels for this guild
		channels, _ := client.Session.GuildChannels(guild.ID)

		for _, c := range channels {
			if c.Name == channel_name {
				return c.ID
			}
		}
	}
	return ""
}
