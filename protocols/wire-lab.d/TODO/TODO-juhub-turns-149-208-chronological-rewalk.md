# TODO-juhub: Turns 149-208 chronological rewalk

## Prior aliases

This TODO is being filed after the TE-mumuv proquint migration locked
2026-05-07, so it is minted directly under `TODO-juhub`. No prior integer or
timestamp alias.

## Status

Open. This TODO owns the one-turn-at-a-time raw-log rewalk for turns 149-208.
`TODO-lilar` remains the historical original walk through turn 192, and
`TODO-jivam` remains the closure gate for overall recovery completion. No turn
may advance until the current turn has been analyzed, reported, and explicitly
approved.

## Decision Intent Log

ID: DI-nagat
Date: 2026-05-10 17:11:45
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: File `TODO-juhub` as the authoritative one-turn-at-a-time rewalk
ledger for turns 149-208. For each turn, read the raw session log, sweep all
later turn logs and later repo artifacts for whether that turn's questions,
decisions, or plans were later settled or changed, report the result, and stop
until explicitly approved to continue.
Intent: Rebuild confidence in the 149-208 recovery from raw evidence instead
of trusting earlier summaries, while keeping the historical walk and the
closure monitor distinct and readable.
Constraints: Do not use `TODO-jivam` as the per-turn ledger. Do not rewrite
existing `TODO-lilar` turn notes; if a correction is needed there, append a
provenance-bearing correction note. Do not advance more than one turn without
explicit approval. Every unresolved finding must be handed to a proper owner
artifact.
Affects: `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`protocols/wire-lab.d/TODO/TODO-jivam-turns-149-170-recovery-completion.md`;
`protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`;
`protocols/wire-lab.d/TODO/TODO.md`;
`/home/stevegt/lab/session-logs/sessions/ea135ce8/149-turn.md` through
`/home/stevegt/lab/session-logs/sessions/ea135ce8/208-turn.md`.

ID: DI-pijun
Date: 2026-05-10 19:37:26
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Promote the turn-173 Usenet-lineage insight into two non-TE docs
with distinct roles: the research note keeps the historical-precedent
evidence, and a new simulation README explores the broader CAS /
Usenet-like / git-like design line. `group-session` is one current specimen in
that broader line, not the identity of the simulation itself.
Intent: Preserve the strongest historical analogy in a design-visible place
without pretending it is already a frozen protocol decision, while giving the
broader architectural framing a simulation-local home that can grow beyond the
current `group-session` specimen.
Constraints: Do not file a new TE for this promotion. Put the research-doc
wording in the Usenet section of
`docs/research/historical-networks-20260503.md`. Use the phrase
`content-addressed Usenet`. Create `simulations/SIM-hugoj-cas-usenetlike-gitlike/README.md`
only; do not create a fuller simulation scaffold in this pass. Add the new
simulation to `simulations/README.md`.
Affects: `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`docs/research/historical-networks-20260503.md`;
`simulations/README.md`;
`simulations/SIM-hugoj-cas-usenetlike-gitlike/README.md`.

ID: DI-gudap
Date: 2026-05-10 13:16:35
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Strengthen the `TODO-juhub` replay discipline so each turn
reconciliation also resolves or explicitly calls out its related loose ends.
Intent: Prevent the 149-208 replay from becoming a second narrative-only pass.
The rewalk must reduce relitigation by reconciling each turn's `TODO-lilar`
`UT-*` fallout against the current owner TODOs and by promoting any
load-bearing insight that is still trapped only in replay artifacts.
Constraints: Do not flip `TODO-lilar` checkboxes directly. Respect the
matrix-as-closure-index rule in `ut-verification-matrix-20260507.md`; close,
retire, or transfer loose ends only in the proper owner artifact. Update
`TODO-lilar` only with additive correction notes when its historical walk note
is actually wrong. Direct design/spec/research/simulation docs are touched only
when a specific turn exposes a load-bearing statement that is still missing
from those docs. Update `DEV-GUIDE-RESOURCES.md` only when such a cited or
relevant source changes.
Affects: `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`protocols/wire-lab.d/TODO/TODO-jivam-turns-149-170-recovery-completion.md`;
`protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`;
`protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md`;
`protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`;
current owner TODOs for the turns being replayed; directly implicated
spec/design/research/simulation docs when needed.

ID: DI-vanak
Date: 2026-05-10 13:36:01
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: During the `TODO-juhub` replay, Steve's single-word reply `turn`
counts as approval of the currently presented turn analysis and authorizes
rewriting that turn's note in the stronger `TODO-juhub` format before
advancing.
Intent: Remove repetitive confirmation chatter so the replay can proceed one
turn at a time while still preserving the explicit approval boundary.
Constraints: `turn` approves only the currently presented turn; it does not
skip turns. It authorizes rewriting the current turn's `TODO-juhub` note and
any already-described turn-local owner, correction-note, or direct-doc cleanup
required by that approved analysis. The replay still stops after the next turn
is presented for review.
Affects: `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
this session's `TODO-juhub` replay cadence.

