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
