package Telegram

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// mention_user_cmd if a telegram user is mentioned it will send a message to discord also mentioning the corresponding user
func mention_user_cmd(client *TelegramClient, update tgbotapi.Update, message_fields []string) {
	channel_id := client.DiscordClient.GetDiscordChannelID(update.Message.Chat.Title)
	if channel_id != "" {
		message_text := client.replaceMentionT2D(strings.Join(message_fields[1:], " "))
		discord_message := fmt.Sprintf("%s says: %s", strings.ToTitle(update.Message.From.FirstName), message_text)
		client.DiscordClient.Session.ChannelMessageSend(channel_id, discord_message)
	} else {
		log.Println("No text channel created for:", update.Message.Chat.Title)
	}
}
