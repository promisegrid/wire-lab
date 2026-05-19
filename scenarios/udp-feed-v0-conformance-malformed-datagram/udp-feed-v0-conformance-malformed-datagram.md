# Malformed datagram

## Scenario ID

udp-feed-v0-conformance-malformed-datagram

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-kuful-udp-feed-v0-conformance/SCENARIOS.md`
- Source simulation: `SIM-kuful-udp-feed-v0-conformance/`
- Source row/title: Malformed datagram
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-kuful-udp-feed-v0-conformance/`.

## Setup

Bob receives arbitrary bytes that do not parse at higher layers.

## Stimulus

Run the candidate simulation against this source test: Whether UDP-feed passes bytes upward unchanged rather than inventing message semantics.

## Expected Pressure

Binding conformance must stay below session semantics.
