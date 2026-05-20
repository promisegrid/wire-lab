# GA Child Generation

Return only JSON with keys `child_id`, `design_delta_summary`, and `files`.
Each file path must be relative to the child simulation root. Include `README.md` and `QUESTION.md`.

Optimization goal: breed a child simulation from exactly two parent simulations, expected to score higher than its parent set on the same rubric and sampled scenarios.
Use the fitness evidence below as training feedback: preserve parent strengths, repair weaknesses, reduce risks, answer or route open questions, and keep changes to one to three bounded design deltas.
Do not merely summarize the parent. The child must make an explicit design move that should improve `fitness.normalized_0_100` while keeping the simulation standalone and auditable.

- Run group ID: `ga-canary-20260520-203341`
- Child ID: `SIM-lavoh-ga-child-0001`
- Child path: `simulations/SIM-lavoh-ga-child-0001/`
- Operation: `breed`
- Parent IDs: `SIM-hugoj-cas-usenetlike-gitlike, SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes`

## Scenario Sample

- `aerospace-design` at `scenarios/aerospace-design/aerospace-design.md`
- `chunking-identity-bakeoff-raw-only-migration` at `scenarios/chunking-identity-bakeoff-raw-only-migration/chunking-identity-bakeoff-raw-only-migration.md`
- `open-source-development` at `scenarios/open-source-development/open-source-development.md`

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
- Score parent cells before child generation when using `tools/ga-runner`.
  `generate` uses completed parent fitness evidence to rank the selected parent
  pool and apply deterministic top-parent plus tournament-diversity parent
  selection. Source: `DI-bukid`.
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

LLM-generated children use one operation, `breed`, with exactly two distinct
parent simulations. The runner must fail or skip generation rather than create a
one-parent child. Source: `DI-sohus`.

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

### `simulations/SIM-hugoj-cas-usenetlike-gitlike/README.md`

```markdown
# SIM-hugoj: CAS / Usenet-like / Git-like lineage

This simulation explores whether part of PromiseGrid should be
understood as a content-addressed-Usenet successor line with git-like
hashed and hash-chained objects and transport separation. It is a
broad design exploration, not a frozen protocol claim and not a claim
that PromiseGrid simply *is* Usenet or git. Source: `DI-pijun`.

`group-session` is the current worked specimen inside this simulation, not the
identity of the simulation itself. The point of the simulation is to ask which
historical invariants age well enough to carry forward into PromiseGrid, which
ones should be adapted, and which ones should be rejected. Source: `DI-pijun`.

## Question

Which parts of the Usenet / FidoNet / email lineage, when combined with
content-addressed storage and git-like transport separation, describe a useful
PromiseGrid design line rather than a misleading analogy? Source: `DI-pijun`.

## Current specimen

The current specimen is `group-session`, because it already exhibits several of
the relevant traits: append-only message growth, DAG/thread structure,
content-addressed message identity, and multiple possible delivery substrates.
That makes it a useful worked example, but not the whole design space. A later
PromiseGrid protocol may preserve some of these traits while diverging sharply
from `group-session` in governance, payload shape, replication, or feed
mechanics. Source: `DI-pijun`.

## What maps cleanly from precedent

- Message identity that stays stable across delivery substrates.
- Store-and-forward replication among named peers or sites.
- Per-instance declaration of how a site exchanges messages with peers.
- Thread/DAG semantics that survive relay across heterogeneous wires.
- Separation between protocol meaning and the substrate that carries bytes.

These are the parts of the turn-173 precedent survey that look structurally
durable across email, Usenet, FidoNet, and modern git-like or gRPC/libp2p-like
systems. See `docs/research/historical-networks-20260503.md`. Source:
`DI-pijun`.

## What does not map cleanly

- Usenet control-message conventions do not automatically become PromiseGrid
  governance rules.
- Git's object model and transport behavior are precedent, not identity; using
  git as a current substrate does not imply PromiseGrid should clone git's
  object types or workflow.
- A content-addressed-Usenet framing does not by itself settle feed naming,
  site manifests, sparse-CAS shape, moderation, or freeze behavior.
- The simulation does not assume that every future PromiseGrid protocol is a
  `group-session` variant.

These are precisely the places where the analogy is useful only if it remains
explicitly exploratory. Source: `DI-pijun`.

## Provenance and authority

The historical grounding for this simulation comes from
`docs/research/historical-networks-20260503.md`, especially its Usenet,
FidoNet, email, and git-adjacent precedent analysis. The immediate trigger for
filing this simulation is the turn-173 replay recovery in
`protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`
and the unresolved follow-on notes in
`protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`.

Nothing in this directory overrides later DRs, DIs, frozen specs, or guide
prose. This simulation is a bounded design workspace for testing one broad
lineage claim and keeping the claim visible until it either graduates or is
rejected. Source: `DI-pijun`.
```

### `scenarios/aerospace-design/aerospace-design.md`

