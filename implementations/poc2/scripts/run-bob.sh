#!/bin/sh
set -u

mkdir -p /work/run
poc2-kernel --node bob --app-listen 127.0.0.1:7001 --peer-listen 0.0.0.0:9000 --evidence /work/run/evidence.jsonl &
kernel_pid=$!

poc2-hello --node bob --kernel 127.0.0.1:7001 --mode receive --to bob --text "Bob receives hello"
app_status=$?

if kill "${kernel_pid}"; then
  wait "${kernel_pid}"
  wait_status=$?
  if [ "${wait_status}" -ne 0 ] && [ "${wait_status}" -ne 143 ]; then
    echo "bob kernel wait status ${wait_status}" >&2
    exit "${wait_status}"
  fi
else
  echo "bob kernel already stopped" >&2
fi

exit "${app_status}"
