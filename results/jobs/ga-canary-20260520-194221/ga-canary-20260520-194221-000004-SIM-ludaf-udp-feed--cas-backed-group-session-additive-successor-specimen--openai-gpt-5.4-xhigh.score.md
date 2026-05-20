# GA Score Cell

Return only JSON with keys `scores`, `fitness`, and `assessment`.
Do not include result identity, source metadata, code fences, or commentary.

## Cell

- Run group ID: `ga-canary-20260520-194221`
- Cell ID: `ga-canary-20260520-194221-000004-SIM-ludaf-udp-feed--cas-backed-group-session-additive-successor-specimen--openai-gpt-5.4-xhigh`
- Simulation ID: `SIM-ludaf-udp-feed`
- Scenario ID: `cas-backed-group-session-additive-successor-specimen`
- Model ID: `openai-gpt-5.4-xhigh`
- Result path: `results/SIM-ludaf-udp-feed/cas-backed-group-session-additive-successor-specimen/openai-gpt-5.4-xhigh/20260520-194221.json`

## Rubric

Score each axis from 0 to 5. Higher is better except `risk_penalty`, where 0 is low risk and 5 is severe risk.
Axes: scenario_fit, promisegrid_alignment, auditability, evolution_safety, layer_boundary_clarity, failure_handling, implementation_plausibility, risk_penalty.

## Source Documents

### `results/RUN-PROTOCOL.md`

