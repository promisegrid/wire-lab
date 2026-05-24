#!/usr/bin/env bash
set -Eeuo pipefail

# Intent: Provide a terminal-friendly canary wrapper that streams observable
# progress to stdout while teeing the complete run transcript to a pasteable
# /tmp log file. Service tier is passed explicitly so the wrapper cannot inherit
# an expensive project/client default. Worker and timeout knobs are also passed
# explicitly so slow provider calls are bounded, while concise output is guided
# by text verbosity instead of hard output caps. The canary also requests
# reasoning summaries and mirrors streamed content to stdout/log for live
# diagnosis. `/tmp/canary-cells` can inject focused sim/scenario selections
# without editing `/tmp/canary.env`, but the run remains an explicit
# budget-governed canary invocation rather than a live queue. Source: DI-simag;
# DI-mopob; DI-juzus; DI-pulap; DI-tufud; DI-vadub; DI-pivuj; DI-suzor;
# DI-guvif; DI-bataj

usage() {
	cat <<'USAGE'
Usage:
  ./run-canary.sh [all|init|score-parents|breed|score-children|help]

Environment overrides:
  GA_CANARY_RUN_GROUP          default: ga-canary-<UTC timestamp>
  GA_CANARY_TIMESTAMP          default: current UTC YYYYMMDD-HHMMSS
  GA_CANARY_SHUFFLE_SEED       default: current UTC YYYYMMDDHHMMSS
  GA_CANARY_MODEL_ID           default: openai-gpt-5.4-medium
  GA_CANARY_API_MODEL          default: gpt-5.4
  GA_CANARY_SCORE_REASONING_EFFORT     default: medium
  GA_CANARY_GENERATE_REASONING_EFFORT  default: medium
  GA_CANARY_REASONING_SUMMARY  default: auto
  GA_CANARY_TEXT_VERBOSITY     default: low
  GA_CANARY_SERVICE_TIER       default: flex
  GA_CANARY_MAX_RUN_COST_USD   default: 10.00
  GA_CANARY_MAX_CELL_USD       default: 0.75
  GA_CANARY_MAX_CHILD_USD      default: 1.00
  GA_CANARY_INCLUDE_SIMS       optional comma/space list of SIM IDs to include
  GA_CANARY_INCLUDE_SCENARIOS  optional comma/space list of scenario IDs to include
  GA_CANARY_PARENT_COUNT       default: 3
  GA_CANARY_SCENARIO_COUNT     default: 3
  GA_CANARY_CHILD_COUNT        default: 2
  GA_CANARY_MAX_PROMOTIONS     default: 1
  /tmp/canary-cells            optional focus file with `sims:` / `scenarios:`
                               sections; entries resolve by unique prefix and
                               merge with GA_CANARY_INCLUDE_* values
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

Subcommands:
  all             default; run init, score-parents, breed, score-children, validate
  init            create the state file and planned cells
  score-parents   score only parent cells in an existing run-group
  breed           generate child simulations in an existing run-group
  score-children  score only child cells in an existing run-group
  help            show this usage text

For subcommands after init, set GA_CANARY_RUN_GROUP explicitly so the wrapper
reuses an existing state file instead of silently creating a fresh run-group.
USAGE
}

subcommand="${1:-all}"
case "$subcommand" in
	""|all|init|score-parents|breed|score-children)
		;;
	help|-h|--help)
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
# Intent: Keep the canary's default state/result lineage aligned with the
# scoring default. The calibration run showed medium scoring is the best
# broad-coverage default, while xhigh should be an explicit escalation path.
# Source: DI-nanor
model_id="${GA_CANARY_MODEL_ID:-openai-gpt-5.4-medium}"
api_model="${GA_CANARY_API_MODEL:-gpt-5.4}"
score_reasoning_effort="${GA_CANARY_SCORE_REASONING_EFFORT:-medium}"
generate_reasoning_effort="${GA_CANARY_GENERATE_REASONING_EFFORT:-medium}"
reasoning_summary="${GA_CANARY_REASONING_SUMMARY:-auto}"
text_verbosity="${GA_CANARY_TEXT_VERBOSITY:-low}"
service_tier="${GA_CANARY_SERVICE_TIER:-flex}"
max_run_cost_usd="${GA_CANARY_MAX_RUN_COST_USD:-10.00}"
max_cell_usd="${GA_CANARY_MAX_CELL_USD:-0.75}"
max_child_usd="${GA_CANARY_MAX_CHILD_USD:-1.00}"
include_sims="${GA_CANARY_INCLUDE_SIMS:-}"
include_scenarios="${GA_CANARY_INCLUDE_SCENARIOS:-}"
focus_file="/tmp/canary-cells"
# Intent: Let focused canaries shrink the parent pool so a known-underweighted
# candidate such as dalor can be forced to breed after a rubric/source repair.
# Source: DI-pozom
parent_count="${GA_CANARY_PARENT_COUNT:-3}"
scenario_count="${GA_CANARY_SCENARIO_COUNT:-3}"
child_count="${GA_CANARY_CHILD_COUNT:-2}"
max_promotions="${GA_CANARY_MAX_PROMOTIONS:-1}"
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

