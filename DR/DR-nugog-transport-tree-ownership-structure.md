# DR-nugog - Transport-tree ownership and structure

DR-ID: DR-nugog
Date: 2026-05-09 11:21:29
Asked by: stevegt@t7a.org (Steve Traugott)
State: implemented
Question: Should the outer `transports/` tree group transport instances under a protocol slug, as `transports/<protocol-slug>/<instance-dir>/`, or keep the current flat/spec-worded shape where each transport instance is directly under `transports/`; and should the cleanup owner be TODO-kugod, TODO-turog, or split between them?
Why this blocks progress: The current drafts and live data disagree with the turn-170 design pressure. The flat tree is easy to leave untouched, but it gives no structural grouping for multiple group-session instances or future transport protocols. The nested tree would clarify ownership and navigation, but changing it affects both the outer apparatus/specimen cleanup and the group-session freeze/migration procedure. This DR records the unresolved question before any spec text or transport path is changed.
Affects: `protocols/wire-lab.d/specs/transport-spec-draft.md`; `simulations/SIM-piloh-turns-149-208-recovery/protocols/group-session.d/specs/group-session-draft.md`; `simulations/SIM-piloh-turns-149-208-recovery/world/transports/wire-lab-devs-draft/`; `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`; `protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md`.
Unblocks: TODO-kugod UT-159.b; TODO-turog freeze-path cleanup; any future edit that reconciles flat `transports/<instance-dir>/` wording with a protocol-grouped tree.
Waiting on: stevegt@t7a.org (Steve Traugott)
Decision: Superseded for the current specimen by `DI-fakin`: do not choose flat-vs-nested root `transports/` for the current `wire-lab-devs-draft` evidence. Move the specimen into `simulations/SIM-piloh-turns-149-208-recovery/world/transports/wire-lab-devs-draft/`, preserve CID evidence, and let any future root group/feed/CAS/transport layout graduate through later DR/DI/spec work.
Linked DI: DI-vopim; DI-fakin
Related commits:
Last updated: 2026-05-10 04:42:05 UTC

## Evidence

- Turn 170 (`/home/stevegt/lab/session-logs/sessions/ea135ce8/170-turn.md`) raised the exact structure concern: if a second named group appears, a flat `transports/draft--wire-lab-devs/` style tree does not show which instance directories belong to the same transport protocol. The turn framed three alternatives: protocol-slug nesting, recursive draft/pCID nesting, or status quo deferral.
- `protocols/wire-lab.d/specs/transport-spec-draft.md` currently states the outer key as flat `transports/<pcid>--<slug>/`, with the pCID as canonical protocol identity and the slug as human-readable instance suffix.
- `protocols/group-session.d/specs/group-session-draft.md` currently states that a group-session transport directory has no subdirectories and that all messages live directly under `transports/<this-pcid>--<slug>/`.
- `transports/wire-lab-devs-draft/README.md` currently documents the live bootstrap instance as a flat directory and says freeze renames it to `transports/wire-lab-devs-<pcid>/`.
- `protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md` owns group-session freeze cleanup, but its earlier freeze-time rename step is already marked unsafe because transport message directories are append-only specimen data.
- `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md` owns the TE-40 residual audit, including UT-159.b's unresolved transport-spec companion audit.

## Candidate decisions

- **Alt-A: protocol-slug grouping.** Use `transports/<protocol-slug>/<instance-dir>/`, for example `transports/group-session/wire-lab-devs-draft/`. TODO-kugod owns outer apparatus/spec wording; TODO-turog owns group-session freeze details below the protocol-slug level.
- **Alt-B: flat current/spec wording.** Keep live/spec wording flat, with transport instances directly under `transports/`. TODO-turog owns freeze-path cleanup; TODO-kugod records UT-159.b as retired or no-op after explicit decision.
- **Alt-C: split/defer.** Keep live data flat for now, but record a future migration rule or compatibility note. TODO-kugod owns the DR/spec-audit question; TODO-turog owns only the group-session freeze procedure once this DR is answered.

## Notes

This DR intentionally does not move `transports/wire-lab-devs-draft/`, rewrite specs, or choose an alternative. It exists so the next edit can cite a concrete decision record instead of relitigating turn 170 from session logs.

## 2026-05-09 update - TE-domat reframes the question

TE-domat (`docs/thought-experiments/TE-domat-transports-groups-reconciliation.md`) shows that this DR's original alternatives are too narrow. The original question asks only whether instances should nest under `transports/<protocol-slug>/` or remain flat. That omits turn 176's explicit `groups` choice and the later TE-nijab decision that `transports/` is lower-layer feed/transport specimen data.

Updated reading: this DR should be answered after DF-domat, or amended to ask the split-tree question directly:

- whether `groups/` should become the L7 group registry/view tree,
- whether `transports/` should remain lower-layer feed/transport simulation evidence,
- and how `transports/wire-lab-devs-draft/` is preserved as historical pre-layered specimen data until an additive migration exists.

TE-domat recommends the split-tree interpretation, not a naive whole-tree `transports/` -> `groups/` rename.