```markdown
# Results Run Protocol

This document defines operational contracts for result evidence under `results/`.
Legacy matrix runs write Markdown files at
`results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md`. GA/search runs
write JSON fitness files at
`results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.json` and checkpoint
state at `results/state/<run-group-id>.json`. Source: `DI-zamin`; `DI-ramar`;
`DI-zanon`; `DI-ruzaj`.

## Purpose

The run protocol keeps full-matrix execution reproducible, auditable, and
comparable across models and timestamps.

## Required Inputs Per Cell

- `sim-id` and simulation docs:
  - `simulations/<sim-id>/README.md`
  - `simulations/<sim-id>/QUESTION.md` (if present)
  - simulation local spec docs (if present)
- `scenario-id` and scenario docs:
  - `scenarios/README.md` for the shared scenario contract
  - `scenarios/<scenario-id>/<scenario-id>.md`
  - `model-id` must be explicit and path-safe
    (for example `openai-gpt-5.4-xhigh`).

## LLM Reasoning Requirement

Full matrix cells must be evaluated by an LLM or human reasoner reading the
cell's source docs. A script may prepare manifests, prompts, queues, and
validation reports, but a script must not generate final verdict prose by
mechanical parsing or keyword heuristics. Source: `DI-moduf`.

Full-matrix execution may be unattended: the queue runner may invoke one
external LLM command per manifest row, checkpoint state after every cell,
and validate the result file. The external runner must still produce the
substantive result file; queue tooling only coordinates work and validation.
Source: `DI-nuhon`; `DI-zamin`.

The preferred runner for legacy Markdown matrix runs is the Go
`tools/matrix-runner` CLI. For API-backed runs it bundles local source document
contents into the provider prompt, because remote APIs cannot read repo-local
paths. Source: `DI-lulom`; `DI-ruzaj`.

The preferred runner for GA/search work is `tools/ga-runner`. It uses JSON
fitness evidence, `promisegrid.ga.state.v1` state, generated child sims under
`simulations/SIM-*`, explicit review via `accept`, and explicit cleanup via
`cull`. Source: `DI-ramar`; `DI-zanon`; `DI-podot`; `DI-kofil`; `DI-ruzaj`.

Root scenario prompts use `scenarios/README.md` for shared scenario contract
context and `scenarios/<scenario-id>/<scenario-id>.md` for scenario-specific
pressure. Per-scenario `README.md` files are intentionally absent so repeated
boilerplate can be cached once rather than bundled with every scenario. Source:
`DI-kizal`.

API-backed runs must use explicit cost controls before any large batch. The
runner defaults to concise result style and a lower output cap, records provider
usage in queue state, prints actual accumulated cost, and can stop before a
cell whose preflight estimate would exceed the configured budget. Source:
`DI-nugiv`.

GA/search provider calls must also send an explicit service tier. The default is
Flex for lower-cost unattended work, `default` must be requested explicitly when
standard processing is desired, and Priority is rejected. Flex `429` and timeout
failures use bounded exponential backoff rather than silently falling back to a
higher-cost tier. Source: `DI-mopob`.

Synchronous GA/search calls should be bounded before Batch mode is available:
raw `tools/ga-runner score` and `generate` default to one worker, five-minute
provider attempts, two attempts per cell or child, and a six-minute retry
elapsed cap. Canary wrappers may raise scoring workers when cost reservations
are active, but they must not dispatch concurrent cells past the configured run
budget. Source: `DI-juzus`.

Unattended canary-style GA runs may continue past unusable individual cells by
using `score -skip-failed-cells` and `generate -skip-failed-children`. Skipped
items keep their failure messages in state, do not create synthetic result JSON,
and do not block child generation or child scoring of generated children. Source:
`DI-zikag`.

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
- API-backed blind runs should use concise evidence prose by default: keep every
  required section, but avoid source restatement and prefer short fit,
  weakness, and open-question bullets. Source: `DI-nugiv`.

## Result File Contract

### Legacy Markdown Matrix Results

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

### GA/Search JSON Fitness Results

GA/search results use schema `promisegrid.ga.result.v1` and path
`results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.json`. The JSON
file is the fitness evidence; there is no separate `results/fitness/` tree.
Required content includes source paths, source hashes, runner/model metadata,
rubric axes, integer rubric scores, normalized fitness, rationale, risks, open
questions, and authority boundary. Source: `DI-ramar`; `DI-zanon`; `DI-pobus`;
`DI-ruzaj`.

Old Markdown canary results remain historical evidence, but `tools/ga-runner`
must ignore `results/**/*.md` when validating or selecting GA fitness evidence.
Source: `DI-ramar`; `DI-pobus`; `DI-ruzaj`.

## Batch Rules

- Generate a deterministic manifest before running a batch.
- Use one `run_group_id` for the manifest and one result timestamp for a
  single run invocation.
- Use manifests with concrete `timestamp`, `result_path`, `ordinal`, and
  `cell_id` fields for unattended runs.
- Use checkpoint state under `results/state/` for any long run that should
  resume without operator prompts.
- Set an explicit `-max-run-cost-usd` for unattended API-backed batches, and use
  estimate-only output-token budgeting with per-cell/per-child estimate caps
  where the runner supports it. Treat hard `-max-output-tokens` caps as an
  explicit emergency fuse, not as the normal budget control. Source: `DI-pulap`.
- Treat `results/` as the only canonical result evidence. Generate result views
  from `results/` instead of committing scenario-side summaries.
- Keep old result files; never rewrite or delete prior runs.

## Failure and Retry Policy

- A cell is `failed` if required input docs are missing or result validation
  fails.
- Retries are explicit reruns with a new timestamp unless the same queue is
  intentionally resumed for a cell that never produced a valid result.
- Do not overwrite result files from a previous attempt.

## Tooling

- Preferred legacy Markdown matrix runner:
  `cd tools/matrix-runner && go run . <subcommand>`
- Preferred GA/search runner:
  `cd tools/ga-runner && go run . <subcommand>`
- Manifest generator:
  `python3 results/tools/generate_matrix_manifest.py --models <model-id>`
- LLM job generator:
  `python3 results/tools/generate_llm_jobs.py --manifest <manifest.csv>`
- Unattended queue runner:
  `python3 results/tools/matrix_queue.py run --manifest <manifest.csv> --runner-command '<command with {prompt_path}>'`
- Generated result view:
  `cd tools/matrix-runner && go run . view -repo-root ../.. -scenario <scenario-id>`
- Result validator:
  `python3 results/tools/validate_results.py --model <model-id> --timestamp <ts>`

## GA/Search Runner Shape

Use `tools/ga-runner` for JSON-fitness GA/search work. Implemented commands are:

- Validate JSON fitness evidence:
  `cd tools/ga-runner && go run . validate -repo-root ../..`
- Preview tracked population and conservative generation sizing:
  `cd tools/ga-runner && go run . init -repo-root ../.. -dry-run -model <model-id> -run-group-id <run-group-id>`
- Create state for one GA/search generation:
  `cd tools/ga-runner && go run . init -repo-root ../.. -model <model-id> -run-group-id <run-group-id>`
- Score parent or child cells with provider-backed reasoning:
  `cd tools/ga-runner && go run . score -repo-root ../.. -run-group-id <run-group-id> -target parents -api-model <provider-model> -reasoning-effort <effort> -service-tier flex -skip-failed-cells -max-run-cost-usd <budget>`
- Generate untracked child simulation trees from the active state:
  `cd tools/ga-runner && go run . generate -repo-root ../.. -run-group-id <run-group-id> -api-model <provider-model> -reasoning-effort <effort> -service-tier flex -skip-failed-children -max-run-cost-usd <budget>`
- Record reviewed promotion candidates without staging or committing:
  `cd tools/ga-runner && go run . accept -repo-root ../.. -run-group-id <run-group-id> -child <SIM-id> -result <json-result-path> -reviewer-note '<note>'`
- Cull rejected generated children and matching result trees through state:
  `cd tools/ga-runner && go run . cull -repo-root ../.. -run-group-id <run-group-id> -child <SIM-id> -reason '<reason>'`

`ga-runner progress` remains planned until its code path is implemented.
Generated children are normal `simulations/SIM-*` trees and are not accepted
merely because they exist on disk. Source: `DI-ramar`; `DI-zanon`; `DI-zohal`;
`DI-zusit`; `DI-podot`; `DI-kofil`; `DI-ruzaj`; `DI-gijom`.

## Recommended Preflight

1. Generate a canary manifest (20-30 cells).
2. Generate canary LLM job prompts or run the API queue with explicit budget
   flags.
3. Validate canary output.
4. Review token/cost fields in the state file and tune text verbosity,
   reasoning effort, or estimate-only output-token budgets.
5. Review drift/comparison report.
6. Run the full manifest only with an explicit budget. Source: `DI-nugiv`.

## Unattended Full-Run Shape

1. Generate a manifest with a fixed model ID and concrete timestamp:
   `cd tools/matrix-runner && go run . manifest -repo-root ../.. -models <model-id>`.
2. Start the queue with OpenAI API-backed execution:
   `cd tools/matrix-runner && go run . run -repo-root ../.. -manifest <manifest.csv> -provider openai -api-model <api-model> -reasoning-effort xhigh -result-style concise -max-output-tokens 6000 -max-run-cost-usd <budget> -max-cell-estimate-usd <cell-cap>`.
3. Let the queue process one cell at a time. Each cell writes or refreshes its
   prompt under `results/jobs/<run-group-id>/`, invokes the runner command,
   validates `result_path`, records token/cost usage, and checkpoints
   `results/state/<run-group-id>.json`. Source: `DI-nuhon`; `DI-lulom`;
   `DI-zamin`; `DI-nugiv`.
4. If the process is interrupted, rerun the same command with the same manifest
   and state path. Completed cells are skipped by default.
5. When the queue completes, validate the manifest:
   `cd tools/matrix-runner && go run . validate -repo-root ../.. -manifest <manifest.csv>`.
6. Generate inspection views from the result tree as needed:
   `cd tools/matrix-runner && go run . view -repo-root ../.. -model <model-id>`.
```

