# Token rotates among writers

## Scenario ID

transport-family-bakeoff-token-rotates-among-writers

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-narok-transport-family-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-narok-transport-family-bakeoff/`
- Source row/title: Token rotates among writers
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-narok-transport-family-bakeoff/`.

## Setup

Alice, Bob, and Carol can publish only while holding a ring token.

## Stimulus

Run the candidate simulation against this source test: Whether token possession is a transport-level ordering promise, a group/session policy, or both.

## Expected Pressure

If token state is load-bearing for delivery semantics, a ring transport deserves its own specimen.