if [ "$subcommand" != "all" ] && [ "$subcommand" != "init" ] && [ -z "${GA_CANARY_RUN_GROUP:-}" ]; then
	echo "GA_CANARY_RUN_GROUP is required for subcommand '$subcommand'." >&2
	exit 2
fi

export GOCACHE="${GOCACHE:-/tmp/wire-lab-gocache}"

if [ "$subcommand" != "init" ] && [ -z "${OPENAI_API_KEY:-}" ]; then
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
if ! [[ "$parent_count" =~ ^[0-9]+$ ]] || [ "$parent_count" -lt 1 ]; then
	echo "GA_CANARY_PARENT_COUNT must be a positive integer." >&2
	exit 2
fi
if ! [[ "$scenario_count" =~ ^[0-9]+$ ]] || [ "$scenario_count" -lt 1 ]; then
	echo "GA_CANARY_SCENARIO_COUNT must be a positive integer." >&2
	exit 2
fi
if ! [[ "$child_count" =~ ^[0-9]+$ ]] || [ "$child_count" -lt 0 ]; then
	echo "GA_CANARY_CHILD_COUNT must be a non-negative integer." >&2
	exit 2
fi
if ! [[ "$max_promotions" =~ ^[0-9]+$ ]] || [ "$max_promotions" -lt 0 ]; then
	echo "GA_CANARY_MAX_PROMOTIONS must be a non-negative integer." >&2
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

append_unique_value() {
	local value="$1"
	local -n values_ref="$2"
	local existing
	for existing in "${values_ref[@]}"; do
		if [ "$existing" = "$value" ]; then
			return 0
		fi
	done
	values_ref+=("$value")
}

append_split_unique_values() {
	local raw_values="$1"
	local target_name="$2"
	local -n values_ref="$target_name"
	local normalized="${raw_values//,/ }"
	local value
	for value in $normalized; do
		if [ -n "$value" ]; then
			append_unique_value "$value" "$target_name"
		fi
	done
}

append_flag_args_from_array() {
	local flag_name="$1"
	local -n values_ref="$2"
	local -n output_args="$3"
	local value
	for value in "${values_ref[@]}"; do
		output_args+=("$flag_name" "$value")
	done
}

