# Minimal Immutable Blob App

## Scenario ID

minimal-immutable-blob-app

## Source / Provenance

- Source type: new harness scenario
- Source path: `/home/stevegt/lab/promisegrid-dev-guide/FEEDBACK.md`
- Source row/title: `FB-vopik - What CAS-facing guarantees are safe for a minimal immutable blob app?`
- Source DI / TODO / TE: `DI-ragaz`; `TODO-rozas`; `DR-tuhaz`; `DR-tumus`

## Purpose

Exercise candidate designs against the smallest useful CAS-facing app: Alice
writes immutable bytes and receives a hash; later Carol presents the hash and
expects to retrieve the same bytes.

## Setup

Alice uploads a blob through Bob's app. Bob stores or publishes a
content-addressed object and returns a hash. Carol later receives only the hash
and partial context. Mallory may withhold storage, replay stale availability
claims, or claim that possession of the hash is enough authorization.

## Stimulus

The original host changes retention policy, a peer cache is incomplete, and Carol
tries to read the blob from a different site years later.

## Expected Pressure

The candidate design must separate content identity from availability,
authorization, ingress, discovery, replication, and retention promises while
still preserving enough evidence for a 100-year audit trail.

## Scenario-Specific Evaluation Questions

- What exactly does `hash in -> blob out` promise, and who makes that promise?
- Is possession of the hash a read capability, an address, or an app-level
  convention?
- What local promise accounting records should Alice, Bob, and Carol keep when
  storage or retrieval fails?
