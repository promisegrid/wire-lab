# Missing pointee

## Scenario ID

cas-backed-group-session-missing-pointee

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-jurar-cas-backed-group-session/SCENARIOS.md`
- Source simulation: `SIM-jurar-cas-backed-group-session/`
- Source row/title: Missing pointee
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-jurar-cas-backed-group-session/`.

## Setup

Bob sees pointer object CID Y but lacks root CID X or some child chunks.

## Stimulus

Run the candidate simulation against this source test: Whether the group view can show pending / unresolved state without treating the message as invalid.

## Expected Pressure

Sparse-CAS behavior must be a normal group-session state.
