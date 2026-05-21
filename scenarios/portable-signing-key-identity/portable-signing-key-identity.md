# Portable Signing-Key Identity

## Scenario ID

portable-signing-key-identity

## Source / Provenance

- Source type: new harness scenario
- Source path: `/home/stevegt/lab/promisegrid-dev-guide/FEEDBACK.md`
- Source row/title: `FB-robif - baseline signing-key identity across browser and plugin host`
- Source DI / TODO / TE: `DI-ragaz`; `TODO-rozas`; `DR-tuhaz`

## Purpose

Test a provisional identity recipe for one user appearing consistently across a
browser tab and a long-running plugin host.

## Setup

Alice has a browser profile and a plugin helper. The browser can use WebCrypto
and IndexedDB. The helper can use Node crypto and a filesystem. Bob observes only
signed protocol artifacts. Carol audits a later key rotation. Mallory injects a
display-name collision and attempts to confuse local usernames with durable
identity.

## Stimulus

Alice rotates from an old signing key to a new signing key and opens a new live
session from the other host embodiment.

## Expected Pressure

The candidate design must distinguish signing-key continuity from presentation
hints, define a pivotable v0 recipe without hardcoding forever cryptography, and
show how constrained-host storage changes the claim.

## Scenario-Specific Evaluation Questions

- Which parts of the key algorithm, rotation, and handshake can be provisional?
- What evidence links old and new keys?
- What should the guide say about browser key storage and XSS risk without
  overstating security?
