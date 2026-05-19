# Mixed-version peer

## Scenario ID

l6-cas-starting-profile-bakeoff-mixed-version-peer

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-bobud-l6-cas-starting-profile-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-bobud-l6-cas-starting-profile-bakeoff/`
- Source row/title: Mixed-version peer
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-bobud-l6-cas-starting-profile-bakeoff/`.

## Setup

Alice uses the minimal profile while Bob experiments with DAG-CBOR Merkle nodes.

## Stimulus

Run the candidate simulation against this source test: Whether peers can reject, store opaquely, or bridge objects whose profile they do not yet implement.

## Expected Pressure

The first profile needs an extension story even if it starts small.
