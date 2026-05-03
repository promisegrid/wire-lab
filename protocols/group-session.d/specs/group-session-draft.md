# Group Transport Spec (DRAFT)

*This is the wire-lab's first transport-protocol spec. It defines the small-finite-closed-group transport-protocol: N≥2 participants, all-see-all visibility, multi-writer DAG of messages, append-only persistence. The Codex↔Perplexity exchange is the N=2 instance, not a separate spec. This file is a draft and is subject to revision; once frozen, its pCID will name this protocol class for all time. See `specs/MANIFEST.md` for freeze status.*

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted. Cross-references to it in other repo files use `specs/group-transport-draft.md` (path) until freeze; after freeze they will use the pCID.

## Purpose

This spec defines the **interior** of a `transports/<this-pcid>--<slug>/` directory: how messages are encoded as files, how they are named on disk, how they reference each other (the DAG of parents), and the v0 receipt mechanism.

The outer convention — that transport directories live under `transports/` keyed by the transport-protocol's pCID, that messages do not declare their transport via a header, and that the code reading a transport directory is the handler for that pCID — is governed by `specs/transport-spec-draft.md` (cited as a path rather than a Markdown link because both files are drafts; once `transport-spec-draft.md` is frozen, this reference will be replaced with its pCID-named successor). This spec inherits all four of those principles.

## Sources

This spec is locked by the conclusions of:

- [TE-24](../docs/thought-experiments/TE-20260430-204108-group-transport-envelope.md): the group-transport envelope: `grid <pcid>` carrier, canonical-bytes rules, explicit-promise body requirement, and `Date`/`From` as conveniences. The TODO 013 carve-out supersedes the original `Prev-Message-CID:`, `IHave:`, and `Kind:` headers from the TE's first drafting; the locked v0 contract uses `Parents:` (DAG links), body-level acknowledgement (receipts), and no `Kind:` header. The `Message-ID:` header described in TE-24 is also dropped: under §2 (filename = CID), the message CID is the message's identifier and a separate human-readable identifier creates two competing identities for the same message and is omitted.
- [TE-26](../docs/thought-experiments/TE-20260430-215624-channel-transport-types-and-threaded-replies.md): the conceptual shift to DAG-shaped message graphs (zero-or-more parents per message).
- [TE-27](../docs/thought-experiments/TE-20260501-021921-transports-rename-and-axes-of-differentiation.md): the per-axis meta-rule that locks small-finite-closed-group with N≥2 as a single protocol class (cardinality is a parameter, not a contract boundary, except at extremes).
- [DR-009](../DR/DR-009-20260430-204108-group-transport-envelope.md): the active decision request governing the group-transport envelope and the freeze gate.

## What this spec covers

A **group transport** under this protocol class is:

- **N≥2 participants.** Cardinality is a parameter; the contract does not change between N=2 and small-finite-closed-N. Very-large or unbounded membership is out of scope (see TE-27 axis A: it crosses a contract boundary).
- **Closed.** Membership is fixed at transport creation. New members do not join an existing instance; if membership changes, a new transport instance is minted.
- **All-see-all visibility.** Every participant observes every message. There is no hub mediation, no subset addressing, no topic filtering. (TE-27 axis B: visibility is a contract boundary; this protocol picks all-see-all.)
- **Multi-writer DAG of messages.** Any participant may post at any time. Messages cite zero or more prior messages as parents, forming a DAG. There is no single-writer ordering. (TE-27 axis F: message-graph shape is a parameter within this spec.)
- **Append-only persistence.** Once a message file is committed to a transport directory, it is not modified or deleted. Compaction, retention bounds, and ephemerality are out of scope. (TE-27 axis E: persistence is a parameter; this protocol picks append-only.)
- **Receipts in message bodies.** Acknowledgement is a kind of message body, not a header. There is no `IHave:` header. (See §6.)

## The v0 message contract

