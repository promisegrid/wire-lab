# Sparse advertisement

## Scenario ID

chunk-feed-replication-sparse-advertisement

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-zazit-chunk-feed-replication/SCENARIOS.md`
- Source simulation: `SIM-zazit-chunk-feed-replication/`
- Source row/title: Sparse advertisement
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-zazit-chunk-feed-replication/`.

## Setup

Alice has a subset of chunks for a Merkle root; Bob has a different subset.

## Stimulus

Run the candidate simulation against this source test: Whether the feed advertises leaves, roots, pointer objects, frontiers, or compact summaries without assuming full replication.

## Expected Pressure

Feed specs must work when no site has all CAS objects.
