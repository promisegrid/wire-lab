# Transport Spec (DRAFT)

*This is the wire-lab's outer spec for the `transports/` directory. It is a draft and is subject to revision; once frozen, its pCID will name this protocol class for all time. See `specs/MANIFEST.md` for freeze status.*

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted. Cross-references to it in other repo files use `specs/transport-spec-draft.md` (path) until freeze; after freeze they will use the pCID.

## Purpose

This rooted harness draft now tracks the **apparatus/governance residue** that
remains after `rusis.11` extracted the specimen-side outer feed convention into
`simulations/SIM-labit-feed-outer/protocols/feed-outer.d/specs/feed-outer-draft.md`.

The extracted specimen-side sections are:

- `## Purpose`
- `## The four locked principles (TE-zalut)`
- `## What this spec does NOT specify`

This rooted file retains the source/governance context and freeze/open-question
memory that still belongs to harness-side ownership. Source: `DI-huzor`.

## Sources

This spec is locked by the conclusions of:

- [TE-hogus](../docs/thought-experiments/TE-hogus-group-transport-envelope.md): `grid <pcid>` as the group-transport envelope; canonical-bytes rules; explicit-promise body requirement. (Source for the group-transport-protocol's contract; not a constraint on this outer spec.)
- [TE-zalut](../docs/thought-experiments/TE-zalut-channel-transport-types-and-threaded-replies.md): transport-protocol types, pCID-keyed transport paths, and DAG message graphs. Establishes the extracted thin-outer principles now tracked in `SIM-labit-feed-outer`.
- [TE-junil](../docs/thought-experiments/TE-junil-transports-rename-and-axes-of-differentiation.md): transports rename and axes of transport-protocol differentiation. Establishes the per-axis meta-rule for distinguishing transport-protocols.
- [DR-009](../DR/DR-009-20260430-204108-group-transport-envelope.md): the active decision request governing the group-transport envelope and its graduation.

## The per-axis meta-rule (TE-junil)

When deciding whether a new transport-protocol warrants a distinct pCID (and therefore a distinct spec doc), the following per-axis rule applies:

| Axis | Distinct pCID per value? | Notes |
|------|--------------------------|-------|
| A. Cardinality (N=2 vs. small-N vs. large-N vs. unbounded) | Parameter, **except at extremes** (large-N, unbounded) | Small-finite-closed-group with N≥2 is one spec; very large or unbounded membership crosses a contract boundary. |
| B. Visibility (all-see-all, hub-mediated, ring-propagated, subset-addressed, topic-filtered, gossip) | **Distinct pCID per class** | Observably different contracts. |
| C. Routing topology (direct, mesh, star, ring, tree, layered, cluster-of-clusters) | **Distinct pCID per class** | Each leaves qualitatively different on-disk artifacts. |
| D. Membership rules (static, invite-only, open-read, open, capability-token) | Parameter, with `capability-token` as a candidate exception | Most are spec parameters; permissioned transport may warrant a distinct pCID. |
| E. Persistence (append-only, bounded retention, compactable, ephemeral) | Parameter, **except `ephemeral`** | Ephemeral crosses a contract boundary because the simulation cannot observe a transport whose messages disappear. |
| F. Message-graph shape (independent, single-writer chain, multi-writer DAG, synchronized frontier, vector-clock) | Parameter | Different parent-link shapes within one transport-protocol. |
| G. Direction (symmetric, hub-asymmetric, multicast, paired) | Parameter | Different direction values within one transport-protocol. |
| H. Reliability / receipts (none, per-message, frontier, cryptographic) | Parameter | Different receipt schemes within one transport-protocol. |

This rule is not exhaustive; it is the working policy that survives until experience teaches a better one.

## Open questions

- **OQ-1 (deferred):** Should the wire-lab spec define a small companion convention for transport-protocols to publish their own pCID on first use, so receivers can discover the protocol-spec from a stranger's first message? Raised and deferred in TE-zalut DF-26.8 Alt-8.B; the locked Alt-8.C (code-as-handler) does not address this. May surface in a future TE.

- **OQ-2 (partially resolved 2026-05-05):** What does it mean for a group of participants to migrate from one transport-protocol-pCID to another transport-protocol-pCID? (S7 of TE-junil.) The migration *invariants* are locked by [TE-numan: Transport-protocol migration invariants](../../../docs/thought-experiments/TE-numan-transport-protocol-migration-semantics.md): audit-trail reconstructibility, no-silent-rewrite, and no-unilateral-abandonment. Any concrete migration contract must satisfy all three. The *operational shape* of migration (close-old-vs-overlap-vs-atomic-swap; back-reference format; message disposition; seal mechanics; group-identity continuity; trigger discipline; authorizing promise) is deferred to a future operational-shape TE tracked as `T-MIG-OPS` in `OPEN-THREADS.md`; that TE is gated on a concrete first migration to design against.

## Freeze gate

This spec graduates to frozen status when:

1. The repo has at least one transport-protocol spec frozen (currently anticipated to be `protocols/group-session.d/specs/group-session-draft.md`).
2. Steve signs a `merge-transport-spec` promise authorizing the freeze.
3. `tools/spec freeze transport-spec` mints the pCID, snapshots the file, and appends the manifest entry.

Until then, the spec lives at `specs/transport-spec-draft.md` and is a working draft.
