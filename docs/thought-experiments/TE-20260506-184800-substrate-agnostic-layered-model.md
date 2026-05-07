# TE-38: Substrate-agnostic layered model (L5/L6/L7) and L6 CAS subtree

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-20260506-184800

## Status

decided

(Vocabulary refactor over the architecture already locked in TE-29.
The locked path scheme, simulated-repo shape, layer
decomposition, and instance-feed-declaration mechanism are all
unchanged from TE-29. This TE renames the layers, cites the
100-year goal as the foundation those layers are answering to,
and extends the path scheme with an L6 CAS subtree so every
chunk lives in content-addressable storage with no
size-conditioned exception.)

## Why this TE now

TE-29 locked the architecture but used the framing "L4-binding
layer" carried over from session-of-2026-05-01 vocabulary. Two
problems surfaced as later TEs and conversations accumulated:

1. **"Binding" was the wrong word.** The L4-binding role is
   really a *feed* — a per-substrate stream of authenticated
   chunks the protocols above can subscribe to. "Binding"
   conflates the spec (what you promise about a substrate)
   with the active subscription (what an instance actually
   reads). Memory-of-record (`projects.wire_lab` and the
   2026-05-04/05/06 sessions) reflects the rename to "feed";
   TE-29 is the last place still saying "binding".

2. **The layered model was never named foundationally.** The
   L5/L6/L7 split — feed protocols / CAS protocols / group
   (forum) protocols — exists in conversation memory and in
   the dependency-sorted TE roster, but no TE cites it as a
   foundational consequence of the 100-year goal (TE-28).
   Without that citation, future-Steve and future-LLM readers
   cannot reconstruct *why* the layers carve where they do
   without rereading the entire conversational history.

In addition, conversation on 2026-05-06 raised one substantive
question over and above vocabulary: TE-29 stops at
five path levels and is silent on where *chunks* live. An
earlier draft (Alt-M.1) tried to make "small messages are leaf
chunks; large messages get pointer files" work, but Steve flagged
on 2026-05-06 17:54 PT that "all chunks should go into CAS;
exceptions make it complicated." This TE locks the no-exception
shape (Alt-M.4 in TODO 22, Q-22.6) as the L6 layer of the
substrate-agnostic model.

## Decision under test

DUT-38: **Adopt L5/L6/L7 as the named, citable layered model
under the 100-year goal, rename "binding" to "feed" in current
and future spec language, and extend TE-29's path scheme with an
L6 CAS subtree (`cas/<cas-protocol-pCID>/<chunk-cid>`) such that
every message file in `transports/` is unconditionally a CBOR
pointer into CAS.**

Three independent sub-decisions, each treated below:

- **DF-38.1** (vocabulary): "binding" → "feed"; "forum" → "group";
  introduce L5/L6/L7 as canonical layer names.
- **DF-38.M** (CAS subtree): all messages are CBOR pointers; no
  size-based exception.
- **DF-38.G** (citation): cite the 100-year goal (TE-28) as the
  load-bearing constraint each layer is answering to.

## Foundational invariants (citing TE-28)

The 100-year goal (TE-28) names a constraint set the layered
model exists to answer. The L5/L6/L7 split is not a conventional
networking-stack layering; it is a substrate-agnostic split
chosen so that each layer can survive independent obsolescence
of every other layer. The invariants the layers honor:

1. **No central registry.** No layer assumes a directory,
   namespace authority, or governance body that can revoke or
   reissue identifiers. Every identifier at every layer is a
   pCID derived from content the parties already have.

