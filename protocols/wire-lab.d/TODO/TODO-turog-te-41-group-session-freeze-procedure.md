# TODO 26: TE-41 group-session freeze procedure

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TODO-26` (integer alias)
- `TODO-20260507-002306` (timestamp alias and pre-migration filename)

## Status

Open. Depends on TE-40 (apparatus-vs-specimen completion) landing.
No twig yet.

## Threads absorbed from OPEN-THREADS.md

### T-GROUP-SESSION-FREEZE (formerly OPEN-THREADS, opened 2026-05-06)

Freeze chain for group-session and outer transport-spec. Carved from
T-021-TODO12 closure.

Scope: TODO 12 subtask 012.8 lists this work but it is wider than
TODO 12. Sequence:
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

Blocking: outer transport-spec-draft.md has its own freeze gate;
nothing else gates this thread today.

Anchor: group-session-draft.md § Freeze gate; transport-spec-draft.md
§ Freeze gate; TODO 12 subtask 012.8.
Disposition-file pointer: `dropped-thread-disposition-20260506.md`
§ TE-41 cluster (15 UTs).

## Question log

(Per AGENTS-ppx Question-logging discipline. No questions logged yet.)

## Decision Intent Log

(Will be populated as DFs lock and product lands.)
