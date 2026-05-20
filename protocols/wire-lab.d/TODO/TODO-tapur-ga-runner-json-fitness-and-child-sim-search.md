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

ID: DI-zohal
Date: 2026-05-19 10:11:37
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Lock the v1 child-generation contract for `tools/ga-runner`: generated
children are materialized directly as normal untracked
`simulations/SIM-<handle>-<slug>/` trees, the runner gives the model exact child
IDs and paths before generation, each generated child has `README.md` and
`QUESTION.md`, and the model may change only bounded design deltas derived from
the selected parent sims and scenario pressure.
Intent: Child proposals must be comparable to their parents as simulation
specimens, not stored as JSON proposal objects. The runner still needs a
machine-readable generation exchange so it can safely write files, hash them,
record provenance, and cull rejected children, but the durable candidate is the
materialized simulation tree.
Constraints: Do not commit generated children until accepted. Do not feed old
Markdown canary results into generation. Generation may use current GA state,
selected parent sim trees, selected scenario files, and JSON fitness evidence
from the active GA run. Generated children must include provenance back to
parent sims, the run group, source result paths when used, design deltas, and an
authority boundary. Child paths must be under `simulations/SIM-*`; generation
must not write into parent sim trees.
Affects: `tools/ga-runner/`; generated untracked `simulations/SIM-*` children;
`results/state/<run-group-id>.json`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-pobus
Date: 2026-05-19 10:18:12
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Implement the first `tools/ga-runner` scaffold as a standalone Go
module with JSON-only fitness result validation, atomic JSON result-writing
helpers for later scoring, `.md` result exclusion, and not-yet-implemented
stubs for the remaining locked command surface.
Intent: `tapur.5` should prove the JSON result contract in code before provider
calls, child generation, population scanning, accept, or cull behavior are
implemented. This keeps the first code pass small and prevents the old Markdown
canary result contract from leaking into GA selection.
Constraints: Do not modify `tools/matrix-runner`. Do not create real result
files, state files, or child sims in this pass. `ga-runner validate` discovers
only `results/<sim>/<scenario>/<model>/<timestamp>.json` files and ignores
`results/**/*.md`. The helper that writes JSON results is available for later
score work but is exercised only in tests during this pass.
Affects: `tools/ga-runner/`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-bagih
Date: 2026-05-19 10:25:10
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Implement `tools/ga-runner init -dry-run` as the first stable
population scanner: it discovers ordinary GA population members from
`git ls-files -z -- simulations`, groups tracked files by
`simulations/SIM-*`, computes deterministic tree hashes, and excludes untracked
child simulation directories by default.
Intent: Generated child sims are written directly under `simulations/`, so
ordinary scans must not treat every directory on disk as accepted population.
Using git-tracked files makes committed/tracked sims the default evaluation
surface while preserving the later ability to include pending children through
the active GA state file.
Constraints: This pass is read-only for GA runs: `init -dry-run` prints the
tracked population and does not write state. Missing/deleted tracked files and
non-`SIM-*` paths are ignored during population grouping. Pending child inclusion
through `results/state/<run-group-id>.json` remains later work.
Affects: `tools/ga-runner/`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-zusit
Date: 2026-05-19 10:35:07
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Implement conservative read-only generation planning in
`tools/ga-runner init -dry-run`: uniformly sample root scenarios, select a small
tracked parent set, plan a small child batch, cap promotions, and leave explicit
scenario tagging as later TODO work before serious GA runs.
Intent: The first GA planning pass should be cheap and deterministic without
pretending that `Source type` is a semantic tag system. Uniform sampling is good
enough for early scaffolding; a later scenario-tag pass can add domain,
pressure, layer, and risk metadata for serious search.
Constraints: Do not write GA state, result files, child sims, or scenario tags
in this pass. Defaults remain conservative: 3 parents, 5 scenarios, 4 children,
and 2 maximum promotions. Reject invalid counts. Scenario sampling is uniform
over `scenarios/<id>/<id>.md` entries and deterministic when `-shuffle-seed` is
provided.
Affects: `tools/ga-runner/`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`;
future scenario-tag TODO work.

ID: DI-podot
Date: 2026-05-19 10:43:01
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Implement `tools/ga-runner accept` as the v1 review/promotion
checkpoint: it reads `promisegrid.ga.state.v1`, verifies selected child sim
paths and selected JSON result evidence, records acceptance in the state file,
and prints exact repo-relative paths for the normal explicit `git add` and
commit workflow without staging or committing on its own.
Intent: Accepted children should be promoted from the same materialized
`simulations/SIM-*` trees that were scored, with selected JSON fitness results
as evidence. The tool should make the review boundary explicit and auditable
without turning existence on disk into acceptance and without bypassing the
repo's normal commit discipline.
Constraints: Reject missing or non-v1 state files, old Markdown canary results,
unknown children, culled children, child paths outside `simulations/SIM-*`,
child tree-hash drift, invalid JSON result files, and result evidence that does
not belong to a selected child. If the v1 state includes cells, selected result
paths must be present in those cells. `accept` may update only the selected
state file and must not create results, child sims, commits, or staged index
entries.
Affects: `tools/ga-runner/`; `results/state/<run-group-id>.json`;
`simulations/SIM-*`; `results/<sim>/<scenario>/<model>/<timestamp>.json`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-kofil
Date: 2026-05-19 10:52:50
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Implement `tools/ga-runner cull` as the explicit rejection cleanup
checkpoint: it reads `promisegrid.ga.state.v1`, verifies selected generated
children from the active state, optionally reports a dry-run plan, deletes only
the selected child sim trees and matching result trees, and records the cull
action in the state file.
Intent: Rejected child sims should not linger as ambiguous untracked candidates
or contaminate later GA work. Culling must still be explicit and state-bound, so
the tool cannot delete arbitrary `simulations/` or `results/` content and cannot
remove accepted children.
Constraints: Reject missing or non-v1 state files, unknown children, accepted
children, already-culled children, child paths outside exact
`simulations/<SIM-id>/`, and unsafe result paths. `cull -dry-run` must validate
and print the deletion plan without deleting files or writing state. Normal
culling may delete only `simulations/<SIM-id>/` and `results/<SIM-id>/` for
selected children, then append a culling record and set child status to `culled`.
Affects: `tools/ga-runner/`; `results/state/<run-group-id>.json`;
`simulations/SIM-*`; `results/SIM-*`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-ruzaj
Date: 2026-05-19 10:57:53
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Document the GA/search runner separately from legacy Markdown matrix
tooling: `results/` keeps both historical Markdown evidence and new JSON
fitness evidence, `tools/ga-runner` owns GA/search JSON result validation,
state-bound review, and culling, and matrix-runner/Python tooling remains
legacy/canary matrix-run support rather than the preferred GA/search path.
Intent: Operators need one clear place to see which result shape and toolchain
apply to GA/search without deleting historical canary evidence or misleading
future runs into using the Markdown matrix contract for JSON fitness work.
Constraints: Preserve old result files and legacy documentation context. Do not
claim unimplemented `ga-runner` modes are operational. Cite the active GA runner
DIs so the docs remain tied to the decision source for JSON fitness, state,
acceptance, and culling.
Affects: `results/README.md`; `results/RUN-PROTOCOL.md`;
`results/tools/README.md`; `tools/ga-runner/README.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-pibuh
Date: 2026-05-19 11:00:13
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Defer explicit scenario-tag design until after the project has some
real GA/search run results to inspect.
Intent: Early tagging would be premature taxonomy work. Initial run evidence
should reveal which domains, pressures, layers, and risks actually matter for
sampling and search, so the tag vocabulary is grounded in observed comparison
needs rather than speculation.
Constraints: Do not let scenario-tagging block the current GA runner path.
Revisit the tag families before serious GA runs or once enough initial
JSON-fitness results exist to guide the taxonomy.
Affects:
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`;
future scenario metadata work.

ID: DI-gijom
Date: 2026-05-19 11:10:59
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Implement the stateful GA runner loop: non-dry-run `init` writes
`promisegrid.ga.state.v1`, `score` evaluates state cells through one provider
and writes `promisegrid.ga.result.v1` JSON fitness evidence, and `generate`
creates untracked child simulation trees from provider-returned file bundles.
Intent: The GA/search runner needs to produce the first real JSON-fitness
evidence and generated child candidates without falling back to the legacy
Markdown matrix contract. The runner should own prompt construction, provider
calls, result validation, state checkpointing, usage/cost recording, and child
tree materialization so a long run can resume safely.
Constraints: V1 provider support is OpenAI-compatible Responses API only.
`score` asks the model for a score payload and the runner fills authoritative
identity/source/rubric fields. `generate` asks for a strict child file bundle and
rejects unsafe paths, missing `README.md`, missing `QUESTION.md`, parent-tree
writes, and malformed JSON. No scenario tags are added in this pass. Real API
calls require explicit operator invocation after implementation; tests use fake
providers only.
Affects: `tools/ga-runner/`; `results/state/<run-group-id>.json`;
`results/<sim>/<scenario>/<model>/<timestamp>.json`; `results/jobs/<run-group-id>/`;
generated untracked `simulations/SIM-*` child trees;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`;
`results/RUN-PROTOCOL.md`; `tools/ga-runner/README.md`.

