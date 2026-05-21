# Grid-envelope draft: Cryptid Multisig signature/proof payloads

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `cryptid-multisig-signature-proofs`.

## Scope

This spec defines one grid-envelope candidate for wire-lab comparison. It is a
specimen inside `SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs`,
not a harness rule and not the canonical PromiseGrid envelope. Source:
`DI-sahiv`.

The design under test is narrow: use Cryptid's Multisig object model as the
signature/proof payload representation while leaving the envelope-level
placement and verification choices open for simulation pressure. The Multisig
source is upstream prior art and is not treated as a frozen PromiseGrid
dependency. Source: `DI-sahiv`.

## Prior-Art Multisig Shape

Cryptid's pre-draft Multisig v0.0.1 encodes a codec-agnostic digital signature
object as:

```text
multisig_sigil signing_codec_sigil message attributes
```

The object starts with Multisig sigil `0x1239`, encoded as varuint `0xb924`.
It then carries a signing-codec sigil, an optional `message` encoded as
varbytes, and a variable number of attributes encoded as a count followed by
attribute-id and varbytes pairs. Source: `DI-sahiv`.

This simulation recognizes these upstream attribute roles as pressure inputs:

- `SigData` for signature bytes;
- `PayloadEncoding` for the signed-message encoding sigil;
- `Scheme` for threshold-signing scheme;
- `Threshold` for the minimum share count required;
- `Limit` for total share count;
- `ShareIdentifier` for the share number or participant-local share label;
- `ThresholdData` for codec-specific threshold material;
- `AlgorithmName` for application-specific or non-standard algorithm naming.

## Envelope Shapes Under Test

The simulation keeps three placement modes alive instead of choosing one:

```text
[pcid, payload, sig_pcid, multisig]
[pcid, payload_with_nested_multisig]
[pcid, combined_multisig]
```

Slots are interpreted positionally only when the selected mode and its `pcid`
or `sig_pcid` define them:

- `pcid` identifies the payload protocol, handler, or proof-bearing schema.
- `payload` is opaque bytes until interpreted by the handler named by `pcid`.
- `sig_pcid`, when present, identifies the signature/proof protocol that says
  the fourth slot is a Cryptid-style Multisig object.
- `multisig` is the exact Multisig object bytes, not decoded fields projected
  into the envelope.

The first mode pressures explicit outer signature dispatch. The second mode
pressures nested payload ownership of signatures. The third mode pressures
combined Multisig objects where the signed message is carried inside the
Multisig `message` varbytes field rather than as a sibling envelope payload.
Source: `DI-sahiv`; `DR-009-20260430-204108`.

## Encoding

The envelope carrier for this specimen is deterministic CBOR positional arrays.
`pcid` and `sig_pcid`, when present, are CIDv1 byte strings or DAG-CBOR links as
defined by the concrete run profile. `payload` and `multisig` are byte strings.
The Multisig object itself keeps its own varuint and varbytes internal encoding;
the envelope does not translate Multisig attributes into CBOR fields.

Canonical envelope bytes are the deterministic CBOR bytes of the selected
outer array. Canonical Multisig bytes are the exact varuint/varbytes bytes
carried in the Multisig slot or nested payload. A verifier must never verify a
re-serialized approximation when the original bytes are available.

## Detached and Combined Signature Policy

A Multisig object with an empty message field is treated as detached. The
verifier must obtain the signed bytes from the selected envelope mode:

- outer explicit mode signs the canonical unsigned prefix `[pcid, payload]`
  unless `sig_pcid` defines stricter associated data;
- nested mode signs the nested bytes selected by the payload protocol named by
  `pcid`;
- share-collection mode signs the same byte string for every share before
  threshold aggregation.

A Multisig object with a non-empty message field is treated as combined. The
verifier must compare the embedded message bytes against the envelope-selected
payload binding before accepting the envelope. If the embedded message and
outer payload disagree, verification fails even if the Multisig cryptographic
check succeeds. Source: `DI-sahiv`.

## Unknown Codec and Attribute Policy

This specimen adopts a conservative unknown-codec rule:

- A receiver that does not understand the envelope `pcid` must not interpret
  `payload`, even if it can skip or parse Multisig framing.
- A receiver that understands `pcid` but not `sig_pcid` may preserve the exact
  Multisig bytes as opaque proof evidence, but must not claim verification.
- A receiver that understands Multisig varuint/varbytes framing but not the
  signing codec may skip the object, index its byte range, and preserve it for
  later tooling, but must not treat unknown `SigData` as valid.
