# Cat-1b handoff: known historical-quotation citations the sweep MUST NOT touch

**Read this before you let `sweep-citations` write to disk.**

**See `tools/README.md` for the full sweep-citations runbook**
(dry-run-first procedure, hazard notes, Cat-1b classification rule,
and recovery steps if you forget `-n`).

`sweep-citations` is a regex-based mechanical Cat-2 tool. It cannot
distinguish a *current pointer* (Cat-2 -- sweep) from a *historical
quotation* (Cat-1b -- leave) by looking at the surrounding text. The
six files listed below contain at least one Cat-1b citation that the
tool will repeatedly try to overwrite on every run. Each entry below
records:

- **path** -- file the tool flags
- **kept-as -- match** -- the literal token the tool would rewrite
- **why** -- the editorial reason the token must stay in its old
  surface form (Cat-1b classification per DI-020-20260502-232651)

When you (the LLM driving the sweep) inspect a `sweep-citations`
dry-run output, you must compare its flagged-files list against this
handoff. For every flagged file:

1. If the file is in this list, AND the dry-run's edit count matches
   the entry's expected count, AND the edit is on the listed line(s)
   touching the listed token, the file is a **known false positive**;
   do NOT let the tool write it. Skip the file (revert any write).
2. If the flagged file is NOT in this list, OR the edit count is
   higher than expected, OR the edit is on a different line, the
   file may have a **new Cat-1b case** the corpus has accumulated
   since this handoff was last revised. Hand-vet that file: read the
   diff, classify each match Cat-1a/Cat-1b/Cat-2 by hand, apply only
   the legitimate Cat-2 edits, and add a new entry to this handoff
   so future sweeps remember the verdict.

When in doubt, leave the token alone. The cost of leaving a stale
integer-alias citation in prose is nearly zero (readers can resolve
it via the `## Prior aliases` section or the `prior_alias` README
column). The cost of distorting a historical quotation is permanent
narrative damage.

---

## The six known false positives (as of 2026-05-07)

### 1. `AGENTS.md` -- 1 edit expected, line ~26

**match:** `(formerly TE-39)` inside the proquint-vocabulary teaching
sentence.

**why:** the sentence deliberately contrasts the new proquint handle
`TE-mumuv` against its prior integer alias `TE-39` so a reader new
to the convention understands which old reference resolves to which
new handle. Mechanical sweep would rewrite to `(formerly TE-mumuv)`,
collapsing the teaching contrast into a tautology
(`TE-mumuv (formerly TE-mumuv)`). Cat-1b: vocabulary-pedagogy
quotation.

### 2. `AGENTS-ppx.md` -- 1 edit expected, line ~218

**match:** `TODO 001` inside an indented (8-space) literal-example
commit-message body that the Onboarding step shows as a template.

**why:** the example demonstrates what a real commit message looked
like at the time the bot performed the bootstrap; the literal text
is the artifact under quotation, not a current pointer. Cat-1b:
quoted commit-message body.

### 3. `docs/thought-experiments/TE-titur-...` -- 6 edits expected

**matches:** `TE-25`, `TE-39`, `TE-25 § S5` (and similar
backward-citation examples), all inside the chunk-D Refinement
section that locks the integer-vs-proquint identity-role split.

**why:** TE-titur's Refinement section *is the document that explains
the proquint migration's identity rules*. It must use the old integer
aliases as literal subject-of-discussion. The `## Prior aliases`
section also lives here (and the tool already skips that section
correctly via `splitByPriorAliases`); the remaining matches are in
prose that quotes the integer aliases as the *thing being talked
about*. Cat-1b: identity-rule pedagogy.

### 4. `docs/thought-experiments/TE-vudaf-...` -- 1 edit expected, line ~97

**match:** `TODO 011` inside a quoted-historical-narrative paragraph
about the path-rename event in the channels branch.

**why:** the paragraph explicitly describes the *act of renaming* and
quotes the old path/integer literally as part of the historical
narrative. Mechanical sweep would distort the historical claim.
Cat-1b: historical-narrative quotation.

### 5. `protocols/wire-lab.d/TODO/TODO-bihon-...` -- 1 edit expected,
line ~10

**match:** `TE-29 OQ-29.9` inside a `Source:` provenance line.

**why:** the line is a provenance citation that anchors the TODO
file to the integer-aliased TE that originated it (TE-29 = TE-vipir).
The integer alias is the form that was canonical at the time of
authoring; rewriting it to TE-vipir would rewrite the provenance
record itself. Cat-1b: provenance citation.

(Note: the file's `## Prior aliases` section is correctly skipped by
the tool's section-aware splitter; this match is NOT in that section,
it's in the body's `Source:` line.)

### 6. `protocols/udp-binding.d/TODO/TODO-jodon-...` -- 3 edits expected

**matches:** `TE-29` (×2 in `Source:` and intro lines) and one inside
a parenthetical historical hedge ("Originally hedged as 'Likely TODO
020' before that slot was claimed ...").

**why:** the `Source:` lines are provenance citations (Cat-1b, same
as TODO-bihon). The hedge paragraphs were rewritten by hand in turn
297 to preserve the speculative-prediction history -- a mechanical
rewrite of "TODO 020" / "TODO 021" inside those hedges would
re-damage the just-corrected text. Cat-1b: provenance + historical
hedge.

---

## Procedure for adding a new entry

When a future sweep flags a file not in this list:

1. Hand-vet the diff.
2. If genuinely Cat-2, let the tool apply (and document nothing here
   -- the tool's job).
3. If Cat-1b, revert the tool's edit (or run with -n, never let it
   write), then add an entry above with: path, expected edit count,
   line number(s), the literal match token, and a one-paragraph
   "why" that future LLMs can verify against the live file.

When a documented Cat-1b case is *resolved* (the prose is rewritten
in a way that no longer matches the regex, or the file is deleted),
remove its entry. Keep this list short and current; a stale entry is
worse than no entry.

## Tool-side contract

`sweep-citations` reads this file at startup and prints its contents
verbatim to stderr before any walking begins. The tool does not
parse the entries; it relies on the LLM driving it to honor them.
The handoff exists at `tools/sweep-citations/CAT-1B-HANDOFF.md`
relative to the repo root specified by `-r`. If the file is missing,
the tool prints a warning and proceeds (never blocks).

## Provenance

Created 2026-05-07 in twig `ppx/sweep-cat1b-handoff` after Steve
ruled that mechanical Cat-1b protection should be a documented LLM
handoff, not a static skip-list inside the tool. The six entries
above are the empirically-derived false-positive set from the TE-39
follow-up sweep (turn 297, commit `f05ab1f`).
