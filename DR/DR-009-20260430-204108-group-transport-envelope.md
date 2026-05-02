# DR-009 - Group-transport envelope via `grid <pcid>`

DR-ID: DR-009-20260430-204108
Date: 2026-04-30 20:41:08
Asked by: stevegt@t7a.org (Steve Traugott)
State: decided
Question: For the wire-lab's first transport-protocol — small-finite-closed-group, N>=2 participants, all-see-all visibility, multi-writer DAG of messages — what envelope shape should select the protocol, support message-CID-linked references between messages, and stay easy for humans and LLMs to read, edit, and reason about?
Why this blocks progress: The repo has converged on the idea that `transports/` is the wire-lab's simulation surface for the network being studied, with `transports/<pcid>--<slug>/` as the on-disk keying convention (transport-spec-draft.md). Without an envelope choice for the first transport-protocol, the work cannot advance beyond chat fragments. The project needs a concrete, auditable way to test protocol-selection-by-pCID, message-CID-linked references, and receipt semantics in real coordination traffic without pretending that the canonical PromiseGrid wire format is already frozen.
Affects: `transports/README.md`; future `transports/*` artifacts; `protocols/group-session.d/TODO/TODO-20260501-045543-group-transport-envelope.md`; `docs/thought-experiments/TE-20260430-204108-group-transport-envelope.md`; `specs/transport-spec-draft.md` (outer rule); `protocols/group-session.d/specs/group-session-draft.md` (the substantive v0 contract); `specs/harness-spec-draft.md` references to the thought experiment.
Unblocks: TODO 012 subtasks; TODO 013 carve-out completion; future group-transport message artifacts; future Codex<->Perplexity message exchange experiments; later decision work on whether the wire-lab's transport envelope should graduate into the canonical wire format.
Waiting on: DI-009-20260430-204108

## Note on rename

Originally filed as `DR-009-20260430-204108-grid-pcid-channel-carrier.md` under the wire-lab's earlier "channel" vocabulary. Per TE-27 (transports rename) and the TODO 013 carve-out, renamed in place to `DR-009-20260430-204108-group-transport-envelope.md` because the load-bearing decisions are properties of the group-transport-protocol, not of the outer wire-lab transport-spec. DR-ID, date, and question identity are preserved.

## Candidate alternatives considered

The full scenario analysis lives in `docs/thought-experiments/TE-20260430-204108-group-transport-envelope.md`. The load-bearing alternatives were:

- An ordinary wrapper header such as `Protocol-CID: <pcid>`.
- A first-line textual carrier selector of the form `grid <pcid>`.
- Encoding the protocol choice in the filename or path rather than in the message body.
- Skipping the textual envelope entirely and jumping straight to a fully structured promise-stack object.

## Decision

For the wire-lab's first transport-protocol — the group-transport-protocol class defined in `protocols/group-session.d/specs/group-session-draft.md` — choose the `grid <pcid>` first-line carrier. Canonicalize the message as UTF-8 text with LF line endings, a fixed header order, exactly one blank line between carrier and headers, exactly one blank line between headers and body, and a trailing newline at EOF. Compute the message CID as CIDv1 (`base32`, `sha2-256`, `raw`) over the full canonical file bytes. Keep `Message-ID` as a human-oriented convenience field; require the body to carry explicit promise prose. Use the `Parents:` header (single line, space-separated message CIDs, always optional) to express DAG links between messages. Express acknowledgement in message bodies, not in envelope headers.

This decision is explicitly scoped to the group-transport-protocol class and does not settle the long-term canonical PromiseGrid wire format. The outer wire-lab transport-spec (`specs/transport-spec-draft.md`) remains silent on envelope shape; other transport-protocol classes (ring, gossip, hub-mediated, large-N, ephemeral, etc.) will produce their own envelope decisions in their own spec docs.

The earlier `Prev-Message-CID:` and `IHave:` headers from the original drafting are NOT part of the locked v0 contract; they are superseded by `Parents:` (for DAG links) and body-level acknowledgement (for receipts), per the TODO 013 carve-out. The earlier `Kind:` header is dropped; message kind is presentational and lives in body convention.

## Linked DI

- `DI-009-20260430-204108` in `protocols/group-session.d/TODO/TODO-20260501-045543-group-transport-envelope.md`

## Related commits

- TE-24, TODO 012, DR-009 originally landed on `ppx/main` in the channel-carrier integration (commit `5115f12` and predecessors).
- TODO 013 carve-out renames and rewrites in progress on `ppx/todo-013-transports-carveout`.

## Last updated

2026-04-30 20:41:08 UTC (original); rewritten in place during TODO 013 carve-out, 2026-04-30 21:46 UTC.
