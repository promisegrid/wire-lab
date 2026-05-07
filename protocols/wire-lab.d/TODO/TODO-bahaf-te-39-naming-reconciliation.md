# TODO-bahaf: TE-mumuv naming reconciliation (proquint filenames)

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TODO-32` (integer alias; never registered in TODO.md because the master cross-list is itself rebuilt by this migration)
- `TODO-20260507-025635` (timestamp alias and pre-migration filename)

Parent TODO for TE-39 work on the wire-lab branch
`ppx/te-20260507-025627-naming-reconciliation`. TE-39 retires the
parallel namespace problem flagged in turn 261 (2026-05-05) where TE-25
introduced timestamped filenames but the corpus continued to use integer
aliases (TE-N, TODO-N) that risked collision under DT3-style insertions.

Replaces both the integer alias and the timestamp prefix with a single
deterministic mint-time identifier: a proquint-1 (5-character
pronounceable, Wilkerson 2009) allocated by `tools/mint-handle/` and
embedded in the filename. The filename itself is the registry; no
HANDLES.md, no central authority, consistent with the 100-year-goal
invariants from TE-28.

## Status

In progress. Twig: `ppx/te-20260507-025627-naming-reconciliation`.

## Foundational invariants (referenced, not redecided)

- **Wire-lab is design notes about a system that has the 100-year
  no-central-authority invariant** (Steve, 2026-05-06 turn 297). The
  invariant binds the protocol substrate; it does not bind wire-lab's
  own corpus shape. So wire-lab MAY use a per-fork registry-equivalent
  (the union of filenames in `docs/thought-experiments/` and
  `protocols/wire-lab.d/TODO/`) for collision-checking handle mints.

- **The 100-year-goal invariants** themselves remain unchanged (TE-dajot;
  see TODO-vunub for cite chain).

- **Parallel-namespace structural risk** (Steve, 2026-05-06 turn 296):
  "the timestamp filenames gain us nothing if we're still risking
  collision with integer numbers. we either need to stop using integer
  numbers, or we need to rename all of the TE files to match their
  integer numbers." TE-mumuv takes the first option, eliminating both the
  integer alias and the timestamp prefix in favor of proquint-1.

## Question log

Per the AGENTS-ppx Question-logging discipline, every question asked of
Steve in TE-mumuv work is logged here at the moment of asking and checked
off only after the resulting product is committed and pushed.

These questions were asked verbally during gameout in turns 297-309
(2026-05-06 PT) before the discipline-required logging was caught up.
They are backfilled here in the order asked, with the answers Steve
gave in chat.

- [x] **Q-23.1** DF-39.A: how to reconcile integer TE numbers with their
      timestamped filenames -- drop integers, rename files to integers,
      or some third option?
    opened: 2026-05-06 19:03 PT (turn 297)
    asked of: stevegt@t7a.org
    blocks: TE-mumuv scope
    alternatives: Alt-A (drop integers, keep timestamp filenames) /
      Alt-B (rename all TE files to integer-N filenames) /
      Alt-F (replace both with a third identifier scheme)
    recommendation: Alt-A
    answered: 2026-05-06 19:05 PT -- pursue Alt-F variants (acronyms,
      abbreviations, collision-resistant pronounceable strings)
    resolved: 2026-05-07 @ this commit
    product: this TODO file (commit pending)

- [x] **Q-23.2** DF-39.B: which Alt-F dictionary -- bip39, proquint,
      curated wordlist, heroku-style adjective-noun, military-concat?
    opened: 2026-05-06 19:08 PT (turn 299)
    asked of: stevegt@t7a.org
    blocks: TE-mumuv handle scheme
    alternatives: Alt-F.bip39 / Alt-F.proquint / Alt-F.curated /
      Alt-F.heroku / Alt-F.milconcat (military-style concat like
      COMSUBPAC) / Alt-F.specops (HAVE BLUE / PAVE LOW)
    recommendation: Alt-F.proquint -- only fully deterministic option;
      curated/milconcat/specops require maintained dictionaries or
      LLM-creative steps that fail the 100-year invariant
    answered: 2026-05-06 19:12 PT -- proquint
    resolved: 2026-05-07 @ this commit
    product: this TODO file (commit pending)

- [x] **Q-23.3** DF-39.C: proquint width -- single quint (5 chars,
      65,536 space) or double quint (10 chars + hyphen, 4 billion
      space)?
    opened: 2026-05-06 19:14 PT (turn 301)
    asked of: stevegt@t7a.org
    blocks: TE-mumuv handle width
    alternatives: pq1 (5ch fixed) / pq2 (11ch fixed)
    recommendation: pq1 -- comfortable for ~300 entries before birthday
      paradox forces retries; wire-lab corpus is 66 entries today and
      unlikely to exceed a few hundred over project lifetime; pq2 is
      conservative if we want to never think about it again
    answered: 2026-05-06 19:55 PT -- pq1
    resolved: 2026-05-07 @ this commit
    product: this TODO file plus tools/mint-handle/main.go (default
      width = 1; -w 2 available for future expansion)

- [x] **Q-23.4** DF-39.D: what does the mint tool hash to produce the
      proquint -- timestamp+slug from filename, last_git_commit_hash,
      time_ns, or something else?
    opened: 2026-05-06 19:17 PT (turn 303)
    asked of: stevegt@t7a.org
    blocks: TE-mumuv mint algorithm
    alternatives: hash(slug+ts) / hash(time_ns) /
      hash(last_commit, time_ns) / random with collision check
    recommendation: random with collision check -- once derivability is
      not required, the entropy source is operationally irrelevant; the
      check is what guarantees uniqueness
    answered: 2026-05-06 19:33 PT -- "we don't need to be able to
      re-derive a handle -- the goal is to simply avoid collisions"
    resolved: 2026-05-07 @ this commit
    product: tools/mint-handle/main.go (uses time.Now().UnixNano() as
      entropy source, retries on collision)

- [x] **Q-23.5** DF-39.E: what does the mint tool check against to
      detect collisions -- repo-wide grep, structured frontmatter scan,
      maintained registry file, or filename glob?
    opened: 2026-05-06 19:43 PT (turn 305)
    asked of: stevegt@t7a.org
    blocks: TE-mumuv mint implementation
    alternatives: repo-wide grep / structured frontmatter scan /
      registry file / filename glob
    recommendation: filename glob -- once the proquint replaces both
      integer and timestamp in the filename itself, the filesystem is
      the registry; no separate scanning needed; no false positives
      from body-text matches
    answered: 2026-05-06 19:46 PT -- replace both integer and timestamp
      with proquint in the filename
    resolved: 2026-05-07 @ this commit
    product: tools/mint-handle/corpus.go (scanCorpus globs scanDirs and
      matches handleFileRE against filenames only)

- [x] **Q-23.6** DF-39.F: keep the slug in the filename
      (`TE-vapoj-substrate-agnostic.md`) or drop it entirely
      (`TE-vapoj.md`)?
    opened: 2026-05-06 19:49 PT (turn 307)
    asked of: stevegt@t7a.org
    blocks: TE-mumuv filename format
    alternatives: keep slug / drop slug
    recommendation: keep slug -- filenames remain human-scannable in
      `ls`; matches prior pattern (TE-<timestamp>-<slug>.md) with only
      the timestamp position swapped for proquint; slug is informative
      not normative (it can change in Cat-3 edits without breaking the
      handle)
    answered: 2026-05-06 19:49 PT -- "yes" (lock recommendation)
    resolved: 2026-05-07 @ this commit
    product: tools/mint-handle/corpus.go handleFileRE accepts
      `(TE|TODO)-<proquint>-<slug>.md` shape

- [x] **Q-23.7** DF-39.G: tool implementation language -- Python,
      shell, Go, or Rust?
    opened: 2026-05-06 19:51 PT (turn 309)
    asked of: stevegt@t7a.org
    blocks: TE-mumuv tool deliverable
    alternatives: Python / shell / Go / Rust
    recommendation: Go -- matches Steve's primary preference and the
      existing tools/spec/ binary's language
    answered: 2026-05-06 19:51 PT -- Go
    resolved: 2026-05-07 @ this commit
    product: tools/mint-handle/{proquint.go, corpus.go, main.go,
      *_test.go, go.mod}

- [x] **Q-23.8** DF-39.H: bulk-rename invocation pattern -- shell loop
      calling the binary 66 times, or add a `-batch N` flag to mint
      multiple handles in one process?
    opened: 2026-05-06 19:55 PT (turn 311)
    asked of: stevegt@t7a.org
    blocks: TE-mumuv migration script
    alternatives: shell loop (66 invocations) / `-batch N` mode
    recommendation: shell loop -- simpler, exercises the same code
      path the eventual single-mint usage will hit, no special-case
      batch logic to test
    answered: 2026-05-06 19:55 PT -- "calling 66 times is fine"
    resolved: 2026-05-07 @ this commit
    product: migration script (commit pending) calls
      `tools/mint-handle/mint-handle` per file

- [x] **Q-23.9** DF-39.I: how to record the prior integer alias for
      historical traceability after the rename -- README index column,
      file-level frontmatter, NOTES file, or no record at all?
    opened: 2026-05-07 02:56 UTC (turn 313)
    asked of: stevegt@t7a.org
    blocks: README index format
    alternatives: README column "prior alias" / per-file frontmatter
      `prior_alias: TE-sihih` / NOTES.md migration table / no record
    recommendation: README column "prior alias" -- single-source-of-
      truth, sortable, easy to grep for stale citations during the
      Cat-2 sweep, can be dropped later when the migration is far in
      the rear-view mirror
    answered: 2026-05-07 03:02 UTC -- "each file should include 'was
      TE-NN' in its frontmatter or text. maybe do both?" -> do BOTH:
      README column AND per-file record.
    resolved: 2026-05-07 @ this commit
    product: README index gains a `prior_alias` column AND each
      renamed file gains a `## Prior aliases` section (see Q-23.10)

