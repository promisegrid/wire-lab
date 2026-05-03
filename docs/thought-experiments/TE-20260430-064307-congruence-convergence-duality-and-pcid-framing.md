# TE-23: Congruence/convergence duality and pCID framing

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-20260430-064307

(First drafted 2026-04-30 06:43:07 UTC.)

## Status

open

## Decision under test

A small architectural question — *should the pCID be the hash of the spec document or the hash of the code that implements it?* — surfaced (in chat 2026-04-29) as the same question that splits the pre-DevOps configuration management world: Traugott's congruence camp (isconf, decomk, "Why Order Matters") versus Burgess's convergence camp (cfengine, Promise Theory, semantic spacetime). The conversation produced a framing essay at [`docs/essays/congruence-convergence-and-the-grid.md`](../essays/congruence-convergence-and-the-grid.md) (committed in `80cb8d3`) which argues, on its own terms:

- The grid's pCID-hashes-the-spec choice is **formalism-neutral** and was already correct on its merits before this conversation made the reasoning explicit.
- The two camps are **dual rather than opposed**; *autonomy is common ground* (the congruence camp's pull-not-push discipline is a Promise-Theory-shaped design that predates Promise Theory's vocabulary), and the genuine disagreement is over what an autonomous agent's promise is *about* — an ordered trajectory (congruence) versus a fixed-point state attractor (convergence).
- The wire-lab harness-spec's phrasing *"promises are assertions of state in the past, present, or future, often conditional"* extends Promise Theory in the direction that lets a journal entry and a host's pull-and-replay promise be two views of the same object — a candidate bridge primitive ("promise-about-trajectory") that may anchor a Church-Turing-equivalent equivalence theorem between the two camps.

This TE asks: *which of those framing claims, if any, should the wire-lab lock as Decision Intents now?*

This is not an implementation TE. There is no code to write, no on-disk machinery to introduce. The decisions under test are about what the wire-lab's foundational vocabulary acknowledges, what stays in essays-and-discussion, and what spawns follow-on TEs.

The four open Decision Framing questions are:

- **DF-23.1 pCID hash content.** Is the grid's top-level pCID the hash of a spec doc, the hash of code, both, or neither?
- **DF-23.2 Inner-CID slot at top level.** Does the harness-spec name an "inner-CID" slot (hCID, fCID, cCID, etc.) at the top level of the grid's vocabulary, or is naming such things deferred to per-protocol specs?
- **DF-23.3 Duality framing in the harness-spec.** Does the harness-spec acknowledge the congruence/convergence duality as foundational framing, with the essay as a referenced source — and if so, with how much weight?
- **DF-23.4 Equivalence-theorem follow-on.** Does the wire-lab open a separate research-flavored TE (TE-24) on the shape of a Church-Turing-equivalent congruence-convergence equivalence theorem, or does that work happen elsewhere (or not at all)?

## Assumptions

- A `Promise` in this repo is an autonomous speech act — an assertion of state in the past, present, or future, often conditional. (Carried from TE-21; locked.)
- A `pCID` is a CIDv1 hash of a spec document's bytes. Two parties claim to "speak protocol pCID X" when each implements the rules in the document whose CIDv1 is X. (Carried from `protocols/wire-lab.d/specs/harness-spec-draft.md` §1 and TE-22.)
- Spec docs are layered promises (TE-21 Alt-E): the doc itself promises future interop conditional on its assumptions/open-questions/known-issues lists, and each peer separately promises to behave as the doc says.
- TE-22 has locked the operational machinery for freezing, hashing, storing, citing, and replacing spec docs. The pCID-hashes-the-spec choice is therefore *already implemented* on disk; this TE asks whether to *acknowledge* that choice as foundational framing rather than incidental.
- The framing essay at `docs/essays/congruence-convergence-and-the-grid.md` exists and is the source-of-truth document for the duality framing. This TE references it; it does not duplicate it.
- Steve has stated (chat 2026-04-29) that hCID/fCID/cCID names "stay inside the protocol that uses them, not promoted to harness-spec top-level vocabulary," following the principle "Let each protocol name its own internals." DF-23.2 records that decision rather than re-litigating it; the question under test is whether the harness-spec should *mention* the principle (as opposed to silently embodying it).
- Steve has stated (chat 2026-04-29) that the pCID is "a protocol CID, not a promise CID" — content hash of a SPEC document. DF-23.1 records that decision; the question under test is whether the harness-spec should formally state this as foundational framing rather than as an in-passing definition.
- The wire-lab's bot identity is `stevegt+ppx@t7a.org` (stevegt-via-perplexity). Steve identity is `stevegt@t7a.org` (Steve Traugott). Where this TE quotes Steve's prior work (Traugott & Brown, isconf, decomk), it does so as a referenced source, not as a self-citation, because the wire-lab is treated as an independent project that happens to be authored partly by Traugott.