ID: DI-nijod
Date: 2026-05-10 13:45:42
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Every per-turn replay response must include an explicit `Work
pending` line near the bottom.
Intent: Make it immediately obvious whether a turn still leaves live `UT-*`,
owner TODO, spec-edit, DR/DI, or other substantive follow-on work, so the
replay steadily burns down loose ends instead of forcing Steve to infer status.
Constraints: `Work pending: yes` when any open `UT`, owner TODO item,
spec-edit, DR/DI, or other substantive work still stems from the turn.
`Work pending: no` only when the turn is fully reconciled and no live owner
work remains because the related loose ends are absent, closed, retired, or
transferred. This line is a report field; it does not itself close work.
Affects: `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
this session's `TODO-juhub` replay responses.

## Why a separate TODO

`TODO-lilar` already records the original historical walk through turn 192 and
explicitly says turns 193+ belong in a successor TODO. `TODO-jivam` is the
bounded recovery closure gate and should stay focused on completion criteria
instead of becoming a long per-turn ledger. `TODO-juhub` therefore owns the
rewalk mechanics while the older artifacts remain evidence and closure logic.

## Per-turn discipline

1. Read the raw session log for the current turn.
2. Sweep every later turn log and every later relevant artifact (`TE`, `DI`,
   `DR`, `TODO`, specs, matrix/disposition docs, essays, and research notes)
   for whether the current turn's questions, decisions, or plans were later
   settled, corrected, superseded, contradicted, or abandoned.
3. Collect every related loose end for that turn from `TODO-lilar`,
   `dropped-thread-disposition-20260506.md`,
   `ut-verification-matrix-20260507.md`, and the current owner TODOs.
4. Determine whether each loose end is already resolved, retired, transferred,
   still open, or missing a proper owner.
5. Compare the turn plus its loose ends against the existing replay artifacts
   and decide whether any historical correction note is actually needed.
6. Report the turn back to Steve using the fixed report format below and stop.
7. After approval, write the turn result here and make any linked owner,
   correction-note, or direct-doc updates required by the turn.
8. Do not advance until each related loose end is either closed, retired, or
   transferred in its proper owner artifact, or explicitly called out here as
   still needing a named decision or work item.

## Interaction shorthand

- `turn` means: approve the currently presented turn, rewrite that turn's note
  in the stronger `TODO-juhub` format, perform any already-described turn-local
  cleanup authorized by the approved analysis, then present the next turn.
- The `Write needed? yes/no` line in the turn report is informational. During
  this replay, it does not require a separate confirmation from Steve before
  rewriting the current turn note in `TODO-juhub`.

## Turn report format

- `Turn N summary`
- `Existing capture`
- `Gaps or contradictions`
- `Related UTs / owners`
- `Owner/doc cleanup`
- `Remaining decisions or work`
- `Work pending: yes|no`
- `Proposed disposition`
- `Write needed? yes/no`
- `Next: wait for approval before turn N+1`

## Loose-end backfill through turn 174

- `Turns 149-154` No related `UT-*` entries were filed for this block. No
  downstream owner or direct-doc cleanup is currently needed beyond the
  already-landed TE-editing-policy and TODO-020 artifacts.
- `Turn 155` `UT-155.a` and `UT-155.b` are now retired in `TODO-kugod` /
  `TODO-rivuk` under `DI-runuh`; no live turn-local owner work remains.
- `Turn 156` `UT-156.a` and `UT-156.b` are retired / resolved-retired.
  `UT-156.c` is resolved under `DI-lajod`, `DI-sujan`, and `DI-kinad`;
  broader TE-40 audit work continues under the later turn-159 rows.
- `Turn 157` `UT-157.b` is retired. `UT-157.a` and `UT-157.c` are resolved in
  the grid-envelope successor owner under `DI-joroh`.
- `Turn 158` `UT-158.b` is resolved; `UT-158.e` and `UT-158.g` are retired;
  `UT-158.a` is resolved-decomposed and `UT-158.h` is resolved-routed under
  `DI-sotuk`; `UT-158.c`, `UT-158.d`, and `UT-158.f` are closed for
  turn-158 scope under `DI-kinad` and `DI-fanah`. Broader TE-40 work begins
  with turn 159 and remains under `UT-159.*`.
- `Turn 159` `UT-159.c` is resolved-retired. `UT-159.a` and `UT-159.d` remain
  open in `TODO-kugod` via `kugod.5`. `UT-159.b` remains open-scoped under
  `TODO-kugod` plus `DR-nugog` via `kugod.8`.
- `Turn 160` `UT-160.b` and `UT-160.c` are answered by the TE-havib follow-on
  verification path. `UT-160.a` remains a spec-edit loose end, and `UT-160.d`
  remains a carry note about the still-unlocked TE-havib DFs. `TODO-lilok`
  still needs its reopened harness-spec-sweep note reconciled against current
  `TODO-kugod` ownership.
- `Turn 161` `UT-161.a` is answered by the TE-havib follow-on disposition.
  Carry items `UT-161.b` and `UT-161.c` remain conversation-only lineage /
  taxonomy notes and still need explicit downstream placement if reused.
- `Turn 162` `UT-162.a` and `UT-162.b` are answered by the later TE-havib
  disposition path. No additional turn-local owner update has landed yet
  because `TODO-lilok` still carries the reopened sweep-handoff note.
- `Turn 163` `UT-163.a` remains a spec-edit loose end for the envelope-agnostic
  `§1.3` wording template. `UT-163.b` remains a carry/procedural note about
  uncaptured TE-havib twig commits.
- `Turn 164` `UT-164.c`, `UT-164.d`, and `UT-164.e` remain open under
  `TODO-turog`. Carry items `UT-164.a` and `UT-164.b` remain procedural /
  historical corrections only.
- `Turn 165` `UT-165.b` remains open under `TODO-turog`. `UT-165.d` and
  `UT-165.e` remain spec-edit loose ends. Carry items `UT-165.a` and
  `UT-165.c` remain process/privacy notes rather than live design work.
- `Turn 166` `UT-166.a`, `UT-166.b`, `UT-166.c`, and `UT-166.e` remain open
  under `TODO-turog`. Carry item `UT-166.d` remains historical metadata only.
- `Turn 167` `UT-167.a`, `UT-167.b`, `UT-167.c`, and `UT-167.e` remain open
  under `TODO-turog`. `UT-167.d` remains an open spec-edit item for the
  missing Cat-2 quotation grep.
- `Turn 168` `UT-168.b`, `UT-168.c`, and `UT-168.d` remain open under
  `TODO-turog`. `UT-168.a`, `UT-168.e`, and `UT-168.f` remain open under
  `TODO-duvuk`.
- `Turn 169` `UT-169.a`, `UT-169.b`, `UT-169.d`, and `UT-169.e` remain open
  under `TODO-duvuk`. Carry item `UT-169.c` remains a procedural cadence note.
- `Turn 170` `UT-170.b`, `UT-170.c`, and `UT-170.d` are retired. `UT-170.a`
  remains live in the TE-sihih / `TODO-vunub` cluster and must close there, not
  in replay history.
- `Turn 171` `UT-171.a` and `UT-171.b` remain live in TE-sihih /
  `TODO-vunub`. Carry items `UT-171.c` and `UT-171.d` remain historical
  cadence / design-process notes.
- `Turn 172` `UT-172.a`, `UT-172.b`, `UT-172.c`, `UT-172.d`, and `UT-172.e`
  remain live in TE-sihih / `TODO-vunub`.
- `Turn 173` `UT-173.a` has partial downstream doc promotion via `DI-pijun`,
  but it is not fully closed as a protocol-identity claim. `UT-173.b`,
  `UT-173.c`, `UT-173.d`, and `UT-173.e` remain live in TE-sihih /
  `TODO-vunub`.
- `Turn 174` `UT-174.a`, `UT-174.b`, `UT-174.c`, `UT-174.d`, and `UT-174.e`
  remain live in TE-sihih / `TODO-vunub`. The historical-lineage correction is
  captured, but the vocabulary/layout fallout still needs substantive owner
  work.

## Subtasks

- [x] juhub.149 Turn 149 raw-log rewalk plus later-turn and later-artifact sweep. **Done 2026-05-10.** Raw turn confirmed the Cat-1a/Cat-1b split recommendation; existing capture was correct; no later contradiction found. See `## Turn notes`.
- [x] juhub.150 Turn 150 raw-log rewalk plus later-turn and later-artifact sweep. **Done 2026-05-10.** Raw turn confirmed the cross-TE quotation-grep safeguard for future Cat-2 vocabulary sweeps; existing capture was correct; no later contradiction found. See `## Turn notes`.
- [x] juhub.151 Turn 151 raw-log rewalk plus later-turn and later-artifact sweep. **Done 2026-05-10.** Raw turn confirmed the Cat-2 cross-TE quotation-grep safeguard; existing capture was correct; no later contradiction found. See `## Turn notes`.
- [x] juhub.152 Turn 152 raw-log rewalk plus later-turn and later-artifact sweep. **Done 2026-05-10.** Raw turn confirmed the top-of-file `Status:` field rule and the immediate unblocking of the follow-on TE-policy sweeps; existing capture was correct; no later contradiction found. See `## Turn notes`.
- [x] juhub.153 Turn 153 raw-log rewalk plus later-turn and later-artifact sweep. **Done 2026-05-10.** Raw turn authorized execution of the already-unblocked TODO-020 rollout work; existing capture was correct; no later contradiction found. See `## Turn notes`.
- [x] juhub.154 Turn 154 raw-log rewalk plus later-turn and later-artifact sweep. **Done 2026-05-10.** Raw turn was a queue-status and recommendation checkpoint, not a new design turn; existing capture was correct; no later contradiction found. See `## Turn notes`.
- [x] juhub.155 Turn 155 raw-log rewalk plus later-turn and later-artifact sweep. **Done 2026-05-10.** Raw turn clarified the stalled DF-1.1 proposal but did not lock it; existing capture was correct; no later contradiction found. See `## Turn notes`.
- [x] juhub.156 Turn 156 raw-log rewalk plus later-turn and later-artifact sweep. **Done 2026-05-10.** Raw turn interrupted the old TE-famar DF path with a scope correction; existing capture was correct; no later contradiction found. See `## Turn notes`.
- [x] juhub.157 Turn 157 raw-log rewalk plus later-turn and later-artifact sweep. **Done 2026-05-10.** Raw turn corrected the transport-vs-envelope confusion and proposed a higher-level envelope-shape TE; existing capture was correct; no later contradiction found. See `## Turn notes`.
- [x] juhub.158 Turn 158 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.159 Turn 159 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.160 Turn 160 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.161 Turn 161 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.162 Turn 162 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.163 Turn 163 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.164 Turn 164 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.165 Turn 165 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.166 Turn 166 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.167 Turn 167 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.168 Turn 168 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.169 Turn 169 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.170 Turn 170 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.171 Turn 171 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.172 Turn 172 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.173 Turn 173 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.174 Turn 174 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.175 Turn 175 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.176 Turn 176 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.177 Turn 177 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.178 Turn 178 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.179 Turn 179 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.180 Turn 180 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.181 Turn 181 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.182 Turn 182 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.183 Turn 183 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.184 Turn 184 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.185 Turn 185 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.186 Turn 186 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.187 Turn 187 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.188 Turn 188 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.189 Turn 189 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.190 Turn 190 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.191 Turn 191 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.192 Turn 192 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.193 Turn 193 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.194 Turn 194 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.195 Turn 195 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.196 Turn 196 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.197 Turn 197 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.198 Turn 198 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.199 Turn 199 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.200 Turn 200 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.201 Turn 201 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.202 Turn 202 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.203 Turn 203 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.204 Turn 204 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.205 Turn 205 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.206 Turn 206 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.207 Turn 207 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.208 Turn 208 raw-log rewalk plus later-turn and later-artifact sweep.

