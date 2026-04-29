# DR-004 - Durable review feedback as contest artifact

DR-ID: DR-004-20260429-033208
Date: 2026-04-29 03:32:08
Asked by: stevegt@t7a.org (Steve Traugott)
State: decided
Question: When Steve wants to request changes to a bot-authored `ppx/{twig}` branch, how should that review feedback be recorded in-repo?
Why this blocks progress: Steve explicitly wants the review reply for `ppx/dr-001-bootstrap` written into a file on `main`. Without a locked artifact shape, the repo would gain a one-off note that does not fit the evolving PromiseGrid proposal vocabulary or scale to future proposal review traffic.
Affects: `proposals/pending/`; Steve's review workflow on `main`; future bot/agent consumption of review outcomes; migration away from GitHub-specific review trails.
Unblocks: `TODO/002-review-feedback-as-contest-artifact.md` (all subtasks); the immediate request-changes response to `ppx/dr-001-bootstrap`.
Waiting on: DI-002-20260429-033208

## Candidate alternatives considered

- (a) Leave request-changes feedback in chat only. Fastest, but not durable, not branch-addressable in-repo, and not aligned with the proposal vocabulary.
- (b) Write an ad hoc markdown note at the repo root on `main`. Durable, but it bypasses the queue shape emerging in `harness-spec.md` §10a.10 and will not scale when multiple proposals are pending.
- (c) Record the feedback as a prose `contest-v1` artifact on `main` under `proposals/pending/<proposal-id>/`. Selected.

## Decision

Choose (c). Steve's request-changes review feedback is recorded as a prose `contest-v1` artifact on `main` under `proposals/pending/<proposal-id>/contest-<timestamp>-steve-traugott.md`.

For the first use of this pattern, `ppx/dr-001-bootstrap` is normalized to `proposals/pending/ppx-dr-001-bootstrap/`, and the artifact names the exact reviewed commit SHA. Chat may still carry a courtesy notification, but the durable review record is the in-repo contest artifact.

## Linked DI

DI-002-20260429-033208 (in `TODO/002-review-feedback-as-contest-artifact.md`)

## Related commits

To be appended after this decided record and its first contest artifact are committed on `main`.

## Last updated

2026-04-29 03:32:08