### `scenarios/README.md`

```markdown
# Root Scenarios

Root `scenarios/` entries are wire-lab comparison apparatus. They describe
pressure that multiple simulations can be run against; they are not PromiseGrid
node layout, production API, shared protocol components, or simulation-local
world state. Source: `DI-faros`; `DI-vabor`; `DI-dimas`; `DI-kizal`.

## Directory Shape

Each root scenario entry uses this shape:

```text
scenarios/
  <scenario-id>/
    <scenario-id>.md
```

- `<scenario-id>` is stable kebab-case.
- Application entries use one `scenarios/<application-id>/` directory per
  application, such as `scenarios/bgp-routing/`.
- Mined simulation rows use one root scenario entry per source row, transformed
  for cross-simulation comparison and linked back to the source
  `simulations/.../SCENARIOS.md` row.
- Per-scenario `README.md` files are intentionally absent. Shared scenario
  contract prose lives here so API prompts can reuse the same cacheable context,
  while `scenarios/<scenario-id>/<scenario-id>.md` files carry only
  scenario-specific pressure. Source: `DI-kizal`.
- Result evidence lives only under
  `results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md`. Generate
  scan tables from that result tree when needed. Source: `DI-zamin`.

## Scenario Entry Template

Use this template for each scenario markdown file:

```markdown
# <Scenario Title>

## Scenario ID

<stable-kebab-case-id>

## Source / Provenance

- Source type: application seed | mined simulation row | new harness scenario
- Source path:
- Source row/title:
- Source DI / TODO / TE:

## Purpose

<What design pressure this scenario applies.>

## Setup

<Initial state, promises, artifacts, sites, policies, and sparse knowledge.>

## Stimulus

<The event or request that exercises the scenario.>

## Expected Pressure

<What the scenario should force a candidate design to explain.>

## Scenario-Specific Evaluation Questions

