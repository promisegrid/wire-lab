#!/usr/bin/env bash
set -Eeuo pipefail

# Intent: Provide a terminal-friendly canary wrapper that streams observable
# progress to stdout while teeing the complete run transcript to a pasteable
# /tmp log file. Service tier is passed explicitly so the wrapper cannot inherit
# an expensive project/client default. Worker and timeout knobs are also passed
# explicitly so slow provider calls are bounded, while concise output is guided
# by text verbosity instead of hard output caps. The canary also requests
# reasoning summaries and mirrors streamed content to stdout/log for live
# diagnosis. Source: DI-simag; DI-mopob; DI-juzus; DI-pulap; DI-tufud; DI-vadub;
# DI-pivuj; DI-suzor; DI-guvif

usage() {
	cat <<'USAGE'
Usage:
  ./run-canary.sh

Environment overrides:
  GA_CANARY_RUN_GROUP          default: ga-canary-<UTC timestamp>
  GA_CANARY_TIMESTAMP          default: current UTC YYYYMMDD-HHMMSS
  GA_CANARY_SHUFFLE_SEED       default: current UTC YYYYMMDDHHMMSS
  GA_CANARY_MODEL_ID           default: openai-gpt-5.4-xhigh
  GA_CANARY_API_MODEL          default: gpt-5.4
  GA_CANARY_SCORE_REASONING_EFFORT     default: xhigh
  GA_CANARY_GENERATE_REASONING_EFFORT  default: medium
  GA_CANARY_REASONING_SUMMARY  default: auto
  GA_CANARY_TEXT_VERBOSITY     default: low
  GA_CANARY_SERVICE_TIER       default: flex
  GA_CANARY_MAX_RUN_COST_USD   default: 10.00
	  GA_CANARY_MAX_CELL_USD       default: 0.75
	  GA_CANARY_MAX_CHILD_USD      default: 1.00
	  GA_CANARY_INCLUDE_SIMS       optional comma/space list of SIM IDs to include
	  GA_CANARY_INCLUDE_SCENARIOS  optional comma/space list of scenario IDs to include
	  GA_CANARY_SCORE_WORKERS      default: 6
  GA_CANARY_GENERATE_WORKERS   default: 1
  GA_CANARY_SCORE_REQUEST_TIMEOUT default: GA_CANARY_REQUEST_TIMEOUT or 5m
  GA_CANARY_GENERATE_REQUEST_TIMEOUT default: 15m
  GA_CANARY_REQUEST_TIMEOUT    legacy score timeout fallback
  GA_CANARY_PROVIDER_ATTEMPTS  default: 2
  GA_CANARY_PROVIDER_ELAPSED   default: 12m
  GA_CANARY_STREAM             default: true
  GA_CANARY_STREAM_IDLE_TIMEOUT default: 2m
  GA_CANARY_STREAM_CONTENT_STDOUT default: true
  GA_CANARY_POLL_SECONDS       default: 30
  GA_CANARY_LOG_FILE           default: /tmp/wire-lab-ga-canary-<run-group>.log

The script prints progress to stdout and writes the same transcript to the log.
It stops on the first ga-runner failure and prints the log filename.
USAGE
}

case "${1:-}" in
	"")
		;;
	-h|--help)
		usage
		exit 0
		;;
	*)
		usage >&2
		exit 2
		;;
esac

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$script_dir/../.." && pwd)"
ga_dir="$repo_root/tools/ga-runner"

