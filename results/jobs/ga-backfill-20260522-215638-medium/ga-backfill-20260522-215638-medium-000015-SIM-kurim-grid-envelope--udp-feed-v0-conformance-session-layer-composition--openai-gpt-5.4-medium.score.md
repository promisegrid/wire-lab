# GA Score Cell

Return only JSON with keys `scores`, `fitness`, and `assessment`.
Do not include result identity, source metadata, code fences, or commentary.

## Cell

- Run group ID: `ga-backfill-20260522-215638-medium`
- Cell ID: `ga-backfill-20260522-215638-medium-000015-SIM-kurim-grid-envelope--udp-feed-v0-conformance-session-layer-composition--openai-gpt-5.4-medium`
- Simulation ID: `SIM-kurim-grid-envelope`
- Scenario ID: `udp-feed-v0-conformance-session-layer-composition`
- Model ID: `openai-gpt-5.4-medium`
- Result path: `results/SIM-kurim-grid-envelope/udp-feed-v0-conformance-session-layer-composition/openai-gpt-5.4-medium/20260523-045646.json`

## Rubric

Score each axis from 0 to 5. Higher is better except `risk_penalty`, where 0 is low risk and 5 is severe risk.
Axes: scenario_fit, promisegrid_alignment, auditability, evolution_safety, layer_boundary_clarity, failure_handling, implementation_plausibility, promise_vocabulary, simplicity_durability, risk_penalty.
- `promise_vocabulary`: reward Promise-Theory-correct, promise-first wording. Prefer payload promises such as "Alice promises this payload meets the protocol specification referred to by this pCID." Penalize normative claim cards, conformance bundles, generic profile support claims, port claims, and central trust-ledger framing.
- `simplicity_durability`: reward small, explicit, deterministic, 100-year durable artifacts that fit small devices. Penalize generic maps, cards, ledgers, bundles, and feature-shopping wrappers.
- The runner recomputes `fitness.raw` and `fitness.normalized_0_100` from your axis scores with normal weighting. Use `fitness.confidence_0_1` for your confidence.

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

### `simulations/SIM-kurim-grid-envelope/README.md`

```markdown
# SIM-kurim: Grid envelope

This simulation is now the parent seed and successor-owner lineage for
positional grid-envelope variants. It keeps the `grid([pcid, payload])`
working hypothesis visible as historical seed material without treating any
single envelope as the settled PromiseGrid wire format. Source: `DI-limom`,
`DI-rugig`, `DI-fanah`.

`DI-fanah` split the turn-158 successor path into 24 standalone positional
grid-envelope simulations across encoding, unknown-pCID handling, and signature
placement. The variant index lives in `../README.md`; this parent lineage keeps
the inventory and owner history rather than preferring one child specimen.
```

### `simulations/SIM-kurim-grid-envelope/QUESTION.md`

```markdown
# Question

Which positional grid-envelope variant or variants, if any, survive comparison
across encoding, unknown-pCID behavior, and signature placement?

Secondary questions from the nested-vs-stacked turn-208 research that any
surviving grid-envelope variant must answer:

- What recursion depth budget should nested `grid([pcid, payload])` messages
  enforce?
- Is `pcid` a pure content hash, or does it carry version / routing metadata?
- How are capability references represented when the payload needs more than
  content references?
- What canonical serialization is used for signing and hashing grid messages?
- Does PromiseGrid need onion-routing layers, and if so what must the outer
  `pcid` reveal to routers?

Source: `DI-limom`, `DI-rugig`, `DI-fanah`, `DI-kabuk`.
```

### `simulations/SIM-kurim-grid-envelope/TODO/TODO-tujad-grid-envelope-successor-owner.md`

