# GA Score Cell

Return only JSON with keys `scores`, `fitness`, and `assessment`.
Do not include result identity, source metadata, code fences, or commentary.

## Cell

- Run group ID: `ga-canary-20260525-hozif-anti-rpc-calibration`
- Cell ID: `ga-canary-20260525-hozif-anti-rpc-calibration-000014-SIM-fozid-child-slim-port-commitment-service-outcomes--promise-accounting-records-kept-storage-promise--openai-gpt-5.4-medium`
- Simulation ID: `SIM-fozid-child-slim-port-commitment-service-outcomes`
- Scenario ID: `promise-accounting-records-kept-storage-promise`
- Model ID: `openai-gpt-5.4-medium`
- Result path: `proposals/ga-canary-20260525-hozif-anti-rpc-calibration/results/SIM-fozid-child-slim-port-commitment-service-outcomes/promise-accounting-records-kept-storage-promise/openai-gpt-5.4-medium/20260526-003243.json`

## Rubric

Score each axis from 0 to 5. Higher is better except `risk_penalty`, where 0 is low risk and 5 is severe risk.
Axes: scenario_fit, promisegrid_alignment, auditability, evolution_safety, layer_boundary_clarity, failure_handling, implementation_plausibility, promise_vocabulary, simplicity_durability, envelope_discipline, kernel_implementation_promises, app_protocol_promise_semantics, risk_penalty.
- Score the candidate at its claimed layer. Do not penalize an envelope-layer design for leaving promise accounting, storage semantics, computation semantics, or application-specific trust updates to the payload protocol when the layer boundary is explicit.
- `promise_vocabulary`: reward Promise-Theory-correct, promise-first wording at the candidate's layer. For envelope sims, reward signed pCID-specific promises such as "Alice promises these payload bytes are shaped according to the protocol specification named by this pCID." Penalize normative claim cards, generic claim headers/cards, conformance bundles, generic profile support claims, port claims, universal statement capsules, central trust-ledger framing, RPC dispatcher framing, service-registry framing, capability-table framing, dispatch authorization, and claims that a kernel globally knows or certifies another agent's abilities.
- `simplicity_durability`: reward small, explicit, deterministic, 100-year durable artifacts that fit small devices. Penalize generic maps, cards, ledgers, bundles, mode matrices, capability maps, service catalogs, feature-shopping wrappers, needless per-record pCID fragmentation, and base-envelope selector-shopping stacks such as `env_pCID`/`sig_pCID`/`payload_pCID`.
- `envelope_discipline`: reward alignment with `DN-jotob` and the current envelope direction: CBOR `grid([42(pCID), payload, ...])`, `pCID` as Protocol CID, protocol-owned slot roles after slot 0, local unknown-pCID behavior, and no universal proof-slot overreach. Treat pCID as a stable protocol-spec identifier, not as a per-message type, payload CID, request ID, operation code, or nonce.
- `kernel_implementation_promises`: reward explicit local kernel implementation promises: apps promise a local kernel which pCIDs they will receive or handle, kernels promise best-effort delivery and local observation records, host assumptions are separated from promises, unsupported pCIDs/features are explicit, pCID adapter mappings are local promises, namespace/reference behavior is voluntary, and the kernel is not treated as an RPC authority, service registry, capability registry, permission issuer, or conformance judge.
- `app_protocol_promise_semantics`: reward higher-layer/app protocols that model storage, computation, send/receive, reciprocal promises, selective sending, promise-as-capability-token behavior, make/break evidence, and local trust updates without command, request/response service, permission, policy-enforcement, or conformance-authority framing. Related promise and observation messages that evolve together should normally be payload `kind` variants under one stable protocol pCID; separate pCIDs are strong only when message families are independently deployable, independently understandable, or separated by a real layer boundary.
- Treat phrases such as "the kernel knows this app supports X", "registered service capability", "authorized dispatch", or "capability table" as weak or invalid unless the design clearly reframes them as local promises made by specific agents plus local observations of kept, broken, refused, or timed-out promises.
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
- `pt_reframe_needed`: technically interesting but drifts into non-PT framing, including RPC-like service dispatch, service-registry vocabulary, capability-table vocabulary, or kernel-known-support claims that could be repaired as local promises and observations; it may survive only as a question-home or rework candidate.
- `pt_invalid`: relies on authority, imposition, global trust, policy enforcement, conformance authority, permission authority, or RPC-style command semantics as a load-bearing design premise; it cannot be promotable.
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