### §1. Subdirectory layout: flat

A transport directory under this protocol has **no subdirectories.** All messages of the transport live as sibling files directly under `transports/<this-pcid>--<slug>/`.

```
transports/<this-pcid>--<slug>/
    <message-cid-1>.txt
    <message-cid-2>.txt
    <message-cid-3>.txt
    ...
```

Rationale: the DAG of `Parents:` links carries all the ordering information any reader needs. Subdirectory structure (per-sender, per-direction, per-date, etc.) would either duplicate the DAG (redundant), pick a privileged axis the protocol does not have (sender or direction), or be presentational (and so belong in a viewer, not in the on-disk format). Flat is honest. (T6 in TODO 013 carve-out, locked Alt-T6.A.)

### §2. Message filename: `<message-cid>.txt`

Each message file is named by its **message CID** (per §3), with the `.txt` suffix appended. The `.txt` extension reflects the canonical-bytes contract from §4 (UTF-8 text with LF line endings); it lets ordinary editors, viewers, and text-search tools work on message files without ceremony.

Four properties follow from naming files by CID:

- **Filename collisions are impossible by construction.** Two messages with the same canonical bytes have the same CID and are the same message; two messages with different canonical bytes have different CIDs and therefore different filenames.
- **Filenames are content-verifiable.** A reader can verify a message's integrity by re-computing the CID over the file's bytes and comparing it to the filename, without consulting any other file or header. The filename and the canonical bytes are mutually self-checking.
- **Two writers cannot collide on the wire.** Under git bindings such as the per-author-branch binding (§9), two participants posting the same logical message produce identical bytes, identical CIDs, and identical filenames; git treats this as a clean union with no conflict. Two participants posting different messages produce different filenames in the same directory, also with no conflict.
- **Append-only is structurally enforced.** Editing an existing message after the fact changes its canonical bytes and therefore its CID, which makes the edit a different file entirely. The original file remains; the edit appears as a sibling. Mutation of an existing message in place is not expressible.

Readers locate messages by CID through the `Parents:` DAG; the filename equals that CID, so a `Parents:` reference to a CID identifies a file directly without requiring an index.

### §3. Message identity: CIDv1 over canonical bytes

The **message CID** is the load-bearing identifier of a message. It is computed as CIDv1 with `base32` multibase, `sha2-256` multihash, and `raw` codec, taken over the full canonical bytes of the message file (every byte from the first byte of `grid ` through the trailing newline, inclusive).

Two messages with the same canonical bytes have the same message CID. Two messages with different canonical bytes (even differing only in whitespace) have different message CIDs. This is the property `Parents:` relies on.

### §4. Message envelope

A message file consists of, in order:

1. The carrier line: `grid <pcid>\n` where `<pcid>` is this transport-protocol's pCID.
2. A blank line: `\n`.
3. A header block: zero or more headers, each of the form `Header-Name: value\n`, in the canonical order specified in §4.7.
4. A blank line: `\n`.
5. The body: free-form UTF-8 text containing at least one explicit `I promise ...` clause (see §5).
6. A trailing newline at EOF: `\n`.

Canonical bytes are UTF-8 with LF (`\n`) line endings. CRLF is not accepted. There is exactly one blank line between the carrier line and the headers, and exactly one blank line between the headers and the body. The file ends with exactly one trailing newline.

#### §4.1 Carrier line

The first line is `grid <pcid>` where `<pcid>` is this transport-protocol's pCID. The carrier line is mandatory. Readers MUST verify it before parsing the rest. This satisfies the outer transport-spec's "messages are dispatched on first line" property without violating "messages do not declare their transport" — the pCID names the protocol-class contract, not which transport-instance the message belongs to. The transport-instance is identified by the directory the file lives in.

#### §4.2 Headers, generally

