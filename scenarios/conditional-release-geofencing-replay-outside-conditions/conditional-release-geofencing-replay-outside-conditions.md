# Replay outside conditions

## Scenario ID

conditional-release-geofencing-replay-outside-conditions

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-zarud-conditional-release-geofencing/SCENARIOS.md`
- Source simulation: `SIM-zarud-conditional-release-geofencing/`
- Source row/title: Replay outside conditions
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-zarud-conditional-release-geofencing/`.

## Setup

Mallory replays a valid old content reference to Dave outside the allowed audience or geography.

## Stimulus

Run the candidate simulation against this source test: Whether receivers, feeds, or group/session state detect stale or unauthorized reuse.

## Expected Pressure

Replay handling determines whether conditions must bind to recipients, epochs, locations, or session context.
