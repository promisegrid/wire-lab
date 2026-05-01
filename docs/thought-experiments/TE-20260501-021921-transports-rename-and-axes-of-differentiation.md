# TE-27: `channels/` → `transports/` rename and axes of transport-protocol differentiation

*Thought experiment, part of the [PromiseGrid Wire Lab](../../specs/harness-spec-draft.md) and the (forthcoming) `specs/transport-spec-draft.md`. This file is content-addressable; its hash is its pCID.*

## TE ID

TE-20260501-021921

## Decision under test

TE-26 locked four principles for what was then called the `channels/` directory: pCID-keyed paths under `channels/<pcid>--<slug>/`, no `Channel:` header in messages, code-as-handler for each pCID, and a thin wire-lab spec deferring directory interiors to per-protocol specs. Drafting moved straight into the carve-out, with `specs/channel-spec-draft.md` (thin) and `specs/codex-perplexity-channel-draft.md` (full TE-24 contract) named as the next deliverables.

Two related observations have surfaced since:

1. **The protocol that two endpoints use is not specific to those two endpoints.** Naming a spec `codex-perplexity-channel-draft.md` anchors a transport-shape ("two named endpoints exchanging messages") to its first instance. The protocol class is generic; the participant pair is an instantiation. The name should describe the class.

2. **What `channels/` actually contains is simulated transports, not channels.** Each subdirectory will need to leave on-disk artifacts that record physical message flow — *how* messages traverse the simulated network — because that flow is what the wire lab is studying. A "channel" in the conventional networking sense is a logical addressing concept layered above a transport; what we have under `channels/` is the transport itself.

These two observations together ask three questions: what should the directory be called, does the word "channel" survive in the wire-lab vocabulary at all, and what differentiates one transport-protocol from another sufficiently to justify a distinct pCID?

## Assumptions

- TE-26's four principles carry forward unchanged in substance; only the vocabulary may shift. pCID-keyed paths, no transport-naming header in messages, code-as-handler, and a thin outer spec all remain locked.
- The wire lab's purpose is to learn what wire-format and topology choices produce desirable behavior in a simulated network. Therefore the directory under study is the simulation surface; its contents are observable artifacts of message flow, not abstractions hiding the network.
- "Let each protocol name its own internals" applies. A transport-protocol-pCID's spec defines its own directory layout, headers, parents, receipts, kinds, persistence rules, etc. The wire-lab spec stays thin.
- The TE-24 v0 contract (parents, receipts, message-id, message kinds, IHave, canonical-bytes) is a candidate transport-protocol's contract, not a wire-lab-global contract. It graduates to the carve-out as one transport-protocol among potentially several.
- Cardinality (N=2 vs. N>2) is a candidate axis of differentiation, but whether different cardinalities warrant distinct transport-pCIDs is itself a question this TE addresses, not one to assume.

## Vocabulary shift

This TE renames `channels/` to `transports/` throughout the repo. The rationale:

- A directory at `transports/<pcid>--<slug>/` represents one simulated transport instance: a particular way bytes-shaped-as-messages traverse a particular set of participants under a particular set of routing rules. The pCID identifies the *class* of transport-protocol; the slug identifies *this instance*.
- "Channel" is reserved (not introduced) as a possible future above-transport concept (logical addressing, multiplexing, group naming) that some later TE may define if a use case forces it. For now there is nothing at the channel layer; the wire-lab vocabulary is **transports** and **messages**.
- TE-26's prose, locked decisions, and follow-on TODO references all migrate from "channel" to "transport." TE-26's substance (no transport-declaration header, pCID-keyed paths, code-as-handler, thin outer spec) carries forward unchanged.

## Axes of transport-protocol differentiation

A transport-protocol's contract is the set of promises its directory and message format imply about how the simulated network behaves. Different values along the axes below produce contracts that are observably different, not merely cosmetically different. The question for each axis is: *does a different value along this axis warrant a distinct transport-protocol-pCID, or is it a parameter within one transport-protocol's spec?*

