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

12. **(Added per TE-32) Create top-level `implementations/` directory
    with stub `README.md`** explaining the spec-side vs
    implementation-side split: `protocols/<slug>.d/` is design (TEs,
    draft specs, design TODOs); `implementations/<impl-name>/` is
    code, test vectors, and conformance fixtures. Each implementation
    is its own free-standing tree with internal shape chosen by the
    implementer; the only universal requirement is a `CHANGELOG.md`
    at the implementation root recording conformance claims against
    upstream spec-doc-CIDs. Multiple implementations of the same
    protocol coexist as siblings; one implementation may implement
    many protocols. Empty otherwise; no implementations exist yet at
    migration time. README points at TE-31 (inversion) and TE-32
    (split) for the rationale.

13. **(Added per TE-32) Retarget TODO 018 (UDP-binding v0 reference
    implementation)** to live under
    `implementations/go-udp-binding-reference/` rather than any path
    under `protocols/udp-binding.d/`. The protocol's spec-side TODOs
    stay in `protocols/udp-binding.d/TODO/`; the implementation work
    is a separate concern that does not migrate with the protocol's
    design. Update the prose of the migrated TODO-018 file to
    reflect this target path.

14. **(Added per TE-32) Retarget TODO 019 (ns-3 harness scaffold)**
    similarly: target `implementations/ns3-harness-fixture/` (or
    whatever final slug is chosen at execution time). The wire-lab
    harness reference impl is a B-side artifact, not a spec-side
    one. Per TE-32 OQ-32.4, the harness has both a design tree
    (already at `protocols/wire-lab.d/`) and a reference
    implementation (at `implementations/wire-lab-harness-reference/`
    or similar). TODO 019 belongs to the latter. Update the migrated
    TODO-019 file's prose accordingly.

15. **(Added per TE-32) Seed empty `CHANGELOG.md` stubs in each
    migrated `protocols/<slug>.d/`** with a placeholder explanatory
    header block describing the A-side semantics (`event: freeze`
    entries authored by spec maintainers; doc-CID published when
    the spec is frozen). Real `freeze` entries get added when the
    first frozen sibling appears for that protocol. The B-side
    `implementations/<impl-name>/CHANGELOG.md` files are not seeded
    at migration time because no implementations exist yet.

## Out of scope for TODO 014

- proposals -> transports/ migration (that's TODO 016, BLOCKED).
- UDP-binding v0 reference implementation (that's TODO 018; this
  TODO retargets its destination per step 13 but does not implement
  it).
- ns-3 harness scaffold (that's TODO 019; this TODO retargets its
  destination per step 14 but does not implement it).
- Defining the precise CHANGELOG.md format (header block syntax
  YAML/fenced/HTML-comment, schema for entries beyond what TE-31
  and TE-32 sketch). Forthcoming TODO 020 will spec the format and
  build a parser; this TODO only seeds empty stubs.

Note: prior versions of this TODO listed TODO 015 (DR/TODO/DI
absorption) and TODO 017 (group-transport -> group-session rename) as
out-of-scope items. Per TE-30, TODO 015 was retired (DR/DI directories
never existed; harness-spec sections 11/12 already absorb decisions
and open questions inline) and TODO 017's rename work was always
covered by step 2 of this TODO. See
`protocols/wire-lab.d/TODO/TODO-20260501-230130-dr-todo-di-absorption-RETIRED.md` and
`protocols/wire-lab.d/TODO/TODO-20260501-230132-group-transport-rename-FOLDED.md` for the records.

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
- Top-level `implementations/` exists with stub `README.md` (per
  TE-32) explaining the A/B split.
- Each `protocols/<slug>.d/` has an empty `CHANGELOG.md` stub (per
  TE-32) ready to record `freeze` events.
- TODO 018 and TODO 019 (after migration into
  `protocols/<slug>.d/TODO/`) name `implementations/...` paths as
  their target rather than `protocols/<slug>.d/...` paths.
- All current in-tree references resolve to new paths (TODO
  references, spec bibliography, TE prose, DR/DI anchors).
- `/tmp/spec check` reports OK.
- README at top level remains byte-identical to origin/main.