## Alternatives

This TE walks four DFs sequentially. Each DF has its own alternatives.

### DF-23.1: pCID hash content

What does a top-level pCID hash?

#### Alt-1.A: Hash of the spec document (status quo, already implemented)

A pCID is the CIDv1 hash of a spec document's bytes (per TE-22). Two peers "speak protocol pCID X" when each promises to behave according to the rules in the document whose CIDv1 is X. Multiple implementations of the same protocol can exist (Go, Rust, hand-rolled C); each is a separate piece of code that promises to honor the spec, not a separate protocol.

- **Easier:** matches the formalism-neutral framing in the essay. The pCID is the agreement; the agreement is a thing autonomous agents can promise about, regardless of whether their promise is about trajectory (congruence-style: "I run the exact code at hash Y") or attractor (convergence-style: "I behave as the spec says, by whatever local means"). Multiple implementations are first-class. Already implemented in TE-22.
- **Harder:** if a specific protocol genuinely requires byte-identical implementations (e.g., a deterministic-replay protocol where any deviation would corrupt the trajectory), that protocol must encode the requirement *inside* its spec, typically by saying "my payloads are addressed to inner code-hashes and my receivers run that exact code." The grid does not enforce code-identity at the top level; the protocol enforces it one level down. This is correct on the merits but is one extra step compared to "the protocol IS the code."

#### Alt-1.B: Hash of the code that implements the protocol

A pCID is the CIDv1 hash of an implementation's compiled bytes. Two peers "speak protocol pCID X" when each runs byte-identical code at hash X. Multiple implementations of the same protocol are not first-class; what exists is multiple protocols that happen to interoperate in practice.

