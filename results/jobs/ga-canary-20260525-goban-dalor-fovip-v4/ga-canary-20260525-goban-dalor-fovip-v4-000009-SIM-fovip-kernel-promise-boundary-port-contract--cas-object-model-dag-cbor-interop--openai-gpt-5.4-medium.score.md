# GA Score Cell

Return only JSON with keys `scores`, `fitness`, and `assessment`.
Do not include result identity, source metadata, code fences, or commentary.

## Cell

- Run group ID: `ga-canary-20260525-goban-dalor-fovip-v4`
- Cell ID: `ga-canary-20260525-goban-dalor-fovip-v4-000009-SIM-fovip-kernel-promise-boundary-port-contract--cas-object-model-dag-cbor-interop--openai-gpt-5.4-medium`
- Simulation ID: `SIM-fovip-kernel-promise-boundary-port-contract`
- Scenario ID: `cas-object-model-dag-cbor-interop`
- Model ID: `openai-gpt-5.4-medium`
- Result path: `results/SIM-fovip-kernel-promise-boundary-port-contract/cas-object-model-dag-cbor-interop/openai-gpt-5.4-medium/20260525-152502.json`

## Rubric

Score each axis from 0 to 5. Higher is better except `risk_penalty`, where 0 is low risk and 5 is severe risk.
Axes: scenario_fit, promisegrid_alignment, auditability, evolution_safety, layer_boundary_clarity, failure_handling, implementation_plausibility, promise_vocabulary, simplicity_durability, envelope_discipline, kernel_implementation_promises, app_protocol_promise_semantics, risk_penalty.
- Score the candidate at its claimed layer. Do not penalize an envelope-layer design for leaving promise accounting, storage semantics, computation semantics, or application-specific trust updates to the payload protocol when the layer boundary is explicit.
- `promise_vocabulary`: reward Promise-Theory-correct, promise-first wording at the candidate's layer. For envelope sims, reward signed pCID-specific promises such as "Alice promises these payload bytes are shaped according to the protocol specification named by this pCID." Penalize normative claim cards, generic claim headers/cards, conformance bundles, generic profile support claims, port claims, universal statement capsules, and central trust-ledger framing.
- `simplicity_durability`: reward small, explicit, deterministic, 100-year durable artifacts that fit small devices. Penalize generic maps, cards, ledgers, bundles, feature-shopping wrappers, and base-envelope selector-shopping stacks such as `env_pCID`/`sig_pCID`/`payload_pCID`.
- `envelope_discipline`: reward alignment with `DN-jotob` and the current envelope direction: CBOR `grid([42(pCID), payload, ...])`, `pCID` as Protocol CID, protocol-owned slot roles after slot 0, local unknown-pCID behavior, and no universal proof-slot overreach.
- `kernel_implementation_promises`: reward explicit local kernel implementation promises: app-facing storage/compute/send/receive/key/lifecycle/dispatch/evidence promises, host assumptions separated from promises, unsupported pCIDs/features, pCID adapter mappings, voluntary namespace/reference behavior, and no kernel-as-authority framing.
- `app_protocol_promise_semantics`: reward higher-layer/app protocols that model storage, computation, send/receive, reciprocal promises, selective sending, promise-as-capability-token behavior, make/break evidence, and local trust updates without command, permission, policy-enforcement, or conformance-authority framing.
- Do not treat higher-layer pCID-owned payload protocols as selector shopping merely because they define their own signed refusal records, exact-byte observation evidence, freeze successor records, transfer semantics, or capability-as-promise-token payload behavior.
- `scenario_fit` may be lower when a scenario asks for higher-layer behavior the candidate intentionally delegates to the payload protocol, but that delegation is not by itself a PT-gate violation when the layer boundary and local trust boundary are clear.
- The runner recomputes `fitness.raw` and `fitness.normalized_0_100` from your axis scores with normal weighting. Use `fitness.confidence_0_1` for your confidence.

