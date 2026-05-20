# GA Child Generation

Return only JSON with keys `child_id`, `design_delta_summary`, and `files`.
Each file path must be relative to the child simulation root. Include `README.md` and `QUESTION.md`.

- Run group ID: `ga-canary-20260520-194221`
- Child ID: `SIM-mofor-ga-child-0001`
- Child path: `simulations/SIM-mofor-ga-child-0001/`
- Operation: `mutation`
- Parent IDs: `SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid`

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

### `simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/README.md`

```markdown
# SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid: Grid-envelope variant

This simulation is a standalone positional grid-envelope specimen. It tests the
combination `enc-dag-cbor`, `unknown-best-effort`, and `sig-wrapper-pcid` without claiming
that this combination is the canonical PromiseGrid wire format. Source: `DI-fanah`.

The local draft spec is
`protocols/grid-envelope.d/specs/grid-envelope-draft.md`.
```

### `simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/QUESTION.md`

```markdown
# Question

Does a positional grid envelope using `enc-dag-cbor`, `unknown-best-effort`, and
`sig-wrapper-pcid` satisfy the wire-lab harness scenarios better than the sibling
variants? Source: `DI-fanah`.
```

### `simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/protocols/grid-envelope.d/CHANGELOG.md`

```markdown
# CHANGELOG: grid-envelope

A-side CHANGELOG (per TE-liviv) for this simulation-local `grid-envelope`
protocol specimen.

This file records freeze events authored by the specimen maintainers. No entries
yet; this protocol specimen has not reached a first freeze.

This protocol tree is a simulation-local specimen created by `DI-fanah`.
```

### `simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/protocols/grid-envelope.d/specs/grid-envelope-draft.md`

