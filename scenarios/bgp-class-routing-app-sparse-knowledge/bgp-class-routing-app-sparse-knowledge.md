# Sparse knowledge

## Scenario ID

bgp-class-routing-app-sparse-knowledge

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-punaz-bgp-class-routing-app/SCENARIOS.md`
- Source simulation: `SIM-punaz-bgp-class-routing-app/`
- Source row/title: Sparse knowledge
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-punaz-bgp-class-routing-app/`.

## Setup

Alice knows only Bob and Carol; Bob knows Dave and Ellen; no peer has the whole graph.

## Stimulus

Run the candidate simulation against this source test: Whether multi-hop discovery can find acceptable paths without requiring full topology replication.

## Expected Pressure

BGP-class pressure must compose with sparse-CAS and sparse relationship knowledge.
