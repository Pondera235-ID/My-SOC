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

type SuricataMessage struct {
    Timestamp   string `json:"timestamp"`
    FlowID      int64  `json:"flow_id"`
    InIface     string `json:"in_iface"`
    EventType   string `json:"event_type"`
    SrcIP       string `json:"src_ip"`
    SrcPort     int    `json:"src_port"`
    DestIP      string `json:"dest_ip"`
    DestPort    int    `json:"dest_port"`
    Proto       string `json:"proto"`
    PktSrc      string `json:"pkt_src"`
    Ether       Ether  `json:"ether"`
    Http        Http   `json:"http"`
    AppProto    string `json:"app_proto"`
    FileInfo    FileInfo `json:"fileinfo"`
}

type Ether struct {
    SrcMac  string `json:"src_mac"`
    DestMac string `json:"dest_mac"`
}

type Http struct {
    Hostname        string `json:"hostname"`
    HttpPort        int    `json:"http_port"`
    URL             string `json:"url"`
    HttpUserAgent   string `json:"http_user_agent"`
    HttpContentType string `json:"http_content_type"`
    HttpMethod      string `json:"http_method"`
    Protocol        string `json:"protocol"`
    Status          int    `json:"status"`
    Length          int    `json:"length"`
}

type FileInfo struct {
    Filename string `json:"filename"`
    Gaps     bool   `json:"gaps"`
    State    string `json:"state"`
    Stored   bool   `json:"stored"`
    Size     int    `json:"size"`
    TxID     int    `json:"tx_id"`
}

type Pattern interface {
    Match(message string) interface{}
    Name() string
}
