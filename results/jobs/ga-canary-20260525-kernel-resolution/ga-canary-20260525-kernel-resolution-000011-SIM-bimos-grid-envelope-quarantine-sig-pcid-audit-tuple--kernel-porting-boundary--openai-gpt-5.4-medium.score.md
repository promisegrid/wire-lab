# GA Score Cell

Return only JSON with keys `scores`, `fitness`, and `assessment`.
Do not include result identity, source metadata, code fences, or commentary.

## Cell

- Run group ID: `ga-canary-20260525-kernel-resolution`
- Cell ID: `ga-canary-20260525-kernel-resolution-000011-SIM-bimos-grid-envelope-quarantine-sig-pcid-audit-tuple--kernel-porting-boundary--openai-gpt-5.4-medium`
- Simulation ID: `SIM-bimos-grid-envelope-quarantine-sig-pcid-audit-tuple`
- Scenario ID: `kernel-porting-boundary`
- Model ID: `openai-gpt-5.4-medium`
- Result path: `results/SIM-bimos-grid-envelope-quarantine-sig-pcid-audit-tuple/kernel-porting-boundary/openai-gpt-5.4-medium/20260525-203836.json`

## Rubric

Score each axis from 0 to 5. Higher is better except `risk_penalty`, where 0 is low risk and 5 is severe risk.
Axes: scenario_fit, promisegrid_alignment, auditability, evolution_safety, layer_boundary_clarity, failure_handling, implementation_plausibility, promise_vocabulary, simplicity_durability, risk_penalty.
- Score the candidate at its claimed layer. Do not penalize an envelope-layer design for leaving promise accounting, storage semantics, computation semantics, or application-specific trust updates to the payload protocol when the layer boundary is explicit.
- `promise_vocabulary`: reward Promise-Theory-correct, promise-first wording at the candidate's layer. For envelope sims, reward signed pCID-specific promises such as "Alice promises these payload bytes are shaped according to the protocol specification named by this pCID." Penalize normative claim cards, generic claim headers/cards, conformance bundles, generic profile support claims, port claims, universal statement capsules, and central trust-ledger framing.
- `simplicity_durability`: reward small, explicit, deterministic, 100-year durable artifacts that fit small devices. Penalize generic maps, cards, ledgers, bundles, feature-shopping wrappers, and base-envelope selector-shopping stacks such as `env_pCID`/`sig_pCID`/`payload_pCID`.
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

Required score-axis checklist: `scenario_fit`, `promisegrid_alignment`, `auditability`, `evolution_safety`, `layer_boundary_clarity`, `failure_handling`, `implementation_plausibility`, `promise_vocabulary`, `simplicity_durability`, `risk_penalty`.
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

### `simulations/SIM-bimos-grid-envelope-quarantine-sig-pcid-audit-tuple/README.md`

```markdown
# SIM-bimos-grid-envelope-quarantine-sig-pcid-audit-tuple: Grid-envelope variant

This simulation is a standalone positional grid-envelope specimen bred from exactly two parent variants.

It combines the parent strengths into one candidate:

- explicit signature-protocol dispatch via a 4-slot envelope;
- strict no-speculative-parse rejection for unknown protocols;
- exact-byte preservation of rejected unknown envelopes as quarantined evidence.

Variant under test: `enc-dag-cbor` / `unknown-quarantine` / `sig-mandatory-sig-pcid-payload`.

This promoted specimen is expected to improve on the parents under mixed-version
IoT fleets, route-evidence exchange, and future CAS object-family rollout because
unknown envelopes are not silently accepted, but they also are not needlessly
discarded.

This specimen was promoted from review-stage child proposal
`SIM-bimos-child-grid-envelope-quarantine-sig-pcid-audit-tuple` from
`ga-canary-20260520-221953` under `DI-dipid`. The ignored proposal artifacts
remain local raw evidence; this directory is the canonical non-child simulation
home.

The local draft spec is `protocols/grid-envelope.d/specs/grid-envelope-draft.md`.
```

### `simulations/SIM-bimos-grid-envelope-quarantine-sig-pcid-audit-tuple/QUESTION.md`

