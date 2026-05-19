# Chunker parameter mismatch

## Scenario ID

cas-object-model-chunker-parameter-mismatch

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-jomag-cas-object-model/SCENARIOS.md`
- Source simulation: `SIM-jomag-cas-object-model/`
- Source row/title: Chunker parameter mismatch
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-jomag-cas-object-model/`.

## Setup

Alice chunks with turn-177 FastCDC-style small averages; Bob chunks with promisebase / pitbase Rabin defaults.

## Stimulus

Run the candidate simulation against this source test: Whether the same logical file produces different leaf CIDs and Merkle roots under different parameters.

## Expected Pressure

TE-43 must lock chunking algorithm and full parameter set or make parameterized chunking explicit in object identity.
