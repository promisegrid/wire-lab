# TE-35: Tabletop simulation of the TE editing policy

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-20260502-232651

## Decision under test

TE-34 locked an editing policy for the TE corpus (DI-020-20260502-213103: categorized edits — Cat-1 path renames in place, Cat-2 vocabulary updates in place with a top-of-file note, Cat-3 / Cat-4 forward / back-pointers in an append-only `## Refinements` section, Cat-5 / Cat-6 / Cat-7 substantive changes via supersedence). It also locked uniform applicability across all TE corpora (DI-020-20260502-213104) and an explicit substantive-vs-mechanical reading split (DI-020-20260502-213105). TE-34's own Refinements section (added 2026-05-02) admits that its scenario analysis was one-paragraph judgments, not tabletop play. This TE runs the tabletop the policy never got.

The decision under test is not whether to edit TEs at all — that question is settled — but whether the locked policy, as written, survives concrete play. Specifically: do the categorical boundaries (Cat-1 / Cat-2 / Cat-3 / Cat-5) hold when actors with different goals interact with the corpus at the same time, when adversarial pressure is applied, and when the policy applies to itself? If the tabletop confirms the policy as locked, AGENTS rollout (TODO 020 subtask 020.5) and the corpus sweeps (020.6, 020.7) unblock. If the tabletop surfaces a refinement, a superseding DI is written first.

This TE simulates the policy under six concrete scenarios with named actors. It does not propose new alternatives to the editing policy as a whole; it stress-tests the locked one. Where a scenario shows the policy bending, the TE names the refinement and recommends whether it rises to a superseding DI or to a Cat-3 Refinement on TE-34.

## Assumptions

- DI-020-20260502-213103 / -213104 / -213105 are in force as written. The tabletop is an audit of those DIs, not a re-litigation. If a scenario surfaces a problem the DIs as written cannot solve, that is a finding; the TE does not pre-empt the DIs.
- TE filenames are immutable per TE-25. The timestamp slug pins the integer alias; renames re-key the corpus. Out of scope for re-test; treated as ground truth.
- DI-003 (contest-artifact immutability for `proposals/pending/<branch>/contest-*.md` and `review-*.md`) is unaffected by the TE editing policy. Out of scope.
- The wire-lab today has 35 TE files (TE-1 through TE-35 once this lands). Most of them are pre-categorization; they were committed before TE-34 existed. The policy applies retroactively, which is itself a scenario-relevant fact.
- Actors in the tabletop follow the cryptography-literature alphabetical convention: Alice, Bob, Carol, Dave, Ellen, Frank are cooperative TE authors and readers; Mallory is the adversary. Steve is named explicitly when his role as repo owner is load-bearing.
- Git is the audit substrate. Any in-place rewrite of a TE preserves the prior version in `git log -p`; the policy can rely on that.

### Threat / trust model

- Cooperative actors are honest and competent but may be working in parallel with stale information about each other's in-flight work.
- Mallory has commit access (or is able to convince a maintainer to merge a patch) and aims to retroactively soften a locked decision without writing a superseding TE.
- A reader joining the project today knows nothing about which edits have happened to which TEs unless the TE itself or its `## Refinements` section says so.
- A reader visiting the corpus a year from now, after Mallory's hypothetical edit, has only what is in the file plus what `git log -p` reveals.

## Candidate alternatives

The policy locked three orthogonal answers. The tabletop walks the locked answer through each scenario and, for contrast, walks one rejected alternative through the same scenario where that contrast is informative.

- **Alt-1.A** (strict immutability — rejected by DI-020-20260502-213103) is walked alongside the locked Alt-1.C wherever an in-place edit is at issue.
- **Alt-1.B** (permissive in-place, no notes — rejected) is walked alongside Alt-1.C wherever the top-of-file note convention is load-bearing.
- **Alt-2.A** (uniform applicability — locked) is walked alongside Alt-2.B (per-protocol — rejected) wherever a per-protocol corpus enters the scenario.
- **Alt-3.C** (explicit substantive / mechanical split — locked) is walked alongside Alt-3.A (single-TE on demand — rejected) wherever a reader with limited time has to decide what to read.

