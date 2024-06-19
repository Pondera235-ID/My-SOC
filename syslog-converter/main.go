package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net"
    "os"
    "path/filepath"
    "plugin"
    "syslog-converter/common"
)

const (
    maxLogFileSize = 5 * 1024 * 1024 // 5 MB
    consolidatedLogFileName = "consolidated.log"
)

func main() {
    addr, err := net.ResolveUDPAddr("udp", ":5140")
    if err != nil {
        log.Fatalf("Failed to resolve UDP address: %v", err)
    }

    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        log.Fatalf("Failed to listen on UDP port 5140: %v", err)
    }
    defer conn.Close()

    log.Println("Listening on UDP port 5140")

    patterns, err := loadPatterns("patterns")
    if err != nil {
        log.Fatalf("Failed to load patterns: %v", err)
    }

    buffer := make([]byte, 2048)

    for {
        n, addr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            log.Printf("Failed to read from UDP: %v", err)
            continue
        }

        message := string(buffer[:n])
        log.Printf("Received message from %s: %s", addr, message)

        for _, pattern := range patterns {
            switch pattern.Name() {
            case "filterlog":
                processFilterlog(pattern, message)
            case "suricata":
                processSuricata(pattern, message)
            }
        }
    }
}

func loadPatterns(dir string) ([]common.Pattern, error) {
    var patterns []common.Pattern
    err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() && filepath.Ext(path) == ".so" {
            p, err := plugin.Open(path)
            if err != nil {
                return err
            }
            sym, err := p.Lookup("Pattern")
            if err != nil {
                return err
            }
            pattern, ok := sym.(common.Pattern)
            if !ok {
                return fmt.Errorf("unexpected type from module symbol")
            }
            patterns = append(patterns, pattern)
        }
        return nil
    })
    return patterns, err
}

func processFilterlog(pattern common.Pattern, message string) {
    convertedMessages := pattern.Match(message).([]common.SyslogMessage)
    for _, convertedMessage := range convertedMessages {
        jsonMessage, err := json.MarshalIndent(convertedMessage, "", "  ")
        if err != nil {
            log.Printf("Failed to convert message to JSON: %v", err)
            continue
        }
        if err := writeLog("filterlog", string(jsonMessage)); err != nil {
            log.Printf("Failed to write log: %v", err)
        }
    }
}

func processSuricata(pattern common.Pattern, message string) {
    convertedMessages := pattern.Match(message).([]common.SuricataMessage)
    for _, convertedMessage := range convertedMessages {
        jsonMessage, err := json.MarshalIndent(convertedMessage, "", "  ")
        if err != nil {
            log.Printf("Failed to convert message to JSON: %v", err)
            continue
        }
        if err := writeLog("suricata", string(jsonMessage)); err != nil {
            log.Printf("Failed to write log: %v", err)
        }
    }
}

func writeLog(source, jsonMessage string) error {
    entry := map[string]interface{}{
        "source": source,
        "log":    jsonMessage,
    }
    consolidatedJSON, err := json.MarshalIndent(entry, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal consolidated log entry: %v", err)
    }

    if err := rotateLogFile(consolidatedLogFileName); err != nil {
        return err
    }

    file, err := os.OpenFile(consolidatedLogFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to open log file %s: %v", consolidatedLogFileName, err)
    }
    defer file.Close()

    if _, err := file.WriteString(string(consolidatedJSON) + "\n"); err != nil {
        return fmt.Errorf("failed to write to log file %s: %v", consolidatedLogFileName, err)
    }
    return nil
}

func rotateLogFile(filename string) error {
    fileInfo, err := os.Stat(filename)
    if err != nil {
        if os.IsNotExist(err) {
            return nil // No file to rotate
        }
        return err
    }

    if fileInfo.Size() < maxLogFileSize {
        return nil // No need to rotate
    }

    // Rename the current log file
    backupName := fmt.Sprintf("%s.bak", filename)
    if err := os.Rename(filename, backupName); err != nil {
        return fmt.Errorf("failed to rotate log file %s: %v", filename, err)
    }

    return nil
}