### Axis A — Cardinality

Number of participants the transport can hold.

- A.1: Exactly 2 (pair).
- A.2: Small finite N (closed group, 3..~50).
- A.3: Large finite N (closed group, ~50..~10⁴).
- A.4: Unbounded (open membership, possibly billions).

A.1 and A.2 share an "every member sees every message" contract. A.3 forces fan-out optimizations (a member may not see every message at near-source latency). A.4 changes the contract qualitatively (no member can know the full membership; visibility is statistical).

### Axis B — Visibility

Who can see which messages.

- B.1: All-see-all (every member sees every message).
- B.2: Hub-mediated (only the hub sees everything; spokes see hub-relayed).
- B.3: Ring-propagated (each node sees messages as they pass; latency proportional to distance).
- B.4: Subset-addressed (sender selects a subset of members per message; non-addressed members do not see).
- B.5: Topic-filtered (members subscribe to topics; see only messages on those topics).
- B.6: Eventual-consistent gossip (every member eventually sees every message, but no per-message guarantee of timely delivery).

B.1 and B.6 differ in *when*, not *whether*, every member sees a message. B.2 changes the contract: spokes have no direct visibility into other spokes. B.3 makes per-message latency a function of topology. B.4 introduces addressing as part of the message envelope. B.5 introduces topics as part of either the message or the subscription.

### Axis C — Routing topology

The graph that messages physically traverse.

- C.1: Direct (every sender writes to a single shared store; every reader reads from it). Simulates a perfect broadcast medium.
- C.2: Mesh (each member can deliver directly to every other member; routing is implicit).
- C.3: Star / hub-and-spoke (all messages route through one designated hub).
- C.4: Ring (messages circulate from one neighbor to the next).
- C.5: Tree (hierarchical parent-child relationships; messages route up and back down).
- C.6: Layered / supernode (a small set of supernodes peer with each other; ordinary nodes attach to a supernode).
- C.7: Cluster-of-clusters (small close-knit groups, each a broadcast domain, with bridges between groups).

C.1 is the simplest simulation: one directory per transport, every member writes-and-reads. C.2-C.7 require the simulation to leave artifacts of *which path* a message traveled, not just *that* it arrived.

### Axis D — Membership rules

How members join, leave, and prove they belong.

- D.1: Static / closed (membership fixed at transport creation; no joins/leaves).
- D.2: Invite-only / closed (existing member must invite; departure recorded).
- D.3: Open-read / restricted-write (anyone can read; only credentialed members write).
- D.4: Open / unauthenticated (anyone can read and write).
- D.5: Capability-token / permissioned (a revocable token grants participation).

### Axis E — Persistence

How long messages live in the transport's directory.

- E.1: Append-only / forever (every message kept indefinitely).
- E.2: Bounded retention (messages older than T are pruned).
- E.3: Compactable (some rule allows replacing N messages with one summary).
- E.4: Ephemeral (messages expire on delivery; no historical record).

E.1 is the easiest to study (every artifact is preserved). E.4 is the hardest to study (the simulation must record the messages externally to observe behavior).

### Axis F — Message-graph shape

The DAG structure messages form within the transport.