## Turn notes

### Turn 149 — 2026-05-03 00:05 UTC

- `Turn 149 summary` Steve's `yes` was a real decision, not a filler
  acknowledgement. The bot had just proposed splitting Cat-1 path renames into
  two cases: current live path references that can be updated mechanically, and
  historical quotations of old paths that must be preserved. Steve's approval
  locked that split.
- `Existing capture` `TODO-lilar` already records the turn correctly: the
  approval led to `DI-020-20260502-232651`, committed as `cd82c19`, merged as
  `d8c3e93`, pushed, then the conversation moved on to DF-35.2. Later policy
  artifacts (`TODO-dinub`, TE-dabol's Refinements, and the TE-editing-policy
  commits) still treat that Cat-1a/Cat-1b split as the live locked rule.
- `Gaps or contradictions` None found. Later turns and later artifacts do not
  retract, narrow, or correct the turn-149 decision. Turn 197's replay
  instructions even classify the early 149-153 block as straightforward
  confirmations of already-landed TODO-020 work.
- `Related UTs / owners` None. The replay/disposition chain begins its `UT-*`
  inventory at turn 155, so turn 149 has no downstream owner TODO to reconcile.
- `Owner/doc cleanup` None needed. The TE-editing-policy artifacts already carry
  this decision, and no later spec, research, simulation, or owner TODO is
  missing a load-bearing statement from this turn.
- `Remaining decisions or work` None.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 150 is next and remains
  pending approval.

### Turn 150 — 2026-05-03 00:23 UTC

- `Turn 150 summary` Steve's `yes` confirmed the second tightening of the TE
  editing policy. The issue here was not about paths anymore; it was about
  quoted words. If one TE quotes another TE's exact wording, and a later
  vocabulary sweep rewrites the quoted term everywhere, then the quote stops
  being historically true. This turn approved the rule that any future Cat-2
  vocabulary sweep must first grep the corpus for quoted uses of the old term
  and leave genuine historical quotations alone.
- `Existing capture` The old replay ledger already has the substantive outcome
  right even though it grouped the approval under the next line item's naming.
  The rule later landed as the Cat-3 Refinement recorded in TE-dabol and TODO
  020, with commit `04126ac` merged as `795a846`. Later artifacts explain the
  effect clearly: Cat-2 sweeps now have two required checks, not one. First,
  the sweeper must name the DIs whose meaning is unchanged. Second, the
  sweeper must grep for old-term-in-quotation contexts before rewriting.
- `Gaps or contradictions` No contradiction found. Later policy docs, TODO 020,
  TE-dabol's Refinements, and the replay notes all continue to describe this as
  a Cat-3 procedural tightening rather than a new superseding DI. The only
  nuance is that `TODO-lilar` phrases the turn boundary in a slightly shifted
  way, because the old walk grouped the DF confirmation block as a tight series
  of approvals. Substantively, the landing is still correct.
- `Related UTs / owners` None. The replay/disposition chain does not file any
  `UT-*` entries for turns 149-154, so there is no downstream owner TODO to
  reconcile for this turn.
- `Owner/doc cleanup` None needed. The TE-editing-policy artifacts already
  carry the rule family this turn confirmed, and no later owner or direct
  design/spec/research/simulation doc appears to be missing a load-bearing
  statement from this turn.
- `Remaining decisions or work` None.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 151 is next and remains
  pending approval.

### Turn 151 — 2026-05-03 00:28 UTC

- `Turn 151 summary` Steve's `yes` approved the next TE-editing safeguard:
  before any Cat-2 vocabulary sweep rewrites a term across the corpus, the bot
  must grep for uses of the old term inside quotation-like contexts and emit
  those matches for human review. In plain English, this is the rule that stops
  a vocabulary cleanup from silently rewriting quoted history and making an
  older TE appear to have said something it never said.
- `Existing capture` `TODO-lilar` already records the turn correctly as the
  confirmation of DF-35.3, the mandatory cross-TE quotation-grep step before a
  Cat-2 sweep. Later artifacts confirm the same result: the rule landed as a
  Cat-3 Refinement on TE-dabol, with the refinement text explaining that a
  Cat-2 sweeper now has a two-step protocol — enumerate unchanged DIs in the
  top-of-file note, then grep and classify quotation-context matches before
  rewriting. Later replay summaries keep the same interpretation.
- `Gaps or contradictions` None found. I did not find any later turn or later
  artifact that narrows, retracts, or contradicts the turn-151 decision. The
  later TE-35 summary still lists DF-35.3 as settled and merged, with the
  quotation-grep safeguard intact.
- `Related UTs / owners` None. The replay/disposition chain does not file any
  `UT-*` entries for turns 149-154, so there is no downstream owner TODO to
  reconcile for this turn.