```markdown
# Question

Does a positional grid envelope using `enc-dag-cbor`, `unknown-quarantine`, and `sig-mandatory-sig-pcid-payload` satisfy the sampled wire-lab scenarios better than both parents by keeping strict acceptance boundaries while allowing exact-byte archival and relay of rejected unknown envelopes as auditable evidence?
```

### `simulations/SIM-bimos-grid-envelope-quarantine-sig-pcid-audit-tuple/protocols/grid-envelope.d/CHANGELOG.md`

```markdown
# CHANGELOG: grid-envelope

A-side CHANGELOG (per TE-liviv) for this simulation-local `grid-envelope`
protocol specimen.

## 2026-05-20

- Promoted specimen
  `SIM-bimos-grid-envelope-quarantine-sig-pcid-audit-tuple` from review-stage
  child proposal `SIM-bimos-child-grid-envelope-quarantine-sig-pcid-audit-tuple`
  under `DI-dipid`.
- Bred from:
  - `SIM-ruzil-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-opaque-bytes`
  - `SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload`
- Introduced `unknown-quarantine` policy: unknown protocol instances are rejected for validity but may be archived and relayed as explicitly rejected opaque evidence.
- Standardized on explicit `sig_pcid` + `sig_payload` signature dispatch.
```

### `simulations/SIM-bimos-grid-envelope-quarantine-sig-pcid-audit-tuple/protocols/grid-envelope.d/specs/grid-envelope-draft.md`

