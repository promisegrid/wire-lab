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

- `Turn N plain-English recap`: summarize the user's prompt and the assistant's
  response for the turn, conclusions reached during the turn, later updates or
  modifications to those conclusions found in later turns, and any loose ends or
  open questions that remained as of the end of the turn.
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
  with turn 159 and is now closed under `DI-mugar`.
- `Turn 159` `UT-159.a`, `UT-159.b`, and `UT-159.d` are resolved in
  `TODO-kugod` under `DI-mugar`; `UT-159.c` remains resolved-retired by the
  TE-havib follow-on verification walk.
- `Turn 160` `UT-160.a` is resolved by the TE-havib Cat-3 refinement and
  `TODO-kugod` / `DI-mugar` nine-item sweep closure. `UT-160.b` and
  `UT-160.c` are answered by the TE-havib follow-on verification path.
  `UT-160.d` is resolved by TE-havib's final all-seven-DF locked status.
  `TODO-lilok`'s former reopened harness-spec-sweep note is now reconciled
  through `TODO-kugod` / `DI-mugar`.
- `Turn 161` `UT-161.a` is answered by the TE-havib follow-on disposition.
  `UT-161.b` and `UT-161.c` are now captured as historical inputs by the
  2026-05-12 TE-havib refinement; they are no longer turn-local live work.
- `Turn 162` `UT-162.a` and `UT-162.b` are answered by the later TE-havib
  disposition path. The former `TODO-lilok` reopened sweep-handoff note is
  now closed through `TODO-kugod` / `DI-mugar`.
- `Turn 163` `UT-163.a` is resolved by the apparatus-level `§1.3` rewrite in
  `harness-spec-draft.md` under `DI-lajod` / `DI-mugar`. `UT-163.b` is closed
  as a future-process rule by `AGENTS-ppx.md` B1 / `DI-021-20260507-212249`;
  the commit-specific concerns it cross-referenced live in their own UT rows.
  `UT-cbf7f41-fallback` is retired by the final OQ-36.6 / DF-36.2 negative
  path plus the active `§1.3` apparatus rewrite.
- `Turn 164` `UT-164.a` and `UT-164.b` are historical corrections, not live
  implementation work. `UT-164.c` is resolved by `TODO-bisur` 012.7's
  four-message round-trip. `UT-164.d` is closed by sim-local `TODO-gapab` /
  `DI-rurab`. `UT-164.e` is closed by `TODO-gapab` / `DI-rurab` and
  `TODO-kakaz` / `DI-bomud`; rewrite-at-freeze remains rejected, and turn 164
  now has zero open successor work.
- `Turn 165` `UT-165.a` is closed as an observational privacy/slug lesson;
  `UT-165.b` is closed by the Steve-authored DI promise shape in `DI-rurab`;
  `UT-165.c` was already closed by the neutral-memory update; `UT-165.d` is
  closed by keeping OQ-G4 deferred while treating m000 as valid specimen
  evidence; `UT-165.e` is closed by the group-session example/freeze-gate
  cleanup under `DI-rurab`.
- `Turn 166` `UT-166.a` is closed for future-process purposes by the current
  decision-first protocol and `DI-vanak`; `UT-166.b`, `UT-166.c`, and
  `UT-166.e` are closed by active `DI-rurab` / specimen wording; `UT-166.d`
  is closed as historical git metadata that must not be rewritten.
- `Turn 167` `UT-167.a` is closed by active filename=CID docs; `UT-167.b` and
  `UT-167.c` are closed by `DI-rurab` branch-membership wording; `UT-167.d`
  is closed by the 2026-05-14 `.msg` corpus audit; `UT-167.e` is closed for
  future-process purposes by the current decision-first protocol and
  `DI-vanak`.
- `Turn 168` `UT-168.a` is closed by active `Message-ID:` compatibility
  wording; `UT-168.b`, `UT-168.c`, and `UT-168.d` are closed by `DI-rurab`,
  current §8/§9 wording, and the infrastructure/message-file distinction;
  `UT-168.e` is closed for future-process purposes by decision-first plus
  `DI-vanak`; `UT-168.f` is closed for active docs because current
  wire-lab-devs docs use post-turn-169 CIDs and stale CIDs are historical-only.
