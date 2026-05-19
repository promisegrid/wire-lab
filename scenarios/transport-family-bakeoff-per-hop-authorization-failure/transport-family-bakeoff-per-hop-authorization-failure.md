# Per-hop authorization failure

## Scenario ID

transport-family-bakeoff-per-hop-authorization-failure

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-narok-transport-family-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-narok-transport-family-bakeoff/`
- Source row/title: Per-hop authorization failure
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-narok-transport-family-bakeoff/`.

## Setup

Bob receives a ring message but is not authorized to forward it to Carol.

## Stimulus

Run the candidate simulation against this source test: Whether authorization failure breaks the ring, skips a hop, records refusal, or reconfigures membership.

## Expected Pressure

Ring semantics need a failure model before they can be compared with gossip.
