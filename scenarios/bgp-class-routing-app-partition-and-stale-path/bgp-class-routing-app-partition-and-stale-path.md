# Partition and stale path

## Scenario ID

bgp-class-routing-app-partition-and-stale-path

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-punaz-bgp-class-routing-app/SCENARIOS.md`
- Source simulation: `SIM-punaz-bgp-class-routing-app/`
- Source row/title: Partition and stale path
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-punaz-bgp-class-routing-app/`.

## Setup

A formerly good path becomes unavailable during intermittent connectivity.

## Stimulus

Run the candidate simulation against this source test: How stale promises, timeouts, and withdrawal notices affect local decisions.

## Expected Pressure

Long-lived routing records need aging and repair without central convergence machinery.
