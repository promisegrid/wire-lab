# Grid-envelope draft: layer pCID with nested signed payload

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `layer-pcid-nested-signed-payload`.

## Scope

This spec defines one grid-envelope candidate for wire-lab comparison. It is a
specimen inside `SIM-janov-grid-envelope-layer-pcid-nested-signed-payload`, not
a harness rule and not the canonical PromiseGrid envelope. Source: `DI-joman`.

## Envelope Shape

The outer envelope shape for this variant is:

```text
[pcid_a, payload_a]
```

Slots are interpreted positionally:

- `pcid_a` identifies a shared layer, ecosystem, or network protocol whose
  participants expect to parse `payload_a`.
- `payload_a` is opaque bytes at the outer envelope layer, but the `pcid_a`
  protocol defines its internal structure.

For this candidate, `pcid_a` defines `payload_a` as the canonical bytes of:

```text
[pcid_b, payload_b, signature_b]
```

Nested slots are interpreted by the `pcid_a` protocol:

- `pcid_b` identifies the actual payload protocol for the application data.
- `payload_b` is the actual application or evidence payload.
- `signature_b` is the Layer-A-required signature for the nested payload.

The nested signature covers the canonical bytes of `[pcid_b, payload_b]` in
this draft. That stricter coverage binds the actual payload bytes to their
payload protocol and avoids replaying the same bytes under a different
`pcid_b`.

## Encoding

Both the outer envelope and the `payload_a` interior are deterministic CBOR
positional arrays. `pcid_a` and `pcid_b` are CIDv1 byte strings. `payload_a`,
`payload_b`, and `signature_b` are byte strings at the layer that carries them.
The outer envelope's canonical bytes are the deterministic CBOR bytes of
`[pcid_a, payload_a]`; the nested signature's canonical bytes are the
deterministic CBOR bytes of `[pcid_b, payload_b]`.

## Unknown pCID Policy

If a receiver lacks a handler for `pcid_a`, it cannot participate in the
Layer-A network semantics. It may preserve the full outer envelope bytes as
opaque evidence, but it MUST NOT claim that `payload_a` has the nested shape or
that `signature_b` verifies.

If a receiver understands `pcid_a` but lacks a handler for `pcid_b`, it can
verify the nested signature over `[pcid_b, payload_b]` if the signature scheme
is known through `pcid_a`, but it MUST mark the actual payload interpretation as
unsupported.

## Signature and Authorship Policy

The outer `[pcid_a, payload_a]` layer has no fixed signature slot. This is a
deliberate pressure point for the simulation: `pcid_a` promises that
`payload_a` has a parseable signed nested structure, but the envelope itself
does not prove who made that conformance promise.

Participants therefore rely on transport identity, local peer-adoption records,
or other surrounding context to decide which agent is promising that
`payload_a` conforms to `pcid_a`. The nested `signature_b` authenticates the
actual payload claim, but it does not by itself authenticate the outer
transport context or the sender's promise that `payload_a` is a valid Layer-A
message.

## Layering-Test Behavior

This variant tests whether a commonly shared `pcid_a` can act as the layer or
network contract that most participating nodes understand:

- Generic outer tooling only needs to parse `[pcid_a, payload_a]`.
- Layer-A nodes can parse `payload_a` and find `pcid_b`, `payload_b`, and
  `signature_b`.
- Variable arity is pushed inside `payload_a`, where `pcid_a` defines it,
  rather than changing the universal outer envelope shape.
- Verification is strong for the nested payload if `signature_b` covers
  `[pcid_b, payload_b]`, but weak for the unsigned outer conformance promise
  unless transport or peer context supplies authorship.

## Open Questions

- Is relying on transport identity for the `pcid_a` conformance promise
  acceptable, or does the outer layer need its own signature/proof slot?
- Should `signature_b` cover only `payload_b`, or should it cover `[pcid_b,
  payload_b]` as this draft proposes?
- Should `pcid_a` identify a broad network/ecosystem shared by most nodes, or
  should it be a narrower layer protocol adopted by a subset of peers?
- Does this pattern improve variable-arity flexibility enough to justify the
  extra nested parsing step?

## Non-Canonical Status

This draft does not declare a winning envelope, does not define a central pCID
registry, and does not constrain sibling simulations. It exists to compare a
layer-pCID nested signed payload against fixed-field and variable-outer-arity
alternatives.
