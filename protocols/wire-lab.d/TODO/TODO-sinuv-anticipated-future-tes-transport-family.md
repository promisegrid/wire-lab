# TODO-sinuv: Anticipated future TEs: ring / cluster-of-clusters / gossip / receipts-at-scale

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TODO-31` (integer alias)
- `TODO-20260507-002306` (timestamp alias and pre-migration filename)

## Status

Parked. Each thread is an anticipated future TE that has been
forward-pointed from existing TEs but is not yet scheduled. Bundled
into one TODO since none are imminent.

## Threads absorbed from OPEN-THREADS.md

### T-RING-TRANSPORT (formerly OPEN-THREADS, opened 2026-05-01)

Ring-transport spec (anticipated future TE). Originally TE-liviv+
forward-pointer in TE-junil.

Scope: token-ring semantics; per-hop authorization; ordering
guarantees stronger than gossip but weaker than centralized.

Anchor: TE-junil Refinements 2026-05-05 (item 2).

### T-CLUSTER-OF-CLUSTERS-TRANSPORT (formerly OPEN-THREADS, opened 2026-05-01)

Cluster-of-clusters transport (anticipated future TE).

Scope: hierarchical transport; cluster as first-class addressable
unit; how messages traverse the inter-cluster boundary.

Anchor: TE-junil Refinements 2026-05-05 (item 3).

### T-GOSSIP-TRANSPORT (formerly OPEN-THREADS, opened 2026-05-01)

Gossip-transport spec (anticipated future TE).

Scope: epidemic propagation; convergence guarantees; how IHave/IWant
generalize.

Anchor: TE-junil Refinements 2026-05-05 (item 4).

### T-RECEIPTS-AT-SCALE (formerly OPEN-THREADS, opened 2026-05-01)

Receipts at scale (anticipated future TE).

Scope: does `IHave: <transport-pcid>:<cid>` need to become a vector
at multi-writer or large-N transports?

Anchor: TE-junil Refinements 2026-05-05 (item 5).

## Question log

(Per AGENTS-ppx Question-logging discipline. No questions logged yet.)

## Decision Intent Log

(Will be populated as DFs lock and product lands.)
