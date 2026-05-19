# Gossip convergence

## Scenario ID

transport-family-bakeoff-gossip-convergence

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-narok-transport-family-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-narok-transport-family-bakeoff/`
- Source row/title: Gossip convergence
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-narok-transport-family-bakeoff/`.

## Setup

Alice emits a message; Bob, Carol, and Dave learn it through epidemic propagation with duplicate paths.

## Stimulus

Run the candidate simulation against this source test: Whether convergence evidence is enough without total ordering or delivery promises.

## Expected Pressure

Gossip is attractive only if duplicate suppression and missing-object repair stay explainable.
