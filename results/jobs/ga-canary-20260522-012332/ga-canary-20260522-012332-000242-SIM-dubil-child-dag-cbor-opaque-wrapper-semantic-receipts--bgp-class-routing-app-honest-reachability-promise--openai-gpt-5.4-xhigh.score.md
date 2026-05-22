# GA Score Cell

Return only JSON with keys `scores`, `fitness`, and `assessment`.
Do not include result identity, source metadata, code fences, or commentary.

## Cell

- Run group ID: `ga-canary-20260522-012332`
- Cell ID: `ga-canary-20260522-012332-000242-SIM-dubil-child-dag-cbor-opaque-wrapper-semantic-receipts--bgp-class-routing-app-honest-reachability-promise--openai-gpt-5.4-xhigh`
- Simulation ID: `SIM-dubil-child-dag-cbor-opaque-wrapper-semantic-receipts`
- Scenario ID: `bgp-class-routing-app-honest-reachability-promise`
- Model ID: `openai-gpt-5.4-xhigh`
- Result path: `proposals/ga-canary-20260522-012332/results/SIM-dubil-child-dag-cbor-opaque-wrapper-semantic-receipts/bgp-class-routing-app-honest-reachability-promise/openai-gpt-5.4-xhigh/20260522-012332.json`

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
fitness evidence, `promisegrid.ga.state.v1` state, generated child sims and
child score evidence under ignored `proposals/<run-group-id>/` trees, explicit
review/promotion, and explicit cleanup via `cull`. Source: `DI-ramar`;
`DI-zanon`; `DI-podot`; `DI-kofil`; `DI-ruzaj`; `DI-fihof`; `DI-lirat`.

Root scenario prompts use `scenarios/README.md` for shared scenario contract
context and `scenarios/<scenario-id>/<scenario-id>.md` for scenario-specific
pressure. Per-scenario `README.md` files are intentionally absent so repeated
boilerplate can be cached once rather than bundled with every scenario. Source:
`DI-kizal`.

GA child-generation prompts keep the parent simulation documents complete but do
not rebundle the root run protocol, root scenario contract, or complete parent
result JSON for every breed call. They include sampled scenario pressure once
and compact score/rationale/risk/open-question summaries from completed parent
fitness results. Source: `DI-dilaf`.

Focused GA canaries may require specific parent simulations or scenarios in the
sample before the remaining slots are filled by deterministic shuffle. Use
`-include-sim` / `-include-scenario` on `tools/ga-runner init`, or
`GA_CANARY_INCLUDE_SIMS` / `GA_CANARY_INCLUDE_SCENARIOS` with the canary wrapper,
when a newly-added sim or scenario must be exercised. Source: `DI-duzur`.

API-backed runs must use explicit cost controls before any large batch. The
runner defaults to concise result style, records provider usage in queue state,
prints actual accumulated cost, and can stop before a cell whose preflight
estimate would exceed the configured budget. Source: `DI-nugiv`; `DI-pulap`.

GA/search provider calls must also send an explicit service tier. The default is
Flex for lower-cost unattended work, `default` must be requested explicitly when
standard processing is desired, and Priority is rejected. OpenAI-compatible
transient failures (`408`, `409`, `429`, `500`, `502`, `503`, `504`, network
timeouts, and request timeouts) use bounded exponential backoff rather than
silently falling back to a higher-cost tier. Source: `DI-mopob`; `DI-tufud`.

Synchronous GA/search calls should be bounded before Batch mode is available:
raw `tools/ga-runner score` and `generate` default to one worker, five-minute
provider attempts, two attempts per cell or child, and a twelve-minute retry
elapsed cap. Streaming is enabled by default for provider liveness logging, with
a two-minute idle timeout for silent stalls. Canary wrappers may raise scoring
and generation workers when cost reservations are active; the terminal canary
defaults to six scoring workers, one child-generation worker, five-minute score
attempts, and fifteen-minute child-generation attempts until a successful
generation phase provides evidence for tighter bounds. It must not dispatch
concurrent cells or children past the configured run budget. Source:
`DI-juzus`; `DI-tufud`; `DI-pivuj`; `DI-suzor`; `DI-guvif`.

