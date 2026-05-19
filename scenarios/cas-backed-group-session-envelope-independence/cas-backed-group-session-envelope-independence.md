# Envelope independence

## Scenario ID

cas-backed-group-session-envelope-independence

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-jurar-cas-backed-group-session/SCENARIOS.md`
- Source simulation: `SIM-jurar-cas-backed-group-session/`
- Source row/title: Envelope independence
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-jurar-cas-backed-group-session/`.

## Setup

The message is wrapped by one candidate grid-envelope variant in one experiment and another in a different experiment.

## Stimulus

Run the candidate simulation against this source test: Whether group-session semantics depend only on resolved payload meaning, not on a chosen grid-envelope winner.

## Expected Pressure

This sim must not backdoor a preferred grid-envelope variant.
