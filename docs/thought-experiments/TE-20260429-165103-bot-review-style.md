# TE-20: Bot review style

*Thought experiment, part of the [PromiseGrid Wire Lab](../../specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-20260429-165103

(First drafted 2026-04-29 16:51:03 UTC, retrospective-style: this TE was authored after the corresponding DI was already locked, in response to the conditions on `review-20260429-162212-steve-traugott.md` for `ppx/dr-001-bootstrap`.)

## Decision under test

How should the bot's proposal-and-review ceremony be shaped so that the decision-and-review record stays in-repo, citable, and forge-portable, while still letting the human principal (Steve) review with whatever tools they prefer?

## Setup

- The bot pushes to `ppx/{twig}` branches; Steve merges to `main`.
- Branch protection is Alt-B from TE-19 (no GitHub PR requirement).
- PromiseGrid is the long-term substrate; the workflow must port to it cleanly.
- `specs/harness-spec-draft.md` already defines DR/DI as the in-repo record format.

## Alternatives

- **Alt-α (DR-on-branch + local review):** Bot pushes `ppx/{twig}` to origin, then writes a Decision Request (DR) on the same branch with `State: open` and `Question: review and converge ppx/{twig}`. Steve reviews using his own tools (Codex, local diff viewers, etc.). When satisfied, Steve merges `ppx/{twig}` into `main` locally and pushes; the bot observes the merge by polling and appends `State: implemented` then `State: closed` to the DR on a follow-up branch. **Easier:** decision-and-review record is fully in-repo (DR files are version-controlled, append-only, citable by ID); ports cleanly to PromiseGrid. **Harder:** requires the bot to be careful about not pushing back to the same branch after merge (avoids divergence); requires a polling cadence for status updates.

- **Alt-β (GitHub-PR centric):** Bot opens a PR; review happens in PR comments; merge happens via the PR UI. **Easier:** uses GitHub's built-in review UI. **Harder:** decision-and-review record lives outside the repo (PR descriptions, comments, approvals are not version-controlled); GitHub-specific; no clean PromiseGrid migration path.

- **Alt-γ (chat-only review):** Bot pushes `ppx/{twig}`; Steve reviews in chat with the bot; merge happens via direct push. No DR file. **Easier:** lightest ceremony. **Harder:** the review reasoning is ephemeral (lives in chat history, not in the repo); audit trails require chat-log archaeology.

- **Alt-δ (review-message-as-promise, used downstream):** Steve writes a review message under `proposals/pending/<proposal-id>/review-<timestamp>-steve-traugott.md` whose load-bearing sentence is a conditional promise. This is *additive* to Alt-α, not a replacement: Alt-α defines how the bot publishes; Alt-δ defines how Steve replies durably. (Locked separately as DI-003-20260429-162212 / DR-005.)

## Selection

Alt-α as the base review style; Alt-δ as the durable-review-reply layer on top of it. The locked decision for Alt-α is `DI-001-20260428-195702` in `protocols/wire-lab.d/TODO/TODO-20260429-030146-perplexity-computer-onboarding.md`, supporting DR `DR/DR-003-20260428-195700-review-style.md`. The Alt-δ layer is `DI-003-20260429-162212` / `DR-005`. Together they form the bot review ceremony in this repo.

## Decision status

`decided` — DI-001-20260428-195702 is already active. This TE is retrospective and exists to satisfy the TE-protocol requirement that non-trivial decisions have a corresponding TE artifact.
