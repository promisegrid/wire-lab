# Grid-envelope draft: Gordian payload/wrapper

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `gordian-payload-wrapper`.

## Scope

This spec defines one grid-envelope candidate for wire-lab comparison. It is a
specimen inside `SIM-fitin-gordian-payload-wrapper`, not a harness rule and not
the canonical PromiseGrid envelope. Source: `DI-nohir`.

## Envelope Shape

The universal outer envelope remains:

```text
[pcid, payload]
```

For this specimen, `pcid` defines `payload` as a Gordian-style payload/wrapper:

```text
[subject_ref, assertion_pcid, assertion_payload, proof_pcid, proof_payload]
```

## Design Move

The point under test is not "make the outer grid Gordian". The point is:

- keep the universal outer parse surface tiny;
- allow one payload family to expose richer subject/assertion/proof structure;
- compare whether that richer payload family helps audit-heavy scenarios enough
  to justify its complexity.

## Encoding

The outer envelope and payload wrapper are deterministic CBOR positional arrays.
References and pCIDs are CIDv1 byte strings or links; assertion and proof bytes
are byte strings.

## Unknown pCID Policy

If `pcid` is unknown, a receiver may preserve or blind-carry the outer bytes,
but it MUST NOT claim to interpret the Gordian payload structure.

If `assertion_pcid` or `proof_pcid` is unknown, it may retain the wrapper bytes
but MUST mark assertion or proof interpretation unsupported.

## Signature and Authorship Policy

This specimen deliberately makes subject/assertion/proof structure a payload-
family concern. Generic grid peers still only know `[pcid, payload]`; Gordian
structure is visible only to peers that support this payload pCID.

## Comparison Pressure

Compared with `SIM-maraz`, this specimen is more openly proof-structured but
does not force a fixed signed outer summary header.

Compared with `SIM-natim`, this specimen makes proof structure a payload-family
feature rather than a mandatory outer attestation.

## Open Questions

- Does this richer wrapper meaningfully improve sparse audit and multi-actor
  attribution?
- Is the extra subject/assertion/proof structure still small-device friendly
  enough when kept out of the universal outer envelope?

## Non-Canonical Status

This draft exists to test Gordian-style structure as a payload family. It does
not claim that the universal outer grid should adopt Gordian semantics.
