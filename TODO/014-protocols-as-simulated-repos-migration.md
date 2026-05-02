# TODO 014 — protocols-as-simulated-repos migration

Source: TE-29 (`docs/thought-experiments/TE-20260501-215027-protocols-as-simulated-repos-and-binding-layer.md`).

Mechanical migration to the locked directory shape from TE-29. Done in
one cohesive change so that the wire-lab settles into its new shape
all at once rather than half-migrating.

## Subtasks

1. Create `protocols/wire-lab.d/` and move `specs/harness-spec-draft.md`
   into `protocols/wire-lab.d/specs/harness-spec-draft.md`. Update all
   in-tree references.

2. Create `protocols/group-session.d/` (new slug; renames
   `group-transport`) and move `specs/group-transport-draft.md` into
   `protocols/group-session.d/specs/group-session-draft.md`. Update
   all in-tree references. DFs T1.A through T6.A do not change; only
   the slug.

3. Decide the fate of `specs/transport-spec-draft.md`. Likely deleted:
   the thin outer rule it captured is now subsumed by per-binding
   specs under `protocols/<binding>.d/`. If any rules remain that
   apply to the wire-lab harness as a whole, fold them into the
   harness spec.

4. Move per-protocol TEs out of top-level
   `docs/thought-experiments/` into the appropriate
   `protocols/<slug>.d/docs/thought-experiments/`:
   - TE-24 (group-transport-envelope), TE-26, TE-27 -> `group-session.d/`
   - TE-28 (100-year goal) and TE-29 (this carve-out) -> `wire-lab.d/`
   - All others stay under wire-lab.d/ unless a clear per-protocol
     home exists.

5. Update `docs/thought-experiments/README.md` to be either a
   harness-level forwarding index (pointing into each protocol's
   `.d/`) or remove it entirely once content has migrated.

6. Update `tools/spec` to walk the new `protocols/<slug>.d/` layout
   and the per-protocol manifests.

7. Update `manifest.json` at the top level to be the registry of
   frozen `(slug, pcid, path, tree_hash, frozen_at)` tuples across
   all protocols.

8. Update `MANIFEST.md` accordingly (or replace with the JSON).

9. Run `/tmp/spec check` until green.

10. **(Added per TE-30) Recover canonical first-drafted timestamps for
    all existing TODOs.** For each file under top-level `TODO/`, find
    the git commit that introduced it and record its author-date as
    the canonical first-drafted-time in `YYYYMMDD-HHMMSS` (UTC) form.
    Where a file was introduced via merge from a long-ago feature
    branch, use the parent-commit date on the originating branch.
    Where unrecoverable, use the earliest plausible date from `git
    log --follow` for that file. Migration-time is never used.
    Output: a mapping of `<old-integer-filename>` ->
    `<new-timestamp-filename>`, recorded as a temporary file in the
    twig and discarded after step 11 lands.

11. **(Added per TE-30) Migrate every TODO into its destination
    `protocols/<slug>.d/TODO/` directory under its new timestamp
    filename.** Per the TE-30 triage table:

    - Harness-level TODOs (001 through 011, 013, 014, 015, 017, 019)
      move to `protocols/wire-lab.d/TODO/TODO-<timestamp>-<slug>.md`.
    - TODO 012 -> `protocols/group-session.d/TODO/`.
    - TODO 016 -> `protocols/ppx-dr.d/TODO/` (BLOCKED status
      preserved verbatim).
    - TODO 018 -> `protocols/udp-binding.d/TODO/`.

    Then:

    - Write each protocol's local `TODO.md` queue listing only that
      protocol's TODOs.
    - Write the master cross-listed `TODO.md` at
      `protocols/wire-lab.d/TODO/TODO.md` listing all TODOs across
      the wire-lab.
    - Delete the top-level `TODO/` directory entirely. No integer
      survivors at the top level.
    - Update all in-tree references to the old integer paths:
      DR-009, DI-009, TE-24, TE-29, TE-30, harness-spec section 8
      bibliography, prior-TODO cross-references inside other TODO
      files (e.g. 014's references to 015/017/018/019).
    - Update the original TODO 014 file's own location: it moves
      from `TODO/014-...md` to
      `protocols/wire-lab.d/TODO/TODO-<timestamp>-protocols-as-simulated-repos-migration.md`
      atomically with the rest of step 11.

## Out of scope for TODO 014

- proposals -> transports/ migration (that's TODO 016, BLOCKED).
- UDP-binding v0 reference implementation (that's TODO 018).
- ns-3 harness scaffold (that's TODO 019).

Note: prior versions of this TODO listed TODO 015 (DR/TODO/DI
absorption) and TODO 017 (group-transport -> group-session rename) as
out-of-scope items. Per TE-30, TODO 015 was retired (DR/DI directories
never existed; harness-spec sections 11/12 already absorb decisions
and open questions inline) and TODO 017's rename work was always
covered by step 2 of this TODO. See
`TODO/015-dr-todo-di-absorption-RETIRED.md` and
`TODO/017-group-transport-rename-FOLDED.md` for the records.

## Done when

- `protocols/` contains live `.d/` directories for wire-lab,
  group-session, ppx-dr, udp-binding, and any other protocols
  enumerated in TE-29.
- Each protocol's `.d/TODO/` directory exists with its protocol's
  TODOs and a local `TODO.md` queue.
- The master cross-listed `TODO.md` exists at
  `protocols/wire-lab.d/TODO/TODO.md` and lists every TODO across
  the wire-lab.
- The top-level `TODO/` directory is gone entirely (per TE-30; no
  integer-named TODO survivors).
- All current in-tree references resolve to new paths (TODO
  references, spec bibliography, TE prose, DR/DI anchors).
- `/tmp/spec check` reports OK.
- README at top level remains byte-identical to origin/main.
