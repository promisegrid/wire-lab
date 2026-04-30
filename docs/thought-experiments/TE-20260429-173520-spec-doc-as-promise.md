# TE-21: Spec doc as promise

*Thought experiment, part of the [PromiseGrid Wire Lab](../../specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-20260429-173520

(First drafted 2026-04-29 17:35:20 UTC.)

## Decision under test

In promise-theoretic terms, what *is* a wire-protocol spec document? Specifically: what kind of speech act does a spec doc perform, what is its promiser, what is its promisee, what is the temporal shape of the promise it encodes, and what are the implications for which content belongs in the spec versus outside it?

This question is upstream of every other vocabulary and pCID question in this repo. If `pCID` is the content hash of "a spec document that defines a wire protocol" (per `specs/harness-spec-draft.md` §1's pCID note), then the nature of *spec document* is the load-bearing concept. Pinning down what a spec doc is determines:

- What content belongs inside the hashed file (and therefore changes the pCID when it changes).
- What content belongs outside in a separate workshop/archaeology artifact (DR, TE, TODO).
- Whether a spec doc with a §11-shaped open-questions section is "complete" or "incomplete" as a promise.
- Whether `specs/harness-spec-draft.md` should be promoted into a sibling-protocol doc store with one pCID, or remain a single workshop doc that carries an embedded protocol.

## Assumptions

- A `Promise` (in this repo's vocabulary, per `specs/harness-spec-draft.md` §1.1 and the `assertion` rename) is an autonomous speech act by an agent: an assertion of state in the past, present, or future, often conditional. The cleanest examples in this repo: "I will deliver this byte stream"; "I attest authorship of this content"; "I promise that I, Steve, will accept revised work after named conditions are met."
- A `pCID` is the content hash of a spec document. Two parties claim to "speak protocol pCID X" when they each implement the rules in the document whose hash is X.
- The Wire Lab is a long-lived design environment, not a one-shot RFC. Open questions, assumptions, and known weaknesses persist throughout the life of any spec because no spec is ever truly finished — only frozen at a hash.
- Steve's framing, surfaced in chat 2026-04-29: "a protocol spec is a prediction (an assertion of future state, a promise that the protocol will work as designed). Not listing assumptions and known issues is an incomplete promise."
- This TE is a vocabulary/philosophy artifact. It does not propose moving `specs/harness-spec-draft.md` to a different directory or splitting it into multiple files. Layout questions are deferred to a follow-on TE.

## Alternatives

Five candidate framings of what a spec doc *is* in promise-theoretic terms.

### Alt-A: Spec doc is just reference content (not a promise at all)

Under this framing, a spec doc is a description — like a dictionary or an encyclopedia entry. It says "here is how the bytes look" but makes no commitment about anything. Implementers can choose to follow it or not; the doc does not bind anyone.

- **Promiser:** none; the doc is descriptive, not performative.
- **Promisee:** none.
- **Temporal shape:** present-tense description.
- **Open-questions section:** unproblematic; a description can be incomplete without losing its descriptive nature.
- **Easier:** matches how RFCs are sometimes read in practice ("informational" RFCs especially).
- **Harder:** fails to explain why parties care which pCID they share. If a spec is just description, two parties using the same pCID have no shared *commitment*, only shared *vocabulary*. But the whole point of pCID-as-port-number is to identify a shared protocol *that both parties have committed to running*. Mere shared vocabulary doesn't give you that.

### Alt-B: Spec doc is a normative declaration by the spec's author

Under this framing, the spec doc is a promise *by its author* that "this is the protocol; if you implement it as written, I declare you to be speaking it." It's a unilateral act of definition.

- **Promiser:** the spec's author.
- **Promisee:** anyone who reads the spec and chooses to implement it.
- **Temporal shape:** present-tense definition with implied future enforcement ("this is what counts as conformance").
- **Open-questions section:** awkward; the author cannot normatively declare "I don't know what should happen here."
- **Easier:** makes the spec author the arbiter, which is how IETF process actually works in practice (the working group authors are the arbiters of what's in the RFC).
- **Harder:** bakes in centralization. PromiseGrid's whole posture is that pCIDs are minted by anyone who writes a spec, not granted by an authority. Alt-B keeps the authoritative-author shape of RFC-land.

### Alt-C: Spec doc is a promise by each implementer (one promise per peer)

Under this framing, the spec doc is a *template* for promises: when peer P claims to "speak pCID X," P is making a promise of the form "I, P, promise to behave as the document with hash X says I will." The spec doc itself is not a promise; only the implementer's adoption of it is.

- **Promiser:** each implementer, separately.
- **Promisee:** every other peer P talks to.
- **Temporal shape:** future commitment by each implementer ("I will behave as ...").
- **Open-questions section:** unproblematic for the spec doc, but each implementer's promise is conditional on the answers they assume to those questions, which creates silent divergence between peers who answered differently.
- **Easier:** clean Promise-Theory shape; each peer's promise is autonomous; matches the "promises are made by agents, not by documents" instinct.
- **Harder:** the spec doc itself loses normative weight; it becomes pure scaffolding for each peer's separate commitment. Makes the §11-style open-questions list part of the *peer's* private interpretation rather than part of the shared protocol. Two peers could disagree about answers to §11 questions and both legitimately claim to speak pCID X.

### Alt-D: Spec doc is a promise by the spec itself (the document as agent)

Under this framing, the document is a promise from the document to its readers: "I, this document, promise that systems implementing the rules I contain will interoperate with each other." The document is treated as an agent in its own right — content-addressed, immutable, signed implicitly by its hash.

- **Promiser:** the spec document, identified by its pCID.
- **Promisee:** any party reading or running the spec.
- **Temporal shape:** prediction about future state ("systems implementing me WILL interoperate"). The spec is fundamentally future-oriented; it is making an empirical claim about a class of future runs of a class of future implementations.
- **Open-questions section:** normative and required. A document promising "systems implementing me will interoperate" cannot honestly omit the assumptions that promise depends on, the known weak spots where the promise is shakiest, and the open questions whose answers would change the promise. Listing them is the document keeping its promise honest.
- **Easier:** matches Steve's chat framing exactly. Lets the doc be a single pCID, including its open-questions section, because the open-questions section is part of what the doc is promising about. Two peers using the same pCID share a single promise — the one made by the doc — and that promise includes its own list of caveats, so divergence is bounded.
- **Harder:** treating a document as an agent is a Promise-Theory unusual move. Standard Promise Theory has agents make promises *about* artifacts; agents are the speech-act bearers. But content-addressed immutable artifacts do have an agent-like property: they cannot lie about themselves (their content is their identity), and they cannot change their mind (they are immutable). So treating a content-addressed spec doc as an agent is at least defensible — and arguably more honest than Alt-C, where each peer privately reinterprets.

### Alt-E: Spec doc is a *layered* promise (Alt-D for the doc, Alt-C for adoption)

Under this framing, both promises happen:

- The spec document, as an agent, promises (Alt-D shape) that systems implementing it will interoperate, conditional on the assumptions and open questions it lists.
- Each peer that adopts the spec promises (Alt-C shape) "I will behave as the document with hash X says I will, accepting the doc's stated assumptions and acknowledging its stated open questions."

The two promises stack: doc-level "the protocol works as designed" plus peer-level "I commit to this interpretation."

- **Promiser:** two layers (the doc; the implementer).
- **Promisee:** two layers (the doc promises to readers; the implementer promises to other peers).
- **Temporal shape:** doc-level is prediction about future; peer-level is commitment about future.
- **Open-questions section:** part of the doc-level promise. Because the doc lists its assumptions, peer-level adoption can cite them ("I implement pCID X; I rely on assumption A1 staying true"). This makes silent divergence visible: a peer that violates A1 can be flagged.
- **Easier:** captures both the doc-as-agent insight and the implementer-makes-a-promise insight. Both feel right; both are described in different parts of `specs/harness-spec-draft.md`.
- **Harder:** more conceptual machinery; two layers means more vocabulary to track. Worth the cost only if the layered structure produces real downstream payoff.

## Scenarios

Each scenario tests how the alternatives behave under one stress.

### S1 (normal): two peers run the same protocol

Peer A and peer B both implement pCID X (they each have a copy of the document with that hash and they each behave as it says).

- **Alt-A:** Two peers with shared vocabulary; no shared commitment. Cannot explain why they expect to interoperate beyond hope.
- **Alt-B:** Both peers obey the spec's author's normative declaration. Interop guaranteed conditional on the author having defined enough.
- **Alt-C:** Each peer made its own promise to behave as X. Interop guaranteed if both promises are well-shaped.
- **Alt-D:** The doc promised they would interoperate. Both peers are inside the doc's promise. Interop guaranteed conditional on the doc's assumptions holding.
- **Alt-E:** Both Alt-C and Alt-D are in force. Interop is doubly guaranteed: the doc predicts it and each peer commits to it.

S1 verdict: A fails outright. B/C/D/E all work but explain interop differently. No discriminator yet.

### S2 (failure): an assumption silently doesn't hold

A and B share pCID X. The spec lists "assumes clocks are within ±5s of each other" as A1. A and B's clocks drift to 30s apart. They mis-interoperate on a TTL-bearing frame.

- **Alt-A:** Cannot diagnose. The spec doesn't promise anything; nothing is broken from the spec's POV.
- **Alt-B:** The author's declaration is silent on what to do when A1 is violated; the failure is a "your problem" rather than a spec failure.
- **Alt-C:** Each peer's promise was conditional on A1. Both peers can claim "I kept my promise; the failure is environmental." Nobody is wrong; nobody is responsible.
- **Alt-D:** The doc's promise was explicitly conditional on A1. A1 was violated, so the doc's promise is in force-majeure territory. Failure is attributable to A1 violation. The diagnostic path is clear: check the assumption list.
- **Alt-E:** Doc-level + peer-level. Same diagnosis as Alt-D plus the peers can re-reach: now that A1 is known to be violated, each peer can re-promise under a tightened assumption ("I now require ±10s") and the doc can be revised, producing a new pCID.

S2 verdict: D and E are the only framings that produce a clean diagnosis-and-recovery story. The assumption list, when normative, is the diagnostic anchor.

### S3 (concurrent): two peers disagree about an open question

Spec X has §11 with question Q7 unanswered. Peer A's implementation guesses Q7 = "yes"; peer B's guesses Q7 = "no". They both claim to speak pCID X. They diverge on a Q7-shaped exchange.

- **Alt-A:** No commitments; nothing to violate; nothing to diagnose.
- **Alt-B:** The author didn't say. A and B can both claim conformance.
- **Alt-C:** Each peer promised "to behave as X says." X doesn't say. Each peer's promise is honest but the promises are silently incompatible. No mechanism flags it.
- **Alt-D:** The doc explicitly listed Q7 as open. The doc's promise was conditional on Q7 being answered ("interop *conditional on the answer to Q7*"). A and B both implemented under the open-question caveat; their interop expectation was correctly downgraded by the doc's own promise.
- **Alt-E:** Same as D, plus: each peer's adoption promise can name *which* answer they're adopting ("I implement pCID X with Q7 = yes"). Peers can detect the mismatch by exchanging adoption metadata. This makes silent divergence into loud divergence.

S3 verdict: D handles the case correctly (the doc warned readers); E adds a recovery path (peers can converge by exchanging answer-choice metadata). The §11 list is doing real work in both.

### S4 (long-horizon): the spec is updated

Year 1: pCID X. Year 5: a question gets answered, an assumption gets tightened, and a new pCID Y is published. Peers using X want to upgrade to Y. Some don't upgrade.

- **Alt-A:** No promises; no upgrade story.
- **Alt-B:** The author publishes a new authoritative declaration. Peers that don't upgrade are "out of date" by the author's authority.
- **Alt-C:** Each peer's promise was tied to a specific pCID. New pCID, new peer-level promise. Old peers can keep their old promise indefinitely; new peers make a new one.
- **Alt-D:** Old doc and new doc are both valid promises; one stronger (more answered questions), one weaker. Peers can choose which to be inside.
- **Alt-E:** Doc-level evolution plus peer-level adoption. A peer can move from adopting X to adopting Y by issuing a new adoption promise. The old promise expires; the new one takes effect.

S4 verdict: D and E both handle multi-version evolution cleanly. C also handles it but loses the doc-level explanation of *why* Y is an improvement over X.

### S5 (trust boundary): an adversary publishes a spec

Adversary M publishes a doc with pCID Z that looks like X but with a hidden weakening. Peer A reads Z and "implements pCID Z" without checking that Z's promise is the same as X's.

- **Alt-A:** No way to compare promises (there are none); A is stuck reading bytes.
- **Alt-B:** A trusts the author. If M is a different author from X's author, A might notice. If M imitates X's author, A is fooled.
- **Alt-C:** A's adoption promise is "I behave as the doc with hash Z says." A is making a clean promise; whether Z is a good spec is orthogonal. A's only defense is to read Z carefully.
- **Alt-D:** Z is the doc's promise. A reads Z and sees what Z promises. A can compare Z's promise against X's promise (both are explicit, including assumptions and open questions) and notice the weakening. The doc-level promise IS the audit surface.
- **Alt-E:** Same as D plus A's adoption promise lets A name Z explicitly, which makes A accountable for the choice.

S5 verdict: D and E provide the strongest audit story because the doc's promise is itself the thing being compared. A provides nothing; B is fooled by author-imitation; C makes A's adoption clean but doesn't help A detect the weakening.

### S6 (scale): hundreds of sibling protocols coexist

Wire Lab grows ten sibling protocol docs (frame format, trust ledger, currency, eval rules, capability tokens, etc.), each with its own pCID. A peer that runs five of them has five separate adoption promises.

- **Alt-A:** Five separate descriptions. No way to compose.
- **Alt-B:** Five separate authoritative declarations. Composition is per-author.
- **Alt-C:** Five peer-level promises. Composition is by the peer's choice of which to adopt.
- **Alt-D:** Five doc-level promises. Each says what its assumptions are. A composing peer can read all five assumption lists and check for compatibility.
- **Alt-E:** Five doc-level promises plus five peer-level adoption promises. The peer's adoption can cite *which combination of pCIDs* it implements, which is exactly the IETF "Updates: 2119" / "Implements: RFCXYZ" idea but with hash-pinned identity.

S6 verdict: D and E scale cleanly because each doc carries its own assumption list and the lists can be machine-checked for compatibility. C scales too but the assumption-tracking has to live somewhere outside the spec, which loses the audit anchor.

## Conclusions

The dominant alternative is **Alt-E (layered: doc-level Alt-D plus peer-level Alt-C)**. Reasons:

1. Alt-E explains *both* the doc-as-agent insight (the spec is making a prediction about future runs of future implementations) and the implementer-makes-a-promise insight (each peer commits separately to running the protocol). Both are visible in the existing repo; Alt-E names them both rather than picking one.

2. Across S2-S5, Alt-D and Alt-E are the only framings that turn the §11-style open-questions list into something doing real work — a load-bearing piece of the spec's promise, not boilerplate. This matches Steve's framing exactly: "not listing assumptions and known issues is an incomplete promise."

3. In S3 and S6, Alt-E's peer-level adoption promise lets peers cite *which* answers they assumed to open questions, which converts silent divergence into loud divergence. That's a real downstream payoff.

4. Alt-C alone is insufficient because it pushes the assumption list out of the spec, which loses the audit anchor (S2, S5). Alt-D alone is insufficient because it doesn't explain why peers feel obligated to behave as the doc says — the peer-level promise is the binding mechanism.

5. Alt-A and Alt-B fail core scenarios. Alt-A has nothing to diagnose; Alt-B bakes in centralization that PromiseGrid is trying to escape.

The conclusion implies several things about how this repo handles spec docs.

### Implications

- **A spec doc IS a promise** (specifically: a doc-level promise of the Alt-D shape). Its promiser is the document itself, identified by its pCID. Its promisee is anyone reading or running it. Its temporal shape is prediction about future state.

- **The promise is conditional.** Every spec doc's promise depends on a list of assumptions; the promise says "systems implementing me will interoperate, conditional on assumptions A1...An." The assumption list is a normative part of the spec, inside the pCID hash.

- **Open questions are also part of the promise** as conditioning factors. The promise says "...and conditional on the answers to open questions Q1...Qm being chosen consistently between peers." The open-questions list is inside the pCID hash.

- **Known issues / weak spots are part of the promise.** The promise says "...and acknowledging that under conditions W1...Wk, the protocol's interop guarantee is weakest." The known-issues list is inside the pCID hash.

- **Therefore: a spec doc with no assumptions list, no open-questions list, and no known-issues list is making a stronger promise than it can keep.** Either those lists are empty (rare, only true after long maturation) or omitting them makes the spec dishonest. The honest move is to write the lists explicitly, even when they say "we don't yet know X."

- **`specs/harness-spec-draft.md` is healthy under this framing.** Its §11 open-questions section is not pollution of the spec; it is part of the spec's promise. Promoting `specs/harness-spec-draft.md` into a sibling-protocol doc store does not require splitting workshop content from protocol content, because workshop content (assumptions, open questions, known issues) IS protocol content under Alt-E.

- **Each peer's adoption is a separate Alt-C promise.** When a peer claims to speak pCID X, the peer is promising "I will behave as the doc with hash X says I will." That promise can additionally cite which choices the peer made for any open questions in X. Peer-level adoption metadata is a real artifact worth designing.

## Decision Framing questions

DF-21.1: Is Alt-E (layered) the right framing, or is Alt-D (doc-only) sufficient and Alt-C-shaped peer adoption is overcomplication?

- (a) Alt-E (recommended). The peer-level promise is real and worth naming because it is what binds an actual peer to actual behavior; the doc-level promise alone would not connect to running code.
- (b) Alt-D only. The doc is the promise; peer adoption is just "you implemented it or you didn't" and doesn't deserve its own promise vocabulary.
- (c) Other framing not yet considered.

DF-21.2: Should the spec's assumption list, open-questions list, and known-issues list be required sections in every wire-lab spec doc, or are they best-practice but not mandatory?

- (a) Required (recommended). A spec doc that omits any of the three is incomplete-as-a-promise, per the analysis. The three sections become structural conventions of the wire-lab spec genre.
- (b) Best practice, not required. A spec doc may omit any of the three if its author judges them empty or irrelevant.
- (c) Required when present, conventional placement; absent allowed only if the doc explicitly states "no known assumptions / no open questions / no known issues" rather than silently omitting.

DF-21.3: When a peer adopts a spec, is the peer's choice of answers to open questions part of its adoption promise (machine-readable), or just commentary?

- (a) Part of the adoption promise, in a structured form (recommended). Peers can exchange adoption metadata and detect divergence on open-question answers.
- (b) Commentary only. Peers' answer choices are private; mismatches surface only when they cause concrete interop failures.
- (c) Required as commentary, optional as structured metadata; structure is a future enhancement.

DF-21.4: Does this TE imply any rename or restructure of `specs/harness-spec-draft.md` itself?

- (a) No (recommended). `specs/harness-spec-draft.md` is already shaped consistently with Alt-E; the TE clarifies what was implicit, but the file's content and location are fine. Layout questions (sibling-protocol doc store, file naming) are a separate TE.
- (b) Yes — `specs/harness-spec-draft.md` should be renamed to make its protocol-doc nature more obvious (e.g., `specs/wire-lab-spec.md`).
- (c) Yes — the §11 open-questions section should be relocated or restructured to make its normative role obvious.

The recommended set is (1.a, 2.a, 3.a, 4.a). Reason: lock the layered framing fully, make the three lists structural, give peer-level adoption first-class promise machinery, and defer layout questions to a follow-on TE.

## Decision status

`needs DF` — awaiting Steve's choice on DF-21.1 through DF-21.4. After DF, the locked decisions become DI entries in a new TODO file (TODO number to be assigned when DF lands).

## Implications for follow-on work

- **TODO 010 (provisional, after DF):** Lock the spec-doc-as-promise vocabulary as DI entries. Update `specs/harness-spec-draft.md` to add (or formalize) the three normative sections: Assumptions, Open Questions, Known Issues. Surface them as required structure for any future sibling spec docs.
- **TE-22 (planned):** Spec-doc-store layout — should sibling protocol docs live in a `specs/` directory? How do specs reference each other by pCID? Does `specs/harness-spec-draft.md` itself migrate? Deferred until TE-21 lands so the vocabulary is settled before the layout work.
- **TODO 006 (existing):** DI-provenance backfill for `specs/harness-spec-draft.md` settled statements. The Alt-E framing reshapes this work: settled statements need DI provenance because they're load-bearing parts of the doc's promise; open questions need DR provenance for the same reason. The two existing TODOs (006 and 007) become two halves of "make the doc-level promise explicit and audit-able."
