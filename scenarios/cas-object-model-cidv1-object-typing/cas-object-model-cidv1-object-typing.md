# CIDv1 object typing

## Scenario ID

cas-object-model-cidv1-object-typing

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-jomag-cas-object-model/SCENARIOS.md`
- Source simulation: `SIM-jomag-cas-object-model/`
- Source row/title: CIDv1 object typing
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-jomag-cas-object-model/`.

## Setup

The same content hash could be interpreted as a raw chunk, Merkle node, or pointer object unless type is bound into identity.

## Stimulus

Run the candidate simulation against this source test: Whether CIDv1 codec / multicodec values carry object type cleanly enough to avoid filename suffixes.

## Expected Pressure

TE-43 must lock object typing through CID codecs or explicitly choose another type-binding mechanism.
