# GA Child Generation

Return only JSON with keys `child_id`, `design_delta_summary`, and `files`.
Each file path must be relative to the child simulation root. Include `README.md` and `QUESTION.md`.

- Run group ID: `ga-canary-20260520-194221`
- Child ID: `SIM-gumuj-ga-child-0002`
- Child path: `simulations/SIM-gumuj-ga-child-0002/`
- Operation: `crossover`
- Parent IDs: `SIM-ludaf-udp-feed, SIM-labit-feed-outer`

## Scenario Sample

- `cas-backed-group-session-additive-successor-specimen` at `scenarios/cas-backed-group-session-additive-successor-specimen/cas-backed-group-session-additive-successor-specimen.md`
- `makerspace-door-access` at `scenarios/makerspace-door-access/makerspace-door-access.md`
- `community-movement-organizing` at `scenarios/community-movement-organizing/community-movement-organizing.md`

## Parent Source Documents

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

### `scenarios/makerspace-door-access/makerspace-door-access.md`

```markdown
# Makerspace Door Access

## Scenario ID

makerspace-door-access

## Source / Provenance

- Source type: application seed
- Source path: `protocols/wire-lab.d/TODO/TODO-dadub-root-scenario-skeleton-and-seed-catalog.md`
- Source row/title: Seed application catalog entry `makerspace-door-access`
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-midif`; `TODO-dadub`

## Purpose

Exercise PromiseGrid design candidates against makerspace-door-access application
pressure: Membership state, door authorization, guest access, emergency override, and
local accountability.

## Setup

Alice depends on an outcome in the Makerspace Door Access domain. Bob makes promises
about membership state, door authorization, guest access, emergency override, and local
accountability. Carol needs enough evidence to rely on Bob's promise without having
complete global state, and Mallory may exploit stale, missing, or ambiguous evidence.

## Stimulus

A routine application event becomes contested or incomplete. Some relevant objects,
signatures, observations, or relationship records are delayed, partitioned, stale, or
disputed, and each peer must decide what to accept, retry, downgrade, or escalate using
only local evidence.

## Expected Pressure

The candidate simulation must show which promises, CAS objects, feeds, identity claims,
names, and promise accounting records are needed for this application pressure, while
avoiding hidden global state or a central authority that would make the result non-
comparable.
```

### `scenarios/community-movement-organizing/community-movement-organizing.md`

```markdown
# Community Movement Organizing

## Scenario ID

community-movement-organizing

## Source / Provenance

- Source type: application seed
- Source path: `protocols/wire-lab.d/TODO/TODO-dadub-root-scenario-skeleton-and-seed-catalog.md`
- Source row/title: Seed application catalog entry `community-movement-organizing`
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-midif`; `TODO-dadub`

## Purpose

Exercise PromiseGrid design candidates against community-movement-organizing application
pressure: Membership, campaigns, working groups, public commitments, moderation, and
schism handling.

## Setup

Alice depends on an outcome in the Community Movement Organizing domain. Bob makes
promises about membership, campaigns, working groups, public commitments, moderation,
and schism handling. Carol needs enough evidence to rely on Bob's promise without having
complete global state, and Mallory may exploit stale, missing, or ambiguous evidence.

## Stimulus

A routine application event becomes contested or incomplete. Some relevant objects,
signatures, observations, or relationship records are delayed, partitioned, stale, or
disputed, and each peer must decide what to accept, retry, downgrade, or escalate using
only local evidence.

## Expected Pressure

The candidate simulation must show which promises, CAS objects, feeds, identity claims,
names, and promise accounting records are needed for this application pressure, while
avoiding hidden global state or a central authority that would make the result non-
comparable.
```

### `simulations/SIM-labit-feed-outer/README.md`

```markdown
# SIM-labit: Feed outer

This simulation reserves an independent lineage for the thin outer-feed
convention currently mixed into rooted transport-side material. It exists so
that later extraction and comparison can happen without assuming the thin outer
feed is already canonical. Source: `DI-limom`, `DI-rugig`.
```

### `simulations/SIM-labit-feed-outer/QUESTION.md`

```markdown
# Question

Is a thin outer-feed convention a useful separable lineage for PromiseGrid once
it is extracted from rooted transport-side material and compared against other
approaches? Source: `DI-limom`, `DI-rugig`.
```

### `simulations/SIM-labit-feed-outer/TODO/TODO-kakaz-feed-outer-freeze-and-cid-cascade-follow-on.md`

```markdown
# TODO-kakaz: Feed-outer freeze and CID-cascade follow-on

## Prior aliases

None. This TODO is created directly as a sim-local successor owner under
`rusis.10`.

## Status

Closed. The feed-outer specimen slice transferred out of rooted
`TODO-turog` and `TODO-duvuk` by `DI-mosor` is resolved by `DI-bomud`.

## Scope

This TODO owns the feed-outer side of the split follow-on work:

- outer transport-spec freeze boundary cleanup and freeze-gate follow-on
  that is not group-session-local;
- transport-level filename/CID-cascade policy and related provenance
  rules for the feed-outer lineage;
- extraction-handoff prep for `rusis.11` from rooted
  `transport-spec-draft.md`.

Out of scope: group-session-local envelope details (tracked in
`TODO-gapab`).

## Closure summary

`DI-bomud` closes this successor owner. Spec freeze publishes a spec pCID and
does not rename, rehash, or rewrite historical feed/transport specimens.
Draft-era directories remain evidence; any frozen successor or derived mirror
is additive and cites its source evidence. Feed-outer remains a thin outer
convention and does not own group-session message filename/CID,
`Message-ID:`, header, body, or reader/writer rules.

## Subtasks

- [x] kakaz.1 Define feed-outer freeze-boundary wording that supersedes
  stale TE-41 step-5 guidance without mutating historical specimen bytes.
  Closed by `DI-bomud`.
- [x] kakaz.2 Define the feed-outer slice of TE-42 Path-A-vs-Path-B and
  deprecation policy handling. Closed by `DI-bomud`: feed-outer owns no
  message-header compatibility rule.
- [x] kakaz.3 Prepare extraction handoff notes for `rusis.11` from
  rooted `protocols/wire-lab.d/specs/transport-spec-draft.md`. Closed by
  `DI-bomud`; active wording now lives in the simulation-local feed-outer
  draft.
- [x] kakaz.4 Back-link resulting decisions and artifacts to rooted
  `TODO-turog` and `TODO-duvuk` historical records. Closed by `DI-bomud`.

## Decision Intent Log

Successor-owner routing into this TODO was locked under `DI-mosor` in
`protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`.

ID: DI-bomud
Date: 2026-05-13 22:48:21
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Close the feed-outer successor slice for turn-164, TE-41, and TE-42
follow-on work. Spec freeze publishes new spec pCIDs but never rewrites
historical transport/feed data; frozen successors or derived views must be
additive and cite prior evidence. Feed-outer remains a thin outer convention
and does not own group-session message filename/CID, `Message-ID:`, header, or
body parsing rules. A Steve-authored DI is the operative
`merge-transport-spec` promise until cryptographic promise tooling exists.
Intent: Retire the stale rewrite-at-freeze plan while preserving a clean
feed-outer boundary from group-session-local message semantics.
Constraints: Do not mutate historical message bytes, old CIDs, or legacy
transport specimens. Keep group-session-local policy under `DI-rurab`. Do not
implement cryptographic signing or `tools/spec freeze` changes in this pass.
Affects: `simulations/SIM-labit-feed-outer/protocols/feed-outer.d/specs/feed-outer-draft.md`;
`simulations/SIM-ludut-wire-lab-devs/world/transports/wire-lab-devs-draft/README.md`;
`simulations/SIM-labit-feed-outer/TODO/TODO-kakaz-feed-outer-freeze-and-cid-cascade-follow-on.md`;
`protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md`;
`protocols/wire-lab.d/TODO/TODO-duvuk-te-42-filename-cid-cascade-policy.md`;
`protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`.
```

