#!/bin/bash

while true; do
    tail -F consolidated.log | nc -u 10.10.2.10 514
    sleep 1
done