- <Only questions unique to this scenario, if any.>
```

Do not duplicate the common actors, applicability, goal checks, result-path
shape, or authority boundary in each scenario file. Those shared rules are part
of this root scenario contract. Source: `DI-kizal`.

## Shared Applicability

Every scenario can be applied to any simulation that claims to address its
application domain, source-row pressure, or underlying PromiseGrid design
question. A result may conclude that the simulation intentionally does not cover
the scenario, but the run should still make the boundary explicit. Source:
`DI-kizal`.

## Actor Convention

Use Alice, Bob, Carol, Dave, Ellen, Frank, and Mallory where named actors help
make promises, failures, sparse knowledge, or adversarial pressure concrete.
Alice normally depends on the outcome, Bob makes or relays promises, Carol
audits or relies on evidence, and Mallory represents adversarial, captured,
stale, or misleading evidence when that pressure is relevant. Source:
`DI-034-20260508-060134`; `DI-kizal`.

## Scenario Quality Gates

Every root scenario run must explicitly test the 100-year PromiseGrid goal and
other overarching PromiseGrid goals. A scenario may be small, but it must not be
narrow in a way that accidentally assumes away the conditions PromiseGrid is
meant to survive. Source: `DI-botup`; `DI-kizal`.

At minimum, each run should address:

- **100-year durability:** Does the scenario still make sense after tools,
  organizations, keys, people, and infrastructure have changed?
- **Sparse, partial knowledge:** Does it avoid assuming any peer has the whole
  graph, all CAS objects, or globally complete state?
- **No central authority or registry:** Does it avoid relying on a central pCID,
  identity, trust, naming, routing, currency, or governance authority unless the
  point of the scenario is to test that failure mode?
- **Peer-local promise accounting:** Does it show what Alice, Bob, Carol, or
  another peer can observe and record locally?
- **Adversarial or failure pressure:** Does it include Mallory, corruption,
  refusal, stale data, partition, capture, default, or another failure mode when
  relevant?
- **Human and LLM auditability:** Can a later person or model understand what
  was promised, what happened, and why the result matters?
- **Migration / evolution path:** Does it expose what happens when protocols,
  keys, names, object shapes, policies, or organizations evolve?

If a gate is not relevant, the result should say why instead of omitting it.

## Common Evaluation Questions

Every scenario run should answer these common questions unless a
scenario-specific question narrows them:

- Which promises, CAS objects, feeds, identity claims, names, promise accounting
  records, or protocol claims does the candidate simulation need?
- What can Alice, Bob, Carol, Mallory, or another peer observe and record locally
  after the stimulus?
- Does the candidate design preserve the expected pressure without appealing to
  hidden global state or central authority?
- Which assumptions would make the run fail the 100-year, sparse-knowledge, or
  no-central-authority goals?
- What DR, DI, frozen spec, TODO, TE, or guide handoff would the evidence inform?

## Result Runs

Result runs live under:

`results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md`

Committed scenario entries are input context, not result summaries. Do not add a
`MATRIX.md` file to a scenario directory. To inspect run evidence, generate a
view from the canonical result tree:

```bash
cd tools/matrix-runner
go run . view -repo-root ../.. -scenario <scenario-id>
```

The generated view is derived evidence navigation. It does not declare a winning
design by itself, and it should not be committed as scenario source state.
Source: `DI-zamin`.

## Authority Boundary

Scenarios and result runs are evidence only. Design authority still graduates
through DR, DI, frozen spec, or PromiseGrid Development Guide handoff. Source:
`DI-faros`; `DI-kizal`.

## Population Plan

`TODO-dadub` owns the initial population plan. The first mining pass created one
root entry per existing sim-local scenario row under `DI-nanih`, the first
application-seed pass created one root entry per seed application under
`DI-midif`, and the cache cleanup pass removed per-scenario boilerplate under
`DI-kizal`. Future entries should follow this shared-contract shape unless a
later DI changes the root scenario contract.
```

### `simulations/SIM-ludaf-udp-feed/README.md`

```markdown
# SIM-ludaf: UDP feed

This simulation carries the renamed UDP feed lineage as an independent feed
family. It preserves the earlier `udp-binding` work as candidate lineage state
while shifting the active framing to `udp-feed`. Source: `DI-limom`,
`DI-rugig`.
```

### `simulations/SIM-ludaf-udp-feed/QUESTION.md`

```markdown
# Question

What parts of the UDP feed lineage remain useful when the lineage is treated as
a feed family rather than as a binding family? Source: `DI-limom`,
`DI-rugig`.
```

### `simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/CHANGELOG.md`

