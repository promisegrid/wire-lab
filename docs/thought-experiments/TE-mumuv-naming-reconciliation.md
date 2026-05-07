# TE-mumuv: Naming reconciliation -- proquint handles for TEs and TODOs

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-mumuv

(Minted 2026-05-07 by `tools/mint-handle/` from `time.Now().UnixNano()`
seed, SHA-256-folded to two bytes, encoded as a Wilkerson 2009
proquint-1. Collision-checked against the corpus filename glob; no
retries needed at corpus size 70.)

## Status

decided

## Decision under test

The wire-lab corpus has carried two parallel identifier namespaces for
every TE and TODO file:

- **Integer aliases** (TE-famar through TE-sihih, TODO-dutaz through TODO-sinuv).
  Stable, terse, easy to cite in conversation and in commit
  messages. Allocated by hand at draft time. The "primary key" used
  in the README index, the master TODO cross-list, and most prose
  citations.
- **Timestamp prefixes** (TE-famar-...) embedded in the
  filename. Introduced by TE-titur ("TE-nibar numbering collision and
  harness-spec path", 2026-04-30) to defuse a real collision between
  two TE-nibar candidates that arrived from independent twigs. Sortable
  by drafting time, immune to re-numbering, but verbose and not
  citable in prose without copying a 16-character substring.

The two namespaces drifted apart immediately. New TEs and TODOs got
both an integer (assigned during drafting and recorded in the README)
and a timestamp (chosen at file-creation time and embedded in the
filename). Citations in TE bodies, AGENTS files, harness-spec, DR/DI
files, and session logs used integers almost exclusively; the
timestamps showed up only in `git mv` operations and in the README's
"First drafted" column.

The collision risk that TE-titur was trying to prevent did not actually
go away. It just moved up one layer: now the integer alias is the
thing that two parallel twigs can collide on, because nothing in the
file or its filename forces a draft-time allocator to pick a unique
integer. TE-titur's resolution depended on the bot remembering to look
at the README every time it picked a number, which is a
precondition that no twig running in isolation can satisfy.

This TE retires the dual-namespace scheme and replaces both the
integer alias and the timestamp prefix with a single mint-time
allocated random label embedded in the filename: a Wilkerson 2009
proquint-1. The filename glob across `docs/thought-experiments/` and
`protocols/*/TODO/` IS the registry; collisions are detected by
filesystem read at mint time and resolved by retry.

## Why this matters (the 100-year-goal lens)

TE-dajot (jamod, "100-year-goal as a load-bearing design constraint")
locks the invariant that the wire protocol must work without any
central authority for at least a century. Wire-lab is design notes
ABOUT a system that has those constraints; wire-lab itself is a
shared-fork repo with a mob-consensus workflow, so it can use any
local convention it likes for its own corpus management.

The dual-namespace problem in this TE is a worked example of what
goes wrong when a system uses centrally-allocated integer
identifiers across a fork-and-merge graph. The proquint resolution
is also a worked example of the design pattern wire-lab is
ultimately documenting: collision-checked random allocation against
a CAS-shaped registry, with the registry being the artifact set
itself rather than a separate authority.

Future wire-lab readers studying the binding/feed-layer design (TE-vipir
gitap) will see the same pattern in this TE's resolution: identifiers
are allocated by minting against a local view of the corpus, and the
collision behavior under fork is the same as it will be under
binding-layer divergence at the protocol level.

## Gameout (full transcript, condensed)

The discussion that produced this TE ran across continuations 19-21
of session 2026-05-06. The summary below preserves the alternatives
considered, the order in which they were raised, and the reason each
was set aside.

### Step 1 -- the dual-namespace problem surfaces

Steve: "we never finished reconciling the integer TE numbers with
their timestamped filenames."

Bot's first reading: the integer-and-timestamp redundancy. Steve's
correction: the actual concern is the collision risk between the two
namespaces -- specifically, the integer-alias namespace collides
across forks because nothing in the file or filename forces uniqueness.
TE-titur had patched the symptom by embedding a sortable timestamp; it
had not addressed the underlying need for a fork-safe identifier.

### Step 2 -- gameout of dictionary-based handles (Alt-F)

Bot proposed several Alt-F variants and surveyed dictionaries:

- **bip39** (2048 words). Pronounceable English; familiar to crypto
  users. Too long for a single-word handle (average 5-7 chars but
  high collision rate per single word; pairs/triples are required to
  hit a useful keyspace).
- **proquint** (Wilkerson 2009). 16 consonants, 4 vowels, CVCVC
  shape. Single proquint-1 = 5 chars, 65 536 distinct values.
  Pronounceable but topic-opaque ("vapoj" tells you nothing about
  what the file is about).
- **curated word list** (handcrafted ~10 000 words). Better topic
  hints but maintenance burden.
- **heroku-style** (adjective-noun-N, e.g. "calm-river-42"). Good
  topic hint when the words are seeded by content; pure entropy when
  not. Long.
- **military concatenated syllables** (milconcat). Variable-length
  syllable concatenation derived from slug words. Topic-bearing but
  variable-length defeats column alignment in indexes; collision
  resolution adds further drift from the "fixed N" promise.
- **specops-style** (HAVE BLUE, COBRA BALL). Two-word allcaps. Too
  long; reserved-list overhead.

Steve's reactions narrowed the field:

- "i like the hint about topic that milconcat provides but i don't
  like the arbitrarily long names." (variable-length penalty)
- "i like the conciseness and pronounceability of proquint but don't
  like that it provides no clue about topic." (proquint-1 alone is
  topic-opaque)

