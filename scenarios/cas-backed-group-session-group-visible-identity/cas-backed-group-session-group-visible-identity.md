# Group-visible identity

## Scenario ID

cas-backed-group-session-group-visible-identity

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-jurar-cas-backed-group-session/SCENARIOS.md`
- Source simulation: `SIM-jurar-cas-backed-group-session/`
- Source row/title: Group-visible identity
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-jurar-cas-backed-group-session/`.

## Setup

Alice posts a group message whose visible identifier could be a pointer-object CID, message-root CID, or envelope CID.

## Stimulus

Run the candidate simulation against this source test: Which CID humans, tools, parent links, and acknowledgements should cite.

## Expected Pressure

The migration must avoid two competing identities for one logical message.