- [x] **Q-23.10** DF-39.J: per-file recording mechanism -- TE files
      have no YAML frontmatter; they have a prose `## TE ID` section.
      Where exactly does "was TE-NN; was TE-<timestamp>" go?
    opened: 2026-05-07 03:03 UTC (turn 313, follow-up to Q-23.9)
    asked of: stevegt@t7a.org
    blocks: file-rename script
    alternatives: extend `## TE ID` block / new `## Prior aliases`
      section / inline in H1 title / README only (skip per-file)
    recommendation: extend `## TE ID` -- one canonical block, mirrors
      existing convention, smallest structural change
    answered: 2026-05-07 03:04 UTC -- "New '## Prior aliases'
      section"
    resolved: 2026-05-07 @ this commit
    product: rename script appends a new `## Prior aliases` section
      (after `## TE ID`) carrying both the integer alias (e.g. TE-sihih)
      and the timestamp alias (e.g. TE-sihih) for each
      renamed file.

## Locked decisions (from Q-23.1 through Q-23.10 above)

1. **Handle scheme:** proquint-1 (5 chars, Wilkerson 2009 alphabet),
   default for all new TE/TODO files.

2. **Filename format:** `(TE|TODO)-<proquint>-<slug>.md`. Proquint
   replaces both the integer alias and the timestamp prefix.