Bot showed a side-by-side of `{5-character milconcat}-{proquint-1}`
and `{6-character milconcat}-{proquint-1}` to test whether a
topic-bearing prefix on a fixed-length proquint suffix would carry.
The milconcat prefix kept colliding on similar slugs, and "fixed N"
became fictional after collision resolution.

### Step 3 -- FAA detour

Steve: "what does FAA use for IFR checkpoints? how do they
generate those?"

Bot researched the FAA 5-letter pronounceable name (5LNC) system
and the ICARD registry (FAA Aeronautical Information Services in
Silver Spring). Key findings:

- 5LNCs are random pronounceable strings; mnemonic
  topic-bearing codes were explicitly REJECTED by the FAA because
  charts get cluttered when controllers and pilots try to read
  meaning into them.
- A central registry resolves collisions before publication.
- The collision-resolution machinery, not the entropy source, is
  what makes the system work.

This collapsed the design space. Topic-bearing prefixes were
re-rejected (FAA had already gone through the same gameout 50 years
ago). The design question became "what is wire-lab's central
registry equivalent?"

### Step 4 -- the registry is the filename glob

Steve: "what if we just use proquint-1 but always check against
existing to ensure no collision?"

Bot raised the central-authority concern: TE-dajot says no central
authority for a century. Steve's clarification:

> "wire-lab is design notes about a system that has those
> constraints. wire-lab itself is a shared-fork repo with a
> mob-consensus workflow."

In other words, the 100-year-no-authority invariant binds the
protocol that wire-lab documents, not wire-lab's own corpus
management. Wire-lab can use any local convention for its own
files, including a per-fork registry. And the per-fork registry can
be implicit -- the filename glob across `docs/thought-experiments/`
and `protocols/*/TODO/` IS the registry.

Steve also closed off the re-derivability question:

> "we don't need to be able to re-derive a handle -- the goal is to
> simply avoid collisions."

This makes the entropy source irrelevant. Time-ns, commit-hash,
SHA-256-of-anything; whatever the seed, the collision check at mint
time is what guarantees uniqueness. The bot tabled
proquint1(time_ns), proquint1(commit_hash, time_ns), and
proquint1(time_ns) with a 70-entry corpus and zero retries.

### Step 5 -- locking the filename grammar

Bot's mistake earlier in the gameout: assuming the proquint would
be added to the filename ALONGSIDE the timestamp, leaving the
timestamp as a redundant secondary key. Steve corrected:

> "i think you're confused -- i mean to replace the TE and TODO
> integers and the timestamps in their filenames with proquint1."

Lock: the filename is `(TE|TODO)-<proquint>-<slug>.md`. The
proquint replaces both the integer alias and the timestamp prefix.
The timestamp survives as a column in the README index (extracted
from `git log` of the introducing commit) and as a "Prior aliases"
note in each renamed file's body, but it leaves the filename.

