# TE-6: Capability-token revocation propagation

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TE-6` (integer alias)
- `TE-20260427-180500` (timestamp alias and pre-migration filename)

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*


## Status

stub
Alice issues tokens that B, C, D delegate among themselves through three more hops. Alice now revokes. How long until the holder at hop 5 finds out? Can a malicious intermediate suppress the revocation? Outcome: a wire-level revocation pattern that resists suppression.
