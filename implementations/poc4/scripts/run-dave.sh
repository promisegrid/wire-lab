#!/bin/sh
set -eu

. /usr/local/bin/process.sh

trap cleanup_bg EXIT INT TERM

start_bg kernel poc4-kernel --node dave --app-listen 127.0.0.1:7201
sleep 1
start_bg relay poc4-relay --node dave --kernel 127.0.0.1:7201 --listen 0.0.0.0:9200 --routes alice=ellen:9200,carol=carol:9200,ellen=ellen:9200
sleep 1
start_bg storage poc4-storage --node dave --app dave-storage-app --kernel 127.0.0.1:7201 --mode serve
storage_pid="$LAST_PID"
start_bg signed poc4-signed --node dave --app dave-signed-app --kernel 127.0.0.1:7201 --mode serve
signed_pid="$LAST_PID"
status=0
if wait_one storage "$storage_pid"; then
	echo "storage promise path completed"
else
	status="$?"
fi
if wait_one signed "$signed_pid"; then
	echo "signed promise path completed"
else
	status="$?"
fi
exit "$status"
