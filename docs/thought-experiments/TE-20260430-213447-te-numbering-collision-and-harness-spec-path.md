# TE-25: Reconciling the TE-21 numbering collision and the `harness-spec.md` path on the channels branch

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-20260430-213447

## Status

decided

## Decision under test

When `origin/stevegt/channels-grid-pcid` is reconciled with `ppx/main`, two surface conflicts arise that are not really about channels at all but about repo conventions. They are joined here because the same merge will resolve both, and they should be settled together so we have a single record:

1. The channels branch labels its TE as **TE-21**, but `ppx/main` already has a different TE-21 (`docs/thought-experiments/TE-20260429-173520-spec-doc-as-promise.md`). One of them has to be renumbered.
2. The channels branch edits a top-level `harness-spec.md`, but on `ppx/main` the spec was renamed to `specs/harness-spec-draft.md` during the genesis-freeze work in TODO 011. The patch cannot apply cleanly without translation.

In addition, the channels branch adds an "Open Question #19" at the bottom of the harness-spec, and that bullet must end up at the correct slot in the renamed file.

## Assumptions

- The TE numbering convention is recorded in `docs/thought-experiments/README.md`: "The TE numbers (TE-1, TE-2, …) are stable identifiers used in the harness-spec; the timestamp slug is what makes the file content-addressable and chronologically sortable on disk."
- The "first drafted" timestamp is the integer-anchoring fact. Two TEs first drafted at different moments cannot share an integer.
- Once a TE integer has appeared in `ppx/main` and been merged through, it is sticky for that TE: renumbering a TE that already exists in committed history is more disruptive than renumbering a TE that is still confined to a working branch.
- The genesis-freeze rename of `harness-spec.md` to `specs/harness-spec-draft.md` is locked as of TODO 011, with the bot-frozen pCID `bafkreigtaivld55rekcswfj26mo26e267m3ytzgflqb2qcclyiicpfzc6i`.
- The channels branch was forked from `aabfa52` (the Apr 29 review-merged tip on `main`), which predates both the rename and TEs 21, 22, and 23.
- The bot owns the integration work on `ppx/main`; Steve owns the eventual merge into `main`.
- Renumbering only changes the integer label and references to it; it does not require renaming the timestamped file (which is the actual content-address anchor).

## Drafting timestamps as the tiebreaker

Sorted by first-drafted timestamp, the relevant TEs are:

| Integer | Timestamp slug | Title | Where it lives today |
|---|---|---|---|
| TE-21 | 20260429-173520 | Spec doc as promise | `ppx/main` |
| TE-22 | 20260429-175530 | Spec-doc store layout and pCID machinery | `ppx/main` |
| TE-23 | 20260430-064307 | Congruence/convergence duality and pCID framing | `ppx/main` |
| (claimed TE-21 on channels branch) | 20260430-204108 | `grid <pcid>` as a repo-local channel-message carrier | `origin/stevegt/channels-grid-pcid` |
| TE-25 | **20260430-213447** | This TE (numbering collision and harness-spec path) | (drafting now) |

By drafting time, the channel-carrier TE is **earlier** than this TE and would naturally land at TE-24 once both are integrated; this TE then lands at TE-25.

## Alternatives

### DF-25.1 — Which TE keeps the integer 21?

#### Alt-1.A: `ppx/main`'s existing TE-21 keeps the number

The TE titled "Spec doc as promise" remains TE-21. The channel-carrier TE is renumbered to its drafting-time-correct position.

- **Easier**: no rewrites of any committed `ppx/main` history; the spec's TE-index keeps its existing 21. The channel-carrier branch is the one with rewrites, but the only changes are the integer label, the in-spec bullet, and the TE-index row. The timestamped filename does not change. Editing one branch is cheaper than editing many committed-and-merged commits across both branches.
- **Harder**: requires touching the channel-carrier branch's TE doc, harness-spec patch, TODO 009 file, DR-009 file, and DI-009 entry to substitute the correct integer. None of these are large edits, but they must all be consistent before the channels branch merges into `ppx/main`.

#### Alt-1.B: The channel-carrier TE keeps TE-21

The channel-carrier TE remains TE-21. The spec-doc-as-promise TE on `ppx/main` is renumbered.

- **Easier**: the channel-carrier branch already references TE-21 in five files (TE itself, harness-spec patch, DR-009, DI-009, TODO 009, and channels/README.md). Leaving those alone means fewer edits on that branch.
- **Harder**: TE-21 ("Spec doc as promise") is referenced not just in `docs/thought-experiments/README.md` but in committed harness-spec text, TE-22, TE-23, and the framing essay. Renumbering it would require rewriting several merged commits or carrying a forward-referencing patch into a new merge commit. It also breaks the "first-drafted-timestamp anchors the integer" invariant by giving a younger TE the older integer.

