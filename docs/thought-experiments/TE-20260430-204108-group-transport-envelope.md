# TE-24: Group-transport envelope: `grid <pcid>` carrier, canonical bytes, and explicit promise body

*Thought experiment, part of the [PromiseGrid Wire Lab](../../specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-20260430-204108

## Note on rename

This TE was originally drafted under "channel" vocabulary (filename `TE-20260430-204108-grid-pcid-channel-carrier.md`, title "`grid <pcid>` as a repo-local channel-message carrier"). Per [TE-27](TE-20260501-021921-transports-rename-and-axes-of-differentiation.md), the wire-lab vocabulary is **transports** and **messages**; "channel" is not used. Per the TODO 013 carve-out, TE-24's load-bearing decisions are properties of the group-transport-protocol, not of the outer wire-lab transport-spec. The file has been renamed and the prose rewritten in place. The TE integer (24) and timestamp (20260430-204108) are preserved.

## Decision under test

For the wire-lab's first transport-protocol — small-finite-closed-group with N≥2 participants, all-see-all visibility, multi-writer DAG of messages — what envelope shape should a message carry so it can:

- declare which protocol governs its parsing,
- support CID-linked references between messages,
- stay legible in git diffs and chat excerpts to humans and LLMs, and
- leave room for protocol evolution without locking the canonical PromiseGrid wire format prematurely?

Specifically:

- Should protocol selection happen in an ordinary wrapper header or in a distinct first-line carrier?
- How much of message identity should come from the message CID versus human-readable convenience fields?
- How far should the wire-lab's first transport-protocol go toward the eventual canonical PromiseGrid wire format, versus stopping at a deliberately temporary envelope?

## Assumptions

- A `pCID` names a protocol spec, not a payload instance or a message instance.
- The canonical PromiseGrid model is still "message as promise stack", not a permanently fixed textual envelope.
- The wire-lab's first transport-protocol exists to exercise protocol-selection-by-pCID, message-CID-based references, and receipt semantics in a human/LLM-readable form before the canonical wire format is frozen.
- Each transport instance under this protocol is a single directory whose interior is a multi-writer DAG of message files; ordering across writers is carried by the DAG of parent links, not by directory layout.
- The body of a transport message must contain an explicit promise so the message remains legible as promise-theory discourse rather than devolving into an envelope whose semantics are hidden in field names alone.

## Alternatives

### Alt-A: Ordinary header field

Use a normal header such as `Protocol-CID: <pcid>` alongside the other headers. The file begins with the usual header block and has no special first-line carrier.

### Alt-B: First-line carrier selector

Begin every transport message with a line of the form `grid <pcid>`, then follow it with a blank line, an ordered header block, another blank line, and the explicit `I promise ...` body.

### Alt-C: Filename/path-selected protocol

Encode the protocol choice in the filename or the transport directory name instead of in the message content.

### Alt-D: Pure structured object

Skip the textual envelope entirely and jump straight to a fully structured promise-stack object (for example CBOR/IPLD) even for the wire-lab's first transport-protocol.

### Alt-E: No explicit pCID in transport messages

Treat transport messages as ad hoc markdown notes whose semantics are inferred from context, with no explicit protocol selector.

## Scenario analysis

### S1 — Normal operation: a sender writes a review-oriented message into a transport

The sender wants a message format that is easy to type, easy to inspect in git diffs, and explicit enough that a later machine parser can recover the intended protocol.

- **Alt-A** keeps the message readable, but the reader must already know the outer header syntax before it can even discover which protocol governs the rest of the file. That weakens the "pCID as protocol selector" test.
- **Alt-B** lets the parser dispatch immediately on the first line, before interpreting the rest of the envelope. This most directly exercises the idea that a protocol selector should be available up front.
- **Alt-C** makes directory layout do semantic work that should belong to the message. That is fragile once messages are copied, quoted, or moved. It also conflicts with the outer wire-lab transport-spec's principle that a message does not declare its transport — a path-encoded pCID in the directory IS the transport identifier; an additional path-encoded *protocol* selector inside is redundant.
- **Alt-D** is closer to an eventual production shape, but it is a poor fit for the current wire-lab need: humans and LLMs need something they can read and edit directly during design work.
- **Alt-E** avoids design effort now, but it gives up the main value of the exercise: testing pCID-selected protocol dispatch in real coordination traffic.

### S2 — Receipt and loss detection: receiver needs to acknowledge prior messages

The receiver wants to say "I have observed these messages" and tie that statement to specific message identities.

- **Alt-A** can support body-citation receipts equally well, but the protocol-selection signal is still buried inside the header block.
- **Alt-B** supports the same receipt mechanics while preserving the "dispatch on first line" property. Combined with a `Parents:` header citing prior message CIDs, it gives a clean DAG of acknowledgements.
- **Alt-C** is weakest here because the path is not part of the hashed message content. The receiver cannot rely on a copied or renamed file preserving its semantics.
- **Alt-D** supports strong message-CID linking in principle, but at the cost of far more machinery than the wire-lab's first transport-protocol needs today.
- **Alt-E** leaves receipt semantics underspecified and makes loss detection dependent on out-of-band conventions.

### S3 — Human and LLM scanability under long-horizon design work

The wire-lab's first transport-protocol is meant to be read in diffs, quotes, commit history, and chat excerpts.

- **Alt-A** is acceptable, but the first line does not advertise that this is a protocol-selected object.
- **Alt-B** is best: the `grid <pcid>` line is compact, visually distinctive, and easy for both people and tools to detect.
- **Alt-C** is awkward because the most important semantic cue disappears when the message is pasted outside its original path.
- **Alt-D** is poor for this stage because it forces every human and LLM reader through a decoding step.
- **Alt-E** is comfortable in the short term but loses protocol discipline exactly where the experiment is meant to build it.

### S4 — Bridging and protocol evolution

Over time, one transport-protocol pCID may be superseded by another (e.g., the group-transport-protocol gains cryptographic signing, or splits into bounded-retention and append-only variants).

- **Alt-A** supports evolution, but bridging logic must understand the envelope before it can discover the protocol selector.
- **Alt-B** gives the cleanest evolutionary story: the first line names the protocol, so a bridge can immediately decide whether translation is needed.
- **Alt-C** pushes too much semantic weight into path conventions and makes cross-repo or copied-message bridging weaker.
- **Alt-D** may be the long-term answer, but choosing it now would conflate "future canonical shape" with "current wire-lab experiment."
- **Alt-E** makes bridging informal and therefore harder to audit.

### S5 — Relationship to the canonical PromiseGrid wire format

The wire-lab's first transport-protocol should inform the canonical format without pretending to settle it prematurely.

- **Alt-A** is modest, but it does not strongly exercise protocol-selection-on-sight.
- **Alt-B** is strong enough to test pCID selection, message-CID linking, and receipt mechanics, while still being explicit that this is a transport-instance envelope rather than the final canonical message object.
- **Alt-C** tests almost nothing about the wire-format questions that matter.
- **Alt-D** would over-lock the canonical direction too early.
- **Alt-E** keeps the canonical question open, but only by failing to learn much in the meantime.

## Conclusions

1. **Alt-B survives best for the wire-lab's first transport-protocol.** The `grid <pcid>` first line gives immediate protocol dispatch, stays readable in git and chat, and exercises the core pCID idea directly.
2. **The message body still needs explicit promise prose.** Otherwise the experiment regresses toward a fixed envelope whose real meaning is hidden in field names.
3. **The message CID should be authoritative for cross-message references, while `Message-ID` stays as a human/threading convenience field.** The wire-lab's first transport-protocol benefits from both machine-stable identity and human-friendly references.
4. **Canonical text serialization should be locked for the wire-lab's first transport-protocol.** Without fixed bytes, message-CID-linked references and acknowledgements are too mushy to test meaningfully.
5. **This TE does not settle the canonical PromiseGrid wire format.** It only recommends an envelope for the wire-lab's first transport-protocol that is strong enough to teach us something before the canonical format is frozen.

## Implications for the repo's open TODOs and pending DIs

- These conclusions are properties of the **group-transport-protocol**, not of the outer wire-lab transport-spec. The outer transport-spec ([`specs/transport-spec-draft.md`](../../specs/transport-spec-draft.md)) is silent on envelope shape, header sets, canonical bytes, and body conventions. The substantive contract — including the locked decisions from the TODO 013 carve-out (`Parents:` header replacing `Prev-Message-CID:`/`IHave:`, no `Kind:` header, ack-in-body, flat subdirectory layout) — lives in [`specs/group-transport-draft.md`](../../specs/group-transport-draft.md).
- `protocols/group-session.d/TODO/TODO-20260501-045543-group-transport-envelope.md` locks this envelope choice for the group-transport-protocol separately from the broader promise-stack DIs tracked in `protocols/wire-lab.d/TODO/TODO-20260429-164955-te-promise-stack-ordering.md`.
- `transports/README.md` points readers at both the outer transport-spec and the group-transport spec; this TE is the source document for the latter.
- `specs/harness-spec-draft.md` acknowledges the TE in its §8 bibliography. The graduation question — does this envelope shape collapse into a more structured canonical wire object eventually? — remains open in `DR-009-20260430-204108`.

## Decision status

`locked for the group-transport-protocol` — the envelope decision for the wire-lab's first transport-protocol is recorded in `DI-009-20260430-204108`. Whether this shape should graduate into the canonical PromiseGrid wire format remains open.
