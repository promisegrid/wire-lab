# TODO-turog: TE-41 group-session freeze procedure

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TODO-26` (integer alias)
- `TODO-20260507-002306` (timestamp alias and pre-migration filename)

## Status

Open. Depends on TE-40 (apparatus-vs-specimen completion) landing and
TE-nijab's transport-layering / freeze-boundary DFs. No twig yet.

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

Note: step 5 is now explicitly under review by TE-nijab. TE-nijab tests
whether this is a category error: specs freeze, while transport message
directories are append-only specimen data. Do not execute step 5 until
DF-nijab.1 through DF-nijab.3 are locked.

Blocking: outer transport-spec-draft.md has its own freeze gate;
nothing else gates this thread today.

Anchor: group-session-draft.md § Freeze gate; transport-spec-draft.md
§ Freeze gate; TODO-bisur subtask 012.8.
Disposition-file pointer: `dropped-thread-disposition-20260506.md`
§ TE-41 cluster (15 UTs).

## Question log

Per AGENTS-ppx Question-logging discipline.

- [ ] **DF-nijab.1** What is `transports/`?
    opened: 2026-05-08 05:29 UTC
    asked of: stevegt@t7a.org (Steve Traugott)
    blocks: DR-suhod B/C/D cleanup and the group-session freeze procedure
    alternatives: 1.A lower-layer simulation surface / 1.B single-protocol directory
    recommendation: 1.A

- [ ] **DF-nijab.2** What happens to `wire-lab-devs-draft`?
    opened: 2026-05-08 05:29 UTC
    asked of: stevegt@t7a.org (Steve Traugott)
    blocks: transport-data migration and any freeze-time rewrite plan
    alternatives: 2.A preserve and supersede / 2.B derived rewrite only
    recommendation: 2.A

- [ ] **DF-nijab.3** How should freeze docs be corrected?
    opened: 2026-05-08 05:29 UTC
    asked of: stevegt@t7a.org (Steve Traugott)
    blocks: TODO-turog step 5, DR-suhod B/C/D, and stale spec references
    alternatives: 3.A correct now after DF / 3.B park behind TODO-pipus/TE-43
    recommendation: 3.A

## Decision Intent Log

(Will be populated as DFs lock and product lands.)
