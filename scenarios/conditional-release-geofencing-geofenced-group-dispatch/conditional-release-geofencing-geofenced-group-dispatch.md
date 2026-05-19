# Geofenced group dispatch

## Scenario ID

conditional-release-geofencing-geofenced-group-dispatch

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-zarud-conditional-release-geofencing/SCENARIOS.md`
- Source simulation: `SIM-zarud-conditional-release-geofencing/`
- Source row/title: Geofenced group dispatch
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-zarud-conditional-release-geofencing/`.

## Setup

Alice permits content only for group members inside a stated region. Carol is a member but is outside the allowed region.

## Stimulus

Run the candidate simulation against this source test: Whether geofence checks are membership checks, per-message dispatch checks, fetch-policy checks, or storage constraints.

## Expected Pressure

The owner layer must explain both refusal and auditability without assuming a central location oracle.
