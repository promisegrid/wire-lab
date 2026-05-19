# Floating exchange rates

## Scenario ID

promise-economy-spectrum-floating-exchange-rates

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-haros-promise-economy-spectrum/SCENARIOS.md`
- Source simulation: `SIM-haros-promise-economy-spectrum/`
- Source row/title: Floating exchange rates
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-haros-promise-economy-spectrum/`.

## Setup

Alice accepts Bob's storage promises at one rate and Carol's at another.

## Stimulus

Run the candidate simulation against this source test: Whether peers can keep local valuation without a central price oracle.

## Expected Pressure

"Everyone is their own central bank" must not become a hidden global currency.
