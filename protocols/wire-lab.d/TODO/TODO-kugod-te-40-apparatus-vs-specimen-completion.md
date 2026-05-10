# TODO-kugod: TE-40 apparatus-vs-specimen completion + TE-famar closure

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TODO-25` (integer alias)
- `TODO-20260507-002306` (timestamp alias and pre-migration filename)

## Status

Open. Promise-stack retirement cascade complete; residual TE-40 recovery
still open. Cat-3 cascade work from TE-havib promise-stack retirement
has been applied to TE-famar, TE-muvuv, and TE-robub. TODO-rivuk and
DR-006 now point readers at the TE-havib DF-36.2 retirement instead of
inviting stale promise-stack DF answers.

## Threads absorbed from OPEN-THREADS.md

### T-PROMSTACK-RETIRE-CASCADE (formerly OPEN-THREADS, opened 2026-05-06)

Cat-3 cascade from TE-havib promise-stack retirement. DF-36.2 Alt-2.A
revised: promise-stack retired as a separate hypothesis; payload-
recursion under per-protocol specs is the answer per TE-lozip § 3.1 +
framing essay § 3.1.

Scope: refine TE-famar, TE-muvuv, TE-robub with Cat-3 entries recasting
promise-stack vocabulary into payload-recursion vocabulary; verify no
other TE in the corpus references promise-stack as a separate layer;
confirm protocols/promise-stack.d/ remains absent from the tree.

Blocking: nothing today; pure refinement work.

Anchor: TE-havib § DF-36.2 Alt-2.A (revised); TE-lozip; framing essay § 3.1.
Disposition-file pointer: `dropped-thread-disposition-20260506.md`
§ TE-40 cluster (18 UTs).

## Question log

(Per AGENTS-ppx Question-logging discipline. No questions logged yet.)

## Decision Intent Log

ID: DI-runuh
Date: 2026-05-08 23:41:08
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Execute TODO-kugod as a narrow Cat-3/Cat-4 retirement cascade. Append forward pointers to TE-famar, TE-muvuv, and TE-robub; close TODO-rivuk and DR-006 as superseded/no-longer-applicable; do not create `protocols/promise-stack.d/`; do not sweep `protocols/wire-lab.d/specs/harness-spec-draft.md`; do not rewrite TE historical body text or TE status fields in this pass.
Intent: Preserve the historical promise-stack analyses while making the current corpus navigable after TE-havib DF-36.2 retired promise-stack as a separate hypothesis. Readers should follow payload-recursion under `grid <pcid> <payload>` via TE-havib, TE-lozip, and the congruence/convergence essay instead of answering stale TE-famar DF-1.* questions.
Constraints: This is a navigational cleanup only. Cat-3 refinements append to `## Refinements`; historical TE bodies remain untouched under TE-dabol/TE-vudaf Cat-1b/Cat-3 rules. DR-006 and TODO-rivuk may be updated to stop stale work queues, but no promise-stack ordering DIs are created.
Affects: `docs/thought-experiments/TE-famar-promise-stack-ordering.md`; `docs/thought-experiments/TE-muvuv-promise-stack-as-zero-knowledge-envelope.md`; `docs/thought-experiments/TE-robub-time-traveling-break-witness.md`; `protocols/wire-lab.d/TODO/TODO-rivuk-te-promise-stack-ordering.md`; `DR/DR-006-20260429-164729-promise-stack-ordering.md`; `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`.

ID: DI-somuj
Date: 2026-05-09 17:43:30
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Reopen TODO-kugod as the TE-40 owner artifact while preserving the completed promise-stack retirement cascade.
Intent: The cascade completed under DI-runuh, but the broader TE-40 recovery cluster still has residual apparatus-vs-specimen work. The TODO status must not imply full closure while unresolved UT-* items remain assigned to TE-40.
Constraints: Do not reopen TODO-rivuk or DR-006 as promise-stack-ordering work; those remain superseded by TE-havib DF-36.2. Keep the completed cascade recorded as complete, and use TODO-kugod for residual TE-40 recovery visibility.
Affects: `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`; `protocols/wire-lab.d/TODO/TODO.md`; `protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`; `protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md`.

