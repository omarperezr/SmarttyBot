package main

import (
	conf "github.com/omarperezr/SmarttyBot/Core/Config"
	dbot "github.com/omarperezr/SmarttyBot/Platforms/Discord"
	mbot "github.com/omarperezr/SmarttyBot/Platforms/Email"
	tbot "github.com/omarperezr/SmarttyBot/Platforms/Telegram"
	wclient "github.com/omarperezr/SmarttyBot/Platforms/Wit"
)

func main() {
	config := conf.LoadConfig(".env")
	config.SetUpGobFiles()

	email_client := mbot.SetUp(config)
	email_client.Init()

	witai_client := wclient.SetUp(config)
	witai_client.Init()

	telegram_client := tbot.SetUp(config, &witai_client)
	telegram_client.Init()

	a_discord_client := dbot.SetUp(config)
	a_discord_client.Init()

	a_discord_client.TelegramClient = telegram_client.Api
	a_discord_client.SendEmail = email_client.SendEmail
	telegram_client.DiscordClient = &a_discord_client
	email_client.DiscordClient = &a_discord_client

	// Register the newDiscordMessage func as a callback for MessageCreate events.
	go telegram_client.StartListening()
	go email_client.StartListening()
	a_discord_client.StartListening()
}
