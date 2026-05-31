# implementations/

B-side trees per [TE-liviv](../docs/thought-experiments/TE-liviv-spec-vs-implementation-split.md).

This top-level directory holds reference implementations of PromiseGrid
protocols. Each implementation lives in its own subdirectory:

```
implementations/<impl-name>/
```

The `<impl-name>` slug is human-readable and chosen by the implementer.
Multiple implementations of the same protocol coexist as siblings; one
implementation may implement many protocols.

New proof-of-concept directories use this naming pattern:

```
implementations/pocN-{slug}/
```

Existing `poc2/` through `poc5/` directories keep their historical names.
Future POCs should include a short descriptive slug, such as
`poc6-dag-cbor-interop`. Source: `DI-sagos`.

## CHANGELOG requirement

Every implementation directory MUST contain a `CHANGELOG.md` at its
root recording **conformance claims** against upstream spec doc-CIDs.
Each entry has the shape:

```changelog-entry
claim:           implements | partially-implements | extends | deprecates
spec:            bafkrei...spec-doc-cid
scope:           full | partial-section-N | etc.
breaking-change: true | false
notes:           prose
```

See [TE-zukug](../docs/thought-experiments/TE-zukug-spec-doc-inversion-and-conformance-changelog.md) for the inversion thesis (conformance reference goes implementation -> spec, never spec -> implementation) and [TE-liviv](../docs/thought-experiments/TE-liviv-spec-vs-implementation-split.md) for the spec-side vs implementation-side split.

## External implementations

This directory is just *our local* collection. Third parties' implementations live in their own external repos with the same `CHANGELOG.md` convention; the shape is location-independent. A registry of known external implementations is out of scope for v0 (see TE-liviv OQ-32.1).

## Current implementations

- `poc2/` — proof-of-concept 2: a minimal two-container app/kernel and kernel/kernel hello flow using pCID-selected `grid([42(pCID), payload, ...])` messages. It is not a final PromiseGrid API. Source: `DI-ratij`; `DI-tijat`.
- `poc3/` — proof-of-concept 3: two containers where each container runs one kernel plus hello, echo, and signed apps. It is executable evidence for same-grid app/kernel messaging, not a final SDK. Source: `DI-horak`.
- `poc4/` — proof-of-concept 4: five containers where relay apps, not kernels, own multi-hop neighbor promises for hello, echo, signed, fibonacci, and storage flows. It is executable evidence, not a final relay API. Source: `DI-ponor`.
- `poc5/` — proof-of-concept 5: local trust and selective sending after Alice observes a broken storage promise. It is executable evidence, not a final trust API. Source: `DI-rarim`; `DI-fofik`.
- `poc6-dag-cbor-interop/` — proof-of-concept 6: standalone Go tests for DAG-CBOR/IPLD interop with CID links, byte strings, tag-42 link encoding, stable CIDs, and local evidence. It is executable evidence for `scenarios/cas-object-model-dag-cbor-interop/`, not a final L6 CAS API. Source: `DI-sagos`.
- `poc7-capability-token-exchange/` — proof-of-concept 7: five-container bounded simulation for signed CBOR `grid(...)`-tagged bearer and non-transferable promise tokens over length-framed TCP, issuer-local redemption/revocation, real storage/compute/data redemption payloads, promise-shaped reciprocal exchange offers, Mallory-to-Dave stale-token circulation, deterministic per-agent local economic decisions, and local trust updates. It is executable evidence for `scenarios/promise-economy-capability-token-exchange/`, not a final token, transport, storage, compute, or economics API. Source: `DI-tugih`; `DI-fibok`; `DI-tanat`; `DI-pabot`; `DI-rodog`; `DI-hanih`.
- `poc8-autonomous-promise-economy/` — proof-of-concept 8: five-container successor to POC7 that keeps one pCID-selected signed CBOR `grid([42(pCID), payload, proof])` promise-economy protocol over length-framed TCP while adding autonomous local need advertisements, offer promises, counter promises, acceptance/refusal decisions, token issue/redemption promises, outcome observations, collateral/stake promises, bearer-for-non-transferable exchange, peer-local exchange-rate quotes, and stale-token trust decay. It is executable evidence, not a final token, transport, trust, storage, compute, or economics API. Source: `DI-sirus`.

TODO-jodon (UDP-binding v0 reference) and TODO-bihon (ns-3 harness fixture) remain future implementation work.
