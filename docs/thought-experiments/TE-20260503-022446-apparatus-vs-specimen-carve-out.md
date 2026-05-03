# TE-36: Apparatus vs. specimen — carving the harness-spec apart from the wire/envelope/ledger hypotheses it studies

**Decision status:** open (DFs awaiting decision).

**First-drafted:** 2026-05-03 02:24:46 UTC.

**Vocabulary note (added during drafting per DF-36.4 framing; revised after Steve's correction):** The harness-spec inherited the words "contest" / `contest-v1` from earlier draft text. Per Steve's standing rule (no "contest"; no "burden"; "discrepancy" rather than "defect") and the directive in DF-36.4 to use Promise Theory (PT) vocabulary, this TE recasts the discourse vocabulary in PT terms throughout. Initial drafting used `imposition` for the proposal primitive; Steve corrected this on 2026-05-03 — in PT canon (Burgess, Bergstra) an imposition is the non-cooperative form (it reduces trustworthiness assessments and is ill-aligned with receiver autonomy). The cooperative form for what proposals/discourse actually do is a **conditional promise** (also called a **reciprocal promise**): a promise whose body is contingent on another agent's promise. A proposer says "I promise X if you (signer) promise to merge." An endorser says "I promise to abide by P if it is merged." A contester says "I promise that observation O conflicts with P; I will keep this assessment under review." Endorsement and contestation collapse into **assessments** of polarity (positive / negative); a counter-proposal is another conditional promise that references a prior one; a hypothesis-result is an assessment of a hypothesis-promise (polarity: confirmed / refuted / unrelated). PT canon does not split endorse / contest into separate primitives — they are both assessments with different polarity. This vocabulary applies to the new pCIDs the carve-out produces. Historical references in older TEs (TE-14, TE-16, TE-29's prose, this TE's tabletop scenarios where actors quote past events) stay under TE-34 Cat-1b (historical-quotation paths, not swept). The frozen DI-003 contest-artifact at `proposals/pending/ppx-dr-001-bootstrap/contest-20260429-033208-steve-traugott.md` is immutable per the rules already in force.

**Decision under test:** The wire-lab harness-spec currently prescribes specific wire shapes (the promise-stack envelope, the `Promise` struct, the `[]Promise` message form, the `promstack` Wrap/Peel/Project library API, the `TrustLedger` field set, the named discourse pCIDs `endorse-v1`/`contest-v1`/`counter-propose-v1`/`hypothesis-result-v1`, the transport-promise as outermost frame). Per Steve's correction on 2026-05-02 / 2026-05-03, the harness-spec is **lab apparatus**: it describes how scenarios run, how actors are named, how transports are simulated, how messages are injected and observed, how trust ledgers are scored, and how outputs are compared. It is NOT prescriptive about wire shape. Each candidate envelope, library API, ledger shape, and discourse vocabulary is a **specimen** — a hypothesis under study — and each must live in its own `protocols/<slug>.d/` directory with its own spec doc and TEs. This TE locks the carve-out shape: how strictly the carve-out is applied, where each currently-named specimen relocates, and how ambiguous or apparatus-adjacent material is handled. It is the prerequisite for sweep edits to harness-spec under the locked TE-34 editing policy (Cat-1a path renames + Cat-2 vocabulary updates citing the DI this TE will produce).

**Source audit:** The classification table in [`protocols/wire-lab.d/docs/audit-20260503-015309-harness-spec-apparatus-vs-specimen.md`](../../protocols/wire-lab.d/docs/audit-20260503-015309-harness-spec-apparatus-vs-specimen.md) (committed 2026-05-03 at `4725b3e`) is this TE's primary input. It walks every section of `harness-spec-draft.md` and classifies as Apparatus / Specimen / Ambiguous with brief rationale.

**Constraints in scope:** C-1 (no central registry), C-2 (multi-generational durability), C-4 (protocol forking is normal), C-5 (trust accrues per-burden), C-6 (signing key is the only structural lock). The carve-out exists precisely to honor C-4: making each envelope/ledger/discourse hypothesis its own forkable protocol object. Conventions or apparatus that pretend a single envelope is canonical violate C-4 silently.

## Background

The wire-lab was originally drafted as one document (`x/wire-v2/...` then promoted to `protocols/wire-lab.d/specs/harness-spec-draft.md`) under the assumption that one promise-stack envelope was the working hypothesis and the harness existed to validate it. That assumption was never explicit but pervaded the prose: §1.1 prescribes the `Promise` struct, §1.2 argues for promise-stack vs. fixed envelopes, §2.1 prescribes the `TrustLedger` shape, §7.1 prescribes the transport-promise as outermost frame, §10a.2/.3/.6 names `endorse-v1`/`contest-v1`/`counter-propose-v1`/`hypothesis-result-v1` as required wire-level vocabulary.

Steve's correction this session is that the wire-lab studies **multiple** hypotheses — at minimum the promise-stack family AND the `grid([pcid, payload])` family AND likely future ones — and the harness-spec must not silently privilege one of them. The point of the harness is to compare them; if the harness-spec prescribes one, the comparison is rigged.

The earlier TE-29 (protocols-as-simrepos) and TE-32 (spec-vs-implementation split) already laid the directory shape: each protocol owns a `protocols/<slug>.d/` simrepo with its own specs/, docs/thought-experiments/, TODO/, DR/ inline, and CHANGELOG.md; implementations live under top-level `implementations/<impl-name>/`. TE-30 locked per-protocol TODO/ directories. TE-31 inverted the doc/repo reference direction (spec-doc-CID upstream, CHANGELOG entries downstream). What TE-29-32 did NOT do is migrate the specimen-bearing material out of the harness-spec; the migration is deferred to whenever the wire-lab carves up its envelope hypotheses, which is now.

The just-finished TODO 014 created `protocols/wire-lab.d/`, `protocols/group-session.d/`, `protocols/ppx-dr.d/`, and `protocols/udp-binding.d/`. None of them is the right home for the promise-stack envelope hypothesis: wire-lab is the apparatus; group-session is a transport (and currently uses a `grid <pcid>` carrier per TE-24); ppx-dr is the proposal protocol; udp-binding is an L4 binding. A new `protocols/promise-stack.d/` (or analogous) is required if the promise-stack hypothesis is to have an apparatus-clean home.

## Why this is hard

1. **Apparatus-vs-specimen is a category boundary, not a syntactic one.** A sentence like "the trust ledger updates if the transport keeps or breaks its promise" is half apparatus (the harness records keep/break) and half specimen (the trust ledger is one shape; "the transport's promise" presumes the promise-stack envelope). Reasonable bots and humans will draw the line in different places.

2. **Some apparatus claims look like specimen claims unless reframed.** "Edges carry promises, not bytes" (§7.1) sounds prescriptive but the apparatus-level claim is "edges annotate inbound traffic with observed metadata sufficient for any candidate envelope to construct its outermost frame." The carve-out has to do work to keep apparatus-level invariants while removing envelope-specific prescriptions.

3. **The promise-stack hypothesis has no protocol directory yet.** Creating one is itself a substantive act: a slug must be chosen, a draft spec doc must be stubbed (even if minimal), TODO/CHANGELOG/MANIFEST scaffolding installed. The TE has to lock not just the carve-out principle but the destination directories.

4. **The harness-spec's "Promise" / "promise-stack" prose is also load-bearing for several existing TEs.** TE-1 is "promise-stack ordering"; TE-12 is "promise-stack as zero-knowledge envelope"; TE-13 is "time-traveling break-witness" and assumes break-witness messages are promise-frames; TE-21 (now TE-24's older drafting) pre-dates the carve-out. These TEs may need to migrate WITH the promise-stack carve-out, or be left where they are with footers pointing at the new home, depending on how strict the carve-out is.

