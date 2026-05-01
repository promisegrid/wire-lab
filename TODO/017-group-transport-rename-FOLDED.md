# TODO 017 — group-transport -> group-session rename — FOLDED INTO TODO 014

Source: TE-29 (`docs/thought-experiments/TE-20260501-215027-protocols-as-simulated-repos-and-binding-layer.md`),
"Migrations triggered by this TE" section, item 4.

**Status: FOLDED INTO TODO 014 BEFORE FILING.** This TODO was named
in TE-29 prose but, on review immediately after TE-29 landed, found
to be already covered by TODO 014 step 2. This file exists so the
integer slot 017 is not orphaned and so a future contributor finds an
explanation rather than a gap.

## Original intent (from TE-29)

> Pure renaming. DFs T1-T6 unchanged.

The premise was that `group-transport` (the v0 carve-out from TODO
013) needed to be renamed to `group-session` because its DFs T1-T6
(Parents header, ack-in-body, flat per-leaf layout) are session-layer
semantics, not L4 framing. TE-29's binding/session/message layer
decomposition makes the slug `group-transport` actively misleading.

## Why folded into TODO 014

TODO 014 step 2 already says:

> Create `protocols/group-session.d/` (new slug; renames
> `group-transport`) and move `specs/group-transport-draft.md` into
> `protocols/group-session.d/specs/group-session-draft.md`. Update
> all in-tree references. DFs T1.A through T6.A do not change; only
> the slug.

That **is** the rename. There is no separate work for TODO 017 to do;
filing it as its own TODO would create two records of one change and
risk drift between them.

## What residual work is not in TODO 014 (and where it goes instead)

TODO 014 step 2 covers the spec file rename and in-tree reference
updates. The following adjacent cleanups are explicitly *not* tracked
under TODO 017:

1. **TE-24's filename retains the `group-transport-envelope` slug.**
   Per the standing rule that TE filenames are not renamed when
   content is later refined (the timestamp pins origin, not last-
   edited; see `docs/thought-experiments/README.md` line 11), the
   filename
   `TE-20260430-204108-group-transport-envelope.md` stays as is.
   The TE's *content* may use updated vocabulary; its filename does
   not change. This is correct, not a defect. If TE-24's filename
   bothers a future contributor, the answer is to read the
   timestamp-anchoring rule, not to rename the file.

2. **Prose references to "group-transport" in TE-26, TE-27, TE-28,
   harness-spec section 8, and DR-009.** These are mostly historical
   or etymological ("the group-transport-envelope work landed under
   TODO 013, and was later renamed to group-session per TE-29").
   Updating them to "group-session" would erase context that helps
   future readers understand the project's history. Keep these as
   is unless a specific reference is actively confusing.

3. **`tools/spec` and any test code** that refers to "group-transport"
   as an identifier. If such code exists at the time TODO 014 lands,
   it is updated under TODO 014 step 6 (`/tmp/spec` walks the new
   `protocols/<slug>.d/` layout). No separate TODO needed.

If, after TODO 014 lands, there is a meaningful residual prose-
vocabulary cleanup that a contributor wants to track separately, file
a new TODO with a fresh integer; do not reuse 017.

## Why this stub exists

Per OQ-100.4 (numbering wrap, stable integers across centuries) and
TE-25 (drafting-time anchoring of TE/TODO numbers), letting an
orphaned integer sit creates exactly the slug-drift / numbering-
confusion problem we want to avoid. A future contributor finding
TODO 017 referenced from TE-29 prose should land here and find a
clear "folded into 014, here's why" record, rather than a missing
file or an integer gap.

## Provenance

- TE-29 first-drafted at 2026-05-01 21:50:27 UTC.
- TODO 017 folding decided 2026-05-01 22:58 UTC (same session).
- The TE-29 prose at "Migrations triggered by this TE" item 4 has
  been annotated to point here.
