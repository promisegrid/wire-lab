# TODO 009 - Repo-local channel carrier via `grid <pcid>`

Capture and lock the repo-local `channels/` carrier choice so Codex and Perplexity can keep iterating on the channel work without prematurely freezing the canonical PromiseGrid wire format.

## Subtasks

- [x] 009.1 Run a thought experiment comparing `grid <pcid>`, ordinary header-based protocol selection, filename/path selection, and a pure structured object.
- [x] 009.2 Lock the repo-local `channels/` carrier decision in a Decision Intent record.
- [x] 009.3 Write a long-form report and v0 contract in `channels/README.md`.
- [x] 009.4 Add the TE to `docs/thought-experiments/README.md`.
- [x] 009.5 Update `harness-spec.md` with a TE summary and a canonical open question about whether the repo-local carrier should graduate into the canonical wire format.
- [ ] 009.6 Create per-direction channel subdirectories only after the v0 contract has been reviewed and future message traffic justifies the extra structure.

## Decision Intent Log

ID: DI-009-20260430-204108
Date: 2026-04-30 20:41:08
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: For the repo-local `channels/` experiment only, a channel message is represented as canonical text whose first line is `grid <pcid>`. The `pCID` selects the channel-message protocol. The message CID is computed as CIDv1 using `base32`, `sha2-256`, and `raw` over the full canonical message-file bytes. `Prev-Message-CID` links each message to its predecessor in the same single-writer channel. Receipts use `IHave: <channel>:<cid>` as a compact channel-qualified acknowledgement. `Message-ID` and `Kind` remain human-oriented convenience fields, while the load-bearing semantic content stays in an explicit `I promise ...` body. This decision does not lock the long-term canonical PromiseGrid wire format.
Intent: Use the repo itself as a low-friction testbed for pCID-selected messaging, CID-chained receipts, and human/LLM-readable coordination so the wire-format discussion can be exercised with real artifacts before the canonical PromiseGrid message format is frozen.
Constraints: The decision applies only to repo-local `channels/` artifacts. Do not create subdirectories under `channels/` in this change. Do not merge or respond to `origin/ppx/main` as part of this work. Keep the distinction clear between protocol CIDs and message/payload CIDs.
Affects: `TODO/TODO.md`; `TODO/009-grid-pcid-channel-carrier.md`; `DR/DR-009-20260430-204108-grid-pcid-channel-carrier.md`; `docs/thought-experiments/TE-20260430-204108-grid-pcid-channel-carrier.md`; `docs/thought-experiments/README.md`; `harness-spec.md`; `channels/README.md`; future repo-local channel messages between Codex and Perplexity.
Linked DR: `DR/DR-009-20260430-204108-grid-pcid-channel-carrier.md`

## Notes

- This TODO intentionally scopes the decision to the repo-local `channels/` experiment. Whether the same carrier shape should survive into the canonical PromiseGrid wire format remains open in `DR/DR-009-20260430-204108-grid-pcid-channel-carrier.md`.
- The long-form rationale lives in `channels/README.md` and the TE document; this TODO is the lock record.
