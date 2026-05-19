# Loopback round trip

## Scenario ID

udp-feed-v0-conformance-loopback-round-trip

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-kuful-udp-feed-v0-conformance/SCENARIOS.md`
- Source simulation: `SIM-kuful-udp-feed-v0-conformance/`
- Source row/title: Loopback round trip
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-kuful-udp-feed-v0-conformance/`.

## Setup

Alice sends a 612-byte payload to Bob over local UDP.

## Stimulus

Run the candidate simulation against this source test: Whether the reference implementation preserves bytes and exposes the expected send/receive API.

## Expected Pressure

A minimal reference may be enough for first v0 evidence if vectors lock the bytes.
