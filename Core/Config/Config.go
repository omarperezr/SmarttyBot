package Config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/omarperezr/SmarttyBot/Utils"
)

var telegram_ids = make(map[string]int64)
var discord_mentions = make(map[string]string)
var telegram_mentions = make(map[string]string)
var from_list = make(map[string]string)
var callback_map = make(map[string]interface{})
var tele_wit_chan = make(chan string)

type Config struct {
	SMTP_Server      string
	SMTP_Port        int
	IMAP_Server      string
	IMAP_Port        int
	Email_Account    string
	Email_Password   string
	WIT_Api_Key      string
	Telegram_Api_Key string
	Discord_Api_Key  string

	Telegram_ids      *map[string]int64
	Discord_mentions  *map[string]string
	Telegram_mentions *map[string]string
	From_list         *map[string]string
	Callback_map      *map[string]interface{}
	Tele_wit_chan     chan string
}

func LoadConfig(filename string) *Config {
	godotenv.Load(filename)

	var conf Config

	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	imapPort, _ := strconv.Atoi(os.Getenv("IMAP_PORT"))
	conf.SMTP_Port = smtpPort
	conf.IMAP_Port = imapPort
	conf.SMTP_Server = os.Getenv("SMTP_SERVER")
	conf.IMAP_Server = os.Getenv("IMAP_SERVER")
	conf.Email_Account = os.Getenv("EMAIL_ACCOUNT")
	conf.Email_Password = os.Getenv("EMAIL_PASWORD")

	conf.WIT_Api_Key = os.Getenv("WIT_API_KEY")
	conf.Telegram_Api_Key = os.Getenv("TELEGRAM_API_KEY")
	conf.Discord_Api_Key = os.Getenv("DISCORD_API_KEY")

	conf.Callback_map = &callback_map
	conf.Tele_wit_chan = tele_wit_chan

	return &conf
}

func (config *Config) SetUpGobFiles() {
	if Utils.FileExists("data/telegram_ids.gob") {
		Utils.DeserializeObject(&telegram_ids, "data/telegram_ids.gob")
	}
	if Utils.FileExists("data/discord_mentions.gob") {
		Utils.DeserializeObject(&discord_mentions, "data/discord_mentions.gob")
	}
	if Utils.FileExists("data/telegram_mentions.gob") {
		Utils.DeserializeObject(&telegram_mentions, "data/telegram_mentions.gob")
	}
	if Utils.FileExists("data/from_list.gob") {
		Utils.DeserializeObject(&from_list, "data/from_list.gob")
	}
	config.Telegram_ids = &telegram_ids
	config.Discord_mentions = &discord_mentions
	config.Telegram_mentions = &telegram_mentions
	config.From_list = &from_list
}
