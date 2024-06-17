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
            convertedMessages := pattern.Match(message)
            for _, convertedMessage := range convertedMessages {
                jsonMessage, err := json.MarshalIndent(convertedMessage, "", "  ")
                if err != nil {
                    log.Printf("Failed to convert message to JSON: %v", err)
                    continue
                }
                if err := writeLog(pattern.Name(), string(jsonMessage)); err != nil {
                    log.Printf("Failed to write log: %v", err)
                }
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

func writeLog(patternName, jsonMessage string) error {
    filename := fmt.Sprintf("%s.log", patternName)
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to open log file %s: %v", filename, err)
    }
    defer file.Close()

    if _, err := file.WriteString(jsonMessage + "\n"); err != nil {
        return fmt.Errorf("failed to write to log file %s: %v", filename, err)
    }
    return nil
}