## 2026-05-09 update - TE-pahah reframes around simulation-first structure

Steve clarified that backward compatibility with the current experimental paths is not required because none of this is production or active use yet. TE-pahah (`docs/thought-experiments/TE-pahah-wire-lab-simulation-first-structure.md`) therefore asks a broader question: what structure best serves wire-lab's simulation and decision-making goals?

Updated reading after TE-pahah: this DR should not lock a root-level `transports/` vs. `groups/` answer until DF-pahah decides whether concrete worlds belong under a top-level `simulations/` tree. TE-pahah recommends `simulations/<sim>/world/{sites,groups,cas,feeds,wires}/` as the primary experiment home, with root-level `groups/` or `transports/` deferred unless a later dogfood/reference layout needs them.

## 2026-05-09 update - TE-vilot combines simulation structure and promise artifacts

TE-vilot (`docs/thought-experiments/TE-vilot-promise-shaped-simulation-artifacts.md`) strengthens the TE-pahah reading. It treats `simulations/` as the safe boundary for testing promise-shaped artifact protocols without forcing every apparatus document to become a PT promise.

Updated reading after TE-vilot: this DR should avoid deciding root-level transport/group layout or universal promise-artifact templates in isolation. If the current `transports/wire-lab-devs-draft/` evidence is migrated, the likely target is a named simulation world. If promise-shaped metadata is desired for that migration, it should be tested as a simulation artifact protocol or commitment-specific template, not as a blanket rewrite of TEs, DRs, TODOs, or guide-resource notes.

## 2026-05-09 update - TE-hirap adds artifact-as-message representation concerns

TE-hirap (`docs/thought-experiments/TE-hirap-artifacts-as-promisegrid-messages.md`) extends TE-vilot from promise-shaped prose to full PromiseGrid-message-shaped artifacts. It distinguishes plain-text `grid <pcid>` specimens, CBOR promise-stack candidates, dual-view identity hazards, and commitment/specimen-only graduation.

Updated reading after TE-hirap: if `transports/wire-lab-devs-draft/` evidence moves into a simulation, the migration should not silently convert every surrounding apparatus file into a PromiseGrid message. Any message-shaped artifact protocol should live inside the named simulation or a commitment/specimen class, state whether text or CBOR is canonical, and avoid two competing CIDs for the same logical artifact.

## 2026-05-09 update - TE-nizor tests whether Pahah is sufficient

TE-nizor (`docs/thought-experiments/TE-nizor-pahah-implementation-sufficiency.md`) examines TE-pahah, TE-vilot, TE-hirap, and turns 149-208 together. It concludes that Pahah can satisfy the recovered concerns only if implemented as `simulations/` plus a minimal simulation contract, not as a bare directory or a root-level `groups/`/`transports` migration.

Updated reading after TE-nizor: this DR should wait for DF-nizor.1 and DF-nizor.2 before any root-level tree migration. If those DFs lock the recommendation, the next transport/group action is likely to seed a named recovery/dogfood simulation from current `transports/wire-lab-devs-draft/` evidence, with any eventual root-level `groups/`, `feeds/`, `cas/`, or `sites/` paths treated as downstream results rather than prerequisites.

## 2026-05-10 update - TE-mupoz narrows physical migration scope

TE-mupoz (`docs/thought-experiments/TE-mupoz-root-protocol-migration-scope-under-simulations.md`) asks how much existing root-level and `protocols/` content should move if TE-nizor is implemented. It now recommends a wire-lab-only root protocols plus specimen-migration policy: keep `protocols/wire-lab.d/` rooted as harness apparatus, move candidate protocol trees such as `group-session`, `udp-binding`, and `ppx-dr` into the first simulation as specimens, move `transports/wire-lab-devs-draft/` into that simulation with source-path/source-commit/CID evidence, and write new concrete world state under `simulations/<sim>/world/`.

Updated reading after TE-mupoz: this DR should not request a broad physical move of apparatus docs under `simulations/`, but it also should not preserve the root `transports/` location or non-wire-lab root protocol trees as compatibility constraints. The live question narrows to the first simulation target path, the migration manifest requirements, and whether future group/transport/feed/CAS/protocol trees appear first inside simulation worlds before any root-level reference layout is adopted.

DF-mupoz.3 is locked by `DI-pakid`: root `protocols/` contains only `wire-lab.d`; candidate PromiseGrid protocols move under simulations as specimens.

## 2026-05-10 update - implemented by DI-fakin

`DI-fakin` implements the current-specimen part of this DR through TE-mupoz:
root `transports/` is no longer an active tracked specimen path, and
`transports/wire-lab-devs-draft/` moved to
`simulations/SIM-piloh-turns-149-208-recovery/world/transports/wire-lab-devs-draft/`
with CID verification and migration manifests. The old root
`transports/README.md` moved to the simulation archive as historical evidence.

This does not freeze a future PromiseGrid node layout. If simulation results
justify root-level `groups/`, `feeds`, `cas`, or transport reference trees, they
must graduate through a later DR, DI, frozen spec, guide prose, or external
PromiseGrid spec corpus.