ID: DI-vopim
Date: 2026-05-09 11:21:29
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Track the remaining TE-40 recovery as an explicit residual checklist in TODO-kugod, and file DR-nugog for the unresolved transport-tree ownership/structure question instead of changing specs or moving transport data in this pass.
Intent: TODO-kugod needs to remain open with a precise map from each TE-40 UT to resolved, retired, transferred, or still-open ownership. The transport-tree question crosses TODO-kugod's outer apparatus cleanup and TODO-turog's group-session freeze work, so it needs a DR before any spec wording or transport path is changed.
Constraints: Do not modify `protocols/wire-lab.d/specs/transport-spec-draft.md`, `protocols/group-session.d/specs/group-session-draft.md`, or `transports/wire-lab-devs-draft/` in this pass. The checklist is a coordination artifact, not a behavior change. DR-nugog asks the question; it does not decide whether the tree becomes `transports/<protocol-slug>/<instance-dir>/` or stays flat.
Affects: `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`; `DR/DR-nugog-transport-tree-ownership-structure.md`; `protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md`; future transport-spec/group-session-spec cleanup.

ID: DI-pakid
Date: 2026-05-10 01:27:27
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Lock DF-mupoz.3 as Alt 3.A: root `protocols/` contains only `wire-lab.d`; candidate PromiseGrid protocols move under simulations as specimens.
Intent: Keep root `protocols/` from mixing wire-lab apparatus with candidate PromiseGrid protocols under test. `protocols/wire-lab.d/` remains the harness apparatus home, while `group-session`, `udp-binding`, `ppx-dr`, and similar candidate protocols are tested inside named simulations until a later graduation path is decided.
Constraints: This locks only DF-mupoz.3. It does not yet lock the first simulation path, the physical migration timing, proposal-record treatment, transport specimen move, or the final graduation destination for candidate protocol specs.
Affects: `docs/thought-experiments/TE-mupoz-root-protocol-migration-scope-under-simulations.md`; `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`; `DR/DR-nugog-transport-tree-ownership-structure.md`; `protocols/wire-lab.d/specs/harness-spec-draft.md`; `DEV-GUIDE-RESOURCES.md`; future simulation migration TODO.

ID: DI-fakin
Date: 2026-05-10 03:56:54
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Lock the remaining TE-mupoz decisions and implementation paths: DF-mupoz.1 Alt 1.A, DF-mupoz.2 Alt 2.A, DF-mupoz.4 Alt 4.C, and DF-mupoz.5 Alt 5.A. Implement the first recovery/dogfood simulation at `simulations/SIM-piloh-turns-149-208-recovery/`; move candidate protocol trees with their TODO queues into the simulation; move `transports/wire-lab-devs-draft/` into the simulation world; archive root `transports/README.md` under `archive/transports/`; move all tracked `proposals/` records under `archive/proposals/`; leave no tracked root `transports/` or `proposals/` path; update `DEV-GUIDE-RESOURCES.md` with detailed relocation guidance.
Intent: Complete the TE-mupoz apparatus/specimen split without preserving obsolete root paths as compatibility commitments. The first simulation becomes the bounded place where turns 149-208 recovery, candidate PromiseGrid protocol specimens, legacy proposal records, and concrete transport evidence can be replayed and evaluated, while rooted wire-lab apparatus remains focused on explaining and governing the experiment.
Constraints: Do not edit `/home/stevegt/lab/promisegrid-dev-guide`. Preserve history with `git mv` for tracked moves and with migration manifests that record old paths, new paths, source commit, CID status where applicable, and graduation limits. Simulation results must feed DR/DI/spec/dev-guide handoff rather than directly becoming authoritative PromiseGrid layout.
Affects: `simulations/SIM-piloh-turns-149-208-recovery/`; `protocols/group-session.d/`; `protocols/ppx-dr.d/`; `protocols/udp-binding.d/`; `transports/`; `proposals/`; `DEV-GUIDE-RESOURCES.md`; `docs/thought-experiments/TE-mupoz-root-protocol-migration-scope-under-simulations.md`; `docs/thought-experiments/README.md`; `protocols/wire-lab.d/specs/harness-spec-draft.md`; `DR/DR-nugog-transport-tree-ownership-structure.md`; `protocols/wire-lab.d/TODO/TODO.md`; `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`.