```markdown
# CHANGELOG: udp-feed

A-side CHANGELOG (per TE-liviv) for the `udp-feed` protocol's spec doc.

This file records **freeze events** authored by the spec maintainers.
Each entry names the doc-CID published at that moment. The format:

```changelog-entry
event:        freeze | withdraw | note
doc-cid:      bafkrei...
title:        Human-readable title for this freeze
workshop-tree: sha256:...   (optional; git tree-hash of this .d/ at freeze)
notes:        prose
```

No entries yet; this protocol has not yet reached a first freeze.

This protocol tree is now a simulation-local specimen per `DI-fakin`.

See [TE-zukug](../../../../docs/thought-experiments/TE-zukug-spec-doc-inversion-and-conformance-changelog.md), [TE-liviv](../../../../docs/thought-experiments/TE-liviv-spec-vs-implementation-split.md), and [TE-potar](../../../../docs/thought-experiments/TE-potar-spec-doc-informative-references.md) for the rationale.
```

### `simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/TODO/TODO-jodon-udp-binding-v0-reference-implementation.md`

```markdown
# TODO-jodon — UDP-feed v0 reference implementation

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TODO-18` (integer alias)
- `TODO-20260501-224805` (timestamp alias and pre-migration filename)

Source: TE-29 (`docs/thought-experiments/TE-20260501-215027-protocols-as-simulated-repos-and-binding-layer.md`).
Spec: `protocols/udp-feed.d/specs/udp-binding-draft.md` (one-page
draft, 10 normative promises, anti-promises, and test-vector
placeholders).

Simulation note: this TODO now lives with the `udp-feed` protocol specimen at
`simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/` (moved in `rusis.6`,
`DI-loluk`). Historical `udp-binding` / `SIM-piloh` paths remain provenance
only unless a line explicitly names a historical reference.

Flesh out UDP-feed v0 from one-page sketch into a usable v0
artifact: Go reference implementation, conformance test vectors, and
a minimal ns-3 harness scaffold that proves end-to-end round trip.
The protocol/specimen question about what evidence should count as a usable v0
conformance surface is now captured in
`simulations/SIM-kuful-udp-feed-v0-conformance/`, without implementing or
closing this TODO. Source: `DI-pukap`.

This is the first concrete binding implementation under the TE-29
layer decomposition, so it doubles as the proving ground for the
binding-layer abstraction itself: anything painful here is an
abstraction defect, not just a UDP issue.

## Subtasks

1. **Reference implementation** at
   `implementations/go-udp-feed-reference/` in Go (per Steve's
   standing language preference, and per TE-liviv which locks B-side
   implementation artifacts under `implementations/<impl-name>/`
   rather than under any `protocols/<slug>.d/` subtree). The
   implementation tree carries its own `CHANGELOG.md` recording
   conformance claims (e.g. `claim: implements` with the upstream
   `udp-binding-draft.md` doc-CID once frozen).

   - `Send(msg []byte, addr Addr) error` honoring promises 1, 2, 3,
     and 7 of the spec (one datagram per message, 1232 max, no
     in-band reliability, DSCP 0).
   - Recv loop with caller-supplied `Handle(msg []byte, src Addr)`
     callback honoring promise 6 (peer-set filter optional).
   - Local-error path for promise 2 size violations.
   - Promise 8: do not disable UDP checksum (`SO_NO_CHECK = 0`).
   - Stateless per promise 9.

2. **Simulation-artifact writer** honoring promise 10. When invoked
   under the wire-lab simulator, write each datagram payload to:

   ```
   transports/udp/<this-binding-pCID>/<session-pCID>/<message-pCID>/<message-id>.txt
   ```

   The session-pCID, message-pCID, and message-id come from the
   session layer above; this binding does not parse them. Behind a
   simulator-mode flag so production sends do not write artifacts.

3. **Test vectors.** Author the placeholder TVs from the spec
   (TV-1 through TV-5):
   - TV-1: 612-byte session message round-trips byte-for-byte over
     loopback UDP/4646.
   - TV-2: 1232-byte message round-trips byte-for-byte.
   - TV-3: 1233-byte message produces a local sender error before
     any datagram leaves the host.
   - TV-4: malformed datagram is handed up to the session layer
     unmodified by the binding.
   - TV-5: simulation-artifact file written for TV-1 contains exactly
     the 612 bytes of TV-1's session message.

   Test vectors live at `tools/udp-feed/testvectors/`.

4. **`/tmp/spec` walks `protocols/`.** Currently the spec checker
   only walks `specs/`. Update it (or write its successor) to walk
   `protocols/<slug>.d/specs/` so this draft is recognized.