ID: DI-simag
Date: 2026-05-19 11:38:51
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Add a terminal wrapper script for the first GA canary shape. The
wrapper runs state init, parent scoring, child generation, child scoring, and
result validation with the locked 3-parent, 3-scenario, 2-child canary defaults,
while streaming progress to stdout and teeing the full transcript to a
pasteable `/tmp/wire-lab-ga-canary-*.log` file.
Intent: The first provider-backed canary can take several minutes and may fail
mid-run due to provider output limits or validation errors. A wrapper gives
Steve a repeatable terminal command, visible progress from the checkpoint state,
and a single `/tmp` log filename that can be pasted back for review.
Constraints: Do not hide `ga-runner` failures; stop on the first failing phase
and print the state summary plus log path. Default to `gpt-5.3-codex`,
`xhigh`, run budget `$5.00`, cell estimate `$0.75`, child estimate `$1.00`,
shuffle seed `20260519`, and uncommitted canary artifacts. The wrapper may warn
about an already-dirty worktree but must not clean, stage, commit, accept, or
cull artifacts.
Affects: `tools/ga-runner/run-canary.sh`;
`results/state/ga-canary-*.json`; `results/jobs/ga-canary-*/`;
`results/<sim>/<scenario>/openai-gpt-5.3-codex-xhigh/<timestamp>.json`;
generated untracked `simulations/SIM-*-ga-child-*`;
`/tmp/wire-lab-ga-canary-*.log`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-mopob
Date: 2026-05-19 20:34:30
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Add an explicit `service-tier` control to provider-backed GA runner
commands and to the canary wrapper. `score` and `generate` default to
`-service-tier flex`; `default` is allowed only when requested explicitly;
`priority` and inherited `auto` behavior are rejected. Flex `429` and timeout
failures retry with bounded exponential backoff for at most five attempts within
a fifteen-minute retry window, with no automatic fallback to `default`.
Intent: Unattended GA/canary runs are cost-sensitive background workloads. They
must not accidentally inherit Priority or another expensive project/client
default, and Flex capacity failures should be handled as retryable transient
conditions rather than forcing Steve to babysit each cell.
Constraints: V1 support remains OpenAI-compatible Responses API only. The public
flag/env names are `-service-tier` and `GA_CANARY_SERVICE_TIER`. State/result
metadata names are `service_tier` for the requested tier and
`served_service_tier` for the provider-reported tier. Retry policy is bounded
Flex-only retry; switching to standard processing requires an explicit
operator-supplied `-service-tier default`.
Affects: `tools/ga-runner/`; `tools/ga-runner/run-canary.sh`;
`results/RUN-PROTOCOL.md`; `results/README.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-mokom
Date: 2026-05-19 20:48:23
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Change the terminal GA canary's default model from
`gpt-5.3-codex` / `openai-gpt-5.3-codex-xhigh` to `gpt-5.4` /
`openai-gpt-5.4-xhigh`.
Intent: The canary log at
`/tmp/wire-lab-ga-canary-ga-canary-20260519-204545.log` showed every parent
score cell failing with `Flex is not available for this model` while using
`gpt-5.3-codex`. The canary should remain a Flex-default unattended workload, so
its default model must be one intended for the Flex service tier rather than
requiring operators to override the service tier or model by hand.
Constraints: Preserve the existing `xhigh` reasoning default, 3x3x2 canary
shape, cost caps, checkpoint paths, `/tmp` log behavior, and explicit
`GA_CANARY_*` overrides. This supersedes only the canary model default recorded
in `DI-simag`; `DI-simag` remains active for the wrapper shape and operational
behavior.
Affects: `tools/ga-runner/run-canary.sh`; `tools/ga-runner/README.md`;
`results/RUN-PROTOCOL.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.
Supersedes: `DI-simag` canary model default only.