- F.1: Independent (no parent links; each message stands alone).
- F.2: Single-writer chain (every message has exactly one parent in its writer's prior log).
- F.3: Multi-writer DAG (each message has zero or more parents drawn from any writer).
- F.4: Synchronized frontier (every new message must name the current frontier as parents; no forks).
- F.5: Vector-clock / causal (parents are a vector of CIDs giving causal context).

F.2 is TE-24's current shape. F.3 is the natural generalization for groups. F.4 is the strongest ordering guarantee. F.5 is what most CRDT-style systems use.

### Axis G — Direction

Who can send.

- G.1: Symmetric (every member can send).
- G.2: Asymmetric / hub (only one designated member sends).
- G.3: Asymmetric / multicast (some members are senders, others receivers).
- G.4: Bidirectional / paired (each member is paired with one other; sends only to that pair).

### Axis H — Reliability / receipts

Whether and how delivery is acknowledged.

- H.1: No receipts (best-effort).
- H.2: Per-message receipts (each receipt acknowledges one CID).
- H.3: Frontier receipts (a receipt acknowledges a contiguous prefix or a frontier vector).
- H.4: Cryptographic receipts (the receipt itself is signed and storable as an independent message).

TE-24 specified H.2/H.4 (`IHave` records signed by the receiver, naming a CID).

## Decision Forks (DFs)

### DF-27.1 — Rename `channels/` to `transports/`?

#### Alt-1.A: Yes, rename now in TE-27's carve-out commits *(LOCKED — chosen)*

The directory becomes `transports/`. All references in `channels/README.md`, `specs/harness-spec-draft.md`, `DR/DR-009-...md`, `TODO/012-...md`, TE-24, and TE-26 migrate. The previously-locked phrasing "`channels/<pcid>--<slug>/`" becomes "`transports/<pcid>--<slug>/`". A redirect note is left in TE-24 and TE-26 explaining the rename so future readers can follow the trail.

- **Easier**: vocabulary is consistent across the repo from the moment TE-27 lands; no half-renamed corpus.
- **Harder**: more files touched in one TE-27 landing.

#### Alt-1.B: Defer the rename to a TODO

Lock the principle here; do the migration in a follow-on TODO commit so TE-27 stays focused on the axis analysis.

#### Alt-1.C: Don't rename; use "channel" and "transport" as synonyms

Keep `channels/` and add a glossary note that "channel" and "transport" are interchangeable in the wire-lab.

### DF-27.2 — Does "channel" survive as an above-transport concept?

#### Alt-2.A: No. The wire-lab vocabulary has only "transports" and "messages." *(LOCKED — chosen)*

If a logical addressing/grouping concept above transports is later needed, a future TE introduces it then with a fresh name.

- **Easier**: one less abstract noun in the vocabulary; fewer places where readers must disambiguate.
- **Harder**: if the project later genuinely needs an above-transport concept, we have to coin a name and explain why "channel" was rejected here.

#### Alt-2.B: Yes, reserved for future above-transport use

"Channel" is held in reserve; not used for now, not used for transports, but documented as a potential future layer.

#### Alt-2.C: Yes, redefine TE-26's locked terms now

Split TE-26's contract: physical routing is "transport," parents/receipts/kinds are "channel." Two specs, two layers. Significant churn.

### DF-27.3 — TE-27 framing scope

#### Alt-3.A: Axes-of-differentiation analysis (broad) *(LOCKED — chosen)*

TE-27 enumerates the axes (A-H above), assigns a per-axis recommendation about whether different values warrant distinct pCIDs or are parameters within one spec, and produces a starter catalogue of named transport-protocols.

- **Easier**: gives the project a meta-rule applicable to every future transport-protocol question.
- **Harder**: longer document; more decisions to make in one place.

#### Alt-3.B: Pair-vs-group only (narrow)

Just answer whether `unicast-channel-draft.md` and `group-channel-draft.md` should both exist. Defer the broader axis analysis to TE-28.

#### Alt-3.C: Two TEs in parallel

TE-27 = rename only. TE-28 = axes analysis. Cleanly separated.

### DF-27.4 — Per-axis: distinct pCID or parameter within a spec?

For each axis A-H, the meta-rule TE-27 establishes:

#### Axis A (Cardinality) → parameter within a spec

Different N within "all members see all messages" contract is operational, not contractual. One spec covers N=2, small N, and small-to-medium N. Large-N (A.3) and unbounded (A.4) cross a contract boundary (no longer all-see-all at near-source latency); these warrant distinct pCIDs.

#### Axis B (Visibility) → distinct pCID per visibility class

B.1 (all-see-all), B.2 (hub-mediated), B.3 (ring-propagated), B.4 (subset-addressed), B.5 (topic-filtered), B.6 (gossip) are observably different contracts. Each merits its own transport-protocol-pCID, possibly with cardinality and topology as parameters.

#### Axis C (Routing topology) → distinct pCID per topology class

A direct-broadcast transport, a mesh, a star, a ring, a tree, a layered overlay, and a cluster-of-clusters each leave qualitatively different on-disk artifacts when simulated. Each warrants its own transport-protocol-pCID. Topologies that are isomorphic (e.g., a 2-node ring is a 2-node mesh is a pair-direct) collapse to whichever spec is more general.

#### Axis D (Membership) → parameter within a spec, with one exception

Static, invite-only, and open-read/restricted-write can be parameters within a transport-protocol's spec (the spec defines a membership-mode). Capability-token/permissioned (D.5) introduces a separate proof obligation and warrants distinct treatment, possibly its own pCID or a shared "permissioned-transports" spec.

#### Axis E (Persistence) → parameter within a spec, with one exception

Append-only, bounded retention, and compactable are parameters. Ephemeral (E.4) crosses a contract boundary because the simulation cannot observe a transport whose messages disappear; ephemeral transports warrant a distinct pCID with explicit "out-of-band observer" semantics.

#### Axis F (Message-graph shape) → parameter within a spec

F.1-F.5 are different message-graph shapes that any transport-protocol may declare in its own spec. Whether parents are absent, single, multi, or vector-clock is a per-transport choice but does not by itself partition transport-protocols into different pCIDs.

#### Axis G (Direction) → parameter within a spec

Symmetric, hub-asymmetric, multicast, and paired are configurations a transport-protocol's spec may declare. Different direction values within one transport-protocol-pCID are common (e.g., a transport may permit asymmetric instances by configuration).

#### Axis H (Reliability / receipts) → parameter within a spec

Receipt scheme is a per-transport choice. Each transport-protocol's spec defines its own receipt format, scope, and frontier semantics.

### DF-27.5 — Starter catalogue of named transport-protocols

Given the per-axis meta-rule — cardinality within all-see-all is a parameter, not a contract boundary — a separate pair-shape spec would violate the meta-rule that TE-27 just locked. The N=2 case is just N=2 in a small-finite-closed-group transport-protocol's spec, not a separate protocol.

#### Alt-5.A: start with two specs now (rejected)

A `specs/unicast-transport-draft.md` (pair-only) and a `specs/group-transport-draft.md` (small N). This duplicates the contract for no reason; rejected because the cardinality meta-rule says they are the same contract.

#### Alt-5.B (LOCKED — chosen): start with one spec

- **`specs/group-transport-draft.md`** — Small-finite-closed-group (Axis A: A.1 and A.2 together; N≥2), all-see-all (B.1), direct (C.1), invite-only or open-read membership (D.2/D.3), append-only (E.1), multi-writer DAG (F.3), symmetric (G.1), per-message or frontier receipts (H.2/H.3/H.4). Inherits TE-24's v0 contract; the spec doc decides parent-header shape, receipt scope, and other interior details. The N=2 case (Codex↔Perplexity, the original TE-24 instance) is a documented common case of this spec, not a separate spec.

#### Alt-5.C: defer all spec carve-outs until a future TODO

TE-27 records the rename and the meta-rule; the carve-out spec comes later.

### DF-27.6 — Topology specs to draft next (after the immediate carve-out)

Anticipated next transport-protocol specs to draft once the group-transport spec lands:

- `specs/ring-transport-draft.md` (Axis C: C.4) — token-passing or ordered-relay; messages circulate; each hop leaves an artifact.
- `specs/star-transport-draft.md` (Axis C: C.3) — hub-and-spoke; the hub sees everything, spokes see only hub-mediated traffic.
- `specs/cluster-cluster-transport-draft.md` (Axis C: C.7) — small close-knit groups bridged into a larger network; matches the project's working hypothesis that the grid may be made of overlapping close-knit groups.
- `specs/gossip-transport-draft.md` (Axis B: B.6) — eventual-consistent flood; per-message latency is statistical.

These are not part of the TE-27 carve-out and are listed here only as the anticipated next TEs/TODOs.

## Scenario analysis

### S1 — Two endpoints establish a transport

Codex and Perplexity create a transport instance. Under the catalogue, this is `group-transport-draft.md` with N=2 and the slug `codex-perplexity`. Path: `transports/<pcid-of-group-transport>--codex-perplexity/`. The directory's interior is whatever the group-transport spec defines. The N=2 case is documented as the common starter case in that spec.

### S2 — A small group of three forms a transport

Codex, Perplexity, and Steve form a transport. Under the catalogue, this is the same `group-transport-draft.md` instantiated at N=3 — same pCID as S1, different participants and different slug. Path: `transports/<pcid-of-group-transport>--<slug>/`. The directory's interior follows the group-transport spec, which uses multi-writer DAG (F.3) for parents.

### S3 — A 2-participant transport adds a third participant

Under Alt-5.B (one spec covers N≥2 in the small-finite-closed-group contract), the transport stays at the same pCID; the third participant joins per the spec's membership rules. No migration to a different transport-protocol-pCID is needed because the contract did not change — only the cardinality parameter did.

This is the case the cardinality meta-rule (DF-27.4) is designed to make uneventful: a 2-participant instance and a 3-participant instance of the same group-transport-pCID share one contract.

### S4 — A ring transport simulates token-passing among five nodes

Path: `transports/<pcid-of-ring-transport>--<slug>/`. The ring-transport spec defines per-node sub-directories (or per-link records, or a single log with per-message routing-trace headers) that record the path each message traveled. The simulation can replay the ring's behavior from the directory contents alone.

### S5 — A cluster-of-clusters transport

Each cluster is itself a group-transport. Bridges between clusters relay select messages across cluster boundaries. The cluster-cluster-transport spec defines how clusters are named, how bridges are declared, and how cross-cluster messages are routed. This is the most general of the catalogued topologies and the most representative of the project's working hypothesis about real-world grid structure.

### S6 — An ephemeral transport (Axis E.4)

Messages disappear on delivery. The transport directory holds only the live message set, not history. The simulation must record messages out-of-band to observe behavior. The ephemeral-transport spec defines the on-delivery deletion rule and the out-of-band observer hook.

### S7 — Migration between transport-protocol-pCIDs *(deferred)*

When a group of participants needs to move from one transport-protocol-pCID to a different one (e.g., from the small-finite-closed-group spec to a future cluster-of-clusters spec because the group has grown beyond what the original spec's contract supports), what does that migration look like operationally? Is it the closing of one transport instance and the opening of another? Does the new transport carry a back-reference to the prior one? Do prior messages migrate in-place, get re-anchored as parents, or stay where they were?

This question is non-trivial and TE-27 does not lock answers to it. It is deferred to a future TE on transport-protocol migration semantics.

## Conclusions

1. **Rename `channels/` to `transports/`** throughout the repo. (DF-27.1 Alt-1.A.)
2. **The wire-lab vocabulary is "transports" and "messages."** "Channel" is not used; if a logical-addressing layer is later needed, a future TE introduces it. (DF-27.2 Alt-2.A.)
3. **Frame this TE as the axes-of-differentiation analysis.** (DF-27.3 Alt-3.A.)
4. **Per-axis meta-rule for distinct-pCID-vs-parameter:** axes B (visibility) and C (routing topology) warrant distinct pCIDs per class. Axis A (cardinality) is a parameter except at large-N or unbounded. Axes D (membership), E (persistence), F (message-graph), G (direction), H (receipts) are parameters within a spec, with narrow exceptions noted. (DF-27.4.)
5. **Start with one transport-protocol spec**: `specs/group-transport-draft.md` (multi-writer DAG, all-see-all, small-finite-closed-group, N≥2). The N=2 case (Codex↔Perplexity) is a documented common case of this spec, not a separate spec; per the cardinality meta-rule a pair-only spec would duplicate the contract for no reason. TE-24's v0 contract migrates into this spec. (DF-27.5 Alt-5.B.)
6. **Subsequent transport-protocol specs anticipated** (each its own future TE or TODO): ring, star, cluster-of-clusters, gossip. (DF-27.6.)
7. **The wire-lab's thin outer spec is `specs/transport-spec-draft.md`** (renamed from the previously-anticipated `channel-spec-draft.md`). It carries TE-26's four locked principles with "channel" → "transport" applied throughout. The TE-26 substance is not redone; only the vocabulary changes.

## Implications

- All TE-26 prose, decisions, and follow-on references migrate from "channel" to "transport." TE-26's locked decisions remain locked under the new vocabulary.
- `channels/README.md` is renamed and rewritten as `transports/README.md` to point at `specs/transport-spec-draft.md` and the catalogue of transport-protocol specs. The transport-shape commitments (single-writer log, 1:1) are removed; that material moves into `specs/unicast-transport-draft.md`.
- Existing TE-24 prose retains its original "channel" wording inline (the document is part of the historical record), with a top note explaining the rename and pointing readers at TE-27.
- TE-26 retains its original "channel" wording inline with a similar top note.
- `specs/harness-spec-draft.md`'s TE-24 §8 bullet is rewritten to use "transport," and a TE-27 §8 entry is added.
- DR-009 and TODO 012 are updated to use "transport."
- The TODO for the carve-out (anticipated TODO 013) covers: creating `specs/transport-spec-draft.md` (thin), `specs/unicast-transport-draft.md` and `specs/group-transport-draft.md` (substantive), and updating all the cross-references.

## Decision status

LOCKED:
- DF-27.1 — Alt-1.A (rename `channels/` → `transports/` now in TE-27's carve-out commits).
- DF-27.2 — Alt-2.A (no "channel" in the wire-lab vocabulary; transport and message only).
- DF-27.3 — Alt-3.A (axes-of-differentiation analysis).
- DF-27.4 — per-axis meta-rule as recorded above.
- DF-27.5 — Alt-5.B (one spec now: group-transport, covering N≥2 small-finite-closed-group).
- DF-27.6 — anticipated next specs: ring, star, cluster-of-clusters, gossip.

Recorded principle: *the wire-lab's `transports/` directory is a simulation surface; each subdirectory is one simulated transport instance keyed by its transport-protocol-pCID; transport-protocols differ along axes of visibility, routing topology, and (at extremes) cardinality, persistence, and membership; other axes are parameters within a single transport-protocol's spec.*

## Implications for follow-on work

- **TODO 013 (anticipated)**: drive these locked alts to a DR; carve out `specs/transport-spec-draft.md` (thin) plus `specs/group-transport-draft.md` (substantive, inheriting TE-24's v0 contract); rename `channels/` → `transports/`; update `channels/README.md` → `transports/README.md`; rewrite TE-26 in place with vocabulary swap; update DR-009, TODO 012, TE-24, and `specs/harness-spec-draft.md` to use the new vocabulary.
- **TODO 014 (anticipated)**: first real message exchange under a group-transport instance (likely N=2, Codex↔Perplexity), exercising the group-transport-draft.md spec end-to-end.
- **TE-28 (anticipated)**: transport-protocol migration semantics — what does it mean for a group of participants to move from one transport-protocol-pCID to another? (Question raised in S7 above; deferred from TE-27.)
- **TE-29 (anticipated)**: ring-transport spec — token-passing semantics, per-hop artifacts, link-failure handling.
- **TE-30 (anticipated)**: cluster-of-clusters transport — small close-knit groups bridged into a larger network. The most representative of the project's working hypothesis about real-world grid structure.
- **TE-31 (anticipated)**: gossip-transport — eventual-consistent flood, statistical latency, infection rules.
- **TE-32 (anticipated)**: receipts at scale — does `IHave: <transport-pcid>:<cid>` need to become a vector at multi-writer or large-N transports? (This was previously listed as TE-28 in TE-26 anticipated work; TE-27's introduction of intervening anticipated TEs renumbers it.)
