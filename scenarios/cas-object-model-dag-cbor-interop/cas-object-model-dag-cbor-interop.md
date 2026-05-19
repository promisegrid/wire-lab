# DAG-CBOR interop

## Scenario ID

cas-object-model-dag-cbor-interop

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-jomag-cas-object-model/SCENARIOS.md`
- Source simulation: `SIM-jomag-cas-object-model/`
- Source row/title: DAG-CBOR interop
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-jomag-cas-object-model/`.

## Setup

Alice stores a Merkle node or pointer object using a DAG-CBOR-compatible representation.

## Stimulus

Run the candidate simulation against this source test: Whether CID links, byte strings, and tags stay compatible with IPFS / IPLD-style tooling without requiring those stacks.

## Expected Pressure

TE-43 must decide whether DAG-CBOR is the default object format or only one allowed profile.