- `Turn 169` `UT-169.a`, `UT-169.b`, and `UT-169.e` are closed by
  `DI-012-20260508-033513` / `DI-rurab` and active group-session §4.3 / §4.7
  compatibility wording; `UT-169.c` is closed for future-process purposes by
  decision-first plus `DI-vanak`; `UT-169.d` is closed as historical branch
  metadata because active twig rules require `ppx/{twig}` kebab-case, not the
  historical `ppx/te-<utc>-<slug>` pattern.
- `Turn 170` `UT-170.a` is closed by TE-sihih / TODO-vunub landing the
  L5/L6/L7 substrate-agnostic model and by DR-nugog / `DI-fakin` superseding
  the original flat-versus-nested root `transports/` question for the current
  specimen. `UT-170.b`, `UT-170.c`, and `UT-170.d` are retired and checked off
  as historical naming / stale-summary / superseded-framing records.
- `Turn 171` `UT-171.a` is closed by TE-sihih / TODO-vunub Q-22.2
  plus `DI-rurab` keeping §9 as the current specimen's inline normative git
  binding; `UT-171.b` is closed by TODO-vunub Q-22.3 retracting the manifest
  field idea in favor of path-as-declaration; `UT-171.c` and `UT-171.d` are
  closed as recorded design-cadence lessons.
- `Turn 172` `UT-172.a` is closed as a framing-stability lesson;
  `UT-172.b` is closed by TE-sihih's forward vocabulary (`feed` with
  `substrate` as prose term); `UT-172.c` is closed by TODO-vunub Q-22.3's
  path-as-declaration retraction of per-instance feed manifests; `UT-172.d`
  is closed by TE-sihih's L5/L6/L7 replacement taxonomy; `UT-172.e` is closed
  by `DI-rurab` keeping §9 inline for the current specimen and deferring any
  future feed-spec extraction to successor work.
- `Turn 173` `UT-173.a` is closed by `DI-pijun`, the historical-networks note,
  and `SIM-hugoj` preserving the content-addressed-Usenet line as exploratory
  design evidence; `UT-173.b` is closed by TE-sihih's `feed` vocabulary;
  `UT-173.c` is closed as a cadence lesson refined by turn 174; `UT-173.d` is
  closed by the bounded negative-precedent check added to the research note;
  `UT-173.e` is closed by the research note's CAS-cardinality entry and
  TE-sihih's L5/L6/L7 split.
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
- `Work pending` `no`
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
- `Work pending` `no`
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
  work moved to later turn-159 rows and is now closed under `DI-mugar`.
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
  TE or DF work. Later residual owner work in `TODO-kugod` carried this turn
  forward as `UT-159.a`, `UT-159.b`, and `UT-159.d`; those rows are now closed
  under `DI-mugar`.
- `Gaps or contradictions` None found. The only important timing nuance is
  that the audit memo existed in the workspace at the end of turn 159 but was
  not committed until Steve's turn-160 authorization. `TODO-lilar` already
  preserves that nuance, and later artifacts align with it.
- `Related UTs / owners` `UT-159.a` is resolved by the `DI-mugar`
  harness-spec sweep; `UT-159.b` is resolved by the `DI-huzor` feed-outer
  extraction plus the 2026-05-12 transport companion audit; `UT-159.c` is
  resolved-retired by the TE-havib follow-on verification walk; `UT-159.d` is
  resolved by the `DI-mugar` treatment of §1.3, §10, and §10a.
- `Owner/doc cleanup` Done. `TODO-kugod` closes `kugod.5`, `kugod.8`, and
  `kugod.9`; `TODO-lilok` no longer needs a separate harness-spec-sweep
  handoff; `TODO-jivam` and the UT verification matrix now point at this
  closure. No correction note is needed in `TODO-lilar`.
- `Remaining decisions or work` None for turn 159 itself. Later turns may
  still carry their own UTs, but the turn-159 audit/sweep residue is closed.
