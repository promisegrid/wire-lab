# TE-fikoj: universal `42(pCID)` envelope shape

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-fikoj

## Status

decided

## Decision under test

`TE-tilir` concluded that PromiseGrid should keep `pCID` as the semantic role of
slot 0 and treat `42(pCID)` as the strongly preferred representation when the
carrier is DAG-CBOR / IPLD. The user has now pushed a stronger claim:

- the extra wrapper cost around `42(pCID)` is small enough to treat as durable
  boilerplate;
- small receivers do not need full IPLD, only a tiny tag-42 / byte-string /
  `00||CID` decoder;
- PromiseGrid wants to interoperate natively with the broader CID / DAG-CBOR /
  IPLD ecosystem;
- therefore PromiseGrid should seriously consider standardizing slot 0 itself as
  `42(pCID)` rather than leaving that choice at the profile-preference layer.

Given that stronger premise, this TE asks:

- should PromiseGrid standardize the outer envelope as a CBOR array whose slot 0
  is always `42(pCID)`; and
- if so, should the outer array land as fixed three-slot
  `grid([42(pCID), payload, sig])`, or as variable-arity
  `grid([42(pCID), payload, ...])` where CBOR array length carries arity and the
  protocol named by `pCID` defines the roles of slots beyond slot 0?

This TE supersedes `TE-tilir`'s framing and recommendation. It does not reopen
the older "payload CID" error, and it does not reopen `TE-pokul`'s broader
envelope-lineage history except where that TE's anti-unnecessary-outer-slot
pressure bears directly on the fixed-vs-variable choice.

## Assumptions

- A `pCID` is a Protocol CID: the CID of a protocol spec document, never the CID
  of a payload object.
- The protocol named by a `pCID` defines payload shape, signable view,
  proof/signature encoding, verification rules, and any slot-role rules that
  apply beyond the slot-0 protocol selector.
- The entire message is a CBOR array, so array length already carries arity.
- Promise Theory remains the semantic layer: the envelope carries evidence that
  a receiver uses to interpret another agent's promise. The envelope does not
  impose behavior or create non-local trust.
- `42(...)` is DAG-CBOR / IPLD Link notation. Freezing `42(pCID)` universally in
  slot 0 would standardize one outer-envelope wire form, but it would not change
  the fact that the semantic role of slot 0 is "Protocol CID."
- A small receiver can parse `42(pCID)` without a full IPLD object model: it
  only needs CBOR, tag `42`, the following byte-string, the `00` sentinel, and
  then CID parsing.
- If IPLD fades, the extra tag-and-sentinel bytes may still survive as stable
  historical boilerplate that later readers can reconstruct locally.
- The design question here is not whether `42(pCID)` is acceptable at all. This
  TE starts from the stronger premise that PromiseGrid is willing to standardize
  it in slot 0 if the outer shape choice makes sense.

## Alternatives

- **Alt-A - Fixed three-slot envelope.** Standardize the outer shape as
  `grid([42(pCID), payload, sig])`. Slot 0 is always the Link-wrapped Protocol
  CID, slot 1 is always the payload, and slot 2 is always the outer signature or
  proof slot.
- **Alt-B - Variable-arity outer array.** Standardize the outer shape as
  `grid([42(pCID), payload, ...])`. Slot 0 is always `42(pCID)`, slot 1 is the
  main payload anchor, CBOR array length carries arity, and the protocol named
  by `pCID` defines the roles of additional outer slots.

## Scenario analysis

### Scenario 1 - Alice on a constrained receiver

Alice runs a tiny receiver on a low-cost controller. She must parse the envelope
with a small CBOR decoder and recover the protocol selector without depending on
live network services.

- **Alt-A** gives Alice a fixed outer contract. Once she recognizes tag `42`
  and recovers the CID, she always knows slot 1 is payload and slot 2 is the
  outer proof slot.
- **Alt-B** asks Alice to learn one more rule: after slot 0 and slot 1, she
  must consult the protocol named by `pCID` to know whether more outer slots
  exist and what they mean.

Finding: fixed three-slot is simpler for the smallest generic receiver, but the
incremental complexity of Alt-B is not extreme if Alice already needs to consult
the protocol spec for payload and proof semantics anyway.

