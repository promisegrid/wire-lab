# TODO-nudav: POC20 timeline pure-function CAS branches

## Status

Planned. Owns a parallel POC20 track for PromiseGrid timeline semantics,
pure-function agents, CAS object-chain branches, and branch-aware
capability-token double-spend behavior. POC20 is parallel to POC19, not a POC19
code-generation blocker. Source: `DI-kakos`; `DI-mokaz`; `DI-lamaz`;
`DI-lulog`; `TE-lodom`.

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

ID: DI-lamaz
Date: 2026-07-09 19:00:56 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Complete the POC20 pre-code lock for items 2-10: commit the prior CAS
event-source work first; use a hybrid of POC16 runtime/pCID/parser/proof/TCP
lessons and POC18 CAS/graph/sync lessons; make the first executable slice a
single unified scenario; use three hashable protocol-family spec docs; use CBOR
map payloads for the first semantic slice; and keep `DESIGN.md` as the single
human-readable design entrypoint.
Intent: POC20 code generation should start from locked semantics rather than
reopening runtime, pCID, payload, scenario, path, projection, double-spend, or
analyzer-gate choices. The pCID docs must be standalone because a pCID names a
whole spec document, while `DESIGN.md` remains the developer handoff document.
Constraints: No POC20 code generation in this batch. Cross-agent communication
in future code must use promise-shaped grid CBOR over TCP, not simulated
in-process transfer. Do not fragment into pCID-per-message-kind; message
variants are payload semantics under three protocol families. All future runtime
state remains rebuildable from local CAS unless a later DI explicitly narrows
the exception.
Affects: `implementations/poc20-timeline-pure-function-cas-branches/docs/DESIGN.md`;
`implementations/poc20-timeline-pure-function-cas-branches/docs/protocols/timeline.md`;
`implementations/poc20-timeline-pure-function-cas-branches/docs/protocols/pure-function.md`;
`implementations/poc20-timeline-pure-function-cas-branches/docs/protocols/capability-token.md`;
`protocols/wire-lab.d/TODO/TODO-nudav-poc20-timeline-pure-function-cas-branches.md`;
`DEV-GUIDE-RESOURCES.md`.

ID: DI-lulog
Date: 2026-07-09 19:19:00 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Model bootstrap Merkle/root CIDs, app/runtime root adoption, and
operator-approved root updates as POC20 timeline promises. Runtime and app root
CIDs that influence computation must be explicit context for pure-function
results.
Intent: POC19's minimum microkernel rule needs a matching semantic model: the
installed binary fetches executable roots, but local adoption of those roots is a
promise event in the agent's CAS timeline. Alice, Bob, or a voluntary group can
adopt different roots, later converge, fork, or reject updates without treating
any root as globally authoritative.
Constraints: Adoption is local and voluntary. A root CID names bytes, not trust,
obligation, or access by itself. Operator approval is a local promise and can be
revoked or superseded by a later timeline event. Changing the POC20 protocol
specs changes their pCIDs, so CID aliases must be recomputed after this DI is
implemented.
Affects: `implementations/poc20-timeline-pure-function-cas-branches/docs/DESIGN.md`;
`implementations/poc20-timeline-pure-function-cas-branches/docs/protocols/timeline.md`;
`implementations/poc20-timeline-pure-function-cas-branches/docs/protocols/pure-function.md`;
`implementations/poc20-timeline-pure-function-cas-branches/docs/protocols/capability-token.md`;
`protocols/wire-lab.d/TODO/TODO-nudav-poc20-timeline-pure-function-cas-branches.md`;
`DEV-GUIDE-RESOURCES.md`; `README.md`.

ID: DI-ruvum
Date: 2026-07-10 14:26:37 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Extend POC20 timeline semantics with replayable root-decision details,
structured projection-conflict records, source-key plus action-hash replay rules,
explicit projection rebuild proof obligations, and sensitive-data shareability
rules for private, encrypted-shareable, and plain-shareable CAS content.
Intent: POC20 should test operational timeline behavior where local decisions can
be replayed from CAS without hidden mutable projection state, without silently
overwriting reviewed values, and without embedding unnecessary secrets or
sensitive personal payloads in broadly replicated trees.
Constraints: Keep all decisions local and voluntary. Use local event language
rather than production monitor language. A source key must come from a stable
upstream fact, root event, or timeline event; generated action hashes are
secondary guards, not primary identity. Sensitive payloads may be private,
encrypted, referenced by opaque handle, or represented by keyed commitment; avoid
plain hashes of guessable sensitive data.
Affects: `protocols/wire-lab.d/TODO/TODO-nudav-poc20-timeline-pure-function-cas-branches.md`;
`implementations/poc20-timeline-pure-function-cas-branches/docs/DESIGN.md`;
`implementations/poc20-timeline-pure-function-cas-branches/docs/protocols/timeline.md`;
`implementations/poc20-timeline-pure-function-cas-branches/docs/protocols/pure-function.md`;
`implementations/poc20-timeline-pure-function-cas-branches/docs/protocols/capability-token.md`;
`DEV-GUIDE-RESOURCES.md`; `README.md`.