#### Alt-1.C: Both TEs keep TE-21, distinguished only by timestamp slug

Treat TE-21 as a non-unique handle and rely on the timestamp slug for disambiguation.

- **Easier**: zero rewrites.
- **Harder**: violates the TE-index README's "stable identifier" promise; collapses the human-friendly integer label; makes spec text and conversation ambiguous; defeats the entire reason integers exist alongside slugs.

### DF-25.2 — What is the new integer for the channel-carrier TE?

Conditional on Alt-1.A.

#### Alt-2.A: The channel-carrier TE becomes TE-25, this TE keeps a tentative TE-24 label

Disrecommended; included only for completeness. The channel-carrier TE was drafted 53 minutes earlier than this TE, so by drafting-time order this assignment is backwards.

#### Alt-2.B: The channel-carrier TE becomes TE-24; this TE becomes TE-25

By strict drafting-time order, the channel-carrier TE was drafted 53 minutes earlier than this one, so it gets the smaller integer.

- **Easier**: directly applies the existing rule. No subjective judgement.
- **Harder**: this TE is being authored with a TE-25 title; the channel-carrier TE's existing references to TE-21 are rewritten to TE-24 during integration. Both edits are mechanical.

#### Alt-2.C: The channel-carrier TE becomes TE-25; this TE keeps a tentative TE-24 label

Treat the reconciliation TE as logically prior because its merge has to land before the channels branch merges, even though it was drafted later.

- **Easier**: matches the order in which the merges actually land on `ppx/main`.
- **Harder**: that "merge order" is a derived fact, not an authoring fact. The README rule is anchored on first drafting, not on merge sequencing.

### DF-25.3 — How is the channels branch reconciled with the renamed spec?

#### Alt-3.A: Translate the patch on the channels branch and force-push it

Rewrite the existing `f9dfd1e Record Channel Carrier V0` commit on `stevegt/channels-grid-pcid` so the harness-spec edit targets `specs/harness-spec-draft.md` instead of `harness-spec.md`, then force-push the branch.

- **Easier**: gives a clean PR diff to review.
- **Harder**: that branch belongs to Steve via Codex; the bot should not be force-pushing other people's branches. Also, force-push to a non-bot-prefixed branch violates the "never force-push" rule.

#### Alt-3.B: Bot translates on a `ppx/` integration twig and merges into `ppx/main`

Cherry-pick `f9dfd1e` onto a `ppx/` working branch, edit it so the harness-spec bullet targets the renamed file, fix the TE integer per DF-25.1 / DF-25.2, then merge `--no-ff` into `ppx/main`. Do not touch the original `stevegt/channels-grid-pcid` branch; let Steve land that branch on `main` whenever he wants. After Steve merges his branch into `main`, the bot reconciles `origin/main` into `ppx/main` again, and any leftover `harness-spec.md`-vs-`specs/harness-spec-draft.md` conflict is resolved in that merge by accepting the renamed file's content.

- **Easier**: respects branch ownership; uses only `ppx/` prefixes for the bot's work; produces a normal `--no-ff` merge.
- **Harder**: introduces some duplicated change content (the bot's translated commit and Steve's original commit will both touch overlapping material). The next `origin/main` → `ppx/main` reconciliation will need to recognize the equivalent content and resolve cleanly.

#### Alt-3.C: Wait until Steve merges the channels branch into `main`, then translate during the next `origin/main` → `ppx/main` reconciliation

Do nothing on `ppx/main` until Steve has decided whether and how to land `stevegt/channels-grid-pcid` on `main`. Then merge `origin/main` into `ppx/main` and resolve the rename conflict at that time.

- **Easier**: avoids any duplication of work; lets Steve's final form drive.
- **Harder**: blocks the channel work entirely from the bot's perspective. The user explicitly asked the bot to "work on channels," so blocking is the wrong default.

#### Alt-3.D: Bot opens a PR back into `stevegt/channels-grid-pcid`

The bot creates a new commit on a `ppx/` branch that, when applied on top of `f9dfd1e`, fixes the path and the TE integer, and proposes that as a PR for Steve to merge.

- **Easier**: keeps Steve's branch as the source of truth and gives him a reviewable patch.
- **Harder**: Steve has not asked for this; a bot PR into a human branch adds review friction that is not yet warranted; and the next `origin/main` → `ppx/main` reconciliation already gives the same review surface naturally.

### DF-25.4 — Where does Open Question #19 land in the renamed spec?

The channel-carrier patch appends a question numbered 19 ("Should the repo-local channels carrier graduate into the canonical PromiseGrid wire format?") to the spec's open-questions list.

