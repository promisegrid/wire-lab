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

ID: DI-lajod
Date: 2026-05-10 13:56:58
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Execute turn-156 cleanup as a narrow apparatus/specimen sweep in `protocols/wire-lab.d/specs/harness-spec-draft.md`: replace the remaining harness-level promise-stack framing in §1.1 and the specimen-specific wording in §1.3 with apparatus-level language, and close `UT-156.c` without claiming the broader TE-40 harness-spec audit is complete.
Intent: Retract the stale turn-156 claim that the harness-spec is "wire-envelope-agnostic" by stating the correct relationship: the harness is apparatus that compares candidate envelopes and adjacent layer choices. This resolves the turn-156 slice while keeping the larger §1.1 / §1.3 carve-out and the wider nine-item audit visible as separate TE-40 work.
Constraints: Do not reopen TE-havib DFs, create new specimen directories, or pretend §1.2, §3.3, §7.1, or the nine `UT-159.a` audit items are finished. Use the current simulation-local specimen surface from `DI-fakin` for forward pointers. Leave `TODO-lilar`, the disposition memo, and the UT verification matrix untouched unless a factual error is found.
Affects: `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`; `protocols/wire-lab.d/specs/harness-spec-draft.md`; `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`.

ID: DI-sujan
Date: 2026-05-10 14:04:03
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Tighten the turn-156 harness-spec cleanup so apparatus-level prose points generically at named simulations under `simulations/` rather than treating `simulations/SIM-piloh-turns-149-208-recovery/` as the harness-wide specimen home. Keep `SIM-piloh` only as a scoped example where the current recovery/dogfood specimen set is relevant.
Intent: Prevent the apparatus spec from silently promoting one current simulation into the preferred or canonical specimen surface. The harness should describe how named simulations provide specimen sets, not imply that one currently-populated simulation is the normative home for all candidate specimens.
Constraints: Do not rewrite historical TE summary sections that name `SIM-piloh` as the current location of specific specimen evidence. Limit this pass to apparatus-level overreach in `harness-spec-draft.md` and the narrow owner/provenance wording that tracks turn-156 cleanup.
Affects: `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`; `protocols/wire-lab.d/specs/harness-spec-draft.md`.
Supersedes: DI-lajod (simulation-specific forward-pointer wording only)

ID: DI-sotuk
Date: 2026-05-12 08:57:44
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Resolve the turn-158 replay bookkeeping gap by decomposing `UT-158.a` into the existing TE-40 owner rows and routing `UT-158.h` to the grid-envelope successor owner. Keep `UT-158.c`, `UT-158.d`, and `UT-158.f` open under their current owners.
Intent: Turn 158 introduced the apparatus-vs-specimen correction and a six-step cleanup sequence. Later work already split that sequence across concrete artifacts, but the replay ledger still showed two carry items without explicit disposition. This decision closes the bookkeeping gap without pretending the remaining apparatus/specimen and grid-envelope protocol/spec work is complete.
Constraints: Do not edit harness specs, TE files, raw logs, `TODO-lilar`, or create new protocol/simulation directories in this pass. Do not make `grid([pcid, payload])` canonical. Do not close `kugod.6` or `tujad.3`.
Affects: `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`; `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`; `simulations/SIM-kurim-grid-envelope/TODO/TODO-tujad-grid-envelope-successor-owner.md`.

