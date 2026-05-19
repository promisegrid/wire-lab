# ns-3 two-node path

## Scenario ID

udp-feed-v0-conformance-ns-3-two-node-path

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-kuful-udp-feed-v0-conformance/SCENARIOS.md`
- Source simulation: `SIM-kuful-udp-feed-v0-conformance/`
- Source row/title: ns-3 two-node path
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-kuful-udp-feed-v0-conformance/`.

## Setup

Alice and Bob communicate through an ns-3-emulated UDP network.

## Stimulus

Run the candidate simulation against this source test: Whether the v0 reference survives non-loopback timing, interface, and packet-capture conditions.

## Expected Pressure

ns-3 may be the evidence that separates a useful specimen from a local toy.
