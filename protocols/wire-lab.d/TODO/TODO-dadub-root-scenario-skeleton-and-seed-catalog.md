# TODO-dadub: Root scenario skeleton and seed catalog

## Prior aliases

None. This TODO was minted after the proquint-handle migration.

## Status

Closed. The first root `scenarios/` design/population pass and the first root
`results/` documentation pass are complete: root templates exist, existing
simulation-local scenario rows have been mined, and seed application entries
have been created. Source: `DI-faros`; `DI-miror`; `DI-vabor`; `DI-dimas`;
`DI-botup`; `DI-nanih`; `DI-midif`.

## Decision Intent Log

ID: DI-vabor
Date: 2026-05-18 18:58:46
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: File a harness TODO for the root scenario skeleton and seed catalog.
The TODO must plan one root `scenarios/<entry-id>/` entry per application, one
root scenario entry per mined row from existing simulation-local `SCENARIOS.md`
tables, transformed root scenarios with source links rather than bare references
or verbatim copies, and a root `results/README.md` / template pass without fake
run-result files.
Intent: TE-dojab locked root scenarios and per-run root results as the comparison
apparatus, but the repo still needs a concrete population plan. The first
scenario population should be broad enough to expose application pressure early,
while still preserving provenance from existing sim-local scenario matrices and
not inventing fake run evidence.
Constraints: Do not create root `scenarios/` or `results/` content in this TODO
write pass. Do not group application entries into broad families by default.
Use stable kebab-case entry IDs. For mined simulation rows, preserve source
provenance and transform each row into a reusable cross-simulation pressure
entry. Root result paths remain governed by `DI-miror`:
`results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md`.
Affects: `protocols/wire-lab.d/TODO/TODO-dadub-root-scenario-skeleton-and-seed-catalog.md`;
`protocols/wire-lab.d/TODO/TODO.md`; future `scenarios/`; future `results/`;
future `simulations/README.md`; future `DEV-GUIDE-RESOURCES.md`.

ID: DI-dimas
Date: 2026-05-18 19:09:34
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Execute the first root scenario/result skeleton pass by creating root
`scenarios/README.md` and `results/README.md`, embedding the scenario-entry,
matrix, and result-run templates there, and updating current-pointer docs after
the skeleton exists.
Intent: The repo now has locked root scenario/result apparatus decisions and a
TODO owner, but agents need concrete root entry and result-run contracts before
they can safely populate application scenarios or mine sim-local scenario rows.
Keeping templates in README files avoids extra root template files and avoids
fake result evidence.
Constraints: Do not create root scenario entries yet. Do not create fake
`results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md` run files.
Keep root scenarios and results as wire-lab harness comparison apparatus with
evidence-only authority.
Affects: `protocols/wire-lab.d/TODO/TODO-dadub-root-scenario-skeleton-and-seed-catalog.md`;
`scenarios/README.md`; `results/README.md`; `simulations/README.md`;
`DEV-GUIDE-RESOURCES.md`.

ID: DI-botup
Date: 2026-05-18 19:14:07
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Require every root scenario to include overarching PromiseGrid goal
checks, including the 100-year durability goal, sparse partial knowledge, no
central authority or registry, peer-local promise accounting, adversarial or
failure pressure, human/LLM auditability, and migration/evolution pressure.
Intent: Root scenarios are supposed to compare simulations against the reasons
PromiseGrid exists, not just against narrow application happy paths. Making the
goal checks explicit prevents scenario authors from accidentally assuming away
long-horizon, sparse, decentralized, adversarial, or human-auditable constraints.
Constraints: Record the gate in `scenarios/README.md` rather than duplicating it
in every current TODO. A scenario may mark a gate not relevant only by saying
why. Do not create scenario entries in this pass.
Affects: `protocols/wire-lab.d/TODO/TODO-dadub-root-scenario-skeleton-and-seed-catalog.md`;
`scenarios/README.md`.

