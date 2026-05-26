#!/bin/sh

PIDS=""
LAST_PID=""

start_bg() {
	name="$1"
	shift
	"$@" &
	pid="$!"
	PIDS="$PIDS $pid"
	LAST_PID="$pid"
	echo "started $name pid=$pid"
}

wait_one() {
	name="$1"
	pid="$2"
	if wait "$pid"; then
		echo "$name exited cleanly"
		return 0
	fi
	status="$?"
	echo "$name exited with status $status"
	return "$status"
}

cleanup_bg() {
	for pid in $PIDS; do
		if kill -0 "$pid" 2>/dev/null; then
			if kill "$pid"; then
				echo "sent termination to pid=$pid"
			else
				status="$?"
				echo "termination for pid=$pid returned status $status"
			fi
		fi
	done
	for pid in $PIDS; do
		if wait "$pid"; then
			echo "cleaned pid=$pid"
		else
			status="$?"
			echo "cleaned pid=$pid status=$status"
		fi
	done
}
