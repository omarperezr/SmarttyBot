package Wit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/omarperezr/SmarttyBot/Core/Config"
	witai "github.com/wit-ai/wit-go"
)

type WitClient struct {
	ApiKey               string
	Client               *witai.Client
	TeleWitChan          chan string
	AwaitingCallbackData *map[string]interface{}
}

func SetUp(config *Config.Config) WitClient {
	instance := WitClient{
		ApiKey:               config.WIT_Api_Key,
		TeleWitChan:          config.Tele_wit_chan,
		AwaitingCallbackData: config.Callback_map,
	}
	return instance
}

// Init creates a new Wit Ai client to parse messages
func (client *WitClient) Init() {
	client.Client = witai.NewClient(client.ApiKey)
}

// Generate_markup generates a markup for the interactive interface that the telegram user gets as a message
func (client *WitClient) Generate_markup(filepath string) *tgbotapi.InlineKeyboardMarkup {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		log.Panic("Error opening markup file:", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	jsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(byteValue), &jsonMap)
	if err != nil {
		log.Panic(err)
	}
	var list_of_buttons [][]tgbotapi.InlineKeyboardButton

	for _, val := range jsonMap["markup"].([]interface{}) {
		row := []tgbotapi.InlineKeyboardButton{}
		for _, button_map := range val.(map[string]interface{}) {
			button_text := button_map.(map[string]interface{})["text"].(string)
			callbackData := button_map.(map[string]interface{})["callbackData"].(string)
			button := tgbotapi.InlineKeyboardButton{
				Text:         button_text,
				CallbackData: &callbackData,
			}
			row = append(row, button)
		}
		list_of_buttons = append(list_of_buttons, row)
	}
	markup := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: list_of_buttons}
	return &markup
}

// add_to_group NOT IMPLEMENTED
func (client *WitClient) add_to_group(parameters []string) (string, interface{}) {
	return fmt.Sprintf("Adding %v to group!\n", parameters), nil
}

// create_group NOT IMPLEMENTED
func (client *WitClient) create_group(parameters []string) (string, interface{}) {
	return fmt.Sprintf("Creating group called %s and adding %v!\n", parameters[0], parameters[1:]), nil
}

// get_parameters Gets a list of parameters from the received telegram message
func (client *WitClient) get_parameters(parameters_list []interface{}) []string {
	var parameters []string

	for _, param_map := range parameters_list {
		parameters = append(parameters, param_map.(map[string]interface{})["value"].(string))
	}

	return parameters
}

// process_telegram_data Stores callback data used to edit or answer telegram messages
func (client *WitClient) process_telegram_data(tele_data string, tele_map map[string]interface{}) {
	if tele_map != nil {
		if tele_map["purpose"] == "report" {
			key := fmt.Sprintf("report|%s", tele_data)
			(*client.AwaitingCallbackData)[key] = tele_map
		}
	}
}

// wait_for_telegram_msg After sending a message to telegram waits for an answer to that message
func (client *WitClient) wait_for_telegram_msg(extra_data map[string]interface{}) {
	telegram_data := <-client.TeleWitChan
	if extra_data != nil {
		client.process_telegram_data(telegram_data, extra_data)
	}
	client.process_telegram_data(telegram_data, nil)
}

// MessageParser Receives a message, parses it and executes a function based on the message intent
func (client *WitClient) MessageParser(message string) (string, interface{}) {
	var parameters []string

	// The functions we are able to execute through witai
	actions := map[string]interface{}{
		"add_to_group":      client.add_to_group,
		"create_group":      client.create_group,
		"get_gitlab_report": get_gitlab_report,
		"get_oldest_issues": get_gitlab_oldest_issue,
		"none":              "none",
	}

	wit_msg, _ := client.Client.Parse(&witai.MessageRequest{
		Query: message,
	})

	if wit_msg.Entities["intent"] == nil {
		log.Println("No intent found")
		go func() {
			<-client.TeleWitChan
		}()
		return "Sorry, I didn't understand that", nil
	}

	intent := wit_msg.Entities["intent"].([]interface{})
	function_name := intent[0].(map[string]interface{})["value"].(string)

	parameters = client.get_parameters(wit_msg.Entities["function_parameters"].([]interface{}))

	switch f := actions[function_name].(type) {

	case func(*WitClient, []string) (string, interface{}):
		string_response, markup := f(client, parameters)
		return string_response, markup
	case func([]string) (string, interface{}):
		string_response, markup := f(parameters)
		return string_response, markup
	default:
		return "No function has been implemented for this yet", nil

	}
}
