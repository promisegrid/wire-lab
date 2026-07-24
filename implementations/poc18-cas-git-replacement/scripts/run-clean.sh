#!/usr/bin/env bash
set -euo pipefail

# Intent: Keep generated POC18 runtime state inside the approved /tmp pattern and
# make cleanup explicit. Source: DI-jifuj
run_root="/tmp/wire-lab-poc18-run"
export GOCACHE="${GOCACHE:-/tmp/wire-lab-poc18-gocache}"
export GOMODCACHE="${GOMODCACHE:-/tmp/wire-lab-poc18-gomodcache}"
export POC18_UID="${POC18_UID:-$(id -u)}"
export POC18_GID="${POC18_GID:-$(id -g)}"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
poc_dir="$(cd "${script_dir}/.." && pwd)"

cd "$poc_dir"

if [ -d "$run_root" ]; then
  if rm -rf "$run_root"; then
    echo "removed existing run root: $run_root"
  else
    status=$?
    echo "host cleanup failed with status ${status}; attempting scoped docker cleanup for $run_root" >&2
    if POC18_RUN_ROOT="$run_root" docker compose run --rm --user root --entrypoint /bin/sh event-collector -c 'rm -rf /run/poc18/*'; then
      echo "docker cleanup complete"
    else
      docker_status=$?
      echo "docker cleanup failed with status ${docker_status}" >&2
      exit "${docker_status}"
    fi
    if rm -rf "$run_root"; then
      echo "removed run root after docker cleanup: $run_root"
    else
      retry_status=$?
      echo "host cleanup retry failed with status ${retry_status}" >&2
      exit "${retry_status}"
    fi
  fi
else
  printf 'run root did not exist: %s\n' "$run_root"
fi

echo "== test POC18 packages =="
go test ./...

echo "== run deterministic POC18 fixture =="
# Intent: Restore the deterministic fixture after package tests so later
# diagnostics can render stable reference-set, snapshot, review, merge, and
# materialization CAS objects before the Docker/TCP phase appends live run
# artifacts. Source: DI-basan
go run ./cmd/poc-sim -run-root "$run_root"

echo "== analyze deterministic POC18 fixture =="
# Intent: `nahop.20` requires the clean regression to fail before Docker if the
# deterministic CAS/VCS fixture no longer proves CID correctness, sparse CAS,
# parent-chain integrity, reference-set semantics, Rabin chunking, Git bridge,
# retention, token economics, and promise vocabulary. Source: DI-bovaf
go run ./cmd/poc-analyze -run-root "$run_root"

mkdir -p "$run_root/observer"

echo "== reset POC18 docker runtime state =="
if POC18_RUN_ROOT="$run_root" POC18_UID="$POC18_UID" POC18_GID="$POC18_GID" docker compose down -v --remove-orphans; then
  echo "reset complete"
else
  status=$?
  echo "docker compose down failed with status ${status}" >&2
  exit "${status}"
fi

echo "== run POC18 TCP agent regression =="
# Intent: The clean regression must exercise real Docker-container agents over
# TCP and must not fall back to in-process peer CAS fixtures. Source: DI-koriz
if POC18_RUN_ROOT="$run_root" POC18_UID="$POC18_UID" POC18_GID="$POC18_GID" docker compose up --build; then
  echo "compose run complete"
else
  status=$?
  echo "docker compose up failed with status ${status}" >&2
  exit "${status}"
fi

echo "== analyze POC18 TCP artifacts =="
observer_run="$run_root/observer/poc18-demo"
go run ./cmd/poc-analyze -run-root "$observer_run"
first_message=""
while IFS= read -r candidate; do
  first_message="$candidate"
  break
done < <(find "$observer_run/message-cas" -type f -name '*.cbor' -print | sort)
if [ -z "$first_message" ]; then
  echo "no collected CBOR message found under $observer_run/message-cas" >&2
  exit 1
fi
go run ./cmd/poc-cbor-diag -file "$first_message" > "$run_root/poc18-first-message.diag.txt"
diagnostic_dir="$run_root/poc18-diagnostics"
# Intent: Produce curated exact-CBOR diagnostic renderings for representative
# POC18 flows so protocol reviews do not depend on hand-picked or paraphrased
# examples. Source: DI-basan
go run ./cmd/poc-cbor-diag -diagnostic-report -run-root "$run_root" -out-dir "$diagnostic_dir"
for diagnostic_flow in reference-set node-version snapshot review-statement merge-snapshot directory-node materialization-object sync-interest object-availability object-retrieval-redemption; do
  diagnostic_path="$diagnostic_dir/${diagnostic_flow}.diag.txt"
  if [ -s "$diagnostic_path" ]; then
    printf 'diagnostic ready: %s\n' "$diagnostic_path"
  else
    echo "missing or empty diagnostic: $diagnostic_path" >&2
    exit 1
  fi
done
if [ -s "$diagnostic_dir/index.json" ]; then
  printf 'diagnostic index ready: %s\n' "$diagnostic_dir/index.json"
else
  echo "missing or empty diagnostic index: $diagnostic_dir/index.json" >&2
  exit 1
fi

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
