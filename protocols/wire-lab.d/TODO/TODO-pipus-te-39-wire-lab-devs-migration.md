# TODO-pipus: TE-mumuv wire-lab-devs migration

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TODO-24` (integer alias)
- `TODO-20260507-002306` (timestamp alias and pre-migration filename)

## Status

Open. TE-sihih has landed, and the current `wire-lab-devs` specimen has moved
into simulation-local world state, but the operational migration from the
pre-CAS inline specimen to the TE-sihih-aligned CAS / pointer / feed shape
remains open. Current blockers are a concrete L6 CAS spec / TE-43 promisebase
adoption path and an additive migration contract; historical specimen messages
must not be rewritten. `SIM-jurar-cas-backed-group-session` now carries the
simulation-facing successor shape for group semantics over CAS roots and pointer
objects, with `SIM-jomag-cas-object-model` and
`SIM-zazit-chunk-feed-replication` carrying the object and feed pressures that
the migration must respect. Their scenario matrices now give TODO-pipus
concrete cases for additive migration, group-visible identity, parent links
through CAS, missing pointees, and feed-side sparse replication. Source:
TE-sihih; `DI-fakin`; `DI-rurab`; `DI-bomud`; `DI-navod`; `DI-pator`. No twig
yet.

## Threads absorbed from OPEN-THREADS.md

### T-MIG-OPS (formerly OPEN-THREADS, opened 2026-05-05)

Operational-shape TE for transport-protocol migration. Deferred
follow-on from TE-numan invariants-only TE.

Scope: 7 DFs (operational shape close-old-vs-overlap-vs-atomic-swap,
back-reference format, message disposition, authorizing promise, seal
mechanics, group-identity continuity, trigger discipline).

Turn-177 simulation input: `SIM-jurar-cas-backed-group-session` is the
successor-specimen charter for CAS-backed group-session migration. It does not
authorize rewriting historical `.txt` evidence; it gives TODO-pipus a
simulation-facing target to design against. Its `SCENARIOS.md` file, plus the
chunk-feed scenarios in `SIM-zazit-chunk-feed-replication`, are the current
scenario inputs for the first additive migration design. Source: `DI-navod`;
`DI-pator`.

Blocking: nothing today; gated on a concrete first migration to design
against. TE-sihih + TE-mumuv close the gate.

Anchor: TE-numan § Implications, transport-spec § OQ-2.

### T-021-CC-Q1 (carried forward, answered)

What goes into TE-numan? Now answered by 2026-05-05 TE-numan merge: TE-numan
became transport-protocol migration invariants; the
promisebase-integration scope moved to TE-43.

## Question log

(Per AGENTS-ppx Question-logging discipline. No questions logged yet.)

## Decision Intent Log

(Will be populated as DFs lock and product lands.)
