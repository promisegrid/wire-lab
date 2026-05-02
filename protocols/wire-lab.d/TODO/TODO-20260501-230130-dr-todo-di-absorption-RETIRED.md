# TODO 015 — DR/TODO/DI absorption — RETIRED BEFORE FILING

Source: TE-29 (`docs/thought-experiments/TE-20260501-215027-protocols-as-simulated-repos-and-binding-layer.md`),
"Migrations triggered by this TE" section, item 2.

**Status: RETIRED.** This TODO was named in TE-29 prose but, on
review immediately after TE-29 landed, found to be substantially
unnecessary. This file exists so the integer slot 015 is not orphaned
and so a future contributor finds an explanation rather than a gap.

## Original intent (from TE-29)

> Inline current `DR/`, `TODO/`, `DI/` content as sections inside the
> relevant protocol spec doc. Per-protocol DRs/TODOs go to that
> protocol's spec; harness-level go to wire-lab harness spec.

The premise was that three top-level directories (`DR/`, `TODO/`,
`DI/`) needed to be absorbed into protocol spec docs to honor TE-29's
locked shape, which puts decisions, open work, and don't-touch
invariants inside the relevant protocol's spec.

## Why retired

On checking the actual repo state immediately after TE-29 landed:

1. **`DR/` does not exist as a top-level directory.** Decision-Record
   identifiers (DR-001, DR-009, etc.) live as filename prefixes
   inside `proposals/`, `docs/thought-experiments/`, and as inline
   prose in `specs/harness-spec-draft.md`. There is nothing to
   relocate.

2. **`DI/` does not exist as a top-level directory.** Don't-touch
   invariants (DI-001, DI-003, etc.) live as filename prefixes inside
   `proposals/` and as inline references in the harness spec. Same
   conclusion: nothing to relocate.

3. **`TODO/` does exist and is load-bearing in its current form.** It
   is harness-level meta-process — how the bot and Steve coordinate
   work session by session — not protocol content. Steve and the bot
   read `TODO/TODO.md` every session as the active work queue.
   Inlining it into the harness spec would either duplicate the
   queue or destroy it. Neither is desirable.

4. **The harness spec already absorbs decisions and open questions
   inline.** Section 11 (Decisions) and section 12 (Open Questions)
   of `specs/harness-spec-draft.md` already serve the role TODO 015
   would have created. They were authored before TE-29 and continue
   to work after TE-29.

Net: the original concern was real (top-level directories should not
mix harness meta-process with protocol content) but the
implementation requirement evaporated once we noticed the directories
in question mostly don't exist, and the one that does (`TODO/`) is
correctly placed as harness-level meta-process anyway.

## What was kept from this concern

- TE-29 step 4 of TODO 014 already moves per-protocol TEs out of the
  top-level `docs/thought-experiments/` into the appropriate
  `protocols/<slug>.d/docs/thought-experiments/`. That captures the
  per-protocol-content-belongs-with-its-protocol intuition that
  motivated TODO 015 in the first place.
- TODO/TODO.md remains at the top level as harness meta-process,
  intentionally. If a future TE re-opens this question (e.g.,
  per-protocol active TODOs become a real need), a new TODO with a
  new integer should be filed; do not reuse 015.

## Provenance

- TE-29 first-drafted at 2026-05-01 21:50:27 UTC.
- TODO 015 retirement decided 2026-05-01 22:58 UTC (same session).
- The TE-29 prose at "Migrations triggered by this TE" item 2 has
  been annotated to point here.