ID: DI-mahim
Date: 2026-05-10 05:49:06
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Treat TE-mupoz and `DI-fakin` as resolving TODO-kugod's current transport-tree and simulation-first blockers for the turns 149-208 recovery specimen, without globally deciding TE-domat, TE-pahah, or TE-nizor. Continue TODO-kugod with a narrower residual path: harness-spec apparatus/specimen cleanup, grid-envelope successor ownership, and transport-spec companion audit.
Intent: Keep the 149-208 replay moving after the simulation migration instead of re-litigating the root-level `transports/` and candidate-protocol placement questions that Mupoz already answered for the current recovery specimen.
Constraints: Do not mark TE-domat, TE-pahah, or TE-nizor globally decided. Do not change harness-spec wording, create grid-envelope files, or run the transport-spec companion audit in this reconciliation pass. Preserve `simulations/SIM-piloh-turns-149-208-recovery/` as the current specimen home under `DI-fakin`.
Affects: `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`; `protocols/wire-lab.d/TODO/TODO.md`; `docs/thought-experiments/TE-domat-transports-groups-reconciliation.md`; `docs/thought-experiments/TE-pahah-wire-lab-simulation-first-structure.md`; `docs/thought-experiments/TE-nizor-pahah-implementation-sufficiency.md`; `docs/thought-experiments/TE-mupoz-root-protocol-migration-scope-under-simulations.md`; `DR/DR-nugog-transport-tree-ownership-structure.md`.

## Residual TE-40 checklist

This checklist maps the TE-40 UT inventory from
`protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md` to
current ownership. TODO-kugod stays open until every open or transferred
row below is closed by a DI, DR, or successor TODO.

| UT | Status | Owner / next artifact | Disposition |
| --- | --- | --- | --- |
| UT-155.a | retired | TODO-rivuk / DI-runuh | TE-famar DF-1.1 is no longer answered as live promise-stack ordering work. |
| UT-155.b | retired | TODO-rivuk / DI-runuh | `Project` / `Peel` / `Wrap` remain historical promise-stack vocabulary, not current apparatus work. |
| UT-156.a | retired | TODO-rivuk / DI-runuh | The abandoned TE-famar structural-role question is superseded by promise-stack retirement. |
| UT-156.b | resolved-retired | TODO-kugod / DI-runuh | TE-famar stays in the historical TE corpus with Cat-3 forward pointers instead of moving. |
| UT-156.c | open | TODO-kugod | The bad "harness-spec is wire-envelope-agnostic" claim still needs harness-spec cleanup in a later scoped pass. |
| UT-157.a | transferred-open | TODO-kugod until grid-envelope successor exists | Candidate envelope inventory belongs with the future grid-envelope/specimen home, not with promise-stack. |
| UT-157.b | retired | TODO-rivuk / DI-runuh | The abandoned TE-famar status-reading DF is no longer live. |
| UT-157.c | transferred-open | TODO-kugod until grid-envelope successor exists | The `grid([pcid, payload])` working-hypothesis prose needs a grid-envelope home or successor TODO. |
| UT-158.b | resolved | TE-havib DF-36.1 | The apparatus-vs-specimen scope is strict carve-out; no new TE-40 scope choice remains. |
| UT-158.c | open | TODO-kugod | Harness-spec §1.1 carve-out remains unexecuted. |
| UT-158.d | open | TODO-kugod | Harness-spec §1.3 apparatus/specimen classification remains unexecuted. |
| UT-158.e | retired | TODO-kugod / DI-runuh | TE-famar is not moved into a promise-stack protocol directory because promise-stack is retired as a separate hypothesis. |
| UT-158.f | transferred-open | TODO-kugod until grid-envelope successor exists | The grid-envelope protocol directory/spec work remains outside this pass and needs a successor owner. |
| UT-158.g | retired | TODO-rivuk / DI-runuh | TODO-rivuk is closed as superseded instead of moved under a promise-stack protocol directory. |
| UT-159.a | open | TODO-kugod | The nine specimen-bearing harness-spec audit items still need a later sweep or explicit retirement. |
| UT-159.b | open-scoped | TODO-kugod plus DR-nugog | The current transport-tree blocker is resolved for this specimen by `DI-fakin`; the remaining work is the transport-spec companion audit against rooted apparatus versus simulation-local specimens. |
| UT-159.c | resolved-retired | TE-havib follow-on verification | The six-scenario mismatch is recorded; no redo is required before residual checklist cleanup proceeds. |
| UT-159.d | open | TODO-kugod | The remaining ambiguous audit areas still need resolution or explicit retirement. |