ID: DI-zikag
Date: 2026-05-19 20:55:19
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Add continuation controls for GA/canary scoring and generation:
retry empty or incomplete OpenAI Responses results through the existing bounded
provider retry loop, add `score -skip-failed-cells`, add
`generate -skip-failed-children`, and make the terminal canary pass both skip
flags so a partial provider anomaly does not stop the whole cycle.
Intent: The canary log at
`/tmp/wire-lab-ga-canary-ga-canary-20260519-184045.log` showed one parent score
cell ending with `openai response contained no output text` after most parent
cells had succeeded. The runner should preserve that failed-cell evidence, skip
the unusable cell after bounded retries, and keep going so child generation and
child scoring can exercise the full GA loop.
Constraints: Do not hide failures silently: skipped cells/children must keep
validation messages in state. Do not fallback from Flex to `default`. Do not
create synthetic fitness JSON for skipped cells. Child scoring should only select
generated or accepted child simulation trees so failed/skipped child-generation
plans do not produce missing-source failures.
Affects: `tools/ga-runner/`; `tools/ga-runner/run-canary.sh`;
`tools/ga-runner/README.md`; `results/RUN-PROTOCOL.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-juzus
Date: 2026-05-19 22:59:20
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Add bounded synchronous throughput controls to `tools/ga-runner`:
provider-backed `score` and `generate` get explicit worker counts, a five-minute
request timeout, a default two-attempt retry budget, and a six-minute retry
elapsed cap; the terminal canary opts into three scoring workers and one child
generation worker while keeping raw commands serial by default.
Intent: A 30-minute single synchronous provider wait makes the GA canary
operationally unusable and hides stalls behind repeated status counts. The sync
runner should provide fast bounded feedback now, while separate Batch-mode work
owns large asynchronous runs.
Constraints: Keep cost controls conservative under concurrency by reserving
estimated cost before launching provider calls. Do not let concurrent workers
write state unsafely. Preserve `-skip-failed-cells` and
`-skip-failed-children` continuation behavior. Do not implement OpenAI Batch in
this TODO pass.
Affects: `tools/ga-runner/`; `tools/ga-runner/run-canary.sh`;
`tools/ga-runner/README.md`; `results/RUN-PROTOCOL.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-pulap
Date: 2026-05-20 12:31:49
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Remove default provider hard output-token caps from `tools/ga-runner`
and the terminal canary. Keep budget controls by using separate estimate-only
output-token counts, send low text verbosity to the provider, and default child
generation to medium reasoning while keeping scoring at xhigh reasoning.
Intent: The canary showed child-generation calls consuming the entire output
cap as hidden reasoning tokens, producing `max_output_tokens` failures after
spending time and budget. The runner should guide concise JSON with soft output
shaping and prompt constraints, while preserving conservative preflight budget
estimates that do not alter provider behavior.
Constraints: Do not remove the explicit `-max-output-tokens` emergency fuse for
manual runs, but default it to omitted. Do not weaken `-max-run-cost-usd`,
`-max-cell-estimate-usd`, or `-max-child-estimate-usd`. Keep result scoring model
identity stable while recording generation reasoning effort in child state.
Structured Outputs are a separate follow-up decision after canary throughput is
healthy.
Affects: `tools/ga-runner/`; `tools/ga-runner/run-canary.sh`;
`tools/ga-runner/README.md`; `results/RUN-PROTOCOL.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-bukid
Date: 2026-05-20 12:56:40
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Before child generation, rank the selected parent pool by completed
parent fitness results, update queued child parent IDs with deterministic
top-parent plus tournament-diversity selection, and strengthen the child prompt
so the model is explicitly optimizing for a higher rubric score.
Intent: Passing parent score JSON into the prompt is useful but insufficient if
the runner still generates from the original uniform parent assignment and the
prompt does not tell the model to preserve strengths, repair weaknesses, reduce
risks, and improve `fitness.normalized_0_100`. The GA loop needs real selection
pressure while preserving diversity to avoid immediate convergence.
Constraints: Do not require a separate pre-generation state file shape or a full
population scoring pass in this change. If no completed parent score evidence is
available, preserve the existing child plan. Keep parent selection deterministic
from state/run inputs so interrupted runs can resume reproducibly.
Affects: `tools/ga-runner/`; `tools/ga-runner/README.md`;
`results/RUN-PROTOCOL.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-sohus
Date: 2026-05-20 13:22:28
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: For LLM-based GA child generation, replace separate `mutation` and
`crossover` operation labels with a single `breed` operation using exactly two
distinct parent simulation IDs.
Intent: LLM generation does not perform byte-level genetic mutation or crossover.
The useful operator is a two-parent design breeding prompt that asks the model to
use both parent simulations and their score evidence to produce one improved,
standalone child. One-parent children hide missing comparison pressure, while
three-or-more-parent prompts inflate context and blur design provenance.
Constraints: New child plans must use `breed` and two distinct parents. Existing
queued or running `mutation`/`crossover` state may be normalized during
generation; historical completed state and result evidence must not be rewritten.
If fewer than two viable parents are available, generation must fail or skip with
clear state evidence rather than silently creating a one-parent child.
Affects: `tools/ga-runner/`; `tools/ga-runner/README.md`;
`results/RUN-PROTOCOL.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-dilaf
Date: 2026-05-20 13:51:45
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Compact GA child-generation prompts by bundling parent simulation
documents once per parent, scenario pressure documents once per sampled
scenario, and score evidence as compact fitness summaries instead of full JSON
result files.
Intent: Canary child generation timed out before response headers while sending
large prompts that repeated root run/scenario boilerplate and embedded complete
parent result JSON. Child generation needs enough context to breed an improved
standalone sim, but it does not need full result source metadata, runner usage,
rubric boilerplate, or repeated root contracts. Compact prompts should reduce
latency and cost without weakening score evidence or parent design context.
Constraints: Keep score prompts source-complete. Do not rewrite historical result
files or current canary state. Keep parent sim docs complete for child
generation; compact only repeated root boilerplate and result evidence. Preserve
scenario-specific pressure, parent IDs, result paths, scores, fitness,
rationale, strengths, weaknesses, risks, and open questions in prompt form.
Affects: `tools/ga-runner/`; `tools/ga-runner/README.md`;
`results/RUN-PROTOCOL.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

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
  `reasoning_effort`, requested `service_tier`, provider-reported
  `served_service_tier`, `request_id`, `response_id`, and token/cost fields when
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
  IDs, generation prompt hash, design-delta summary, service-tier metadata, tree
  hash, and status.
