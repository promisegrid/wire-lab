#!/bin/sh
set -eu

. /usr/local/bin/process.sh

install_traps

start_bg kernel poc5-kernel --node dave --app-listen 127.0.0.1:7201
sleep 1
start_bg relay poc5-relay --node dave --kernel 127.0.0.1:7201 --listen 0.0.0.0:9200 --routes alice=ellen:9200,carol=carol:9200,ellen=ellen:9200
sleep 1
start_bg storage poc5-storage --node dave --app dave-storage-app --kernel 127.0.0.1:7201 --mode serve
storage_pid="$LAST_PID"
poc5-signed --node dave --app dave-signed-app --kernel 127.0.0.1:7201 --mode idle
status=0
if wait_one storage "$storage_pid"; then
	echo "storage promise path completed"
else
	status="$?"
fi
if [ "$status" -eq 0 ]; then
	mark_done dave
	wait_done_all alice bob carol dave ellen
fi
exit "$status"
