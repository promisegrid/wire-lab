# Channels

This directory is the repo-local coordination layer for channel-style message traffic between Codex and Perplexity. It is intentionally narrower than the main `proposals/` workflow. `proposals/` remains the governance and review queue; `channels/` is where the repo can exercise message-carrier ideas directly, with committed artifacts that are readable in git and understandable to humans and LLMs. Source: `DI-009-20260430-204108`.

## Status

The current choice is intentionally scoped and intentionally provisional:

- For the repo-local `channels/` experiment only, the carrier starts with `grid <pcid>`. Source: `DI-009-20260430-204108`.
- The carrier is a testbed for protocol selection, CID chaining, and receipts. It is **not** yet the canonical PromiseGrid wire format. Source: `DI-009-20260430-204108`.
- This change does **not** create any subdirectories under `channels/` yet. The directory is being used first as a design/reporting surface so the protocol can be reviewed before the repo starts emitting real traffic here. Source: `DI-009-20260430-204108`.

The thought experiment behind this choice lives in `docs/thought-experiments/TE-20260430-204108-grid-pcid-channel-carrier.md`. The decision request is `DR/DR-009-20260430-204108-grid-pcid-channel-carrier.md`. The lock record is in `TODO/012-grid-pcid-channel-carrier.md`.

## What problem `channels/` is solving

The repo already has one established queue: `proposals/pending/` and `proposals/approved/`. That queue is good for governance and review artifacts, but it is not a great place to experiment with message-carrier details such as:

- how a pCID selects a protocol in practice,
- how a message acknowledges prior messages without mutable "read" ledgers,
- how a human/LLM-readable message can still be content-addressed,
- how a repo-local convention can teach the wire-format discussion something real before the canonical format is frozen.

`channels/` exists to test those questions with low ceremony. Source: `DI-009-20260430-204108`.

## Current v0 carrier

The current repo-local carrier is:

```text
grid <pcid>

Message-ID: MSG-20260430-103000Z-receipt-review
Date: 2026-04-30 10:30:00 UTC
From: stevegt+ppx@t7a.org (stevegt-via-perplexity)
To: codex@local (Codex)
Prev-Message-CID: <previous-message-cid>
Kind: receipt
IHave: codex-perplexity:<cid-of-latest-codex-message-seen>

I promise I have observed the `codex-perplexity` message chain through
CID `<cid-of-latest-codex-message-seen>`.
```

### Why the first line is `grid <pcid>`

The first line is the protocol selector. This is the strongest part of the v0 decision:

- the parser can dispatch immediately on the first line,
- the pCID is visible before any header interpretation happens,
- the message stays readable in diffs and chat excerpts,
- the carrier exercises the "pCID selects a protocol spec" idea directly rather than only talking about it.

This is the main reason the current choice prefers `grid <pcid>` over an ordinary wrapper header like `Protocol-CID: ...`. Source: `DI-009-20260430-204108`.

### Why the body still says `I promise ...`

The message is not meant to become a hollow envelope whose real semantics are hidden inside field names. The body still carries the load-bearing promise in plain prose. Headers such as `Kind:` and `IHave:` are conveniences for indexing and parsing; they are not meant to erase the underlying promise-theory framing. Source: `DI-009-20260430-204108`.

## CID and canonicalization rules

The repo-local channel experiment needs strong enough byte discipline that a CID means something real. The current v0 choice is:

- encode the message as UTF-8 text,
- use LF line endings only,
- keep a fixed header order,
- use exactly one blank line between the header block and the body,
- end the file with a trailing newline,
- omit absent optional headers entirely,
- compute the message CID as CIDv1 using `base32`, `sha2-256`, and `raw` over the full canonical file bytes.

That means the file itself is the hashed artifact. The CID is not based only on the body, not based on the filename, and not based on git metadata. Source: `DI-009-20260430-204108`.

This choice is strong enough to test real chaining and receipts without prematurely forcing the whole repo onto a binary or schema-first representation. It also lines up with the existing `tools/spec` direction: content-address the artifact itself, not the path around it. Source: `DI-009-20260430-204108`.

## Message identity, chaining, and receipts

The carrier distinguishes between human-friendly references and the machine-stable chain:

- `Message-ID` is kept as a human/threading convenience field.
- The message CID is the authoritative identity for chaining and receipts.
- `Prev-Message-CID` links each message to the previous message in the same single-writer channel.
- `IHave: <channel>:<cid>` is the compact cumulative receipt form.

The important semantic point is that `IHave` is **not** a per-message "read" mark. It means "I have the contiguous chain through this CID for the named channel." That avoids mutable read ledgers while still making loss detection and acknowledgement explicit. Source: `DI-009-20260430-204108`.

## Why `IHave` instead of a read ledger

The repo-local experiment considered a mutable list of "read messages." That path was rejected because it creates shared mutable state, duplicates the message log, and does not fit an append-only git-backed workflow well.

`IHave` is better because:

- it is append-only,
- it lives in the reverse direction as another message,
- it composes naturally with `Prev-Message-CID`,
- it lets the receiver acknowledge a contiguous prefix rather than each message one by one.

The current compact form is `IHave: codex-perplexity:<cid>`, not separate `IHave-Channel:` and `IHave:` fields. That keeps the receipt small while still naming the chain being acknowledged. Source: `DI-009-20260430-204108`.

## Why this is only a repo-local carrier

The project is not ready to declare that the canonical PromiseGrid wire format is a text file that starts with `grid <pcid>`. That would be over-locking.

The repo-local carrier is deliberately narrower:

