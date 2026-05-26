#!/bin/sh
set -eu

. /usr/local/bin/process.sh

trap cleanup_bg EXIT INT TERM

start_bg kernel poc4-kernel --node alice --app-listen 127.0.0.1:7201
sleep 1
start_bg relay poc4-relay --node alice --kernel 127.0.0.1:7201 --listen 0.0.0.0:9200 --routes bob=bob:9200,carol=bob:9200,dave=ellen:9200,ellen=ellen:9200
sleep 4

poc4-hello --node alice --app alice-hello-app --kernel 127.0.0.1:7201 --mode ask-signed --target-node dave --target-app dave-signed-app --text "hello from Alice"
poc4-fibonacci --node alice --app alice-fibonacci-client --kernel 127.0.0.1:7201 --mode client --target-node carol --target-app carol-fibonacci-app --n 10
sleep 2
