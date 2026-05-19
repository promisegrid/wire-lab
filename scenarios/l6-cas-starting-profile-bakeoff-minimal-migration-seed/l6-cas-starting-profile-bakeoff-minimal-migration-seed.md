# Minimal migration seed

## Scenario ID

l6-cas-starting-profile-bakeoff-minimal-migration-seed

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-bobud-l6-cas-starting-profile-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-bobud-l6-cas-starting-profile-bakeoff/`
- Source row/title: Minimal migration seed
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-bobud-l6-cas-starting-profile-bakeoff/`.

## Setup

Alice migrates one historical inline group-session message into a pointer object plus raw CAS bytes.

## Stimulus

Run the candidate simulation against this source test: Whether raw chunks plus a minimal pointer object are enough to test sparse fetch, verification, and preservation of historical bytes.

## Expected Pressure

If the minimal profile proves the migration path, TE-43 can defer chunked Merkle complexity.