- **Easier:** maximally deterministic. No daylight between protocol identity and implementation identity. Naturally fits congruence-camp intuitions (the trajectory IS the agreement). For deterministic-replay applications (the function-call shape Steve's original instinct hinted at), this is the obvious primitive.
- **Harder:** the grid's top-level addressing becomes a vote against multi-implementation protocols. Convergence-camp consumers (anyone who wants "the spec describes desired behavior; implementations may vary") are told their primitive is second-class. This is exactly the partisan choice the framing essay argues against. Also: implementations have platform-specific variation (compiler, libc, target arch) that protocol identity should usually transcend; tying pCID to byte-identical compiled code requires either reproducible-builds discipline at the protocol layer or weakening the meaning of "byte-identical."

#### Alt-1.C: Both, at different levels — pCID hashes spec, hCID/fCID/cCID hashes code

Top-level pCID hashes the spec document (Alt-1.A). The harness-spec acknowledges that *specific protocols* may declare an "inner code-hash" addressing scheme inside their own payloads — call it hCID, fCID, cCID, whatever the protocol's spec wants. The inner CID is named per-protocol, not at the harness-spec level.

- **Easier:** captures both worldviews without privileging either. Lets a "function-call protocol" exist as a normal grid protocol whose spec says "treat me congruently — my payloads are addressed to inner code-hashes and my receivers run that exact code." Lets a "desired-state protocol" exist as a normal grid protocol whose spec says "treat me convergently — my payloads are state predicates and my receivers run whatever local logic gets them there." Both are valid pCIDs at the top. (This is also what the essay argues for, structurally.)
- **Harder:** introduces a two-level addressing distinction (top-level = spec-as-pCID; inner = per-protocol). Readers must understand both. Mitigated by the fact that most readers only care about the top level; the inner layer is invisible unless a specific protocol requires it.

#### Alt-1.D: pCID hashes a manifest that points at both a spec doc and reference code

Top-level pCID hashes a small manifest file that names both the spec document's CIDv1 and (optionally) reference implementation code's CIDv1. Two peers "speak protocol pCID X" when each peer's manifest hash equals X, which transitively pins both spec and reference code.

- **Easier:** records both the spec and the canonical implementation in a single content-addressed object. A reader can resolve from pCID to either prose or code without ambiguity.
- **Harder:** the manifest layer adds a level of indirection that doesn't pay for itself. Most protocols won't have a single canonical reference implementation (or won't want to lock to one); the optional-code field is then either always-empty or always-pointing-at-one-arbitrary-implementation. The two-level addressing in Alt-1.C achieves the same expressivity without requiring a manifest object.

### DF-23.2: Inner-CID slot at top level

Does the harness-spec name an "inner-CID" slot (hCID, fCID, cCID, etc.) at the top level of the grid's vocabulary?

#### Alt-2.A: Yes, name it explicitly in the harness-spec vocabulary

The harness-spec adds a named slot — say, "inner-CID" or "implementation-CID" — alongside pCID in its vocabulary section. Specific protocols then have a pre-named hook to fill in.

- **Easier:** uniform vocabulary across protocols. A protocol author looking for "where do I put the code-hash" finds the answer in the harness-spec.
- **Harder:** premature standardization. Different protocols may want different inner-CID semantics: a function-call protocol's inner CID is code; an interpreter protocol's inner CID is a sub-spec; a content-addressed-RPC protocol's inner CID is yet another flavor. Forcing them all into one named slot privileges one shape over the others. Also: the harness-spec's vocabulary stays cleaner if it only names what the harness *itself* uses, not what protocols *might* use.

#### Alt-2.B: No, defer naming entirely to per-protocol specs

The harness-spec does not introduce an inner-CID name. Each protocol's spec is free to define its own inner-CID concept (or not) under whatever name fits its semantics. The principle "Let each protocol name its own internals" is the rule.

- **Easier:** keeps the harness-spec's vocabulary minimal. Lets protocols evolve diverse inner-CID semantics without coordinating with the harness. Matches Steve's stated decision in chat 2026-04-29.
- **Harder:** there is no single name a reader can search for; readers learn each protocol's inner-CID name from its spec. This is the same situation as TCP's "port number" not being a concept HTTP uses — and it works fine because protocol-specific vocabulary is appropriately scoped.

#### Alt-2.C: No, but the harness-spec acknowledges the pattern as a footnote

The harness-spec does not introduce an inner-CID name, but does carry a non-normative note observing that some protocols may content-address their payloads' implementations and may give such addresses protocol-specific names. The note exists to head off the "where do I put the code-hash?" reader question without committing to a slot.

- **Easier:** preserves the minimal-vocabulary discipline (Alt-2.B) while signaling to confused readers that the pattern is recognized. Avoids a future reader concluding the harness has no opinion.
- **Harder:** non-normative notes can grow into normative cruft over time if not policed. Mitigated by clear marking ("informative" or "non-normative" status).

#### Alt-2.D: No, and the harness-spec is silent on the pattern

The harness-spec does not introduce an inner-CID name and says nothing about the pattern. Protocol authors discover the pattern (or invent it independently) when they need it. The principle "Let each protocol name its own internals" is implicit in the harness's minimalism.

- **Easier:** smallest harness-spec footprint. Most protocols never need an inner-CID concept at all (most protocols are not function-call protocols, content-addressed-RPC, or meta-circular interpreters); for those, an inner-CID footnote is dead weight.
- **Harder:** the function-call-shaped reader has to discover the pattern themselves. May lead to inconsistent inner-CID conventions across protocols if protocol authors don't share a reference. (Counterargument: TCP's `<port, IP>` shape was discovered by application authors, not mandated by IP, and converged anyway because the shape is forced by the semantics.)

### DF-23.3: Duality framing in the harness-spec

Does the harness-spec acknowledge the congruence/convergence duality as foundational framing, with the framing essay as a referenced source? If so, with how much weight?

#### Alt-3.A: Full framing — the harness-spec opens with a duality section

The harness-spec adds an early prose section (right after the §1 opening, before §2 vocabulary) titled something like "Foundational framing: congruence and convergence" that summarizes the essay's argument in 1-2 paragraphs and links to the essay for full treatment. The harness-spec's design choices (spec-as-pCID, autonomy-as-promise, past/present/future-tensed promises) are explicitly motivated by this framing.

- **Easier:** every reader of the harness-spec encounters the framing immediately. Future design choices have an explicit foundation to test against ("does this proposal preserve formalism-neutrality? does it preserve autonomy as common ground?"). Anchors the wire-lab's identity to the reconciliation ambition.
- **Harder:** the harness-spec becomes longer and more philosophical. Readers who came for protocol mechanics have to wade through framing they may not need. The reconciliation ambition is *aspirational* (the equivalence theorem is conjectural); locking it as foundational framing risks promising more than the wire-lab can deliver, and the framing may need revision later.

#### Alt-3.B: Brief acknowledgement — one paragraph + link to the essay

The harness-spec adds a single short paragraph (in §1 or in a "Background" section near the top) noting that the design is informed by the congruence/convergence duality, with a link to the essay. The paragraph is non-normative and explicitly framed as background context, not as a constraint on future design choices.

- **Easier:** acknowledges the framing without inflating the harness-spec's scope. Lets the essay carry the full argument; the harness-spec only carries the pointer. Future readers know where to find the framing if they want it.
- **Harder:** the framing's influence on design choices is implicit rather than explicit; readers must infer the connection. For a foundational document, that's a missed teaching opportunity.

#### Alt-3.C: Vocabulary-only — duality acknowledged through term choice, no explicit section

The harness-spec does not call out the duality explicitly, but its vocabulary choices reflect it: pCID hashes a spec doc (Alt-1.A); promises span past/present/future (already locked); peers are autonomous (already locked); spec-as-promise framing (TE-21). The essay exists at `docs/essays/congruence-convergence-and-the-grid.md` and is discoverable by anyone who reads the repo top-down, but the harness-spec does not link to it.

- **Easier:** keeps the harness-spec terse and protocol-mechanics-focused. The duality is "in the design" but not "on the page." Readers who want the framing find the essay through the docs/ tree.
- **Harder:** readers of the harness-spec in isolation (which is how most spec readers arrive) get no clue that the framing exists. The deliberate formalism-neutrality of pCID-as-spec-hash looks arbitrary rather than principled. Misses the chance to plant a flag.

#### Alt-3.D: Silent — no acknowledgement at all

The harness-spec does not acknowledge the duality. The essay exists in the repo but is not linked from any normative document; it lives as a research note that influenced the design but is not part of the design's stated foundation.

- **Easier:** maximally minimal. Treats the duality framing as a private notebook entry that influenced design choices but does not bind future design.
- **Harder:** the framing essay's influence is invisible to anyone not in the original conversation. Future maintainers of the wire-lab may unintentionally violate formalism-neutrality (or autonomy-as-common-ground) without realizing they're departing from a considered position. The reconciliation ambition becomes folklore rather than design.

### DF-23.4: Equivalence-theorem follow-on

Does the wire-lab open a separate research-flavored TE (TE-24) on the shape of a Church-Turing-equivalent congruence-convergence equivalence theorem?

#### Alt-4.A: Open TE-24 now, research-flavored

A follow-on TE — TE-24 — is opened immediately to work out: (i) candidate definitions of "effectively administrable infrastructure"; (ii) candidate formal models for the two formalisms (ordered journals replayed under host promise; autonomous-agent promise systems with convergence dynamics); (iii) the smallest cases in which the equivalence claim could be tested in simulation; (iv) the obstructions that would, if found, refute the equivalence.

- **Easier:** keeps the thread alive while the conversation is fresh. Aligns research effort with implementation work (TE-24 can run in parallel with concrete wire-lab development). Plants a stake the wire-lab can be held to.
- **Harder:** research-flavored TEs may not converge to a DI for a long time, possibly never. They risk distracting from concrete wire-lab work. The wire-lab may not have the resources or the right contributors to push a theorem-shaped question to completion.

#### Alt-4.B: Open TE-24 later, when the wire-lab has more concrete protocols to test against

TE-24 is on the roadmap but not opened now. The reasoning: a formal equivalence claim is most testable when there are concrete protocols (one congruence-shaped, one convergence-shaped) to model in simulation. Until both shapes exist, the equivalence is a thought experiment without test material.

- **Easier:** lets concrete wire-lab work generate the test material that TE-24 will need. Avoids opening a TE that has nothing actionable in it yet.
- **Harder:** the thread may go cold. New conversations may overwrite the framing. The researcher who picks TE-24 up later must reload the essay and the original conversations.

#### Alt-4.C: Don't open TE-24 in the wire-lab; do this work elsewhere

The equivalence theorem is too large for a wire-lab TE. It belongs in a separate research project — a paper, a thesis, a different repo — that cites the wire-lab as a substrate but does the formal work on its own footing.

- **Easier:** keeps the wire-lab focused on what it can actually deliver (a working substrate). Avoids the impedance mismatch between TE-flavored design conversations and theorem-shaped formal work.
- **Harder:** the equivalence work loses the wire-lab's test material (concrete protocols, real promise vocabulary) as a starting point. The researcher who eventually does this work has to re-derive the same framing the essay already articulates.

#### Alt-4.D: Don't open TE-24 at all; treat the equivalence as folklore

The framing essay exists; the equivalence claim is in it; that's enough. No formal follow-on. If someone later proves or refutes the equivalence, that work happens organically in whatever venue makes sense at the time.

- **Easier:** lowest-effort. Lets the framing inform design without committing to follow-on research.
- **Harder:** abandons the most concrete reason to think the equivalence theorem is reachable (the promise-about-trajectory thread). Treats the essay's load-bearing claim #5 ("a Church-Turing-equivalent might be reachable from there") as decoration rather than a research direction.

## Scenarios

Six scenarios, each played against the alternatives.

### S1 (genesis reading): a new contributor reads the harness-spec for the first time

A new contributor (human or bot) reads `protocols/wire-lab.d/specs/harness-spec-draft.md` for the first time. They have no prior context.

What does the harness-spec teach them about why pCID hashes a spec doc and not code?

- **DF-23.1:** Alt-1.A is the status quo; its rationale is invisible without DF-23.3 carrying the framing.
- **DF-23.2:** Alt-2.B (no inner-CID slot) means readers do not encounter inner-CID vocabulary in the harness; this is fine if they're not building a function-call protocol.
- **DF-23.3:** Alt-3.A (full framing) puts the rationale on the first page; Alt-3.B (brief ack) gives them a link to follow if curious; Alt-3.C (vocabulary-only) leaves them to infer; Alt-3.D (silent) leaves them with no explanation.
- **DF-23.4:** Independent of this scenario.

S1 verdict: Alt-3.A or Alt-3.B serve a new reader best. Alt-3.A spends harness-spec real estate on philosophy; Alt-3.B keeps the harness-spec lean while still surfacing the framing. Alt-3.C and Alt-3.D leave the genesis-reader confused or unaware.

### S2 (function-call protocol author): someone wants to define a congruent protocol on the grid

A contributor wants to define a function-call protocol on the grid: payloads name an inner code-hash, receivers run that exact code. They open the harness-spec to learn how to declare their inner-CID concept.

- **DF-23.1:** Alt-1.A or Alt-1.C tell them their protocol's pCID hashes their spec; the inner code-hash lives inside their protocol.
- **DF-23.2:** Alt-2.A pre-names a slot for them (uniform but possibly mis-shaped); Alt-2.B leaves naming to them ("hCID" feels natural); Alt-2.C gives them a non-normative pointer that acknowledges the pattern; Alt-2.D leaves them to invent the convention.
- **DF-23.3:** Alt-3.A or Alt-3.B explicitly tell them their protocol is one of the two valid shapes ("congruent"); Alt-3.C and Alt-3.D leave them to figure that out from context.
- **DF-23.4:** Independent.

S2 verdict: Alt-2.B or Alt-2.C handle this contributor well. Alt-2.A risks forcing their inner-CID into a slot whose semantics may not fit their protocol. Alt-3.A or Alt-3.B reassure them that congruent-protocol-on-this-grid is a recognized shape, not a hack.

### S3 (desired-state protocol author): someone wants to define a convergent protocol on the grid

A contributor wants to define a convergent protocol: payloads describe desired state, receivers run whatever local logic gets them there. Implementations may vary across peers.

- **DF-23.1:** Alt-1.A is exactly what they need — the pCID is the spec; implementations are first-class.
- **DF-23.2:** Their protocol does not need an inner-CID at all; Alt-2.B, Alt-2.C, or Alt-2.D leave it out, which is correct.
- **DF-23.3:** Alt-3.A or Alt-3.B tell them their shape is recognized; Alt-3.C and Alt-3.D leave them to infer.
- **DF-23.4:** Independent.

S3 verdict: this scenario is the easy case for Alt-1.A — convergence-shaped protocols are the shape pCID-as-spec-hash *naturally* fits. Alt-3.A or Alt-3.B add the explicit acknowledgement that convergence is one of the two recognized shapes; Alt-3.C and Alt-3.D leave the contributor to infer.

### S4 (LLM-orchestration protocol author): the deterministic-code-calls-LLM case

A contributor wants to define a protocol where deterministic client code (referenceable by hCID) sends a system message (a pCID — a behavioral spec) plus a user prompt (payload) to an LLM and treats the model's response as the protocol's reply. This is the pCID-inside-hCID nesting case from §3.1 of the essay.

- **DF-23.1:** Alt-1.A (top-level pCID hashes the spec) and Alt-1.C (inner CIDs allowed per-protocol) together capture this case cleanly. Alt-1.A says the protocol-as-a-whole is named by its spec's pCID; Alt-1.C says the protocol is free to declare that its payloads are themselves pCIDs (the system message) wrapped in an hCID context (the client code).
- **DF-23.2:** Alt-2.B or Alt-2.C let the protocol author name its internals freely (some hCID-ish slot for the client code, the existing pCID for the system message).
- **DF-23.3:** Alt-3.A or Alt-3.B explicitly recognize this as a hybrid shape — congruent code orchestrating a convergent receiver. Alt-3.C and Alt-3.D do not.
- **DF-23.4:** This shape is one of the most concrete test cases for the equivalence-theorem framing (a deterministic harness orchestrating a non-deterministic responder, observably tracing the same agreement two different ways). Alt-4.A or Alt-4.B keep that test material reachable from a future TE-24.

S4 verdict: this case shows that DF-23.1's Alt-1.C (allowing inner CIDs per-protocol) is doing real work even when the inner CID is something other than "byte-identical implementation code." Alt-2.B or Alt-2.C handle the protocol's internal naming; Alt-3.A or Alt-3.B make the hybrid shape recognizable.

### S5 (future maintainer departs from formalism-neutrality): someone proposes a code-as-pCID protocol at the top level

A future maintainer of the wire-lab proposes that the harness-spec be amended so pCIDs name code rather than specs (or that a parallel "cCID" be introduced at the top level for byte-identical-code addressing).

- **DF-23.1:** Alt-1.A is the locked answer; the proposal is a deviation that requires an explicit DI to override.
- **DF-23.2:** Alt-2.B or Alt-2.C kept inner-CID concepts out of the harness vocabulary, so the proposal cannot trivially add a "cCID" by extending an existing slot — it must defend itself as a top-level vocabulary addition.
- **DF-23.3:** Alt-3.A puts the formalism-neutrality framing on the harness-spec's front page, so the proposal is immediately testable against it. Alt-3.B requires a click-through to the essay. Alt-3.C and Alt-3.D leave the proposal's framing-violation invisible.
- **DF-23.4:** Independent.

S5 verdict: Alt-3.A is the strongest defense against drift. Alt-3.B is acceptable defense. Alt-3.C and Alt-3.D leave future maintainers without a stated principle to test proposals against.

### S6 (someone wants to formally test the equivalence claim)

A researcher (could be Steve, could be a future contributor, could be a graduate student) wants to formally test or refute the equivalence theorem outlined in the essay. They want to know what the wire-lab's stated position is and what test material exists.

- **DF-23.1, DF-23.2, DF-23.3:** these affect what the researcher finds in the harness-spec, but the researcher's primary object of study is the essay and any follow-on TEs.
- **DF-23.4:** Alt-4.A gives them a TE-24 to walk into; Alt-4.B tells them to wait until concrete protocols exist (and gives them a roadmap stake to track); Alt-4.C tells them to do this work outside the wire-lab; Alt-4.D leaves them to reload the framing from the essay alone.

S6 verdict: Alt-4.A is the most welcoming-to-research; Alt-4.B is a reasonable defer; Alt-4.C and Alt-4.D push the researcher away from the wire-lab's substrate. Given that the wire-lab has not yet produced concrete protocols of either shape, Alt-4.B may be the most honest answer: the test material does not yet exist.

## Conclusions

Across S1-S6, a coherent recommended set emerges:

- **DF-23.1: Alt-1.A** — pCID hashes the spec document. This is already implemented (TE-22 and prior) and is correct on the merits; this TE locks it as a foundational framing position rather than an incidental machinery choice. **Plus an explicit acknowledgement of Alt-1.C** as the way the grid hosts congruent protocols: the harness-spec records that specific protocols may declare inner code-hash addressing inside their own payloads (without naming the inner-CID at the harness level — see DF-23.2). Rejected: Alt-1.B (code-as-pCID) is a partisan choice the framing essay argues against; Alt-1.D (manifest pointing at both) adds an indirection layer that doesn't pay for itself.

- **DF-23.2: Alt-2.B** — no inner-CID slot at the top level; defer naming entirely to per-protocol specs. Matches Steve's stated decision in chat 2026-04-29 ("Let each protocol name its own internals"). Rejected: Alt-2.A (pre-named slot) privileges one inner-CID semantics over others; Alt-2.C and Alt-2.D are weaker forms of the same defer-to-protocols answer, but Alt-2.B is the cleanest. Note: the essay's nesting taxonomy (pCID-in-pCID, hCID-in-pCID, hCID-in-hCID, pCID-in-hCID) is preserved as discoverable reasoning in `docs/essays/congruence-convergence-and-the-grid.md`; the harness-spec does not duplicate it.

- **DF-23.3: Alt-3.B** — brief acknowledgement of the duality framing in the harness-spec (one paragraph in §1 or a short "Background" subsection), with a link to `docs/essays/congruence-convergence-and-the-grid.md`. The paragraph is explicitly framed as background context, not as a normative constraint; the design choices it informs (pCID hashes spec, autonomy-as-promise, past/present/future tenses) are themselves locked elsewhere (TE-21, TE-22, this TE). Rejected: Alt-3.A inflates the harness-spec with philosophy and risks promising more than the wire-lab can deliver; Alt-3.C and Alt-3.D leave the framing invisible to readers who arrive through the harness-spec rather than through the essay. The brief acknowledgement is the smallest move that defends against future drift (S5) while keeping the harness-spec's main text protocol-mechanics-focused.

- **DF-23.4: Alt-4.B** — open TE-24 *later*, when the wire-lab has at least one congruence-shaped protocol and one convergence-shaped protocol to model in simulation. Track it as a roadmap stake in this TE's "Implications for follow-on work" section so the thread does not go cold; reload from the essay when ready. Rejected: Alt-4.A opens a research-flavored TE without test material, which risks stalling; Alt-4.C pushes the work out of the wire-lab and loses the test material the wire-lab will eventually generate; Alt-4.D abandons the most concrete reason to think the equivalence is reachable (the promise-about-trajectory thread) and treats the essay's load-bearing claim #5 as decoration.

The full recommended set across the four DFs is **(1.a + 1.c-as-permission, 2.b, 3.b, 4.b)**.

### Implications

- **Harness-spec gets a brief Background subsection.** A new short subsection — proposed title "Background: congruence and convergence" — appears near the top of `protocols/wire-lab.d/specs/harness-spec-draft.md` (after §1 opening, before §2 vocabulary). One paragraph (3-5 sentences) summarizes the duality framing and links to the essay. The subsection is non-normative; it carries the explicit marker "*This subsection is informative, not normative. The design choices it describes are locked in TE-21, TE-22, and TE-23.*" When the next freeze of the harness-spec is cut, the Background subsection becomes part of the frozen pCID's content; the link to the essay is by relative path (`../docs/essays/congruence-convergence-and-the-grid.md`) so it survives a host migration.

- **Harness-spec vocabulary stays small.** No "inner-CID," "hCID," "fCID," "cCID" entries are added. The harness-spec carries `pCID` (already present) and nothing else CID-shaped at the top level. Per-protocol specs are free to introduce their own inner-CID concepts under whatever name fits.

- **Essay becomes a normatively-cited source from the harness-spec.** The Background subsection's link to `docs/essays/congruence-convergence-and-the-grid.md` makes the essay a referenced source, not a private note. The essay is therefore subject to the strict-cross-ref rule from TE-22 (a draft cross-reference must cite a frozen pCID; drafts may not cite other drafts) once the harness-spec next freezes. **Open question (becomes part of TE-23's TODO):** does an essay file in `docs/essays/` need to be frozen the same way a spec file does? Tentative answer: yes, treated as a non-protocol referenced source, freezing produces `docs/essays/congruence-convergence-and-the-grid-{cidv1}.md` and the manifest entry records `kind: essay` to distinguish from `kind: spec`. Sub-question: does the essay need its own pCID class, or does CIDv1 suffice for any referenced source? Recommendation: CIDv1 suffices; `kind` field in the manifest distinguishes essays from specs.

