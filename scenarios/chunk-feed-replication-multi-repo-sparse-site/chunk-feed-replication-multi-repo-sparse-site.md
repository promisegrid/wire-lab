# Multi-repo sparse site

## Scenario ID

chunk-feed-replication-multi-repo-sparse-site

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-zazit-chunk-feed-replication/SCENARIOS.md`
- Source simulation: `SIM-zazit-chunk-feed-replication/`
- Source row/title: Multi-repo sparse site
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-zazit-chunk-feed-replication/`.

## Setup

Alice's site state, Bob's site state, and a large shared corpus live in separate repos or fixtures.

## Stimulus

Run the candidate simulation against this source test: Whether feed promises and CAS object references remain meaningful when the harness orchestrates multiple storage roots.

## Expected Pressure

Turn 178's multi-repo question should be explored without assuming one repo contains every site's content.
