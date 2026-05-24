# TE-vujaj: Grid-envelope DAG-CBOR link wording

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-vujaj

## Status

needs DF

## Decision under test

When PromiseGrid talks about a three-slot envelope candidate, should the repo
say the shape is `grid([42(pCID), payload, sig])`, or should it say the
conceptual shape is `grid([pCID, payload, sig])` and reserve `42(pCID)` for
DAG-CBOR-specific encoding examples?

This TE is about wording discipline at the boundary between semantic protocol
selection and wire-profile serialization. It does not reopen `TE-pokul`'s
two-slot-versus-three-slot question, and it does not by itself freeze either
CBOR or DAG-CBOR as the one PromiseGrid carrier profile.

## Assumptions

- `pCID` means Protocol CID: the CID of the protocol spec document, not the CID
  of a payload object. The protocol named by that `pCID` defines payload shape,
  signable view, proof encoding, verification rules, and related semantics.
- `42(...)` is DAG-CBOR link encoding. It means the enclosed CID is being
  serialized as a DAG-CBOR / IPLD Link value. By itself, `42(...)` does not say
  whether the CID names a protocol spec, a content object, a schema bundle, or
  some other referenced artifact.
- PromiseGrid still has an open encoding question. The current guide-resource
  snapshot says CBOR is simple and DAG-CBOR remains attractive, but encoding is
  not final.
- `TE-pokul` already narrowed the envelope-family question to minimal
  `[pCID, payload]` versus `[pCID, payload, proof]` where the protocol named by
  `pCID` defines the proof semantics. This TE must not silently turn an encoding
  notation choice into a new outer-envelope family.
- Existing DAG-CBOR positional specimens such as `SIM-gojot` can legitimately
  encode `pCID` as a DAG-CBOR Link. That specimen fact does not automatically
  answer how the abstract PromiseGrid shape should be described in guide prose,
  TE prose, or cross-profile comparisons.
- Promise Theory remains the interpretation layer: the important semantic fact is
  that the current sender uses a protocol-selected envelope to make an evidence-
  carrying promise about payload bytes. The wire encoding supports that promise;
  it does not replace the promise with a codec marker.

## Alternatives

- **Alt-A — Conceptual shape first.** Say the conceptual outer shape is
  `grid([pCID, payload, sig])`. When a DAG-CBOR profile is being shown, write the
  encoding example as `grid([42(pCID), payload, sig])` or equivalent profile-
  specific notation.
- **Alt-B — DAG-CBOR link notation as canonical phrase.** Say the outer shape is
  `grid([42(pCID), payload, sig])` even in general repo prose, on the theory that
  a pCID is normally carried as a CID link and DAG-CBOR is the clearest current
  representation of that fact.
- **Alt-C — Explicitly split abstract contract from wire profiles.** Avoid one
  single canonical phrase. Say the abstract contract is "slot 0 is the Protocol
  CID" and then give per-profile examples such as deterministic-CBOR byte string,
  DAG-CBOR Link `42(pCID)`, or another future carrier representation.

## Scenario analysis

### Scenario 1 - Cross-profile portability

Alice is writing guide prose that must remain correct whether the carrier profile
is deterministic CBOR with CID bytes, DAG-CBOR with Link values, or some future
profile that still preserves the same protocol-selection meaning.

- **Alt-A** travels well because it keeps the semantic statement independent of
  one current wire notation. The phrase says what slot 0 means, not how a
  specific codec renders it.
- **Alt-B** risks accidentally making DAG-CBOR look like the universal contract
  rather than one strong carrier candidate. A future reader may infer that
  PromiseGrid depends on DAG-CBOR even if a protocol family uses plain CBOR bytes
  or another carrier profile.
- **Alt-C** is the most explicit about abstraction boundaries, but it is heavier
  prose and may feel evasive if overused in simple discussions.

### Scenario 2 - DAG-CBOR / IPLD interoperability

Bob is implementing a DAG-CBOR/IPLD-compatible specimen and wants the repo to
say plainly that a Protocol CID is carried as a Link when the wire profile is
DAG-CBOR.

- **Alt-A** handles this cleanly by allowing DAG-CBOR examples to say
  `42(pCID)` exactly where the encoding matters.
- **Alt-B** gives Bob the strongest immediate notation match, but it also exports
  a profile-local detail into every higher-level discussion.
- **Alt-C** gives Bob precise language: "abstract slot 0 is the Protocol CID;
  under DAG-CBOR it is encoded as `42(pCID)`." This is accurate, but more verbose
  than Bob needs in every envelope comparison.

### Scenario 3 - Small-device and 100-year readability

Carol arrives decades later with only partial knowledge of IPLD and wants to
understand the protocol-selection meaning before she learns any one codec's
conventions.

- **Alt-A** is easiest for Carol because the phrase reads like a semantic
  contract: protocol selector, payload, signature.
- **Alt-B** makes Carol learn a codec-specific notation before she can even parse
  the conceptual claim. If she does not already know DAG-CBOR, the notation can
  look like the semantics are "a special wrapper around the protocol name."
- **Alt-C** is the most durable if written carefully, because it teaches the
  distinction directly. Its cost is that every mention becomes longer.