- `cells`: scoring cells with cell ID, sim ID, scenario ID, expected JSON result
  path, status, attempts, service-tier metadata, usage/cost fields, and
  validation message.
- `acceptance`: accepted child IDs, selected result paths, reviewer note, and
  acceptance timestamp.
- `culling`: culled child IDs, deleted sim paths, deleted result paths, cull
  reason, and cull timestamp.

Child statuses are `generated`, `queued`, `running`, `scored`, `accepted`,
`culled`, `failed`, or `skipped`. Cell statuses are `queued`, `running`, `done`,
`failed`, or `skipped`.

### `tools/ga-runner` command surface

The v1 CLI commands are:

- `init`: create `results/state/<run-group-id>.json`, discover the stable
  committed/tracked sim population, choose or record a scenario sample, and
  initialize parent-selection state.
- `score`: evaluate manifest cells with one model, write JSON result files under
  `results/<sim>/<scenario>/<model>/<timestamp>.json`, validate each result, and
  checkpoint state after every cell. Provider-backed scoring sends explicit
  `-service-tier flex` by default; `default` requires explicit operator choice,
  and `priority` is rejected. `-skip-failed-cells` marks unusable cells as
  `skipped` after retries so later GA phases can continue.
- `generate`: use selected parent sims and scenario pressure to write normal
  untracked child sim trees directly under `simulations/SIM-<handle>-<slug>/`,
  then record their paths and tree hashes in state. Provider-backed child
  generation uses the same explicit service-tier policy as scoring.
  `-skip-failed-children` marks unusable child plans as `skipped` after retries
  so child scoring can proceed for generated children.
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

