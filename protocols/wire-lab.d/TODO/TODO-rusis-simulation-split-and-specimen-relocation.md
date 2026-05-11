# TODO-rusis: Simulation split and specimen relocation

## Prior aliases

This TODO is being filed after the TE-mumuv proquint migration locked
2026-05-07, so it is minted directly under `TODO-rusis`. No prior
integer or timestamp alias.

## Status

Open. Tracks the post-Mupoz split of the mixed recovery simulation into
content-named sims, the retirement of active `ppx-dr`, the `udp-binding`
-> `udp-feed` rename, and the relocation of specimen-owned work out of
`protocols/wire-lab.d/`. The split treats sims as independent evolving
lineages guided by the rooted harness, not as shared protocol homes.

## Decision Intent Log

ID: DI-tugit
Date: 2026-05-10 23:42:35
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Dissolve `simulations/SIM-piloh-turns-149-208-recovery/` and replace it with separate content-named sims for `wire-lab-devs`, `group-session`, `feed-outer`, `udp-feed`, and `grid-envelope`. Retire `ppx-dr` as an active protocol tree, rename `udp-binding` now to `udp-feed`, move the specimen half of rooted `transport-spec-draft.md` into a new `feed-outer` sim-local protocol tree, and redistribute the remaining rooted apparatus residue into existing harness docs and owner TODOs rather than creating a new rooted replacement draft.
Intent: The current mixed `SIM-piloh` tree is process-named and conflates concrete world evidence, candidate protocol specimens, legacy proposal archive, and rooted apparatus residue. The repo should expose simulation content by what it contains, keep rooted wire-lab apparatus focused on harness governance, and make specimen ownership explicit enough that replay cleanup does not have to re-litigate the same placement questions.
Constraints: Dissolve the mixed `SIM-piloh` tree completely rather than leaving an umbrella successor. Keep `wire-lab-devs` as the only concrete world simulation in this pass. Retire active `ppx-dr` trees but preserve proposal history under archive/provenance. Update active paths and current-pointer prose to `udp-feed`, not `udp-binding`. Do not create a new rooted harness replacement for `protocols/wire-lab.d/specs/transport-spec-draft.md`; apparatus residue stays distributed in existing rooted artifacts. Split specimen-owned work now mixed into `TODO-turog` and `TODO-duvuk` into sim-local successors where appropriate.
Affects: `protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`; `protocols/wire-lab.d/TODO/TODO.md`; `simulations/`; `protocols/wire-lab.d/specs/harness-spec-draft.md`; `protocols/wire-lab.d/specs/transport-spec-draft.md`; `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`; `protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md`; `protocols/wire-lab.d/TODO/TODO-duvuk-te-42-filename-cid-cascade-policy.md`; `DEV-GUIDE-RESOURCES.md`.

ID: DI-rubad
Date: 2026-05-11 02:39:21
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Treat the post-Mupoz sims as independent evolving candidate lineages rather than shared protocol homes or live dependencies. Narrow `rusis.1` to minimal `README.md` scaffolds only. Before any existing file moves, create a complete disposition table for every tracked artifact under `SIM-piloh` and for each rooted mixed owner artifact in scope; nothing is deleted without being moved, archived, left rooted, or superseded with provenance.
Intent: The harness should guide comparison, selection, mutation, and transfer of ideas across sims. Sims themselves should evolve independently, including local duplicates or variants when useful. The scaffold should not introduce extra sources of truth such as local protocol inventories, local decision logs, empty placeholder trees, or sim-to-sim dependency declarations before a lineage actually needs them.
Constraints: Do not create `QUESTION.md`, `protocol-set.md`, `decisions.md`, empty `protocols/`, empty `world/`, empty `archive/`, or empty `seed/` as part of `rusis.1`. Do not encode sim-to-sim live references as part of the scaffold. Preserve every existing TODO, TE, DR, DI, proposal/review record, migration note, message file, and spec draft by moving it with `git mv`, archiving it, leaving it rooted, or adding an explicit supersession/provenance note before removing its old active role.
Affects: `protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`; `protocols/wire-lab.d/TODO/TODO.md`; future `simulations/SIM-*-wire-lab-devs/README.md`; future `simulations/SIM-*-group-session/README.md`; future `simulations/SIM-*-feed-outer/README.md`; future `simulations/SIM-*-udp-feed/README.md`; future `simulations/SIM-*-grid-envelope/README.md`.
Supersedes: DI-tugit (shared-home model, `rusis.1` scaffold depth, and insufficient no-deletion discipline only)

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
design work can compare independent sim lineages without re-opening the
basic placement questions.

