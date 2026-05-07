# TE-9: Two communities, two pCIDs, same intent

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TE-9` (integer alias)
- `TE-20260427-180800` (timestamp alias and pre-migration filename)

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*


## Status

stub
Community A's spec for compute-call hashes to pCID `0xAAAA…`. Community B independently writes a spec for the same purpose; theirs hashes to pCID `0xBBBB…`. The two pCIDs are different bytes (different specs), even though they aim at the same use case. The two communities meet. How does bridging happen — a gateway agent that holds both specs and translates, dual-emitting handlers that publish both pCIDs, runtime spec-comparison? Outcome: a documented interop pattern. (This is exactly the no-central-registry tension: anyone can mint a pCID, so collisions of intent — same purpose, different specs — are inevitable in a long-lived open system.)
