# TODO-jivam: Turns 149-208 recovery completion monitor

## Status

Open. This TODO shall not close until all recovery represented by the
turns 149-208 bounded slice is complete: each source turn is accounted
for, each UT or equivalent recovery item is resolved, retired, or
explicitly transferred under its owner artifact, and the owner artifacts
that block recovery are closed or superseded under DI/DR provenance. The
original turns 149-170 monitor report remains preserved below as history;
`DI-zufar` expands the live closure gate through turn 208 without
renaming this file.

## Decision Intent Log

ID: DI-jivam
Date: 2026-05-09 17:55:15
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Create a dedicated open TODO that preserves the full turns 149-170 recovery monitor report and stays open until that recovery is complete.
Intent: Keep the Sunday context-loss recovery visible as a first-class work item instead of relying on TODO-lilar alone or on transient agent reports. The monitor TODO is a closure guard: it prevents a completed replay walk or a completed sub-cascade from being mistaken for complete recovery while TE-40, TE-41, TE-42, and adjacent UT owners still carry unresolved recovery work.
Constraints: This TODO is a monitoring and closure-gate artifact. It does not close UTs directly, does not flip TODO-lilar checkboxes, and does not supersede the matrix-as-closure-index rule locked by `DI-021-20260507-210204`. Closure requires evidence that all turns 149-170 recovery items are represented and complete through their proper owner artifacts.
Affects: `protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`; `protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md`; `protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`; `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`; `protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md`; `protocols/wire-lab.d/TODO/TODO-duvuk-te-42-filename-cid-cascade-policy.md`; `protocols/wire-lab.d/TODO/TODO.md`.

ID: DI-zufar
Date: 2026-05-09 11:24:51
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Expand TODO-jivam from a turns 149-170 recovery monitor to a turns 149-208 recovery monitor while preserving the original 149-170 report as historical evidence.
Intent: Steve corrected the stopping point: turn 170 was not the right boundary for the recovery monitor. The live monitor must cover the later recovery/session-log boundary sequence through turn 208, including the turn-195 context-loss correction, the turn-196 72-hour ledger request, the turn-197 instruction that created TODO-lilar and bounded its chronological walk to 149-192, and the turn-208 instruction to keep the ledger in TODO 021 while new promise-stack provenance work was logged. The expanded monitor keeps the complete 149-208 recovery slice visible without rewriting the prior 149-170 snapshot.
Constraints: Preserve the existing filename unless a later explicit decision authorizes a rename. Do not edit specs, transports, TODO-kugod, or DR-nugog for this correction. Use TODO-jivam as a monitoring and closure-gate artifact only: it does not directly close UTs, does not flip TODO-lilar or disposition checkboxes, and does not supersede the matrix-as-closure-index mechanism. Treat turns 193-208 as session-log boundary evidence unless or until a later owner artifact accounts for them under DI/DR provenance.
Affects: `protocols/wire-lab.d/TODO/TODO-jivam-turns-149-170-recovery-completion.md`; `protocols/wire-lab.d/TODO/TODO.md`; `/home/stevegt/lab/session-logs/sessions/ea135ce8/149-turn.md` through `/home/stevegt/lab/session-logs/sessions/ea135ce8/208-turn.md`; `protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`; `protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md`; `protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`.
Supersedes: (none — `DI-jivam` remains the creation decision and historical 149-170 monitor snapshot; this DI expands the live closure boundary.)

## Closure rule

This TODO may close only after all of the following are true:

