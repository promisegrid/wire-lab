#!/bin/sh

PIDS=""
LAST_PID=""
RUN_ID="${POC4_RUN_ID:-manual}"
DONE_DIR="${POC4_DONE_DIR:-/run/poc4}/$RUN_ID"
DONE_MARKED=0

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
		remove_pid "$pid"
		echo "$name exited cleanly"
		return 0
	fi
	wait_status="$?"
	remove_pid "$pid"
	echo "$name exited with status $wait_status"
	return "$wait_status"
}

remove_pid() {
	pid_to_remove="$1"
	remaining=""
	for pid in $PIDS; do
		if [ "$pid" = "$pid_to_remove" ]; then
			:
		else
			remaining="$remaining $pid"
		fi
	done
	PIDS="$remaining"
}

cleanup_bg() {
	for pid in $PIDS; do
		if kill -0 "$pid" 2>/dev/null; then
			if kill "$pid"; then
				echo "sent termination to pid=$pid"
			else
				kill_status="$?"
				echo "termination for pid=$pid returned status $kill_status"
			fi
		fi
	done
	for pid in $PIDS; do
		if wait "$pid"; then
			echo "cleaned pid=$pid"
		else
			cleanup_status="$?"
			echo "cleaned pid=$pid status=$cleanup_status"
		fi
	done
}

install_traps() {
	trap on_exit EXIT
	trap on_signal INT TERM
}

on_exit() {
	exit_status="$?"
	trap - EXIT INT TERM
	cleanup_bg
	exit "$exit_status"
}

on_signal() {
	trap - EXIT INT TERM
	cleanup_bg
	if [ "$DONE_MARKED" -eq 1 ]; then
		exit 0
	fi
	exit 143
}

mark_done() {
	node="$1"
	# Intent: This marker is Compose process coordination so
	# --abort-on-container-exit can mean "the bounded demo is complete" instead
	# of "one intentionally short-lived app returned." Source: DI-rinuv.
	if mkdir -p "$DONE_DIR"; then
		:
	else
		marker_status="$?"
		echo "could not create done directory $DONE_DIR status=$marker_status"
		return "$marker_status"
	fi
	if printf '%s\n' "$node" >"$DONE_DIR/$node.done"; then
		DONE_MARKED=1
		echo "marked $node done in $DONE_DIR"
		return 0
	fi
	marker_status="$?"
	echo "could not mark $node done status=$marker_status"
	return "$marker_status"
}

wait_done_all() {
	for node in "$@"; do
		while [ ! -f "$DONE_DIR/$node.done" ]; do
			sleep 1
		done
		echo "observed $node done"
	done
}