- `Owner/doc cleanup` None needed. The TE-editing-policy corpus already carries
  this safeguard, and no later owner or direct design/spec/research/simulation
  doc appears to be missing a load-bearing statement from this turn.
- `Remaining decisions or work` None.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 152 is next and remains
  pending approval.

### Turn 152 — 2026-05-03 00:35 UTC

- `Turn 152 summary` Steve's `yes` closed the last open DF from TE-vudaf. The
  substance was the new top-of-file `Status:` field on every TE, placed where a
  reader sees it before reading the body. The problem being solved was stale
  supersedence discoverability: if the only "this TE is superseded" marker lives
  at the very bottom, a reader can miss it and act on stale reasoning. This
  turn approved the fix and immediately reframed the next work as execution
  tasks, not more policy debate.
- `Existing capture` `TODO-lilar` already records the turn correctly as the
  confirmation of DF-35.4 Alt-4.a: a uniform top-of-file `Status:` header field
  on every TE. The raw turn itself is mostly a summary table, but it clearly
  states the outcome: DF-35.4 landed as a Cat-3 Refinement on TE-dabol, the
  retrofit is subtask 020.10, and subtasks 020.5 / 020.6 / 020.7 / 020.10 moved
  from deferred to ready-to-execute. Later artifacts preserve exactly that
  framing. TE-dabol's Refinements define the field shape and purpose, and
  `TODO-dinub` records that 020.10 later added the field to all 35 existing TEs.
- `Gaps or contradictions` None found. The later corpus consistently treats
  this turn as the point where the TE-editing policy became fully settled: four
  DIs plus four Cat-3 Refinements, with the remaining work reduced to rollout
  and mechanical sweeps. No later artifact disputes the turn boundary or the
  meaning of the `Status:` field decision.
- `Related UTs / owners` None. The replay/disposition chain does not file any
  `UT-*` entries for turns 149-154, so there is no downstream owner TODO to
  reconcile for this turn.
- `Owner/doc cleanup` None needed. The TE-editing-policy corpus already carries
  this decision, and no later owner or direct design/spec/research/simulation
  doc appears to be missing a load-bearing statement from this turn.
- `Remaining decisions or work` None.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 153 is next and remains
  pending approval.

### Turn 153 — 2026-05-03 00:44 UTC

- `Turn 153 summary` Steve's `do it` was not another policy choice. It was
  authorization to execute the rollout work that turn 152 had just unblocked.
  In plain English, the TE-editing policy debate was over, and this turn is the
  bot reporting what it actually landed from that now-settled queue.
- `Existing capture` `TODO-lilar` already captures the important nuance
  correctly. The raw turn itself explicitly described two landed items: 020.7
  (the TE-famar `## Refinements` forward pointer) and 020.5 (the AGENTS
  rollout). But the same raw turn also said all four previously unblocked
  subtasks were now done, and the later TODO-020 state confirms the full batch
  included 020.6 (the Cat-1a/Cat-1b path-reference sweep) and 020.10 (the
  top-of-file `Status:` retrofit across all 35 TEs). So the raw turn's prose
  named two items in detail, while the later artifacts confirm the larger
  completion set.
- `Gaps or contradictions` None found. Later artifacts preserve the same
  interpretation: after this turn, only 020.8 remained open in TODO-020. I did
  not find any later correction saying the turn-153 batch was narrower than
  `TODO-lilar` records.
- `Related UTs / owners` None. The replay/disposition chain does not file any
  `UT-*` entries for turns 149-154, so there is no downstream owner TODO to
  reconcile for this turn.
- `Owner/doc cleanup` None needed. TODO-020 already carries the landing state
  of the rollout batch, and no later owner or direct
  design/spec/research/simulation doc appears to be missing a load-bearing
  statement from this turn.
- `Remaining decisions or work` None for turn 153 itself.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 154 is next and remains
  pending approval.

### Turn 154 — 2026-05-03 01:17 UTC

- `Turn 154 summary` Steve's `next?` did not open a new design thread by
  itself. The bot replied with a queue-status snapshot: here are the open
  TODOs, here is the recommended next move, and here is the suggested order.
  In plain English, this was a planning checkpoint between the now-finished
  TE-editing-policy rollout and the next substantive decision work.
- `Existing capture` `TODO-lilar` already records the turn correctly as a pure
  queue / recommendation boundary. The raw turn lists 020.8 as the small
  cleanup item, TODO-rivuk as the high-leverage next DF work, and several other
  open items ranging from DI/DR backfill to implementation scaffolding. The
  later replay notes correctly preserve two important nuances: this turn is not
  the start of TE-havib, and TODO-bisur was still visibly alive in the queue
  here with two open subtasks.
- `Gaps or contradictions` None found. Later notes explicitly correct the turn
  boundary in the same way: turn 154 is the queue-status turn, while the
  TE-havib / apparatus-vs-specimen sequence begins only after Steve's next
  response and the scope challenge that follows. I found no later artifact that
  reclassifies turn 154 as a substantive design decision.
- `Related UTs / owners` None. The replay/disposition chain does not file any
  `UT-*` entries for turns 149-154, so there is no downstream owner TODO to
  reconcile for this turn.
- `Owner/doc cleanup` None needed. This turn is a queue snapshot rather than a
  turn where a missing design/spec/research/simulation statement needs to be
  promoted, and there is no turn-local owner cleanup beyond the already-listed
  queue artifacts themselves.
- `Remaining decisions or work` None for turn 154 itself. The queue items named
  here remained open, but that is not a replay loose end created by this turn.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 155 is next and remains
  pending approval.

### Turn 155 — 2026-05-03 01:38 UTC

- `Turn 155 summary` Steve picked the TE-famar path from the queue, but he did
  not lock DF-1.1. Instead he pushed back on two missing pieces: what exactly
  `Project` means, and where any resulting decision would actually live. In
  plain English, this turn is the bot trying to rescue the old promise-stack
  DF flow by explaining the terms more concretely before asking for a lock.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  answer defines `Project(msg, predicate)` as a non-consuming, order-independent
  query over frames, contrasts it with `Peel(msg)`, and says the old intended
  lifecycle was: lock DI-005-1.1 in TODO-rivuk, later cite the resulting DIs
  from `harness-spec-draft.md §1.1`, then close DR-006. Later artifacts also
  preserve the crucial outcome: none of this actually locked, because turn 156
  immediately reframed the scope before Steve answered yes.
- `Gaps or contradictions` None found. The later corpus consistently treats two
  things as the important leftovers from this turn: first, DF-1.1 was never
  locked and the old TODO-rivuk queue was later superseded; second, the clearest
  `Project / Peel / Wrap` definitions in the corpus still live only in this
  conversation. The wording "projection mode" is also recognized later as a bad
  phrase that should not propagate into committed text; the real issue was
  whether `Project` is part of the spec contract.
- `Related UTs / owners` `UT-155.a` and `UT-155.b` are the turn-local loose
  ends. Both are now retired in `TODO-kugod`: DF-1.1 is no longer live
  promise-stack work, and the `Project / Peel / Wrap` vocabulary remains
  historical rather than active apparatus work.