- it is a human/LLM-friendly experiment,
- it exercises pCID-selected protocol dispatch in real coordination work,
- it teaches the team what parts of the shape feel stable and what parts feel like scaffolding,
- it leaves open the possibility that the canonical format later becomes a more structured promise-stack object.

That distinction matters because the harness-spec still treats the long-term message as a promise stack and explicitly resists freezing parser-level envelopes too early. Source: `DI-009-20260430-204108`. The open graduation question remains tracked in `DR/DR-009-20260430-204108-grid-pcid-channel-carrier.md`.

## Alternatives that were considered

### 1. Ordinary wrapper header: `Protocol-CID: <pcid>`

**Pros**

- easy to explain,
- compatible with the current markdown-header idiom used in proposal/review artifacts,
- keeps the whole message in one conventional header block.

**Cons**

- the parser must already know the wrapper syntax before it can discover which protocol governs the message,
- it weakens the test of "pCID as first-class protocol selector",
- it feels more like metadata on a note than like a selected protocol object.

### 2. Filename/path-selected protocol

**Pros**

- little in-file overhead,
- directory names can still communicate direction or ownership.

**Cons**

- the protocol signal disappears when the message is copied or quoted,
- the path is not part of the hashed content,
- directory layout ends up doing semantic work that belongs to the message itself.

This was rejected as too weak for a real protocol experiment. Source: `DI-009-20260430-204108`.

### 3. Pure structured object now

Examples include CBOR, IPLD, or a more direct serialization of the promise stack.

**Pros**

- closer to a likely eventual production representation,
- easier to make mechanically unambiguous,
- easier to hand to pure programs later.

**Cons**

- too much structure too early,
- weaker fit for human/LLM editing during active design work,
- risks conflating "temporary repo-local carrier" with "canonical PromiseGrid wire format."

This remains a serious long-term alternative. It was not chosen for the repo-local v0 because the experiment still benefits from being plainly readable and editable. Source: `DI-009-20260430-204108`.

### 4. No explicit pCID in `channels/`

**Pros**

- least ceremony,
- fastest way to start exchanging notes.

**Cons**

- fails to test the pCID idea,
- leaves protocol selection implicit,
- gives the project almost no empirical signal about how a pCID-selected carrier actually feels.

This option was rejected because it gives up the main value of introducing `channels/` at all. Source: `DI-009-20260430-204108`.

## Related thought experiments and why they matter here

The `channels/` work is not isolated. It touches several earlier TEs:

- `TE-1` (`docs/thought-experiments/TE-20260427-180000-promise-stack-ordering.md`) matters because it established that protocol semantics should not be confused with one global fixed frame placement. The channel carrier is therefore explicitly scoped to repo-local transport, not elevated to the canonical message object.
- `TE-9` (`docs/thought-experiments/TE-20260427-180800-two-communities-two-pcids-same-intent.md`) matters because `channels/` is a direct place to test how protocol selection by pCID behaves when more than one message-carrier idea may exist over time.
- `TE-12` (`docs/thought-experiments/TE-20260427-181100-promise-stack-as-zero-knowledge-envelope.md`) matters because it reinforces that a promise frame can carry richer semantics later; the repo-local textual carrier should not foreclose that path.
- `TE-14` (`docs/thought-experiments/TE-20260428-080000-harness-spec-change-walks-through-unified-flow.md`) matters because it pushed the project away from over-locked envelopes and toward suggested shapes that can earn structure over time.
- `TE-15` (`docs/thought-experiments/TE-20260428-094500-should-this-design-become-promisegrid-readme.md`) matters because it sharpened the distinction between a repo-local experimental artifact and a canonical public-facing protocol statement.
- `TE-24` (`docs/thought-experiments/TE-20260430-204108-grid-pcid-channel-carrier.md`) is the direct comparison of the channel-carrier alternatives and is the immediate rationale for the current v0 decision.

That TE bundle is explanatory context for the repo-local v0 choice locked in `DI-009-20260430-204108`.

## What this document is, and what it is not

This document is:

- the long-form report you can hand to Perplexity so the channel work can continue from a committed in-repo artifact,
- the v0 contract for repo-local channel messages,
- the bridge between the earlier chat discussion and the repo's DR/DI/TE protocol.

This document is **not**:

- a declaration that the canonical PromiseGrid wire format is textual,
- a declaration that `channels/` is now fully operational with per-direction logs,
- a commitment that the eventual production wire format will keep `Message-ID`, `Kind`, or the exact textual header names found here.

Source: `DI-009-20260430-204108`.

## Open questions

These are intentionally left open and should be treated as DR-backed uncertainties, not settled facts:

- Should the repo-local `grid <pcid>` carrier survive into the canonical PromiseGrid wire format, or eventually collapse into a more structured promise-stack object? Open question: `DR/DR-009-20260430-204108-grid-pcid-channel-carrier.md`.
- When `channels/` begins carrying real traffic, should the repo create per-direction subdirectories such as `codex-perplexity/` and `perplexity-codex/`, or is a flatter initial layout still better? Open question: `DR/DR-009-20260430-204108-grid-pcid-channel-carrier.md`.
- Should the future tool support reuse the existing CID primitive from `tools/spec` directly, or split it into a smaller shared helper before channel tooling begins? Open question: `DR/DR-009-20260430-204108-grid-pcid-channel-carrier.md`.

## Next steps

The next useful work after this report is:

1. Let Perplexity continue from the pushed branch rather than trying to merge or answer `origin/ppx/main`.
2. Decide whether the next change should be real channel-message examples under `channels/`, or tooling that computes/validates the message CIDs first.
3. Once real traffic starts, revisit whether `channels/` should split into per-direction append-only logs.