ID: DI-nanih
Date: 2026-05-18 19:17:55
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Execute the first mining pass by transforming every row from the
existing simulation-local `SCENARIOS.md` sources into a root scenario entry, one
directory per source row, with README, MATRIX, and scenario markdown files.
Intent: Existing simulation scenario rows already encode concrete design
pressure and provenance. Moving them into root scenario entries gives future
model runs an apples-to-apples comparison surface without deleting or rewriting
the simulation-local sources.
Constraints: Preserve source path and source row/title in every mined scenario.
Apply the `DI-botup` overarching goal checks to every scenario. Keep every
matrix in `not-run` state until a real run exists. Do not create fake result
run files. Do not populate application seed entries in this pass.
Affects: `protocols/wire-lab.d/TODO/TODO-dadub-root-scenario-skeleton-and-seed-catalog.md`;
`scenarios/`.

ID: DI-midif
Date: 2026-05-18 19:22:41
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Execute the seed application pass by creating one root
`scenarios/<entry-id>/` entry for every application listed in the seed
application catalog, with README, MATRIX, and scenario markdown files for each
entry.
Intent: The root comparison apparatus needs application pressure, not just
protocol-internal or simulation-local pressure. Lightweight application seeds
let candidate simulations be compared against real-world domains early while
keeping result evidence separate until actual model runs exist.
Constraints: Preserve one directory per catalog entry. Keep entries lightweight
but include application pressure, actors, expected PromiseGrid stress points,
the `DI-botup` overarching goal checks, and the evidence-only authority
boundary. Keep every matrix in `not-run` state until a real run exists. Do not
create fake result run files.
Affects: `protocols/wire-lab.d/TODO/TODO-dadub-root-scenario-skeleton-and-seed-catalog.md`;
`scenarios/`; `scenarios/README.md`; `simulations/README.md`;
`DEV-GUIDE-RESOURCES.md`.

ID: DI-moduf
Date: 2026-05-18 20:15:39
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Matrix cell result files must be produced by LLM or human reasoning,
not by mechanical document parsing or keyword heuristics. Scripts may generate
manifests, LLM prompts, queues, and validation/comparison reports; scripted
prototype outputs are plumbing tests only and are excluded from result
validation and comparison by default.
Intent: Root results are supposed to preserve deeper reasoning about how a
simulation behaves under a scenario's 100-year, sparse-knowledge,
no-central-authority, auditability, and migration pressures. Mechanical parsing
can verify shape and prepare work, but treating parser-generated prose as
evidence would create fake confidence and undermine the apples-to-apples
simulation matrix.
Constraints: Preserve any already-written prototype files for audit history
rather than deleting them. Mark prototype matrix rows explicitly. Keep
`results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md` as the run
artifact path from `DI-miror`, but require real result-producing modes to be
`codex-manual-blind` or `llm-doc-eval-blind` unless Steve explicitly opts into
prototype handling for tooling tests.
Affects: `protocols/wire-lab.d/TODO/TODO-dadub-root-scenario-skeleton-and-seed-catalog.md`;
`results/RUN-PROTOCOL.md`; `results/README.md`; `results/tools/`;
`results/comparisons/`; `scenarios/*/MATRIX.md`.

ID: DI-nuhon
Date: 2026-05-18 20:57:50
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Full-matrix runs must be able to execute unattended through a
checkpointed queue: each manifest row carries a concrete timestamp, result path,
ordinal, and cell ID; a queue runner invokes an external LLM command once per
cell; each cell is validated and then used to update the scenario matrix; queue
state is persisted after every cell so a 7000-cell run can resume without
interactive keypresses.
Intent: The matrix is too large to drive manually prompt-by-prompt, but result
files still need deeper LLM/human reasoning under `DI-moduf`. The correct split
is for scripts to coordinate durable work units and validation while delegating
the actual verdict prose to a real model runner.
Constraints: Do not synthesize result verdict prose mechanically. Preserve
`results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md` as the result
artifact path. Keep the runner command external and explicit so Codex, OpenAI
API, Anthropic API, Perplexity, or another evaluator can be swapped without
changing result semantics. State and prompts are operational artifacts under
`results/state/` and `results/jobs/`; result files remain the evidence source.
Affects: `results/tools/matrix_common.py`; `results/tools/generate_matrix_manifest.py`;
`results/tools/generate_llm_jobs.py`; `results/tools/matrix_queue.py`;
`results/tools/update_matrix_rows.py`; `results/tools/validate_results.py`;
`results/tools/README.md`; `results/RUN-PROTOCOL.md`; `scenarios/*/MATRIX.md`.

