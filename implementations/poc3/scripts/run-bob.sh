#!/bin/sh
set -u

poc3-kernel --node bob --app-listen 127.0.0.1:7101 --peer-listen 0.0.0.0:9100 --peer alice:9100 &
kernel_pid=$!

sleep 1

poc3-hello --node bob --kernel 127.0.0.1:7101 --mode receive --to bob --text "Bob receives hello" &
hello_pid=$!

poc3-echo --node bob --kernel 127.0.0.1:7101 --mode serve --to bob --text "Bob serves echo" &
echo_pid=$!

poc3-signed --node bob --kernel 127.0.0.1:7101 --mode receive --to bob --text "Bob receives signed" &
signed_pid=$!

status=0
if wait "${hello_pid}"; then
  hello_status=0
else
  hello_status=$?
  status="${hello_status}"
fi
if wait "${echo_pid}"; then
  echo_status=0
else
  echo_status=$?
  status="${echo_status}"
fi
if wait "${signed_pid}"; then
  signed_status=0
else
  signed_status=$?
  status="${signed_status}"
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
  echo "bob kernel already stopped" >&2
  kernel_status=0
fi

if [ "${status}" -ne 0 ]; then
  exit "${status}"
fi
exit "${kernel_status}"
