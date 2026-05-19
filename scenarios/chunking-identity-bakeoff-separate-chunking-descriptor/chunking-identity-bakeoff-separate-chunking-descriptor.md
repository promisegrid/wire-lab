# Separate chunking descriptor

## Scenario ID

chunking-identity-bakeoff-separate-chunking-descriptor

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-gobaz-chunking-identity-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-gobaz-chunking-identity-bakeoff/`
- Source row/title: Separate chunking descriptor
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-gobaz-chunking-identity-bakeoff/`.

## Setup

Alice stores a chunked object whose root points at a chunking descriptor, provisionally called a chunking CID or cCID candidate.

## Stimulus

Run the candidate simulation against this source test: Whether peers can verify chunks and compare roots without overloading the object pCID.

## Expected Pressure

A separate descriptor may isolate chunker evolution but adds another identity object.
