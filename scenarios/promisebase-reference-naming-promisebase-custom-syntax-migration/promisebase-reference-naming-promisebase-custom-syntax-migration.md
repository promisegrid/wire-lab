# Promisebase custom syntax migration

## Scenario ID

promisebase-reference-naming-promisebase-custom-syntax-migration

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-ligan-promisebase-reference-naming/SCENARIOS.md`
- Source simulation: `SIM-ligan-promisebase-reference-naming/`
- Source row/title: Promisebase custom syntax migration
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-ligan-promisebase-reference-naming/`.

## Setup

A promisebase-era reference uses non-CID custom syntax for a root.

## Stimulus

Run the candidate simulation against this source test: Whether migration wraps it, rejects it, or maps it into CID-backed reference objects.

## Expected Pressure

Prior-art adoption must be deliberate and not preserve known-bad syntax by accident.
