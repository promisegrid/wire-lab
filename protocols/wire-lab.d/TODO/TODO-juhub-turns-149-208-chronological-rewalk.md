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
- [ ] juhub.153 Turn 153 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.154 Turn 154 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.155 Turn 155 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.156 Turn 156 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.157 Turn 157 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.158 Turn 158 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.159 Turn 159 raw-log rewalk plus later-turn and later-artifact sweep.
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