The locked DI is the policy under test. The rejected alternatives are foils that make the locked DI's behavior visible.

## Scenario analysis

### S1 — Alice fixes a typo while Bob is mid-sweeping a vocabulary change

**Setup.** Alice notices a missing period in TE-12's S2. She opens a one-line patch on twig `ppx/te-fix-typo`. At the same time, Bob is on twig `ppx/te-vocabulary-sweep` rewriting "channel" to "transport" across the entire TE corpus, including TE-12; his sweep includes a top-of-file note on each touched file pointing at the TE-26 / TE-27 vocabulary lock that drove the rewrite. Alice and Bob commit and push their twigs within an hour of each other.

**Locked policy (Alt-1.C).** Alice's typo fix is Cat-1 by analogy: it is mechanical, the diff is self-explanatory, no top-of-file note is required. Bob's sweep is Cat-2: each touched file gets a top-of-file note. When the twigs land in order — say Alice first, then Bob — Bob's sweep cleanly merges over Alice's typo. When they land Bob first, then Alice, Alice's twig still applies cleanly because the typo line is unchanged by the vocabulary sweep (different bytes). What about the top-of-file note Bob added? Alice's twig was forked off ppx/main *before* Bob's note existed, so her diff does not touch the note region; the merge is clean. Either order works.

What does a reader see afterward? A TE-12 that has Bob's vocabulary note at the top, the typo fixed in S2, and a `git log -p` history showing two commits in some order. Both edits are recoverable. Cat-1 + Cat-2 compose cleanly.

**Alt-1.A foil.** Under strict immutability, Alice cannot fix the typo without writing a new "TE-12-typo-fix" TE that supersedes TE-12. Bob cannot run the sweep without writing 30+ "vocabulary update" TEs that supersede the originals. The corpus doubles in TE count for two trivial passes; readers must follow `superseded by` chains to find the current text. Cat-1 + Cat-2 under Alt-1.A is impractical; the locked Alt-1.C is doing real work here.

**Alt-1.B foil.** Under permissive-no-notes, Bob's sweep silently rewrites 30+ TEs with no marker that vocabulary has changed. Alice's typo fix is indistinguishable from a substantive edit. A reader visiting TE-12 a year later sees current text but no marker that the wording was changed; the only way to know is `git log -p`. The Cat-2 top-of-file note is what makes Alt-1.C better than Alt-1.B in this scenario.

**Finding.** Locked policy survives. The TE-26 / TE-27 precedent is what is being formalized; this scenario reproduces that precedent and confirms it.

### S2 — Three-deep supersedence and a reader who follows a citation

**Setup.** TE-N (a real-shaped example: pretend N=10) locks a decision via DI-N. Five TEs later TE-(N+5) supersedes TE-N's analysis with a wider scope; DI-(N+5) supersedes DI-N. Four more TEs later TE-(N+9) supersedes TE-(N+5)'s analysis again, after a new constraint surfaces; DI-(N+9) supersedes DI-(N+5). Carol, a new contributor reading TE-(N-2)'s "Implications" section a year later, sees a reference to TE-N. She clicks through.

**Locked policy (Alt-1.C).** TE-N's body is unchanged — substantive edits do not edit the original — except for its `Decision status` line, which under Alt-1.C reads `superseded by TE-(N+5) / DI-(N+5)`. TE-(N+5)'s `Decision status` reads `superseded by TE-(N+9) / DI-(N+9)`. Carol clicks TE-N → reads its analysis → sees `superseded by` → clicks TE-(N+5) → reads → sees `superseded by` → clicks TE-(N+9) → reads → no further supersedence → this is current. Three hops.

Is Carol confused? She reads three TEs to reach the current decision. Each tells her: "this analysis is preserved; the decision moved here." That is intentional under Alt-1.C: the original analysis is preserved as evidence of what was considered when the decision was made, not as the current decision. The three-hop chain is the cost of full provenance.

What if the chain is even longer — five or ten supersedences deep? The policy as written has no upper bound. At some chain depth, a reader is doing archaeology, not orientation. A reader who only wants the *current* decision should read the latest TE in the chain, not start from the oldest.

