package Telegram

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/omarperezr/SmarttyBot/Utils"
)

// register_user_cmd registers the username that sends the telegram message to the discord username sent as a message
func register_user_cmd(client *TelegramClient, update tgbotapi.Update, message_fields []string) {
	if len(message_fields) == 2 {
		telegram_user_mention := fmt.Sprintf("@%s", update.Message.From.UserName)
		(*client.TelegramMentions)[telegram_user_mention] = strings.ToUpper(message_fields[1])
		Utils.SerializeObject(client.TelegramMentions, "data/telegram_mentions.gob")
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You used the command incorrectly, the correct usage is: -r @TelegramUser discordUser")
		if message_sent, err := client.Api.Send(msg); err != nil {
			log.Panic("Error sending message to telegram:", err, message_sent)
		}
	}
}