5. **TODO-bihon — ns-3 harness scaffold for UDP-feed v0.** A
   minimal 2-node ns-3 scenario that proves round-trip works through
   the Go reference implementation talking over an ns-3-emulated
   UDP wire. See `protocols/wire-lab.d/TODO/TODO-bihon-ns3-harness-scaffold.md` for full
   subtasks. Tracked as a sibling TODO rather than a subtask of this
   one because the scaffold has its own follow-on lifecycle (loss
   scenarios, multi-binding scenarios) that long outlives the v0
   binding implementation. TODO-jodon is done when the scaffold from
   TODO-bihon successfully runs the v0 reference implementation
   end-to-end at least once.

6. **Update `protocols/udp-feed.d/specs/udp-binding-draft.md`**
   to reference (a) the implementation path under `tools/`, (b) the
   test-vector files, (c) the ns-3 harness scenario name. Replace
   "to be added in TODO-jodon" placeholders with concrete pointers.

## Out of scope for TODO-jodon

- Multiple concurrent UDP bindings on different ports (use case:
  one binding per protocol stack on a host). Possible v1 extension.
- IPv6 specifics beyond the 1232-byte size derivation. Should Just
  Work but is not exercised by v0 test vectors.
- NAT traversal. Explicit anti-promise in the spec.
- Path MTU discovery beyond "error out below 1232." Implementations
  may add it; not required by v0.
- Multicast. Permitted by the spec but not exercised by v0 test
  vectors; receivers can join groups manually for ad-hoc testing.

## Out of scope, but flagged as next likely TODO

- A second feed/binding (TCP feed v0 or WebSocket feed v0) so the
  binding-layer abstraction is exercised across at least two
  qualitatively different real-world transports. This is when the
  per-binding-pCID forking property at C-4 actually starts paying
  off. (Originally hedged as "Likely TODO 020" before the TODO 020 slot was claimed by an unrelated TE-editing-policy item; this is a future TODO not yet filed.)
- A Go reference implementation of group-session v0 to ride above
  UDP-feed v0 and prove the layer composition. (Originally hedged as "Likely TODO 021" before that slot was claimed by session-replay-cleanup; this is a future TODO not yet filed.)

## Dependencies