**Alt-3.A foil (single-TE on demand).** Carol reads only TE-N because the citation pointed her there. She reads `superseded by TE-(N+5) / DI-(N+5)` and stops, treating the supersedence note as advisory. She acts on TE-N's locked decision — which is no longer current. Alt-3.A is wrong here; the locked Alt-3.C says she scans the corpus for shared assumptions and would notice TE-(N+5) and TE-(N+9) before acting.

**Locked policy (Alt-3.C).** Carol's question is substantive (she wants to know what decision is in force). Per Alt-3.C she scans TE titles + Decision-under-test sections across the corpus, notices that TE-(N+5) and TE-(N+9) share the same decision space, deep-reads them, finds the chain, and acts on TE-(N+9). Three TE reads, but ordered correctly: read the latest first, walk back only as far as needed.

**Finding with refinement candidate.** The locked Alt-1.C (substantive supersedence) and Alt-3.C (substantive reading scans corpus first) compose correctly. But a reader of TE-N has no forward-pointer to TE-(N+5) inside TE-N's body — only the `Decision status` line. If Carol misses the `Decision status` line (it is at the bottom of the file), she may read 200 lines of analysis under the impression that it is current. **Recommended refinement:** make the `Decision status` line position-invariant — either always near the top (as a header field) or always conspicuous. This is a Cat-3 / Cat-4 navigational tweak to TE-1 through TE-35, not a DI revision; it can land as part of subtask 020.7 (the TE-1 Refinements work) extended to also relocate `Decision status` to a top-of-file header field across all TEs, with a Cat-1 sweep. **Filed as future-work** in this TE's Implications section; no DI change.

### S3 — Path migration during an in-flight twig

**Setup.** Dave has been working on twig `ppx/te-new-decision` for three days. His twig contains a new TE that references `specs/harness-spec-draft.md` (the pre-migration path). Meanwhile, on ppx/main, TODO 014's protocols-as-simrepos migration lands and renames the spec to `protocols/wire-lab.d/specs/harness-spec-draft.md`. Dave rebases his twig on the new ppx/main.

**Locked policy (Alt-1.C, Cat-1).** Dave's rebase produces a merge conflict: his TE references the old path; ppx/main has the new path. Cat-1 says the path-rename is mechanical, edited in place, no note required. Dave updates the path in his TE and continues. When Ellen reviews the twig, she sees one path-rename in the diff alongside Dave's substantive new content; the path-rename is self-explanatory in `git log -p`.

What if Dave's TE is *already merged* and the path migration is the in-flight work? That is the actual TODO 014 case. Under Alt-1.C, the migration includes a Cat-1 sweep over all TEs that reference the old path. The sweep is permissible; its scope is "every reference inside any TE." The sweep itself is what subtask 020.6 enacts.

**Alt-1.A foil.** Under strict immutability, Dave cannot update the path in his TE without writing a new TE. The path migration cannot sweep TEs without writing 30+ "path-update" TEs. The migration is impractically expensive. The TODO 014 commit log already shows this: the migration was performed but the TE path-references were intentionally left stale because the bot was applying the now-retired "do not back-edit" rule. Under Alt-1.A those references would have to remain stale forever (or trigger 30+ supersedence TEs); under Alt-1.C they can be swept.

**Locked policy survives, but a sub-question surfaces.** When the sweep happens, what about TEs that quote the old path *as a historical fact* — e.g., TE-30 (`TE-20260430-213447-te-numbering-collision-and-harness-spec-path`) literally describes "the channels branch edits a top-level `harness-spec.md`, but on `ppx/main` the spec was renamed to `specs/harness-spec-draft.md` during the genesis-freeze work in TODO 011." That sentence is about the *act of renaming*; the literal string `specs/harness-spec-draft.md` is the object of historical discussion, not a current pointer. A naive Cat-1 sweep would update it to `protocols/wire-lab.d/specs/harness-spec-draft.md`, which would distort the historical claim ("the spec was renamed to" the post-TODO-014 path is anachronistic).

