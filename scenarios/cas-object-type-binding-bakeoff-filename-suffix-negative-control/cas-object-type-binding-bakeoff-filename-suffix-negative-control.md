# Filename suffix negative control

## Scenario ID

cas-object-type-binding-bakeoff-filename-suffix-negative-control

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-kohad-cas-object-type-binding-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-kohad-cas-object-type-binding-bakeoff/`
- Source row/title: Filename suffix negative control
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-kohad-cas-object-type-binding-bakeoff/`.

## Setup

Carol renames a local file from `.ptr` to `.raw` without changing bytes or CID.

## Stimulus

Run the candidate simulation against this source test: Whether path suffixes can safely carry type meaning in sparse replication, export/import, and archival storage.

## Expected Pressure

If suffix changes alter interpretation without changing content identity, suffixes are unsuitable as the primary discriminator.
