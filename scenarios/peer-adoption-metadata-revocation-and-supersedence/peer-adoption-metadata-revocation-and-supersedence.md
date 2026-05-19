# Revocation and supersedence

## Scenario ID

peer-adoption-metadata-revocation-and-supersedence

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-dihiz-peer-adoption-metadata/SCENARIOS.md`
- Source simulation: `SIM-dihiz-peer-adoption-metadata/`
- Source row/title: Revocation and supersedence
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-dihiz-peer-adoption-metadata/`.

## Setup

Bob stops following pCID X and adopts pCID Y.

## Stimulus

Run the candidate simulation against this source test: Whether the old claim remains auditable while the current claim becomes discoverable.

## Expected Pressure

The design needs explicit current-pointer, supersedence, or freshness semantics.