- `Owner/doc cleanup` None needed now. The owner chain already retired the old
  TE-famar line correctly, and there is no direct-doc promotion needed from
  this turn unless that vocabulary is deliberately revived under a new active
  owner later.
- `Remaining decisions or work` None for turn 155 itself. The conversation text
  remains useful historical evidence, but there is no still-open replay owner
  action attached to this turn after the later retirement.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 156 is next and remains
  pending approval.

### Turn 156 — 2026-05-03 01:39 UTC

- `Turn 156 summary` Steve cut off the old promise-stack DF flow before any
  lock happened. His point was that the topic under discussion was not a
  harness-level invariant; it was only one candidate wire-envelope design. In
  plain English, this turn is where the bot first admits the old TE-famar /
  TODO-rivuk framing is scoped wrong.
- `Existing capture` `TODO-lilar` already records the turn correctly as a
  mid-DF scope correction. The raw turn shows the bot responding by offering
  three structural choices for what to do with TE-famar's analysis: make it a
  harness-level default, demote it to a per-envelope concern, or split it into
  harness-level vocabulary plus per-envelope lockings. The later corpus also
  preserves the key outcome: none of those three options was actually locked,
  because turn 157 reframed the question again before Steve chose among them.
- `Gaps or contradictions` None found. Later artifacts preserve two critical
  follow-on facts consistently: first, the Option 1 / 2 / 3 menu was abandoned
  rather than answered; second, the bot's wording that the harness-spec should
  be "wire-envelope-agnostic" is itself later corrected as wrong. The later
  apparatus-vs-specimen framing says the harness is not envelope-agnostic; it is
  the apparatus that compares candidate envelopes and other layer choices.
- `Related UTs / owners` `UT-156.a` is retired and `UT-156.b` is
  resolved-retired under `TODO-rivuk` / `DI-runuh`; `UT-156.c` is now resolved
  in `TODO-kugod` under `DI-lajod`, which rewrites `harness-spec-draft.md §1.1`
  and `§1.3` at apparatus level and explicitly retracts the stale
  "wire-envelope-agnostic" wording.
- `Owner/doc cleanup` Done. `TODO-kugod` now records `UT-156.c` as resolved, and
  `harness-spec-draft.md` now states that the harness compares candidate
  envelopes rather than defining one canonical envelope. No correction note is
  needed in `TODO-lilar`; the former turn-158 slice is now closed, and the
  remaining broader sweep lives under the later `UT-159.*` rows in
  `TODO-kugod`.
- `Remaining decisions or work` None for turn `156` itself. The broader
  apparatus/specimen cleanup is still open, but it belongs to later turns and
  their owner items rather than to the turn-156 stale-claim residue.
- `Work pending` `no`
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 157 is next and remains
  pending approval.

### Turn 157 — 2026-05-03 01:42 UTC

- `Turn 157 summary` Steve corrected the bot again: transports and envelopes
  are different things, the goal is one transport-agnostic message envelope,
  and `grid([pcid, payload])` is only the current working hypothesis, not a
  proven answer. In plain English, this turn forces the bot to step back from
  the old promise-stack-only framing and admit that the envelope decision
  itself is still under study.
- `Existing capture` `TODO-lilar` already records the turn correctly. The raw
  turn names five candidate envelope shapes, explains why TE-famar jumped ahead
  by assuming the promise-stack family was already the right abstraction, and
  recommends "Reading 2": file a higher-level envelope-shape TE, treat TE-famar
  as misframed input to that larger decision, and gate the old TODO-rivuk DF
  queue behind the new envelope-shape work. Later owner cleanup also preserves
  the two load-bearing carry-forwards from this turn: the candidate-envelope
  inventory and the `grid([pcid, payload])` working-hypothesis prose move to the
  future grid-envelope successor work rather than staying attached to
  promise-stack.
- `Gaps or contradictions` None found. Later artifacts preserve this as a
  transitional correction, not as a final settled architecture. The five
  candidate envelopes and the Reading-2 recommendation are kept as historical
  evidence, but turn 158 immediately rejects the remaining assumption that the
  envelope belongs in the harness-spec as a single harness-wide decision. I did
  not find any later artifact claiming the turn-157 framing was itself the final
  answer.
- `Related UTs / owners` `UT-157.a` and `UT-157.c` are resolved in
  `TODO-kugod` by `DI-joroh`, with the candidate-envelope inventory and
  `grid([pcid, payload])` working-hypothesis prose captured in
  `simulations/SIM-kurim-grid-envelope/TODO/TODO-tujad-grid-envelope-successor-owner.md`.
  `UT-157.b` is retired under `TODO-rivuk` / `DI-runuh`.
- `Owner/doc cleanup` Done. The abandoned Reading-1/2/3 question is retired, the
  five candidate envelopes are captured in the grid-envelope successor owner,
  and the "working but not yet proven" grid framing is recorded as a candidate
  hypothesis rather than a canonical harness rule. No correction note is needed
  in `TODO-lilar`.
- `Remaining decisions or work` None for turn 157 itself. The later
  grid-envelope protocol directory/spec work that had remained under turn 158 /
  `UT-158.f` and `tujad.3` is now closed for turn-158 scope by `DI-fanah`.
- `Work pending` `no`
- `Proposed disposition` `reconciled after successor-owner capture`
- `Write needed? yes/no` `yes` for this rewalk update in `TODO-juhub`; `no`
  correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 158 is next and remains
  pending approval.

### Turn 158 — 2026-05-03 01:46 UTC

- `Turn 158 summary` This is the real apparatus-vs-specimen break. Steve
  pointed out that even calling the envelope "harness-wide" smuggles in the
  wrong assumption, because wire-lab exists to test multiple hypotheses at all
  layers, not to bake one answer into the harness. In plain English, this is
  the turn where the bot finally accepts that the harness-spec is lab
  apparatus and candidate envelopes are specimens under study.
- `Existing capture` `TODO-lilar` already captures the turn correctly as the
  foundational apparatus-vs-specimen reframe. The raw turn lays out a six-step
  sequence: audit the harness-spec, file the harness-level TE on the split,
  give each candidate envelope its own protocol home, sweep specimen material
  out of the harness-spec, reframe the old promise-stack TODO under protocol
  ownership, and file a parallel TODO for the grid-envelope hypothesis.
- `Gaps or contradictions` The insight itself stands; no later artifact reverts
  to the old claim that the harness-spec should define a single envelope
  specimen. The gap was bookkeeping: the six-step sequence and the parallel
  grid-hypothesis TODO were still represented as loose carry items even though
  later artifacts had decomposed, routed, or materialized them as successor
  simulations.
- `Related UTs / owners` `UT-158.a` is resolved-decomposed in `TODO-kugod`
  under `DI-sotuk`; `UT-158.b` is resolved by TE-havib DF-36.1; `UT-158.e`
  and `UT-158.g` are retired under `TODO-rivuk` / `DI-runuh`; `UT-158.h` is
  resolved-routed to `simulations/SIM-kurim-grid-envelope/TODO/TODO-tujad-grid-envelope-successor-owner.md`
  under `DI-sotuk`. `UT-158.c` and `UT-158.d` are resolved for turn-158 scope
  under `DI-kinad`; `UT-158.f` is resolved-transferred under `DI-fanah`.
