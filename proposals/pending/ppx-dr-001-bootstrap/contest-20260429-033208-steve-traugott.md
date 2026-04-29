# contest-v1 - Request changes for `ppx/dr-001-bootstrap`

Contest-ID: CONTEST-20260429-033208
Date: 2026-04-29 03:32:08 UTC
From: stevegt@t7a.org (Steve Traugott)
Target proposal: `ppx/dr-001-bootstrap`
Target commit: `db9cd416b35088988be2f0a1a8fe76a95c9dbcb8`
Target promiser: stevegt+ppx@t7a.org (stevegt-via-perplexity)
State: contested
Vocabulary: prose `contest-v1`
Decision provenance: `DI-002-20260429-033208`
Related DR: `DR/DR-004-20260429-033208-review-feedback-as-contest-artifact.md`

## Claim contested

`ppx/dr-001-bootstrap` is not ready to merge yet.

## Conflicting observations

1. The branch locks three non-trivial decisions while also listing multiple plausible alternatives, but it adds no TE artifact under `docs/thought-experiments/`.
2. `DR/DR-002-20260428-195700-drop-require-pr.md` uses free-text `Waiting on` content instead of either a DI ID or a person identity in `email (FirstName)` format.
3. `TODO/TODO.md` marks TODO 001 complete while subtask `001.4` remains open, and while `001.5` and `001.6` are still listed under the same TODO.

## Requested changes

- Add the TE artifact(s) needed to cover the bot identity, branch-protection, and review-style decisions before those decisions remain locked in DI form.
- Fix `Waiting on` in `DR/DR-002-20260428-195700-drop-require-pr.md` so it uses either a DI ID or `stevegt@t7a.org (Steve Traugott)`.
- Fix TODO completion state so TODO 001 is not marked done while open subtasks remain; if `001.5` and `001.6` are follow-on work, move them into separate TODOs.
- Optionally rebase or refresh the branch onto current `origin/main` so the diff is easier to review.

## Disposition

This proposal remains pending. Re-review after the bot pushes an updated branch.
