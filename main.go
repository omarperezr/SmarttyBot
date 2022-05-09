package main

import (
	conf "github.com/omarperezr/SmarttyBot/Core/Config"
	dbot "github.com/omarperezr/SmarttyBot/Platforms/Discord"
	mbot "github.com/omarperezr/SmarttyBot/Platforms/Email"
	tbot "github.com/omarperezr/SmarttyBot/Platforms/Telegram"
	wclient "github.com/omarperezr/SmarttyBot/Platforms/Wit"
	"github.com/omarperezr/SmarttyBot/Utils"
)

var telegram_ids = make(map[string]int64)
var discord_mentions = make(map[string]string)
var telegram_mentions = make(map[string]string)
var from_list = make(map[string]string)

func main() {
	config := conf.LoadConfig(".env")
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

	email_client := mbot.EmailClient{
		Email:      config.Email_Account,
		Password:   config.Email_Password,
		SMTPServer: config.SMTP_Server,
		SMTPPort:   config.SMTP_Port,
		IMAPServer: config.IMAP_Server,
		IMAPPort:   config.IMAP_Port,
		From_list:  &from_list,
	}
	email_client.Init()

	callback_map := make(map[string]interface{})
	tele_wit_chan := make(chan string)
	a_witai_client := wclient.WitClient{
		ApiKey:               config.WIT_Api_Key,
		TeleWitChan:          tele_wit_chan,
		AwaitingCallbackData: callback_map,
	}
	a_witai_client.Init()

	a_telegram_client := tbot.TelegramClient{
		ApiKey:               config.Telegram_Api_Key,
		TelegramIDs:          &telegram_ids,
		TelegramMentions:     &telegram_mentions,
		DiscordMentions:      &discord_mentions,
		WitClient:            &a_witai_client,
		TeleWitChan:          tele_wit_chan,
		AwaitingCallbackData: callback_map,
	}
	a_telegram_client.Init()

	a_discord_client := dbot.DiscordClient{
		ApiKey:          config.Discord_Api_Key,
		TelegramIDs:     &telegram_ids,
		DiscordMentions: &discord_mentions,
		Email_from_list: &from_list,
	}
	a_discord_client.Init()

	a_discord_client.TelegramClient = a_telegram_client.Api
	a_discord_client.SendEmail = email_client.SendEmail
	a_telegram_client.DiscordClient = &a_discord_client
	email_client.DiscordClient = &a_discord_client

	// Register the newDiscordMessage func as a callback for MessageCreate events.
	go a_telegram_client.StartListening()
	go email_client.StartListening()
	a_discord_client.StartListening()
}
