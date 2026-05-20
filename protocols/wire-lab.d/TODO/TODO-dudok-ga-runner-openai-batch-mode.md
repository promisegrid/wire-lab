# TODO-dudok: GA runner OpenAI Batch mode and large-run orchestration

## Prior aliases

None. This TODO was minted after the proquint-handle migration.

## Status

Open. This TODO owns asynchronous OpenAI Batch support for large GA runs. It is
separate from `TODO-tapur`, which owns the current synchronous GA runner.

## Decision Intent Log

ID: DI-nizam
Date: 2026-05-19 22:59:20
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Track OpenAI Batch support as a separate harness TODO instead of
folding it into the synchronous `TODO-tapur` throughput fix.
Intent: Batch runs have different orchestration, persistence, retry, and ingest
boundaries than the sync canary path. Keeping them separate prevents the urgent
timeout/throughput fix from becoming a large asynchronous-run redesign.
Constraints: Preserve the GA JSON result and state contracts from `TODO-tapur`.
Do not revive the legacy Markdown matrix-runner contract. Final Batch command
shape, state schema extension, storage paths, and cost policy still require
later DF/DI before implementation.
Affects: `protocols/wire-lab.d/TODO/TODO-dudok-ga-runner-openai-batch-mode.md`;
`protocols/wire-lab.d/TODO/TODO.md`; future `tools/ga-runner/` Batch commands;
future `results/jobs/` or `results/batches/` Batch artifacts.

## Scope

This TODO covers OpenAI Batch support for `tools/ga-runner`: batch request
preparation, submission, polling, ingest, cancellation, retry/resume, docs, and
validation for large GA runs.

It does not reopen the old Markdown matrix-runner canary path and does not
replace the synchronous canary controls in `TODO-tapur`.

## Predecessor context

- `TODO-tapur` owns `tools/ga-runner`, JSON fitness results, GA state,
  provider-backed synchronous scoring and generation, child-sim materialization,
  review, promotion, and culling.
- `results/RUN-PROTOCOL.md` documents the current JSON run evidence contract.
- Batch execution should preserve the same durable outputs as sync mode:
  `results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.json`, generated
  child simulation trees, and `results/state/<run-group-id>.json`.

## DF/DI decisions needed before implementation

- Command shape: decide whether Batch support is a `batch` subcommand group
  such as `prepare`, `submit`, `poll`, `ingest`, and `cancel`, or separate
  top-level commands. Source: `DI-nizam`.
- State schema: decide whether Batch metadata extends
  `promisegrid.ga.state.v1` or requires a new state schema version. Source:
  `DI-nizam`.
- Batch artifact paths: decide whether JSONL inputs, provider output files, and
  ingest logs live under `results/jobs/<run-group-id>/batch/`,
  `results/batches/<run-group-id>/`, or another bounded results path. Source:
  `DI-nizam`.
- Prompt/source strategy: decide whether v1 Batch requests bundle the same
  source-complete prompts as sync scoring/generation, or whether a later
  server-side file strategy is worth designing first. Source: `DI-nizam`.
- Phase policy: decide how the runner records the dependency chain from parent
  scoring to child generation to child scoring, and how it prevents child phases
  from starting before prerequisite Batch outputs are ingested. Source:
  `DI-nizam`.
- Cost policy: decide how Batch pricing, estimates, usage metadata, and budget
  stop conditions interact with the existing sync cost controls. Source:
  `DI-nizam`.
- Failure policy: decide how failed, expired, missing, cancelled, malformed, or
  partially ingested Batch rows become retryable state without overwriting prior
  valid results or accepted child sims. Source: `DI-nizam`.

## Subtasks

- [ ] dudok.1 Run DF for Batch command shape, state schema extension, storage
  paths, prompt/source strategy, phase policy, and cost policy. Source:
  `DI-nizam`.
- [ ] dudok.2 Define the phase model for parent scoring, child generation, and
  child scoring; avoid one monolithic batch because child phases depend on prior
  phase outputs. Source: `DI-nizam`.
- [ ] dudok.3 Define deterministic JSONL request and `custom_id` mapping back to
  GA state cells and child-generation plans. Source: `DI-nizam`.
- [ ] dudok.4 Define Batch metadata in GA state, including batch IDs, file IDs,
  status, request counts, output files, error files, and retry lineage. Source:
  `DI-nizam`.
- [ ] dudok.5 Implement prepare, submit, poll, ingest, and cancel behavior after
  DF locks command and path decisions. Source: `DI-nizam`.
- [ ] dudok.6 Preserve sync-mode durable outputs during ingest: JSON fitness
  result files, state updates, generated child simulation trees, and validation
  messages. Source: `DI-nizam`.
- [ ] dudok.7 Add retry/resume semantics for failed, expired, or missing Batch
  rows without overwriting existing valid result files or accepted child sims.
  Source: `DI-nizam`.
- [ ] dudok.8 Document operator workflow, cost controls, large-run acceptance
  criteria, and how Batch differs from sync canary runs. Source: `DI-nizam`.

## Validation and acceptance criteria

- Batch preparation is deterministic for a fixed GA state and writes a JSONL
  request set whose line identifiers map back to exactly one GA state cell or
  child-generation plan. Source: `DI-nizam`.
- Batch submit, poll, ingest, and cancel behavior is covered by tests using
  fake provider responses; normal tests must not require network access. Source:
  `DI-nizam`.
- Ingest produces the same durable artifacts as sync mode: validated JSON
  fitness result files, checkpointed GA state updates, generated child sim
  trees, usage/cost metadata when available, and validation messages. Source:
  `DI-nizam`.
- Resume and retry never overwrite existing valid result files, accepted child
  simulation trees, or prior Batch attempt evidence. Source: `DI-nizam`.
- Batch-mode documentation explains the operator workflow, phase boundaries,
  cost controls, retry/resume behavior, and the distinction from the current
  synchronous canary path. Source: `DI-nizam`.
- Implementation validation should include `go test ./...`, `go vet ./...`,
  `errcheck ./...`, `bash -n tools/ga-runner/run-canary.sh`, and a dry-run or
  fake-provider Batch cycle that exercises prepare through ingest. Source:
  `DI-nizam`.
