# Opaque lower-layer carriage

## Scenario ID

conditional-release-geofencing-opaque-lower-layer-carriage

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-zarud-conditional-release-geofencing/SCENARIOS.md`
- Source simulation: `SIM-zarud-conditional-release-geofencing/`
- Source row/title: Opaque lower-layer carriage
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-zarud-conditional-release-geofencing/`.

## Setup

Bob's node stores encrypted content whose condition vocabulary it cannot parse.

## Stimulus

Run the candidate simulation against this source test: Whether lower layers can safely carry opaque condition references while avoiding accidental promise violations.

## Expected Pressure

If opaque carriage is acceptable, the condition object must be verifiable without every layer understanding its semantics.
