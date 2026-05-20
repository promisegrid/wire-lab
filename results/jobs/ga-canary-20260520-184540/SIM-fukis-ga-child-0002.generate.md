# GA Child Generation

Return only JSON with keys `child_id`, `design_delta_summary`, and `files`.
Each file path must be relative to the child simulation root. Include `README.md` and `QUESTION.md`.

- Run group ID: `ga-canary-20260520-184540`
- Child ID: `SIM-fukis-ga-child-0002`
- Child path: `simulations/SIM-fukis-ga-child-0002/`
- Operation: `crossover`
- Parent IDs: `SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload, SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid`

## Scenario Sample

- `insurance-claims` at `scenarios/insurance-claims/insurance-claims.md`
- `state-governance` at `scenarios/state-governance/state-governance.md`
- `website-backend-hosting` at `scenarios/website-backend-hosting/website-backend-hosting.md`

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
  `-max-cell-estimate-usd` plus `-max-output-tokens` to prevent accidentally
  starting cells whose worst-case estimate is too high.
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
4. Review token/cost fields in the state file and tune output caps or result
   style.
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

### `simulations/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/README.md`

```markdown
# SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload: Grid-envelope variant

This simulation is a standalone positional grid-envelope specimen. It tests the
combination `enc-cbor`, `unknown-best-effort`, and `sig-mandatory-sig-pcid-payload` without claiming
that this combination is the canonical PromiseGrid wire format. Source: `DI-fanah`.

The local draft spec is
`protocols/grid-envelope.d/specs/grid-envelope-draft.md`.
```

### `simulations/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/QUESTION.md`

```markdown
# Question

Does a positional grid envelope using `enc-cbor`, `unknown-best-effort`, and
`sig-mandatory-sig-pcid-payload` satisfy the wire-lab harness scenarios better than the sibling
variants? Source: `DI-fanah`.
```

### `simulations/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/protocols/grid-envelope.d/CHANGELOG.md`

```markdown
# CHANGELOG: grid-envelope

A-side CHANGELOG (per TE-liviv) for this simulation-local `grid-envelope`
protocol specimen.

This file records freeze events authored by the specimen maintainers. No entries
yet; this protocol specimen has not reached a first freeze.

This protocol tree is a simulation-local specimen created by `DI-fanah`.
```

### `simulations/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/protocols/grid-envelope.d/specs/grid-envelope-draft.md`

```markdown
# Grid Envelope Variant Spec (DRAFT)

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `enc-cbor` / `unknown-best-effort` / `sig-mandatory-sig-pcid-payload`.
> Source: `DI-fanah`.

## Purpose

This spec defines one full positional grid-envelope candidate for wire-lab
comparison. It is a specimen inside `SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload`, not a harness rule and not the
canonical PromiseGrid wire format.

## Positional Envelope Shape

The envelope shape for this variant is:

```text
[pcid, payload, sig_pcid, sig_payload]
```

Slots are interpreted positionally:

- `pcid` identifies the protocol/spec/handler that interprets `payload`.
- `payload` is opaque bytes until interpreted by the handler named by `pcid`.
- `sig_pcid` identifies the signature or proof protocol.
- `sig_payload` is opaque bytes for the handler named by `sig_pcid`.

A `payload` may itself be the canonical bytes of another grid envelope when the
protocol named by `pcid` specifies recursive nesting. The outer grid envelope
does not prescribe the payload's internal organization beyond the bytes boundary.

## Encoding

This variant encodes the envelope as deterministic CBOR positional arrays. Slot values use definite-length encodings. `pcid` and `sig_pcid`, when present, are CIDv1 byte strings; `payload`, `signature`, and `sig_payload` are byte strings. The canonical bytes for signing and hashing are the deterministic CBOR bytes of the exact positional array under this spec.

## Unknown pCID Policy

If a receiver lacks a handler for `pcid`, it may expose `payload` bytes to generic tooling for inspection or salvage. Any such result MUST be marked unsupported and unverified; best-effort inspection does not count as interpretation under the missing `pcid` rules.

## Signature and Authorship Policy

The third and fourth positional slots are mandatory. `sig_pcid` identifies the signature or proof protocol; `sig_payload` is opaque bytes interpreted by that signature protocol. The signature payload covers the canonical unsigned prefix `[pcid, payload]` under this variant's encoding unless `sig_pcid` publishes stricter rules.

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

### `scenarios/insurance-claims/insurance-claims.md`

```markdown
# Insurance Claims

## Scenario ID

insurance-claims

## Source / Provenance

- Source type: application seed
- Source path: `protocols/wire-lab.d/TODO/TODO-dadub-root-scenario-skeleton-and-seed-catalog.md`
- Source row/title: Seed application catalog entry `insurance-claims`
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-midif`; `TODO-dadub`

## Purpose

Exercise PromiseGrid design candidates against insurance-claims application pressure:
Claim evidence, adjuster authority, fraud pressure, payments, and appeal promises.

## Setup

Alice depends on an outcome in the Insurance Claims domain. Bob makes promises about
claim evidence, adjuster authority, fraud pressure, payments, and appeal promises. Carol
needs enough evidence to rely on Bob's promise without having complete global state, and
Mallory may exploit stale, missing, or ambiguous evidence.

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

### `scenarios/state-governance/state-governance.md`

```markdown
# State Governance

