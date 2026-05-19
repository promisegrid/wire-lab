# Carrier independence

## Scenario ID

chunk-feed-replication-carrier-independence

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-zazit-chunk-feed-replication/SCENARIOS.md`
- Source simulation: `SIM-zazit-chunk-feed-replication/`
- Source row/title: Carrier independence
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-zazit-chunk-feed-replication/`.

## Setup

The same chunk exchange is attempted over UDP, git, libp2p, IPFS, or ATPROTO-adjacent carriers.

## Stimulus

Run the candidate simulation against this source test: Which semantics belong to the feed role and which are carrier mechanics.

## Expected Pressure

The simulation should preserve turn-177's claim that feeds move chunks independent of substrate.
