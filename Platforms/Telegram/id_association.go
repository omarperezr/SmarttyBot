package Telegram

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/omarperezr/SmarttyBot/Utils"
)

// associate_existing_group if the bot is added to a new group this function associates the name of the chat to the id
func associate_existing_group(client *TelegramClient, update tgbotapi.Update) {
	for _, user := range *(update.Message.NewChatMembers) {
		if user.UserName == "smartty_bridge_bot" {
			(*client.TelegramIDs)[strings.ToLower(update.Message.Chat.Title)] = update.Message.Chat.ID
			Utils.SerializeObject(client.TelegramIDs, "data/telegram_ids.gob")
			break
		}
	}
}

// associate_created_group if the bot creates a group makes chat name-chat id association here
func associate_created_group(client *TelegramClient, update tgbotapi.Update) {
	(*client.TelegramIDs)[strings.ToLower(update.Message.Chat.Title)] = update.Message.Chat.ID
	Utils.SerializeObject(client.TelegramIDs, "data/telegram_ids.gob")
}