### `proposals/ga-canary-20260525-hozif-anti-rpc-calibration/simulations/SIM-fozid-child-slim-port-commitment-service-outcomes/README.md`

```markdown
# SIM-fozid-child-slim-port-commitment-service-outcomes: Slim port commitment + service outcomes

This child breeds `SIM-fovip-kernel-promise-boundary-port-contract` with `SIM-dalor-grid-envelope-protocol-owned-signature-slot`.

It keeps the strong kernel promise boundary from `SIM-fovip`, but makes three bounded moves:

1. replace the wide self-describing kernel record with a slimmer pCID-selected `kernel_port_commitment` payload that has no inner `record_pcid` field and requires explicit per-surface promiser entries;
2. add a small `service_outcome` payload family for kept/refused/unavailable/broken/timeout observations plus receiver-local next-step effects for pull/keep/advertise choices;
3. use `grid([pCID, payload, signature])` only as an optional signed-carriage profile for exact-byte audit evidence, not as a universal kernel requirement.

## Why this should score higher

- Preserves local trust, explicit promisers, unsupported-pCID refusal, host-assumption separation, and exact-byte evidence.
- Improves app-protocol promise semantics for kept-storage and honest-refusal scenarios.
- Removes inner-selector confusion and trims MCU or split-service complexity by letting small ports publish only the surfaces they actually expose.
- Keeps the envelope contribution narrow and auditable instead of turning it into a universal policy layer.

## Files

- `protocols/kernel-port-commitment.d/specs/kernel-port-commitment-draft.md`
- `protocols/service-outcome.d/specs/service-outcome-draft.md`
- `protocols/signed-grid-carriage-profile.d/specs/protocol-owned-signed-carriage-profile.md`
- `examples/scenario-walkthroughs.md`
```

### `proposals/ga-canary-20260525-hozif-anti-rpc-calibration/simulations/SIM-fozid-child-slim-port-commitment-service-outcomes/QUESTION.md`

```markdown
# Question

Can PromiseGrid improve the kernel porting boundary and the kept/refused service-accounting scenarios at once by publishing a slim pCID-selected `kernel_port_commitment`, pairing it with a small `service_outcome` payload family, and using protocol-owned signed carriage only as an optional exact-byte evidence profile rather than a universal kernel rule?
```

### `proposals/ga-canary-20260525-hozif-anti-rpc-calibration/simulations/SIM-fozid-child-slim-port-commitment-service-outcomes/examples/scenario-walkthroughs.md`

```markdown
# Scenario walkthroughs

## Kernel porting boundary

Example `kernel_port_commitment` payload for a browser/WASM adapter:

```text
[
  'alice-browser-port',
  'browser-wasm',
  [
    ['dispatch', 'alice-adapter', 'dispatches supported app/kernel messages', ['pCID-store', 'pCID-fetch'], 'refuse unsupported pCID'],
    ['storage', 'alice-indexeddb-adapter', 'stores while browser quota and lifecycle allow', ['pCID-store'], 'record unavailable on quota or lifecycle denial'],
    ['network-send', 'alice-fetch-adapter', 'sends requests; does not promise remote delivery', ['pCID-fetch'], 'refuse unsupported roles']
  ],
  ['pCID-store', 'pCID-fetch'],
  ['pCID-key-use'],
  ['browser storage quota', 'browser network stack', 'tab lifecycle'],
  'keep exact request and response bytes for refusal, timeout notes, unsupported-pCID carriage, and broken-promise evidence',
  [['putChunk()', 'pCID-store'], ['getChunk()', 'pCID-fetch']],
  'voluntary namespaces only; list maintaining promisers when exposed',
  'CID-rooted references mount into the local view only',
  'resource state is reconstructed from a selected frontier'
]
```

A sender that wants portable evidence may carry that payload as:

```text
grid([pCID(kernel_port_commitment), payload_bytes, signature_bytes])
```

## Kept storage promise

Bob records a kept outcome locally:

```text
[
  'pCID-store',
  'cid:C',
  'alice-storage',
  'bob',
  'kept',
  'Alice served exact bytes before expiry',
  ['prefer', 'prefer', 'allow', 'kept once; continue pulling'],
  'req-bytes:cid1 resp-bytes:cid2'
]
```

## Refused service

Honest refusal stays distinct from timeout:

```text
[
  'pCID-store',
  'cid:C',
  'alice-storage',
  'bob',
  'refused',
  'policy refusal: missing group authorization',
  ['retry-later', 'neutral', 'do-not-advertise', 'refusal is not corruption'],
  'signed-refusal-bytes:cid3'
]
```

```text
[
  'pCID-store',
  'cid:C',
  'alice-storage',
  'bob',
  'timeout',
  'no answer before local deadline',
  ['decrease', 'neutral', 'neutral', 'absence evidence only'],
  'deadline-log:cid4'
]
```
```