- `Work pending` `no`
- `Proposed disposition` `reconciled after apparatus/specimen sweep`
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
- `Related UTs / owners` `UT-160.a` is resolved by the 2026-05-12 Cat-3
  refinement in TE-havib plus `TODO-kugod` / `DI-mugar`, which confirms the
  nine audit items were swept or explicitly retired. `UT-160.b` is answered by
  the TODO-lilok verification walk as procedural-meta, not a live DF split.
  `UT-160.c` is answered by the same verification walk as wrong on inspection.
  `UT-160.d` is resolved by the final TE-havib state: all seven DFs are locked
  after the Alt-B re-presentation path.
- `Owner/doc cleanup` Done. TE-havib now has a Cat-3 refinement that preserves
  the historical "eight" wording but points readers at the nine-item audit
  count and `DI-mugar` closure. `TODO-lilok` is closed; `TODO-lilar` remains
  append-only historical evidence and does not need a correction note.
- `Remaining decisions or work` None for turn 160 itself. Later turns still
  carry their own UTs, but the turn-160 count mismatch, PT/tabletop concerns,
  and end-of-corpus DF-lock concern are reconciled.
- `Work pending` `no`
- `Proposed disposition` `reconciled after TE-havib closure`
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
- `Gaps or contradictions` None that overturn the existing capture. The two
  main carry-forwards from this turn are now reconciled: the asymmetry is
  resolved by TE-havib DF-36.2's retirement of promise-stack as a separate
  envelope hypothesis, and the nine-axis comparison plus richer open-set
  assertion taxonomy are captured as historical inputs in the TE-havib
  2026-05-12 refinement.
- `Related UTs / owners` `UT-161.a` is answered by TE-havib DF-36.2 and the
  TODO-lilok verification path: promise-stack is retired as a separate envelope
  hypothesis, so the asymmetry concern is moot as a live OQ. `UT-161.b` and
  `UT-161.c` are captured by the 2026-05-12 TE-havib refinement as historical
  inputs rather than live downstream requirements.
- `Owner/doc cleanup` Done. TE-havib now records the nine-axis comparison's
  load-bearing result and the assertion-taxonomy examples. No correction note
  is needed in `TODO-lilar`.
- `Remaining decisions or work` None for turn 161 itself. Future assertion
  taxonomy work can still be filed if a later protocol needs it, but that would
  be new work rather than a turn-161 loose end.
- `Work pending` `no`
- `Proposed disposition` `reconciled after TE-havib refinement`
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
  and the promise-stack-home path is later re-presented under Alt-B rather than
  left as a cleanly settled direction. The final TE-havib lock resolves the
  suspicion in the negative by retiring `promise-stack` as a separate protocol
  hypothesis instead of treating it as the active specimen home.
- `Gaps or contradictions` None remaining for turn 162. The two former
  residual concerns are now reconciled: `OQ-36.6` is visibly resolved in
  TE-havib, and `DF-36.2` was re-presented under Alt-B as Alt-2.A revised
  rather than left as a provisional promise-stack-home decision.
- `Related UTs / owners` `UT-162.a` is resolved by TE-havib's `OQ-36.6`
  negative-resolution path and the TODO-lilok verification walk. `UT-162.b` is
  resolved by the Alt-B re-presentation of `DF-36.2` and the final Alt-2.A
  revised lock.
- `Owner/doc cleanup` Done. TE-havib already carries the final `OQ-36.6` and
  `DF-36.2` resolution text; TODO-lilok is closed; the verification matrix now
  has a turn-162 closure pointer.
- `Remaining decisions or work` None for turn 162 itself.
- `Work pending` `no`
- `Proposed disposition` `reconciled after TE-havib final lock`
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
- `Gaps or contradictions` None remaining for turn 163. The turn's
  envelope-agnostic `§1.3` template no longer lives only in conversation:
  active `harness-spec-draft.md` now carries an apparatus-level
  layering-scenarios section under `DI-lajod`, later completed for the wider
  turn-159 audit by `DI-mugar`. The uncaptured-commit concern is also no longer
  a turn-local work item: its future-process lesson is captured in
  `AGENTS-ppx.md` B1, while any commit-specific residue is carried by the
  separately named UT rows.
