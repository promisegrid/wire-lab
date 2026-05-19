# Offline verification

## Scenario ID

peer-adoption-metadata-offline-verification

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-dihiz-peer-adoption-metadata/SCENARIOS.md`
- Source simulation: `SIM-dihiz-peer-adoption-metadata/`
- Source row/title: Offline verification
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-dihiz-peer-adoption-metadata/`.

## Setup

Alice receives Bob's adoption record through Dave while Bob is offline.

## Stimulus

Run the candidate simulation against this source test: Whether third-party carriage works without a central registry.

## Expected Pressure

Content-addressed adoption objects get stronger if offline verification is a first-class scenario.