## Promise Theory Fundamentals

Apply these Mark Burgess reference notes while scoring:
- Agents are autonomous.
- A promise is a scoped declaration of intent.
- No agent can make a promise on behalf of another agent.
- Promises do not guarantee outcomes.
- Trust is a local assessment of whether a promise will be kept.
- Promises to receive or use are not equivalent to obligations, impositions, or promises to give.
Reference notes: Mark Burgess, *In Search of Certainty*; *Promise Theory: Principles and Applications*; *Thinking in Promises*.

## PT Gate

Classify the design as exactly one of `pt_clean`, `pt_reframe_needed`, or `pt_invalid`.
- `pt_clean`: promise-first and locally trustworthy enough to compete normally.
- `pt_reframe_needed`: technically interesting but drifts into non-PT framing; it may survive only as a question-home or rework candidate.
- `pt_invalid`: relies on authority, imposition, global trust, or RPC-style command semantics; it cannot be promotable.
Complete every PT rule check in `assessment.pt_gate`, and explain violations in `assessment.pt_gate.violations`.

Required score-axis checklist: `scenario_fit`, `promisegrid_alignment`, `auditability`, `evolution_safety`, `layer_boundary_clarity`, `failure_handling`, `implementation_plausibility`, `promise_vocabulary`, `simplicity_durability`, `envelope_discipline`, `kernel_implementation_promises`, `app_protocol_promise_semantics`, `risk_penalty`.
A response missing any required `scores` axis or `assessment.pt_gate` is invalid.

## Source Documents

### `results/RUN-PROTOCOL.md`

