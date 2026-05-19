# Pointer object identity

## Scenario ID

cas-object-model-pointer-object-identity

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-jomag-cas-object-model/SCENARIOS.md`
- Source simulation: `SIM-jomag-cas-object-model/`
- Source row/title: Pointer object identity
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-jomag-cas-object-model/`.

## Setup

Alice creates pointer object CID Y that points at root CID X; Bob has Y but not X.

## Stimulus

Run the candidate simulation against this source test: Whether pointer objects are verifiable CAS objects in their own right and how sparse-CAS resolution behaves when the pointee is absent.

## Expected Pressure

TODO-pipus migration must preserve the distinction between pointer-object identity and pointed-at root identity.
