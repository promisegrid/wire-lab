#!/bin/sh
set -u

mkdir -p /work/run
poc2-kernel --node alice --app-listen 127.0.0.1:7001 --peer-listen 0.0.0.0:9000 --peer bob:9000 --evidence /work/run/evidence.jsonl &
kernel_pid=$!

poc2-hello --node alice --kernel 127.0.0.1:7001 --mode send --to bob --text "hello from Alice"
app_status=$?

if kill "${kernel_pid}"; then
  wait "${kernel_pid}"
  wait_status=$?
  if [ "${wait_status}" -ne 0 ] && [ "${wait_status}" -ne 143 ]; then
    echo "alice kernel wait status ${wait_status}" >&2
    exit "${wait_status}"
  fi
else
  echo "alice kernel already stopped" >&2
fi

exit "${app_status}"