timestamp="${GA_CANARY_TIMESTAMP:-$(date -u +%Y%m%d-%H%M%S)}"
shuffle_seed="${GA_CANARY_SHUFFLE_SEED:-$(date -u +%Y%m%d%H%M%S)}"
run_group="${GA_CANARY_RUN_GROUP:-ga-canary-$timestamp}"
model_id="${GA_CANARY_MODEL_ID:-openai-gpt-5.4-xhigh}"
api_model="${GA_CANARY_API_MODEL:-gpt-5.4}"
score_reasoning_effort="${GA_CANARY_SCORE_REASONING_EFFORT:-xhigh}"
generate_reasoning_effort="${GA_CANARY_GENERATE_REASONING_EFFORT:-medium}"
reasoning_summary="${GA_CANARY_REASONING_SUMMARY:-auto}"
text_verbosity="${GA_CANARY_TEXT_VERBOSITY:-low}"
service_tier="${GA_CANARY_SERVICE_TIER:-flex}"
max_run_cost_usd="${GA_CANARY_MAX_RUN_COST_USD:-10.00}"
max_cell_usd="${GA_CANARY_MAX_CELL_USD:-0.75}"
max_child_usd="${GA_CANARY_MAX_CHILD_USD:-1.00}"
include_sims="${GA_CANARY_INCLUDE_SIMS:-}"
include_scenarios="${GA_CANARY_INCLUDE_SCENARIOS:-}"
score_workers="${GA_CANARY_SCORE_WORKERS:-6}"
generate_workers="${GA_CANARY_GENERATE_WORKERS:-1}"
score_request_timeout="${GA_CANARY_SCORE_REQUEST_TIMEOUT:-${GA_CANARY_REQUEST_TIMEOUT:-5m}}"
generate_request_timeout="${GA_CANARY_GENERATE_REQUEST_TIMEOUT:-15m}"
provider_attempts="${GA_CANARY_PROVIDER_ATTEMPTS:-2}"
provider_elapsed="${GA_CANARY_PROVIDER_ELAPSED:-12m}"
stream="${GA_CANARY_STREAM:-true}"
stream_idle_timeout="${GA_CANARY_STREAM_IDLE_TIMEOUT:-2m}"
stream_content_stdout="${GA_CANARY_STREAM_CONTENT_STDOUT:-true}"
poll_seconds="${GA_CANARY_POLL_SECONDS:-30}"
log_file="${GA_CANARY_LOG_FILE:-/tmp/wire-lab-ga-canary-$run_group.log}"
state_file="$repo_root/results/state/$run_group.json"

export GOCACHE="${GOCACHE:-/tmp/wire-lab-gocache}"

if [ -z "${OPENAI_API_KEY:-}" ]; then
	echo "OPENAI_API_KEY is required for provider-backed canary runs." >&2
	exit 2
fi

if ! [[ "$poll_seconds" =~ ^[0-9]+$ ]] || [ "$poll_seconds" -lt 1 ]; then
	echo "GA_CANARY_POLL_SECONDS must be a positive integer." >&2
	exit 2
fi
if ! [[ "$score_workers" =~ ^[0-9]+$ ]] || [ "$score_workers" -lt 1 ]; then
	echo "GA_CANARY_SCORE_WORKERS must be a positive integer." >&2
	exit 2
fi
if ! [[ "$generate_workers" =~ ^[0-9]+$ ]] || [ "$generate_workers" -lt 1 ]; then
	echo "GA_CANARY_GENERATE_WORKERS must be a positive integer." >&2
	exit 2
fi
case "$stream" in
	true|false)
		;;
	*)
		echo "GA_CANARY_STREAM must be true or false." >&2
		exit 2
		;;
esac
case "$stream_content_stdout" in
	true|false)
		;;
	*)
		echo "GA_CANARY_STREAM_CONTENT_STDOUT must be true or false." >&2
		exit 2
		;;
esac
if ! [[ "$shuffle_seed" =~ ^[0-9]+$ ]]; then
	echo "GA_CANARY_SHUFFLE_SEED must be a decimal integer." >&2
	exit 2
fi
if ! [[ "$provider_attempts" =~ ^[0-9]+$ ]] || [ "$provider_attempts" -lt 1 ]; then
	echo "GA_CANARY_PROVIDER_ATTEMPTS must be a positive integer." >&2
	exit 2
fi

append_repeat_flags() {
	local flag_name="$1"
	local raw_values="$2"
	local -n output_args="$3"
	local normalized="${raw_values//,/ }"
	for value in $normalized; do
		output_args+=("$flag_name" "$value")
	done
}

