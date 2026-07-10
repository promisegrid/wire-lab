# TODO-nudav: POC20 timeline pure-function CAS branches

## Status

Planned. Owns a parallel POC20 track for PromiseGrid timeline semantics,
pure-function agents, CAS object-chain branches, and branch-aware
capability-token double-spend behavior. POC20 is parallel to POC19, not a POC19
code-generation blocker. Source: `DI-kakos`; `DI-mokaz`; `TE-lodom`.

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

ID: DI-bibah
Date: 2026-07-09 15:13:53 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Add an implementation-local POC20 design document at
`implementations/poc20-timeline-pure-function-cas-branches/docs/DESIGN.md`
while keeping `TE-lodom` and `TODO-nudav` in their canonical global locations.
Intent: The POC20 implementation directory should have a local design entrypoint
before any code generation starts, but the cross-cutting TE and TODO should
remain discoverable through the global TE corpus and harness TODO queue.
Constraints: The design document is not executable code, not a frozen protocol
spec, and not a POC19 blocker. It must point back to `TE-lodom`, `TODO-nudav`,
and `DI-kakos`; it must preserve the durable-CAS-timeline model with local
ledgers as derived indexes; it must not define a global ledger or global branch
authority.
Affects: `implementations/poc20-timeline-pure-function-cas-branches/docs/DESIGN.md`;
`protocols/wire-lab.d/TODO/TODO-nudav-poc20-timeline-pure-function-cas-branches.md`;
`DEV-GUIDE-RESOURCES.md`.

ID: DI-mokaz
Date: 2026-07-09 18:42:55 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Lock POC20's durable state model around local CAS as the
chronological event source and source of truth. Derived local indexes,
projections, caches, JSON files, SQLite tables, or in-memory summaries may
exist only as rebuildable views over local CAS.
Intent: The POC20 plan should not imply hidden mutable projection-only state, a
global ledger, or any projection that can outrank the CAS event stream. PromiseGrid
agents need to explain token state, branch state, trust state, compute context,
and group-timeline state from CAS objects that can be replayed, compared,
shared, withheld, or encrypted under local promises.
Constraints: Local CAS is not automatically public. Some local CAS objects may
remain private forever, some may be shareable only after encryption, and some
may be plain-shareable. Sharing any CAS object remains a local promise decision.
If a derived projection disagrees with local CAS, local CAS wins. Do not start
POC20 code generation until the executable slice is separately locked by DF.
Affects: `docs/thought-experiments/TE-lodom-promise-timeline-pure-function-cas-branches.md`;
`protocols/wire-lab.d/TODO/TODO-nudav-poc20-timeline-pure-function-cas-branches.md`;
`implementations/poc20-timeline-pure-function-cas-branches/docs/DESIGN.md`;
`DEV-GUIDE-RESOURCES.md`.
Supersedes: `DI-bibah` derived-projection terminology only; `DI-bibah` remains active
for the implementation-local design-document placement decision.

## Tasks

- [x] nudav.1 Write `TE-lodom` covering promise timeline assertions,
  pure-function agents, CAS branches, local/group timelines, and double-spend
  branch semantics.
- [x] nudav.2 Record `DI-kakos` locking POC20 as a parallel semantic-model track,
  not a POC19 blocker.
- [x] nudav.3 Add POC19 cross-links so the production-shaped path avoids choices
  that would prevent POC20 branch-based promise histories.
- [x] nudav.4 Review `TE-lodom` and decide whether Alternative C is accepted:
  durable local CAS event streams with derived local projections and indexes.
  Source: `DI-mokaz`.
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
- [x] nudav.9 Add an implementation-local POC20 design document at
  `implementations/poc20-timeline-pure-function-cas-branches/docs/DESIGN.md`.
  Source: `DI-bibah`.
- [x] nudav.10 Lock the CAS event-source invariant: local CAS is the source of
  truth; all local indexes, caches, and projections are rebuildable from local
  CAS; and object sharing is a separate local promise decision. Source:
  `DI-mokaz`.

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
- Local CAS is the chronological event source and source of truth for each
  agent; current state is always rebuildable by replaying local CAS event
  objects.
- Any SQLite table, JSON file, in-memory map, or other local projection is
  disposable and rebuildable from local CAS. If a projection and local CAS
  disagree, local CAS wins.
- Local CAS objects can be private, encrypted-shareable, or plain-shareable, and
  the existence of an object in CAS never implies that the agent promises to send
  it.
