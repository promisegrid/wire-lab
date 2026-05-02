# TE-34: TE editing policy and the TE corpus as one document with facets

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-20260502-212810

## Decision under test

Two related questions about how the TE corpus is treated as a body of writing:

1. **Editing policy.** Under what conditions may a TE's contents be modified after it is first committed? The repo currently has uneven precedent: TE-24 and TE-26 were rewritten in place when vocabulary changed (with top-of-file notes); other TEs have been left strictly alone. The bot has at times claimed there is a blanket "do not back-edit TEs" rule, but no such rule actually exists in any AGENTS file or in `docs/thought-experiments/README.md`. The only locked invariant is from TE-25: TE filenames are immutable because the timestamp slug is the content-address anchor that pins the integer alias.

2. **Reading discipline.** Should the bot read each TE as a self-contained artifact, or should it read the TE corpus holistically — as one large multi-faceted thought experiment in which any new TE can refine, contradict, or extend any earlier TE? The user has now stated explicitly: *"don't be afraid to view the entire corpus of TEs as one big TE with lots of facets — be holistic."* This TE locks that framing and works out the consequences for how TEs are written, edited, and cited.

These are joined here because the answers compose: if the corpus is one document with facets, then editing policy must say what it means to "refine" a facet versus to "rewrite" one, and what kind of marker the corpus needs to remain navigable.

## Assumptions

- TE filenames are immutable. This is locked by TE-25's "timestamp slugs are content-addressable identifiers" and "first-drafted-timestamp anchors the integer." Renaming a TE file would re-key the corpus; we do not do that.
- TE integer aliases (TE-1, TE-2, …) are also stable per TE-25, anchored on first-drafted timestamp. Renumbering happens only when two TEs collide on the same integer; the younger one moves.
- There is no general rule that TE *contents* are immutable. The two existing precedents (TE-24, TE-26) explicitly rewrite contents in place with top-of-file notes, and TE-27 lists those rewrites as part of its own conclusions.
- The contest-artifact immutability rule (DI-003: `proposals/pending/<branch>/contest-*.md` and `review-*.md` are immutable) applies to adversarial review records. TEs are not adversarial; they are the author's own thinking. The rationale for contest immutability does not transfer.
- TE-30 made the harness itself a protocol-as-simrepo and TE-32 split spec-side from implementation-side. The TE corpus is spec-side material for the wire-lab harness specifically. Per-protocol simrepos under `protocols/<slug>.d/docs/thought-experiments/` would have their own corpora, each with its own editing policy following the rule this TE locks.
- The user's "100-year goal" framing (TE-28) and the audit-trail framing (TE-21 + TE-23) both demand that any rewrite be traceable: a future reader must be able to reconstruct what a TE said before any later refinement.
- Git is the audit substrate. Any in-place rewrite of a TE preserves the prior version in git history; `git log -p <file>` gives byte-for-byte diff. The rewrite policy can therefore lean on git rather than requiring duplicate copies inside the file.

## What "one big TE with facets" actually means

The user's instruction is short but consequential. Working out what it means in practice:

**A facet is a TE.** Each TE is a focused look at one decision under test. The corpus is not the union of independent essays; it is the simulation of a single design problem (a wire lab that must survive 100 years) viewed from many sides. TEs that look unrelated at first reading often turn out to share assumptions and constraints — TE-1 (promise-stack ordering) and TE-29 (protocols-as-simrepos and the L4-binding layer) are both, fundamentally, about *who names what and where the boundary lives*. TE-21 (spec doc as promise) and TE-31 (spec-doc inversion) are the same question asked at different scales.

**Reading one TE in isolation is a category error.** When I (the bot) want to know whether TE-1's recommended set still stands, I should not just re-read TE-1; I should ask "what later facets touch this question?" That is what I did when the user asked about TE-1, and the answer was "very little touches the substance, but the layout work makes locking TE-1 more pressing, not less." That kind of cross-TE reading should be the default, not an extra step.

**Citation crosses freely.** Later TEs cite earlier ones explicitly; earlier TEs do not need to forward-cite later ones to remain accurate. But a *path* in an early TE that points at a now-moved file is a real defect against the corpus's navigability — the corpus is meant to be readable end-to-end at any point in time, and stale paths break that. Same for vocabulary that has been deprecated repo-wide.