Headers are line-oriented. Each header is `Name: value\n`. Header names are case-sensitive. Header values do not span lines (no continuation lines). Trailing whitespace on a header value is significant and should be avoided; canonical encoders SHOULD strip trailing whitespace before computing the message CID.

#### §4.3 No `Message-ID:` header

This protocol does NOT have a `Message-ID:` header. The message CID (§3), which is also the filename (§2), is the message's identifier; a separate human-readable identifier creates two competing identities for the same message and is therefore omitted. Author identity lives in `From:` (§4.5); authoring time lives in `Date:` (§4.4); navigation aids that humans want (such as a short slug) belong in body prose, not in the envelope.

#### §4.4 `Date:` (mandatory, single-valued)

UTC timestamp in `YYYY-MM-DDTHH:MM:SSZ` format. The `Date:` header records when the sender claims the message was authored; it is not authoritative for ordering (which comes from the DAG).

#### §4.5 `From:` (mandatory, single-valued)

The sender's identity, as a free-form printable-ASCII string. The protocol does not specify the identity scheme (no requirement that this be a key, an email, or a pCID); future TEs may tighten this.

#### §4.6 `Parents:` (optional, single line, space-separated message CIDs)

Identifies the message CIDs of zero or more prior messages this message acknowledges as direct ancestors in the DAG.

- The header is **optional.** A message with no `Parents:` header has no parents named (it is a root, or its author chose not to cite any). Absence and an empty header are NOT distinguished; canonical encoders SHOULD omit the header entirely when there are no parents. (DF-T3, Alt-T3.B.)
- The header is **single-line, space-separated.** All parent CIDs appear on one physical line, separated by single ASCII spaces. There is no multi-line `Parents:` form. (DF-T2, Alt-T2.A.)
- Each value is a base32-encoded CIDv1 message CID (no `cid:` or other prefix).
- Order is significant for canonical-bytes purposes (the bytes hash differently if the order changes), but is NOT semantically privileged: readers MUST treat the parent set as a multiset for graph-walking purposes.
- The header name is `Parents:` (plural). There is no `Prev-Message-CID:` header in this protocol; single-parent messages use `Parents: <one-cid>`. (DF-T4, Alt-T4.A.)

The `Parents:` mechanism is how the DAG is realized on disk. A reader that walks back through `Parents:` from any message can reconstruct that message's causal past, bounded by the transport's first message (which has no parents).

#### §4.7 Canonical header order

In canonical bytes, headers MUST appear in this order, and any header that is absent simply does not appear (no placeholder line):

1. `Date:`
2. `From:`
3. `Parents:` (if present)

Future versions of this spec MAY add headers; they will be inserted at locked positions in this order list. Unknown headers (from a hypothetical future revision) MUST NOT be silently dropped; readers MUST reject messages they cannot fully parse, since the message CID covers the unknown bytes and a reader that strips them changes the message identity.

#### §4.8 No `Kind:` header

This protocol does NOT have a `Kind:` header. The original TE-24 sketch included one as a human-oriented convenience; on review, message kind is presentational and varies per use case, so it is left to body convention rather than an envelope field. (Q1 in TODO 013 carve-out, locked Alt-Q1.A.)

#### §4.9 No `IHave:` header

This protocol does NOT have an `IHave:` header. Acknowledgement is a body concern (see §6), not an envelope concern.

#### §4.10 No `Transport:` header

Per the outer transport-spec, the message envelope contains no `Transport:` or `Transport-Type:` header. The carrier line's pCID names the protocol class; the directory the file lives in identifies the transport instance.

### §5. Body: explicit promise prose

The message body MUST contain at least one explicit `I promise ...` clause. This requirement preserves the message's legibility as promise-theory discourse and prevents the envelope from devolving into a fixed schema whose semantics hide in field names.

