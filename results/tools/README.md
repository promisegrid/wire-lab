# Results Tooling

This directory contains legacy Python matrix-run preflight and execution
helpers. The preferred implementation is now the Go CLI under
`tools/matrix-runner/`; these scripts remain as reference and rollback tools
until a later retirement decision. Source: `DI-lulom`.

## Preferred Go Runner

From `tools/matrix-runner/`:

```bash
go run . manifest -repo-root ../.. \
  -models openai-gpt-5.3-codex-xhigh

go run . run -repo-root ../.. \
  -manifest ../../results/manifests/<manifest>.csv \
  -provider openai \
  -api-model <openai-api-model> \
  -reasoning-effort xhigh

go run . validate -repo-root ../.. \
  -manifest ../../results/manifests/<manifest>.csv \
  -strict-matrix
```

The Go runner bundles local source documents for API-backed runs, checkpoints
under `results/state/`, and updates scenario matrices after validation. Source:
`DI-lulom`.

## Legacy Python Tools

## Generate Manifest

Generate full matrix for one model:

```bash
python3 results/tools/generate_matrix_manifest.py \
  --models openai-gpt-5.3-codex-xhigh
```

The manifest includes concrete `timestamp`, `result_path`, `ordinal`, and
`cell_id` fields so a long run can checkpoint and resume one cell at a time.
Source: `DI-nuhon`.

Generate deterministic canary manifest (30 cells):

```bash
python3 results/tools/generate_matrix_manifest.py \
  --models openai-gpt-5.3-codex-xhigh \
  --shuffle-seed 42 \
  --limit-cells 30
```

## Generate LLM Jobs

Generate prompt files from a manifest:

```bash
python3 results/tools/generate_llm_jobs.py \
  --manifest results/manifests/<manifest>.csv
```

The generator does not write final result verdicts. It writes prompt files under
`results/jobs/<run-group-id>/` for Codex or another LLM runner to evaluate.
Source: `DI-moduf`.

Generate only part of a manifest:

```bash
python3 results/tools/generate_llm_jobs.py \
  --manifest results/manifests/<manifest>.csv \
  --start-index 500
```

## Run Unattended Queue

Run a manifest through an external noninteractive LLM command:

```bash
python3 results/tools/matrix_queue.py run \
  --manifest results/manifests/<manifest>.csv \
  --runner-command '<llm-runner> {prompt_path}'
```

The queue writes one prompt per cell, invokes the runner command, validates the
result path, updates `scenarios/<scenario-id>/MATRIX.md`, and checkpoints
`results/state/<run-group-id>.json` after each cell. The runner command must
write the result file requested in the prompt; the queue does not generate
verdict prose. Source: `DI-nuhon`.

Resume the same queue after interruption:

```bash
python3 results/tools/matrix_queue.py run \
  --manifest results/manifests/<manifest>.csv \
  --runner-command '<llm-runner> {prompt_path}'
```

Inspect progress:

```bash
python3 results/tools/matrix_queue.py progress \
  --manifest results/manifests/<manifest>.csv
```

Runner-command placeholders:

- `{prompt_path}`
- `{result_path}`
- `{cell_id}`
- `{sim_id}`
- `{scenario_id}`
- `{model_id}`
- `{run_group_id}`

The command is split with `shlex`, not run through a shell. Put shell pipelines
or complex orchestration behind a small wrapper script when needed.

## Update Scenario Matrix Rows

Copy validated result evidence into the matching scenario matrix:

```bash
python3 results/tools/update_matrix_rows.py \
  --result results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md
```

## Validate Batch

Validate all results for a model/timestamp and enforce matrix links:

```bash
python3 results/tools/validate_results.py \
  --model openai-gpt-5.3-codex-xhigh \
  --timestamp <YYYYMMDD-HHMMSS> \
  --strict-matrix
```

Validate every row in a concrete manifest:

```bash
python3 results/tools/validate_results.py \
  --manifest results/manifests/<manifest>.csv \
  --strict-matrix
```

Prototype plumbing outputs are rejected by default. Use `--allow-prototype` only
when deliberately checking preserved scripted-output fixtures. Source:
`DI-moduf`.

## Compare Models

```bash
python3 results/comparisons/compare_model_results.py \
  --old-model openai-gpt-5.5-xhigh \
  --new-model openai-gpt-5.3-codex-xhigh
```

Comparison excludes scripted prototype outputs by default; pass
`--include-prototype` only for tooling audits. Source: `DI-moduf`.
