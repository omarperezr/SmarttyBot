package Config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	SMTP_Server    string
	SMTP_Port      int
	IMAP_Server    string
	IMAP_Port      int
	Email_Account  string
	Email_Password string

	WIT_Api_Key      string
	Telegram_Api_Key string
	Discord_Api_Key  string
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

	return &conf
}