Canary wrappers may request `reasoning.summary=auto` and print one no-newline
progress dot per supported reasoning-summary text-delta event while mirroring
reasoning-summary part-done events to stdout/log. The dot is a liveness
diagnostic, not raw hidden reasoning or reasoning-summary delta content;
part-added events and visible-output deltas are intentionally quiet, and raw
commands should keep stdout stream content off unless explicitly requested.
Source: `DI-vadub`; `DI-babik`; `DI-vajut`; `DI-sakam`; `DI-fupob`;
`DI-ramun`.

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
  pool and apply deterministic linear-rank weighted high-parent plus uniform
  random scored non-high parent selection. Source: `DI-tufud`; `DI-puhog`.
- Keep child-generation prompts compact: parent simulation documents once,
  scenario-specific pressure once, and compact fitness summaries instead of full
  parent result JSON. Source: `DI-dilaf`.
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
Generated children are ignored review-stage
`proposals/<run-group-id>/simulations/<SIM-id>/` trees with matching ignored
`proposals/<run-group-id>/results/<SIM-id>/` score evidence. They are not
accepted merely because they exist on disk. A review/promotion pass should
rename selected children to final descriptive non-child `SIM-*` names, fill any
missing standing simulation files, and copy selected score evidence into
canonical `results/` before commit. The detailed operator procedure is
`tools/ga-runner/PROMOTION.md`, used when Steve says
`promote <child-proquint> ...`. Source: `DI-ramar`; `DI-zanon`; `DI-zohal`;
`DI-zusit`; `DI-podot`; `DI-kofil`; `DI-ruzaj`; `DI-gijom`; `DI-fihof`;
`DI-lirat`; `DI-dikoh`.

LLM-generated children use one operation, `breed`, with exactly two distinct
parent simulations. The runner must fail or skip generation rather than create a
one-parent child. Source: `DI-sohus`.

Child breed prompts use compact fitness summaries rather than full parent result
JSON documents. The durable result files remain the source of truth in
`results/`; the prompt carries only the score, rationale, strength, weakness,
risk, and open-question evidence needed for the child to improve the rubric.
Source: `DI-dilaf`.

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

### `proposals/ga-canary-20260522-012332/simulations/SIM-dubil-child-dag-cbor-opaque-wrapper-semantic-receipts/README.md`

```markdown
# SIM-dubil-child-dag-cbor-opaque-wrapper-semantic-receipts: Grid-envelope variant

This simulation breeds `SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid` with `SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0`.

It keeps the minimal positional `[pcid, payload]` envelope and DAG-CBOR encoding, but makes three bounded design moves:

- unknown envelopes are escrowable and forwardable as exact opaque bytes with mandatory local receipts instead of hard drop only
- audited uses standardize on pCID-selected signature wrappers instead of an unsigned-only base story
- applications may derive a stable `semantic_id` from the resolved `(pcid, payload)` pair so meaning can survive outer wrapper or carrier changes

The intent is to preserve the parents' strong layer boundaries and sparse-knowledge friendliness while improving auditability, mixed-version evolution safety, and envelope independence.

The local draft spec is `protocols/grid-envelope.d/specs/grid-envelope-draft.md`.
```

### `proposals/ga-canary-20260522-012332/simulations/SIM-dubil-child-dag-cbor-opaque-wrapper-semantic-receipts/QUESTION.md`

```markdown
# Question

Does a positional DAG-CBOR grid envelope with opaque unknown-envelope escrow, mandatory unsupported or malformed local receipts, and signature wrappers that bind a stable `semantic_id` satisfy the wire-lab scenarios better than either parent on auditability, evolution safety, and envelope independence while keeping the base envelope minimal?
```

### `proposals/ga-canary-20260522-012332/simulations/SIM-dubil-child-dag-cbor-opaque-wrapper-semantic-receipts/protocols/grid-envelope.d/CHANGELOG.md`

```markdown
# CHANGELOG: grid-envelope

A-side CHANGELOG for this simulation-local `grid-envelope` protocol specimen.

This file records freeze events authored by the specimen maintainers. No entries yet; this protocol specimen has not reached a first freeze.

This protocol tree is a simulation-local specimen bred from two parent specimens under the GA run.
```

