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

ID: DI-humam
Date: 2026-05-10 21:29:53
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Redefine `rusis.1` as the split-aware disposition pass. The first deliverable is not minimal sim scaffolds; it is a complete disposition table for every tracked `SIM-piloh` artifact and each rooted mixed-owner artifact in scope, with explicit handling for files that split across multiple successor locations or ownership domains.
Intent: The current blocker is not missing empty sim roots. It is the lack of an auditable map from each mixed umbrella artifact to its continuing lineage, archive/provenance home, rooted harness home, or supersession path. The table must be precise enough to prevent accidental stuffing of unrelated material into `wire-lab-devs` or any other successor lineage, while still allowing later move passes to stay mechanical and reviewable.
Constraints: Do not create new sim directories or files in `rusis.1`. Do not move, rename, archive, or delete any existing tracked artifact in this pass. Mixed files may use adaptive split units rather than a fixed section- or paragraph-only granularity, but each split unit must map to a concrete later owner or destination. Keep `simulations/README.md` untouched in this pass so it does not become a second live inventory. Preserve `DI-rubad` for the independent-lineage and no-deletion rules; this DI supersedes only the old scaffold-first interpretation of `rusis.1`.
Affects: `protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`.
Supersedes: DI-rubad (`rusis.1` scaffold wording only)

ID: DI-rugig
Date: 2026-05-10 22:26:38
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Lock `rusis.2` as the first concrete split-pass blueprint. The first execution pass creates all five successor sims with root `README.md` and `QUESTION.md`, but it moves substantive lineage content immediately only for `wire-lab-devs`, `group-session`, and `udp-feed`. The `feed-outer` and `grid-envelope` sims also exist in the first pass, but they start only with root docs plus `seed/extraction-sources.md`. Dissolved umbrella provenance and retired `ppx-dr` material go to a new rooted archive under `protocols/wire-lab.d/archive/`. The rooted mixed artifacts stay rooted in this first pass.
Intent: The disposition table is complete enough to lock the first execution wave without improvising what moves now, what stays rooted, and what only gets seeded. This keeps the first split mechanical, preserves provenance, and prevents the first pass from over-extracting `feed-outer` or `grid-envelope` before their source material is cleanly separable.
Constraints: Create all five successor sims in the first pass. Put `README.md` and `QUESTION.md` at each successor sim root. Route concrete `wire-lab-devs` world/evidence, `group-session.d`, and the renamed `udp-feed.d` lineage tree in the first pass. Rename the moved UDP lineage tree immediately from `protocols/udp-binding.d/` to `protocols/udp-feed.d/`, while preserving old-name provenance in seed notes and historical text. Create `feed-outer` and `grid-envelope` now, but seed them only with root docs and `seed/extraction-sources.md`; do not create protocol trees for them in the first pass. Preserve dissolved umbrella provenance and retired `ppx-dr` under `protocols/wire-lab.d/archive/`, including a migration subtree for `SIM-piloh` and a retired subtree for `ppx-dr`. Leave `protocols/wire-lab.d/specs/transport-spec-draft.md`, `protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md`, `protocols/wire-lab.d/TODO/TODO-duvuk-te-42-filename-cid-cascade-policy.md`, and `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md` rooted in this first pass. Do not update `simulations/README.md` until `SIM-piloh` is fully retired so the index does not advertise a transitional double-state.
Affects: `protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`; future `protocols/wire-lab.d/archive/`; future `simulations/SIM-*-wire-lab-devs/`; future `simulations/SIM-*-group-session/`; future `simulations/SIM-*-udp-feed/`; future `simulations/SIM-*-feed-outer/`; future `simulations/SIM-*-grid-envelope/`.
Supersedes: DI-tugit (`rusis.2` execution layout only)

ID: DI-limom
Date: 2026-05-11 10:12:24
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Execute `rusis.3` by creating the five successor sim roots and the rooted archive skeleton only. Use `README.md` and `QUESTION.md` at each sim root. Use `README.md` marker files at `protocols/wire-lab.d/archive/`, `protocols/wire-lab.d/archive/migrations/SIM-piloh-turns-149-208-recovery/`, and `protocols/wire-lab.d/archive/retired/ppx-dr/`. Do not move lineage content, split rooted mixed artifacts, or update `simulations/README.md` in this pass.
Intent: The first filesystem pass should make the future destinations concrete and reviewable without mixing destination creation with content movement or extraction. That keeps `rusis.3` structural and leaves the substantive lineage moves to the later `rusis.*` tasks already locked by `DI-rugig`.
Constraints: Create `simulations/SIM-ludut-wire-lab-devs/`, `simulations/SIM-rakot-group-session/`, `simulations/SIM-ludaf-udp-feed/`, `simulations/SIM-labit-feed-outer/`, and `simulations/SIM-kurim-grid-envelope/`. Keep each new sim root minimal in this pass: only `README.md` and `QUESTION.md`. Create only the rooted archive marker files needed to track the archive skeleton in Git. Do not create `seed/`, `protocols/`, `world/`, or archive-content files yet. Do not edit existing `SIM-piloh` content or the rooted mixed artifacts in this pass.
Affects: `protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`; `protocols/wire-lab.d/archive/README.md`; `protocols/wire-lab.d/archive/migrations/SIM-piloh-turns-149-208-recovery/README.md`; `protocols/wire-lab.d/archive/retired/ppx-dr/README.md`; `simulations/SIM-ludut-wire-lab-devs/README.md`; `simulations/SIM-ludut-wire-lab-devs/QUESTION.md`; `simulations/SIM-rakot-group-session/README.md`; `simulations/SIM-rakot-group-session/QUESTION.md`; `simulations/SIM-ludaf-udp-feed/README.md`; `simulations/SIM-ludaf-udp-feed/QUESTION.md`; `simulations/SIM-labit-feed-outer/README.md`; `simulations/SIM-labit-feed-outer/QUESTION.md`; `simulations/SIM-kurim-grid-envelope/README.md`; `simulations/SIM-kurim-grid-envelope/QUESTION.md`.
Supersedes: DI-rugig (`rusis.3` execution details only)