ID: DI-kinad
Date: 2026-05-12 09:23:15
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Close the remaining turn-158 rows `UT-158.c`, `UT-158.d`, and
`UT-158.f` after the apparatus-level §1.1/§1.3 wording was already rewritten
under `DI-lajod` and `DI-sujan`, and after the grid-envelope successor path was
materialized as 24 standalone positional variant simulations under `DI-fanah`.
Intent: The replay should not keep turn-158-local residue open after the work
has been either landed or transferred into concrete successor artifacts. Broader
TE-40 audit work remains visible under the later turn-159 rows, but turn 158 no
longer needs to carry its own loose-end queue.
Constraints: Close only turn-158 residue. Do not close `UT-159.a`,
`UT-159.b`, `UT-159.d`, `kugod.5`, `kugod.8`, `kugod.9`, or the broader TE-40
backlog. Do not declare a winning grid-envelope variant.
Affects: `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`;
`protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`simulations/SIM-kurim-grid-envelope/TODO/TODO-tujad-grid-envelope-successor-owner.md`;
`simulations/SIM-*-grid-envelope-enc-<cbor|dag-cbor>-unknown-<opaque|hard-reject|best-effort>-sig-<wrapper-pcid|unsigned-v0|mandatory-opaque-bytes|mandatory-sig-pcid-payload>/`.

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
| UT-156.c | resolved | TODO-kugod / DI-lajod + DI-sujan + DI-kinad | The stale "wire-envelope-agnostic" claim is superseded by the apparatus-level §1.1 / §1.3 rewrite in `harness-spec-draft.md`, with simulation references tightened so the harness does not privilege `SIM-piloh`; later turn-158 local residue is now closed, while broader TE-40 audit work continues under the turn-159 rows. |
| UT-157.a | resolved | `simulations/SIM-kurim-grid-envelope/TODO/TODO-tujad-grid-envelope-successor-owner.md` / DI-joroh | Candidate envelope inventory is captured in the grid-envelope successor owner, not promise-stack. |
| UT-157.b | retired | TODO-rivuk / DI-runuh | The abandoned TE-famar status-reading DF is no longer live. |
| UT-157.c | resolved | `simulations/SIM-kurim-grid-envelope/TODO/TODO-tujad-grid-envelope-successor-owner.md` / DI-joroh | The `grid([pcid, payload])` working-hypothesis prose is captured in the concrete grid-envelope successor TODO. |
| UT-158.a | resolved-decomposed | TODO-kugod / DI-sotuk + DI-kinad | The raw six-step sequence is no longer a standalone work item: the audit memo and TE-havib captured the apparatus/specimen split, promise-stack relocation was retired under `DI-runuh`, grid-envelope ownership is routed through `TODO-tujad`, and the remaining turn-158 rows are now closed by `DI-kinad` and `DI-fanah`. |
| UT-158.b | resolved | TE-havib DF-36.1 | The apparatus-vs-specimen scope is strict carve-out; no new TE-40 scope choice remains. |
| UT-158.c | resolved | TODO-kugod / DI-lajod + DI-sujan + DI-kinad | Harness-spec §1.1 has the apparatus-level rewrite and neutral simulation wording required for the turn-158 slice; broader nine-item audit work remains under `UT-159.a` and `kugod.5`. |
| UT-158.d | resolved | TODO-kugod / DI-lajod + DI-sujan + DI-kinad | Harness-spec §1.3 has the apparatus-level layering-test rewrite required for the turn-158 slice; broader classification and specimen-audit work remains under `UT-159.a`, `UT-159.d`, `kugod.5`, and `kugod.8`. |
| UT-158.e | retired | TODO-kugod / DI-runuh | TE-famar is not moved into a promise-stack protocol directory because promise-stack is retired as a separate hypothesis. |
| UT-158.f | resolved-transferred | `simulations/SIM-kurim-grid-envelope/TODO/TODO-tujad-grid-envelope-successor-owner.md` / DI-fanah + DI-kinad | The grid-envelope protocol/spec successor path is materialized as 24 standalone positional variant simulations; no single preferred specimen is selected. |
| UT-158.g | retired | TODO-rivuk / DI-runuh | TODO-rivuk is closed as superseded instead of moved under a promise-stack protocol directory. |
| UT-158.h | resolved-routed | `simulations/SIM-kurim-grid-envelope/TODO/TODO-tujad-grid-envelope-successor-owner.md` / DI-sotuk + DI-joroh + DI-fanah | The requested parallel TODO for the grid hypothesis is satisfied by the grid-envelope successor owner; the concrete protocol/spec successor path is now closed for turn-158 scope by the 24 variant simulations. |
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
- [x] kugod.6 Finish the harness-spec §1.1 / §1.3 apparatus/specimen cleanup for `UT-158.c` and `UT-158.d`; `UT-156.c`'s stale turn-156 phrasing is already superseded by `DI-lajod`. Closed for turn-158 scope by `DI-kinad`; broader turn-159 audit work remains under `kugod.5` and `kugod.8`.
- [x] kugod.7 Create or identify the successor owner for grid-envelope work, candidate envelope inventory, and the `grid([pcid, payload])` working hypothesis from UT-157.a, UT-157.c, and UT-158.f. Closed by routing to `simulations/SIM-kurim-grid-envelope/TODO/TODO-tujad-grid-envelope-successor-owner.md` in `DI-mosor`.
- [ ] kugod.8 Run the transport-spec companion audit after the Mupoz migration, distinguishing rooted wire-lab apparatus claims from simulation-local transport specimens for UT-159.b.
- [ ] kugod.9 Update TODO-jivam and the TE-40 disposition matrix after kugod.5 through kugod.8 are resolved, transferred, or explicitly retired.
