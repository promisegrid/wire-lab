# Dropped-thread disposition — TODO 21 walk notes (turns 155-192)

**Date:** 2026-05-06
**Author:** stevegt-via-perplexity (proposed; awaits Steve review)
**Source:** `protocols/wire-lab.d/TODO/TODO-20260504-170746-session-replay-cleanup.md` — 186 UT entries spanning turns 155-192

## Purpose

This file is the Phase 1 output of the dropped-thread classification pass approved by Steve at turn 288. Each UT (unfinished thread) gets a proposed disposition. Phase 2 — drafting TEs in dependency order, starting with TE-38 — is gated on Steve's review of the dispositions below.

**Vocabulary correction (turn 287):** "plurality" / "pluralism" is bot coinage from turn 172; never endorsed by Steve. Steve endorsed **"feed"** and **"substrate"** at turn 175. Going forward use **"substrate-agnostic"**. Walk notes are not retroactively edited — they are append-only history.

## Summary

| Disposition | Count |
| --- | ---: |
| TE-38 | 52 |
| TE-40 | 18 |
| TE-41 | 15 |
| TE-42 | 7 |
| TE-43 | 25 |
| TE-45 | 1 |
| TE-36-followon | 5 |
| Spec-edit | 5 |
| Retire | 3 |
| Carry | 55 |
| **Total** | **186** |

## Proposed TE roster (Phase 2 work)

Drafting order (dependency-sorted):