### Step 6 -- tooling order

Steve: "aren't you going to write a minting program first? (i'd
prefer go.)"

Bot built `tools/mint-handle/` in Go (proquint encoder, corpus
scanner, retry loop, CLI flags, 12 unit tests). The corpus scan
walks `scanLiteralDirs` plus `scanGlobs` (currently
`protocols/*/TODO`) and treats any file matching the new
`(TE|TODO)-<proquint>-<slug>.md` regex as owning its handle.

Steve confirmed the bulk-rename invocation pattern:

> "calling 66 times is fine"

so no `-batch` flag was needed. The migration tool was implemented
as a separate program (`tools/migrate-handles/`) that calls the
encoder directly with an in-memory claimed set so the 70-file batch
cannot self-collide.

### Step 7 -- prior-alias recording

Bot's question (Q-23.9): how to record the prior integer alias for
historical traceability after the rename?

Steve: "each file should include 'was TE-NN' in its frontmatter or
text. maybe do both?"

Bot's follow-up (Q-23.10): TE files have no YAML frontmatter; they
have a prose `## TE ID` section. Where exactly should the
prior-alias text go?

Steve: "New '## Prior aliases' section."

Lock: dual record. (a) `docs/thought-experiments/README.md` and
`protocols/wire-lab.d/TODO/TODO.md` gain a `prior_alias` column
(integer plus timestamp). (b) Each renamed file body gains a new
`## Prior aliases` section immediately after `## TE ID`, listing the
integer (TE-NN) and the timestamp (TE-YYYYMMDD-HHMMSS) it carried
before this commit.

## Locked decisions (summary)

1. **Handle scheme:** proquint-1 (5 chars, Wilkerson 2009 alphabet
   `bdfghjklmnprstvz` consonants and `aiou` vowels, CVCVC shape,
   65 536 distinct values).

2. **Filename grammar:** `(TE|TODO)-<proquint>-<slug>.md`. The
   proquint replaces both the integer alias and the timestamp
   prefix. The slug remains in lower-kebab and may include digits.

3. **Mint algorithm:** mint-time-allocated random label.
   Entropy = `time.Now().UnixNano()`, folded through SHA-256, first
   2 bytes encoded as a proquint-1. Retry on collision with a fresh
   nanosecond reading.

4. **Registry:** the union of filenames in `docs/thought-experiments/`
   and `protocols/*/TODO/` IS the registry. No HANDLES.md, no
   frontmatter `handle:` field, no separate authority. The filename
   is canonical.

5. **Mint-date provenance:** retained as a column in the README
   index (extracted from `git log` of the introducing commit, not
   from the filename, since the filename no longer encodes time).

6. **Cross-fork behavior:** forks diverge naturally. Two forks may
   independently mint `vapoj` for different files after the fork
   point; this is a feature, not a bug, consistent with PromiseGrid's
   no-central-authority stance applied recursively to its own design
   corpus. Reconciliation at merge time picks one occupant of the
   handle and re-mints the other (mechanically straightforward
   because the handle is just a filename component).

7. **Saturation threshold:** if the corpus exceeds 32 000 handles
   (50 percent of pq1 space) we revisit and likely upgrade new
   mints to proquint-2 (11 chars, 4 294 967 296 distinct values).
   Not a concern for any plausible wire-lab evolution.

8. **Prior-alias recording:** dual record. README index column plus
   per-file `## Prior aliases` section. Pre-migration session-logs
   in `wire-lab-logs` repo keep verbatim integer and timestamp
   references; the dual record makes them resolvable indefinitely.

## Alternatives considered and rejected

- **Alt-A: drop integers, keep timestamps.** Rejected because
  timestamps are not citable in prose, and the timestamp granularity
  (HHMMSS) does not actually guarantee uniqueness across fast-
  drafting batches (the `2026-05-07-002306` cluster in the existing
  corpus contains 9 TODOs that all share the same timestamp).
- **Alt-B: keep integers, drop timestamps.** Rejected because it
  re-introduces the cross-fork collision problem TE-titur was patching.
