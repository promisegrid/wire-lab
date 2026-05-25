# TE-tilir: slot-0 IPLD interop and bootstrap supersedence

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-tilir

## Status

superseded by TE-fikoj / DI-kafat

## Decision under test

`TE-titaj` concluded that `42(pCID)` should not be the general long-horizon
standard for slot 0 and that PromiseGrid should keep `pCID` as the semantic
role while leaving open whether to freeze a tiny envelope-level bootstrap rule
for exact bytes. Since that TE landed, one missing pressure packet has become
explicit:

- PromiseGrid wants to interoperate well with the wider CID / IPLD / DAG-CBOR
  ecosystem that now includes ATProto / Bluesky-shaped stacks and other
  content-addressed projects.

Another earlier concern also deserves recalibration:

- the wire-level difference between a raw-CID byte string and a DAG-CBOR
  `42(pCID)` Link is small enough that it may be reasonable to treat the extra
  bytes as durable bootstrap boilerplate rather than as a prohibitive burden.

This TE asks: given those two pressures, what should PromiseGrid now treat as
the right relationship between the **semantic meaning** of slot 0 and the
**preferred wire/profile representation** of slot 0?

This TE supersedes `TE-titaj`'s framing and recommendation. It does not reopen
`TE-pokul`'s two-slot-versus-three-slot family question, and it does not make
`TE-vujaj`'s wording facet disappear. `TE-vujaj` still owns whether ordinary
repo prose should say `pCID` or `42(pCID)` in general text.

## Assumptions

- A `pCID` is a Protocol CID: the CID of a protocol spec document, never the
  CID of a payload object.
- The protocol named by `pCID` defines payload shape, signable view,
  proof/signature encoding, verification rules, and related semantics.
- The whole message is CBOR. Array arity is carried by the CBOR array header.
- Promise Theory remains the semantic layer: slot 0 helps a receiver interpret
  another agent's evidence-bearing promise; it does not impose behavior or
  create non-local trust.
- All trust judgments remain local.
- A receiver must recover the `pCID` before it can consult the protocol spec
  named by that `pCID`. Therefore the protocol spec cannot bootstrap slot 0.
- `42(...)` is DAG-CBOR/IPLD Link notation. It is one representation of a CID,
  not the semantic role of the CID.
- In the common binary case discussed here:
  - a raw CID carried as a CBOR byte string looks like a CBOR `bstr` header
    followed by CID bytes;
  - a DAG-CBOR Link adds tag `42`, a `bstr` wrapper, and the IPLD Link `00`
    sentinel before the same CID bytes.
- Rich runtimes can usually carry either form easily. The design pressure is
  whether the extra wrapper burden is acceptably small for constrained and
  archival receivers given the interop it buys.

## Alternatives

- **Alt-A - Semantic-only standard.** PromiseGrid standard prose says the slot-0
  shape is `grid([pCID, payload, sig])` and leaves profile preference open.
  `42(pCID)` is permitted in DAG-CBOR/IPLD profiles, but not preferred.
- **Alt-B - Semantic standard plus preferred DAG-CBOR/IPLD profile.** PromiseGrid
  standard prose keeps slot 0 semantically as the Protocol CID, while also
  stating that `42(pCID)` is the strongly preferred representation when the
  carrier profile is DAG-CBOR/IPLD-oriented.
- **Alt-C - Universal `42(pCID)` standard.** PromiseGrid standard prose says
  the canonical slot-0 form is `grid([42(pCID), payload, sig])`, effectively
  making the DAG-CBOR/IPLD Link wrapper the default contract for the base
  envelope.

## Scenario analysis

### Scenario 1 - Alice on a constrained receiver revisited

Alice runs a tiny receiver on a low-cost controller. She has a small CBOR
parser and basic CID support or locally archived CID decoding notes. `TE-titaj`
treated `42(pCID)` as meaningfully heavier than raw CID bytes.

- **Alt-A** gives Alice the lightest normative burden. She can keep using a
  plain byte-string bootstrap if that is what her chosen profile does.
