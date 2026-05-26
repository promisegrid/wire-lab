#!/bin/sh
set -eu

. /usr/local/bin/process.sh

install_traps

start_bg kernel poc5-kernel --node alice --app-listen 127.0.0.1:7201
sleep 1
start_bg relay poc5-relay --node alice --kernel 127.0.0.1:7201 --listen 0.0.0.0:9200 --routes bob=bob:9200,carol=bob:9200,dave=ellen:9200,ellen=ellen:9200
sleep 4

poc5-storage --node alice --app alice-trust-client --kernel 127.0.0.1:7201 --mode trust-client --target-node bob --target-app bob-storage-app --fallback-node dave --fallback-app dave-storage-app --key poc5-sensitive-key --value poc5-sensitive-value
mark_done alice
wait_done_all alice bob carol dave ellen
