# TODO-bisur - Group-transport envelope via `grid <pcid>`

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TODO-12` (integer alias)
- `TODO-20260501-045543` (timestamp alias and pre-migration filename)

Capture and lock the wire-lab's first transport-protocol envelope so Codex and Perplexity can keep iterating on the group-transport work without prematurely freezing the canonical PromiseGrid wire format.

## Note on rename

Originally tracked as `TODO/012-grid-pcid-channel-carrier.md` under the wire-lab's earlier "channel" vocabulary. Per TE-junil (transports rename) and the TODO-motof carve-out, renamed in place to `protocols/group-session.d/TODO/TODO-bisur-group-transport-envelope.md`. The TODO integer (012), DI-ID (`DI-009-20260430-204108`), and original intent are preserved.

## Simulation note

This TODO now lives with the `group-session` protocol specimen at
`simulations/SIM-rakot-group-session/protocols/group-session.d/` per
`DI-firiv`. Historical path references below (including earlier
`SIM-piloh` locations) remain evidence of the pre-split layout unless a
line explicitly names the simulation-local path.

## Subtasks

- [x] 012.1 Run a thought experiment comparing `grid <pcid>`, ordinary header-based protocol selection, filename/path selection, and a pure structured object.
- [x] 012.2 Lock the wire-lab's first transport-protocol envelope decision in a Decision Intent record.
- [x] 012.3 Write a long-form report and v0 contract.
- [x] 012.4 Add the TE to `docs/thought-experiments/README.md`.
- [x] 012.5 Update `specs/harness-spec-draft.md` with a TE summary and a canonical open question about whether the wire-lab's transport envelope should graduate into the canonical wire format.
- [x] 012.6 Carve the v0 contract out of `transports/README.md` into the substantive group-transport spec so the README points at the spec rather than embedding the contract. (Performed during TODO-motof. Note: the spec was named `specs/group-transport-draft.md` at the time of carve-out; per TE-vipir + TE-liviv the substantive contract now lives at `protocols/group-session.d/specs/group-session-draft.md`. Per DF-38.5 locked verbally in turn 176 and written 2026-05-06, the draft state-suffix convention is `<slug>-<state>` so the file remains `group-session-draft.md`. The carve-out itself is the locked work; the file path drift is a pure rename.)
- [x] 012.7 First real round-trip on the wire-lab-devs instance: a four-message DAG with two distinct senders, exercising §3 (CID computation, all four CIDs verified by `tools/spec cid`), §4 (envelope, all four messages), §4.6 (`Parents:` DAG link, three uses), §6 (body-as-receipt, three uses), and §7 (append-only, four files coexisting unmodified) of `protocols/group-session.d/specs/group-session-draft.md`. Closed 2026-05-06 via Alt-F2 (DF-021-TODO12.2): m003 from `alice` (TE-character mock second sender) cites m002 in `Parents:` and acknowledges in body. Artifacts at `transports/wire-lab-devs-draft/`. The §9 per-author-branch multi-LLM convergence story is explicitly non-normative (not a freeze-gate requirement) and remains an open §9 OQ.
- [ ] 012.8 Freeze `protocols/wire-lab.d/specs/transport-spec-draft.md` (outer rule) and `protocols/group-session.d/specs/group-session-draft.md` (substantive contract). 012.7 has produced sufficient on-disk artifacts to validate the v0 envelope. Remaining freeze-gate conditions are wider than this TODO: outer transport-spec-draft.md must freeze first (its own gate), and Steve must sign a `merge-group-transport-spec` promise. Tracked under T-GROUP-SESSION-FREEZE in OPEN-THREADS.

## Decision Intent Log

ID: DI-009-20260430-204108
Date: 2026-04-30 20:41:08
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: For the wire-lab's first transport-protocol — the group-transport-protocol class defined in `specs/group-transport-draft.md` — a transport message is represented as canonical text whose first line is `grid <pcid>`. The pCID selects the protocol class. The message CID is computed as CIDv1 using `base32`, `sha2-256`, and `raw` over the full canonical message-file bytes. DAG links to prior messages are expressed via a `Parents:` header (single line, space-separated message CIDs, always optional). Receipts are expressed in message bodies, not in envelope headers. `Message-ID` remains a human-oriented convenience field; the load-bearing semantic content stays in an explicit `I promise ...` body. This decision does not lock the long-term canonical PromiseGrid wire format. Earlier sketches that included `Prev-Message-CID:`, `IHave:`, and `Kind:` headers are superseded by the TODO-motof carve-out: `Parents:` replaces single-parent chaining, body-level acknowledgement replaces `IHave:`, and `Kind:` is dropped entirely.
Intent: Use the wire-lab itself as a low-friction testbed for pCID-selected messaging, message-CID-linked DAG references, and human/LLM-readable coordination so the wire-format discussion can be exercised with real artifacts before the canonical PromiseGrid message format is frozen.
Constraints: The decision applies only to the group-transport-protocol class. Other transport-protocol classes (ring, gossip, hub-mediated, large-N, ephemeral, etc.) will produce their own envelope decisions in their own spec docs. Keep the distinction clear between protocol CIDs and message CIDs.
Affects: `TODO/TODO.md`; `protocols/group-session.d/TODO/TODO-bisur-group-transport-envelope.md`; `DR/DR-009-20260430-204108-group-transport-envelope.md`; `docs/thought-experiments/TE-hogus-group-transport-envelope.md`; `docs/thought-experiments/README.md`; `specs/harness-spec-draft.md`; `specs/transport-spec-draft.md`; `specs/group-transport-draft.md`; `transports/README.md`; future group-transport messages between Codex and Perplexity.
Linked DR: `DR/DR-009-20260430-204108-group-transport-envelope.md`

ID: DI-012-20260508-033513
Date: 2026-05-08 03:35:13
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Group-session v0 canonical writers must not emit `Message-ID:`. The message CID remains the only protocol identifier. Readers may accept exactly one legacy `Message-ID:` header when it appears before `Date:`, but they must ignore its value semantically.
Intent: Keep identity content-addressed while preserving enough backward compatibility to read the early `wire-lab-devs` transport messages that were drafted before the `Message-ID:` conflict was settled.
Constraints: The compatibility allowance is reader-only; canonical encoders must omit `Message-ID:`. If a legacy `Message-ID:` header is present, its bytes are still part of the message file and therefore part of CID verification. The header must not be used for `Parents:`, receipts, indexing, deduplication, API identity, or any other load-bearing protocol role. Unknown headers remain rejected except for this explicitly tolerated legacy slot.
Affects: `protocols/group-session.d/specs/group-session-draft.md`; `transports/wire-lab-devs-draft/README.md`; `DR/DR-suhod-mihip-merge-blockers-partial-fix.md`; future group-session readers and canonical encoders.
Supersedes: DI-009-20260430-204108 (`Message-ID` convenience-field clause only)
Linked DR: `DR/DR-suhod-mihip-merge-blockers-partial-fix.md`

## Notes

- This TODO scopes the decision to the group-transport-protocol class. Whether the same envelope shape should survive into the canonical PromiseGrid wire format remains open in `DR/DR-009-20260430-204108-group-transport-envelope.md`.
- The long-form rationale lives in TE-hogus and the v0 contract lives in `specs/group-transport-draft.md`; this TODO is the lock record.
- The TODO-motof carve-out (separate file `protocols/wire-lab.d/TODO/TODO-motof-transports-carveout.md`) split the original combined "channels/" material into the outer transport-spec, the group-transport spec, and this updated TODO. See TODO-motof for the carve-out details.
