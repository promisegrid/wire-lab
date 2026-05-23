# Grid-Envelope Signature Prior-Art Note

Date: 2026-05-22
Owner TODO: `TODO-tugoz`
Source: `DI-dunat`; `DI-kafot`

## Purpose

Capture the external signed-envelope patterns most relevant to the next
signable-view and Gordian comparison sims, so later TODOs and simulations can
cite a repo-owned note instead of re-deriving these lessons from chat memory.

## Scope

This note is descriptive prior art. It does not declare current PromiseGrid
consensus and does not choose the winning envelope.

## Ceramic

Observed pattern:

- Ceramic signed events use DAG-JOSE / DagJWS rather than putting raw signature
  bytes directly into the application payload.
- The signed object references the actual event content by CID/link; the JWS
  payload is a base64url-encoded CID link to the event, and the signed event
  object also exposes a `link` field for traversal.

Why it matters here:

- This is strong evidence for the **wrapper-proof** family represented locally
  by specimens like `SIM-gojot`.
- Ceramic shows a clean way to avoid self-referential signing: the signature
  wrapper points at content rather than embedding a signature into the exact same
  object bytes being signed.

PromiseGrid relevance:

- Good comparison point for "proof as wrapper selected by pCID".
- Weak comparison point for "proof carried inside a minimal payload contract",
  because Ceramic deliberately externalizes the proof wrapper.

## AT Protocol

Observed pattern:

- atproto signs deterministic CBOR/DRISL bytes for authenticated objects.
- For signed labels, the `sig` field is stored in the object, but signing and
  verification are defined over the object **without** the `sig` field.
- atproto repository commits likewise sign deterministic object bytes and rely
  on canonical CBOR-like encodings for stable authenticated representations.

Why it matters here:

- This is the clearest direct prior art for the missing PromiseGrid
  **signable-view** specimen.
- atproto demonstrates that "proof lives in the object, but the signable view
  excludes the proof field" is workable when the exclusion rule is explicit and
  deterministic.

PromiseGrid relevance:

- Strong comparison point for the new explicit signable-view sim.
- Supports the requirement that the payload contract name the signable view
  explicitly rather than relying on vague "sign the object except the proof"
  folklore.

## UCAN

Observed pattern:

- UCAN separates signed content from the envelope; the UCAN envelope includes a
  signed payload and a signature component.
- The spec also allows presentation/storage in IPLD forms while converting to a
  canonical form for signature validation.

Why it matters here:

- UCAN is another clean example of **explicit envelope separation** between what
  is signed and the proof material that carries the signature.
- It reinforces that canonicalization and signature-validation shape need to be
  defined together, not left implicit.

PromiseGrid relevance:

- Useful comparison point for both outer-wrapper and payload-wrapper proof
  designs.
- Less direct than atproto for the exact inside-payload signable-view case.

## Gordian Envelope

Observed pattern from prior analysis:

- Gordian-style envelope thinking emphasizes explicit subject/assertion/proof
  structure and richer wrapper semantics.
- That is attractive for audit-heavy workflows, selective disclosure, and more
  explicit proof-bearing documents.

Why it matters here:

- It is a natural comparison family for **payload/wrapper-focused** PromiseGrid
  sims.
- It is also useful as a **negative control** for a universal outer envelope:
  if the wrapper becomes too semantically rich, it risks fighting the current
  PromiseGrid pressure toward a tiny positional outer grid.

PromiseGrid relevance:

- Worth simulating as a payload/wrapper family.
- Not yet evidence that the universal outer grid should become Gordian-like.

## Shared Lessons

Across Ceramic, atproto, and UCAN, the common lessons are:

1. **Canonical bytes must be explicit.**
   Signature rules work when the exact signed bytes are defined, not when they
   are implied by implementation folklore.
2. **Self-reference must be avoided deliberately.**
   Systems either sign a linked payload (Ceramic) or define the signable object
   as the object-without-proof (atproto), or separate payload from envelope
   explicitly (UCAN).
3. **Wrapper and payload are different design levers.**
   External proof wrappers and inside-object proofs are both viable, but they
   create different audit, parsing, and small-device costs.

## What This Suggests For The Next Sims

- The next dedicated PromiseGrid signable-view sim should use the **atproto-like
  lesson**: explicit deterministic signable projection for a proof carried in
  the payload contract.
- The Gordian comparison batch should stay outside current consensus and test
  Gordian ideas primarily as **payload/wrapper** structures, with a separate
  negative-control specimen that asks whether a Gordian-like universal outer
  envelope is too heavy.

## Sources

- Ceramic event log and DAG-JOSE documentation:
  - `https://developers.ceramic.network/docs/protocol/js-ceramic/streams/event-log`
  - `https://developers.ceramic.network/docs/protocol/js-ceramic/accounts/object-capabilities`
- AT Protocol data model and label-signing documentation:
  - `https://atproto.com/specs/data-model`
  - `https://atproto.com/specs/label`
  - `https://atproto.com/guides/data-repos`
- UCAN specification:
  - `https://github.com/ucan-wg/spec`