```markdown
# TODO-tujad: Grid-envelope successor owner

## Prior aliases

None. This TODO is created directly as a sim-local successor owner under
`rusis.10`.

## Status

Closed for the turn-157/turn-158 successor-owner scope. This TODO remains the
parent seed and historical owner for TE-40 transferred rows `UT-157.a`,
`UT-157.c`, `UT-158.f`, and `UT-158.h`; the concrete successor path is now
materialized as 24 standalone positional grid-envelope simulations under
`DI-fanah`. Successor-owner routing into this TODO was locked under `DI-mosor`
in `protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`.
Seed anchors were established earlier under `DI-nijon`.

## Decision Intent Log

ID: DI-joroh
Date: 2026-05-12 08:44:53
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Resolve the turn-157 grid-envelope replay cleanup by capturing the
candidate envelope inventory and `grid([pcid, payload])` working hypothesis in
this sim-local successor owner, without locking a canonical envelope or creating
a protocol tree.
Intent: Turn 157 contains load-bearing design material that should not remain
only in replay memory. Capturing the alternatives and hypothesis here gives the
grid-envelope lineage a concrete owner for the open replay residue while
preserving the fact that the hypothesis remains unproven.
Constraints: Preserve `Env-1` through `Env-5` as candidate inventory. Do not
decide the final PromiseGrid envelope. Do not create `protocols/grid-envelope.d/`
or draft a grid-envelope spec in this cleanup pass. Leave turn-158 protocol/spec
work open under `tujad.3`.
Affects: `simulations/SIM-kurim-grid-envelope/TODO/TODO-tujad-grid-envelope-successor-owner.md`;
`protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`;
`protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`.

ID: DI-fanah
Date: 2026-05-12 09:22:37
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Split the turn-158 grid-envelope protocol/spec successor work into
24 standalone positional grid-envelope variant simulations instead of choosing
one draft or keeping the work only in this parent owner.
Intent: Turn 158 requires candidate envelopes to behave as specimens that can
evolve and compete independently. The grid-envelope successor therefore needs a
variant matrix that exposes the open design axes without making the harness or
the parent simulation prefer one answer.
Constraints: All envelope shapes are positional. Include the encoding axis
`cbor` versus `dag-cbor`; the unknown-pCID axis `opaque`, `hard-reject`, and
`best-effort`; and the signature axis `wrapper-pcid`, `unsigned-v0`,
`mandatory-opaque-bytes`, and `mandatory-sig-pcid-payload`. Do not create
map-shaped envelope variants, choose a winning variant, edit raw replay logs,
edit historical TE bodies, or make grid-envelope canonical.
Affects: `simulations/README.md`;
`simulations/SIM-kurim-grid-envelope/README.md`;
`simulations/SIM-kurim-grid-envelope/QUESTION.md`;
`simulations/SIM-kurim-grid-envelope/TODO/TODO-tujad-grid-envelope-successor-owner.md`;
`protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`;
`protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`simulations/SIM-*-grid-envelope-enc-<cbor|dag-cbor>-unknown-<opaque|hard-reject|best-effort>-sig-<wrapper-pcid|unsigned-v0|mandatory-opaque-bytes|mandatory-sig-pcid-payload>/`.

ID: DI-joman
Date: 2026-05-20 16:02:50
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Add two standalone grid-envelope arity probes alongside the 24
positional variants: one where the first `pcid` defines the outer field count
and field types, and one where the outer layer stays `[pcid_a, payload_a]` but
`pcid_a` defines a nested signed payload structure inside `payload_a`.
Intent: The current grid-envelope matrix does not directly test whether
variable arity belongs in the outer envelope or inside the pCID-defined payload
layer. These probes preserve both hypotheses as independently runnable
specimens without changing the existing positional matrix or declaring a
canonical PromiseGrid envelope.
Constraints: Do not mark either probe as preferred. Keep them standalone under
`simulations/`. Keep the existing 24 positional variants intact. The
layer-pCID nested-signed-payload probe must explicitly preserve the concern
that an unsigned outer envelope relies on transport or local context to
identify the agent promising `payload_a` conforms to `pcid_a`.
Affects: `simulations/SIM-sajar-grid-envelope-variable-arity-pcid-defined-fields/`;
`simulations/SIM-janov-grid-envelope-layer-pcid-nested-signed-payload/`;
`simulations/README.md`;
`DEV-GUIDE-RESOURCES.md`;
`simulations/SIM-kurim-grid-envelope/TODO/TODO-tujad-grid-envelope-successor-owner.md`.