- `Related UTs / owners` `UT-163.a` is resolved by the `§1.3` apparatus-level
  rewrite in `protocols/wire-lab.d/specs/harness-spec-draft.md`. `UT-163.b` is
  closed as a procedural rule by `AGENTS-ppx.md` B1 / `DI-021-20260507-212249`.
  `UT-cbf7f41-fallback` is retired for turn-163 purposes by OQ-36.6's negative
  resolution, DF-36.2's promise-stack retirement, and the active `§1.3`
  apparatus rewrite.
- `Owner/doc cleanup` Done. The verification matrix now has a turn-163 closure
  pointer. No `TODO-lilar` checkbox is flipped; the matrix remains the closure
  index.
- `Remaining decisions or work` None for turn 163 itself. Commit-specific rows
  not owned by this turn remain under their own UT identifiers if later replay
  reaches them.
- `Work pending` `no`
- `Proposed disposition` `reconciled after §1.3 rewrite and B1 transfer`
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
  later carry-forwards are now split by owner and closed: the "4-DF-locked"
  overclaim and TE-numan revert/merge discrepancy are historical corrections;
  the group-session round-trip coverage issue was later closed by the
  four-message `TODO-bisur` 012.7 artifact; membership-by-posting is closed by
  `DI-rurab`; freeze-rewrite is rejected by `DI-bomud`.
- `Related UTs / owners` `UT-164.a` is a historical correction against the
  current TE-havib status. `UT-164.b` is a historical correction; TE-numan
  scope was later settled by the TE-39 / TODO-lilar cross-cutting disposition.
  `UT-164.c` is resolved by `TODO-bisur` 012.7. `UT-164.d` is closed by
  sim-local `TODO-gapab` / `DI-rurab`. `UT-164.e` is closed by
  `TODO-gapab` / `DI-rurab` and `TODO-kakaz` / `DI-bomud`; the no-rewrite
  freeze boundary is now explicit in the feed-outer and wire-lab-devs docs.
- `Owner/doc cleanup` Done. The verification matrix has a turn-164 closure
  pointer. Rooted `TODO-turog` and `TODO-duvuk` remain historical coordination
  memory; active specimen-owned successor work is closed in the sim-local
  TODOs.
- `Remaining decisions or work` None for turn 164.
- `Work pending` `no`
- `Proposed disposition` `reconciled after successor closure`
- `Write needed? yes/no` `no`
- `Next` Turn 165 processed by the next chronological replay step.

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
- `Gaps or contradictions` Captured and now resolved for turn-165 cleanup
  purposes. The "reverted cleanly" claim remains owned by `UT-164.b`, already
  reconciled as historical evidence. `UT-165.a` is closed as an observational
  privacy/slug lesson because the active specimen uses the generic
  `wire-lab-devs` slug and current group-session examples use neutral
  Alice/Bob prose. `UT-165.b` is closed by `DI-rurab`, which defines the
  interim `merge-group-transport-spec` shape as a Steve-authored DI until
  cryptographic promise tooling exists. `UT-165.c` was already closed by the
  neutral-memory update. `UT-165.d` is closed without a spec edit: OQ-G4 remains
  deferred, and m000 is one valid first-message pattern rather than a v0
  genesis-message mandate. `UT-165.e` is closed by the group-session spec
  cleanup that removed `codex-perplexity` examples and names the
  wire-lab-devs specimen as the freeze-gate evidence.
- `Related UTs / owners` `UT-165.a`, `UT-165.b`, `UT-165.d`, and
  `UT-165.e` are now checked off in `TODO-lilar`; `UT-165.c` was already
  checked off. No active owner TODO remains for turn 165.
- `Owner/doc cleanup` Done. `TODO-lilar` UT rows are closed; `TODO-juhub`
  carries the stronger turn note; `group-session-draft.md` and
  `wire-lab-devs-draft/README.md` already carry the resulting active wording
  from the `DI-rurab` / `DI-bomud` cleanup.
