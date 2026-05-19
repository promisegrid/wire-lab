# pCID-driven object

## Scenario ID

chunking-identity-bakeoff-pcid-driven-object

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-gobaz-chunking-identity-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-gobaz-chunking-identity-bakeoff/`
- Source row/title: pCID-driven object
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-gobaz-chunking-identity-bakeoff/`.

## Setup

Alice stores an object whose governing pCID defines the chunker. Bob has the object root but not the spec text cached.

## Stimulus

Run the candidate simulation against this source test: Whether chunk verification and fetch planning require resolving the pCID first.

## Expected Pressure

pCID-driven chunking couples object interpretation to protocol-spec availability.