ID: DI-lilov
Date: 2026-05-11 10:18:23
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Execute `rusis.4` by moving only the concrete `wire-lab-devs` `world/` subtree and `seed/wire-lab-devs-draft-migration.md` from `SIM-piloh` into `SIM-ludut-wire-lab-devs/`. Defer `SIM-piloh/archive/transports/*` again for a later `rusis.*` pass instead of moving it now.
Intent: The world/evidence slice is already a clean lineage-owned unit, but the transport archive still has unresolved placement between lineage-local provenance and the rooted migration archive. `rusis.4` should move the concrete world now without forcing an archive decision prematurely.
Constraints: Use `git mv` for all moved tracked files. Create only the destination `seed/` directory needed for the moved migration note. Do not touch `SIM-piloh/archive/transports/*`, rooted mixed artifacts, or `simulations/README.md` in this pass. Preserve the moved `.txt` message files byte-for-byte.
Affects: `protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`; `simulations/SIM-ludut-wire-lab-devs/world/`; `simulations/SIM-ludut-wire-lab-devs/seed/wire-lab-devs-draft-migration.md`; `simulations/SIM-piloh-turns-149-208-recovery/world/`; `simulations/SIM-piloh-turns-149-208-recovery/seed/wire-lab-devs-draft-migration.md`.
Supersedes: DI-rugig (`rusis.4` transport-archive timing only)

ID: DI-firiv
Date: 2026-05-11 10:29:25
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Execute `rusis.5` by moving the complete `group-session.d` lineage tree from `SIM-piloh` into `SIM-rakot-group-session/protocols/group-session.d/`, and add one lineage-local migration note at `simulations/SIM-rakot-group-session/seed/protocol-migration.md`.
Intent: The `group-session` tree is already a coherent lineage-owned unit in the `rusis.1` table. `rusis.5` should relocate that unit intact, then preserve the old mixed-sim provenance in a short local note so later passes do not depend on the obsolete umbrella inventory.
Constraints: Use `git mv` for the entire tracked `group-session.d` tree. Create only the destination parent `protocols/` and `seed/` directories plus the new `seed/protocol-migration.md` file. Do not edit rooted mixed artifacts, `simulations/README.md`, or any non-`group-session` lineage content in this pass.
Affects: `protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`; `simulations/SIM-rakot-group-session/protocols/group-session.d/`; `simulations/SIM-rakot-group-session/seed/protocol-migration.md`; `simulations/SIM-piloh-turns-149-208-recovery/protocols/group-session.d/`.
Supersedes: DI-rugig (`rusis.5` execution details only)

ID: DI-loluk
Date: 2026-05-11 10:38:24
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Execute `rusis.6` by moving the complete UDP lineage tree from `SIM-piloh` into `SIM-ludaf-udp-feed/protocols/udp-feed.d/`, renaming the active lineage tree from `udp-binding.d` to `udp-feed.d` during the move, and adding one local migration note at `simulations/SIM-ludaf-udp-feed/seed/protocol-migration.md`.
Intent: The `rusis.1` table already treats the UDP tree as a coherent lineage unit, and `DI-rugig` locked the active rename to `udp-feed`. `rusis.6` should perform the physical tree rename now while keeping the old `udp-binding` name alive only as historical provenance.
Constraints: Use `git mv` for the whole tracked UDP tree. Create only the destination parent `protocols/` and `seed/` directories plus the new `seed/protocol-migration.md` file. Do not rewrite internal file titles or prose in this pass; preserve the old `udp-binding` name there until a later cleanup pass. Do not edit rooted mixed artifacts, `simulations/README.md`, or any non-UDP lineage content in this pass.
Affects: `protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`; `simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/`; `simulations/SIM-ludaf-udp-feed/seed/protocol-migration.md`; `simulations/SIM-piloh-turns-149-208-recovery/protocols/udp-binding.d/`.
Supersedes: DI-rugig (`rusis.6` execution details only)

ID: DI-ludaz
Date: 2026-05-11 22:38:57
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Execute `rusis.7` by creating `simulations/SIM-labit-feed-outer/seed/extraction-sources.md` as the feed-outer extraction source note, while keeping `protocols/wire-lab.d/specs/transport-spec-draft.md` rooted and unchanged until the later extraction pass.
Intent: The feed-outer lineage needs a concrete seed record before the rooted transport draft is split. This note makes the later extraction source explicit without prematurely copying, rewriting, or moving rooted apparatus material.
Constraints: Create only the `seed/` directory and `seed/extraction-sources.md` under `SIM-labit-feed-outer`. Do not create a `protocols/feed-outer.d/` tree in this pass. Do not edit `transport-spec-draft.md`, rooted mixed TODOs, `simulations/README.md`, or other lineage content.
Affects: `protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`; `simulations/SIM-labit-feed-outer/seed/extraction-sources.md`.
Supersedes: DI-rugig (`rusis.7` execution details only)

ID: DI-nijon
Date: 2026-05-11 23:07:39
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Execute `rusis.8` by creating `simulations/SIM-kurim-grid-envelope/seed/extraction-sources.md` as the grid-envelope extraction source note, while keeping protocol-tree creation and specimen prose extraction deferred for a later pass.
Intent: Make the grid-envelope successor owner explicit for the transferred-open `TODO-kugod` rows (`UT-157.a`, `UT-157.c`, `UT-158.f`) without prematurely copying or rewriting rooted or cross-lineage protocol text.
Constraints: Create only `seed/` and `seed/extraction-sources.md` under `SIM-kurim-grid-envelope` in this pass. Do not create `protocols/grid-envelope.d/` yet. Do not edit rooted harness docs, `TODO-kugod`, `group-session` draft text, or `simulations/README.md` in this pass.
Affects: `protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`; `simulations/SIM-kurim-grid-envelope/seed/extraction-sources.md`.
Supersedes: DI-rugig (`rusis.8` execution details only)

