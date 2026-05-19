# Session-layer composition

## Scenario ID

udp-feed-v0-conformance-session-layer-composition

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-kuful-udp-feed-v0-conformance/SCENARIOS.md`
- Source simulation: `SIM-kuful-udp-feed-v0-conformance/`
- Source row/title: Session-layer composition
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-kuful-udp-feed-v0-conformance/`.

## Setup

A minimal group/session message rides above UDP-feed v0.

## Stimulus

Run the candidate simulation against this source test: Whether UDP-feed's API is sufficient for the next layer without leaking binding details.

## Expected Pressure

If composition is required, TODO-jodon's done criteria must include more than UDP round trip.
