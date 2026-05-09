# TODO-jivam: Turns 149-170 recovery completion monitor

## Status

Open. This TODO shall not close until all recovery represented by the
turns 149-170 slice is complete: each source turn is accounted for, each
UT is resolved, retired, or explicitly transferred under its owner
artifact, and the owner artifacts that block recovery are closed or
superseded under DI/DR provenance.

## Decision Intent Log

ID: DI-jivam
Date: 2026-05-09 17:55:15
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Create a dedicated open TODO that preserves the full turns 149-170 recovery monitor report and stays open until that recovery is complete.
Intent: Keep the Sunday context-loss recovery visible as a first-class work item instead of relying on TODO-lilar alone or on transient agent reports. The monitor TODO is a closure guard: it prevents a completed replay walk or a completed sub-cascade from being mistaken for complete recovery while TE-40, TE-41, TE-42, and adjacent UT owners still carry unresolved recovery work.
Constraints: This TODO is a monitoring and closure-gate artifact. It does not close UTs directly, does not flip TODO-lilar checkboxes, and does not supersede the matrix-as-closure-index rule locked by `DI-021-20260507-210204`. Closure requires evidence that all turns 149-170 recovery items are represented and complete through their proper owner artifacts.
Affects: `protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`; `protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md`; `protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`; `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`; `protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md`; `protocols/wire-lab.d/TODO/TODO-duvuk-te-42-filename-cid-cascade-policy.md`; `protocols/wire-lab.d/TODO/TODO.md`.

## Closure rule

This TODO may close only after all of the following are true:

- [ ] jivam.1 Source coverage for turns 149-170 remains verified against `/home/stevegt/lab/session-logs/sessions/ea135ce8/149-turn.md` through `/home/stevegt/lab/session-logs/sessions/ea135ce8/170-turn.md`.
- [ ] jivam.2 Turns 149-154 remain verified as walked with no leftover threads, or any later contradiction is filed as a new UT/DR with provenance.
- [ ] jivam.3 All turn-155-through-170 UTs in `dropped-thread-disposition-20260506.md` have owner artifacts and closure evidence in the verification matrix.
- [ ] jivam.4 TODO-kugod closes or transfers all residual TE-40 recovery work after the apparatus-vs-specimen sweep, grid-envelope home, transport-spec companion audit, and stale TE-famar status/path issues are resolved or explicitly retired.
- [ ] jivam.5 TODO-turog closes or transfers all TE-41 group-session freeze work after its TE-40 and migration-design blockers are resolved.
- [ ] jivam.6 TODO-duvuk closes or transfers all TE-42 Message-ID / filename / CID-cascade policy work after TE-41 no longer blocks it.
- [ ] jivam.7 Adjacent turn-149-through-170 items outside TE-40/41/42 are resolved or explicitly transferred: TE-havib follow-on, TE-sihih `UT-170.a`, Spec-edit, Retire, and Carry entries.
- [ ] jivam.8 TODO-lilar remains open until every relevant downstream owner has completed, retired, or transferred its UTs under the matrix-as-closure-index rule.
- [ ] jivam.9 A final recovery monitor pass confirms no closed owner artifact still contains unfinished recovery work.

## Monitor report snapshot

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