```markdown
# Aerospace Design

## Scenario ID

aerospace-design

## Source / Provenance

- Source type: application seed
- Source path: `protocols/wire-lab.d/TODO/TODO-dadub-root-scenario-skeleton-and-seed-catalog.md`
- Source row/title: Seed application catalog entry `aerospace-design`
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-midif`; `TODO-dadub`

## Purpose

Exercise PromiseGrid design candidates against aerospace-design application pressure:
Requirements, analyses, design reviews, configuration control, and long-lived
engineering evidence.

## Setup

Alice depends on an outcome in the Aerospace Design domain. Bob makes promises about
requirements, analyses, design reviews, configuration control, and long-lived
engineering evidence. Carol needs enough evidence to rely on Bob's promise without
having complete global state, and Mallory may exploit stale, missing, or ambiguous
evidence.

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

### `scenarios/chunking-identity-bakeoff-raw-only-migration/chunking-identity-bakeoff-raw-only-migration.md`

```markdown
# Raw-only migration

## Scenario ID

chunking-identity-bakeoff-raw-only-migration

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-gobaz-chunking-identity-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-gobaz-chunking-identity-bakeoff/`
- Source row/title: Raw-only migration
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-gobaz-chunking-identity-bakeoff/`.

## Setup

Alice migrates one historical message as a single raw chunk behind a pointer object.

## Stimulus

Run the candidate simulation against this source test: Whether the first CAS migration can proceed without chunked Merkle roots.

## Expected Pressure

Raw-only may unblock migration but does not answer large-object replication.
```

### `scenarios/open-source-development/open-source-development.md`

```markdown
# Open Source Development

## Scenario ID

open-source-development

## Source / Provenance

- Source type: application seed
- Source path: `protocols/wire-lab.d/TODO/TODO-dadub-root-scenario-skeleton-and-seed-catalog.md`
- Source row/title: Seed application catalog entry `open-source-development`
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-midif`; `TODO-dadub`

## Purpose

Exercise PromiseGrid design candidates against open-source-development application
pressure: Issues, patches, review promises, release artifacts, maintainer authority, and
contributor reputation.

## Setup

Alice depends on an outcome in the Open Source Development domain. Bob makes promises
about issues, patches, review promises, release artifacts, maintainer authority, and
contributor reputation. Carol needs enough evidence to rely on Bob's promise without
having complete global state, and Mallory may exploit stale, missing, or ambiguous
evidence.

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

### `simulations/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/README.md`

```markdown
# SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes: Grid-envelope variant

This simulation is a standalone positional grid-envelope specimen. It tests the
combination `enc-cbor`, `unknown-best-effort`, and `sig-mandatory-opaque-bytes` without claiming
that this combination is the canonical PromiseGrid wire format. Source: `DI-fanah`.

The local draft spec is
`protocols/grid-envelope.d/specs/grid-envelope-draft.md`.
```

### `simulations/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/QUESTION.md`

```markdown
# Question

Does a positional grid envelope using `enc-cbor`, `unknown-best-effort`, and
`sig-mandatory-opaque-bytes` satisfy the wire-lab harness scenarios better than the sibling
variants? Source: `DI-fanah`.
```

### `simulations/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/protocols/grid-envelope.d/CHANGELOG.md`

```markdown
# CHANGELOG: grid-envelope

A-side CHANGELOG (per TE-liviv) for this simulation-local `grid-envelope`
protocol specimen.

This file records freeze events authored by the specimen maintainers. No entries
yet; this protocol specimen has not reached a first freeze.

This protocol tree is a simulation-local specimen created by `DI-fanah`.
```

### `simulations/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/protocols/grid-envelope.d/specs/grid-envelope-draft.md`

