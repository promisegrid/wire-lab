# Same bytes, different chunkers

## Scenario ID

chunking-identity-bakeoff-same-bytes-different-chunkers

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-gobaz-chunking-identity-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-gobaz-chunking-identity-bakeoff/`
- Source row/title: Same bytes, different chunkers
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-gobaz-chunking-identity-bakeoff/`.

## Setup

Alice chunks a file with a small FastCDC-style target; Bob chunks the same file with promisebase / pitbase-style Rabin defaults.

## Stimulus

Run the candidate simulation against this source test: Whether leaf CIDs and Merkle roots diverge, and where the differing chunking parameters are visible.

## Expected Pressure

TE-43 must either lock chunking parameters or bind them into object identity explicitly.
