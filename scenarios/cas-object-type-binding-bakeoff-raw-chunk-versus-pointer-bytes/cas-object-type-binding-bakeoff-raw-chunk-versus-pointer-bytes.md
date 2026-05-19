# Raw chunk versus pointer bytes

## Scenario ID

cas-object-type-binding-bakeoff-raw-chunk-versus-pointer-bytes

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-kohad-cas-object-type-binding-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-kohad-cas-object-type-binding-bakeoff/`
- Source row/title: Raw chunk versus pointer bytes
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-kohad-cas-object-type-binding-bakeoff/`.

## Setup

Alice receives bytes whose hash is known, but the local filename is missing.

## Stimulus

Run the candidate simulation against this source test: Whether CID codec identity alone tells Bob whether to parse the bytes as a pointer object or treat them as raw payload.

## Expected Pressure

Object type must survive transport without relying on local paths.
