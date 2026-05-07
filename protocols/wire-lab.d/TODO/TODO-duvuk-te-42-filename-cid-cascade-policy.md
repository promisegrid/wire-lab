# TODO 27: TE-42 filename / CID-cascade policy

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TODO-27` (integer alias)
- `TODO-20260507-002306` (timestamp alias and pre-migration filename)

## Status

Open. Depends on TE-41 (group-session freeze) to establish the
freeze-and-rename mechanic. No twig yet.

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

## Question log

(Per AGENTS-ppx Question-logging discipline. No questions logged yet.)

## Decision Intent Log

(Will be populated as DFs lock and product lands.)