5. **C-4 cuts both ways.** "Protocol forking is normal" means the wire-lab must support multiple competing promise-stack variants, multiple competing grid-envelope variants, and multiple competing trust-ledger shapes — all of them simultaneously. The carve-out shape must not implicitly forbid forks within a candidate family (e.g. promise-stack-v1 and a community-forked promise-stack-alt should both fit cleanly).

## Tabletop scenarios

Six concrete scenarios with named actors. Each scenario tests one aspect of how a carve-out alternative behaves in practice.

### Scenario S1: Alice the new contributor reads the spec corpus today

Alice arrives at `github.com/promisegrid/wire-lab` with no context. She opens README.md, follows the link to `protocols/wire-lab.d/specs/harness-spec-draft.md`, reads §1.1, and concludes "the wire format is `[]Promise`, encoded as a CBOR array, accessed via `promstack.Wrap`/`Peel`/`Project`." She begins implementing. Three weeks later she discovers TE-24 ("group-transport envelope: `grid <pcid>` carrier") and is confused: which is canonical?

**Test:** under each carve-out alternative, does Alice still come away with a single canonical envelope after reading harness-spec? She must NOT — the harness studies multiple. The alternative is correct only if Alice instead concludes "there are multiple candidate envelopes; the harness compares them" and follows links to each candidate's protocol directory.

### Scenario S2: Bob the second-implementer attempts a clean-room rebuild