### `simulations/SIM-labit-feed-outer/TODO/TODO.md`

```markdown
# TODO queue: feed-outer

Per `rusis.10`, this queue tracks feed-outer lineage follow-on work that
was split out of rooted mixed-owner TODOs.

## Index

| Handle | Mint date | Title | Prior alias |
|---|---|---|---|
| [TODO-kakaz](TODO-kakaz-feed-outer-freeze-and-cid-cascade-follow-on.md) | 2026-05-11 | Feed-outer freeze + CID-cascade follow-on owner split from rooted TODOs **(closed: DI-bomud)** | none |
```

### `simulations/SIM-labit-feed-outer/protocols/feed-outer.d/specs/feed-outer-draft.md`

```markdown
# Feed-Outer Spec (DRAFT)

*This simulation-local draft is the extracted specimen-side outer feed
convention from `protocols/wire-lab.d/specs/transport-spec-draft.md`. It is a
draft and is subject to revision; once frozen, its pCID will name this
protocol class for all time.*

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Source extraction: `DI-huzor` (`rusis.11`).

## Purpose

This spec defines the **outer convention** for the wire-lab's `transports/`
directory: how transport instances are named on disk, the rule that messages do
not declare their transport via a header, and the requirement that each
transport-protocol's pCID names a separate spec defining the directory's
interior.

This spec is intentionally **thin.** It does not define:

- header sets,
- parent-link semantics,
- receipt formats,
- message-kind vocabulary,
- subdirectory layout inside a transport,
- canonical-bytes rules for messages.

All of those are properties of individual transport-protocols, defined in their
own spec docs (the first being
`simulations/SIM-rakot-group-session/protocols/group-session.d/specs/group-session-draft.md`).

## Freeze boundary and historical data

Spec freeze publishes a pCID for a spec. It does not rename, rehash, or rewrite
existing transport/feed specimens. Draft-era directories remain historical
evidence; any frozen successor or derived mirror is additive and must cite the
source evidence it derives from. Source: `DI-bomud`.

Until cryptographic promise tooling exists, a Steve-authored DI is the
operative `merge-transport-spec` promise for this feed-outer lineage. Source:
`DI-bomud`.

## The four locked principles (TE-zalut)

### Principle 1: a message does not declare its transport

A message envelope contains no `Transport:` header, no `Transport-Type:` header,
and no per-message reference to which transport it belongs to. The transport
carrying a message is identified by the transport itself: in the repo-local
case, by the directory the message file lives in.

Asking a message to declare its transport is layer inversion. The transport is
the carrier; the message is the cargo; the cargo does not name the carrier.

If a message needs to *reference* a different transport (e.g., a receipt
acknowledging a message that arrived on another transport), the referencing
protocol's spec defines how to do that. This outer spec is silent.

### Principle 2: transport directories are keyed `transports/<pcid>--<slug>/`

The directory name is structured:

```
transports/<pcid>--<slug>/
```

where:

- **`<pcid>`** is the canonical pCID of the transport-protocol that transport
  speaks. This is the load-bearing identifier: it tells any reader which
  protocol's contract governs the directory's interior.
- **`<slug>`** is a human-readable suffix that tools ignore (or use only for
  display). It exists so humans can navigate `transports/` without parsing
  pCIDs and so commit-log entries are legible.
- **`--`** (double hyphen) separates the two. The double hyphen is unlikely to
  appear inside a CIDv1 base32 string.

The pCID is canonical; the slug is a convenience. Two directories with the same
pCID and different slugs are **two different transport instances** of the same
protocol. Two directories with different pCIDs are different
transport-protocols and may have entirely different interior structure.
Draft or pre-freeze specimens may use explicit draft-state names such as
`wire-lab-devs-draft`; those names are not rewritten when a spec pCID is
minted. Source: `DI-bomud`.

### Principle 3: each transport-protocol-pCID names a spec defining its directory's interior

The pCID *is* the protocol's identity. The protocol gets to define everything
inside `transports/<its-pcid>--<slug>/`:

- subdirectory layout (flat, per-direction, per-participant, sharded by date,
  etc.),
- message file naming conventions,
- header set,
- parent-link semantics (whether messages cite parents at all, what header names
  them, how multiple parents serialize, optionality),
- receipt format,
- message-kind vocabulary,
- canonical-bytes rules,
- persistence rules (append-only, bounded retention, compactable, ephemeral),
- visibility rules (all-see-all, hub-mediated, ring-propagated, etc.),
- membership rules (closed, open, invite-only, capability-token, etc.).

The feed-outer spec does not constrain any of these. They live in the
transport-protocol's own spec doc.

### Principle 4: code-as-handler

The code that reads a transport directory's structure *is* the handler for that
pCID. Each transport-protocol comes with its own reader/writer code; the pCID
identifies the protocol; the protocol identifies (by convention or naming) the
code that speaks it. There is no machine-readable companion file (no
`transport.yaml` schema). The frozen markdown spec is the human-readable
contract that the code must implement.

Tools that want to display N transport-instances of M different
transport-protocols need M handlers. That is the cost of polymorphism, not a
flaw of this design.

## What this spec does NOT specify

- The first line of a message (`grid <pcid>` is one carrier choice; not all
  transport-protocols must use it).
- Header names (`Message-ID`, `Date`, `From`, `To`, `Parents`, `IHave`, etc. —
  all defined per-protocol).
- Canonical-bytes encoding (UTF-8/LF discipline is one choice; not all protocols
  must use it).
- File-naming inside a transport directory.
- Message-CID cascade rules, legacy `Message-ID:` compatibility, and any
  reader-side rehash/deprecation policy.
- Subdirectory structure inside a transport directory.

If a future reader asks "where do I find out how to write a message for this
transport?" the answer is always: read the spec named by that transport's pCID.
The feed-outer spec is silent on the message format. Source: `DI-bomud`.
```