### `proposals/ga-canary-20260522-012332/simulations/SIM-dubil-child-dag-cbor-opaque-wrapper-semantic-receipts/protocols/grid-envelope.d/specs/grid-envelope-draft.md`

```markdown
# Grid Envelope Variant Spec (DRAFT)

> Status: DRAFT. Not frozen.
> Variant: `enc-dag-cbor` / `unknown-opaque-receipt` / `sig-wrapper-pcid` / `semantic-id-v1`.

## Purpose

This spec defines one full positional grid-envelope candidate for wire-lab comparison. It is a specimen inside `SIM-dubil-child-dag-cbor-opaque-wrapper-semantic-receipts`, not a harness rule and not the canonical PromiseGrid wire format.

The design goal is to keep the parents' narrow envelope boundary while repairing two recurring failures in the fitness evidence:

- hard rejection of unknown pCIDs lost useful evidence and hurt migration
- unsigned minimal carriage weakened later audit, attribution, and replay analysis

## Positional Envelope Shape

The base envelope shape remains:

```text
[pcid, payload]
```

Slots are interpreted positionally:

- `pcid` identifies the protocol or handler that interprets `payload`
- `payload` is opaque bytes until interpreted by the handler named by `pcid`

The base envelope intentionally keeps no fixed signature field, route field, policy field, or transport field.

## Encoding

This variant encodes the envelope as DAG-CBOR-compatible positional arrays.

- `pcid` is a DAG-CBOR Link value
- `payload` is a byte string
- canonical outer bytes are the DAG-CBOR bytes of the exact two-slot array

A payload may itself be the canonical bytes of another grid envelope when the protocol named by `pcid` specifies recursive nesting.

## Unknown pCID Policy

If a receiver lacks a handler for `pcid`, it MUST NOT parse `payload` speculatively.

If `pcid` is unknown, the receiver:

- MAY store the exact outer bytes as opaque content
- MAY forward the exact outer bytes unchanged
- MUST surface an explicit unsupported-pCID result locally
- MUST NOT treat the message as semantically understood

This is neither silent acceptance nor hard evidence loss. It is opaque escrow with explicit local audit.

## Mandatory Local Receipt for Unsupported Envelopes

Whenever a receiver observes a well-formed outer envelope whose `pcid` is unsupported, it MUST create a local receipt.

Minimum local receipt fields:

```text
receipt_type: unsupported-envelope-v1
envelope_cid: <cid or hash of exact outer bytes>
outer_pcid: <unsupported pcid>
source_peer: <local peer label or null>
action: stored | forwarded | dropped
observed_at: <local timestamp or monotonic event id>
note: <optional local diagnostic>
```

This receipt is local evidence only. It does not upgrade an unknown envelope into a valid understood message.

## Malformed Outer Bytes

If received bytes do not decode as the required two-slot DAG-CBOR array, the receiver MUST NOT invent a `pcid`, MUST NOT fabricate payload meaning, and MUST treat the bytes as not-a-valid-envelope under this variant.

The receiver MAY retain the raw bytes unchanged for audit or abuse handling. If it does, or if it explicitly drops after observation, it MUST create a local malformed-envelope receipt.

Minimum local receipt fields:

```text
receipt_type: malformed-envelope-v1
raw_bytes_hash: <hash of observed bytes>
source_peer: <local peer label or null>
action: retained | dropped
observed_at: <local timestamp or monotonic event id>
note: <parse failure summary>
```

This keeps malformed datagram handling below session semantics while still preserving local evidence.

## Signature and Authorship Policy

The base envelope stays minimal and does not embed a fixed signature slot.

For accountable uses, this specimen standardizes a local expectation for wrapper protocols selected by `pcid`:

- audited traffic SHOULD be wrapped by a signature wrapper chosen by its own pCID
- a conforming signature wrapper for this specimen MUST bind both exact inner-envelope bytes and the derived `semantic_id`
- the wrapper SHOULD also bind signer identity and MAY bind contextual assertions such as audience, epoch, route scope, hop policy, or session label

Minimum bound items for a conforming wrapper profile:

```text
inner_envelope_bytes_or_cid
semantic_id
signer
sig_pcid
signature
optional_assertions
```

Binding both exact bytes and `semantic_id` gives two useful audit views at once:

- exact-byte accountability for what was actually sent
- cross-wrapper and cross-carrier comparability for what the sender meant to assert

Peers that do not understand the wrapper pCID fall back to the unknown-pCID receipt behavior above instead of mandatory hard drop with no retained evidence.

## Semantic ID Rule

For any successfully resolved application envelope with tuple `(pcid, payload)`, the receiver MAY derive:

```text
semantic_id = CID(hash(DAG-CBOR([pg-semantic-v1, pcid, payload])))
```

Rules:

- `semantic_id` is derived from the resolved `(pcid, payload)` pair
- outer carrier framing is excluded
- outer signature-wrapper bytes are excluded
- re-encoding of the same meaning under another compatible outer envelope variant does not change `semantic_id` if the same `(pcid, payload)` pair is recovered
- any change to `pcid` or `payload` changes `semantic_id`

Applications that need envelope independence, deduplication, replay tables, or cross-carrier comparison SHOULD key local state on `semantic_id` rather than on outer wrapper bytes alone.

## Layering-Test Behavior

This variant answers the comparison pressures as follows.

- Carrier independence: carriers move exact bytes; carrier choice does not change `semantic_id`.
- Envelope independence: group or session logic can compare recovered `semantic_id` values instead of treating outer envelope bytes as message meaning.
- Unknown typed objects: unknown `pcid` messages can be stored and forwarded opaquely without speculative parsing.
- Raw chunk versus pointer bytes: transported meaning is attached to in-band `pcid`, not to local filenames or container codec alone.
- Sparse knowledge: peers may relay opaque bytes they do not understand while still recording explicit unsupported receipts locally.
- Routing and duplicate advertisements: route offers, refusals, and chunk advertisements still live in payload protocols, but signed wrappers now give a common accountability hook without bloating the base envelope.
- Conditional release and replay pressure: the base envelope does not solve audience or geography binding by itself, but a conforming signature wrapper can bind those assertions and local replay tables can key them to `semantic_id`.
- UDP malformed datagrams: invalid outer bytes never gain invented semantics, yet a local malformed receipt can preserve evidence.

## Non-Goals

This draft does not:

- declare a winning envelope for PromiseGrid
- define a central pCID registry
- require every message to be signed
- claim that opaque forwarding alone is semantic acceptance
- replace the need for payload protocols that define routing, feed, session, or conditional-release meaning

## Why This Should Outperform the Parents

Relative to the hard-reject parent, this variant keeps unsupported traffic auditable and relayable in mixed-version networks.

Relative to the unsigned parent, this variant gives accountable uses a clearer signature-wrapper contract that binds both exact bytes and a stable message meaning identifier.

Relative to both parents, the `semantic_id` rule directly addresses envelope-independence pressure instead of leaving cross-variant equivalence implicit.

## Freeze Gate

This draft can freeze only after at least one simulation run compares it against sibling positional grid-envelope variants and a maintainer signs a merge or freeze promise for this specific specimen.
```

