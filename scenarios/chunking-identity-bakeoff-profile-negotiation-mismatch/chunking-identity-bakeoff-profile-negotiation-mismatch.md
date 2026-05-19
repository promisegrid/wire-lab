# Profile negotiation mismatch

## Scenario ID

chunking-identity-bakeoff-profile-negotiation-mismatch

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-gobaz-chunking-identity-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-gobaz-chunking-identity-bakeoff/`
- Source row/title: Profile negotiation mismatch
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-gobaz-chunking-identity-bakeoff/`.

## Setup

Alice advertises profile `fastcdc-small`; Bob supports only `rabin-large`.

## Stimulus

Run the candidate simulation against this source test: Whether peers can refuse, bridge, or request raw bytes without silently accepting mismatched roots.

## Expected Pressure

Negotiated profiles need failure behavior and may still need persistent identity binding.
