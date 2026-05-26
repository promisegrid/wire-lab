#!/bin/sh

PIDS=""
LAST_PID=""
RUN_ID="${POC4_RUN_ID:-manual}"
DONE_DIR="${POC4_DONE_DIR:-/run/poc4}/$RUN_ID"

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

mark_done() {
	node="$1"
	# Intent: This marker is Compose process coordination so
	# --abort-on-container-exit can mean "the bounded demo is complete" instead
	# of "one intentionally short-lived app returned." Source: DI-rinuv.
	if mkdir -p "$DONE_DIR"; then
		:
	else
		status="$?"
		echo "could not create done directory $DONE_DIR status=$status"
		return "$status"
	fi
	if printf '%s\n' "$node" >"$DONE_DIR/$node.done"; then
		echo "marked $node done in $DONE_DIR"
		return 0
	fi
	status="$?"
	echo "could not mark $node done status=$status"
	return "$status"
}

wait_done_all() {
	for node in "$@"; do
		while [ ! -f "$DONE_DIR/$node.done" ]; do
			sleep 1
		done
		echo "observed $node done"
	done
}
