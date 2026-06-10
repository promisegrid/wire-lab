# TODO-juput: POC13 Run-Scoped Durability Promises

## Decision Intent Log

ID: DI-sunuf
Date: 2026-06-08 18:44:56
Status: active
Decision: Extend POC13 with run-scoped durable CAS/evidence state, promise-shaped retention/GC, backpressure, rate-limit, replay-protection, restart/recovery, analyzer, and fuzz/chaos evidence.
Intent: POC13 should test production-relevant durability and operational pressure without turning the POC into cross-run persistent infrastructure. CAS bytes, compute cache, capability tokens, replay windows, and evidence journals may survive a process restart inside one run, while clean-run reset remains the boundary between experiments. Retention, deletion, backpressure, rate limits, and replay handling are local promises and observations, not central policy or command enforcement.
Constraints: Keep the single top-level `promise` action; do not add workflow-specific top-level action kinds; keep trust local and per-agent; do not make GC, rate limits, or replay protection global authorities; keep state under the current POC13 run root; clean-run scripts may remove all runtime state between experiments; preserve the POC11/POC12 superset behavior locked by `DI-sinur`.
Affects: implementations/poc13-cas-compute-functions/; protocols/wire-lab.d/TODO/TODO.md; DEV-GUIDE-RESOURCES.md.

ID: DI-vahan
Date: 2026-06-09 06:56:44
Status: active
Decision: In POC13, Alice's second arbitrary compute request should follow local trust evidence and use Dave after Carol exposes malformed bad-result evidence.
Intent: Negative compute evidence should have visible routing consequences without deleting the alternate-function coverage gate. Alice may still test arbitrary pCID-defined compute payloads, but she should not keep sending fresh compute work to Carol immediately after locally reducing trust in Carol.
Constraints: Keep all behavior promise-first; do not introduce RPC-style dispatch or global trust; keep the single top-level `promise` action; preserve `compute_alternate_function_executed` coverage through a peer that currently has enough local trust for Alice to contact.
Affects: implementations/poc13-cas-compute-functions/runtime/node.go; implementations/poc13-cas-compute-functions/README.md; DEV-GUIDE-RESOURCES.md.

## Tasks

- [x] juput.1 Add run-scoped durable CAS/evidence stores with a clean-run reset invariant.
- [x] juput.2 Model retention as explicit promises: retain-until, delete-after, token-expiry, disk-pressure, and superseded-checkpoint conditions.
- [x] juput.3 Record GC as local evidence for promise-ended, object-removed, object-retained, and promise-broken cases.
- [x] juput.4 Add restart/recovery tests inside one POC13 run using run-scoped stores.
- [x] juput.5 Model backpressure as capacity promises between apps, kernels, and peers.
- [x] juput.6 Model rate limits as reciprocal self-promises: sender send-rate promises and receiver accept-rate promises.
- [x] juput.7 Add analyzer gates for retention, GC, backpressure, and rate-limit promise evidence.
- [x] juput.8 Add replay protection for exact envelopes and capability tokens as promise/evidence semantics.
- [x] juput.9 Add fuzz/chaos tests for CBOR parsing, pCID routing, crashes, delayed ACKs, and partial writes.
