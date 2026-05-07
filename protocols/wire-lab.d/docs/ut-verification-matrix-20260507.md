# UT verification matrix — TODO-lilar 186 UTs (read-only, index-only)

**Date:** 2026-05-07
**Author:** stevegt-via-perplexity (bot)
**Pass type:** read-only, index-only (Alt-V.4). No UT checkboxes edited.
**Trigger:** verification pass over the dropped-thread classification in
`protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md`,
following user approval of Alt-V.4 (index-only verification) and
DF-V.1 Alt-C (no checkbox edits this pass).

## Purpose and scope

This file is a verification cross-reference, not a disposition rewrite.
It classifies each of the 186 UT entries in TODO-lilar (also enumerated
verbatim in the disposition memo) into a closure-readiness state, and
records owners (TE / TODO / spec-edit / cluster) for each item.

**Hard rules carried forward from the read-only pass:**

- **TEs are analysis artifacts; they do not own work.** A TE narrows
  alternatives and feeds DF questions. Per AGENTS.md ("Decision-First
  Specification and Compliance Protocol") DI entries — recorded in a
  TODO file — are the records that capture the locked decision. A UT
  that is "transferred to a TE" is **not** resolved by that transfer;
  it is queued for analysis whose output is a DI in some TODO file.
- **TODOs own work.** Closure of an open question requires a DI entry,
  a spec edit, or an explicit retire/supersede note in a TODO file.
- **Cluster assignment is not resolution.** Sorting a UT into the
  TE-sihih cluster (or any other cluster) is bookkeeping, not closure.
- **Only resolved or retired UTs are closeable.** Transferred, blocked,
  and unclear UTs remain open regardless of where they have been
  filed.

**Pass discipline (Alt-V.4):**

- The 186 UT bullets in `dropped-thread-disposition-20260506.md` are
  preserved as written; no checkboxes flipped, no bullets edited.
- The walk-note `- [x]` markers in TODO-lilar (which mark *turn walks*,
  not UT closures) are also untouched.
- This file is the only artifact added by this pass; pointers from
  TODO-lilar or the disposition memo are deliberately omitted under the
  same "no checkbox edits" gate (the disposition memo's bullet text
  cannot be safely augmented without creating a closure-flavored edit;
  TODO-lilar is append-only history). Future readers reach this matrix
  by directory listing of `protocols/wire-lab.d/docs/`.

## State definitions

Each UT lands in exactly one state. Resolved / retired are closeable;
the rest are not.

- **Resolved.** A locked decision (DI / spec edit / merged commit) on
  `ppx/main` answers the UT's substantive question. Closure-ready.
- **Retired.** Superseded by other landed work or by a later vocabulary
  / scope correction; the UT is no longer actionable. Closure-ready.
- **Transferred.** Owned by a future TE, a future TODO subtask, or a
  spec-edit not yet drafted. The UT remains open; the cluster name is
  a bookkeeping label, not a resolution.
- **Blocked.** Owned by a TE / TODO that is itself blocked on prior
  work landing first. Open.
- **Unclear.** Either the UT's substantive question is ambiguous in
  the walk-note text, or the verifier cannot tell from current corpus
  state whether the UT is resolved, retired, or transferred. Open.

## Summary

| State | Count | Closeable? |
| --- | ---: | --- |
| Resolved  | 38 | yes |
| Retired   |  7 | yes |
| Transferred | 114 | no |
| Blocked   |  5 | no |
| Unclear   | 22 | no |
| **Total** | **186** | — |

**Closeable now: 45 (38 resolved + 7 retired).**
**Not closeable: 141 (114 transferred + 5 blocked + 22 unclear).**

This pass does not flip any checkbox. The 45 closeable UTs are
flagged for a future closure pass that selects a closure mechanism
(per-UT note in the disposition memo, group-close in a TODO, or
explicit retire in a successor TE). That mechanism is itself a DF and
is out of scope for the index-only pass.

## Per-cluster breakdown

The disposition memo groups the 186 UTs into ten clusters
(TE-havib follow-on; TE-sihih; TE-40; TE-41; TE-42; TE-43; TE-45;
Spec-edit; Retire; Carry). For matrix purposes those clusters are
relabeled A–J:

| Code | Cluster (per disposition memo) | Owner type | Items |
| --- | --- | --- | ---: |
| A | TE-havib follow-on (OQ-36.6 + tabletop) | TE-havib (already decided 2026-05-05) | 5 |
| B | TE-sihih (substrate-agnostic layered model) | TE-sihih (drafting; depends on TE-mumuv) | 52 |
| C | TE-40 (apparatus-vs-specimen completion) | TE-40 (drafting) | 18 |
| D | TE-41 (group-session freeze procedure) | TE-41 (depends on TE-40) | 15 |
| E | TE-42 (filename / CID-cascade policy) | TE-42 (depends on TE-41) | 7 |
| F | TE-43 (promisebase prior-art adoption) | TE-43 (drafting) | 25 |
| G | TE-45 (conditional-release / geofencing) | TE-45 (orthogonal) | 1 |
| H | Spec-edit (small direct edits, no TE) | direct edit | 5 |
| I | Retire (superseded; no action) | retire | 3 |
| J | Carry (procedural / AGENTS / cadence) | AGENTS-ppx.md / per-turn discipline | 55 |
| | | **Total** | **186** |

State per cluster:

| Cluster | Resolved | Retired | Transferred | Blocked | Unclear | Total |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| A — TE-havib follow-on | 5 |  0 |  0 | 0 |  0 |  5 |
| B — TE-sihih           | 6 |  3 | 30 | 1 | 12 | 52 |
| C — TE-40              | 0 |  0 | 14 | 4 |  0 | 18 |
| D — TE-41              | 0 |  0 | 15 | 0 |  0 | 15 |
| E — TE-42              | 0 |  0 |  7 | 0 |  0 |  7 |
| F — TE-43              | 2 |  0 | 23 | 0 |  0 | 25 |
| G — TE-45              | 0 |  0 |  1 | 0 |  0 |  1 |
| H — Spec-edit          | 0 |  1 |  4 | 0 |  0 |  5 |
| I — Retire             | 0 |  3 |  0 | 0 |  0 |  3 |
| J — Carry              | 25 |  0 | 20 | 0 | 10 | 55 |
| **Total**              | **38** | **7** | **114** | **5** | **22** | **186** |

## Caveats and contested counts

The cluster tallies above match the read-only pass and the disposition
memo's per-cluster headcount, but two counts are contested or
sensitive to a single semantic call. Verifiers should know which.

### Caveat 1 — Cluster J (Carry) "resolved" depends on
*captured-in-walk-note* counting as resolution

25 of the 55 Carry items are counted "resolved" in this matrix. Carry
items are procedural / AGENTS-rule / cadence observations whose
"resolution" is generally that the observation is itself the
resolution: the lesson lands in AGENTS-ppx.md, in a per-turn
discipline note, or simply in the walk-note's own text. Whether that
counts as resolved depends on a closure semantic that is **not yet
locked**:

- **Tight reading:** a Carry UT is resolved only if a corresponding
  AGENTS-ppx.md rule, DI entry, or repo-wide convention has actually
  landed. Under this reading some of the 25 may slip from resolved to
  transferred or unclear (the lesson exists only in the walk note,
  not yet codified).
- **Loose reading (used in this matrix):** a Carry UT is resolved if
  the walk note itself is the durable record — i.e., the walk note
  states the procedural lesson clearly enough that no follow-up TE or
  spec edit is needed, and the lesson does not require enforcement
  beyond the existing DF / TE / comment-preservation protocols.

The 38-resolved / 7-retired / 45-closeable tally in the summary uses
the **loose reading** for Cluster J consistency with the read-only
pass. A future closure pass should pick a closure semantic explicitly
(this is itself a DF) before any Carry UT is checked off. Under the
tight reading the closeable total drops below 45; the exact number is
not estimated here because per-Carry-UT walk-note review is out of
scope for an index-only pass.

This caveat affects only Cluster J. Clusters A–I use the same
resolved / retired definitions under both readings.

### Caveat 2 — Cluster A (TE-havib follow-on) "resolved" rests on the
2026-05-07 verification walk

All 5 UTs in Cluster A are marked resolved. The disposition memo
records them as **closed 2026-05-07** under the
`te-havib-followon-verification-walk-20260507.md` walk, which itself
verified that all six tabletop scenarios hold against the locked
state of TE-havib (DF-36.1 / .2 / .3 / .7 with OQ-36.6 resolved in
the negative). The resolution status here is contingent on that walk
remaining valid. If a future TE supersedes any of TE-havib's locks,
the affected UTs revert to transferred or unclear pending re-walk.

### Caveat 3 — Line-number citations are deliberately omitted

The disposition memo, TODO-lilar, and TE-havib are all editable
artifacts; line numbers shift. This matrix cites file paths and
section headings (e.g. "§ TE-sihih: Substrate-agnostic layered model
(52 UTs)" in `dropped-thread-disposition-20260506.md`) but does not
quote line numbers, to avoid stale citations. UT IDs (`UT-NNN.x`) are
the stable identifiers and are sufficient to locate each item.

## Owner and pointer table

The disposition memo is the canonical bullet-by-bullet ledger; this
matrix does not duplicate it. Per-cluster owners and onward pointers:

| Cluster | Disposition-memo section | Owner artifact (current) | Status of owner |
| --- | --- | --- | --- |
| A — TE-havib follow-on | "TE-havib follow-on: OQ-36.6 + tabletop walk (5 UTs)" | `docs/thought-experiments/TE-havib-apparatus-vs-specimen-carve-out.md` and verification walk `protocols/wire-lab.d/docs/te-havib-followon-verification-walk-20260507.md` | TE decided 2026-05-05; walk landed 2026-05-07 |
| B — TE-sihih | "TE-sihih: Substrate-agnostic layered model (52 UTs)" | `protocols/wire-lab.d/TODO/TODO-vunub-te-38-substrate-agnostic-layered-model.md` | TE drafting; depends on TE-mumuv |
| C — TE-40 | "TE-40: Apparatus-vs-specimen completion (18 UTs)" | `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md` | TE drafting |
| D — TE-41 | "TE-41: Group-session freeze procedure (15 UTs)" | `protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md` | depends on TE-40 |
| E — TE-42 | "TE-42: Filename/CID-cascade policy (7 UTs)" | `protocols/wire-lab.d/TODO/TODO-duvuk-te-42-filename-cid-cascade-policy.md` | depends on TE-41 |
| F — TE-43 | "TE-43: Promisebase prior-art adoption (25 UTs)" | `protocols/wire-lab.d/TODO/TODO-kituj-te-43-promisebase-prior-art-adoption.md` | TE drafting |
| G — TE-45 | "TE-45: Conditional-release / geofencing / recursive promise-graph (1 UTs)" | `protocols/wire-lab.d/TODO/TODO-ralud-te-45-conditional-release-geofencing.md` | TE drafting (orthogonal) |
| H — Spec-edit | "Spec-edit (small direct edits, no TE) (5 UTs)" | direct edits to `protocols/wire-lab.d/specs/*` and disposition memo | not started |
| I — Retire | "Retire (superseded; no action) (3 UTs)" | n/a (superseded by other landed work) | retired |
| J — Carry | "Carry (procedural / AGENTS-rule / cadence notes) (55 UTs)" | `AGENTS-ppx.md` and per-turn discipline | mixed: lessons recorded; codification varies — see Caveat 1 |

The owner artifacts above are TODOs and TEs, not the UTs themselves.
A UT moves from transferred → resolved only when its owner files a
DI entry (per AGENTS.md "Decision Intent Log") or lands a spec edit
that answers the substantive question. Cluster reassignment is
**not** sufficient.

## What this pass does not do

- Does **not** flip any UT checkbox in TODO-lilar.
- Does **not** edit any UT bullet in `dropped-thread-disposition-20260506.md`.
- Does **not** mint a new TE handle. (TEs are analysis artifacts; this
  is a verification artifact — analogous to the precedent at
  `protocols/wire-lab.d/docs/te-havib-followon-verification-walk-20260507.md`,
  which is not a TE either.)
- Does **not** introduce new closure semantics beyond Alt-V.4.
- Does **not** fix the contested Cluster J reading (Caveat 1) or the
  Carry tight/loose closure DF.
- Does **not** add a back-pointer from TODO-lilar or the disposition
  memo. Both files are append-only history; the disposition memo's
  bullet structure cannot be augmented without a closure-flavored
  edit, which Alt-V.4 forbids in this pass. Future readers find this
  matrix via directory listing of `protocols/wire-lab.d/docs/` or via
  search for "ut-verification-matrix".

## Inputs

- `protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`
  (the 186 UT walk notes, append-only history)
- `protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md`
  (the disposition memo; § per-cluster bullets are the verbatim UT
  inventory used here)
- `docs/thought-experiments/TE-havib-apparatus-vs-specimen-carve-out.md`
  (decided 2026-05-05; basis for Cluster A resolution)
- `protocols/wire-lab.d/docs/te-havib-followon-verification-walk-20260507.md`
  (precedent for verification-walk artifact shape and location;
  basis for Cluster A resolution)
- `AGENTS.md` Decision-First Specification and Compliance Protocol;
  Decision Intent Log conventions; TE Protocol; Comment Preservation
  Protocol (used to constrain the pass shape).

## Provenance and follow-on DFs

This pass is verification, not decision-making. No DI entry is filed
because no decision is locked here that was not already locked by
the user's Alt-V.4 + DF-V.1 Alt-C approval. Two DFs remain open and
are listed for the next pass:

- **DF-V.2 (Carry closure semantic).** Loose vs tight reading for
  Cluster J resolution (see Caveat 1). Affects whether 25 Carry items
  are closeable.
- **DF-V.3 (closure mechanism for the 45 closeable UTs).** Per-UT
  notes in the disposition memo, group-close in a successor TODO, or
  explicit retire-by-supersession entries in the relevant TEs.
  Out of scope for the index-only pass.

When DF-V.2 and DF-V.3 are answered, a follow-on closure pass can
flip the 45 closeable UT bullets in the disposition memo (and add the
matching entries to TODO-lilar's tail or to a new TODO) under whatever
mechanism the user picks.

## DF-V.2 lock — Alt C (split Carry-cluster semantics) — 2026-05-07

This section is an additive amendment to the read-only pass. It does
not flip any UT checkbox, does not edit any earlier section of this
file, does not touch `dropped-thread-disposition-20260506.md` bullets,
and does not touch TODO-lilar walk notes. Per AGENTS.md the locked
decision is recorded as `DI-021-20260507-204144` in
`protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`
("Decision Intent Log"); this subsection is the docs-side reference
back to that DI.

**Locked decision (DF-V.2 Alt C, "split Carry-cluster semantics"):**
The Carry cluster (cluster J, 55 UTs) does not have a single uniform
closure semantic. It contains two distinct sub-classes that must be
classified separately by any future closure pass:

1. **Sub-class J1 — in-turn self-corrections and purely historical
   observations.** UTs whose substantive content is fully captured by
   the walk note itself; the note is the durable record and no
   follow-up rule needs to land. These are **closeable** as resolved
   or retired under the loose reading of Caveat 1.
2. **Sub-class J2 — durable cross-session process rules.** UTs that
   record procedural lessons, AGENTS-style cadence rules,
   collaborator-anonymity discipline, foreground-DONE confirmation
   discipline, redaction discipline, twig-naming discipline, etc.
   These are **not closeable** by their walk-note record alone; they
   remain in the **transferred** state, owned by an explicit
   procedure home (typically `AGENTS-ppx.md` or the relevant
   per-protocol procedure file).

**What this lock authorizes (only this; nothing more):**

- A future closure pass may classify each of the 55 Carry UTs into
  J1 or J2 using the criteria above.
- A future `AGENTS-ppx.md` (or other named procedure-home) edit may,
  under a separate DF, transfer the J2 rules into that procedure
  home using whatever wording the per-rule transfer-list DF locks.

**What this lock does NOT authorize:**

- It does **not** flip any UT checkbox in TODO-lilar or in
  `dropped-thread-disposition-20260506.md`. Closure-mechanism choice
  remains DF-V.3 and is still open.
- It does **not** edit any UT bullet in
  `dropped-thread-disposition-20260506.md`. Those bullets are the
  append-only ledger and are out of scope for this DF.
- It does **not** edit any walk note in TODO-lilar. TODO-lilar is
  append-only history; the only edit landed by this DF in that file
  is the new DI entry under "Decision Intent Log".
- It does **not** retroactively alter the 25 resolved / 20 transferred
  / 10 unclear Carry counts recorded by the read-only Alt-V.4 pass.
  Those counts stand as the loose-reading snapshot. Any recount under
  the J1/J2 split is a verification pass and is gated on DF-V.3.
- It does **not** authorize promotion of any specific rule into
  `AGENTS-ppx.md`. Which J2 rules transfer, in what wording, with
  what enforcement, is a separate DF (not yet filed) that must lock
  the transfer list before any AGENTS-ppx edit is made.
- It does **not** infer that all Carry UTs are resolved just because
  they were noted; that is the failure mode this DF rules out.

**Caveat-1 update (additive):** Caveat 1 in the "Caveats and
contested counts" section above remains as written for historical
fidelity. The DF-V.2 Alt C lock supersedes the loose-vs-tight
either/or framing of that caveat with a per-UT split: J1 follows the
loose reading, J2 follows the tight reading. The recount is gated on
DF-V.3 and is not performed here.

**Outstanding DFs after this lock:**

- **DF-V.3 (closure mechanism).** Still open. The closure pass that
  applies the J1/J2 split and flips the 45 (loose-reading)
  closeable-flagged UT bullets is gated on this DF.
- **DF-V.4 (Carry-J2 transfer-list to AGENTS-ppx.md).** New DF
  exposed by DF-V.2 Alt C. Locks which J2 rules transfer to
  `AGENTS-ppx.md` (or another named procedure home), in what
  wording, with what enforcement and ownership. Not yet filed; the
  per-rule transfer is not authorized until this DF lands.

Provenance: locked by user 2026-05-07; recorded as
`DI-021-20260507-204144` in TODO-lilar.

## DF-V.3 lock — Alt C (matrix-as-closure-index) — 2026-05-07

This section is an additive amendment to the read-only pass and to the
DF-V.2 lock subsection above. It does not flip any UT checkbox, does
not edit any earlier section of this file, does not touch
`dropped-thread-disposition-20260506.md` bullets, and does not touch
TODO-lilar walk notes. Per AGENTS.md the locked decision is recorded
as `DI-021-20260507-210204` in
`protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`
("Decision Intent Log"); this subsection is the docs-side reference
back to that DI.

**Locked decision (DF-V.3 Alt C, "matrix-as-closure-index"):**
The closure mechanism for the 45 (loose-reading) closeable-flagged UTs
is: no UT checkbox flips now in any file; this matrix IS the closure
index of record; future cluster-owner TODOs (already enumerated in
"Owner and pointer table" above) cite this matrix when their
substantive work lands and close their own UTs as part of their own
DI / spec-edit / retire entry. No bulk closure pass is performed
against the disposition memo or against TODO-lilar.

**What this lock authorizes (only this; nothing more):**

- Future cluster-owner TODOs (TODO-vunub for Cluster B / TE-sihih,
  TODO-kugod for Cluster C / TE-40, TODO-turog for Cluster D / TE-41,
  TODO-duvuk for Cluster E / TE-42, TODO-kituj for Cluster F / TE-43,
  TODO-ralud for Cluster G / TE-45, plus future Spec-edit / Retire /
  Carry-J2-transfer TODOs) MAY, when they land their substantive work
  under their own DIs, mark the UTs that the matrix already attributes
  to their cluster as closed in that successor TODO's own ledger,
  citing back to this matrix and to the DI that lands the work.
- This matrix MAY accumulate Cat-3 / Cat-4 forward-pointer entries in
  a future `## Refinements` section if specific UT closures need
  navigational pointers; substantive recounts, classification edits,
  or summary-table changes remain out of scope and require a
  superseding TE per the TE editing policy.

**What this lock does NOT authorize:**

- It does **not** flip any UT checkbox in TODO-lilar or in
  `dropped-thread-disposition-20260506.md`. Both files remain
  append-only history.
- It does **not** edit any UT bullet in
  `dropped-thread-disposition-20260506.md`. The disposition memo's
  bullet structure is the append-only ledger and stays as written.
- It does **not** edit, retitle, or recount the body of this matrix.
  The 45-closeable / 141-not-closeable summary, the per-cluster
  breakdown, and the "Owner and pointer table" remain as recorded by
  the read-only Alt-V.4 pass.
- It does **not** authorize a single bulk closure pass parallel to
  the dispositions memo. Closure is bundled into each successor
  TODO's substantive landing; there is no separate closure-only
  artifact.
- It does **not** authorize cross-cluster closures. Each successor
  TODO closes only the UTs the matrix already attributes to its
  cluster (Clusters A–J).
- It does **not** supersede `DI-021-20260507-204144` (DF-V.2 Alt C).
  The J1/J2 Carry split still governs Cluster J classification at
  closure time; Alt C moves that classification into each Carry-rule
  transfer (or explicit no-transfer note) instead of a global recount
  pass.
- It does **not** address DF-V.4 (Carry-J2 transfer-list to
  `AGENTS-ppx.md`, exposed by DF-V.2 Alt C). Individual J2 UTs close
  only when their rule actually transfers under DF-V.4's eventual
  lock.

**Outstanding DFs after this lock:**

- **DF-V.4 (Carry-J2 transfer-list to `AGENTS-ppx.md`).** Still open;
  exposed by DF-V.2 Alt C. Locks which J2 rules transfer to
  `AGENTS-ppx.md` (or another named procedure home), in what wording,
  with what enforcement and ownership.

Provenance: locked by user 2026-05-07; recorded as
`DI-021-20260507-210204` in TODO-lilar.
