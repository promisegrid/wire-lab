# Grid-envelope draft: tag-42 selector with protocol-owned slot-2 varsig

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `tag42-protocol-owned-slot2-varsig`.

## Scope

This spec defines one direct grid-envelope specimen for wire-lab comparison. It
is a specimen inside
`SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig`, not a harness rule
and not the canonical PromiseGrid envelope. It uses `pCID` only as Protocol
CID: the content identifier of the protocol specification document, never the
content identifier of payload bytes. Source: `DI-sisak`; `DI-mabit`.

This specimen implements one concrete member of the broader locked family:

```text
grid([42(pCID), payload, ...])
```

For this specimen, the protocol named by `pCID` chooses one later outer slot:

```text
grid([42(pCID), payload, varsig])
```

That is a protocol-owned choice for this specimen, not a universal requirement
that every PromiseGrid protocol use a third-slot `varsig`. Source: `DI-sisak`;
`DI-mabit`.

## Envelope Shape

The outer envelope shape is:

```text
[42(pCID), payload, varsig]
```

Slots are interpreted positionally:

- slot `0` is the tagged protocol selector, currently `42(pCID)`;
- slot `1` is opaque payload bytes until interpreted by the protocol named by
  `pCID`;
- slot `2` is this protocol's `varsig` proof over the signable view named by
  the same `pCID`.

The key design move under test is:

- PromiseGrid fixes the selector position and the primary payload anchor;
- the protocol named by `pCID` owns whether later outer slots exist;
- this specimen uses that freedom to place one `varsig` proof in slot `2`
  without introducing a second selector such as `sig_pcid`.

## Signable Bytes

The signable view for this specimen is the canonical bytes of:

```text
[42(pCID), payload]
```

The `varsig` in slot `2` is evidence over that exact prefix unless the protocol
named by `pCID` later refines associated-data rules more narrowly. This binds
both the tagged selector and the payload bytes without adding outer
selector-shopping machinery. Source: `DI-mabit`.

## Encoding

The outer envelope is a deterministic CBOR positional array. Slot `0` is the
tagged selector `42(pCID)`. Slots `1` and `2` are byte strings at the carrier
layer. The CBOR array header carries arity; this specimen does not add a second
arity field.

Small receivers do not need a full IPLD object model. To recover the selector
they need only:

- CBOR parsing;
- tag `42`;
- the following byte string;
- the leading `00` sentinel;
- CID parsing.

## Unknown pCID Policy

If a receiver lacks a handler for `pCID`, it may preserve or blind-carry the
exact outer bytes as uninterpreted evidence under local policy, but it MUST NOT
claim to parse the payload or verify the `varsig`.

This keeps the Promise Theory boundary explicit: bytes may survive as evidence,
but semantic acceptance remains local and protocol-dependent. Carriage is not
acceptance. Source: `DI-sisak`; `DI-mabit`.

## `varsig` Policy

This specimen has no separate `sig_pcid`, `env_pcid`, or `payload_pcid`. The
single `pCID` defines:

- what `varsig` encoding is valid in slot `2`;
- what signer binding and signer identity rules apply;
- whether freshness, delegation, threshold, or revocation semantics exist;
- whether any associated data beyond canonical `[42(pCID), payload]` bytes is
  required.

The universal envelope itself enforces only three things:

- slot `0` is the tagged selector;
- slot `1` is the primary payload anchor;
- later outer-slot roles are owned by the protocol named by `pCID`.

## Comparison Pressure

Compared with `SIM-dalor`, this specimen keeps a visible outer proof slot but
also makes the tagged selector `42(pCID)` part of the direct specimen.

Compared with `SIM-pobod`, this specimen keeps the outer shape smaller and
avoids pushing explicit nested structure into the base-envelope design.

Compared with `SIM-jufag`, this specimen removes `sig_pcid` and keeps one-pCID
discipline: one protocol selector names payload shape and slot-2 proof
semantics together. Source: `DI-mabit`.

## Open Questions

- Does one protocol-owned `varsig` slot preserve enough generic audit clarity
  without reintroducing selector shopping?
- Does this specimen outperform the fixed-three-slot `dalor` branch on the same
  scenario slice while remaining simpler than explicit-`sig_pcid` designs?
- Is slot `2 = varsig` a strong direct specimen of the broader
  `grid([42(pCID), payload, ...])` family, or does it still freeze too much
  proof structure too early?

## Non-Canonical Status

This draft does not declare a winning universal slot-2 rule. It exists to test
one direct member of the locked tagged-selector family against nearby
three-slot, nested-payload, and explicit-selector alternatives.