```markdown
# Grid Envelope Variant Spec (DRAFT)

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `enc-dag-cbor` / `unknown-quarantine` / `sig-mandatory-sig-pcid-payload`.
> Bred from two simulation-local parent specimens. Source lineage: `DI-fanah` plus GA child synthesis promoted under `DI-dipid`.

## Purpose

This spec defines one full positional grid-envelope candidate for wire-lab
comparison. It is a specimen inside `SIM-bimos-grid-envelope-quarantine-sig-pcid-audit-tuple`, not a harness rule
and not the canonical PromiseGrid wire format.

The design goal is to keep the stricter parent's explicit signature dispatch and
hard validity boundary, while importing the other parent's better mixed-version
evolution behavior by preserving rejected unknown envelopes as exact opaque
evidence.

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

This variant encodes the envelope as DAG-CBOR-compatible positional arrays.
`pcid` and `sig_pcid` are DAG-CBOR Link values. `payload` and `sig_payload` are
byte strings. The envelope remains positional: no map/object envelope fields are
introduced.

The canonical bytes for hashing, archival, signing input, and evidence relay are
the DAG-CBOR bytes of the exact positional array under this spec.

## Unknown Protocol Policy: `unknown-quarantine`

If a receiver lacks a handler for `pcid` or `sig_pcid`, the envelope enters
**quarantine**.

Quarantine means all of the following:

- The receiver MUST NOT accept the envelope as a valid interpreted message under
  this variant.
- The receiver MUST NOT parse `payload` speculatively without the handler named
  by `pcid`.
- The receiver MUST NOT treat `sig_payload` as verified without the handler named
  by `sig_pcid`.
- The receiver MUST preserve the exact canonical bytes if it stores the object.
- The receiver MAY archive the bytes together with local rejection metadata.
- The receiver MAY relay the exact bytes onward only as **quarantined opaque
  evidence**, and MUST preserve the fact that local interpretation failed.

Relay of quarantined evidence is not acceptance. A relaying peer is forwarding
raw evidence for later interpretation by a different peer, version, or archive.

This policy is intended to avoid the main failure modes seen in the parents:

- unlike `unknown-hard-reject`, useful future evidence is not forced to vanish;
- unlike bare `unknown-opaque`, unknown content is not mistaken for locally valid
  understood traffic.

## Signature and Authorship Policy

The third and fourth positional slots are mandatory.

- `sig_pcid` identifies the signature or proof protocol.
- `sig_payload` is opaque bytes interpreted by that signature protocol.

Unless `sig_pcid` publishes stricter rules, the signature or proof covers the
canonical unsigned prefix:

```text
[pcid, payload]
```

Signing `pcid` together with `payload` reduces type-confusion risk and makes the
dispatch rule part of the auditable signed evidence.

This envelope layer enforces presence, position, and byte-shape. Signer
identity, authority, rotation, revocation, delegation, and freshness policy are
still defined by the protocol ecosystems named by `pcid` and `sig_pcid`.

## Local Audit Tuple

To improve long-term auditability without adding more wire-level fields,
implementations SHOULD retain a local audit tuple for each received envelope:

- exact canonical envelope bytes;
- observed receive time and source context;
- extracted `pcid` and `sig_pcid` bytes;
- local interpretation result: accepted, rejected, or quarantined;
- local verification result and verifier version, if verification was possible;
- any locally known content-addressed retrieval path for the specs named by
  `pcid` and `sig_pcid`.

This tuple is local evidence, not a global registry requirement and not a new
consensus rule.

## Layering-Test Behavior

This variant answers the layering pressure as follows:

- Ordering disagreements are handled by the protocol named by `pcid`; the grid
  envelope preserves only dispatch identity and exact bytes.
- Forwarding or relay evidence can be represented by wrapper envelopes, by the
  payload protocol, or by quarantine relay of exact rejected bytes.
- External or content-addressed body references live inside `payload` under the
  protocol named by `pcid`.
- Incompatible interpretation rules fail visibly at the `pcid` or `sig_pcid`
  boundary, but evidence can still survive for later audit or upgraded peers.

## Scenario Pressure Notes

### IoT fleet maintenance

This variant does not define device identity, maintenance history, firmware
approval, telemetry, or access-control objects by itself. It does improve the
envelope substrate for long-lived mixed fleets by allowing unknown future update
or telemetry envelopes to be retained and escalated instead of dropped, while
still preventing accidental acceptance.

### BGP routing

This variant does not define route, withdrawal, freshness, or leak-policy
semantics by itself. It does improve survivability of unfamiliar route evidence
across sparse peers: newer route objects can be quarantined and relayed for
later inspection rather than black-holed.

### CAS application object families

The first binding rule remains extensible because `pcid` identifies the payload
family without requiring reinterpretation of old bytes. Unknown future families
can be preserved as quarantined evidence, which softens rollout brittleness while
keeping local validity boundaries explicit.

## Non-Goals

This draft does not:

- declare a winning envelope;
- define a central pCID registry;
- freeze a final PromiseGrid signing scheme;
- define application payloads for IoT, routing, or CAS families;
- claim that quarantine relay settles trust or freshness.

## Freeze Gate

This draft can freeze only after at least one simulation run compares it against
sibling positional grid-envelope variants and a maintainer signs a merge/freeze
promise for this specific specimen.
```

### `scenarios/kernel-porting-boundary/kernel-porting-boundary.md`

