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

## Out of scope for TODO 014

- DR/TODO/DI absorption into spec docs (that's TODO 015).
- proposals -> transports/ migration (that's TODO 016).
- group-transport -> group-session rename in code/comments (that's
  the spec rename only; ledger/test references handled in TODO 017).
- UDP-binding v0 reference implementation (that's TODO 018).

## Done when

- `protocols/` contains live `.d/` directories for wire-lab,
  group-session, ppx-dr, udp-binding, and any other protocols
  enumerated in TE-29.
- All current in-tree references resolve to new paths.
- `/tmp/spec check` reports OK.
- README at top level remains byte-identical to origin/main.