- **TE-24 is on the roadmap but not opened now.** The TODO directory carries an entry — say, `TODO/012-te-24-equivalence-theorem-shape.md` — that records: (i) the framing essay's promise-about-trajectory thread is the load-bearing claim that motivates TE-24; (ii) TE-24 cannot productively open until the wire-lab has at least one congruence-shaped protocol (function-call shape) and one convergence-shaped protocol (desired-state shape) to model in simulation; (iii) the trigger to open TE-24 is the first concrete protocol of either shape reaching freeze. Until then, the framing essay carries the thread.

- **The framing essay does not require revision.** The essay's claims (1) through (5) align with this TE's recommended decisions. No edits to `docs/essays/congruence-convergence-and-the-grid.md` are needed as a consequence of TE-23.

- **Future spec readers have a stated principle to test proposals against.** "Does this proposal preserve formalism-neutrality?" and "Does this proposal preserve autonomy as common ground?" become readable, citable tests. A future maintainer who proposes code-as-pCID at the top level (S5) is asked to defend against the framing the harness-spec's Background subsection now carries.

- **No tooling change.** TE-22's machinery (single Go binary `tools/spec` with `freeze`, `check`, `cid`, `ls` subcommands) requires no extension. The Background subsection is just additional spec text. The essay-as-referenced-source open question above may eventually require `tools/spec freeze` to accept a `--kind=essay` flag or accept any path under `specs/` or `docs/essays/`; that is deferred to the TODO that captures the open question.