```markdown
# Kernel Porting Boundary

## Scenario ID

kernel-porting-boundary

## Source / Provenance

- Source type: new harness scenario
- Source path: `/home/stevegt/lab/promisegrid-dev-guide/FEEDBACK.md`
- Source row/title: `FB-vitih`, `FB-mulum`, and `FB-potin`
- Source DI / TODO / TE: `DI-ragaz`; `TODO-rozas`; `DR-davod`

## Purpose

Exercise the developer-facing boundary for a first real PromiseGrid port while
the kernel/runtime terminology remains unsettled. Source: `DI-ragaz`;
`DI-fidot`.

## Setup

Alice wants to port PromiseGrid infrastructure to a new host environment. Bob
offers a native service, Carol offers a browser/WASM adapter, Dave offers an
MCU/header-only profile, Ellen offers split local services, and Mallory claims
that copying the wire-lab harness or one microkernel shape is the porting
target. The available specs are drafts, and only some future frozen pCIDs will
become obligations.

## Stimulus

Alice writes a first porting plan and a port promise record. The plan must say
what it implements now, which pCID-selected messages it exposes, what draft
evidence it follows, what it refuses or cannot promise, what host assumptions it
depends on, what evidence it records, and what it defers until `DR-davod`
decides the guide-facing boundary.

## Expected Pressure

The candidate design must separate harness apparatus from porting target,
identify which binding/session/message/CAS/runtime obligations are provisional,
and preserve a clear path to future frozen-spec implementation promises.

It must also show whether:

- app/kernel operations are pCID-selected `grid([42(pCID), payload, ...])`
  messages, even when local APIs wrap them;
- storage, compute, network send/receive, key use, device access, lifecycle,
  dispatch, refusal, receipt, evidence, namespace, reference, and checkpoint
  operations each name their promiser;
- host/runtime assumptions are separate from PromiseGrid promises;
- unsupported pCIDs and unsupported roles are direct refusals or non-promises;
- exact bytes are preserved where needed for proof, replay, unsupported-pCID
  carriage, or broken-promise evidence;
- voluntary group namespaces are reciprocal promises, not global truth;
- CID-rooted references let Alice share a resource that Bob maps into Bob's own
  local view;
- file/resource current state can be reconstructed as a checkpoint over a
  selected promise-log frontier.

## Scenario Variants

- **Native node:** Bob's service promises storage, dispatch, networking, keys,
  lifecycle, and evidence, but must name every pCID it supports and every role it
  does not promise.
- **Browser/WASM:** Carol's adapter depends on browser storage, network, clocks,
  and lifecycle. Carol can promise adapter behavior, but not that the browser
  will keep promises the browser has not made.
- **Mobile sandbox:** Dave promises work only while the OS offers foreground or
  background execution. Unavailable background work must be recorded as an
  unavailable promise, refusal, or host assumption rather than hidden success.
- **MCU/header-only:** Erin supports one actuator pCID, one bounded evidence
  store, and no general namespace service. The port is credible only if it says
  exactly what it cannot promise.
- **Split local services:** Frank separates dispatch, storage, keys, networking,
  and evidence among local promisers. The record must show which service makes
  each promise and how evidence moves between them.
- **Voluntary namespace:** Alice, Bob, and Carol maintain `/project/report` by
  reciprocal namespace promises. Mallory's lookalike namespace is rejected unless
  a local agent trusts Mallory's promise history.
- **CID-rooted reference:** Alice sends Bob a reference rooted at a CID with
  pCID, selector/path, frontier, promise body, and evidence. Bob chooses whether
  and where to mount it.
- **Checkpointed resource:** Grace reconstructs a file from old pCID specs,
  promise/event logs, and a selected frontier. A different branch may produce a
  different current file because it selects a different promise history.

## Scenario-Specific Evaluation Questions

- Should the guide say kernel, runtime, dispatcher, handler host, or library?
- What is the minimum viable porting target before final freeze?
- Which K1-K5 ingress, feed, CAS, session, and app-layer details should remain
  blocked versus provisional orientation?
- Does the candidate preserve local trust, autonomous promisers, and
  make/break evidence?
- Does the candidate avoid global permission, global namespace authority, and
  universal process-shape assumptions?
- Are local APIs faithful adapters over pCID-selected messages, or do they hide
  the promise boundary?
```

## Required JSON Shape

{"scores":{"scenario_fit":0,"promisegrid_alignment":0,"auditability":0,"evolution_safety":0,"layer_boundary_clarity":0,"failure_handling":0,"implementation_plausibility":0,"promise_vocabulary":0,"simplicity_durability":0,"risk_penalty":0},"fitness":{"raw":0,"normalized_0_100":0,"confidence_0_1":0.0},"assessment":{"rationale":"","strengths":[],"weaknesses":[],"risks":[],"open_questions":[],"authority_boundary":"Evidence only; does not settle PromiseGrid design.","pt_gate":{"status":"pt_clean","autonomous_agents":{"status":"pass","note":""},"scoped_intent":{"status":"pass","note":""},"no_promises_for_others":{"status":"pass","note":""},"no_guaranteed_outcomes":{"status":"pass","note":""},"local_trust_assessment":{"status":"pass","note":""},"accept_use_not_obligation":{"status":"pass","note":""},"violations":[]}}}
