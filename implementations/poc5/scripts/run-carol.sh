#!/bin/sh
set -eu

. /usr/local/bin/process.sh

install_traps

start_bg kernel poc5-kernel --node carol --app-listen 127.0.0.1:7201
sleep 1
start_bg relay poc5-relay --node carol --kernel 127.0.0.1:7201 --listen 0.0.0.0:9200 --routes alice=bob:9200,bob=bob:9200,dave=dave:9200
sleep 1
poc5-hello --node carol --app carol-hello-app --kernel 127.0.0.1:7201 --mode idle
poc5-signed --node carol --app carol-signed-app --kernel 127.0.0.1:7201 --mode idle
mark_done carol
wait_done_all alice bob carol dave ellen
