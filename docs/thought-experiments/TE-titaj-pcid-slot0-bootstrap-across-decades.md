# TE-titaj: pCID slot-0 bootstrap across decades

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-titaj

## Status

superseded by TE-tilir / DI-lasah

## Decision under test

What slot-0 standard should PromiseGrid choose now if the deciding criterion
is not present-day elegance but the widest range of future receivers that can
still parse a message, recover the protocol selector, and interpret the rest of
the evidence-bearing promise across decades, devices, runtimes, and identifier
ecosystem drift?

This TE compares three alternatives for a three-slot envelope candidate:

- `grid([42(pCID), payload, sig])`
- `grid([pCID, payload, sig])`
- `grid([slot0_bytes, payload, sig])`, where the PromiseGrid envelope standard
  defines how a receiver deterministically recovers `pCID` from `slot0_bytes`
  before consulting the protocol spec named by that `pCID`

This TE does not reopen `TE-pokul`'s separate family question about whether the
long-term outer envelope should land at two slots or three. It is strictly
about slot-0 meaning, representation, and bootstrap.

## Assumptions

- A `pCID` is a Protocol CID: the CID of a protocol spec document, never the
  CID of a payload object.
- The protocol named by `pCID` defines payload shape, signable view,
  proof/signature encoding, verification rules, and related semantics.
- The whole message is CBOR. Array arity is carried by the CBOR array header.
- Promise Theory remains the semantic layer: slot 0 helps a receiver interpret
  another agent's evidence-bearing promise; it does not impose behavior or
  create non-local trust.
- All trust judgments remain local. A receiver can only decide, from its own
  evidence and history, whether the sender's promise is intelligible,
  believable, useful, or worth relaying.
- A receiver must recover the `pCID` before it can consult the protocol spec
  named by that `pCID`. Therefore the protocol spec cannot be the thing that
  bootstraps slot 0.
- `42(...)` is DAG-CBOR/IPLD Link notation. It is one representation of a CID,
  not the meaning of the CID.
- CID and multihash standards may remain healthy, may fragment, or may become
  historical formats that later receivers must reconstruct from archived text
  and source code.

## Alternatives

- **Alt-A - Semantic `pCID` standard.** PromiseGrid standard prose says the
  shape is `grid([pCID, payload, sig])`. Carrier profiles may render `pCID`
  differently, but the standard is role-first: slot 0 is the Protocol CID.
- **Alt-B - DAG-CBOR Link standard.** PromiseGrid standard prose says the shape
  is `grid([42(pCID), payload, sig])`. DAG-CBOR/IPLD Link representation is
  treated as the normal slot-0 form, not merely one profile-local rendering.
- **Alt-C - Envelope-bootstrap bytes standard.** PromiseGrid standard prose
  says the shape is `grid([slot0_bytes, payload, sig])`. The PromiseGrid
  envelope standard defines a small, durable, deterministic rule for recovering
  the `pCID` from slot-0 bytes. The protocol spec named by that `pCID` governs
  everything after bootstrap.

## Scenario analysis

### Scenario 1 - Alice on a constrained device today

Alice runs a tiny receiver on a makerspace door controller. She has a small
CBOR parser, limited flash, and no appetite for carrying a full IPLD stack.
She still needs to extract slot 0, decide whether the message names a protocol
she knows, and then decide whether to store, forward, reject, or quarantine the
message as local evidence.

- **Alt-A** is conceptually clear to Alice, but by itself it does not tell her
  exactly what byte-level forms she must accept in slot 0. She still needs a
  profile or envelope rule for concrete bootstrap.
- **Alt-B** gives Alice exactness, but only by requiring DAG-CBOR Link
  knowledge in her smallest receiver. That is a stronger requirement than
  "understand slot 0 as a protocol selector."
- **Alt-C** can work well for Alice if the envelope bootstrap rule is tiny and
  stable. It fails badly if `slot0_bytes` means "whatever future profiles feel
  like doing."

Receiver lesson: the smallest receivers want a minimal deterministic bootstrap
rule, not a large profile ecosystem.

### Scenario 2 - Bob on a rich server runtime in twenty years

Bob runs a large relay on a server with full CBOR, DAG-CBOR, CID, multihash,
and archival tooling. He wants exact-byte fidelity, good debugging, and
interoperability with content-addressed toolchains.

- **Alt-A** is easy for Bob to live with; he can support many carrier profiles
  while still treating slot 0 semantically as the Protocol CID.
- **Alt-B** is comfortable for Bob because his tooling already understands
  `42(...)`. The danger is that Bob's comfort gets mistaken for a universal
  assumption about all receivers.
