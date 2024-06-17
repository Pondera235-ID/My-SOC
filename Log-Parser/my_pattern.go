package main

import (
    "regexp"
    "strings"
    "syslog-converter/common"
)

type MyPattern struct{}

func (p MyPattern) Match(message string) []common.SyslogMessage {
    pattern := regexp.MustCompile(`<\d+>(?P<LogTimestamp>\w+ \d+ \d+:\d+:\d+) (?P<Process>[\w\[\]\d]+): \d+,,,\d+,(?P<NetworkInterface>\w+),(?P<Action>[\w-]+,\w+),(?P<Direction>\w+),.*?,.*?,.*?,.*?,.*?,(?P<Protocol>\w+),.*?,(?P<SrcIP>\d+\.\d+\.\d+\.\d+),(?P<DstIP>\d+\.\d+\.\d+\.\d+),(?P<SrcPort>\d+),(?P<DstPort>\d+)`)
    matches := pattern.FindAllStringSubmatch(message, -1)

    var convertedMessages []common.SyslogMessage
    fieldNames := pattern.SubexpNames()

    for _, match := range matches {
        var msg common.SyslogMessage
        for i, name := range fieldNames {
            if i != 0 && name != "" {
                switch name {
                case "LogTimestamp":
                    msg.LogTimestamp = match[i]
                case "Process":
                    processDetails := strings.Split(match[i], "[")
                    msg.ProcessName = processDetails[0]
                    if len(processDetails) > 1 {
                        msg.ProcessID = strings.TrimSuffix(processDetails[1], "]")
                    }
                case "Action":
                    actions := strings.Split(match[i], ",")
                    if len(actions) > 1 {
                        msg.Action = actions[0] + "-" + actions[1]
                    }
                case "NetworkInterface":
                    msg.NetworkInterface = match[i]
                case "Direction":
                    msg.Direction = match[i]
                case "Protocol":
                    msg.Protocol = match[i]
                case "SrcIP":
                    msg.SrcIP = match[i]
                case "DstIP":
                    msg.DstIP = match[i]
                case "SrcPort":
                    msg.SrcPort = match[i]
                case "DstPort":
                    msg.DstPort = match[i]
                }
            }
        }
        convertedMessages = append(convertedMessages, msg)
    }

    return convertedMessages
}

func (p MyPattern) Name() string {
    return "filterlog"
}

var Pattern MyPattern