### Scenario 4 - Protocol-CID discipline

Dave wants to preserve the correction that `pCID` names the protocol spec, not a
payload object, and he wants the wording to reinforce that slot 0 is a protocol
selector.

- **Alt-A** reinforces the correction strongly. The phrase keeps the eye on
  `pCID` as a semantic role.
- **Alt-B** is still compatible with the correction, but it visually shifts
  attention toward the DAG-CBOR wrapper. The reader can come away remembering the
  wrapper more than the Protocol-CID role.
- **Alt-C** is strongest if the repo is actively teaching the distinction, but it
  may be overkill for every ordinary envelope mention.

### Scenario 5 - Mixed carriers and generic relays

Ellen is writing cross-simulation notes comparing a DAG-CBOR specimen, a
deterministic-CBOR specimen, and a future profile that carries CID bytes without
IPLD Link notation. She needs one phrase that does not unfairly bias the
comparison.

- **Alt-A** gives Ellen a neutral comparison phrase and lets each specimen spell
  out its own on-wire representation.
- **Alt-B** bakes one specimen's representation into the shared comparison
  vocabulary and risks implying the other specimens are deviations from the
  canonical form.
- **Alt-C** also works, but again at the cost of more words whenever the repo
  wants to discuss the envelope family at a high level.

### Scenario 6 - TE-pokul and the surviving envelope-family question

Frank is reading `TE-pokul` and this TE together. He needs to understand whether
`42(pCID)` introduces a new alternative, a sub-variant, or just notation.

- **Alt-A** keeps the relationship clean: `TE-pokul` owns the family question
  `[pCID, payload]` versus `[pCID, payload, proof]`, while this TE owns whether
  DAG-CBOR link notation should appear in the generic phrase.
- **Alt-B** risks making the notation look like a new family-level choice, even
  though the real distinction is still semantic slots versus profile encoding.
- **Alt-C** makes the layering explicit and therefore minimizes accidental
  duplication with `TE-pokul`.

## Cross-cutting findings

### `42(...)` is representation, not role

The repo's Protocol-CID correction makes one thing especially important: a CID's
**role** and a CID's **serialization form** are different questions. `pCID`
answers the role question. `42(...)` answers one encoding question for one CBOR
family. Treating them as if they live at the same layer invites confusion.

### DAG-CBOR examples remain valuable

Nothing in this TE argues against writing `42(pCID)` when the text is actually
about DAG-CBOR bytes, IPLD interoperability, Link traversal, or a DAG-CBOR
specimen such as `SIM-gojot`. The question is whether that notation should be
promoted to the generic repo phrase for the outer envelope.

### One phrase may not fit every audience

Guide prose, TE prose, specimen-local spec text, and low-level encoding examples
serve different readers. A durable corpus may need a neutral conceptual phrase
for cross-profile reasoning and a profile-local phrase for exact bytes.

## Conclusions

- Reject any reading where `42(pCID)` is "just the real name of the slot-0
  semantics." It is a DAG-CBOR link representation of whatever CID occupies the
  slot.
- Reject any move that would silently turn a DAG-CBOR profile detail into the
  canonical abstract statement of the PromiseGrid envelope before the repo has
  even frozen encoding.
- Keep the distinction sharp between:
  - the **semantic role**: slot 0 is the Protocol CID;
  - the **wire/profile representation**: under DAG-CBOR, that CID may be encoded
    as `42(pCID)`.
- The surviving wording choices are therefore:
  - **Alt-A:** conceptual shape first, DAG-CBOR example second;
  - **Alt-C:** explicitly separate abstract contract from profile examples.
- **Alt-B** survives only if the repo intentionally wants DAG-CBOR notation to be
  the canonical cross-profile teaching phrase. Current evidence in this corpus
  does not yet make that move feel neutral.

## Recommendation

Recommend **Alt-A** as the default repo wording:

- say the conceptual candidate shape as `grid([pCID, payload, sig])`;
- when discussing a DAG-CBOR profile or specimen, show the on-wire example as
  `grid([42(pCID), payload, sig])` or equivalent exact-byte notation.

This keeps the Protocol-CID correction legible, preserves DAG-CBOR precision
where needed, and avoids implying that DAG-CBOR is already the canonical
PromiseGrid carrier.

## Implications for open work

- `TE-pokul` remains the owner of the two-slot-versus-three-slot family
  question. `TE-vujaj` should be read as a wording/representation facet, not as
  a replacement for that TE.
- Future corrected successors under `TODO-mopob` should prefer wording that says
  the protocol named by `pCID` defines semantics, then describe DAG-CBOR Link
  notation only where the carrier profile is in scope.
- If a future DI freezes DAG-CBOR as the canonical outer carrier profile, the
  corpus may later revisit whether `42(pCID)` becomes part of the ordinary
  envelope phrase. That is not locked by current evidence.

## Decision status

`needs DF` - the TE recommends Alt-A, but the repo has not yet locked whether
general PromiseGrid prose should prefer conceptual `grid([pCID, payload, sig])`
wording or a stricter abstract-versus-profile split.
