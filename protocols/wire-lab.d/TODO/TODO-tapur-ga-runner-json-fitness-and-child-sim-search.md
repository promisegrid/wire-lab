# TODO-tapur: GA runner JSON fitness and child sim search

## Prior aliases

None. This TODO was minted after the proquint-handle migration.

## Status

Open. This TODO owns the fresh `tools/ga-runner` path for JSON fitness results,
proposal-stage child-sim generation under ignored `proposals/`, GA/search
orchestration, review, promotion, and culling. It uses `TODO-dadub` as
predecessor context for root scenarios and result evidence, but does not reopen
`TODO-dadub`. Source: `DI-ramar`; `DI-lirat`.

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

ID: DI-tufud
Date: 2026-05-20 14:20:19
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Strengthen synchronous GA provider reliability before the next
canary: parent selection uses one highest-scoring parent plus one deterministic
uniform random scored parent, OpenAI-compatible calls retry transient provider
failures including `408`, `409`, `429`, `500`, `502`, `503`, `504`, network
timeouts, and request timeouts, child/provider retry windows get enough elapsed
budget for a real second attempt, and Responses API streaming is enabled by
default with event-progress diagnostics and an idle timeout.
Intent: The canary logs showed child generation hanging or timing out without
useful liveness evidence, while parent selection still used tournament behavior
after Steve asked for one fit parent plus one random parent. The runner should
create visible progress, recover from ordinary transient provider failures, and
apply simple deterministic selection pressure without overfitting the second
parent choice.
Constraints: Keep `service_tier` explicit and do not fallback to Priority.
Preserve `max_output_tokens` as an opt-in emergency fuse only. Treat
`max_output_tokens` exhaustion as deterministic non-retryable unless the
operator changes the request shape. Keep the first parent sorted by completed
fitness evidence; choose the second parent uniformly from other scored parents
with a deterministic seed. Do not implement Batch in this change.
Affects: `tools/ga-runner/`; `tools/ga-runner/run-canary.sh`;
`tools/ga-runner/README.md`; `results/RUN-PROTOCOL.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.
Supersedes: `DI-bukid` tournament-diversity selection only; `DI-juzus`
six-minute elapsed retry default only; `DI-mopob` Flex-only transient retry
status list only.

ID: DI-vadub
Date: 2026-05-20 14:35:53
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Make the terminal GA canary request OpenAI reasoning summaries and
print streamed reasoning-summary and visible-output deltas to stdout/log when
streaming is enabled.
Intent: Steve wants to see live provider content instead of only event counters
while canary cells run. OpenAI does not expose raw reasoning tokens, so the
canary should request the supported `reasoning.summary` output and print those
summary deltas, plus visible output deltas, as diagnostic stream content.
Constraints: Do not claim or attempt to expose hidden raw reasoning tokens. Keep
stdout content opt-in for raw `score` and `generate` commands so their normal
status output stays readable; make the canary opt in by default. Preserve JSON
result parsing and state checkpointing semantics.
Affects: `tools/ga-runner/`; `tools/ga-runner/run-canary.sh`;
`tools/ga-runner/README.md`; `results/RUN-PROTOCOL.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-pivuj
Date: 2026-05-20 14:43:07
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Raise the terminal canary wrapper's default sync parallelism to six
score workers and two child-generation workers while keeping raw `score` and
`generate` command defaults serial.
Intent: The 2026-05-20 canary log showed parent scoring progressing normally but
too slowly: three score workers left parent cells queued while individual cells
completed in roughly one-to-three minutes. The canary should use more of the
available sync provider throughput before Batch mode exists, but raw commands
should remain conservative unless the operator opts into concurrency.
Constraints: Preserve the existing cost-reservation gate before concurrent
provider dispatch; do not increase the default run budget; keep worker counts
overridable by `GA_CANARY_SCORE_WORKERS` and `GA_CANARY_GENERATE_WORKERS`; keep
generation parallelism modest because children write simulation trees.
Affects: `tools/ga-runner/run-canary.sh`; `tools/ga-runner/README.md`;
`results/RUN-PROTOCOL.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-suzor
Date: 2026-05-20 14:47:50
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Keep the terminal canary's default child-generation worker count at
one while raising only scoring parallelism to six workers.
Intent: The 2026-05-20 canary evidence covered parent scoring only: parent cells
completed cleanly, queued work remained, and no provider retries or timeouts
occurred. The same run stopped before child generation because the wrapper
referenced an unset `reasoning_summary` variable in that historical execution.
Until a generation phase completes, raising generation parallelism would be an
untested change in the phase that writes child simulation trees.
Constraints: Keep `GA_CANARY_GENERATE_WORKERS` overridable for explicit tests;
use `GA_CANARY_SCORE_WORKERS=6` for the next canary; consider
`GA_CANARY_GENERATE_WORKERS=2` only after at least one successful generate phase
shows generation is safe and slow enough to justify parallel writes.
Affects: `tools/ga-runner/run-canary.sh`; `tools/ga-runner/README.md`;
`results/RUN-PROTOCOL.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.
Supersedes: `DI-pivuj` only for the default child-generation worker count.

ID: DI-babik
Date: 2026-05-20 14:55:52
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: For terminal canary stream-content output, do not print
`response.reasoning_summary_text.delta` event names or reasoning-summary text;
print one `.` with no newline for each reasoning-summary text delta event
instead.
Intent: Reasoning-summary text is useful for proving the provider is alive, but
printing the event name and full summary content makes canary stdout/logs noisy
and hard to scan. Progress dots preserve the liveness signal without exposing or
interleaving summary content with JSON output diagnostics.
Constraints: Keep visible `response.output_text.delta` mirroring available for
JSON-output diagnosis. Do not alter the parsed provider response text used for
result validation. Suppress reasoning-summary part text as well so no summary
content reaches the canary content stream.
Affects: `tools/ga-runner/openai.go`; `tools/ga-runner/ga_runner_test.go`;
`tools/ga-runner/README.md`; `results/RUN-PROTOCOL.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.
Supersedes: `DI-vadub` only for printing reasoning-summary event names and
summary text to stdout/log.

ID: DI-vajut
Date: 2026-05-20 14:58:21
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Keep `response.reasoning_summary_text.delta` output as no-newline
progress dots, but print `response.reasoning_summary_part.done` event names and
content to the canary stream-content output.
Intent: Delta events are too noisy when printed verbosely, but part-done events
are coarse enough to be useful diagnostic summaries. The canary should show
fine-grained liveness with dots while still preserving completed reasoning
summary parts for operator review.
Constraints: Do not print delta event names or delta content. Keep visible
`response.output_text.delta` mirroring unchanged. Do not append
reasoning-summary text to the parsed provider response used for result
validation.
Affects: `tools/ga-runner/openai.go`; `tools/ga-runner/ga_runner_test.go`;
`tools/ga-runner/README.md`; `results/RUN-PROTOCOL.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.
Supersedes: `DI-babik` only for suppressing
`response.reasoning_summary_part.done` event names and content.

