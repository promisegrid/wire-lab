# TE-pudiv: App/kernel grid message boundary

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-pudiv

## Status

needs DF

## Decision under test

Should the local app/kernel boundary use the same PromiseGrid message shape as
the wire boundary: `grid([42(pCID), payload, ...])`?

`TE-mazop` identified the missing app/kernel/host promise surface and minimum
credible first-port contract. It allowed pCID-selected envelope handling, but it
did not make one alternative explicit: apps, kernels, hosts, and local services
may all speak ordinary PromiseGrid messages to each other, even when the message
never leaves the device.

This TE tests that alternative directly. It does not close `DR-davod`.

## Assumptions and constraints

- The current outer envelope remains `grid([42(pCID), payload, ...])`.
- `pCID` means Protocol CID. The pCID names the spec that defines payload shape,
  later outer slots, proof semantics, and handler behavior.
- Local messages are still promises between autonomous agents or local roles.
  A kernel does not command an app, and an app does not command a kernel.
- A message crossing a local process boundary, host boundary, or service boundary
  can be evidence even if it is never sent over a network.
- Host/runtime facts remain local assumptions unless the host adapter is itself
  an explicit promiser.
- Unknown pCID handling remains local: preserve exact bytes where promised or
  refuse with evidence; do not semantically accept unsupported messages.

## Alternatives

### Alt 1 - Separate local API

Apps call a local API such as function calls, syscalls, host imports, or RPC
methods. The kernel/runtime translates those calls into PromiseGrid messages
only when it needs to speak on the wire.

**Easier:** Natural for conventional SDKs and host runtimes. Strong tooling
support in browsers, native daemons, WASM hosts, and mobile OS APIs.

**Harder:** Creates two semantic surfaces: local API semantics and wire-message
semantics. Future readers must reconstruct how a local operation mapped to a
pCID promise.

**New obligations:** The guide must define a stable local API or explain why
every implementation-specific API is only an adapter over pCID messages.

### Alt 2 - Same grid message boundary

Apps and kernel roles exchange `grid([42(pCID), payload, ...])` messages
locally. Storage requests, compute requests, network sends, key-use requests,
device requests, refusals, receipts, and evidence records are all pCID-selected
messages. A rich daemon, WASM host, MCU library, or split local service may
transport those messages differently, but the semantic object is the same.

**Easier:** One semantic surface for local and remote interactions. Exact bytes
can be logged, replayed, forwarded, or audited. pCID specs define both local
and wire interpretation.

**Harder:** Some runtimes still need ergonomic wrappers. A small device may not
want to allocate full message objects for every internal call. Guide prose must
avoid implying that every local function call is necessarily serialized.

**New obligations:** Define when exact grid bytes must exist, when an in-memory
representation is sufficient, and how local wrappers prove they are faithful to
the pCID-selected message.

### Alt 3 - Hybrid adapter

The promised boundary is pCID-selected grid messaging, but implementations may
offer local APIs as adapters. The adapter promises that each exposed operation
has a corresponding pCID message shape, evidence record, and refusal behavior.

**Easier:** Preserves one semantic model while allowing ergonomic SDKs and
host-specific APIs.

**Harder:** Adapter promise fidelity becomes a real burden. If the adapter
hides the underlying pCID messages, the system drifts back toward opaque local
APIs.

**New obligations:** The port promise record must say which local API calls
correspond to which pCID messages and what exact evidence is recorded.

## Scenario analysis

### S1 - Native daemon with local apps

Alice's app asks Bob's local daemon to store bytes, send a message, and use a
signing key. Under Alt 1, the daemon exposes SDK calls and later emits wire
messages. Under Alt 2, Alice's app sends local grid messages to Bob's daemon.
Under Alt 3, the SDK call is a wrapper over a specific pCID message.

Alt 2 gives the cleanest evidence: Bob can record the exact request, refusal, or
receipt bytes. Alt 3 may be the practical developer surface if the adapter keeps
the pCID mapping explicit. Alt 1 is weakest for long-horizon replay unless the
local API transcript is also normalized into evidence messages.

### S2 - Browser/WASM host

