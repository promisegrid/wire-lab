# Arbitrary body shape

## Scenario ID

cas-backed-group-session-arbitrary-body-shape

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-jurar-cas-backed-group-session/SCENARIOS.md`
- Source simulation: `SIM-jurar-cas-backed-group-session/`
- Source row/title: Arbitrary body shape
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-jurar-cas-backed-group-session/`.

## Setup

Alice's message body is a CBOR text string, CBOR map, encrypted blob, signed payload, or large file root.

## Stimulus

Run the candidate simulation against this source test: Whether group semantics can stay stable while body bytes vary.

## Expected Pressure

TODO-pipus and TE-43 must agree on what the group layer sees versus what CAS stores.