- [ ] jivam.1 Source coverage for turns 149-208 remains verified against `/home/stevegt/lab/session-logs/sessions/ea135ce8/149-turn.md` through `/home/stevegt/lab/session-logs/sessions/ea135ce8/208-turn.md`.
- [ ] jivam.2 Turns 149-154 remain verified as walked with no leftover threads, or any later contradiction is filed as a new UT/DR with provenance.
- [ ] jivam.3 All turn-155-through-192 UTs in `dropped-thread-disposition-20260506.md` have owner artifacts and closure evidence in the verification matrix.
- [ ] jivam.4 TODO-kugod closes or transfers all residual TE-40 recovery work after the apparatus-vs-specimen sweep, grid-envelope home, transport-spec companion audit, and stale TE-famar status/path issues are resolved or explicitly retired.
- [ ] jivam.5 TODO-turog closes or transfers all TE-41 group-session freeze work after its TE-40 and migration-design blockers are resolved.
- [ ] jivam.6 TODO-duvuk closes or transfers all TE-42 Message-ID / filename / CID-cascade policy work after TE-41 no longer blocks it.
- [ ] jivam.7 Adjacent turn-149-through-192 items outside TE-40/41/42 are resolved or explicitly transferred: TE-havib follow-on, TE-sihih, TODO-kituj/TE-43, TODO-ralud/TE-45, Spec-edit, Retire, and Carry entries.
- [ ] jivam.8 TODO-lilar remains open until every relevant downstream owner has completed, retired, or transferred its UTs under the matrix-as-closure-index rule.
- [ ] jivam.9 Turns 193-208 are explicitly accounted for as post-lilar session-log boundary material: either no additional recovery item is required, or each item is filed in an owner artifact with DI/DR provenance.
- [ ] jivam.10 Boundary turns 195, 196, 197, and 208 receive explicit final-review evidence because they define the context-loss correction, the 72-hour ledger request, the TODO-lilar creation boundary, and the later "keep the ledger in TODO 021" instruction.
- [ ] jivam.11 A final recovery monitor pass confirms no closed owner artifact still contains unfinished recovery work for turns 149-208.

## Monitor report snapshot

The following snapshot is preserved as the original turns 149-170 report.
It is historical evidence, not the current live scope after `DI-zufar`.

Read-only monitoring report received 2026-05-09:

> **Monitoring Report**
> - Read-only status: no files edited. Current working tree already has uncommitted changes in `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md` and `protocols/wire-lab.d/TODO/TODO.md`.
> - Source coverage: turn files exist from `/home/stevegt/lab/session-logs/sessions/ea135ce8/149-turn.md:1` through `/home/stevegt/lab/session-logs/sessions/ea135ce8/170-turn.md:1`; TODO-lilar’s mandate covers turns 149–192 at `protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md:14`.
> - Turns 149–154 are walked and explicitly have “No leftover threads,” so they are not in the disposition UT inventory: `protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md:23`–`protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md:28`.
> - The disposition file covers 186 UTs spanning turns 155–192, not 149–154: `protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md:1` and `protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md:5`.
>
> **UT Coverage**
> - Total represented from turns 149–170: 66 UTs, all from turns 155–170.
> - TE-havib follow-on: `UT-160.b`, `UT-160.c`, `UT-161.a`, `UT-162.a`, `UT-162.b` at `protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md:59`–`protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md:63`.
> - TE-40 / TODO-kugod: `UT-155.a/b`, `UT-156.a/b/c`, `UT-157.a/b/c`, `UT-158.b/c/d/e/f/g`, `UT-159.a/b/c/d` at `protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md:126`–`protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md:143`.
> - TE-41 / TODO-turog: `UT-164.c/d/e`, `UT-165.b`, `UT-166.a/b/c/e`, `UT-167.a/b/c/e`, `UT-168.b/c/d` at `protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md:149`–`protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md:163`.
> - TE-42 / TODO-duvuk: `UT-168.a/e/f`, `UT-169.a/b/d/e` at `protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md:169`–`protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md:175`.
> - Other represented UTs: TE-sihih `UT-170.a` at `protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md:69`; Spec-edit UTs at `protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md:215`–`protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md:219`; Retire `UT-170.b/c/d` at `protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md:223`–`protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md:225`; Carry UTs from this range at `protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md:229`–`protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md:240`.
>
> **Owner Status**
> - Closure mechanism: the UT verification matrix is the closure index; no bulk checkbox flips are authorized: `protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md:353`–`protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md:360`.
> - Complete: TE-havib follow-on / TODO-lilok is closed verified-superseded: `protocols/wire-lab.d/TODO/TODO-lilok-te-36-followon-oq-and-tabletop-walk.md:10`–`protocols/wire-lab.d/TODO/TODO-lilok-te-36-followon-oq-and-tabletop-walk.md:16`.
> - Open: TODO-lilar remains open while downstream UTs remain unresolved: `protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md:12`.
> - Open: TODO-kugod is now reopened for residual TE-40 work, while the promise-stack retirement cascade is complete: `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md:10`–`protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md:16`.
> - Open: TODO-turog is blocked on TE-40 plus TODO-pipus/TE-43 migration design, with freeze-doc cleanup parked: `protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md:10`–`protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md:15`.
> - Open: TODO-duvuk depends on TE-41 and has no DI/product yet: `protocols/wire-lab.d/TODO/TODO-duvuk-te-42-filename-cid-cascade-policy.md:10`–`protocols/wire-lab.d/TODO/TODO-duvuk-te-42-filename-cid-cascade-policy.md:40`.
>
> **Blockers / Inconsistencies**
> - TE-40/TODO-kugod: current on-disk status is open via `DI-somuj`, but this is in uncommitted working-tree changes; the previous committed wording said closed. Evidence for the new intent is `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md:53`–`protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md:60`.
> - TE-41/TODO-turog: step 5 is explicitly unsafe; TE-nijab says transport message directories are append-only specimen data and freeze-doc cleanup is parked behind TODO-pipus/TE-43: `protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md:37`–`protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md:41` and `DR/DR-suhod-mihip-merge-blockers-partial-fix.md:306`–`DR/DR-suhod-mihip-merge-blockers-partial-fix.md:326`.
> - TE-42/TODO-duvuk: blocked behind TE-41; it still must resolve Path-A-vs-Path-B legacy rehash, Message-ID prohibition vs deprecation, strict-reader interactions, and CID-filename stability: `protocols/wire-lab.d/TODO/TODO-duvuk-te-42-filename-cid-cascade-policy.md:17`–`protocols/wire-lab.d/TODO/TODO-duvuk-te-42-filename-cid-cascade-policy.md:27`.
> - Stale status risk: TE-sihih itself is `decided` at `docs/thought-experiments/TE-sihih-substrate-agnostic-layered-model.md:16`–`docs/thought-experiments/TE-sihih-substrate-agnostic-layered-model.md:18`, but TODO-vunub still says “In progress” at `protocols/wire-lab.d/TODO/TODO-vunub-te-38-substrate-agnostic-layered-model.md:15`–`protocols/wire-lab.d/TODO/TODO-vunub-te-38-substrate-agnostic-layered-model.md:18`.
> - Session-replay cleanup: TODO-topit/transcript-snapshot work is deferred indefinitely and agents must not write session-log artifacts unless reactivated: `protocols/wire-lab.d/TODO/TODO-topit-transcript-snapshot-procedure.md:39`–`protocols/wire-lab.d/TODO/TODO-topit-transcript-snapshot-procedure.md:46` and `protocols/wire-lab.d/TODO/TODO-topit-transcript-snapshot-procedure.md:150`–`protocols/wire-lab.d/TODO/TODO-topit-transcript-snapshot-procedure.md:160`.

## Turns 171-208 expansion snapshot

Added 2026-05-09 under `DI-zufar`. This section expands the live
monitor from the historical 149-170 snapshot to the explicit bounded
149-208 scope. It is a monitor snapshot only: it records where evidence
currently lives and what still blocks closure, without editing
TODO-lilar, the disposition memo, the verification matrix, specs,
transports, TODO-kugod, or DR-nugog.

### Source coverage

- Turns 171-192 are inside TODO-lilar's completed chronological walk.
  TODO-lilar records `021.171` through `021.192` as walked, then records
  `021.192-end` as the end of TODO-lilar's original in-scope walk.
