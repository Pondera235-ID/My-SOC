package common

type SyslogMessage struct {
    LogTimestamp     string `json:"log_timestamp"`
    ProcessName      string `json:"process_name"`
    ProcessID        string `json:"process_id"`
    Action           string `json:"action"`
    SrcIP            string `json:"src_ip"`
    DstIP            string `json:"dst_ip"`
    SrcPort          string `json:"src_port"`
    DstPort          string `json:"dst_port"`
    NetworkInterface string `json:"network_interface"`
    Direction        string `json:"direction"`
    Protocol         string `json:"protocol"`
}

type Pattern interface {
    Match(message string) []SyslogMessage
    Name() string
}