ID: DI-bujiv
Date: 2026-05-18 21:05:50
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: The unattended matrix queue must persist a cell's `running` state
before launching the external LLM runner, not only after the runner exits.
Intent: If a long unattended run is killed, loses power, or the terminal dies
while an LLM cell is in flight, the checkpoint file should identify the exact
in-flight cell instead of making the run appear to have stopped cleanly before
that cell began.
Constraints: Keep checkpointing at cell granularity; do not attempt partial
checkpointing inside a single LLM response. On restart, `done` cells remain
skipped by default, while interrupted `running` cells may be retried explicitly
by the queue's normal retry path once reviewed.
Affects: `results/tools/matrix_queue.py`;
`protocols/wire-lab.d/TODO/TODO-dadub-root-scenario-skeleton-and-seed-catalog.md`.

ID: DI-lulom
Date: 2026-05-18 21:30:12
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Port the result matrix tooling into a standalone Go CLI at
`tools/matrix-runner/`, keep the Python scripts as legacy/reference for now,
and implement OpenAI API-backed unattended execution as the first real provider
path.
Intent: A durable full-matrix run needs one typed runner that owns manifests,
prompt bundling, provider calls, checkpointing, validation, matrix updates, and
comparison reports without shelling through a separate LLM command per cell.
The API runner must bundle local source document contents because remote model
APIs cannot read repository paths directly.
Constraints: Preserve the existing result path shape and manifest CSV schema.
Do not synthesize verdict prose mechanically; API output must provide the result
body. Use raw Go HTTP for the OpenAI Responses API in this pass and leave
Anthropic/Perplexity as future provider adapters. Do not delete the Python
scripts until a later retirement decision.
Affects: `tools/matrix-runner/`; `tools/README.md`; `results/tools/README.md`;
`results/RUN-PROTOCOL.md`; `results/README.md`;
`protocols/wire-lab.d/TODO/TODO-dadub-root-scenario-skeleton-and-seed-catalog.md`.

ID: DI-zamin
Date: 2026-05-18 22:18:33
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Delete committed `scenarios/*/MATRIX.md` files and make `results/`
the canonical source of truth for completed scenario/simulation/model runs.
Matrix-style summaries must be generated from `results/` on demand rather than
maintained as committed scenario state.
Intent: Scenario trees should be input context and result trees should be output
evidence. Committed matrices duplicate result paths and verdict notes, create
placeholder `not-run` churn, and risk being mistaken for prompt context.
Generated views preserve scanability without creating a second source of truth.
Constraints: Preserve root `scenarios/` and root per-run `results/` from
TE-dojab. Do not delete result files. Do not feed generated result views or
prior result summaries into blind run prompts. Keep legacy Python scripts as
reference, but mark matrix-update and strict-matrix behavior obsolete.
Affects: `docs/thought-experiments/TE-lizuh-results-canonical-no-scenario-matrix.md`;
`docs/thought-experiments/TE-dojab-simulation-run-model-and-scenario-result-matrix.md`;
`scenarios/`; `scenarios/README.md`; `results/README.md`;
`results/RUN-PROTOCOL.md`; `results/tools/README.md`; `results/tools/`;
`tools/matrix-runner/`; `tools/README.md`;
`protocols/wire-lab.d/TODO/TODO-dadub-root-scenario-skeleton-and-seed-catalog.md`.

## Scope