```markdown
# Results Run Protocol

This document defines operational contracts for result evidence under `results/`.
Legacy matrix runs write Markdown files at
`results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md`. GA/search runs
keep historical `promisegrid.ga.result.v1` JSON fitness files, write new
`promisegrid.ga.result.v2` JSON fitness files at
`results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.json`, and
checkpoint state at `results/state/<run-group-id>.json`. Source: `DI-zamin`;
`DI-ramar`; `DI-zanon`; `DI-ruzaj`; `DI-roruj`.

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
review/promotion, explicit cleanup via `cull`, and audit-first rubric-v2
backfill planning via `audit` / `backfill-init`. Source: `DI-ramar`;
`DI-zanon`; `DI-podot`; `DI-kofil`; `DI-ruzaj`; `DI-fihof`; `DI-lirat`;
`DI-higot`; `DI-roruj`.

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
when a newly-added sim or scenario must be exercised. The canary wrapper also
accepts `/tmp/canary-cells` as a plain-text focus file with `sims:` and
`scenarios:` sections; selectors resolve by unique prefix against current
sim/scenario IDs, merge with any `GA_CANARY_INCLUDE_*` values, and must be
valid before any provider-backed work begins. Source: `DI-duzur`; `DI-bataj`.

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

GA rubric-v2 score calls now default to provider-enforced structured outputs via
an explicit score output contract, `json_schema_strict`. This is transport
hardening, not a rubric change: new score results should record
`runner.output_contract` so `prompt_json` and provider-enforced schema runs can
coexist in the corpus without rewriting history. The explicit fallback path is
`score -output-contract prompt_json`. Adopting structured outputs does not by
itself authorize a full-corpus rescore; first rerun any affected failed cells
and use a small calibration slice to check for material rank drift. Source:
`DI-fogop`.

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

GA/search results use schema `promisegrid.ga.result.v1` or
`promisegrid.ga.result.v2` and path
`results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.json`. The JSON
file is the fitness evidence; there is no separate `results/fitness/` tree.
Required content includes source paths, source hashes, runner/model metadata,
rubric axes, integer rubric scores, normalized fitness, rationale, risks, open
questions, and authority boundary. Historical v1 evidence remains append-only;
new promise-first vocabulary scoring writes v2. Source: `DI-ramar`; `DI-zanon`;
`DI-pobus`; `DI-ruzaj`; `DI-roruj`.

Promoted canonical results may preserve historical `source.*` proposal paths
even after the proposal tree is deleted. Audit/backfill code may compare those
historical source hashes against a byte-identical canonical `simulations/<sim-id>/`
tree as a current exact-match fallback, but the stored JSON evidence is not
rewritten. Source: `DI-zobur`.

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
- Audit canonical v1 evidence before targeted vocabulary-aware backfill:
  `cd tools/ga-runner && go run . audit -repo-root ../..`
- Create a targeted rubric-v2 backfill state from canonical v1 evidence:
  `cd tools/ga-runner && go run . backfill-init -repo-root ../.. -run-group-id <run-group-id>`
  Optional staged override for honest multi-stage rescoring:
  `cd tools/ga-runner && go run . backfill-init -repo-root ../.. -run-group-id <run-group-id> -staged-model-id openai-gpt-5.4-high -staged-reasoning-effort high`
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
accepted merely because they exist on disk. A review/promotion pass must not
rewrite scored result JSONs or scored simulation files in place. If canonical
homes are needed, use a byte-identical move/copy path or leave the scored
artifacts under `proposals/` until that path is locked. The detailed operator
procedure is `tools/ga-runner/PROMOTION.md`, used when Steve says
`promote <child-proquint> ...`. Source: `DI-ramar`; `DI-zanon`; `DI-zohal`;
`DI-zusit`; `DI-podot`; `DI-kofil`; `DI-ruzaj`; `DI-gijom`; `DI-fihof`;
`DI-lirat`; `DI-dikoh`; `DI-zadik`; `DI-higot`; `DI-roruj`; `DI-hijub`.

`backfill-init` preserves historical source model lineage by default. When a
targeted rerun must score under a new model ID or reasoning default, pass
`-staged-model-id` and optionally `-staged-reasoning-effort`. The emitted
state, planned cell `model_id`, result directories, and derived provider
`api_model` then match the new stage instead of reusing historical lineage from
the source evidence. Source: `DI-hijub`.

Broad GA parent scoring now defaults to `medium` reasoning effort. Use `xhigh`
explicitly for tie-breaks, promotion candidates, and design-state-sensitive
comparisons where the extra cost is justified. Source: `DI-nanor`.

GA scorer prompts must enumerate the full rubric-v2 axis list and treat any
missing required score axis as an invalid response. Under
`-output-contract json_schema_strict`, schema adherence is delegated to the
provider and missing-axis responses should surface as provider or validation
failures rather than a local schema-correction loop. Under the explicit
fallback `-output-contract prompt_json`, `score` may send one targeted
schema-correction retry, accumulate both call costs in state, and still reject
the cell if the retry remains incomplete. It must not auto-fill missing scores
locally. Source: `DI-kibuf`; `DI-fogop`.

Audit-first rubric-v2 backfill should compare historical source hashes against
current sim/scenario bytes while reporting root-contract drift separately. This
lets v2 reruns focus on stable sim/scenario evidence even when the root scoring
contract docs themselves have evolved. Targeted backfill includes exact-match
hard-hit vocabulary families plus a clean grid-envelope calibration slice before
any wider rerun is attempted. Source: `DI-roruj`.

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

For audit-first rubric-v2 rescoring, the comparison report is a derived Markdown
artifact under `results/reports/<run-group-id>-comparison.md`. It compares each
completed v2 backfill cell against the latest exact-match canonical
`promisegrid.ga.result.v1` record for the same `sim_id` + `scenario_id`,
preferring the same `runner.api_model` when available. Same-model historical
reruns collapse to the latest record and are reported separately from true
ambiguous historical pairings. The report summarizes sim-rank drift plus large
per-cell deltas before any broader rerun. Source: `DI-zuzup`; `DI-sirir`.

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

### `simulations/SIM-fovip-kernel-promise-boundary-port-contract/README.md`

```markdown
# SIM-fovip: Kernel promise boundary port contract

