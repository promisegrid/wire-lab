# TE-19: Branch-protection posture for `main`

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-20260429-165102

(First drafted 2026-04-29 16:51:02 UTC, retrospective-style: this TE was authored after the corresponding DI was already locked, in response to the conditions on `review-20260429-162212-steve-traugott.md` for `ppx/dr-001-bootstrap`.)

## Status

decided

## Decision under test

For a repo that will eventually migrate from GitHub to PromiseGrid, what branch-protection posture on `main` minimizes forge-specific lockin while keeping the human principal in control of merges?

## Setup

- `github.com/promisegrid/wire-lab` `main` is the canonical pointer. PromiseGrid will eventually replace GitHub as the substrate.
- The human principal is Steve. The agent is Perplexity Computer (the bot), operating on Steve's behalf.
- Bot pushes only to `ppx/{twig}` branches; never to `main`.
- The bot's PAT scope is deliberately limited; admin (ruleset-modification) permission is excluded.

## Alternatives

- **Alt-A (require GitHub PR):** Keep "Require a pull request before merging" plus "Require approvals". **Easier:** uses GitHub's built-in review UI. **Harder:** ties the workflow to GitHub-specific PR semantics. PromiseGrid migration would have to invent a `git push`-shaped equivalent. The PR ceremony also lives outside the repo (PR descriptions, comments, approvals are not version-controlled).

- **Alt-B (Git-native push restriction):** Drop "Require a pull request before merging". Keep "Restrict who can push to matching branches → stevegt only", "Block force pushes", and "Restrict deletions". Steve merges by `git merge --no-ff` plus `git push origin main`. **Easier:** the merge mechanism is plain Git; it ports cleanly to any future substrate where `main` is just a signed pointer. The "only Steve can push" rule maps to PromiseGrid's eventual "the canonical pointer follows Steve's signing key" semantics (`protocols/wire-lab.d/specs/harness-spec-draft.md` §10a.8). Decision-and-review record stays in-repo as DR/DI files. **Harder:** Steve has no GitHub UI walking him through the merge; he runs `git` commands locally.

- **Alt-C (no protection):** Anyone with write access can push `main`. **Easier:** simplest. **Harder:** removes the only mechanism that keeps the canonical pointer following Steve's intent. Bot bugs or stray pushes can corrupt `main`.

- **Alt-D (signed-commit requirement):** Require that all commits on `main` be GPG- or SSH-signed by Steve's key. **Easier:** signed pointer matches PromiseGrid's eventual model exactly. **Harder:** GitHub's signed-commit enforcement has UX gaps; signing every merge commit adds friction. Could be added later as an additive policy on top of Alt-B.

## Selection

Alt-B. The locked decision is recorded as `DI-001-20260428-195701` in `protocols/wire-lab.d/TODO/TODO-20260429-030146-perplexity-computer-onboarding.md`, and the supporting DR is `DR/DR-002-20260428-195700-drop-require-pr.md`. Steve confirmed in chat on 2026-04-29 that the "Require PR" rule has been removed from the `main` ruleset.

Alt-D may be revisited later as an additive layer — leaving Alt-B as the base posture and adding signed-commit enforcement when PromiseGrid's signing-key model is closer to landing.

## Decision status

`decided` — DI-001-20260428-195701 is already active. This TE is retrospective and exists to satisfy the TE-protocol requirement that non-trivial decisions have a corresponding TE artifact.
