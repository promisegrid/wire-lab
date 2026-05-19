# Results Run Protocol

This document defines the operational contract for matrix runs under
`results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md`.

## Purpose

The run protocol keeps full-matrix execution reproducible, auditable, and
comparable across models and timestamps.

## Required Inputs Per Cell

- `sim-id` and simulation docs:
  - `simulations/<sim-id>/README.md`
  - `simulations/<sim-id>/QUESTION.md` (if present)
  - simulation local spec docs (if present)
- `scenario-id` and scenario docs:
  - `scenarios/<scenario-id>/<scenario-id>.md`
  - `scenarios/<scenario-id>/README.md`
- `model-id` must be explicit and path-safe
  (for example `openai-gpt-5.3-codex-xhigh`).

## LLM Reasoning Requirement

Full matrix cells must be evaluated by an LLM or human reasoner reading the
cell's source docs. A script may prepare manifests, prompts, queues, and
validation reports, but a script must not generate final verdict prose by
mechanical parsing or keyword heuristics. Source: `DI-moduf`.

Full-matrix execution may be unattended: the queue runner may invoke one
external LLM command per manifest row, checkpoint state after every cell,
validate the result file, and update the scenario matrix. The external runner
must still produce the substantive result file; queue tooling only coordinates
work and validation. Source: `DI-nuhon`.

Allowed result-producing run modes:

- `codex-manual-blind`
- `llm-doc-eval-blind`

Excluded prototype mode:

- `scripted-doc-eval-blind`

Prototype mode may be used only to test plumbing. Prototype files must not be
treated as design evidence and must not be included in model comparison reports
unless a tool is explicitly invoked with a prototype opt-in flag.
Source: `DI-moduf`.

## Blindness Rules

- Blind mode must not read prior result files for the same cell before
  generating a new verdict.
- Blind mode must not copy old verdict text.
- Every blind run result file should include:
  - `- Run mode: codex-manual-blind` for manually produced Codex runs.
  - `- Run mode: llm-doc-eval-blind` for external LLM API/batch runs.

## Result File Contract

Each result must contain these sections:

- `## Result ID`
- `## Scenario`
- `## Simulation`
- `## Runner`
- `## Prompt / Procedure`
- `## Observed Behavior`
- `## Verdict`
- `## Evidence Links`
- `## Open Questions`
- `## Handoff Target`
- `## Authority Boundary`

Each result must include one line starting with `Evidence verdict:`.

## Batch Rules

- Generate a deterministic manifest before running a batch.
- Use one `run_group_id` for the manifest and one result timestamp for a
  single run invocation.
- Use manifests with concrete `timestamp`, `result_path`, `ordinal`, and
  `cell_id` fields for unattended runs.
- Use checkpoint state under `results/state/` for any long run that should
  resume without operator prompts.
- Update scenario `MATRIX.md` rows when a result is written.
- Keep old result files; never rewrite or delete prior runs.

## Failure and Retry Policy

- A cell is `failed` if required input docs are missing or result validation
  fails.
- Retries are explicit reruns with a new timestamp unless the same queue is
  intentionally resumed for a cell that never produced a valid result.
- Do not overwrite result files from a previous attempt.

## Tooling

- Manifest generator:
  `python3 results/tools/generate_matrix_manifest.py --models <model-id>`
- LLM job generator:
  `python3 results/tools/generate_llm_jobs.py --manifest <manifest.csv>`
- Unattended queue runner:
  `python3 results/tools/matrix_queue.py run --manifest <manifest.csv> --runner-command '<command with {prompt_path}>'`
- Matrix row updater:
  `python3 results/tools/update_matrix_rows.py --result <result.md>`
- Result validator:
  `python3 results/tools/validate_results.py --model <model-id> --timestamp <ts> --strict-matrix`

## Recommended Preflight

1. Generate a canary manifest (20-30 cells).
2. Generate canary LLM job prompts.
3. Run canary prompts through an LLM evaluator.
4. Validate canary output.
5. Review drift/comparison report.
6. Run the full manifest.

## Unattended Full-Run Shape

1. Generate a manifest with a fixed model ID and concrete timestamp:
   `python3 results/tools/generate_matrix_manifest.py --models <model-id>`.
2. Start the queue with an explicit noninteractive runner command:
   `python3 results/tools/matrix_queue.py run --manifest <manifest.csv> --runner-command '<runner> {prompt_path}'`.
3. Let the queue process one cell at a time. Each cell writes or refreshes its
   prompt under `results/jobs/<run-group-id>/`, invokes the runner command,
   validates `result_path`, updates the scenario `MATRIX.md`, and checkpoints
   `results/state/<run-group-id>.json`. Source: `DI-nuhon`.
4. If the process is interrupted, rerun the same command with the same manifest
   and state path. Completed cells are skipped by default.
5. When the queue completes, run strict validation over the manifest:
   `python3 results/tools/validate_results.py --manifest <manifest.csv> --strict-matrix`.
