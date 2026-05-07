# TODO-lilok: TE-havib follow-on: OQ-36.6 + tabletop walk

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TODO-23` (integer alias)
- `TODO-20260507-002306` (timestamp alias and pre-migration filename)

## Status

Closed 2026-05-07 — **verified-superseded** under Alt-lilok.1.B disposition path. See verification walk memo at [`protocols/wire-lab.d/docs/te-havib-followon-verification-walk-20260507.md`](../../wire-lab.d/docs/te-havib-followon-verification-walk-20260507.md).

All five UTs in this cluster were captured 2026-05-06 against the pre-Alt-B state of TE-havib. The 2026-05-05 Alt-B disposition path re-presentation locked DF-36.1, .2, .3, .7 (and resolved OQ-36.6 in the negative) before the cluster was filed. The verification walk confirmed all six tabletop scenarios hold against the locked state and all five UTs are answered by the locks (UT-160.b is procedural-meta, UT-160.c was wrong on inspection, UT-161.a is moot under DF-36.2 retirement, UT-162.a is resolved by OQ-36.6 negative resolution, UT-162.b is resolved by Alt-B re-presentation).

No live work remains. The actual harness-spec sweep edits (§1.1, §2.1, §7.1, §10a.2/.3/.6, §10 table) are tracked under TODO-vuhuj's leftover sweep per DF-36.7 lock annotation, not here.

## Threads absorbed from OPEN-THREADS.md

### T-TE36-FOLLOWON (formerly OPEN-THREADS, opened 2026-05-06)

OQ-36.6 deferred investigation (provisional DF-36.2 lock pending full
evidence review); audit-memo-style Alice-through-Mallory tabletop walk
that TE-havib did not run (UT-159.c / UT-160.c).

5 UTs in the dropped-thread-disposition file under the "TE-havib follow-on"
cluster.

Blocking: nothing today; pure follow-on refinement.

Anchor: TE-havib § OQ-36.6; audit memo at `4725b3e`.
Disposition-file pointer: `dropped-thread-disposition-20260506.md`
§ TE-havib follow-on cluster (5 UTs).

## Question log

(Per AGENTS-ppx Question-logging discipline. No questions logged yet.)

## Decision Intent Log

(Will be populated as DFs lock and product lands.)
