# TODO-kituj: TE-43 promisebase prior-art adoption

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TODO-28` (integer alias)
- `TODO-20260507-002306` (timestamp alias and pre-migration filename)

## Status

Open. Depends on TE-sihih (substrate-agnostic layered model) L6
substrate definition. No twig yet.

## Threads absorbed from OPEN-THREADS.md

### T-PROMISEBASE-ADOPTION (formerly OPEN-THREADS, opened 2026-05-06)

Largest cluster after TE-sihih (25 UTs).

Scope:
  - Ratify the prototype-not-canon stance (Steve turn-191: "prototype
    at best; prefer wire-lab in conflict").
  - Document pitbase as L6 substrate prior art with 17/17 chunker-
    merkle tests green (UT-183.c).
  - Reconcile the kv branch on the remote (UT-190.d -- undiscovered
    through turn 192).
  - Resolve the Rabin-vs-FastCDC chunking parameter mismatch (UT-181.b:
    pitbase 512 KiB min / 8 MiB max vs turn-177 ~16 KiB average).
  - Fix the import-path error (UT-181.a: t7a/pitbase ->
    stevegt/promisebase).
  - Address fuse/ test failures and cmd/pb Docker SDK rot (UT-184.e/g).

Blocking: TE-sihih (substrate-agnostic layered model) L6 substrate
definition must land first.

Anchor: turn-191 canon rule; pitbase main + kv branch on
stevegt/promisebase.
Disposition-file pointer: `dropped-thread-disposition-20260506.md`
§ TE-43 cluster (25 UTs).

### T-021-CC-Q5 (carried forward, open)

promisebase canonicality: "wire-lab is canon, promisebase is
prototype" -- needs to be either ratified by a TE or formalized as a
canon rule somewhere readers will find. TE-43 is the natural home
since it locks the prototype-not-canon stance.

## Question log

(Per AGENTS-ppx Question-logging discipline. No questions logged yet.)

## Decision Intent Log

(Will be populated as DFs lock and product lands.)