This simulation is the active evidence surface for the `DR-davod` kernel design
packet. It follows `TE-mazop`, which found that `TE-jimar` was not enough to
close the kernel-developer porting boundary, and now incorporates the follow-on
`TE-pudiv`, `TE-dunas`, and `TE-gakoh` questions. Source: `DI-funaf`;
`DI-fidot`.

## Question

What minimum PromiseGrid kernel implementation promises can be claimed across
rich native nodes, browser/WASM hosts, mobile sandboxes, MCU/header-only ports,
and split local service graphs without pretending that one process shape,
namespace, API, or prior-art pattern is universal?

## Candidate specimen

The specimen under test is a kernel implementation promise record: a local
implementation promises its own behavior by publishing:

- profile name and runtime class;
- supported pCIDs and unsupported-pCID behavior;
- app-facing promises for storage, compute, network send/receive, key use,
  device access, lifecycle, pCID dispatch, refusal, receipt, evidence,
  namespace, reference, and checkpoint behavior;
- host/runtime assumptions that the port depends on but does not itself promise;
- explicitly unsupported features;
- exact evidence records for kept, refused, unavailable, and broken promises;
- adapter promises when local APIs wrap pCID-selected grid messages;
- voluntary namespace promises when the port projects group namespaces;
- CID-rooted promise-bound reference behavior for cross-agent path sharing.

## Runtime classes

- **Native node:** Bob runs a local service with storage, network, keys, feed,
  CAS, dispatch, lifecycle, and evidence roles.
- **Browser/WASM:** Alice runs inside a host that owns storage, network, clocks,
  keys, and lifecycle.
- **Mobile sandbox:** Dave can run only while the OS permits background work and
  network access.
- **MCU/header-only:** Carol supports one or two pCIDs, bounded evidence, and a
  hardware/device promise.
- **Split local services:** Ellen separates dispatch, storage, keys, networking,
  and evidence into local services with separate promises.

## Basic principles under test

- Kernel is a role/profile set, not a ruler.
- Everything useful is a promise: app/kernel operations, resources, namespaces,
  references, and kernel implementation promise records all help agents make or
  evaluate promises.
- The app/kernel boundary is a promise boundary; exposed operations are
  pCID-selected `grid([42(pCID), payload, ...])` messages, even when local APIs
  provide ergonomic adapters.
- A kernel implementation promise record is not a global certificate. Alice,
  Bob, Carol, and later agents evaluate the record and make/break history
  locally.
- Host assumptions are not implementation promises unless the host is also an
  explicit promiser.
- Voluntary group namespaces may exist inside trust relationships, but imposed
  universal namespaces are rejected.
- File-like resources are promise-log projections or checkpoints, not evidence
  that PromiseGrid is filesystem-first.

## Evidence axes

The simulation should let reviewers ask whether the candidate:

- names the local promiser for each storage, compute, network, key, device,
  lifecycle, dispatch, namespace, reference, and evidence promise;
- maps every exposed app/kernel operation to a pCID-selected message or explains
  why the operation is outside the PromiseGrid boundary;
- records exact bytes when needed for proof, replay, unsupported-pCID carriage,
  or broken-promise evidence;
- states host/runtime assumptions separately from the port's own promises;
- names unsupported pCIDs and unsupported roles directly;
- keeps trust local and avoids global permission, namespace, conformance, or
  policy authority;
- treats V, Amoeba, Plan 9, and Hurd as pattern pressure, not imported design
  authority;
- supports voluntary group namespaces and CID-rooted promise-bound references
  without treating Alice's local path as global truth;
- represents file/resource state as checkpoints or projections over selected
  promise/event log frontiers.

## Boundaries

This simulation does not close `DR-davod` and does not define a final
PromiseGrid kernel API. It tests whether kernel implementation promises give
guide writers enough evidence to discuss kernel developers without promising a
daemon, microkernel, browser host, mobile runtime, MCU library, namespace
protocol, or SDK as the single correct implementation shape.

