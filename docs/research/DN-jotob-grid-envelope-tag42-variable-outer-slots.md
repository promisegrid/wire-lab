# Grid Envelope: tagged selector and variable outer slots

This note explains the current PromiseGrid outer-envelope direction in plain
English. It is not a TE and it is not a frozen protocol spec. It is a
developer-facing explanation of the direction locked by `TE-fikoj` and
`DI-sisak`, cross-checked against the current `DEV-GUIDE-RESOURCES.md`
snapshot, and now made explicit about the sim lineage that fed that decision.
Source: `DI-bumon`; `DI-zozov`.

## The current direction

PromiseGrid's outer envelope is currently treated as a CBOR array with this
general shape:

```text
grid([42(pCID), payload, ...])
```

That means:

- the message is a CBOR array;
- slot `0` carries the protocol selector;
- slot `1` carries the primary payload bytes;
- later outer slots may exist, but only when the protocol named by `pCID`
  defines them.

The key point is that the universal envelope stays small, while the protocol
named by `pCID` owns any extra outer structure.

## What slot `0` means

`pCID` means **Protocol CID**. It is the CID of the protocol specification
document, not the CID of a particular payload object.

The current standard way to carry that selector in slot `0` is
`42(pCID)`.  In IPLD / IPFS / Bluesky contexts, the `42` tag is a
common way to indicate that the following bytes are a CID-based link
to another document, in this case the specification document for the
protocol the message is part of.  Using this IPLD link tag creates
compatibility with those other ecosystems and their tooling. If those
ecosystems later become less relevant, the `42` tag could later be
succeeded by another tag.

In practical terms, that compatibility buys:

- easier reuse of CID / DAG-CBOR tooling;
- lower friction for developers already working with IPLD, IPFS, or Bluesky;
- a more natural fit for archival, CAS, and adjacent protocol experiments that
  already expect tagged CID-bearing objects.

It does **not** mean:

- IPLD now defines PromiseGrid semantics;
- slot `0` semantically means "DAG-CBOR wrapper" rather than Protocol CID;
- trust becomes global, centrally authorized, or delegated to outside tooling;
- later outer-slot roles stop being owned by the protocol named by `pCID`.

So the compatibility choice is about wire-level interop and developer
comfort. The extra wrapper bytes are accepted as durable boilerplate
because that interop value is worth the cost.

The important semantic point is unchanged: slot `0` tells the receiver which
protocol specification names the payload contract and any later outer-slot
rules.

## What slot `1` means

Slot `1` is the primary payload anchor. A receiver that understands the `pCID`
uses that protocol to interpret the payload bytes in slot `1`.

A receiver that does **not** understand the `pCID` may still keep or relay the
exact bytes under local policy, but that is only carriage or evidence
preservation. It is not semantic acceptance.

## What later outer slots mean

PromiseGrid does not currently freeze one universal rule for slots `2..N`.
Instead, the protocol named by `pCID` decides whether later outer slots exist
and what they mean.

That keeps the base envelope from hard-coding one proof format, one witness
layout, or one summary-header pattern too early.

One direct specimen of this idea is now `SIM-zukis`, where the
protocol defines:

```text
grid([42(pCID), payload, varsig])
```

In that specimen, slot `2` is a `varsig` proof. The important point is
that this is a **protocol-owned example**, not a universal PromiseGrid
law that slot `2` must always be `varsig`.  Other signature formats,
proof formats, or later-slot roles are still possible in other
protocols named by other `pCID`s.

## What a receiver does

At a high level, a receiver does this:

1. Parse the CBOR array.
2. Recover the `pCID` from slot `0`.
3. If the receiver supports that `pCID`, use the named protocol to interpret
   slot `1` and any later outer slots.
4. If the receiver does not support that `pCID`, it may preserve the exact
   bytes as evidence under local policy, but it does not claim to understand or
   accept the message semantically.

This keeps trust local. The envelope helps peers recognize and evaluate
promises, but it does not create a central authority, global permission system,
or universal trust ledger.

## Why this direction was chosen

The main alternative was to freeze a universal fixed three-slot shape:

```text
grid([42(pCID), payload, sig])
```

That would make one outer proof slot globally regular, which is attractive for
generic readers. But it also freezes one universal outer-slot story too early.

The current direction keeps the selector and payload anchor stable while letting
the named protocol decide whether it needs an outer proof slot, a witness set,
or no later outer slots at all.

In short:

- fixed three-slot gives more universal regularity;
- variable outer slots give more protocol-owned evolution without changing the
  envelope family.

The current direction chose the second path.

## Design heritage

This direction has a clear lineage. It did not appear from nowhere at the TE
layer.

`SIM-dalor` is the most important direct ancestor because it carried the key
idea that an outer proof slot can be **protocol-owned** rather than frozen as a
universal envelope law. That pressure survives into the current direction,
which keeps slot `0` and slot `1` stable while letting the protocol named by
`pCID` decide whether later outer slots exist and what they mean.

The nearby sims matter too:

- `SIM-pobod` kept pressure on the relationship between outer-envelope evidence
  and nested payload structure;
- `SIM-jufag` kept explicit proof-format pressure visible without making
  `sig_pcid` part of the chosen universal envelope rule;
- `SIM-zukis` is the direct post-lock specimen that demonstrates one concrete
  member of the chosen family: `grid([42(pCID), payload, varsig])`.

So the current direction should be read as a locked TE/DI conclusion with clear
simulation heritage: `dalor` and the related envelope sims supplied the design
pressure, while `TE-fikoj` and `DI-sisak` supplied the actual locking
authority. Source: `DI-zozov`.

## What is still open

Several things are still intentionally not frozen at the universal envelope
layer:

- whether proof bytes live inside payloads, in later outer slots, or both;
- unknown-`pCID` retention limits and relay policy;
- whether tag `42` remains the current selector instance forever or is later
  succeeded by another family-level selector tag while preserving the same
  semantic role;
- payload-level canonicalization and detailed proof rules.

What **is** currently fixed is the direction of travel:

- CBOR array outer envelope;
- tagged protocol selector in slot `0`, currently `42(pCID)`;
- stable payload anchor in slot `1`;
- later outer slots defined by the protocol named by `pCID`.

Sources: `TE-fikoj`; `DI-sisak`; `DI-bumon`; `DI-zozov`;
`DEV-GUIDE-RESOURCES.md`.
