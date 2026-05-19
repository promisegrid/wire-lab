# Name collision

## Scenario ID

promisebase-reference-naming-name-collision

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-ligan-promisebase-reference-naming/SCENARIOS.md`
- Source simulation: `SIM-ligan-promisebase-reference-naming/`
- Source row/title: Name collision
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-ligan-promisebase-reference-naming/`.

## Setup

Alice and Bob both publish `release` pointing at different roots.

## Stimulus

Run the candidate simulation against this source test: Whether names are scoped by peer, group, pCID, site, or another authority.

## Expected Pressure

Human-readable names need scope rules or they become a new central registry problem.
