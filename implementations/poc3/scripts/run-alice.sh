#!/bin/sh
set -u

poc3-kernel --node alice --app-listen 127.0.0.1:7101 --peer-listen 0.0.0.0:9100 --peer bob:9100 &
kernel_pid=$!

sleep 2

poc3-echo --node alice --kernel 127.0.0.1:7101 --mode ask --to bob --text "echo this promise" &
echo_pid=$!

poc3-hello --node alice --kernel 127.0.0.1:7101 --mode send --to bob --text "hello from Alice"
hello_status=$?

poc3-signed --node alice --kernel 127.0.0.1:7101 --mode send --to bob --text "signed hello from Alice"
signed_status=$?

if wait "${echo_pid}"; then
  echo_status=0
else
  echo_status=$?
fi

if kill "${kernel_pid}"; then
  if wait "${kernel_pid}"; then
    kernel_status=0
  else
    kernel_status=$?
    if [ "${kernel_status}" -eq 143 ]; then
      kernel_status=0
    fi
  fi
else
  echo "alice kernel already stopped" >&2
  kernel_status=0
fi

if [ "${hello_status}" -ne 0 ]; then
  exit "${hello_status}"
fi
if [ "${signed_status}" -ne 0 ]; then
  exit "${signed_status}"
fi
if [ "${echo_status}" -ne 0 ]; then
  exit "${echo_status}"
fi
exit "${kernel_status}"