**Holistic does not mean unified.** The corpus does not need a single thesis. TEs are allowed to disagree, even with later TEs. When they disagree, the resolution lives in the later TE's conclusions and DI; the earlier TE keeps its alternative analysis intact. What gets refined is the *framing*, not the *judgment*.

**The corpus has a temporal axis.** A TE was written at a moment with a particular set of locked-in assumptions. Reading TE-1 today, the reader needs to know that TE-1 was written before TE-29 reframed protocols-as-simrepos. The corpus carries this implicitly via timestamps, but a reader following links from `harness-spec-draft.md §8` to TE-1 may not know that the §1.1 path reference inside TE-1's conclusions section is from before the migration. A path-update sweep is the cheap fix; a vocabulary-shift sweep is the expensive one.

## Categories of post-hoc edit

Working through what kinds of changes a TE might need after first commit, in increasing order of caution required:

### Cat-1: Path or filename references

A TE refers to `specs/harness-spec-draft.md` but the file is now at `protocols/wire-lab.d/specs/harness-spec-draft.md`. The TE's claim is unchanged; only the pointer is stale. **Recommended treatment**: sweep in place during the migration that moved the file. No top-of-file note required; the change is mechanical and the diff is self-explanatory in `git log -p`.

### Cat-2: Vocabulary updates

A TE uses "channel" but the repo-wide term is now "transport." The TE's analysis and decisions are unchanged in substance, but the words have shifted. **Recommended treatment**: rewrite in place with a top-of-file vocabulary note (the TE-26 pattern). The note explains the rewrite, points at the TE that drove the rename, and asserts that structure / DF labels / locked decisions are unchanged. Future readers see the new vocabulary plus a marker that lets them go look up the old wording in git.

### Cat-3: Reference updates to deferred follow-ons

A TE's "Implications" or "Future work" section names a downstream DR that has since been filed, or names an open question that has since been closed. **Recommended treatment**: add a "## Status as of YYYY-MM-DD" or "## Refinements" section at the bottom of the TE — append-only. Each refinement names the later TE / DR / DI that resolved or restated the item. The original future-work list stays intact; the addendum is a forward-pointer.

### Cat-4: Adding cross-references the original draft did not have

After a later TE is written, an earlier TE could benefit from a back-pointer ("see also TE-N"). **Recommended treatment**: same append-only "Refinements" section. The earlier TE's body is not modified; the refinements section is where forward-pointers accumulate.

### Cat-5: Recasting the alternative analysis

The TE's Alt-A through Alt-E framing turned out to miss a real alternative, or one of the analyzed alternatives was scoped wrong. **Recommended treatment**: do not rewrite. Write a new TE that supersedes the old one's framing; the new TE cites the old one, names what was missing, and produces the corrected analysis. The earlier TE's "Decision status" line gets updated to `superseded by TE-N` (a one-line edit) and that is the only contents change.

### Cat-6: Recasting the locked decision

The TE locked DF-N.M.X but the locked answer is now wrong. **Recommended treatment**: do not rewrite. Write a new DI in the relevant TODO file that supersedes the old DI (the existing DI append-only convention from AGENTS.md / DI-NNN supersedes machinery handles this). Update the TE's "Decision status" line to `superseded by DI-N-...`. Do not modify the original DF questions or the original answer.

### Cat-7: Recasting the assumptions

A TE assumes something about the repo or the world that is no longer true. **Recommended treatment**: same as Cat-3 + Cat-5 combined. Add a refinements section noting which assumption was invalidated and by which later TE; if the invalidated assumption changes the conclusions, write a superseding TE.

The pattern is: **mechanical updates (Cat-1, Cat-2) edit in place; substantive updates (Cat-5, Cat-6, Cat-7) write a new TE / DI; navigational updates (Cat-3, Cat-4) append a refinements section to the existing TE.**

## Alternatives

### DF-34.1 — How permissive should TE editing be?

#### Alt-1.A: Strict immutability — no contents changes after first commit