- `Owner/doc cleanup` Done. `TODO-kugod` now has explicit `UT-158.a`,
  `UT-158.c`, `UT-158.d`, `UT-158.f`, and `UT-158.h` disposition rows, and
  `TODO-tujad` now closes `tujad.3` by pointing to the 24 standalone
  positional grid-envelope successor simulations.
  No correction note is needed in `TODO-lilar`.
- `Remaining decisions or work` None for turn 158 itself. Broader TE-40 audit
  work remains under later turn-159 rows (`UT-159.a`, `UT-159.b`, `UT-159.d`,
  `kugod.5`, and `kugod.8`), not as turn-158-local residue.
- `Work pending` `no`
- `Proposed disposition` `reconciled after positional variant split`
- `Write needed? yes/no` `yes` for this rewalk update in `TODO-juhub` and the
  owner-routing updates in `TODO-kugod` / `TODO-tujad`; `no` correction note is
  needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 159 is next and remains
  pending approval.

### Turn 159 — 2026-05-03 01:50 UTC

- `Turn 159 summary` Steve confirmed the apparatus-vs-specimen reframe and
  told the bot to proceed with step 1 of the six-step plan. In plain English,
  this is the turn where the bot turns the reframe into a concrete inventory:
  it audits the harness-spec section by section and sorts what stays in the
  harness, what must move out as specimen-specific material, and what remains
  ambiguous.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The
  raw turn shows the bot produced an audit memo as an untracked working
  document, not yet a commit, and that memo listed the concrete specimen-
  bearing material in the harness-spec plus the ambiguous areas needing later
  TE or DF work. Later residual owner work in `TODO-kugod` directly carries
  this turn forward: the nine specimen-bearing audit items remain live as
  `UT-159.a`, the transport-spec companion audit remains live as `UT-159.b`,
  and the unresolved ambiguous areas remain live as `UT-159.d`.
- `Gaps or contradictions` None found. The only important timing nuance is
  that the audit memo existed in the workspace at the end of turn 159 but was
  not committed until Steve's turn-160 authorization. `TODO-lilar` already
  preserves that nuance, and later artifacts align with it.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 160 is next and remains
  pending approval.

### Turn 160 — 2026-05-03 02:24 UTC

- `Turn 160 summary` Steve's `do it and continue` authorized two separate
  actions: first, commit the audit memo from turn 159, and second, continue
  into step 2 by drafting the apparatus-vs-specimen TE. In plain English, this
  is the turn where the work stops being just an audit and becomes a formal DF
  program on a dedicated TE twig.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  turn says the audit memo was committed at `4725b3e`, the new TE-havib work
  was drafted on twig `ppx/te-20260503-022446-apparatus-vs-specimen`, and the
  TE already contained seven drafted DFs plus six tabletop scenarios. Just as
  importantly, the answer text presented only `DF-36.1` to Steve in this turn,
  following the standing one-DF-at-a-time rule; the other six drafted DFs
  existed in the TE file but were not yet exposed in conversation.
- `Gaps or contradictions` None that overturn the existing capture. The later
  carry-forward items already preserve the important weaknesses introduced here:
  the audit count in the answer text said eight specimen-bearing items while
  the audit itself listed nine; the PT vocabulary collapse was baked into
  `DF-36.4` rather than framed as its own decision; the TE's six scenarios only
  partially align with the audit's recommended scenario set; and most of the
  seven drafted DFs remained unlocked at end-of-corpus.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 161 is next and remains
  pending approval.

### Turn 161 — 2026-05-03 02:57 UTC

- `Turn 161 summary` Steve asked the key redundancy question: what is the
  actual difference between `promise-stack` and `grid-pcid-payload`? In plain
  English, this is where the bot is forced to explain whether these are truly
  two peer envelope hypotheses or whether one is really just a special case of
  the other.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  turn contains the most detailed side-by-side comparison in this whole replay
  slice: `promise-stack` is presented as a recursive CBOR array of promise
  frames whose trust and layering semantics live in the envelope itself, while
  `grid-pcid-payload` is presented as a thin dispatch envelope where the pCID
  is the top-level selector and the payload shape is left to the protocol it
  names. The most important point is the asymmetry the bot exposed: a
  `grid-pcid-payload` message can carry a promise-stack inside its payload, but
  a promise-stack message does not cleanly host `grid-pcid-payload` as a peer
  outer envelope. That asymmetry is the conceptual seed of `OQ-36.6`.
- `Gaps or contradictions` None that overturn the existing capture. Later
  artifacts already preserve the two main carry-forwards from this turn: the
  asymmetry may not have survived cleanly into the later committed OQ text, and
  the nine-axis comparison table plus the richer open-set assertion taxonomy
  still live only in conversation rather than in a committed doc.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 162 is next and remains
  pending approval.

### Turn 162 — 2026-05-03 03:12 UTC

- `Turn 162 summary` Steve sharpened the redundancy concern from turn 161 into
  an explicit suspicion: promise-stack may be overcomplicated machinery
  invented from a misunderstanding of how nested messages already work inside
  `grid-pcid-payload`. He then gave two procedural instructions: note that
  concern for later, and keep going on `DF-36.5`. In plain English, this is
  the turn where the redundancy issue stops being just an implication and
  becomes an explicitly parked open question.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  turn's visible reply is minimal, but the substantive effect is preserved in
  later artifacts: the concern is recorded on the TE-havib twig as `OQ-36.6`,
  and from this point on the promise-stack-home path becomes provisional rather
  than cleanly settled. If the later investigation concludes that promise-stack
  is really just one payload-shape under `grid-pcid-payload`, then the earlier
  promise-stack-home direction either retires or stays only as a minimal
  placeholder pending that decision.
- `Gaps or contradictions` None that overturn the existing capture. The later
  residual notes already preserve the two important nuances from this turn:
  the committed `OQ-36.6` text may not preserve the cleaner asymmetry from turn
  161 as clearly as the conversation did, and the provisional nature of the
  earlier promise-stack-home direction may not be visually obvious to a later
  reader of the TE twig.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 163 is next and remains
  pending approval.

### Turn 163 — 2026-05-03 03:34 UTC

- `Turn 163 summary` Steve rejected the prior presentation of `DF-36.5` as
  unreadable and told the bot to format it better. In plain English, this is
  not just a meta-turn about presentation; it is the first readable,
  substantive walk of `DF-36.5`.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  turn re-presents the real decision about `§1.3` of `harness-spec`: whether
  the four simulator tests written in promise-stack vocabulary should stay in
  the harness envelope-agnostically, move wholesale to a specimen spec, or be
  split across the two abstraction levels. The bot recommends `Alt-5.C`: keep
  an apparatus-level summary in the harness and move specimen-specific details
  out. This is also the first DF in TE-havib where the parked `OQ-36.6`
  uncertainty is built explicitly into the recommendation: the apparatus-level
  summary survives either way, but the specimen-side destination depends on
  whether promise-stack later survives as a distinct specimen.
- `Gaps or contradictions` None that overturn the existing capture. The main
  carry-forward from this turn is that the bot effectively created an
  envelope-agnostic rewrite template for the four `§1.3` tests and implied
  that the same pattern should apply to other ambiguous sections such as `§3.3`
  and `§7.1`. That template lives only in conversation, not in a committed
  artifact.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 164 is next and remains
  pending approval.