**Refinement candidate.** Cat-1 needs a "historical-quotation exemption." A path string used as a current pointer (e.g., a markdown link `[Wire Lab](../../specs/harness-spec-draft.md)`) is updated; a path string used as the literal object of historical discourse is not. **The boundary:** if updating the string changes a true historical claim into a false one, do not update it. Example: in TE-30, the sentence "the spec was renamed *to* `specs/harness-spec-draft.md`" should not be re-rewritten to say "the spec was renamed to `protocols/wire-lab.d/specs/harness-spec-draft.md`" because the rename in question was the genesis-freeze rename, not the protocols-as-simrepos rename.

This is a meaningful refinement of Cat-1. **Recommended treatment:** Cat-1 is split into Cat-1a (current-pointer paths — sweep) and Cat-1b (historical-quotation paths — leave). The sweep tool (a regex-replace would be wrong) must be a careful manual or semi-manual review per match. **This rises to a superseding DI** because it materially changes how subtask 020.6 will be executed.

**Finding.** Locked policy bends. File DI-020-20260502-232651-cat1a-cat1b refining Cat-1 to distinguish current-pointer paths from historical-quotation paths. Add to the Cat-1 description in TE-34's Refinements section. Subtask 020.6 must enumerate matches and review each before sweeping.

### S4 — Per-protocol corpus with a stricter regime

**Setup.** Frank is the maintainer of a hypothetical security-sensitive protocol (`protocols/identity-attestation.d/`) under the wire-lab umbrella. His protocol's TEs may be cited in dispute resolution when an attestation is challenged. He wants Alt-1.A (strict immutability) for his corpus: once a TE is committed, even a typo fix requires a superseding TE.

**Locked policy (Alt-2.A, uniform).** DI-020-20260502-213104 says one policy for all corpora. To get Alt-1.A, Frank must write a TE that supersedes DI-020-20260502-213104 *for his corpus only*. The locked policy explicitly affords this: "A future protocol that judges its threat model demands a stricter regime [...] may write its own superseding TE that locks a different policy for that corpus only; absent such a superseding TE, this policy is in force."

Frank writes `protocols/identity-attestation.d/docs/thought-experiments/TE-NN-stricter-editing-policy.md` framing the question, citing DI-020-20260502-213104 as the override target, walking his threat model through the same six scenarios as this TE, and locking Alt-1.A for his corpus via a new DI in his protocol's TODO. The override is per-corpus; the wire-lab harness corpus and other per-protocol corpora are unaffected.

**Alt-2.B foil.** Under per-protocol-by-default, Frank does not need to write an override TE — Alt-1.A would just be the choice he makes when he sets up his corpus. But other protocol authors who do not think about editing policy at all end up with no policy, which is worse than having to override.

**Walking it.** Under Alt-2.A, Frank's path is: (1) detect the friction (his corpus needs Alt-1.A); (2) write the override TE; (3) lock the override DI; (4) operate his corpus under Alt-1.A. Under Alt-2.B, Frank's path is: (1) when setting up his corpus, decide on a policy; (2) document it in his corpus's setup TE; (3) operate under that policy. Alt-2.B is one step shorter for Frank but two steps longer (or undefined) for a protocol maintainer who does not have an opinion.

**Mid-stream.** What if Frank's protocol operates under Alt-2.A's default for two years, and in year two he realizes his threat model needs Alt-1.A retroactively? Under Alt-2.A he writes the override TE now; it applies prospectively. TEs already in his corpus were written under Alt-1.C and remain editable under Alt-1.C until the override TE locks Alt-1.A; from that point forward, his corpus's TEs are immutable. This is unambiguous under Alt-2.A — there is one switch point and the override TE pins it.

**Finding.** Locked policy survives. The override path works for Frank's case and produces a clean retroactive switch.

### S5 — Mallory's adversarial Cat-2 vocabulary update