The body is free-form UTF-8 text. Markdown is conventional but not required. There is no maximum length (within the transport-instance's filesystem limits). The body MAY contain additional structured content (code blocks, quoted prior messages, etc.) so long as the canonical-bytes rules are observed.

### §6. Receipts: acknowledge in the body, not the header

A message that acknowledges another message does so by saying so in the body, not by adding an envelope header.

A v0 acknowledgement has the shape of a normal message whose body explicitly cites the message CID(s) being acknowledged and contains a promise to that effect. Recommended body convention:

```
I promise that I have observed and accepted the following message(s):

- bafk...abc
- bafk...def
```

The acknowledgement message itself participates in the DAG: it cites the messages it acknowledges in its `Parents:` header (so graph-walkers see the relationship structurally) AND in its body prose (so humans and LLM readers see the relationship in plain text). The two MUST be consistent — every CID in the body acknowledgement list MUST also appear in `Parents:`, and conversely. (DF-T5.)

This v0 receipt scheme is per-message: each acknowledgement explicitly lists the message CIDs it accepts. There is no compact "I have everything up to frontier F" form in v0; cumulative-prefix or frontier-style acknowledgement is deferred to a future TE (Q2 in TODO 013 carve-out, deferred).

### §7. Persistence

A transport directory under this protocol is **append-only.** Once a message file is committed (in git terms: once it appears in a commit on the canonical branch), it is not modified or deleted. Editing a message after-the-fact would change its canonical bytes and therefore its message CID, breaking every `Parents:` reference to it.

The protocol does not specify retention bounds; transports under this protocol are presumed durable for the lifetime of the repo. Compactable, bounded-retention, or ephemeral variants are different transport-protocols with different pCIDs.

### §8. Membership

Membership is **closed and fixed at transport creation.** The set of `From:` values that may appear in a transport instance is determined by the social/organizational context of the transport's creation; the protocol does not enforce a membership list cryptographically in v0. A transport's slug typically names the participants (e.g. `codex-perplexity`).

If membership changes — a participant leaves, a new participant joins — a new transport instance is created (a new directory under `transports/`, with a fresh slug). The old transport remains as immutable history.

### §9. Per-author-branch git binding with content-addressed merge (non-normative)

This section describes the conventional git binding used by transport instances of this protocol that ride a shared git remote as their wire. It is non-normative: the protocol's contract is the on-disk shape of `transports/<pcid>--<slug>/` plus the canonical bytes of each message file. The git binding is one way of placing those files in front of the participants such that the append-only and DAG-of-parents semantics are preserved across multiple writers.

#### §9.1 Branch ownership

- **Each participant has their own write branch named `<author-id>/main`.** Examples: `alice/main`, `bob/main`, `carol/main`. The branch name is per-participant and stable across the lifetime of the transport.
- **A participant authors message files only on their own `<author-id>/main` branch.** No participant ever authors a new message file on another participant's branch.
- **A participant may, and is expected to, propagate other participants' message files onto their own branch.** This is not authoring; the propagated files are byte-identical to their origin (and therefore have the same CID and the same filename). Propagation is the merge step described in §9.3.

#### §9.2 Filename = CID makes merges trivial

Under §2, every message file's name is its message CID. This has three consequences for git merges:

- **No participant can name a file the same as another participant's distinct file.** Distinct canonical bytes produce distinct CIDs which produce distinct filenames.
- **Two participants who independently obtain the same message produce byte-identical files at the same path.** Git treats the two adds as a no-conflict union and stores one copy.
- **Forwarding a message from one branch to another is a verbatim file copy.** No content edit, no metadata change, no per-branch annotation; the file's CID-named path is its identity.

#### §9.3 Receive-merge-push-then-optionally-post cycle

A participant's transport interaction proceeds in two phases per cycle. The merge phase is mandatory whenever new messages are observed; the post phase is optional.

**Merge phase (mandatory when new messages are observed):**

1. `git fetch --all` to retrieve every known participant's `<author-id>/main`.
2. For each known participant's branch other than the agent's own, list `*.txt` files under `transports/<pcid>--<slug>/` that are not present on the agent's own branch.
3. For each such file, verify integrity per §3 (re-compute CID, compare to filename) and per §4 (envelope structure). Files that fail verification are rejected from the merge.
4. Copy the verified files into the agent's working tree on the agent's own branch.
5. Commit with a message of the form `transport: merge <count> messages from <branches>`.
6. Push the agent's own branch.

**Post phase (optional):**

7. If the agent has a new message to author, write it as a fresh `*.txt` file under `transports/<pcid>--<slug>/` per §4 (envelope), §5 (explicit promise), §6 (body-as-receipt where applicable), and §4.6 (`Parents:` set to the message CIDs of direct ancestors).
8. Compute the file's CID per §3 and rename the file accordingly so that filename equals CID per §2.
9. Commit with a message of the form `transport: post <message-cid>`.
10. Push the agent's own branch.

The merge phase precedes the post phase because the agent's `Parents:` references in any new message must be CIDs the agent has already verified and propagated onto its own branch; otherwise readers fetching only the agent's branch may see a `Parents:` reference to a message they do not have. Posting a message whose ancestors are not yet propagated is permitted but discouraged: it requires readers to fetch the original author's branch to resolve the `Parents:` reference.

#### §9.4 Convergence and consistency

Under steady-state operation, every participant's `<author-id>/main` branch eventually contains every message anyone has ever posted to the transport. The transport state is **eventually consistent**: any divergence between branches is bounded by the time since the last fetch-merge-push cycle on each branch.

The git commit graph is incidental to the message graph. Readers walk `Parents:` links per §4.6 to reconstruct causal ordering, and ignore the commit graph for that purpose. A merge commit on `<author-id>/main` that propagates twelve messages from another branch does not, by itself, mean those twelve messages are causally ordered before any subsequent message on `<author-id>/main`; the `Parents:` headers are the only authority for causal order.

#### §9.5 Infrastructure files

Files inside `transports/<pcid>--<slug>/` that are not message files (typically a `README.md` navigation aid) are infrastructure, not protocol messages. The merge phase of §9.3 does NOT propagate infrastructure files; it propagates only `*.txt` message files. Infrastructure files propagate via the branch participants forked from when joining the transport. Edits to infrastructure files are coordinated out-of-band and do not participate in the DAG.

#### §9.6 Append-only is structurally enforced

Per §7 and the filename-equals-CID rule of §2, mutation of an existing message in place is not expressible: any byte change produces a different CID and therefore a different filename, leaving the original file untouched. The only way to functionally remove a message is to refuse to propagate it, which other participants can detect by comparing what they have on their own branches to what the refusing participant's branch carries. The standing "never force-push" rule on `<author-id>/main` branches makes accidental file deletion recoverable from the remote's reflog and intentional deletion observable.

#### §9.7 Membership

Membership under this binding is the set of `<author-id>/main` branches a participant is configured to fetch and propagate from. The closed-and-fixed property of §8 is satisfied when participants share a fixed list of branches; an unrecognized branch name is, by convention, ignored on read until membership is explicitly extended.

#### §9.8 Other bindings remain compatible

Other git bindings are possible and remain compatible with the protocol's on-disk contract. For example, a single shared write branch with merge-on-pull is workable but produces fast-forward rejections under concurrent writes; a per-thread branch binding is workable but loses the all-see-all property unless threads are merged. The per-author-branch binding above is the recommended starting point because it eliminates merge conflicts by construction (§9.2) and reaches consistency through propagation rather than coordination (§9.3-9.4).

## Worked example

A two-participant transport between Codex and Perplexity, freshly created, before any messages exist:

```
transports/<group-transport-pcid>--codex-perplexity/
```

After Codex posts the first message:

```
transports/<group-transport-pcid>--codex-perplexity/
    bafkreigtaivld55rekcswfj26mo26e267m3ytzgflqb2qcclyiicpfzc6i.txt
```

(The CID in this filename is a placeholder for illustration; the real value is computed over canonical bytes.)

Where the file `bafkreigtaivld55rekcswfj26mo26e267m3ytzgflqb2qcclyiicpfzc6i.txt` contains:

```
grid <group-transport-pcid>

Date: 2026-04-30T20:31:14Z
From: codex@promisegrid.example

Hello, Perplexity. I promise to coordinate with you on the
group-transport-spec carve-out work and to record decisions in
DI-009-derived intent records.
```

After Perplexity replies citing Codex's message as a parent:

```
transports/<group-transport-pcid>--codex-perplexity/
    bafkreigtaivld55rekcswfj26mo26e267m3ytzgflqb2qcclyiicpfzc6i.txt
    bafkreih5xxxxx2example2cid2for2reply2filename2placeholder2onlyaa.txt
```

Where the second file (the reply) contains:

```
grid <group-transport-pcid>

Date: 2026-04-30T20:37:14Z
From: perplexity@promisegrid.example
Parents: bafkreigtaivld55rekcswfj26mo26e267m3ytzgflqb2qcclyiicpfzc6i

I promise that I have observed and accepted the following message(s):

- bafkreigtaivld55rekcswfj26mo26e267m3ytzgflqb2qcclyiicpfzc6i

I promise to begin work on the group-transport-draft v0 contract.
```

(The CID in this example is a placeholder; the real value is computed over canonical bytes.)

## What this spec does NOT specify

- The pCID itself. The pCID for this protocol is minted at freeze (see §Freeze gate).
- Any concrete protocol-instance pCID under this class. Each transport instance gets a slug for human navigation; the protocol-class pCID is shared across all instances.
- Message-graph algorithms (frontier computation, lowest-common-ancestor, conflict resolution between concurrent writers). Those are reader concerns, not on-disk-format concerns.
- Cryptographic signing of `From:`, message bodies, or message CIDs. Future revisions or successor protocols may add this; v0 does not.
- Retention or compaction behavior beyond "append-only" (see §7).
- Membership-change semantics beyond "create a new transport instance" (see §8).
- Cumulative-prefix or frontier-style receipts (deferred to a future TE).

## Open questions

- **OQ-G1 (deferred):** Should `From:` be tightened to a key, a pCID, or some other structured identity in a future revision? Raised but not closed in v0.
- **OQ-G2 (deferred):** What does cumulative-prefix or frontier acknowledgement look like under this protocol's DAG model? Q2 in TODO 013 carve-out, deferred to its own future TE.
- **OQ-G3 (deferred):** When two writers concurrently extend the same parent set, the DAG fans out; subsequent messages typically cite both leaves to converge. Should the protocol prescribe any fan-in obligation, or is this entirely a writer convention? v0 leaves it to convention.
- **OQ-G4 (deferred):** Should there be a canonical "transport-creation" or "genesis" message at the root of every transport's DAG, naming the participants and the slug? v0 does not require one; first message is whatever the first writer produces.
- **OQ-G5 (deferred):** For N>2 instances, are there observability or fairness considerations (e.g., should the protocol require that every participant's `From:` actually appear before some milestone)? v0 is silent.

## Freeze gate

This spec graduates to frozen status when:

1. `specs/transport-spec-draft.md` is itself frozen (this spec depends on its outer rules).
2. At least one real transport instance under this protocol has been created and exchanged at least one round-trip of messages, exercising §3 (CID computation), §4 (envelope), §4.6 (`Parents:` DAG link), §6 (body-as-receipt), and §7 (append-only). The codex-perplexity instance is the anticipated first.
3. Steve signs a `merge-group-transport-spec` promise authorizing the freeze.
4. `tools/spec freeze group-transport-spec` mints the pCID, snapshots the file, and appends the manifest entry.

Until then, the spec lives at `specs/group-transport-draft.md` and is a working draft.
