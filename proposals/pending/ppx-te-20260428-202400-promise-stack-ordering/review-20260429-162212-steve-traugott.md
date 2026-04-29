# Review message for `ppx/te-20260428-202400-promise-stack-ordering`

Review-ID: REVIEW-20260429-162212
Date: 2026-04-29 16:22:12 UTC
From: stevegt@t7a.org (Steve Traugott)
Target proposal: `ppx/te-20260428-202400-promise-stack-ordering`
Target commit: `0ec32d9e7d0025bfebe2592c9b863b0000f850c8`
Target promiser: stevegt+ppx@t7a.org (stevegt-via-perplexity)
Queue state: pending
Linked DI: `DI-003-20260429-162212`
Related DR: `DR/DR-005-20260429-162212-review-reply-as-promise.md`

Promise:
I, Steve Traugott, promise to review a revised form of `ppx/te-20260428-202400-promise-stack-ordering` after the proposer publishes an updated branch that satisfies the conditions listed below. Until then, I do not promise to merge the current revision.

## Observed blockers

1. The branch is non-trivial TE work, but it adds no DR file at all.
2. The TE is not tracked in any `TODO/*.md`, even though the TE protocol requires required TEs to be tracked there.
3. The unresolved DF questions remain free-floating and still reference the placeholder path `TODO/00X-promise-stack-ordering.md` rather than a real TODO/DR scaffold.

## Conditions for re-review

- Add a DR for the non-trivial branch and link it to the promise-stack-ordering work.
- Add a real TODO file that tracks the TE and the future DF/DI work.
- Replace the `TODO/00X` placeholder with the real TODO path and anchor the unresolved questions to concrete DR/DI planning.
- Optionally refresh or rebase the branch onto current `origin/main` so the diff is clean and reviewable.

## Disposition

This proposal remains pending. Re-review after the proposer publishes an updated branch that satisfies the listed conditions.
