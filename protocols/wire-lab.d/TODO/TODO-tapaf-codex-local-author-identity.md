# TODO-tapaf: Codex local author identity

## Prior aliases

None. This TODO was minted after the proquint-handle migration.

## Status

Implemented. This TODO records the local Codex overlay change that stops
impersonating Steve by default and instead uses the real local git author
identity configured in the clone. It leaves the existing `ppx` bot identity and
branch-prefix rules unchanged. Source: `DI-nufit`.

## Decision Intent Log

ID: DI-nufit
Date: 2026-05-23 23:27:04
Status: active
Author: jessica@t7a.org (Jessica)
Decision: Change `AGENTS-codex.md` so local Codex sessions use the actual local
git author identity configured in the current clone, rather than always acting
as Steve. Keep the existing `ppx` bot identity, bot branch-prefix rules, and
Steve-only `main` push rules unchanged.
Intent: The local Codex overlay should support real multi-author work in one
repo without falsifying commit attribution. The local operator may be Steve,
Jessica, or another authorized human using this clone. The local overlay should
therefore require an explicitly configured real local identity and should stop
rewriting every local action into Steve's authorship.
Constraints: Do not change the Perplexity bot identity locked by
`DI-001-20260428-195700`. Do not relax the `main` branch protection or
Steve-only direct-push posture locked by `DI-001-20260428-195701`. Limit this
change to local Codex identity guidance; do not turn it into a broad branch or
governance rewrite. Local sessions must never use the bot identity unless the
work is actually being performed under the bot workflow.
Affects: `AGENTS-codex.md`; `protocols/wire-lab.d/TODO/TODO.md`;
`DR/DR-pijum-codex-local-author-identity.md`;
`docs/thought-experiments/TE-puzud-codex-local-author-identity.md`.

## Subtasks

- [x] tapaf.1 Compare the existing Steve-only local Codex identity rule against
  alternatives that use the actual local git author identity or a separate
  Codex-only identity.
- [x] tapaf.2 Record the selected local-author rule as a DI and linked DR/TE.
- [x] tapaf.3 Update `AGENTS-codex.md` so it requires a real configured local
  author identity and explicitly forbids using the `ppx` bot identity for local
  work.