The current envelope shape `grid([42(pCID), payload, ...])` is input evidence,
not a reopened decision.

## Related evidence

- `docs/research/DN-lujad-promisegrid-kernel-role-profile.md`
- `docs/thought-experiments/TE-jimar-kernel-runtime-portability-boundary.md`
- `docs/thought-experiments/TE-mazop-kernel-promise-boundary-and-minimum-port-contract.md`
- `docs/thought-experiments/TE-pudiv-app-kernel-grid-message-boundary.md`
- `docs/thought-experiments/TE-dunas-prior-art-influence-on-promisegrid-kernel.md`
- `docs/thought-experiments/TE-gakoh-local-views-over-promise-event-hypergraph.md`
- `DR/DR-davod-promisegrid-kernel-dev-porting-boundary.md`
- `scenarios/kernel-porting-boundary/kernel-porting-boundary.md`
- `simulations/SIM-funas-kernel-porting-boundary/`
```

### `simulations/SIM-fovip-kernel-promise-boundary-port-contract/QUESTION.md`

```markdown
# Question

Which minimum PromiseGrid kernel implementation promises can guide writers
describe without turning one runtime shape into a false universal kernel, and
without treating kernel, host, namespace, app API, or prior-art patterns as
external authority?

The answer must be concrete enough to test app-facing promises, host/runtime
assumptions, unsupported features, pCID coverage, exact-byte evidence records,
broken-promise handling, app/kernel pCID messages, voluntary group namespaces,
CID-rooted promise-bound references, and file/resource checkpoints over promise
logs across native, browser/WASM, mobile, MCU, and split local-service
deployments. Source: `DI-funaf`; `DI-fidot`; `DR-davod`; `TE-mazop`.
```

### `simulations/SIM-fovip-kernel-promise-boundary-port-contract/protocols/kernel-port.d/specs/kernel-port-contract-draft.md`

```markdown
# Kernel implementation promise record draft

This draft is a simulation-local specimen for `SIM-fovip`. It is not a frozen
PromiseGrid protocol spec. Its purpose is to make the `DR-davod` question
testable by forcing a port to publish the promises, assumptions, unsupported
features, and evidence records that make a first implementation credible. Source:
`DI-funaf`; `DI-fidot`.

## Role

A kernel implementation promise record says what a PromiseGrid implementation
promises to local apps and operators, what it depends on from the host runtime,
and what it explicitly does not promise.

The record is not a global certificate, permission, or authority. Each receiving
agent still assesses the record and the port's make/break history locally.

## Record shape

```text
kernel_implementation_promise_record = [
  record_pcid,
  port_identity,
  profile,
  supported_pcids,
  app_facing_promises,
  host_assumptions,
  unsupported_features,
  evidence_policy,
  adapter_promises,
  namespace_promises,
  reference_promises
]
```

## Fields

- `record_pcid` is the Protocol CID for this promise-record shape.
- `port_identity` names the local implementation or agent making the promises.
- `profile` names the runtime class: native node, browser/WASM, mobile sandbox,
  MCU/header-only, split local services, or another pCID-defined profile.
- `supported_pcids` lists the pCID-selected protocols the implementation promises to
  parse, dispatch, validate, or preserve.
- `app_facing_promises` states what the implementation promises for storage,
  compute, network send/receive, key use, device access, lifecycle, pCID
  dispatch, refusal, receipt, evidence recording, namespace projection,
  reference resolution, and resource checkpoint behavior.
- `host_assumptions` states what the port depends on from a browser, OS, mobile
  sandbox, language runtime, hardware platform, or local service graph.
- `unsupported_features` states what the port refuses or cannot perform.
- `evidence_policy` states what exact bytes and local records the implementation promises
  to keep for kept, refused, unavailable, and broken promises.
- `adapter_promises` states which local APIs wrap which pCID-selected messages
  and what evidence the adapter records.
