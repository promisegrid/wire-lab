# Corrupt chunk

## Scenario ID

chunk-feed-replication-corrupt-chunk

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-zazit-chunk-feed-replication/SCENARIOS.md`
- Source simulation: `SIM-zazit-chunk-feed-replication/`
- Source row/title: Corrupt chunk
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-zazit-chunk-feed-replication/`.

## Setup

Mallory advertises or sends bytes whose hash does not match CID C.

## Stimulus

Run the candidate simulation against this source test: Whether rejection, accounting, and retry behavior are local enough to avoid central enforcement.

## Expected Pressure

Feed behavior must compose with CAS hash verification and peer-local accounting records.