- **Alt-C** also works for Bob if the bootstrap rule is explicit. Rich runtimes
  do not object to a smaller shared rule; they simply have more conveniences on
  top.

Receiver lesson: rich receivers can tolerate all three; they are not the right
baseline for a 100-year interoperability decision.

### Scenario 3 - Carol decades later with partial ecosystem continuity

Carol arrives in fifty years with a new language runtime. She has CBOR support,
some archived PromiseGrid specs, and perhaps partial CID libraries. She does not
know whether DAG-CBOR/IPLD is still culturally central.

- **Alt-A** helps Carol remember the important invariant: slot 0 names the
  protocol spec. But she still needs a concrete decoding rule from elsewhere.
- **Alt-B** makes Carol first recover a DAG-CBOR/IPLD convention before she can
  recover the `pCID`. If that ecosystem has faded, the notation itself becomes a
  historical dependency.
- **Alt-C** gives Carol the best chance if the bootstrap rule is archived in a
  tiny, stable PromiseGrid envelope rule that says how slot-0 bytes become a
  `pCID` without assuming a large live ecosystem.

Receiver lesson: long-horizon readers need the semantic role made explicit and
the bootstrap rule made small.

### Scenario 4 - Dave during an IPLD-dominant century

Dave lives in a future where IPLD survives and becomes dominant because it keeps
good content-addressed tooling, selectors, archival software, and broad
interchange support. Most receivers already know DAG-CBOR links.

- **Alt-A** still works; profile-local prose can say that the common
  representation under this future is DAG-CBOR `42(pCID)`.
- **Alt-B** feels natural in this world, and the cost of standardizing the
  notation looks small.
- **Alt-C** still works, but it may feel redundant if the envelope bootstrap
  rule collapses to "slot 0 carries a CID in a form every receiver already
  knows."

Receiver lesson: IPLD dominance makes Alt-B look attractive, but only under one
future branch.

### Scenario 5 - Ellen after IPLD fades but CID/multihash remains legible

Ellen lives in a future where IPLD faded the way XML faded: too heavy for the
common case, culturally displaced by simpler tools, but not completely erased.
CID and multihash encodings still exist in archives and some libraries.

- **Alt-A** ages well because it did not tie the meaning of slot 0 to one
  profile-local notation.
- **Alt-B** now forces Ellen to recover a notation family that is no longer the
  common cultural default, even though the underlying CID concept still lives.
- **Alt-C** still works if the PromiseGrid envelope bootstrap rule says
  something small like "slot 0 bytes carry a canonical CID byte form" or a
  similarly tiny closed set of allowed forms.

Receiver lesson: when ecosystems fade unevenly, the notation-heavy standard
ages worse than the role-first or small-bootstrap alternatives.

### Scenario 6 - Frank after CID and multihash themselves fade or fork

Frank lives in a harder future. CID and multihash are no longer universal.
Several successor identifier schemes exist. PromiseGrid messages survive only
because old specs, source trees, and transcripts were archived carefully.

- **Alt-A** still preserves the key human-facing statement: slot 0 is the
  Protocol CID. But it does not by itself tell Frank how a byte string becomes
  that identifier.
- **Alt-B** now burdens Frank with two layers of archaeology: first CID, then
  DAG-CBOR/IPLD link syntax around CID.
- **Alt-C** is strongest in this scenario if and only if PromiseGrid had frozen
  a very small envelope bootstrap rule that Frank can reconstruct from local
  archives. Without that, Alt-C collapses into "opaque bytes, good luck."

Receiver lesson: the harsher the historical break, the more important it is
that PromiseGrid itself own a tiny bootstrap rule instead of inheriting one
implicitly from a large external ecosystem.

### Scenario 7 - Forensic replay and exact-byte debugging

Grace receives a disputed message years later. She must show exactly what slot 0
contained on the wire, what she decoded from it, and why she trusted or
rejected the rest.

- **Alt-A** is best for explaining meaning but needs a companion profile rule
  for exact bytes.
- **Alt-B** is precise for one profile family, but risks making the forensic
  story look universal when it is really profile-specific.
- **Alt-C** can be excellent for forensics if the envelope bootstrap rule is
  deterministic and compact. Then Grace can say: these slot-0 bytes decode to
  this `pCID`; this `pCID` names this protocol spec; this is how I interpreted
  the rest.

Receiver lesson: forensics benefits from a clear separation between slot-0
 bytes, bootstrap rule, and protocol semantics.

### Scenario 8 - Mixed receivers on one network

Alice's small device, Bob's large relay, Carol's future runtime, and Frank's
archival reader all receive the same family of messages. PromiseGrid wants the
same message to remain understandable, even if not everyone supports the same
tooling conveniences.

