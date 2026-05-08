# TODO-turog: TE-41 group-session freeze procedure

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TODO-26` (integer alias)
- `TODO-20260507-002306` (timestamp alias and pre-migration filename)

## Status

Open. Depends on TE-40 (apparatus-vs-specimen completion) landing and
the TODO-pipus/TE-43 CAS/feed migration design. TE-nijab's
transport-layering / freeze-boundary DFs are locked, but DF-nijab.3
parks freeze-doc cleanup instead of executing it now.

## Threads absorbed from OPEN-THREADS.md

### T-GROUP-SESSION-FREEZE (formerly OPEN-THREADS, opened 2026-05-06)

Freeze chain for group-session and outer transport-spec. Carved from
T-021-TODO12 closure.

Scope: TODO-bisur subtask 012.8 lists this work but it is wider than
TODO-bisur. Sequence:
  1. Freeze `protocols/wire-lab.d/specs/transport-spec-draft.md` (its
     own freeze gate must clear first).
  2. Freeze `protocols/group-session.d/specs/group-session-draft.md`
     (012.8).
  3. Steve signs `merge-group-transport-spec` promise.
  4. `tools/spec freeze` mints both pCIDs, snapshots, manifest entries.
  5. Rename `transports/wire-lab-devs-draft/` to
     `transports/wire-lab-devs-<group-session-pcid>/` and rewrite every
     message's grid envelope from `grid draft:group-session` to
     `grid <pcid>` in a single mechanical commit.

Note: step 5 is now explicitly blocked by TE-nijab's locked result.
TE-nijab found that this is a category error: specs freeze, while
transport message directories are append-only specimen data. Do not
execute step 5 as written; DF-nijab.3 parks freeze-doc cleanup behind
TODO-pipus/TE-43 instead of rewriting this checklist now.

Blocking: outer transport-spec-draft.md has its own freeze gate;
nothing else gates this thread today.

Anchor: group-session-draft.md § Freeze gate; transport-spec-draft.md
§ Freeze gate; TODO-bisur subtask 012.8.
Disposition-file pointer: `dropped-thread-disposition-20260506.md`
§ TE-41 cluster (15 UTs).

## Question log

Per AGENTS-ppx Question-logging discipline.

- [x] **DF-nijab.1** What is `transports/`?
    opened: 2026-05-08 05:29 UTC
    answered: 2026-05-08 05:47 UTC
    asked of: stevegt@t7a.org (Steve Traugott)
    blocks: DR-suhod B/C/D cleanup and the group-session freeze procedure
    alternatives: 1.A lower-layer simulation surface / 1.B single-protocol directory
    recommendation: 1.A
    answer: 1.A lower-layer simulation surface
    linked DI: DI-026-20260508-054722

- [x] **DF-nijab.2** What happens to `wire-lab-devs-draft`?
    opened: 2026-05-08 05:29 UTC
    answered: 2026-05-08 05:47 UTC
    asked of: stevegt@t7a.org (Steve Traugott)
    blocks: transport-data migration and any freeze-time rewrite plan
    alternatives: 2.A preserve and supersede / 2.B derived rewrite only
    recommendation: 2.A
    answer: 2.A preserve and supersede
    linked DI: DI-026-20260508-054723

- [x] **DF-nijab.3** How should freeze docs be corrected?
    opened: 2026-05-08 05:29 UTC
    answered: 2026-05-08 05:47 UTC
    asked of: stevegt@t7a.org (Steve Traugott)
    blocks: TODO-turog step 5, DR-suhod B/C/D, and stale spec references
    alternatives: 3.A correct now after DF / 3.B park behind TODO-pipus/TE-43
    recommendation: 3.A
    answer: 3.B park behind TODO-pipus/TE-43
    linked DI: DI-026-20260508-054724

## Decision Intent Log

ID: DI-026-20260508-054722
Date: 2026-05-08 05:47:22
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: `transports/` is lower-layer network/feed simulation data, not a namespace owned by one frozen higher-layer protocol.
Intent: Preserve the OSI-style layering model tested by TE-nijab: one lower-layer carrier/feed can carry many higher-layer protocols without becoming any one of them.
Constraints: Do not rewrite existing transport message bytes or treat spec freeze as transport-data mutation. Future path/metadata details remain owned by TODO-pipus/TE-43 and related CAS/feed work.
Affects: `docs/thought-experiments/TE-nijab-transport-layering-and-freeze-boundaries.md`; `protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md`; `DR/DR-suhod-mihip-merge-blockers-partial-fix.md`; future `transports/` specs and data conventions.

ID: DI-026-20260508-054723
Date: 2026-05-08 05:47:23
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: `transports/wire-lab-devs-draft/` remains historical pre-layered specimen data; future migration or supersession must be additive.
Intent: Keep old message CIDs and bytes reconstructible while allowing future layered/CAS transport data to supersede the bootstrap specimen without falsifying history.
Constraints: Do not edit existing `.txt` transport messages as part of this decision. Any derived rewrite must be clearly non-authoritative unless a later DI gives it a distinct role.
Affects: `transports/wire-lab-devs-draft/`; `docs/thought-experiments/TE-nijab-transport-layering-and-freeze-boundaries.md`; `protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md`; `DR/DR-suhod-mihip-merge-blockers-partial-fix.md`.

ID: DI-026-20260508-054724
Date: 2026-05-08 05:47:24
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Freeze-doc cleanup is parked behind TODO-pipus/TE-43; TODO-turog step 5 stays marked suspect and blocked instead of being rewritten now.
Intent: Avoid locking premature wording for CAS/feed migration before the migration design exists, while making the current freeze checklist unsafe to execute as written.
Constraints: Do not update `group-session-draft.md`, `transport-spec-draft.md`, `transports/README.md`, or transport data in this pass. DR-suhod B/C/D remain reframed but not closed.
Affects: `protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md`; `DR/DR-suhod-mihip-merge-blockers-partial-fix.md`; future TODO-pipus/TE-43 cleanup.