- TODO-vuhuj (protocols-as-simulated-repos migration) ideally lands
  first so the spec lives at its final path. If TODO-vuhuj has not
  landed, this TODO uses `protocols/udp-feed.d/specs/udp-binding-draft.md`
  (already created in TE-vipir's commit).
- TODO-losoh (group-session rename) is not required; UDP-feed does
  not depend on the session protocol's slug.

## Done when

- `tools/udp-feed/` Go package builds and passes all five test
  vectors.
- `protocols/udp-feed.d/specs/udp-binding-draft.md` is updated
  with concrete pointers replacing "to be added" placeholders.
- TODO-bihon's ns-3 harness scenario runs the Go reference
  implementation end-to-end at least once and produces matching
  PCAP and `.txt` artifacts.
- `/tmp/spec check` (or its successor) reports OK with the spec
  recognized at its new path under `protocols/udp-feed.d/`.
```

### `simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/TODO/TODO.md`

```markdown
# TODO queue: udp-feed

Per-protocol TODO queue (per TE-magup). This queue moved with the
`udp-feed` protocol specimen into
`simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/` in `rusis.6`
(`DI-loluk`). Earlier `udp-binding` paths remain historical provenance
only. Items in this file touch only this simulation-local
protocol specimen. Anything broader is harness-level and lives at
`../../../../../protocols/wire-lab.d/TODO/TODO.md`.

Per TE-mumuv (TE-39, locked 2026-05-07), each TODO is addressable
by its proquint handle (TODO-<handle>). Prior integer / timestamp
aliases survive in the `Prior alias` column and in each file's
`## Prior aliases` section.

## Index

| Handle | Mint date | Title | Prior alias |
|---|---|---|---|
| [TODO-jodon](TODO-jodon-udp-binding-v0-reference-implementation.md) | 2026-05-01 | TODO 018 — UDP-feed v0 reference implementation (historical title used `UDP-binding`) | `TODO-18` / `TODO-20260501-224805` |
```

### `simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/specs/udp-binding-draft.md`

```markdown
# UDP-feed v0 (draft)

## Status

DRAFT. One-page sketch authored alongside TE-vipir
(`docs/thought-experiments/TE-vipir-protocols-as-simulated-repos-and-binding-layer.md`).
This is the first concrete L4 feed/binding spec under the layer
decomposition locked in TE-vipir. Subject to revision before freeze.

Simulation note: this draft now lives at
`simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/` (moved in
`rusis.6`, `DI-loluk`). Earlier `protocols/udp-binding.d/` and
`SIM-piloh` paths remain historical provenance only. The path examples
below remain draft evidence under test, not final PromiseGrid API
layout.

## Abstract

This document specifies how PromiseGrid messages are carried over
UDP (RFC 768) in version 0 of the binding. It is intentionally
minimal: one PromiseGrid session message per UDP datagram, no
fragmentation, default port 4646, best-effort delivery promises only.

This binding does not redefine UDP. It defines exactly how
PromiseGrid uses UDP and what a conformant implementation promises.

## Layer position

UDP-feed v0 occupies level 2 in the five-level stack defined by
TE-vipir:

```
transports/udp/udp-feed-bafkrei...U1/<session-pCID>/<message-pCID>/<message-id>.txt
                ^^^^^^^^^^^^^^^^^^^^^^^^
                this spec
```

Level 1 (`udp`) names the real-world transport (RFC 768), which this
binding does not modify. Levels 3 and above (session, message) are
opaque to this binding.

## Promises (normative)

I promise that:

1. **One message per datagram.** Each PromiseGrid session message
   handed to this binding for sending is transmitted in exactly one
   UDP datagram. The datagram payload equals the session message
   bytes verbatim, with no prefix, suffix, framing, magic number, or
   length field added by this binding.

2. **Maximum size 1232 bytes.** Session messages exceeding 1232
   bytes (the IPv6 minimum MTU minus IPv6 and UDP headers) MUST be
   rejected at the sender with a local error. This binding does not
   fragment. Senders MUST NOT rely on IP fragmentation; if Path MTU
   Discovery reports a path MTU below 1232, the sender MUST surface
   a local error rather than send a fragmented datagram.

3. **No delivery, ordering, or deduplication promises.** This binding
   inherits UDP's best-effort semantics. Datagrams may be lost,
   duplicated, reordered, or delivered to unintended ports. Higher
   layers (session protocol) are responsible for any guarantees
   beyond best-effort.

4. **Default port 4646.** Conformant implementations default to UDP
   port 4646 for both send and receive. Any other port is allowed by
   mutual configuration. This document does not allocate the port; no
   central allocator exists. 4646 is a convention, not a reservation.

5. **Address shape `host:port`.** Peer addresses presented to the
   binding take the form `host:port` where `host` is an IPv4 or IPv6
   literal or a DNS name resolvable at send time, and `port` is a UDP
   port number. Multicast group addresses are permitted; behavior on
   multicast send is identical (one datagram per recipient join).

6. **Receive contract.** A conformant receiver binds a UDP socket to
   the configured port, calls `recvfrom` (or equivalent) with a
   buffer of at least 1232 bytes, and hands each datagram payload up
   to the session layer as one complete message. Datagrams whose
   source address is not in the configured peer set MAY be dropped
   silently or logged; either behavior is conformant.

7. **DSCP default zero.** Senders SHOULD set IP DSCP to 0 (default
   forwarding). Other DSCP values are permitted by configuration but
   are out of scope for this spec.

8. **No checksums beyond UDP's.** This binding does not add a
   checksum. UDP's own checksum (RFC 768) is sufficient. Senders MUST
   NOT disable the UDP checksum (zero-checksum optimization is
   prohibited at this binding).

9. **No connection state.** This binding is stateless. No handshake,
   no keep-alive, no reconnect. Senders may send to any address at
   any time; receivers may receive from any address at any time
   (subject to the peer-set filter at promise 6).

10. **Simulation artifact format.** When this binding runs inside the
    wire-lab simulator, it writes one file per datagram seen on the
    wire to:

    ```
    transports/udp/<this-binding-pCID>/<session-pCID>/<message-pCID>/<message-id>.txt
    ```

    File contents equal the exact datagram payload bytes (the same
    bytes handed up to the session layer, byte-for-byte). The session
    and message-protocol pCIDs in the path are derived by the session
    layer above; this binding does not parse them. The `<message-id>`
    in the filename is supplied by the session layer (typically a
    content hash of the session message bytes).

## Anti-promises (non-normative clarifications)

This binding does not promise:

- **Reliability.** Datagrams may be lost. Use a session protocol that
  handles retransmit if you need reliability.
- **Order.** Datagrams may arrive in any order.
- **Privacy or authenticity.** UDP is unencrypted. This binding does
  not add encryption or signatures. Use session/message-layer
  cryptography (per C-6) for authenticity.
- **Spam resistance.** A receiver bound to UDP/4646 will receive any
  datagram any sender sends to it. Filtering and rate-limiting are
  out of scope for v0.
- **NAT traversal.** Pure UDP without rendezvous works for hosts with
  reachable addresses. NAT traversal is a separate concern (probably
  a future binding or a session-layer feature).
- **Path MTU discovery.** Implementations MAY perform PMTUD; this
  binding only requires that they error out below 1232 rather than
  fragment.

## Test vectors (placeholder)

To be added in TODO-jodon. At minimum:

- TV-1: a 612-byte session message round-trips byte-for-byte through
  loopback UDP/4646.
- TV-2: a 1232-byte session message round-trips byte-for-byte.
- TV-3: a 1233-byte session message produces a local sender error
  before any datagram leaves the host.
- TV-4: a malformed datagram (e.g., truncated by transport-level
  loss simulation) is handed up to the session layer without
  modification by this binding (it is not this binding's job to
  validate session-layer content).
- TV-5: simulation-artifact file written for TV-1 contains exactly
  the 612 bytes of TV-1's session message.

## Reference implementation

To be authored in TODO-jodon under
`tools/udp-feed/` with at minimum:

- `Send(msg []byte, addr Addr) error`
- a recv loop that reads from a configured UDP socket and invokes a
  caller-supplied `Handle(msg []byte, src Addr)` callback.

Language: Go (per Steve's standing preference).

## Forking and versioning

Per C-4 (forking is normal, TE-dajot), any author may publish a UDP
binding with different choices (e.g., different default port,
different size limit, added framing). Such a fork takes a different
pCID and lives as a sibling under `protocols/`. Two parties wishing
to interoperate must agree on the same UDP-feed pCID; they cannot
silently mix bindings.

This feed may evolve to v1 if a load-bearing change is required
(e.g., explicit fragmentation support, NAT-traversal hooks).
Migration from v0 to v1 follows the rules in TE-dajot OQ-100.2.

## Bibliography

- RFC 768 (UDP)
- RFC 8200 (IPv6, for the 1232-byte size derivation)
- TE-vipir (this binding's layer position)
- TE-dajot (load-bearing constraints C-1 through C-6)

## Open questions

- OQ-UDP-1: Should the default port 4646 be parameterized in some
  registry-free way (e.g., derived from a hash of the binding pCID)?
  Lean: no, fixed convention is simpler.
- OQ-UDP-2: Should the 1232-byte limit be raised for known-IPv4-only
  paths? Lean: no, uniformity is more valuable than the few extra
  bytes.
- OQ-UDP-3: Multicast semantics in the simulation. How does
  `transports/udp/...` represent a single multicast send arriving at N
  receivers? Lean: one file per (sender, receiver) pair; deduplication
  by content hash makes this storage-cheap.
```

### `simulations/SIM-ludaf-udp-feed/seed/protocol-migration.md`

```markdown
# UDP lineage migration

The UDP lineage tree moved from
`SIM-piloh-turns-149-208-recovery/protocols/udp-binding.d/` into
`SIM-ludaf-udp-feed/protocols/udp-feed.d/` during `rusis.6`.

The active lineage name is now `udp-feed`. The older mixed umbrella records in
`protocol-set.md` and `seed/protocol-tree-migrations.md` referred to this same
lineage as `udp-binding`. That old name remains historical provenance only; it
is no longer the active directory name. Source: `DI-loluk`, `DI-rugig`.
```

### `scenarios/cas-backed-group-session-additive-successor-specimen/cas-backed-group-session-additive-successor-specimen.md`

```markdown
# Additive successor specimen

## Scenario ID

cas-backed-group-session-additive-successor-specimen

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-jurar-cas-backed-group-session/SCENARIOS.md`
- Source simulation: `SIM-jurar-cas-backed-group-session/`
- Source row/title: Additive successor specimen
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-jurar-cas-backed-group-session/`.

## Setup

Existing `.txt` messages remain historical while a new CAS-backed specimen is added beside them.

## Stimulus

Run the candidate simulation against this source test: Whether migration can be additive without rewriting historical message bytes or invalidating existing CIDs.

## Expected Pressure

TODO-pipus must design overlap / successor mechanics rather than mutate old evidence.
```

## Required JSON Shape

{"scores":{"scenario_fit":0,"promisegrid_alignment":0,"auditability":0,"evolution_safety":0,"layer_boundary_clarity":0,"failure_handling":0,"implementation_plausibility":0,"risk_penalty":0},"fitness":{"raw":0,"normalized_0_100":0,"confidence_0_1":0.0},"assessment":{"rationale":"","strengths":[],"weaknesses":[],"risks":[],"open_questions":[],"authority_boundary":"Evidence only; does not settle PromiseGrid design."}}
