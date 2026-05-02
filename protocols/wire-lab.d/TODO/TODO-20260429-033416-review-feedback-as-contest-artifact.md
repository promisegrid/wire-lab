# TODO 002 - Durable review feedback as contest artifact

Record Steve's request-changes review feedback for bot-authored proposal branches in a durable, in-repo artifact that fits the PromiseGrid proposal vocabulary instead of leaving it as ephemeral chat.

## Subtasks

- [x] 002.1 Run a thought experiment comparing the viable artifact shapes.
- [x] 002.2 Lock the durable review-feedback model in a Decision Intent record.
- [x] 002.3 Publish the first `contest-v1` review artifact for `ppx/dr-001-bootstrap` on `main`.

## Decision Intent Log

ID: DI-002-20260429-033208
Date: 2026-04-29 03:32:08
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Steve's request-changes feedback for a bot-authored proposal branch is persisted as a durable prose `contest-v1` artifact on `main` under `proposals/pending/<proposal-id>/contest-<timestamp>-steve-traugott.md`. The first such artifact targets `ppx/dr-001-bootstrap`.
Intent: Fit Steve's review feedback into the evolving PromiseGrid discourse model so the durable record lives in-repo, can be read by humans or agents later, and avoids devolving into an ad hoc repo-root note or a GitHub-specific PR comment trail.
Constraints: The artifact must identify the proposal branch and, when practical, the reviewed commit SHA. Prose-only content is acceptable because the proposal envelope is intentionally not locked yet. Chat may still notify the bot, but `main` is the durable source of truth for the contest record. The first artifact path is normalized from `ppx/dr-001-bootstrap` to `proposals/pending/ppx-dr-001-bootstrap/`.
Affects: `proposals/pending/`; Steve's review workflow for `ppx/{twig}` branches; future bot polling/parsing of review outcomes; how request-changes feedback is preserved across GitHub migration.
Linked DR: DR/DR-004-20260429-033208-review-feedback-as-contest-artifact.md

## Notes

Steve approved this model in chat after reviewing the alternatives narrowed by `TE-20260429-033208-review-feedback-as-contest-artifact.md`.