#### Alt-4.A: Append at the next free integer in the renamed spec

Renumber to whatever the next free integer is in `specs/harness-spec-draft.md` (the renamed file may or may not have the same numbering tail as the original). Append the bullet there.

- **Easier**: matches existing convention.
- **Harder**: requires reading the renamed file's open-question list to confirm the next free integer.

#### Alt-4.B: Drop the question for now and let Steve add it via the channels branch when he merges to `main`

Bot does not carry the open question on `ppx/main`; it lands organically when Steve merges his branch into `main` and the next `origin/main` → `ppx/main` reconciliation pulls it in.

- **Easier**: avoids the bot authoring open-question text on Steve's behalf.
- **Harder**: the question is genuinely useful for the channel work the bot is about to do; deferring it leaves a gap in the spec while channel work is in progress.

#### Alt-4.C: Carry the question on `ppx/main` as Open Question #19, scoped to repo-local

Add the bullet at the next free integer with explicit scoping: "For the repo-local `channels/` experiment only…" plus a pointer to TE-24 (the channel-carrier TE) and DR-009.

- **Easier**: keeps the question visible while channel work continues.
- **Harder**: when Steve merges his branch into `main` later, his version of the question will collide with this one and need to be deduped.

## Scenario analysis

### S1 — Channel work starts immediately after this TE lands

The bot needs to begin channel work today.

- **Alt-1.A + Alt-2.B + Alt-3.B + Alt-4.C**: bot translates the channel-carrier commit onto a `ppx/` twig, renumbers TE-21 to TE-24, fixes the path to `specs/harness-spec-draft.md`, adds the open question scoped to repo-local. Channel work can begin as soon as that lands.
- **Alt-3.C**: bot is blocked.
- **Alt-1.B**: bot has to renumber three TEs of committed history for low gain, then begin work.

### S2 — Steve merges his branch into `main` next week

Steve eventually lands `stevegt/channels-grid-pcid` on `main`. The bot then reconciles `origin/main` into `ppx/main`.

- **Alt-3.B + Alt-4.C**: the merge will see overlapping content. Bot resolves by keeping the renamed-path version; the open question appears once on each side and is deduped to one. Mechanical, low-risk.
- **Alt-3.D**: Steve's branch already carries the bot's translation, so the merge is clean.
- **Alt-3.C**: bot has done no channel work in the meantime.

### S3 — A future TE refers to the channel-carrier TE by integer

Someone writes "see TE-21" or "see TE-24" in a future TE or essay.

- **Alt-1.A + Alt-2.B**: the integer is unambiguous and matches the drafting-time invariant.
- **Alt-1.C**: the reference is ambiguous and the reader has to disambiguate by slug or context.

### S4 — TE-NN integer is needed across forks

Two long-running branches each draft their own TE-NN before either merges.

- The drafting-time rule is durable: whichever TE was drafted earlier gets the smaller integer, and the other is renumbered at merge time. This is exactly the workflow this TE is recording.

## Conclusions

1. **TE numbering invariant**: integers are anchored on first-drafted timestamp, not on merge order or branch of origin. When two branches collide on the same integer, the later-drafted TE is renumbered, never the earlier one. This generalizes from the present case.
2. **For the present case**: the spec-doc-as-promise TE on `ppx/main` keeps TE-21. The channel-carrier TE becomes TE-24 (it was drafted at 20260430-204108, before this reconciliation TE at 20260430-213447). This reconciliation TE is TE-25.
3. **Path translation**: the bot does not edit Steve's branch. Instead, the bot brings the channel-carrier change into `ppx/main` via a `ppx/` integration twig that translates the patch onto `specs/harness-spec-draft.md` and applies the renumbering. Steve's `stevegt/channels-grid-pcid` branch is left untouched and remains his to land on `main` whenever he chooses.
4. **Future reconciliation**: when Steve eventually merges his branch into `main`, the next `origin/main` → `ppx/main` merge will see overlapping content. The resolution rule is "prefer the renamed-path version; dedupe duplicate spec bullets and TE-index rows; keep `ppx/main`'s renumbered integer."
5. **Open Question #19 is carried on `ppx/main` immediately**, scoped to repo-local channels and pointing at TE-24 and DR-009. Duplicate-removal at next reconciliation is acceptable.

## Implications

