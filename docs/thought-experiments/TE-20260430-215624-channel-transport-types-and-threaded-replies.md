# TE-26: Channel transport types, pCID-keyed channel paths, and DAG message graphs

*Thought experiment, part of the [PromiseGrid Wire Lab](../../specs/harness-spec-draft.md) and the (forthcoming) `specs/channel-spec-draft.md`. This file is content-addressable; its hash is its pCID.*

## TE ID

TE-20260430-215624

## Decision under test

The current `channels/` v0 contract (TE-24, `DR-009-20260430-204108`, `channels/README.md`) implicitly assumes one transport shape: a 1:1 message-based flow simulating something like a websocket between two named endpoints (Codex and Perplexity), with a single-writer per-direction append-only log and a `Prev-Message-CID` linking to the previous message **by the same sender**. Several pieces of that shape need to be re-examined before subdirectory layout is locked in:

1. **Transport-type taxonomy.** What other transport shapes are first-class enough to deserve their own structural support in `channels/`? At minimum: 1:1 unicast, group/broadcast, subgroup multicast, pub/sub by topic, anycast, request-reply, broadcast-with-receipts, single-writer log, and general DAG.
2. **Reply graph shape.** The current `Prev-Message-CID` is a single CID linking only to the previous message **by the same sender** in the same channel. That suffices for a single-writer log, but it cannot represent **threaded discussion** — i.e., a message that has multiple ancestors, possibly from multiple senders. The honest framing is a git-like DAG of messages within a channel, where any message can have zero or more parents; "reply" is just one reason a message might cite a parent.
3. **Layer placement.** Should a message *declare* its channel via a header? No: the channel is a transport carrying the message, and asking a message which channel it is on is layer inversion. Channel identity belongs to the transport (in the repo-local case, the directory path), not to the envelope.
4. **Path keying.** If a channel's directory is the binding between message and channel, what keys that directory? A human-friendly slug carries no protocol identity and forces the wire-lab spec to define structure for every channel. A pCID, by contrast, *is* a protocol identity: each channel is keyed by the pCID of the protocol it speaks, and that protocol gets to define everything inside the directory. The wire-lab spec can then stay small.

Settling these together is appropriate because each one's answer reshapes the others.

## Assumptions