2. **Multi-generational durability.** Every layer must be
   reconstructible from on-disk artifacts after every named
   participant has left the project. Specs are content-
   addressable; instances declare themselves through paths
   (TE-29's path-as-declaration); chunks are addressable
   independent of the wire that delivered them.

3. **Adversarial-by-default.** Every layer assumes Mallory
   participates. A feed (L5) authenticates frames so a hostile
   wire cannot inject. A CAS (L6) authenticates chunks by hash
   so a hostile feed cannot substitute. A group protocol (L7)
   authenticates assertions by signing key so a hostile CAS
   cannot fabricate authorship.

4. **Protocol forking is normal.** Two parties may author
   competing L5 feeds, competing L6 CAS schemes, or competing
   L7 group protocols without forking the layers above or
   below. The path scheme makes the choice mechanical: each
   level is a pCID that names exactly which spec is in force.

5. **Trust accrues per-burden.** TE-29's guarantee that
   trust attaches to a (party, promise) pair, not to a party
   alone, holds across all three layers. Each layer asserts
   only what it can keep.

6. **Signing key is the only structural lock.** No layer
   depends on hardware, vendor identity, hostname, or DNS;
   the only durable identity primitive is a signing keypair
   the operator controls.

These six invariants are the load-bearing reason layers carve
where they do. They are restated here, not introduced; the
canonical home is TE-28.

## Locked shape: L5/L6/L7

### L7 — group (forum) protocols

L7 protocols define **what a message means inside a community
of agents**: assertion types, vocabulary, conditional-promise /
reciprocal-promise pairings, group-membership conventions, the
shape of the discrepancy-handling discourse. The wire-lab's
own dev process is an L7 protocol (the `wire-lab-devs` group);
PromiseGrid contracts and consensus conversations are also
L7. The renaming of "forum" to "group" in this TE matches
the memory record from the 2026-05-04 conversation; "forum"
remains acceptable in human prose where "group" is awkward,
but path components and spec filenames use "group".

L7 is responsible for nothing about transport, framing, or
chunking. An L7 protocol that runs over UDP feeds, file-drop
feeds, or hand-carried USB drives is the same L7 protocol.

### L6 — CAS (content-addressable storage) protocols

L6 protocols define **how chunks are addressed and stored**.
The foundational L6 promise is *hash-based addressing*: every
chunk has a pCID derived from its bytes under a named hash
function and codec. Concrete L6 specs include:

- Rabin chunking (or equivalent content-defined slicing) —
  defines how a byte stream is split.
- Merkle-tree assembly — defines how chunks are composed
  back into messages.
- CIDv1 codec selection — defines which hash function and
  multibase encoding the pCIDs use.

These three are commonly bundled, but the layered model treats
them as separable: a community can swap CIDv1 for a successor
codec without changing chunking, or swap the chunker without
changing the codec. The L6 protocol pCID names exactly which
bundle is in force for any given chunk.

### L5 — feed protocols

L5 protocols define **how authenticated chunks reach
participants over a specific substrate**. The substrate may
be UDP, TCP, WebSocket, MQTT, file-drop on a shared filesystem,
SMTP, an out-of-band channel like LoRa, or a physically
delivered medium. The L5 spec promises:

- A framing convention sized to the substrate.
- An authentication convention (per-frame signing or
  group-key MAC) so a hostile substrate cannot inject.
- A retransmission / liveness convention if applicable.

L5 sits between L6 and the underlying wire. It does not own
chunking (that is L6) and it does not own the meaning of
messages (that is L7).

## Locked shape: extended path scheme with L6 CAS subtree

TE-29 locked the message-side path:

    transports/<wire>/<feed-pCID>/<group-pCID>/<msg-pCID>/<msg-id>.txt

(TE-29 wrote `<binding-pCID>`; this TE renames the level to
`<feed-pCID>` to match the L5 / "feed" vocabulary. The shape
is otherwise unchanged.)

This TE adds the L6 CAS subtree at the repo root, beside
`transports/`:

    cas/<cas-protocol-pCID>/<chunk-cid>

and locks the relationship between the two trees:

> **Every message file in `transports/` is unconditionally a
> CBOR pointer of the form
> `{cas: <cas-protocol-pCID>, root: <chunk-cid>}`. Every chunk
> is stored under `cas/<cas-protocol-pCID>/<chunk-cid>`. There
> is no size-based exception, no inline-leaf branch, no
> "small messages stay in transports/" carve-out.**

The L6 indirection is uniform. A reader walking
`transports/<wire>/<feed-pCID>/<group-pCID>/<msg-pCID>/<msg-id>.txt`
parses the CBOR pointer, reads the named L6 spec from the
`cas/<cas-protocol-pCID>/` tree (which itself contains the
spec.d directory if the CAS protocol is defined inside the
wire-lab), and resolves the chunk. The walk is the same for
a single-frame message and a 100-MB streamed message; the
pointer-file shape does not change.

### Why no exception

The earlier Alt-M.1 / Alt-M.3 drafts allowed small messages to
remain inline, on the theory that the indirection cost was
wasted for messages that fit in one frame. Steve flagged this
on 2026-05-06 17:54 PT: exceptions make every reader (human or
bot) carry a branch in their head ("is this file a chunk or a
pointer?"), and that branch is exactly the kind of short-
horizon convenience the 100-year goal forbids. Alt-M.4 — the
locked shape — has one rule everywhere.

### Migration of today's transports/wire-lab-devs-draft/

Today's `transports/wire-lab-devs-draft/` directory contains
inline message files written before any L6 CAS spec existed.
Those files are pre-CAS draft state and remain readable as
ordinary files in the repo for now. They will be migrated to
the locked pointer-and-CAS shape when the first L6 CAS spec
lands (anticipated TE-43, "promisebase adoption"). No L6
spec exists in the wire-lab yet to migrate them under; until
one does, the inline files are a known transient.

## Vocabulary table (rename map)

| Old (in TE-29 and earlier) | New (TE-38 and forward) | Notes |
|---|---|---|
| L4 / L4 binding | L5 / feed | "L4" was inherited from OSI counting and never quite matched; layers in this model are L5/L6/L7 by explicit choice, not by analogy to OSI. |
| binding spec | feed spec | The spec sense. |
| `<binding-pCID>` (in path) | `<feed-pCID>` | Path level 2. |
| `udp-binding.d/` | `udp-feed.d/` | Spec directory naming. Existing draft directories keep their old names until their first re-cut; no retroactive rename of locked DI quotations (Cat-1b). |
| forum | group | "forum" remains acceptable in human prose but "group" is the path/spec word. |
| `<forum-pCID>` (in path) | `<group-pCID>` | Path level 3. |
| (no name) | L6 CAS protocol | New level introduced by this TE. |
| (no name) | `cas/` repo subtree | New top-level directory. |

The Cat-1b quotation rule applies: any TE that already locked
a DI using the old vocabulary keeps the old word inside the
quotation. Forward-going text uses the new vocabulary.

## What changes in TE-29

This TE does not supersede TE-29. The architecture TE-29
locked is the architecture this TE inherits. Three navigational
forward-pointers belong on TE-29 (filed as a Cat-3 Refinement
in this TE's commit set):

- The path-level-2 vocabulary is "feed" not "binding"; the
  same level, same role, renamed.
- The path-level-3 vocabulary is "group" not "forum".
- The L6 CAS subtree (`cas/<cas-protocol-pCID>/<chunk-cid>`)
  is the canonical chunk-storage shape; messages in
  `transports/` are CBOR pointers.

TE-29's body text, locked DIs, and recommendation are unchanged.

## Implications and future work

1. **TE-39 through TE-45 inherit this vocabulary.** The
   dependency-sorted TE roster (in TODO 22 and the
   dropped-thread-disposition file) uses "feed" and "group"
   throughout. No TE is blocked on this rename; this TE
   exists so a single-TE reader of any successor TE can
   find the vocabulary lock without reading the conversation
   transcript.

2. **TE-43 (anticipated, "promisebase adoption") owns the
   migration.** The first concrete L6 CAS spec — Rabin
   chunking + Merkle assembly + CIDv1 codec, bundled as
   `promisebase`-equivalent — lands in TE-43 and triggers
   the rewrite of `transports/wire-lab-devs-draft/` into
   pointer-and-CAS form. This TE does not specify the
   rewrite; it locks only the target shape.

3. **Spec directory renames** (`udp-binding.d/` →
   `udp-feed.d/`, etc.) are deferred to the per-feed TEs
   (TE-39 family). The current draft directories keep
   their old names until those TEs re-cut them under the
   new vocabulary; this avoids a sweep that would touch
   files in flight.

4. **Tooling.** Any future tooling that walks the
   `transports/` tree must follow CBOR pointers into `cas/`
   to materialize message bytes. This is recorded as an
   anticipated TODO under TE-43.

## Decision status

`decided`. Sub-decisions:

- DF-38.1 = Alt-1.A (L5 protocols are top-level under
  `protocols/`, not nested) — locked in TODO 22 Q-22.2.
- DF-38.M = Alt-M.4 (all messages are CBOR pointers; no
  exception) — locked in TODO 22 Q-22.6 on 2026-05-07.
- DF-38.G = adopted (TE-28 cited as foundation in the
  "Foundational invariants" section above).

Retracted sub-decisions (logged in TODO 22 with `[~]`):

- DF-38.A — layered-model foundation as a separate question;
  retracted because L5/L6/L7 was already settled in memory.
- DF-38.2 — instance feed-declaration mechanism; retracted
  because TE-29 already locked path-as-declaration.

## Reference to load-bearing constraints

- TE-28 (`docs/thought-experiments/TE-20260501-202713-100-year-goal-as-design-constraint.md`)
  — the 100-year goal that the L5/L6/L7 split answers.
- TE-29 (`docs/thought-experiments/TE-20260501-215027-protocols-as-simulated-repos-and-binding-layer.md`)
  — the architecture this TE renames and extends.
- TE-34 / TE-35 — TE editing policy under which this TE's
  Cat-3 Refinement on TE-29 is filed.
- TODO 22 (`protocols/wire-lab.d/TODO/TODO-20260506-184800-te-38-substrate-agnostic-layered-model.md`)
  — parent TODO with the question log (Q-22.1 through
  Q-22.8).

## Recommendation

Lock the vocabulary and the L6 CAS subtree as written.
Apply the Cat-3 Refinement to TE-29 in the same commit set.
Defer spec-directory renames and the
`transports/wire-lab-devs-draft/` rewrite to the per-feed TEs
and TE-43 respectively.