- Define the root scenario-entry contract implied by `DI-faros`: a root scenario
  entry is a `scenarios/<entry-id>/` directory with a local `README.md`, local
  `MATRIX.md`, and one or more scenario markdown files.
- Populate root scenario entries broadly in the next implementation pass: one
  entry per application and one entry per mined sim-local scenario row.
- Transform existing sim-local `SCENARIOS.md` rows into root scenario entries
  with source links, cross-sim applicability notes, and authority-boundary
  language.
- Create root `results/README.md` and result templates only; do not create fake
  run files until a real model run exists.
- Keep root scenarios and results as wire-lab harness comparison apparatus, not
  PromiseGrid node layout, production API, or simulation-local world state.

## Subtasks

- [x] dadub.1 Define the root scenario-entry template, including required
  fields, source/provenance field, applicability field, evaluation questions,
  and authority boundary. Done in `scenarios/README.md` under `DI-dimas`.
- [x] dadub.2 Create `scenarios/README.md` explaining that root scenarios are
  shared comparison apparatus governed by `DI-faros`, not shared protocol
  components or PromiseGrid deployment layout. Done under `DI-dimas`.
- [x] dadub.3 Create one root scenario entry for every application listed in the
  seed application catalog below, using one `scenarios/<entry-id>/` directory per
  application rather than broad application-family grouping. Done for 52
  application entries under `DI-midif`.
- [x] dadub.4 Mine every row from every existing simulation-local `SCENARIOS.md`
  source listed below and create one transformed root scenario entry per row.
  Done for 95 rows under `DI-nanih`.
- [x] dadub.5 For each mined scenario entry, record the source simulation path,
  source row title, original decision pressure, and what cross-sim comparisons
  the root entry is meant to enable. Done under `DI-nanih`.
- [x] dadub.6 Create `results/README.md` documenting the locked path
  `results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md`, UTC
  timestamp rule, model-ID rule, and evidence-only authority boundary. Source:
  `DI-miror`; `DI-faros`. Done under `DI-dimas`.
- [x] dadub.7 Add a result-run template without creating any fake run-result
  files. The template may live in `results/README.md` or a template file chosen
  in the implementation pass. Done in `results/README.md` under `DI-dimas`.
- [x] dadub.8 Update `simulations/README.md` and `DEV-GUIDE-RESOURCES.md` after
  root `scenarios/` and `results/` skeletons exist, so readers know where
  cross-sim scenario matrices and run evidence live. Done under `DI-dimas`.
- [x] dadub.9 Validate that every root scenario entry preserves provenance,
  avoids direct design authority claims, has a stable kebab-case entry ID, and
  addresses the overarching PromiseGrid goal checks from `DI-botup`. Done for
  mined entries under `DI-nanih` and application seed entries under `DI-midif`.
- [x] dadub.10 Run stale-layout checks for old one-file-per-cell result paths and
  run `git diff --check`. Done for the scenario population passes under
  `DI-nanih` and `DI-midif`.
- [x] dadub.11 Add unattended full-matrix run support with concrete manifest
  paths, checkpointed queue state, per-cell LLM prompt generation, validation,
  and scenario matrix row updates. Done under `DI-nuhon`.
- [x] dadub.12 Persist `running` queue state before invoking the external LLM
  runner so interrupted unattended runs expose the in-flight cell. Done under
  `DI-bujiv`.
- [x] dadub.13 Port matrix/result tooling into `tools/matrix-runner/`, including
  manifest, jobs, run, progress, validate, update-matrix, compare, and OpenAI
  API-backed prompt bundling. Done under `DI-lulom`; the update-matrix clause
  was superseded by dadub.14 / `DI-zamin`.
- [x] dadub.14 Remove committed scenario `MATRIX.md` state, make `results/`
  canonical for completed runs, and replace matrix updates with generated views.
  Done under `DI-zamin`.

## Candidate scenario entries

### Existing simulation-local scenario sources to mine

The next implementation pass must mine every table row in these files. Each row
becomes its own root scenario entry, transformed for cross-simulation use rather
than copied verbatim.