ID: DI-sakam
Date: 2026-05-20 15:01:24
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Print both `response.reasoning_summary_part.added` and
`response.reasoning_summary_part.done` event names and content to canary
stream-content output, while keeping `response.reasoning_summary_text.delta`
events as no-newline progress dots only.
Intent: Part-added and part-done events are useful coarse-grained summary
markers, while text-delta events are too frequent and noisy. The canary should
show progress for deltas without printing their event names or text, and should
show summary-part boundaries and content when the provider emits them.
Constraints: Do not print reasoning-summary text-delta event names or delta
content. Keep visible `response.output_text.delta` mirroring unchanged. Do not
append reasoning-summary text to the parsed provider response used for result
validation.
Affects: `tools/ga-runner/openai.go`; `tools/ga-runner/ga_runner_test.go`;
`tools/ga-runner/README.md`; `results/RUN-PROTOCOL.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.
Supersedes: `DI-vajut` only for omitting
`response.reasoning_summary_part.added` event names and content.

ID: DI-guvif
Date: 2026-05-20 15:13:57
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Split terminal canary request-timeout defaults by phase: scoring uses
`GA_CANARY_SCORE_REQUEST_TIMEOUT` with a five-minute default, child generation
uses `GA_CANARY_GENERATE_REQUEST_TIMEOUT` with a fifteen-minute default, and
child generation keeps `GA_CANARY_GENERATE_REASONING_EFFORT=medium` by default
while scoring remains `xhigh`.
Intent: The 2026-05-20 canary evidence showed parent scoring finishing cleanly
under six `xhigh` workers and five-minute request attempts, while child
generation with an explicit `xhigh` override spent the whole timeout budget in
reasoning-summary events without emitting child-bundle output. Generation needs
a longer per-attempt window and a lower default reasoning effort without
weakening the scoring phase.
Constraints: Preserve legacy `GA_CANARY_REQUEST_TIMEOUT` as the score timeout
fallback for existing operator habits, but do not let it silently shorten the
new generation timeout unless `GA_CANARY_GENERATE_REQUEST_TIMEOUT` is set.
Leave provider elapsed and attempt count unchanged for now.
Affects: `tools/ga-runner/run-canary.sh`; `tools/ga-runner/README.md`;
`results/RUN-PROTOCOL.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-fupob
Date: 2026-05-20 15:17:11
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Print `response.reasoning_summary_part.done` event names and content
to canary stream-content output, but do not print
`response.reasoning_summary_part.added` event names or content.
Intent: Part-added events are repetitive structural markers and do not carry the
completed summary content Steve wants to inspect. Printing only part-done keeps
the useful coarse summary while reducing stdout/log noise during long child
generation calls.
Constraints: Keep `response.reasoning_summary_text.delta` as no-newline dots
only. Keep visible `response.output_text.delta` mirroring unchanged. Do not
append reasoning-summary text to the parsed provider response used for result
validation.
Affects: `tools/ga-runner/openai.go`; `tools/ga-runner/ga_runner_test.go`;
`tools/ga-runner/README.md`; `results/RUN-PROTOCOL.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.
Supersedes: `DI-sakam` only for printing
`response.reasoning_summary_part.added` event names and content.

ID: DI-ramun
Date: 2026-05-20 15:18:25
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Do not print `response.output_text.delta` event names or content to
canary stdout/log; keep those deltas only inside the internal response buffer
used for JSON parsing.
Intent: Visible-output deltas are high-volume and can interleave with progress
dots and monitor lines, making the canary transcript hard to scan. The canary
still needs output-text deltas for assembling the provider response, but the
operator-facing transcript should show only liveness dots and completed
reasoning-summary parts.
Constraints: Preserve response assembly and final JSON parsing. Keep
`response.reasoning_summary_text.delta` as no-newline dots only, keep
`response.reasoning_summary_part.done` event names and content, and keep
`response.reasoning_summary_part.added` quiet.
Affects: `tools/ga-runner/openai.go`; `tools/ga-runner/ga_runner_test.go`;
`tools/ga-runner/README.md`; `results/RUN-PROTOCOL.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-fihof
Date: 2026-05-20 16:16:53
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Treat generated GA children as ignored review-stage proposals named
`SIM-<handle>-child-<descriptive-slug>` with matching ignored child result
trees, rather than as commit-ready `SIM-<handle>-ga-child-<ordinal>` simulation
names.
Intent: The first canary-generated children demonstrated useful design moves,
but generic `ga-child` names and visible untracked result trees made review
state look like canonical simulation state. Generated children should be easy
to inspect and score, but should not be accidentally committed before a human
review/codex promotion pass gives them final simulation names and any missing
standing files.
Constraints: Keep unreviewed child proposals under `simulations/` so scoring
can reuse existing simulation source loading, but ignore `SIM-*-child-*` trees
and matching `results/SIM-*-child-*` evidence. The generator must propose a
descriptive child ID under the planned handle prefix. Existing generated
children from `ga-canary-20260520-221953` are renamed in-place as review-stage
children; their proposal result metadata may be rewritten because the proposal
artifacts remain uncommitted review staging evidence.
Affects: `.gitignore`; `tools/ga-runner/population.go`;
`tools/ga-runner/generate.go`; `tools/ga-runner/ga_runner_test.go`;
`tools/ga-runner/README.md`; `results/RUN-PROTOCOL.md`;
`simulations/SIM-bimos-child-grid-envelope-quarantine-sig-pcid-audit-tuple/`;
`simulations/SIM-jufag-child-grid-envelope-quarantine-sig-pcid-outcomes/`;
`results/SIM-bimos-child-grid-envelope-quarantine-sig-pcid-audit-tuple/`;
`results/SIM-jufag-child-grid-envelope-quarantine-sig-pcid-outcomes/`;
`results/state/ga-canary-20260520-221953.json`.

