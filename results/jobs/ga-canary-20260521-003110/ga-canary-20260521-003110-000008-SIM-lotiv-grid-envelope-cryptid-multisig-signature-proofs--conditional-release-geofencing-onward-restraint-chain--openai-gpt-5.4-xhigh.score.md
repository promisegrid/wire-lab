# GA Score Cell

Return only JSON with keys `scores`, `fitness`, and `assessment`.
Do not include result identity, source metadata, code fences, or commentary.

## Cell

- Run group ID: `ga-canary-20260521-003110`
- Cell ID: `ga-canary-20260521-003110-000008-SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs--conditional-release-geofencing-onward-restraint-chain--openai-gpt-5.4-xhigh`
- Simulation ID: `SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs`
- Scenario ID: `conditional-release-geofencing-onward-restraint-chain`
- Model ID: `openai-gpt-5.4-xhigh`
- Result path: `results/SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs/conditional-release-geofencing-onward-restraint-chain/openai-gpt-5.4-xhigh/20260521-003110.json`

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
  pool and apply deterministic highest-scoring-parent plus uniform random scored
  non-top parent selection. Source: `DI-tufud`.
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
missing standing simulation files, and move selected score evidence into
canonical `results/` before commit. Source: `DI-ramar`; `DI-zanon`; `DI-zohal`;
`DI-zusit`; `DI-podot`; `DI-kofil`; `DI-ruzaj`; `DI-gijom`; `DI-fihof`;
`DI-lirat`.

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

### `simulations/SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs/README.md`

```markdown
# SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs: Grid-envelope Cryptid Multisig signature/proof probe

This simulation is a standalone, non-child grid-envelope specimen. It tests
whether envelope signatures and proofs can use Cryptid's Multisig object model
as the signature/proof payload representation while keeping PromiseGrid's
signature-placement and verification questions unresolved. Source: `DI-sahiv`.

The upstream prior-art source is Cryptid's Multisig Specification v0.0.1,
currently marked pre-draft. This simulation treats that format as design input:
the Multisig object starts with the Multisig sigil `0x1239` encoded as varuint
`0xb924`, then carries a signing-codec sigil, optional message bytes, and a
counted sequence of attributes. Source: `DI-sahiv`.

The local draft spec is
`protocols/grid-envelope.d/specs/grid-envelope-draft.md`.

## Design Pressure

- **Detached versus combined:** the Multisig message field may be empty
  (detached) or present (combined), so the sim can compare signatures over
  outer envelope bytes, nested payload bytes, or in-object message bytes.
- **Envelope versus nested payload:** the same Multisig bytes can occupy an
  outer signature/proof slot or live inside a payload protocol's nested object.
- **Variable arity:** Multisig's counted attributes let one signature object
  carry extra codec-specific proof material without changing the outer envelope
  shape.
- **pCID interaction:** the envelope `pcid` still selects payload semantics, and
  a `sig_pcid` or payload schema decides how Multisig verification is invoked.
- **Unknown codecs:** generic tools that understand varuint and varbytes can
  skip unknown Multisig signing codecs or unknown attributes without claiming
  verification.
- **Threshold shares:** attributes such as `Scheme`, `Threshold`, `Limit`,
  `ShareIdentifier`, and `ThresholdData` let the sim test individual shares,
  accumulation, and final aggregate verification.
- **Verifier obligations:** verifiers must bind the exact message bytes, signing
  codec, verifying key material, required attributes, threshold policy, and
  payload interpretation before accepting an envelope.

## Non-Canonical Status

This simulation does not choose a final PromiseGrid envelope shape, does not
freeze a pCID, does not require Cryptid Multisig as a PromiseGrid dependency,
and does not supersede the existing positional, arity, nested-signature, or
generated child grid-envelope specimens. Source: `DI-sahiv`; open long-term
envelope pressure remains represented by `DR-009-20260430-204108`.
```

### `simulations/SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs/QUESTION.md`

```markdown
# Question

Can grid-envelope signatures and proofs use Cryptid's Multisig object model as
their signature/proof payload representation while preserving unresolved
PromiseGrid choices about detached versus combined signatures, outer versus
nested placement, variable arity, pCID binding, unknown-codec handling,
threshold-share aggregation, and verifier obligations? Source: `DI-sahiv`;
`DR-009-20260430-204108`.
```

### `simulations/SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs/protocols/grid-envelope.d/CHANGELOG.md`

```markdown
# CHANGELOG: grid-envelope

A-side CHANGELOG (per TE-liviv) for this simulation-local `grid-envelope`
protocol specimen.

This file records freeze events authored by the specimen maintainers. No entries
yet; this protocol specimen has not reached a first freeze.

This protocol tree is a simulation-local specimen created by `DI-sahiv`.
```

### `simulations/SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs/protocols/grid-envelope.d/specs/grid-envelope-draft.md`