init_include_args=()
append_repeat_flags "-include-sim" "$include_sims" init_include_args
append_repeat_flags "-include-scenario" "$include_scenarios" init_include_args

log_dir="$(dirname "$log_file")"
mkdir -p "$log_dir"
: > "$log_file"
exec > >(tee -a "$log_file") 2>&1

print_state_summary() {
	echo "[progress] log=$log_file"
	if [ ! -f "$state_file" ]; then
		echo "[progress] state not created yet: $state_file"
		return 0
	fi
	python3 - "$state_file" <<'PY'
import json
import sys
from collections import Counter
from datetime import datetime, timezone

state_path = sys.argv[1]
with open(state_path, encoding="utf-8") as state_file:
    state = json.load(state_file)

cells = state.get("cells", [])
children = state.get("children", [])
cell_counts = Counter(cell.get("status", "") for cell in cells)
child_counts = Counter(child.get("status", "") for child in children)
cost = sum(cell.get("cost_usd", 0) or 0 for cell in cells)
cost += sum(child.get("cost_usd", 0) or 0 for child in children)

print(f"[progress] state={state_path}")
print(f"[progress] cells={dict(sorted(cell_counts.items()))} children={dict(sorted(child_counts.items()))} cost_usd={cost:.6f}")

def age_text(updated_at):
    if not updated_at:
        return "unknown"
    try:
        value = updated_at.replace("Z", "+00:00")
        started = datetime.fromisoformat(value)
    except ValueError:
        return "unknown"
    seconds = max(0, int((datetime.now(timezone.utc) - started).total_seconds()))
    minutes, seconds = divmod(seconds, 60)
    return f"{minutes}m{seconds:02d}s"

running = [
    ("cell", cell.get("cell_id", ""), cell.get("updated_at", ""), cell.get("validation_message", ""))
    for cell in cells
    if cell.get("status") == "running"
]
running += [
    ("child", child.get("child_id") or child.get("sim_id", ""), child.get("updated_at", ""), child.get("validation_message", ""))
    for child in children
    if child.get("status") == "running"
]
if running:
    print("[progress] running:")
    for item_type, item_id, updated_at, message in running[:10]:
        print(f"[progress] - {item_type} {item_id} age={age_text(updated_at)} message={message}")

failures = [
    (cell.get("cell_id", ""), cell.get("validation_message", ""))
    for cell in cells
    if cell.get("status") == "failed"
]
failures += [
    (child.get("child_id") or child.get("sim_id", ""), child.get("validation_message", ""))
    for child in children
    if child.get("status") == "failed"
]
if failures:
    print("[progress] failures:")
    for item_id, message in failures[:10]:
        print(f"[progress] - {item_id}: {message}")
PY
}

run_step() {
	local label="$1"
	shift
	echo
	echo "== $label =="
	echo "+ $*"
	"$@"
	print_state_summary
}

run_step_with_monitor() {
	local label="$1"
	shift
	echo
	echo "== $label =="
	echo "+ $*"
	"$@" &
	local pid=$!
	while kill -0 "$pid" 2>/dev/null; do
		print_state_summary
		sleep "$poll_seconds"
	done
	local status=0
	if wait "$pid"; then
		status=0
	else
		status=$?
	fi
	print_state_summary
	return "$status"
}

on_error() {
	local status=$?
	local line="$1"
	echo
	echo "Canary failed at line $line with status $status."
	print_state_summary
	echo "Log file: $log_file"
	exit "$status"
}

trap 'on_error "$LINENO"' ERR

echo "GA canary run group: $run_group"
echo "Timestamp: $timestamp"
echo "Shuffle seed: $shuffle_seed"
echo "Repo root: $repo_root"
echo "Log file: $log_file"
echo "GOCACHE: $GOCACHE"
echo "Service tier: $service_tier"
echo "Score reasoning effort: $score_reasoning_effort"
echo "Generate reasoning effort: $generate_reasoning_effort"
echo "Reasoning summary: $reasoning_summary"
echo "Text verbosity: $text_verbosity"
echo "Score workers: $score_workers"
echo "Generate workers: $generate_workers"
echo "Score request timeout: $score_request_timeout"
echo "Generate request timeout: $generate_request_timeout"
echo "Provider attempts: $provider_attempts"
echo "Provider elapsed: $provider_elapsed"
echo "Provider streaming: $stream"
echo "Stream idle timeout: $stream_idle_timeout"
echo "Stream content stdout: $stream_content_stdout"
if [ -n "$include_sims" ]; then
	echo "Included sims: $include_sims"
