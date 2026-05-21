# SIM-tizad-scoped-conformance-citation-ledger: Scoped conformance and citation records

This promoted simulation tests a bounded guide pattern for early PromiseGrid apps that need honest partial conformance, heterogeneous embodiments, and optional off-grid live behavior without overclaiming. The pattern adds two small artifacts: a publishable layer-scoped conformance manifest and peer-local promise-accounting records.

## Question

How should an app say what is local, what is live-only, and what is durably promised so Bob can interoperate honestly and Carol can audit failures years later?

## Design Move

- **Layer-scoped conformance manifest:** publish separate claims for `local-implementation`, `live-transport`, `audit-publication`, and optional `device-effect` scope, all anchored to one authoritative protocol-boundary contract or contract family.
- **Durable citation rule:** any durable audit claim must cite a concrete CAS milestone object such as a blob CID, snapshot CID, save blob, receipt, or break-witness; it must not treat a live session as the durable object.
- **Peer-local promise-accounting records:** each peer keeps minimal records for identity, availability, authorization, and non-idempotent effect claims, including failures.

## Decision Axes

- **Authoritative boundary identity:** which contract reference is the durable app identity across embodiments.
- **Claim scope separation:** what is safe to claim at local, live, audit, and effect layers.
- **Promise class separation:** content identity is distinct from availability, authorization, and physical-effect execution.
- **Durable object choice:** what exact object an audit message cites for blobs, CRDT milestones, key rotation, or device receipts.
- **Failure evidence:** what Alice, Bob, and Carol must retain when retrieval, replay handling, or interop fails.

## Related Root Scenarios

- `scenarios/app-semantics-partial-conformance/app-semantics-partial-conformance.md`
- `scenarios/live-crdt-audit-publication/live-crdt-audit-publication.md`
- `scenarios/minimal-immutable-blob-app/minimal-immutable-blob-app.md`
- `scenarios/multi-embodiment-app-identity/multi-embodiment-app-identity.md`
- `scenarios/device-bound-agent-physical-effect/device-bound-agent-physical-effect.md`
- `scenarios/portable-signing-key-identity/portable-signing-key-identity.md`

## Why this should score better

This design keeps the parents' strong boundary discipline but repairs three recurring gaps: lack of a concrete partial-conformance claim shape, lack of a precise durable object to cite, and lack of explicit local records for failures and replay disputes. It should improve auditability, failure handling, and implementation plausibility without pretending that live transport, signatures, or device delegation are already frozen.

## Boundaries

This simulation does not define a final live binding, final envelope/signature format, universal capability token, central registry, or final device delegation standard. It is guide-safe orientation for claim scoping and evidence capture while upstream decisions remain open.

This specimen was promoted from review-stage child proposal
`SIM-tizad-child-scoped-conformance-citation-ledger` from
`ga-canary-20260521-011601` under `DI-fihub`. The ignored proposal artifacts
remain local raw evidence; this directory is the canonical non-child simulation
home.
