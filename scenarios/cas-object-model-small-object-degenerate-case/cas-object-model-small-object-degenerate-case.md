# Small-object degenerate case

## Scenario ID

cas-object-model-small-object-degenerate-case

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-jomag-cas-object-model/SCENARIOS.md`
- Source simulation: `SIM-jomag-cas-object-model/`
- Source row/title: Small-object degenerate case
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-jomag-cas-object-model/`.

## Setup

A small group message fits in one chunk.

## Stimulus

Run the candidate simulation against this source test: Whether it still uses the same pointer / root / object-typing rules as a large object.

## Expected Pressure

The model should avoid special cases that create a second identity path for small messages.
