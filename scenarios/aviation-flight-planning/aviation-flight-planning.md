# Aviation Flight Planning

## Scenario ID

aviation-flight-planning

## Source / Provenance

- Source type: application seed
- Source path: `protocols/wire-lab.d/TODO/TODO-dadub-root-scenario-skeleton-and-seed-catalog.md`
- Source row/title: Seed application catalog entry `aviation-flight-planning`
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-midif`; `TODO-dadub`

## Purpose

Exercise PromiseGrid design candidates against aviation-flight-planning application
pressure: Flight plans, weather, aircraft status, crew qualifications, airspace
constraints, and dispatch release.

## Setup

Alice depends on an outcome in the Aviation Flight Planning domain. Bob makes promises
about flight plans, weather, aircraft status, crew qualifications, airspace constraints,
and dispatch release. Carol needs enough evidence to rely on Bob's promise without
having complete global state, and Mallory may exploit stale, missing, or ambiguous
evidence.

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
