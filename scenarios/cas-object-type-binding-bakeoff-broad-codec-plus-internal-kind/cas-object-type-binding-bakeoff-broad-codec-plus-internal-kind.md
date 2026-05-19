# Broad codec plus internal kind

## Scenario ID

cas-object-type-binding-bakeoff-broad-codec-plus-internal-kind

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-kohad-cas-object-type-binding-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-kohad-cas-object-type-binding-bakeoff/`
- Source row/title: Broad codec plus internal kind
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-kohad-cas-object-type-binding-bakeoff/`.

## Setup

Alice stores a CBOR object under a broad PromiseGrid object codec, and the object bytes include `kind = pointer`.

## Stimulus

Run the candidate simulation against this source test: Whether internal kind improves forward compatibility or merely duplicates a CID-level type claim.

## Expected Pressure

TE-43 must avoid two independent type authorities unless the split has a clear rule.