```markdown
# Grid-envelope draft: Cryptid Multisig signature/proof payloads

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `cryptid-multisig-signature-proofs`.

## Scope

This spec defines one grid-envelope candidate for wire-lab comparison. It is a
specimen inside `SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs`,
not a harness rule and not the canonical PromiseGrid envelope. Source:
`DI-sahiv`.

The design under test is narrow: use Cryptid's Multisig object model as the
signature/proof payload representation while leaving the envelope-level
placement and verification choices open for simulation pressure. The Multisig
source is upstream prior art and is not treated as a frozen PromiseGrid
dependency. Source: `DI-sahiv`.

## Prior-Art Multisig Shape

Cryptid's pre-draft Multisig v0.0.1 encodes a codec-agnostic digital signature
object as:

```text
multisig_sigil signing_codec_sigil message attributes
```

The object starts with Multisig sigil `0x1239`, encoded as varuint `0xb924`.
It then carries a signing-codec sigil, an optional `message` encoded as
varbytes, and a variable number of attributes encoded as a count followed by
attribute-id and varbytes pairs. Source: `DI-sahiv`.

This simulation recognizes these upstream attribute roles as pressure inputs:

- `SigData` for signature bytes;
- `PayloadEncoding` for the signed-message encoding sigil;
- `Scheme` for threshold-signing scheme;
- `Threshold` for the minimum share count required;
- `Limit` for total share count;
- `ShareIdentifier` for the share number or participant-local share label;
- `ThresholdData` for codec-specific threshold material;
- `AlgorithmName` for application-specific or non-standard algorithm naming.

## Envelope Shapes Under Test

The simulation keeps three placement modes alive instead of choosing one:

```text
[pcid, payload, sig_pcid, multisig]
[pcid, payload_with_nested_multisig]
[pcid, combined_multisig]
```

Slots are interpreted positionally only when the selected mode and its `pcid`
or `sig_pcid` define them:

- `pcid` identifies the payload protocol, handler, or proof-bearing schema.
- `payload` is opaque bytes until interpreted by the handler named by `pcid`.
- `sig_pcid`, when present, identifies the signature/proof protocol that says
  the fourth slot is a Cryptid-style Multisig object.
- `multisig` is the exact Multisig object bytes, not decoded fields projected
  into the envelope.

The first mode pressures explicit outer signature dispatch. The second mode
pressures nested payload ownership of signatures. The third mode pressures
combined Multisig objects where the signed message is carried inside the
Multisig `message` varbytes field rather than as a sibling envelope payload.
Source: `DI-sahiv`; `DR-009-20260430-204108`.

## Encoding

The envelope carrier for this specimen is deterministic CBOR positional arrays.
`pcid` and `sig_pcid`, when present, are CIDv1 byte strings or DAG-CBOR links as
defined by the concrete run profile. `payload` and `multisig` are byte strings.
The Multisig object itself keeps its own varuint and varbytes internal encoding;
the envelope does not translate Multisig attributes into CBOR fields.

Canonical envelope bytes are the deterministic CBOR bytes of the selected
outer array. Canonical Multisig bytes are the exact varuint/varbytes bytes
carried in the Multisig slot or nested payload. A verifier must never verify a
re-serialized approximation when the original bytes are available.

## Detached and Combined Signature Policy

A Multisig object with an empty message field is treated as detached. The
verifier must obtain the signed bytes from the selected envelope mode:

- outer explicit mode signs the canonical unsigned prefix `[pcid, payload]`
  unless `sig_pcid` defines stricter associated data;
- nested mode signs the nested bytes selected by the payload protocol named by
  `pcid`;
- share-collection mode signs the same byte string for every share before
  threshold aggregation.

A Multisig object with a non-empty message field is treated as combined. The
verifier must compare the embedded message bytes against the envelope-selected
payload binding before accepting the envelope. If the embedded message and
outer payload disagree, verification fails even if the Multisig cryptographic
check succeeds. Source: `DI-sahiv`.

## Unknown Codec and Attribute Policy

This specimen adopts a conservative unknown-codec rule:

- A receiver that does not understand the envelope `pcid` must not interpret
  `payload`, even if it can skip or parse Multisig framing.
- A receiver that understands `pcid` but not `sig_pcid` may preserve the exact
  Multisig bytes as opaque proof evidence, but must not claim verification.
- A receiver that understands Multisig varuint/varbytes framing but not the
  signing codec may skip the object, index its byte range, and preserve it for
  later tooling, but must not treat unknown `SigData` as valid.
- A receiver may ignore unknown non-critical attributes only if the signing
  codec or `sig_pcid` says they are non-critical; otherwise unknown attributes
  keep the verification result unsupported.

This policy intentionally separates structural skippability from cryptographic
acceptance. Skipping unknown signing codecs helps storage, relay, and future
audit; it is not a validity decision. Source: `DI-sahiv`.

## Threshold and Multi-Payload Pressure

Threshold runs use the Multisig attributes as follows:

- every share carries `Scheme`, `Threshold`, `Limit`, and `ShareIdentifier`
  values that must agree with the threshold policy named by `sig_pcid` or the
  payload schema;
- `SigData` carries the share or aggregate signature data as defined by the
  signing codec;
- `ThresholdData` carries codec-specific accumulation material when the codec
  needs more than raw signature bytes;
- `PayloadEncoding` records the signed-message encoding when the signature
  codec requires that value for replay-safe verification;
- `AlgorithmName` is advisory unless the selected `sig_pcid` makes it part of
  the verification policy.

The simulation should penalize designs that let two shares over different
message bytes, different `pcid` values, or different threshold policies
aggregate as if they belonged to the same proof.

## pCID Interaction

The envelope `pcid` continues to name the payload protocol. Multisig does not
replace pCID-selected payload semantics. A `sig_pcid`, nested payload schema, or
future frozen profile must define:

- whether the signed bytes include `pcid`, `payload`, both, or an enclosing
  transcript;
- whether `PayloadEncoding` duplicates, complements, or constrains the pCID;
- which Multisig signing codecs are acceptable for the payload protocol;
- how verifier key material is found, authenticated, rotated, and revoked;
- whether unknown attributes are reject, ignore, quarantine, or relay-only
  evidence.

Until those choices are frozen, this specimen is evidence for verifier
obligations rather than a PromiseGrid validity rule. Source: `DI-sahiv`;
`DR-009-20260430-204108`.

## Verifier Obligations

A verifier accepts an envelope under this specimen only after all of the
following hold:

- the envelope shape is recognized by `pcid`, `sig_pcid`, or a local run
  profile;
- the exact signed byte string is determined without ambiguity;
- detached and combined-message bindings agree with the selected envelope mode;
- the signing codec is understood and allowed by the selected signature policy;
- required attributes are present, canonical, non-duplicated unless the codec
  permits duplicates, and internally consistent;
- threshold shares all bind to the same message, pCID context, threshold policy,
  and signer set before aggregation;
- unknown codecs or critical attributes produce an unsupported or quarantined
  result rather than a successful verification;
- local audit records retain exact envelope bytes, exact Multisig bytes, the
  verification profile, and the reason for accept, reject, unsupported, or
  quarantine.

## Scenario Pressure Notes

### Normal detached signature

Alice sends `[pcid, payload, sig_pcid, multisig]` where `multisig` has an empty
message field. Bob verifies that the signature covers the canonical unsigned
prefix `[pcid, payload]`, not only `payload`, so replay under a different
payload protocol fails.

### Combined signature mismatch

Carol sends `[pcid, payload, sig_pcid, multisig]` where the Multisig message
field contains bytes that differ from `payload`. Dave must reject the envelope
because the embedded message and envelope-selected payload binding disagree.

### Unknown signing codec

Ellen understands the envelope `pcid` and Multisig framing but lacks the
signing-codec implementation. Ellen may keep and relay the exact proof bytes as
unsupported evidence, but she must not mark the envelope verified.

### Threshold share collection

Frank receives three BLS-style shares that claim a threshold of three out of
four. He aggregates only shares whose message bytes, `pcid` binding, scheme,
threshold, limit, and signer-set policy match; mixed-context shares remain
separate unsupported evidence.

### Nested payload proof

Alice sends `[pcid, payload_with_nested_multisig]`. Bob can verify only after
the `pcid` payload schema identifies the nested Multisig byte range and the
exact bytes the nested proof signs.

## Non-Goals

This draft does not:

- declare a winning envelope;
- freeze a pCID;
- require a central pCID registry;
- decide detached versus combined signatures globally;
- decide whether signatures live in the envelope or nested payload;
- define final PromiseGrid key discovery, revocation, freshness, or authority;
- claim that Cryptid Multisig is stable or normative for PromiseGrid.

## Freeze Gate

This draft can freeze only after simulation runs compare it against sibling
grid-envelope signature-placement specimens, at least one verifier profile
fully specifies pCID binding and unknown-codec behavior, and a maintainer signs
a merge/freeze promise for this specific specimen.
```

