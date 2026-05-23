# Grid-envelope draft: UCAN-like signed-body envelope

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `signed-body-envelope-ucan-like`.

## Scope

This spec defines one grid-envelope candidate for wire-lab comparison. It is a
specimen inside `SIM-mipum-grid-envelope-signed-body-envelope-ucan-like`, not a
harness rule and not the canonical PromiseGrid envelope. Source: `DI-nohir`.

## Envelope Shape

The outer envelope shape is:

```text
[pcid, payload]
```

For this specimen, `pcid` defines `payload` as the canonical bytes of:

```text
[body_pcid, body_payload, proof_pcid, proof_payload]
```

Slots are interpreted positionally:

- `body_pcid` identifies the signed body contract.
- `body_payload` is the canonical body bytes.
- `proof_pcid` identifies the signature or proof contract.
- `proof_payload` is the proof bytes interpreted by `proof_pcid`.

## Signed Bytes

The proof validates the canonical bytes of:

```text
[body_pcid, body_payload]
```

The payload contract does not require a reserved internal proof slot inside the
body itself. Instead it explicitly separates the signed body from the sibling
proof material.

## Encoding

The outer envelope and payload object are deterministic CBOR positional arrays.
`pcid`, `body_pcid`, and `proof_pcid` are CIDv1 byte strings or links.
`body_payload` and `proof_payload` are byte strings.

## Unknown pCID Policy

If `pcid` is unknown, a receiver may preserve or blind-carry the outer bytes,
but it MUST NOT claim to parse the body/proof split.

If `pcid` is understood but `body_pcid` is unknown, a receiver may retain
`body_payload` and `proof_payload`, but it MUST mark body interpretation
unsupported.

If `proof_pcid` is unknown, it MUST mark proof validation unsupported even if
the body shape is known.

## Signature and Authorship Policy

This specimen keeps proof material inside the payload contract, but it does not
hide the body/proof separation inside one object with a reserved proof slot.

The payload contract itself says:

- these bytes are the signed body;
- these bytes are the proof material;
- proof validates the canonical body pair `[body_pcid, body_payload]`.

## Comparison Pressure

Compared with `SIM-pamap`, this specimen avoids a named
`payload_without_sig` projection and instead defines an explicit sibling
body/proof split.

Compared with `SIM-riliz`, this specimen still avoids universal outer proof
slots.

Compared with `SIM-jumav`, this specimen keeps body and proof together in one
payload contract instead of using a linked wrapper-proof object.

## Open Questions

- Is an explicit body/proof split easier to teach and audit than an inside-
  payload reserved proof slot?
- Does this shape preserve enough small-device simplicity while keeping proof
  semantics explicit?
- Should a future version require proof to bind the outer `pcid` as associated
  data, or is binding `[body_pcid, body_payload]` enough?

## Non-Canonical Status

This draft exists to compare one payload-local signed-body envelope shape
against both universal proof tuples and inside-payload explicit signable-view
rules. It does not declare a winning PromiseGrid envelope.
