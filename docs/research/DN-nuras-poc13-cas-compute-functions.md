# POC13 CAS storage and CID-named compute functions

This design note defines the intended POC13 shape before implementation. It is
not a frozen PromiseGrid storage, CAS, compute, cache, token, or kernel API.
Source: `DI-bibom`.

## Purpose

POC13 should test whether PromiseGrid can express storage and computation as
local promises over pCID-selected messages without drifting into RPC, global
storage authority, global compute authority, or hidden runtime state.

The eventual implementation path is:

```text
implementations/poc13-cas-compute-functions/
```

## Storage model

Storage is decentralized sparse CAS. A CID identifies bytes or a structured
object, but the CID does not by itself promise:

- who currently has the bytes;
- who will retain the bytes;
- who will replicate the bytes;
- who will serve the bytes;
- who may receive the bytes;
- whether a given peer locally trusts the bytes or the promiser.

Those are promises by specific agents. Alice may store bytes locally, Bob may
promise retention for a time window, Carol may promise replication, and Dave may
promise serving if local trust and capability terms are satisfied. Mallory may
send bytes that fail CID verification; that becomes local evidence about
Mallory's promise, not a failure of the CID concept.

The provisional storage protocol pCID name is `cas_storage_v1`.

## Compute model

Compute is a promise to execute CAS-named code under stated terms. The envelope
pCID names the compute protocol; the function code identity is a payload-level
`function_cid` that refers to CAS-stored code or a code object.

The provisional compute protocol pCID name is `cid_compute_v1`.

Pure function results may be cached only when the full identity matches:

- compute protocol pCID;
- `function_cid`;
- exact input CIDs or scalar values;
- declared context object CIDs;
- ABI/version or execution-profile identity;
- exact result bytes and evidence.

Impure functions should be made pure-after-the-fact where possible. Timestamp,
randomness, sensor reads, filesystem observations, network observations, and
peer observations should be explicit context objects. A result can then promise:
"this result was produced from this function, these inputs, and these context
objects," rather than hiding ambient state.

## App/kernel pressure

POC13 should preserve the current kernel role/profile model. A local storage
role, compute role, cache role, key/signing role, transport role, or evidence
role may be one process, several processes, objects in one runtime, WASM host
functions, or MCU/library code. Each role promises only its own behavior.
Source: `DI-gopag`.

## Non-goals

POC13 must not freeze:

- final CAS storage API;
- final compute API;
- global function registry;
- global content availability registry;
- global cache authority;
- permission/conformance semantics;
- hidden RPC verbs;
- provider-backed scoring or GA behavior.

## Implementation requirements for the later POC

- Keep one top-level semantic act: `promise`.
- Keep message variants inside pCID-owned payloads.
- Keep pCID as protocol identity, not function identity.
- Keep content/function CIDs as payload-level values.
- Record local evidence for kept, broken, malformed, unavailable, duplicate, and
  receiver-non-commitment outcomes.
- Treat cache hits as local promise evidence, not proof of global truth.
