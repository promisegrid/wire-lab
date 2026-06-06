#!/usr/bin/env bash
set -u

# Intent: Provide a repo-local clean regression runner that resets Docker state,
# runs the POC13 TCP scenario, and immediately analyzes the evidence. Source:
# DI-hohuf
script_dir="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
script_status="$?"
if [ "$script_status" -ne 0 ]; then
  printf 'failed to resolve script directory status=%s\n' "$script_status" >&2
  exit "$script_status"
fi
repo_dir="$(CDPATH= cd -- "$script_dir/.." && pwd)"
repo_status="$?"
if [ "$repo_status" -ne 0 ]; then
  printf 'failed to resolve repo directory status=%s\n' "$repo_status" >&2
  exit "$repo_status"
fi
config_path="${POC13_CONFIG:-./config.example.json}"
run_dir="/run/poc13/poc13-demo"

run_step() {
  step_name="$1"
  shift
  printf '\n== %s ==\n' "$step_name"
  "$@"
  status="$?"
  if [ "$status" -ne 0 ]; then
    printf 'step failed: %s status=%s\n' "$step_name" "$status" >&2
    exit "$status"
  fi
}

if [ ! -d "$repo_dir" ]; then
  printf 'missing repo dir: %s\n' "$repo_dir" >&2
  exit 1
fi

cd "$repo_dir"
cd_status="$?"
if [ "$cd_status" -ne 0 ]; then
  printf 'failed to enter repo dir: %s status=%s\n' "$repo_dir" "$cd_status" >&2
  exit "$cd_status"
fi
export POC13_CONFIG="$config_path"

printf 'repo: %s\n' "$repo_dir"
printf 'config: %s\n' "$POC13_CONFIG"
if [ -n "${OPENAI_API_KEY:-}" ]; then
  printf 'OPENAI_API_KEY: present\n'
else
  printf 'OPENAI_API_KEY: missing; scripted fallback will be used\n'
fi

run_step 'reset docker volume state' docker compose down -v --remove-orphans
run_step 'build and run poc13' docker compose up --build
run_step 'analyze poc13' docker compose run --rm --entrypoint /usr/local/bin/poc13-analyze alice-bob "$run_dir"
