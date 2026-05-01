# TODO 016 — proposals as transport messages — BLOCKED

Source: TE-29 (`docs/thought-experiments/TE-20260501-215027-protocols-as-simulated-repos-and-binding-layer.md`),
"Migrations triggered by this TE" section, item 3.

**Status: BLOCKED on four upstream prerequisites.** Filed as a stub
because the work is real and the integer slot deserves to be reserved,
but is not actionable until the prerequisites land. Do not start
without resolving each.

## Goal

Move `proposals/pending/ppx-dr-001-bootstrap/` content (the bootstrap
proposal and Steve's frozen contest reply) from its ad-hoc
`proposals/` location to the canonical TE-29 transport-message path:

```
transports/<bootstrap-transport-slug>/<bootstrap-binding-pCID>/<group-session-pCID>/<ppx-dr-pCID>/<message-id>.msg
```

This realizes the TE-29 insight that "a proposal is a message on the
ppx-dr transport, and a contest is a reply." Treating these files as
wire messages in a transport leaf directory makes the model uniform:
every wire artifact is a message in some transport.

## Prerequisites (all four must land first)

1. **TODO 014 — protocols-as-simulated-repos migration.**
   `protocols/` must exist as the canonical location for protocol
   specs before this TODO has anywhere to point its pCIDs at.

2. **A frozen UDP-binding (or other bootstrap-binding) pCID.**
   UDP-binding is currently a draft at
   `protocols/udp-binding.d/specs/udp-binding-draft.md`. The path in
   `transports/` requires a frozen binding pCID, not a draft slug.
   This depends on:
   - TODO 018 (UDP-binding v0 reference implementation) landing
   - The freeze ceremony itself being decided (OQ-29.1, currently
     deferred pending Steve's answers on doc-centric vs repo-centric
     and on spec-doc-self-sufficiency)

3. **A frozen group-session pCID.** Same logic. The session-protocol
   layer of the path requires a frozen pCID for group-session v0,
   which depends on the same freeze ceremony and on group-session
   having a reference implementation.

4. **A drafted-and-frozen ppx-dr message-protocol spec.** ppx-dr does
   not yet have a draft spec. The proposal/contest payload schema
   that the existing `proposals/` files implicitly embody must be
   written down explicitly, drafted, reviewed, and frozen before
   those files can be relocated under a frozen ppx-dr pCID.

Until all four prerequisites land, this TODO is not actionable.

## Subtasks (when unblocked)

1. **Decide the bootstrap-transport slug** at level 1 of the path.
   Real-world transport that carries the bootstrap proposal could be
   `udp`, `file-drop`, `git`, or a new bootstrap-specific slug.
   Likely `file-drop` since the proposal/contest pair currently lives
   as files in a git tree, and `file-drop` is the binding closest to
   that semantics. Worth a small TE if the choice is non-obvious.

2. **Rename `.md` to `.msg`** on each file. This is a one-time edit
   that affects byte-content (filename is not part of the bytes, but
   tools that key off extension may behave differently).

3. **DI-003 anchor migration.** DI-003 (don't-touch list)
   currently anchors on
   `proposals/pending/ppx-dr-001-bootstrap/contest-20260429-033208-steve-traugott.md`.
   The protection follows the bytes, not the path: the same content
   at a new path is still under DI-003. The DI-003 entry must be
   updated to reference the new path and filename in the same commit
   that does the move, not after, to avoid any window where the
   contest file is unprotected.

4. **Compute a `<message-id>` for each file.** Per TE-29 OQ-29.2 lean,
   the message-id is a content hash of the message bytes. Use the
   same hashing convention that the eventual ppx-dr spec mandates;
   if ppx-dr v0 picks a different hash, this work re-runs.

5. **Create `transports/<slug>/<binding-pCID>/<session-pCID>/<ppx-dr-pCID>/`**
   and move the (renamed) files into it. Verify byte-identity post-
   move: `git diff` of pre-move and post-move file contents must be
   empty (the path changed; the bytes did not).

6. **Update all in-tree references** to the old `proposals/` path:
   harness spec, TE prose, DR/DI references, README if any.

7. **Decide the fate of `proposals/`** as a directory. Likely
   removed entirely; if any not-yet-canonical-message content remains,
   it migrates to a more appropriate path or stays under a renamed
   parent.

## Out of scope for TODO 016

- Drafting the ppx-dr message-protocol spec (that is a prerequisite,
  to be filed as its own TODO when it becomes the next-actionable
  item).
- Designing the freeze ceremony (OQ-29.1; its own future TE).
- Generalizing this pattern to other future ad-hoc files outside
  `proposals/`. If new ad-hoc message files appear, each gets its
  own TODO under the same pattern.

## Done when

- `proposals/pending/ppx-dr-001-bootstrap/` is empty or removed.
- The two files (bootstrap proposal, Steve's contest reply) live at
  their canonical `transports/...` paths with byte-identical
  content.
- DI-003 anchors on the new paths and the don't-touch list is
  re-verified.
- All in-tree references resolve.
- `/tmp/spec check` (or its successor) reports OK.

## Why this is filed despite being blocked

Per OQ-100.4 (numbering wrap, stable integers across centuries) and
TE-25 (drafting-time anchoring of TE/TODO numbers), filing the integer
now keeps TODO numbering stable in fact, not just in promise. A
future contributor finding TODO 016 referenced from TE-29 prose
should land here and find a clear "blocked, why, what unblocks it,
what to do" record, rather than a missing file or an integer gap.

## Provenance

- TE-29 first-drafted at 2026-05-01 21:50:27 UTC.
- TODO 016 filed as BLOCKED stub at 2026-05-01 22:58 UTC (same
  session).
- Will be reactivated when all four prerequisites land.
