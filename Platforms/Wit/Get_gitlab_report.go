package Wit

import (
	"fmt"
	"log"

	"github.com/omarperezr/SmarttyBot/Utils"
)

func get_gitlab_report(client *WitClient, parameters []string) (string, interface{}) {
	jsonMap, err := Utils.Execute_system_command("python", "gitlab_report.py", "--report", parameters[0])
	if err != nil {
		log.Panic(err)
	}

	home_message := "You can check the current status of the project from this message!\n\nJust select what you want to check"

	// Data that is returned when buttons are pressed, also data to identify this map
	report_map := make(map[string]interface{})
	report_map["current"] = "home"
	report_map["purpose"] = "report"
	report_map["home"] = home_message
	report_map["status"] = fmt.Sprintf("The current phase is %f%% done\n", jsonMap["status"])
	report_map["oldest_todo"] = jsonMap["oldest_todo"]
	report_map["oldest_in_progress"] = jsonMap["oldest_in_progress"]

	issues := map[string]interface{}{
		"issues_todo":        jsonMap["issues_todo"],
		"issues_in_progress": jsonMap["issues_in_progress"],
		"issues_done":        jsonMap["issues_done"],
	}

	report_map["issues_todo"] = format_issues(issues)["issues_todo"]
	report_map["issues_in_progress"] = format_issues(issues)["issues_in_progress"]
	report_map["issues_done"] = format_issues(issues)["issues_done"]

	markup := client.Generate_markup("Markups/gitlab_report.json")
	report_map["markup"] = markup
	go client.wait_for_telegram_msg(report_map)

	report_message := home_message
	return report_message, markup
}