- `Remaining decisions or work` None for turn 165.
- `Work pending` `no`
- `Proposed disposition` `reconciled after lilar UT closure`
- `Write needed? yes/no` `no`
- `Next` Turn 166 is next.

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
- `Gaps or contradictions` Captured and now resolved for turn-166 cleanup
  purposes. `UT-166.a` is closed for future-process purposes by the current
  decision-first protocol plus `DI-vanak`'s explicit replay approval shorthand;
  historical bootstrap commits remain preserved and are not rewritten.
  `UT-166.b` is closed by `DI-rurab`: active membership is the fixed configured
  set of exact `<author-id>/main` branches, so guessed actors are not enrolled
  by speculation. `UT-166.c` is closed by active specimen docs that use
  `stevegt-via-perplexity` as the committed `From:` identity. `UT-166.d` is
  closed as historical git metadata; rewriting the stale twig name would be a
  history rewrite. `UT-166.e` is closed by `DI-rurab`, which supersedes
  membership-by-posting with fixed configured branch membership, passive
  observer non-membership, and no self-admission from unknown branches.
- `Related UTs / owners` `UT-166.a`, `UT-166.b`, `UT-166.c`, `UT-166.d`, and
  `UT-166.e` are now checked off in `TODO-lilar`. Later turn rows that cite
  `UT-166.a` remain future-turn work only where their own turn-specific
  execution pattern still needs reconciliation.
- `Owner/doc cleanup` Done. `TODO-lilar` turn-166 UT rows are closed;
  stale future-row references to `UT-166.a` as "pending" were updated to point
  at the resolved process baseline; active group-session and wire-lab-devs
  docs already carry the `DI-rurab` membership and identity-safe wording.
- `Remaining decisions or work` None for turn 166.
- `Work pending` `no`
- `Proposed disposition` `reconciled after lilar UT closure`
- `Write needed? yes/no` `no`
- `Next` Turn 167 is next.

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
- `Gaps or contradictions` Captured and now resolved for turn-167 cleanup
  purposes. `UT-167.a` is closed because current active docs use filename = CID
  and no active sequential `m<N>-...` filename rule remains. `UT-167.b` is
  closed by `DI-rurab`, which rejects membership-by-posting for active wording
  and uses fixed configured `<author-id>/main` membership. `UT-167.c` is
  closed by the same `DI-rurab` wording: `{name}` now means author-id, not
  group name. `UT-167.d` is closed by the 2026-05-14 corpus audit for
  `\.msg\b`; remaining matches are historical replay/disposition notes or
  message-body evidence, not active spec guidance. `UT-167.e` is closed for
  future-process purposes by the current decision-first protocol plus
  `DI-vanak`; later execute-on-directive rows keep their own turn-specific
  reconciliation work.
- `Related UTs / owners` `UT-167.a`, `UT-167.b`, `UT-167.c`, `UT-167.d`, and
  `UT-167.e` are now checked off in `TODO-lilar`. The active group-session
  and wire-lab-devs docs already carry the filename, branch-membership, and
  passive-observer wording needed for this turn.
- `Owner/doc cleanup` Done. `TODO-lilar` turn-167 UT rows are closed;
  `TODO-juhub` carries this stronger note; the verification matrix has a
  turn-167 closure pointer. No active `.msg` spec-edit target was found.
- `Remaining decisions or work` None for turn 167.
- `Work pending` `no`
- `Proposed disposition` `reconciled after lilar UT closure and .msg audit`
- `Write needed? yes/no` `no`
- `Next` Turn 168 is next.

### Turn 168 — 2026-05-03 11:45 UTC

- `Turn 168 summary` Steve corrected the new branch-binding model from turn
  167 in two important ways: message files must not use sequential numbers, and
  before posting, each agent must first merge all observed messages from all
  branches into the directory on their own branch and push that merged state.
  In plain English, each participant's branch becomes "my replicated view of
  the message set, plus optionally a new post," and the message CID becomes the
  stable on-disk identity.
- `Existing capture` `TODO-lilar` captures the turn correctly. The bot treated
  both directives as executable: it rewrote filename rules to CID filenames,
  expanded §9 into a receive/merge/push/optionally-post cycle, renamed the
  existing bootstrap files to then-current CIDs, and authored the m2 on-wire
  ratification message. The turn fixes the global-sequence mistake from turn
  167 and defines the first actual replication model for the Git-backed
  specimen.
