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
| B | TE-sihih (substrate-agnostic layered model) | TE-sihih decided; TODO-vunub closed | 52 |
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
| B — TE-sihih | "TE-sihih: Substrate-agnostic layered model (52 UTs)" | `protocols/wire-lab.d/TODO/TODO-vunub-te-38-substrate-agnostic-layered-model.md` | TE-sihih decided; TODO-vunub closed; successor migration/root-layout questions tracked separately |
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

## DF-V.4 lock — Alt B (per-bundle transfer to AGENTS-ppx.md) — 2026-05-07

This section is an additive amendment to the read-only pass and to the
DF-V.2 / DF-V.3 lock subsections above. It does not flip any UT
checkbox, does not edit any earlier section of this file, does not
touch `dropped-thread-disposition-20260506.md` bullets, and does not
touch TODO-lilar walk notes. Per AGENTS.md the locked decisions are
recorded as one DI per bundle in
`protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`
("Decision Intent Log"); this subsection is the docs-side reference
back to those DIs and the per-bundle transfer-status ledger
authorized by the DF-V.3 Alt C "matrix-as-closure-index" mechanism.

**Locked decision (DF-V.4 Alt B, "per-bundle transfer to
AGENTS-ppx.md"):** the J2 sub-class of the Carry cluster (Cluster J)
identified by DF-V.2 Alt C is transferred to `AGENTS-ppx.md` one
bundle at a time, with one DI in TODO-lilar per bundle. Seven bundles
were locked in this pass, in the priority order Steve approved:
B1 → B3 → B4 → B2 → B5 → B6 → B7. All seven bundles have a single
named destination: `AGENTS-ppx.md`, under the new top-level section
"Carry-J2 procedural discipline (durable cross-session rules)".

**Per-bundle transfer status ledger.** Per the DF-V.3 Alt C
mechanism, each bundle's substantive landing under its own DI is
what marks the bundle as transferred for matrix purposes. The
seven bundles below are recorded as transferred by this pass:

| Bundle | Title | Destination | DI in TODO-lilar | Status |
| --- | --- | --- | --- | --- |
| B1 | Foreground authorization (separate authorization from execution) | `AGENTS-ppx.md` § "Carry-J2 procedural discipline" → "B1 — Foreground authorization (separate authorization from execution)" | `DI-021-20260507-212249` | transferred |
| B3 | Collaborator anonymity / non-mention | `AGENTS-ppx.md` § "Carry-J2 procedural discipline" → "B3 — Collaborator anonymity / non-mention" | `DI-021-20260507-212250` | transferred |
| B4 | PAT redaction and credential hygiene | `AGENTS-ppx.md` § "Carry-J2 procedural discipline" → "B4 — PAT redaction and credential hygiene" | `DI-021-20260507-212251` | transferred |
| B2 | Foreground DONE confirmation | `AGENTS-ppx.md` § "Carry-J2 procedural discipline" → "B2 — Foreground DONE confirmation" | `DI-021-20260507-212252` | transferred |
| B5 | One-DF-at-a-time discipline | `AGENTS-ppx.md` § "Carry-J2 procedural discipline" → "B5 — One-DF-at-a-time discipline" | `DI-021-20260507-212253` | transferred |
| B6 | Apologize, audit, invalidate, propose | `AGENTS-ppx.md` § "Carry-J2 procedural discipline" → "B6 — Apologize, audit, invalidate, propose" | `DI-021-20260507-212254` | transferred |
| B7 | Ground-truthing before citation | `AGENTS-ppx.md` § "Carry-J2 procedural discipline" → "B7 — Ground-truthing before citation" | `DI-021-20260507-212255` | transferred |

**What this lock authorizes (only this; nothing more):**

- The seven AGENTS-ppx.md edits enumerated above (one new top-level
  section with seven subsections, one per bundle).
- The seven DI entries in TODO-lilar enumerated above.
- This subsection's own existence as the per-bundle transfer-status
  ledger of record, parallel in shape to the DF-V.2 and DF-V.3 lock
  subsections.

**What this lock does NOT authorize:**

- It does **not** flip any UT checkbox in TODO-lilar or in
  `dropped-thread-disposition-20260506.md`. Both files remain
  append-only history. The matrix-as-closure-index mechanism locked
  by `DI-021-20260507-210204` (DF-V.3 Alt C) governs closure; the
  per-bundle transfer status above is the index entry for these
  bundles, not a closure-flavored edit to either append-only file.
- It does **not** edit any UT bullet in
  `dropped-thread-disposition-20260506.md`, or any walk note in
  TODO-lilar. The disposition memo's bullet structure stays as
  written; TODO-lilar's only edits in this pass are the seven new
  DI entries.
- It does **not** edit, retitle, or recount the body of this matrix
  above this subsection. The 45-closeable / 141-not-closeable
  summary, the per-cluster breakdown, and the "Owner and pointer
  table" remain as recorded by the read-only Alt-V.4 pass and the
  DF-V.2 / DF-V.3 locks.
- It does **not** authorize per-J2-UT recount of the 55 Carry items
  into B1–B7 buckets. The DF-V.3 Alt C "matrix-as-closure-index"
  mechanism explicitly defers per-UT classification into each
  Carry-rule transfer; this DF-V.4 lock records bundle-level
  transfer status, not per-UT closure.
- It does **not** authorize edits to AGENTS.md, AGENTS-codex.md, or
  any per-protocol procedure file. The destination locked by Alt B
  is `AGENTS-ppx.md` for all seven bundles.
- It does **not** authorize back-pointers from TODO-lilar walk notes
  or from `dropped-thread-disposition-20260506.md` to this
  subsection. Both are append-only; the seven new DI entries in
  TODO-lilar's Decision Intent Log are the only landed edits to
  TODO-lilar by this pass.
- It does **not** supersede `DI-021-20260507-204144` (DF-V.2 Alt C)
  or `DI-021-20260507-210204` (DF-V.3 Alt C). The J1/J2 split and
  the matrix-as-closure-index mechanism remain in force.

**Outstanding DFs after this lock:**

- None exposed by this DF that are required for the seven bundles
  above to land. Per-J2-UT classification (which of the 55 Carry
  walk-note UTs each bundle covers) remains a future verification
  pass and is gated on the matrix-as-closure-index mechanism rather
  than on a new DF.

Provenance: locked by user 2026-05-07; recorded as
`DI-021-20260507-212249` (B1), `DI-021-20260507-212250` (B3),
`DI-021-20260507-212251` (B4), `DI-021-20260507-212252` (B2),
`DI-021-20260507-212253` (B5), `DI-021-20260507-212254` (B6), and
`DI-021-20260507-212255` (B7) in TODO-lilar.

## Refinements

### 2026-05-12 — TE-40 owner closure pointer

Cluster C / TE-40 remains counted as originally recorded above; this refinement
does not recount or edit the read-only matrix body. The owner artifact
`protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`
now records TE-40 residual closure under `DI-mugar`, including the turn-159
`UT-159.a`, `UT-159.b`, and `UT-159.d` rows and the previously verified
`UT-159.c` row. This is the Cat-3 / Cat-4 forward pointer authorized by the
matrix-as-closure-index rule.

### 2026-05-12 — Turn-160 TE-havib closure pointer

`docs/thought-experiments/TE-havib-apparatus-vs-specimen-carve-out.md` now
has a Cat-3 refinement preserving the historical "eight specimen-bearing items"
wording while pointing readers at the nine-item audit count and the `DI-mugar`
harness-spec sweep closure. That closes `UT-160.a`. The same TE's final status
records all seven DFs locked after the Alt-B re-presentation path, closing
`UT-160.d`. `UT-160.b` and `UT-160.c` remain covered by the earlier
TODO-lilok verification walk.

### 2026-05-12 — Turn-161 TE-havib closure pointer

`docs/thought-experiments/TE-havib-apparatus-vs-specimen-carve-out.md` now
also records the turn-161 nine-axis `promise-stack` vs `grid-pcid-payload`
comparison and the assertion-taxonomy examples as historical inputs to the
final DF-36.2 / OQ-36.6 retirement path. That closes `UT-161.b` and
`UT-161.c` as captured historical inputs. `UT-161.a` remains covered by the
earlier TODO-lilok verification walk.

### 2026-05-12 — Turn-162 TE-havib closure pointer

The TODO-lilok verification walk records `UT-162.a` as resolved by TE-havib's
`OQ-36.6` negative-resolution path and `UT-162.b` as resolved by the Alt-B
re-presentation of `DF-36.2`. TE-havib's current status and final `DF-36.2`
lock make both outcomes visible, so no separate turn-162 work remains.

### 2026-05-13 — Turn-163 DF-36.5 closure pointer

`UT-163.a` is resolved by the active `§1.3` apparatus-level layering-scenarios
rewrite in `protocols/wire-lab.d/specs/harness-spec-draft.md` under
`DI-lajod`, with the wider turn-159 harness sweep later closed by `DI-mugar`.
`UT-163.b` is closed for future-process purposes by `AGENTS-ppx.md` B1
(`DI-021-20260507-212249`); commit-specific residue remains with the named
commit UTs rather than with turn 163. `UT-cbf7f41-fallback` is retired for the
active `§1.3` path because OQ-36.6 resolved in the negative, DF-36.2 retired
promise-stack as a separate hypothesis, and the active harness prose now keeps
the layering scenarios at apparatus level without a specimen-side fallback
target.

### 2026-05-13 — Turn-164 routing pointer

Turn 164 initially routed successor work as follows. `UT-164.a` and
`UT-164.b` are historical corrections rather than live implementation work:
current TE-havib status corrects the lock count, and the TE-numan scope
question was later settled by the TE-39 / TODO-lilar cross-cutting
disposition. `UT-164.c` is resolved by `TODO-bisur` 012.7's four-message
round-trip, which exercises §3, §4, §4.6, §6, and §7 of the group-session
draft. The same-day closure pointer below supersedes the initial `UT-164.d`
and `UT-164.e` routing state.

### 2026-05-13 — Turn-164 successor closure pointer

This supersedes the routing pointer immediately above for turn-164 closure
state. `UT-164.d` is closed by sim-local `TODO-gapab` / `DI-rurab`, which
locks fixed configured `<author-id>/main` membership, passive-observer
non-membership, git-binding scope, and legacy `Message-ID:` reader/writer
behavior for the group-session specimen. `UT-164.e` is closed by `TODO-gapab`
/ `DI-rurab` plus `TODO-kakaz` / `DI-bomud`, which keeps spec freeze additive,
rejects historical transport-message rewrites, and keeps message-level CID and
header policy out of feed-outer. Turn 164 has no remaining open UT work.

### 2026-05-14 — Turn-165 closure pointer

Turn 165 has no remaining open UT work. `UT-165.a` is closed as an
observational privacy/slug lesson because the active specimen uses the generic
`wire-lab-devs` slug and active examples use neutral Alice/Bob prose.
`UT-165.b` is closed by `DI-rurab`, which defines the interim
`merge-group-transport-spec` shape as a Steve-authored DI until cryptographic
promise tooling exists. `UT-165.c` was already closed by the neutral-memory
update. `UT-165.d` is closed by keeping OQ-G4 deferred while treating m000 as
valid specimen evidence rather than a v0 genesis-message mandate. `UT-165.e`
is closed by the group-session example and freeze-gate cleanup under
`DI-rurab`.

### 2026-05-14 — Turn-166 closure pointer

Turn 166 has no remaining open UT work. `UT-166.a` is closed for
future-process purposes by the current decision-first protocol plus
`DI-vanak`'s explicit replay approval shorthand; the historical bootstrap
commits are preserved. `UT-166.b` is closed by `DI-rurab`: active membership is
the fixed configured set of exact `<author-id>/main` branches, so guessed
actors are not enrolled by speculation. `UT-166.c` is closed by active specimen
docs using `stevegt-via-perplexity` as the committed `From:` identity.
`UT-166.d` is closed as historical git metadata that must not be rewritten.
`UT-166.e` is closed by `DI-rurab`, which supersedes membership-by-posting with
fixed configured branch membership, passive-observer non-membership, and no
self-admission from unknown branches.

### 2026-05-14 — Turn-167 closure pointer

Turn 167 has no remaining open UT work. `UT-167.a` is closed by active
filename=CID docs: no active sequential `m<N>-...` filename rule remains.
`UT-167.b` and `UT-167.c` are closed by `DI-rurab`, which defines fixed
configured `<author-id>/main` branch membership and resolves `{name}` as
author-id for the git-bound specimen. `UT-167.d` is closed by the
2026-05-14 corpus audit for `\.msg\b`; remaining matches are historical
replay/disposition notes or message-body evidence, not active spec guidance.
`UT-167.e` is closed for future-process purposes by the current decision-first
protocol plus `DI-vanak`; later execute-on-directive rows keep their own
turn-specific reconciliation work.

### 2026-05-14 — Turn-168 closure pointer

Turn 168 has no remaining open UT work. `UT-168.a` is closed by active
`Message-ID:` compatibility wording in group-session §4.3 / §4.7. `UT-168.b`
is closed by `DI-rurab` and current §9 wording, which makes the git binding
normative for the wire-lab-devs specimen without making it the only possible
future binding. `UT-168.c` is closed by the passive-observer / configured-member
split in current §8 / §9.3. `UT-168.d` is closed by the active message-file
verification and infrastructure-file distinction in §2, §4, §9.3, and §9.5.
`UT-168.e` is closed for future-process purposes by decision-first plus
`DI-vanak`. `UT-168.f` is closed for active docs because the wire-lab-devs README
uses the post-turn-169 CIDs; turn-168-era CIDs remain historical-only, with the
broader rehash-continuity question retained under turn 169.

### 2026-05-15 — Turn-169 closure pointer

Turn 169 has no remaining open UT work. `UT-169.a` is closed by explicitly
recording the Path-A-vs-Path-B reasoning/action divergence and by the active
settled policy in `DI-012-20260508-033513` / `DI-rurab`: canonical writers omit
`Message-ID:`, while readers may tolerate exactly one legacy pre-`Date:` header
without giving it semantic identity. `UT-169.b` is closed by the explicit §4.7
legacy-header carve-out. `UT-169.c` is closed for future-process purposes by the
current decision-first stop rule plus `DI-vanak`. `UT-169.d` is closed as
historical branch metadata because active twig rules require `ppx/{twig}`
kebab-case, not `ppx/te-<utc>-<slug>`. `UT-169.e` is closed by the active
writer-prohibition / reader-tolerance split in §4.3.

### 2026-05-15 — Turn-170 closure pointer

Turn 170 has no remaining open UT work. `UT-170.a` is closed because the
original `DF-37.1` flat-versus-nested root `transports/` framing was superseded
by the later substrate/feed/CAS/site/simulation reframing: TE-sihih landed the
L5/L6/L7 model, TE-domat and DR-nugog reframed the tree question, and `DI-fakin`
implemented the current-specimen answer by moving the evidence into a simulation
world rather than choosing any root-level turn-170 alternative. `UT-170.b` is
retired as orphan DF-numbering history, `UT-170.c` is retired as stale
continuity-summary history for the abandoned `git-file-transport` plan, and
`UT-170.d` is retired because the substrate axis is now explicit in later owner
artifacts.

### 2026-05-16 — Turn-171 closure pointer

Turn 171 has no remaining open UT work. `UT-171.a` is closed by the combination
of TE-sihih / TODO-vunub Q-22.2, which makes L5 feeds first-class, and
`DI-rurab`, which keeps group-session §9 inline as the normative git binding for
the current wire-lab-devs specimen. `UT-171.b` is closed by TODO-vunub Q-22.3,
which retracts the per-instance manifest-field idea in favor of path-as-
declaration. `UT-171.c` is closed as a recorded design-cadence lesson about
YAGNI during active framing, and `UT-171.d` is closed as a recorded positive
pattern for two-part action-plus-framing questions.

### 2026-05-16 — Turn-172 closure pointer

Turn 172 has no remaining open UT work. `UT-172.a` is closed as a documented
framing-stability lesson: the old `git-file-transport` wording was neither a
transport-protocol name to keep nor a meaningless working label, but an early
misclassified hint of the later L5 feed axis. `UT-172.b` is closed by TE-sihih's
forward vocabulary (`feed`, with `substrate` as the descriptive prose term).
`UT-172.c` is closed by TODO-vunub Q-22.3 and TE-vipir path-as-declaration,
which retract the proposed `bindings/` / `messages/` instance layouts. `UT-172.d`
is closed by TE-sihih's L5/L6/L7 taxonomy replacing the turn-172
carrier/transport/binding sketch. `UT-172.e` is closed for current recovery by
`DI-rurab`, which keeps group-session §9 inline for the current specimen while
deferring any future feed-spec extraction to successor work.