A TE's contents are byte-frozen once committed. Path renames, vocabulary changes, and stale references are accepted as artifacts of when the TE was written. New TEs supersede old ones for any substantive change.

- **Easier**: clean audit trail; no risk of silently rewriting design history; no top-of-file notes accumulate; no rule about what counts as "mechanical" vs. "substantive" is needed.
- **Harder**: the corpus accumulates stale paths and deprecated vocabulary that increase friction for every future reader. TE-26 and TE-24 would have to be reverted or split into superseding TEs. The forty-cell vocabulary update from "channel" → "transport" would have produced four new TEs that say nothing new.
- **New obligation**: a discipline of writing new "vocabulary update" or "path update" TEs whenever any sweep happens. Most of those TEs would carry no real content beyond "this TE updates references; it locks no new decision."

#### Alt-1.B: Permissive — any in-place edit is fine, git carries the history

A TE's contents are editable freely; git log preserves the prior versions. No top-of-file notes required.

- **Easier**: no rule overhead; sweeping is cheap; corpus stays current.
- **Harder**: a reader cannot see at a glance whether the TE in front of them is the original draft, a vocabulary refresh, or a substantive rewrite. The audit trail exists in git but is not visible in the rendered file. Cross-TE citations become ambiguous: does TE-29 cite TE-1's original framing or its refreshed framing?
- **New obligation**: every reader must check git log to know whether they are reading the original.

#### Alt-1.C: Categorized — mechanical edits in place, substantive edits via supersedence, navigational edits via append-only refinements (the proposed policy in this TE)

Split the question by category. Cat-1 and Cat-2 edit in place (with a vocabulary note for Cat-2). Cat-3 and Cat-4 append. Cat-5 through Cat-7 write a new TE or DI. The corpus is current and navigable; substantive history is preserved by writing new artifacts rather than overwriting.

- **Easier**: matches the precedents already in the repo (TE-24 and TE-26 are Cat-2 edits with vocabulary notes; TE-25 is a Cat-5 superseding rewrite that produced a new TE rather than rewriting the channels-branch TE). The categories give the bot a checklist for any proposed edit.
- **Harder**: requires a category judgment for every edit. Edge cases between Cat-2 (vocabulary) and Cat-5 (recasting analysis) will need adjudication. The user is the ultimate arbiter when a category is in doubt.
- **New obligation**: a top-of-file vocabulary note for Cat-2 edits (one-paragraph header pointing at the TE / TODO that drove the rewrite); a "## Refinements" append-only section convention for Cat-3 and Cat-4.

#### Alt-1.D: Time-windowed mutability — TEs are mutable for some interval after drafting, then immutable

A TE is freely editable for, say, 7 days after first commit; afterward it locks. The interval lets the author refine the draft after first review; the lock prevents long-after-the-fact rewrites.

- **Easier**: clean rule; matches some traditional editorial practice.
- **Harder**: the wire-lab does not have a publication step that anchors "draft" vs. "published." TEs are committed when first drafted, not when finalized; the time window would have to be measured against an arbitrary moment. The 100-year framing makes any wall-clock interval feel arbitrary.
- **New obligation**: a clock to enforce the window; a rule about what happens to a TE that is rediscovered after the window expires.

### DF-34.2 — What about the wire-lab harness's TEs versus per-protocol TEs?

The wire-lab harness has a TE corpus at `docs/thought-experiments/`. Per-protocol simrepos under `protocols/<slug>.d/docs/thought-experiments/` will eventually have their own corpora.

#### Alt-2.A: One policy applies to all corpora

Whatever DF-34.1 locks applies uniformly across the harness corpus and all per-protocol corpora.

- **Easier**: one policy to remember; future protocol authors do not have to relitigate this question.
- **Harder**: assumes the wire-lab's editing policy is universally appropriate. A protocol whose threat model demands stricter immutability (e.g., a security-sensitive protocol whose TEs may be cited as evidence in dispute resolution) might want Alt-1.A.

#### Alt-2.B: Each protocol picks its own policy in its protocol-specific equivalent of this TE

The wire-lab harness locks DF-34.1's answer for its own corpus. Each per-protocol simrepo writes its own version of TE-34 (or cites this one) when it sets up its `docs/thought-experiments/`.

