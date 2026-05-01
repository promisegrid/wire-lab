# Transport Spec (DRAFT)

*This is the wire-lab's outer spec for the `transports/` directory. It is a draft and is subject to revision; once frozen, its pCID will name this protocol class for all time. See `specs/MANIFEST.md` for freeze status.*

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted. Cross-references to it in other repo files use `specs/transport-spec-draft.md` (path) until freeze; after freeze they will use the pCID.

## Purpose

This spec defines the **outer convention** for the wire-lab's `transports/` directory: how transport instances are named on disk, the rule that messages do not declare their transport via a header, and the requirement that each transport-protocol's pCID names a separate spec defining the directory's interior.

This spec is intentionally **thin.** It does not define:
- header sets,
- parent-link semantics,
- receipt formats,
- message-kind vocabulary,
- subdirectory layout inside a transport,
- canonical-bytes rules for messages.

All of those are properties of individual transport-protocols, defined in their own spec docs (the first being `specs/group-transport-draft.md`).

## Sources

This spec is locked by the conclusions of:

- [TE-24](../docs/thought-experiments/TE-20260430-204108-group-transport-envelope.md): `grid <pcid>` as the group-transport envelope; canonical-bytes rules; explicit-promise body requirement. (Source for the group-transport-protocol's contract; not a constraint on this outer spec.)
- [TE-26](../docs/thought-experiments/TE-20260430-215624-channel-transport-types-and-threaded-replies.md): transport-protocol types, pCID-keyed transport paths, and DAG message graphs. Establishes the four principles below.
- [TE-27](../docs/thought-experiments/TE-20260501-021921-transports-rename-and-axes-of-differentiation.md): transports rename and axes of transport-protocol differentiation. Establishes the per-axis meta-rule for distinguishing transport-protocols.
- [DR-009](../DR/DR-009-20260430-204108-group-transport-envelope.md): the active decision request governing the group-transport envelope and its graduation.

## The four locked principles (TE-26)

### Principle 1: a message does not declare its transport

A message envelope contains no `Transport:` header, no `Transport-Type:` header, and no per-message reference to which transport it belongs to. The transport carrying a message is identified by the transport itself: in the repo-local case, by the directory the message file lives in.

Asking a message to declare its transport is layer inversion. The transport is the carrier; the message is the cargo; the cargo does not name the carrier.

If a message needs to *reference* a different transport (e.g., a receipt acknowledging a message that arrived on another transport), the referencing protocol's spec defines how to do that. This outer spec is silent.

### Principle 2: transport directories are keyed `transports/<pcid>--<slug>/`

The directory name is structured:

```
transports/<pcid>--<slug>/
```

where:

- **`<pcid>`** is the canonical pCID of the transport-protocol that transport speaks. This is the load-bearing identifier: it tells any reader which protocol's contract governs the directory's interior.
- **`<slug>`** is a human-readable suffix that tools ignore (or use only for display). It exists so humans can navigate `transports/` without parsing pCIDs and so commit-log entries are legible.
- **`--`** (double hyphen) separates the two. The double hyphen is unlikely to appear inside a CIDv1 base32 string.

The pCID is canonical; the slug is a convenience. Two directories with the same pCID and different slugs are **two different transport instances** of the same protocol. Two directories with different pCIDs are different transport-protocols and may have entirely different interior structure.

### Principle 3: each transport-protocol-pCID names a spec defining its directory's interior

The pCID *is* the protocol's identity. The protocol gets to define everything inside `transports/<its-pcid>--<slug>/`:
- subdirectory layout (flat, per-direction, per-participant, sharded by date, etc.),
- message file naming conventions,
- header set,
- parent-link semantics (whether messages cite parents at all, what header names them, how multiple parents serialize, optionality),
- receipt format,
- message-kind vocabulary,
- canonical-bytes rules,
- persistence rules (append-only, bounded retention, compactable, ephemeral),
- visibility rules (all-see-all, hub-mediated, ring-propagated, etc.),
- membership rules (closed, open, invite-only, capability-token, etc.).

The wire-lab transport-spec does not constrain any of these. They live in the transport-protocol's own spec doc.

### Principle 4: code-as-handler

The code that reads a transport directory's structure *is* the handler for that pCID. Each transport-protocol comes with its own reader/writer code; the pCID identifies the protocol; the protocol identifies (by convention or naming) the code that speaks it. There is no machine-readable companion file (no `transport.yaml` schema). The frozen markdown spec is the human-readable contract that the code must implement.

Tools that want to display N transport-instances of M different transport-protocols need M handlers. That is the cost of polymorphism, not a flaw of this design.

## The per-axis meta-rule (TE-27)

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

## What this spec does NOT specify

- The first line of a message (`grid <pcid>` is one carrier choice; not all transport-protocols must use it).
- Header names (`Message-ID`, `Date`, `From`, `To`, `Parents`, `IHave`, etc. — all defined per-protocol).
- Canonical-bytes encoding (UTF-8/LF discipline is one choice; not all protocols must use it).
- File-naming inside a transport directory.
- Subdirectory structure inside a transport directory.

If a future reader asks "where do I find out how to write a message for this transport?" the answer is always: read the spec named by that transport's pCID. The wire-lab transport-spec is silent on the message format.

## Open questions

- **OQ-1 (deferred):** Should the wire-lab spec define a small companion convention for transport-protocols to publish their own pCID on first use, so receivers can discover the protocol-spec from a stranger's first message? Raised and deferred in TE-26 DF-26.8 Alt-8.B; the locked Alt-8.C (code-as-handler) does not address this. May surface in a future TE.

- **OQ-2 (deferred):** What does it mean for a group of participants to migrate from one transport-protocol-pCID to another transport-protocol-pCID? (S7 of TE-27.) Deferred to a future TE on transport-protocol migration semantics.

## Freeze gate

This spec graduates to frozen status when:

1. The repo has at least one transport-protocol spec frozen (currently anticipated to be `specs/group-transport-draft.md`).
2. Steve signs a `merge-transport-spec` promise authorizing the freeze.
3. `tools/spec freeze transport-spec` mints the pCID, snapshots the file, and appends the manifest entry.

Until then, the spec lives at `specs/transport-spec-draft.md` and is a working draft.