- TODO-lilar explicitly says turns 193 and beyond are out of its 149-192
  mandate and require a successor artifact if Steve wants them walked.
  That out-of-scope statement is the reason TODO-jivam now monitors
  193-208 as boundary evidence instead of assuming TODO-lilar already
  owns the later turns.
- Session log files exist for the boundary checks inspected here:
  `/home/stevegt/lab/session-logs/sessions/ea135ce8/171-turn.md`,
  `/home/stevegt/lab/session-logs/sessions/ea135ce8/192-turn.md`,
  `/home/stevegt/lab/session-logs/sessions/ea135ce8/195-turn.md`,
  `/home/stevegt/lab/session-logs/sessions/ea135ce8/196-turn.md`,
  `/home/stevegt/lab/session-logs/sessions/ea135ce8/197-turn.md`, and
  `/home/stevegt/lab/session-logs/sessions/ea135ce8/208-turn.md`.

### Existing owner evidence for turns 171-192

- The disposition memo covers 186 UT entries spanning turns 155-192, so
  turns 171-192 are represented there even though turns 193-208 are not.
- The disposition inventory shows turn-171-through-192 entries across
  TE-sihih, TE-43, TE-45, Spec-edit, Retire, and Carry clusters, in
  addition to the earlier TE-40/41/42 clusters captured by the
  historical 149-170 report.
- The verification matrix remains the closure index of record. Its
  DF-V.3 lock says no bulk checkbox flips occur in TODO-lilar or the
  disposition memo; future cluster-owner TODOs close their own UTs when
  substantive work lands under their own DIs.
- TODO-lilar remains open because the completed chronological walk is
  not equivalent to closed recovery: unresolved UTs stay owned by
  downstream TEs/TODOs until resolved, retired, or explicitly
  transferred.

### Boundary turns requiring final review

- **Turn 195:** Steve corrected the directory-name/context-loss problem
  and instructed a detailed review from TE-35 through the first
  promisebase mention. This is the direct context-loss boundary that
  motivated the replay work.
- **Turn 196:** Steve requested a one-chunk-at-a-time review of the last
  72 hours with topic status and related TE numbers. The answer
  identifies verbal rules not yet committed, including the
  `<slug>-draft` rename, `grid envelope` vocabulary, the transports to
  groups direction, first-class sites, CBOR/CIDv1 message format,
  wire-lab-as-canon versus promisebase-as-prototype, and dogfooding
  pressure.
- **Turn 197:** Steve gave the concrete recovery procedure: create a TODO
  file with one line item per turn starting at 149, walk chronologically,
  write outstanding questions/answers/unfinished threads to the repo
  before changing threads, and check off line items. The answer created
  the predecessor TODO-lilar artifact with a 149-192 scope.
- **Turn 208:** Steve corrected the storage location with "keep the
  ledger in TODO 021" and then requested promise-stack provenance plus
  nested-vs-stacked research. The answer says `UT-PSTK-origin` was
  logged in the ledger and that the research subagent was running. This
  makes turn 208 a recovery/session-log boundary point even though it is
  outside TODO-lilar's original 149-192 chronological walk.

### Open monitor questions

- Whether turns 193-208 need their own successor walk remains unresolved
  for TODO-jivam closure. The closure rule above requires each of those
  turns to be explicitly accounted as no-op boundary material or filed in
  an owner artifact with DI/DR provenance.
- `UT-PSTK-origin` is reported in the turn-208 answer as logged in TODO
  021, but the current disposition/matrix artifacts cited by the
  historical monitor are centered on turns 155-192. A final pass must
  verify whether that post-192 entry is already represented in the
  appropriate owner artifact or needs a new pointer.
- The turn-195/196 correction sequence and turn-197 recovery-procedure
  sequence should be reviewed together during final closure, because
  they define the handoff from session-log evidence to TODO-lilar rather
  than a normal in-scope UT cluster.