ID: DI-lirat
Date: 2026-05-20 16:28:30
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Stage generated GA child simulation trees and their child score
evidence under root `proposals/<run-group-id>/`, and ignore the entire root
`proposals/` tree until a human review/promotion pass moves selected artifacts
into canonical `simulations/` and `results/` homes.
Intent: Generated children are not canonical simulation specimens and generated
child score files are not canonical result evidence until reviewed. Keeping both
under one ignored run-scoped proposal tree prevents accidental commits, keeps
each child run reviewable as a unit, and avoids hiding accidental child writes
under canonical `simulations/` or `results/` paths.
Constraints: Parent score results remain under canonical `results/`; child
simulation proposals use
`proposals/<run-group-id>/simulations/<child-sim-id>/`; child score evidence uses
`proposals/<run-group-id>/results/<child-sim-id>/<scenario>/<model>/<timestamp>.json`.
Promotion remains a separate review step that creates final non-child `SIM-*`
names and canonical result paths as needed.
Affects: `.gitignore`; `tools/ga-runner/`; `tools/ga-runner/README.md`;
`results/RUN-PROTOCOL.md`; generated `proposals/<run-group-id>/`;
`results/state/<run-group-id>.json`.
Supersedes: DI-fihof; DI-podot

ID: DI-duzur
Date: 2026-05-20 16:47:07
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Let GA canary planning require specific parent simulations and
scenarios through explicit include lists, while filling the rest of the parent
and scenario sample by the existing deterministic shuffle.
Intent: Canary runs often need to exercise a newly-added simulation or a
specific scenario without abandoning small randomized coverage. Required
inclusions make focused canaries reproducible and prevent a new design point
from being missed by the uniform sample.
Constraints: Includes must validate against discovered tracked simulations and
root scenarios; duplicate includes are ignored; include counts may consume or
expand the effective sample size rather than silently dropping requested items.
The canary wrapper exposes comma/space-separated environment variables and the
`init` command exposes repeatable flags.
Affects: `tools/ga-runner/planning.go`; `tools/ga-runner/population.go`;
`tools/ga-runner/run-canary.sh`; `tools/ga-runner/README.md`;
`tools/ga-runner/ga_runner_test.go`; `results/RUN-PROTOCOL.md`.

