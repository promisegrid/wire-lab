# App Semantics Partial Conformance

## Scenario ID

app-semantics-partial-conformance

## Source / Provenance

- Source type: new harness scenario
- Source path: `/home/stevegt/lab/promisegrid-dev-guide/FEEDBACK.md`
- Source row/title: `FB-dodos`, `FB-hisis`, `FB-kutub`, `FB-gomod`, and `FB-tahof`
- Source DI / TODO / TE: `DI-ragaz`; `TODO-rozas`; `DR-tuhaz`

## Purpose

Test whether an app can publish honest provisional semantics and conformance
claims without pretending to implement a final PromiseGrid app contract.

## Setup

Alice ships a useful first app slice. It uses local IDs internally, signs current
messages with an adapter-local carriage choice, and supports append/read behavior
but not the full merge, authority, capability-token, or break-witness story.
Carol needs to know what can interoperate. Mallory benefits if the app overclaims.

## Stimulus

Alice publishes a B-side conformance claim and Bob tries to interoperate using
only that claim, the draft spec path, and locally observed wire artifacts.

## Expected Pressure

The candidate design must distinguish local implementation shortcuts from
protocol-boundary identity and conformance, and must say which semantic claims
are provisional, blocked, or safe orientation.

## Scenario-Specific Evaluation Questions

- What wording makes a partial-conformance claim honest?
- Which identity is authoritative at the protocol boundary?
- How should provisional signature carriage and capability/witness language be
  described without freezing them?