## Decision Framing questions

DF-23.1: pCID hash content.

- (a) Alt-1.A — pCID hashes the spec document (recommended).
- (b) Alt-1.B — pCID hashes the implementation code.
- (c) Alt-1.C — pCID hashes the spec; per-protocol inner CIDs allowed (recommended as a permission complementing Alt-1.A).
- (d) Alt-1.D — pCID hashes a manifest pointing at both spec and reference code.

DF-23.2: Inner-CID slot at top level.

- (a) Alt-2.A — yes, name an inner-CID slot in the harness-spec vocabulary.
- (b) Alt-2.B — no, defer naming entirely to per-protocol specs (recommended).
- (c) Alt-2.C — no, but the harness-spec carries a non-normative footnote acknowledging the pattern.
- (d) Alt-2.D — no, and the harness-spec is silent on the pattern.

DF-23.3: Duality framing in the harness-spec.

- (a) Alt-3.A — full framing section in the harness-spec (1-2 paragraphs), with the design choices explicitly motivated by the duality.
- (b) Alt-3.B — brief acknowledgement (one paragraph + link to the essay), explicitly non-normative (recommended).
- (c) Alt-3.C — vocabulary-only acknowledgement; no explicit framing section.
- (d) Alt-3.D — silent; the essay exists but is not linked from the harness-spec.