ID: DI-zizab
Date: 2026-07-11 18:15:37 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Remove `-v1` from current POC20 protocol-family spec names and
supersede rollback-root wording with prior/corrective Merkle/CAS root wording.
The current POC20 protocol specs are `timeline.md`, `pure-function.md`, and
`capability-token.md`; their pCID symlinks are recomputed after content changes.
Intent: A pCID already versions a protocol by content identity, so a manual `v1`
suffix is redundant and misleading. POC20 root-decision records should align with
POC19 and POC21 by preserving prior/corrective root history instead of implying
rollback safety.
Constraints: Do not rewrite historical DI bodies. Current prose, acceptance
criteria, spec paths, and pCID aliases must use unsuffixed protocol-family names.
Changing spec bytes requires recomputing the CIDv1 base32 symlink aliases.
Affects: `protocols/wire-lab.d/TODO/TODO-nudav-poc20-timeline-pure-function-cas-branches.md`;
`implementations/poc20-timeline-pure-function-cas-branches/docs/DESIGN.md`;
`implementations/poc20-timeline-pure-function-cas-branches/docs/protocols/timeline.md`;
`implementations/poc20-timeline-pure-function-cas-branches/docs/protocols/pure-function.md`;
`implementations/poc20-timeline-pure-function-cas-branches/docs/protocols/capability-token.md`;
`README.md`; `DEV-GUIDE-RESOURCES.md`.
Supersedes: `DI-lamaz` spec-path naming only; `DI-lulog` Merkle/root wording
only; `DI-ruvum` rollback-root wording only.

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
- [x] nudav.5 Lock the first executable POC20 slice with DF before code
  generation: one unified scenario with Alice, Bob, Carol, Dave, Ellen, and
  Mallory exercising timeline, pure-function, token double-spend, projection
  rebuild, and CAS shareability behavior. Source: `DI-lamaz`.
- [x] nudav.6 Plan POC20 protocol specs for promise assertions, pure-function
  result promises, timeline branches, group timeline agreements, token issue,
  token redemption, and branch merge/non-merge records. Source: `DI-lamaz`.
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
- [x] nudav.11 Lock future POC20 runtime architecture as POC16 runtime/pCID/parser
  lessons plus POC18 CAS/graph/sync lessons. Source: `DI-lamaz`.
- [x] nudav.12 Lock three protocol-family pCIDs for POC20: timeline,
  pure-function, and capability-token. Source: `DI-lamaz`.
- [x] nudav.13 Lock CBOR map payloads for the first semantic POC20 slice. Source:
  `DI-lamaz`.
- [x] nudav.14 Lock CAS object profiles and derived projection rebuild rules in
  the implementation-local design. Source: `DI-lamaz`.
- [x] nudav.15 Lock double-spend and pure-function scenario details before code
  generation. Source: `DI-lamaz`.
- [x] nudav.16 Lock future package, command, runtime path, and diagnostic path
  names for code generation. Source: `DI-lamaz`.
- [x] nudav.17 Define analyzer gates before code generation while leaving their
  implementation in `nudav.8`. Source: `DI-lamaz`.
- [x] nudav.18 Lock bootstrap root adoption and root update semantics as local
  timeline promises, with runtime/app roots included in pure-function context
  when they affect computation. Source: `DI-lulog`.
- [x] nudav.19 Lock replayable root-decision detail: adopted, rejected,
  superseded, prior-root-retained, corrective-root-adopted,
  full-state-restore-attempted, full-state-restore-performed, still-evaluating,
  approving role, impact summary CID, capability-change summary, local reason,
  and prior/corrective Merkle/CAS root CIDs. Source: `DI-ruvum`; `DI-zizab`.
- [x] nudav.20 Lock structured projection-conflict and replay semantics: conflicts
  are local timeline records, source keys are stable upstream/timeline facts, and
  generated action hashes are secondary guards requiring new local decisions when
  they change. Source: `DI-ruvum`.
- [x] nudav.21 Lock privacy-compatible storage requirements: sensitive payloads use
  private or encrypted CAS namespaces, opaque handles, encrypted-object CIDs, or
  keyed commitments; broad summaries must avoid plaintext secrets and unnecessary
  sensitive personal payloads. Source: `DI-ruvum`.
- [x] nudav.22 Remove redundant `-v1` suffixes from current POC20 protocol-family
  spec names and recompute their pCID symlink aliases. Source: `DI-zizab`.

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
- Three standalone hashable spec docs exist for the first protocol families, and
  each has a CIDv1 base32 symlink alias under the POC20 `docs/protocols/`
  directory.
- Future cross-agent communication uses promise-shaped grid CBOR over TCP rather
  than in-process simulated transfer.
- Bootstrap roots and app/runtime update roots are adopted by local promise
  events. A root CID names a Merkle/CAS object graph; it does not by itself
  create trust, authority, or an obligation to execute.
- Pure-function result records include the relevant app/runtime root CIDs in
  explicit context when executable code or runtime behavior affects the result.
- Root-decision records include local decision state, approving role, impact
  summary CID, capability-change summary, local reason, prior Merkle/CAS root CID
  when one exists, any corrective Merkle/CAS root CID, and explicit
  full-state-restore attempts where applicable.
- Projection conflicts are represented as local timeline records that can be
  replayed from CAS and do not silently overwrite reviewed values or stale
  projections.
- Replay keys use stable source facts, root events, or timeline events as primary
  identity and generated action hashes only as secondary same-action guards.
- Sensitive examples prove that broad replicated summaries avoid plaintext
  secrets and unnecessary sensitive personal payloads.