Bob (per §9 "Realism Suggestions") writes a second `promstack` implementation in a different language from a separate reading of the spec. He needs an authoritative spec-doc CID to hash and reference in his B-side `implementations/bob-promstack/CHANGELOG.md` per TE-32. Today the only candidate is `harness-spec-draft.md` itself, which forces him to claim conformance to the entire harness — including ToC scenarios, ingress models, agent profiles he does not implement.

**Test:** under each carve-out alternative, can Bob's CHANGELOG entry be tightly scoped — claiming conformance only to the promise-stack envelope spec, not to "the wire-lab"? The alternative is correct if and only if the promise-stack hypothesis has its own doc-CID Bob can reference.

### Scenario S3: Carol the second envelope hypothesis

Carol arrives with a proposal for an alternative envelope: `[pCID, salted-canonical-CBOR-payload]` with mandatory canonicalization rules and an optional outer signature wrapper. She wants to file it as a peer hypothesis to the promise-stack family, run her envelope through the same C1-C7 ToC scenarios, and see how it scores.

**Test:** under each carve-out alternative, where does Carol's spec doc live? Does she have to fork the harness-spec? Does she merely add a sibling `protocols/<carol-envelope>.d/`? Does she edit one shared registry-of-envelopes? The alternative is correct if Carol's act is additive and does not require touching files Steve owns alone.

### Scenario S4: Dave the promise-stack forker

