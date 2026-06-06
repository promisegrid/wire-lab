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

ID: DI-notig
Date: 2026-06-06 02:04:59
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Implement POC13 as a self-contained Go proof of concept under
`implementations/poc13-cas-compute-functions/`, with commands
`poc13-supervisor`, `poc13-agent`, and `poc13-analyze`. The first executable
version uses LLM-backed local agents Alice, Bob, Carol, Dave, Ellen, Frank,
Grace, and Mallory across Docker containers, with deterministic protocol
validation in Go before any LLM decision is trusted.
Intent: POC13 should pressure-test CAS storage and CID-named compute as
PromiseGrid promises rather than RPC calls. The POC needs enough live autonomy
to expose local decision and malformed-input pressure, while keeping wire shape,
pCID handling, signatures, CID verification, cache keys, and analyzer gates
deterministic enough to diagnose.
Constraints: Keep POC13 self-contained and do not import POC12 code or extract
a shared library yet. Use only provisional protocol pCID names
`cas_storage_v1` and `cid_compute_v1`; message variants, content CIDs, function
CIDs, and context CIDs live inside pCID-owned payloads. Keep one top-level
semantic act, `promise`, and do not add global registries, central storage or
compute authority, hidden RPC verbs, permission/conformance framing, or stored
API keys.
Affects: `implementations/poc13-cas-compute-functions/**`;
`implementations/README.md`;
`DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-godad-poc13-cas-compute-functions.md`.

## Prior aliases

- None.

## Status

Implemented as first executable evidence. POC13 remains provisional and does
not freeze final storage, compute, cache, provider, kernel, or app APIs.

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
  so POC13 is visible as storage/compute evidence, not a stable API.
- [x] godad.7 Implement `implementations/poc13-cas-compute-functions/` as a
  self-contained first executable POC with analyzer gates.

## Acceptance criteria

- POC13's storage and compute protocols keep one top-level semantic act:
  `promise`.
- `cas_storage_v1` and `cid_compute_v1` are protocol pCID names, not message
  type names.
- Function code identity is payload-level CAS identity, not envelope pCID
  identity.
- Content identity never implies availability, retention, access, or serving
  authority.
- Impure compute inputs are explicit context promises or context objects so
  result evidence can be replayed and audited locally.