ID: DI-dipid
Date: 2026-05-20 17:21:31
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Promote both `ga-canary-20260520-221953` Jufag and Bimos proposal
children as canonical non-child simulation specimens, with copied canonical JSON
score evidence and preserved proposal-source provenance.
Intent: Steve chose to promote both reviewed children rather than only the
higher-scoring Jufag variant, so the quarantine/profiled-signature-outcomes and
quarantine/audit-tuple variants can continue competing as independent simulation
specimens. Canonical result files should be discoverable under root `results/`,
but must not falsely claim that the LLM evaluated the post-promotion tree; the
scored proposal path remains the source evidence.
Constraints: Final simulation IDs drop only the `child` marker:
`SIM-jufag-grid-envelope-quarantine-sig-pcid-outcomes` and
`SIM-bimos-grid-envelope-quarantine-sig-pcid-audit-tuple`. Copy proposal trees and
result JSONs rather than moving or culling ignored proposal artifacts. Update
canonical result storage identity fields while preserving `source.*` fields that
point at the original proposal tree, and add explicit promotion metadata.
Affects: `proposals/ga-canary-20260520-221953/`;
`results/state/ga-canary-20260520-221953.json`;
`simulations/SIM-jufag-grid-envelope-quarantine-sig-pcid-outcomes/`;
`simulations/SIM-bimos-grid-envelope-quarantine-sig-pcid-audit-tuple/`;
`results/SIM-jufag-grid-envelope-quarantine-sig-pcid-outcomes/`;
`results/SIM-bimos-grid-envelope-quarantine-sig-pcid-audit-tuple/`;
`simulations/README.md`; `DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-dikoh
Date: 2026-05-21 13:33:26
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Store the repeatable GA child-promotion operator procedure in
`tools/ga-runner/PROMOTION.md` and use it when Steve says
`promote <child-proquint> [<child-proquint> ...]`.
Intent: Jufag and Bimos proved the promotion pattern, but the working procedure
was spread across state records, TODO provenance, canonical result metadata, and
operator memory. A dedicated procedure lets future promotions preserve the same
review, final naming, canonical result, and proposal-source provenance discipline
without re-deriving the steps each time.
Constraints: Promotion remains a review step, not automatic acceptance from disk.
The agent must resolve proquints through GA state, record a promotion DI before
canonical edits, copy proposal artifacts instead of moving or deleting them,
preserve `source.*` as the exact scored proposal evidence, add `promotion`
metadata to copied result JSON, update public indexes, and leave culling to an
explicit later command.
Affects: `tools/ga-runner/PROMOTION.md`; `tools/ga-runner/README.md`;
`results/RUN-PROTOCOL.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-lunuv
Date: 2026-05-22 00:10:12
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Promote `SIM-nonol-child-unknown-quarantine-blind-carry` from
`ga-canary-20260522-012332` as canonical
`SIM-nonol-unknown-quarantine-blind-carry`, and reject
`SIM-dubil-child-dag-cbor-opaque-wrapper-semantic-receipts` from the same run.
The `nonol` promotion keeps the proposal artifacts as scored provenance, copies
the proposal simulation tree into `simulations/`, copies all selected `nonol`
JSON fitness results into canonical `results/`, rewrites only canonical storage
identity fields plus promotion metadata, and preserves `source.*` fields as the
exact proposal tree scored by the LLM. The `dubil` rejection is recorded through
state-bound culling because its `semantic_id` mutation is unnecessary beside the
canonical `[pcid, payload]` envelope CID and adds avoidable protocol machinery.
Intent: Keep the promoted grid-envelope line aligned with pCID determinism,
Burgess Promise Theory simplicity, small-device constraints, and 100-year
durability. `nonol` preserves the useful middle path between opaque acceptance
and hard rejection: unknown pCIDs may be quarantined or blind-carried as exact
bytes without being treated as semantically accepted. `dubil` scored higher, but
its extra `semantic_id` layer is not accepted as current design consensus.
Constraints: Final canonical simulation ID is
`SIM-nonol-unknown-quarantine-blind-carry`. Do not move or rewrite scored
`nonol` proposal artifacts during promotion. Cull only the rejected `dubil`
proposal sim/result roots named by `results/state/ga-canary-20260522-012332.json`.
Update `results/state/ga-canary-20260522-012332.json`,
`simulations/README.md`, and `DEV-GUIDE-RESOURCES.md`; add canonical result
evidence under
`results/SIM-nonol-unknown-quarantine-blind-carry/<scenario>/openai-gpt-5.4-xhigh/20260522-012332.json`.
Affects: `results/state/ga-canary-20260522-012332.json`;
`proposals/ga-canary-20260522-012332/simulations/SIM-dubil-child-dag-cbor-opaque-wrapper-semantic-receipts/`;
`proposals/ga-canary-20260522-012332/results/SIM-dubil-child-dag-cbor-opaque-wrapper-semantic-receipts/`;
`proposals/ga-canary-20260522-012332/simulations/SIM-nonol-child-unknown-quarantine-blind-carry/`;
`proposals/ga-canary-20260522-012332/results/SIM-nonol-child-unknown-quarantine-blind-carry/`;
`simulations/SIM-nonol-unknown-quarantine-blind-carry/`;
`results/SIM-nonol-unknown-quarantine-blind-carry/`;
`simulations/README.md`; `DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-zadik
Date: 2026-05-22 00:15:37
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Change GA child promotion from copy-preserve to move-cleanup. Future
promotions move accepted proposal simulation trees and selected proposal result
trees from `proposals/<run-group-id>/` into canonical `simulations/` and
`results/` homes, then apply canonical cleanup in place. Existing promoted
children whose canonical simulations/results already exist should have their
accepted proposal simulation/result roots removed so `proposals/` contains only
unreviewed or still-reviewable proposals.
Intent: `proposals/` is ignored staging, not an archive. Keeping promoted child
trees in both proposal and canonical homes creates clutter and duplicate sources
of truth. The durable evidence after promotion is the canonical simulation tree,
canonical JSON result evidence, the result `source.*` hashes and historical
scored-source paths, the run state acceptance/culling records, and any job
prompts retained under `results/jobs/`.
Constraints: Do not remove unpromoted proposal children. Before cleaning an
accepted proposal root, verify the corresponding canonical simulation and
canonical result root exist. Existing canonical result `source.*` fields may
continue to name the historical proposal path and scored file hashes; those
fields are provenance metadata, not a promise that the ignored proposal path
still exists after cleanup. Promotion metadata must say whether artifacts were
moved or cleaned after earlier copy-style promotion.
Affects: `tools/ga-runner/PROMOTION.md`; `tools/ga-runner/README.md`;
`results/RUN-PROTOCOL.md`; `simulations/README.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`;
promoted accepted proposal roots under `proposals/ga-canary-*`.
Supersedes: DI-dikoh; DI-lunuv

ID: DI-higot
Date: 2026-05-22 10:57:43
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Treat scored artifact bytes as append-only. After a simulation README,
simulation-local specimen file, or result JSON has been scored or promoted as
scored evidence, do not rewrite that artifact in place to add cleanup metadata,
canonical-home notes, or other post-score edits. If promotion or cleanup needs
documentation, record it in TODOs, procedures, indexes, and surrounding docs
instead. Revert the current uncommitted scored-artifact edits, including the
uncommitted `SIM-nonol-unknown-quarantine-blind-carry` sim/result trees and the
post-score edits to already-scored promoted sim/result artifacts.
Intent: Preserve scored evidence as exact historical bytes. Post-score edits to
the artifact itself blur the distinction between what the LLM actually scored
and what later operators wanted the canonical tree to say about that score.
This repo can still evolve promotion procedure and provenance policy, but those
changes belong in surrounding apparatus docs rather than inside the scored
artifact bytes.
Constraints: Do not rewrite already-scored result JSONs or already-scored
simulation specimen/readme files to add provenance notes. If a promoted result
or specimen needs a stable home, move or copy the artifact bytes unchanged, or
leave them under `proposals/` until a non-mutating promotion path is agreed.
When reverting uncommitted scored-artifact changes, also fix any surrounding doc
references that would otherwise point at reverted uncommitted artifact paths.
Affects: `protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`;
`results/RUN-PROTOCOL.md`; `tools/ga-runner/PROMOTION.md`;
`tools/ga-runner/README.md`; `simulations/README.md`;
`DEV-GUIDE-RESOURCES.md`; currently modified `results/SIM-*/`;
currently modified `simulations/SIM-*/`; uncommitted
`results/SIM-nonol-unknown-quarantine-blind-carry/`; uncommitted
`simulations/SIM-nonol-unknown-quarantine-blind-carry/`.

ID: DI-fihub
Date: 2026-05-21 13:45:38
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Promote the Maraz, Natim, Savak, and Tizad GA child proposals as
canonical non-child simulation specimens with copied canonical JSON score
evidence and preserved proposal-source provenance.
Intent: Steve selected `maraz`, `natim`, `savak`, and `tizad` after the
proposal review/evaluation pass. Maraz and Natim extend the grid-envelope
arity/signature/proof exploration, while Savak and Tizad extend the guide-facing
claim/conformance evidence family. Promotion should make those specimens
discoverable and runnable without pretending the LLM scored the cleaned
canonical trees rather than the ignored proposal trees.
Constraints: Promote `SIM-maraz-child-signed-summary-header-nested-schema` from
`ga-canary-20260521-003110` as
`SIM-maraz-grid-envelope-signed-summary-header-nested-schema`; promote
`SIM-natim-child-nested-payload-outer-attestation-multisig` from
`ga-canary-20260521-003110` as
`SIM-natim-grid-envelope-nested-payload-outer-attestation-multisig`; promote
`SIM-savak-child-scoped-claim-card-audit-ledger` from
`ga-canary-20260521-011601` as
`SIM-savak-scoped-claim-card-audit-ledger`; promote
`SIM-tizad-child-scoped-conformance-citation-ledger` from
`ga-canary-20260521-011601` as
`SIM-tizad-scoped-conformance-citation-ledger`. Copy proposal trees and result
JSONs rather than moving or culling ignored proposal artifacts. Update canonical
result storage identity fields while preserving `source.*` fields that point at
the exact proposal tree scored by the LLM, and add explicit promotion metadata
with this DI. Add missing grid-envelope manifests to the promoted Maraz and
Natim canonical trees so they match comparable promoted grid-envelope specimens.
Affects: `proposals/ga-canary-20260521-003110/`;
`proposals/ga-canary-20260521-011601/`;
`results/state/ga-canary-20260521-003110.json`;
`results/state/ga-canary-20260521-011601.json`;
`simulations/SIM-maraz-grid-envelope-signed-summary-header-nested-schema/`;
`simulations/SIM-natim-grid-envelope-nested-payload-outer-attestation-multisig/`;
`simulations/SIM-savak-scoped-claim-card-audit-ledger/`;
`simulations/SIM-tizad-scoped-conformance-citation-ledger/`;
`results/SIM-maraz-grid-envelope-signed-summary-header-nested-schema/`;
`results/SIM-natim-grid-envelope-nested-payload-outer-attestation-multisig/`;
`results/SIM-savak-scoped-claim-card-audit-ledger/`;
`results/SIM-tizad-scoped-conformance-citation-ledger/`;
`simulations/README.md`; `DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-puhog
Date: 2026-05-21 13:54:46
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Replace GA child-generation parent selection from deterministic
top-parent plus uniform random scored parent to deterministic weighted-high plus
uniform random scored parent. The high parent uses linear rank weights over the
fitness-ranked scored parent pool, with weights `n, n-1, ..., 1` from best rank
to worst rank; the second parent remains uniformly selected from the other
scored parents using the existing deterministic seed inputs.
Intent: Steve asked for `high + random` instead of `top + random`: the first
parent should still be biased toward better completed parent fitness evidence,
but it should not always be the single highest-scoring parent. Linear rank
weighting is bounded, keeps every scored parent reachable, preserves
deterministic resumability from run/child/selection inputs, and avoids adding
public flags before canary evidence proves a need for tuning knobs.
Constraints: Preserve the existing parent ranking criterion: average completed
parent `fitness.normalized_0_100`, descending, with `SimID` as the tie-breaker.
Keep exactly two distinct `breed` parents. Preserve the second parent's uniform
random diversity over scored non-high parents. Use private helper names for the
implementation (`weightedHighParent` and `deterministicUint64`) and avoid new
public flags or state schema fields. Do not touch in-progress promotion
artifacts outside the GA-runner implementation, directly related tests/docs,
and this append-only DI.
Affects: `tools/ga-runner/generate.go`; `tools/ga-runner/ga_runner_test.go`;
`tools/ga-runner/README.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.
Supersedes: `DI-tufud` highest-scoring first-parent selection only.

ID: DI-lanuz
Date: 2026-05-21 16:07:41
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Promote the Virim and Zifoj GA child proposals as canonical non-child
guide-feedback simulation specimens with copied canonical JSON score evidence
and preserved proposal-source provenance.
Intent: Steve selected `virim` and `zifoj` after the
`ga-canary-20260521-210902` run. Virim explores a concrete App Manifest,
Embodiment Claim, Identity Continuity Receipt, and Savepoint Audit Envelope
profile. Zifoj scored highest in the run and explores a general boundary claim
package with scoped claim cards, field status tags, and peer-local
promise-accounting records. Promotion should make both specimens discoverable
without claiming that the LLM scored the cleaned canonical trees instead of the
ignored proposal trees.
Constraints: Promote `SIM-virim-child-manifested-embodiment-savepoint-receipts`
from `ga-canary-20260521-210902` as
`SIM-virim-manifested-embodiment-savepoint-receipts`; promote
`SIM-zifoj-child-boundary-claim-ledger` from `ga-canary-20260521-210902` as
`SIM-zifoj-boundary-claim-promise-accounting-records`. Copy proposal trees and
result JSONs rather than moving or culling ignored proposal artifacts. Update
canonical result storage identity fields while preserving `source.*` fields
that point at the exact proposal tree scored by the LLM, and add explicit
promotion metadata with this DI. Normalize canonical Zifoj prose away from the
ambiguous `ledger` term and toward peer-local promise-accounting records.
Affects: `proposals/ga-canary-20260521-210902/`;
`results/state/ga-canary-20260521-210902.json`;
`simulations/SIM-virim-manifested-embodiment-savepoint-receipts/`;
`simulations/SIM-zifoj-boundary-claim-promise-accounting-records/`;
`results/SIM-virim-manifested-embodiment-savepoint-receipts/`;
`results/SIM-zifoj-boundary-claim-promise-accounting-records/`;
`simulations/README.md`; `DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

