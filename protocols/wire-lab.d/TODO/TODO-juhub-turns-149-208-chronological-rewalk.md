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
3. Compare the result against existing replay artifacts, especially
   `TODO-lilar`, `dropped-thread-disposition-20260506.md`,
   `ut-verification-matrix-20260507.md`, and any current owner TODOs.
4. Classify the turn as exactly one of: `clean`, `already captured`,
   `captured but correction needed`, `missing owner artifact`, or
   `new contradiction`.
5. Report the turn back to Steve using the fixed report format below and stop.
6. After approval, write the turn result here and make any linked owner or
   correction-note updates before moving to the next turn.

## Turn report format

- `Turn N summary`
- `Existing capture`
- `Gaps or contradictions`
- `Proposed disposition`
- `Write needed? yes/no`
- `Next: wait for approval before turn N+1`

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
- [ ] juhub.160 Turn 160 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.161 Turn 161 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.162 Turn 162 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.163 Turn 163 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.164 Turn 164 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.165 Turn 165 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.166 Turn 166 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.167 Turn 167 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.168 Turn 168 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.169 Turn 169 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.170 Turn 170 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.171 Turn 171 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.172 Turn 172 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.173 Turn 173 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.174 Turn 174 raw-log rewalk plus later-turn and later-artifact sweep.
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
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 158 is next and remains
  pending approval.

### Turn 158 — 2026-05-03 01:46 UTC

- `Turn 158 summary` This is the real apparatus-vs-specimen break. Steve
  pointed out that even calling the envelope "harness-wide" smuggles in the
  wrong assumption, because wire-lab exists to test multiple hypotheses at all
  layers, not to bake one answer into the harness. In plain English, this is
  the turn where the bot finally accepts that the harness-spec is the lab
  apparatus and the candidate envelopes are specimens under study.
- `Existing capture` `TODO-lilar` already captures the turn correctly as the
  foundational apparatus-vs-specimen reframe. The raw turn lays out the six
  step sequence that follows from that correction: audit the harness-spec,
  file the harness-level TE on the split, give each candidate envelope its own
  protocol home, sweep specimen material out of the harness-spec, reframe the
  old promise-stack TODO under protocol ownership, and file a parallel TODO for
  the grid-envelope hypothesis. Later residual owner work also preserves the
  two load-bearing carry-forwards from this turn: the harness-spec carve-outs
  still open under `TODO-kugod`, and the future grid-envelope successor work
  also still open there.
- `Gaps or contradictions` None found. Later artifacts continue to treat this
  as the foundational TE-havib turn. The insight itself stands, while the
  downstream cleanup remains only partially complete. I did not find any later
  artifact that reverts to the old claim that the harness-spec should define a
  single envelope specimen.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
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