```markdown
# Grid Envelope Variant Spec (DRAFT)

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `enc-dag-cbor` / `unknown-best-effort` / `sig-wrapper-pcid`.
> Source: `DI-fanah`.

## Purpose

This spec defines one full positional grid-envelope candidate for wire-lab
comparison. It is a specimen inside `SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid`, not a harness rule and not the
canonical PromiseGrid wire format.

## Positional Envelope Shape

The envelope shape for this variant is:

```text
[pcid, payload]
```

Slots are interpreted positionally:

- `pcid` identifies the protocol/spec/handler that interprets `payload`.
- `payload` is opaque bytes until interpreted by the handler named by `pcid`.

A `payload` may itself be the canonical bytes of another grid envelope when the
protocol named by `pcid` specifies recursive nesting. The outer grid envelope
does not prescribe the payload's internal organization beyond the bytes boundary.

## Encoding

This variant encodes the envelope as DAG-CBOR-compatible positional arrays. `pcid` and `sig_pcid`, when present, are DAG-CBOR Link values; `payload`, `signature`, and `sig_payload` are byte strings. The envelope remains positional: no map/object envelope fields are introduced. The canonical bytes for signing and hashing are the DAG-CBOR bytes of the exact positional array under this spec.

## Unknown pCID Policy

If a receiver lacks a handler for `pcid`, it may expose `payload` bytes to generic tooling for inspection or salvage. Any such result MUST be marked unsupported and unverified; best-effort inspection does not count as interpretation under the missing `pcid` rules.

## Signature and Authorship Policy

The base envelope has no fixed signature slot. Signatures, encryption, authorship, or hop evidence are represented by outer or inner grid envelopes whose own `pcid` selects the relevant signature or evidence protocol. This keeps the envelope shape minimal and tests whether pCID-selected wrapper protocols are enough for authorship and integrity.

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

## Existing Fitness Evidence From This Run

### `results/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/cas-backed-group-session-additive-successor-specimen/openai-gpt-5.4-xhigh/20260520-194221.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid-cas-backed-group-session-additive-successor-specimen-openai-gpt-5.4-xhigh-20260520-194221",
  "run_group_id": "ga-canary-20260520-194221",
  "cell_id": "ga-canary-20260520-194221-000001-SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid--cas-backed-group-session-additive-successor-specimen--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid",
  "scenario_id": "cas-backed-group-session-additive-successor-specimen",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-194221",
  "result_path": "results/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/cas-backed-group-session-additive-successor-specimen/openai-gpt-5.4-xhigh/20260520-194221.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_f6bec73a064b43278be62d98bccb37d0",
    "response_id": "resp_0bed908d3a6ab2a3006a0e0ea063c8819aaa796a460367b8a8",
    "input_tokens": 6064,
    "output_tokens": 7671,
    "cost_usd": 0.118006
  },
  "source": {
    "repo_commit": "d79e4669c410f21c4c0abf616422de8e169da5c3",
    "sim_path": "simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/",
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
        "path": "simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/README.md",
        "sha256": "caf1ee91081288fb0196938506b45525f5ed257ffdc0694d1d1ee30dd7c42efa"
      },
      {
        "path": "simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/QUESTION.md",
        "sha256": "4a1249953f1687b3a8e60505b36073e09b23c64438a583dd0af0ecdde4550c26"
      },
      {
        "path": "simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/protocols/grid-envelope.d/CHANGELOG.md",
        "sha256": "aef2b7915320af9607f415cef649397530561e4f0f9372fda1bbf2691dd3f9b1"
      },
      {
        "path": "simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/protocols/grid-envelope.d/specs/grid-envelope-draft.md",
        "sha256": "70ccda1d3f064c237b55cfa517cdf00538315c6ba5aa3a39668acc096478d071"
      },
      {
        "path": "scenarios/cas-backed-group-session-additive-successor-specimen/cas-backed-group-session-additive-successor-specimen.md",
        "sha256": "e1eaad54488dd151065a9023e1d348e7cd0cbc1c3475ffa5bc19fe0f5b1e1887"
      }
    ],
    "simulation_tree_hash": "cfea4f21c240950e568c6e59cdc5e33ecb81989d8d42c9188a1b52c94462cd6c"
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
    "auditability": 3,
    "evolution_safety": 4,
    "layer_boundary_clarity": 5,
    "failure_handling": 3,
    "implementation_plausibility": 4,
    "risk_penalty": 3
  },
  "fitness": {
    "raw": 26,
    "normalized_0_100": 65,
    "confidence_0_1": 0.73
  },
  "assessment": {
    "rationale": "This envelope variant is a plausible lower-layer aid for additive migration because it can preserve historical bytes, wrap old content, and add new evidence without rewriting prior material. But it does not itself define the successor, overlap, CID-linking, or group-session rules that this scenario primarily tests, so the fit is only partial.",
    "strengths": [
      "Clear [pcid, payload] boundary keeps legacy bytes opaque and untouched while allowing new envelopes around them.",
      "Wrapper-based signatures/evidence support additive provenance instead of mutating historical content.",
      "Unknown-pCID handling allows salvage while explicitly marking results unsupported and unverified."
    ],
    "weaknesses": [
      "The spec does not define explicit successor-specimen or overlap mechanics for historical .txt messages and existing CIDs.",
      "Migration meaning depends on external payload protocols and wrapper conventions rather than the envelope itself.",
      "The specimen is still draft and not frozen, reducing long-term stability of the exact migration anchor."
    ],
    "risks": [
      "Best-effort inspection could be mistaken for verified interpretation if tooling blurs the unsupported/unverified boundary.",
      "Wrapper-centric signature/authorship rules may fragment verification behavior or enable omission/downgrade mistakes.",
      "Different implementations may choose inconsistent legacy-reference patterns, weakening additive-migration guarantees."
    ],
    "open_questions": [
      "What canonical payload or wrapper protocol links historical .txt message CIDs to successor CAS-backed session objects?",
      "How should auditors verify lineage and authorship across nested wrappers when the base envelope has no fixed signature slot?",
      "What freezes the pCID and successor conventions so future peers can still resolve and audit them?",
      "How are peer-local promise-accounting records represented during mixed old/new deployments?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

### `results/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/makerspace-door-access/openai-gpt-5.4-xhigh/20260520-194221.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid-makerspace-door-access-openai-gpt-5.4-xhigh-20260520-194221",
  "run_group_id": "ga-canary-20260520-194221",
  "cell_id": "ga-canary-20260520-194221-000002-SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid--makerspace-door-access--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid",
  "scenario_id": "makerspace-door-access",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-194221",
  "result_path": "results/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/makerspace-door-access/openai-gpt-5.4-xhigh/20260520-194221.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_04b41a46d5bb4536b1d3111360e76034",
    "response_id": "resp_0cf785c4755144de006a0e0e9edd74819bb7940e46ffb8541c",
    "input_tokens": 6136,
    "output_tokens": 10589,
    "cost_usd": 0.158984
  },
  "source": {
    "repo_commit": "d79e4669c410f21c4c0abf616422de8e169da5c3",
    "sim_path": "simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/",
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
        "path": "simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/README.md",
        "sha256": "caf1ee91081288fb0196938506b45525f5ed257ffdc0694d1d1ee30dd7c42efa"
      },
      {
        "path": "simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/QUESTION.md",
        "sha256": "4a1249953f1687b3a8e60505b36073e09b23c64438a583dd0af0ecdde4550c26"
      },
      {
        "path": "simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/protocols/grid-envelope.d/CHANGELOG.md",
        "sha256": "aef2b7915320af9607f415cef649397530561e4f0f9372fda1bbf2691dd3f9b1"
      },
      {
        "path": "simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/protocols/grid-envelope.d/specs/grid-envelope-draft.md",
        "sha256": "70ccda1d3f064c237b55cfa517cdf00538315c6ba5aa3a39668acc096478d071"
      },
      {
        "path": "scenarios/makerspace-door-access/makerspace-door-access.md",
        "sha256": "18baa89e4f187886d95c27d7a7bfe55db1277d9be7d9da2efa821587c69c994c"
      }
    ],
    "simulation_tree_hash": "cfea4f21c240950e568c6e59cdc5e33ecb81989d8d42c9188a1b52c94462cd6c"
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
    "layer_boundary_clarity": 4,
    "failure_handling": 2,
    "implementation_plausibility": 4,
    "risk_penalty": 3
  },
  "fitness": {
    "raw": 22,
    "normalized_0_100": 55,
    "confidence_0_1": 0.72
  },
  "assessment": {
    "rationale": "Good wire-level fit for sparse-knowledge, no-central-registry exchange of contested door-access evidence: [pcid,payload] is simple, canonical DAG-CBOR bytes aid stable hashing, and unknown pCIDs fail explicitly. But this is only an envelope specimen; membership, guest, override, revocation, freshness, and audit semantics remain outside the design, so local access decisions depend on additional wrapper and payload protocols.",
    "strengths": [
      "Clear minimal boundary with explicit pCID dispatch and canonical DAG-CBOR encoding.",
      "Unknown-pCID handling requires unsupported/unverified labeling, which helps safe salvage under partial knowledge.",
      "Wrapper-based signatures/evidence can evolve without hard-coding a single central registry or fixed signing scheme."
    ],
    "weaknesses": [
      "Does not define the application objects and rules needed for membership, guest access, emergency override, or local promise accounting.",
      "No fixed signature or audit slot makes consistent accountability depend on extra wrapper conventions.",
      "Draft status and unminted pCID weaken present-day 100-year durability confidence."
    ],
    "risks": [
      "Operators or tooling may over-trust best-effort inspection of unknown payloads as authorization evidence.",
      "Different wrapper conventions may produce inconsistent door decisions or audit trails across devices.",
      "Fail-closed behavior on unknown or rotated pCIDs could cause lockouts during migration or partitions."
    ],
    "open_questions": [
      "Which pCIDs and payload schemas represent grants, revocations, guest passes, and emergency overrides?",
      "What local deny/degrade/escalate policy should a door controller use when wrappers, handlers, or freshness proofs are missing or stale?",
      "How are replay protection and door-event audit records linked when signatures live only in wrapper protocols?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

### `results/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/community-movement-organizing/openai-gpt-5.4-xhigh/20260520-194221.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid-community-movement-organizing-openai-gpt-5.4-xhigh-20260520-194221",
  "run_group_id": "ga-canary-20260520-194221",
  "cell_id": "ga-canary-20260520-194221-000003-SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid--community-movement-organizing--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid",
  "scenario_id": "community-movement-organizing",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-194221",
  "result_path": "results/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/community-movement-organizing/openai-gpt-5.4-xhigh/20260520-194221.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_50b41118bb3f46278d924df874838719",
    "response_id": "resp_0ea42f792aa2555f006a0e0e9ee06881999bf2bbbe7cf417ec",
    "input_tokens": 6142,
    "output_tokens": 7812,
    "cost_usd": 0.1201165
  },
  "source": {
    "repo_commit": "d79e4669c410f21c4c0abf616422de8e169da5c3",
    "sim_path": "simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/",
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
        "path": "simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/README.md",
        "sha256": "caf1ee91081288fb0196938506b45525f5ed257ffdc0694d1d1ee30dd7c42efa"
      },
      {
        "path": "simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/QUESTION.md",
        "sha256": "4a1249953f1687b3a8e60505b36073e09b23c64438a583dd0af0ecdde4550c26"
      },
      {
        "path": "simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/protocols/grid-envelope.d/CHANGELOG.md",
        "sha256": "aef2b7915320af9607f415cef649397530561e4f0f9372fda1bbf2691dd3f9b1"
      },
      {
        "path": "simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/protocols/grid-envelope.d/specs/grid-envelope-draft.md",
        "sha256": "70ccda1d3f064c237b55cfa517cdf00538315c6ba5aa3a39668acc096478d071"
      },
      {
        "path": "scenarios/community-movement-organizing/community-movement-organizing.md",
        "sha256": "c1b45264fd4c188f8a0c860b24691e4926c04fd2de12e4a7c7779054e012db7e"
      }
    ],
    "simulation_tree_hash": "cfea4f21c240950e568c6e59cdc5e33ecb81989d8d42c9188a1b52c94462cd6c"
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
    "failure_handling": 3,
    "implementation_plausibility": 4,
    "risk_penalty": 3
  },
  "fitness": {
    "raw": 18,
    "normalized_0_100": 51,
    "confidence_0_1": 0.68
  },
  "assessment": {
    "rationale": "Useful as a minimal wire boundary for disputed local evidence, but weak as direct support for community movement organizing because membership, moderation, public commitments, schism handling, and promise accounting are all delegated to unspecified payload and wrapper protocols. The variant is reasonably aligned with PromiseGrid decentralization and evolution goals through pCID dispatch, opaque payloads, and explicit unsupported/unverified handling, but draft status and wrapper-dependent signatures limit auditability for contested social records.",
    "strengths": [
      "Very clear layering: the base envelope only dispatches by pCID and leaves payload semantics to the selected protocol.",
      "Canonical DAG-CBOR bytes and explicit unsupported/unverified marking help prevent silent misinterpretation.",
      "Wrapper-based signatures avoid hardwiring a central registry or single signing scheme."
    ],
    "weaknesses": [
      "The spec does not define the identity claims, feeds, names, moderation records, or promise-accounting artifacts the scenario needs.",
      "No fixed signature/authorship slot means public commitments depend on extra wrapper conventions.",
      "The spec is still a draft with no freeze and no minted pCID, reducing long-term comparability."
    ],
    "risks": [
      "Best-effort inspection of unknown payloads could mislead users if tooling blurs the unsupported/unverified boundary.",
      "Divergent wrapper conventions may fragment verification across factions or schisms.",
      "Long-term audit depends on preserving external pCID/spec material outside the envelope bytes."
    ],
    "open_questions": [
      "Which payload and wrapper pCIDs carry membership, moderation, and schism evidence for this application?",
      "How do Alice and Carol record locally verifiable outcomes when signatures live in optional wrappers rather than the base envelope?",
      "How are pCID definitions and handler specs discovered and preserved for 100-year audit without reintroducing central authority?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

## Required JSON Shape

{"child_id":"SIM-mofor-ga-child-0001","design_delta_summary":"one to three bounded design deltas","files":[{"path":"README.md","content":"# ..."},{"path":"QUESTION.md","content":"# ..."}]}