- **Alt-B** asks more of Alice, but not a full IPLD stack. She needs to
  recognize tag `42`, parse the following `bstr`, verify the `00` sentinel, and
  then parse CID bytes. That is a small decoder increment, not a large object
  model.
- **Alt-C** makes Alice carry that increment unconditionally, even where a
  deployment gains little from DAG-CBOR/IPLD-native interop.

Finding: the small-device objection is weaker than `TE-titaj` implied, but it
does not disappear. The wrapper is small, not free.

### Scenario 2 - Bob in an IPLD-native implementation

Bob runs a PromiseGrid implementation that already stores and transmits data in
CID / DAG-CBOR / IPLD-native structures. He wants slot 0 to fit naturally into
those tools.

- **Alt-A** works, but it treats Bob's ecosystem as incidental rather than as a
  first-class interoperability target.
- **Alt-B** matches Bob's environment well: PromiseGrid still speaks in
  semantic-role terms, but the preferred DAG-CBOR/IPLD rendering is explicit.
- **Alt-C** is also comfortable for Bob, but it exports his preferred profile
  into the universal semantics of slot 0.

Finding: there is real value in telling IPLD-native implementations that
`42(pCID)` is not merely tolerated but preferred when they are in that profile.

### Scenario 3 - ATProto / Bluesky ecosystem interop

Carol wants PromiseGrid artifacts to feel legible and interoperable inside a
broader community that already treats CID-bearing DAG-CBOR structures as normal
engineering practice. She is not asking PromiseGrid to become ATProto; she is
asking PromiseGrid not to isolate itself from the surrounding CID/DAG-CBOR
world unnecessarily.

- **Alt-A** preserves conceptual purity, but it undersignals this interop goal.
  A reader from the broader CID/DAG-CBOR ecosystem may see PromiseGrid as
  tolerating their conventions rather than actively aligning with them.
- **Alt-B** says something more useful: PromiseGrid's semantic contract remains
  `pCID`, and when a DAG-CBOR/IPLD profile is in scope, the preferred
  realization is `42(pCID)`. This aligns with the surrounding ecosystem without
  collapsing role into representation.
- **Alt-C** maximizes that interop signal, but at the cost of making a profile
  convention look like the timeless meaning of slot 0.

Finding: the ATProto / Bluesky style of ecosystem spread is a real argument for
giving `42(pCID)` more design weight than `TE-titaj` gave it.

### Scenario 4 - Byte-level bootstrap cost

Dave is doing forensic work on old messages. He compares two common wire forms:

- raw CID in a CBOR byte string;
- DAG-CBOR Link tag `42` carrying `00 || CID`.

The difference is small: the CID bytes are the same payload either way, and the
Link form adds only a short wrapper. The receiver distinction is not "full IPLD
stack versus no parsing at all." The receiver distinction is:

- parse CBOR `bstr`, then CID; or
- parse tag `42`, then `bstr`, then `00`, then CID.

- **Alt-A** makes this small delta a reason not to prefer the Link form.
- **Alt-B** treats the delta as small enough that ecosystem interop can justify
  preferring the Link form in DAG-CBOR/IPLD profiles.
- **Alt-C** treats the delta as so small that the Link form should become the
  universal standard.

Finding: the practical byte-level penalty is small enough to weaken the
anti-`42(pCID)` argument, but not so small that it forces universalization.

### Scenario 5 - IPLD survives for a century

Ellen lives in the future branch where CID, DAG-CBOR, and IPLD remain healthy
and broadly used.

- **Alt-A** still works, but it leaves interop value under-described.
- **Alt-B** looks robust: the semantic role remains durable while the preferred
  profile representation aligns with the dominant ecosystem.
- **Alt-C** also works, and in this branch it may look attractive. But its
  success depends on the continued health of one representation ecosystem.

Finding: if IPLD remains healthy, `42(pCID)` should be more than an afterthought
in PromiseGrid guidance.

### Scenario 6 - IPLD fades but CID-like identifiers remain understandable

Frank lives in the future branch where IPLD as a full ecosystem faded, but CID
history is still recoverable from local docs and source archives.