## Scope

- Define and implement a new GA/search runner without changing
  `tools/matrix-runner` as part of this TODO.
- Treat root scenarios and committed sims as the stable evaluation surface from
  `TODO-dadub`, while letting GA runs create ignored proposal-stage child sims as
  temporary candidates.
- Make JSON fitness result files the canonical output for reviewed GA evidence,
  with child score evidence staged under ignored `proposals/<run-group-id>/`
  until promotion.
- Keep old Markdown result files as historical canary evidence, outside the
  GA-runner input set.
- Plan child generation, scoring, review, promotion, and culling in one owner so
  the GA loop has no hidden side channel.

## Locked GA runner contracts

### JSON fitness result schema

GA-runner parent fitness results are JSON files at
`results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.json`. Child score
evidence is the same JSON schema staged at
`proposals/<run-group-id>/results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.json`
until review/promotion. The required schema ID is `promisegrid.ga.result.v1`.

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
proposal children and for safe culling.

Required top-level fields:

- `schema`, `run_group_id`, `created_at`, `updated_at`, `repo_commit`, and
  `model_id`.
- `population`: committed/tracked parent sims available at run initialization,
  each with sim ID, path, and tree hash.
- `scenario_sample`: scenario IDs chosen for this generation, sample policy, and
  source paths/hashes.
- `parents`: selected parent sim IDs and selection rationale.
- `children`: generated child sim IDs, paths under
  `proposals/<run-group-id>/simulations/<child-sim-id>/`, parent IDs, generation
  prompt hash, design-delta summary, service-tier metadata, tree hash, and
  status.
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
  initialize parent-selection state. Optional `-include-sim` and
  `-include-scenario` flags guarantee focused coverage before deterministic
  shuffle fills the remaining sample slots. Source: `DI-duzur`.