## Scenario ID

state-governance

## Source / Provenance

- Source type: application seed
- Source path: `protocols/wire-lab.d/TODO/TODO-dadub-root-scenario-skeleton-and-seed-catalog.md`
- Source row/title: Seed application catalog entry `state-governance`
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-midif`; `TODO-dadub`

## Purpose

Exercise PromiseGrid design candidates against state-governance application pressure:
Cross-agency coordination, benefits, licensing, legislative records, and regional
authority.

## Setup

Alice depends on an outcome in the State Governance domain. Bob makes promises about
cross-agency coordination, benefits, licensing, legislative records, and regional
authority. Carol needs enough evidence to rely on Bob's promise without having complete
global state, and Mallory may exploit stale, missing, or ambiguous evidence.

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

### `scenarios/website-backend-hosting/website-backend-hosting.md`

```markdown
# Website Backend Hosting

## Scenario ID

website-backend-hosting

## Source / Provenance

- Source type: application seed
- Source path: `protocols/wire-lab.d/TODO/TODO-dadub-root-scenario-skeleton-and-seed-catalog.md`
- Source row/title: Seed application catalog entry `website-backend-hosting`
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-midif`; `TODO-dadub`

## Purpose

Exercise PromiseGrid design candidates against website-backend-hosting application
pressure: Deployments, service ownership, secret rotation, incident response, and uptime
promises.

## Setup

Alice depends on an outcome in the Website Backend Hosting domain. Bob makes promises
about deployments, service ownership, secret rotation, incident response, and uptime
promises. Carol needs enough evidence to rely on Bob's promise without having complete
global state, and Mallory may exploit stale, missing, or ambiguous evidence.

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

### `simulations/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/README.md`

```markdown
# SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid: Grid-envelope variant

This simulation is a standalone positional grid-envelope specimen. It tests the
combination `enc-dag-cbor`, `unknown-hard-reject`, and `sig-wrapper-pcid` without claiming
that this combination is the canonical PromiseGrid wire format. Source: `DI-fanah`.

The local draft spec is
`protocols/grid-envelope.d/specs/grid-envelope-draft.md`.
```

### `simulations/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/QUESTION.md`

```markdown
# Question

Does a positional grid envelope using `enc-dag-cbor`, `unknown-hard-reject`, and
`sig-wrapper-pcid` satisfy the wire-lab harness scenarios better than the sibling
variants? Source: `DI-fanah`.
```

### `simulations/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/protocols/grid-envelope.d/CHANGELOG.md`

```markdown
# CHANGELOG: grid-envelope

A-side CHANGELOG (per TE-liviv) for this simulation-local `grid-envelope`
protocol specimen.

This file records freeze events authored by the specimen maintainers. No entries
yet; this protocol specimen has not reached a first freeze.

This protocol tree is a simulation-local specimen created by `DI-fanah`.
```

### `simulations/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/protocols/grid-envelope.d/specs/grid-envelope-draft.md`

```markdown
# Grid Envelope Variant Spec (DRAFT)

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `enc-dag-cbor` / `unknown-hard-reject` / `sig-wrapper-pcid`.
> Source: `DI-fanah`.

## Purpose

This spec defines one full positional grid-envelope candidate for wire-lab
comparison. It is a specimen inside `SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid`, not a harness rule and not the
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

If a receiver lacks a handler for `pcid`, the envelope is rejected at the envelope layer. The receiver may keep local diagnostics, but it MUST NOT accept, store, or forward the message as a valid grid-envelope message under this variant.

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

## Existing Fitness Evidence From This Run

