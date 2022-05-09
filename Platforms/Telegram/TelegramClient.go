package Telegram

import (
	"log"
	"strings"

	"github.com/omarperezr/SmarttyBot/Core/Config"
	dClient "github.com/omarperezr/SmarttyBot/Platforms/Discord"
	wclient "github.com/omarperezr/SmarttyBot/Platforms/Wit"

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
	AwaitingCallbackData *map[string]interface{}
}

func SetUp(config *Config.Config, witai_client *wclient.WitClient) TelegramClient {
	instance := TelegramClient{
		ApiKey:               config.Telegram_Api_Key,
		TelegramIDs:          config.Telegram_ids,
		TelegramMentions:     config.Telegram_mentions,
		DiscordMentions:      config.Discord_mentions,
		WitClient:            witai_client,
		TeleWitChan:          config.Tele_wit_chan,
		AwaitingCallbackData: config.Callback_map,
	}
	return instance
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
			if update.Message.NewChatMembers != nil {
				associate_existing_group(client, update)

			} else if update.Message.GroupChatCreated {
				associate_created_group(client, update)

			} else if len(message_fields) > 1 {
				if message_fields[0] == "-m" || message_fields[0] == "-mention" {
					mention_user_cmd(client, update, message_fields)
				} else if message_fields[0] == "-p" || message_fields[0] == "-parse" {
					natural_language_cmd(client, update, message_fields)
				} else if message_fields[0] == "-r" || message_fields[0] == "-register" {
					register_user_cmd(client, update, message_fields)
				}
			}

		} else if update.CallbackQuery != nil { // Callbacks, messages from pressing buttons and stuff like that
			gitlab_report_callback(client, update)
		}
	}
}