- **Alt-A** is safe here because it never overcommitted to the wrapper.
- **Alt-B** still works: the semantic standard remains role-first, and the
  preferred profile guidance for DAG-CBOR/IPLD deployments remains intelligible
  as a profile note rather than as a universal law.
- **Alt-C** now looks more brittle because the historical wrapper convention is
  over-promoted into the base standard.

Finding: profile preference ages better than universal wrapper commitment.

### Scenario 7 - IPLD fades and the wrapper becomes historical boilerplate

Grace is reading century-old PromiseGrid packets in a world where DAG-CBOR is
no longer a living daily tool. She can still parse the packet because the
wrapper around the CID is short and stable enough to be treated as archival
boilerplate.

- **Alt-A** offers no special help here; it simply does not emphasize the
  wrapper.
- **Alt-B** still behaves well because the wrapper was only ever preferred
  within one profile family. Grace can treat it as recoverable historic
  boilerplate while keeping slot 0's semantic meaning clear.
- **Alt-C** asks Grace to treat that historic wrapper as the standard meaning of
  slot 0, which is a stronger and less necessary burden.

Finding: "the extra bytes can become durable boilerplate" is a good argument
against overstating the risk of `42(pCID)`, but not a good argument for making
it the universal semantic standard.

## Cross-cutting findings

### `TE-titaj` understated the interop upside

`TE-titaj` was right to separate semantic role from representation, but it
weighted the ecosystem cost of `42(pCID)` more heavily than the ecosystem value
of aligning with CID / DAG-CBOR communities that PromiseGrid wants to interwork
with.

### The wrapper burden is real but small

The extra tag-and-sentinel machinery in `42(pCID)` is a real bootstrap burden,
but it is not a demand for a full IPLD object stack. For many receivers it is a
small decoder extension rather than a system-wide commitment.

### Universalization still goes too far

Even after rebalancing, the case for making `42(pCID)` the universal timeless
meaning of slot 0 remains weaker than the case for keeping `pCID` as the
semantic role and treating `42(pCID)` as the preferred representation in the
appropriate profile family.

## Conclusions

- Reject the claim that `42(pCID)` is merely too heavy to be preferred anywhere.
  That overstates the burden.
- Reject the claim that `42(pCID)` should become the universal semantic meaning
  of slot 0. That still confuses representation with role.
- **Alt-A** survives only as a conservative baseline. It is no longer the best
  fit for PromiseGrid's stated interop ambitions.
- **Alt-B** is now the strongest surviving alternative:
  - keep slot 0 semantically as the Protocol CID;
  - strongly prefer `42(pCID)` when the carrier profile is DAG-CBOR / IPLD.
- **Alt-C** remains too strong. It makes one ecosystem's representation look
  universal across all profiles and all futures.

## Recommendation

Recommend **Alt-B**:

- PromiseGrid standard prose should continue to say that slot 0 means the
  Protocol CID.
- PromiseGrid should also say explicitly that when a DAG-CBOR / IPLD profile is
  used, `42(pCID)` is the strongly preferred slot-0 representation.
- PromiseGrid should not yet freeze `42(pCID)` as the universal representation
  across every profile and future carrier.

This keeps the semantic role clear, aligns with the surrounding CID /
DAG-CBOR/IPLD ecosystem that PromiseGrid wants to interoperate with, and treats
the extra wrapper cost as small enough to be worthwhile in the profiles that
benefit from it.

## Implications for open work

- `TE-titaj` is superseded for current slot-0 framing and recommendation.
- `TE-vujaj` remains the wording-focused facet for ordinary repo prose.
- Future DF should ask whether the repo wants to lock Alt-B now:
  - semantic `pCID` standard plus strongly preferred `42(pCID)` in DAG-CBOR /
    IPLD profiles; or
  - keep that profile preference advisory only.
- If a future DF ever revisits universal wrapper commitment, it should do so as
  a separate question with a stronger interop case than is currently needed.

## Decision status

`superseded by TE-fikoj / DI-kafat` - `TE-fikoj` accepts the stronger premise
that slot 0 itself may be standardized as `42(pCID)`, then moves the live
question to fixed three-slot versus variable-arity outer arrays.
