# Matrix Cell API Evaluation

Return only the complete Markdown result file. Do not wrap it in a code fence. Do not ask for confirmation.

## Required Output Coordinates

- Result path: `results/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/conditional-release-geofencing-opaque-lower-layer-carriage/openai-gpt-5.3-codex-xhigh/20260519-054830.md`
- Header: `# Result: SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid / conditional-release-geofencing-opaque-lower-layer-carriage / openai-gpt-5.3-codex-xhigh / 20260519-054830`
- Result ID: `SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid-conditional-release-geofencing-opaque-lower-layer-carriage-openai-gpt-5.3-codex-xhigh-20260519-054830`
- Scenario ID: `conditional-release-geofencing-opaque-lower-layer-carriage`
- Scenario path: `scenarios/conditional-release-geofencing-opaque-lower-layer-carriage/conditional-release-geofencing-opaque-lower-layer-carriage.md`
- Simulation ID: `SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid`
- Simulation path: `simulations/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/`
- Simulation commit: `f5cd6086ed89e653eb439a8950e218aa708dd013`
- Model ID: `openai-gpt-5.3-codex-xhigh`
- Run timestamp UTC: `20260519-054830`

## Required Result Contract

The Markdown result must contain exactly the required sections from `results/RUN-PROTOCOL.md`, include `- Run mode: llm-doc-eval-blind`, include one line starting with `Evidence verdict:`, and include an explicit `Authority Boundary` section.

## Evaluation Task

Evaluate the simulation against the scenario using deeper reasoning. Explain what the simulation covers, what it pushes to another layer, how the scenario's 100-year durability, sparse knowledge, no central authority, auditability, and migration pressures expose weaknesses, and which open questions remain.

## Source Documents

### `results/RUN-PROTOCOL.md`

```markdown
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

Full-matrix execution may be unattended: the queue runner may invoke one
external LLM command per manifest row, checkpoint state after every cell,
and validate the result file. The external runner must still produce the
substantive result file; queue tooling only coordinates work and validation.
Source: `DI-nuhon`; `DI-zamin`.

The preferred runner is now the Go `tools/matrix-runner` CLI. For API-backed
runs it bundles local source document contents into the provider prompt, because
remote APIs cannot read repo-local paths. Source: `DI-lulom`.

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
- Use manifests with concrete `timestamp`, `result_path`, `ordinal`, and
  `cell_id` fields for unattended runs.
- Use checkpoint state under `results/state/` for any long run that should
  resume without operator prompts.
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

- Preferred Go runner:
  `cd tools/matrix-runner && go run . <subcommand>`
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

## Recommended Preflight

1. Generate a canary manifest (20-30 cells).
2. Generate canary LLM job prompts.
3. Run canary prompts through an LLM evaluator.
4. Validate canary output.
5. Review drift/comparison report.
6. Run the full manifest.

## Unattended Full-Run Shape

1. Generate a manifest with a fixed model ID and concrete timestamp:
   `cd tools/matrix-runner && go run . manifest -repo-root ../.. -models <model-id>`.
2. Start the queue with OpenAI API-backed execution:
   `cd tools/matrix-runner && go run . run -repo-root ../.. -manifest <manifest.csv> -provider openai -api-model <api-model> -reasoning-effort xhigh`.
3. Let the queue process one cell at a time. Each cell writes or refreshes its
   prompt under `results/jobs/<run-group-id>/`, invokes the runner command,
   validates `result_path`, and checkpoints `results/state/<run-group-id>.json`.
   Source: `DI-nuhon`; `DI-lulom`; `DI-zamin`.
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
world state. Source: `DI-faros`; `DI-vabor`; `DI-dimas`.

## Directory Shape

Each root scenario entry uses this shape:

```text
scenarios/
  <entry-id>/
    README.md
    <scenario-id>.md
```

- `<entry-id>` is stable kebab-case.
- Application entries use one `scenarios/<application-id>/` directory per
  application, such as `scenarios/bgp-routing/`.
- Mined simulation rows use one root scenario entry per source row, transformed
  for cross-simulation comparison and linked back to the source
  `simulations/.../SCENARIOS.md` row.
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

## Applies To

<Which simulations, protocol specimens, or candidate design families should be
compared against this scenario.>

## Actors

Use Alice, Bob, Carol, Dave, Ellen, Frank, and Mallory where named actors help
make promises, failures, or adversarial pressure concrete.

## Setup

<Initial state, promises, artifacts, sites, policies, and sparse knowledge.>

## Stimulus

<The event or request that exercises the scenario.>

## Expected Pressure

<What the scenario should force a candidate design to explain.>

## Overarching Goal Checks

- 100-year durability:
- Sparse, partial knowledge:
- No central authority or registry:
- Peer-local trust / promise accounting:
- Adversarial or failure pressure:
- Human and LLM auditability:
- Migration / evolution path:

## Evaluation Questions

- <Question 1>
- <Question 2>
- <Question 3>

## Result Runs

Result runs live under:

`results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md`

## Authority Boundary

This scenario and its result runs are evidence only. Design authority still
graduates through DR, DI, frozen spec, or PromiseGrid Development Guide handoff.
```

## Scenario Quality Gates

Every root scenario must explicitly state how it exercises the 100-year
PromiseGrid goal and other overarching PromiseGrid goals. A scenario may be
small, but it must not be narrow in a way that accidentally assumes away the
conditions PromiseGrid is meant to survive. Source: `DI-botup`.