Carol runs a browser app. The browser owns storage, network access, clocks, and
lifecycle. Carol's PromiseGrid code can ask a host adapter to store bytes or
send a message.

Alt 2 says the request to the host adapter is itself a grid message. The browser
still is not a PromiseGrid authority; the adapter can promise only what it will
attempt and record. Alt 3 fits if JavaScript exposes ergonomic calls that are
specified as wrappers over pCID messages. Alt 1 risks making browser API quirks
part of the promise surface.

### S3 - MCU/header-only port

Dave's MCU supports one actuator pCID. It can parse a compact grid message,
verify the pCID-specific proof, and write a tiny evidence record. It may not
serialize every internal function call.

Alt 2 is viable if the message boundary is at ingress, egress, and evidence
records rather than every internal helper call. Alt 3 is strongest here: the
firmware may use local C functions internally, but any exposed PromiseGrid
operation is specified as a pCID-selected grid message.

### S4 - Split local services

Ellen splits storage, key use, network send, evidence recording, and dispatch
into local services. If each service speaks grid messages, evidence handoff is
uniform and each local service can be its own promiser. If the services use a
private RPC surface, the port promise record must document how that RPC maps to
pCID messages.

Alt 2 best preserves Promise Theory boundaries: each local service promises its
own behavior through messages. Alt 3 remains acceptable if the private RPC is
only an adapter and the evidence record preserves the pCID-level promise.

### S5 - Long-horizon replay

Frank reimplements the port decades later. Source languages, OS APIs, and SDKs
have changed. The old pCID specs and exact message bytes remain.

Alt 2 is easiest to replay because local and wire evidence share one format.
Alt 3 is acceptable if the adapter mapping was recorded. Alt 1 is fragile unless
the local API has its own durable spec and evidence format.

## Cross-cutting findings

- The same grid envelope is a serious candidate for the app/kernel semantic
  boundary.
- "Same message format" does not require every internal function call to be
  serialized. It means exposed app/kernel operations have pCID-selected message
  semantics and exact evidence when they cross a promise boundary.
- Local APIs can still exist, but they should be treated as adapters whose
  promise is to faithfully emit, accept, or record pCID-selected messages.
- This framing fits Promise Theory better than a kernel-command API because
  every app, kernel role, host adapter, or local service remains an autonomous
  promiser.

## Conclusions

- `TE-mazop` should not be rewritten; it correctly identified the missing
  evidence but did not center this alternative.
- `DR-davod` should include a same-grid app/kernel boundary question before it
  asks for a final kernel definition.
- `SIM-fovip` should either be extended or followed by a successor sim that
  tests app/kernel requests, refusals, receipts, and evidence records as
  `grid([42(pCID), payload, ...])` messages.
- The strongest surviving shape is Alt 3: **grid-message semantics as the
  promised boundary, with local APIs allowed as adapters**. Alt 2 remains the
  cleanest pure form. Alt 1 should be treated as an implementation convenience,
  not the likely PromiseGrid semantic boundary.

## Recommended next DF packet for DR-davod

Before deciding the stable kernel-developer porting boundary, answer:

1. Are app/kernel operations first-class pCID-selected grid messages?
2. If local APIs exist, must each API operation map to a pCID message and
   evidence record?
3. Which app/kernel operations need first-class pCIDs: storage, compute,
   network send, key use, device access, lifecycle, refusal, receipt, and
   evidence recording?
4. When must exact grid bytes be preserved locally, and when may an in-memory
   representation stand in for them?
5. Should `SIM-fovip` be extended in place for this question, or should a
   successor sim test the same-grid app/kernel boundary separately?

## Implications for open work

- `DR-davod` remains open.
- `TE-mazop` gets a forward pointer to this TE rather than a body rewrite.
- `DEV-GUIDE-RESOURCES.md` should describe the same-grid app/kernel boundary as
  active evidence, not settled guidance.
- A future sim should test storage, compute, send, key, device, lifecycle,
  refusal, receipt, and evidence operations as local grid messages.

## Decision status

`needs DF` - this TE makes the same-grid app/kernel boundary first-class but
does not decide the final kernel-developer porting boundary.
