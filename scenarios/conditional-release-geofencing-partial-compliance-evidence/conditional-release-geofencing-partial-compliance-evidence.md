# Partial compliance evidence

## Scenario ID

conditional-release-geofencing-partial-compliance-evidence

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-zarud-conditional-release-geofencing/SCENARIOS.md`
- Source simulation: `SIM-zarud-conditional-release-geofencing/`
- Source row/title: Partial compliance evidence
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-zarud-conditional-release-geofencing/`.

## Setup

Bob accepts onward-restraint but Carol's node cannot express the same condition vocabulary.

## Stimulus

Run the candidate simulation against this source test: Whether the protocol rejects forwarding, downgrades the condition, records partial compliance, or asks for a translation promise.

## Expected Pressure

Mixed-version behavior exposes whether the design needs explicit condition-version negotiation.