### `simulations/SIM-labit-feed-outer/seed/extraction-sources.md`

```markdown
# Feed-outer extraction sources

The `feed-outer` lineage was seeded from the rooted draft at
`protocols/wire-lab.d/specs/transport-spec-draft.md` in `rusis.7`, and the
specimen-side sections were extracted in `rusis.11`. Source: `DI-ludaz`,
`DI-rugig`, `DI-huzor`.

Extracted into
`simulations/SIM-labit-feed-outer/protocols/feed-outer.d/specs/feed-outer-draft.md`:

- `## Purpose`
- `## The four locked principles (TE-zalut)`
- `## What this spec does NOT specify`

Kept rooted in
`protocols/wire-lab.d/specs/transport-spec-draft.md` as harness-side
apparatus/governance residue:

- `## Sources`
- `## The per-axis meta-rule (TE-junil)`
- `## Open questions`
- `## Freeze gate`

This note now records completed extraction provenance plus rooted-residue
boundaries. Source: `DI-huzor`.
```

## Existing Fitness Evidence From This Run

### `results/SIM-ludaf-udp-feed/cas-backed-group-session-additive-successor-specimen/openai-gpt-5.4-xhigh/20260520-194221.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-ludaf-udp-feed-cas-backed-group-session-additive-successor-specimen-openai-gpt-5.4-xhigh-20260520-194221",
  "run_group_id": "ga-canary-20260520-194221",
  "cell_id": "ga-canary-20260520-194221-000004-SIM-ludaf-udp-feed--cas-backed-group-session-additive-successor-specimen--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-ludaf-udp-feed",
  "scenario_id": "cas-backed-group-session-additive-successor-specimen",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-194221",
  "result_path": "results/SIM-ludaf-udp-feed/cas-backed-group-session-additive-successor-specimen/openai-gpt-5.4-xhigh/20260520-194221.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_074c1bb8095f4d5690dd619d8e30f9c6",
    "response_id": "resp_0021587c34d8931c006a0e0f0f40a8819b947179a4512e641c",
    "input_tokens": 9584,
    "output_tokens": 10246,
    "cost_usd": 0.160216
  },
  "source": {
    "repo_commit": "d79e4669c410f21c4c0abf616422de8e169da5c3",
    "sim_path": "simulations/SIM-ludaf-udp-feed/",
    "scenario_path": "scenarios/cas-backed-group-session-additive-successor-specimen/cas-backed-group-session-additive-successor-specimen.md",
    "root_contract_paths": [
      "results/RUN-PROTOCOL.md",
      "scenarios/README.md"
    ],
    "files": [
      {
        "path": "results/RUN-PROTOCOL.md",
        "sha256": "9bf7fd58adfffcb4aeebb87bca14576ac9016681ed3aac0750f2ddafb00792a6"
      },
      {
        "path": "scenarios/README.md",
        "sha256": "406c4c7f400df14788d1caea61406f83adbc474160806ce4f1aa6a88d409d483"
      },
      {
        "path": "simulations/SIM-ludaf-udp-feed/README.md",
        "sha256": "6238e2781f9658feae574df275356e4e25a6d22eba5b8a2f12ae22b1561ef252"
      },
      {
        "path": "simulations/SIM-ludaf-udp-feed/QUESTION.md",
        "sha256": "9136722a6e7b982ace707418665d088de2e9f229e099ea4a646ad6c5196e5463"
      },
      {
        "path": "simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/CHANGELOG.md",
        "sha256": "59df64a3749a92b706582ebe7041e19ba6545cf4840eaf2715c9cedef55e07aa"
      },
      {
        "path": "simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/TODO/TODO-jodon-udp-binding-v0-reference-implementation.md",
        "sha256": "e63864119b92646905ff3519a4173074c3ce11256fc9d7ae3dcbbf99b7de7dd4"
      },
      {
        "path": "simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/TODO/TODO.md",
        "sha256": "cc9141e03e0e508984e6ce60108ce2e20d3fd440172275251303c3fd4a0445d1"
      },
      {
        "path": "simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/specs/udp-binding-draft.md",
        "sha256": "5305152a5372c37f26f81c77b0d02fb9b3e82dd431d452ad1f11bdd46c4f0c7d"
      },
      {
        "path": "simulations/SIM-ludaf-udp-feed/seed/protocol-migration.md",
        "sha256": "d835ef1b97297872e7ef2700464a9be5b68fdd31aca4b05483dbc11d98003669"
      },
      {
        "path": "scenarios/cas-backed-group-session-additive-successor-specimen/cas-backed-group-session-additive-successor-specimen.md",
        "sha256": "e1eaad54488dd151065a9023e1d348e7cd0cbc1c3475ffa5bc19fe0f5b1e1887"
      }
    ],
    "simulation_tree_hash": "1f626da6f8c3305ac8a02bdbbbb9f66cdcfdba59ebd91dff2433acf4ec107c93"
  },
  "rubric": {
    "rubric_version": "ga-rubric-20260519-v1",
    "score_scale": "0..5",
    "score_meanings": {
      "0": "no fit or absent",
      "5": "strong fit",
      "risk_penalty": "0 low risk, 5 severe risk"
    },
    "axes": [
      "scenario_fit",
      "promisegrid_alignment",
      "auditability",
      "evolution_safety",
      "layer_boundary_clarity",
      "failure_handling",
      "implementation_plausibility",
      "risk_penalty"
    ]
  },
  "scores": {
    "scenario_fit": 1,
    "promisegrid_alignment": 2,
    "auditability": 4,
    "evolution_safety": 2,
    "layer_boundary_clarity": 4,
    "failure_handling": 2,
    "implementation_plausibility": 2,
    "risk_penalty": 2
  },
  "fitness": {
    "raw": 15,
    "normalized_0_100": 43,
    "confidence_0_1": 0.82
  },
  "assessment": {
    "rationale": "Useful only as indirect evidence. The candidate documents lineage-preserving rename/migration and clear transport/session boundaries, but it does not model the core scenario problem: adding a CAS-backed group-session successor beside historical .txt evidence without rewriting bytes or invalidating existing CIDs.",
    "strengths": [
      "Preserves earlier udp-binding lineage as historical provenance instead of destructive rewrite.",
      "Strong feed/session boundary: higher-layer session and message pCIDs stay opaque to the binding.",
      "Spec, TODO, and migration notes are easy to audit and frame evolution without a central registry."
    ],
    "weaknesses": [
      "No CAS-backed specimen, group-session behavior, or successor record format.",
      "No explicit overlap mechanism linking historical .txt artifacts to new CAS objects while keeping old evidence unchanged.",
      "No peer-local accounting trail showing what a later peer or auditor can record about successor status."
    ],
    "risks": [
      "Transport-lineage migration could be overread as evidence for session-layer additive-successor design.",
      "No adversarial model for stale or misleading successor claims.",
      "Draft status means artifact paths and promises may still change before freeze."
    ],
    "open_questions": [
      "What record links a historical .txt artifact to a new CAS-backed specimen while preserving existing CIDs?",
      "How would a 100-year-later auditor distinguish additive successor evidence from a rewrite?",
      "Which DI or TODO owns overlap/successor mechanics at the group-session layer?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

### `results/SIM-ludaf-udp-feed/makerspace-door-access/openai-gpt-5.4-xhigh/20260520-194221.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-ludaf-udp-feed-makerspace-door-access-openai-gpt-5.4-xhigh-20260520-194221",
  "run_group_id": "ga-canary-20260520-194221",
  "cell_id": "ga-canary-20260520-194221-000005-SIM-ludaf-udp-feed--makerspace-door-access--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-ludaf-udp-feed",
  "scenario_id": "makerspace-door-access",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-194221",
  "result_path": "results/SIM-ludaf-udp-feed/makerspace-door-access/openai-gpt-5.4-xhigh/20260520-194221.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_4d901411ad324bd686ec8fa0193535a2",
    "response_id": "resp_0c92a08b55e469ba006a0e0f1196f4819b87542fba4b07d69c",
    "input_tokens": 9656,
    "output_tokens": 13534,
    "cost_usd": 0.206374
  },
  "source": {
    "repo_commit": "d79e4669c410f21c4c0abf616422de8e169da5c3",
    "sim_path": "simulations/SIM-ludaf-udp-feed/",
    "scenario_path": "scenarios/makerspace-door-access/makerspace-door-access.md",
    "root_contract_paths": [
      "results/RUN-PROTOCOL.md",
      "scenarios/README.md"
    ],
    "files": [
      {
        "path": "results/RUN-PROTOCOL.md",
        "sha256": "9bf7fd58adfffcb4aeebb87bca14576ac9016681ed3aac0750f2ddafb00792a6"
      },
      {
        "path": "scenarios/README.md",
        "sha256": "406c4c7f400df14788d1caea61406f83adbc474160806ce4f1aa6a88d409d483"
      },
      {
        "path": "simulations/SIM-ludaf-udp-feed/README.md",
        "sha256": "6238e2781f9658feae574df275356e4e25a6d22eba5b8a2f12ae22b1561ef252"
      },
      {
        "path": "simulations/SIM-ludaf-udp-feed/QUESTION.md",
        "sha256": "9136722a6e7b982ace707418665d088de2e9f229e099ea4a646ad6c5196e5463"
      },
      {
        "path": "simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/CHANGELOG.md",
        "sha256": "59df64a3749a92b706582ebe7041e19ba6545cf4840eaf2715c9cedef55e07aa"
      },
      {
        "path": "simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/TODO/TODO-jodon-udp-binding-v0-reference-implementation.md",
        "sha256": "e63864119b92646905ff3519a4173074c3ce11256fc9d7ae3dcbbf99b7de7dd4"
      },
      {
        "path": "simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/TODO/TODO.md",
        "sha256": "cc9141e03e0e508984e6ce60108ce2e20d3fd440172275251303c3fd4a0445d1"
      },
      {
        "path": "simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/specs/udp-binding-draft.md",
        "sha256": "5305152a5372c37f26f81c77b0d02fb9b3e82dd431d452ad1f11bdd46c4f0c7d"
      },
      {
        "path": "simulations/SIM-ludaf-udp-feed/seed/protocol-migration.md",
        "sha256": "d835ef1b97297872e7ef2700464a9be5b68fdd31aca4b05483dbc11d98003669"
      },
      {
        "path": "scenarios/makerspace-door-access/makerspace-door-access.md",
        "sha256": "18baa89e4f187886d95c27d7a7bfe55db1277d9be7d9da2efa821587c69c994c"
      }
    ],
    "simulation_tree_hash": "1f626da6f8c3305ac8a02bdbbbb9f66cdcfdba59ebd91dff2433acf4ec107c93"
  },
  "rubric": {
    "rubric_version": "ga-rubric-20260519-v1",
    "score_scale": "0..5",
    "score_meanings": {
      "0": "no fit or absent",
      "5": "strong fit",
      "risk_penalty": "0 low risk, 5 severe risk"
    },
    "axes": [
      "scenario_fit",
      "promisegrid_alignment",
      "auditability",
      "evolution_safety",
      "layer_boundary_clarity",
      "failure_handling",
      "implementation_plausibility",
      "risk_penalty"
    ]
  },
  "scores": {
    "scenario_fit": 1,
    "promisegrid_alignment": 3,
    "auditability": 2,
    "evolution_safety": 3,
    "layer_boundary_clarity": 5,
    "failure_handling": 1,
    "implementation_plausibility": 2,
    "risk_penalty": 4
  },
  "fitness": {
    "raw": 18,
    "normalized_0_100": 45,
    "confidence_0_1": 0.88
  },
  "assessment": {
    "rationale": "Useful as a clearly bounded transport substrate, but not as a makerspace door-access design. The simulation documents UDP carriage, local wire artifacts, and versioning/forking expectations, yet it leaves identity, authorization, revocation, emergency override, and accountable door decisions to higher layers.",
    "strengths": [
      "Excellent layer boundary clarity: udp-feed is transport-only and says so explicitly.",
      "Registry-free, peer-local transport behavior and simulator artifact paths fit PromiseGrid basics.",
      "Rename/fork/version notes provide a workable, if still draft, evolution story."
    ],
    "weaknesses": [
      "No membership, guest, override, revocation, or door-policy objects are defined.",
      "Raw UDP provides no authenticity, privacy, reliability, ordering, or replay protection.",
      "Transport artifacts alone are not enough for Carol to audit a contested door decision.",
      "The protocol is still draft with no freeze or conformance evidence."
    ],
    "risks": [
      "If misapplied directly, spoofed, stale, or lost datagrams could cause unsafe unlock or deny outcomes.",
      "Best-effort delivery can hide missed revocations or emergency updates during partitions.",
      "Opaque byte logs may preserve wire evidence without preserving long-term human accountability."
    ],
    "open_questions": [
      "Which signed membership, guest-access, and emergency-override objects sit above udp-feed?",
      "How should a local door controller behave when authorization evidence is stale, missing, or disputed?",
      "What durable audit bundle links a received datagram to a justified door decision decades later?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

### `results/SIM-ludaf-udp-feed/community-movement-organizing/openai-gpt-5.4-xhigh/20260520-194221.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-ludaf-udp-feed-community-movement-organizing-openai-gpt-5.4-xhigh-20260520-194221",
  "run_group_id": "ga-canary-20260520-194221",
  "cell_id": "ga-canary-20260520-194221-000006-SIM-ludaf-udp-feed--community-movement-organizing--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-ludaf-udp-feed",
  "scenario_id": "community-movement-organizing",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-194221",
  "result_path": "results/SIM-ludaf-udp-feed/community-movement-organizing/openai-gpt-5.4-xhigh/20260520-194221.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_2bc6d99845364e53b7bdbe62f8b6968a",
    "response_id": "resp_07577195d8b401ee006a0e0f4b580c8198b2fe64b7ec970476",
    "input_tokens": 9662,
    "output_tokens": 10676,
    "cost_usd": 0.1663725
  },
  "source": {
    "repo_commit": "d79e4669c410f21c4c0abf616422de8e169da5c3",
    "sim_path": "simulations/SIM-ludaf-udp-feed/",
    "scenario_path": "scenarios/community-movement-organizing/community-movement-organizing.md",
    "root_contract_paths": [
      "results/RUN-PROTOCOL.md",
      "scenarios/README.md"
    ],
    "files": [
      {
        "path": "results/RUN-PROTOCOL.md",
        "sha256": "9bf7fd58adfffcb4aeebb87bca14576ac9016681ed3aac0750f2ddafb00792a6"
      },
      {
        "path": "scenarios/README.md",
        "sha256": "406c4c7f400df14788d1caea61406f83adbc474160806ce4f1aa6a88d409d483"
      },
      {
        "path": "simulations/SIM-ludaf-udp-feed/README.md",
        "sha256": "6238e2781f9658feae574df275356e4e25a6d22eba5b8a2f12ae22b1561ef252"
      },
      {
        "path": "simulations/SIM-ludaf-udp-feed/QUESTION.md",
        "sha256": "9136722a6e7b982ace707418665d088de2e9f229e099ea4a646ad6c5196e5463"
      },
      {
        "path": "simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/CHANGELOG.md",
        "sha256": "59df64a3749a92b706582ebe7041e19ba6545cf4840eaf2715c9cedef55e07aa"
      },
      {
        "path": "simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/TODO/TODO-jodon-udp-binding-v0-reference-implementation.md",
        "sha256": "e63864119b92646905ff3519a4173074c3ce11256fc9d7ae3dcbbf99b7de7dd4"
      },
      {
        "path": "simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/TODO/TODO.md",
        "sha256": "cc9141e03e0e508984e6ce60108ce2e20d3fd440172275251303c3fd4a0445d1"
      },
      {
        "path": "simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/specs/udp-binding-draft.md",
        "sha256": "5305152a5372c37f26f81c77b0d02fb9b3e82dd431d452ad1f11bdd46c4f0c7d"
      },
      {
        "path": "simulations/SIM-ludaf-udp-feed/seed/protocol-migration.md",
        "sha256": "d835ef1b97297872e7ef2700464a9be5b68fdd31aca4b05483dbc11d98003669"
      },
      {
        "path": "scenarios/community-movement-organizing/community-movement-organizing.md",
        "sha256": "c1b45264fd4c188f8a0c860b24691e4926c04fd2de12e4a7c7779054e012db7e"
      }
    ],
    "simulation_tree_hash": "1f626da6f8c3305ac8a02bdbbbb9f66cdcfdba59ebd91dff2433acf4ec107c93"
  },
  "rubric": {
    "rubric_version": "ga-rubric-20260519-v1",
    "score_scale": "0..5",
    "score_meanings": {
      "0": "no fit or absent",
      "5": "strong fit",
      "risk_penalty": "0 low risk, 5 severe risk"
    },
    "axes": [
      "scenario_fit",
      "promisegrid_alignment",
      "auditability",
      "evolution_safety",
      "layer_boundary_clarity",
      "failure_handling",
      "implementation_plausibility",
      "risk_penalty"
    ]
  },
  "scores": {
    "scenario_fit": 1,
    "promisegrid_alignment": 3,
    "auditability": 4,
    "evolution_safety": 3,
    "layer_boundary_clarity": 5,
    "failure_handling": 2,
    "implementation_plausibility": 3,
    "risk_penalty": 4
  },
  "fitness": {
    "raw": 17,
    "normalized_0_100": 49,
    "confidence_0_1": 0.86
  },
  "assessment": {
    "rationale": "A strong transport-layer PromiseGrid specimen but only a weak direct fit for community movement organizing. It clarifies decentralized feed semantics, local wire artifacts, and composition boundaries, yet leaves membership, moderation, public commitments, identity disputes, and schism handling to higher layers that are not present here.",
    "strengths": [
      "Clear layer contract with explicit promises, anti-promises, and simulator artifact paths.",
      "Good PromiseGrid discipline around decentralization, peer-local observation, and pCID-based forking/migration.",
      "Concrete TODO, test-vector, and harness plans make the draft auditable and plausibly buildable."
    ],
    "weaknesses": [
      "It does not model the application objects or decision rules needed for organizing scenarios.",
      "Authenticity, privacy, spam resistance, and durable promise accounting are explicitly out of scope.",
      "It remains a draft with no freeze event and no completed reference implementation or harness evidence."
    ],
    "risks": [
      "If misused as more than transport, lossy or spoofable UDP traffic could distort contested commitments or moderation events.",
      "host:port and DNS assumptions do not by themselves satisfy 100-year durability or registry-free discovery goals.",
      "Raw payload artifacts may be too weak for local dispute resolution without signed higher-layer records."
    ],
    "open_questions": [
      "Which session, identity, and message-layer specs above udp-feed would make organizing commitments auditable and authentic?",
      "How should moderation and schism evidence be represented so peers can decide locally under partition or stale data?",
      "What long-term discovery and addressing approach replaces fragile host:port assumptions?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

### `results/SIM-labit-feed-outer/cas-backed-group-session-additive-successor-specimen/openai-gpt-5.4-xhigh/20260520-194221.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-labit-feed-outer-cas-backed-group-session-additive-successor-specimen-openai-gpt-5.4-xhigh-20260520-194221",
  "run_group_id": "ga-canary-20260520-194221",
  "cell_id": "ga-canary-20260520-194221-000007-SIM-labit-feed-outer--cas-backed-group-session-additive-successor-specimen--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-labit-feed-outer",
  "scenario_id": "cas-backed-group-session-additive-successor-specimen",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-194221",
  "result_path": "results/SIM-labit-feed-outer/cas-backed-group-session-additive-successor-specimen/openai-gpt-5.4-xhigh/20260520-194221.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_94c8608b2e844a53ad508b42446ab6ad",
    "response_id": "resp_02fffdd686613aa5006a0e0fabb43c8198b32aba59ae5f5675",
    "input_tokens": 7882,
    "output_tokens": 12907,
    "cost_usd": 0.1944915
  },
  "source": {
    "repo_commit": "d79e4669c410f21c4c0abf616422de8e169da5c3",
    "sim_path": "simulations/SIM-labit-feed-outer/",
    "scenario_path": "scenarios/cas-backed-group-session-additive-successor-specimen/cas-backed-group-session-additive-successor-specimen.md",
    "root_contract_paths": [
      "results/RUN-PROTOCOL.md",
      "scenarios/README.md"
    ],
    "files": [
      {
        "path": "results/RUN-PROTOCOL.md",
        "sha256": "9bf7fd58adfffcb4aeebb87bca14576ac9016681ed3aac0750f2ddafb00792a6"
      },
      {
        "path": "scenarios/README.md",
        "sha256": "406c4c7f400df14788d1caea61406f83adbc474160806ce4f1aa6a88d409d483"
      },
      {
        "path": "simulations/SIM-labit-feed-outer/README.md",
        "sha256": "5df982634d6bf2c3911480fbf73f3ba7448222db5aad97081904fb76a991f959"
      },
      {
        "path": "simulations/SIM-labit-feed-outer/QUESTION.md",
        "sha256": "11331c6cd26618a9a2ebdecdacd90340da5ad45006d2d793c7754136df1afedc"
      },
      {
        "path": "simulations/SIM-labit-feed-outer/TODO/TODO-kakaz-feed-outer-freeze-and-cid-cascade-follow-on.md",
        "sha256": "f52dfefcc991894afe9491b049cd6fb4525df75399bd8512fd40cd6cfd028a6e"
      },
      {
        "path": "simulations/SIM-labit-feed-outer/TODO/TODO.md",
        "sha256": "3bb3f6416756d1d8ff9ba9805c4c1fb581086a207f8e9c9187b928b7621434cc"
      },
      {
        "path": "simulations/SIM-labit-feed-outer/protocols/feed-outer.d/specs/feed-outer-draft.md",
        "sha256": "c628890da6e3ee82585ea262452392cfe3d38ff6f9ca142275c3f00a850d6faa"
      },
      {
        "path": "simulations/SIM-labit-feed-outer/seed/extraction-sources.md",
        "sha256": "9a832d0836848fb42307935420bb940f6cccc3fb43e054873fbb115fa22290b3"
      },
      {
        "path": "scenarios/cas-backed-group-session-additive-successor-specimen/cas-backed-group-session-additive-successor-specimen.md",
        "sha256": "e1eaad54488dd151065a9023e1d348e7cd0cbc1c3475ffa5bc19fe0f5b1e1887"
      }
    ],
    "simulation_tree_hash": "060eb3fe3a7c63ee4e88f66a0c9af758b9e7d2219b68b87c5b132ff2fb7c50ff"
  },
  "rubric": {
    "rubric_version": "ga-rubric-20260519-v1",
    "score_scale": "0..5",
    "score_meanings": {
      "0": "no fit or absent",
      "5": "strong fit",
      "risk_penalty": "0 low risk, 5 severe risk"
    },
    "axes": [
      "scenario_fit",
      "promisegrid_alignment",
      "auditability",
      "evolution_safety",
      "layer_boundary_clarity",
      "failure_handling",
      "implementation_plausibility",
      "risk_penalty"
    ]
  },
  "scores": {
    "scenario_fit": 3,
    "promisegrid_alignment": 4,
    "auditability": 4,
    "evolution_safety": 5,
    "layer_boundary_clarity": 5,
    "failure_handling": 2,
    "implementation_plausibility": 3,
    "risk_penalty": 2
  },
  "fitness": {
    "raw": 24,
    "normalized_0_100": 73,
    "confidence_0_1": 0.84
  },
  "assessment": {
    "rationale": "The simulation is strong evidence that evolution can stay additive and preserve historical transport/feed evidence, but it only partially answers this scenario because the actual CAS-backed group-session overlap mechanics are outside this lineage's scope.",
    "strengths": [
      "The freeze boundary is explicitly additive: historical transport/feed specimens remain evidence and are not renamed, rehashed, or rewritten.",
      "Any frozen successor or derived mirror must cite its source evidence, which helps migration auditability.",
      "Layer boundaries are very clear: feed-outer owns transport naming and outer convention, not message-level CID, header, or body rules."
    ],
    "weaknesses": [
      "It does not define how old text specimens and new CAS-backed specimens are linked or co-read inside a group-session.",
      "Message-CID cascade, legacy compatibility, and reader-side migration policy are explicitly out of scope.",
      "Failure, stale-reader, and adversarial-mirror handling are only lightly specified."
    ],
    "risks": [
      "Teams may overread the outer-layer additive rule as a complete migration recipe when the message layer is still unresolved.",
      "The spec is still draft, so the final frozen pCID and operational migration path are not yet locked.",
      "Temporary reliance on a Steve-authored DI as the operative promise weakens decentralization until stronger promise tooling exists."
    ],
    "open_questions": [
      "Which transport-protocol spec records the old .txt to new CAS linkage and local audit trail?",
      "How should readers validate additive derived mirrors across draft and frozen phases?",
      "What local evidence proves equivalence without relying on a central registry or a single trusted maintainer?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

### `results/SIM-labit-feed-outer/makerspace-door-access/openai-gpt-5.4-xhigh/20260520-194221.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-labit-feed-outer-makerspace-door-access-openai-gpt-5.4-xhigh-20260520-194221",
  "run_group_id": "ga-canary-20260520-194221",
  "cell_id": "ga-canary-20260520-194221-000008-SIM-labit-feed-outer--makerspace-door-access--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-labit-feed-outer",
  "scenario_id": "makerspace-door-access",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-194221",
  "result_path": "results/SIM-labit-feed-outer/makerspace-door-access/openai-gpt-5.4-xhigh/20260520-194221.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_6300b9d3374d47cbad812cebbc708325",
    "response_id": "resp_02be4b73cad9b6d7006a0e0fdf0e508198a2eea8e21c6af79a",
    "input_tokens": 7954,
    "output_tokens": 10240,
    "cost_usd": 0.1572795
  },
  "source": {
    "repo_commit": "d79e4669c410f21c4c0abf616422de8e169da5c3",
    "sim_path": "simulations/SIM-labit-feed-outer/",
    "scenario_path": "scenarios/makerspace-door-access/makerspace-door-access.md",
    "root_contract_paths": [
      "results/RUN-PROTOCOL.md",
      "scenarios/README.md"
    ],
    "files": [
      {
        "path": "results/RUN-PROTOCOL.md",
        "sha256": "9bf7fd58adfffcb4aeebb87bca14576ac9016681ed3aac0750f2ddafb00792a6"
      },
      {
        "path": "scenarios/README.md",
        "sha256": "406c4c7f400df14788d1caea61406f83adbc474160806ce4f1aa6a88d409d483"
      },
      {
        "path": "simulations/SIM-labit-feed-outer/README.md",
        "sha256": "5df982634d6bf2c3911480fbf73f3ba7448222db5aad97081904fb76a991f959"
      },
      {
        "path": "simulations/SIM-labit-feed-outer/QUESTION.md",
        "sha256": "11331c6cd26618a9a2ebdecdacd90340da5ad45006d2d793c7754136df1afedc"
      },
      {
        "path": "simulations/SIM-labit-feed-outer/TODO/TODO-kakaz-feed-outer-freeze-and-cid-cascade-follow-on.md",
        "sha256": "f52dfefcc991894afe9491b049cd6fb4525df75399bd8512fd40cd6cfd028a6e"
      },
      {
        "path": "simulations/SIM-labit-feed-outer/TODO/TODO.md",
        "sha256": "3bb3f6416756d1d8ff9ba9805c4c1fb581086a207f8e9c9187b928b7621434cc"
      },
      {
        "path": "simulations/SIM-labit-feed-outer/protocols/feed-outer.d/specs/feed-outer-draft.md",
        "sha256": "c628890da6e3ee82585ea262452392cfe3d38ff6f9ca142275c3f00a850d6faa"
      },
      {
        "path": "simulations/SIM-labit-feed-outer/seed/extraction-sources.md",
        "sha256": "9a832d0836848fb42307935420bb940f6cccc3fb43e054873fbb115fa22290b3"
      },
      {
        "path": "scenarios/makerspace-door-access/makerspace-door-access.md",
        "sha256": "18baa89e4f187886d95c27d7a7bfe55db1277d9be7d9da2efa821587c69c994c"
      }
    ],
    "simulation_tree_hash": "060eb3fe3a7c63ee4e88f66a0c9af758b9e7d2219b68b87c5b132ff2fb7c50ff"
  },
  "rubric": {
    "rubric_version": "ga-rubric-20260519-v1",
    "score_scale": "0..5",
    "score_meanings": {
      "0": "no fit or absent",
      "5": "strong fit",
      "risk_penalty": "0 low risk, 5 severe risk"
    },
    "axes": [
      "scenario_fit",
      "promisegrid_alignment",
      "auditability",
      "evolution_safety",
      "layer_boundary_clarity",
      "failure_handling",
      "implementation_plausibility",
      "risk_penalty"
    ]
  },
  "scores": {
    "scenario_fit": 0,
    "promisegrid_alignment": 3,
    "auditability": 3,
    "evolution_safety": 4,
    "layer_boundary_clarity": 5,
    "failure_handling": 1,
    "implementation_plausibility": 3,
    "risk_penalty": 2
  },
  "fitness": {
    "raw": 13,
    "normalized_0_100": 37,
    "confidence_0_1": 0.89
  },
  "assessment": {
    "rationale": "Strong transport-layer boundary evidence, but not a makerspace door-access design. It preserves local, additive protocol evidence and avoids central registry assumptions, yet leaves membership, authorization, guest/override policy, identity, and contested-decision handling to unspecified higher layers, so overall fitness is capped by near-zero direct scenario coverage.",
    "strengths": [
      "Very clear outer-vs-interior layer boundary: transport naming is separated from message and application semantics.",
      "Good evolution posture: freeze is additive, historical specimens are not rewritten, and provenance is explicit.",
      "Plausible low-level convention for local evidence storage and handler dispatch via protocol pCID."
    ],
    "weaknesses": [
      "Does not model membership state, door authorization, guest access, or emergency override.",
      "Does not specify scenario-level identity claims, promise accounting records, or local access-decision rules under stale or disputed evidence.",
      "Failure handling is mostly out of scope; accept/retry/downgrade/escalate behavior is undefined."
    ],
    "risks": [
      "A reviewer could over-read this thin outer-feed convention as sufficient for authorization when it only defines transport-level structure.",
      "Code-as-handler without a machine-readable schema may make long-term cross-implementation interoperability harder."
    ],
    "open_questions": [
      "Which higher-layer protocol would carry door-access authorization, guest grants, and audit records on top of this outer feed?",
      "How should Alice or Carol decide locally when membership evidence is stale, partitioned, or contested?",
      "How are emergency override actions recorded so later auditors can distinguish valid override from abuse without central authority?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

### `results/SIM-labit-feed-outer/community-movement-organizing/openai-gpt-5.4-xhigh/20260520-194221.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-labit-feed-outer-community-movement-organizing-openai-gpt-5.4-xhigh-20260520-194221",
  "run_group_id": "ga-canary-20260520-194221",
  "cell_id": "ga-canary-20260520-194221-000009-SIM-labit-feed-outer--community-movement-organizing--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-labit-feed-outer",
  "scenario_id": "community-movement-organizing",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-194221",
  "result_path": "results/SIM-labit-feed-outer/community-movement-organizing/openai-gpt-5.4-xhigh/20260520-194221.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_41e670c9c41142359045b80e44401990",
    "response_id": "resp_0d8f8888c93324f6006a0e0ff7728c819a9e517c99b2cf1ac8",
    "input_tokens": 7960,
    "output_tokens": 14028,
    "cost_usd": 0.210322
  },
  "source": {
    "repo_commit": "d79e4669c410f21c4c0abf616422de8e169da5c3",
    "sim_path": "simulations/SIM-labit-feed-outer/",
    "scenario_path": "scenarios/community-movement-organizing/community-movement-organizing.md",
    "root_contract_paths": [
      "results/RUN-PROTOCOL.md",
      "scenarios/README.md"
    ],
    "files": [
      {
        "path": "results/RUN-PROTOCOL.md",
        "sha256": "9bf7fd58adfffcb4aeebb87bca14576ac9016681ed3aac0750f2ddafb00792a6"
      },
      {
        "path": "scenarios/README.md",
        "sha256": "406c4c7f400df14788d1caea61406f83adbc474160806ce4f1aa6a88d409d483"
      },
      {
        "path": "simulations/SIM-labit-feed-outer/README.md",
        "sha256": "5df982634d6bf2c3911480fbf73f3ba7448222db5aad97081904fb76a991f959"
      },
      {
        "path": "simulations/SIM-labit-feed-outer/QUESTION.md",
        "sha256": "11331c6cd26618a9a2ebdecdacd90340da5ad45006d2d793c7754136df1afedc"
      },
      {
        "path": "simulations/SIM-labit-feed-outer/TODO/TODO-kakaz-feed-outer-freeze-and-cid-cascade-follow-on.md",
        "sha256": "f52dfefcc991894afe9491b049cd6fb4525df75399bd8512fd40cd6cfd028a6e"
      },
      {
        "path": "simulations/SIM-labit-feed-outer/TODO/TODO.md",
        "sha256": "3bb3f6416756d1d8ff9ba9805c4c1fb581086a207f8e9c9187b928b7621434cc"
      },
      {
        "path": "simulations/SIM-labit-feed-outer/protocols/feed-outer.d/specs/feed-outer-draft.md",
        "sha256": "c628890da6e3ee82585ea262452392cfe3d38ff6f9ca142275c3f00a850d6faa"
      },
      {
        "path": "simulations/SIM-labit-feed-outer/seed/extraction-sources.md",
        "sha256": "9a832d0836848fb42307935420bb940f6cccc3fb43e054873fbb115fa22290b3"
      },
      {
        "path": "scenarios/community-movement-organizing/community-movement-organizing.md",
        "sha256": "c1b45264fd4c188f8a0c860b24691e4926c04fd2de12e4a7c7779054e012db7e"
      }
    ],
    "simulation_tree_hash": "060eb3fe3a7c63ee4e88f66a0c9af758b9e7d2219b68b87c5b132ff2fb7c50ff"
  },
  "rubric": {
    "rubric_version": "ga-rubric-20260519-v1",
    "score_scale": "0..5",
    "score_meanings": {
      "0": "no fit or absent",
      "5": "strong fit",
      "risk_penalty": "0 low risk, 5 severe risk"
    },
    "axes": [
      "scenario_fit",
      "promisegrid_alignment",
      "auditability",
      "evolution_safety",
      "layer_boundary_clarity",
      "failure_handling",
      "implementation_plausibility",
      "risk_penalty"
    ]
  },
  "scores": {
    "scenario_fit": 1,
    "promisegrid_alignment": 4,
    "auditability": 4,
    "evolution_safety": 5,
    "layer_boundary_clarity": 5,
    "failure_handling": 1,
    "implementation_plausibility": 4,
    "risk_penalty": 3
  },
  "fitness": {
    "raw": 21,
    "normalized_0_100": 60,
    "confidence_0_1": 0.87
  },
  "assessment": {
    "rationale": "The feed-outer lineage is strong substrate evidence for 100-year durability, sparse-knowledge operation, and non-central transport identification, but it is intentionally too thin for community movement organizing. It cleanly states where transport duties stop while leaving identity, membership, campaigns, moderation, public commitments, schism handling, and dispute resolution to higher-layer protocols.",
    "strengths": [
      "Peers can infer transport from local directory placement instead of trusting per-message transport claims.",
      "pCID-named protocol directories and additive freeze rules support durable, non-central evolution without rewriting history.",
      "The spec is explicit about what it does not own, giving very clear outer-vs-inner layer boundaries."
    ],
    "weaknesses": [
      "It does not define identity claims, membership or campaign objects, moderation actions, public commitments, or promise accounting.",
      "It provides no peer-local accept/retry/downgrade/escalate behavior for delayed, stale, or disputed community events.",
      "It is still draft-state and depends on per-protocol handlers rather than a shared machine-readable schema."
    ],
    "risks": [
      "Application builders may mistake a thin transport convention for sufficient movement-governance semantics.",
      "The interim Steve-authored DI as the operative merge promise is a temporary authority concentration until cryptographic promises exist.",
      "Long-term replay and audit may depend on preserving handler code and readable specs for each transport pCID."
    ],
    "open_questions": [
      "Which higher-layer protocol on top of feed-outer carries membership, campaigns, working groups, moderation, and public commitments?",
      "What local CAS objects and promise-accounting records would Alice and Carol retain for contested votes, membership changes, or moderation actions?",
      "How are schisms or competing lineages represented so peers can reason locally without hidden global state?",
      "How will the design replace DI-based merge authority with cryptographic promises while preserving 100-year auditability?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

## Required JSON Shape

{"child_id":"SIM-gumuj-ga-child-0002","design_delta_summary":"one to three bounded design deltas","files":[{"path":"README.md","content":"# ..."},{"path":"QUESTION.md","content":"# ..."}]}
