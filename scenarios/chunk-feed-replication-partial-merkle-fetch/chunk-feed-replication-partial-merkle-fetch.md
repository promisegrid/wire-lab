# Partial Merkle fetch

## Scenario ID

chunk-feed-replication-partial-merkle-fetch

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-zazit-chunk-feed-replication/SCENARIOS.md`
- Source simulation: `SIM-zazit-chunk-feed-replication/`
- Source row/title: Partial Merkle fetch
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-zazit-chunk-feed-replication/`.

## Setup

Bob wants root R but only some children are locally available.

## Stimulus

Run the candidate simulation against this source test: Whether the feed can request missing children without understanding group-session message semantics.

## Expected Pressure

L5 should remain meaning-oblivious while still serving L6 CAS repair.
