#!/bin/sh
set -eu

. /usr/local/bin/process.sh

trap cleanup_bg EXIT INT TERM

start_bg kernel poc4-kernel --node carol --app-listen 127.0.0.1:7201
sleep 1
start_bg relay poc4-relay --node carol --kernel 127.0.0.1:7201 --listen 0.0.0.0:9200 --routes alice=bob:9200,bob=bob:9200,dave=dave:9200
sleep 1
start_bg fibonacci poc4-fibonacci --node carol --app carol-fibonacci-app --kernel 127.0.0.1:7201 --mode serve
fibonacci_pid="$LAST_PID"
sleep 3
poc4-storage --node carol --app carol-storage-client --kernel 127.0.0.1:7201 --mode client --target-node dave --target-app dave-storage-app --key poc4-key --value poc4-value
wait_one fibonacci "$fibonacci_pid"
sleep 2
