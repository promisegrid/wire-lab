# TODO-godad: POC13 CAS storage and CID-named compute functions

## Decision Intent Log

ID: DI-bibom
Date: 2026-06-06 01:12:50
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Define POC13 as `implementations/poc13-cas-compute-functions/`, a
future executable proof of concept for decentralized CAS storage promises and
CID-named function-call compute promises. POC13 should use two provisional
protocol pCID names, `cas_storage_v1` and `cid_compute_v1`, under the same
`grid([42(pCID), ...protocol-defined-slots])` envelope discipline. A pCID names
the protocol spec; message types and function CIDs live inside pCID-owned
payloads.
Intent: POC12 proves a production-ish app/kernel/device workflow, but it does
not yet prove PromiseGrid's core storage/compute shape. POC13 should focus on a
decentralized sparse CAS where content identity is separate from availability,
retention, replication, access, and serving promises, plus compute where code is
stored in CAS and called by CID. Pure results can be cached by exact function,
input, protocol, and context identity; impure work must externalize timestamp,
randomness, sensor reads, or other ambient inputs as explicit context objects so
the run is replayable and pure-after-the-fact.
Constraints: This TODO defines POC13 but does not implement the executable POC
yet. Do not add separate pCIDs for message types, global registries, central
storage authority, RPC verbs, permission/conformance framing, hidden ambient
inputs, or provider-backed runs. Keep storage, compute, capability, retention,
and cache behavior expressed as promises and local evidence.
Affects: `protocols/wire-lab.d/TODO/TODO-godad-poc13-cas-compute-functions.md`;
`protocols/wire-lab.d/TODO/TODO.md`;
`docs/research/DN-nuras-poc13-cas-compute-functions.md`;
`implementations/README.md`;
`DEV-GUIDE-RESOURCES.md`; storage/compute scenarios; future
`implementations/poc13-cas-compute-functions/**`.

## Prior aliases

- None.

## Status

Planned. This TODO locks the POC13 shape and scenario pressure. The executable
POC13 implementation remains a later task.

## Tasks

- [x] godad.1 Record the POC13 definition DI and cross-list this TODO.
- [x] godad.2 Write `docs/research/DN-nuras-poc13-cas-compute-functions.md`.
- [x] godad.3 Document decentralized CAS storage promises: content identity,
  availability, retention, replication, access, serving, corrupt-byte evidence,
  and partial availability.
- [x] godad.4 Document CID-named compute promises: `function_cid`, inputs,
  declared context objects, result cache keys, pure-after-the-fact impure calls,
  and broken/malformed result evidence.
- [x] godad.5 Expand or add scenarios for CAS storage and CID-named function
  compute pressure before implementing the executable POC.
- [x] godad.6 Update `DEV-GUIDE-RESOURCES.md` and `implementations/README.md`
  so POC13 is visible as planned evidence, not a stable API.
- [ ] godad.7 Implement `implementations/poc13-cas-compute-functions/` in a
  later batch after this definition is reviewed.

## Acceptance criteria

- POC13's planned storage and compute protocols keep one top-level semantic act:
  `promise`.
- `cas_storage_v1` and `cid_compute_v1` are protocol pCID names, not message
  type names.
- Function code identity is payload-level CAS identity, not envelope pCID
  identity.
- Content identity never implies availability, retention, access, or serving
  authority.
- Impure compute inputs are explicit context promises or context objects so
  result evidence can be replayed and audited locally.
