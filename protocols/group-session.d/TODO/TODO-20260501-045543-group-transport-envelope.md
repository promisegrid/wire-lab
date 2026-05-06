# TODO 012 - Group-transport envelope via `grid <pcid>`

Capture and lock the wire-lab's first transport-protocol envelope so Codex and Perplexity can keep iterating on the group-transport work without prematurely freezing the canonical PromiseGrid wire format.

## Note on rename

Originally tracked as `TODO/012-grid-pcid-channel-carrier.md` under the wire-lab's earlier "channel" vocabulary. Per TE-27 (transports rename) and the TODO 013 carve-out, renamed in place to `protocols/group-session.d/TODO/TODO-20260501-045543-group-transport-envelope.md`. The TODO integer (012), DI-ID (`DI-009-20260430-204108`), and original intent are preserved.

## Subtasks

- [x] 012.1 Run a thought experiment comparing `grid <pcid>`, ordinary header-based protocol selection, filename/path selection, and a pure structured object.
- [x] 012.2 Lock the wire-lab's first transport-protocol envelope decision in a Decision Intent record.
- [x] 012.3 Write a long-form report and v0 contract.
- [x] 012.4 Add the TE to `docs/thought-experiments/README.md`.
- [x] 012.5 Update `specs/harness-spec-draft.md` with a TE summary and a canonical open question about whether the wire-lab's transport envelope should graduate into the canonical wire format.
- [x] 012.6 Carve the v0 contract out of `transports/README.md` into the substantive group-transport spec so the README points at the spec rather than embedding the contract. (Performed during TODO 013. Note: the spec was named `specs/group-transport-draft.md` at the time of carve-out; per TE-29 + TE-32 the substantive contract now lives at `protocols/group-session.d/specs/group-session-draft.md`. Per DF-38.5 locked verbally in turn 176 and written 2026-05-06, the draft state-suffix convention is `<slug>-<state>` so the file remains `group-session-draft.md`. The carve-out itself is the locked work; the file path drift is a pure rename.)
- [x] 012.7 First real round-trip on the wire-lab-devs instance: a four-message DAG with two distinct senders, exercising §3 (CID computation, all four CIDs verified by `tools/spec cid`), §4 (envelope, all four messages), §4.6 (`Parents:` DAG link, three uses), §6 (body-as-receipt, three uses), and §7 (append-only, four files coexisting unmodified) of `protocols/group-session.d/specs/group-session-draft.md`. Closed 2026-05-06 via Alt-F2 (DF-021-TODO12.2): m003 from `alice` (TE-character mock second sender) cites m002 in `Parents:` and acknowledges in body. Artifacts at `transports/wire-lab-devs-draft/`. The §9 per-author-branch multi-LLM convergence story is explicitly non-normative (not a freeze-gate requirement) and remains an open §9 OQ.
- [ ] 012.8 Freeze `protocols/wire-lab.d/specs/transport-spec-draft.md` (outer rule) and `protocols/group-session.d/specs/group-session-draft.md` (substantive contract). 012.7 has produced sufficient on-disk artifacts to validate the v0 envelope. Remaining freeze-gate conditions are wider than this TODO: outer transport-spec-draft.md must freeze first (its own gate), and Steve must sign a `merge-group-transport-spec` promise. Tracked under T-GROUP-SESSION-FREEZE in OPEN-THREADS.

## Decision Intent Log

ID: DI-009-20260430-204108
Date: 2026-04-30 20:41:08
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: For the wire-lab's first transport-protocol — the group-transport-protocol class defined in `specs/group-transport-draft.md` — a transport message is represented as canonical text whose first line is `grid <pcid>`. The pCID selects the protocol class. The message CID is computed as CIDv1 using `base32`, `sha2-256`, and `raw` over the full canonical message-file bytes. DAG links to prior messages are expressed via a `Parents:` header (single line, space-separated message CIDs, always optional). Receipts are expressed in message bodies, not in envelope headers. `Message-ID` remains a human-oriented convenience field; the load-bearing semantic content stays in an explicit `I promise ...` body. This decision does not lock the long-term canonical PromiseGrid wire format. Earlier sketches that included `Prev-Message-CID:`, `IHave:`, and `Kind:` headers are superseded by the TODO 013 carve-out: `Parents:` replaces single-parent chaining, body-level acknowledgement replaces `IHave:`, and `Kind:` is dropped entirely.
Intent: Use the wire-lab itself as a low-friction testbed for pCID-selected messaging, message-CID-linked DAG references, and human/LLM-readable coordination so the wire-format discussion can be exercised with real artifacts before the canonical PromiseGrid message format is frozen.
Constraints: The decision applies only to the group-transport-protocol class. Other transport-protocol classes (ring, gossip, hub-mediated, large-N, ephemeral, etc.) will produce their own envelope decisions in their own spec docs. Keep the distinction clear between protocol CIDs and message CIDs.
Affects: `TODO/TODO.md`; `protocols/group-session.d/TODO/TODO-20260501-045543-group-transport-envelope.md`; `DR/DR-009-20260430-204108-group-transport-envelope.md`; `docs/thought-experiments/TE-20260430-204108-group-transport-envelope.md`; `docs/thought-experiments/README.md`; `specs/harness-spec-draft.md`; `specs/transport-spec-draft.md`; `specs/group-transport-draft.md`; `transports/README.md`; future group-transport messages between Codex and Perplexity.
Linked DR: `DR/DR-009-20260430-204108-group-transport-envelope.md`

## Notes

- This TODO scopes the decision to the group-transport-protocol class. Whether the same envelope shape should survive into the canonical PromiseGrid wire format remains open in `DR/DR-009-20260430-204108-group-transport-envelope.md`.
- The long-form rationale lives in TE-24 and the v0 contract lives in `specs/group-transport-draft.md`; this TODO is the lock record.
- The TODO 013 carve-out (separate file `protocols/wire-lab.d/TODO/TODO-20260501-045544-transports-carveout.md`) split the original combined "channels/" material into the outer transport-spec, the group-transport spec, and this updated TODO. See TODO 013 for the carve-out details.
