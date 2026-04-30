# DR-009 - Repo-local `channels/` carrier via `grid <pcid>`

DR-ID: DR-009-20260430-204108
Date: 2026-04-30 20:41:08
Asked by: stevegt@t7a.org (Steve Traugott)
State: decided
Question: For the repo-local `channels/` experiment, what carrier shape should select the channel-message protocol and support CID-chained receipts while remaining easy for humans and LLMs to read, edit, and reason about?
Why this blocks progress: The repo has converged on the idea of an explicit `channels/` area for Codex↔Perplexity coordination, but without a carrier choice the work cannot advance beyond chat fragments. The project needs a concrete, auditable way to test protocol-selection-by-pCID, CID-linked message chains, and receipt semantics in-repo without pretending that the canonical PromiseGrid wire format is already frozen.
Affects: `channels/README.md`; future `channels/*` artifacts; `TODO/009-grid-pcid-channel-carrier.md`; `docs/thought-experiments/TE-20260430-204108-grid-pcid-channel-carrier.md`; `harness-spec.md` references to the thought experiment and its remaining open question.
Unblocks: TODO 009 subtasks; future channel-message artifacts; future Codex↔Perplexity message exchange experiments; later decision work on whether the repo-local carrier should graduate into the canonical wire format.
Waiting on: DI-009-20260430-204108

## Candidate alternatives considered

The full scenario analysis lives in `docs/thought-experiments/TE-20260430-204108-grid-pcid-channel-carrier.md`. The load-bearing alternatives were:

- An ordinary wrapper header such as `Protocol-CID: <pcid>`.
- A first-line textual carrier selector of the form `grid <pcid>`.
- Encoding the protocol choice in the filename or path rather than in the message body.
- Skipping the textual wrapper entirely and jumping straight to a fully structured promise-stack object.

## Decision

For the repo-local `channels/` experiment only, choose the `grid <pcid>` first-line carrier. Canonicalize the message as UTF-8 text with LF line endings, a fixed header order, exactly one blank line before the body, and a trailing newline at EOF. Compute the message CID as CIDv1 (`base32`, `sha2-256`, `raw`) over the full canonical file bytes. Keep `Message-ID` and `Kind` as human-oriented convenience fields, but require the body to carry explicit promise prose. Use `Prev-Message-CID` for the single-writer message chain and `IHave: <channel>:<cid>` for compact channel-qualified receipts.

This decision is explicitly scoped to repo-local `channels/` artifacts and does not settle the long-term canonical PromiseGrid wire format.

## Linked DI

- `DI-009-20260430-204108` in `TODO/009-grid-pcid-channel-carrier.md`

## Related commits

- Pending first commit on `stevegt/channels-grid-pcid`

## Last updated

2026-04-30 20:41:08 UTC
