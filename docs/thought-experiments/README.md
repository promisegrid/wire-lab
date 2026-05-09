# Thought Experiments

Each thought experiment is a falsifiable mental run of a Wire Lab design choice. Each lives in its own file.

## Filename convention

```
TE-<proquint>-<slug>.md
```

The proquint handle (5 lowercase characters from the alphabet `bdfghjklmnprstvz` x `aiou` x `bdfghjklmnprstvz` x `aiou` x `bdfghjklmnprstvz`, per Wilkerson 2009) is the stable identifier of a TE. Handles are minted by `tools/mint-handle` from `time_ns -> SHA-256 -> first 2 bytes -> proquint-1`, with collision-retry against the directory glob. The handle is short, pronounceable, fork-stable, and assigned at file-creation time. The slug is a kebab-case rendering of the TE's title; it is informational and may be edited.

The proquint handle replaces both the integer alias (TE-1, TE-2, ...) and the timestamp slug (TE-YYYYMMDD-HHMMSS) used before the TE-39 (TE-mumuv) migration. The drafting timestamp is preserved in git history (the file's first commit) and surfaces in the index below as the **Mint date** column. Each migrated file carries a `## Prior aliases` section recording the integer + timestamp aliases it had before the migration.

## Index

| Handle | Mint date | Title | Prior alias |
|---|---|---|---|
| [TE-famar](TE-famar-promise-stack-ordering.md) | 2026-04-28 | Promise-stack ordering | `TE-1` / `TE-20260427-180000` |
| [TE-jovoj](TE-jovoj-trust-ledger-merge-after-partition.md) | 2026-04-28 | Trust-ledger merge after partition | `TE-2` / `TE-20260427-180100` |
| [TE-pubum](TE-pubum-currency-exchange-rate-equilibration.md) | 2026-04-28 | Currency exchange-rate equilibration | `TE-3` / `TE-20260427-180200` |
| [TE-himug](TE-himug-sybil-under-double-auction.md) | 2026-04-28 | Sybil under double auction | `TE-4` / `TE-20260427-180300` |
| [TE-jikaf](TE-jikaf-kernel-as-handler-vs-classical-kernel.md) | 2026-04-28 | Kernel-as-handler vs. classical kernel | `TE-5` / `TE-20260427-180400` |
| [TE-fijub](TE-fijub-capability-token-revocation-propagation.md) | 2026-04-28 | Capability-token revocation propagation | `TE-6` / `TE-20260427-180500` |
| [TE-kuhog](TE-kuhog-human-novice-onboarding-under-k4.md) | 2026-04-28 | Human-novice onboarding under K4 | `TE-7` / `TE-20260427-180600` |
| [TE-sigan](TE-sigan-generational-handoff.md) | 2026-04-28 | Generational handoff | `TE-8` / `TE-20260427-180700` |
| [TE-morid](TE-morid-two-communities-two-pcids-same-intent.md) | 2026-04-28 | Two communities, two pCIDs, same intent | `TE-9` / `TE-20260427-180800` |
| [TE-mokut](TE-mokut-slow-mover-survival.md) | 2026-04-28 | Slow-mover survival | `TE-10` / `TE-20260427-180900` |
| [TE-rotim](TE-rotim-ostroms-principles-audit.md) | 2026-04-28 | Ostrom's principles audit | `TE-11` / `TE-20260427-181000` |
| [TE-muvuv](TE-muvuv-promise-stack-as-zero-knowledge-envelope.md) | 2026-04-28 | Promise-stack as zero-knowledge envelope | `TE-12` / `TE-20260427-181100` |
| [TE-robub](TE-robub-time-traveling-break-witness.md) | 2026-04-28 | Time-traveling break-witness | `TE-13` / `TE-20260427-181200` |
| [TE-botom](TE-botom-harness-spec-change-walks-through-unified-flow.md) | 2026-04-28 | A harness-spec change walks through the unified flow | `TE-14` / `TE-20260428-080000` |
| [TE-dodaf](TE-dodaf-should-this-design-become-promisegrid-readme.md) | 2026-04-28 | Should this design become `promisegrid/promisegrid/README.md`? | `TE-15` / `TE-20260428-094500` |
| [TE-lodar](TE-lodar-review-feedback-as-contest-artifact.md) | 2026-04-28 | Thought Experiment - Durable review feedback as contest artifact | `TE-16` / `TE-20260429-033208` |
| [TE-mirah](TE-mirah-review-reply-as-promise.md) | 2026-04-29 | Thought Experiment - Review reply as promise | `TE-17` / `TE-20260429-162212` |
| [TE-lusut](TE-lusut-bot-identity-and-branch-prefix.md) | 2026-04-29 | Bot identity and branch prefix | `TE-18` / `TE-20260429-165101` |
| [TE-mipat](TE-mipat-branch-protection-posture.md) | 2026-04-29 | Branch-protection posture for `main` | `TE-19` / `TE-20260429-165102` |
| [TE-tibas](TE-tibas-bot-review-style.md) | 2026-04-29 | Bot review style | `TE-20` / `TE-20260429-165103` |
| [TE-nibar](TE-nibar-spec-doc-as-promise.md) | 2026-04-29 | Spec doc as promise | `TE-21` / `TE-20260429-173520` |
| [TE-rujak](TE-rujak-spec-doc-store-and-pcid-machinery.md) | 2026-04-29 | Spec-doc store layout and pCID machinery | `TE-22` / `TE-20260429-175530` |
| [TE-lozip](TE-lozip-congruence-convergence-duality-and-pcid-framing.md) | 2026-04-30 | Congruence/convergence duality and pCID framing | `TE-23` / `TE-20260430-064307` |
| [TE-hogus](TE-hogus-group-transport-envelope.md) | 2026-05-01 | Group-transport envelope: `grid <pcid>` carrier, canonical bytes, and explicit promise body | `TE-24` / `TE-20260430-204108` |
| [TE-titur](TE-titur-te-numbering-collision-and-harness-spec-path.md) | 2026-04-30 | Reconciling the TE-nibar numbering collision and the `harness-spec.md` path on the channels branch | `TE-25` / `TE-20260430-213447` |
| [TE-zalut](TE-zalut-channel-transport-types-and-threaded-replies.md) | 2026-05-01 | Transport-protocol types, pCID-keyed transport paths, and DAG message graphs | `TE-26` / `TE-20260430-215624` |
| [TE-junil](TE-junil-transports-rename-and-axes-of-differentiation.md) | 2026-05-01 | `channels/` → `transports/` rename and axes of transport-protocol differentiation | `TE-27` / `TE-20260501-021921` |
| [TE-dajot](TE-dajot-100-year-goal-as-design-constraint.md) | 2026-05-01 | The 100-year goal as a load-bearing design constraint | `TE-28` / `TE-20260501-202713` |
| [TE-vipir](TE-vipir-protocols-as-simulated-repos-and-binding-layer.md) | 2026-05-01 | Protocols as simulated repos, and the L4-binding layer | `TE-29` / `TE-20260501-215027` |
| [TE-magup](TE-magup-todo-numbering-and-per-protocol-shape.md) | 2026-05-02 | TODO numbering and per-protocol TODO shape | `TE-30` / `TE-20260502-002548` |
| [TE-zukug](TE-zukug-spec-doc-inversion-and-conformance-changelog.md) | 2026-05-02 | Spec-doc as upstream, simrepo as implementation — inverting the conformance reference | `TE-31` / `TE-20260502-004924` |
| [TE-liviv](TE-liviv-spec-vs-implementation-split.md) | 2026-05-02 | Spec-side vs implementation-side split, and the `implementations/` top-level | `TE-32` / `TE-20260502-014525` |
| [TE-potar](TE-potar-spec-doc-informative-references.md) | 2026-05-02 | Spec-doc Informative References to its workshop, RFC-shaped | `TE-33` / `TE-20260502-020439` |
| [TE-dabol](TE-dabol-te-editing-policy-and-holistic-corpus.md) | 2026-05-02 | TE editing policy and the TE corpus as one document with facets | `TE-34` / `TE-20260502-212810` |
| [TE-vudaf](TE-vudaf-editing-policy-tabletop.md) | 2026-05-02 | Tabletop simulation of the TE editing policy | `TE-35` / `TE-20260502-232651` |
| [TE-havib](TE-havib-apparatus-vs-specimen-carve-out.md) | 2026-05-06 | Apparatus vs. specimen — carving the harness-spec apart from the wire/envelope/ledger hypotheses it studies | `TE-36` / `TE-20260503-022446` |
| [TE-numan](TE-numan-transport-protocol-migration-semantics.md) | 2026-05-06 | Transport-protocol migration invariants | `TE-37` / `TE-20260506-041241` |
| [TE-sihih](TE-sihih-substrate-agnostic-layered-model.md) | 2026-05-07 | Substrate-agnostic layered model (L5/L6/L7) and L6 CAS subtree | `TE-38` / `TE-20260506-184800` |
| [TE-mumuv](TE-mumuv-naming-reconciliation.md) | 2026-05-07 | Naming reconciliation -- proquint handles for TEs and TODOs | `TE-39` (newly minted in TE-39 twig) |
| [TE-nijab](TE-nijab-transport-layering-and-freeze-boundaries.md) | 2026-05-08 | Transport layering and freeze boundaries |  |
| [TE-david](TE-david-promisegrid-dev-guide-resources.md) | 2026-05-08 | PromiseGrid dev-guide resources from wire-lab evidence |  |

The proquint handle is **both** the stable identifier and the display nickname. It is collision-free at mint time, fork-stable across branches (each fork mints its own handles; collisions at merge time are handled by re-minting), and short enough to use directly in prose ("per TE-titur S5"). DF / DI / DR descendant numbering still uses the handle root: DF-titur.1, DI-titur-..., DR-009 (DR has its own numbering scheme). Backward citations to integer aliases (e.g., "per TE-25 S5") remain valid; readers may consult the cited file's `## Prior aliases` section or the `Prior alias` column above to recover the integer.

Forward-pointers to a not-yet-drafted TE MUST use a thread-id (T-...) recorded in `OPEN-THREADS.md`. Naming an unminted *future* proquint is impossible (proquints are minted, not predicted), and naming a future *integer* is forbidden (the construction that produced the DT3 drift, locked closed by the 2026-05-07 Cat-3 Refinement on TE-titur).

## Editing policy

TE filenames are mostly immutable: once the proquint handle is minted, the file keeps that handle through any future title or content edits. The exception is the TE-39 corpus migration itself, which renamed 70 files in a single commit (`85766f0`); each migrated file's `## Prior aliases` section preserves the audit trail.

TE contents are edited under a categorized policy locked in [TE-dabol](TE-dabol-te-editing-policy-and-holistic-corpus.md) and refined by [TE-vudaf](TE-vudaf-editing-policy-tabletop.md). The locked DIs are `DI-020-20260502-213103` (categorized regimes), `DI-020-20260502-213104` (uniform applicability across all TE corpora), and `DI-020-20260502-213105` (holistic reading by default; single-TE reading only for obviously mechanical questions). The Cat-1 clause of `DI-020-20260502-213103` was superseded on 2026-05-02 by `DI-020-20260502-232651` (Cat-1a / Cat-1b split). Several Cat-3 navigational refinements appear in TE-dabol's `## Refinements` section. The canonical statement of the policy lives in `AGENTS.md` under "TE Editing Policy (Required)"; the seven categories in summary:

- **Cat-1a (current-pointer paths).** Mechanical sweep in place; no top-of-file note.
- **Cat-1b (historical-quotation paths).** Left untouched. Path references inside markdown blockquotes, attributed to another TE ("TE-N states ..."), in past tense, inside `## Refinements` sections, supersedence notes, or `Decision status` lines are Cat-1b. When in doubt, treat as Cat-1b.
- **Cat-2 (vocabulary updates).** Edit in place, with a top-of-file note pointing at the driving TE or TODO. The note must enumerate by ID every DI that lives in the affected TE, paired with an explicit promise that the rewrite preserves each DI's meaning. Mandatory pre-step: grep the corpus for the old term inside quotation contexts and classify each match Cat-2 (sweep) or Cat-2-historical (leave) before sweeping.
- **Cat-3 (navigational forward pointers).** Append a dated entry to the TE's `## Refinements` section (created if absent, placed after `## Decision status`). The TE body above is unchanged. No DI is filed.
- **Cat-4 (resolved-implication forward pointers).** Same shape as Cat-3, used when an Implications-and-future-work item resolves (a TODO filed; a DR opened; a downstream TE landed).
- **Cat-5 / Cat-6 / Cat-7 (substantive supersedence).** Not edits. Write a new TE that supersedes the old one and a new DI that supersedes the old DI. Update the older TE's `## Decision status` to `superseded by TE-<handle>` and its top-of-file `## Status` field to `superseded by TE-<handle> / DI-<id>`; otherwise leave the body untouched.

Every TE carries a top-of-file `## Status` field placed immediately after the TE ID line. Canonical values: `needs DF`, `decided`, `decided, refined`, `superseded by TE-<handle> / DI-<id>`, `withdrawn`. Legacy values preserved during the 2026-05-02 retrofit: `stub`, `open`, `recommended for immediate adoption`, `locked for the <protocol>`. New TEs prefer canonical values.

The corpus is read holistically by default: the TE corpus is one document with many facets, not a collection of independent essays. When any TE is in scope, the first move is to scan TE titles plus `## Status` fields and `Decision under test` sections across the corpus to find facets that share assumptions, vocabulary, or decisions. Single-TE reading is reserved for obviously mechanical questions (a single typo; a path that has demonstrably moved; a `## Status` field retrofit) and only after the holistic read has confirmed the question is mechanical.

Applicability is uniform across every TE corpus in this repository, whether the TE lives at the harness level (this directory) or inside a per-protocol directory (`protocols/<slug>.d/`). Per-protocol corpora may add stricter rules but may not relax these rules.

## Adding a new TE

1. Decide the title.
2. Render the title to kebab-case for the slug.
3. Mint a proquint handle: `cd tools/mint-handle && go run . -w 1`. The tool scans `docs/thought-experiments/` and `protocols/*/TODO/` for existing handles and retries until it finds a collision-free value.
4. Create `TE-<handle>-<slug>.md` in this directory. Include a top-of-file `## Status` field placed immediately after the TE ID line, with the appropriate initial value (`needs DF` for a TE in DF state; `decided` for a TE that locks DIs in the same commit). Use canonical values; reserve legacy values for the retrofit corpus.
5. While drafting, write any forward-pointers (to TEs that do not yet exist) using a thread-id from `OPEN-THREADS.md`. Naming an unminted future proquint is impossible (proquints are minted, not predicted); naming a future integer alias is forbidden (the DT3 drift class, locked closed by the 2026-05-07 Cat-3 Refinement on TE-titur).
6. Add a one-line summary to `../../protocols/wire-lab.d/specs/harness-spec-draft.md` Section 8 with a link.
7. Add the row to this index. The mint date is the date the file is committed (column auto-fills from `git log --format=%ad --date=short -- <path> | tail -1`); leave the `Prior alias` column blank for new TEs that never carried an integer alias.
8. Open a PR.

## Migration history

- **2026-05-07.** TE-mumuv (TE-39, formerly TE-39) locks proquint handles. Migration tool `tools/migrate-handles/` renamed 38 TEs and 32 TODOs in a single commit (`85766f0`); each renamed file gained a `## Prior aliases` section. Citation sweep `tools/sweep-citations/` rewrote 1,451 references across 94 files. The integer alias (TE-N) and the timestamp alias (TE-YYYYMMDD-HHMMSS) survive only as `Prior alias` entries in this index and in each file's `## Prior aliases` section.
- **2026-05-05.** TE-titur Cat-3 Refinement (chained on the same TE) split TE identity into a stable identifier (timestamp+slug filename) and a display nickname (integer TE-N). Superseded on 2026-05-07 by the proquint adoption above; the 2026-05-05 forward-pointer rule is reaffirmed and tightened in the 2026-05-07 Refinement.
- **2026-04-30.** TE-titur (formerly TE-25) locked the integer-alias display-nickname convention.