ID: DI-sahiv
Date: 2026-05-20 16:25:11
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Add a standalone non-child grid-envelope simulation named
`SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs` to test whether
Cryptid's Multisig object model can serve as the envelope signature/proof
payload model without prematurely choosing detached versus combined signatures,
outer versus nested signature placement, fixed versus variable arity, pCID
binding, unknown-codec handling, threshold-share aggregation, or verifier
obligations.
Intent: The existing positional, arity, and nested-signature grid-envelope
specimens test where signature bytes may live, but they do not directly pressure
a codec-agnostic signature object that can carry detached or combined payloads,
skippable unknown signing codecs, and threshold-share attributes. A dedicated
Cryptid Multisig specimen preserves that design space as runnable evidence while
keeping the unresolved PromiseGrid envelope questions explicit.
Constraints: Keep the simulation standalone under `simulations/`, not generated,
not a child/proposal sim, and not canonical PromiseGrid wire format. Include
normal local simulation files and a simulation-local `protocols/grid-envelope.d/`
draft. Treat the upstream Multisig source as pre-draft prior art, not as a
normative PromiseGrid dependency. Do not delete or overwrite sibling simulations
or unrelated uncommitted edits.
Affects:
`simulations/SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs/`;
`simulations/README.md`;
`DEV-GUIDE-RESOURCES.md`;
`simulations/SIM-kurim-grid-envelope/TODO/TODO-tujad-grid-envelope-successor-owner.md`.

## Scope

This TODO owns grid-envelope follow-on that was previously parked under
`TODO-kugod` as "until grid-envelope successor exists":

- candidate envelope inventory ownership;
- `grid([pcid, payload])` working-hypothesis prose ownership;
- turn-158 parallel grid-hypothesis TODO ownership;
- concrete successor planning for grid-envelope protocol directory/spec
  work in this lineage.

Anchor seed note:
`simulations/SIM-kurim-grid-envelope/seed/extraction-sources.md`.

## Candidate Envelope Inventory

Turn 157 named these candidate envelope alternatives as inputs to later
grid-envelope and envelope-bakeoff work. This inventory records the candidates;
it does not select a winner. Source: `DI-joroh`.

- `Env-1`: `grid([pcid, payload])`. A two-element CBOR array where the first
  element is a pCID identifying which protocol, handler, or assertion type
  interprets the second element. This is the current working hypothesis. A
  payload may itself be another `grid([pcid, payload])` value if recursion is
  needed.
- `Env-2`: Promise stack of grid frames. A CBOR sequence of
  `grid([pcid, payload])` frames where stack semantics apply at the sequence
  level and the grid shape applies at each frame. This candidate reconciles the
  grid hypothesis with the earlier TE-famar promise-stack work.
- `Env-3`: Bare CBOR with no shared envelope. Each protocol chooses its own
  message shape. This is maximally permissive, but it may leave the harness
  without a shared parser to exercise across candidate transports.
- `Env-4`: Capability-port triplet. A direct `(promiser, assertion, body)`
  structure with no grid indirection. This is closest to the older
  harness-spec `Promise` shape.
- `Env-5`: Tagged union over `Env-1` and `Env-2`. Single-frame messages use
  `grid([pcid, payload])`; multi-frame messages use a stack of grid frames; a
  top-level tag distinguishes the two cases.

## Grid Envelope Working Hypothesis

`grid([pcid, payload])` is the current working hypothesis for a
transport-agnostic message envelope, but turn 157 explicitly says it has not
been proven. This simulation owns that hypothesis as a candidate specimen to be
tested against alternatives, not as a settled harness rule. Source: `DI-joroh`.

The harness may use this inventory to construct later bakeoffs, but any final
canonical-envelope decision still needs its own TE/DF/DI path. Turn 158's
apparatus-vs-specimen correction remains controlling: the harness compares
candidate envelopes rather than declaring this one canonical in advance.

## Successor Variant Simulations

`DI-fanah` closes the concrete successor path by creating 24 standalone
positional grid-envelope simulations. Each variant carries its own local draft
under `protocols/grid-envelope.d/` so future evolution can compare specimens
without relying on this parent lineage as a shared protocol bundle.

