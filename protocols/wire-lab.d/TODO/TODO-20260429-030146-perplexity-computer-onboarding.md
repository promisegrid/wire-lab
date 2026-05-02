# TODO 001 - Perplexity Computer onboarding

Establish how Perplexity Computer (an LLM-driven agent operating on Steve's behalf) participates in this repo's mob-consensus-style workflow. This TODO records the locked decisions that bootstrap the bot's identity, branch namespace, branch-protection posture, and review style.

The decisions below were made in chat on 2026-04-28; this TODO file captures them as Decision Intent Log entries so that future work has citable DI IDs.

## Subtasks

- [x] 001.1 Lock the bot's email identity and `<user>` branch prefix.
- [x] 001.2 Lock the GitHub branch-protection posture for `main`.
- [x] 001.3 Lock the bot's review/convergence style.
- [x] 001.4 Steve to remove the "Require a pull request before merging" rule from the `main` ruleset in repo Settings. (Done; confirmed by Steve in chat 2026-04-29.)

The two follow-on items originally listed as 001.5 and 001.6 have been moved out of this TODO so that TODO 001 closes cleanly when its bootstrap subtasks land. They now live as:

- TODO 006 — `protocols/wire-lab.d/TODO/TODO-20260429-165252-backfill-di-provenance-harness-spec.md` (originally 001.5).
- TODO 007 — `protocols/wire-lab.d/TODO/TODO-20260429-165253-backfill-dr-harness-spec-section-11.md` (originally 001.6).

## Decision Intent Log

ID: DI-001-20260428-195700
Date: 2026-04-28 19:57:00
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: The Perplexity Computer bot identifies as `stevegt+ppx@t7a.org (stevegt-via-perplexity)`. All bot-authored branches use the prefix `ppx/{twig}` (not `stevegt+ppx/{twig}` as the mob-consensus default would derive from the email local-part).
Intent: Make it unambiguous that the bot is acting on Steve's behalf (subaddress routes to Steve's inbox; parenthetical names the actor) while keeping a short, memorable branch prefix that matches Steve's stated preference. The `ppx/` prefix is a concession against the strict mob-consensus convention because the literal `stevegt+ppx/` prefix is awkward.
Constraints: The bot must set `git config user.email = stevegt+ppx@t7a.org` and `user.name = stevegt-via-perplexity` in this repo's clone before committing. Branch names must always start `ppx/`.
Affects: bot's `.git/config` in any clone; all future branch names from the bot; all `Author` and `Asked by` fields in DR/DI records the bot writes; commit author metadata.
Linked DR: DR/DR-001-20260428-195700-bot-identity.md

ID: DI-001-20260428-195701
Date: 2026-04-28 19:57:01
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Drop the "Require a pull request before merging" rule from the `main` branch ruleset on github.com/promisegrid/wire-lab. Keep the "Restrict who can push to matching branches → stevegt only" rule, plus "Block force pushes" and "Restrict deletions". Steve merges by direct push to `main`; the bot pushes only to `ppx/{twig}` branches.
Intent: Avoid GitHub-specific lockin in the workflow. PromiseGrid will eventually replace GitHub as the substrate; the merge mechanism should be plain `git push` to a protected branch, not a forge-specific PR ceremony. The "only Steve can push to main" rule maps cleanly to PromiseGrid's eventual "the canonical pointer follows Steve's signing key" semantics (see `specs/harness-spec-draft.md` §10a.8).
Constraints: Steve must remove the "Require PR" rule in repo Settings → Rules → Rulesets. The bot has no permission to modify rulesets (PAT scope deliberately excludes admin permissions). The bot must never attempt to push to `main` directly; only to `ppx/{twig}` branches.
Affects: github.com/promisegrid/wire-lab branch ruleset on `main`; the bot's push targets; the convergence ceremony (no GitHub PR clicks).
Linked DR: DR/DR-002-20260428-195700-drop-require-pr.md

ID: DI-001-20260428-195702
Date: 2026-04-28 19:57:02
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Adopt review style α: the bot pushes `ppx/{twig}` to origin, then writes a Decision Request (DR) on the same branch with `State: open` and `Question: review and converge ppx/{twig}`. Steve reviews using his own tools (Codex, local diff viewers, etc.) — not the GitHub PR UI. When satisfied, Steve merges `ppx/{twig}` into `main` locally and pushes; the bot observes the merge by polling and appends `State: implemented` then `State: closed` to the DR on a follow-up branch.
Intent: Keep the decision-and-review record entirely in-repo (DR files are version-controlled, append-only, citable by ID) instead of relying on GitHub PR descriptions and comments (out-of-repo, mutable, GitHub-locked). Reviewing by local diff plus DR file matches the mob-consensus model and ports cleanly to PromiseGrid.
Constraints: Every non-trivial bot change must include a DR file authored on the same branch. Trivial fixes (typo, broken link, formatting) may use a `ppx/chore-{twig}` branch and skip the DR if Steve has separately approved a "chore exception" pattern (not approved as of this DI). The bot must not open GitHub PRs; if it accidentally does, Steve closes them without merging.
Affects: every future branch from the bot; the bot's task-completion handoff format; the DR/ directory layout.
Linked DR: DR/DR-003-20260428-195700-review-style.md

## Notes

This bootstrap was created by Perplexity Computer on `ppx/dr-001-bootstrap` and pushed for Steve to review and merge. The DR files for DR-001, DR-002, DR-003 are committed alongside this TODO file. Steve answered all three DR questions in chat before any file was written, so each DR is created with `State: decided` from the start.

The "TODO 025 migration" referenced in `AGENTS.md` is from a prior repo and does not apply here; this repo starts TODO numbering at 001.