**Setup.** Mallory wants to retroactively soften DI-005 (a hypothetical promise-stack ordering decision locking Alt-A), which she opposes. She cannot delete the DI — DIs are append-only by AGENTS.md DR/DI Source-of-Truth Protocol. She cannot edit DR-005's body — DRs are append-only too. But she can target TE-1 (the TE that fed DI-005) with what looks like a Cat-2 vocabulary update: she rewrites "outermost-first" to "first-listed" throughout TE-1's body, claiming this is a vocabulary clarification. The two phrasings are subtly different — "outermost-first" pins ordering to nesting depth; "first-listed" admits any sequence the reader interprets as a list. A reader of post-Mallory TE-1 who was undecided between Alt-A and Alt-D could now read TE-1 as supporting Alt-D.

**Locked policy (Alt-1.C, Cat-2).** Mallory's diff carries a top-of-file note (otherwise it fails the Cat-2 form): "Vocabulary updated 2026-MM-DD; structure, DF labels, and locked decisions unchanged." The note is a load-bearing assertion about what was *not* changed. Steve, on review, reads the note and the diff. Two outcomes:

**Outcome 5.a — Steve catches it.** The note claims locked decisions are unchanged. Steve reads the diff, notices that "outermost-first" → "first-listed" changes the meaning of Alt-A's premise, and flags the patch. He blocks the merge. The Cat-2 form's note is what made the bad-faith claim explicit and falsifiable — Mallory had to *assert* that nothing substantive changed, and Steve could check the assertion. Without the note (Alt-1.B), Mallory's diff would have looked like any other in-place edit, and Steve would have had to compare meanings without a marker.

**Outcome 5.b — Steve does not catch it.** The diff merges. A year later, Carol reads TE-1, finds it supportive of Alt-D, and acts. Some time after that, the discrepancy is noticed. Recovery: write a new TE that re-locks Alt-A by re-deriving it from scenario play (Mallory cannot rewrite scenarios without making the rewrite obvious — the actors and their actions are concrete); supersede the corrupted DI if necessary; the corrupted text remains in TE-1 with its Cat-2 note and with `git log -p` revealing the rewrite. The corruption is visible-in-hindsight even after the fact, because the note exists.

**Alt-1.B foil.** Without the top-of-file note, Mallory's edit is indistinguishable from any in-place edit. There is nothing to claim and nothing to falsify. The corruption is harder to detect; recovery requires comparing against `git log -p` blindly across all TE-1 edits. The note is doing real adversarial-resistance work.

**Alt-1.A foil.** Mallory could not have edited TE-1 at all; she would have had to write "TE-MM-vocabulary-clarification" as a superseding TE. That superseding TE is a fresh artifact that carries the bad-faith semantic shift; Steve reviews it as a new TE, with full DF / Decision-status framing. The bad-faith pressure has nowhere to hide. Alt-1.A is *more* adversarial-resistant than Alt-1.C in this scenario.

**Refinement candidate.** Cat-2 vocabulary updates need a stronger guard. The top-of-file note is necessary but not sufficient; Cat-2 in adversarial conditions wants either (a) a co-signer (the maintainer plus one other reviewer must both endorse), or (b) a "this affects no DI" check in the note that explicitly enumerates the DIs the rewrite does not touch, with a falsifiable claim per DI. **Recommended treatment:** Cat-2 notes must enumerate the DIs the rewrite asserts unchanged, by ID, so that Steve can check each. **This is a Cat-3 refinement to TE-34, not a DI revision** — it tightens the Cat-2 procedure without changing which category an edit falls in.

**Finding.** Locked policy survives but Cat-2 needs a procedural tightening. File the tightening as a Cat-3 Refinement on TE-34.

### S6 — Holistic reading on a quoting TE

**Setup.** TE-X (pretend X=22) quotes TE-Y (pretend Y=21) verbatim in its Assumptions section to anchor a definition. Later, a Cat-2 vocabulary sweep edits TE-Y in place. The quoted phrase in TE-X now does not match TE-Y's current text.

**Locked policy (Alt-3.C, holistic for substantive).** Carol reads TE-X. Under Alt-3.C she scans TE titles + Decision-under-test sections across the corpus and notices TE-Y is in scope (TE-X cites it explicitly). She deep-reads TE-Y, finds the Cat-2 top-of-file note, and reads the current TE-Y text. She notices that TE-X's quote uses the old vocabulary. Two outcomes:

**Outcome 6.a — the discrepancy is incidental.** TE-X's quote uses old vocabulary but the surrounding analysis was unaffected by the rename. Carol reads through; understanding survives. No action needed; or, optionally, a Cat-2 tweak to TE-X under the same vocabulary sweep.

**Outcome 6.b — the discrepancy is substantive.** The vocabulary change altered the meaning of the quoted phrase. TE-X's analysis no longer follows from its premises. This is a cross-TE substantive impact; Cat-5 / Cat-6 territory. Carol writes a new TE re-doing TE-X's analysis under the current vocabulary, supersedes TE-X, and updates TE-X's `Decision status` line.

**Alt-3.A foil (single-TE).** Carol reads only TE-X. She acts on TE-X's analysis with the old-vocabulary quote and never realizes the discrepancy. A year later, an inconsistency surfaces; tracing it back through the corpus is expensive.

**Refinement candidate.** When a Cat-2 sweep updates TE-Y, the sweeper must check whether any *other* TE quotes TE-Y verbatim. If so, the sweeper either (a) updates the quote in the quoting TE (Cat-2 cascade) or (b) leaves the quote intact and adds a Cat-3 Refinement note in the quoting TE pointing at TE-Y's current state. **Recommended treatment:** Cat-2 sweeps must include a "find quoting TEs" sub-step. Tooling (a `grep` for verbatim quoted phrases) makes this cheap. **This is a Cat-3 refinement to TE-34**, parallel to S5's Cat-2 tightening.

**Finding.** Locked policy survives but Cat-2 needs a cross-TE check. File the cross-TE check alongside the S5 tightening as a single Cat-3 Refinement on TE-34.

## What each alternative makes easier, harder, and what it newly demands

| Scenario | Locked policy holds? | Refinement needed? | Severity |
|---|---|---|---|
| S1 — concurrent typo + sweep | Yes | None | n/a |
| S2 — three-deep supersedence | Yes (Alt-3.C catches) | Move `Decision status` to top-of-file header | Cat-1 sweep, future TODO |
| S3 — path migration | Yes, with split | Cat-1 -> Cat-1a / Cat-1b | **Superseding DI** |
| S4 — per-protocol stricter | Yes | None | n/a |
| S5 — Mallory's Cat-2 attack | Partial; note is necessary, not sufficient | Cat-2 note must enumerate DIs asserted unchanged | Cat-3 Refinement on TE-34 |
| S6 — quoting TE under Cat-2 sweep | Partial; reader catches it under Alt-3.C | Cat-2 sweeps must check for verbatim quoting TEs | Cat-3 Refinement on TE-34 |

## Surviving alternatives

Three of the four locked DIs survive intact: DI-020-20260502-213104 (uniform applicability) and DI-020-20260502-213105 (substantive / mechanical reading split) are confirmed by S2 / S4 / S6.

DI-020-20260502-213103 (categorized editing) is confirmed at the macro level — Alt-1.A is too restrictive (S1, S3) and Alt-1.B is too permissive (S1, S5) — but Cat-1 specifically needs to be split into Cat-1a (current-pointer path updates — sweep) and Cat-1b (historical-quotation paths — leave). This rises to a superseding DI because it materially changes how subtask 020.6 is executed.

Cat-2 also needs procedural tightening (S5 enumeration-of-unchanged-DIs in the note; S6 cross-TE quoting check), but neither rises to a DI revision; both are Cat-3 Refinements to TE-34 that tighten the procedure without changing the category boundaries.

## Conclusions

1. **The locked editing policy substantially survives tabletop play.** The categorical scheme (Cat-1 / Cat-2 / Cat-3 / Cat-5) holds in all six scenarios. No category boundary is wrong; one category needs sub-division and one category needs procedural tightening.

2. **Cat-1 must split into Cat-1a (current-pointer paths) and Cat-1b (historical-quotation paths).** A naive regex-replace sweep over the corpus would distort the historical claims in TE-30, TE-32, TE-33, and other TEs that quote old paths as the literal object of historical discussion. The sweep tool must enumerate matches and the reviewer must classify each before applying. Subtask 020.6's plan changes accordingly.

