# Long-horizon reprofile

## Scenario ID

l6-cas-starting-profile-bakeoff-long-horizon-reprofile

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-bobud-l6-cas-starting-profile-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-bobud-l6-cas-starting-profile-bakeoff/`
- Source row/title: Long-horizon reprofile
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-bobud-l6-cas-starting-profile-bakeoff/`.

## Setup

A future implementation wants to replace the first profile with a richer object graph.

## Stimulus

Run the candidate simulation against this source test: Whether old pointer objects and raw chunks remain addressable and explainable after a later profile lands.

## Expected Pressure

The starting profile should avoid identity choices that become migration debt.
