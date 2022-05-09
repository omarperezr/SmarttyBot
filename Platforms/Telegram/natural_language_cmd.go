package Telegram

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// natural_language_cmd sends a request to WitAI to parse a message using nlp, uses a channel to send a message to witAI_client and
// execute a corresponding action based on the message
func natural_language_cmd(client *TelegramClient, update tgbotapi.Update, message_fields []string) {
	var message_sent tgbotapi.Message
	var err error

	telegram_message, markup := client.WitClient.MessageParser(strings.Join(message_fields[1:], " "))
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, telegram_message)
	if markup != nil {
		msg.ReplyMarkup = markup
	}

	if message_sent, err = client.Api.Send(msg); err != nil {
		log.Panic("Error sending message to telegram:", err, message_sent)
	}
	client.TeleWitChan <- fmt.Sprintf("%d", message_sent.MessageID)
}
