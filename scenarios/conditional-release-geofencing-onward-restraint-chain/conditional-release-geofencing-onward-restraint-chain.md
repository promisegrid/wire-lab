# Onward-restraint chain

## Scenario ID

conditional-release-geofencing-onward-restraint-chain

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-zarud-conditional-release-geofencing/SCENARIOS.md`
- Source simulation: `SIM-zarud-conditional-release-geofencing/`
- Source row/title: Onward-restraint chain
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-zarud-conditional-release-geofencing/`.

## Setup

Alice sends content to Bob only if Bob promises to forward it only to recipients who make the same promise. Bob wants to forward to Carol.

## Stimulus

Run the candidate simulation against this source test: Whether the recursive promise graph is represented at group-session, conditional-release, transport/feed, or CAS-object level.

## Expected Pressure

If the graph is central to dispatch semantics, group/session ownership gets stronger; if it composes across sessions, a separate family gets stronger.
