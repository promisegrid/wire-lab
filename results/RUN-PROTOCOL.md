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
- Update scenario `MATRIX.md` rows when a result is written.
- Keep old result files; never rewrite or delete prior runs.

## Failure and Retry Policy

- A cell is `failed` if required input docs are missing or result validation
  fails.
- Retries are explicit reruns with a new timestamp.
- Do not overwrite result files from a previous attempt.

## Tooling

- Manifest generator:
  `python3 results/tools/generate_matrix_manifest.py --models <model-id>`
- LLM job generator:
  `python3 results/tools/generate_llm_jobs.py --manifest <manifest.csv>`
- Result validator:
  `python3 results/tools/validate_results.py --model <model-id> --timestamp <ts> --strict-matrix`

## Recommended Preflight

1. Generate a canary manifest (20-30 cells).
2. Generate canary LLM job prompts.
3. Run canary prompts through an LLM evaluator.
4. Validate canary output.
5. Review drift/comparison report.
6. Run the full manifest.
