# Refused service

## Scenario ID

promise-accounting-records-refused-service

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-rusap-promise-accounting-records/SCENARIOS.md`
- Source simulation: `SIM-rusap-promise-accounting-records/`
- Source row/title: Refused service
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-rusap-promise-accounting-records/`.

## Setup

Alice refuses to send C because of policy, cost, group context, or missing authorization.

## Stimulus

Run the candidate simulation against this source test: Whether refusal is recorded differently from failure, corruption, or timeout.

## Expected Pressure

Promise accounting must support honest refusal instead of treating every refusal as misbehavior.
