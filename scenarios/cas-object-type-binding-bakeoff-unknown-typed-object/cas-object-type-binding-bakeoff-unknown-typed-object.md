# Unknown typed object

## Scenario ID

cas-object-type-binding-bakeoff-unknown-typed-object

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-kohad-cas-object-type-binding-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-kohad-cas-object-type-binding-bakeoff/`
- Source row/title: Unknown typed object
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-kohad-cas-object-type-binding-bakeoff/`.

## Setup

Dave receives a CID whose codec he does not implement.

## Stimulus

Run the candidate simulation against this source test: Whether the peer can store, advertise, and forward the object opaquely while avoiding unsafe parsing.

## Expected Pressure

Type binding must define unknown-type behavior for long-lived mixed-version networks.
