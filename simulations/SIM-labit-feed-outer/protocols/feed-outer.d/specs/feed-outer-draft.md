# Feed-Outer Spec (DRAFT)

*This simulation-local draft is the extracted specimen-side outer feed
convention from `protocols/wire-lab.d/specs/transport-spec-draft.md`. It is a
draft and is subject to revision; once frozen, its pCID will name this
protocol class for all time.*

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Source extraction: `DI-huzor` (`rusis.11`).

## Purpose

This spec defines the **outer convention** for the wire-lab's `transports/`
directory: how transport instances are named on disk, the rule that messages do
not declare their transport via a header, and the requirement that each
transport-protocol's pCID names a separate spec defining the directory's
interior.

This spec is intentionally **thin.** It does not define:

- header sets,
- parent-link semantics,
- receipt formats,
- message-kind vocabulary,
- subdirectory layout inside a transport,
- canonical-bytes rules for messages.

All of those are properties of individual transport-protocols, defined in their
own spec docs (the first being
`simulations/SIM-rakot-group-session/protocols/group-session.d/specs/group-session-draft.md`).

## Freeze boundary and historical data

Spec freeze publishes a pCID for a spec. It does not rename, rehash, or rewrite
existing transport/feed specimens. Draft-era directories remain historical
evidence; any frozen successor or derived mirror is additive and must cite the
source evidence it derives from. Source: `DI-bomud`.

Until cryptographic promise tooling exists, a Steve-authored DI is the
operative `merge-transport-spec` promise for this feed-outer lineage. Source:
`DI-bomud`.

## The four locked principles (TE-zalut)

### Principle 1: a message does not declare its transport

A message envelope contains no `Transport:` header, no `Transport-Type:` header,
and no per-message reference to which transport it belongs to. The transport
carrying a message is identified by the transport itself: in the repo-local
case, by the directory the message file lives in.

Asking a message to declare its transport is layer inversion. The transport is
the carrier; the message is the cargo; the cargo does not name the carrier.

If a message needs to *reference* a different transport (e.g., a receipt
acknowledging a message that arrived on another transport), the referencing
protocol's spec defines how to do that. This outer spec is silent.

### Principle 2: transport directories are keyed `transports/<pcid>--<slug>/`

The directory name is structured:

```
transports/<pcid>--<slug>/
```

where:

- **`<pcid>`** is the canonical pCID of the transport-protocol that transport
  speaks. This is the load-bearing identifier: it tells any reader which
  protocol's contract governs the directory's interior.
- **`<slug>`** is a human-readable suffix that tools ignore (or use only for
  display). It exists so humans can navigate `transports/` without parsing
  pCIDs and so commit-log entries are legible.
- **`--`** (double hyphen) separates the two. The double hyphen is unlikely to
  appear inside a CIDv1 base32 string.

The pCID is canonical; the slug is a convenience. Two directories with the same
pCID and different slugs are **two different transport instances** of the same
protocol. Two directories with different pCIDs are different
transport-protocols and may have entirely different interior structure.
Draft or pre-freeze specimens may use explicit draft-state names such as
`wire-lab-devs-draft`; those names are not rewritten when a spec pCID is
minted. Source: `DI-bomud`.

### Principle 3: each transport-protocol-pCID names a spec defining its directory's interior

The pCID *is* the protocol's identity. The protocol gets to define everything
inside `transports/<its-pcid>--<slug>/`:

- subdirectory layout (flat, per-direction, per-participant, sharded by date,
  etc.),
- message file naming conventions,
- header set,
- parent-link semantics (whether messages cite parents at all, what header names
  them, how multiple parents serialize, optionality),
- receipt format,
- message-kind vocabulary,
- canonical-bytes rules,
- persistence rules (append-only, bounded retention, compactable, ephemeral),
- visibility rules (all-see-all, hub-mediated, ring-propagated, etc.),
- membership rules (closed, open, invite-only, capability-token, etc.).

The feed-outer spec does not constrain any of these. They live in the
transport-protocol's own spec doc.

### Principle 4: code-as-handler

The code that reads a transport directory's structure *is* the handler for that
pCID. Each transport-protocol comes with its own reader/writer code; the pCID
identifies the protocol; the protocol identifies (by convention or naming) the
code that speaks it. There is no machine-readable companion file (no
`transport.yaml` schema). The frozen markdown spec is the human-readable
contract that the code must implement.

Tools that want to display N transport-instances of M different
transport-protocols need M handlers. That is the cost of polymorphism, not a
flaw of this design.

## What this spec does NOT specify

- The first line of a message (`grid <pcid>` is one carrier choice; not all
  transport-protocols must use it).
- Header names (`Message-ID`, `Date`, `From`, `To`, `Parents`, `IHave`, etc. —
  all defined per-protocol).
- Canonical-bytes encoding (UTF-8/LF discipline is one choice; not all protocols
  must use it).
- File-naming inside a transport directory.
- Message-CID cascade rules, legacy `Message-ID:` compatibility, and any
  reader-side rehash/deprecation policy.
- Subdirectory structure inside a transport directory.

If a future reader asks "where do I find out how to write a message for this
transport?" the answer is always: read the spec named by that transport's pCID.
The feed-outer spec is silent on the message format. Source: `DI-bomud`.
