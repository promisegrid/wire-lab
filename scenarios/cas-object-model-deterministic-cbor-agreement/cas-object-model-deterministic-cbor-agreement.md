# Deterministic CBOR agreement

## Scenario ID

cas-object-model-deterministic-cbor-agreement

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-jomag-cas-object-model/SCENARIOS.md`
- Source simulation: `SIM-jomag-cas-object-model/`
- Source row/title: Deterministic CBOR agreement
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-jomag-cas-object-model/`.

## Setup

Alice and Bob encode the same pointer object with independent implementations.

## Stimulus

Run the candidate simulation against this source test: Whether map ordering, integer/string choices, tags, and byte-string boundaries produce identical bytes and therefore identical CIDs.

## Expected Pressure

TE-43 must lock a precise CBOR profile rather than saying "use CBOR" generically.