- `namespace_promises` states whether the port projects voluntary group
  namespaces and which promisers maintain those namespace frontiers.
- `reference_promises` states how the port handles CID-rooted promise-bound
  references and local path mounting.

## Promise rules

- An implementation promises only its own behavior.
- An implementation may cite host/runtime assumptions, but it does not promise
  that the host will keep them unless the host is also an explicit promiser.
- Unsupported features must be named directly. They must not be hidden behind a
  generic "partial implementation" label.
- Evidence records are local. They help Alice, Bob, Carol, and future agents
  update their own trust judgments; they are not a global trust authority.
- Local APIs are adapters. If an API call crosses a PromiseGrid promise
  boundary, the record must identify the corresponding pCID-selected message or
  state that the operation is outside the PromiseGrid boundary.
- Voluntary group namespaces are promises among agents. The record must not
  describe a namespace as universal truth.
- File-like resources are projections over promises, logs, and checkpoints. The
  record must say what evidence is preserved for the selected frontier.

## Minimum credible first port

A first credible port is allowed to be small. It must still publish:

- at least one supported pCID or a bounded exact-byte carriage profile;
- clear unsupported-pCID behavior;
- app-facing promises for every operation it exposes;
- host/runtime assumptions for every operation it delegates;
- evidence records for kept, refused, unavailable, and broken promises;
- adapter, namespace, reference, and checkpoint promises where the port exposes
  those surfaces;
- an implementation promise record that can be compared with later behavior.

## Scenario pressure

The same record shape must be tested against:

- a native node with broad local services;
- a browser/WASM host with delegated storage, network, key, and lifecycle
  behavior;
- a mobile sandbox with restricted background execution;
- an MCU/header-only port with one pCID and bounded evidence;
- a split local service graph with multiple local promisers;
- a voluntary group namespace maintained by Alice, Bob, and Carol;
- a CID-rooted promise-bound reference from Alice that Bob mounts locally;
- a file/resource checkpoint reconstructed from a selected promise-log frontier.
```

### `scenarios/cas-object-model-dag-cbor-interop/cas-object-model-dag-cbor-interop.md`

```markdown
# DAG-CBOR interop

## Scenario ID

cas-object-model-dag-cbor-interop

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-jomag-cas-object-model/SCENARIOS.md`
- Source simulation: `SIM-jomag-cas-object-model/`
- Source row/title: DAG-CBOR interop
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-jomag-cas-object-model/`.

## Setup

Alice stores a Merkle node or pointer object using a DAG-CBOR-compatible representation.

## Stimulus

Run the candidate simulation against this source test: Whether CID links, byte strings, and tags stay compatible with IPFS / IPLD-style tooling without requiring those stacks.

## Expected Pressure

TE-43 must decide whether DAG-CBOR is the default object format or only one allowed profile.
```

## Required JSON Shape

{"scores":{"scenario_fit":0,"promisegrid_alignment":0,"auditability":0,"evolution_safety":0,"layer_boundary_clarity":0,"failure_handling":0,"implementation_plausibility":0,"promise_vocabulary":0,"simplicity_durability":0,"envelope_discipline":0,"kernel_implementation_promises":0,"app_protocol_promise_semantics":0,"risk_penalty":0},"fitness":{"raw":0,"normalized_0_100":0,"confidence_0_1":0.0},"assessment":{"rationale":"","strengths":[],"weaknesses":[],"risks":[],"open_questions":[],"authority_boundary":"Evidence only; does not settle PromiseGrid design.","pt_gate":{"status":"pt_clean","autonomous_agents":{"status":"pass","note":""},"scoped_intent":{"status":"pass","note":""},"no_promises_for_others":{"status":"pass","note":""},"no_guaranteed_outcomes":{"status":"pass","note":""},"local_trust_assessment":{"status":"pass","note":""},"accept_use_not_obligation":{"status":"pass","note":""},"violations":[]}}}
