#!/usr/bin/env bash
set -euo pipefail

# Intent: Give operators one repeatable command for the POC14 clean regression:
# reset runtime state, run compose, discover the fresh evidence directory, and
# let poc14-analyze enforce the clean-run gates. Source: DI-jidah; DI-sinur;
# DI-punib; DI-sunuf
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
poc_dir="$(cd "${script_dir}/.." && pwd)"

cd "${poc_dir}"

echo "== reset POC14 docker runtime state =="
if docker compose down -v --remove-orphans; then
	echo "reset complete"
else
	status=$?
	echo "docker compose down failed with status ${status}" >&2
	exit "${status}"
fi

echo "== run POC14 compose regression =="
if docker compose up --build --abort-on-container-exit; then
	echo "compose run complete"
else
	status=$?
	echo "docker compose up failed with status ${status}" >&2
	exit "${status}"
fi

echo "== discover fresh POC14 run directory =="
if run_id="$(docker compose run --rm --entrypoint sh dave -c '
count=0
latest=""
for dir in /run/poc14/*; do
	if [ -f "${dir}/monitor-report.json" ]; then
		count=$((count + 1))
		latest="${dir##*/}"
	fi
done
if [ "${count}" -ne 1 ]; then
	echo "expected exactly one /run/poc14/<run_id> directory, found ${count}" >&2
	exit 1
fi
echo "${latest}"
')"; then
	echo "run id: ${run_id}"
else
	status=$?
	echo "run directory discovery failed with status ${status}" >&2
	exit "${status}"
fi

echo "== analyze POC14 clean regression =="
if docker compose run --rm --entrypoint /usr/local/bin/poc14-analyze dave "/run/poc14/${run_id}"; then
	echo "POC14 clean regression PASS"
else
	status=$?
	echo "POC14 clean regression FAIL with status ${status}" >&2
	exit "${status}"
fi