| Simulation | Encoding | Unknown-pCID policy | Signature policy |
|---|---|---|---|
| `../SIM-mahih-grid-envelope-enc-cbor-unknown-opaque-sig-wrapper-pcid/` | CBOR | Opaque store/forward | Wrapper pCID |
| `../SIM-gasus-grid-envelope-enc-cbor-unknown-opaque-sig-unsigned-v0/` | CBOR | Opaque store/forward | Unsigned v0 |
| `../SIM-vutar-grid-envelope-enc-cbor-unknown-opaque-sig-mandatory-opaque-bytes/` | CBOR | Opaque store/forward | Mandatory opaque bytes |
| `../SIM-vamaz-grid-envelope-enc-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload/` | CBOR | Opaque store/forward | Mandatory signature pCID + payload |
| `../SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid/` | CBOR | Hard reject | Wrapper pCID |
| `../SIM-gazan-grid-envelope-enc-cbor-unknown-hard-reject-sig-unsigned-v0/` | CBOR | Hard reject | Unsigned v0 |
| `../SIM-hupir-grid-envelope-enc-cbor-unknown-hard-reject-sig-mandatory-opaque-bytes/` | CBOR | Hard reject | Mandatory opaque bytes |
| `../SIM-kovis-grid-envelope-enc-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload/` | CBOR | Hard reject | Mandatory signature pCID + payload |
| `../SIM-vivus-grid-envelope-enc-cbor-unknown-best-effort-sig-wrapper-pcid/` | CBOR | Best-effort inspection | Wrapper pCID |
| `../SIM-fonig-grid-envelope-enc-cbor-unknown-best-effort-sig-unsigned-v0/` | CBOR | Best-effort inspection | Unsigned v0 |
| `../SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/` | CBOR | Best-effort inspection | Mandatory opaque bytes |
| `../SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/` | CBOR | Best-effort inspection | Mandatory signature pCID + payload |
| `../SIM-gojot-grid-envelope-enc-dag-cbor-unknown-opaque-sig-wrapper-pcid/` | DAG-CBOR | Opaque store/forward | Wrapper pCID |
| `../SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0/` | DAG-CBOR | Opaque store/forward | Unsigned v0 |
| `../SIM-ruzil-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-opaque-bytes/` | DAG-CBOR | Opaque store/forward | Mandatory opaque bytes |
| `../SIM-riliz-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload/` | DAG-CBOR | Opaque store/forward | Mandatory signature pCID + payload |
| `../SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/` | DAG-CBOR | Hard reject | Wrapper pCID |
| `../SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0/` | DAG-CBOR | Hard reject | Unsigned v0 |
| `../SIM-tohol-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-opaque-bytes/` | DAG-CBOR | Hard reject | Mandatory opaque bytes |
| `../SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload/` | DAG-CBOR | Hard reject | Mandatory signature pCID + payload |
| `../SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/` | DAG-CBOR | Best-effort inspection | Wrapper pCID |
| `../SIM-zifik-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-unsigned-v0/` | DAG-CBOR | Best-effort inspection | Unsigned v0 |
| `../SIM-fonol-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/` | DAG-CBOR | Best-effort inspection | Mandatory opaque bytes |
| `../SIM-rakir-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/` | DAG-CBOR | Best-effort inspection | Mandatory signature pCID + payload |

## Arity and Nested-Signature Probe Simulations

`DI-joman` adds two arity-focused probes that are intentionally outside the
24-row positional matrix. They test whether arity should be a property of the
outer envelope or of a pCID-defined nested payload layer.

| Simulation | Probe question |
|---|---|
| `../SIM-sajar-grid-envelope-variable-arity-pcid-defined-fields/` | Can the first `pcid` safely define how many outer fields follow it and what each field means? |
| `../SIM-janov-grid-envelope-layer-pcid-nested-signed-payload/` | Can a shared layer/network pCID define a parseable signed nested payload while leaving the outer envelope unsigned? |

## Signature/Proof Object Probe Simulations

`DI-sahiv` adds a standalone non-child probe for using Cryptid's Multisig object
model as the bytes carried by a grid-envelope signature/proof slot or by a
payload protocol's nested proof. The probe keeps detached versus combined
signatures, outer versus nested placement, variable arity, pCID binding,
unknown-codec handling, threshold shares, and verifier obligations open for
comparison.