### `results/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/insurance-claims/openai-gpt-5.4-xhigh/20260520-184540.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload-insurance-claims-openai-gpt-5.4-xhigh-20260520-184540",
  "run_group_id": "ga-canary-20260520-184540",
  "cell_id": "ga-canary-20260520-184540-000004-SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload--insurance-claims--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload",
  "scenario_id": "insurance-claims",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-184540",
  "result_path": "results/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/insurance-claims/openai-gpt-5.4-xhigh/20260520-184540.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_28b95ad70b1b4baab4f93d3b3b676c13",
    "response_id": "resp_016df51201e5e546006a0e01b978448198b6246df9fac146db",
    "input_tokens": 6175,
    "output_tokens": 8389,
    "cost_usd": 0.12825225
  },
  "source": {
    "repo_commit": "bf8f3ec6b37af889c7f7b101cdf6c7b6ee00edaa",
    "sim_path": "simulations/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/",
    "scenario_path": "scenarios/insurance-claims/insurance-claims.md",
    "root_contract_paths": [
      "results/RUN-PROTOCOL.md",
      "scenarios/README.md"
    ],
    "files": [
      {
        "path": "results/RUN-PROTOCOL.md",
        "sha256": "577f1f424926d38817eb0dbface9c3cc1a1d70378a07629c848eb0ecd8351bc6"
      },
      {
        "path": "scenarios/README.md",
        "sha256": "406c4c7f400df14788d1caea61406f83adbc474160806ce4f1aa6a88d409d483"
      },
      {
        "path": "simulations/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/README.md",
        "sha256": "85676851ce384edf8fc4546252b69b524b2c57bc5d7a72e603dba3e63c24c90c"
      },
      {
        "path": "simulations/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/QUESTION.md",
        "sha256": "50fa82530ea929fde8df7856a9f2fa67758e32dfc2ff081c7f9a23d231395ace"
      },
      {
        "path": "simulations/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/protocols/grid-envelope.d/CHANGELOG.md",
        "sha256": "aef2b7915320af9607f415cef649397530561e4f0f9372fda1bbf2691dd3f9b1"
      },
      {
        "path": "simulations/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/protocols/grid-envelope.d/specs/grid-envelope-draft.md",
        "sha256": "5a83134788f5fd6bebb92119be82086f34baa3756ec72971fc8da65456cce1b2"
      },
      {
        "path": "scenarios/insurance-claims/insurance-claims.md",
        "sha256": "2a9fe6a4c9a84a35f032d93ae9c36d2adec30a8175f821a34cf8592c73810ab1"
      }
    ],
    "simulation_tree_hash": "70c6a17b236d0e9f1ba89621ed8e390e4169976db54f9ebcc2337ee9d9646834"
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
    "evolution_safety": 3,
    "layer_boundary_clarity": 4,
    "failure_handling": 2,
    "implementation_plausibility": 4,
    "risk_penalty": 3
  },
  "fitness": {
    "raw": 23,
    "normalized_0_100": 57.5,
    "confidence_0_1": 0.76
  },
  "assessment": {
    "rationale": "This variant is a credible low-level carrier for contested claim records because it gives deterministic CBOR bytes, explicit pCID-based dispatch, and mandatory signature slots, but it does not itself model claim evidence types, adjuster delegation, payment or appeal promises, or peer-local promise accounting. Its best contribution to the insurance-claims scenario is visible failure at the interpretation boundary and long-term salvage of unknown payloads; its main gap is that the application semantics and authority chain still have to live elsewhere.",
    "strengths": [
      "Deterministic CBOR arrays provide stable canonical bytes for hashing and signing, which helps later audit and replay of claim artifacts.",
      "Mandatory signature slots reduce the chance that important claim envelopes are treated as unsigned by default and make provenance expectations explicit.",
      "pCID-based dispatch and the unknown-pCID rule keep unsupported interpretations visible instead of silently guessing semantics.",
      "The spec is clear that ordering, forwarding, and external body references belong to payload protocols or wrapper envelopes, which helps preserve layer boundaries."
    ],
    "weaknesses": [
      "The specimen does not define insurance-domain payloads for claim evidence, adjuster authority, fraud review, payment promises, or appeals.",
      "Peer-local promise accounting, names, feeds, and identity/delegation chains are outside this layer, so the scenario's core decision logic remains underspecified.",
      "Auditability is limited by positional, opaque byte slots; humans and future tools still need the referenced pCID handlers or salvage tooling.",
      "The variant is still draft and unfrozen, with no minted stable spec pCID or freeze history."
    ],
    "risks": [
      "Best-effort inspection of unknown payloads could be overtrusted in a disputed claim, even though the spec says such output must remain unsupported and unverified.",
      "Signature presence is mandatory, but signature meaning and signer authority depend on external sig_pcid and identity protocols that may diverge or age poorly.",
      "Because the specimen is not frozen, long-lived records using it may face interoperability or archival interpretation drift."
    ],
    "open_questions": [
      "Which payload protocols would carry claim evidence bundles, adjuster delegation, payment commitments, and appeal records in a way peers can audit locally?",
      "How are signer identity, delegation, revocation, and time-bounded adjuster authority represented without hidden central registries?",
      "Should insurance-grade use forbid or further constrain best-effort salvage so unsupported payload views cannot be confused with admissible evidence?",
      "Do contested-claims workflows need additional wrapper conventions for timestamps, receipts, countersignatures, or relay provenance beyond this base envelope?",
      "How should peers retry, downgrade, or escalate when referenced CAS objects or signature handlers are missing for long periods?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

### `results/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/state-governance/openai-gpt-5.4-xhigh/20260520-184540.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload-state-governance-openai-gpt-5.4-xhigh-20260520-184540",
  "run_group_id": "ga-canary-20260520-184540",
  "cell_id": "ga-canary-20260520-184540-000005-SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload--state-governance--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload",
  "scenario_id": "state-governance",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-184540",
  "result_path": "results/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/state-governance/openai-gpt-5.4-xhigh/20260520-184540.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_2e2f74642df249a5acdce2d9c4f1716a",
    "response_id": "resp_09b0bfc3e995b041006a0e01c3bd488198bebfde7344d36486",
    "input_tokens": 6181,
    "output_tokens": 8794,
    "cost_usd": 0.13393275
  },
  "source": {
    "repo_commit": "bf8f3ec6b37af889c7f7b101cdf6c7b6ee00edaa",
    "sim_path": "simulations/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/",
    "scenario_path": "scenarios/state-governance/state-governance.md",
    "root_contract_paths": [
      "results/RUN-PROTOCOL.md",
      "scenarios/README.md"
    ],
    "files": [
      {
        "path": "results/RUN-PROTOCOL.md",
        "sha256": "577f1f424926d38817eb0dbface9c3cc1a1d70378a07629c848eb0ecd8351bc6"
      },
      {
        "path": "scenarios/README.md",
        "sha256": "406c4c7f400df14788d1caea61406f83adbc474160806ce4f1aa6a88d409d483"
      },
      {
        "path": "simulations/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/README.md",
        "sha256": "85676851ce384edf8fc4546252b69b524b2c57bc5d7a72e603dba3e63c24c90c"
      },
      {
        "path": "simulations/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/QUESTION.md",
        "sha256": "50fa82530ea929fde8df7856a9f2fa67758e32dfc2ff081c7f9a23d231395ace"
      },
      {
        "path": "simulations/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/protocols/grid-envelope.d/CHANGELOG.md",
        "sha256": "aef2b7915320af9607f415cef649397530561e4f0f9372fda1bbf2691dd3f9b1"
      },
      {
        "path": "simulations/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/protocols/grid-envelope.d/specs/grid-envelope-draft.md",
        "sha256": "5a83134788f5fd6bebb92119be82086f34baa3756ec72971fc8da65456cce1b2"
      },
      {
        "path": "scenarios/state-governance/state-governance.md",
        "sha256": "6cbfa949fc3eed42fd81bb6820d73ec96ac9c71669b3117d6e46673ddba38360"
      }
    ],
    "simulation_tree_hash": "70c6a17b236d0e9f1ba89621ed8e390e4169976db54f9ebcc2337ee9d9646834"
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
    "evolution_safety": 3,
    "layer_boundary_clarity": 5,
    "failure_handling": 2,
    "implementation_plausibility": 4,
    "risk_penalty": 3
  },
  "fitness": {
    "raw": 24,
    "normalized_0_100": 60,
    "confidence_0_1": 0.78
  },
  "assessment": {
    "rationale": "This draft envelope is a credible wire-level substrate for state-governance evidence because it gives deterministic CBOR bytes, explicit payload and signature protocol dispatch, mandatory signature slots, and clear layer boundaries without assuming a central registry. Unknown-pCID best-effort salvage is useful for cross-agency partial interoperability, but the specimen stops at transport semantics: it does not define the identity, delegation, naming, feed, revocation, multi-authority, or local decision rules needed when records are delayed, stale, or disputed.",
    "strengths": [
      "Deterministic CBOR encoding and mandatory signature metadata support stable hashing, verification, and later audit of exactly what bytes were signed.",
      "The spec is explicit about envelope vs payload vs signature responsibilities, which is strong for layer-boundary clarity under sparse knowledge.",
      "Using pCID and sig_pCID avoids a required central registry and fits peer-local verification better than globally coordinated dispatch.",
      "The four-slot positional structure is straightforward to implement and compare against sibling wire variants."
    ],
    "weaknesses": [
      "The specimen is only an envelope draft; it does not define the state-governance payloads, authority records, feeds, or promise-accounting artifacts the scenario asks peers to reason about.",
      "Single signature slots are not enough by themselves to express countersignatures, delegation chains, regional overrides, or revocation-heavy public-sector workflows.",
      "Failure handling stops at visible unsupported or unverified status; it does not specify acceptance, retry, downgrade, escalation, or stale-evidence policy for contested cases.",
      "The spec is still unfrozen and has no minted specimen pCID or freeze history, so long-term durability is not yet demonstrated."
    ],
    "risks": [
      "Operators may over-trust best-effort inspection of unknown payloads even though the spec says such results are unsupported and unverified.",
      "Cross-agency interoperability can fragment if different agencies mint incompatible payload or signature protocol identifiers without clear migration guidance.",
      "Long-lived governance records need time, revocation, and authority-change semantics that are outside this envelope and could be handled inconsistently downstream."
    ],
    "open_questions": [
      "What payload protocol carries agency identity, jurisdiction, delegation, and promise-accounting records for this scenario?",
      "How are multiple signatures or chained attestations represented in practice: nested envelopes, a richer signature protocol, or payload-level structures?",
      "What local policy tells Alice or Carol when unsupported best-effort inspection may inform triage but must never justify acceptance?",
      "How would this variant handle decades-long migration of keys, pCIDs, and legal authority changes while preserving auditability?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

### `results/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/website-backend-hosting/openai-gpt-5.4-xhigh/20260520-184540.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload-website-backend-hosting-openai-gpt-5.4-xhigh-20260520-184540",
  "run_group_id": "ga-canary-20260520-184540",
  "cell_id": "ga-canary-20260520-184540-000006-SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload--website-backend-hosting--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload",
  "scenario_id": "website-backend-hosting",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-184540",
  "result_path": "results/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/website-backend-hosting/openai-gpt-5.4-xhigh/20260520-184540.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_33b128bd6fbc488aae7fb5bf277d4795",
    "response_id": "resp_0c9308ad2c6848ca006a0e01de5640819a873a536433e85332",
    "input_tokens": 6196,
    "output_tokens": 10459,
    "cost_usd": 0.157269
  },
  "source": {
    "repo_commit": "bf8f3ec6b37af889c7f7b101cdf6c7b6ee00edaa",
    "sim_path": "simulations/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/",
    "scenario_path": "scenarios/website-backend-hosting/website-backend-hosting.md",
    "root_contract_paths": [
      "results/RUN-PROTOCOL.md",
      "scenarios/README.md"
    ],
    "files": [
      {
        "path": "results/RUN-PROTOCOL.md",
        "sha256": "577f1f424926d38817eb0dbface9c3cc1a1d70378a07629c848eb0ecd8351bc6"
      },
      {
        "path": "scenarios/README.md",
        "sha256": "406c4c7f400df14788d1caea61406f83adbc474160806ce4f1aa6a88d409d483"
      },
      {
        "path": "simulations/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/README.md",
        "sha256": "85676851ce384edf8fc4546252b69b524b2c57bc5d7a72e603dba3e63c24c90c"
      },
      {
        "path": "simulations/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/QUESTION.md",
        "sha256": "50fa82530ea929fde8df7856a9f2fa67758e32dfc2ff081c7f9a23d231395ace"
      },
      {
        "path": "simulations/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/protocols/grid-envelope.d/CHANGELOG.md",
        "sha256": "aef2b7915320af9607f415cef649397530561e4f0f9372fda1bbf2691dd3f9b1"
      },
      {
        "path": "simulations/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/protocols/grid-envelope.d/specs/grid-envelope-draft.md",
        "sha256": "5a83134788f5fd6bebb92119be82086f34baa3756ec72971fc8da65456cce1b2"
      },
      {
        "path": "scenarios/website-backend-hosting/website-backend-hosting.md",
        "sha256": "410d1a5fce15529c59260d9a9b20479dcb0eacb406ab9c480d2a6279f73fe18f"
      }
    ],
    "simulation_tree_hash": "70c6a17b236d0e9f1ba89621ed8e390e4169976db54f9ebcc2337ee9d9646834"
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
    "promisegrid_alignment": 4,
    "auditability": 3,
    "evolution_safety": 3,
    "layer_boundary_clarity": 4,
    "failure_handling": 2,
    "implementation_plausibility": 4,
    "risk_penalty": 3
  },
  "fitness": {
    "raw": 19,
    "normalized_0_100": 54,
    "confidence_0_1": 0.78
  },
  "assessment": {
    "rationale": "This variant is a competent low-level evidence envelope for contested website-hosting events: Alice and Carol get deterministic CBOR bytes, explicit payload/signature dispatch, and a clear unsupported/unverified path when payload handlers are missing, all without assuming a central registry. However, it remains only a substrate for the scenario. Deployment, service-ownership, secret-rotation, incident-response, and uptime semantics; freshness and revocation rules; peer-local promise accounting; and explicit unknown-signature-protocol handling are left outside the envelope, so the scenario fit is only partial and operational risk remains moderate.",
    "strengths": [
      "Deterministic CBOR plus signing over [pcid, payload] gives stable bytes for local audit of contested records.",
      "Separate pcid and sig_pcid support sparse-knowledge dispatch and independent evolution of payload and signature protocols.",
      "Unknown payload handling requires unsupported/unverified labeling instead of silent reinterpretation.",
      "The spec is unusually explicit about layer boundaries: payload protocols, wrappers, and the envelope each have distinct responsibilities."
    ],
    "weaknesses": [
      "The draft does not define the website-hosting payload schemas, identity claims, names, feeds, or promise-accounting records the scenario actually needs.",
      "It has no built-in freshness, expiry, revocation, or anti-replay semantics for stale but signed deployment or ownership evidence.",
      "Receiver behavior for unknown sig_pcid is not specified as clearly as unknown payload pcid behavior.",
      "Positional opaque byte fields are less human-auditable without resolver tooling and payload/spec access.",
      "The specimen is still draft, unfrozen, and its own spec pCID is not yet minted, which weakens current 100-year durability claims."
    ],
    "risks": [
      "Operators under incident pressure may over-trust best-effort inspection of unknown payloads despite the required unsupported/unverified label.",
      "Mandatory signature verification can become an availability problem during key rotation or signature-protocol migration if sig_pcid handlers lag.",
      "Replay or stale-record risk remains unless payload protocols add explicit freshness and revocation evidence."
    ],
    "open_questions": [
      "Which payload pCIDs and CAS objects carry deployment approvals, service ownership assertions, secret-rotation attestations, incident timelines, and uptime evidence?",
      "How do Alice and Carol determine currentness, revocation status, and authorized signer sets using only local evidence?",
      "What is the required downgrade path when sig_pcid is unknown or verification cannot be completed?",
      "For website-hosting operations, should hop-local incident or relay evidence live in wrapper envelopes, payload objects, or sig_payload?",
      "What peer-local promise-accounting record supports retry, downgrade, and escalation decisions during partition?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

### `results/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/insurance-claims/openai-gpt-5.4-xhigh/20260520-184540.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid-insurance-claims-openai-gpt-5.4-xhigh-20260520-184540",
  "run_group_id": "ga-canary-20260520-184540",
  "cell_id": "ga-canary-20260520-184540-000007-SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid--insurance-claims--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid",
  "scenario_id": "insurance-claims",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-184540",
  "result_path": "results/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/insurance-claims/openai-gpt-5.4-xhigh/20260520-184540.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_8ce1565797414a4ab08e5add79d306c9",
    "response_id": "resp_0daf72cf0ce7a0a3006a0e023afdf8819ab7882f2b0f617a7f",
    "input_tokens": 6088,
    "output_tokens": 10472,
    "cost_usd": 0.157262
  },
  "source": {
    "repo_commit": "bf8f3ec6b37af889c7f7b101cdf6c7b6ee00edaa",
    "sim_path": "simulations/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/",
    "scenario_path": "scenarios/insurance-claims/insurance-claims.md",
    "root_contract_paths": [
      "results/RUN-PROTOCOL.md",
      "scenarios/README.md"
    ],
    "files": [
      {
        "path": "results/RUN-PROTOCOL.md",
        "sha256": "577f1f424926d38817eb0dbface9c3cc1a1d70378a07629c848eb0ecd8351bc6"
      },
      {
        "path": "scenarios/README.md",
        "sha256": "406c4c7f400df14788d1caea61406f83adbc474160806ce4f1aa6a88d409d483"
      },
      {
        "path": "simulations/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/README.md",
        "sha256": "37ecc4b11d4fb7bdc0436801830411c65c739e69957ee1c0ac72fec86ec20d7b"
      },
      {
        "path": "simulations/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/QUESTION.md",
        "sha256": "ce6ab563ed7d413b3f2027ee8863184e9c0f28e59a1d2c074ab124a7f4310754"
      },
      {
        "path": "simulations/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/protocols/grid-envelope.d/CHANGELOG.md",
        "sha256": "aef2b7915320af9607f415cef649397530561e4f0f9372fda1bbf2691dd3f9b1"
      },
      {
        "path": "simulations/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/protocols/grid-envelope.d/specs/grid-envelope-draft.md",
        "sha256": "65b435945ee4d952230dec56867d56bb533ec3e533452bcaa54861949a9b3ba8"
      },
      {
        "path": "scenarios/insurance-claims/insurance-claims.md",
        "sha256": "2a9fe6a4c9a84a35f032d93ae9c36d2adec30a8175f821a34cf8592c73810ab1"
      }
    ],
    "simulation_tree_hash": "90a6dd3091fcccc5e5fef0a39e859ae00100ea2a89b212ac10fc2fc2aca1d42e"
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
    "promisegrid_alignment": 2,
    "auditability": 3,
    "evolution_safety": 1,
    "layer_boundary_clarity": 5,
    "failure_handling": 2,
    "implementation_plausibility": 4,
    "risk_penalty": 3
  },
  "fitness": {
    "raw": 21,
    "normalized_0_100": 52.5,
    "confidence_0_1": 0.77
  },
  "assessment": {
    "rationale": "Good wire-layer specimen for exact bytes, dispatch, and visible failure boundaries, but only a partial fit to insurance-claims application pressure. The unknown-hard-reject rule protects against silent misinterpretation, yet it weakens 100-year and sparse-knowledge handling because disputed claim evidence may need to be preserved, relayed, and later reinterpreted before every peer has the relevant payload or signature-wrapper pCID.",
    "strengths": [
      "Very clear [pcid, payload] boundary keeps envelope dispatch separate from claim-domain semantics.",
      "Canonical DAG-CBOR bytes support reproducible hashing/signing and stable local evidence records.",
      "Explicit reject-on-unknown behavior avoids silent misinterpretation of claim data.",
      "Wrapper-selected signature protocols allow authorship and integrity schemes to evolve without changing the base envelope.",
      "The specimen does not require a central registry in the envelope shape itself; handler knowledge remains peer-local."
    ],
    "weaknesses": [
      "The specimen does not define claim evidence, adjuster delegation, payment, appeal, or promise-accounting structures.",
      "Unknown pCID hard reject prevents treating unrecognized envelopes as valid first-class messages, which is brittle for delayed or cross-version evidence.",
      "Authorship and integrity audit depend on understanding external signature-wrapper pCIDs rather than a fixed envelope-level signature slot.",
      "The spec is still draft and unfrozen, with no minted pCID or mature freeze history."
    ],
    "risks": [
      "Version skew can make otherwise important claim evidence operationally unusable at the moment a dispute needs preservation.",
      "Mallory could exploit pCID or wrapper mismatch to cause selective evidence blindness or denial of service.",
      "Operators may compensate with central translation gateways or exceptions, weakening PromiseGrid's no-central-authority goals."
    ],
    "open_questions": [
      "Can unknown envelopes be safely quarantined or archived for later review without violating the hard-reject rule?",
      "What common payload and signature-wrapper set would be required for claim intake, adjuster authority, fraud review, payment, and appeal?",
      "How are identity claims and adjuster delegation bound across nested wrappers in a way Carol can audit locally?",
      "Would a preserve-but-not-interpret policy outperform hard reject for long-lived claims evidence?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

### `results/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/state-governance/openai-gpt-5.4-xhigh/20260520-184540.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid-state-governance-openai-gpt-5.4-xhigh-20260520-184540",
  "run_group_id": "ga-canary-20260520-184540",
  "cell_id": "ga-canary-20260520-184540-000008-SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid--state-governance--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid",
  "scenario_id": "state-governance",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-184540",
  "result_path": "results/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/state-governance/openai-gpt-5.4-xhigh/20260520-184540.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_ef969e05b9ac4fe4bae1734f956616b4",
    "response_id": "resp_04ae555d3b2916f3006a0e0242bd50819bb3974ff3b1c98eea",
    "input_tokens": 6094,
    "output_tokens": 11085,
    "cost_usd": 0.1658545
  },
  "source": {
    "repo_commit": "bf8f3ec6b37af889c7f7b101cdf6c7b6ee00edaa",
    "sim_path": "simulations/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/",
    "scenario_path": "scenarios/state-governance/state-governance.md",
    "root_contract_paths": [
      "results/RUN-PROTOCOL.md",
      "scenarios/README.md"
    ],
    "files": [
      {
        "path": "results/RUN-PROTOCOL.md",
        "sha256": "577f1f424926d38817eb0dbface9c3cc1a1d70378a07629c848eb0ecd8351bc6"
      },
      {
        "path": "scenarios/README.md",
        "sha256": "406c4c7f400df14788d1caea61406f83adbc474160806ce4f1aa6a88d409d483"
      },
      {
        "path": "simulations/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/README.md",
        "sha256": "37ecc4b11d4fb7bdc0436801830411c65c739e69957ee1c0ac72fec86ec20d7b"
      },
      {
        "path": "simulations/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/QUESTION.md",
        "sha256": "ce6ab563ed7d413b3f2027ee8863184e9c0f28e59a1d2c074ab124a7f4310754"
      },
      {
        "path": "simulations/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/protocols/grid-envelope.d/CHANGELOG.md",
        "sha256": "aef2b7915320af9607f415cef649397530561e4f0f9372fda1bbf2691dd3f9b1"
      },
      {
        "path": "simulations/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/protocols/grid-envelope.d/specs/grid-envelope-draft.md",
        "sha256": "65b435945ee4d952230dec56867d56bb533ec3e533452bcaa54861949a9b3ba8"
      },
      {
        "path": "scenarios/state-governance/state-governance.md",
        "sha256": "6cbfa949fc3eed42fd81bb6820d73ec96ac9c71669b3117d6e46673ddba38360"
      }
    ],
    "simulation_tree_hash": "90a6dd3091fcccc5e5fef0a39e859ae00100ea2a89b212ac10fc2fc2aca1d42e"
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
    "evolution_safety": 1,
    "layer_boundary_clarity": 4,
    "failure_handling": 2,
    "implementation_plausibility": 4,
    "risk_penalty": 4
  },
  "fitness": {
    "raw": 20,
    "normalized_0_100": 50,
    "confidence_0_1": 0.78
  },
  "assessment": {
    "rationale": "Useful as a lower-layer specimen: the positional DAG-CBOR envelope gives clear canonical bytes and explicit dispatch, and wrapper-selected signature protocols keep layer boundaries clean. But it is only a wire-format slice of the state-governance problem, not a full account of disputed records, names, delegations, or promise accounting. The decisive drawback is the unknown-hard-reject rule, which makes cross-agency upgrades and multi-decade evolution brittle.",
    "strengths": [
      "Minimal [pcid, payload] structure is straightforward to implement and gives a stable byte sequence for hashing and signing.",
      "Putting signatures and hop/authorship evidence in pCID-selected wrappers preserves a clean boundary between envelope transport and higher-layer evidence semantics.",
      "Peers can make local accept/reject decisions from the explicit pCID boundary, and unsupported handlers fail visibly rather than being silently misinterpreted."
    ],
    "weaknesses": [
      "The specimen does not define the application objects, names, identity claims, feeds, or promise-accounting records that the state-governance scenario needs.",
      "Unknown-hard-reject prevents graceful archive, relay, or deferred interpretation when agencies upgrade at different times.",
      "Audit and security depend on external wrapper conventions that are not specified here, especially for signer binding, delegation, and hop evidence."
    ],
    "risks": [
      "Version skew across agencies can create hard interoperability failures as soon as a new wrapper or payload pCID appears.",
      "Long-lived governance records may be dropped or stranded because older peers cannot treat unknown envelopes as valid messages for retention or forwarding.",
      "Divergent signature-wrapper conventions could create authorship disputes, stripping/rewrapping ambiguity, or uneven verifier behavior."
    ],
    "open_questions": [
      "Can a peer safely archive or relay an unknown envelope as opaque bytes without violating the hard-reject contract?",
      "How are pCID definitions discovered, authenticated, cached, and migrated across agencies over decades without sliding into de facto central governance?",
      "What exact wrapper rules prevent signature stripping or wrapper-order ambiguity in contested cases?",
      "Where do agency names, regional authority claims, and delegation chains live so Carol can audit a disputed decision from local evidence alone?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

### `results/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/website-backend-hosting/openai-gpt-5.4-xhigh/20260520-184540.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid-website-backend-hosting-openai-gpt-5.4-xhigh-20260520-184540",
  "run_group_id": "ga-canary-20260520-184540",
  "cell_id": "ga-canary-20260520-184540-000009-SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid--website-backend-hosting--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid",
  "scenario_id": "website-backend-hosting",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-184540",
  "result_path": "results/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/website-backend-hosting/openai-gpt-5.4-xhigh/20260520-184540.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_47bd9cca2b5f4f27b885c1c173c6e25c",
    "response_id": "resp_0882e56c8adcd440006a0e0289bef881999c5ff3406b6d2fb7",
    "input_tokens": 6109,
    "output_tokens": 10597,
    "cost_usd": 0.15904875
  },
  "source": {
    "repo_commit": "bf8f3ec6b37af889c7f7b101cdf6c7b6ee00edaa",
    "sim_path": "simulations/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/",
    "scenario_path": "scenarios/website-backend-hosting/website-backend-hosting.md",
    "root_contract_paths": [
      "results/RUN-PROTOCOL.md",
      "scenarios/README.md"
    ],
    "files": [
      {
        "path": "results/RUN-PROTOCOL.md",
        "sha256": "577f1f424926d38817eb0dbface9c3cc1a1d70378a07629c848eb0ecd8351bc6"
      },
      {
        "path": "scenarios/README.md",
        "sha256": "406c4c7f400df14788d1caea61406f83adbc474160806ce4f1aa6a88d409d483"
      },
      {
        "path": "simulations/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/README.md",
        "sha256": "37ecc4b11d4fb7bdc0436801830411c65c739e69957ee1c0ac72fec86ec20d7b"
      },
      {
        "path": "simulations/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/QUESTION.md",
        "sha256": "ce6ab563ed7d413b3f2027ee8863184e9c0f28e59a1d2c074ab124a7f4310754"
      },
      {
        "path": "simulations/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/protocols/grid-envelope.d/CHANGELOG.md",
        "sha256": "aef2b7915320af9607f415cef649397530561e4f0f9372fda1bbf2691dd3f9b1"
      },
      {
        "path": "simulations/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/protocols/grid-envelope.d/specs/grid-envelope-draft.md",
        "sha256": "65b435945ee4d952230dec56867d56bb533ec3e533452bcaa54861949a9b3ba8"
      },
      {
        "path": "scenarios/website-backend-hosting/website-backend-hosting.md",
        "sha256": "410d1a5fce15529c59260d9a9b20479dcb0eacb406ab9c480d2a6279f73fe18f"
      }
    ],
    "simulation_tree_hash": "90a6dd3091fcccc5e5fef0a39e859ae00100ea2a89b212ac10fc2fc2aca1d42e"
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
    "evolution_safety": 1,
    "layer_boundary_clarity": 5,
    "failure_handling": 2,
    "implementation_plausibility": 4,
    "risk_penalty": 4
  },
  "fitness": {
    "raw": 14,
    "normalized_0_100": 45,
    "confidence_0_1": 0.77
  },
  "assessment": {
    "rationale": "The specimen is strong as a minimal wire-layer envelope: [pcid, payload], canonical DAG-CBOR bytes, and wrapper-selected signature protocols make dispatch and hashing rules explicit. But for website-backend-hosting it only supplies a transport shell, not the deployment, ownership, secret-rotation, incident-response, or uptime records the scenario needs. Its unknown-pCID hard-reject rule gives visible failure, yet it is brittle for rolling upgrades, partial knowledge, and outage-time evidence preservation.",
    "strengths": [
      "Very clear separation between envelope dispatch and payload semantics.",
      "Canonical DAG-CBOR bytes support reproducible hashing and signature verification.",
      "Avoids assuming a central pCID registry and allows signatures or other evidence to live in wrapper protocols.",
      "Incompatible interpretation fails at an explicit dispatch boundary instead of silently misparsing bytes."
    ],
    "weaknesses": [
      "Does not define the application-level promises, feeds, or CAS objects needed for backend hosting operations.",
      "Base envelope has no built-in authorship or relay-evidence slot, so auditors must interpret extra wrapper protocols.",
      "Unknown pCID messages are hard rejected rather than gracefully downgraded, relayed, or preserved as valid envelope traffic.",
      "The spec is still a draft specimen, so long-term interoperability and freeze discipline are unproven."
    ],
    "risks": [
      "Rolling upgrades or key-rotation changes could surface as availability failures when some peers lack a new pCID handler.",
      "Incident responders may lose useful in-band evidence if unknown envelopes cannot be carried forward as valid messages.",
      "Mallory can exploit unknown-pCID rejection or wrapper inconsistency to trigger drops or dispute handling during contested events."
    ],
    "open_questions": [
      "What payload protocols represent deployments, service ownership, secret rotation, incident response, and uptime promises?",
      "Can peers preserve unknown envelopes as raw audit artifacts without undermining the hard-reject safety goal?",
      "How are pCIDs discovered, cached, and migrated over long time horizons without a central registry?",
      "What wrapper conventions provide consistent authorship, hop evidence, encryption, and key-rotation semantics across services?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

## Required JSON Shape

{"child_id":"SIM-fukis-ga-child-0002","design_delta_summary":"one to three bounded design deltas","files":[{"path":"README.md","content":"# ..."},{"path":"QUESTION.md","content":"# ..."}]}
