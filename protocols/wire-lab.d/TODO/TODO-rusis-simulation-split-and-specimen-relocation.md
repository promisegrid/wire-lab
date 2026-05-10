# TODO-rusis: Simulation split and specimen relocation

## Prior aliases

This TODO is being filed after the TE-mumuv proquint migration locked
2026-05-07, so it is minted directly under `TODO-rusis`. No prior
integer or timestamp alias.

## Status

Open. Tracks the post-Mupoz split of the mixed recovery simulation into
content-named sims, the retirement of active `ppx-dr`, the `udp-binding`
-> `udp-feed` rename, and the relocation of specimen-owned work out of
`protocols/wire-lab.d/`.

## Decision Intent Log

ID: DI-tugit
Date: 2026-05-10 23:42:35
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Dissolve `simulations/SIM-piloh-turns-149-208-recovery/` and replace it with separate content-named sims for `wire-lab-devs`, `group-session`, `feed-outer`, `udp-feed`, and `grid-envelope`. Retire `ppx-dr` as an active protocol tree, rename `udp-binding` now to `udp-feed`, move the specimen half of rooted `transport-spec-draft.md` into a new `feed-outer` sim-local protocol tree, and redistribute the remaining rooted apparatus residue into existing harness docs and owner TODOs rather than creating a new rooted replacement draft.
Intent: The current mixed `SIM-piloh` tree is process-named and conflates concrete world evidence, candidate protocol specimens, legacy proposal archive, and rooted apparatus residue. The repo should expose simulation content by what it contains, keep rooted wire-lab apparatus focused on harness governance, and make specimen ownership explicit enough that replay cleanup does not have to re-litigate the same placement questions.
Constraints: Dissolve the mixed `SIM-piloh` tree completely rather than leaving an umbrella successor. Keep `wire-lab-devs` as the only concrete world simulation in this pass. Retire active `ppx-dr` trees but preserve proposal history under archive/provenance. Update active paths and current-pointer prose to `udp-feed`, not `udp-binding`. Do not create a new rooted harness replacement for `protocols/wire-lab.d/specs/transport-spec-draft.md`; apparatus residue stays distributed in existing rooted artifacts. Split specimen-owned work now mixed into `TODO-turog` and `TODO-duvuk` into sim-local successors where appropriate.
Affects: `protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`; `protocols/wire-lab.d/TODO/TODO.md`; `simulations/`; `protocols/wire-lab.d/specs/harness-spec-draft.md`; `protocols/wire-lab.d/specs/transport-spec-draft.md`; `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`; `protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md`; `protocols/wire-lab.d/TODO/TODO-duvuk-te-42-filename-cid-cascade-policy.md`; `DEV-GUIDE-RESOURCES.md`.

## Context

The current simulation layout still reflects the recovery process that
created it rather than the stable content boundaries that later work
needs. The mixed tree at
`simulations/SIM-piloh-turns-149-208-recovery/` currently combines:

- concrete `wire-lab-devs` world evidence;
- candidate protocol trees for `group-session`, `udp-binding`, and
  `ppx-dr`;
- archive material for old proposal and transport design surfaces;
- recovery-only scaffolding and migration notes.

The rooted harness side still has specimen-owned residue:

- `protocols/wire-lab.d/specs/transport-spec-draft.md` mixes a
  specimen-side thin outer contract with rooted apparatus/meta residue;
- `TODO-turog` and `TODO-duvuk` still mix group-session and outer-feed
  ownership;
- several rooted current-pointer docs still point at `SIM-piloh`,
  rooted `protocols/grid-envelope.d/`, rooted `protocols/ppx-dr.d/`,
  or the old `udp-binding` name as if those were still the active
  specimen homes.

This TODO coordinates the content-named split so future replay and
design work can refer to stable specimen homes without re-opening the
basic placement questions.

## Target shape

- `simulations/SIM-<handle>-wire-lab-devs/` holds the concrete world
  evidence and archive/provenance that belong to the dogfood/developer
  coordination case.
- `simulations/SIM-<handle>-group-session/` holds the active
  `group-session.d` specimen tree and its TODO queue.
- `simulations/SIM-<handle>-udp-feed/` holds the renamed `udp-feed.d`
  specimen tree and its TODO queue.
- `simulations/SIM-<handle>-feed-outer/` holds the sim-local thin outer
  feed contract that currently lives specimen-side inside rooted
  `transport-spec-draft.md`.
- `simulations/SIM-<handle>-grid-envelope/` holds the
  `grid([pcid, payload])` working-hypothesis specimen tree and its owner
  TODO.
- No active `ppx-dr.d` tree remains after the split. Proposal history is
  archive-only provenance under the `wire-lab-devs` simulation.

## Subtasks

- [ ] rusis.1 Mint the simulation handles and scaffold the five
  content-named sims: `wire-lab-devs`, `group-session`, `feed-outer`,
  `udp-feed`, and `grid-envelope`.
- [ ] rusis.2 Dissolve
  `simulations/SIM-piloh-turns-149-208-recovery/` using `git mv` so no
  mixed umbrella sim remains.
- [ ] rusis.3 Move the concrete `wire-lab-devs` world state, proposal
  archive, and dogfood provenance into the `wire-lab-devs` sim.
- [ ] rusis.4 Move `group-session.d` into its own sim and keep
  `TODO-bisur` with the specimen.
- [ ] rusis.5 Move `udp-binding.d` into its own sim, rename it now to
  `udp-feed.d`, and keep `TODO-jodon` with the renamed specimen.
- [ ] rusis.6 Create the `grid-envelope` sim-local specimen home,
  draft/spec scaffold, and owner TODO; transfer the current grid-envelope
  handoff rows there.
- [ ] rusis.7 Create the `feed-outer` sim-local specimen home, move the
  specimen half of rooted `transport-spec-draft.md` there, and remove
  the rooted active draft.
- [ ] rusis.8 Retire the active `ppx-dr` tree, preserve only
  archive/provenance material, and update any current-pointer docs that
  still treat `ppx-dr` as an active specimen home.
- [ ] rusis.9 Split specimen-owned follow-on work out of rooted
  `TODO-turog` and `TODO-duvuk`, and update `TODO-kugod` so open UT rows
  point to the new sim-local owners instead of placeholders.
- [ ] rusis.10 Update rooted current-pointer docs, indexes, and guide
  resources to the new sim-local paths and names while leaving historical
  quotations untouched.
- [ ] rusis.11 Validate that no active `SIM-piloh` or `ppx-dr` tree
  remains, that active docs use `udp-feed`, and that `git diff --check`
  passes.
