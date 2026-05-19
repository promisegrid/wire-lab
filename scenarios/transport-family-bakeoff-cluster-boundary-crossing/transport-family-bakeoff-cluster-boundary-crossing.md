# Cluster boundary crossing

## Scenario ID

transport-family-bakeoff-cluster-boundary-crossing

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-narok-transport-family-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-narok-transport-family-bakeoff/`
- Source row/title: Cluster boundary crossing
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-narok-transport-family-bakeoff/`.

## Setup

Alice sends from cluster A to Dave in cluster B through cluster representatives.

## Stimulus

Run the candidate simulation against this source test: Whether clusters are addressable principals, routing hints, policy scopes, or aggregation boundaries.

## Expected Pressure

Cluster-of-clusters transport needs explicit promises at each boundary.