3. **Mint algorithm:** mint-time-allocated random label.
   `tools/mint-handle/` draws entropy from `time.Now().UnixNano()`,
   folds through SHA-256, encodes the first 2 bytes as proquint-1,
   checks against the corpus filename set, retries on collision.

4. **Registry:** the union of filenames in `docs/thought-experiments/`
   and `protocols/wire-lab.d/TODO/` IS the registry. No HANDLES.md.
   No frontmatter `handle:` field. The filename is canonical.

5. **Mint date provenance:** retained as a column in the README index
   (extracted from git log of the introducing commit, not from the
   filename, since the filename no longer encodes time). Per-file
   frontmatter MAY also carry a `mint_date:` field for offline
   readability, but the README index is canonical.

6. **Cross-fork behavior:** forks diverge naturally. Two forks may
   independently mint `vapoj` for different TEs after the fork point;
   this is a feature, not a bug, consistent with PromiseGrid's
   no-central-authority stance applied recursively to its own design
   corpus. (See Steve, 2026-05-06 turn 297.)

7. **Saturation threshold:** if `len(corpus) > 32_000` (50% of pq1
   space) we revisit and likely upgrade new mints to proquint-2. Not a
   concern for any plausible wire-lab evolution.

8. **Prior-alias recording:** dual record. (a) README index gains a
   `prior_alias` column carrying the integer alias and the timestamp.
   (b) Each renamed file body gains a new `## Prior aliases` section
   immediately after `## TE ID`, listing the integer (TE-NN) and the
   timestamp (TE-YYYYMMDD-HHMMSS) it carried before this commit.
   Rationale: README is the grep target for migration sweeps;
   per-file section keeps each file self-describing for offline
   readers (Q-23.9, Q-23.10).

## Migration scope

Bounded but corpus-wide. Touches:

- 38 TE files (rename + body integer-alias replacement)
- 28 TODO files (rename + body integer-alias replacement)
- `docs/thought-experiments/README.md` (index format change)
- `protocols/wire-lab.d/TODO/TODO.md` (master cross-list format change)
- `AGENTS.md` and `AGENTS-ppx.md` (TE-editing-policy and convention
  updates referencing integer aliases)
- `protocols/wire-lab.d/specs/harness-spec-draft.md` (TE-N citations)
- DR-009 / DI-009 / DI-011 / DI-020 and any other DR/DI files citing
  TE-N or TODO-N
- `protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md`
  (cites TE-N)
- Pre-migration audit reports under `protocols/wire-lab.d/TODO/`
  (kept verbatim; they are historical record)

Estimated edit count: ~200 individual citation rewrites plus 66 file
renames.

## Backward-compatibility note

Pre-migration filenames in session-logs (under wire-lab-logs repo)
remain verbatim. The integer-alias-to-proquint mapping lives in two
places (Q-23.9, Q-23.10):

- `docs/thought-experiments/README.md` `prior_alias` column (and the
  master TODO list equivalent) -- the grep target.
- A `## Prior aliases` section in each renamed file body -- the
  offline / self-describing record.

Future readers of historical session logs can cross-reference either.

## Driving TE

TE-mumuv is the driving TE for this work. Filename will be
`TE-<proquint>-naming-reconciliation.md`, minted by tools/mint-handle/
during this twig's drafting commit.

## Cross-references

- Predecessor: TE-titur (taluj) "te-numbering-collision-and-harness-spec-
  path" (2026-04-30) introduced the dual-identifier scheme that this
  TE retires. A Cat-3 Refinement on TE-titur records the proquint
  adoption.
- Foundation: TE-dajot (fivas) "100-year-goal-as-design-constraint" --
  cited for the no-central-authority invariant that wire-lab itself
  is design notes ABOUT (not bound by).
- Vocabulary touch: AGENTS.md and AGENTS-ppx.md gain a "handle"
  vocabulary line and lose any "integer alias" guidance.
