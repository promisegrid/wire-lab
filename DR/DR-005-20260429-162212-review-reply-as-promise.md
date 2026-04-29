# DR-005 - Review reply as promise

DR-ID: DR-005-20260429-162212
Date: 2026-04-29 16:22:12
Asked by: stevegt@t7a.org (Steve Traugott)
State: decided
Question: In promise-theory terms, how should Steve's durable review reply to a bot-authored `ppx/{twig}` proposal branch be recorded on `main`?
Why this blocks progress: The previously merged `contest` framing captured the durability goal, but not the right semantics for Steve's role. Steve's actual request-changes reply is a conditional commitment about his own future review behavior, so the repo needs a clearer promise-shaped record before more branches accumulate under the wrong model.
Affects: `proposals/pending/`; Steve's review workflow on `main`; future bot/agent consumption of review outcomes; interpretation of the existing `ppx/dr-001-bootstrap` review artifact.
Unblocks: `TODO/003-review-reply-as-promise.md` (all subtasks); the immediate durable reply for `ppx/te-20260428-202400-promise-stack-ordering`; the retrofit of `ppx/dr-001-bootstrap`.
Waiting on: DI-003-20260429-162212

## Candidate alternatives considered

- (a) Keep using `contest-v1`. Durable, but semantically mismatched for Steve's conditional review commitment.
- (b) Record a final approve/reject decision. Clear, but too strong for a request-changes state where the proposal remains pending.
- (c) Record a simple prose review message whose core content is a conditional promise to review again after named conditions are satisfied. Selected.
- (d) Respond with a counter-proposal. Useful when Steve wants to replace the proposal outright, but too heavy for ordinary request-changes review.

## Decision

Choose (c). Steve's durable review reply is a prose review message whose load-bearing sentence is a conditional promise: if the proposer publishes a revised branch that satisfies named conditions, Steve promises to review it again; until then, Steve does not promise to merge the current revision.

These review messages live under `proposals/pending/<proposal-id>/review-<timestamp>-steve-traugott.md`. The reviewed proposal remains in `pending/` while those conditions are outstanding. The earlier `contest` artifact for `ppx/dr-001-bootstrap` remains in history, but a new promise-shaped review message supersedes it operationally.

## Linked DI

DI-003-20260429-162212 (in `TODO/003-review-reply-as-promise.md`)

## Related commits

To be appended after the DI, TE, and review-message artifacts are committed on `main`.

## Last updated

2026-04-29 16:22:12