### Scenario 2 - Bob in a rich IPLD-native implementation

Bob runs a CID / DAG-CBOR-native implementation and wants PromiseGrid messages
to feel natural in that ecosystem.

- **Alt-A** is comfortable for Bob. It gives him one fixed outer proof slot that
  generic tools can expect.
- **Alt-B** is also comfortable because CBOR arrays with stable slot 0 and
  protocol-defined later slots are ordinary structured objects for his runtime.

Finding: rich IPLD-native stacks do not strongly prefer fixed three-slot over
variable arity. Both interoperate well once `42(pCID)` is standardized.

### Scenario 3 - ATProto / Bluesky ecosystem interop

Carol wants PromiseGrid to interoperate gracefully with a wider ecosystem that
already treats CID-bearing DAG-CBOR structures as normal engineering practice.

- **Alt-A** advertises a narrower universal contract: slot 0, payload, signature.
  This can make PromiseGrid easier to explain to adjacent ecosystems that expect
  a concise fixed object skeleton.
- **Alt-B** advertises a more protocol-owned object family: slot 0 stays fixed,
  but additional outer slots may appear when the protocol named by `pCID` says
  they should.

Finding: both interoperate. Alt-A is easier to summarize in one sentence; Alt-B
fits better with the broader content-addressed habit of letting the schema own
object shape.

### Scenario 4 - Future envelope evolution without outer-shape churn

Dave wants the outer envelope to survive new proof families, selector needs,
summary slots, or migration hints without forcing PromiseGrid to redesign the
base envelope again.

- **Alt-A** resists outer-shape churn by pushing almost everything into the
  payload or into the fixed outer signature slot. That is clean if the fixed
  proof slot is enough.
- **Alt-B** uses CBOR arity as the escape hatch PromiseGrid already has. If a
  future protocol truly needs an extra outer slot, it can add one without
  pretending the entire wire family has changed.

Finding: once slot 0 is frozen as `42(pCID)`, CBOR array arity becomes a strong
argument for Alt-B. It preserves future maneuvering room without inventing a new
escape hatch.

### Scenario 5 - Generic relays and archival readers

Ellen is a relay or archival peer. She may not fully understand the protocol
named by `pCID`, but she wants a durable generic understanding of what the
envelope is carrying.

- **Alt-A** helps Ellen because the outer contract is globally fixed: protocol
  selector, payload, signature. Even partial readers know where the proof lives.
- **Alt-B** gives Ellen a stable slot 0 and slot 1, but she may need protocol
  knowledge before interpreting later outer slots.

Finding: generic partial understanding favors Alt-A. But Ellen already cannot
interpret payload semantics without protocol knowledge, so Alt-B's extra burden
is not entirely new.

### Scenario 6 - IPLD fades but the wrapper survives as boilerplate

Frank lives in a future where DAG-CBOR is no longer fashionable, but archived
software and notes still make tag `42` plus `00||CID` recoverable.

- **Alt-A** keeps the rest of the outer contract simpler once the historical
  wrapper is stripped or reconstructed.
- **Alt-B** treats the wrapper as boilerplate and also treats any later outer
  slots as protocol-owned structure that can be rediscovered from the spec named
  by `pCID`.

Finding: if the wrapper is acceptable as archival boilerplate, Alt-B also
becomes acceptable as archival boilerplate. The difference shifts from
bootstrap-burden to how much outer structure PromiseGrid wants globally fixed.

### Scenario 7 - Signature slot pressure

Grace argues that a fixed outer signature slot makes one thing universally
legible: there is always an outer sender-proof location.

- **Alt-A** satisfies Grace directly.
- **Alt-B** can still satisfy Grace if most protocols define slot 2 as the outer
  proof slot, but it does not force that forever. The cost is weaker universal
  expectations; the benefit is that PromiseGrid does not freeze one proof-slot
  story too early if some protocols later need different outer structure.

Finding: fixed three-slot buys immediate regularity. Variable arity buys
protocol-owned freedom at the cost of that regularity.

### Scenario 8 - Reserved space and future selectors