### `scenarios/bgp-class-routing-app-honest-reachability-promise/bgp-class-routing-app-honest-reachability-promise.md`

```markdown
# Honest reachability promise

## Scenario ID

bgp-class-routing-app-honest-reachability-promise

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-punaz-bgp-class-routing-app/SCENARIOS.md`
- Source simulation: `SIM-punaz-bgp-class-routing-app/`
- Source row/title: Honest reachability promise
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-punaz-bgp-class-routing-app/`.

## Setup

Alice advertises reachability to Carol through Bob, and Bob later forwards as promised.

## Stimulus

Run the candidate simulation against this source test: What Alice, Bob, and Carol each record locally after the path works.

## Expected Pressure

Route-like promises need observable kept/broken outcomes without a global route authority.
```

## Required JSON Shape

{"scores":{"scenario_fit":0,"promisegrid_alignment":0,"auditability":0,"evolution_safety":0,"layer_boundary_clarity":0,"failure_handling":0,"implementation_plausibility":0,"risk_penalty":0},"fitness":{"raw":0,"normalized_0_100":0,"confidence_0_1":0.0},"assessment":{"rationale":"","strengths":[],"weaknesses":[],"risks":[],"open_questions":[],"authority_boundary":"Evidence only; does not settle PromiseGrid design."}}
