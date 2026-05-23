# Grid-envelope draft: Gordian selective disclosure

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `gordian-selective-disclosure`.

## Scope

This spec defines one grid-envelope candidate for wire-lab comparison. It is a
specimen inside `SIM-vizan-gordian-selective-disclosure`, not a harness rule
and not the canonical PromiseGrid envelope. Source: `DI-nohir`.

## Envelope Shape

The universal outer envelope remains:

```text
[pcid, payload]
```

For this specimen, `pcid` defines `payload` as:

```text
[subject_ref, claim_root, disclosure_root, proof_pcid, proof_payload]
```

Slots are interpreted positionally:

- `subject_ref` identifies the subject of the disclosed statement.
- `claim_root` identifies the full underlying claim graph or body root.
- `disclosure_root` identifies the disclosed subset or reveal package.
- `proof_pcid` identifies the proof profile.
- `proof_payload` is the proof bytes for the disclosed subset.

## Design Move

This specimen keeps Gordian-like structure inside a payload family, but it adds
one more pressure point: selective disclosure and reveal packages are first-
class inside the payload contract.

## Encoding

The outer envelope and payload object are deterministic CBOR positional arrays.
References and pCIDs are CIDv1 byte strings or links; proof bytes are byte
strings.

## Unknown pCID Policy

If `pcid` is unknown, a receiver may preserve or blind-carry the outer bytes,
but it MUST NOT claim to parse claim/disclosure structure.

If `proof_pcid` is unknown, it may preserve `claim_root` and `disclosure_root`
as uninterpreted references, but it MUST mark proof validation unsupported.

## Signature and Authorship Policy

This specimen binds proof to the disclosed subset rather than to the whole claim
graph alone. It is therefore testing whether selective-disclosure structure
needs to be explicit in a PromiseGrid payload family.

## Comparison Pressure

Compared with `SIM-fitin`, this specimen adds disclosure/reveal structure and
therefore more complex audit surfaces.

Compared with `SIM-suzuf`, this specimen keeps that complexity out of the
universal outer envelope.

## Open Questions

- Does selective-disclosure structure belong in PromiseGrid payload families at
  all, or is it too specialized for the baseline envelope family work?
- Do sparse peers gain enough from explicit disclosure roots to justify the
  added wrapper complexity?

## Non-Canonical Status

This draft exists to test selective-disclosure pressure inside a Gordian-style
payload family. It does not recommend a richer universal outer envelope.
