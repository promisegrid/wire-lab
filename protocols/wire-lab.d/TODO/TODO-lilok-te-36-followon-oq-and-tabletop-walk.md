# TODO-lilok: TE-havib follow-on: OQ-36.6 + tabletop walk

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TODO-23` (integer alias)
- `TODO-20260507-002306` (timestamp alias and pre-migration filename)

## Status

Closed. Reopened 2026-05-09 because the file's prior closure depended on
a claimed TODO-vuhuj handoff for harness-spec sweep edits that was not
found in TODO-vuhuj during follow-up audit. That missing handoff is now
resolved by `TODO-kugod` / `DI-mugar`, which lands the harness-spec sweep
and closes the turn-159 residual rows that carried the sweep work.

All five UTs in this cluster were captured 2026-05-06 against the pre-Alt-B state of TE-havib. The 2026-05-05 Alt-B disposition path re-presentation locked DF-36.1, .2, .3, .7 (and resolved OQ-36.6 in the negative) before the cluster was filed. The verification walk confirmed all six tabletop scenarios hold against the locked state and all five UTs are answered by the locks (UT-160.b is procedural-meta, UT-160.c was wrong on inspection, UT-161.a is moot under DF-36.2 retirement, UT-162.a is resolved by OQ-36.6 negative resolution, UT-162.b is resolved by Alt-B re-presentation).

The original five-UT cluster has no live work remaining, and the actual
harness-spec sweep edits (§1.1, §2.1, §7.1, §10a.2/.3/.6, §10 table) are
now tracked and landed through `TODO-kugod` / `DI-mugar`.

## Threads absorbed from OPEN-THREADS.md

### T-TE36-FOLLOWON (formerly OPEN-THREADS, opened 2026-05-06)

OQ-36.6 deferred investigation (provisional DF-36.2 lock pending full
evidence review); audit-memo-style Alice-through-Mallory tabletop walk
that TE-havib did not run (UT-159.c / UT-160.c).

5 UTs in the dropped-thread-disposition file under the "TE-havib follow-on"
cluster.

Blocking: none. Harness-spec sweep ownership is resolved by `TODO-kugod` /
`DI-mugar`.

Anchor: TE-havib § OQ-36.6; audit memo at `4725b3e`.
Disposition-file pointer: `dropped-thread-disposition-20260506.md`
§ TE-havib follow-on cluster (5 UTs).

## Question log

(Per AGENTS-ppx Question-logging discipline. No questions logged yet.)

## Decision Intent Log

ID: DI-fahak
Date: 2026-05-09 10:58:02
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Reopen TODO-lilok despite the verified-superseded disposition of its original five-UT cluster.
Intent: A TODO should not be marked closed when it still points to unfinished harness-spec sweep work and the claimed downstream tracker is not evident in TODO-vuhuj.
Constraints: Preserve the 2026-05-07 verification-walk result as historical evidence; do not mark the original five UTs live again unless a later audit finds that their locks failed; track or retire the harness-spec sweep explicitly before closing TODO-lilok again.
Affects: `protocols/wire-lab.d/TODO/TODO-lilok-te-36-followon-oq-and-tabletop-walk.md`; `protocols/wire-lab.d/TODO/TODO.md`; `protocols/wire-lab.d/TODO/TODO-vuhuj-protocols-as-simulated-repos-migration.md`; `protocols/wire-lab.d/specs/harness-spec-draft.md`.