### Child-generation contract

`tools/ga-runner generate` creates untracked child simulation trees directly
under `simulations/`. The child tree, not a JSON proposal file, is the generated
candidate.

The runner must prepare each generation prompt with:

- exact child sim ID and target path under `simulations/SIM-<handle>-<slug>/`;
- selected parent sim IDs, parent paths, and parent tree hashes;
- selected scenario sample and scenario pressure summaries;
- relevant JSON fitness results from the active GA run when available;
- compact result-path, score, fitness, rationale, strength, weakness, risk, and
  open-question summaries instead of complete fitness-result JSON documents;
- required operation type: `breed`;
- exactly two distinct parent simulation IDs;
- a bounded design-delta budget of one to three substantive changes;
- a requirement that the child remain a standalone simulation tree.

The model's generation response may use a strict machine-readable file-bundle
envelope so the runner can write files deterministically, but that envelope is
only transport for generation. The durable child artifact is the materialized
simulation directory.

Each generated child must contain:

- `README.md`, describing the candidate design, parentage, design deltas, and
  authority boundary;
- `QUESTION.md`, stating the decision question the child simulation tests;
- optional `SCENARIOS.md`, only when the child adds simulation-local scenario
  pressure not already represented by root `scenarios/`;
- optional local protocol/spec directories when the design needs concrete local
  specimen files;
