# Grid-envelope draft: Gordian universal-envelope negative control

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `gordian-universal-envelope-negative-control`.

## Scope

This spec defines one negative-control grid-envelope candidate for wire-lab
comparison. It is a specimen inside
`SIM-suzuf-gordian-universal-envelope-negative-control`, not a harness rule and
not the canonical PromiseGrid envelope. Source: `DI-nohir`.

## Envelope Shape

This specimen deliberately expands the universal outer envelope to:

```text
[subject_ref, assertion_pcid, assertion_payload, proof_pcid, proof_payload]
```

There is no minimal universal `[pcid, payload]` prefix in this negative control.

## Design Move

This specimen exists to pressure the current PromiseGrid preference for a tiny
positional outer grid. It asks whether making subject/assertion/proof structure
universal produces better enough audit behavior to justify the extra global
parser burden.

## Encoding

The outer envelope is a deterministic CBOR positional array. References and
pCIDs are CIDv1 byte strings or links; assertion and proof bytes are byte
strings.

## Unknown pCID Policy

Because there is no universal `pcid` dispatch slot, generic peers must know
enough of the global outer contract to parse where subject, assertion, and
proof fields live. Unknown assertion or proof pCIDs are therefore more invasive
here than in the minimal-grid family.

## Signature and Authorship Policy

This negative control makes proof-bearing structure universal rather than local
to one payload family. It is intentionally testing a higher parser and policy
burden than the current minimal-grid consensus favors.

## Comparison Pressure

Compared with `SIM-fitin`, this specimen makes Gordian semantics everybody's
outer-envelope concern.

Compared with the minimal-grid family, this specimen removes the tiny universal
dispatch surface in exchange for richer built-in structure.

## Open Questions

- Does universal proof-bearing structure help enough domains to justify the
  global parse burden?
- Does the lack of a minimal outer dispatch slot create worse mixed-version and
  unknown-protocol behavior than the current minimal-grid direction allows?

## Non-Canonical Status

This draft is a negative control only. It exists to make the costs of a richer
universal outer envelope explicit, not to recommend that shape as current
PromiseGrid consensus.