- `simulations/SIM-bobud-l6-cas-starting-profile-bakeoff/SCENARIOS.md`
- `simulations/SIM-bohof-group-session-freeze-promise/SCENARIOS.md`
- `simulations/SIM-dihiz-peer-adoption-metadata/SCENARIOS.md`
- `simulations/SIM-gobaz-chunking-identity-bakeoff/SCENARIOS.md`
- `simulations/SIM-haros-promise-economy-spectrum/SCENARIOS.md`
- `simulations/SIM-jomag-cas-object-model/SCENARIOS.md`
- `simulations/SIM-jurar-cas-backed-group-session/SCENARIOS.md`
- `simulations/SIM-kohad-cas-object-type-binding-bakeoff/SCENARIOS.md`
- `simulations/SIM-kuful-udp-feed-v0-conformance/SCENARIOS.md`
- `simulations/SIM-ligan-promisebase-reference-naming/SCENARIOS.md`
- `simulations/SIM-narok-transport-family-bakeoff/SCENARIOS.md`
- `simulations/SIM-punaz-bgp-class-routing-app/SCENARIOS.md`
- `simulations/SIM-ranib-spec-requirement-sections/SCENARIOS.md`
- `simulations/SIM-rusap-promise-accounting-records/SCENARIOS.md`
- `simulations/SIM-zarud-conditional-release-geofencing/SCENARIOS.md`
- `simulations/SIM-zazit-chunk-feed-replication/SCENARIOS.md`

### Seed application entries

The next implementation pass should create one root `scenarios/<entry-id>/`
entry per application below. The initial entry can be minimal, but it must state
the application pressure, actors, expected PromiseGrid stress points, and result
authority boundary.