## Simulation model

Each simulation is an independently evolving candidate lineage. A sim is
not a shared protocol home, not a library, and not a live dependency of
another sim. Local files inside a sim are part of that candidate's state.
If two sims carry similar material, that duplication is expected: they
are separate lineages until the harness later compares, selects,
mutates, retires, or transfers an idea.

The rooted harness owns cross-sim comparison, replay cleanup,
selection, and provenance. Harness artifacts may record relationships
between sims, but a sim README should state only that sim's hypothesis
or lineage intent. It should not maintain a directory inventory or
point at another sim as a source of current truth.

## Preservation rule

No existing artifact is deleted as a cleanup shortcut. Every existing
tracked path in the mixed `SIM-piloh` tree and every rooted mixed owner
artifact in scope must receive one of these dispositions before move
work starts:

- **Move with lineage:** the artifact is local state for a continuing
  sim lineage and moves with `git mv`.
- **Archive:** the artifact is retired historical evidence and moves to
  an archive/provenance location.
- **Stay rooted:** the artifact is harness memory or cross-sim
  governance and remains under rooted harness paths.
- **Supersede:** the artifact remains available but gains a successor or
  supersession note before its old active role is retired.

Message `.txt` files are a stricter case: they move only byte-for-byte
with `git mv`, because their filenames depend on their byte content.

## Target shape

- `simulations/SIM-<handle>-wire-lab-devs/` is the concrete dogfood
  lineage for developer coordination evidence and any local design state
  needed by that lineage.
- `simulations/SIM-<handle>-group-session/` is an independent lineage
  exploring group-session design choices.
- `simulations/SIM-<handle>-udp-feed/` is an independent lineage
  exploring the renamed UDP feed design family.
- `simulations/SIM-<handle>-feed-outer/` is an independent lineage
  exploring the thin outer feed convention currently mixed into rooted
  `transport-spec-draft.md`.
- `simulations/SIM-<handle>-grid-envelope/` is an independent lineage
  exploring the `grid([pcid, payload])` working hypothesis.
- No active `ppx-dr.d` tree remains after the split. Proposal history is
  preserved as archive/provenance, not deleted.

## Subtasks

- [ ] rusis.1 Mint the simulation handles and scaffold only minimal
  `README.md` files for the five content-named sims: `wire-lab-devs`,
  `group-session`, `feed-outer`, `udp-feed`, and `grid-envelope`.
  Each README states the lineage intent and evaluation question; it does
  not list directory contents, declare shared homes, or point to another
  sim as current truth.
- [ ] rusis.2 Produce a complete file-by-file disposition table for
  every tracked path under
  `simulations/SIM-piloh-turns-149-208-recovery/`, plus each rooted
  mixed owner artifact in scope. No move, archive, supersession, or
  active-role retirement starts until this table exists.
- [ ] rusis.3 Dissolve
  `simulations/SIM-piloh-turns-149-208-recovery/` using `git mv` so no
  mixed umbrella sim remains and every file follows its recorded
  disposition.
- [ ] rusis.4 Move the concrete `wire-lab-devs` world state, proposal
  archive, and dogfood provenance according to the disposition table.
- [ ] rusis.5 Move `group-session.d` and `TODO-bisur` according to the
  disposition table, preserving them as local lineage state rather than
  treating them as a shared active protocol home.
- [ ] rusis.6 Move `udp-binding.d` into the UDP feed lineage, rename it
  to `udp-feed.d`, and keep `TODO-jodon` with the renamed local state.
- [ ] rusis.7 Create the `grid-envelope` lineage content and owner TODO
  only when the disposition table identifies the exact material to seed
  there.
- [ ] rusis.8 Create the `feed-outer` lineage content by extracting the
  specimen half of rooted `transport-spec-draft.md`; preserve the rooted
  draft's apparatus residue as rooted harness memory or supersession
  provenance.
- [ ] rusis.9 Retire the active `ppx-dr` tree, preserve only
  archive/provenance material, and update any current-pointer docs that
  still treat `ppx-dr` as an active specimen home.
- [ ] rusis.10 Split specimen-owned follow-on work out of rooted
  `TODO-turog` and `TODO-duvuk`, and update `TODO-kugod` so open UT rows
  point to concrete lineage/disposition records instead of placeholders.
- [ ] rusis.11 Update rooted current-pointer docs, indexes, and guide
  resources to the new sim-local paths and names while leaving historical
  quotations untouched.
- [ ] rusis.12 Validate that no active `SIM-piloh` or `ppx-dr` tree
  remains, that active docs use `udp-feed`, and that `git diff --check`
  passes.
