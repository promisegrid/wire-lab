# TODO 29: TE-44 wire-lab/promisebase merge trajectory

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TODO-29` (integer alias)
- `TODO-20260507-002306` (timestamp alias and pre-migration filename)

## Status

Open. Depends on TE-43 (promisebase prior-art adoption). No twig
yet.

## Threads absorbed from OPEN-THREADS.md

### T-WIRELAB-PROMISEBASE-MERGE (formerly OPEN-THREADS, opened 2026-05-06)

Zero direct UTs but UT-192.f's "will gradually merge" framing requires
its own DF lock.

Scope: surface the implicit "wire-lab and promisebase will be in the
same codebase eventually" architectural commitment as a deliberate DF;
name the merge endpoint, the staging mechanic, and the canonical-vs-
prototype policy at merge time.

Blocking: TE-43 (promisebase prior-art adoption) lands first.

Anchor: UT-192.f; TE-43 future cross-reference.
Disposition-file pointer: `dropped-thread-disposition-20260506.md`
§ Proposed TE roster (Phase 2 work).

## Question log

(Per AGENTS-ppx Question-logging discipline. No questions logged yet.)

## Decision Intent Log

(Will be populated as DFs lock and product lands.)
