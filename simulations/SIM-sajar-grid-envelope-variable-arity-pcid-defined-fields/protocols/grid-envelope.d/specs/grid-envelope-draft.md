# Grid-envelope draft: variable-arity pCID-defined fields

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `variable-arity-pcid-defined-fields`.

## Scope

This spec defines one grid-envelope candidate for wire-lab comparison. It is a
specimen inside `SIM-sajar-grid-envelope-variable-arity-pcid-defined-fields`,
not a harness rule and not the canonical PromiseGrid envelope. Source:
`DI-joman`.

## Envelope Shape

The envelope shape for this variant is:

```text
[pcid, field_1, field_2, ..., field_n]
```

Slots are interpreted positionally:

- `pcid` identifies the protocol/spec/handler that defines the remaining
  outer fields.
- `field_1` through `field_n` are interpreted only by the handler named by
  `pcid`.

The first `pcid` defines the arity, field order, field types, canonical
signature coverage, and any distinction between payload bytes, proof bytes,
routing evidence, or application data. Generic tooling can identify the first
slot and preserve the full encoded array, but it cannot know which later field
is payload or proof without the `pcid` handler.

## Encoding

This variant encodes the envelope as a deterministic CBOR positional array.
The first slot is a CIDv1 byte string. Later slots are typed by the `pcid`
spec; this draft does not impose a universal byte-string-only rule on them.
The canonical bytes for hashing and signing are the deterministic CBOR bytes of
the exact array accepted by the `pcid` handler.

## Unknown pCID Policy

If a receiver lacks a handler for `pcid`, it may store or forward the full
envelope bytes as opaque evidence, but it MUST NOT interpret any later field,
claim that a payload field exists, claim that a signature verifies, or expose
application data as understood. Unknown-variable-arity envelopes are
unsupported byte evidence until their `pcid` spec is available.

## Signature and Authorship Policy

This variant has no universal signature slot. A signature exists only when the
`pcid`-named schema defines one. That schema must define which fields are
covered, whether the first `pcid` is signed as associated data, and how signer
identity is represented. A scenario should penalize schemas that make payload
or type-substitution attacks possible.

## Layering-Test Behavior

This variant intentionally moves more responsibility into the first `pcid` than
the fixed-field variants do:

- Ordering disagreements are handled by the `pcid` schema because field count
  and field order are not shared across all envelopes.
- Unknown-pCID behavior is conservative because generic tooling cannot know
  which fields are safe to inspect.
- Signature verification is schema-local, so generic envelope tooling cannot
  verify messages without the `pcid` handler.
- Evolution can add new per-pCID schemas without changing the outer envelope
  spec, but each schema must solve arity, typing, signature coverage, and
  migration itself.

## Open Questions

- Does variable outer arity make generic routing, indexing, and audit too weak
  compared with fixed `[pcid, payload]` or `[pcid, payload, sig_pcid,
  sig_payload]` envelopes?
- Should each `pcid` schema publish machine-readable arity and field-type
  metadata, or is prose plus handler code acceptable?
- What maximum field count or encoded-size limits are required to avoid denial
  of service from adversarial arrays?
- Can signature coverage be specified reliably enough when there is no shared
  signature slot?

## Non-Canonical Status

This draft does not declare a winning envelope, does not define a central pCID
registry, and does not constrain sibling simulations. It exists to let the
runner compare variable outer arity against fixed-field and nested-payload
alternatives.
