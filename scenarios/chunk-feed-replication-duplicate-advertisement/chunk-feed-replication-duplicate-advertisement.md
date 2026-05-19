# Duplicate advertisement

## Scenario ID

chunk-feed-replication-duplicate-advertisement

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-zazit-chunk-feed-replication/SCENARIOS.md`
- Source simulation: `SIM-zazit-chunk-feed-replication/`
- Source row/title: Duplicate advertisement
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-zazit-chunk-feed-replication/`.

## Setup

Alice and Carol both advertise chunk C.

## Stimulus

Run the candidate simulation against this source test: Whether duplicate offers are harmless and how Bob chooses between peers.

## Expected Pressure

Promise accounting can influence peer choice without making the feed a central reputation service.
