# Cross-layer decision

## Scenario ID

promise-accounting-records-cross-layer-decision

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-rusap-promise-accounting-records/SCENARIOS.md`
- Source simulation: `SIM-rusap-promise-accounting-records/`
- Source row/title: Cross-layer decision
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-rusap-promise-accounting-records/`.

## Setup

L7 group policy says Bob values a root; L6 knows missing chunks; L5 sees offers from Alice and Carol.

## Stimulus

Run the candidate simulation against this source test: What information flows between layers when Bob decides which chunks to pull.

## Expected Pressure

The turn-178 "decides" issue should be made explicit without collapsing all accounting into one layer.
