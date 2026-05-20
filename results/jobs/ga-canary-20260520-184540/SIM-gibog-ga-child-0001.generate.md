# GA Child Generation

Return only JSON with keys `child_id`, `design_delta_summary`, and `files`.
Each file path must be relative to the child simulation root. Include `README.md` and `QUESTION.md`.

- Run group ID: `ga-canary-20260520-184540`
- Child ID: `SIM-gibog-ga-child-0001`
- Child path: `simulations/SIM-gibog-ga-child-0001/`
- Operation: `mutation`
- Parent IDs: `SIM-hugoj-cas-usenetlike-gitlike`

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

### `simulations/SIM-hugoj-cas-usenetlike-gitlike/README.md`

```markdown
# SIM-hugoj: CAS / Usenet-like / Git-like lineage

This simulation explores whether part of PromiseGrid should be
understood as a content-addressed-Usenet successor line with git-like
hashed and hash-chained objects and transport separation. It is a
broad design exploration, not a frozen protocol claim and not a claim
that PromiseGrid simply *is* Usenet or git. Source: `DI-pijun`.

`group-session` is the current worked specimen inside this simulation, not the
identity of the simulation itself. The point of the simulation is to ask which
historical invariants age well enough to carry forward into PromiseGrid, which
ones should be adapted, and which ones should be rejected. Source: `DI-pijun`.

## Question

Which parts of the Usenet / FidoNet / email lineage, when combined with
content-addressed storage and git-like transport separation, describe a useful
PromiseGrid design line rather than a misleading analogy? Source: `DI-pijun`.

## Current specimen

The current specimen is `group-session`, because it already exhibits several of
the relevant traits: append-only message growth, DAG/thread structure,
content-addressed message identity, and multiple possible delivery substrates.
That makes it a useful worked example, but not the whole design space. A later
PromiseGrid protocol may preserve some of these traits while diverging sharply
from `group-session` in governance, payload shape, replication, or feed
mechanics. Source: `DI-pijun`.

## What maps cleanly from precedent

- Message identity that stays stable across delivery substrates.
- Store-and-forward replication among named peers or sites.
- Per-instance declaration of how a site exchanges messages with peers.
- Thread/DAG semantics that survive relay across heterogeneous wires.
- Separation between protocol meaning and the substrate that carries bytes.

These are the parts of the turn-173 precedent survey that look structurally
durable across email, Usenet, FidoNet, and modern git-like or gRPC/libp2p-like
systems. See `docs/research/historical-networks-20260503.md`. Source:
`DI-pijun`.

## What does not map cleanly

- Usenet control-message conventions do not automatically become PromiseGrid
  governance rules.
- Git's object model and transport behavior are precedent, not identity; using
  git as a current substrate does not imply PromiseGrid should clone git's
  object types or workflow.
- A content-addressed-Usenet framing does not by itself settle feed naming,
  site manifests, sparse-CAS shape, moderation, or freeze behavior.
- The simulation does not assume that every future PromiseGrid protocol is a
  `group-session` variant.

These are precisely the places where the analogy is useful only if it remains
explicitly exploratory. Source: `DI-pijun`.

## Provenance and authority

The historical grounding for this simulation comes from
`docs/research/historical-networks-20260503.md`, especially its Usenet,
FidoNet, email, and git-adjacent precedent analysis. The immediate trigger for
filing this simulation is the turn-173 replay recovery in
`protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`
and the unresolved follow-on notes in
`protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`.

Nothing in this directory overrides later DRs, DIs, frozen specs, or guide
prose. This simulation is a bounded design workspace for testing one broad
lineage claim and keeping the claim visible until it either graduates or is
rejected. Source: `DI-pijun`.
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

## Existing Fitness Evidence From This Run

### `results/SIM-hugoj-cas-usenetlike-gitlike/insurance-claims/openai-gpt-5.4-xhigh/20260520-184540.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-hugoj-cas-usenetlike-gitlike-insurance-claims-openai-gpt-5.4-xhigh-20260520-184540",
  "run_group_id": "ga-canary-20260520-184540",
  "cell_id": "ga-canary-20260520-184540-000001-SIM-hugoj-cas-usenetlike-gitlike--insurance-claims--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-hugoj-cas-usenetlike-gitlike",
  "scenario_id": "insurance-claims",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-184540",
  "result_path": "results/SIM-hugoj-cas-usenetlike-gitlike/insurance-claims/openai-gpt-5.4-xhigh/20260520-184540.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_89babdf3de624b07a7d098ac29acda77",
    "response_id": "resp_09a5c6dd802c11e6006a0e01567b1c8199aee7707f0230806b",
    "input_tokens": 5776,
    "output_tokens": 7660,
    "cost_usd": 0.117348
  },
  "source": {
    "repo_commit": "bf8f3ec6b37af889c7f7b101cdf6c7b6ee00edaa",
    "sim_path": "simulations/SIM-hugoj-cas-usenetlike-gitlike/",
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
        "path": "simulations/SIM-hugoj-cas-usenetlike-gitlike/README.md",
        "sha256": "0506d32e30f1fe4903d1c0c1f2e6a7f25a3c3a438d6cfa0782d3e99c9448bdb8"
      },
      {
        "path": "scenarios/insurance-claims/insurance-claims.md",
        "sha256": "2a9fe6a4c9a84a35f032d93ae9c36d2adec30a8175f821a34cf8592c73810ab1"
      }
    ],
    "simulation_tree_hash": "144951c8a39cbb610ac9d305fe823308c60e01e271c2f3b222a32276d4a196ef"
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
    "failure_handling": 2,
    "implementation_plausibility": 2,
    "risk_penalty": 4
  },
  "fitness": {
    "raw": 17,
    "normalized_0_100": 49,
    "confidence_0_1": 0.64
  },
  "assessment": {
    "rationale": "The candidate is a plausible substrate for contested evidence exchange: content-addressed objects, append-only history, and transport separation help with durable references, sparse knowledge, and partitioned delivery. It is only a partial fit for insurance claims, though, because the simulation does not define claim-specific authority, privacy, fraud, payment, or appeal semantics; it mainly addresses the evidence-log layer rather than the full application workflow.",
    "strengths": [
      "Stable content-addressed and hash-linked records support long-lived evidence references and tamper-evident audit trails.",
      "Store-and-forward replication and substrate independence fit sparse knowledge and delayed or partitioned evidence exchange.",
      "The simulation clearly distinguishes historical precedent from protocol identity, which keeps layer boundaries explicit."
    ],
    "weaknesses": [
      "It does not define insurance-specific objects or promises for adjuster authority, payment authorization, denials, or appeals.",
      "Failure handling is mostly transport-level; dispute resolution, retry, downgrade, and escalation rules are not worked through.",
      "The current worked specimen is a general message/thread model, not a claims-processing workflow."
    ],
    "risks": [
      "Applying a Usenet-like or broadly replicated CAS pattern directly to insurance claims could expose sensitive or regulated data without a clear privacy boundary.",
      "The design may under-specify institutional authority, compliance, and policy-driven decision rights that insurance processing depends on.",
      "Tamper-evident storage helps provenance but does not by itself resolve fraud, conflicting testimony, or evidentiary truth."
    ],
    "open_questions": [
      "How are adjuster, reviewer, and payer authority claims represented, delegated, and revoked without hidden central-state assumptions?",
      "What selective-disclosure, access-control, retention, or redaction mechanism prevents uncontrolled spread of sensitive claim evidence?",
      "How are payments, denials, and appeal transitions recorded as peer-local, auditable promise-accounting events rather than ad hoc messages?",
      "What local decision rules should peers use when required evidence is stale, delayed, contradictory, or missing?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

### `results/SIM-hugoj-cas-usenetlike-gitlike/state-governance/openai-gpt-5.4-xhigh/20260520-184540.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-hugoj-cas-usenetlike-gitlike-state-governance-openai-gpt-5.4-xhigh-20260520-184540",
  "run_group_id": "ga-canary-20260520-184540",
  "cell_id": "ga-canary-20260520-184540-000002-SIM-hugoj-cas-usenetlike-gitlike--state-governance--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-hugoj-cas-usenetlike-gitlike",
  "scenario_id": "state-governance",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-184540",
  "result_path": "results/SIM-hugoj-cas-usenetlike-gitlike/state-governance/openai-gpt-5.4-xhigh/20260520-184540.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_33c5b8e9f00a405c9523b32fc56ac1d4",
    "response_id": "resp_0bb65ede1f22b08e006a0e01569b248198a4acc5bcfe659526",
    "input_tokens": 5782,
    "output_tokens": 6276,
    "cost_usd": 0.0979825
  },
  "source": {
    "repo_commit": "bf8f3ec6b37af889c7f7b101cdf6c7b6ee00edaa",
    "sim_path": "simulations/SIM-hugoj-cas-usenetlike-gitlike/",
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
        "path": "simulations/SIM-hugoj-cas-usenetlike-gitlike/README.md",
        "sha256": "0506d32e30f1fe4903d1c0c1f2e6a7f25a3c3a438d6cfa0782d3e99c9448bdb8"
      },
      {
        "path": "scenarios/state-governance/state-governance.md",
        "sha256": "6cbfa949fc3eed42fd81bb6820d73ec96ac9c71669b3117d6e46673ddba38360"
      }
    ],
    "simulation_tree_hash": "144951c8a39cbb610ac9d305fe823308c60e01e271c2f3b222a32276d4a196ef"
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
    "auditability": 3,
    "evolution_safety": 4,
    "layer_boundary_clarity": 5,
    "failure_handling": 2,
    "implementation_plausibility": 3,
    "risk_penalty": 4
  },
  "fitness": {
    "raw": 17,
    "normalized_0_100": 49,
    "confidence_0_1": 0.77
  },
  "assessment": {
    "rationale": "This simulation is a plausible lower-layer ingredient for the scenario because it emphasizes content-addressed, hash-chained, append-only messages, store-and-forward replication, sparse knowledge, and transport separation. Those traits help with durability, local evidence retention, and relay under partitions. But as scored against state-governance pressure, it remains far from sufficient: it does not yet specify the identity claims, naming, jurisdiction/regional-authority rules, promise-accounting records, or accept/retry/downgrade/escalate policies needed for contested government records. The README is notably strong at keeping substrate, specimen, and precedent boundaries explicit, which improves evolution safety and reduces analogy drift, but the actual governance layer is still missing.",
    "strengths": [
      "Stable content-addressed identities plus append-only/hash-chained history support later audit of disputed records.",
      "Store-and-forward replication and transport separation fit sparse knowledge, partitioned operation, and long-lived evolution better than a centralized registry model.",
      "The simulation explicitly distinguishes precedent from identity and substrate from governance, which gives it unusually clear layer boundaries.",
      "Named peers or sites and per-instance exchange declarations could serve as a base for cross-agency evidence relay."
    ],
    "weaknesses": [
      "It does not define regional authority, delegation, jurisdiction, or cross-agency conflict-resolution semantics.",
      "The scenario's required identity claims, names, signatures, and peer-local promise-accounting records are not specified.",
      "Failure handling is mostly substrate-level; there is no clear policy for when a peer should accept, retry, downgrade, or escalate contested state.",
      "No domain model is provided for benefits, licensing, legislative records, or other state-governance artifacts."
    ],
    "risks": [
      "A message-centric lineage could be mistaken for a governance model, hiding unresolved authority and policy semantics.",
      "Without stronger naming, identity, and reconciliation rules, stale or ambiguous cross-agency records could be treated as trustworthy evidence.",
      "Government-state use amplifies the harm of under-specified escalation and revocation behavior."
    ],
    "open_questions": [
      "What authority, delegation, and revocation objects would sit above this CAS/message substrate for regional or agency governance?",
      "How would peers represent jurisdictional boundaries and escalation paths without introducing a hidden central registry?",
      "What local promise-accounting record would let Carol justify accepting or rejecting a contested application event?",
      "How would schema evolution, moderation/freeze rules, and retention policy work for legislative or licensing records over decades?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

### `results/SIM-hugoj-cas-usenetlike-gitlike/website-backend-hosting/openai-gpt-5.4-xhigh/20260520-184540.json`

```json
{
  "schema": "promisegrid.ga.result.v1",
  "result_id": "SIM-hugoj-cas-usenetlike-gitlike-website-backend-hosting-openai-gpt-5.4-xhigh-20260520-184540",
  "run_group_id": "ga-canary-20260520-184540",
  "cell_id": "ga-canary-20260520-184540-000003-SIM-hugoj-cas-usenetlike-gitlike--website-backend-hosting--openai-gpt-5.4-xhigh",
  "sim_id": "SIM-hugoj-cas-usenetlike-gitlike",
  "scenario_id": "website-backend-hosting",
  "model_id": "openai-gpt-5.4-xhigh",
  "timestamp_utc": "20260520-184540",
  "result_path": "results/SIM-hugoj-cas-usenetlike-gitlike/website-backend-hosting/openai-gpt-5.4-xhigh/20260520-184540.json",
  "runner": {
    "tool": "ga-runner",
    "provider": "openai",
    "api_model": "gpt-5.4",
    "reasoning_effort": "xhigh",
    "service_tier": "default",
    "served_service_tier": "default",
    "request_id": "req_16172f68cc4649c7a3f5ae7644604627",
    "response_id": "resp_026544e5df9d0eba006a0e01567c10819abe4f4d6563a1dcbe",
    "input_tokens": 5797,
    "output_tokens": 7225,
    "cost_usd": 0.11129475
  },
  "source": {
    "repo_commit": "bf8f3ec6b37af889c7f7b101cdf6c7b6ee00edaa",
    "sim_path": "simulations/SIM-hugoj-cas-usenetlike-gitlike/",
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
        "path": "simulations/SIM-hugoj-cas-usenetlike-gitlike/README.md",
        "sha256": "0506d32e30f1fe4903d1c0c1f2e6a7f25a3c3a438d6cfa0782d3e99c9448bdb8"
      },
      {
        "path": "scenarios/website-backend-hosting/website-backend-hosting.md",
        "sha256": "410d1a5fce15529c59260d9a9b20479dcb0eacb406ab9c480d2a6279f73fe18f"
      }
    ],
    "simulation_tree_hash": "144951c8a39cbb610ac9d305fe823308c60e01e271c2f3b222a32276d4a196ef"
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
    "promisegrid_alignment": 4,
    "auditability": 4,
    "evolution_safety": 4,
    "layer_boundary_clarity": 5,
    "failure_handling": 2,
    "implementation_plausibility": 2,
    "risk_penalty": 3
  },
  "fitness": {
    "raw": 19,
    "normalized_0_100": 54,
    "confidence_0_1": 0.79
  },
  "assessment": {
    "rationale": "The simulation offers a strong PromiseGrid-aligned evidence substrate—content-addressed objects, append-only and hash-chained history, store-and-forward replication, and transport separation—but only an indirect fit for website backend hosting. It helps with durable, locally auditable records under sparse knowledge, yet it does not define the concrete operational promises and decision rules needed for deployments, service ownership, secret rotation, uptime, and incident response.",
    "strengths": [
      "Content-addressed and hash-chained records support durable audit trails for disputed deployment or incident evidence.",
      "Store-and-forward replication across named peers or sites fits sparse, partial-knowledge conditions without assuming a central authority.",
      "The simulation is explicit about precedent versus identity and protocol meaning versus transport, which improves migration and long-term evolution safety."
    ],
    "weaknesses": [
      "The current specimen is message/feed oriented, not a website-backend hosting model with explicit deployment, ownership, uptime, or secret-rotation promises.",
      "No concrete object, feed, or authorization shapes are given for service ownership changes, rollback, secret revocation, or incident escalation.",
      "The README does not specify local decision rules for accepting, retrying, downgrading, or escalating when operational evidence is stale or conflicting."
    ],
    "risks": [
      "The Usenet/git analogy could be overextended, treating replicated evidence as sufficient for live operational control.",
      "Secret rotation and revocation require freshness and authorization semantics that are not present in the current simulation.",
      "Delayed propagation may be acceptable for audit logs but insufficient for uptime and incident-response commitments."
    ],
    "open_questions": [
      "What CAS objects and feeds would represent deployments, service owners, secret epochs, health observations, and incident decisions?",
      "How are freshness, expiry, supersession, and revocation modeled without a central registry or hidden global state?",
      "What locally observable threshold lets Alice or Carol accept, retry, downgrade, or escalate a contested backend action?",
      "How are long-lived service names and backend identities bound across migrations, key changes, and organizational turnover?"
    ],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}
```

## Required JSON Shape

{"child_id":"SIM-gibog-ga-child-0001","design_delta_summary":"one to three bounded design deltas","files":[{"path":"README.md","content":"# ..."},{"path":"QUESTION.md","content":"# ..."}]}
