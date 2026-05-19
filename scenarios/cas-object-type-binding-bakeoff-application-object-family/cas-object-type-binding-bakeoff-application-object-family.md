# Application object family

## Scenario ID

cas-object-type-binding-bakeoff-application-object-family

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-kohad-cas-object-type-binding-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-kohad-cas-object-type-binding-bakeoff/`
- Source row/title: Application object family
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-kohad-cas-object-type-binding-bakeoff/`.

## Setup

Ellen proposes a future application-level CAS object distinct from raw chunks, Merkle nodes, and pointer objects.

## Stimulus

Run the candidate simulation against this source test: Whether the chosen binding model leaves room for new object families without reinterpreting old bytes.

## Expected Pressure

The first type-binding rule should be extensible without changing old CIDs.
