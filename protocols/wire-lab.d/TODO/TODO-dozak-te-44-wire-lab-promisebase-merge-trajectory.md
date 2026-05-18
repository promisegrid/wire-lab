# TODO-dozak: TE-44 wire-lab/promisebase merge trajectory

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

Turn 191's canon rule applies to merge-time decisions too: any convergence,
merge, or long-lived divergence path must preserve wire-lab as design canon
unless Steve explicitly locks a different outcome, and must document conflicts
rather than silently importing promisebase design. Source: `DI-sapiv`.

Turn 179's proposal to add PromiseGrid narrative material to promisebase docs is
treated here as a scope-creep warning, not an authorization. Cross-repo
promisebase documentation, merge, or convergence work needs its own later DF or
explicit task; it must not be inferred from TE-sihih or from the invalidated
wholesale-adoption pivot. Source: `DI-vabij`.

Blocking: TE-43 (promisebase prior-art adoption) lands first.

Anchor: UT-192.f; TE-43 future cross-reference.
Disposition-file pointer: `dropped-thread-disposition-20260506.md`
§ Proposed TE roster (Phase 2 work).

## Question log

- 2026-05-17: Turn 179's cross-repo promisebase-docs proposal is routed here as
  merge-trajectory / scope-boundary pressure. This does not authorize any
  promisebase edit. Source: `DI-vabij`.
- 2026-05-17: Turn 191's prototype-not-canon rule is routed here as the
  merge-time canon boundary: convergence must document conflicts and keep
  wire-lab authoritative unless Steve locks an exception. Source: `DI-sapiv`.

## Subtasks

- [ ] dozak.1 Decide whether wire-lab and promisebase should converge, merge, or
  stay independent before any cross-repo promisebase documentation or code
  trajectory is treated as PromiseGrid plan of record. Source: `DI-vabij`;
  `DI-sapiv`.

## Decision Intent Log

(Will be populated as DFs lock and product lands.)