fi
if [ -n "$include_scenarios" ]; then
	echo "Included scenarios: $include_scenarios"
fi

if git -C "$repo_root" status --short | grep -q .; then
	echo "[warning] worktree has uncommitted or untracked files; continuing because canary outputs are expected to be uncommitted."
fi

cd "$ga_dir"

run_step "init state" \
	go run . init \
		-repo-root "$repo_root" \
		-model "$model_id" \
		-run-group-id "$run_group" \
		-timestamp "$timestamp" \
			-shuffle-seed "$shuffle_seed" \
			-parent-count 3 \
			-scenario-count 3 \
			-child-count 2 \
			-max-promotions 1 \
			"${init_include_args[@]}"

run_step_with_monitor "score parent cells" \
	go run . score \
		-repo-root "$repo_root" \
		-run-group-id "$run_group" \
		-target parents \
		-api-model "$api_model" \
		-reasoning-effort "$score_reasoning_effort" \
		-reasoning-summary "$reasoning_summary" \
		-text-verbosity "$text_verbosity" \
		-service-tier "$service_tier" \
		-workers "$score_workers" \
		-request-timeout "$score_request_timeout" \
		-provider-max-attempts "$provider_attempts" \
		-provider-max-elapsed "$provider_elapsed" \
		-stream="$stream" \
		-stream-idle-timeout "$stream_idle_timeout" \
		-stream-content-stdout="$stream_content_stdout" \
		-skip-failed-cells \
		-max-run-cost-usd "$max_run_cost_usd" \
		-max-cell-estimate-usd "$max_cell_usd"

run_step_with_monitor "generate child simulations" \
	go run . generate \
		-repo-root "$repo_root" \
		-run-group-id "$run_group" \
		-api-model "$api_model" \
		-reasoning-effort "$generate_reasoning_effort" \
		-reasoning-summary "$reasoning_summary" \
		-text-verbosity "$text_verbosity" \
		-service-tier "$service_tier" \
		-workers "$generate_workers" \
		-request-timeout "$generate_request_timeout" \
		-provider-max-attempts "$provider_attempts" \
		-provider-max-elapsed "$provider_elapsed" \
		-stream="$stream" \
		-stream-idle-timeout "$stream_idle_timeout" \
		-stream-content-stdout="$stream_content_stdout" \
		-skip-failed-children \
		-max-run-cost-usd "$max_run_cost_usd" \
		-max-child-estimate-usd "$max_child_usd"

run_step_with_monitor "score child cells" \
	go run . score \
		-repo-root "$repo_root" \
		-run-group-id "$run_group" \
		-target children \
		-api-model "$api_model" \
		-reasoning-effort "$score_reasoning_effort" \
		-reasoning-summary "$reasoning_summary" \
		-text-verbosity "$text_verbosity" \
		-service-tier "$service_tier" \
		-workers "$score_workers" \
		-request-timeout "$score_request_timeout" \
		-provider-max-attempts "$provider_attempts" \
		-provider-max-elapsed "$provider_elapsed" \
		-stream="$stream" \
		-stream-idle-timeout "$stream_idle_timeout" \
		-stream-content-stdout="$stream_content_stdout" \
		-skip-failed-cells \
		-max-run-cost-usd "$max_run_cost_usd" \
		-max-cell-estimate-usd "$max_cell_usd"

run_step "validate timestamp results" \
	go run . validate \
		-repo-root "$repo_root" \
		-model "$model_id" \
		-timestamp "$timestamp"

echo
echo "Canary completed successfully."
echo "State file: $state_file"
echo "Log file: $log_file"