ID: DI-turud
Date: 2026-05-11 23:20:11
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Execute `rusis.9` by archiving whole-file `SIM-piloh` historical artifacts into the rooted migration archive and retiring the active `ppx-dr` lineage tree into the rooted retired archive.
Intent: Retire `SIM-piloh` umbrella-history and `ppx-dr` active-lineage residue without deleting evidence, while keeping mixed files that need section-level decomposition for later passes.
Constraints: Move only whole-file archive targets in this pass: `QUESTION.md`, `archive/proposals/*`, and `archive/transports/*` from `SIM-piloh` into `protocols/wire-lab.d/archive/migrations/SIM-piloh-turns-149-208-recovery/`; move `protocols/ppx-dr.d/*` from `SIM-piloh` into `protocols/wire-lab.d/archive/retired/ppx-dr/`. Do not split `README.md`, `concerns.md`, `decisions.md`, `protocol-set.md`, `seed/protocol-tree-migrations.md`, or rooted mixed-owner TODO/spec artifacts in this pass.
Affects: `protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`; `protocols/wire-lab.d/archive/migrations/SIM-piloh-turns-149-208-recovery/`; `protocols/wire-lab.d/archive/retired/ppx-dr/`; `simulations/SIM-piloh-turns-149-208-recovery/`.
Supersedes: DI-rugig (`rusis.9` execution details only)

ID: DI-nubun
Date: 2026-05-11 23:22:16
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Clarify `rusis.13` validation wording so it checks deprecated-name cleanup and neutrality constraints rather than implying any preferred simulation home.
Intent: Keep the split model explicitly lineage-neutral while still enforcing the `udp-binding` -> `udp-feed` active-name migration outside historical quotations.
Constraints: Update only `rusis.13` wording in this pass. Preserve the existing `rusis.13` goals about retiring active `SIM-piloh` and `ppx-dr` trees and running `git diff --check`.
Affects: `protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`.

ID: DI-mosor
Date: 2026-05-11 23:24:03
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Execute `rusis.10` by routing specimen-owned follow-on work from rooted `TODO-turog` and `TODO-duvuk` into concrete sim-local successor TODO artifacts, and by updating `TODO-kugod` open/transferred UT rows to reference concrete successor records instead of placeholders.
Intent: Keep rooted TODOs as historical coordination memory while making active lineage-owned work explicit inside independent simulation trees, so later replay does not depend on ambiguous "until successor exists" wording.
Constraints: Keep rooted `TODO-turog` and `TODO-duvuk` files as historical records; do not delete their thread history or DI logs. Create successor owner artifacts under `SIM-rakot-group-session`, `SIM-labit-feed-outer`, and `SIM-kurim-grid-envelope`. Update only ownership/status pointers in `TODO-kugod`; do not claim those UTs are resolved by this routing pass.
Affects: `protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`; `protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md`; `protocols/wire-lab.d/TODO/TODO-duvuk-te-42-filename-cid-cascade-policy.md`; `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`; `simulations/SIM-rakot-group-session/protocols/group-session.d/TODO/TODO.md`; `simulations/SIM-rakot-group-session/protocols/group-session.d/TODO/TODO-gapab-group-session-freeze-and-cid-cascade-follow-on.md`; `simulations/SIM-labit-feed-outer/TODO/TODO.md`; `simulations/SIM-labit-feed-outer/TODO/TODO-kakaz-feed-outer-freeze-and-cid-cascade-follow-on.md`; `simulations/SIM-kurim-grid-envelope/TODO/TODO.md`; `simulations/SIM-kurim-grid-envelope/TODO/TODO-tujad-grid-envelope-successor-owner.md`.

ID: DI-huzor
Date: 2026-05-12 08:29:13
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Execute `rusis.11` by extracting specimen-side outer-feed sections from rooted `transport-spec-draft.md` into `SIM-labit-feed-outer` as `protocols/feed-outer.d/specs/feed-outer-draft.md`, while preserving rooted apparatus/governance residue in place.
Intent: Complete the locked apparatus-vs-specimen split for the rooted transport draft without deleting provenance, and make the feed-outer lineage carry its own candidate-specimen draft.
Constraints: Move the specimen-side units (`## Purpose`, `## The four locked principles`, `## What this spec does NOT specify`) into the feed-outer sim-local draft. Keep rooted units (`## Sources`, `## The per-axis meta-rule`, `## Open questions`, `## Freeze gate`) in rooted `transport-spec-draft.md` with explicit supersession notes. Do not treat the rooted draft as feed-outer's active specimen home after extraction.
Affects: `protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`; `protocols/wire-lab.d/specs/transport-spec-draft.md`; `simulations/SIM-labit-feed-outer/protocols/feed-outer.d/specs/feed-outer-draft.md`; `simulations/SIM-labit-feed-outer/seed/extraction-sources.md`.

ID: DI-jugub
Date: 2026-05-12 08:29:13
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Execute `rusis.12` as a deep split of the remaining `SIM-piloh` shell files: move those mixed shell artifacts into rooted migration archive storage, route live ownership to concrete successor artifacts, and update rooted current-pointer docs/indexes/resources to content-named lineage homes.
Intent: Retire the mixed umbrella simulation as an active specimen surface while preserving all historical evidence and avoiding any implied preferred simulation home.
Constraints: No deletion of historical shell content; preserve by `git mv` under rooted archive migration paths plus explicit supersession mapping. Update active pointers in `simulations/README.md`, rooted TODO index links, guide-resource mappings, and harness current-pointer prose; leave historical quotations and archived evidence wording intact.
Affects: `protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`; `protocols/wire-lab.d/archive/migrations/SIM-piloh-turns-149-208-recovery/`; `simulations/README.md`; `protocols/wire-lab.d/TODO/TODO.md`; `DEV-GUIDE-RESOURCES.md`; `protocols/wire-lab.d/specs/harness-spec-draft.md`; successor lineage seed/docs touched by supersession routing.

