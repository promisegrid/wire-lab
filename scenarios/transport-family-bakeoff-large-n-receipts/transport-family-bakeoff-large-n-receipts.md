# Large-N receipts

## Scenario ID

transport-family-bakeoff-large-n-receipts

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-narok-transport-family-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-narok-transport-family-bakeoff/`
- Source row/title: Large-N receipts
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-narok-transport-family-bakeoff/`.

## Setup

One thousand peers may acknowledge possession or processing of the same object.

## Stimulus

Run the candidate simulation against this source test: Whether receipts become vectors, compact summaries, per-peer promise records, or separate content-addressed proofs.

## Expected Pressure

Receipt scale may be a cross-family extension rather than a transport family by itself.