Henry argues that the `42(...)` wrapper itself is already a tiny layer of
boilerplate and could act like durable reserved space if future selector
machinery ever needs to grow nearby.

- **Alt-A** says that reserved-space thinking is not enough reason to freeze
  exactly three slots globally. If a protocol needs more than `[slot0, payload,
  sig]`, it should use payload structure or later reopen the outer-envelope
  question.
- **Alt-B** says CBOR arity plus protocol-defined slot roles already gives a
  cleaner reserved-space story: keep slot 0 fixed, keep the array extensible,
  let the protocol define later slots when truly needed.

Finding: if the goal is "durable reserved space," Alt-B expresses that more
directly than freezing exactly three slots now.

## Cross-cutting findings

### Universal `42(pCID)` is now viable

Given the interop goal and the small byte-level overhead, the stronger premise
behind this TE is coherent: PromiseGrid can plausibly standardize slot 0 itself
as `42(pCID)` without demanding a full IPLD object stack from every receiver.

### Fixed three-slot regularity is real

Alt-A's strongest virtue is universal legibility. Any partial reader knows where
the protocol selector, payload, and outer proof live.

### CBOR arity plus protocol-owned slot roles is also real

Alt-B's strongest virtue is that PromiseGrid already has an outer extensibility
mechanism: CBOR array length. Once slot 0 is fixed and the protocol spec owns
the meaning of later slots, PromiseGrid can avoid repeated redesign of the base
shape.

### The fixed-vs-variable question is no longer about future-proofing bytes

The hard part is no longer whether `42(pCID)` is too heavy. The hard part is
whether PromiseGrid wants more **global universal regularity** or more
**protocol-owned outer evolution**.

## Conclusions

- Reject any return to the idea that `42(pCID)` is too costly to standardize in
  slot 0. Under current assumptions, it is acceptable.
- Keep `pCID` as the semantic role of slot 0 even if PromiseGrid freezes the
  wire form as `42(pCID)`.
- **Alt-A** remains plausible because fixed three-slot buys global regularity and
  easier generic partial interpretation.
- **Alt-B** is the stronger long-horizon fit:
  - it keeps slot 0 universally fixed as `42(pCID)`;
  - it uses CBOR array length as the built-in extensibility mechanism;
  - it lets the protocol named by `pCID` define later outer-slot roles instead
    of forcing PromiseGrid to decide every future outer slot today.

## Recommendation

Recommend **Alt-B**, with a narrower reading than "anything goes":

- standardize the outer envelope as a CBOR array whose slot 0 is always
  `42(pCID)`;
- keep slot 1 as the main payload anchor;
- allow additional outer slots when the protocol named by `pCID` defines them;
- do not freeze a universal exactly-three-slot rule unless a later TE shows that
  the benefits of one globally fixed proof slot outweigh the value of CBOR-array
  extensibility.

In short: if PromiseGrid is going to universalize `42(pCID)` at all, it should
probably also lean into CBOR-array arity and pCID-owned slot roles, rather than
freezing the whole family at exactly three slots too early.

## Implications for open work

- `TE-tilir` becomes the historical predecessor for the interop and byte-cost
  argument.
- `TE-vujaj` remains the wording/history facet for semantic role versus
  representation language.
- `TE-pokul` remains relevant as the older anti-unnecessary-outer-slot pressure,
  but its likely landing-shape analysis should now be read through this stronger
  universal-`42(pCID)` premise.
- Future DF should ask whether PromiseGrid now wants to lock:
  - universal slot-0 `42(pCID)` plus variable-arity outer arrays with
    pCID-defined slot roles; or
  - preserve fixed three-slot as the universal outer contract despite this TE's
    recommendation.

## Decision status

Locked by `DI-sisak` - the current PromiseGrid direction is the variable-arity
CBOR outer envelope `grid([42(pCID), payload, ...])`. Slot 0 belongs to a
tagged CBOR protocol-selector family, with tag `42` as the current standard
instance; CBOR array length carries arity; and the protocol named by `pCID`
defines later outer-slot roles. Fixed three-slot
`grid([42(pCID), payload, sig])` remains considered but unchosen.
