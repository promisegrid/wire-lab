# TODO-tavit: Turns 209-311 quality audit

## Prior aliases

None. This TODO was minted after the proquint migration and has no integer or
timestamp alias.

## Status

Open. This TODO owns a targeted quality audit for session `ea135ce8` turns
209-311. It is not an extension of the completed `TODO-juhub` 149-208 recovery
rewalk unless this audit discovers evidence that a full successor rewalk is
needed.

## Decision Intent Log

ID: DI-duluk
Date: 2026-05-17 23:33:32
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: File `TODO-tavit` as a standalone quality-audit owner for turns
209-311, using a fast triage pass first and escalating only suspicious turns
into full per-turn rewalk notes.
Intent: Treat the post-208 slice as higher risk because Steve suspects the
Perplexity computer-service agent may have shifted to a lower-quality Claude
path around the turn-190-through-194 period, while avoiding an automatic
103-turn rewalk that would distract from closing the already-bounded 149-208
recovery work.
Constraints: Do not reopen or extend `TODO-juhub` by default. Do not presume
every turn 209-311 artifact is wrong. Preserve historical notes append-only.
Route PromiseGrid design loose ends into simulation questions or DRs, and route
harness/process loose ends into TODOs or DRs. If a turn's loose ends are fully
captured in simulation questions for PromiseGrid design or TODOs for harness
work, mark that turn's work as not pending.
Affects: `protocols/wire-lab.d/TODO/TODO-tavit-turns-209-311-quality-audit.md`;
`protocols/wire-lab.d/TODO/TODO.md`;
`/home/stevegt/lab/session-logs/sessions/ea135ce8/209-turn.md` through
`/home/stevegt/lab/session-logs/sessions/ea135ce8/311-turn.md`.

ID: DI-rapom
Date: 2026-05-18 10:19:51
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Add a source-fidelity gate before the fast quality screen for turns
209-311. The audit must classify each turn source as raw, reconstructed,
summary-only, bundled, or indirect, and must resolve required bundle pointers
before interpreting any reconstructed turn.
Intent: A file named `NNN-turn.md` can exist while still being insufficient
evidence. Turns 308-311 are the concrete hazard: their per-turn files are
reconstructed summaries that point at the loose bundled response
`/home/stevegt/lab/session-logs/turns/turn-314.md`. The quality audit must not
mistake such summaries for full raw exchange evidence.
Constraints: Do not treat source-file presence as source fidelity. If a turn
points at a bundle, transcript, loose file, or other indirect source, verify the
target path exists, verify which turns it covers, and use it alongside the
per-turn summary before classifying quality or closure. If a required source is
missing, mark the turn source-blocked instead of safe.
Affects: `protocols/wire-lab.d/TODO/TODO-tavit-turns-209-311-quality-audit.md`;
`/home/stevegt/lab/session-logs/sessions/ea135ce8/209-turn.md` through
`/home/stevegt/lab/session-logs/sessions/ea135ce8/311-turn.md`;
`/home/stevegt/lab/session-logs/turns/turn-314.md`.

## Decision lock summary

- `D-tavit-architecture`: use a standalone TODO owner rather than extending
  `TODO-juhub` or `TODO-jivam`.
- `D-tavit-behavior`: perform a fast quality screen first; escalate only
  suspicious turns into full per-turn rewalk notes.
- `D-tavit-implementation`: track phases and escalation criteria here, then
  route any discovered loose ends to the correct downstream owner artifacts.
- `D-tavit-path-1`: create
  `protocols/wire-lab.d/TODO/TODO-tavit-turns-209-311-quality-audit.md` as a
  permanent harness TODO artifact.
- `D-tavit-path-2`: update `protocols/wire-lab.d/TODO/TODO.md` to cross-list
  `TODO-tavit`.

## Audit question

Are turns 209-311 low-quality enough that they need a full chronological rewalk,
or can the repo safely close only the specific weak claims, stale pointers,
unowned decisions, and overreaches found by a targeted audit?

## Source-fidelity gate

Before `tavit.2` fast-screening begins, `tavit.1` must build a source-fidelity
table for every turn 209-311:

- Classify the available evidence as `raw`, `reconstructed`, `summary-only`,
  `bundled`, or `indirect`.
- Record the exact source path or bundle path used for the turn.
- For reconstructed or summary-only turns, identify the full source that the
  summary points to and verify that it exists before interpreting the turn.
- For bundled turns, record the covered turn range and verify alignment between
  the per-turn summary and the bundle.
- If a required source is absent or ambiguous, mark the turn `source-blocked`
  and escalate it instead of treating it as safe.

Known source-fidelity hazard: turns 308-311 are reconstructed summaries under
`/home/stevegt/lab/session-logs/sessions/ea135ce8/` and point to the loose
bundle `/home/stevegt/lab/session-logs/turns/turn-314.md`; the audit must read
that bundle before classifying those turns.

## Quality screen criteria

Flag a turn or artifact for escalation when it shows any of these patterns:

- A decision is asserted as settled but no DI, DR, TODO, simulation question, or
  design doc owns it.
- A protocol/specimen conclusion is stored only in a harness/process TODO.
- A simulation is described as an "active" or preferred specimen home instead
  of one independently evolving candidate.
- A doc or TODO routes work to a stale path, stale sim name, or generic sim when
  a standalone sim should own the question.
- A closure statement says work is done even though the remaining question is
  not routed to a simulation question, TODO, or DR.
- A model/process claim is plausible but unsupported, especially around the
  suspected Perplexity model-quality transition.

## Subtasks

- [ ] tavit.1 Verify source coverage and source fidelity for turns 209-311:
  record missing files, reconstructed summaries, summary-only turns, bundled
  sources, indirect-source pointers, bundle coverage ranges, and any
  source-blocked turns before interpreting the slice.
- [ ] tavit.2 Fast-screen turns 209-311 for low-quality reasoning, unsupported
  closure, stale routing, specimen/harness conflation, and unowned design
  decisions.
- [ ] tavit.3 Inventory repo artifacts created or materially changed by the
  turn-209-through-311 slice, with special attention to TODO, TE, DR, DI,
  simulation, and design-doc outputs.
- [ ] tavit.4 Route every discovered PromiseGrid design loose end to a
  simulation question or DR, and every harness/process loose end to a TODO or
  DR.
- [ ] tavit.5 Escalate any suspicious turn into a full per-turn note using the
  stronger `TODO-juhub` report format, including plain-English recap, later
  updates, conclusions, loose ends, and `Work pending: yes/no`.
- [ ] tavit.6 Produce a closeout summary that lists safe turns, escalated turns,
  downstream owner artifacts, and whether a full 209-311 chronological rewalk is
  still warranted.
