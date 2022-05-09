package Wit

import (
	"log"

	"github.com/omarperezr/SmarttyBot/Utils"
)

func get_gitlab_oldest_issue(client *WitClient, parameters []string) (string, interface{}) {
	jsonMap, err := Utils.Execute_system_command("python", "gitlab_report.py", "--oldest_issue", parameters[0])
	if err != nil {
		log.Panic(err)
	}

	home_message := "Oldest issue"
	report_map := make(map[string]interface{})
	report_map["current"] = "home"

	report_map["home"] = home_message
	report_map["oldest_todo"] = jsonMap["oldest_todo"]
	report_map["oldest_in_progress"] = jsonMap["oldest_in_progress"]

	markup := client.Generate_markup("Markups/gitlab_report.json")
	report_map["markup"] = markup
	go client.wait_for_telegram_msg(report_map)

	return home_message, markup
}