DF-23.4: Equivalence-theorem follow-on.

- (a) Alt-4.A — open TE-24 now, research-flavored.
- (b) Alt-4.B — open TE-24 later, when the wire-lab has at least one congruence-shaped and one convergence-shaped protocol to model (recommended).
- (c) Alt-4.C — don't open TE-24; do this work elsewhere.
- (d) Alt-4.D — don't open TE-24 at all; treat the equivalence as folklore.

The recommended set is **(1.a with 1.c-as-permission, 2.b, 3.b, 4.b)**. The pattern is:

- *Lock what is already implemented and correct on the merits* (Alt-1.A: pCID hashes the spec; this is what TE-22 implemented).
- *Acknowledge the per-protocol freedom* that lets congruent protocols nest inside (Alt-1.C as a permission, not a requirement).
- *Keep the harness-spec vocabulary minimal* (Alt-2.B: no inner-CID at the top).
- *Plant the framing flag without inflating the harness-spec* (Alt-3.B: brief acknowledgement + link).
- *Defer research work until test material exists* (Alt-4.B: TE-24 on the roadmap, not opened yet).

Already-locked decisions, not under DF in this TE (carried from chat 2026-04-29 and prior TEs):

- **pCID format = CIDv1.** Multibase + multihash + codec wrap. (TE-22.)
- **Spec doc as layered promise.** (TE-21 Alt-E.)
- **"Let each protocol name its own internals."** Steve's stated principle in chat 2026-04-29 (this TE locks it as DF-23.2 Alt-2.B).
- **pCID is "protocol CID, not promise CID."** Steve's correction in chat 2026-04-29 (this TE locks it as DF-23.1 Alt-1.A's framing).
- **Promises are assertions of state in the past, present, or future, often conditional.** (Locked in harness-spec vocabulary; central to the essay's promise-about-trajectory thread.)

## Decision status

`open` — recommendations stated, awaiting Steve's decision on the four DFs. The expected lock pattern is (1.a + 1.c-as-permission, 2.b, 3.b, 4.b) but Steve may diverge — particularly on DF-23.3 (the harness-spec footprint cost of acknowledging the framing) and DF-23.4 (whether to open TE-24 now versus later).

Amendment history:

- 2026-04-30 (initial): TE drafted with four DFs (pCID hash content, inner-CID slot, duality framing, equivalence-theorem follow-on). Recommended set (1.a + 1.c-as-permission, 2.b, 3.b, 4.b).

## Implications for follow-on work

- **TODO 012 (provisional):** Once Steve locks DF-23.3 as Alt-3.B (or stronger), draft the Background subsection for `protocols/wire-lab.d/specs/harness-spec-draft.md` and add the link to the essay. Re-freeze the harness-spec on the next freeze cycle. The Background subsection lands in the next pCID; the previous pCID remains as historical evidence per TE-22's append-only-log discipline.

- **TODO 013 (provisional):** Capture the essay-as-referenced-source open question raised under "Implications" above. Decide whether `docs/essays/*.md` files freeze the same way `specs/*.md` files do, and whether the manifest distinguishes `kind: spec` from `kind: essay`. Tentative answer in the conclusions; needs explicit decision.

- **TODO 014 (provisional, becomes TE-24 trigger):** Track the wire-lab's first congruence-shaped protocol (function-call shape) and first convergence-shaped protocol (desired-state shape) reaching freeze. When both exist, open TE-24 to formally examine the promise-about-trajectory equivalence claim against concrete test material. Until then, the framing essay carries the thread.

- **No change required to TE-21 or TE-22.** TE-21's spec-doc-as-promise framing and TE-22's pCID machinery are both consistent with this TE's recommendations; TE-23 layers framing on top, not corrections to either.

- **Bot-side change:** when the bot drafts new TEs going forward, it should test design proposals against the two stated principles (formalism-neutrality and autonomy-as-common-ground) where relevant. This is a posture instruction, not a code change.
