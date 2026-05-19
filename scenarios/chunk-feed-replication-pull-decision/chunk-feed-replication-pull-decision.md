# Pull decision

## Scenario ID

chunk-feed-replication-pull-decision

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-zazit-chunk-feed-replication/SCENARIOS.md`
- Source simulation: `SIM-zazit-chunk-feed-replication/`
- Source row/title: Pull decision
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-zazit-chunk-feed-replication/`.

## Setup

Bob receives an advertisement for chunk C and has peer-local promise accounting records about Alice.

## Stimulus

Run the candidate simulation against this source test: Which inputs decide whether Bob pulls, delays, refuses, or asks another peer.

## Expected Pressure

The "decides" step needs an explicit cross-layer interface instead of a hand wave.
