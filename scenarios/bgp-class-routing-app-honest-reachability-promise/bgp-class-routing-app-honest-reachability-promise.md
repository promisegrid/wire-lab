# Honest reachability promise

## Scenario ID

bgp-class-routing-app-honest-reachability-promise

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-punaz-bgp-class-routing-app/SCENARIOS.md`
- Source simulation: `SIM-punaz-bgp-class-routing-app/`
- Source row/title: Honest reachability promise
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-punaz-bgp-class-routing-app/`.

## Setup

Alice advertises reachability to Carol through Bob, and Bob later forwards as promised.

## Stimulus

Run the candidate simulation against this source test: What Alice, Bob, and Carol each record locally after the path works.

## Expected Pressure

Route-like promises need observable kept/broken outcomes without a global route authority.