- **Easier**: matches the standing rule "let each protocol name its own internals" (TE-23 DI-23.2 Alt-2.B). A security-critical protocol can pick Alt-1.A; an exploratory protocol can pick Alt-1.C; a long-stable protocol can change its mind by writing its own superseding TE.
- **Harder**: more decisions to make per protocol; new authors may forget to make this decision and end up with no policy.

### DF-34.3 — How should the bot read the corpus?

#### Alt-3.A: One TE at a time, on demand

When a TE is referenced, read just that TE. Cross-references are followed only when the user asks.

- **Easier**: less context loaded; faster responses; matches how most readers approach long technical corpora.
- **Harder**: the bot misses the "TE-1 is more pressing now, not less" kind of insight that requires holding several TEs in working memory at once. The user has now explicitly named this failure mode.

#### Alt-3.B: Holistic reading by default

Treat the corpus as one document with facets. When any TE is referenced, the bot's first move is to identify which other TEs share assumptions, decisions, or vocabulary with it; reading is at corpus scope, not TE scope.

- **Easier**: matches the user's stated framing. Catches the "later work strengthens / weakens / supersedes earlier work" dynamic without prompting. Surfaces the "one big TE with facets" structure in answers, not just in this TE's prose.
- **Harder**: more work per question; risk of paralysis if the corpus grows large; risk of overweighting tangentially related TEs. The bot needs a triage discipline: read all TE titles + Decision-under-test sections first, then deep-read only the ones that share assumptions with the question at hand.

#### Alt-3.C: Holistic for substantive questions, single-TE for mechanical questions

Use Alt-3.B when the question is about decisions, framing, or design tradeoffs; use Alt-3.A when the question is mechanical (e.g., "what does TE-N's filename reference look like").

- **Easier**: matches what efficient reading actually looks like.
- **Harder**: requires the bot to classify the question first, which is itself a holistic-reading task. In practice this collapses to "always start holistic; narrow if the question is obviously mechanical."

## Scenario analysis

### S1 — Vocabulary sweep ("channel" → "transport")

**Alt-1.A** forces four new TEs that say nothing beyond "this TE renames N references." **Alt-1.B** rewrites silently and the next reader has no marker that vocabulary has changed. **Alt-1.C** rewrites with a top-of-file note (which is what the repo already did for TE-24 and TE-26). Alt-1.C is the actual current practice and matches the user's existing instinct.

### S2 — Path migration (top-level `specs/` → `protocols/wire-lab.d/specs/`)