### `proposals/ga-canary-20260525-hozif-anti-rpc-calibration/simulations/SIM-fozid-child-slim-port-commitment-service-outcomes/protocols/kernel-port-commitment.d/specs/kernel-port-commitment-draft.md`

```markdown
# Kernel port commitment draft

> Status: DRAFT. Simulation-local specimen.

## Role

`kernel_port_commitment` is a pCID-selected payload for a port's own implementation promises.
The outer message pCID selects this spec; the payload does not repeat a `record_pcid` field.

This keeps the app/kernel boundary crisp: the kernel port promise record is one higher-layer payload protocol, not a second selector stack inside itself.

## Shape

```text
kernel_port_commitment = [
  port_id,
  profile,
  surface_promises,
  supported_pcids,
  unsupported_pcids,
  host_assumptions,
  evidence_policy,
  adapter_bindings,
  namespace_behavior,
  reference_behavior,
  checkpoint_behavior
]

surface_promise = [
  surface_name,
  promiser_id,
  promise_summary,
  message_pcids,
  refusal_or_nonpromise
]
```

## Rules

- List only exposed surfaces. Omitted surface means no promise.
- Every listed surface names one promiser explicitly. Split local services add more `surface_promise` entries instead of collapsing into one kernel-global actor.
- `supported_pcids` and `unsupported_pcids` cover the app/kernel message families this port handles, preserves, or directly refuses.
- `host_assumptions` names browser, OS, runtime, hardware, or local-service dependencies that the port relies on but does not itself promise.
- `evidence_policy` must say what exact bytes or bounded local records are kept for kept, refused, unavailable, broken, timeout, and unsupported-pCID carriage cases.
- `adapter_bindings` maps local API calls to pCID-selected messages, or marks a call as outside the PromiseGrid promise boundary.
- `namespace_behavior`, `reference_behavior`, and `checkpoint_behavior` stay narrow: if the port exposes them, describe them; if not, omit them or state non-promise directly.

## Minimum credible tiny port

A small MCU or header-only port may publish only:

- one `surface_promise`;
- one supported pCID or bounded exact-byte carriage promise;
- direct unsupported-pCID refusal or non-promise behavior;
- bounded evidence retention;
- explicit non-promises for everything else.

## Promise discipline

- A port promises only its own behavior.
- A host assumption is not silently upgraded into a kernel promise.
- Voluntary namespaces remain reciprocal promises, not global truth.
- CID-rooted references mount into a receiver's local view, not a universal path.
- Resource current state is reconstructed from a selected frontier or checkpoint, not treated as outside-promise truth.
```

### `proposals/ga-canary-20260525-hozif-anti-rpc-calibration/simulations/SIM-fozid-child-slim-port-commitment-service-outcomes/protocols/service-outcome.d/specs/service-outcome-draft.md`

```markdown
# Service outcome draft

> Status: DRAFT. Simulation-local specimen.

## Role

`service_outcome` is a small pCID-owned payload family for local promise-accounting observations.
It records what one observer saw another promiser do about a concrete service request or resource.
It is not a global verdict, global score, or authority record.

## Shape

```text
service_outcome = [
  subject_pcid,
  subject_ref,
  promiser_id,
  observer_id,
  outcome_kind,
  detail,
  local_effect,
  evidence_ref
]

