# Grid-envelope draft: protocol-owned signature slot

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `protocol-owned-signature-slot`.

## Scope

This spec defines one grid-envelope candidate for wire-lab comparison. It is a
specimen inside `SIM-dalor-grid-envelope-protocol-owned-signature-slot`, not a
harness rule and not the canonical PromiseGrid envelope. Source: `DI-kukuk`.

## Envelope Shape

The outer envelope shape is:

```text
[pcid, payload, signature]
```

Slots are interpreted positionally:

- `pcid` identifies the payload protocol and the proof rules for the third slot.
- `payload` is opaque bytes until interpreted by the handler named by `pcid`.
- `signature` is mandatory proof bytes over the canonical unsigned prefix.

This is the key design move under test: the proof remains a sibling outer slot,
but there is no separate outer proof-profile selector. If a protocol wants a
varsig, multisig, or other proof family, that choice is part of the protocol
named by `pcid`.

## Signable Bytes

The signed bytes are the canonical bytes of:

```text
[pcid, payload]
```

This binds both the payload bytes and the payload protocol name without adding a
second protocol selector to the outer envelope.

## Encoding

The outer envelope is a deterministic CBOR positional array. `pcid` is a CIDv1
byte string or link as defined by the carrier profile. `payload` and
`signature` are byte strings at the carrier layer.

## Unknown pCID Policy

If a receiver lacks a handler for `pcid`, it may preserve or blind-carry the
exact outer bytes as uninterpreted evidence, but it MUST NOT claim to parse the
payload shape or verify the signature.

This specimen intentionally stays close to the current simple-envelope direction:
unsupported `pcid` means “bytes may survive, meaning does not.”

## Signature and Authorship Policy

This specimen has no universal outer `sig_pcid` slot. Instead, the protocol
named by `pcid` defines:

- whether the third slot is interpreted as varsig, multisig, or another proof
  family;
- signer binding and signer identity rules;
- delegation, freshness, revocation, and threshold semantics, if any;
- whether extra associated data beyond canonical `[pcid, payload]` bytes is
  required.

The outer envelope itself enforces only three things:

- there is a third proof slot;
- the signable baseline is canonical `[pcid, payload]`;
- proof semantics are owned by the protocol named by `pcid`.

## Comparison Pressure

Compared with `SIM-jufag`, this specimen removes the separate outer proof
selector and asks whether one payload pCID is enough to own both payload and
proof semantics.

Compared with `SIM-pamap`, this specimen keeps the proof as an outer sibling
slot rather than moving it into the payload contract.

Compared with `SIM-jumav`, this specimen avoids wrapper indirection and signs
the carried payload directly.

Compared with `SIM-kurim`, this specimen adds one universal outer proof slot
while still trying to keep the outer surface small.

## Open Questions

- Is one payload pCID enough to keep proof-family evolution legible, or does a
  separate outer proof selector age better?
- Does this design force too many proof-profile changes to mint new payload
  pCIDs?
- Do generic peers lose too much audit clarity when the outer slot is explicit
  but its proof family is hidden behind `pcid`?

## Non-Canonical Status

This draft does not declare a winning envelope and does not constrain sibling
simulations. It exists to compare the protocol-owned three-slot outer proof idea
against minimal outer envelopes, explicit outer proof selectors, payload-owned
proofs, and wrapper-proof designs.