Dave likes the promise-stack family but wants a variant: same wire shape, different acceptance semantics (e.g. mandatory innermost-signature instead of TE-1's still-open ordering). Per C-4, this is a normal fork.

**Test:** under each carve-out alternative, can Dave fork only the promise-stack hypothesis without re-stating the harness or the trust ledger? The alternative is correct if `protocols/promise-stack.d/` (or wherever the hypothesis lives) is the unit of fork, not the harness.

### Scenario S5: Mallory the bad-faith carve-out attacker

Mallory submits a PR claiming "the apparatus-vs-specimen carve-out makes §1.1 obsolete; here is a deletion of §1.1 with no replacement." Without the carve-out, §1.1 is the entire definition of the wire format and deleting it leaves new contributors with nothing.

**Test:** under each carve-out alternative, does the carve-out itself force the prerequisite (the destination protocol directory must exist with the migrated content) before the harness-spec sweep is allowed to run? The alternative is correct if the editing-policy DI it produces forbids harness-spec deletions whose target directory has not yet absorbed the content.

### Scenario S6: Ellen the 30-years-later contributor

Ellen arrives in 2056 with no living mentor. She finds `harness-spec-draft.md` (or its successor at some pCID), reads it, and tries to understand what PromiseGrid's wire format actually is. The harness has been forked and re-forked many times; multiple promise-stack variants exist; multiple grid-envelope variants exist; some have been retired.

**Test:** under each carve-out alternative, does Ellen's reading of harness-spec lead her to "look at the live `protocols/<various>/` to see which envelopes are still under study"? Or does it leave her assuming whatever is in §1.1 is still canonical? The alternative is correct if the harness-spec is honestly silent on envelope identity — Ellen has to look at the protocol directories to learn what the live specimens are.

## Decision points

Each DF is described in prose, alternatives enumerated, the bot's recommendation called out per Steve's standing rule. Steve responds "yes" to confirm the recommendation or names an alt to override.

### DF-36.1 — Carve-out strictness

The first decision is how aggressively the carve-out runs. Three plausible strictnesses, ordered from loosest to tightest.

**Alt-1.A: Strict carve-out.** Harness-spec mentions no candidate envelope, library API, struct shape, or named pCID by name. All eight specimen-bearing items from the audit relocate. Apparatus-level invariants stay, but they are stated envelope-agnostically: "any candidate envelope must support out-of-order receiver handling," not "the receiver consumes promises top-down." The reader of harness-spec leaves with no concrete envelope in mind — they must follow links to per-protocol specs to learn what the live specimens are. S1, S3, S4, S6 cleanly pass; S2 cleanly passes once `protocols/promise-stack.d/` exists; S5's attack is structurally prevented because there is nothing left in §1.1 to delete after the carve-out runs.

**Alt-1.B: Mixed / worked-example.** Harness-spec uses ONE named candidate (the promise-stack family) as a worked example throughout, with an explicit top-of-section banner: "the following is one candidate hypothesis under study; see `protocols/promise-stack.d/` for the live spec, and see §8 TE-? for peer hypotheses." Apparatus-level invariants stay alongside the worked example. The reader leaves with the promise-stack hypothesis well-understood and the awareness that other hypotheses exist. S2 passes weakly (Bob's CHANGELOG must point at `protocols/promise-stack.d/specs/...`, not the harness); S3 weakly passes if Carol's hypothesis lives in a sibling directory but reads as "the alternate"; S4 passes; S6 fails Ellen if the worked example outlasts the live promise-stack hypothesis (the example becomes stale-canon).

**Alt-1.C: Registry-of-knobs.** Harness-spec lists every studied dimension (envelope, ledger shape, discourse vocabulary, ordering rule, criticality model) as a knob with a placeholder; candidate values fill placeholders elsewhere. Closest to TE-29's directory shape applied to apparatus-vs-specimen: harness-spec is itself a registry, not a doc. The prose recasts every concrete shape as `<envelope-pcid>`, `<ledger-pcid>`, `<discourse-pcid>`. Reader of harness-spec leaves with a complete picture of the dimensions being studied but no concrete shape anywhere. S1, S3, S4, S6 pass cleanly; S2 passes; S5 cannot attack because the harness-spec has nothing to delete. Cost: substantial rewrite, much higher than Alt-1.A; risk that the registry abstraction is itself a guess about the future.

**Bot recommendation: Alt-1.A (Strict).** Alt-1.B silently privileges the promise-stack family and re-creates the problem TE-36 exists to solve, just with a banner that readers will skim past. Alt-1.C is structurally the cleanest but is a year of work and presumes knob shapes we have not yet seen used in anger. Alt-1.A is the minimum carve-out that satisfies all six tabletop scenarios and is the lightest sweep through the existing harness-spec text. The Cat-2 vocabulary updates per TE-34 editing policy are well-understood for Alt-1.A: replace `Promise` / `[]Promise` / `promstack` references with envelope-agnostic apparatus phrasing, link to per-protocol directories for the candidates. Alt-1.A also pairs cleanly with the per-DF migration decisions below — each specimen-bearing item gets one concrete destination, no shared abstraction layer required.

### DF-36.2 — Promise-stack hypothesis home

Where does the promise-stack family live? §1.1 prescribes the `Promise` struct, the `[]Promise` shape, the CBOR array encoding, and the `promstack` library (Wrap/Peel/Project). §1.2 argues for it. §1.3 lists tests over it. §2.2 break-witness is a promise-frame. §7.1 transport-promise is the outermost promise-frame. TE-1, TE-12, TE-13 reference it. None of the four protocol directories created by TODO 014 is the right home for it.

**Alt-2.A: New `protocols/promise-stack.d/`.** Standard simrepo shape per TE-29. Spec doc at `protocols/promise-stack.d/specs/promise-stack-draft.md` carries the migrated §1.1/.2/.3 plus the Promise struct + promstack API. TODO/, docs/thought-experiments/ (TE-1 relocates), MANIFEST.md, CHANGELOG.md scaffolding installed. The slug is meaningful (matches the existing `promstack` library name) and forward-compatible with future variants (`promise-stack-v2.d/`, `promise-stack-zk.d/` for TE-12's ZK variant).

**Alt-2.B: New `protocols/wire-envelope.d/` neutral umbrella.** A single directory holds all envelope hypotheses (promise-stack, grid, future Carols) as sub-specs. Promotes apparatus-vs-specimen at one level higher (one envelope-comparison protocol with sub-hypotheses). Cost: one more layer of indirection; conflicts with TE-29's "let each protocol name its own internals" because different envelopes will want different internal shapes; collapses Carol and Dave's forks into PRs against one shared directory rather than independent forks (violates C-4 in spirit).

**Alt-2.C: Defer creation; treat §1.1-§1.3 as orphan text.** Leave the promise-stack prose in harness-spec under a "this material will migrate" banner; do not create the destination directory until someone actively works on the hypothesis. Cost: harness-spec stays muddled for an unknown period; S1/S2/S3/S4/S6 keep failing until the deferral is undone; the carve-out TE produces no concrete state change.

**Bot recommendation: Alt-2.A.** The slug `promise-stack` is meaningful (already used by the library name `promstack`), forward-compatible with variants and forks, and matches TE-29's per-protocol-simrepo shape. Alt-2.B over-couples competing hypotheses and silently violates C-4. Alt-2.C punts the carve-out without solving anything. Creating `protocols/promise-stack.d/` is a single skeleton commit (mirrors what TODO 014 did for the other four protocol directories) and unblocks the harness-spec sweep.

### DF-36.3 — Grid-envelope hypothesis home

`grid([pcid, payload])` is Steve's "we *think* it's the right shape" working hypothesis. Currently `grid <pcid>` appears as a textual carrier in `protocols/group-session.d/specs/group-session-draft.md` per TE-24, but the envelope-shape claim is broader than the group-session transport: `grid([pcid, payload])` is also a candidate for non-group-session transports, including the eventual canonical PromiseGrid wire format.

**Alt-3.A: New `protocols/grid-envelope.d/`.** Parallel sibling to `protocols/promise-stack.d/`. The spec doc holds the `grid([pcid, payload])` envelope hypothesis as its own object, transport-agnostic. The group-session spec then references the grid-envelope spec as the envelope it carries (per TE-32 inversion: implementation/use-side references upstream spec-CID).

**Alt-3.B: Folded into `protocols/group-session.d/`.** The grid envelope is currently only used by group-session, so leave it there until a second transport adopts it. Cost: when the second transport (UDP-binding-running-grid? a hypothetical broadcast-over-MQTT?) adopts grid, the envelope spec has to migrate out anyway, and the migration breaks group-session's existing references. Also conflates transport (group-session is one transport) with envelope (grid is one envelope), which violates the transport-≠-envelope distinction Steve stated this session.

**Alt-3.C: New `protocols/pcid-payload.d/` with a more neutral slug.** The slug `grid` is overloaded with PromiseGrid the project; `pcid-payload` is shape-descriptive. Cost: humans and existing TEs already use the word "grid" for this envelope (per TE-24/26/27 vocabulary); a slug rename would force a Cat-2 vocabulary sweep across the existing TE corpus.

**Alt-3.D: Defer.** Leave the envelope shape inside group-session until a second transport demands separation.

**Bot recommendation: Alt-3.A.** Honors the transport-≠-envelope distinction Steve made this session. The grid envelope is a hypothesis on its own merits and may be adopted by transports that aren't group-session (per Steve's "we need to be able to experiment with different message envelopes" — the experiment is over the envelope axis, not the transport axis). The slug `grid-envelope` is unambiguous within the wire-lab context (no clash with `grid-poc`, no clash with PromiseGrid). Alt-3.B silently forecloses the transport-agnostic experiments. Alt-3.C creates vocabulary churn for marginal gain. Alt-3.D punts. After Alt-3.A lands, group-session's spec doc relocates its envelope claims to "carries the envelope specified at `protocols/grid-envelope.d/specs/grid-envelope-draft.md`" via TE-31's spec-doc-CID reference convention.

### DF-36.4 — Apparatus-level vocabulary for proposals/discourse (PT-recast)

§10a.2 names a `proposal-checklist-v0` pCID. §10a.3 names `endorse-v1`, `contest-v1`, `counter-propose-v1`. §10a.6 names `hypothesis-result-v1`. These are wire-level discourse vocabulary embedded in harness-spec. Per the vocabulary note at the top of this TE, the new pCID names will use Promise Theory primitives — `conditional-promise` (also called `reciprocal-promise`: a promise whose body is contingent on another agent's promise) and `assessment` (with polarity parameter) — rather than the legacy endorse/contest pair, and explicitly **not** `imposition` (which in PT canon is the non-cooperative form that reduces trustworthiness).

Why conditional/reciprocal rather than imposition: a proposal in the discourse is cooperative — the proposer is offering to keep some promise *if* the deciding signer (or community) accepts. PT models this as a conditional promise of the form "I promise X if you promise Y" (canonical notation: A1 — +X|Y → A2, where Y is the prerequisite promise from A2). An imposition, by contrast, is one agent attempting to induce behavior in another without a reciprocal commitment — the wrong shape for cooperative protocol design and explicitly flagged in PT literature as trust-eroding.

The PT mapping the carve-out adopts:

| Legacy harness-spec name | PT-vocabulary recast | One-line semantic |
|---|---|---|
| `proposal` (§10a.2) | `conditional-promise` (reciprocal) | Q promises to abide by change X if signer promises to merge |
| `proposal-checklist-v0` | `conditional-promise-checklist` (community convention) | Suggested prose shape for conditional promises in this discourse |
| `endorse-v1` (§10a.3) | `assessment` (polarity: positive) | Q assesses conditional-promise P as one Q would also keep |
| `contest-v1` (§10a.3) | `assessment` (polarity: negative) | Q assesses conditional-promise P as conflicting with Q's observations |
| `counter-propose-v1` (§10a.3) | `conditional-promise` (with `references:` pointer) | A new reciprocal promise explicitly referencing a prior one |
| `hypothesis-result-v1` (§10a.6) | `assessment` (of a hypothesis-promise; polarity: confirmed / refuted / unrelated) | The harness reports its assessment of a hypothesis the agent promised |

The `conditional-promise` and `assessment` primitives collapse what were five named pCIDs into two parameterized PT primitives. This is more honest about what the discourse actually is in PT terms: proposals are reciprocal promises ("I will accept X if you sign"), not impositions on the signer's autonomy; endorsements and contests are both polarities of the same assessment primitive.

**Alt-4.A: Move PT-recast vocabulary to `protocols/ppx-dr.d/`.** ppx-dr is already created (TODO 014) for "proposals as messages on a transport." It is the natural home for the conditional-promise + assessment vocabulary. ppx-dr's spec doc absorbs both primitives as the discourse layer carried by the protocol.

**Alt-4.B: New `protocols/discourse.d/` separate from ppx-dr.** Conditional-promise + assessment is a distinct pattern from "messages on a transport"; ppx-dr could specialize to the transport while a separate `discourse.d/` carries the PT primitives generically. Cost: more directories; the two protocols are tightly coupled in practice.

**Alt-4.C: Stay in harness-spec, PT-recast in place.** Harness-spec declares the PT primitives at apparatus level ("the harness studies systems in which agents make conditional promises and emit assessments"). Cost: conflicts with DF-36.1 Alt-1.A Strict — apparatus mentions PT primitives but specific wire encodings still constitute specimen prescription if left in harness-spec.

**Alt-4.D: Defer decision.** Leave §10a.2/.3/.6 as the only home for the vocabulary until ppx-dr's spec is fleshed out.

**Bot recommendation: Alt-4.A with PT-recast vocabulary (conditional-promise + assessment).** ppx-dr exists as a per-protocol simrepo; PT primitives fit its scope (conditional promises are proposals on the wire; assessments are receiver-side polarity-bearing observations). Apparatus-level invariant retained in harness-spec: "the harness records conditional promises and assessments per PT canon; specific wire encodings are specimen-level and live in `protocols/ppx-dr.d/`."

### DF-36.5 — §1.3 simulator-tests-about-layering disposition

§1.3 lists four simulator tests phrased in promise-stack vocabulary: out-of-order promise stacks, forwarding-node promise-stripping, promise about a missing inner body (merkle-reference), two agents disagreeing about which promise frame to evaluate first. The audit memo flagged this section as ambiguous: each item is partly an apparatus claim ("the harness exercises any candidate envelope under X scenario") and partly a specimen claim (presumes promise-stack vocabulary).

**Alt-5.A: Reframe at apparatus-level.** Rewrite §1.3 envelope-agnostically: "The harness exercises any candidate envelope under (1) out-of-order receiver handling, (2) forwarding-node intermediate-frame stripping, (3) deferred-body / fetch-elsewhere references, (4) ordering disagreement between two receivers; each candidate envelope's behavior under these scenarios is reported as part of its acceptance results." The §1.3 content stays in harness-spec as apparatus.

**Alt-5.B: Migrate with §1.1.** Move §1.3 verbatim to `protocols/promise-stack.d/specs/...`; harness-spec drops the section entirely. Apparatus equivalents (the four scenarios stated envelope-agnostically) appear elsewhere in harness-spec or in a new section.

**Alt-5.C: Both.** Apparatus-level summary stays in harness-spec under a new title (e.g. "Apparatus-level layering-test scenarios"); specimen-specific details (promise-frame vocabulary, top-down accept/defer/reject semantics) migrate to promise-stack spec.

**Bot recommendation: Alt-5.C.** Alt-5.A loses information that promise-stack readers want (the specimen-specific tests are real and testable). Alt-5.B leaves a gap in harness-spec where the apparatus-level claim should be. Alt-5.C honors both: apparatus-level invariants (every candidate envelope must survive these four scenarios) live in harness-spec; specimen-specific test details live with the specimen. This pattern generalizes to other ambiguous sections (§3.3 bullet 1, §7.1).

**LOCKED: Alt-5.C Both** (Steve, 2026-05-03). Apparatus-level summary of the four layering-test scenarios stays in harness-spec under an envelope-agnostic title; specimen-specific detail (promise-frame vocabulary, top-down accept/defer/reject semantics) migrates to the promise-stack spec. Specimen-side migration is gated on OQ-36.6 — if promise-stack proves redundant with grid-pcid-payload, the specimen-detail target shifts accordingly, but the apparatus-level summary in harness-spec stands either way.

### DF-36.6 — TrustLedger struct shape disposition

§2.1 prescribes a specific `TrustLedger` field set: `first_seen_ns`, `interactions`, `kept`, `broken`, `evidence_chain`, `open_promises`, `score`, `score_components`, `reputation_imports`, `relationship_age_ns`, `last_drift_ns`. The audit flagged this as half-apparatus / half-specimen: existence of a per-peer trust ledger is apparatus; the specific field set is one candidate shape.

**Alt-6.A: Lift to apparatus-level invariants only.** Harness-spec drops the struct, retains the apparatus-level claim: "any candidate trust ledger must (i) record kept/broken/open evidence per peer pair, (ii) support per-assertion-type vector scoring per C-5, (iii) support reputation imports as down-weighted third-party promises." No struct, no field set — the invariants alone. No new protocol directory needed yet.

**Alt-6.B: Migrate struct as-is to `protocols/trust-ledger.d/`.** New protocol directory (parallel to promise-stack and grid-envelope). The struct becomes one candidate shape; future variants (trust-ledger-eigen, trust-ledger-bayesian) live as siblings. Apparatus invariants stay in harness-spec; struct migrates.

**Alt-6.C: Both.** Apparatus invariants in harness-spec; struct in `protocols/trust-ledger.d/` as one candidate shape; harness-spec links to the directory for live specimens. Same pattern as DF-36.5 Alt-5.C.

**Alt-6.D: Lift to apparatus invariants; create the protocol directory only when a second candidate shape is filed.** Lazy carve-out: keep the apparatus invariants only, do not create `protocols/trust-ledger.d/` until needed. The struct from §2.1 becomes "an example shape worth further study" recorded in the audit, not a live spec doc.

**Bot recommendation: Alt-6.D.** Trust-ledger shape is less developed as a hypothesis than promise-stack or grid-envelope; only one struct exists, no peer hypotheses are queued, the struct itself was drafted as illustrative rather than load-bearing. Creating `protocols/trust-ledger.d/` now would be busywork — it would hold one example shape and no substantive content. The apparatus invariants (record kept/broken/open per peer pair; per-assertion vector; imports as down-weighted) are sufficient to constrain harness behavior without naming a struct. When a second candidate shape arrives (e.g. a Bayesian-only ledger that drops `evidence_chain`, or an Eigen-only ledger that drops local scores entirely), `protocols/trust-ledger.d/` materializes and the §2.1 struct migrates with the new candidate as company. Until then, harness-spec stays apparatus-clean and no empty directory exists.

**LOCKED: Alt-6.D Lazy carve-out** (Steve, 2026-05-03). Lift §2.1 to apparatus-level invariants only in harness-spec; do not create `protocols/trust-ledger.d/` yet. The §2.1 struct field set is recorded in the audit memo as an example shape worth further study; sweep step 4 will not migrate it. The directory materializes when a second candidate shape is filed.

### DF-36.7 — §10 grid-poc directory table disposition

§10 gives a directory table mapping wire-lab artifacts to `x/wire-v2/promstack/`, `x/wire-v2/trust/`, `x/sims/simkit/...`, `x/sims/scenarios/...`, etc. After the carve-out, the rows for `promstack` and `trust` (and any specimen-specific row) become **implementations** (per TE-32 B-side `implementations/<impl-name>/`) rather than apparatus components.

**Alt-7.A: Split the table — apparatus rows stay; specimen rows migrate.** The `promstack` row migrates to `protocols/promise-stack.d/`'s spec doc as "reference Go implementation lives at `implementations/<impl-name>/...`". Same for `trust` if Alt-6.B/C is chosen; otherwise `trust` row is removed entirely under Alt-6.A/D. Apparatus rows (`simkit`, `scenarios`, `viewer`, `realruntimes`, `reports`) stay in harness-spec.

**Alt-7.B: Delete the table entirely.** Per TE-32, each protocol's CHANGELOG names its implementations; the harness-spec doesn't need a directory table at all. Cost: loses a useful at-a-glance map of where things live; readers have to walk all CHANGELOG files to reconstruct it. Mitigation: harness-spec links to a programmatically-generated implementation-map under `tools/`.

**Alt-7.C: Keep table as-is with an explanatory note.** Add a top-of-table note: "specimen rows below are example implementations of candidate hypotheses; see per-protocol CHANGELOGs for the live conformance set." Cost: silently outdated as new candidates land; pretends the table is canonical when it is convention.

**Bot recommendation: Alt-7.A.** Cleanest split that preserves the apparatus rows (which are useful as a fast orientation aid for anyone reading harness-spec) while moving specimen rows to where they belong per TE-32. The `promstack` row gains a destination once `protocols/promise-stack.d/` exists (DF-36.2 Alt-2.A); the `trust` row can be deferred under DF-36.6 Alt-6.D until needed. Alt-7.B is more architecturally pure but costs orientation aid for marginal gain. Alt-7.C silently outdates.

## Summary of recommended alternatives

| DF | Recommended Alt | One-line rationale |
|---|---|---|
| 36.1 | Alt-1.A (Strict) | Minimum carve-out that satisfies all six tabletop scenarios; lightest sweep. |
| 36.2 | Alt-2.A (`protocols/promise-stack.d/`) | Meaningful slug, fork-friendly, matches TE-29 simrepo shape. |
| 36.3 | Alt-3.A (`protocols/grid-envelope.d/`) | Honors transport-≠-envelope distinction; transport-agnostic. |
| 36.4 | Alt-4.A (move discourse vocabulary to `protocols/ppx-dr.d/`) | Vocabulary fits ppx-dr's scope; ppx-dr already exists. |
| 36.5 | Alt-5.C (Both apparatus + specimen) | Preserves apparatus invariants AND specimen-specific tests in their right homes. |
| 36.6 | Alt-6.D (Apparatus invariants only; lazy directory creation) | No empty directory; apparatus invariants suffice for now; trust-ledger directory materializes when a second candidate arrives. |
| 36.7 | Alt-7.A (Split — apparatus stays, specimen rows migrate) | Preserves orientation aid; honors TE-32 B-side. |

## Open questions surfaced (OQ-36.x)

- **OQ-36.1 — Which existing TEs migrate with the promise-stack carve-out?** TE-1 (promise-stack ordering) clearly belongs in `protocols/promise-stack.d/docs/thought-experiments/`. TE-12 (promise-stack as ZK envelope) likely too. TE-13 (time-traveling break-witness) presumes promise-frame vocabulary but is broader. Migration discipline (filename retention per TE-25, Cat-1a path-rename sweeps per TE-34) needs a concrete checklist before TODO 5 reframes.

- **OQ-36.2 — Apparatus-level break-witness mechanism without naming a frame.** §2.2 requires a break-witness mechanism. Apparatus-level reframe needs phrasing that does not presume promise-frame vocabulary while still capturing the recursive-trust property (recipient applies own ledger to witness before believing). Sketch: "the harness studies break-witness mechanisms in which observation-of-broken-promise becomes itself an artifact subject to receiver-side trust evaluation."

- **OQ-36.3 — Apparatus-level transport-edge annotation.** §7.1's "edges carry promises" needs apparatus-level reframe: edges annotate inbound traffic with observed metadata sufficient for any candidate envelope to construct its outermost frame. Concrete metadata set (source endpoint, time, integrity class, liveness) is uncontroversial; how that metadata is exposed to handlers depends on candidate envelope; harness-spec stays at the metadata level.

- **OQ-36.4 — Forward-pointer hygiene during sweep.** Per TE-34, Cat-3 navigational additions (forward-pointers from old TEs to the new TE-36 / DI / target protocol) are append-only `## Refinements` entries. Sweep must ensure every TE that referenced specimen-level material in harness-spec gains a Refinement pointing at the new home. Cheap with grep; checklist needed.

- **OQ-36.6 — Is promise-stack actually distinct from grid-pcid-payload, or a misreading of nested messages?** Steve flagged on 2026-05-03 that he suspects `promise-stack` (the `[]Promise` envelope hypothesis from harness-spec §1.1) is over-complicated and may have been invented by the bot based on a misunderstanding of how nested messages work inside `grid-pcid-payload` (where any `<payload>` slot can itself contain a `grid <pcid> <inner-payload>` recursively, with the per-protocol rules at each level coming from each pCID's own spec). Under that reading, the recursive structure that promise-stack tries to express in the envelope is already available for free in grid-pcid-payload's payload-recursion, and the `Promise` struct + `Wrap`/`Peel`/`Project` library API is invented machinery solving a problem that grid-pcid-payload does not have. Deferred for separate investigation. The DF-36.2 lock (`protocols/promise-stack.d/`) stands provisionally; if a follow-on TE concludes promise-stack is a degenerate case of grid-pcid-payload, the directory either retires (CHANGELOG `deprecates` entry) or never gets populated past skeleton. The harness-spec sweep in step 4 of the 6-step plan should NOT migrate §1.1 / §1.2 / §1.3 content into `protocols/promise-stack.d/specs/` until this OQ is resolved — the carve-out creates the directory but leaves it minimal pending the investigation.

- **OQ-36.5 — `transport-spec-draft.md` companion audit.** Audit memo flagged `protocols/wire-lab.d/specs/transport-spec-draft.md` (≈9 KB) as needing the same apparatus-vs-specimen audit before step 4 sweep. Likely lighter than the harness-spec audit (transport-spec was carved at TE-26/27 to be apparatus-level for transport-protocol differentiation), but unverified. Either folds into TODO 014's leftover sweep or files a new follow-on TODO.

## Refinements

(Reserved for forward-pointer / back-pointer additions per TE-34 Cat-3 policy.)

## Source

This TE is the formal carve-out decision corresponding to step 2 of the 6-step corrected plan from session 2026-05-02. Step 1 (audit) is the input at [`protocols/wire-lab.d/docs/audit-20260503-015309-harness-spec-apparatus-vs-specimen.md`](../../protocols/wire-lab.d/docs/audit-20260503-015309-harness-spec-apparatus-vs-specimen.md). Steps 3-6 (per-protocol directory creation, harness-spec sweep, TODO 5 reframe, parallel TODO for grid envelope) are gated on the DI this TE produces.