1. **TE-38** Substrate-agnostic layered model — foundational; vocabulary, layer numbers, substrate naming flow downstream from this.
2. **TE-39** wire-lab-devs migration plan — depends on TE-38 layer naming. (Carved out at turn 177; not in this UT inventory because turns 175-178's migration DFs were folded into TE-38 walk notes.)
3. **TE-40** Apparatus-vs-specimen completion — independent of TE-38; closes TE-1 limbo and finishes the harness-spec sweep.
4. **TE-41** Group-session freeze procedure — depends on TE-40 (apparatus/specimen distinction); ratifies the freeze gate, membership-by-posting, merge-group-transport-spec signature, §8/§9 detail.
5. **TE-42** Filename / CID-cascade policy — depends on TE-41 (filename convention); when do we cascade vs deprecate vs grandfather? Covers Message-ID removal precedent and future Rabin-chunking cascade.
6. **TE-43** Promisebase prior-art adoption — depends on TE-38 layered model (decides which promisebase pieces wire-lab depends on at L6).
7. **TE-44** wire-lab/promisebase merge trajectory — depends on TE-43; ratifies or rejects Steve's turn-192 implicit "will gradually merge" commitment.
8. **TE-45** Conditional-release / geofencing / recursive promise-graph — orthogonal to TE-38 layered model; can land independently.
9. **TE-36 follow-on** — OQ-36.6 nine-axis comparison + Alice-through-Mallory tabletop walk; closes TE-36's deferred work.

**Spec-edits** are small enough to land directly without a TE.
**Carry** items are procedural lessons / AGENTS-rule reminders / cadence notes — they belong in AGENTS-ppx.md or per-turn discipline, not in TEs.
**Retire** items are superseded by other landed work.

---

## Disposition by cluster

### TE-36 follow-on: OQ-36.6 + tabletop walk (5 UTs)

- **UT-160.b** (T160): DF-36.4 PT-recast collapsed five named pCIDs into two PT primitives, baked into Alt-4.A rather than presented as its own DF.
- **UT-160.c** (T160): TE-36's six tabletop scenarios used different framing than the audit memo's recommended Alice-through-Mallory walk.
- **UT-161.a** (T161): The asymmetry between the two envelope hypotheses -- grid-pcid-payload can host promise-stack as one possible payload, but promise-stack cannot cleanly host grid-pcid-payload at envelope level -- is t...
- **UT-162.a** (T162): The deferred OQ-36.6 investigation is itself a major unfinished thread.
- **UT-162.b** (T162): DF-36.2 is provisionally locked pending OQ-36.6, but the conditionality may not survive visual inspection of the TE file.

### TE-38: Substrate-agnostic layered model (52 UTs)

- **UT-170.a** (T170): (no bolded summary)
- **UT-171.a** (T171): (no bolded summary)
- **UT-171.b** (T171): Bot's argument-against #5 sketches an alternative encoding -- substrate as instance-manifest field rather than directory level -- but does not propose where that manifest would live, what its schema w...
- **UT-172.a** (T172): The substrate-pluralism reframe explicitly retracts turn 171's argument-against #4 about the "git-file-transport" working name -- the bot's two-turn position on the same TODO-list artifact has flipped...
- **UT-172.b** (T172): DF-38.5 "is 'binding' the right name?" is opened with three named alternatives ('substrate,' 'carrier-medium,' 'delivery-mechanism') but no recommendation -- and turn 174 will challenge "binding" spec...
- **UT-172.c** (T172): Layout A vs Layout B (two presentations of how `bindings/` attaches to instances) is offered as a sub-decision under DF-38.2 / DF-38.4 with a tentative recommendation toward Layout B ("equivalently an...
- **UT-172.d** (T172): (no bolded summary)
- **UT-172.e** (T172): The bot's "clean spec evolution" claim -- that §9 can be replaced with a single short clause once binding protocols exist -- assumes the binding-protocol family is in place before any spec edit lands;...
- **UT-173.a** (T173): The wire-lab-as-content-addressed-Usenet framing is load-bearing for many downstream decisions but has not been ratified.
- **UT-173.b** (T173): DF-38.5's "binding" recommendation rests on W3C/WSDL precedent and conflict-avoidance with existing wire-lab vocabulary; turn 174 rejects it on different grounds (Usenet-gateway-vocabulary precedent),...
- **UT-173.c** (T173): Bot's closing question dropped its framing-component (single-part "want me to draft TE-38 now?") rather than the two-part action+framing pattern of turns 171 and 172. This is a cadence shift, not a ca...
- **UT-173.d** (T173): (no bolded summary)
- **UT-173.e** (T173): The git-itself-as-precedent observation is striking and may have non-obvious implications -- if git's own object/transport separation is the precedent, then wire-lab is in some sense building "git sem...
- **UT-174.a** (T174): The historical-networks-20260503 research doc just committed in the turn-173 cluster (commit `6890592`) used "binding" as one of multiple candidate terms with explicit non-recommendation, so it remain...
- **UT-174.b** (T174): `udp-binding` is now retroactively misnamed in wire-lab's own vocabulary -- a previously-merged spec name is in queue for rename, and the spec's frozen pCID under wire-lab manifest may complicate the ...
- **UT-174.c** (T174): Option γ's INSTANCE.md absorbs membership-and-feeds in one file, which collides with §8's still-open membership-management question -- INSTANCE.md may need to be a wire artifact, or governance manifes...
- **UT-174.d** (T174): The encapsulation-not-layering reframe is structurally important but introduces vocabulary ambiguity: encapsulation has its own connotations including OSI's PDU-encapsulation and tunneling-protocol en...
- **UT-174.e** (T174): Turn-173's three-pattern map A/B/C onto the wire-lab tree partially needs revision -- Pattern B mapped to a `bindings/` subdirectory which is now rejected; the new mapping is Pattern B -> INSTANCE.md ...
- **UT-175.a** (T175): m000 transport-creation will undergo a SECOND CID cascade triggered by Steve's carrier-line -> grid-envelope correction -- replaying the strict-reader contradiction first surfaced in UT-169.b.
- **UT-175.b** (T175): TE-38 scope doubled in one turn (5 DFs in turn 172 -> 10 DFs in turn 175); bot's single-TE lean carries real coupling-vs-decomposition risk and Steve's answer to meta-question 1 is the gating decision...
- **UT-175.c** (T175): Bot's recommendation to use "headers" for the From/Date/Parents block (replacing inner-block "envelope") imports email/HTTP convention but those have MIME-loaded specifications wire-lab may not want t...
- **UT-175.d** (T175): Steve's joint endorsement of "feed" + "substrate" + "forum" + "site" pulls the wire-lab vocabulary firmly into the Usenet/FidoNet lineage and ratifies UT-173.a's wire-lab-as-content-addressed-Usenet f...
- **UT-175.e** (T175): (no bolded summary)
- **UT-175.f** (T175): Bot's silent reversal of the encapsulation-not-layering frame is itself a procedural concern, paralleling UT-169.a/c/e and UT-167.e/168.e/169.c.
- **UT-175.g** (T175): Steve's joint endorsement "Feed is okay, but you also used substrate, which I also like" is structurally underspecified -- bot did not push back to clarify which word is canonical-in-spec-text vs desc...
- **UT-176.a** (T176): The slug-state naming convention's affected scope is wider than bot enumerated -- bot identified only `transports/draft--wire-lab-devs/` as needing rename, but the rule "<slug>-<state>" should be chec...
- **UT-176.b** (T176): Layer numbering not fully resolved -- bot offered three options (internal 1/2/3, OSI 8/9/10, original 0-3 corrected to 4-7) and recommended internal 1/2/3 but Steve's "4-7" suggestion was not directly...
- **UT-176.c** (T176): "groups" vocabulary collides with `group-session` protocol family name -- bot acknowledged in turn 175 ("meshes with group-session but conflates protocol-name with directory-name") that this risks con...
- **UT-176.d** (T176): PromiseGrid-as-Usenet+git successor framing changes the project's stated scope and goals -- this may be a TE-by-itself decision rather than background material that TE-38 inherits.
- **UT-176.e** (T176): The body-can-be-anything realization affects §4.7 strict-reader rule again -- if body can be a nested grid envelope, strict readers must recursively parse, which is a third strict-reader-rule complica...
- **UT-176.f** (T176): Foundational CID promise as constitutive vs voluntary -- bot characterizes it as "a site cannot be a CAS without making this promise" but Promise Theory typically treats promises as voluntary, not con...
- **UT-176.h** (T176): Symlinks-into-CAS deferred to later TE creates a future-coupling problem -- if TE-38 lands without symlinks and a later TE adds them, the transition requires moving messages from groups/<slug>-draft/<...
- **UT-177.a** (T177): Layer inversion contradicts the turn-175 layered-stack model that was just locked in.
- **UT-177.c** (T177): CBOR adoption is now baked into TE-38 (DF-38.11) but the implications cascade further than the DF text acknowledges.
- **UT-177.d** (T177): Rabin chunking creates a third potential CID cascade beyond the foundational-CID-promise concern.
- **UT-177.e** (T177): Promise economy reframed from background framing to foundational invariant -- the practical implication is that every protocol spec must now contain a promise vocabulary section.
- **UT-177.f** (T177): The 100-year goal is now an explicit design constraint but no DF operationalizes it.
- **UT-177.g** (T177): Bot's "game it out" reversal on the layer order was explicit and well-handled, contrary to the cadence concern of UT-175.f.
- **UT-177.h** (T177): Pointer files plus chunking together make the existing `<cid>.txt` file structure structurally wrong, not just terminologically wrong.
- **UT-177.i** (T177): "Easy mental models" is now a stated 100-year requirement but no DF treats it.
- **UT-178.a** (T178): Sparse-CAS is now a foundational design assumption but no DF treats it.
- **UT-178.b** (T178): The 'decides' word is the load-bearing point of the entire layered model and bot's turn-177 description glossed it.
- **UT-178.c** (T178): Layered accounting question is unanswered: does promise accounting live at L7, or per-layer at L5/L6/L7?
- **UT-178.d** (T178): PromiseGrid-replaces-BGP is a major application Steve flagged as future consideration.
- **UT-178.e** (T178): TE-38 opening section is now mandated to be the capture narrative, not background.
- **UT-178.f** (T178): Anonymity-is-pathology is now a stated design stance with implications for group identity.
- **UT-178.g** (T178): Dogfooding pressure now constrains TE-38's deliverable.
- **UT-178.h** (T178): Promisebase surfacing is the largest single context expansion of the session and changes TE-38's nature.
- **UT-178.i** (T178): Interop constraint with libp2p/IPFS/ATPROTO restricts CIDv1/multihash/multibase choices.
- **UT-178.j** (T178): Multi-repo question is unanswered and structurally significant.
- **UT-178.k** (T178): Bot's turn-178 answer was unusually terse (4 lines, three meta-questions) -- a cadence inversion from prior turns.
- **UT-178.l** (T178): CIDv1 codec field for merkle-vs-content distinction is the right answer to UT-177.h.

### TE-40: Apparatus-vs-specimen completion (18 UTs)

- **UT-155.a** (T155): DF-1.1 (TE-1 promise-stack ordering, Alt-E projection-mode) was never locked.
- **UT-155.b** (T155): The most explicit `Project / Peel / Wrap` operation definitions in the corpus live only in the conversation, not in any committed file.
- **UT-156.a** (T156): The Option 1 / Option 2 / Option 3 question for TE-1's structural role was abandoned, not answered.
- **UT-156.b** (T156): TE-1 is miscategorized in the corpus and the misclassification has not been repaired.
- **UT-156.c** (T156): The bot's "harness-spec is wire-envelope-agnostic" claim is wrong and must be retracted, not committed.
- **UT-157.a** (T157): The five candidate envelopes (Env-1 through Env-5) are named in the conversation but uncommitted.
- **UT-157.b** (T157): Reading 1 / Reading 2 / Reading 3 question for TE-1's status was abandoned, not answered.
- **UT-157.c** (T157): The "`grid([pcid, payload])` is the working hypothesis but not yet proven" framing is asserted but uncommitted.
- **UT-158.b** (T158): Three apparatus-vs-specimen TE scopes were named but never picked.
- **UT-158.c** (T158): Harness-spec §1.1 carve-out is identified but not yet executed.
- **UT-158.d** (T158): Harness-spec §1.3 invariants classification (apparatus vs specimen) was named but not committed.
- **UT-158.e** (T158): TE-1 should move from `docs/thought-experiments/` (harness-level) to a promise-stack protocol directory (specimen-level).
- **UT-158.f** (T158): `grid([pcid, payload])` needs its own protocol directory and spec doc.
- **UT-158.g** (T158): TODO 5 should move to the promise-stack protocol's TODO directory and be reframed.
- **UT-159.a** (T159): The audit at `4725b3e` identified NINE specimen-bearing items in harness-spec; UT-158.c/d only flagged §1.1 and §1.3.
- **UT-159.b** (T159): `protocols/wire-lab.d/specs/transport-spec-draft.md` companion audit was flagged but not done.
- **UT-159.c** (T159): The audit memo recommended six tabletop scenarios for the apparatus-vs-specimen step-2 TE; TE-36 (turns 160-167) did not use them.
- **UT-159.d** (T159): Two of the three audit-flagged ambiguous areas have not been resolved by any TE.

### TE-41: Group-session freeze procedure (15 UTs)

- **UT-164.c** (T164): Genesis message m000 covers §3/§4/§5/§7 of group-session spec but NOT §4.6 Parents: or §6 body-as-receipt; turn-165's recap of the freeze gate requires §3/§4/§4.6/§6/§7.
- **UT-164.d** (T164): The membership-by-posting convention -- "the membership list is whatever set of From: values appears in the transport's commit history before the spec is frozen" -- is a load-bearing semi-formal rule ...
- **UT-164.e** (T164): The rename-at-freeze operation -- "directory will be renamed from `transports/draft--wire-lab-devs/` to `transports/<pcid>--wire-lab-devs/` in a single mechanical commit, and every message's carrier l...
- **UT-165.b** (T165): `merge-group-transport-spec` signature mechanic is referenced in the freeze sequence but never operationally defined.
- **UT-166.a** (T166): Bot executed the bootstrap commit (`a1c85f3`) and merge (`a1ecc72`) AS PART OF turn 166's response, without a separate Steve "yes" turn between recommendation and execution.
- **UT-166.b** (T166): Bot's answer text named `stevegt-codex (presumably)` as a likely wire-lab-devs participant -- but Steve never authorized assuming Codex is enrolled in this transport instance.
- **UT-166.c** (T166): Bot's answer text uses `stevegt-perplexity` as the bot's identifier; the committed `From:` line in m000/m001/m002 uses `stevegt-via-perplexity` (with "via").
- **UT-166.e** (T166): Bot raised the §8 membership-pinning question ("closed and fixed at transport creation means we should pin the participant set up front") and then declined to address it; the bootstrap landed in the s...
- **UT-167.a** (T167): The sequential-numbering convention `m<N>-<author-id>-<utc>-<slug>.txt` was codified in committed spec text and exercised in m001 ONE TURN BEFORE Steve's turn 168 directive will reverse it ("do not se...
- **UT-167.b** (T167): Spec §9 (per-author-branch binding) silently resolves the §8 membership question with extensible-by-posting -- the same raise-and-silently-resolve pattern as UT-166.e but now codified in committed spe...
- **UT-167.c** (T167): Steve's prompt "Members should fetch from all branches but only post on {name}/main" carries an ambiguity in `{name}` that §9 silently resolved one way without flagging the other.
- **UT-167.e** (T167): Same raise-and-resolve / execute-without-explicit-yes pattern as turn 166: Steve's prompt was a terse directive ("Use .txt instead of .msg. Members should fetch...") and the bot landed three commits (...
- **UT-168.b** (T168): §9 of group-session-draft.md grew from a one-paragraph non-normative note (committed in turn 167 at `463af44`) to an eight-subsection (§9.1-9.8) detailed binding spec at `19e9b37`, while staying label...
- **UT-168.c** (T168): The "merge-then-optionally-post" cycle adds a third state to the previously-binary post / no-post participation model, and the consequences for non-merging readers were not surfaced.
- **UT-168.d** (T168): §9.5's "infrastructure files NOT propagated by merge" rule introduces a new file-class distinction that is not formalized: which files are infrastructure?

### TE-42: Filename/CID-cascade policy (7 UTs)

- **UT-168.a** (T168): The Message-ID header was explicitly retained in turn 168's spec rewrite ("the Message-ID: header (§4.3) is retained as a human-readable convenience inside the envelope but is not the filename") in co...
- **UT-168.e** (T168): m2 (CID-filename + merge-cycle ratification, then-CID `bafkreihnvl...`, current-CID `bafkreidef4b4qdc4xjvkjrern7jm4ta75q55ed2u2ilwcrkxqhn7n4fjce` after Message-ID removal) was authored without an expl...
- **UT-168.f** (T168): All three CIDs claimed by turn 168's answer text and recorded in commit messages (`bafkreifmjs...`, `bafkreihqx6...`, `bafkreihnvl...`) were retroactively invalidated by turn 169's Message-ID removal ...
- **UT-169.a** (T169): The reasoning text explicitly recommended Path A ("deprecate but not forbid"; preserve legacy bytes; "Path A is the right move") -- the commit executed Path B (hard remove + rehash all three legacy me...
- **UT-169.b** (T169): The bot's own reasoning flagged §4.7's "strict readers MUST reject messages they cannot fully parse" rule as a blocking concern ("This is a hard problem. If we remove Message-ID and a strict reader re...
- **UT-169.d** (T169): Twig name `ppx/transport-remove-message-id` breaks the prevailing `ppx/te-<utc>-<slug>` pattern used by other twigs in this session.
- **UT-169.e** (T169): §4.3's rewrite changes Message-ID from "mandatory" to "prohibited" -- a stronger enforcement than the deprecation the reasoning recommended.

### TE-43: Promisebase prior-art adoption (25 UTs)

- **UT-179.c** (T179): Conditional-release promises ("send only if recipient promises onward-restraint") imply a recursive promise-graph that the protocol must track.
- **UT-181.a** (T181): Bot's import-path claim github.com/t7a/pitbase/db is wrong; correct path is github.com/stevegt/promisebase/db.
- **UT-181.b** (T181): Pitbase's Rabin chunking defaults (512 KiB min / 8 MiB max) are dramatically larger than turn-177's recommended FastCDC ~16 KiB average; the parameter mismatch is now an unresolved DF.
- **UT-181.c** (T181): The two failing tests (TestPutStreamBig, TestPutStreamSmall) are in the chunker-merkle-builder integration -- the exact path wire-lab needs.
- **UT-181.d** (T181): The pitbase worm class\n header pattern ("hash includes class header byte") is structurally analogous to the foundational CID promise but at the file format layer.
- **UT-181.e** (T181): Pitbase's Stream-as-symlink design is the prior-art-prototype of the pointer files concept from turn 177.
- **UT-181.f** (T181): Bot's offer to "dig into chunker_test.go:106" was an in-line decision-point that Steve answered terse-positive in turn 182 ("That used to work"); cadence preserved.
- **UT-182.b** (T182): Bot's "open a PR / leave the commit ready for you to push" phrasing is procedurally loose given bot has no promisebase PAT yet.
- **UT-183.a** (T183): Bot prepared a commit (d98b5d3) on a local-only promisebase clone with no persistence path if the session reset before PAT grant.
- **UT-183.b** (T183): Bot's recommendation "option 2: read more of promisebase before drafting TE-38" sets a new procedural default that should be made canonical.
- **UT-183.c** (T183): The 17/17 test result establishes pitbase's chunker-merkle-builder integration path as operationally green and unblocks the L6-substrate working framing for TE-38.
- **UT-184.d** (T184): The kv/fs/ refactor (Oct 2025, scanner-to-Options-struct) interacts non-trivially with the L6-substrate framing -- if kv/fs/ becomes pitbase's new bottom layer with db/'s CAS logic moving up, the wire...
- **UT-184.e** (T184): The Docker SDK rot in cmd/pb/ (6 errors against newer Docker SDK, types renamed in Docker v28) is a second regression that may also need fixing on a separate twig before TE-38 drafting.
- **UT-184.f** (T184): RFC-1005 Option 2 (content-addressable test-driven fabric: test tree CID + executable tree CID + args; cache only when test passes) is the explicit prior-art seed of the promise-economy and is dated J...
- **UT-184.g** (T184): The fuse/ test failures and the cmd/pb Docker rot together suggest pitbase's main is currently in a partial-rot state -- multiple modules are broken, only db/ and kv/fs/ are green.
- **UT-187.c** (T187): Bot's proposed twig name promisebase-adoption-and-federation-layer omits two of the three substantive scope axes from the proposed title.
- **UT-190.a** (T190): Bot reported remote-branch enumeration as a narrated finding ("only main exists on the remote") rather than as raw command output, and the narrated finding was wrong -- the kv branch on stevegt/promis...
- **UT-190.b** (T190): Bot concluded "the TE-38 plan I laid out doesn't need to account for any active collaborator branches in promisebase" based on a wrong enumeration -- the kv branch contains 7 commits of in-progress wo...
- **UT-190.c** (T190): Steve's question used the plural "branches" -- "Can you examine the other promisebase branches?" -- which itself implied an expectation that multiple branches exist; bot answered "only main exists" wi...
- **UT-190.d** (T190): The kv branch on stevegt/promisebase is undiscovered in the conversation history through turn 192; turn-190 walk-time enumeration is the first record of its existence in this collaboration's working m...
- **UT-191.a** (T191): Bot's response to Steve's canon rule ("prototype at best; prefer wire-lab in conflict") was a scope-level self-correction ("My earlier scope was tilted toward 'adopt promisebase wholesale.' That was w...
- **UT-191.c** (T191): Bot's revised twig name `ppx/te-20260504-NNNNNN-promisebase-as-prototype-source` is on the long side relative to the apparatus-vs-specimen kebab-case standard from prior twigs.
- **UT-191.d** (T191): Bot's turn-191 response treats kv/fs/ as "mid-refactor" within main's tree but does not know about the separate kv branch on the remote (UT-190.d).
- **UT-192.b** (T192): Bot's turn-192 reframe still does not surface the kv-branch finding (UT-190.d / UT-191.d).
- **UT-192.f** (T192): Bot's turn-192 framing introduces the implicit commitment that wire-lab and promisebase "will gradually merge" ("wire-lab and promisebase will be in the same codebase eventually") -- a major architect...

### TE-45: Conditional-release / geofencing / recursive promise-graph (1 UTs)

- **UT-179.d** (T179): Geofencing requirement adds a constraint dimension orthogonal to group membership.

### Spec-edit (small direct edits, no TE) (5 UTs)

- **UT-160.a** (T160): Bot's answer text said the audit identified "eight specimen-bearing items"; the actual audit at `4725b3e` lists NINE items.
- **UT-163.a** (T163): The four envelope-agnostic rephrasings of §1.3 simulator tests constitute a vocabulary template that lives only in the bot's reply text, not in any committed file.
- **UT-165.d** (T165): OQ-G4 ("Should there be a canonical 'transport-creation' or 'genesis' message at the root of every transport's DAG?") was implicitly answered "yes" by construction for the wire-lab-devs instance witho...
- **UT-165.e** (T165): Group-session-draft.md still uses `codex-perplexity` in five places as the example/canonical slug, but the live transport instance is `wire-lab-devs`.
- **UT-167.d** (T167): The .msg->.txt sweep across six unfrozen documents was performed without the cross-TE quotation grep that DF-35.3 / TE-34 Cat-2 sweep policy requires.

### Retire (superseded; no action) (3 UTs)

- **UT-170.b** (T170): "DF-37.1" uses TE-37 numbering for a TE document that does not exist.
- **UT-170.c** (T170): The continuity-summary "completed items" list at the top of recent session-summaries shows three items as completed -- "Create protocols/git-file-transport.d/ skeleton," "Draft specs/git-file-transpor...
- **UT-170.d** (T170): DF-37.1 framing is structurally incomplete: it asks "how should the protocol-slug level be added?" but the actual missing dimension -- the substrate axis (git, rsync, uucp, udp, ...) -- is not surface...

### Carry (procedural / AGENTS-rule / cadence notes) (55 UTs)

- **UT-158.a** (T158): The bot's six-step plan was offered as an ordered sequence with dependencies; only step 1 (audit) executed before the conversation pivoted into TE-36 directly.
- **UT-158.h** (T158): A parallel TODO for the grid hypothesis (same shape as the reframed TODO 5) was named but not filed.
- **UT-160.d** (T160): At end of session corpus, five of TE-36's seven DFs remain unlocked.
- **UT-161.b** (T161): The nine-axis comparison table between promise-stack and grid-pcid-payload lives only in conversation, not in any committed file.
- **UT-161.c** (T161): Bot's turn-161 list of `assertion` types is a finer-grained taxonomy than the five named pCIDs from harness-spec section 10a.
- **UT-163.b** (T163): FOUR commits on the TE-36 twig land in conversation-file gaps with no captured Steve authorization turn.
- **UT-164.a** (T164): Bot's pause text said "a 4-DF-locked TE" but only 2 DFs were actually locked at the time.
- **UT-164.b** (T164): Turn 165's answer text said TE-37 git-file-transport twig was "reverted cleanly to 4725b3e" but the reflog shows the twig was MERGED into ppx/main as `a1ecc72` and then deleted.
- **UT-165.a** (T165): The bot's slug-naming proposals violated the collaborator-anonymity rule the same turn the rule was installed.
- **UT-165.c** (T165): The bot's "Saving [redacted-collaborator] collaborator context" memory write inferred a pronoun ("she") not present in Steve's prompt.
- **UT-166.d** (T166): Stale branch name `ppx/te-20260503-112348-git-file-transport` was used for the bootstrap commit and is preserved in reflog and merge-commit metadata.
- **UT-169.c** (T169): The reasoning text explicitly ended with intent to ASK Steve before acting ("Let me give him the reasoning and a recommendation, with one short clarifying question (Path A vs Path B), so he can answer...
- **UT-171.c** (T171): Bot's caveat (b) "defer until a second binding actually exists" applies a YAGNI rule whose triggering condition turn 172 immediately satisfies -- so the YAGNI deferral lasts exactly three minutes.
- **UT-171.d** (T171): Bot's recommendation question at the end ("Want me to proceed with Alt-1.A as the move, or does the substrate question still feel unresolved?") is structurally a two-part question that foregrounds Alt...
- **UT-176.g** (T176): TE-38 now has 12 DFs and bot has asked three meta-questions in three consecutive turns (174/175/176) -- the DF list keeps growing because Steve's cadence is correction-plus-elaboration, not draft-then...
- **UT-176.i** (T176): The duplicate `021.176` line in this TODO was introduced by MY OWN turn-175 walk commit `47009f1`, not by a "prior continuity summary defect" as turn-176 commit message `09f92d7` claimed.
- **UT-177.b** (T177): TE-38 DF count has grown 5 -> 10 -> 12 -> 15 across four turns.
- **UT-178.m** (T178): Procedural defect during the turn-178 walk -- spurious duplicate forward-preview block introduced and self-corrected.
- **UT-179.a** (T179): Bot's wholesale-adoption pivot was based on reading design documentation, not implementation code, and bot acknowledged the doc-vs-code uncertainty in passing without acting on it.
- **UT-179.b** (T179): The promise-economy spectrum (pure-social trust scoring to typed-fungible-capability-token marketplace) is now stated explicitly but the protocol-stays-agnostic constraint is not yet captured as a DF.
- **UT-179.e** (T179): Bot's TE-38 sketch named "DF-38.7: Vocabulary fixes (grid envelope vs carrier line, groups vs transports, slug-state naming)" as a single DF, but each of these is a separate locked correction from ear...
- **UT-179.f** (T179): Bot's three-meta-questions ending recovers the cadence pattern that was inverted in turn 178.
- **UT-179.g** (T179): (no bolded summary)
- **UT-180.a** (T180): Bot's turn-179 wholesale-adoption pivot was a structurally bad design decision and the corrective procedure is now established as canon.
- **UT-180.b** (T180): The two-paths framing (shared-vocabulary vs read-code-first) under-represents a third path that Steve never selected explicitly: independent design with optional convergence.
- **UT-180.c** (T180): The pending-line summary for 021.180 was wrong about the conversational dynamic.
- **UT-180.d** (T180): Bot's apology language ("You're right -- I owe you an apology") is the right tone but the procedural pattern of "apologize, audit, invalidate, propose paths" should be elevated to a standing response ...
- **UT-182.a** (T182): Bot's turn-182 answer presupposes diagnostic work performed between turns 181 and 182 that is not shown in the answer text.
- **UT-182.c** (T182): Steve's "That used to work" carries two implicit confirmations -- regression hypothesis confirmed AND implicit yes to dig in -- and the bot correctly inferred both.
- **UT-184.a** (T184): Bot's turn-184 answer lists ten DFs flat without recommendations, Cat-1/2/3 framing, or consideration paragraphs -- violating the standing one-at-a-time DF discipline.
- **UT-184.b** (T184): Bot's question to Steve mentions "[redacted-collaborator]" by name alongside other 2021 collaborator names (Matt, Angela, Jessica, Ryan) -- borderline collaborator-non-mention rule violation.
- **UT-184.c** (T184): Bot's claim "x/discussion.md is dumped grokker chat sessions (153 occurrences of grokker boilerplate in 5466 lines)" treats "grokker boilerplate" as a known marker without defining what it looks like.
- **UT-185.a** (T185): Bot's PAT-handling response did not mention fine-grained-PAT expiry as a defense-in-depth mitigation against the carry-over-context-summary verbatim-echo concern.
- **UT-185.b** (T185): Bot's offer to mark tokens read-only via filename suffix (gh-pat-promisebase-readonly) is a convention not enforced by the filesystem -- a misnamed file silently violates the convention.
- **UT-185.c** (T185): Bot's caveat (1) acknowledges PAT verbatim-echo in carry-over context summaries but does not propose a redaction mechanism; the carry-over summaries continue to carry live PAT bytes.
- **UT-186.a** (T186): First walk note applying the UT-185.c PAT-redaction discipline -- Steve's turn-186 query is the literal PAT bytes; this walk note redacts them as `{{SECRET:gh-pat-promisebase}}` rather than echoing ve...
- **UT-186.b** (T186): Bot's GOTOOLCHAIN=auto diagnostic is correct, useful, and not written down anywhere persistent yet.
- **UT-186.c** (T186): Bot's pivot-back paragraph re-quotes turn 184's question that named [redacted-collaborator]; the collaborator-mention propagates from the original turn into a follow-up turn.
- **UT-187.a** (T187): Bot's three confirmation questions Q1/Q2/Q3 in turn 187 were never answered because turn 188 hijacked the thread with the push-status check.
- **UT-187.b** (T187): Bot's Q2 ("shorter or fuller per-DF recommendations") asks Steve to override the standing one-DF-at-a-time-with-paragraphs rule rather than applying it.
- **UT-187.d** (T187): Bot's parenthetical "placeholder phrasing -- not naming [redacted-collaborator]" is the first time bot has explicitly named the collaborator-non-mention discipline in a chat answer; the meta-acknowled...
- **UT-188.a** (T188): Visibility-of-status gap -- bot's turn-186 push of the randStream fix was buried in a side-channel observation paragraph rather than foregrounded as a DONE confirmation, causing Steve to ask 40 minute...
- **UT-188.b** (T188): Bot's mid-sentence self-correction on the force-push rule ("d98b5d3 -> reset main back to aedaac2, force-push... actually wait, you said 'never force push' -- so the cleaner revert would be a git reve...
- **UT-188.c** (T188): Bot's offer-to-revert ("if you'd rather I had left it on the twig for you to review/merge, say the word and I'll revert") is a healthy collaboration pattern for unilateral cross-repo state changes.
- **UT-189.a** (T189): Bot's response to a truncated/fragment message was "Answer skipped" with no acknowledgement to Steve that the fragment had been received.
- **UT-189.b** (T189): The string "Answer skipped" cannot be definitively attributed to bot-generated text vs platform/transcript convention from the conversation log evidence alone.
- **UT-189.c** (T189): Walk-note convention for fragment+completion turn-pairs: treat the fragment as procedural-only with cross-reference to the completion turn for substance.
- **UT-191.b** (T191): TE-38's DF count reduction (10 -> 6, or 5 with fold) is the structural consequence of bot's scope-level self-correction in UT-191.a -- and the first time in this session's TE-38 history that the DF co...
- **UT-191.e** (T191): Bot's claim "RFCs 1003-1007 are 5 years old and predate the promise framing" should be ground-truthed before TE-38 cites it.
- **UT-191.f** (T191): Bot interpreted Steve's "discuss the conflict ... should prefer wire-lab" as requiring documentation of the discussion in the relevant TE.
- **UT-192.a** (T192): Bot performed its third TE-38 reframe in three consecutive turns (187 wholesale-adoption -> 191 salvage-source -> 192 active-prototype-graduating-into-promisegrid), each time re-deriving the DF list f...
- **UT-192.c** (T192): Steve's word "ref" in "ref, factor, modernize, and use" is ambiguous and bot did not flag the ambiguity.
- **UT-192.d** (T192): Multiple DF lists for TE-38 coexist in the conversation without explicit supersession.
- **UT-192.e** (T192): Twig-name proposals for TE-38 have accumulated across turns without resolution.
- **UT-192.x** (T192): The pre-walk placeholder line for 021.192 conflated turns 192 and 193 into a single quotation.

---

## Other active TODOs (open items not in TODO 21)

In addition to TODO 21's 186 UTs, five other TODO files have open items totaling ~24. These pre-date the TE-38 cluster (most are from April 29) and reference earlier TE work (TE-1, TE-21, TE-22).

### TODO-20260429-164955-te-promise-stack-ordering.md (7 open)

DF-1.1 through DF-1.4 (TE-1 promise-stack ordering decisions). All four DFs await Steve's "yes" — none locked.

**Disposition:** **TE-40 / Apparatus-vs-specimen completion.** TE-1 is in limbo per UT-155.a/156.a/156.b/158.e. TE-40 should explicitly close TE-1 — either lock the four DFs as Alt-E hybrid (bot's recommendations) or formally retire TE-1 and move promise-stack-ordering work to a successor TE.

### TODO-20260429-165252-backfill-di-provenance-harness-spec.md (4 open)

Four steps to backfill DI (Decision Implementation) provenance into harness-spec for settled statements that lack it.

**Disposition:** **TE-40 / Apparatus-vs-specimen completion.** The harness-spec sweep that TE-40 would do should also do this DI-provenance backfill in the same pass — both touch the same file, both classify the same content.

### TODO-20260429-165253-backfill-dr-harness-spec-section-11.md (4 open)

Four steps to lift harness-spec §11 open-questions into DR (Decision Register) entries.

**Disposition:** **TE-40 / Apparatus-vs-specimen completion.** Same rationale: ride along with the harness-spec sweep.

### TODO-20260429-173837-te-spec-doc-as-promise.md (8 open)

DF-21.1 through DF-21.4 (TE-21 spec-doc-as-promise) plus TE-22 / TE-23 placeholder follow-ons. All four DF-21 decisions await Steve's "yes" — none locked.

**Disposition:** **Carry / Re-evaluate.** TE-21 is older work whose conclusions may be partially absorbed by the TE-38 substrate-agnostic layered model and the eventual freeze procedure (TE-41). Bot recommends: do not pick up TE-21 work until TE-38 lands; then re-evaluate whether TE-21 is still coherent or has been overtaken.

### TODO-20260429-180020-te-spec-doc-store-and-pcid-machinery.md (1 open)

Item 011.10 — follow-on TE on peer-level adoption metadata. "Missing half of TE-21 Alt-E that TE-22 did not address."

**Disposition:** **Carry.** Same logic as TE-21 above — wait for TE-38 to land.

---

## Phase 2 next step

Steve reviews this file, corrects dispositions where the bot got it wrong, and approves the TE drafting order. Then TE-38 drafting begins (already-accumulated 30+ DFs from turns 172-178; needs scope reduction to 3-5 anchor DFs per the apparatus-vs-specimen / TE-editing-policy precedent).