- A receiver may ignore unknown non-critical attributes only if the signing
  codec or `sig_pcid` says they are non-critical; otherwise unknown attributes
  keep the verification result unsupported.

This policy intentionally separates structural skippability from cryptographic
acceptance. Skipping unknown signing codecs helps storage, relay, and future
audit; it is not a validity decision. Source: `DI-sahiv`.

## Threshold and Multi-Payload Pressure

Threshold runs use the Multisig attributes as follows:

- every share carries `Scheme`, `Threshold`, `Limit`, and `ShareIdentifier`
  values that must agree with the threshold policy named by `sig_pcid` or the
  payload schema;
- `SigData` carries the share or aggregate signature data as defined by the
  signing codec;
- `ThresholdData` carries codec-specific accumulation material when the codec
  needs more than raw signature bytes;
- `PayloadEncoding` records the signed-message encoding when the signature
  codec requires that value for replay-safe verification;
- `AlgorithmName` is advisory unless the selected `sig_pcid` makes it part of
  the verification policy.

The simulation should penalize designs that let two shares over different
message bytes, different `pcid` values, or different threshold policies
aggregate as if they belonged to the same proof.

## pCID Interaction

The envelope `pcid` continues to name the payload protocol. Multisig does not
replace pCID-selected payload semantics. A `sig_pcid`, nested payload schema, or
future frozen profile must define:

- whether the signed bytes include `pcid`, `payload`, both, or an enclosing
  transcript;
- whether `PayloadEncoding` duplicates, complements, or constrains the pCID;
- which Multisig signing codecs are acceptable for the payload protocol;
- how verifier key material is found, authenticated, rotated, and revoked;
- whether unknown attributes are reject, ignore, quarantine, or relay-only
  evidence.

Until those choices are frozen, this specimen is evidence for verifier
obligations rather than a PromiseGrid validity rule. Source: `DI-sahiv`;
`DR-009-20260430-204108`.

## Verifier Obligations

A verifier accepts an envelope under this specimen only after all of the
following hold:

- the envelope shape is recognized by `pcid`, `sig_pcid`, or a local run
  profile;
- the exact signed byte string is determined without ambiguity;
- detached and combined-message bindings agree with the selected envelope mode;
- the signing codec is understood and allowed by the selected signature policy;
- required attributes are present, canonical, non-duplicated unless the codec
  permits duplicates, and internally consistent;
- threshold shares all bind to the same message, pCID context, threshold policy,
  and signer set before aggregation;
- unknown codecs or critical attributes produce an unsupported or quarantined
  result rather than a successful verification;
- local audit records retain exact envelope bytes, exact Multisig bytes, the
  verification profile, and the reason for accept, reject, unsupported, or
  quarantine.

## Scenario Pressure Notes

### Normal detached signature

Alice sends `[pcid, payload, sig_pcid, multisig]` where `multisig` has an empty
message field. Bob verifies that the signature covers the canonical unsigned
prefix `[pcid, payload]`, not only `payload`, so replay under a different
payload protocol fails.

### Combined signature mismatch

Carol sends `[pcid, payload, sig_pcid, multisig]` where the Multisig message
field contains bytes that differ from `payload`. Dave must reject the envelope
because the embedded message and envelope-selected payload binding disagree.

### Unknown signing codec

Ellen understands the envelope `pcid` and Multisig framing but lacks the
signing-codec implementation. Ellen may keep and relay the exact proof bytes as
unsupported evidence, but she must not mark the envelope verified.

### Threshold share collection

Frank receives three BLS-style shares that claim a threshold of three out of
four. He aggregates only shares whose message bytes, `pcid` binding, scheme,
threshold, limit, and signer-set policy match; mixed-context shares remain
separate unsupported evidence.

### Nested payload proof

Alice sends `[pcid, payload_with_nested_multisig]`. Bob can verify only after
the `pcid` payload schema identifies the nested Multisig byte range and the
exact bytes the nested proof signs.

## Non-Goals

This draft does not:

- declare a winning envelope;
- freeze a pCID;
- require a central pCID registry;
- decide detached versus combined signatures globally;
- decide whether signatures live in the envelope or nested payload;
- define final PromiseGrid key discovery, revocation, freshness, or authority;
- claim that Cryptid Multisig is stable or normative for PromiseGrid.

## Freeze Gate

This draft can freeze only after simulation runs compare it against sibling
grid-envelope signature-placement specimens, at least one verifier profile
fully specifies pCID binding and unknown-codec behavior, and a maintainer signs
a merge/freeze promise for this specific specimen.
