#!/usr/bin/env bash
set -euo pipefail

# Intent: Give operators one repeatable command for the POC16 clean regression:
# reset runtime state, run compose to natural completion, and let the
# observer-only collector provide analyzer input. Source: DI-jidah; DI-sinur;
# DI-punib; DI-sunuf; DI-dirat
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
poc_dir="$(cd "${script_dir}/.." && pwd)"

cd "${poc_dir}"

echo "== reset POC16 docker runtime state =="
if docker compose down -v --remove-orphans; then
	echo "reset complete"
else
	status=$?
	echo "docker compose down failed with status ${status}" >&2
	exit "${status}"
fi

echo "== run POC16 compose regression =="
if docker compose up --build; then
	echo "compose run complete"
else
	status=$?
	echo "docker compose up failed with status ${status}" >&2
	exit "${status}"
fi

echo "== analyze POC16 clean regression =="
if docker compose run --rm --entrypoint /usr/local/bin/poc16-analyze event-collector /run/poc16/poc16-demo; then
	echo "POC16 clean regression PASS"
else
	status=$?
	echo "POC16 clean regression FAIL with status ${status}" >&2
	exit "${status}"
fi