```markdown
# Grid Envelope Variant Spec (DRAFT)

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `enc-cbor` / `unknown-best-effort` / `sig-mandatory-opaque-bytes`.
> Source: `DI-fanah`.

## Purpose

This spec defines one full positional grid-envelope candidate for wire-lab
comparison. It is a specimen inside `SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes`, not a harness rule and not the
canonical PromiseGrid wire format.

## Positional Envelope Shape

The envelope shape for this variant is:

```text
[pcid, payload, signature]
```

Slots are interpreted positionally:

- `pcid` identifies the protocol/spec/handler that interprets `payload`.
- `payload` is opaque bytes until interpreted by the handler named by `pcid`.
- `signature` is mandatory opaque bytes over the canonical unsigned prefix.

A `payload` may itself be the canonical bytes of another grid envelope when the
protocol named by `pcid` specifies recursive nesting. The outer grid envelope
does not prescribe the payload's internal organization beyond the bytes boundary.

## Encoding

This variant encodes the envelope as deterministic CBOR positional arrays. Slot values use definite-length encodings. `pcid` and `sig_pcid`, when present, are CIDv1 byte strings; `payload`, `signature`, and `sig_payload` are byte strings. The canonical bytes for signing and hashing are the deterministic CBOR bytes of the exact positional array under this spec.

## Unknown pCID Policy

If a receiver lacks a handler for `pcid`, it may expose `payload` bytes to generic tooling for inspection or salvage. Any such result MUST be marked unsupported and unverified; best-effort inspection does not count as interpretation under the missing `pcid` rules.

## Signature and Authorship Policy

The third positional slot, `signature`, is mandatory opaque bytes. The bytes cover the canonical unsigned prefix `[pcid, payload]` under this variant's encoding. The envelope layer enforces presence and byte-string shape; signature algorithm, signer identity, and verification semantics are determined by the protocol ecosystem being tested with this variant.

## Layering-Test Behavior

This variant answers the harness §1.3 layering scenarios as follows:

- Ordering disagreements are handled by the protocol named by `pcid`; the grid
  envelope only preserves the bytes and dispatch identity needed to make failures
  explicit.
- Forwarding, relay, or hop-local evidence is represented either by wrapper
  grid envelopes, by the payload protocol, or by the signature slots available in
  this variant.
- External or content-addressed body references live inside `payload` under the
  protocol named by `pcid`; the envelope can carry those payload bytes without
  understanding them.
- Incompatible interpretation rules fail visibly at the `pcid` dispatch boundary
  or under this variant's unknown-pCID policy.

## Non-Goals

This draft does not declare a winning envelope, does not define a central pCID
registry, does not freeze a final PromiseGrid signing scheme, and does not make
sibling grid-envelope variants obsolete.

## Freeze Gate

This draft can freeze only after at least one simulation run compares it against
sibling positional grid-envelope variants and Steve signs a merge/freeze promise
for this specific specimen.
```

## Existing Fitness Evidence From This Run