- `score`: evaluate manifest cells with one model, write parent JSON result files
  under `results/<sim>/<scenario>/<model>/<timestamp>.json`, write child score
  evidence under `proposals/<run-group-id>/results/<sim>/<scenario>/<model>/<timestamp>.json`,
  validate each result, and checkpoint state after every cell. Provider-backed
  scoring sends explicit `-service-tier flex` by default; `default` requires
  explicit operator choice, and `priority` is rejected. `-skip-failed-cells`
  marks unusable cells as `skipped` after retries so later GA phases can
  continue.
- `generate`: use selected parent sims and scenario pressure to write normal
  proposal-stage child sim trees under
  `proposals/<run-group-id>/simulations/<SIM-id>/`, then record their paths and
  tree hashes in state. Provider-backed child generation uses the same explicit
  service-tier policy as scoring. `-skip-failed-children` marks unusable child
  plans as `skipped` after retries so child scoring can proceed for generated
  children.
- `validate`: validate GA state, child sim tree shape, JSON result path shape,
  schema fields, source hashes, and score ranges; ignore all `results/**/*.md`
  files.
- `progress`: print state counts, cost totals, generated children, scored cells,
  accepted children, and culled children.
- `accept`: verify selected proposal child and proposal result hashes, record
  acceptance in state, and print exact paths for the later promotion workflow.
- `cull`: delete only generated child sim trees named in the active state file
  and their matching `proposals/<run-group-id>/results/<child-sim-id>/` trees,
  then record the cull action in state.

### Child-generation contract

`tools/ga-runner generate` creates child simulation trees under ignored
`proposals/<run-group-id>/simulations/`. The child tree, not a JSON proposal
file, is the generated candidate.

The runner must prepare each generation prompt with:

- exact child sim ID and target path under
  `proposals/<run-group-id>/simulations/SIM-<handle>-child-<slug>/`;
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
- generating outside `proposals/<run-group-id>/simulations/SIM-*`;
- creating a broad "best of everything" child with no bounded deltas;
- importing old Markdown canary result prose as evidence;
- treating generated children as accepted merely because they exist on disk.

After writing a child tree, the runner records child ID, path, parent IDs,
operation type, prompt hash, response hash, per-file hashes, tree hash, and
status in the GA state file.

