# TE-1: Promise-stack ordering

*Thought experiment, part of the [PromiseGrid Wire Lab](../../specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-20260427-180000

(Original TE-1 from `specs/harness-spec-draft.md` v2; first drafted at the time of v2 creation. Expanded into full scenario form on 2026-04-28 by the bot.)

## Decision under test

How is a promise stack ordered? Specifically:

- **D1 (encoding order):** Are frames serialized innermost-first or outermost-first on the wire?
- **D2 (evaluation order):** Does a receiver evaluate the stack outermost-first ("peel until I find what I trust"), innermost-first ("verify the kernel of meaning first, then ratify outer claims"), or both/either?
- **D3 (canonical frame placement):** For a given assertion type (transport-promise, signature, capability-issuance, content-promise), is there a canonical position in the stack, or is position arbitrary?
- **D4 (convergence on partial peeling):** When two endpoints have peeled different amounts of the same logical message, what semantics define their views as still-the-same-message?

D1 and D2 are independent (encoding order and evaluation order can disagree if framing carries length prefixes). D3 and D4 are downstream of D2. This TE attempts to narrow all four; each may end up with its own DI.

## Assumptions

- The spec's `Promise` shape (§1.1 of `specs/harness-spec-draft.md`) is roughly right: each frame has a `promiser`, a `assertion`, optional `body`, optional `evidence`, optional `ttl`. We are not testing the shape; we are testing how a sequence of these is laid out and consumed.
- The `promstack` library has three operations: `Wrap`, `Peel`, `Project` (§1.1).
- A receiver is permitted to **accept**, **defer**, or **reject** each promise individually. Stack-level accept/reject is the conjunction of frame-level decisions plus a policy.
- The transport may be a stream (TCP, TLS, Noise, in-process channel) or a datagram (UDP, capability port, hyperedge wire).
- Routing intermediaries may add their own frames mid-flight; this is explicitly desired (§1.2 chain-of-custody).
- Promiser identities span agents, transports, runtimes, and kernels — any of which can be a frame's promiser.
- The trust ledger (§2) is per-assertion-type, per-promiser. Different orderings may make different assertions easier or harder to track.

### Threat / trust model

- Any frame may be forged in isolation; cryptographic evidence within a frame is the only protection for that frame.
- An intermediary can append its own frames truthfully or fraudulently; the receiver's trust in the intermediary determines weight.
- An intermediary can also **strip** frames (drop them on egress). The model must distinguish "this frame was never present" from "this frame was stripped".
- Replay is a real concern; frames carry no implicit ordering proof relative to each other unless an explicit promise asserts it.
- Parties may be running mixed versions of `promstack` and may disagree about how to interpret a frame they don't recognize.

## Candidate alternatives

### Alt-A: Outermost-first on wire and in evaluation

- Wire order: frames serialized outermost-first. The first bytes a reader gets describe the most recently added (outermost) promise.
- Evaluation: receiver peels outermost, decides accept/defer/reject, then descends.
- Canonical placement: transport-promise outermost (it was the last thing added by the inbound edge); content-promise innermost (it was the first thing the original author wrote); signature can sit at any layer depending on what it's signing.
- Partial-peel semantics: a receiver that has peeled N frames "is reading the message at depth N"; the remaining suffix is still a valid message rooted at the inner content.

### Alt-B: Innermost-first on wire, outermost-first in evaluation

- Wire order: frames serialized innermost-first. The first bytes are the original content-promise from the author; subsequent bytes are routing, transport, and outer wrappings.
- Evaluation: same as Alt-A — receiver peels outermost-first, which means seeking to the *end* of the byte stream first, then walking back.
- Canonical placement: same as Alt-A.
- Partial-peel semantics: same as Alt-A but tricky for streaming receivers because they cannot evaluate until they have the whole thing.

### Alt-C: Innermost-first on wire and in evaluation

- Wire order: innermost-first.
- Evaluation: receiver evaluates innermost-first — verifies the content-promise (and its signature, if signature is innermost) before ratifying any outer promises about it.
- Canonical placement: content-promise innermost, signature wrapping the content directly outside it, transport-promise outermost. The structure is "what was promised, who promised it, how it got here", read in that order.
- Partial-peel semantics: a streaming receiver verifies content as bytes arrive and can begin acting on the inner payload before the outer transport-promise has been fully consumed.

### Alt-D: Order-agnostic — Project, not Peel

- Wire order: implementation-defined; both Alt-A and Alt-C interoperate via length-prefixed framing or CBOR sequence boundaries.
- Evaluation: receivers ignore order entirely. They use `Project(msg, predicate)` to extract all frames matching some predicate (e.g., "all signatures by Alice", "all transport-promises", "all capability-issuances") and apply policy across the projected set. Order is informational only.
- Canonical placement: there is no canonical position; frames are an unordered multiset with provenance.
- Partial-peel semantics: there is no peeling; receivers consume frames as a set. Streaming is supported by accepting frames as they arrive and re-evaluating policy after each.

### Alt-E: Hybrid — wire-order canonical, evaluation-order policy-driven

- Wire order: outermost-first by convention (matches Alt-A; matches the "transport is the most recent author" intuition).
- Evaluation: receivers may use `Peel` (sequential) or `Project` (set-wise) freely. The choice is policy, not protocol.
- Canonical placement: weakly canonical. Assertion types have *typical* positions, but a receiver that finds them out of typical position must still accept them — `Project` will find them regardless.
- Partial-peel semantics: a peeling receiver that has consumed N frames has a well-defined remaining-suffix message; a projecting receiver has a multiset of frames it has examined so far.

## Scenario analysis

The five alternatives are evaluated against AGENTS.md's prescribed scenario dimensions. Each cell describes what the alternative makes easier, what it makes harder, and what new obligations it creates.

### S1 — Normal operation: signed message over TLS through one router

A typical message: Alice authors a content-promise, signs it, sends over TLS through a router R, who appends a transport-promise. Bob receives.

- **Alt-A (outermost-first wire+eval):** Easiest. Bob reads transport-promise (R asserted), then signature-promise (Alice asserted), then content-promise. Each is verifiable in encounter order. **Easier:** receiver can early-reject if it doesn't trust R. **Harder:** Bob doesn't see Alice's content until he's processed two outer frames he might not care about. **New obligation:** routers must reliably append (not prepend); the wire reflects time-of-application.
- **Alt-B (innermost-first wire, outermost-first eval):** Same evaluation, worse wire layout. Bob must buffer to the end of the message before evaluating anything. **Easier:** matches some legacy serializations where the original payload comes first. **Harder:** streaming is impossible; all messages are effectively datagrams. **New obligation:** length-prefix or end-marker discipline.
- **Alt-C (innermost-first both):** Bob verifies Alice's content first, then Alice's signature, then R's transport-promise as ratification. **Easier:** content can begin processing as soon as inner bytes arrive (streaming friendly for trusted-content). **Harder:** Bob processes content from a sender he hasn't authenticated yet; if R is hostile and Alice's signature is forged, Bob has wasted work. **New obligation:** receivers must defer side-effects until outer frames are also accepted.
- **Alt-D (Project-only):** Bob runs `Project(msg, signature-predicate)` to find Alice's signature, `Project(msg, transport-predicate)` to find R's transport-promise, evaluates both, then decides. **Easier:** order-independent; matches receiver-policy diversity. **Harder:** receiver must walk the entire stack before deciding anything. **New obligation:** every frame must self-identify its assertion type unambiguously; predicates must be cheap.
- **Alt-E (hybrid):** Defaults match Alt-A; receivers that prefer Alt-D semantics use `Project`. **Easier:** doesn't force a single style. **Harder:** two valid styles double the surface to test and document. **New obligation:** the spec must declare which positions are conventional but non-binding.

### S2 — Failure / corruption / incomplete writes

A message arrives truncated, or a router crashes mid-append, or bytes are flipped.

- **Alt-A:** A truncation that drops the inner content but leaves outer transport-promise intact looks like a "message that says the network said something arrived, but the something is missing". The receiver can detect this from the missing inner frame. **Easier:** outer-first means the most recent action (the inbound transport) is most likely to be intact; trust ledger updates against R correctly even if Alice's content was already lost upstream. **Harder:** an attacker can truncate and produce a valid-but-empty-content message; receivers must require a non-trivial inner content unless explicitly waived. **New obligation:** "minimum useful prefix" rule; receivers reject `<transport-promise> <empty>`.
- **Alt-B:** Truncation removes the *outer* frames first, leaving inner content visible but with no chain-of-custody. Receiver thinks it has Alice's content with no transport context. **Easier:** content always survives partial transmission. **Harder:** chain-of-custody is the first thing lost; trust ledger has nothing to update against R because R's promise is gone. **New obligation:** receiver must treat truncated messages as low-trust regardless of inner content quality.
- **Alt-C:** Same wire as Alt-B but evaluation depends on inner-first verification. If outer frames are missing, receiver lacks transport context to weight the verified inner content. **Easier:** content verifies before transport — distinguishes "bad content" from "bad transport". **Harder:** same as Alt-B for chain-of-custody loss. **New obligation:** distinct rejection reasons for "content invalid" vs "context missing".
- **Alt-D:** Truncation removes whatever frames were last on the wire; receiver's projection just returns fewer frames in some categories. **Easier:** graceful degradation; receivers can still extract whatever survived. **Harder:** the receiver may not notice truncation at all if the projected categories happen to be present. **New obligation:** every message must declare its expected frame inventory in the innermost frame, so the receiver can detect missing frames.
- **Alt-E:** Inherits Alt-A failure semantics by default; Alt-D's degradation is available as policy. **Easier:** sites with high-loss links can choose graceful degradation; sites with strong integrity needs use peeling. **Harder:** the wire bytes are the same in both modes, so a degraded receiver cannot tell that a peeling receiver would have rejected the message.

### S3 — Concurrent actors / mixed-version nodes

Two routers along the path; one runs `promstack` v1, one runs v2 with a new assertion type. A frame uses an unrecognized assertion tag.

- **Alt-A:** v1 receiver peels v2 frame, doesn't recognize assertion type, must defer or reject. Choice: skip-and-continue (best-effort interop) or stop (strict). **Easier:** the peeling sequence makes the unrecognized frame immediately visible; the v1 receiver can flag it and pass it through to a v2 inner peer. **Harder:** if the v1 receiver is the *terminal* receiver, it cannot evaluate the v2 frame at all and must defer indefinitely. **New obligation:** every frame must carry a `criticality` flag — "skip me if you don't understand me" vs "fail closed if you don't understand me".
- **Alt-B:** Same evaluation as Alt-A; same trade-offs.
- **Alt-C:** Inner-first evaluation means the v1 receiver may ratify the inner content based on its limited understanding of the outer v2 frames. If the v2 frame is a revocation or compliance assertion, the v1 receiver acts on stale information. **Harder:** silent under-evaluation; v1 doesn't know what it doesn't know. **New obligation:** inner frames must be self-contained enough that ratification by outer frames is *additive*, never *constraining*. (This may be impossible for some assertion types, in which case Alt-C is wrong for them.)
- **Alt-D:** v1 receiver projects all frames it understands, ignores the rest. **Easier:** unknown frames don't block known-frame evaluation. **Harder:** silent skipping of unknown critical frames is the default. **New obligation:** same `criticality` flag as Alt-A, but checked before projection chooses to ignore.
- **Alt-E:** Assertion types declare in their pCID-spec whether they are "advisory" (skip if unknown) or "critical" (reject if unknown). Both peeling and projecting receivers honor this.

### S4 — Long-horizon evolution and migration

Year 1: messages have 3 frames typical. Year 50: messages have 8-frame stacks because new assertion types (capability tokens, ZK proofs from TE-12, jurisdiction tags) have proliferated. Old archives must still be readable.

- **Alt-A:** New assertion types are added at the outermost layer (most recent additions). Old messages have a short prefix of well-understood frames, and any "future-knowledge" frames sit outside. A year-1 receiver replaying a year-50 message can peel the recognizable inner frames and ignore the outer additions. **Easier:** old archives remain interpretable by old code. **Harder:** if a year-50 outer frame revokes a year-1 inner frame's authority, the year-1 replay reaches the wrong conclusion. **New obligation:** revocation frames must propagate backward into archives, not just forward.
- **Alt-B:** Same archive interpretability as Alt-A but worse during streaming because outer additions sit at the beginning of the wire bytes — old code reading the start of a year-50 message hits unfamiliar frames first.
- **Alt-C:** Old code verifies inner content first (which it understands) and then encounters outer ratification frames it may not understand. Inner-first evaluation makes archives interpretable. **Easier:** inner content is always first-class. **Harder:** outer frames carry the temporal context (when was this added, by whom?), and inner-first evaluation leaves that as background until end. Old code may act on inner content without knowing it has been superseded.
- **Alt-D:** Projection over a multiset means new assertion types simply add new categories. Old code projects over the categories it knows; new categories are invisible to it. **Easier:** maximal additivity. **Harder:** invisibility of new categories means old code cannot detect that something it should care about has happened.
- **Alt-E:** Combines Alt-A's archive interpretability with Alt-D's additivity. The `criticality` flag from S3 is the key mechanism — old code that encounters a critical-unknown frame must defer.

### S5 — Trust-boundary changes

A router that was trusted (high keep-ratio) is later discovered to have signed false transport-promises. The trust ledger updates. How do prior messages whose interpretation depended on that router's promises get re-evaluated?

- **Alt-A:** Outer transport-promises are visually "on top". When trust drops, the receiver can re-walk recent messages and downgrade their effective trust. The peeling order means the transport-promise is already at hand. **Easier:** mechanical re-walking of recent messages. **Harder:** historical messages may have already been acted on; re-walking only updates future evaluations, not past actions. **New obligation:** trust events must include a "re-evaluation horizon" (how far back to re-walk).
- **Alt-B:** Same as Alt-A.
- **Alt-C:** Inner-first evaluation already privileged inner content over outer transport. A trust drop in the transport may have *less* effect because the inner content's trust was already weighted independently. **Easier:** less re-evaluation churn on transport trust changes. **Harder:** the transport's role in surfacing the message at all is under-credited; a hostile transport that was trusted may have selectively suppressed or delayed messages, and inner-first evaluation has no record of that.
- **Alt-D:** Trust events update per-frame trust scores; projection re-runs against the updated scores. **Easier:** orthogonal per-frame updates. **Harder:** re-running projection over historical archives is unbounded work; needs caching or windowing.
- **Alt-E:** Choice of re-evaluation strategy is per-deployment. Trust events include hints about which re-evaluation strategy is appropriate.

### S6 — Scale (storage, bandwidth, CPU, operational complexity)

A high-throughput link with millions of messages per second. A long-horizon archive of trillions of messages. A streaming receiver with bounded memory.

- **Alt-A:** Streaming-friendly: receiver can decide-and-discard outer frames as they arrive without buffering. Bounded memory possible. **Easier:** wire-order matches evaluation-order; no seek. **Harder:** if the receiver is interested in the inner content but doesn't trust the outer transport, it has done work to reject the whole message before reaching the content — sometimes wasted, sometimes the right early-bail.
- **Alt-B:** Streaming-hostile: receiver must buffer the whole message to find outer frames at the end. Bounded memory requires a maximum-message-size policy. **Easier:** for batch processing where messages are read whole, both orders are equivalent. **Harder:** streaming and bounded memory.
- **Alt-C:** Streaming-friendly for inner content; outer ratification arrives later but receiver has already begun side-effect-free verification. **Easier:** parallelizable: inner verification on one core, outer ratification on another as bytes arrive. **Harder:** receivers must defer side-effects until outer ratification, which constrains the streaming pipeline. Effectively makes the "first useful work" point further from the "first byte arrives" point.
- **Alt-D:** Multiset semantics need at least one full pass before policy can be applied. Streaming is awkward unless predicates are cheap and incremental. **Easier:** order-independent encoders/decoders; transport-layer flexibility. **Harder:** memory cost is at least the message size; CPU cost is per-predicate-per-frame.
- **Alt-E:** Streaming receivers default to peeling (Alt-A semantics); batch receivers may prefer projection. Same wire bytes; different consumer choices.

## What each alternative makes easier, harder, and what it newly demands

| Alt | Easier | Harder | New obligation |
|-----|--------|--------|----------------|
| A | streaming, early-rejection on transport, archive interpretability, mechanical trust-event re-walking | inner content delayed; outer frames must be reliably append-only | routers must append, not prepend; criticality flag for unknown frames |
| B | none unique | streaming, bounded memory, chain-of-custody on truncation | length-prefix discipline; minimum-useful-prefix rule |
| C | inner-content-first evaluation; parallelizable verification | silent under-evaluation when outer frames revoke; harder trust-event re-walking | additivity-only constraint on outer frames |
| D | order-independence; graceful degradation; receiver-policy diversity | full-message inspection required; silent skipping default | per-frame assertion self-identification; criticality declaration; expected-inventory frame |
| E | accommodates both peeling and projecting receivers; future-flexible | doubled spec surface; subtle per-deployment behavior differences | clear declaration of conventional positions; criticality flag |

## Surviving alternatives

After scenario analysis:

- **Alt-B is rejected.** It carries Alt-A's evaluation costs and Alt-C's wire awkwardness with none of the wins of either. The only argument for it (matches some legacy serializations where payload comes first) is an aesthetic preference, not a design constraint.
- **Alt-A and Alt-E are nearly equivalent.** Alt-A is Alt-E with the projection mode disabled. Alt-E strictly dominates in optionality; the cost is doubled spec surface. Survives.
- **Alt-C survives but with a major caveat.** Inner-first evaluation requires the additivity-only constraint on outer frames (S3, S4). This constraint is plausibly impossible for revocation, capability-revocation, and compliance-assertion assertion types. So Alt-C may be the right answer for *some* assertion types but not all. This implies a per-assertion-type ordering policy rather than a global one — which leans toward Alt-E with declarative position conventions.
- **Alt-D survives as a complement, not an alternative.** Pure projection without ordering is too permissive (silent skipping by default) and too memory-hungry (full-message inspection). But projection-as-secondary-mechanism alongside peeling is exactly what Alt-E offers.

**Surviving alternatives: Alt-A and Alt-E.**

## Conclusions

1. **Wire order should be outermost-first.** Both Alt-A and Alt-E agree. Reasons: streaming-friendly, matches the temporal model (most recent additions visible first), enables early rejection on transport-promise without buffering, and aligns archive interpretability (old code reads inner content after current outer wrappers).

2. **Evaluation order should default to outermost-first peeling.** Both Alt-A and Alt-E agree. The peeling sequence is the cheapest correct evaluation in S1, S2, S6.

3. **Receivers should also have access to projection (`Project`).** This is what distinguishes Alt-E from Alt-A. The cost (slightly larger spec) is much smaller than the gain (per-assertion-type policy expressiveness; receiver-side flexibility for trust-event re-evaluation; clean handling of ZK and capability assertions that don't fit a strict ordering).

4. **There is no globally canonical frame placement.** The conventional positions are: transport-promise outermost (last-applied, by the inbound edge); content-promise innermost (first-applied, by the original author); signature wraps whatever it signs (varies by protocol). But these are *conventions*, not parser rules. A receiver that finds a signature in an unexpected position must still accept it — `Project(msg, signature-predicate)` will find it.

5. **Every frame must carry a `criticality` flag** (or its assertion-type-spec must declare criticality globally for that type). This emerges from S3, S4, and Alt-D's failure modes. Without criticality, a mixed-version receiver cannot distinguish "skip me if you don't understand me" from "fail closed".

6. **Per-assertion-type ordering policy is a real concept** that deserves its own DR. Alt-C's analysis shows that some assertions (revocation, compliance) demand outer-first to take effect; others (content authentication) work fine inner-first. The policy belongs to the assertion-type spec (the pCID), not to the wire format.

## Implications for the repo's open TODOs and pending DIs

- **OQ-1 (Promise structure / §1.1):** This TE adds two required fields:
  - `criticality: "advisory" | "critical"` per frame, OR per assertion-type spec
  - `position-convention: "outermost" | "innermost" | "any"` declared in the assertion-type spec, not enforced by the wire library
- **TE-12 (zero-knowledge promise frames):** Compatible with Alt-E — a ZK proof frame is just another promise type discovered by `Project`. The fact that it can sit anywhere in the stack is a feature.
- **TE-13 (time-traveling break-witness):** Compatible with Alt-E — break-witnesses are critical-criticality outer frames. The "re-evaluation horizon" the TE-13 conclusion will need is not in conflict with Alt-E's peeling/projection model.
- **§1.3 of specs/harness-spec-draft.md:** The four invariants (out-of-order handling, frame stripping, missing-body promises, loud failure on disagreement) all remain. This TE refines #1: out-of-order is handled by `Project`, not by attempted re-ordering at the wire level.
- **Future work (downstream of locking this TE):**
  - DR for "criticality flag location: per-frame or per-assertion-type-spec?"
  - DR for "per-assertion-type position convention: who declares it and where?"
  - DR for "minimum-useful-prefix rule for truncated messages"
  - DR for "expected-inventory frame: required, optional, or implicit from the innermost content-promise?"

## Decision Framing — questions for the user

After this TE, the surviving alternatives are Alt-A and Alt-E. Alt-E strictly dominates Alt-A unless the spec-surface cost is prohibitive.

**DF-1.1**: Adopt Alt-E (peeling + projection both first-class, with conventional but non-binding frame positions)?
- (a) Yes — adopt Alt-E. (Recommended.)
- (b) Adopt Alt-A only — peeling is mandatory, projection is a downstream library convenience but not part of the wire-spec semantics.
- (c) Defer; want to see one more scenario or alternative first.

**DF-1.2**: Where does the `criticality` flag live?
- (a) On every frame, declared by the promiser. Receivers honor it directly.
- (b) In the assertion-type spec (the pCID). Receivers look it up by assertion type. Frames carry only assertion-type pCID; criticality is global per type.
- (c) Hybrid: assertion-type spec declares default, frame may override.
- (d) Defer; needs its own TE before locking.

**DF-1.3**: Wire encoding direction — outermost-first byte order?
- (a) Yes — outermost-first on the wire (matches Alt-A and Alt-E).
- (b) Innermost-first on the wire, with length-prefix to enable seek-to-end.
- (c) Defer; this should be its own DR.

**DF-1.4**: Are conventional frame positions (transport outermost, content innermost) normative or merely descriptive?
- (a) Descriptive only — receivers must accept any position; `Project` is authoritative.
- (b) Normative for senders, accepting for receivers — senders must follow conventions, receivers must accept anything.
- (c) Normative for both — non-conventional placements are a protocol error.
- (d) Per-assertion-type — each assertion's spec declares whether its position is normative.

The recommended set is (1.1.a, 1.2.c, 1.3.a, 1.4.d). Reason: maximally additive, leaves room for assertion-type specialization, doesn't lock down what isn't yet known.

## Decision status

`needs DF` — awaiting user choice on DF-1.1 through DF-1.4. After DF, the locked decisions become DI entries in `protocols/wire-lab.d/TODO/TODO-20260429-164955-te-promise-stack-ordering.md`.