### Turn 164 — 2026-05-03 11:20 UTC

- `Turn 164 summary` This is the hard pivot away from the TE-havib DF walk and
  toward an urgent operational problem: Steve needs file-based transport
  working so he and another human collaborator can use the repo itself to
  collaborate. In plain English, the raw turn does not implement the bootstrap
  yet; it marks the urgency boundary and the point where the bot pauses to
  confirm scope before switching threads.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  answer is a pause, not a solution: the bot says it should not assume who the
  collaborator is and explicitly notes that pivoting away from TE-havib is a
  consequential choice. Later artifacts correctly treat this as the major
  thread boundary where TE-havib becomes parked and the transport-bootstrap
  work begins immediately afterward.
- `Gaps or contradictions` None that overturn the existing capture. The key
  later carry-forwards are already preserved elsewhere: the bot overclaimed
  TE-havib as "4-DF-locked" when only two DFs were actually locked, it drafted
  the unnecessary TE-numan in the gap immediately after this turn, and the
  remaining unfinished TE-havib threads had to be carried forward before the
  transport work took over.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 165 is next and remains
  pending approval.

### Turn 165 — 2026-05-03 11:21 UTC

- `Turn 165 summary` Steve disclosed that the collaborator is a real human
  using Claude in another clone and imposed a privacy constraint: the
  collaborator's name must not appear in committed docs. In plain English, this
  is the turn where the transport-bootstrap problem becomes a real coordination
  problem with an explicit anonymity rule attached to it.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  answer stacks four distinct moves: it records the collaborator context as a
  side effect, corrects the bot's architectural misread by admitting that a new
  git-file-transport TE was unnecessary because the existing transport and
  group-session drafts already cover the use case, claims the mistaken TE twig
  was "reverted cleanly," and then lays out the concrete bootstrap
  chicken-and-egg with three options and a recommendation for the draft-
  directory approach (`Option A`).
- `Gaps or contradictions` None that overturn the existing capture. The later
  residual notes already preserve the important problems introduced here: the
  "reverted cleanly" claim conflicts with the later reflog evidence, the slug
  examples in the answer were still not fully aligned with the anonymity rule,
  the `merge-group-transport-spec` step was treated as operationally defined
  when it was not, and the turn may have triggered a memory write that inferred
  facts not actually present in Steve's prompt.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 166 is next and remains
  pending approval.

### Turn 166 — 2026-05-03 11:29 UTC

- `Turn 166 summary` Steve corrected two things at once: the transport group is
  at least three developer agents, not a two-party collaboration, and the slug
  must be generic and identity-free rather than derived from people-names. In
  plain English, this is the turn where the bootstrap stops being a hypothetical
  plan and becomes the executed `wire-lab-devs` instance under urgency.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  answer updates the group model, proposes generic slug candidates, chooses
  `wire-lab-devs`, raises the membership-pinning issue from `§8`, and then
  reports the bootstrap as already live and pushed. Later timeline evidence
  confirms that the bootstrap commit and merge landed as part of this same
  response, not before it, so this is a high-effects execution turn rather than
  just a naming correction.
- `Gaps or contradictions` None that overturn the existing capture. The later
  residual notes already preserve the important problems introduced here: the
  bot speculated about one participant who Steve had not explicitly enrolled,
  raised but did not resolve the `§8` membership-pinning requirement before
  executing, reported through a stale branch name, and offered two reply-paths
  while immediately executing one of them.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 167 is next and remains
  pending approval.

### Turn 167 — 2026-05-03 11:37 UTC

- `Turn 167 summary` Steve gave two directives at once: switch message files
  from `.msg` to `.txt`, and require members to fetch all branches but post
  only on their own `<author-id>/main`. In plain English, this is both a
  presentational cleanup turn and the first committed transport-binding turn
  that maps the abstract group-session protocol onto an actual Git branch
  discipline.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  answer shows that the bot treated both directives as immediately executable:
  it rewrote the docs, renamed the bootstrap message file without changing its
  bytes, added new non-normative `§9` to describe the per-author-branch
  binding, and authored `m001` as the first on-wire ratification of that
  branch-binding rule. This is the turn where the transport instance becomes
  "in flight" with an explicit Git discipline: read from all known branches,
  write only to your own branch, and let `Parents:` rather than branch topology
  carry ordering.
- `Gaps or contradictions` None that overturn the existing capture. The later
  residual notes already preserve the important problems introduced here: the
  sequential `m<N>` naming convention is locked here and reversed one turn
  later, `§9` quietly leans toward extensible-by-posting despite the stronger
  `§8` language, the meaning of `{name}` was not fully nailed down in Steve's
  prompt but the bot chose an interpretation in the spec, the six-document
  sweep appears to have happened without the quotation-aware Cat-2 grep
  discipline, and this is another execute-on-directive turn rather than a
  recommend-and-wait turn.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 168 is next and remains
  pending approval.

### Turn 168 — 2026-05-03 11:45 UTC

- `Turn 168 summary` Steve corrected the new branch-binding model from turn
  167 in two important ways: message files must not use sequential numbers, and
  before posting, each agent must first merge all observed messages from all
  branches into the directory on their own branch and push that merged state.
  In plain English, this turns each participant's branch from "my outbound
  posts" into "my current replicated view of the message set, plus optionally a
  new post," and it makes the message CID, not a human numbering scheme, the
  stable on-disk identity.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  answer shows the bot treated both directives as immediately executable: it
  rewrote the spec so filenames are CIDs, expanded `§9` into a fuller
  receive/merge/push/optionally-post cycle, renamed the existing bootstrap
  files to their CIDs, and authored a new on-wire ratification message for the
  new rules. This is the turn that fixes the global-sequence mistake from turn
  167 and defines the first actual replication model for the Git-backed
  channel.
- `Gaps or contradictions` None that overturn the existing capture. The later
  residual notes already preserve the important follow-ons from this turn: the
  `Message-ID` header is still retained here and then challenged immediately in
  turn 169, `§9` grows substantially while still labeled non-normative, and the
  CIDs reported in this turn are historically correct for this moment but do
  not remain the current on-disk CIDs after turn 169's rehash.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 169 is next and remains
  pending approval.

### Turn 169 — 2026-05-03 11:54 UTC

- `Turn 169 summary` Steve asked whether `Message-ID` is still needed now that
  the canonical identifier is the message CID. In plain English, the raw turn
  is a careful reasoning memo, not an execution turn: the bot walks through
  every plausible use of `Message-ID`, concludes that it mostly creates
  confusion once filename = CID, and then reasons about how to remove it
  without breaking the existing bootstrap messages.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  reasoning lands on the conservative answer: keep legacy messages parseable,
  deprecate `Message-ID` for new messages, and ask Steve to choose between soft
  deprecation and hard removal. The key compatibility issue the bot itself
  identifies is `§4.7`: if strict readers reject unknown headers, then old
  bootstrap messages containing `Message-ID` become a problem under a spec that
  simply forbids the header. The turn matters because it recommends `Path A`
  and explicitly says it should ask before acting.
