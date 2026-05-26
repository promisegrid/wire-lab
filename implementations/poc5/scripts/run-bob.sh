#!/bin/sh
set -eu

. /usr/local/bin/process.sh

install_traps

start_bg kernel poc5-kernel --node bob --app-listen 127.0.0.1:7201
sleep 1
start_bg relay poc5-relay --node bob --kernel 127.0.0.1:7201 --listen 0.0.0.0:9200 --routes alice=alice:9200,carol=carol:9200,ellen=alice:9200
sleep 1
poc5-signed --node bob --app bob-signed-app --kernel 127.0.0.1:7201 --mode idle
start_bg storage poc5-storage --node bob --app bob-storage-app --kernel 127.0.0.1:7201 --mode serve-break
storage_pid="$LAST_PID"
wait_one storage "$storage_pid"
mark_done bob
wait_done_all alice bob carol dave ellen
