package Telegram

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/omarperezr/SmarttyBot/Utils"
)

func gitlab_report_callback(client *TelegramClient, update tgbotapi.Update) {
	possible_queries := []string{"report"}
	query := strings.Split(update.CallbackQuery.Data, ":")[0]

	if Utils.StringInSlice(query, possible_queries) {
		main_key := fmt.Sprintf("%s|%d", query, (*update.CallbackQuery).Message.MessageID)
		secondary_key := strings.Split(update.CallbackQuery.Data, ":")[1]
		if main_dict, ok := (*client.AwaitingCallbackData)[main_key]; ok {
			if secondary_key != main_dict.(map[string]interface{})["current"].(string) {
				(*client.AwaitingCallbackData)[main_key].(map[string]interface{})["current"] = secondary_key
				edit := tgbotapi.NewEditMessageText(
					update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
					main_dict.(map[string]interface{})[secondary_key].(string),
				)
				edit.BaseEdit.ReplyMarkup = main_dict.(map[string]interface{})["markup"].(*tgbotapi.InlineKeyboardMarkup)
				if _, err := client.Api.Send(edit); err != nil {
					log.Panic("Error sending message to telegram:", err)
				}
			}
		}
	}
}
