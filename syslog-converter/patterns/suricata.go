package main

import (
    "encoding/json"
    "strings"
    "syslog-converter/common"
)

type SuricataPattern struct{}

func (p SuricataPattern) Match(message string) interface{} {
    var messages []common.SuricataMessage

    entries := strings.Split(message, "<45>")
    for _, entry := range entries {
        if entry == "" {
            continue
        }

        start := strings.Index(entry, "{")
        if start == -1 {
            continue
        }

        jsonStr := entry[start:]
        var msg common.SuricataMessage
        if err := json.Unmarshal([]byte(jsonStr), &msg); err != nil {
            continue
        }

        messages = append(messages, msg)
    }

    return messages
}

func (p SuricataPattern) Name() string {
    return "suricata"
}

var Pattern SuricataPattern
