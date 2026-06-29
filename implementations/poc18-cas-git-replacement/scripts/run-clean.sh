#!/usr/bin/env bash
set -euo pipefail

# Intent: Keep generated POC18 runtime state inside the approved /tmp pattern and
# make cleanup explicit. Source: DI-jifuj
run_root="/tmp/wire-lab-poc18-run"
export GOCACHE="${GOCACHE:-/tmp/wire-lab-poc18-gocache}"
export GOMODCACHE="${GOMODCACHE:-/tmp/wire-lab-poc18-gomodcache}"

if [ -d "$run_root" ]; then
  rm -rf "$run_root"
else
  printf 'run root did not exist: %s\n' "$run_root"
fi

go test ./...
go run ./cmd/poc-sim -run-root "$run_root"
go run ./cmd/poc-analyze -run-root "$run_root"
go run ./cmd/poc-cbor-diag -run-root "$run_root"