| Entry ID | Application pressure |
|---|---|
| `bgp-routing` | Route-like reachability promises, hijacks, leaks, stale paths, and sparse topology knowledge. |
| `erp` | Multi-department operational promises spanning orders, inventory, purchasing, accounting, and fulfillment. |
| `accounting` | Audit trails, payment promises, reconciliations, corrections, and multi-party evidence. |
| `shipping-label-printing` | Carrier selection, label authorization, address data, payment, and chain-of-custody promises. |
| `nonprofit-management` | Donations, restricted funds, volunteers, governance, program delivery, and public accountability. |
| `open-source-development` | Issues, patches, review promises, release artifacts, maintainer authority, and contributor reputation. |
| `makerspace-door-access` | Membership state, door authorization, guest access, emergency override, and local accountability. |
| `machine-equipment-access` | Equipment authorization tied to training, certification, maintenance status, and safety constraints. |
| `training-certification` | General credential issuance, renewal, revocation, prerequisites, and verifier trust. |
| `manufacturing-order-flow` | Work orders, materials, machine availability, quality checks, and shipment handoff. |
| `website-backend-hosting` | Deployments, service ownership, secret rotation, incident response, and uptime promises. |
| `tuition-attendance` | Tuition payment, class enrollment, attendance evidence, refunds, and credential eligibility. |
| `makerspace-dues-access` | Dues payment, membership standing, door access, tool access, and suspension/reinstatement. |
| `municipal-governance` | Local ordinances, service requests, permits, budgets, public records, and citizen participation. |
| `state-governance` | Cross-agency coordination, benefits, licensing, legislative records, and regional authority. |
| `national-governance` | Large-scale public authority, federal records, procurement, identity, and checks against capture. |
| `community-movement-organizing` | Membership, campaigns, working groups, public commitments, moderation, and schism handling. |
| `climate-coordination` | Global coordination of mitigation/adaptation efforts, funding promises, measurements, and accountability. |
| `nuclear-nonproliferation` | High-stakes verification promises, inspection evidence, sanctions, and trust under adversarial incentives. |
| `internet-governance` | Protocol evolution, registry-like pressure, multi-stakeholder legitimacy, and schism resilience. |
| `devops` | Deploy, rollback, monitoring, incident response, access control, and postmortem promises. |
| `system-administration` | User accounts, configuration drift, backups, patching, access review, and break-glass authority. |
| `aviation-flight-planning` | Flight plans, weather, aircraft status, crew qualifications, airspace constraints, and dispatch release. |
| `airspace-management` | Airspace reservations, restrictions, deconfliction, notices, and multi-authority coordination. |
| `airfield-management` | Runway/taxiway status, maintenance closures, ground services, safety reports, and local authority. |
| `air-traffic-control` | Separation, handoff, clearances, stale data, emergency priority, and real-time coordination. |
| `flight-training` | Student progress, instructor endorsements, aircraft checkout, weather minimums, and certification evidence. |
| `aerospace-design` | Requirements, analyses, design reviews, configuration control, and long-lived engineering evidence. |
| `aerospace-development` | Program execution, supplier evidence, test results, safety gates, and certification dependencies. |
| `aerospace-project-funding` | Milestone funding, investor/grantor promises, deliverable evidence, and governance under uncertainty. |
| `decentralized-manufacturing` | Distributed fabrication, machine capability claims, quality evidence, and local fulfillment. |
| `decentralized-supply-chain` | Supplier promises, provenance, substitutions, delays, recalls, and sparse knowledge of the graph. |
| `logistics` | Routing, custody, capacity, exceptions, proof of delivery, and multi-carrier coordination. |
| `securities-trading` | Orders, settlement, compliance evidence, counterparty risk, and market manipulation pressure. |
| `multicurrency-transactions` | Exchange-rate promises, settlement, local valuation, accounting, and currency-risk evidence. |
| `personal-currencies` | Peer-issued value, local exchange rates, default history, trust portability, and speculation failure. |
| `disaster-response` | Resource requests, responder credentials, logistics, donations, and rapidly changing authority. |
| `mutual-aid` | Local needs/offers, trust, delivery evidence, conflict resolution, and relationship-based accounting. |
| `healthcare-referrals` | Referral promises, consent, provider credentials, scheduling, and outcome evidence. |
| `healthcare-records-consent` | Patient consent, revocation, delegated access, audit trails, and emergency exceptions. |
| `insurance-claims` | Claim evidence, adjuster authority, fraud pressure, payments, and appeal promises. |
| `scientific-data-collaboration` | Dataset provenance, authorship, access promises, reproducibility, and embargoes. |
| `journalism-source-provenance` | Source protection, evidence chains, publication promises, corrections, and adversarial disinformation. |
| `legal-contracting` | Offers, acceptance, amendments, signatures, jurisdiction, evidence, and dispute resolution. |
| `energy-grid-coordination` | Generation/load promises, demand response, outage repair, safety constraints, and market signals. |
| `agriculture-food-traceability` | Farm inputs, harvest lots, cold chain, inspections, recalls, and certification evidence. |
| `federated-identity-credentials` | Issuer trust, credential presentation, revocation, selective disclosure, and verifier policy. |
| `software-supply-chain-attestation` | Build provenance, dependency promises, signing, vulnerability response, and reproducibility. |
| `iot-fleet-maintenance` | Device identity, maintenance history, firmware updates, telemetry, and access control. |
| `emergency-communications` | Intermittent connectivity, priority messages, responder trust, and degraded-mode operation. |
| `housing-coop-management` | Membership, dues, maintenance obligations, governance, disputes, and shared-resource access. |
| `public-procurement` | Solicitations, bids, award promises, delivery evidence, conflict-of-interest checks, and auditability. |

## Implementation notes

- BGP should appear twice in the broader scenario work: once as the existing
  `SIM-punaz-bgp-class-routing-app` source to mine, and once as a root
  `scenarios/bgp-routing/` application entry.
- Application entries are intentionally lightweight in the first pass. The goal
  is to make the pressure visible and comparable, not to fully model each
  industry before the first runs.
- Mined sim-local rows should keep their original source path and row title so
  future agents can audit whether the transformation preserved intent.
- Root result files should not be invented as examples. Until a real run occurs,
  result documentation should stop at README/template level.