ID: DI-roruj
Date: 2026-05-22 12:15:37
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Add GA rubric/result v2 for promise-first vocabulary scoring, keep
historical v1 scored artifacts append-only, and use an audit-first targeted
backfill instead of an immediate full-corpus rescore. The new result schema is
`promisegrid.ga.result.v2`, the new rubric version is
`ga-rubric-20260522-v2`, the two new axes are `promise_vocabulary` and
`simplicity_durability`, and v2 scoring uses normal weighting across all axes.
The first targeted backfill audits canonical v1 results, treats sim/scenario
source-byte matches as the exact-match gate while allowing root-contract drift
to be reported separately, includes all exact-match hard-hit sims, adds a clean
grid-envelope calibration slice, preserves original model IDs by default, and
writes only new timestamped v2 result files.
Intent: Full-corpus rescoring is expensive, and the new vocabulary rules are
meant to change scoring going forward without falsifying historical evidence.
The runner therefore needs a distinct v2 scoring contract, a cheap audit to
find the sims most likely to move, and a targeted backfill path that can be
reviewed before any broader rerun is justified.
Constraints: Do not rewrite or delete any scored v1 sim/result bytes. Keep
`promisegrid.ga.state.v1` as the current run-state schema. Preserve v1
validation for historical evidence. Audit current sim/scenario bytes, not the
mutated v2 root-contract docs, when deciding whether a historical result is an
exact-match backfill candidate. Keep comparison/reporting follow-on work open
until targeted v2 evidence exists.
Affects: `tools/ga-runner/`; `results/RUN-PROTOCOL.md`;
`tools/ga-runner/README.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-bataj
Date: 2026-05-22 19:35:22
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Add `/tmp/canary-cells` as an optional plain-text focus file for
`tools/ga-runner/run-canary.sh`. The file contains `sims:` and/or
`scenarios:` sections, one selector per line, ignores blank lines and `#`
comments, resolves selectors by unique prefix expansion against current
`simulations/SIM-*` and root scenario IDs, merges the resolved IDs with any
`GA_CANARY_INCLUDE_*` values, and fails fast on malformed, missing, or
ambiguous selectors before any provider calls. Do not add a live scoring queue
or change the canary budget model.
Intent: Targeted canary slices should be easy to edit without mutating
`/tmp/canary.env`, but unattended spending must remain attached to explicit
canary invocations rather than an always-on queue.
Constraints: Implement this in `tools/ga-runner/run-canary.sh` only. Keep
existing budget, service-tier, worker, and timeout controls unchanged. Use the
fixed focus-file path `/tmp/canary-cells`. Merge resolved file entries with the
existing `GA_CANARY_INCLUDE_SIMS` / `GA_CANARY_INCLUDE_SCENARIOS` values rather
than replacing them.
Affects: `tools/ga-runner/run-canary.sh`; `tools/ga-runner/README.md`;
`results/RUN-PROTOCOL.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-zobur
Date: 2026-05-22 20:07:19
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Make `tools/ga-runner audit` and `backfill-init` promotion-aware
without rewriting historical scored artifacts. When a canonical promoted result
still preserves historical `source.sim_path` and sim-root source-file paths
under a deleted `proposals/...` tree, audit may fall back to the current
canonical `simulations/<sim-id>/` tree, but it counts as `exact_match` only if
the canonical sim files and tree hash still match the historical scored bytes
exactly. Missing historical proposal roots must not abort audit; unresolved
sources become non-exact and are excluded from targeted backfill.
Intent: Historical `source.*` provenance must stay append-only, but targeted
rubric-v2 rescoring is blocked unless the runner can compare promoted canonical
results against the byte-identical canonical sim trees that replaced deleted
proposal trees.
Constraints: Do not rewrite any scored result JSON or simulation bytes. Keep
`source.*` proposal provenance authoritative. Use canonical fallback only for
current audit/backfill comparison. Report source resolution in `audit` output.
Keep `tapur.36` as a later comparison-report step after actual v2 results
exist.
Affects: `tools/ga-runner/result.go`; `tools/ga-runner/validate.go`;
`tools/ga-runner/ga_runner_test.go`; `tools/ga-runner/README.md`;
`results/RUN-PROTOCOL.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-hijub
Date: 2026-05-22 20:23:22
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Let `backfill-init` optionally override the staged rescoring
`model_id` and default `reasoning_effort` for the generated v2 state while
preserving historical model lineage by default when no override is provided.
This is the path for honest multi-stage rescoring such as a `high` first pass
followed by an `xhigh` pass with distinct result paths and state files.
Intent: Two-stage rescoring is useful, but the generated state/result paths must
not keep historical `openai-gpt-5.4-xhigh` model IDs when a staged pass is
actually being scored at `high` or another effort. Stage overrides need to be
explicit, additive, and auditable.
Constraints: Default behavior stays unchanged: preserve original model IDs and
derived API models when no stage override is provided. Do not rewrite existing
state or result files. Result paths, `cell.model_id`, and `state.model_id` must
match the staged override when one is provided. `cell.api_model` should derive
from the staged `model_id` unless explicitly overridden later by the scoring
command.
Affects: `tools/ga-runner/population.go`; `tools/ga-runner/ga_runner_test.go`;
`tools/ga-runner/README.md`; `results/RUN-PROTOCOL.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-sirih
Date: 2026-05-22 21:53:06
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Treat the 5-cell `gpt-5.4` reasoning-effort calibration as evidence
that broad GA parent scoring should prefer `medium` by default, with `xhigh`
reserved for tie-breaks, promotion candidates, and design-state-sensitive
comparisons, but do not change runner defaults until that default-change choice
is explicitly implemented.
Intent: The calibration slice showed `medium` matching `xhigh` much better than
`low` or `high` while costing far less than `xhigh`. Record the evidence now so
the repo has a durable recommendation, without silently changing live runner
behavior in the same step.
Constraints: Calibration evidence is the 5-cell run-group family
`ga-calib-20260523-033216-{low,medium,high,xhigh}` under `results/state/` and
their matching result trees. Do not change `tools/ga-runner` defaults under this
DI alone. A later explicit implementation step may update defaults or canary
wrapper policy after review.
Affects: `protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`;
`results/state/ga-calib-20260523-033216-low.json`;
`results/state/ga-calib-20260523-033216-medium.json`;
`results/state/ga-calib-20260523-033216-high.json`;
`results/state/ga-calib-20260523-033216-xhigh.json`.

ID: DI-nanor
Date: 2026-05-22 21:54:04
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Change the default GA parent scoring effort from `xhigh` to `medium`
in `tools/ga-runner`, and change the canary wrapper's default scoring `model_id`
to `openai-gpt-5.4-medium` so the emitted state/result lineage matches the new
default. Keep `xhigh` as an explicit escalation path for tie-breaks, promotion
candidates, and design-state-sensitive comparisons.
Intent: The 5-cell calibration run showed `medium` matching the `xhigh` ranking
materially better than `low` or `high` while costing far less than `xhigh`.
Make the cheaper default real in runner behavior and keep the higher-effort path
explicit rather than implicit.
Constraints: Do not change child generation defaults. Do not remove `xhigh`
support. Update code comments, docs, and default-expectation tests together so
the new default is explicit and auditable.
Affects: `tools/ga-runner/provider.go`; `tools/ga-runner/run-canary.sh`;
`tools/ga-runner/README.md`; `tools/ga-runner/ga_runner_test.go`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-zuzup
Date: 2026-05-22 22:18:11
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Add `ga-runner compare-backfill` to generate targeted rubric-v2
v1-vs-v2 comparison reports. The command writes a derived Markdown report under
`results/reports/<run-group-id>-comparison.md` by default and compares each
completed v2 backfill cell against the latest exact-match canonical
`promisegrid.ga.result.v1` result for the same `sim_id` + `scenario_id`,
preferring the same `runner.api_model` when available.
Intent: `tapur.36` needs a durable, reviewable artifact that shows whether the
vocabulary-aware rubric actually changes rankings before broader rescoring is
authorized.
Constraints: Do not rewrite any scored artifact bytes. Treat the report as
derived evidence only. Include sim-rank drift, family highlights for
grid-envelope and conformance-family sims when present, and the largest
per-cell score deltas. Report unmatched and ambiguous historical pairings
instead of hiding them.
Affects: `tools/ga-runner/main.go`; `tools/ga-runner/compare.go`;
`tools/ga-runner/ga_runner_test.go`; `tools/ga-runner/README.md`;
`results/RUN-PROTOCOL.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-guhar
Date: 2026-05-22 22:41:54
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Deduplicate targeted backfill selection by `sim_id` +
`scenario_id` before writing the v2 state. When multiple exact-match canonical
v1 candidates survive for the same pair, keep one deterministic winner using
the same ordering preference as `compare-backfill`: prefer a record with
`runner.api_model`, then newer `timestamp_utc`, then lexical `model_id`, then
path. Future `backfill-init` states should not emit repeated v2 cells that
write the same result path.
Intent: The first live comparison report surfaced repeated state rows for the
same sim/scenario pair, which inflated ambiguity and made the targeted backfill
state noisier than the actual result corpus.
Constraints: Do not rewrite historical v1 or completed v2 result bytes. Keep
the ambiguity visible in `compare-backfill` when the historical corpus still
contains multiple exact-match v1 candidates; this DI only deduplicates future
backfill selection/state materialization.
Affects: `tools/ga-runner/validate.go`; `tools/ga-runner/ga_runner_test.go`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

ID: DI-sirir
Date: 2026-05-22 23:33:34
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: `compare-backfill` should collapse same-model historical reruns after
`runner.api_model` narrowing, keep the latest rerun deterministically, and
report the collapsed rerun groups separately from true ambiguous historical
pairings. A historical rerun group is a set of exact-match canonical
`promisegrid.ga.result.v1` records for one `sim_id` + `scenario_id` that share
the same `model_id` and `runner.api_model`.
Intent: The remaining `7` "ambiguous" pairings in the first deduped backfill
comparison are repeated same-model xhigh historical reruns for the
`SIM-robot-app-semantics-conformance` slice, not true mixed-lineage ambiguity.
The report should keep real ambiguity visible without overstating rerun noise.
Constraints: Do not rewrite or delete any historical result artifacts. Preserve
deterministic latest-rerun selection. Keep true ambiguity visible when more than
one narrowed lineage remains after rerun collapse.
Affects: `tools/ga-runner/compare.go`; `tools/ga-runner/ga_runner_test.go`;
`tools/ga-runner/README.md`; `results/RUN-PROTOCOL.md`;
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.

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
- [x] tapur.4 Define child-generation prompts that produce normal simulation
  shaped trees under `proposals/<run-group-id>/simulations/SIM-<handle>-child-<slug>/`
  with `README.md`, `QUESTION.md` when needed, optional `SCENARIOS.md`, optional
  local protocol/spec dirs, provenance back to parent sims, and bounded design
  deltas. Source: `DI-ramar`; `DI-zohal`; `DI-lirat`.
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
- [x] tapur.8 Implement review and promotion: accepted children are reviewed from
  ignored `proposals/<run-group-id>/` paths before selected designs are promoted
  into canonical non-child `simulations/SIM-*` paths with selected JSON result
  evidence; rejected children remain uncommitted. Source: `DI-ramar`; `DI-podot`;
  `DI-lirat`; `DI-dipid`.