ID: DI-patoz
Date: 2026-05-12 08:29:13
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Execute `rusis.13` validation as a hard gate over active docs and trees after `rusis.11`/`rusis.12`: no active `SIM-piloh`/`ppx-dr` tree, no deprecated `udp-binding` in active contexts except historical/archive usage, no privileged simulation-home wording, and `git diff --check` pass.
Intent: Close the split pass with explicit neutrality and naming-hygiene checks so replay cleanup does not regress into process-named or privileged-home phrasing.
Constraints: Historical TE narrative, DI logs, and archive/provenance files may retain legacy names and paths as historical evidence. Active pointers and owner records must use current lineage homes.
Affects: `protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`; all active current-pointer docs swept in `rusis.12`/`rusis.13`.

ID: DI-navod
Date: 2026-05-17 00:06:19
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Backfill turn-177 successor work into four standalone, root-only
simulation charters: `SIM-jomag-cas-object-model`,
`SIM-zazit-chunk-feed-replication`,
`SIM-jurar-cas-backed-group-session`, and
`SIM-rusap-promise-accounting-records`.
Intent: Turn 177's conclusions should evolve as simulation specimens instead
of remaining only in replay notes or being forced into existing grid-envelope,
group-session, feed, or wire-lab-devs lineages. Separate simulations preserve
independent evolution across CAS object modeling, chunk-feed replication,
CAS-backed group semantics, and peer-local promise accounting.
Constraints: Create only root `README.md` and `QUESTION.md` files in this pass.
Do not create protocol trees, world fixtures, manifests, local TODO queues, or
shared protocol bundles. Do not edit the 24 grid-envelope variants. Treat the
new simulations as provisional design-point workspaces, not canonical
PromiseGrid APIs or preferred outcomes.
Affects: `protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`;
`simulations/README.md`; `DEV-GUIDE-RESOURCES.md`;
`simulations/SIM-jomag-cas-object-model/README.md`;
`simulations/SIM-jomag-cas-object-model/QUESTION.md`;
`simulations/SIM-zazit-chunk-feed-replication/README.md`;
`simulations/SIM-zazit-chunk-feed-replication/QUESTION.md`;
`simulations/SIM-jurar-cas-backed-group-session/README.md`;
`simulations/SIM-jurar-cas-backed-group-session/QUESTION.md`;
`simulations/SIM-rusap-promise-accounting-records/README.md`;
`simulations/SIM-rusap-promise-accounting-records/QUESTION.md`;
`protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`protocols/wire-lab.d/TODO/TODO-kituj-te-43-promisebase-prior-art-adoption.md`;
`protocols/wire-lab.d/TODO/TODO-pipus-te-39-wire-lab-devs-migration.md`;
`protocols/wire-lab.d/TODO/TODO-kulih-te-spec-doc-as-promise.md`;
`protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`.

ID: DI-pator
Date: 2026-05-17 09:22:51
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Add docs-only `SCENARIOS.md` matrices to the four turn-177
simulation charters created under `DI-navod`.
Intent: The root charters make the turn-177 design points visible, but cleanup
still needs concrete scenario pressure so successor owners can act without
re-litigating what each simulation is for. Scenario files make CAS object
typing, chunk-feed replication, CAS-backed group semantics, and peer-local
promise accounting testable while preserving the actual DF/DI decisions for
TODO-kituj / TE-43, TODO-pipus, and TODO-kulih / TE-nibar.
Constraints: Documentation-only pass. Create only `SCENARIOS.md` files under
the four existing turn-177 simulation roots and update rooted owner pointers.
Do not create protocol trees, world fixtures, manifests, local TODO queues, or
shared protocol bundles. Do not edit the 24 grid-envelope variants. Do not
choose CBOR profile, chunking parameters, CID codec assignments, migration
mechanics, or spec-doc requirements in this pass.
Affects: `protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`;
`simulations/SIM-jomag-cas-object-model/SCENARIOS.md`;
`simulations/SIM-zazit-chunk-feed-replication/SCENARIOS.md`;
`simulations/SIM-jurar-cas-backed-group-session/SCENARIOS.md`;
`simulations/SIM-rusap-promise-accounting-records/SCENARIOS.md`;
`protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`protocols/wire-lab.d/TODO/TODO-kituj-te-43-promisebase-prior-art-adoption.md`;
`protocols/wire-lab.d/TODO/TODO-pipus-te-39-wire-lab-devs-migration.md`;
`protocols/wire-lab.d/TODO/TODO-kulih-te-spec-doc-as-promise.md`;
`protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`.

ID: DI-davov
Date: 2026-05-17 09:28:02
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: File three scoped DRs for the remaining turn-177 cleanup decisions:
`DR-tumus` for concrete L6 CAS adoption, `DR-gabif` for additive
CAS-backed group-session migration, and `DR-robon` for turn-177 spec-shape
requirements.
Intent: The turn-177 simulation charters and scenarios make the design pressure
visible, but the remaining work is now blocked on explicit decisions. Filing
DRs gives TODO-kituj / TE-43, TODO-pipus, and TODO-kulih / TE-nibar concrete
question records so replay cleanup does not have to re-derive the same loose
ends from raw turn notes.
Constraints: Do not answer the DRs in this pass. Do not mark turn 177 fully
closed. Preserve existing `TODO-lilar` closure state; the DRs are successor
decision records, not reopened UT rows.
Affects: `DR/DR-tumus-turn-177-l6-cas-adoption.md`;
`DR/DR-gabif-turn-177-cas-backed-group-session-migration.md`;
`DR/DR-robon-turn-177-spec-shape-requirements.md`;
`protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`;
`protocols/wire-lab.d/TODO/TODO-kituj-te-43-promisebase-prior-art-adoption.md`;
`protocols/wire-lab.d/TODO/TODO-pipus-te-39-wire-lab-devs-migration.md`;
`protocols/wire-lab.d/TODO/TODO-kulih-te-spec-doc-as-promise.md`;
`protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`;
`DEV-GUIDE-RESOURCES.md`.

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
- Each successor sim begins with root `README.md` and `QUESTION.md`.
- `protocols/wire-lab.d/archive/` is the rooted archive home for
  dissolved umbrella provenance and retired active lineages.
- No active `ppx-dr.d` tree remains after the split. Proposal history is
  preserved as archive/provenance, not deleted.

## Subtasks