trim_line() {
	local text="$1"
	text="${text#"${text%%[![:space:]]*}"}"
	text="${text%"${text##*[![:space:]]}"}"
	printf '%s\n' "$text"
}

list_available_sims() {
	local path
	for path in "$repo_root"/simulations/SIM-*; do
		if [ -d "$path" ]; then
			basename "$path"
		fi
	done
}

list_available_scenarios() {
	local path scenario_id
	for path in "$repo_root"/scenarios/*; do
		if [ ! -d "$path" ]; then
			continue
		fi
		scenario_id="$(basename "$path")"
		if [ -f "$path/$scenario_id.md" ]; then
			printf '%s\n' "$scenario_id"
		fi
	done
}

resolve_unique_prefix() {
	local kind="$1"
	local selector="$2"
	shift 2
	local candidates=("$@")
	local candidate normalized exact_matches=() prefix_matches=()
	for candidate in "${candidates[@]}"; do
		normalized="$candidate"
		if [ "$kind" = "sim" ]; then
			normalized="${candidate#SIM-}"
		fi
		if [ "$candidate" = "$selector" ] || [ "$normalized" = "$selector" ]; then
			exact_matches+=("$candidate")
		fi
	done
	if [ "${#exact_matches[@]}" -eq 1 ]; then
		printf '%s\n' "${exact_matches[0]}"
		return 0
	fi
	for candidate in "${candidates[@]}"; do
		normalized="$candidate"
		if [ "$kind" = "sim" ]; then
			normalized="${candidate#SIM-}"
		fi
		if [[ "$candidate" == "$selector"* || "$normalized" == "$selector"* ]]; then
			prefix_matches+=("$candidate")
		fi
	done
	if [ "${#prefix_matches[@]}" -eq 1 ]; then
		printf '%s\n' "${prefix_matches[0]}"
		return 0
	fi
	if [ "${#prefix_matches[@]}" -eq 0 ]; then
		echo "Focus file selector '$selector' matched no $kind IDs." >&2
		return 1
	fi
	echo "Focus file selector '$selector' is ambiguous for $kind IDs: ${prefix_matches[*]}" >&2
	return 1
}

parse_focus_file() {
	local path="$1"
	if [ ! -f "$path" ]; then
		return 0
	fi
	# Intent: Resolve `/tmp/canary-cells` against the current repo before any
	# provider work starts so focused canaries fail fast on malformed or
	# ambiguous selectors. Source: DI-bataj
	local sims_name="$2"
	local scenarios_name="$3"
	local -n sims_ref="$sims_name"
	local -n scenarios_ref="$scenarios_name"
	local line raw section=""
	local resolved
	local available_sims=()
	local available_scenarios=()
	mapfile -t available_sims < <(list_available_sims)
	mapfile -t available_scenarios < <(list_available_scenarios)
	while IFS= read -r raw || [ -n "$raw" ]; do
		line="$(trim_line "$raw")"
		if [ -z "$line" ] || [[ "$line" == \#* ]]; then
			continue
		fi
		case "$line" in
			sims:|simulations:)
				section="sims"
				continue
				;;
			scenarios:)
				section="scenarios"
				continue
				;;
		esac
		if [[ "$line" == -* ]]; then
			line="$(trim_line "${line#-}")"
		fi
		if [ -z "$section" ]; then
			echo "Focus file $path is malformed: entry '$line' must appear under a sims: or scenarios: section." >&2
			return 1
		fi
		case "$section" in
			sims)
				resolved="$(resolve_unique_prefix sim "$line" "${available_sims[@]}")" || return 1
				append_unique_value "$resolved" "$sims_name"
				;;
			scenarios)
				resolved="$(resolve_unique_prefix scenario "$line" "${available_scenarios[@]}")" || return 1
				append_unique_value "$resolved" "$scenarios_name"
				;;
			*)
				echo "Focus file $path is malformed: unknown section '$section'." >&2
				return 1
				;;
		esac
	done < "$path"
}

init_include_sims=()
init_include_scenarios=()
focus_file_sims=()
focus_file_scenarios=()
append_split_unique_values "$include_sims" init_include_sims
append_split_unique_values "$include_scenarios" init_include_scenarios
parse_focus_file "$focus_file" focus_file_sims focus_file_scenarios
for value in "${focus_file_sims[@]}"; do
	append_unique_value "$value" init_include_sims
done
for value in "${focus_file_scenarios[@]}"; do
	append_unique_value "$value" init_include_scenarios
done
init_include_args=()
append_flag_args_from_array "-include-sim" init_include_sims init_include_args
append_flag_args_from_array "-include-scenario" init_include_scenarios init_include_args

log_dir="$(dirname "$log_file")"
mkdir -p "$log_dir"
if [ "$subcommand" = "all" ] || [ "$subcommand" = "init" ]; then
	: > "$log_file"
else
	touch "$log_file"
fi
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

run_init_stage() {
	run_step "init state" \
		go run . init \
			-repo-root "$repo_root" \
			-model "$model_id" \
			-run-group-id "$run_group" \
			-timestamp "$timestamp" \
			-shuffle-seed "$shuffle_seed" \
			-parent-count "$parent_count" \
			-scenario-count "$scenario_count" \
			-child-count "$child_count" \
			-max-promotions "$max_promotions" \
			"${init_include_args[@]}"
}

run_score_parents_stage() {
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
}

run_breed_stage() {
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
}

run_score_children_stage() {
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
}

run_validate_stage() {
	run_step "validate timestamp results" \
		go run . validate \
			-repo-root "$repo_root" \
			-model "$model_id" \
			-timestamp "$timestamp"
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
echo "Subcommand: $subcommand"
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
echo "Parent count: $parent_count"
echo "Scenario count: $scenario_count"
echo "Child count: $child_count"
echo "Max promotions: $max_promotions"
echo "Score request timeout: $score_request_timeout"
echo "Generate request timeout: $generate_request_timeout"
echo "Provider attempts: $provider_attempts"
echo "Provider elapsed: $provider_elapsed"
echo "Provider streaming: $stream"
echo "Stream idle timeout: $stream_idle_timeout"
echo "Stream content stdout: $stream_content_stdout"
if [ -f "$focus_file" ]; then
	echo "Focus file: $focus_file"
fi
if [ "${#focus_file_sims[@]}" -gt 0 ]; then
	echo "Focus-file sims: ${focus_file_sims[*]}"
fi
if [ "${#focus_file_scenarios[@]}" -gt 0 ]; then
	echo "Focus-file scenarios: ${focus_file_scenarios[*]}"
fi
if [ "${#init_include_sims[@]}" -gt 0 ]; then
	echo "Included sims: ${init_include_sims[*]}"
fi
if [ "${#init_include_scenarios[@]}" -gt 0 ]; then
	echo "Included scenarios: ${init_include_scenarios[*]}"
fi

if git -C "$repo_root" status --short | grep -q .; then
	echo "[warning] worktree has uncommitted or untracked files; continuing because canary outputs are expected to be uncommitted."
fi

cd "$ga_dir"

case "$subcommand" in
	all)
		run_init_stage
		run_score_parents_stage
		run_breed_stage
		run_score_children_stage
		run_validate_stage
		;;
	init)
		run_init_stage
		;;
	score-parents)
		run_score_parents_stage
		;;
	breed)
		run_breed_stage
		;;
	score-children)
		run_score_children_stage
		;;
	*)
		echo "unknown subcommand: $subcommand" >&2
		exit 2
		;;
esac

echo
echo "Canary command completed successfully."
echo "State file: $state_file"
echo "Log file: $log_file"
