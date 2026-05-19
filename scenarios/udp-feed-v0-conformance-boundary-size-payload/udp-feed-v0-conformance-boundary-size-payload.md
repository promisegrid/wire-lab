# Boundary-size payload

## Scenario ID

udp-feed-v0-conformance-boundary-size-payload

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-kuful-udp-feed-v0-conformance/SCENARIOS.md`
- Source simulation: `SIM-kuful-udp-feed-v0-conformance/`
- Source row/title: Boundary-size payload
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-kuful-udp-feed-v0-conformance/`.

## Setup

Alice sends exactly 1232 bytes, then 1233 bytes.

## Stimulus

Run the candidate simulation against this source test: Whether the implementation honors the size promise and errors locally before oversize send.

## Expected Pressure

Size behavior should be in vectors before wider conformance claims.