- The repo-local carrier remains `grid <pcid>` first-line + canonical text headers + explicit `I promise ...` body (TE-24, locked).
- The carrier is meant to teach the project something about wire-format and transport choices before the canonical PromiseGrid wire format is frozen. Over-locking transport assumptions in v0 has the same downside as over-locking parser-level envelopes did in TE-14.
- "Single-writer append-only log" is a useful invariant for some transports but not all. Whatever shape replaces it must still let a receiver acknowledge a contiguous prefix without a mutable read ledger (TE-24's argument for `IHave`).
- A message's CID is computed over its full canonical bytes and is the authoritative identity. Anything else (`Message-ID`, headers, file path) is a convenience layer built on top.
- Promise-theory framing: the body still says `I promise ...`. Headers are conveniences for indexing, dispatch, and receipt math; they are not where the load-bearing semantics live.
- The repo's prior `channels/codex-perplexity/` directional layout has not actually been created yet — `channels/README.md` says "this change does not create any subdirectories under `channels/` yet." So the layout question is genuinely open and not blocked by existing committed files under per-direction paths.
- "Let each protocol name its own internals" (standing rule from TE-14) applies recursively: just as the wire-lab carrier names its own internals via its pCID, each channel-protocol names its own internals via its own pCID.

## Transport-type catalogue

The list below is ordered roughly from "narrowest topology" to "broadest reach." Note that under the conclusions below, transport type does **not** drive the on-disk path layout — the per-channel-pCID protocol does. The catalogue is preserved here as raw material that any future channel-protocol-pCID can draw on.

### T1 — 1:1 unicast (the current implicit shape)

One sender, one receiver, two single-writer append-only logs (one per direction). Each direction's log is internally totally ordered by parent links. Cross-direction ordering is partial (only the receipts via `IHave` create cross-direction links).

### T2 — group broadcast (small, named group, all-to-all visibility)

Several senders, all participants see all messages. No one is privileged. The "log" is no longer single-writer; the channel is a directed acyclic graph of messages where each new message names one or more prior messages (the **frontier**) of the group as parents.

### T3 — subgroup multicast (delivery to a named subset)

A sender directs a message to a named subset of a group. Receipts must distinguish "I have this message" from "I am in the named subset." Subset membership itself is a promise that can be made or revoked.

### T4 — pub/sub by topic

Senders publish to a topic; receivers subscribe to topics. Decoupled identities; the topic is the rendezvous point.

### T5 — anycast

A request is delivered to "any one of N addressed peers"; the responder identifies itself in its reply.

### T6 — request/reply (synchronous-feeling RPC)

A logical pair: a request expects exactly one reply (or one reply per addressed responder). Distinct from generic threaded discussion because the reply-or-not state is structurally tracked.

### T7 — broadcast-with-receipts (gossip / epidemic)

A message is published to a channel; all participants are expected to acknowledge with `IHave` over the message's CID.

### T8 — single-writer append-only log (special case of T1)

One writer, many readers. The log is totally ordered, and each message has at most one parent. Useful for status feeds, build artifacts, hash chains of design state.

### T9 — directed acyclic graph (general threaded discussion)

The most general case. Each message has a list of zero or more parent CIDs in the same channel. Reduces to T1 or T8 when each message has exactly one parent in its own log; reduces to T2 when every message names the channel frontier as parents; etc.

## Three principles emerging from the transport survey

1. **A message does not declare its channel.** The channel is the transport that carries it; asking a message which channel it is on is layer inversion. In the repo-local carrier, the channel is identified by the directory the message file lives in. The envelope has no `Channel:` header.

2. **Channel identity is a pCID, not a slug.** The directory under `channels/` is keyed by the pCID of the channel-protocol that channel speaks. The pCID *is* the protocol's identity, and the protocol gets to define the directory's internal structure, headers, parent semantics, receipt format, and everything else. A human-readable slug may be appended to the directory name as a discovery convenience, but the pCID is canonical.

3. **The wire-lab channel-spec stays thin.** It defines: the outer convention (channels live at `channels/<pcid>--<slug>/`), the rule that messages do not declare their channel via a header, and the requirement that each channel-protocol's pCID names a spec defining everything inside the directory. It does **not** define `Parents`, `In-Reply-To`, header grammar for parent lists, transport vocabulary, or receipt format — those are named by individual channel-protocols.

These three principles together collapse several of the original DFs (26.1, 26.2, 26.3, 26.5, 26.6) into delegations to the per-channel-pCID protocol. The DFs that remain at the wire-lab level concern (a) the path-keying convention, (b) the timing of carve-out, (c) the operational meaning of "the pCID defines the structure," and (d) what to do with TE-24's existing v0 contract. Within the codex-perplexity channel-protocol *specifically*, an additional DF concerns the parent-header naming.

## Decision Forks (DFs)

### DF-26.1 — `channels/` subdirectory layout *(withdrawn — delegated to each channel-pCID)*

The wire-lab spec does not pick a layout for a channel's interior. Each channel-protocol's pCID names a spec that does. The wire-lab spec's only layout commitment is the outer keying convention (DF-26.7).

### DF-26.2 — Parent-link header *(withdrawn — delegated to each channel-pCID)*

Whether a parent-link header exists at all, what it is named, and whether it accepts one or many CIDs are all properties of each channel-protocol's spec. TE-26 surfaces only the conceptual shift toward DAG parent semantics; it does not pick header shapes for any specific channel-protocol.

### DF-26.3 — How a message declares its channel *(withdrawn — layer inversion)*

A message does not declare its channel. The transport (in repo-local terms, the directory path) carries it. No `Channel:` header, no `Channel-Type:` header, no per-message manifest reference.

### DF-26.4 — When does this TE's recommendation get applied?

#### Alt-4.A: Apply to `specs/channel-spec-draft.md` immediately, before any real channel traffic *(LOCKED — recommended)*

The new channel-spec carve-out, when written, embeds these conclusions directly: pCID-keyed channel paths, no `Channel:` header, parent-link semantics delegated to each channel-pCID. The TE-24 v0 contract (parents, receipts, message-id, etc.) is reframed at the same time as the contract of one specific channel-protocol-pCID, documented in a separate `specs/codex-perplexity-channel-draft.md` (or similar).

- **Easier**: the new wire-lab channel-spec doesn't carry vestigial 1:1 / single-writer assumptions or vestigial header definitions that would need to be moved later.
- **Harder**: the channel-spec carve-out is now blocked behind these decisions, but those decisions are now made.

#### Alt-4.B: Carve out the channel-spec now (using current TE-24 v0 contract), apply this TE's recommendations as the channel-spec's first revision

Land the channel-spec carve-out as the existing v0 (Prev-Message-CID, no transport taxonomy, no pCID keying). Then immediately revise the draft to apply this TE's recommendations. Two commits.

- **Easier**: each commit has a single concern; the carve-out is mechanical.
- **Harder**: a brief window where the channel-spec is on disk with the soon-to-be-replaced shape. Anyone reading the draft in that window sees a stale contract.

#### Alt-4.C: Defer this TE entirely until after first real traffic

Don't change anything until at least one real message has been written under the v0 contract. Then revisit.

- **Easier**: avoid rewriting on speculation.
- **Harder**: zero real messages exist today, and the assumptions being challenged (1:1 is the only transport, channel-identity is a slug, header set is fixed by wire-lab) affect every directory-layout decision that has to be made before the first message lands.

### DF-26.5 — Multi-CID parent-list serialization *(withdrawn — delegated to each channel-pCID)*

How a parent-link header (whatever it is named) serializes one or many CIDs is a property of each channel-protocol's spec. The wire-lab spec does not constrain it, and TE-26 does not survey the alternatives.

### DF-26.6 — Receipts (`IHave`) cross-channel naming *(withdrawn — delegated to each channel-pCID)*

The TE-24 contract's `IHave: <channel>:<cid>` form references a possibly-different channel from the one carrying the receipt. Whether the receipt names the referenced channel by pCID, slug, or something else is a property of the channel-protocol that defines the receipt format. The wire-lab spec does not constrain it.

### DF-26.7 — How is a channel directory keyed under `channels/`?

The path needs to be both (a) canonical in terms of protocol identity (so different channel-protocols can never collide) and (b) navigable by humans who know the channel by name.

#### Alt-7.A: Pure pCID paths, with a `channels/INDEX.md` mapping slugs to pCIDs

`channels/<pcid>/`. A separate index file maps human slugs to pCIDs.

- **Easier**: paths are unambiguously canonical; one source of truth (the directory name).
- **Harder**: humans cannot navigate `channels/` directly without the index; tab-completion shows opaque hashes; commit log entries (`channels/bafkrei.../...`) are unreadable.

#### Alt-7.B: Symlink convention (`channels/<slug>` → `channels/<pcid>/`)

Canonical directory is `channels/<pcid>/`; a symlink at `channels/<slug>` points at it for human navigation.

- **Easier**: humans use the symlink; tools and the canonical path stay unambiguous; symlinks are git-tracked.
- **Harder**: symlinks have known portability issues on some filesystems and clones; the symlink is a second piece of state that can drift from the canonical name.

#### Alt-7.C: Combined `channels/<pcid>--<slug>/` *(LOCKED — chosen)*

Directory name carries both. The pCID is canonical; the slug is a human-readable suffix that tools ignore (or use only for display).

- **Easier**: one directory, no symlinks, no separate index. Path is simultaneously canonical and human-readable. Tab-completion produces the slug naturally. Commit log entries are legible.
- **Harder**: long directory names; the pCID and slug are coupled in the name even though only the pCID is load-bearing.

### DF-26.8 — Operational meaning of "the pCID defines the structure"

When we say "the pCID defines the channel-protocol," what artifact actually does that work?

#### Alt-8.A: A frozen spec markdown is the human-readable contract; tooling reads it textually

Each channel-protocol-pCID corresponds to a frozen `specs/<...>-channel-<pcid>.md`. Tooling that wants to behave protocol-specifically must hand-implement that protocol's reader/writer.

#### Alt-8.B: Spec markdown plus a machine-readable companion (e.g., `channel.yaml`) describing layout/headers

The frozen spec is the human contract; an adjacent machine-readable file declares layout/headers/etc. for tooling to consume programmatically.

- **Easier**: tooling can be channel-protocol-agnostic, parameterized by the companion file.
- **Harder**: introduces a second source of truth; the companion file's schema is itself a wire-lab-level decision; YAML imposes new lock-in.

#### Alt-8.C: The code that reads the directory structure *is* the handler for that pCID *(LOCKED — chosen)*

Each channel-protocol comes with its own reader/writer code. The pCID identifies the protocol; the protocol identifies (via convention or naming) the code that speaks it. The frozen markdown is the human-readable contract that the code must implement. There is no machine-readable companion file at the wire-lab level.

- **Easier**: no second source of truth; no schema to maintain; matches "let each protocol name its own internals" recursively. Code-as-handler is honest about what protocol-specific behavior actually is.
- **Harder**: tools that want to display N channels of M different protocols need M handlers. That cost is real but it is the cost of polymorphism, not the cost of this design choice.

### DF-26.9 — Parent-link header *(withdrawn — delegated to each channel-pCID)*

Whether a channel-protocol exposes a parent-link header at all, what it is named, what it accepts, how it serializes, and whether it is optional are all properties of that channel-protocol's own spec. The wire-lab channel-spec does not mention parent-link headers. TE-26 surfaces only the *conceptual* shift toward a DAG-shaped message graph; it does not pick header names.

### DF-26.10 — What happens to TE-24's existing v0 contract?

#### Alt-10.A: Reframe TE-24's contract now as one specific channel-protocol's contract *(LOCKED — chosen)*

The wire-lab `specs/channel-spec-draft.md` is the *thin* outer rule (pCID-keyed paths, no `Channel:` header, code-as-handler principle). The TE-24 v0 contract (parents, receipts, message-id, message kinds, `IHave`, canonical-bytes) is documented separately as the contract of one specific channel-protocol — the codex-perplexity channel — in a draft spec doc. That spec doc has its own pCID once frozen.

- **Easier**: clean separation between wire-lab-level rules and channel-protocol-level contracts. Future channel-protocols can reuse the wire-lab outer rules without inheriting codex-perplexity-specific assumptions.
- **Harder**: two new spec docs to maintain (the thin wire-lab one and the codex-perplexity one) instead of one fat one.

#### Alt-10.B: Keep TE-24's contract as wire-lab-global default

Keep the contract at the wire-lab level until a second channel-protocol shows up demanding different rules.

- **Easier**: one document; no premature splitting.
- **Harder**: bakes in single-channel-protocol assumptions; future second-protocol arrival is a more disruptive split.

#### Alt-10.C: Reframe TE-24 now (Alt-10.A) AND freeze the codex-perplexity channel-protocol pCID immediately

Eat the dogfood.

- **Easier**: forces the freeze discipline to be exercised.
- **Harder**: premature; we have not sent any real traffic on this channel; freezing now locks decisions we will revisit once we do.

## Scenario analysis

### S1 — Codex publishes a status update intended for both Perplexity and Steve

A "status update" is conceptually broadcast (T2 / T7), not unicast.

- Under the locked design, Codex's status feed is its own channel keyed by its own channel-protocol-pCID, e.g., `channels/<pcid-for-status-feed-protocol>--codex-status/`. The codex-perplexity channel-protocol does not enter into this; the status-feed channel-protocol is a separate spec, with its own (possibly different) parent-header rules and receipt rules.

### S2 — Perplexity proposes a design change in a dedicated review channel; Codex and Steve both reply

Three senders in one channel, each citing prior messages by others.

- Under DAG parent semantics: Codex's reply names Perplexity's message as a parent. Steve's reply names Perplexity's message as a parent and possibly also Codex's. The DAG handles this naturally without a single-writer-log assumption. The exact header used is a property of the codex-perplexity channel-protocol's spec, not TE-26.

### S3 — A long-running 1:1 between Codex and Perplexity occasionally gets a third participant

- Under flat-per-channel layout *within* the codex-perplexity channel-protocol: the channel admits the third sender naturally. No structural change. (The channel-protocol's spec defines its own layout under `channels/<pcid>--codex-perplexity/`; the wire-lab spec is silent on this.)

### S4 — Two senders publish concurrently to the same multi-writer channel; neither saw the other's message before sending

Standard concurrent-edit case.

- Under DAG parent semantics: both messages name the same parent(s) (pointing at the previous frontier). The DAG has a fork; subsequent messages can name both as parents, healing the fork. No structural problem.

### S5 — Channel-protocol evolves from "1:1-only" to "multi-writer" over time

- Under pCID-keyed channel paths: the channel-protocol-pCID is *frozen* at freeze time. If a channel needs different rules, a *new* channel-protocol-pCID is minted, and migration happens by writing under a new directory. The frozen pCID guarantees that the protocol's rules are stable for as long as that channel exists.
- A channel that wants to evolve in place would need to be using a channel-protocol-pCID whose contract permits the evolution (e.g., a generic-DAG channel-protocol that allows arbitrary sender sets from the start).

### S6 — Tooling needs to render a thread

Reader wants to display "this message and the messages it is replying to."

- Under code-as-handler (DF-26.8 Alt-8.C): the tool dispatches on the channel-protocol-pCID to that protocol's reader. The reader walks whatever parent-link representation that protocol defines, recursively. No wire-lab-level threading code exists.

### S7 — An observer joins the channel mid-stream and asks "have I seen everything?"

Standard "completeness check" case.

- The channel-protocol's spec defines the answer. For a DAG-based protocol, the observer fetches all files in the channel directory, verifies every parent CID is present, and computes the channel frontier as messages no other message names as parent. `IHave: <channel-pcid>:<cid>` (or whatever the protocol calls it) names a contiguous prefix-frontier acknowledgement.

## Conclusions

1. **A message does not declare its channel.** The channel is the transport carrying it. No `Channel:` header. (DF-26.3, withdrawn.)
2. **Channel directories under `channels/` are keyed by `<pcid>--<slug>`** (DF-26.7 Alt-7.C). The pCID is canonical protocol identity; the slug is a human-readable suffix.
3. **Each channel-protocol-pCID names a spec that defines the directory's interior:** layout, headers (including any parent-link header), parent semantics, receipt format, message-kind vocabulary, etc. The wire-lab channel-spec does not define these. (DF-26.1, 26.2, 26.5, 26.6 — all withdrawn / delegated.)
4. **The code that reads the directory structure is the handler for that pCID** (DF-26.8 Alt-8.C). There is no machine-readable companion file at the wire-lab level.
5. **A channel's message graph is conceptually a DAG.** Each message can have zero or more parent messages within the same channel. The single-writer-log shape (every message has exactly one parent in its own log) is a special case. How any channel-protocol expresses this — header names, serialization, optionality, or whether a parent-link header exists at all — is delegated to that channel-protocol's spec.
6. **TE-24's v0 contract is reframed now as the codex-perplexity channel-protocol's contract** (DF-26.10 Alt-10.A), in a separate draft spec. The wire-lab channel-spec ships thin.
7. **Apply these conclusions in the channel-spec carve-out itself** (DF-26.4 Alt-4.A). No two-step revision dance; no vestigial 1:1 assumption in the new wire-lab spec.

## Implications

- **The single-CID `Prev-Message-CID` of TE-24 is conceptually subsumed by DAG parent semantics**, but its concrete shape (header name, serialization, list-vs-singleton, optionality) in the codex-perplexity channel-protocol's contract is a decision belonging to that channel-protocol's spec, not to TE-26.
- **No `Channel:` header anywhere.** Existing references in `channels/README.md`, the harness-spec-draft, and DR-009 must be removed in the carve-out commit.
- **Channels live at `channels/<pcid>--<slug>/`.** The codex-perplexity channel's directory will, once its channel-protocol-pCID is minted, become `channels/<that-pcid>--codex-perplexity/`. Until then, no on-disk directory is created (consistent with TE-24's "no subdirectories yet" stance).
- **Wire-lab channel-spec is thin.** It defines: pCID-keyed directory naming under `channels/`, the absence-of-`Channel:`-header rule, the requirement that each channel-protocol-pCID names a spec defining everything inside the directory, and the code-as-handler principle. That is roughly all.
- **codex-perplexity channel-protocol draft spec is its own document** (e.g., `specs/codex-perplexity-channel-draft.md`), inheriting the full TE-24 v0 contract. Whether and how it adopts DAG parent semantics (replacing the single-CID `Prev-Message-CID`) is a decision that lives inside that spec doc, not in TE-26.
- **Channel-spec freeze and codex-perplexity-channel-spec freeze are independent** future events, each minting their own pCID under Steve's signature on a corresponding `merge-<slug>-spec` promise.
- **TODO 012's scope expands** to cover the carve-out plus the codex-perplexity-channel draft spec (subtasks 012.7+).

## Decision status

LOCKED:
- DF-26.1 — withdrawn (delegated to each channel-pCID).
- DF-26.2 — withdrawn (delegated to each channel-pCID).
- DF-26.3 — withdrawn (a message does not declare its channel; layer inversion).
- DF-26.4 — Alt-4.A (apply in the channel-spec carve-out).
- DF-26.5 — withdrawn (delegated to each channel-pCID).
- DF-26.6 — withdrawn (delegated to each channel-pCID).
- DF-26.7 — Alt-7.C (`channels/<pcid>--<slug>/`).
- DF-26.8 — Alt-8.C (code-as-handler; no machine-readable companion).
- DF-26.9 — withdrawn (delegated to each channel-pCID; TE-26 surfaces the DAG concept only).
- DF-26.10 — Alt-10.A (reframe TE-24's contract now as the codex-perplexity channel-protocol's contract).

Recorded principle: *channel identity, layout, and message structure are named by the pCID; the wire-lab spec defines only the outer envelope and the `channels/<pcid>--<slug>/` convention. The handler for a pCID is the code that reads its directory structure.*

## Implications for follow-on work

- **TODO 013 (anticipated)**: drive these locked alts to a DR; carve out `specs/channel-spec-draft.md` (thin) and `specs/codex-perplexity-channel-draft.md` (the full TE-24 v0 contract, with any DAG-related revisions decided inside that spec, not by TE-26); update `channels/README.md`; remove channel material from `specs/harness-spec-draft.md`; update DR-009 and TODO 012.
- **TODO 014 (anticipated)**: first real channel-message exchange under the new contract, exercising the codex-perplexity channel-protocol's parent-link mechanism in both single-writer and multi-writer paths.
- **TE-27 (anticipated)**: should the wire-lab spec define a small companion convention for channel-protocols to publish their own pCID on first use, so receivers can discover the protocol-spec from a stranger's first message? This is the question DF-26.8 Alt-8.B raised and the locked Alt-8.C deferred.
- **TE-28 (anticipated)**: receipts under multi-writer channels — does `IHave: <channel-pcid>:<cid>` need to become a vector to acknowledge a frontier rather than a single tip? (Decision belongs in each channel-protocol's spec, but the question is general enough to warrant a TE.)
