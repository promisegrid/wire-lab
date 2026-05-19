# Published mutable ref

## Scenario ID

promisebase-reference-naming-published-mutable-ref

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-ligan-promisebase-reference-naming/SCENARIOS.md`
- Source simulation: `SIM-ligan-promisebase-reference-naming/`
- Source row/title: Published mutable ref
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-ligan-promisebase-reference-naming/`.

## Setup

Alice publishes `project/latest` first pointing at root X and later at root Y.

## Stimulus

Run the candidate simulation against this source test: How signed update history, replay protection, and reader expectations work.

## Expected Pressure

Mutable refs are not CAS roots; they need explicit update semantics if adopted.
