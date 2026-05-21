# SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs: Grid-envelope Cryptid Multisig signature/proof probe

This simulation is a standalone, non-child grid-envelope specimen. It tests
whether envelope signatures and proofs can use Cryptid's Multisig object model
as the signature/proof payload representation while keeping PromiseGrid's
signature-placement and verification questions unresolved. Source: `DI-sahiv`.

The upstream prior-art source is Cryptid's Multisig Specification v0.0.1,
currently marked pre-draft. This simulation treats that format as design input:
the Multisig object starts with the Multisig sigil `0x1239` encoded as varuint
`0xb924`, then carries a signing-codec sigil, optional message bytes, and a
counted sequence of attributes. Source: `DI-sahiv`.

The local draft spec is
`protocols/grid-envelope.d/specs/grid-envelope-draft.md`.

## Design Pressure

- **Detached versus combined:** the Multisig message field may be empty
  (detached) or present (combined), so the sim can compare signatures over
  outer envelope bytes, nested payload bytes, or in-object message bytes.
- **Envelope versus nested payload:** the same Multisig bytes can occupy an
  outer signature/proof slot or live inside a payload protocol's nested object.
- **Variable arity:** Multisig's counted attributes let one signature object
  carry extra codec-specific proof material without changing the outer envelope
  shape.
- **pCID interaction:** the envelope `pcid` still selects payload semantics, and
  a `sig_pcid` or payload schema decides how Multisig verification is invoked.
- **Unknown codecs:** generic tools that understand varuint and varbytes can
  skip unknown Multisig signing codecs or unknown attributes without claiming
  verification.
- **Threshold shares:** attributes such as `Scheme`, `Threshold`, `Limit`,
  `ShareIdentifier`, and `ThresholdData` let the sim test individual shares,
  accumulation, and final aggregate verification.
- **Verifier obligations:** verifiers must bind the exact message bytes, signing
  codec, verifying key material, required attributes, threshold policy, and
  payload interpretation before accepting an envelope.

## Non-Canonical Status

This simulation does not choose a final PromiseGrid envelope shape, does not
freeze a pCID, does not require Cryptid Multisig as a PromiseGrid dependency,
and does not supersede the existing positional, arity, nested-signature, or
generated child grid-envelope specimens. Source: `DI-sahiv`; open long-term
envelope pressure remains represented by `DR-009-20260430-204108`.
