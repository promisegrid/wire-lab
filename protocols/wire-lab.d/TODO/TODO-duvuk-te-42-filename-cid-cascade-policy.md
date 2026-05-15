# TODO-duvuk: TE-42 filename / CID-cascade policy

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TODO-27` (integer alias)
- `TODO-20260507-002306` (timestamp alias and pre-migration filename)

## Status

Closed as historical coordination memory. The active group-session successor
slice is closed by `TODO-gapab` / `DI-rurab`; the active feed-outer successor
slice is closed by `TODO-kakaz` / `DI-bomud`. This rooted file preserves the
original TE-42 thread narrative and provenance; it is no longer an active
execution owner.

## Threads absorbed from OPEN-THREADS.md

### T-FILENAME-CID-CASCADE (formerly OPEN-THREADS, opened 2026-05-06)

Codify the Path-A-vs-Path-B policy for legacy-message rehash on
envelope-changing edits (UT-169.a recommended Path A but commit
executed Path B); revisit Message-ID prohibition vs deprecation
(UT-169.e); resolve interaction with strict-reader rule § 4.7
(UT-169.b) and recursive-body parsing (UT-176.e); document the
CID-as-filename convention's stability properties.

Blocking: TE-41 (group-session freeze) lands first to establish the
freeze-and-rename mechanic this policy generalizes.

Anchor: transport-spec-draft.md § 4.3, § 4.7; group-session-draft.md
§ 9.
Disposition-file pointer: `dropped-thread-disposition-20260506.md`
§ TE-42 cluster (7 UTs).

Closure note, 2026-05-13: group-session message filename/CID,
`Message-ID:`, and reader/writer compatibility are closed by `DI-rurab`.
Feed-outer explicitly owns no message-header or rehash/deprecation rule and
keeps historical specimens immutable under `DI-bomud`.

## Rusis.10 successor owner routing

- Group-session envelope-side follow-on now lives at
  `simulations/SIM-rakot-group-session/protocols/group-session.d/TODO/TODO-gapab-group-session-freeze-and-cid-cascade-follow-on.md`.
- Feed-outer transport-side follow-on now lives at
  `simulations/SIM-labit-feed-outer/TODO/TODO-kakaz-feed-outer-freeze-and-cid-cascade-follow-on.md`.
- This rooted file keeps the original thread narrative and provenance
  as a historical record; active specimen-side ownership is transferred
  to those sim-local artifacts.

## Question log

(Per AGENTS-ppx Question-logging discipline. No questions logged yet.)

## Decision Intent Log

(Will be populated as DFs lock and product lands.)
