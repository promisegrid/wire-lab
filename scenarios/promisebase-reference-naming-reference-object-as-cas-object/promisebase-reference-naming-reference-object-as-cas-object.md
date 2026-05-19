# Reference object as CAS object

## Scenario ID

promisebase-reference-naming-reference-object-as-cas-object

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-ligan-promisebase-reference-naming/SCENARIOS.md`
- Source simulation: `SIM-ligan-promisebase-reference-naming/`
- Source row/title: Reference object as CAS object
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-ligan-promisebase-reference-naming/`.

## Setup

Alice writes a CBOR reference object whose CID R points at root X and includes a human-readable label.

## Stimulus

Run the candidate simulation against this source test: Whether the reference object is L6 CAS, L7 metadata, or a separate protocol object.

## Expected Pressure

Reference-object identity must not be confused with the target root identity.
