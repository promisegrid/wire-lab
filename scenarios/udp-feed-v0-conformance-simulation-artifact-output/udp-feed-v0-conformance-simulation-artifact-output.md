# Simulation artifact output

## Scenario ID

udp-feed-v0-conformance-simulation-artifact-output

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-kuful-udp-feed-v0-conformance/SCENARIOS.md`
- Source simulation: `SIM-kuful-udp-feed-v0-conformance/`
- Source row/title: Simulation artifact output
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-kuful-udp-feed-v0-conformance/`.

## Setup

A simulator-mode send writes an artifact file for the transmitted bytes.

## Stimulus

Run the candidate simulation against this source test: Whether artifact output proves promise 10 without becoming production behavior.

## Expected Pressure

The artifact contract needs to be testable and explicitly scoped.
