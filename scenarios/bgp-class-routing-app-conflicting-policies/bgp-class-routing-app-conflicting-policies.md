# Conflicting policies

## Scenario ID

bgp-class-routing-app-conflicting-policies

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-punaz-bgp-class-routing-app/SCENARIOS.md`
- Source simulation: `SIM-punaz-bgp-class-routing-app/`
- Source row/title: Conflicting policies
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-punaz-bgp-class-routing-app/`.

## Setup

Alice prefers paths that avoid a jurisdiction; Carol prefers cheapest path; Bob has both offers.

## Stimulus

Run the candidate simulation against this source test: Whether route choice can be policy-relative instead of globally best.

## Expected Pressure

A PromiseGrid routing app should support peer-specific preference and refusal semantics.