- `Gaps or contradictions` Captured and now resolved for active docs. The
  `Message-ID:` retention is closed by the active §4.3/§4.7 compatibility
  rule, §9's old non-normative status is closed by `DI-rurab`, passive-reader
  ambiguity is closed by current §8/§9.3, infrastructure/message-file
  boundaries are explicit enough for the current specimen, and stale turn-168
  CIDs are historical-only after turn 169's rehash.
- `Related UTs / owners` `UT-168.a` through `UT-168.f` are checked off in
  `TODO-lilar`. The only deeper rehash-continuity questions continue under the
  turn-169 rows, not as turn-168 work.
- `Owner/doc cleanup` Done. Updated the active `TODO-bisur` 012.7 note so it no
  longer calls §9 explicitly non-normative or keeps a stale open §9 OQ. No
  transport message files were changed.
- `Remaining decisions or work` None for turn 168.
- `Work pending` no.
- `Proposed disposition` `reconciled after lilar UT closure and active-doc CID/status audit`
- `Write needed? yes/no` `no` further turn-168 write is needed after this pass.
- `Next` Turn 169 is next.

### Turn 169 — 2026-05-03 11:54 UTC

- `Turn 169 summary` Steve asked whether `Message-ID:` is still needed now that
  the canonical identifier is the message CID. In plain English, the raw turn is
  a careful reasoning memo: the bot audits every plausible use of `Message-ID:`,
  concludes that it creates competing identity once filename = CID, and then
  reasons about compatibility for the already-authored bootstrap messages.
- `Existing capture` `TODO-lilar` captures the turn correctly. The raw reasoning
  recommended the conservative Path A (deprecate / preserve legacy bytes / ask
  Steve), while later repo history executed Path B (hard remove / rehash). The
  later active policy is no longer the unqualified Path-B-only state: current
  group-session §4.3 / §4.7, under `DI-012-20260508-033513`, says canonical
  writers MUST NOT emit `Message-ID:`, but readers MAY tolerate exactly one
  legacy pre-`Date:` header and MUST ignore its value semantically.
- `Gaps or contradictions` Captured and now resolved for active docs. The
  reasoning/action inversion is explicitly recorded; the strict-reader problem
  has an explicit legacy-header carve-out; the writer-side prohibition versus
  reader-side deprecation split is now deliberate; and the twig-name issue is
  historical metadata because active branch rules require only `ppx/{twig}` with
  a short kebab-case task phrase.
- `Related UTs / owners` `UT-169.a` through `UT-169.e` are checked off in
  `TODO-lilar`. `TODO-duvuk` is closed as historical coordination memory and no
  longer owns active execution for these Message-ID / CID-cascade rows.
- `Owner/doc cleanup` Done. Updated `TODO-duvuk` so its original
  T-FILENAME-CID-CASCADE scope is visibly historical and updated `TODO-jivam`
  so its live TE-42 closure condition no longer appears open. Historical quoted monitor snapshots remain unchanged as evidence.
- `Remaining decisions or work` None for turn 169. Later nested-body and
  grid-envelope cascade questions remain with their own later-turn rows, not as
  turn-169 work.
- `Work pending` no.
- `Proposed disposition` `reconciled after lilar UT closure and active Message-ID compatibility audit`
- `Write needed? yes/no` `no` further turn-169 write is needed after this pass.
- `Next` Turn 170 is next.

### Turn 170 — 2026-05-03 16:53 UTC

- `Turn 170 summary` Steve shifted back from execution to design review and
  asked whether the flat `transports/draft--wire-lab-devs/` layout should gain
  a protocol/grouping layer, especially if a second named group appears. In
  plain English, this opened the directory-axis question that later turns
  broadened into protocol identity, feed/substrate, site, CAS, and simulation
  placement.
