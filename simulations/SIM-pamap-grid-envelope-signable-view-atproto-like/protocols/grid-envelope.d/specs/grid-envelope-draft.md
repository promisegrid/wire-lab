# Grid-envelope draft: atproto-like explicit signable view

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `signable-view-atproto-like`.

## Scope

This spec defines one grid-envelope candidate for wire-lab comparison. It is a
specimen inside `SIM-pamap-grid-envelope-signable-view-atproto-like`, not a
harness rule and not the canonical PromiseGrid envelope. Source: `DI-nohir`.

## Envelope Shape

The outer envelope shape is:

```text
[pcid, payload]
```

For this specimen, `pcid` defines `payload` as the canonical bytes of:

```text
[payload_core, sig_pcid, sig_payload]
```

Slots are interpreted positionally:

- `payload_core` is the actual application or evidence body for this payload
  pCID.
- `sig_pcid` identifies the proof profile used by this payload contract.
- `sig_payload` is the proof bytes carried in the reserved proof slot.

## Explicit Signable View

This specimen names an explicit projection:

```text
payload_without_sig = [payload_core, sig_pcid]
```

The signed bytes are the canonical bytes of:

```text
grid([pcid, payload_without_sig])
```

This is the key design move under test: the proof is carried inside the payload
contract, but the signable view excludes the proof bytes and still binds the
payload to the outer `pcid`.

## Encoding

The outer envelope and `payload_without_sig` are deterministic CBOR positional
arrays. `pcid` and `sig_pcid` are CIDv1 byte strings or links as defined by the
carrier profile. `payload_core` and `sig_payload` are byte strings at the layer
that carries them.

## Unknown pCID Policy

If a receiver lacks a handler for `pcid`, it may preserve or blind-carry the
exact outer bytes as uninterpreted evidence, but it MUST NOT claim to parse the
payload shape or verify the proof.

If a receiver understands `pcid` but lacks a handler for `sig_pcid`, it may
parse `payload_core` and retain `sig_payload`, but it MUST mark proof
verification unsupported.

## Signature and Authorship Policy

This specimen has no universal outer proof slots. The payload contract itself
requires the proof slot and the named signable view.

The proof therefore authenticates the promiser's statement:

- that these are the exact outer `pcid` and payload-core bytes;
- that the reserved proof slot exists under this payload contract;
- that the proof profile choice `sig_pcid` is part of what is signed.

## Comparison Pressure

Compared with `SIM-riliz`, this specimen avoids universal outer proof slots.

Compared with `SIM-gojot`, this specimen does not require a separate wrapper
protocol to carry the proof.

Compared with `SIM-janov`, this specimen avoids an extra layer pCID and tests
the signable-view rule directly on one `grid([pcid, payload])` object.

## Open Questions

- Is signing `grid([pcid, payload_without_sig])` simpler and more durable than
  always requiring a separate wrapper-proof object?
- Does reserving a proof slot inside the payload contract keep enough audit
  clarity for sparse peers, or does it hide too much from generic tooling?
- Is `sig_pcid` worth signing as part of `payload_without_sig`, or does that
  make profile migration harder than necessary?

## Non-Canonical Status

This draft does not declare a winning envelope and does not constrain sibling
simulations. It exists to compare an explicit named signable view against outer
proof tuples, wrapper-proof designs, and nested signed payload layering.