### `scenarios/conditional-release-geofencing-onward-restraint-chain/conditional-release-geofencing-onward-restraint-chain.md`

```markdown
# Onward-restraint chain

## Scenario ID

conditional-release-geofencing-onward-restraint-chain

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-zarud-conditional-release-geofencing/SCENARIOS.md`
- Source simulation: `SIM-zarud-conditional-release-geofencing/`
- Source row/title: Onward-restraint chain
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-zarud-conditional-release-geofencing/`.

## Setup

Alice sends content to Bob only if Bob promises to forward it only to recipients who make the same promise. Bob wants to forward to Carol.

## Stimulus

Run the candidate simulation against this source test: Whether the recursive promise graph is represented at group-session, conditional-release, transport/feed, or CAS-object level.

## Expected Pressure

If the graph is central to dispatch semantics, group/session ownership gets stronger; if it composes across sessions, a separate family gets stronger.
```

## Required JSON Shape

{"scores":{"scenario_fit":0,"promisegrid_alignment":0,"auditability":0,"evolution_safety":0,"layer_boundary_clarity":0,"failure_handling":0,"implementation_plausibility":0,"risk_penalty":0},"fitness":{"raw":0,"normalized_0_100":0,"confidence_0_1":0.0},"assessment":{"rationale":"","strengths":[],"weaknesses":[],"risks":[],"open_questions":[],"authority_boundary":"Evidence only; does not settle PromiseGrid design."}}