- `Existing capture` `TODO-lilar` captures the turn correctly. The raw answer
  framed `DF-37.1` with three alternatives: protocol-slug nesting, recursive
  draft/pCID nesting, or status quo deferral. It recommended protocol-slug
  nesting, but no implementation happened and Steve did not answer that DF as
  posed.
- `Gaps or contradictions` Captured and now resolved. `DF-37.1` is not still an
  actionable open DF: turns 171-176 reframed the missing axis, TE-sihih / TODO-
  vunub landed the L5/L6/L7 substrate-agnostic model, TE-domat and DR-nugog
  reframed the root `transports/` / `groups/` question, and `DI-fakin` resolved
  the current specimen by moving it into a simulation world instead of choosing
  any of the original root flat/nested alternatives.
- `Related UTs / owners` `UT-170.a` through `UT-170.d` are checked off in
  `TODO-lilar`. `TODO-vunub` is now marked closed because TE-sihih is decided;
  DR-nugog is implemented for the current specimen and has an append-only note
  naming the current `SIM-ludut-wire-lab-devs` path after the later rusis split.
- `Owner/doc cleanup` Done. Updated TODO-vunub status, the master TODO row, the
  verification matrix owner/closure notes, and DR-nugog current-path provenance.
- `Remaining decisions or work` None for turn 170. Later root/reference-layout
  questions, if any, are downstream graduation questions rather than turn-170
  recovery loose ends.
- `Work pending` no.
- `Proposed disposition` `reconciled after TE-sihih / DR-nugog / simulation-path audit`
- `Write needed? yes/no` `no` further turn-170 write is needed after this pass.
- `Next` Turn 171 is next.

### Turn 171 — 2026-05-03 16:56 UTC

- `Turn 171 summary` Steve refined the turn-170 tree question from "should
  there be protocol grouping?" to "should there also be a separate path layer
  meaning git file transfer?" In plain English, this is where the replay first
  separates protocol identity from delivery substrate, before turn 172 expands
  the substrate axis beyond git.
- `Existing capture` `TODO-lilar` captures the turn correctly. The raw answer
  rejects a `git/` path layer for the current one-substrate case, keeps the
  protocol-slug recommendation from turn 170, and says substrate facts should
  become another axis or per-instance metadata if they later need first-class
  representation.
- `Gaps or contradictions` Captured and now resolved. The "§9 captures git"
  comfort is superseded by TE-sihih's L5 feed model, but the active
  group-session specimen is also intentionally allowed to keep §9 inline as the
  normative wire-lab-devs git binding under `DI-rurab`. The manifest-field hook
  is closed by TODO-vunub Q-22.3, which retracts manifest schema work in favor
  of TE-vipir path-as-declaration.
- `Related UTs / owners` `UT-171.a` through `UT-171.d` are checked off in
  `TODO-lilar`. `TODO-vunub` records the retired turn-171 substrate/manifest
  framing; cadence-only rows are preserved as lessons, not live implementation
  owners.
- `Owner/doc cleanup` Done. Updated TODO-vunub's retired-question and DI-log
  notes. No active spec or transport-message bytes were changed.
- `Remaining decisions or work` None for turn 171. Later feed-vocabulary,
  groups/transports, and nested-envelope questions remain with their own later
  turns, not with turn 171.
- `Work pending` no.
- `Proposed disposition` `reconciled after TE-sihih / TODO-vunub / DI-rurab audit`
- `Write needed? yes/no` `no` further turn-171 write is needed after this pass.
- `Next` Turn 172 is next.

### Turn 172 — 2026-05-03 17:01 UTC

- `Turn 172 summary` Steve blew up the narrow `git` question by listing
  multiple peer substrates: `rsync`, `unison`, `uucp`, `udp`, `svn`, `cvs`,
  and `git`. In plain English, this is the turn where the replay stops
  treating delivery substrate as a side detail and starts treating it as a
  first-class design axis.
- `Existing capture` `TODO-lilar` captures the turn correctly. The raw answer
  proposes a new TE/DF program because if one `group-session` can move over
  multiple byte-moving substrates, then git-specific rules cannot be treated as
  timeless core group semantics. The answer's specific words and layout are
  transitional: it says `binding`, sketches `bindings/` and `messages/`
  subdirectories, and imagines extracting `§9` into a separate git-specific
  spec.