**Alt-1.A** would have required either leaving stale paths inside TEs (degrading navigability) or writing a rename TE for every TE touched (boilerplate). **Alt-1.B** would have updated paths silently, indistinguishable from a substantive rewrite. **Alt-1.C** treats path renames as Cat-1 and updates in place with no header note required, because the diff is mechanical and self-explanatory. The just-completed TODO 014 migration intentionally *did not* update paths inside TEs because of my (the bot's) overgeneralized "do not back-edit" rule; under Alt-1.C those edits would have been correct.

### S3 — Decision is later overturned

A TE locks DF-N.M.X; a later TE shows DF-N.M.X was wrong. **Alt-1.A** is silent on what to do (the TE cannot be rewritten). **Alt-1.B** is dangerous: rewriting the TE to say the new answer destroys the audit trail of what was originally locked and why. **Alt-1.C** says: do not rewrite; write a new TE that supersedes the original, write a new DI that supersedes the original DI, and update the original TE's "Decision status" line to `superseded by TE-N / DI-N-…`. The locked decision history is preserved; readers following links to the old TE see the supersedence marker.

### S4 — A TE's "Implications and future work" section names downstream DRs

The TE was written before the DRs were filed; now they exist. **Alt-1.A** leaves the future-work list as a now-stale forecast. **Alt-1.B** rewrites the section, hiding the original forecast. **Alt-1.C** appends a "## Refinements" section that maps each future-work bullet to the DR / TE / DI that resolved it. The forecast is preserved alongside the resolution.

### S5 — Holistic reading catches a non-obvious refinement

The user asks "is TE-N still right?" Under **Alt-3.A** the bot re-reads TE-N and answers from inside that TE's frame. Under **Alt-3.B** the bot also reads TE-(N+k) and TE-(N+m) that share assumptions with TE-N, and notices that TE-(N+m) tightened an assumption TE-N relied on. The Alt-3.B answer is materially better. The user has explicitly named this failure mode by saying "be holistic"; that names the failure mode of Alt-3.A.

### S6 — The bot proposes a TE edit; the user wants the original wording back

Under **Alt-1.A** this can never happen because no edit is allowed. Under **Alt-1.B** the user must check git log to know what the original said. Under **Alt-1.C** the answer is in the file: Cat-1 / Cat-2 edits have a top-of-file note pointing at the rewriting TE; Cat-3 / Cat-4 edits are in a clearly-marked refinements section the user can ignore; Cat-5+ edits did not modify the TE. Reverting an Alt-1.C edit is cheap because the categories make the rewrite scope visible.

## Surviving alternatives

DF-34.1 Alt-1.A is rejected: it produces too much boilerplate and contradicts the existing TE-24 and TE-26 precedents already in the repo. Alt-1.B is rejected: it loses the navigational marker that lets readers tell mechanical updates from substantive ones. Alt-1.D is rejected: the wire-lab has no natural "draft" → "published" transition. **Alt-1.C survives.**

DF-34.2 Alt-2.A and Alt-2.B both work; Alt-2.B aligns with "let each protocol name its own internals." **Alt-2.B is recommended.** This TE locks the policy for the wire-lab harness corpus; per-protocol simrepos can adopt or override.

DF-34.3 Alt-3.A is the failure mode the user named. Alt-3.B is the explicit instruction. Alt-3.C is what Alt-3.B looks like in practice (no question is purely mechanical, but obvious mechanical questions can short-circuit). **Alt-3.B is recommended, with Alt-3.C as the day-to-day implementation.**

## Conclusions

1. **TE filenames remain immutable.** TE-25's invariant stands.

2. **TE contents are edited under three regimes**, by category:
   - **Mechanical edits** (path renames, vocabulary updates) edit in place. Vocabulary updates carry a top-of-file note pointing at the TE / TODO that drove the rewrite (the TE-24 / TE-26 pattern). Path renames need no note.
   - **Navigational edits** (forward-pointers to later TEs / DRs / DIs that resolved a future-work bullet, or back-pointers added after the fact) go in an append-only "## Refinements" section at the bottom of the TE.
   - **Substantive edits** (recasting alternative analysis, reversing a locked decision, invalidating an assumption) are not edits. Write a new TE that supersedes the old one; write a new DI that supersedes the old DI. Update only the old TE's "Decision status" line to point at the supersedence.

3. **The TE corpus is read holistically by default.** When any TE is in scope, the bot's first move is to scan TE titles + "Decision under test" sections across the corpus to find facets that share assumptions, vocabulary, or decisions. Single-TE reading is reserved for obviously mechanical questions.

4. **Per-protocol corpora set their own editing policy.** Each `protocols/<slug>.d/docs/thought-experiments/` directory may adopt this TE's policy by reference, or write its own superseding TE when stricter or looser policy is appropriate.

5. **The previous bot rule "do not back-edit TEs" is retired.** It was overgeneralized from the contest-artifact immutability rule (DI-003), which applies to adversarial review records, not to the author's own thinking. The current policy is the per-category rule above.

## Implications for the repo's open TODOs and pending DIs

- **TODO 5 (TE-1 promise-stack ordering, drive to DI):** When this lands, the DI work updates TE-1's "Decision status" line from `needs DF` to `decided per DI-005-NN-…`, and the "Implications" section gets a "## Refinements" append-only block listing which downstream items were filed (DR-006 normalization, the harness-spec §1.1 update, etc.). TE-1's body is otherwise untouched.
- **TODO 014 (just landed):** The path references inside TE files were intentionally left stale because the bot was applying the now-retired "do not back-edit" rule. Under TE-34's Cat-1, those paths should be updated in place. A follow-on TODO under `protocols/wire-lab.d/TODO/` will sweep them. (To be filed alongside the DI for this TE.)
- **TODO 6 + TODO 7 (provenance backfill):** The retrospective DIs / DRs under these TODOs do not modify TEs; they modify the harness-spec and DR/. Unaffected by this policy.
- **All future TEs:** The "## Refinements" section convention is new; it does not appear in TEs 1–33. Adding it to existing TEs is itself a Cat-3 / Cat-4 edit and is permitted under this TE's own rule. The convention will be exercised when TODO 5 lands its DI.
- **AGENTS.md / AGENTS-codex.md / AGENTS-ppx.md:** Need a new "## TE editing policy" section that codifies the conclusions above. This TE's DI, when locked, drives that AGENTS update.
- **`docs/thought-experiments/README.md`:** The current text is silent on contents (it only mentions filenames). Add a "## Editing policy" section pointing at TE-34, summarizing the three regimes, and naming the "## Refinements" convention.

## Decision Framing — questions for the user

Each DF question is a single paragraph below, with multiple-choice answers, recommended choice marked.

**DF-34.1**: How permissive is TE editing after first commit? **Alt-1.A** treats TE contents as byte-frozen — any change requires a new superseding TE, even path renames and vocabulary updates. This is clean but produces boilerplate and contradicts the TE-24 / TE-26 precedent already in the repo. **Alt-1.B** allows any in-place edit and lets git carry the history; this is cheap but loses the navigational marker that distinguishes mechanical edits from substantive ones. **Alt-1.C** (recommended) categorizes edits: mechanical edits (paths, vocabulary) in place, with a vocabulary note for Cat-2; navigational edits (forward-pointers) in an append-only "## Refinements" section; substantive edits (recasting analysis, reversing decisions) write new TEs and supersede. **Alt-1.D** would make TEs mutable for a fixed wall-clock window then lock; the wire-lab has no natural "draft → published" boundary that anchors such a window.

- (a) Alt-1.A — strict immutability.
- (b) Alt-1.B — permissive in-place edits; git carries history.
- (c) Alt-1.C — categorized: mechanical-in-place, navigational-append, substantive-supersede. **(Recommended.)**
- (d) Alt-1.D — time-windowed mutability.

**DF-34.2**: Does this TE's editing policy apply to per-protocol simrepos' TE corpora, or do they pick their own? **Alt-2.A** locks one policy for all corpora — simpler to remember but assumes the wire-lab's policy is universally appropriate. **Alt-2.B** (recommended) lets each protocol set its own policy in a protocol-specific equivalent of this TE, matching the standing "let each protocol name its own internals" rule (TE-23 DI-23.2 Alt-2.B); the wire-lab harness locks the policy in DF-34.1's answer for its own corpus.

- (a) Alt-2.A — uniform policy across all TE corpora.
- (b) Alt-2.B — per-protocol; wire-lab harness locks its own; others adopt or override. **(Recommended.)**

**DF-34.3**: How should the bot read the corpus by default? **Alt-3.A** is one TE at a time, on demand — fast but misses cross-TE refinements (the failure mode the user named). **Alt-3.B** (recommended) treats the corpus as one document with facets: when any TE is in scope, the bot first scans TE titles + "Decision under test" sections to identify shared-assumption / shared-vocabulary / shared-decision neighbors, then deep-reads selectively. **Alt-3.C** is Alt-3.B's day-to-day form: holistic for substantive questions, narrowed for obviously mechanical ones.

- (a) Alt-3.A — single-TE reading on demand.
- (b) Alt-3.B — holistic by default; the corpus is one document. **(Recommended; Alt-3.C is its practical form.)**
- (c) Alt-3.C — holistic for substantive, narrowed for mechanical.

The recommended set is **(34.1.c, 34.2.b, 34.3.b)**: categorized editing policy, per-protocol applicability, holistic reading by default.

## Decision status

`needs DF` — awaiting user choice on DF-34.1 through DF-34.3. After DF, the locked decisions become DI entries in `protocols/wire-lab.d/TODO/TODO-<timestamp>-te-editing-policy-and-holistic-corpus.md` (to be filed alongside this TE), and AGENTS.md / AGENTS-codex.md / AGENTS-ppx.md gain a new "TE editing policy" section.
