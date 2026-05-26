#!/bin/sh
set -eu

. /usr/local/bin/process.sh

trap cleanup_bg EXIT INT TERM

start_bg kernel poc4-kernel --node bob --app-listen 127.0.0.1:7201
sleep 1
start_bg relay poc4-relay --node bob --kernel 127.0.0.1:7201 --listen 0.0.0.0:9200 --routes alice=alice:9200,carol=carol:9200,ellen=alice:9200
sleep 1
poc4-signed --node bob --app bob-signed-app --kernel 127.0.0.1:7201 --mode idle
sleep 3
poc4-echo --node bob --app bob-echo-app --kernel 127.0.0.1:7201 --mode client --target-node ellen --target-app ellen-echo-app --text "echo from Bob"
mark_done bob
wait_done_all alice bob carol dave ellen