### `results/SIM-hugoj-cas-usenetlike-gitlike/aerospace-design/openai-gpt-5.4-xhigh/20260520-203341.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-hugoj-cas-usenetlike-gitlike-aerospace-design-openai-gpt-5.4-xhigh-20260520-203341",
  "run_group_id": "ga-canary-20260520-203341",
  "cell_id": "ga-canary-20260520-203341-000001-SIM-hugoj-cas-usenetlike-gitlike--aerospace-design--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-hugoj-cas-usenetlike-gitlike",
  "scenario_id": "aerospace-design",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-203341",
  "result_path": "results/SIM-hugoj-cas-usenetlike-gitlike/aerospace-design/openai-gpt-5.4-xhigh/20260520-203341.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_4fe906b9eecf4b0db37117eefcbd098f",
    "response_id": "resp_00652f0bee6de283006a0e1aa7736c81988a77bed3583cfa70",
    "input_tokens": 5907,
    "output_tokens": 8635,
    "cost_usd": 0.13122725
  },
  "source": {
    "repo_commit": "39aac816b4af60ab2490c21c49fcc18fdbd68771",
    "sim_path": "simulations/SIM-hugoj-cas-usenetlike-gitlike/",
    "scenario_path": "scenarios/aerospace-design/aerospace-design.md",
    "root_contract_paths": [
      "results/RUN-PROTOCOL.md",
      "scenarios/README.md"
    ],
    "files": [
      {
        "path": "results/RUN-PROTOCOL.md",
        "sha256": "c6e5388b5635afb39230dc34c5577b5bec8f0171d13868da2c660a615b58eb18"
      },
      {
        "path": "scenarios/README.md",
        "sha256": "406c4c7f400df14788d1caea61406f83adbc474160806ce4f1aa6a88d409d483"
      },
      {
        "path": "simulations/SIM-hugoj-cas-usenetlike-gitlike/README.md",
        "sha256": "0506d32e30f1fe4903d1c0c1f2e6a7f25a3c3a438d6cfa0782d3e99c9448bdb8"
      },
      {
        "path": "scenarios/aerospace-design/aerospace-design.md",
        "sha256": "b0b74ad6985b28c92cdb3c57323432c7e0125f80b1c9ba6f676f089e30d75d9a"
      }
    ],
    "simulation_tree_hash": "144951c8a39cbb610ac9d305fe823308c60e01e271c2f3b222a32276d4a196ef"
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
    "auditability": 3,
    "evolution_safety": 4,
    "layer_boundary_clarity": 5,
    "failure_handling": 3,
    "implementation_plausibility": 4,
    "risk_penalty": 3
  },
  "fitness": {
    "raw": 28,
    "normalized_0_100": 70,
    "confidence_0_1": 0.75
  },
  "assessment": {
    "rationale": "Strong substrate-level fit for durable, partition-tolerant engineering evidence, but still under-specifies aerospace-specific traceability, authority, and configuration-control semantics.",
    "strengths": [
      "Content-addressed, hash-linked objects support long-lived evidence and configuration history.",
      "Store-and-forward replication and transport separation fit sparse knowledge, partitions, and 100-year migration pressure.",
      "The simulation is unusually clear about boundaries between precedent, protocol meaning, and transport layers."
    ],
    "weaknesses": [
      "The current specimen is message/thread oriented rather than an explicit requirements-analysis-review model.",
      "Signature workflows, delegated authority, baselines, and freeze rules are not defined here.",
      "Peer-local accept/retry/downgrade/escalate logic for disputed evidence is under-specified."
    ],
    "risks": [
      "A Usenet/git analogy could be overextended into governance or certification semantics it does not actually provide.",
      "Durable storage could be mistaken for sufficient aerospace configuration control without stronger traceability and approval records."
    ],
    "open_questions": [
      "What CAS object and feed shapes represent requirements, analyses, review dispositions, and baselines?",
      "How are signatures, delegated authority, and local decision policies recorded during partitions or stale-evidence windows?",
      "What freeze and migration rules preserve decades-long auditability without central registries?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

### `results/SIM-hugoj-cas-usenetlike-gitlike/chunking-identity-bakeoff-raw-only-migration/openai-gpt-5.4-xhigh/20260520-203341.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-hugoj-cas-usenetlike-gitlike-chunking-identity-bakeoff-raw-only-migration-openai-gpt-5.4-xhigh-20260520-203341",
  "run_group_id": "ga-canary-20260520-203341",
  "cell_id": "ga-canary-20260520-203341-000002-SIM-hugoj-cas-usenetlike-gitlike--chunking-identity-bakeoff-raw-only-migration--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-hugoj-cas-usenetlike-gitlike",
  "scenario_id": "chunking-identity-bakeoff-raw-only-migration",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-203341",
  "result_path": "results/SIM-hugoj-cas-usenetlike-gitlike/chunking-identity-bakeoff-raw-only-migration/openai-gpt-5.4-xhigh/20260520-203341.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_99d3fb7979f643ab88a2af1dc11d2cd8",
    "response_id": "resp_0170a072f3c6172c006a0e1aa779748199b51fd843fac92821",
    "input_tokens": 5858,
    "output_tokens": 9795,
    "cost_usd": 0.1473815
  },
  "source": {
    "repo_commit": "39aac816b4af60ab2490c21c49fcc18fdbd68771",
    "sim_path": "simulations/SIM-hugoj-cas-usenetlike-gitlike/",
    "scenario_path": "scenarios/chunking-identity-bakeoff-raw-only-migration/chunking-identity-bakeoff-raw-only-migration.md",
    "root_contract_paths": [
      "results/RUN-PROTOCOL.md",
      "scenarios/README.md"
    ],
    "files": [
      {
        "path": "results/RUN-PROTOCOL.md",
        "sha256": "c6e5388b5635afb39230dc34c5577b5bec8f0171d13868da2c660a615b58eb18"
      },
      {
        "path": "scenarios/README.md",
        "sha256": "406c4c7f400df14788d1caea61406f83adbc474160806ce4f1aa6a88d409d483"
      },
      {
        "path": "simulations/SIM-hugoj-cas-usenetlike-gitlike/README.md",
        "sha256": "0506d32e30f1fe4903d1c0c1f2e6a7f25a3c3a438d6cfa0782d3e99c9448bdb8"
      },
      {
        "path": "scenarios/chunking-identity-bakeoff-raw-only-migration/chunking-identity-bakeoff-raw-only-migration.md",
        "sha256": "8d0616b2a8d070f1066fb55cdea847a46fdb5db98fd280cd8169d0b199564d77"
      }
    ],
    "simulation_tree_hash": "144951c8a39cbb610ac9d305fe823308c60e01e271c2f3b222a32276d4a196ef"
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
    "promisegrid_alignment": 3,
    "auditability": 4,
    "evolution_safety": 4,
    "layer_boundary_clarity": 5,
    "failure_handling": 2,
    "implementation_plausibility": 4,
    "risk_penalty": 3
  },
  "fitness": {
    "raw": 23,
    "normalized_0_100": 66,
    "confidence_0_1": 0.71
  },
  "assessment": {
    "rationale": "The simulation can support a raw-only first CAS migration because it centers stable content-addressed message identity across changing delivery substrates, but it does not specify pointer-object semantics, chunking/Merkle upgrades, or sparse-CAS behavior, so it only partially covers the scenario and leaves large-object migration unresolved.",
    "strengths": [
      "Substrate-independent message identity supports long-lived migration of historical messages.",
      "Hash-addressed and hash-chained objects plus DAG/thread structure give strong audit anchors.",
      "Protocol meaning vs transport separation is unusually clear, which helps later evolution."
    ],
    "weaknesses": [
      "No explicit chunking or Merkle-root strategy is defined.",
      "Sparse-CAS shape and pointer/wrapper object semantics are left unsettled.",
      "Peer-local failure and adversarial evidence handling are only thinly described."
    ],
    "risks": [
      "A raw-only bootstrap could become accidental long-term identity policy without a clean upgrade path.",
      "Large-object replication and partial retrieval may diverge across implementations.",
      "Migration evidence may be hard to compare if wrapper objects are defined inconsistently."
    ],
    "open_questions": [
      "What exact pointer object links the historical message to its raw CAS blob?",
      "How should a later chunked form relate to the first raw-only identity?",
      "What local records let Alice and Carol audit migration without any central registry?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

### `results/SIM-hugoj-cas-usenetlike-gitlike/open-source-development/openai-gpt-5.4-xhigh/20260520-203341.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-hugoj-cas-usenetlike-gitlike-open-source-development-openai-gpt-5.4-xhigh-20260520-203341",
  "run_group_id": "ga-canary-20260520-203341",
  "cell_id": "ga-canary-20260520-203341-000003-SIM-hugoj-cas-usenetlike-gitlike--open-source-development--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-hugoj-cas-usenetlike-gitlike",
  "scenario_id": "open-source-development",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-203341",
  "result_path": "results/SIM-hugoj-cas-usenetlike-gitlike/open-source-development/openai-gpt-5.4-xhigh/20260520-203341.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_6e2f6269eba741eea9c3ad7ebad16c09",
    "response_id": "resp_04e72d58fb1cfc6a006a0e1aa778dc819bad0a589a53ea93b6",
    "input_tokens": 5910,
    "output_tokens": 8979,
    "cost_usd": 0.1360485
  },
  "source": {
    "repo_commit": "39aac816b4af60ab2490c21c49fcc18fdbd68771",
    "sim_path": "simulations/SIM-hugoj-cas-usenetlike-gitlike/",
    "scenario_path": "scenarios/open-source-development/open-source-development.md",
    "root_contract_paths": [
      "results/RUN-PROTOCOL.md",
      "scenarios/README.md"
    ],
    "files": [
      {
        "path": "results/RUN-PROTOCOL.md",
        "sha256": "c6e5388b5635afb39230dc34c5577b5bec8f0171d13868da2c660a615b58eb18"
      },
      {
        "path": "scenarios/README.md",
        "sha256": "406c4c7f400df14788d1caea61406f83adbc474160806ce4f1aa6a88d409d483"
      },
      {
        "path": "simulations/SIM-hugoj-cas-usenetlike-gitlike/README.md",
        "sha256": "0506d32e30f1fe4903d1c0c1f2e6a7f25a3c3a438d6cfa0782d3e99c9448bdb8"
      },
      {
        "path": "scenarios/open-source-development/open-source-development.md",
        "sha256": "0965a20ce8888de4b4b3f3ec8bc07d0efa493a54f319f792d59f244231000c7f"
      }
    ],
    "simulation_tree_hash": "144951c8a39cbb610ac9d305fe823308c60e01e271c2f3b222a32276d4a196ef"
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
    "evolution_safety": 4,
    "layer_boundary_clarity": 5,
    "failure_handling": 3,
    "implementation_plausibility": 4,
    "risk_penalty": 3
  },
  "fitness": {
    "raw": 24,
    "normalized_0_100": 69,
    "confidence_0_1": 0.74
  },
  "assessment": {
    "rationale": "Strong substrate-level fit for decentralized open-source evidence flows: durable CAS identities, append-only DAG/thread history, store-and-forward replication, and transport separation match sparse, multi-site development well. The fit is only partial at the application layer because maintainer authority, contributor reputation, review commitments, and release trust are not yet modeled as explicit peer-local promises.",
    "strengths": [
      "Durable content-addressed and hash-linked history supports replayable audit trails across transport changes.",
      "Store-and-forward replication and substrate independence fit partitioned, multi-site collaboration without requiring central hosting.",
      "The simulation is clear about boundaries between historical precedent, current specimen, and any future PromiseGrid protocol, which helps evolution."
    ],
    "weaknesses": [
      "Maintainer authority, contributor reputation, and review promises are not defined as first-class local records.",
      "The current specimen is discussion-centric and does not spell out patch, merge, or release workflows.",
      "Peer-local accept, retry, downgrade, and escalation rules for stale or disputed evidence are not explicit."
    ],
    "risks": [
      "The Usenet/git analogy could under-specify governance, moderation, and release policy needed for real open-source projects.",
      "Unresolved naming, feed, and sparse-CAS decisions could push deployments back toward central forges or registries.",
      "If authority and reputation remain implicit, different peers may make incompatible trust decisions with weak interoperability."
    ],
    "open_questions": [
      "How are issues, patches, reviews, and release artifacts represented as distinct CAS objects or feeds?",
      "How does a peer record maintainer authority and contributor reputation locally without a central forge or registry?",
      "What retry, downgrade, or escalation rules apply when signatures, reviews, or release evidence are delayed, stale, or disputed?",
      "How do maintainer and release claims survive key, organization, or hosting changes over long time spans?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

### `results/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/aerospace-design/openai-gpt-5.4-xhigh/20260520-203341.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes-aerospace-design-openai-gpt-5.4-xhigh-20260520-203341",
  "run_group_id": "ga-canary-20260520-203341",
  "cell_id": "ga-canary-20260520-203341-000004-SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes--aerospace-design--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes",
  "scenario_id": "aerospace-design",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-203341",
  "result_path": "results/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/aerospace-design/openai-gpt-5.4-xhigh/20260520-203341.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_618f7e87811c4f7fa7f3f0a3e2949886",
    "response_id": "resp_0acacd3c42fcdb53006a0e1b254584819b92dcec19145d43ec",
    "input_tokens": 6260,
    "output_tokens": 13436,
    "cost_usd": 0.199059
  },
  "source": {
    "repo_commit": "39aac816b4af60ab2490c21c49fcc18fdbd68771",
    "sim_path": "simulations/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/",
    "scenario_path": "scenarios/aerospace-design/aerospace-design.md",
    "root_contract_paths": [
      "results/RUN-PROTOCOL.md",
      "scenarios/README.md"
    ],
    "files": [
      {
        "path": "results/RUN-PROTOCOL.md",
        "sha256": "c6e5388b5635afb39230dc34c5577b5bec8f0171d13868da2c660a615b58eb18"
      },
      {
        "path": "scenarios/README.md",
        "sha256": "406c4c7f400df14788d1caea61406f83adbc474160806ce4f1aa6a88d409d483"
      },
      {
        "path": "simulations/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/README.md",
        "sha256": "577753a742867694e2a8ee84c229e8264f04fa109942f444c3af2feda056a45c"
      },
      {
        "path": "simulations/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/QUESTION.md",
        "sha256": "eaebc9d28329655b84215f609f1d60f7a6267bccf424f314a2be1aef0b1f764f"
      },
      {
        "path": "simulations/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/protocols/grid-envelope.d/CHANGELOG.md",
        "sha256": "aef2b7915320af9607f415cef649397530561e4f0f9372fda1bbf2691dd3f9b1"
      },
      {
        "path": "simulations/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/protocols/grid-envelope.d/specs/grid-envelope-draft.md",
        "sha256": "30aadb2f95a3906e1c5c27acb96ea3687a88346a2d6e2d171db2d90c55bbe656"
      },
      {
        "path": "scenarios/aerospace-design/aerospace-design.md",
        "sha256": "b0b74ad6985b28c92cdb3c57323432c7e0125f80b1c9ba6f676f089e30d75d9a"
      }
    ],
    "simulation_tree_hash": "d31569a6a684e574940f7156736bbf02b5789289cb2520469fd57807458c1a18"
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
    "scenario_fit": 2,
    "promisegrid_alignment": 3,
    "auditability": 2,
    "evolution_safety": 3,
    "layer_boundary_clarity": 5,
    "failure_handling": 3,
    "implementation_plausibility": 4,
    "risk_penalty": 3
  },
  "fitness": {
    "raw": 24,
    "normalized_0_100": 60,
    "confidence_0_1": 0.72
  },
  "assessment": {
    "rationale": "Unknown-best-effort and deterministic CBOR make this a useful recovery-friendly shell under sparse knowledge, but aerospace-design disputes need durable signer, review, and configuration semantics that this draft leaves to opaque higher layers.",
    "strengths": [
      "Deterministic CBOR gives stable canonical bytes for storage, hashing, and local dispute comparison.",
      "Unknown-pCID best-effort salvage helps long-horizon recovery when a peer lacks a handler, while requiring unsupported/unverified labeling.",
      "The envelope cleanly separates transport bytes and dispatch identity from payload-level domain semantics.",
      "The draft avoids dependence on a central pCID registry, which fits no-central-authority goals."
    ],
    "weaknesses": [
      "It does not itself model requirements traceability, analyses, design-review decisions, configuration baselines, or promise-accounting records.",
      "Signature bytes are mandatory but opaque; algorithm, signer role, identity, and verification semantics are external.",
      "A single envelope signature slot is thin for multi-party aerospace approvals, countersignatures, timestamps, and partial attestations.",
      "The specimen is still draft and unfrozen, weakening century-scale assurance."
    ],
    "risks": [
      "Operators may over-trust best-effort inspection of unknown payloads during contested evidence cases.",
      "If payload handlers or signature conventions are lost, long-lived artifacts may survive as bytes but not as usable proof.",
      "Critical control semantics may get buried in opaque payload conventions, reducing cross-peer comparability."
    ],
    "open_questions": [
      "Should the envelope carry an explicit signature-protocol identifier to preserve verification across tool and algorithm changes?",
      "What payload conventions provide traceability links, review outcomes, configuration control, and local promise-accounting records?",
      "How should peers downgrade, retry, or escalate when required objects or signatures exist only as stale or unverified bytes?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

### `results/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/chunking-identity-bakeoff-raw-only-migration/openai-gpt-5.4-xhigh/20260520-203341.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes-chunking-identity-bakeoff-raw-only-migration-openai-gpt-5.4-xhigh-20260520-203341",
  "run_group_id": "ga-canary-20260520-203341",
  "cell_id": "ga-canary-20260520-203341-000005-SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes--chunking-identity-bakeoff-raw-only-migration--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes",
  "scenario_id": "chunking-identity-bakeoff-raw-only-migration",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-203341",
  "result_path": "results/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/chunking-identity-bakeoff-raw-only-migration/openai-gpt-5.4-xhigh/20260520-203341.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_062573b79e2d4593bf78e24b24523c98",
    "response_id": "resp_0e4e3c233f18f42f006a0e1b2d09ac819aa51de91871235187",
    "input_tokens": 6211,
    "output_tokens": 15356,
    "cost_usd": 0.22585325
  },
  "source": {
    "repo_commit": "39aac816b4af60ab2490c21c49fcc18fdbd68771",
    "sim_path": "simulations/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/",
    "scenario_path": "scenarios/chunking-identity-bakeoff-raw-only-migration/chunking-identity-bakeoff-raw-only-migration.md",
    "root_contract_paths": [
      "results/RUN-PROTOCOL.md",
      "scenarios/README.md"
    ],
    "files": [
      {
        "path": "results/RUN-PROTOCOL.md",
        "sha256": "c6e5388b5635afb39230dc34c5577b5bec8f0171d13868da2c660a615b58eb18"
      },
      {
        "path": "scenarios/README.md",
        "sha256": "406c4c7f400df14788d1caea61406f83adbc474160806ce4f1aa6a88d409d483"
      },
      {
        "path": "simulations/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/README.md",
        "sha256": "577753a742867694e2a8ee84c229e8264f04fa109942f444c3af2feda056a45c"
      },
      {
        "path": "simulations/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/QUESTION.md",
        "sha256": "eaebc9d28329655b84215f609f1d60f7a6267bccf424f314a2be1aef0b1f764f"
      },
      {
        "path": "simulations/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/protocols/grid-envelope.d/CHANGELOG.md",
        "sha256": "aef2b7915320af9607f415cef649397530561e4f0f9372fda1bbf2691dd3f9b1"
      },
      {
        "path": "simulations/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/protocols/grid-envelope.d/specs/grid-envelope-draft.md",
        "sha256": "30aadb2f95a3906e1c5c27acb96ea3687a88346a2d6e2d171db2d90c55bbe656"
      },
      {
        "path": "scenarios/chunking-identity-bakeoff-raw-only-migration/chunking-identity-bakeoff-raw-only-migration.md",
        "sha256": "8d0616b2a8d070f1066fb55cdea847a46fdb5db98fd280cd8169d0b199564d77"
      }
    ],
    "simulation_tree_hash": "d31569a6a684e574940f7156736bbf02b5789289cb2520469fd57807458c1a18"
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
    "promisegrid_alignment": 3,
    "auditability": 2,
    "evolution_safety": 3,
    "layer_boundary_clarity": 5,
    "failure_handling": 3,
    "implementation_plausibility": 4,
    "risk_penalty": 2
  },
  "fitness": {
    "raw": 21,
    "normalized_0_100": 65,
    "confidence_0_1": 0.72
  },
  "assessment": {
    "rationale": "At the wire-envelope layer this specimen can carry a pointer object that references a single raw CAS chunk without requiring chunked Merkle roots, but it leaves pointer schema, chunk identity, replication, and signature verification to other protocols, so it is a partial rather than complete fit for the migration scenario.",
    "strengths": [
      "It can wrap a pointer object that references one raw CAS chunk without imposing chunked Merkle roots at the envelope layer.",
      "Deterministic CBOR and explicit pCID dispatch let peers record exact local evidence without hidden global state or a central registry.",
      "Unknown-pCID best-effort salvage offers a bounded migration fallback while still requiring unsupported/unverified labeling."
    ],
    "weaknesses": [
      "The specimen does not define the pointer-object protocol, chunk identity rules, or any answer for large-object replication.",
      "Signature bytes are mandatory but semantically opaque, so signer identity and verification durability depend on external ecosystems.",
      "The draft is not frozen and the specimen pCID is not yet minted, so long-term migration claims remain provisional."
    ],
    "risks": [
      "Operators may mistake wire-level transport support for a complete chunking-identity migration design.",
      "Best-effort inspection of unknown pCIDs can be over-trusted despite the required unsupported/unverified boundary.",
      "Mandatory signatures may complicate first migration of historical objects with weak or absent provenance."
    ],
    "open_questions": [
      "What payload protocol carries the pointer object and raw-chunk identity claims?",
      "If the source artifact is unsigned or provenance-poor, what should the mandatory envelope signature mean?",
      "How are signature algorithms, signer identities, and verification rules preserved over 100-year horizons?",
      "How does the design evolve from raw-only migration to chunked large-object replication without breaking old envelopes?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

### `results/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/open-source-development/openai-gpt-5.4-xhigh/20260520-203341.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes-open-source-development-openai-gpt-5.4-xhigh-20260520-203341",
  "run_group_id": "ga-canary-20260520-203341",
  "cell_id": "ga-canary-20260520-203341-000006-SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes--open-source-development--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes",
  "scenario_id": "open-source-development",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-203341",
  "result_path": "results/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/open-source-development/openai-gpt-5.4-xhigh/20260520-203341.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_a8f1513f1a234402b1fc62945db9f279",
    "response_id": "resp_0eba01522760a9be006a0e1b37e758819ab43d754a742045a8",
    "input_tokens": 6263,
    "output_tokens": 17857,
    "cost_usd": 0.26095825
  },
  "source": {
    "repo_commit": "39aac816b4af60ab2490c21c49fcc18fdbd68771",
    "sim_path": "simulations/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/",
    "scenario_path": "scenarios/open-source-development/open-source-development.md",
    "root_contract_paths": [
      "results/RUN-PROTOCOL.md",
      "scenarios/README.md"
    ],
    "files": [
      {
        "path": "results/RUN-PROTOCOL.md",
        "sha256": "c6e5388b5635afb39230dc34c5577b5bec8f0171d13868da2c660a615b58eb18"
      },
      {
        "path": "scenarios/README.md",
        "sha256": "406c4c7f400df14788d1caea61406f83adbc474160806ce4f1aa6a88d409d483"
      },
      {
        "path": "simulations/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/README.md",
        "sha256": "577753a742867694e2a8ee84c229e8264f04fa109942f444c3af2feda056a45c"
      },
      {
        "path": "simulations/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/QUESTION.md",
        "sha256": "eaebc9d28329655b84215f609f1d60f7a6267bccf424f314a2be1aef0b1f764f"
      },
      {
        "path": "simulations/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/protocols/grid-envelope.d/CHANGELOG.md",
        "sha256": "aef2b7915320af9607f415cef649397530561e4f0f9372fda1bbf2691dd3f9b1"
      },
      {
        "path": "simulations/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/protocols/grid-envelope.d/specs/grid-envelope-draft.md",
        "sha256": "30aadb2f95a3906e1c5c27acb96ea3687a88346a2d6e2d171db2d90c55bbe656"
      },
      {
        "path": "scenarios/open-source-development/open-source-development.md",
        "sha256": "0965a20ce8888de4b4b3f3ec8bc07d0efa493a54f319f792d59f244231000c7f"
      }
    ],
    "simulation_tree_hash": "d31569a6a684e574940f7156736bbf02b5789289cb2520469fd57807458c1a18"
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
    "scenario_fit": 2,
    "promisegrid_alignment": 3,
    "auditability": 2,
    "evolution_safety": 4,
    "layer_boundary_clarity": 5,
    "failure_handling": 3,
    "implementation_plausibility": 4,
    "risk_penalty": 3
  },
  "fitness": {
    "raw": 25,
    "normalized_0_100": 62.5,
    "confidence_0_1": 0.74
  },
  "assessment": {
    "rationale": "Useful only as a lower-layer carrier for open-source artifacts with mandatory signature bytes: deterministic CBOR, explicit pcid dispatch, and unsupported/unverified downgrade behavior work under partial knowledge. It does not model maintainer authority, contributor reputation, promise accounting, freshness, or dispute semantics, so application fit remains limited.",
    "strengths": [
      "Deterministic CBOR plus a mandatory signature slot over a fixed prefix give stable canonical bytes with attached signature material for local storage, hashing, and archival.",
      "Explicit pcid dispatch and the lack of a central registry fit sparse-knowledge and no-central-authority goals.",
      "Unknown-pcid best-effort salvage lets older tools retain and inspect new artifacts while forcing unsupported/unverified labeling.",
      "Envelope, payload, and signature responsibilities are separated clearly."
    ],
    "weaknesses": [
      "It does not define issue, patch, review promise, release artifact, maintainer authority, or contributor reputation objects.",
      "Mandatory signature bytes do not carry in-envelope signer identity, algorithm, or verification rules.",
      "Opaque payloads and positional CBOR reduce generic human/LLM auditability without the referenced handler spec.",
      "No peer-local promise-accounting, freshness, revocation, retry, or escalation rules are provided for contested development workflows."
    ],
    "risks": [
      "Consumers may over-trust the mere presence of signature bytes when the relevant verification regime is missing or disputed.",
      "Best-effort inspection of unknown payloads can produce misleading partial interpretations if tooling hides the unsupported/unverified state.",
      "Long-term verification depends on preserving external handler and signature-verification specs, not just the envelope bytes."
    ],
    "open_questions": [
      "What payload protocol encodes review state, maintainer authority, and contributor reputation for this application?",
      "How are key rotation, revocation, and project-specific trust roots represented across decades?",
      "The encoding text mentions sig_pcid and sig_payload when present; are extra signature conventions or wrappers intended beyond the three-slot envelope?",
      "Are wrapper conventions needed for timestamps, countersignatures, relay evidence, or merge/release attestations?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

## Required JSON Shape

{"child_id":"SIM-lavoh-ga-child-0001","design_delta_summary":"one to three bounded design deltas","files":[{"path":"README.md","content":"# ..."},{"path":"QUESTION.md","content":"# ..."}]}
