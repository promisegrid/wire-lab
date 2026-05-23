# Grid-envelope draft: Ceramic-like wrapper proof

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `wrapper-proof-ceramic-like`.

## Scope

This spec defines one grid-envelope candidate for wire-lab comparison. It is a
specimen inside `SIM-jumav-grid-envelope-wrapper-proof-ceramic-like`, not a
harness rule and not the canonical PromiseGrid envelope. Source: `DI-nohir`.

## Envelope Shape

The outer envelope shape is:

```text
[pcid, payload]
```

For this specimen, `pcid` defines `payload` as the canonical bytes of a wrapper
object:

```text
[content_ref, proof_pcid, proof_payload]
```

Slots are interpreted positionally:

- `content_ref` is a CID or content-addressed link to the actual content object.
- `proof_pcid` identifies the wrapper-proof profile.
- `proof_payload` is the detached proof bytes interpreted by `proof_pcid`.

## Signed Bytes

The wrapper proof signs canonical bytes derived from the linked-content
reference and wrapper profile, not the exact bytes of the wrapper object that
contains the proof.

This keeps proof and signed content structurally separate. The wrapper payload
therefore says: "this proof profile attests to this content reference", instead
of "this object signs itself minus one field".

## Encoding

The outer envelope and wrapper payload are deterministic DAG-CBOR positional
arrays. `pcid`, `content_ref`, and `proof_pcid` are CIDv1 byte strings or links.
`proof_payload` is a byte string.

## Unknown pCID Policy

If `pcid` is unknown, a receiver may preserve or blind-carry the outer bytes as
uninterpreted evidence, but it MUST NOT parse the wrapper or claim proof
verification.

If `pcid` is understood but `proof_pcid` is unknown, a receiver may retain the
content reference and wrapper bytes, but it MUST mark proof validation
unsupported.

## Signature and Authorship Policy

This specimen treats proof carriage as a wrapper-family concern rather than a
universal outer-envelope concern. A receiver can often reason about the wrapper
and the linked content separately:

- the wrapper identifies what content is being attested;
- the proof validates that attestation under `proof_pcid`;
- the actual content object may be retrieved, cached, or audited independently.

## Comparison Pressure

Compared with `SIM-pamap`, this specimen avoids an explicit inside-payload
signable view by keeping proof material detached from the referenced content.

Compared with `SIM-gojot`, this specimen makes the wrapper shape explicit rather
than leaving proof-wrapper meaning to a generic "wrapper pCID" family.

Compared with `SIM-maraz`, this specimen uses linked content and proof wrapper
separation instead of a fixed signed outer audit header.

## Open Questions

- Does proof-over-link improve 100-year auditability enough to justify the extra
  object/reference layer?
- Do sparse peers benefit more from linked wrapper proofs or from summary-header
  approaches like `SIM-maraz`?
- Is linked-proof indirection too heavy for small devices compared with inside-
  payload signable views?

## Non-Canonical Status

This draft exists to compare wrapper-proof behavior against inside-payload
signable views and heavier audit-header designs. It does not declare a winning
envelope or wrapper family.
