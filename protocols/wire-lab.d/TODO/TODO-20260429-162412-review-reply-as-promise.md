# TODO 003 - Review reply as promise

Record Steve's durable review replies for bot-authored proposal branches as simple prose messages whose load-bearing content is a conditional promise, rather than as `contest-v1` artifacts.

## Subtasks

- [x] 003.1 Run a thought experiment comparing the viable review-reply shapes.
- [x] 003.2 Lock the promise-shaped review-reply model in a Decision Intent record.
- [x] 003.3 Publish the first promise-shaped review reply for `ppx/te-20260428-202400-promise-stack-ordering` on `main`.
- [x] 003.4 Publish a superseding promise-shaped review reply for `ppx/dr-001-bootstrap` on `main`.

## Decision Intent Log

ID: DI-003-20260429-162212
Date: 2026-04-29 16:22:12
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Steve's durable review reply for a bot-authored proposal branch is a prose review message on `main` whose core promise is conditional: if the proposer publishes a revised branch that satisfies named conditions, Steve promises to review it again; until then, Steve does not promise to merge the current revision. These review messages live under `proposals/pending/<proposal-id>/review-<timestamp>-steve-traugott.md`.
Intent: Keep the review flow aligned with the repo's promise-stack model by framing Steve's reply as a simple promise message, not as a peer-style contest. This preserves a durable in-repo record while keeping the semantics easy for humans and future agents to read.
Constraints: The review message must name the proposal branch and, when practical, the reviewed commit SHA. The semantics live in the `Promise:` paragraph rather than in a special envelope schema. The proposal remains in `proposals/pending/` unless Steve makes a final rejection or merge decision. Earlier `contest` artifacts are not deleted; a newer review message supersedes them operationally.
Affects: `proposals/pending/`; Steve's review workflow for `ppx/{twig}` branches; future bot polling/parsing of review outcomes; interpretation of the earlier `ppx/dr-001-bootstrap` review artifact.
Supersedes: DI-002-20260429-033208
Linked DR: DR/DR-005-20260429-162212-review-reply-as-promise.md

## Notes

This correction was approved in chat after re-framing the problem around the repo's promise-stack model and the conditional nature of a request-changes reply.
