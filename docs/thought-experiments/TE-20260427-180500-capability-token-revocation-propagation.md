# TE-6: Capability-token revocation propagation

*Thought experiment, part of the [PromiseGrid Wire Lab](../../harness-spec.md). This file is content-addressable; its hash is its pCID.*

Alice issues tokens that B, C, D delegate among themselves through three more hops. Alice now revokes. How long until the holder at hop 5 finds out? Can a malicious intermediate suppress the revocation? Outcome: a wire-level revocation pattern that resists suppression.