3. **Cat-2 vocabulary updates need a strengthened top-of-file note.** The note must enumerate, by DI ID, the DIs the rewrite asserts unchanged. This makes the bad-faith Cat-2 attack (S5) explicitly falsifiable — the rewriter has to commit to a list of DIs they claim are unaffected, and the reviewer can check each. This is a Cat-3 Refinement on TE-34, not a DI revision.

4. **Cat-2 sweeps must include a cross-TE quoting check.** When TE-Y is rewritten Cat-2, any TE that quotes TE-Y verbatim must be either swept in cascade or annotated with a Cat-3 Refinement. Cheap to enforce with `grep`. This is a Cat-3 Refinement on TE-34, not a DI revision.

5. **The `Decision status` line should be relocated to a top-of-file header field across the corpus.** Currently it is at the bottom of the file, where a reader who acts on the body's analysis may never see the supersedence marker. This is a Cat-1 sweep (sub-class Cat-1a, current-pointer-equivalent navigational marker), filed as future-work; it does not block this TE from concluding.

6. **AGENTS rollout (subtask 020.5) unblocks after the superseding DI for Cat-1a / Cat-1b lands.** The two Cat-3 Refinements on TE-34 (Cat-2 strengthening, cross-TE quoting check) can land alongside the superseding DI or in a follow-up pass; they do not block AGENTS rollout because they tighten a procedure rather than change a rule's shape.

## Implications for the repo's open TODOs and pending DIs

- **TODO 020 subtask 020.5 (AGENTS rollout):** unblocks after the Cat-1a / Cat-1b superseding DI lands. AGENTS.md sections must reflect the split.
- **TODO 020 subtask 020.6 (TE path-reference sweep):** plan changes. The sweep is no longer a regex-replace; it is an enumeration + manual review. Each match is classified Cat-1a (sweep) or Cat-1b (leave) before edit. Estimated effort triples but accuracy demands it.
- **TODO 020 subtask 020.7 (TE-1 Refinements):** unaffected by the tabletop's findings. Proceeds as planned.
- **TODO 020 subtask 020.8 (README reconfirm):** unaffected.
- **New TODO subtask 020.10:** relocate `Decision status` to a top-of-file header field across the corpus. Cat-1a sweep. Filed as a future subtask under TODO 020.
- **Cat-3 Refinements to TE-34:** S5 strengthens the Cat-2 note (enumeration of unchanged DIs); S6 adds a cross-TE quoting check. Both land as a single edit to TE-34's `## Refinements` section once this TE is committed.
- **DI-003 (contest-artifact immutability):** unaffected. Confirmed by reasoning in TE-34, untested in this tabletop.

## Decision Framing — questions for the user

**DF-35.1**: Cat-1 needs to split into Cat-1a (current-pointer paths, swept) and Cat-1b (historical-quotation paths, left). Should this land as a superseding DI for DI-020-20260502-213103, or as a Cat-3 Refinement on TE-34 that re-interprets Cat-1 in place? **Alt-35.1.a (recommended)** writes a new DI (DI-020-20260502-232651-cat1a-cat1b) that supersedes DI-020-20260502-213103, since the change materially affects subtask 020.6's execution and a future reader of the AGENTS.md TE Editing Policy section needs to see the split as a first-class rule, not as an exception buried in a Refinements section. **Alt-35.1.b** treats it as a Cat-3 procedural Refinement on TE-34 — cheaper, but the Refinements section then carries a substantive boundary change that some readers will miss. **Alt-35.1.c** does both: file the new DI and also add the Refinement on TE-34 pointing at the new DI; redundant but maximally visible.

- (a) Alt-35.1.a — superseding DI; Cat-1a / Cat-1b promoted to first-class rule. **(Recommended.)**
- (b) Alt-35.1.b — Cat-3 Refinement on TE-34 only.
- (c) Alt-35.1.c — both (DI plus a Refinement pointer on TE-34).

