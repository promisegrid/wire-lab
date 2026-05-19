# Route hijack

## Scenario ID

bgp-class-routing-app-route-hijack

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-punaz-bgp-class-routing-app/SCENARIOS.md`
- Source simulation: `SIM-punaz-bgp-class-routing-app/`
- Source row/title: Route hijack
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-punaz-bgp-class-routing-app/`.

## Setup

Mallory advertises a short attractive path to Carol but cannot deliver traffic or chunks.

## Stimulus

Run the candidate simulation against this source test: How peers detect failed promises and locally downgrade future route choices.

## Expected Pressure

PromiseGrid routing apps need hijack costs that are local and relationship-specific.
