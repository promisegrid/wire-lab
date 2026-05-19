# TODO-tapur: GA runner JSON fitness and child sim search

## Prior aliases

None. This TODO was minted after the proquint-handle migration.

## Status

Open. This TODO owns the fresh `tools/ga-runner` path for JSON fitness results,
direct child-sim generation under `simulations/`, GA/search orchestration,
review, promotion, and culling. It uses `TODO-dadub` as predecessor context for
root scenarios and result evidence, but does not reopen `TODO-dadub`. Source:
`DI-ramar`.

## Decision Intent Log

ID: DI-ramar
Date: 2026-05-19 10:05:55
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Build a fresh `tools/ga-runner` for GA/search work instead of extending
`tools/matrix-runner`; use one model per run; write new fitness evidence as JSON
at `results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.json`; generate
child sims directly as untracked `simulations/SIM-<handle>-<slug>/` trees; commit
only accepted children and selected JSON results; cull rejected children by
deleting their untracked sim tree and matching result tree.
Intent: The matrix-runner canary path proved useful, but GA/search needs a
cleaner result contract and child-sim workflow. Child sims must have the same
directory/file shape as parent sims, not a JSON proposal format. Fitness should
be the result, not a separate `results/fitness/` tree, and the v1 model should
both reason about the cell and emit the structured score.
Constraints: Preserve old Markdown canary result files but make `tools/ga-runner`
ignore them. Do not create a separate `results/fitness/` tree. Ordinary
population scans should treat committed/tracked sims as the stable population;
pending untracked children are included through the current GA run manifest.
Generated children are not accepted merely because they exist on disk. The
runner must make culling explicit so stale rejected children and their result
trees do not contaminate later runs.
Affects: `tools/ga-runner/`; `simulations/`; `results/`; `results/state/`;
`results/RUN-PROTOCOL.md`; `results/README.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`;
`protocols/wire-lab.d/TODO/TODO.md`.

