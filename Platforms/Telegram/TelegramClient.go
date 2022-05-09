package Telegram

import (
	"fmt"
	"log"
	"strings"

	dClient "github.com/omarperezr/SmarttyBot/Platforms/Discord"
	wclient "github.com/omarperezr/SmarttyBot/Platforms/Wit"
	"github.com/omarperezr/SmarttyBot/Utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramClient struct {
	ApiKey               string
	Api                  *tgbotapi.BotAPI
	TelegramIDs          *map[string]int64
	TelegramMentions     *map[string]string
	DiscordMentions      *map[string]string
	DiscordClient        *dClient.DiscordClient
	WitClient            *wclient.WitClient
	TeleWitChan          chan string
	AwaitingCallbackData map[string]interface{}
}

func (client *TelegramClient) Init() {
	// Create a new Telegram client session using the provided client token.
	var telegram_err error
	client.Api, telegram_err = tgbotapi.NewBotAPI(client.ApiKey)
	if telegram_err != nil {
		log.Panic("Error creating Telegram session,", telegram_err)
	}
	client.Api.Debug = false
}

func (client *TelegramClient) replaceMentionT2D(messageText string) string {
	message_fields := strings.Fields(messageText)
	new_message := messageText
	for _, word := range message_fields {
		if strings.HasPrefix(word, "@") {
			if discord_mention_key, ok := (*client.TelegramMentions)[word]; ok {
				new_message = strings.Replace(new_message, word, (*client.DiscordMentions)[discord_mention_key], 1)
			}
		}
	}
	return new_message
}

func (client *TelegramClient) StartListening() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := client.Api.GetUpdatesChan(u)
	if err != nil {
		log.Panic("Error setting updates channel for telegram:", err)
	}

	for update := range updates {
		if update.Message != nil { // Messages and group updates
			message_fields := strings.Fields(update.Message.Text)
			// If a new group is created with the bot it will add the id with the name
			if update.Message.NewChatMembers != nil {
				for _, user := range *(update.Message.NewChatMembers) {
					if user.UserName == "smartty_bridge_bot" {
						(*client.TelegramIDs)[strings.ToLower(update.Message.Chat.Title)] = update.Message.Chat.ID
						Utils.SerializeObject(client.TelegramIDs, "data/telegram_ids.gob")
						break
					}
				}
			} else if update.Message.GroupChatCreated {
				(*client.TelegramIDs)[strings.ToLower(update.Message.Chat.Title)] = update.Message.Chat.ID
				Utils.SerializeObject(client.TelegramIDs, "data/telegram_ids.gob")

			} else if len(message_fields) > 1 && (message_fields[0] == "-m" || message_fields[0] == "-mention") {
				channel_id := client.DiscordClient.GetDiscordChannelID(update.Message.Chat.Title)
				if channel_id != "" {
					message_text := client.replaceMentionT2D(strings.Join(message_fields[1:], " "))
					discord_message := fmt.Sprintf("%s says: %s", strings.ToTitle(update.Message.From.FirstName), message_text)
					client.DiscordClient.Session.ChannelMessageSend(channel_id, discord_message)
				} else {
					log.Println("No text channel created for:", update.Message.Chat.Title)
				}

			} else if len(message_fields) > 1 && (message_fields[0] == "-p" || message_fields[0] == "-parse") {
				log.Println("HERE", message_fields)
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

			} else if len(message_fields) == 2 && (message_fields[0] == "-r" || message_fields[0] == "-register") {
				telegram_user_mention := fmt.Sprintf("@%s", update.Message.From.UserName)
				(*client.TelegramMentions)[telegram_user_mention] = strings.ToUpper(message_fields[1])
				Utils.SerializeObject(client.TelegramMentions, "data/telegram_mentions.gob")
			}
		} else if update.CallbackQuery != nil { // Callbacks, messages from pressing buttons and stuff like that
			possible_queries := []string{"report"}
			query := strings.Split(update.CallbackQuery.Data, ":")[0]

			if Utils.StringInSlice(query, possible_queries) {
				main_key := fmt.Sprintf("%s|%d", query, update.CallbackQuery.Message.MessageID)
				secondary_key := strings.Split(update.CallbackQuery.Data, ":")[1]
				if main_dict, ok := client.AwaitingCallbackData[main_key]; ok {
					if secondary_key != main_dict.(map[string]interface{})["current"].(string) {
						client.AwaitingCallbackData[main_key].(map[string]interface{})["current"] = secondary_key
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
	}
}