- **Alt-A** gives the cleanest shared language, but needs a separate answer for
  concrete bootstrap.
- **Alt-B** makes the richest receiver's favorite notation the universal
  requirement.
- **Alt-C** has the best chance of giving every receiver the same minimum parse
  contract, provided the bootstrap rule remains narrow and durable.

Receiver lesson: the common denominator is not "whatever the richest receiver
likes best." It is "the smallest stable rule every receiver can still
reconstruct."

## Cross-cutting findings

### The bootstrap problem is load-bearing

Any slot-0 standard that says "the protocol spec defines how to interpret slot
0" is circular. The receiver needs the `pCID` before it can find that protocol
spec. Therefore the bootstrap rule for slot 0 must live either in the envelope
standard or in a carrier profile that the envelope standard explicitly adopts.

### `42(...)` is a representation, not the role

The more time passes, the more dangerous it is to confuse a particular wire
notation with the semantic role "slot 0 names the protocol spec." `42(pCID)` is
one way to carry a CID. It is not the universal meaning of slot 0.

### The lowest common denominator is small

Across all scenarios, the minimum credible future receiver can be expected to:

- parse a CBOR array;
- know that slot 0 is the protocol selector;
- apply a small deterministic bootstrap rule that recovers a `pCID` from slot-0
  bytes;
- then consult archived or local knowledge about the protocol named by that
  `pCID`.

The minimum credible future receiver cannot safely be assumed to carry a full
IPLD ecosystem, a live registry service, or culturally current knowledge of
every historic CID wrapper convention.

### Role-first prose and bootstrap-first wire rules are different jobs

This TE keeps two distinctions sharp:

- a **semantic statement** for humans and design prose: what slot 0 means;
- a **bootstrap rule** for receivers: how slot-0 bytes become a `pCID`.

PromiseGrid may need both. Conflating them makes the standard harder to read and
harder to preserve.

## Conclusions

- Reject any standard that depends on the protocol spec to bootstrap slot 0.
  That is circular.
- Reject any standard that silently treats `42(pCID)` as the timeless meaning of
  slot 0. That bakes one ecosystem's notation into the semantic contract.
- **Alt-B** is the weakest long-horizon choice. It works well if IPLD dominates
  for a century, but it ages poorly if DAG-CBOR/IPLD stops being universal.
- **Alt-A** survives as the best statement of semantic meaning: slot 0 is the
  Protocol CID.
- **Alt-C** survives as the strongest candidate for exact receiver bootstrap,
  but only in a narrowed form: PromiseGrid itself must define a very small,
  explicit, durable rule for turning slot-0 bytes into a `pCID`.
- Therefore the live question is not "Alt-A or Alt-C as pure rivals?" The live
  question is whether PromiseGrid should freeze:
  - only the semantic role in the standard and leave concrete byte bootstrap to
    profiles; or
  - both the semantic role and one tiny envelope-level bootstrap rule.

## Recommendation

Recommend the following interpretation of the current evidence:

- use **Alt-A** for conceptual PromiseGrid prose and cross-profile reasoning:
  `grid([pCID, payload, sig])`;
- reject **Alt-B** as the general standard wording;
- carry **Alt-C** forward only in its narrowed form as the candidate for an
  eventual exact bootstrap rule, not as "arbitrary bytes."

In plain terms: PromiseGrid should say now that slot 0 means the Protocol CID,
and it should remain open whether the envelope standard later freezes a tiny
byte-level bootstrap rule for that `pCID`. If such a rule is frozen, it should
be justified as the lowest-common-denominator receiver contract, not as a
special devotion to DAG-CBOR.

## Implications for open work

- `TE-vujaj` remains the wording-focused facet: whether repo prose should say
  `pCID` or `42(pCID)` in ordinary discussions.
- `TE-pokul` remains the owner of the separate two-slot vs three-slot family
  question.
- Future DF for this TE should ask whether PromiseGrid wants to standardize only
  the semantic role of slot 0 now, or also standardize a tiny envelope-level
  bootstrap rule for exact slot-0 bytes.
- Any future bootstrap rule should be written as a small archival contract that
  constrained devices, relays, and long-horizon forensic readers can all
  reconstruct locally.

## Decision status

`superseded by TE-tilir / DI-lasah` - `TE-tilir` adds the ATProto / Bluesky
ecosystem-interop scenario and the small-byte-delta bootstrap analysis, then
rebalances the recommendation toward semantic `pCID` plus strongly preferred
`42(pCID)` in DAG-CBOR / IPLD profiles.