ID: DI-zanon
Date: 2026-05-19 10:09:59
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Lock the first `tools/ga-runner` contracts before implementation:
JSON fitness results use `promisegrid.ga.result.v1`, GA state uses
`promisegrid.ga.state.v1` at `results/state/<run-group-id>.json`, and the v1
command surface is `init`, `score`, `generate`, `validate`, `progress`,
`accept`, and `cull`.
Intent: The GA runner needs a stable machine-readable result shape and run-state
shape before code exists, so child generation, scoring, culling, and later
review do not drift into ad hoc files. The command set separates population
setup, scoring, generation, validation, progress, acceptance, and destructive
culling without reusing the Markdown-oriented matrix-runner contract.
Constraints: Preserve `DI-ramar`: one model per run, JSON fitness evidence in
the normal `results/<sim>/<scenario>/<model>/<timestamp>.json` tree, generated
children as normal untracked `simulations/SIM-*` trees, no `results/fitness/`,
and old `.md` canary results ignored by GA-runner. `accept` records and reports
promotion candidates but does not auto-commit. `cull` deletes only generated
children named in the active GA state and their matching result trees.
Affects: `tools/ga-runner/`; `results/`; `results/state/`; `simulations/`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`;
future `results/RUN-PROTOCOL.md`; future `results/README.md`.

## Scope

- Define and implement a new GA/search runner without changing
  `tools/matrix-runner` as part of this TODO.
- Treat root scenarios and committed sims as the stable evaluation surface from
  `TODO-dadub`, while letting GA runs create untracked child sims as temporary
  candidates.
- Make JSON fitness result files the canonical output for GA runs.
- Keep old Markdown result files as historical canary evidence, outside the
  GA-runner input set.
- Plan child generation, scoring, review, promotion, and culling in one owner so
  the GA loop has no hidden side channel.

## Locked GA runner contracts

### JSON fitness result schema

GA-runner fitness results are JSON files at
`results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.json`. The
required schema ID is `promisegrid.ga.result.v1`.

Required top-level fields:

- `schema`: literal `promisegrid.ga.result.v1`.
- `result_id`: stable result identifier derived from sim, scenario, model, and
  timestamp.
- `run_group_id`: GA run group that produced the result.
- `cell_id`: state-file cell identifier for this sim/scenario/model/timestamp.
- `sim_id`, `scenario_id`, `model_id`, `timestamp_utc`, and `result_path`.
- `runner`: object containing `tool`, `provider`, `api_model`,
  `reasoning_effort`, `request_id`, `response_id`, and token/cost fields when
  available.
- `source`: object containing repo commit, sim path, scenario path,
  root contract paths, per-file `sha256` entries, and a simulation tree hash.
- `rubric`: object containing `rubric_version`, score scale, score meanings, and
  the ordered score axes.
- `scores`: object with integer `0..5` values for `scenario_fit`,
  `promisegrid_alignment`, `auditability`, `evolution_safety`,
  `layer_boundary_clarity`, `failure_handling`,
  `implementation_plausibility`, and `risk_penalty`.
- `fitness`: object containing `raw`, `normalized_0_100`, and `confidence_0_1`.
- `assessment`: object containing concise rationale, strengths, weaknesses,
  risks, open questions, and authority boundary.

The v1 result has one reasoning model. It does not contain separate combiner,
judge, evaluator, or scorer-model fields. If later work adds second-pass judging,
that requires a new schema version.

### GA state manifest

Each GA run has one state file at `results/state/<run-group-id>.json` using
schema `promisegrid.ga.state.v1`. This file is the authority for pending
untracked children and for safe culling.

Required top-level fields:

- `schema`, `run_group_id`, `created_at`, `updated_at`, `repo_commit`, and
  `model_id`.
- `population`: committed/tracked parent sims available at run initialization,
  each with sim ID, path, and tree hash.
- `scenario_sample`: scenario IDs chosen for this generation, sample policy, and
  source paths/hashes.
- `parents`: selected parent sim IDs and selection rationale.
- `children`: generated child sim IDs, paths under `simulations/SIM-*`, parent
  IDs, generation prompt hash, design-delta summary, tree hash, and status.
- `cells`: scoring cells with cell ID, sim ID, scenario ID, expected JSON result
  path, status, attempts, usage/cost fields, and validation message.
- `acceptance`: accepted child IDs, selected result paths, reviewer note, and
  acceptance timestamp.
- `culling`: culled child IDs, deleted sim paths, deleted result paths, cull
  reason, and cull timestamp.

Child statuses are `generated`, `queued`, `running`, `scored`, `accepted`,
`culled`, or `failed`. Cell statuses are `queued`, `running`, `done`, `failed`,
or `skipped`.

### `tools/ga-runner` command surface

The v1 CLI commands are:

- `init`: create `results/state/<run-group-id>.json`, discover the stable
  committed/tracked sim population, choose or record a scenario sample, and
  initialize parent-selection state.
- `score`: evaluate manifest cells with one model, write JSON result files under
  `results/<sim>/<scenario>/<model>/<timestamp>.json`, validate each result, and
  checkpoint state after every cell.
- `generate`: use selected parent sims and scenario pressure to write normal
  untracked child sim trees directly under `simulations/SIM-<handle>-<slug>/`,
  then record their paths and tree hashes in state.
- `validate`: validate GA state, child sim tree shape, JSON result path shape,
  schema fields, source hashes, and score ranges; ignore all `results/**/*.md`
  files.
- `progress`: print state counts, cost totals, generated children, scored cells,
  accepted children, and culled children.
- `accept`: verify selected child and result hashes, record acceptance in state,
  and print exact paths to stage for the normal repo commit workflow.
- `cull`: delete only generated child sim trees named in the active state file
  and their matching `results/<child-sim-id>/` trees, then record the cull action
  in state.

## Subtasks

- [x] tapur.1 Define the canonical JSON fitness result schema, including source
  paths, source hashes, model ID, rubric version, rubric scores, normalized
  fitness, rationale, risks, open work, and run metadata. Source: `DI-ramar`;
  `DI-zanon`.
- [x] tapur.2 Define the GA run manifest under `results/state/<run-group-id>.json`,
  including parent sim IDs, generated child sim IDs, child paths, scenario
  sample, expected JSON result paths, source hashes, statuses, accept state, and
  cull state. Source: `DI-ramar`; `DI-zanon`.
- [x] tapur.3 Specify `tools/ga-runner` commands for manifest generation, child
  generation, scoring, validation, progress/resume, accept, and cull. Source:
  `DI-ramar`; `DI-zanon`.
- [ ] tapur.4 Define child-generation prompts that produce normal
  `simulations/SIM-<handle>-<slug>/` trees with `README.md`, `QUESTION.md` when
  needed, optional `SCENARIOS.md`, optional local protocol/spec dirs, provenance
  back to parent sims, and bounded design deltas. Source: `DI-ramar`.
- [ ] tapur.5 Implement JSON-only fitness result writing and validation for
  `tools/ga-runner`, and make the runner ignore `results/**/*.md` canary files.
  Source: `DI-ramar`.
- [ ] tapur.6 Implement stable-population scanning so ordinary scans use
  committed/tracked `simulations/SIM-*` trees, while pending untracked children
  are included only through the active GA manifest. Source: `DI-ramar`.
- [ ] tapur.7 Implement conservative generation sizing: score existing sims,
  choose a small parent set, generate a small child batch, score each child
  against a stratified scenario sample, and promote at most a small number of
  children per generation. Source: `DI-ramar`.
- [ ] tapur.8 Implement review and promotion: accepted children are staged from
  their existing `simulations/SIM-*` paths and committed with selected JSON
  result evidence; rejected children remain uncommitted. Source: `DI-ramar`.
- [ ] tapur.9 Implement culling: rejected child sim trees and matching
  `results/<child-sim-id>/` trees are deleted only through an explicit cull
  command that records the action in the GA state file. Source: `DI-ramar`.
- [ ] tapur.10 Update `results/RUN-PROTOCOL.md`, `results/README.md`, and tool
  docs so GA-runner JSON results, Markdown canary-result exclusion, child-sim
  generation, review, promotion, and culling are documented from the same
  decision source. Source: `DI-ramar`.

## Predecessor context

- `TODO-dadub` owns the completed root scenario/result skeleton, scenario corpus,
  old canary/matrix-runner path, and source-of-truth decision that `results/`
  holds run evidence.
- `DI-moduf` requires real result-producing modes to use LLM or human reasoning
  rather than mechanical verdict synthesis.
- `DI-nuhon`, `DI-bujiv`, and `DI-lulom` established checkpointing, unattended
  execution, and provider-backed run lessons that `tools/ga-runner` should reuse
  conceptually without inheriting the old Markdown result contract.
- `DI-zamin` makes `results/` canonical evidence and generated views preferable
  to committed matrix summaries.
- `DI-nugiv` requires cost controls before large unattended API-backed runs.
- `DI-kizal` keeps scenario files compact by centralizing shared scenario
  boilerplate in `scenarios/README.md`, which should remain part of GA-runner
  prompt bundling.
