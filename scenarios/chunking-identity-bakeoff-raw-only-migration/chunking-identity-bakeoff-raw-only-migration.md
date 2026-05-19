# Raw-only migration

## Scenario ID

chunking-identity-bakeoff-raw-only-migration

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-gobaz-chunking-identity-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-gobaz-chunking-identity-bakeoff/`
- Source row/title: Raw-only migration
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-gobaz-chunking-identity-bakeoff/`.

## Setup

Alice migrates one historical message as a single raw chunk behind a pointer object.

## Stimulus

Run the candidate simulation against this source test: Whether the first CAS migration can proceed without chunked Merkle roots.

## Expected Pressure

Raw-only may unblock migration but does not answer large-object replication.