- provenance text naming parent sims, run group ID, generation model, source
  scenario sample, source JSON fitness results when used, and generation time.

Allowed generation operation:

- `breed`: combine two distinct parent simulations into one standalone child
  with one to three explicit design deltas and preserved provenance for both
  parents.

Forbidden generation operations:

- rewriting a parent in place;
- generating outside `simulations/SIM-*`;
- creating a broad "best of everything" child with no bounded deltas;
- importing old Markdown canary result prose as evidence;
- treating generated children as accepted merely because they exist on disk.

After writing a child tree, the runner records child ID, path, parent IDs,
operation type, prompt hash, response hash, per-file hashes, tree hash, and
status in the GA state file.

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
- [x] tapur.4 Define child-generation prompts that produce normal
  `simulations/SIM-<handle>-<slug>/` trees with `README.md`, `QUESTION.md` when
  needed, optional `SCENARIOS.md`, optional local protocol/spec dirs, provenance
  back to parent sims, and bounded design deltas. Source: `DI-ramar`;
  `DI-zohal`.
- [x] tapur.5 Implement JSON-only fitness result writing and validation for
  `tools/ga-runner`, and make the runner ignore `results/**/*.md` canary files.
  Source: `DI-ramar`; `DI-pobus`.
- [x] tapur.6 Implement stable-population scanning so ordinary scans use
  committed/tracked `simulations/SIM-*` trees, while pending untracked children
  are included only through the active GA manifest. Source: `DI-ramar`;
  `DI-bagih`.
