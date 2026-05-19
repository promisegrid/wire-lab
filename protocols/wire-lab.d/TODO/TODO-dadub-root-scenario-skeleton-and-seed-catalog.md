# TODO-dadub: Root scenario skeleton and seed catalog

## Prior aliases

None. This TODO was minted after the proquint-handle migration.

## Status

Open. Owns the first root `scenarios/` design/population pass and the first
root `results/` documentation pass after TE-dojab locked the comparison
apparatus shape. This TODO does not create root `scenarios/` or `results/`
content by itself; it defines the work that the next implementation pass should
execute. Source: `DI-faros`; `DI-miror`; `DI-vabor`.

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
- [ ] dadub.3 Create one root scenario entry for every application listed in the
  seed application catalog below, using one `scenarios/<entry-id>/` directory per
  application rather than broad application-family grouping.
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
  mined entries under `DI-nanih`.
- [x] dadub.10 Run stale-layout checks for old one-file-per-cell result paths and
  run `git diff --check`. Done for the mined-entry pass under `DI-nanih`.

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
- Application entries should be lightweight at first. The goal is to make the
  pressure visible and comparable, not to fully model each industry in the first
  pass.
- Mined sim-local rows should keep their original source path and row title so
  future agents can audit whether the transformation preserved intent.
- Root result files should not be invented as examples. Until a real run occurs,
  result documentation should stop at README/template level.