- [x] tapur.9 Implement culling: rejected child sim trees and matching
  `proposals/<run-group-id>/results/<child-sim-id>/` trees are deleted only
  through an explicit cull command that records the action in the GA state file.
  Source: `DI-ramar`; `DI-kofil`; `DI-lirat`.
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
  improve over parent scores while pairing the best scored parent with one
  deterministic uniform random scored non-top parent. Source: `DI-bukid`;
  `DI-tufud`.
- [x] tapur.21 Compact child-generation prompts so breed calls include parent
  documents once, scenario pressure once, and summarized fitness evidence instead
  of repeated root boilerplate and full result JSON. Source: `DI-dilaf`.
- [x] tapur.22 Add Responses API streaming liveness diagnostics, broader
  transient-provider retries, a twelve-minute default retry elapsed budget, and
  top-plus-random scored parent selection for the next GA canary. Source:
  `DI-tufud`.
- [x] tapur.23 Make the terminal canary request reasoning summaries and print
  supported reasoning-summary and visible-output stream deltas to stdout/log
  without claiming access to hidden raw reasoning tokens. Source: `DI-vadub`.
- [x] tapur.24 Raise terminal canary default sync parallelism to six score
  workers and two child-generation workers while preserving serial raw command
  defaults and cost-reservation gates. Source: `DI-pivuj`.
- [x] tapur.25 Narrow the terminal canary default child-generation worker count
  back to one until a successful generation phase provides evidence for
  parallel child writes. Source: `DI-suzor`.
- [x] tapur.26 Replace canary reasoning-summary stream-content text with one
  no-newline progress dot per `response.reasoning_summary_text.delta` event.
  Source: `DI-babik`.
- [x] tapur.27 Restore canary `response.reasoning_summary_part.done` event names
  and content while keeping text-delta events as no-newline progress dots.
  Source: `DI-vajut`.
- [x] tapur.28 Print canary `response.reasoning_summary_part.added` event names
  and content as well, while keeping text-delta events as no-newline progress
  dots. Source: `DI-sakam`.
- [x] tapur.29 Split canary score and child-generation request timeout defaults
  so scoring remains `xhigh` / `5m` and child generation defaults to `medium` /
  `15m`. Source: `DI-guvif`.
- [x] tapur.30 Stop printing canary `response.reasoning_summary_part.added`
  event names and content; keep `response.reasoning_summary_part.done` output.
  Source: `DI-fupob`.
- [x] tapur.31 Stop printing canary `response.output_text.delta` event names and
  content while keeping output deltas for internal JSON response assembly.
  Source: `DI-ramun`.
- [x] tapur.32 Define the rubric/result v2 contract: `promisegrid.ga.result.v2`,
  `ga-rubric-20260522-v2`, and the new `promise_vocabulary` +
  `simplicity_durability` axes while keeping historical v1 evidence valid and
  append-only. Source: `DI-roruj`.
- [x] tapur.33 Add a deterministic GA audit mode that classifies canonical v1
  results by sim/scenario exact-match status, root-contract drift, and
  vocabulary hard-hit / soft-hit / clean status before any targeted backfill.
  Source: `DI-roruj`.
- [x] tapur.34 Add the audit-first targeted backfill selection policy: include
  exact-match hard-hit sims first, then add a clean grid-envelope calibration
  slice instead of paying for a full-corpus rerun up front. Source: `DI-roruj`.
- [x] tapur.35 Add targeted backfill state initialization that writes a fresh
  `results/state/<run-group-id>.json` with queued v2 result paths and preserves
  original model IDs by default. Source: `DI-roruj`.
- [x] tapur.36 Add a v1-vs-v2 comparison report once targeted rubric-v2 results
  exist, with explicit rank-delta reporting for envelope contenders and the
  claim/conformance family. Implemented as `ga-runner compare-backfill` under
  `DI-zuzup`; the first report is
  `results/reports/ga-backfill-20260522-215638-medium-comparison.md`. Source:
  `DI-roruj`; `DI-zuzup`.
- [x] tapur.37 Update surrounding GA-runner docs so operators can see the v1/v2
  coexistence contract, the new `audit` / `backfill-init` commands, and the
  targeted-backfill-first policy. Source: `DI-roruj`.
- [x] tapur.38 Add `/tmp/canary-cells` focus-file parsing to
  `tools/ga-runner/run-canary.sh`, merging resolved sim/scenario prefixes with
  `GA_CANARY_INCLUDE_*` so focused canaries no longer require editing
  `/tmp/canary.env` and no live scoring queue is introduced. Source:
  `DI-bataj`.
- [x] tapur.39 Make `audit` and `backfill-init` promotion-aware so canonical
  promoted results can use byte-identical canonical sim trees as audit fallback
  when historical proposal roots are gone, while unresolved sources stay
  non-exact instead of aborting the run. Source: `DI-zobur`.
- [x] tapur.40 Add staged backfill model/effort overrides so multi-stage rescoring
  can emit truthful state/result-path model IDs instead of reusing the
  historical model lineage strings when a new stage intentionally changes
  effort. Source: `DI-hijub`.
- [x] tapur.41 Change the default GA parent scoring effort from `xhigh` to
  `medium`, using calibration evidence from
  `ga-calib-20260523-033216-{low,medium,high,xhigh}` and keeping `xhigh` as the
  explicit escalation path for tie-breaks and promotions. Source: `DI-nanor`.
- [x] tapur.42 Refine `compare-backfill` ambiguity accounting so same-model
  historical reruns collapse to one deterministic winner and are reported
  separately from true ambiguous historical pairings. Source: `DI-sirir`.

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
