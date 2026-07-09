# TODO-nudav: POC20 timeline pure-function CAS branches

## Status

Planned. Owns a parallel POC20 track for PromiseGrid timeline semantics,
pure-function agents, CAS object-chain branches, and branch-aware
capability-token double-spend behavior. POC20 is parallel to POC19, not a POC19
code-generation blocker. Source: `DI-kakos`; `TE-lodom`.

## Decision Intent Log

ID: DI-kakos
Date: 2026-07-09 12:17:37 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Create a parallel POC20 planning track, starting with `TE-lodom`, to
explore promises as timeline assertions, deterministic pure-function agents,
CAS-backed local and group timelines, and branch-aware capability-token
double-spend handling.
Intent: POC19 should remain the production-shaped plumbing path, but the current
POC lineage has not yet fully tested the deeper semantic model where every
promise is durable timeline state and double-spend is visible branch evidence
rather than hidden mutable ledger state. POC20 should test that model without
blocking POC19 scaffolding.
Constraints: Do not start POC20 code until the TE is reviewed and follow-up DF
locks the executable slice. Keep Promise Theory framing: local trust remains
local, no branch is globally authoritative, pure-function services must make
ambient context explicit, and capability-token semantics remain pCID-defined.
Affects: `docs/thought-experiments/TE-lodom-promise-timeline-pure-function-cas-branches.md`;
`protocols/wire-lab.d/TODO/TODO-nudav-poc20-timeline-pure-function-cas-branches.md`;
`protocols/wire-lab.d/TODO/TODO-vumas-poc19-production-shape.md`;
`implementations/poc19-production-shape/docs/DESIGN.md`;
`DEV-GUIDE-RESOURCES.md`.

## Tasks

- [x] nudav.1 Write `TE-lodom` covering promise timeline assertions,
  pure-function agents, CAS branches, local/group timelines, and double-spend
  branch semantics.
- [x] nudav.2 Record `DI-kakos` locking POC20 as a parallel semantic-model track,
  not a POC19 blocker.
- [x] nudav.3 Add POC19 cross-links so the production-shaped path avoids choices
  that would prevent POC20 branch-based promise histories.
- [ ] nudav.4 Review `TE-lodom` and decide whether Alternative C is accepted:
  durable CAS timelines with local ledgers as derived indexes.
- [ ] nudav.5 Lock the first executable POC20 slice with DF before code
  generation.
- [ ] nudav.6 Plan POC20 protocol specs for promise assertions, pure-function
  result promises, timeline branches, group timeline agreements, token issue,
  token redemption, and branch merge/non-merge records.
- [ ] nudav.7 Implement a small POC20 run with Alice, Bob, Carol, Dave, Ellen,
  and Mallory exercising local timelines, group timelines, deterministic
  function results, and double-spend branches.
- [ ] nudav.8 Add analyzer/regression gates proving that token double-spend is
  represented as visible branch evidence rather than hidden mutable state.

## Acceptance criteria for the future executable POC20

- CAS objects form visible parent-linked local and group timelines.
- Pure-function results are reproducible from explicit function, input, and
  context CIDs.
- Ambient inputs such as timestamps, randomness, sensor reads, model versions,
  peer quotes, and exchange rates are explicit context objects.
- Token issue, transfer, redemption, double-spend, and merge/non-merge behavior
  are represented as pCID-defined promise objects on branches.
- Receivers can keep, reject, merge, refuse to merge, or compensate branches
  using only local promise history and local trust judgments.
- No run requires a global ledger, global trust authority, hidden monitor, or
  out-of-band side channel to explain token state.
