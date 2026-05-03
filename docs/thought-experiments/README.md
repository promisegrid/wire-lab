# Thought Experiments

Each thought experiment is a falsifiable mental run of a Wire Lab design choice. Each lives in its own file, addressable by content hash (its pCID).

## Filename convention

```
TE-YYYYMMDD-HHMMSS-some-phrase.md
```

The timestamp is the moment the TE was first drafted (or, for TEs that pre-date this convention, the revision of `protocols/wire-lab.d/specs/harness-spec-draft.md` in which they first appeared). The slug is a kebab-case rendering of the TE's title. Files are not renamed when the experiment is later refined — the timestamp pins origin, not last-edited.

## Index

| # | First drafted | Title |
|---|---|---|
| TE-1 | 2026-04-27 18:00 | [Promise-stack ordering](TE-20260427-180000-promise-stack-ordering.md) |
| TE-2 | 2026-04-27 18:01 | [Trust-ledger merge after partition](TE-20260427-180100-trust-ledger-merge-after-partition.md) |
| TE-3 | 2026-04-27 18:02 | [Currency exchange-rate equilibration](TE-20260427-180200-currency-exchange-rate-equilibration.md) |
| TE-4 | 2026-04-27 18:03 | [Sybil under double auction](TE-20260427-180300-sybil-under-double-auction.md) |
| TE-5 | 2026-04-27 18:04 | [Kernel-as-handler vs. classical kernel](TE-20260427-180400-kernel-as-handler-vs-classical-kernel.md) |
| TE-6 | 2026-04-27 18:05 | [Capability-token revocation propagation](TE-20260427-180500-capability-token-revocation-propagation.md) |
| TE-7 | 2026-04-27 18:06 | [Human-novice onboarding under K4](TE-20260427-180600-human-novice-onboarding-under-k4.md) |
| TE-8 | 2026-04-27 18:07 | [Generational handoff](TE-20260427-180700-generational-handoff.md) |
| TE-9 | 2026-04-27 18:08 | [Two communities, two pCIDs, same intent](TE-20260427-180800-two-communities-two-pcids-same-intent.md) |
| TE-10 | 2026-04-27 18:09 | [Slow-mover survival](TE-20260427-180900-slow-mover-survival.md) |
| TE-11 | 2026-04-27 18:10 | [Ostrom's principles audit](TE-20260427-181000-ostroms-principles-audit.md) |
| TE-12 | 2026-04-27 18:11 | [Promise-stack as zero-knowledge envelope](TE-20260427-181100-promise-stack-as-zero-knowledge-envelope.md) |
| TE-13 | 2026-04-27 18:12 | [Time-traveling break-witness](TE-20260427-181200-time-traveling-break-witness.md) |
| TE-14 | 2026-04-28 08:00 | [A harness-spec change walks through the unified flow](TE-20260428-080000-harness-spec-change-walks-through-unified-flow.md) |
| TE-15 | 2026-04-28 09:45 | [Should this design become `promisegrid/promisegrid/README.md`?](TE-20260428-094500-should-this-design-become-promisegrid-readme.md) |
| TE-16 | 2026-04-29 03:32 | [Durable review feedback as contest artifact](TE-20260429-033208-review-feedback-as-contest-artifact.md) |
| TE-17 | 2026-04-29 16:22 | [Review reply as promise](TE-20260429-162212-review-reply-as-promise.md) |
| TE-18 | 2026-04-29 16:51 | [Bot identity and branch prefix](TE-20260429-165101-bot-identity-and-branch-prefix.md) |
| TE-19 | 2026-04-29 16:51 | [Branch-protection posture for `main`](TE-20260429-165102-branch-protection-posture.md) |
| TE-20 | 2026-04-29 16:51 | [Bot review style](TE-20260429-165103-bot-review-style.md) |
| TE-21 | 2026-04-29 17:35 | [Spec doc as promise](TE-20260429-173520-spec-doc-as-promise.md) |
| TE-22 | 2026-04-29 17:55 | [Spec-doc store layout and pCID machinery](TE-20260429-175530-spec-doc-store-and-pcid-machinery.md) |
| TE-23 | 2026-04-30 06:43 | [Congruence/convergence duality and pCID framing](TE-20260430-064307-congruence-convergence-duality-and-pcid-framing.md) |
| TE-24 | 2026-04-30 20:41 | [Group-transport envelope: `grid <pcid>` carrier, canonical bytes, and explicit promise body](TE-20260430-204108-group-transport-envelope.md) (rewritten in place per TODO 013; original drafting used "channel" vocabulary) |
| TE-25 | 2026-04-30 21:34 | [TE-21 numbering collision and harness-spec path](TE-20260430-213447-te-numbering-collision-and-harness-spec-path.md) |
| TE-26 | 2026-04-30 21:56 | [Transport-protocol types, pCID-keyed transport paths, and DAG message graphs](TE-20260430-215624-channel-transport-types-and-threaded-replies.md) (filename retains "channel" from original drafting; rewritten in place per TE-27) |
| TE-27 | 2026-05-01 02:19 | [Transports rename and axes of transport-protocol differentiation](TE-20260501-021921-transports-rename-and-axes-of-differentiation.md) |
| TE-28 | 2026-05-01 20:27 | [The 100-year goal as a load-bearing design constraint](TE-20260501-202713-100-year-goal-as-design-constraint.md) |
| TE-29 | 2026-05-01 21:50 | [Protocols as simulated repos, and the L4-binding layer](TE-20260501-215027-protocols-as-simulated-repos-and-binding-layer.md) |
| TE-30 | 2026-05-02 00:25 | [TODO numbering and per-protocol TODO shape](TE-20260502-002548-todo-numbering-and-per-protocol-shape.md) |
| TE-31 | 2026-05-02 00:49 | [Spec-doc as upstream, simrepo as implementation: inverting the conformance reference](TE-20260502-004924-spec-doc-inversion-and-conformance-changelog.md) |
| TE-32 | 2026-05-02 01:45 | [Spec-side vs implementation-side split, and the `implementations/` top-level](TE-20260502-014525-spec-vs-implementation-split.md) |
| TE-33 | 2026-05-02 02:04 | [Spec-doc Informative References to its workshop, RFC-shaped](TE-20260502-020439-spec-doc-informative-references.md) |
| TE-34 | 2026-05-02 21:28 | [TE editing policy and the TE corpus as one document with facets](TE-20260502-212810-te-editing-policy-and-holistic-corpus.md) |
| TE-35 | 2026-05-02 23:26 | [Tabletop simulation of the TE editing policy](TE-20260502-232651-editing-policy-tabletop.md) |
| TE-36 | 2026-05-03 02:24 | [Apparatus vs. specimen — carving the harness-spec apart from the wire/envelope/ledger hypotheses it studies](TE-20260503-022446-apparatus-vs-specimen-carve-out.md) |

The TE numbers (TE-1, TE-2, …) are stable identifiers used in the harness-spec; the timestamp slug is what makes the file content-addressable and chronologically sortable on disk.

## Editing policy

TE filenames are immutable: the timestamp slug is the content-address anchor that pins the integer alias (TE-1, TE-2, ...), locked in [TE-25](TE-20260430-213447-te-numbering-collision-and-harness-spec-path.md).

TE contents are edited under a categorized policy locked in [TE-34](TE-20260502-212810-te-editing-policy-and-holistic-corpus.md) and refined by [TE-35](TE-20260502-232651-editing-policy-tabletop.md). The locked DIs are `DI-020-20260502-213103` (categorized regimes), `DI-020-20260502-213104` (uniform applicability across all TE corpora), and `DI-020-20260502-213105` (holistic reading by default; single-TE reading only for obviously mechanical questions). The Cat-1 clause of `DI-020-20260502-213103` was superseded on 2026-05-02 by `DI-020-20260502-232651` (Cat-1a / Cat-1b split). Four Cat-3 navigational refinements appear in TE-34's `## Refinements` section. The canonical statement of the policy lives in `AGENTS.md` under "TE Editing Policy (Required)"; the seven categories in summary:

- **Cat-1a (current-pointer paths).** Mechanical sweep in place; no top-of-file note.
- **Cat-1b (historical-quotation paths).** Left untouched. Path references inside markdown blockquotes, attributed to another TE ("TE-N states ..."), in past tense, inside `## Refinements` sections, supersedence notes, or `Decision status` lines are Cat-1b. When in doubt, treat as Cat-1b.
- **Cat-2 (vocabulary updates).** Edit in place, with a top-of-file note pointing at the driving TE or TODO. The note must enumerate by ID every DI that lives in the affected TE, paired with an explicit promise that the rewrite preserves each DI's meaning. Mandatory pre-step: grep the corpus for the old term inside quotation contexts and classify each match Cat-2 (sweep) or Cat-2-historical (leave) before sweeping.
- **Cat-3 (navigational forward pointers).** Append a dated entry to the TE's `## Refinements` section (created if absent, placed after `## Decision status`). The TE body above is unchanged. No DI is filed.
- **Cat-4 (resolved-implication forward pointers).** Same shape as Cat-3, used when an Implications-and-future-work item resolves (a TODO filed; a DR opened; a downstream TE landed).
- **Cat-5 / Cat-6 / Cat-7 (substantive supersedence).** Not edits. Write a new TE that supersedes the old one and a new DI that supersedes the old DI. Update the older TE's `## Decision status` to `superseded by TE-<id>` and its top-of-file `## Status` field to `superseded by TE-<id> / DI-<id>`; otherwise leave the body untouched.

Every TE carries a top-of-file `## Status` field placed immediately after the TE ID line. Canonical values: `needs DF`, `decided`, `decided, refined`, `superseded by TE-<id> / DI-<id>`, `withdrawn`. Legacy values preserved during the 2026-05-02 retrofit: `stub`, `open`, `recommended for immediate adoption`, `locked for the <protocol>`. New TEs prefer canonical values.

The corpus is read holistically by default: the TE corpus is one document with many facets, not a collection of independent essays. When any TE is in scope, the first move is to scan TE titles plus `## Status` fields and `Decision under test` sections across the corpus to find facets that share assumptions, vocabulary, or decisions. Single-TE reading is reserved for obviously mechanical questions (a single typo; a path that has demonstrably moved; a `## Status` field retrofit) and only after the holistic read has confirmed the question is mechanical.

Applicability is uniform across every TE corpus in this repository, whether the TE lives at the harness level (this directory) or inside a per-protocol directory (`protocols/<slug>.d/`). Per-protocol corpora may add stricter rules but may not relax these rules.

## Adding a new TE

1. Decide the title.
2. Pick a UTC timestamp — typically `date -u +%Y%m%d-%H%M%S`.
3. Render the title to kebab-case for the slug.
4. Create `TE-YYYYMMDD-HHMMSS-slug.md` in this directory. Include a top-of-file `## Status` field placed immediately after the TE ID line, with the appropriate initial value (`needs DF` for a TE in DF state; `decided` for a TE that locks DIs in the same commit). Use canonical values; reserve legacy values for the retrofit corpus.
5. Add a one-line summary to `../../protocols/wire-lab.d/specs/harness-spec-draft.md` §8 with a link.
6. Add the row to this index.
7. Open a PR.