- Future TEs that draft in parallel branches must be renumbered at merge time using drafting-time order. The bot will check `docs/thought-experiments/README.md` before assigning an integer; if the highest committed integer on `ppx/main` is N, a fresh draft is N+1 unless a higher-numbered TE has been drafted on another branch and is still in flight.
- The path-rename pattern (`harness-spec.md` → `specs/harness-spec-draft.md`) will recur whenever Steve's branches forked from before the rename land. The bot's reconciliation rule is "translate to the renamed path on `ppx/main`; let `main` keep whatever shape Steve chooses; reconcile at the next `origin/main` → `ppx/main` merge."
- Open spec questions added on parallel branches should be scoped explicitly when they cite branch-local work (e.g., "for the repo-local `channels/` experiment only"), to make later deduplication unambiguous.

## Decision Forks (DFs)

### DF-25.1 — Which TE keeps the integer 21?

The TE labeled TE-21 on `ppx/main` is "Spec doc as promise" (drafted 20260429-173520). The TE labeled TE-21 on `origin/stevegt/channels-grid-pcid` is "`grid <pcid>` as a repo-local channel-message carrier" (drafted 20260430-204108). Only one can keep the integer.

- **Alt-1.A**: `ppx/main`'s existing TE-21 keeps the integer; channel-carrier TE is renumbered. Recommended.
- **Alt-1.B**: channel-carrier TE keeps the integer; `ppx/main`'s TE-21 is renumbered. Disrecommended (would rewrite committed history and violate the drafting-time invariant).
- **Alt-1.C**: both keep TE-21, distinguished by slug only. Disrecommended (defeats the integer-as-handle convention).

### DF-25.2 — What integer does the channel-carrier TE become?

Conditional on Alt-1.A.

- **Alt-2.A**: TE-25. Disrecommended (drafted earlier than this reconciliation TE).
- **Alt-2.B**: TE-24. Recommended (matches drafting-time order).
- **Alt-2.C**: TE-25, with this TE staying at TE-24. Disrecommended (introduces double-renumbering churn).

### DF-25.3 — How is the channels branch reconciled with the renamed spec?

- **Alt-3.A**: bot force-pushes Steve's branch. Disrecommended (violates branch ownership and never-force-push rule).
- **Alt-3.B**: bot integrates via a `ppx/` twig that translates the path; Steve's branch is untouched. Recommended.
- **Alt-3.C**: wait until Steve merges to `main`; do nothing on `ppx/main` until then. Disrecommended (blocks channel work).
- **Alt-3.D**: bot proposes a PR into Steve's branch. Disrecommended (unsolicited; the next `origin/main` → `ppx/main` merge already provides the review surface).

### DF-25.4 — Where does Open Question #19 land?

- **Alt-4.A**: append at the next free integer with the channel-carrier branch's existing wording. Acceptable but slightly worse than C.
- **Alt-4.B**: drop the question on `ppx/main`; let it land via Steve's branch later. Disrecommended (leaves a gap during active channel work).
- **Alt-4.C**: carry as Open Question #19 with explicit "repo-local only" scoping and a pointer to TE-24 and DR-009. Recommended.

## Recommended set

**Alt-1.A + Alt-2.B + Alt-3.B + Alt-4.C.**

Rationale: this set is the only set that simultaneously honors the drafting-time invariant for TE numbering (Alt-1.A, Alt-2.B), respects branch ownership and the never-force-push rule (Alt-3.B), and avoids leaving a gap in the spec's open-questions list while active channel work is underway (Alt-4.C).

## Decision status

**Locked 2026-04-30 by Steve Traugott.** Recommended set adopted: **Alt-1.A + Alt-2.B + Alt-3.B + Alt-4.C.**

- DF-25.1 = Alt-1.A: `ppx/main`'s existing TE-21 ("Spec doc as promise") keeps the integer.
- DF-25.2 = Alt-2.B: the channel-carrier TE becomes TE-24.
- DF-25.3 = Alt-3.B: bot integrates via a `ppx/` twig that translates the path; `origin/stevegt/channels-grid-pcid` is not touched.
- DF-25.4 = Alt-4.C: Open Question #19 is carried on `ppx/main` with explicit "repo-local channels only" scoping and pointers to TE-24 and DR-009.

## Implications for follow-on work

- After this TE lands, the next bot work item is to translate `f9dfd1e Record Channel Carrier V0` from `origin/stevegt/channels-grid-pcid` onto a `ppx/` integration twig with the integer renumbered (TE-21 → TE-24), the spec edit retargeted to `specs/harness-spec-draft.md`, the TE-index row updated, the open-question integer set, and the open-question text scoped to repo-local. Then the bot can begin actual channel work.
- TODO 011 and the `tools/spec` workflow stay unchanged. The spec on `ppx/main` will be re-frozen once the channel-carrier change is integrated, producing a new pCID for the post-channels spec.
- A follow-on TE (TE-26 or later) may be needed if Steve's eventual merge into `main` lands a meaningfully different version of the channel-carrier text.
