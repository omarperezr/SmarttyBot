package Wit

import (
	"fmt"
	"strings"
)

// Format_issues returns a list of text formatted gitlab issues
func format_issues(jsonMap map[string]interface{}) map[string]string {
	formatted_issues := make(map[string]string)

	for key, val := range jsonMap {
		list_issues := make([]string, len(val.([]interface{})))
		for i, v := range val.([]interface{}) {
			list_issues[i] = fmt.Sprint(v)
		}

		formatted_issues[key] = strings.Join(list_issues, "\n")
	}

	return formatted_issues
}
