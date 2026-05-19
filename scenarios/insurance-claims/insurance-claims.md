# Insurance Claims

## Scenario ID

insurance-claims

## Source / Provenance

- Source type: application seed
- Source path: `protocols/wire-lab.d/TODO/TODO-dadub-root-scenario-skeleton-and-seed-catalog.md`
- Source row/title: Seed application catalog entry `insurance-claims`
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-midif`; `TODO-dadub`

## Purpose

Exercise PromiseGrid design candidates against insurance-claims application pressure:
Claim evidence, adjuster authority, fraud pressure, payments, and appeal promises.

## Setup

Alice depends on an outcome in the Insurance Claims domain. Bob makes promises about
claim evidence, adjuster authority, fraud pressure, payments, and appeal promises. Carol
needs enough evidence to rely on Bob's promise without having complete global state, and
Mallory may exploit stale, missing, or ambiguous evidence.

## Stimulus

A routine application event becomes contested or incomplete. Some relevant objects,
signatures, observations, or relationship records are delayed, partitioned, stale, or
disputed, and each peer must decide what to accept, retry, downgrade, or escalate using
only local evidence.

## Expected Pressure

The candidate simulation must show which promises, CAS objects, feeds, identity claims,
names, and promise accounting records are needed for this application pressure, while
avoiding hidden global state or a central authority that would make the result non-
comparable.
