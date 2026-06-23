#!/usr/bin/env bash
set -euo pipefail

# Intent: Give POC17 one deterministic clean command that resets generated
# artifacts, runs the Go behavior simulator, and checks analyzer gates. Source:
# DI-pobir
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
poc_dir="$(cd "${script_dir}/.." && pwd)"
run_dir="/tmp/wire-lab-poc17/poc17-demo"
export GOCACHE="${GOCACHE:-/tmp/wire-lab-go-build}"

cd "${poc_dir}"

echo "== reset POC17 runtime artifacts =="
if rm -rf "${run_dir}"; then
	echo "reset complete"
else
	status=$?
	echo "runtime artifact reset failed with status ${status}" >&2
	exit "${status}"
fi

echo "== run POC17 simulator =="
if go run ./cmd/poc17-sim -config config.json; then
	echo "simulator run complete"
else
	status=$?
	echo "poc17 simulator failed with status ${status}" >&2
	exit "${status}"
fi

echo "== analyze POC17 artifacts =="
if go run ./cmd/poc17-analyze "${run_dir}"; then
	echo "POC17 clean regression PASS"
else
	status=$?
	echo "POC17 clean regression FAIL with status ${status}" >&2
	exit "${status}"
fi