local_effect = [
  next_pull,
  next_keep,
  next_advertise,
  note
]
```

## Outcome kinds

- `kept`
- `refused`
- `unavailable`
- `broken`
- `timeout`
- `corrupt_observation`

## Rules

- `refused` means an explicit honest refusal. It is different from `timeout`, `broken`, or `corrupt_observation`.
- `timeout` means no timely answer was observed.
- `unavailable` means the promiser reported temporary inability without claiming success.
- `broken` means the observer has local evidence that a promise was not kept.
- `local_effect` records only the observer's own next-step hints for pull, keep, or advertise choices. It is not shared authority.
- `evidence_ref` points to exact local bytes, preserved absence evidence, or a bounded local evidence-store reference.
- Higher-layer protocols may subtype `detail`, but they reuse this stable outcome vocabulary so kept/refused/unavailable distinctions stay comparable.

## Intended pressure coverage

This specimen is intentionally small but directly answers two recurring gaps:

- Bob can record a kept storage promise in a reusable way.
- Bob can record an honest refusal differently from timeout, corruption, or silent absence.
```

### `proposals/ga-canary-20260525-hozif-anti-rpc-calibration/simulations/SIM-fozid-child-slim-port-commitment-service-outcomes/protocols/signed-grid-carriage-profile.d/specs/protocol-owned-signed-carriage-profile.md`

```markdown
# Protocol-owned signed carriage profile

> Status: DRAFT. Optional carriage profile for this simulation.

## Role

This profile gives `kernel_port_commitment` and `service_outcome` payloads a compact exact-byte signed transport when a sender wants portable audit evidence.
It is not a universal PromiseGrid requirement and not a claim that every kernel port or local API must always use an outer signature slot.

## Shape

```text
grid([pCID, payload, signature])
```

## Rules

- The signed bytes are canonical `[pCID, payload]` bytes.
- `pCID` is Protocol CID: it names the payload protocol spec and the proof-family rules for `signature`.
- There is no outer `sig_pCID` selector.
- The signature is evidence only of the current sender's own scoped promise about these bytes.
- If a receiver does not support `pCID`, it may preserve exact bytes, but it must not claim to parse the payload or verify proof semantics.
- Other unsigned or differently carried local profiles remain allowed when another protocol or local boundary says so.

## Why kept optional here

Making this carriage profile optional preserves the envelope discipline and auditability gains from the parent envelope specimen without forcing a universal proof slot onto tiny ports, local adapters, or scenario cases that only need local records.
```

### `scenarios/promise-accounting-records-kept-storage-promise/promise-accounting-records-kept-storage-promise.md`

```markdown
# Kept storage promise

## Scenario ID

promise-accounting-records-kept-storage-promise

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-rusap-promise-accounting-records/SCENARIOS.md`
- Source simulation: `SIM-rusap-promise-accounting-records/`
- Source row/title: Kept storage promise
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-rusap-promise-accounting-records/`.

## Setup

Alice promises to store chunk C for a period; Bob later asks and Alice serves it.

## Stimulus

Run the candidate simulation against this source test: What Bob records locally and how that affects future pull / keep / advertise choices.

## Expected Pressure

Specs may need promise-vocabulary sections that name observable promise outcomes.
```

## Required JSON Shape

{"scores":{"scenario_fit":0,"promisegrid_alignment":0,"auditability":0,"evolution_safety":0,"layer_boundary_clarity":0,"failure_handling":0,"implementation_plausibility":0,"promise_vocabulary":0,"simplicity_durability":0,"envelope_discipline":0,"kernel_implementation_promises":0,"app_protocol_promise_semantics":0,"risk_penalty":0},"fitness":{"raw":0,"normalized_0_100":0,"confidence_0_1":0.0},"assessment":{"rationale":"","strengths":[],"weaknesses":[],"risks":[],"open_questions":[],"authority_boundary":"Evidence only; does not settle PromiseGrid design.","pt_gate":{"status":"pt_clean","autonomous_agents":{"status":"pass","note":""},"scoped_intent":{"status":"pass","note":""},"no_promises_for_others":{"status":"pass","note":""},"no_guaranteed_outcomes":{"status":"pass","note":""},"local_trust_assessment":{"status":"pass","note":""},"accept_use_not_obligation":{"status":"pass","note":""},"violations":[]}}}
