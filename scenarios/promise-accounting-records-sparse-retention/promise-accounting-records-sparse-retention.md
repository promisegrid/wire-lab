# Sparse retention

## Scenario ID

promise-accounting-records-sparse-retention

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-rusap-promise-accounting-records/SCENARIOS.md`
- Source simulation: `SIM-rusap-promise-accounting-records/`
- Source row/title: Sparse retention
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-rusap-promise-accounting-records/`.

## Setup

Bob cannot keep all chunks and must choose what to retain.

## Stimulus

Run the candidate simulation against this source test: Whether local records can include promises, costs, group value, and peer reliability.

## Expected Pressure

Sparse-CAS makes retention a policy decision, not a background storage detail.