- [x] rusis.1 Produce a complete split-aware disposition table for
  every tracked path under
  `simulations/SIM-piloh-turns-149-208-recovery/`, plus each rooted
  mixed owner artifact in scope. The table must allow whole-file and
  adaptive split-unit rows, and no move, archive, supersession, or
  active-role retirement starts until it exists.
- [x] rusis.2 Lock the first concrete split-pass blueprint: create all
  five successor sims with root `README.md` and `QUESTION.md`; move
  substantive lineage content immediately only for `wire-lab-devs`,
  `group-session`, and `udp-feed`; create `feed-outer` and
  `grid-envelope` in the first pass with root docs plus
  `seed/extraction-sources.md`; preserve dissolved umbrella provenance
  and retired `ppx-dr` under `protocols/wire-lab.d/archive/`; keep the
  rooted mixed artifacts rooted; and defer `simulations/README.md`
  updates until the active sim set is no longer transitional.
- [x] rusis.3 Create the rooted archive skeleton and the five successor
  sim roots described by `rusis.2`, without yet splitting rooted mixed
  artifacts.
- [x] rusis.4 Move the concrete `wire-lab-devs` world state and local
  provenance into the `wire-lab-devs` lineage according to the
  disposition table, and defer `archive/transports/*` again for a later
  `rusis.*` pass.
- [x] rusis.5 Move `group-session.d` and `TODO-bisur` into the
  `group-session` lineage, preserving them as local lineage state rather
  than treating them as a shared protocol home.
- [x] rusis.6 Move `udp-binding.d` into the UDP feed lineage, rename it
  to `udp-feed.d`, keep `TODO-jodon` with the renamed local state, and
  record the old-name provenance in lineage seed notes.
- [x] rusis.7 Create the `feed-outer` lineage root docs and
  `seed/extraction-sources.md`, but defer extraction from rooted
  `transport-spec-draft.md` until the later split pass.
- [x] rusis.8 Create the `grid-envelope` lineage root docs and
  `seed/extraction-sources.md`, but defer creation of a protocol tree
  until the later extraction pass.
- [x] rusis.9 Archive the dissolved umbrella-sim material and retire the
  active `ppx-dr` tree under `protocols/wire-lab.d/archive/`, preserving
  the mixed umbrella files whole where the disposition table calls for
  historical archive rather than per-lineage decomposition.
- [x] rusis.10 Split specimen-owned follow-on work out of rooted
  `TODO-turog` and `TODO-duvuk`, and update `TODO-kugod` so open UT rows
  point to concrete lineage/disposition records instead of placeholders.
- [x] rusis.11 Extract the specimen-side `feed-outer` material out of
  rooted `transport-spec-draft.md`, while preserving or superseding the
  rooted apparatus residue in place.
- [x] rusis.12 Update rooted current-pointer docs, indexes, and guide
  resources after `SIM-piloh` is retired and the active sim set is no
  longer transitional, leaving historical quotations untouched.
- [x] rusis.13 Validate that no active `SIM-piloh` or `ppx-dr` tree
  remains, that active docs do not use deprecated `udp-binding` except
  historical quotations, that active docs do not imply a privileged
  simulation home, and that `git diff --check` passes.
- [x] rusis.14 Backfill turn-177 successor work into standalone
  simulation charters for CAS object modeling, chunk-feed replication,
  CAS-backed group-session semantics, and peer-local promise accounting.
- [x] rusis.15 Add scenario matrices to the four turn-177 simulations so
  TODO-kituj / TE-43, TODO-pipus, and TODO-kulih / TE-nibar have concrete
  simulation pressure without moving decisions into the simulation charters.
- [x] rusis.16 File scoped DRs for the three remaining turn-177 decision
  clusters: concrete L6 CAS adoption, additive CAS-backed group-session
  migration, and spec-shape requirements.

## `rusis.2` first concrete split-pass blueprint

The first concrete split pass creates all five successor sims and the
rooted archive home, but it does not yet perform every extraction
identified by the disposition table.

- **Create now**
  - `simulations/SIM-<handle>-wire-lab-devs/`
  - `simulations/SIM-<handle>-group-session/`
  - `simulations/SIM-<handle>-udp-feed/`
  - `simulations/SIM-<handle>-feed-outer/`
  - `simulations/SIM-<handle>-grid-envelope/`
  - `protocols/wire-lab.d/archive/`
- **Put at each successor sim root now**
  - `README.md`
  - `QUESTION.md`
- **Move substantive lineage content now**
  - `wire-lab-devs`: move the concrete `world/` subtree,
    `seed/wire-lab-devs-draft-migration.md`, and the transport-side
    archive/provenance material that belongs with the dogfood lineage
  - `group-session`: move the whole `protocols/group-session.d/` tree
    and add a lineage migration seed note
  - `udp-feed`: move the current UDP lineage tree, rename it immediately
    to `protocols/udp-feed.d/`, and add a lineage migration seed note
- **Seed now, extract later**
  - `feed-outer`: create root docs plus `seed/extraction-sources.md`;
    do not create a protocol tree in the first pass
  - `grid-envelope`: create root docs plus `seed/extraction-sources.md`;
    do not create a protocol tree in the first pass
- **Archive now**
  - dissolved umbrella provenance under
    `protocols/wire-lab.d/archive/migrations/SIM-piloh-turns-149-208-recovery/`
  - retired `ppx-dr` material under
    `protocols/wire-lab.d/archive/retired/ppx-dr/`
- **Keep rooted in the first pass**
  - `protocols/wire-lab.d/specs/transport-spec-draft.md`
  - `protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md`
  - `protocols/wire-lab.d/TODO/TODO-duvuk-te-42-filename-cid-cascade-policy.md`
  - `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`
- **Defer**
  - section-level extraction of rooted mixed artifacts
  - successor-owner routing for the mixed rooted TODOs
  - `simulations/README.md` updates until the active sim catalog is no
    longer transitional

## `rusis.1` split-aware disposition table

Use the smallest sensible unit for mixed files. Whole-file rows are
used where the file is already coherent. Split rows are used where one
file contains material that will later belong to different successor
lineages, rooted harness memory, or archive/provenance.

