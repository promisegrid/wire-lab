# Minimal Immutable Blob App

## Scenario ID

minimal-immutable-blob-app

## Source / Provenance

- Source type: new harness scenario
- Source path: `/home/stevegt/lab/promisegrid-dev-guide/FEEDBACK.md`
- Source row/title: `FB-vopik - What CAS-facing guarantees are safe for a minimal immutable blob app?`
- Source DI / TODO / TE: `DI-ragaz`; `TODO-rozas`; `DR-tuhaz`; `DR-tumus`;
  `DI-bibom`; `TODO-godad`

## Purpose

Exercise candidate designs against the smallest useful decentralized CAS-facing
app: Alice writes immutable bytes and receives a content identity; later Carol
presents that identity and asks one or more peers to keep a serving promise for
the same exact bytes.

## Setup

Alice uploads a blob through Bob's app. Bob may promise to store, retain,
replicate, advertise, or serve the resulting content-addressed object, but each
of those is a separate promise. Carol later receives only the CID and partial
context. Dave has a cache that might contain the bytes but has made no retention
promise. Mallory may withhold storage, replay stale availability claims, provide
bytes that do not match the CID, or claim that possession of the CID is enough
evidence of an access promise.

## Stimulus

The original host changes retention policy, a peer cache is incomplete, and Carol
tries to read the blob from a different site years later.

## Expected Pressure

The candidate design must separate content identity from availability, access,
ingress, discovery, replication, retention, and serving promises while still
preserving enough exact-byte evidence for a 100-year audit trail.

## Scenario-Specific Evaluation Questions

- What exactly does `hash in -> blob out` promise, and who makes that promise?
- Is possession of the hash a read capability, an address, or an app-level
  convention?
- What local promise accounting records should Alice, Bob, and Carol keep when
  storage or retrieval fails?
- How does Carol distinguish "this CID names these bytes" from "Bob or Dave
  currently promises to serve these bytes"?
- What evidence should be recorded when Mallory serves bytes that fail CID
  verification, or when Bob once stored the object but no longer promises
  retention?
