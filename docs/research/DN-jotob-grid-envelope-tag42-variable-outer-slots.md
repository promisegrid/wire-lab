# Grid Envelope: tagged selector and variable outer slots

This note explains the current PromiseGrid outer-envelope direction in plain
English. It is not a TE and it is not a frozen protocol spec. It is a
developer-facing explanation of the direction locked by `TE-fikoj` and
`DI-sisak`, refined by `TE-lamun` / `DI-rojij` / `DI-punam` / `DI-sizoh`, cross-checked against the
current `DEV-GUIDE-RESOURCES.md` snapshot, and now made explicit about the sim
lineage that fed that decision. Source: `DI-bumon`; `DI-zozov`; `DI-rojij`;
`DI-punam`; `DI-sizoh`.

## The current direction

PromiseGrid's outer envelope is currently treated as a CBOR array with this
formal shape:

```text
grid([42(pCID), ...protocol-defined-slots])
```

The recommended example profile remains:

```text
grid([42(pCID), payload, ...])
```

That means:

- the message is a CBOR array;
- slot `0` carries the protocol selector;
- the protocol named by `pCID` defines the meaning, count, order, signable view,
  validation rules, and failure behavior for every following slot;
- most protocols should still use slot `1` as the primary payload/body anchor
  unless their protocol spec has a specific reason not to.

The key point is that the universal envelope stays small. Slot `0` gets the
receiver to the right protocol spec, and that spec owns the remaining slot
vector. `grid([42(pCID), payload, ...])` is the ordinary profile, not a universal
law that every protocol must make slot `1` be payload.

## What slot `0` means

`pCID` means **Protocol CID**. It is the CID of the protocol specification
document, not the CID of a particular payload object.

We wrap the protocol CID in a CBOR `42` tag.  In the IPLD / IPFS /
Bluesky ecosystem, the `42` tag is a common way to indicate that the
following bytes are a CID-based link to another document.  In this
case, it's a link to the specification document for the protocol the
message is part of.  Using this IPLD link tag creates compatibility
with those other ecosystems and their tooling. 

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
protocol specification names the contract for every following slot.

If the IPLD-related ecosystems later become less relevant, the `42`
tag could later be replaced by another tag.  The tag itself is only a
few bytes of overhead, and in the worst case it can be ignored by
receivers that don't understand it.

## What slot `1` means

Slot `1` is the primary payload/body anchor in the recommended example profile.
A receiver that understands the `pCID` uses that protocol to interpret slot `1`
and all later slots.

The formal rule is slightly broader: the protocol named by `pCID` may assign a
different role to slot `1`, such as proof metadata, a negotiation record, or a
compact selector, if the protocol spec explicitly defines and justifies that
shape. That deviation should be uncommon. Most protocol specs should keep slot
`1` boring and payload-like.

A receiver that does **not** understand the `pCID` may still keep or relay the
exact bytes under local policy, but that is only carriage or evidence
preservation. It is not semantic acceptance.

## What later outer slots mean

PromiseGrid does not currently freeze one universal rule for slots after slot
`0`. Instead, the protocol named by `pCID` decides whether following slots exist
and what they mean. The example profile still calls slot `1` payload and treats
later slots as optional protocol-defined material.

That keeps the base envelope from hard-coding one proof format, one witness
layout, or one summary-header pattern too early.

The current hierarchy is:

- **Envelope:** `grid([42(pCID), ...protocol-defined-slots])`.
- **Payload example profile:** `grid([42(pCID), payload, ...])`.
- **Signed-message example profile:** `grid([42(pCID), payload, proof])`.
- **Single-signer proof example:** one varsig-style proof over the signable view
  defined by the pCID-named spec.
- **Multisig example form:** a pCID-defined proof set or proof chain made of
  multiple single-signer proofs. This is not a universal envelope requirement.

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

Plainly: a normal signed-message example may use a visible proof slot, but the
outer envelope does not command every protocol to use one. The pCID-named spec
must say what `proof` means, what bytes are signed, how freshness is represented,
what signer identity or key reference is used, and how a receiver records its
own local keep/break judgment. A varsig-style proof is a good compact
single-signer example. A multisig design should usually be modeled as several
agents each making their own observable promise: either an unordered proof set
when order does not matter, or an ordered proof chain when countersigning or
witness sequence matters. Threshold and aggregate signatures may compress that
evidence, but the pCID spec still has to preserve enough participant metadata
for local trust accounting.

## What a receiver does

At a high level, a receiver does this:

1. Parse the CBOR array.
2. Recover the `pCID` from slot `0`.
3. If the receiver supports that `pCID`, use the named protocol to interpret
   slots `1..N`.
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

The current direction keeps the selector stable and keeps the payload anchor as
an example profile, while letting the named protocol decide the full slot
vector when it has a justified reason to deviate.

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
which keeps slot `0` stable and preserves slot `1` as the example payload
profile while letting the protocol named by `pCID` decide the full following
slot vector.

The nearby sims matter too:

- `SIM-pobod` kept pressure on the relationship between outer-envelope evidence
  and nested payload structure;
- `SIM-jufag` kept explicit proof-format pressure visible without making
  `sig_pcid` part of the chosen universal envelope rule;
- `SIM-zukis` is the direct post-lock specimen that demonstrates one concrete
  member of the recommended example profile:
  `grid([42(pCID), payload, varsig])`.

So the current direction should be read as a locked TE/DI conclusion with clear
simulation heritage: `dalor` and the related envelope sims supplied the design
pressure, `TE-fikoj` / `DI-sisak` supplied the `42(pCID)` and variable-arity
lock, and `TE-lamun` / `DI-rojij` refined slot `1` from universal payload law
to recommended profile; `DI-punam` later narrowed the wording to recommended
example profile. Source: `DI-zozov`; `DI-rojij`; `DI-punam`.

## What is still open

Several things are still intentionally not frozen at the universal envelope
layer:

- whether proof bytes live inside payloads, in later outer slots, or both;
- unknown-`pCID` retention limits and relay policy;
- whether tag `42` remains the current selector instance forever or is later
  succeeded by another family-level selector tag while preserving the same
  semantic role;
- payload-level canonicalization and detailed proof rules;
- how strongly future specs should justify rare non-payload slot `1` layouts.

What **is** currently fixed is the direction of travel:

- CBOR array outer envelope;
- tagged protocol selector in slot `0`, currently `42(pCID)`;
- slots `1..N` defined by the protocol named by `pCID`;
- stable payload/body anchor in slot `1` as the recommended example profile.

Sources: `TE-fikoj`; `TE-lamun`; `DI-sisak`; `DI-rojij`; `DI-punam`;
`DI-sizoh`; `DI-bumon`; `DI-zozov`;
`DEV-GUIDE-RESOURCES.md`.