| Current path | Current unit | Current role | Disposition | Target destination / successor owner | Later-phase note / blocker |
|---|---|---|---|---|---|
| `simulations/SIM-piloh-turns-149-208-recovery/README.md` | whole file | Umbrella simulation identity + obsolete contract summary | split | see child rows below | Old umbrella sim is being dissolved; this file cannot move whole |
| `simulations/SIM-piloh-turns-149-208-recovery/README.md` | intro paragraphs above `## Contract` | Historical description of the first recovery world and its purpose | archive | `SIM-piloh` provenance note under rooted or archived migration history | Keep as historical evidence of the umbrella simulation rather than current lineage charter |
| `simulations/SIM-piloh-turns-149-208-recovery/README.md` | `## Contract` | Obsolete contract for the umbrella sim (`QUESTION.md`, `protocol-set.md`, `decisions.md`, etc.) | supersede | `TODO-rusis` + later successor sim file-shape decisions | Do not copy this contract forward blindly into successor sims |
| `simulations/SIM-piloh-turns-149-208-recovery/QUESTION.md` | whole file | Umbrella recovery question for the mixed sim | archive whole | `SIM-piloh` provenance note under rooted or archived migration history | Later successor sims get their own question/charter artifacts if needed |
| `simulations/SIM-piloh-turns-149-208-recovery/archive/proposals/MANIFEST.md` | whole file | Manifest for proposal archive contents | archive whole | retired proposal/provenance archive | Keep with archived proposal records; not a continuing lineage artifact |
| `simulations/SIM-piloh-turns-149-208-recovery/archive/proposals/approved/ppx-merge-all-20260429-164729/review-20260429-170826-steve-traugott.md` | whole file | Historical approved proposal review | archive whole | retired proposal/provenance archive | Preserve exactly; no active sim ownership implied |
| `simulations/SIM-piloh-turns-149-208-recovery/archive/proposals/pending/ppx-dr-001-bootstrap/contest-20260429-033208-steve-traugott.md` | whole file | Historical proposal/contest evidence | archive whole | retired `ppx-dr` archive/provenance location | Preserve with retired `ppx-dr` history |
| `simulations/SIM-piloh-turns-149-208-recovery/archive/proposals/pending/ppx-dr-001-bootstrap/review-20260429-162212-steve-traugott.md` | whole file | Historical proposal/review evidence | archive whole | retired `ppx-dr` archive/provenance location | Preserve with retired `ppx-dr` history |
| `simulations/SIM-piloh-turns-149-208-recovery/archive/proposals/pending/ppx-te-20260428-202400-promise-stack-ordering/review-20260429-162212-steve-traugott.md` | whole file | Historical TE review evidence | archive whole | retired proposal/provenance archive | Stays archival; not a continuing sim lineage artifact |
| `simulations/SIM-piloh-turns-149-208-recovery/archive/transports/MANIFEST.md` | whole file | Manifest for archived transport-side provenance | archive whole | `wire-lab-devs` lineage provenance archive or rooted migration archive | Exact destination chosen in move pass after archive layout is fixed |
| `simulations/SIM-piloh-turns-149-208-recovery/archive/transports/README.md` | whole file | Historical transport archive explanation | archive whole | `wire-lab-devs` lineage provenance archive or rooted migration archive | Keep as history; do not treat as live feed/group spec |
| `simulations/SIM-piloh-turns-149-208-recovery/concerns.md` | whole file | Mixed concerns matrix spanning rooted recovery governance and sim evidence | split | see child rows below | File mixes rooted recovery concerns with local specimen evidence |
| `simulations/SIM-piloh-turns-149-208-recovery/concerns.md` | recovery-closure row | Replay/recovery closure concern tied to `TODO-jivam` / future observations | stay rooted | rooted recovery ownership (`TODO-jivam` / `TODO-juhub` chain) | This is harness closure governance, not successor sim-local state |
| `simulations/SIM-piloh-turns-149-208-recovery/concerns.md` | apparatus-vs-specimen row | Harness/specimen boundary concern | stay rooted | rooted apparatus/specimen ownership (`TODO-kugod`) | Governs the harness; not a sim-local charter |
| `simulations/SIM-piloh-turns-149-208-recovery/concerns.md` | legacy-proposal row | Guide-writer/process concern anchored in archives and guide resources | split | retired proposal/provenance archive + rooted guide-resource ownership | Archive evidence survives; active workflow pointers stay rooted |
| `simulations/SIM-piloh-turns-149-208-recovery/concerns.md` | transport-evidence row | Byte-verifiable message-evidence concern | move with lineage | `wire-lab-devs` lineage provenance / evidence notes | Concerns the concrete transport specimen directly |
| `simulations/SIM-piloh-turns-149-208-recovery/concerns.md` | provenance-graduation row | Rule that simulation results graduate only via later provenance | stay rooted | rooted harness / migration-governance memory | Applies across sims, not just one successor |
| `simulations/SIM-piloh-turns-149-208-recovery/decisions.md` | whole file | Umbrella-sim decision/graduation summary | split | see child rows below | Mixes rooted migration decisions with umbrella-sim-specific graduation wording |
| `simulations/SIM-piloh-turns-149-208-recovery/decisions.md` | `## Locked inputs` | Summary of already-rooted DFs and DIs (`DI-pakid`, `DI-fakin`) | stay rooted | rooted migration provenance (`TODO-kugod` / `TODO-rusis`) | Duplicate summary; should not become successor sim-local truth |
| `simulations/SIM-piloh-turns-149-208-recovery/decisions.md` | `## Graduation rule` | Cross-sim graduation rule | stay rooted | rooted harness provenance / migration guidance | Graduation remains a harness rule, not a single lineage file |
| `simulations/SIM-piloh-turns-149-208-recovery/events/README.md` | whole file | Placeholder for umbrella-sim event logging | supersede whole | later successor-specific event artifact, if any | Do not create empty event logs by default in successor sims |
| `simulations/SIM-piloh-turns-149-208-recovery/observations/README.md` | whole file | Placeholder for umbrella-sim observations | supersede whole | later successor-specific observation artifact, if any | Do not create empty observation logs by default in successor sims |
| `simulations/SIM-piloh-turns-149-208-recovery/protocol-set.md` | whole file | Mixed protocol inventory for three lineages under one umbrella | split | see child rows below | This inventory cannot move whole once lineages become independent |
| `simulations/SIM-piloh-turns-149-208-recovery/protocol-set.md` | `group-session` row | Current `group-session` specimen mapping | move with lineage | `group-session` lineage provenance note | Keep provenance with the lineage, not with an umbrella registry |
| `simulations/SIM-piloh-turns-149-208-recovery/protocol-set.md` | `udp-binding` row | Current UDP specimen mapping under old name | move with lineage | `udp-feed` lineage provenance note | Rename lineage later; preserve old-name provenance |
| `simulations/SIM-piloh-turns-149-208-recovery/protocol-set.md` | `ppx-dr` row | Retired proposal/review specimen mapping | archive | retired `ppx-dr` archive/provenance location | Preserve as retirement provenance, not an active lineage registry |
| `simulations/SIM-piloh-turns-149-208-recovery/protocols/group-session.d/CHANGELOG.md` | whole file | Local lineage changelog | move whole | `group-session` lineage | Moves with the local lineage tree |
| `simulations/SIM-piloh-turns-149-208-recovery/protocols/group-session.d/TODO/TODO-bisur-group-transport-envelope.md` | whole file | Local lineage TODO owner | move whole | `group-session` lineage | Owner file moves with lineage tree |
| `simulations/SIM-piloh-turns-149-208-recovery/protocols/group-session.d/TODO/TODO.md` | whole file | Local lineage TODO index | move whole | `group-session` lineage | Moves with lineage tree |
| `simulations/SIM-piloh-turns-149-208-recovery/protocols/group-session.d/manifest.json` | whole file | Local lineage manifest | move whole | `group-session` lineage | Moves with lineage tree |
| `simulations/SIM-piloh-turns-149-208-recovery/protocols/group-session.d/specs/group-session-draft.md` | whole file | Local lineage draft spec | move whole | `group-session` lineage | Moves with lineage tree |
| `simulations/SIM-piloh-turns-149-208-recovery/protocols/ppx-dr.d/CHANGELOG.md` | whole file | Retired lineage changelog | archive whole | retired `ppx-dr` archive/provenance location | Preserve as retirement history |
| `simulations/SIM-piloh-turns-149-208-recovery/protocols/ppx-dr.d/TODO/TODO-pozig-proposals-as-transport-messages-BLOCKED.md` | whole file | Retired lineage/blocker record | archive whole | retired `ppx-dr` archive/provenance location | Preserve blocked state as historical evidence |
| `simulations/SIM-piloh-turns-149-208-recovery/protocols/ppx-dr.d/TODO/TODO.md` | whole file | Retired lineage TODO index | archive whole | retired `ppx-dr` archive/provenance location | Preserve as historical evidence only |
| `simulations/SIM-piloh-turns-149-208-recovery/protocols/ppx-dr.d/manifest.json` | whole file | Retired lineage manifest | archive whole | retired `ppx-dr` archive/provenance location | Preserve with retired lineage |
| `simulations/SIM-piloh-turns-149-208-recovery/protocols/udp-binding.d/CHANGELOG.md` | whole file | Local UDP lineage changelog | move whole | `udp-feed` lineage | Moves with lineage tree; later rename preserves provenance |
| `simulations/SIM-piloh-turns-149-208-recovery/protocols/udp-binding.d/TODO/TODO-jodon-udp-binding-v0-reference-implementation.md` | whole file | Local UDP lineage TODO owner | move whole | `udp-feed` lineage | Move with lineage; later path rename records old name |
| `simulations/SIM-piloh-turns-149-208-recovery/protocols/udp-binding.d/TODO/TODO.md` | whole file | Local UDP lineage TODO index | move whole | `udp-feed` lineage | Move with lineage; later path rename records old name |
| `simulations/SIM-piloh-turns-149-208-recovery/protocols/udp-binding.d/manifest.json` | whole file | Local UDP lineage manifest | move whole | `udp-feed` lineage | Move with lineage; later path rename records old name |
| `simulations/SIM-piloh-turns-149-208-recovery/protocols/udp-binding.d/specs/udp-binding-draft.md` | whole file | Local UDP lineage draft spec | move whole | `udp-feed` lineage | Move with lineage; later path rename records old name |
| `simulations/SIM-piloh-turns-149-208-recovery/results/README.md` | whole file | Placeholder for umbrella-sim results | supersede whole | later successor-specific result artifact, if any | Do not create empty result logs by default in successor sims |
| `simulations/SIM-piloh-turns-149-208-recovery/seed/protocol-tree-migrations.md` | whole file | Mixed migration provenance for three protocol trees | split | see child rows below | Per-lineage provenance should separate rather than travel as one mixed file |
| `simulations/SIM-piloh-turns-149-208-recovery/seed/protocol-tree-migrations.md` | intro paragraphs | General explanation of protocol-tree migration method | stay rooted | rooted migration provenance (`TODO-rusis` or successor rooted note) | Cross-lineage migration method is harness memory |
| `simulations/SIM-piloh-turns-149-208-recovery/seed/protocol-tree-migrations.md` | `group-session` row | `group-session` migration provenance | move with lineage | `group-session` lineage provenance note | Keep with lineage |
| `simulations/SIM-piloh-turns-149-208-recovery/seed/protocol-tree-migrations.md` | `udp-binding` row | UDP migration provenance under old name | move with lineage | `udp-feed` lineage provenance note | Preserve old-name provenance on rename |
| `simulations/SIM-piloh-turns-149-208-recovery/seed/protocol-tree-migrations.md` | `ppx-dr` row | Retired `ppx-dr` migration provenance | archive | retired `ppx-dr` archive/provenance location | Preserve retirement provenance |
| `simulations/SIM-piloh-turns-149-208-recovery/seed/wire-lab-devs-draft-migration.md` | whole file | Migration provenance for concrete wire-lab-devs transport evidence | move whole | `wire-lab-devs` lineage provenance / evidence notes | Belongs with the concrete evidence lineage |
| `simulations/SIM-piloh-turns-149-208-recovery/world/README.md` | whole file | World-level intro for concrete specimen state | move whole | `wire-lab-devs` lineage | This world subtree is concrete evidence, not shared apparatus |
| `simulations/SIM-piloh-turns-149-208-recovery/world/cas/README.md` | whole file | Placeholder CAS world note | move whole | `wire-lab-devs` lineage | Move now; later kept or superseded once actual local CAS state exists |
| `simulations/SIM-piloh-turns-149-208-recovery/world/feeds/README.md` | whole file | Placeholder feed world note | move whole | `wire-lab-devs` lineage | Move now; later kept or superseded once actual local feed state exists |
| `simulations/SIM-piloh-turns-149-208-recovery/world/groups/README.md` | whole file | Placeholder group world note | move whole | `wire-lab-devs` lineage | Move now; later kept or superseded once actual local group state exists |
| `simulations/SIM-piloh-turns-149-208-recovery/world/sites/README.md` | whole file | Placeholder site world note | move whole | `wire-lab-devs` lineage | Move now; later kept or superseded once actual local site state exists |
| `simulations/SIM-piloh-turns-149-208-recovery/world/transports/wire-lab-devs-draft/README.md` | whole file | README for concrete wire-lab-devs transport evidence | move whole | `wire-lab-devs` lineage | Preserve as part of concrete evidence |
| `simulations/SIM-piloh-turns-149-208-recovery/world/transports/wire-lab-devs-draft/bafkreia46vxsahmeicugfxmc7natorkstc3mdaz4r5d3zz46whjwpvqwta.txt` | whole file | Message evidence | move whole | `wire-lab-devs` lineage | Must move byte-for-byte only |
| `simulations/SIM-piloh-turns-149-208-recovery/world/transports/wire-lab-devs-draft/bafkreidef4b4qdc4xjvkjrern7jm4ta75q55ed2u2ilwcrkxqhn7n4fjce.txt` | whole file | Message evidence | move whole | `wire-lab-devs` lineage | Must move byte-for-byte only |
| `simulations/SIM-piloh-turns-149-208-recovery/world/transports/wire-lab-devs-draft/bafkreihhuejiefrqrm7zgw2jsdqc37lwmbvfkw5uqbnjx3wsobcxh3y7ni.txt` | whole file | Message evidence | move whole | `wire-lab-devs` lineage | Must move byte-for-byte only |
| `simulations/SIM-piloh-turns-149-208-recovery/world/transports/wire-lab-devs-draft/bafkreihnonvsf3vmcagukqcxwoh35255eduulvwwx3kax6ty4iidklk5vu.txt` | whole file | Message evidence | move whole | `wire-lab-devs` lineage | Must move byte-for-byte only |
| `simulations/SIM-piloh-turns-149-208-recovery/world/wires/README.md` | whole file | Placeholder wires world note | move whole | `wire-lab-devs` lineage | Move now; later kept or superseded once actual local wire state exists |
| `protocols/wire-lab.d/specs/transport-spec-draft.md` | whole file | Rooted mixed artifact: thin outer feed contract + rooted apparatus/meta residue | split | see child rows below | Exact extraction waits for later `feed-outer` move pass |
| `protocols/wire-lab.d/specs/transport-spec-draft.md` | `## Purpose` + principles 1-4 + `## What this spec does NOT specify` | Candidate thin outer feed convention | supersede | later `feed-outer` lineage draft / provenance note | Extract into feed-outer lineage without silently changing meaning in this pass |
| `protocols/wire-lab.d/specs/transport-spec-draft.md` | `## Sources` + `## The per-axis meta-rule` + `## Open questions` + `## Freeze gate` | Rooted apparatus/meta governance around the old outer draft | stay rooted | rooted harness migration/governance ownership (`TODO-kugod` / `TODO-turog` / later owner) | Needs later cleanup and successor routing; do not move now |
| `protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md` | whole file | Rooted mixed owner TODO spanning group-session freeze and old outer transport-freeze framing | split | see child rows below | Historical DI log stays rooted; live ownership later splits |
| `protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md` | prior aliases / status / DI log | Rooted coordination memory | stay rooted | rooted harness TODO | Historical ownership record stays rooted |
| `protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md` | `T-GROUP-SESSION-FREEZE` thread body and related parked cleanup references | mixed live ownership | supersede | later split between `group-session` lineage owner and rooted feed/transport migration owner | Requires successor TODO routing in later pass |
| `protocols/wire-lab.d/TODO/TODO-duvuk-te-42-filename-cid-cascade-policy.md` | whole file | Rooted mixed owner TODO spanning message-level and outer-feed filename/CID policy concerns | split | see child rows below | Historical shell stays rooted; live policy ownership later splits |
| `protocols/wire-lab.d/TODO/TODO-duvuk-te-42-filename-cid-cascade-policy.md` | prior aliases / status / question log / DI log shell | Rooted coordination memory | stay rooted | rooted harness TODO | Keep rooted as historical owner record |
| `protocols/wire-lab.d/TODO/TODO-duvuk-te-42-filename-cid-cascade-policy.md` | `T-FILENAME-CID-CASCADE` thread body | mixed live ownership | supersede | later split between `group-session` lineage owner and `feed-outer` lineage owner | Exact split depends on later move pass |
| `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md` | whole file | Rooted TE-40 owner TODO with both continuing rooted work and references to successor lineage placement | split | see child rows below | Rooted owner stays open, but some residual rows will point outward |
| `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md` | prior aliases / status / DI log | Rooted apparatus/specimen governance | stay rooted | rooted harness TODO | Must remain rooted |
| `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md` | residual checklist rows about harness-spec / rooted apparatus cleanup | stay rooted | rooted harness TODO | These are still apparatus-side tasks |
| `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md` | residual rows that currently point at future `grid-envelope`, `feed-outer`, or moved specimen homes | supersede | later successor owners in the relevant lineage sims | Update these rows only after successor owners exist |