- `Gaps or contradictions` None that overturn the existing capture. The main
  carry-forward is the reasoning-versus-action split already preserved in
  `TODO-lilar`: the raw turn recommends deprecate-and-ask, while the later repo
  history performs hard removal plus rehash. That later divergence belongs to
  the subsequent commit history, not to the content of this raw turn itself.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 170 is next and remains
  pending approval.

### Turn 170 — 2026-05-03 16:53 UTC

- `Turn 170 summary` Steve shifted back from execution to design review and
  asked whether the flat `transports/draft--wire-lab-devs/` layout should gain
  a protocol-grouping layer so a second group or a second transport protocol
  would not make the tree ambiguous. In plain English, this is a pure
  DF-opening turn, not an implementation turn.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  answer identifies two structural weaknesses in the flat layout: it gives no
  protocol grouping by path, and it leaves future multi-instance layouts human-
  ambiguous even if the freeze-time naming rules still technically work. The
  bot then frames `DF-37.1` with three alternatives -- a protocol-slug layer,
  recursive draft/pCID nesting, or defer/status quo -- and recommends the
  protocol-slug solution (`Alt-1.A`) because protocol slugs are already the
  stable human-readable handles used elsewhere in the repo.
- `Gaps or contradictions` None that overturn the existing capture. The later
  residual notes already preserve the important follow-ons from this turn: no
  implementation happened here, the original `DF-37.1` was never directly
  answered as posed, and later turns reframed the problem along the substrate
  axis instead of resolving the protocol-grouping question cleanly on its own
  terms.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 171 is next and remains
  pending approval.

### Turn 171 — 2026-05-03 16:56 UTC

- `Turn 171 summary` Steve refined the tree question from "should there be
  protocol grouping?" to "should there also be a separate path layer meaning
  git file transfer?" In plain English, this is the first turn where the bot
  cleanly separates protocol identity from delivery substrate in the directory
  discussion.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  answer argues that `group-session` answers "what protocol governs these wire
  bytes?" while `git` answers "how do those bytes move between collaborators?"
  and that putting both on the same tree axis would mix two different kinds of
  things. On that basis, the bot rejects a `git/` path layer for now, keeps the
  protocol-slug recommendation from turn 170, and says that if substrate later
  needs first-class representation it should become another axis or per-instance
  metadata rather than another protocol-tree path segment.
- `Gaps or contradictions` None that overturn the existing capture. The later
  residual notes already preserve the important limits of this turn: it is the
  clearest early statement of the protocol-vs-substrate distinction, but it is
  only locally stable because turn 172 immediately expands the substrate
  question from one substrate (`git`) to many and partially overturns the
  "§9 already captures it" comfort.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 172 is next and remains
  pending approval.

### Turn 172 — 2026-05-03 17:01 UTC

- `Turn 172 summary` Steve blew up the narrow `git` question by listing
  multiple peer substrates: `rsync`, `unison`, `uucp`, `udp`, `svn`, `cvs`,
  and `git`. In plain English, this is the turn where the bot stops treating
  substrate as a side detail and starts treating it as a first-class design
  axis.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  answer responds with a substrate-pluralism reframe: if the same
  `group-session` instance might run over several different ways of moving
  bytes, then git-specific rules do not belong buried inside `group-session` as
  if they were part of the transport protocol itself. The bot's big structural
  move is to propose three families -- carrier protocol, transport protocol,
  and substrate-mapping protocol -- and to say that `§9` should eventually move
  out of `group-session` into its own substrate-specific spec. The bot does not
  implement anything here; it proposes a new TE/DF program because the change
  is too large to patch directly.
- `Gaps or contradictions` None that overturn the existing capture. The later
  residual notes already preserve the important nuances: some vocabulary and
  layout specifics from this turn do not survive unchanged, but the deeper idea
  introduced here -- substrate as its own first-class axis rather than a hidden
  appendix inside `group-session` -- remains the load-bearing contribution of
  the turn.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 173 is next and remains
  pending approval.

### Turn 173 — 2026-05-03 17:06 UTC

- `Turn 173 summary` Steve did not answer the proposed TE directly. Instead he
  asked for precedent: is this architecture grounded in real systems, RFCs, and
  historical networks? In plain English, this is a validation turn. Steve is
  testing whether the substrate-pluralism / separate-feed-family idea is a real
  historical pattern or just a fresh bot abstraction.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  answer responds with a long precedent survey -- email over SMTP/UUCP/X.400,
  Usenet over NNTP/UUCP, FidoNet, CORBA GIOP/IIOP, SOAP/WSDL, and modern systems
  such as gRPC, libp2p, Matrix, and git itself -- and extracts three recurring
  patterns: separate substrate specs parallel to the protocol, per-instance
  substrate declarations, and message identity invariant across substrates. The
  strongest framing move is the Usenet comparison: the bot argues that
  `group-session` can be read as a very small content-addressed-Usenet-like
  design, with Usenet's substrate pluralism as the closest historical analog.
  The strongest later non-TE design-visible carry-forward is now
  `docs/research/historical-networks-20260503.md`, and the broader design line
  is now explicitly explored in
  `simulations/SIM-hugoj-cas-usenetlike-gitlike/README.md`.
- `Gaps or contradictions` None that overturn the existing capture. The exact
  phrase from the raw turn was initially preserved only in replay artifacts and
  not in a design doc; that gap has now been narrowed by the later doc
  promotion work under `DI-pijun`. The remaining important limit is that the
  analogy is still exploratory rather than a frozen protocol identity claim, so
  the later docs preserve the lineage and precedent without silently locking
  "PromiseGrid is Usenet" as settled fact.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 174 is next and remains
  pending approval.

### Turn 174 — 2026-05-03 17:13 UTC

- `Turn 174 summary` Steve does not accept the turn-173 proposal as ready to
  draft. Instead he raises three objections at once: `binding` does not sound
  like real Usenet/email vocabulary, the idea seems to conflict with OSI, and
  putting a substrate-spec file inside the messages directory feels inverted.
  In plain English, this is the turn where the proposal is pulled away from a
  W3C/RPC-style framing and re-anchored in the messaging lineage Steve had in
  mind.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  answer accepts all three objections and revises the proposal in three
  matching ways: `binding` is rejected in favor of `feed`; the relationship is
  reframed as alternative delivery/encapsulation rather than OSI-style vertical
  layering; and the nested `bindings/` idea is replaced with an instance-root
  `INSTANCE.md` manifest. This is also the turn where the bot explicitly
  admits that it imported `binding` from WSDL/CORBA/RPC lineage and that the
  more honest historical lineage for this problem is Usenet/FidoNet/email.
  Later residual notes correctly preserve the follow-on consequences: the
  historical-networks note stays accurate but needs vocabulary-aware reading,
  `udp-binding` now looks retroactively misnamed, `INSTANCE.md` may be
  overloading feed facts and governance facts, and the turn-173 Pattern-B map
  must change from `bindings/` to an instance-manifest shape.
- `Gaps or contradictions` None that overturn the existing capture. The main
  later limit is that the encapsulation-not-layering framing introduced here is
  itself not stable for long: turn 175 will immediately swing back toward a
  four-layer model. But that later reversal is already preserved in the UT
  ledger; it does not mean turn 174 was miscaptured.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 175 is next and remains
  pending approval.