- [x] tapur.7 Implement conservative generation sizing: score existing sims,
  choose a small parent set, generate a small child batch, score each child
  against a uniform scenario sample, and promote at most a small number of
  children per generation. Source: `DI-ramar`; `DI-zusit`.
- [x] tapur.8 Implement review and promotion: accepted children are staged from
  their existing `simulations/SIM-*` paths and committed with selected JSON
  result evidence; rejected children remain uncommitted. Source: `DI-ramar`;
  `DI-podot`.
- [x] tapur.9 Implement culling: rejected child sim trees and matching
  `results/<child-sim-id>/` trees are deleted only through an explicit cull
  command that records the action in the GA state file. Source: `DI-ramar`;
  `DI-kofil`.
- [x] tapur.10 Update `results/RUN-PROTOCOL.md`, `results/README.md`, and tool
  docs so GA-runner JSON results, Markdown canary-result exclusion, child-sim
  generation, review, promotion, and culling are documented from the same
  decision source. Source: `DI-ramar`; `DI-ruzaj`.
- [ ] tapur.11 Deferred until initial GA/search run results exist: add explicit
  scenario tags before serious GA runs. Candidate tag families: `domain` (for
  example logistics, governance, aviation, CAS, group-session, promisebase),
  `pressure` (sparse knowledge, adversarial trust, migration, auditability,
  naming, transport loss), `layer` (application, promise/accounting,
  group/session, CAS, envelope, transport), and `risk` (safety-critical,
  financial, governance, privacy, low-stakes). Source: `DI-zusit`; `DI-pibuh`.
- [x] tapur.12 Implement stateful non-dry-run `init`, provider-backed `score`,
  and provider-backed `generate` for the GA/search loop. Source: `DI-ramar`;
  `DI-zanon`; `DI-zohal`; `DI-gijom`.
- [x] tapur.13 Add a terminal canary wrapper that streams state progress to
  stdout and writes a pasteable `/tmp` transcript for review. Source:
  `DI-gijom`; `DI-simag`.
- [x] tapur.14 Add explicit service-tier controls and bounded Flex retry handling
  so GA/canary runs default to `flex`, reject `priority`, and never inherit an
  expensive tier by accident. Source: `DI-mopob`.
- [x] tapur.15 Change the terminal canary's default model to `gpt-5.4` /
  `openai-gpt-5.4-xhigh` after the `gpt-5.3-codex` canary failed because Flex
  was not available for that model. Source: `DI-mokom`.
- [x] tapur.16 Add retry/skip continuation for provider anomalies so the canary
  can finish parent scoring, child generation, child scoring, and validation even
  when individual cells or children are unusable after bounded retries. Source:
  `DI-zikag`.
- [x] tapur.17 Add bounded timeout, retry-budget, worker-count, and progress
  controls so sync GA canaries do not block for 30 minutes per provider call and
  can score multiple cells concurrently before Batch mode exists. Source:
  `DI-juzus`.
- [x] tapur.18 Remove default hard output caps from GA provider calls, add
  estimate-only output-token budgeting, send low text verbosity, and split the
  canary's score/generate reasoning defaults. Source: `DI-pulap`.
- [ ] tapur.19 Evaluate OpenAI Structured Outputs for GA score and child-bundle
  responses after the uncapped canary completes. Pros: schema-constrained JSON,
  fewer parser retries, and shorter formatting prompts. Cons: OpenAI-specific
  schema plumbing, stricter Markdown-in-JSON escaping, and no direct fix for
  hidden reasoning-token consumption. Source: `DI-pulap`.
- [x] tapur.20 Add fitness-ranked parent selection before child generation and
  strengthen child prompts so generated children are explicitly expected to
  improve over parent scores while preserving tournament diversity. Source:
  `DI-bukid`.
- [x] tapur.21 Compact child-generation prompts so breed calls include parent
  documents once, scenario pressure once, and summarized fitness evidence instead
  of repeated root boilerplate and full result JSON. Source: `DI-dilaf`.

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
  boilerplate in `scenarios/README.md`; score prompts keep using that root
  contract, while child-generation prompts use scenario-specific pressure and
  compact fitness evidence to avoid repeated boilerplate.
