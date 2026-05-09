# DR-nugog - Transport-tree ownership and structure

DR-ID: DR-nugog
Date: 2026-05-09 11:21:29
Asked by: stevegt@t7a.org (Steve Traugott)
State: open
Question: Should the outer `transports/` tree group transport instances under a protocol slug, as `transports/<protocol-slug>/<instance-dir>/`, or keep the current flat/spec-worded shape where each transport instance is directly under `transports/`; and should the cleanup owner be TODO-kugod, TODO-turog, or split between them?
Why this blocks progress: The current drafts and live data disagree with the turn-170 design pressure. The flat tree is easy to leave untouched, but it gives no structural grouping for multiple group-session instances or future transport protocols. The nested tree would clarify ownership and navigation, but changing it affects both the outer apparatus/specimen cleanup and the group-session freeze/migration procedure. This DR records the unresolved question before any spec text or transport path is changed.
Affects: `protocols/wire-lab.d/specs/transport-spec-draft.md`; `protocols/group-session.d/specs/group-session-draft.md`; `transports/wire-lab-devs-draft/`; `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`; `protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md`.
Unblocks: TODO-kugod UT-159.b; TODO-turog freeze-path cleanup; any future edit that reconciles flat `transports/<instance-dir>/` wording with a protocol-grouped tree.
Waiting on: stevegt@t7a.org (Steve Traugott)
Decision:
Linked DI: DI-vopim
Related commits:
Last updated: 2026-05-09 18:21:29 UTC

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