### 2026-05-09/10 transport-tree and post-Mupoz tracking

- [x] kugod.1 Resolve the current TE-domat (`docs/thought-experiments/TE-domat-transports-groups-reconciliation.md`) and DR-nugog blocker before closing UT-159.b. Closed for the current recovery specimen by `DI-fakin` and `DI-mahim`; TE-domat remains useful background for future root-level transport or group-tree decisions.
- [x] kugod.2 Resolve the current TE-pahah (`docs/thought-experiments/TE-pahah-wire-lab-simulation-first-structure.md`) blocker before implementing any root-level `transports/` or `groups/` migration. Closed for the current recovery specimen by `DI-fakin` and `DI-mahim`; the candidate artifacts moved under `simulations/SIM-piloh-turns-149-208-recovery/` instead of becoming root-level layer trees.
- [x] kugod.3 Resolve the current TE-nizor (`docs/thought-experiments/TE-nizor-pahah-implementation-sufficiency.md`) blocker before closing residual recovery for turns 149-208. Closed for the current recovery specimen by `DI-fakin` and `DI-mahim`; TE-nizor remains open as broader analysis for root-level layer trees, artifact-message shape, CBOR/text representation, and CAS/feed details.
- [x] kugod.4 Resolve TE-mupoz (`docs/thought-experiments/TE-mupoz-root-protocol-migration-scope-under-simulations.md`) before moving existing root-level or `protocols/` content under `simulations/`. TE-mupoz tests whether root `protocols/` should retain only `wire-lab.d`, and whether candidate protocol drafts, proposal records, and transport specimens should move into the first simulation as specimens when TE-nizor is implemented. Closed by `DI-pakid` and `DI-fakin`; implementation target is `simulations/SIM-piloh-turns-149-208-recovery/`.
- [ ] kugod.5 Inventory the harness-spec claims that still mix harness apparatus with candidate PromiseGrid specimens, including the nine UT-159.a audit items and the ambiguous UT-159.d areas.
- [ ] kugod.6 Amend the harness spec so §1.1 and §1.3 state the apparatus/specimen boundary cleanly and remove or supersede the stale "wire-envelope-agnostic" claim from UT-156.c, UT-158.c, and UT-158.d.
- [ ] kugod.7 Create or identify the successor owner for grid-envelope work, candidate envelope inventory, and the `grid([pcid, payload])` working hypothesis from UT-157.a, UT-157.c, and UT-158.f.
- [ ] kugod.8 Run the transport-spec companion audit after the Mupoz migration, distinguishing rooted wire-lab apparatus claims from simulation-local transport specimens for UT-159.b.
- [ ] kugod.9 Update TODO-jivam and the TE-40 disposition matrix after kugod.5 through kugod.8 are resolved, transferred, or explicitly retired.
