# Outer/feed split

## Scenario ID

group-session-freeze-promise-outer-feed-split

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-bohof-group-session-freeze-promise/SCENARIOS.md`
- Source simulation: `SIM-bohof-group-session-freeze-promise/`
- Source row/title: Outer/feed split
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-bohof-group-session-freeze-promise/`.

## Setup

The group-session semantics are stable, but a lower outer/feed rule remains unsettled.

## Stimulus

Run the candidate simulation against this source test: Whether the freeze promise can split scopes without pretending both layers are solved.

## Expected Pressure

The merge promise must name exactly what is frozen and what remains provisional.