| Simulation | Probe question |
|---|---|
| `../SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs/` | Can Cryptid Multisig carry grid-envelope signature/proof bytes without prematurely settling PromiseGrid signature design choices? |

## Subtasks

- [x] tujad.1 Materialize the candidate envelope inventory owner record
  for `UT-157.a`.
- [x] tujad.2 Materialize the `grid([pcid, payload])`
  working-hypothesis owner record for `UT-157.c`.
- [x] tujad.3 Define and track the concrete successor path for
  grid-envelope protocol directory/spec work for `UT-158.f`. Closed by
  `DI-fanah` via the 24 standalone positional successor simulations.
- [x] tujad.4 Back-link resulting decisions and artifacts to
  `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`.
- [x] tujad.5 Record that this TODO satisfies `UT-158.h` as the
  turn-158 parallel grid-hypothesis TODO; `tujad.3` is now closed by
  the `DI-fanah` successor split.
```

### `simulations/SIM-kurim-grid-envelope/TODO/TODO.md`

```markdown
# TODO queue: grid-envelope

Per `rusis.10`, this queue tracks grid-envelope lineage follow-on work
that was still parked in rooted TE-40 owner rows.

## Index

| Handle | Mint date | Title | Prior alias |
|---|---|---|---|
| [TODO-tujad](TODO-tujad-grid-envelope-successor-owner.md) | 2026-05-11 | Grid-envelope successor owner for TE-40 transferred-open rows | none |
```

### `simulations/SIM-kurim-grid-envelope/seed/extraction-sources.md`

```markdown
# Grid-envelope extraction sources

The `grid-envelope` lineage is seeded from rooted carve-out decisions and the
current `grid <pcid>` carrier text still living inside the `group-session`
lineage draft. `rusis.8` records those anchors without extracting or rewriting
protocol prose. Source: `DI-nijon`, `DI-rugig`.

Primary anchors for the later extraction/owner-routing pass:

- `docs/thought-experiments/TE-havib-apparatus-vs-specimen-carve-out.md`
  (`DF-36.3`, lock `Alt-3.A`).
- `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`
  (`UT-157.a`, `UT-157.c`, `UT-158.f`, and open task `kugod.7`).
- `simulations/SIM-rakot-group-session/protocols/group-session.d/specs/group-session-draft.md`
  (`## Sources`, `§4 Message envelope`, and `§4.1 Carrier line`).

Deferred in this pass:

- creating `protocols/grid-envelope.d/` under this sim;
- moving or rewriting `grid <pcid>` envelope prose out of the
  `group-session` lineage draft;
- editing rooted harness docs or rooted TODO owner routing.

This note is provenance only. It does not declare any source file as the live
authority for this sim and does not create a protocol tree yet. Source:
`DI-nijon`.
```

### `scenarios/udp-feed-v0-conformance-session-layer-composition/udp-feed-v0-conformance-session-layer-composition.md`

```markdown
# Session-layer composition

## Scenario ID

udp-feed-v0-conformance-session-layer-composition

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-kuful-udp-feed-v0-conformance/SCENARIOS.md`
- Source simulation: `SIM-kuful-udp-feed-v0-conformance/`
- Source row/title: Session-layer composition
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-kuful-udp-feed-v0-conformance/`.

## Setup

A minimal group/session message rides above UDP-feed v0.

## Stimulus

Run the candidate simulation against this source test: Whether UDP-feed's API is sufficient for the next layer without leaking binding details.

## Expected Pressure

If composition is required, TODO-jodon's done criteria must include more than UDP round trip.
```

## Required JSON Shape

{"scores":{"scenario_fit":0,"promisegrid_alignment":0,"auditability":0,"evolution_safety":0,"layer_boundary_clarity":0,"failure_handling":0,"implementation_plausibility":0,"promise_vocabulary":0,"simplicity_durability":0,"risk_penalty":0},"fitness":{"raw":0,"normalized_0_100":0,"confidence_0_1":0.0},"assessment":{"rationale":"","strengths":[],"weaknesses":[],"risks":[],"open_questions":[],"authority_boundary":"Evidence only; does not settle PromiseGrid design."}}
