# Corrupt data

## Scenario ID

promise-accounting-records-corrupt-data

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-rusap-promise-accounting-records/SCENARIOS.md`
- Source simulation: `SIM-rusap-promise-accounting-records/`
- Source row/title: Corrupt data
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-rusap-promise-accounting-records/`.

## Setup

Mallory sends bytes that fail CID verification.

## Stimulus

Run the candidate simulation against this source test: What observation is recorded and how future local decisions change.

## Expected Pressure

CAS verification and peer-local records should compose without central enforcement.
