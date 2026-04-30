# TODO 009 - ppx/main as long-lived bot integration branch

Establish `ppx/main` as a long-lived bot integration branch and update `AGENTS-ppx.md` to document the new bot workflow: bot creates `ppx/{twig}` working branches off `ppx/main`, merges them into `ppx/main` (no-ff), pushes `ppx/main` to origin, and deletes the `ppx/{twig}` branches both locally and on origin. Steve (via Codex) merges `origin/ppx/main` into `main` when ready. The bot keeps `ppx/main` current by periodically merging `origin/main` INTO `ppx/main`, never via rebase.

## Subtasks

- [x] 009.1 Create `ppx/main` on origin, branched off current `main`.
- [x] 009.2 Update `AGENTS-ppx.md` to document the new branch model: add `ppx/main` to the branch-model section; update Kind 1/2/3/4 procedures so working branches are merged into `ppx/main` and deleted instead of pushed to origin and left.
- [x] 009.3 Update `AGENTS-ppx.md` "First action of every session" to fetch and check `ppx/main` alongside `main`.
- [x] 009.4 Update `AGENTS-ppx.md` "Things that are forbidden" to add an explicit no-force-push rule covering `ppx/main` and `ppx/{twig}` working branches.
- [x] 009.5 Update `AGENTS-ppx.md` glossary entry for `pCID` to use the locked `Protocol CID` definition (per Steve's correction; `pCID` is the hash of a spec document, not a message or payload).

## Decision Intent Log

ID: DI-009-20260429-173358
Date: 2026-04-29 17:33:58
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: The bot uses a long-lived `ppx/main` integration branch. Working branches `ppx/{twig}` are created off `ppx/main`, merged back into `ppx/main` with `--no-ff`, and then deleted both locally and on origin. The bot pushes `ppx/main` to origin after each integration. Steve (via Codex) merges `origin/ppx/main` into `main`. The bot keeps `ppx/main` current by periodically merging `origin/main` INTO `ppx/main` (never the other direction; never via rebase, since rebase would require force-push).
Intent: Reduce cleanup work on Steve's side. Under the prior workflow, every `ppx/{twig}` branch left a residue on origin that Steve had to delete manually after merging to `main`. With `ppx/main` as the integration target, only one bot-side branch ever appears on origin, and the bot owns its lifecycle.
Constraints: The bot never pushes to `main`. The bot never force-pushes any branch, including `ppx/main` itself. The bot never rebases `ppx/main`. When `origin/main` advances, the bot brings the new commits into `ppx/main` by `git merge --no-ff origin/main` (a merge commit, not a rewrite). Working branches `ppx/{twig}` may be pushed to origin if the bot wants a backup or wants to share work-in-progress, but the integration target is `ppx/main`, and the working branch is deleted (locally and on origin) once merged.
Affects: `AGENTS-ppx.md` (branch-model section, Kind 1-4 procedures, first-action checklist, forbidden-actions list); the bot's per-session workflow; how `ppx/{twig}` working branches are named and disposed of; how `ppx/main` is kept current with `main`; what Codex sees on origin.
Linked DR: none directly (chat-instructed workflow change). The closest related DRs are DR-001 (bot identity), DR-002 (drop require-PR), and DR-003 (review style); this DI extends them by adding the integration-branch layer between `ppx/{twig}` and `main`. Does NOT supersede DI-001-20260428-195700/195701/195702 \u2014 those remain in force.

ID: DI-009-20260429-173359
Date: 2026-04-29 17:33:59
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: `pCID` is the abbreviation for "Protocol CID", not "Promise Content ID". A pCID is the content hash of a spec document that defines a wire protocol; it is analogous to a TCP/UDP port number but with no central registry, because the spec's hash IS the port number. A pCID is NOT the hash of any particular message, payload, or promise body.
Intent: Lock the canonical definition so that future docs and bot output do not drift back to the old "Promise Content ID" gloss. The specs/harness-spec-draft.md \u00a71 note already carries this definition; AGENTS-ppx.md glossary previously used the wrong gloss and is now corrected.
Constraints: Future references to `pCID` in any artifact (docs, TEs, DRs, commit messages, chat) must use the "Protocol CID" form. The bot must not write "Promise CID" or "Promise Content ID" anywhere.
Affects: `AGENTS-ppx.md` glossary; future documentation that defines `pCID`; bot's vocabulary use in chat and commits.
Linked DR: none (chat-directed correction; same status as DI-004 burden\u2192assertion).

## Notes

- `ppx/main` was created on origin as part of subtask 009.1 before the AGENTS-ppx.md update landed. The first work to land on the new `ppx/main` was the DR-006 normalization (twig `ppx/dr-006-normalize`, merge `4bb80db`).
- Subtasks 009.2-009.5 land together on twig `ppx/agents-ppx-workflow-update`.