**DF-35.2**: Cat-2 strengthening — should the top-of-file note for a Cat-2 vocabulary update be required to enumerate the DIs the rewrite asserts unchanged, by DI ID? **Alt-35.2.a (recommended)** makes the enumeration mandatory: the note must list the DI IDs the rewrite does not touch, so a reviewer can check each. The S5 attack is explicitly defeated. Cost: a few extra minutes per Cat-2 sweep. **Alt-35.2.b** makes the enumeration recommended-not-required: maintainers may include the list when they think it useful. The S5 attack is partially defeated only when maintainers happen to include the list. **Alt-35.2.c** keeps the note as locked (pointer to the rewriting TE / TODO; assertion that structure, DF labels, and locked decisions are unchanged) without enumeration; relies on reviewer diligence. The S5 attack relies on reviewer diligence; the note is necessary, not sufficient.

- (a) Alt-35.2.a — enumeration mandatory in Cat-2 notes. **(Recommended.)**
- (b) Alt-35.2.b — enumeration recommended but not required.
- (c) Alt-35.2.c — no change; current note convention stands.

**DF-35.3**: Cat-2 cross-TE quoting check — when TE-Y is rewritten Cat-2, must the sweeper check whether any other TE quotes TE-Y verbatim and either sweep-cascade or annotate? **Alt-35.3.a (recommended)** makes the check mandatory: every Cat-2 sweep includes a `grep` pass for verbatim quotes; matches are either swept in cascade or annotated with a Cat-3 Refinement. **Alt-35.3.b** makes the check recommended-not-required. **Alt-35.3.c** drops the check; relies on holistic reading (Alt-3.C) to catch discrepancies at read time, not write time. The S6 reader does catch the discrepancy under Alt-3.C, but only if she is reading substantively; mechanical reads miss it.

- (a) Alt-35.3.a — cross-TE quoting check mandatory in Cat-2 sweeps. **(Recommended.)**
- (b) Alt-35.3.b — recommended but not required.
- (c) Alt-35.3.c — no check; rely on Alt-3.C reader to catch.

**DF-35.4**: `Decision status` relocation. Currently `Decision status` is at the bottom of every TE. S2 surfaced that a reader who acts on the body of a superseded TE may never see the supersedence marker. Should this be addressed by a Cat-1a sweep? **Alt-35.4.a (recommended)** files a new TODO subtask (020.10) to add a top-of-file header field (e.g., a YAML or simple `Status: ...` line near the TE ID) on every TE; sweep across the corpus; existing bottom-of-file `Decision status` section either stays as-is or is replaced. **Alt-35.4.b** leaves `Decision status` at the bottom; relies on Alt-3.C reader to scan to end. **Alt-35.4.c** addresses it differently — e.g., requires every TE that supersedes another to begin with an explicit "supersedes TE-N" line at the top, but does not relocate the existing `Decision status` section.

- (a) Alt-35.4.a — relocate to top-of-file header field; Cat-1a sweep across corpus. **(Recommended.)**
- (b) Alt-35.4.b — no change.
- (c) Alt-35.4.c — top-of-file "supersedes TE-N" marker only on superseding TEs; original `Decision status` stays at bottom.

The recommended set is **(35.1.a, 35.2.a, 35.3.a, 35.4.a)**: promote Cat-1a / Cat-1b to a superseding DI; make Cat-2 enumeration and cross-TE quoting checks mandatory; relocate `Decision status` to a top-of-file header.

## Decision status

`partially decided` — DF-35.1 answered Alt-1.C on 2026-05-02 (bot recommendation under Steve's `make a recommendation` directive); DI-020-20260502-232651 in TODO 020 supersedes DI-020-20260502-213103's Cat-1 clause, splitting Cat-1 into Cat-1a (current-pointer, swept) and Cat-1b (historical-quotation, left). DF-35.2 answered Alt-2.a on 2026-05-02 (same directive); recorded as a Cat-3 procedural Refinement on TE-34 mandating that Cat-2 top-of-file notes enumerate unchanged DIs by ID (no separate DI). DF-35.3 / -35.4 remain pending. The locked decisions from TE-34 (DI-020-20260502-213104 / -213105 plus the residual Cat-2/3/4/5/6/7 clauses of -213103) remain in force.