- `Gaps or contradictions` Captured and now resolved. The load-bearing
  insight survives, but the exact turn-172 taxonomy does not. TE-sihih replaces
  `binding` with L5 `feed`, adds L6 CAS and L7 group as the citable layer
  split, retracts per-instance feed manifests in favor of TE-vipir
  path-as-declaration, and leaves current wire-lab-devs git rules inline in
  group-session `§9` under `DI-rurab`.
- `Related UTs / owners` `UT-172.a` through `UT-172.e` are checked off in
  `TODO-lilar`. `TODO-vunub` records the turn-172 proposal as retired by
  TE-sihih / Q-22.2 / Q-22.3 / Q-22.6, with future feed-spec extraction owned
  by successor work rather than this replay turn.
- `Owner/doc cleanup` Done. Updated `TODO-lilar`, `TODO-vunub`, this turn note,
  and the UT verification matrix. No active spec or transport-message bytes
  were changed.
- `Remaining decisions or work` None for turn 172. Later historical-analog,
  feed-vocabulary, CAS/site, and simulation-layout details belong to their own
  later turns and owner artifacts.
- `Work pending` no.
- `Proposed disposition` `reconciled after TE-sihih / TODO-vunub / DI-rurab audit`
- `Write needed? yes/no` `no` further turn-172 write is needed after this pass.
- `Next` Turn 173 is next.

### Turn 173 — 2026-05-03 17:06 UTC

- `Turn 173 plain-English recap` Steve did not answer the proposed TE directly.
  He asked whether the substrate-pluralism idea had precedent in practice,
  RFCs, and historical networks. The assistant answered with a precedent survey
  covering email over SMTP/UUCP/X.400, Usenet over NNTP/UUCP, FidoNet, CORBA
  GIOP/IIOP, SOAP/WSDL, modern pluggable-transport systems, and git itself. The
  turn's conclusion was that the broad concept is well precedented: message
  identity can stay stable while different substrates move the same bytes. The
  strongest plain-English insight was that `group-session` looks like a tiny
  content-addressed Usenet, but later turns modify the details: `binding` is
  rejected in favor of `feed`, the `bindings/` layout sketch is retracted, and
  the Usenet line is captured as exploratory design evidence rather than a
  frozen claim that PromiseGrid is Usenet. At the end of the turn, the open
  questions were naming, layout, negative counter-precedent, and git/CAS
  cardinality.
- `Existing capture` `TODO-lilar` captures the raw turn and its five residual
  rows. Later promotion work under `DI-pijun` puts the historical survey in
  `docs/research/historical-networks-20260503.md` and the broader design line in
  `simulations/SIM-hugoj-cas-usenetlike-gitlike/README.md`.
- `Gaps or contradictions` Captured and now resolved for turn-173 recovery. The
  exact Usenet analogy is now design-visible, the later vocabulary correction is
  recorded, and the negative-precedent gap has been narrowed by a bounded check
  in the research note. The analogy remains deliberately exploratory; future
  TEs may import or reject specific Usenet mechanisms one at a time.
- `Related UTs / owners` `UT-173.a` through `UT-173.e` are checked off in
  `TODO-lilar`. `TODO-vunub` records the turn-173 historical-precedent questions
  as retired by `DI-pijun`, the historical-networks note, `SIM-hugoj`, and
  TE-sihih's feed/CAS/layering decisions.
- `Owner/doc cleanup` Done. Updated `TODO-lilar`, `TODO-vunub`, the
  historical-networks research note, this turn note, and the UT verification
  matrix. No active spec or transport-message bytes were changed.
- `Remaining decisions or work` None for turn 173. Later feed-vocabulary,
  site/CAS, and simulation-layout decisions remain with their own later turns
  and owner artifacts.
- `Work pending` no.
- `Proposed disposition` `reconciled after DI-pijun / historical-networks / SIM-hugoj / TE-sihih audit`
- `Write needed? yes/no` `no` further turn-173 write is needed after this pass.
- `Next` Turn 174 is next.

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
