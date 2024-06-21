#!/usr/bin/env python3

import socket
import json
import datetime

LOGSTASH_HOST = "localhost"
LOGSTASH_PORT = 514
LISTEN_PORT = 5140
LOG_FILE = "/var/log/reformat_logs.log"

def log_message(message):
    with open(LOG_FILE, "a") as log_file:
        log_file.write(f"{datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')} - {message}\n")

def reformat_filterlog(line):
    parts = line.split(",")
    try:
        log_data = {
            "source": "filterlog",
            "rule_id": parts[3],
            "interface": parts[4],
            "action": parts[6],
            "direction": parts[7],
            "ip_version": parts[8],
            "tos": parts[9],
            "ttl": parts[10],
            "id": parts[11],
            "offset": parts[12],
            "ip_flags": parts[13],
            "proto_id": parts[14],
            "protocol": parts[15],
            "length": parts[16],
            "src_ip": parts[17],
            "dest_ip": parts[18],
            "src_port": parts[19],
            "dest_port": parts[20],
            "data_length": parts[21],
            "tcp_flags": parts[22],
            "seq": parts[23],
            "window": parts[25],
            "options": parts[26]
        }
        return json.dumps(log_data)
    except IndexError:
        return None

def reformat_suricata(line):
    try:
        json_start = line.find("{")
        json_data = line[json_start:]
        log_data = json.loads(json_data)
        log_data["source"] = "suricata"
        return json.dumps(log_data)
    except json.JSONDecodeError:
        return None

def main():
    listen_sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    listen_sock.bind(("0.0.0.0", LISTEN_PORT))

    send_sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

    log_message("Listening for UDP packets on port 5140...")

    while True:
        data, addr = listen_sock.recvfrom(65535)
        log_line = data.decode("utf-8").strip()

        if "filterlog" in log_line:
            formatted_log = reformat_filterlog(log_line)
        elif "suricata" in log_line:
            formatted_log = reformat_suricata(log_line)
        else:
            formatted_log = None
            log_message(f"Unrecognized log format: {log_line}")

        if formatted_log:
            log_message(f"Successfully reformatted: {log_line}")
            send_sock.sendto(formatted_log.encode("utf-8"), (LOGSTASH_HOST, LOGSTASH_PORT))

if __name__ == "__main__":
    main()