- **Alt-C: keep both.** The status quo. Rejected because the dual
  namespace is confusing and the integer-alias allocator has no
  fork-safe path.
- **Alt-D: heroku-style adjective-noun-N.** Rejected on length and
  topic-opacity grounds (the topic hint requires content-aware
  seeding the bot does not have).
- **Alt-E: bip39 single word.** Rejected on length variance and
  collision rate at single-word granularity.
- **Alt-F.1: milconcat (military concatenated syllables).**
  Rejected because variable length defeats column alignment and
  collision resolution drifts the "fixed N" promise.
- **Alt-F.2: milconcat prefix + proquint-1 suffix.** Rejected
  because the milconcat prefix still collides on similar slugs and
  carries the same length-drift problem.
- **Alt-F.3: specops-style allcaps two-word.** Rejected on length
  and reserved-list maintenance grounds.
- **Alt-G: derived handle (e.g. proquint-1(SHA-256(content))).**
  Tabled. The bot raised it as a "self-certifying handle" option
  (the handle re-derives from the file content, no registry
  needed). Steve closed it off: "we don't need to be able to
  re-derive a handle -- the goal is to simply avoid collisions."
  Re-derivability would also collide with the editing-policy TEs
  (TE-dabol, TE-vudaf) which permit in-place edits, breaking the handle's
  stability under any content-bound derivation.
- **Alt-H: FAA-style 5LNC with central registry.** Rejected because
  wire-lab has no central authority equivalent. The registry-as-
  filename-glob trick (locked above) achieves the same collision
  guarantee using the corpus itself as the registry.

## Migration scope (for the record)

Bounded. This twig executes the migration in four chunks on the
same branch:

- **Chunk A (committed):** rename 70 files from
  `(TE|TODO)-<timestamp>-<slug>.md` to
  `(TE|TODO)-<proquint>-<slug>.md` and inject a `## Prior aliases`
  section in each renamed file body.
- **Chunk B (this commit):** TE-mumuv itself.
- **Chunk C (pending):** corpus-wide citation sweep replacing
  integer and timestamp references with proquint handles. Touches
  ~200 individual references across TE bodies, TODOs, AGENTS files,
  the harness-spec, DR/DI files, and the README.
- **Chunk D (pending):** TE-titur (titur) gets a Cat-3 Refinement
  recording the proquint adoption; the README index gains a
  `prior_alias` column and uses the proquint as the primary key;
  the master TODO cross-list is rebuilt the same way; the twig
  merges to ppx/main.

Pre-migration audit reports under `protocols/wire-lab.d/TODO/`
(`pre18-audit-report-20260505.md`, `pre149-audit-report-20260505.md`,
`dropped-thread-disposition-20260506.md`) are kept verbatim as
historical record. Their references to integer and timestamp
aliases stay; future readers cross-reference via the README
`prior_alias` column or the per-file `## Prior aliases` section.

## Cross-references

- **Predecessor:** TE-titur (titur), "TE-nibar numbering collision and
  harness-spec path" (2026-04-30), introduced the dual-identifier
  scheme that this TE retires. A Cat-3 Refinement on TE-titur (added
  in chunk D) records the proquint adoption.
- **Foundation:** TE-dajot (jamod), "100-year goal as a load-bearing
  design constraint", cited above for the no-central-authority
  invariant that wire-lab's documentary corpus is design notes
  ABOUT (not bound by).
- **Sibling:** TE-vipir (gitap), "Protocols as simulated repos and the
  L4-binding layer", uses the same collision-checked allocation
  pattern at the protocol level. The handle scheme in this TE is
  the corpus-management instance of that pattern.
- **Editing-policy interaction:** TE-dabol (dabol) and TE-vudaf (vudaf)
  permit in-place edits of TE bodies. Handles must therefore be
  stable under content edits, which rules out content-derived
  handles (Alt-G above).
- **Vocabulary:** AGENTS.md and AGENTS-ppx.md gain a "handle"
  vocabulary line in chunk D, replacing any "integer alias"
  guidance.

## Decision status

Decided. Implementation in flight on twig
`ppx/te-20260507-025627-naming-reconciliation`. Chunks A and B
committed; chunks C and D pending.