At minimum, each scenario should address:

- **100-year durability:** Does the scenario still make sense after tools,
  organizations, keys, people, and infrastructure have changed?
- **Sparse, partial knowledge:** Does it avoid assuming any peer has the whole
  graph, all CAS objects, or globally complete state?
- **No central authority or registry:** Does it avoid relying on a central pCID,
  identity, trust, naming, routing, currency, or governance authority unless the
  point of the scenario is to test that failure mode?
- **Peer-local trust / promise accounting:** Does it show what Alice, Bob, Carol,
  or another peer can observe and record locally?
- **Adversarial or failure pressure:** Does it include Mallory, corruption,
  refusal, stale data, partition, capture, default, or another failure mode when
  relevant?
- **Human and LLM auditability:** Can a later person or model understand what
  was promised, what happened, and why the result matters?
- **Migration / evolution path:** Does it expose what happens when protocols,
  keys, names, object shapes, policies, or organizations evolve?

If a gate is not relevant, the scenario should say why instead of omitting it.

## Generated Result Views

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

## Population Plan

`TODO-dadub` owns the initial population plan. The first mining pass created
one root entry per existing sim-local scenario row under `DI-nanih`, and the
first application-seed pass created one root entry per seed application under
`DI-midif`. Future entries should follow the same templates unless a later DI
changes the root scenario contract.
```

### `scenarios/conditional-release-geofencing-opaque-lower-layer-carriage/README.md`

```markdown
# Opaque lower-layer carriage

This root scenario entry was mined from `simulations/SIM-zarud-conditional-release-geofencing/SCENARIOS.md`. It is shared wire-lab
comparison apparatus, not PromiseGrid node layout or final design authority.
Source: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`.

## Scenario Files

- [`conditional-release-geofencing-opaque-lower-layer-carriage.md`](conditional-release-geofencing-opaque-lower-layer-carriage.md)

## Source

- Source simulation: `SIM-zarud-conditional-release-geofencing/`
- Source file: `simulations/SIM-zarud-conditional-release-geofencing/SCENARIOS.md`
- Source row title: Opaque lower-layer carriage
```

### `scenarios/conditional-release-geofencing-opaque-lower-layer-carriage/conditional-release-geofencing-opaque-lower-layer-carriage.md`

```markdown
# Opaque lower-layer carriage

## Scenario ID

conditional-release-geofencing-opaque-lower-layer-carriage

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-zarud-conditional-release-geofencing/SCENARIOS.md`
- Source simulation: `SIM-zarud-conditional-release-geofencing/`
- Source row/title: Opaque lower-layer carriage
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-zarud-conditional-release-geofencing/`.

## Applies To

- Primary source simulation: `SIM-zarud-conditional-release-geofencing/`
- Other simulations that claim to address `conditional-release-geofencing` or the same PromiseGrid
  pressure should record result runs under `results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md`; generated views can summarize completed runs when needed.

## Actors

Use Alice, Bob, Carol, Dave, Ellen, Frank, and Mallory as needed when running
this scenario. Preserve any actor roles implied by the source row instead of
inventing a central coordinator.

## Setup

Bob's node stores encrypted content whose condition vocabulary it cannot parse.

## Stimulus

Run the candidate simulation against this source test: Whether lower layers can safely carry opaque condition references while avoiding accidental promise violations.

## Expected Pressure

If opaque carriage is acceptable, the condition object must be verifiable without every layer understanding its semantics.

## Overarching Goal Checks

- 100-year durability: The result must explain whether this pressure remains
  understandable after tooling, organizations, keys, people, and infrastructure
  change.
- Sparse, partial knowledge: The run must not assume any peer has all CAS
  objects, all relationship records, or a complete global graph unless that
  assumption is explicitly the failure being tested.
- No central authority or registry: The run must avoid relying on a central pCID,
  identity, trust, naming, routing, currency, or governance authority unless the
  scenario is explicitly testing that dependency.
- Peer-local trust / promise accounting: The run should identify what each peer
  can observe and record locally after the stimulus.
- Adversarial or failure pressure: The run should preserve the source row's
  failure, refusal, corruption, ambiguity, scale, or migration pressure; add
  Mallory only when adversarial behavior makes that pressure clearer.
- Human and LLM auditability: The result should be readable enough for a later
  person or model to reconstruct what was promised, what happened, and why the
  evidence matters.
- Migration / evolution path: The run should expose whether protocol, object,
  key, policy, naming, or organizational evolution changes the answer.

## Evaluation Questions

- Which promises or protocol claims does the candidate simulation need in order
  to handle the setup and stimulus?
- What local observations can Alice, Bob, Carol, or another peer record after the
  run?
- Does the candidate design preserve the expected pressure without appealing to
  hidden global state or central authority?
- What DR, DI, frozen spec, TODO, TE, or guide handoff would this evidence inform?

## Result Runs

Result runs live under:

`results/<sim-id>/conditional-release-geofencing-opaque-lower-layer-carriage/<model-id>/<YYYYMMDD-HHMMSS>.md`

## Authority Boundary

This scenario and its result runs are evidence only. Design authority still
graduates through DR, DI, frozen spec, or PromiseGrid Development Guide handoff.
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

