# TE-24: `grid <pcid>` as a repo-local channel-message carrier

*Thought experiment, part of the [PromiseGrid Wire Lab](../../specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-20260430-204108

## Decision under test

For the repo-local `channels/` experiment, how should a channel message identify the protocol it follows and carry enough structure to support CID-linked receipts, human/LLM readability, and future protocol evolution?

Specifically:

- Should protocol selection happen in an ordinary wrapper header or in a distinct first-line carrier?
- How much of the message identity should come from the message CID versus human-readable convenience fields?
- How far should the repo-local channel experiment go toward the eventual canonical PromiseGrid wire format, versus stopping at a deliberately temporary carrier?

## Assumptions

- A `pCID` names a protocol spec, not a payload instance or a message instance.
- The canonical PromiseGrid model is still "message as promise stack", not a permanently fixed textual envelope.
- The repo-local `channels/` experiment exists to exercise protocol selection, hash chaining, and receipts in a human/LLM-readable form before the canonical wire format is frozen.
- Each direction of message flow will eventually be a single-writer append-only log, even if the initial change stops at `channels/README.md` and does not yet create per-direction subdirectories.
- The body of a channel message must still contain an explicit promise so the message remains legible as promise-theory discourse rather than devolving into an envelope whose semantics are hidden in field names alone.

## Alternatives

### Alt-A: Ordinary header field

Use a normal header such as `Protocol-CID: <pcid>` alongside the other headers. The file begins with the usual header block and has no special first-line carrier.

### Alt-B: First-line carrier selector

Begin every repo-local channel message with a line of the form `grid <pcid>`, then follow it with a blank line, an ordered header block, another blank line, and the explicit `I promise ...` body.

### Alt-C: Filename/path-selected protocol

Encode the protocol choice in the filename or the channel directory name instead of in the message content.

### Alt-D: Pure structured object

Skip the textual carrier entirely and jump straight to a fully structured promise-stack object (for example CBOR/IPLD) even for the repo-local `channels/` experiment.

### Alt-E: No explicit pCID in repo-local channels

Treat `channels/` as ad hoc markdown notes whose semantics are inferred from context, with no explicit protocol selector.

## Scenario analysis

### S1 — Normal operation: Codex writes a review-oriented message to Perplexity

The sender wants a message format that is easy to type, easy to inspect in git diffs, and explicit enough that a later machine parser can recover the intended protocol.

- **Alt-A** keeps the message readable, but the reader must already know the outer header syntax before it can even discover which protocol governs the rest of the file. That weakens the "pCID as protocol selector" test.
- **Alt-B** lets the parser dispatch immediately on the first line, before interpreting the rest of the wrapper. This most directly exercises the idea that a protocol selector should be available up front.
- **Alt-C** makes directory layout do semantic work that should belong to the message. That is fragile once messages are copied, quoted, or moved.
- **Alt-D** is closer to an eventual production shape, but it is a poor fit for the current repo-local need: humans and LLMs need something they can read and edit directly during design work.
- **Alt-E** avoids design effort now, but it gives up the main value of the exercise: testing pCID-selected protocol dispatch in real coordination traffic.

### S2 — Receipt and loss detection: receiver needs to acknowledge a contiguous prefix

The receiver wants to say "I have the entire message chain through this point" without maintaining a mutable read ledger.

- **Alt-A** can support `Prev-Message-CID` and `IHave`, but the protocol-selection signal is still buried inside the header block.
- **Alt-B** supports the same receipt mechanics while preserving the "dispatch on first line" property. Combined with `Prev-Message-CID`, it gives a clean single-writer hash chain.
- **Alt-C** is weakest here because the path is not part of the hashed message content. The receiver cannot rely on a copied or renamed file preserving its semantics.
- **Alt-D** supports strong chaining in principle, but at the cost of far more machinery than the repo-local channel experiment needs today.
- **Alt-E** leaves receipt semantics underspecified and makes loss detection dependent on out-of-band conventions.

### S3 — Human and LLM scanability under long-horizon design work

The repo-local channel is meant to be read in diffs, quotes, commit history, and chat excerpts.

- **Alt-A** is acceptable, but the first line does not advertise that this is a protocol-selected object.
- **Alt-B** is best: the `grid <pcid>` line is compact, visually distinctive, and easy for both people and tools to detect.
- **Alt-C** is awkward because the most important semantic cue disappears when the message is pasted outside its original path.
- **Alt-D** is poor for this stage because it forces every human and LLM reader through a decoding step.
- **Alt-E** is comfortable in the short term but loses protocol discipline exactly where the experiment is meant to build it.

### S4 — Bridging and protocol evolution

Over time, one channel-message protocol pCID may be superseded by another.

- **Alt-A** supports evolution, but bridging logic must understand the wrapper before it can discover the protocol selector.
- **Alt-B** gives the cleanest evolutionary story: the first line names the protocol, so a bridge can immediately decide whether translation is needed.
- **Alt-C** pushes too much semantic weight into path conventions and makes cross-repo or copied-message bridging weaker.
- **Alt-D** may be the long-term answer, but choosing it now would conflate "future canonical shape" with "current repo-local experiment."
- **Alt-E** makes bridging informal and therefore harder to audit.

### S5 — Relationship to the canonical PromiseGrid wire format

The repo-local experiment should inform the canonical format without pretending to settle it prematurely.

- **Alt-A** is modest, but it does not strongly exercise protocol-selection-on-sight.
- **Alt-B** is strong enough to test pCID selection, CID chaining, and receipt mechanics, while still being explicit that this is a repo-local carrier rather than the final canonical message object.
- **Alt-C** tests almost nothing about the wire-format questions that matter.
- **Alt-D** would over-lock the canonical direction too early.
- **Alt-E** keeps the canonical question open, but only by failing to learn much in the meantime.

## Conclusions

1. **Alt-B survives best for the repo-local `channels/` experiment.** The `grid <pcid>` first line gives immediate protocol dispatch, stays readable in git and chat, and exercises the core pCID idea directly.
2. **The message body still needs explicit promise prose.** Otherwise the experiment regresses toward a fixed envelope whose real meaning is hidden in field names.
3. **The message CID should be authoritative for chaining and receipts, while `Message-ID` stays as a human/threading convenience field.** The repo-local experiment benefits from both machine-stable identity and human-friendly references.
4. **Canonical text serialization should be locked for the repo-local experiment.** Without fixed bytes, CID-linked receipts are too mushy to test meaningfully.
5. **This TE does not settle the canonical PromiseGrid wire format.** It only recommends a repo-local carrier that is strong enough to teach us something before the canonical format is frozen.

## Implications for the repo's open TODOs and pending DIs

- `TODO/012-grid-pcid-channel-carrier.md` should lock the repo-local decision separately from the broader promise-stack DIs tracked in `TODO/005-te-promise-stack-ordering.md`.
- `channels/README.md` should carry the long-form rationale and the v0 contract so Perplexity can continue the work from a committed in-repo artifact.
- `specs/harness-spec-draft.md` should acknowledge the TE and keep the graduation question open: does the repo-local `grid <pcid>` carrier remain a local convenience, or does it eventually collapse into a more structured canonical wire object?

## Decision status

`locked for repo-local channels only` — the repo-local carrier decision is recorded in `DI-009-20260430-204108`. Whether this shape should graduate into the canonical PromiseGrid wire format remains open.
