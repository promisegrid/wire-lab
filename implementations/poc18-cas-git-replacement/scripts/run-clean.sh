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
mkdir -p "$run_root"

go test ./...
go run ./cmd/poc-sim -run-root "$run_root"
go run ./cmd/poc-analyze -run-root "$run_root"
go run ./cmd/poc-cbor-diag -run-root "$run_root"
go build -o "$run_root/grid" ./cmd/grid
# Intent: Exercise the normal repo-local CLI path so `.grid/state.json`,
# `grid status`, `grid log`, `grid track`, and `grid untrack` stay covered by
# the clean POC run. Running this after poc-sim preserves both fixture JSON logs
# and the smoke repo's `.grid` directory for inspection. Source: DI-bikif;
# DI-kiram; DI-jokav
grid_repo="$run_root/grid-cli"
mkdir -p "$grid_repo"
printf 'POC18 grid CLI clean run\n' > "$grid_repo/README.md"
printf 'local-only smoke data\n' > "$grid_repo/local.log"
(
  cd "$grid_repo"
  "$run_root/grid" init
  "$run_root/grid" snapshot -out "$run_root/grid-snapshot-initial.json"
  "$run_root/grid" untrack local.log
  "$run_root/grid" status -out "$run_root/grid-status-tracking-removed.json"
  "$run_root/grid" snapshot -out "$run_root/grid-snapshot.json"
  "$run_root/grid" status -out "$run_root/grid-status.json"
  "$run_root/grid" log -out "$run_root/grid-log.json"
  "$run_root/grid" track local.log
  "$run_root/grid" status -out "$run_root/grid-status-tracking-added.json"
  "$run_root/grid" snapshot -out "$run_root/grid-snapshot-tracked.json"
  "$run_root/grid" status -out "$run_root/grid-status-tracked.json"
)
