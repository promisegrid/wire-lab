# TE-29: Protocols as simulated repos, and the L4-binding layer

## Status

DRAFT. Captures conversation locks from session of 2026-05-01 covering
(a) the directory shape under which each protocol becomes its own
simulated repo, (b) layered transport/session/message decomposition
under `transports/`, and (c) the role of the L4-binding layer. Some
sub-questions are explicitly deferred (see "Open questions" below) and
will be answered in follow-on TEs.

## Why this TE

TE-26 and TE-27 established that "transport" was the right vocabulary
and surfaced axes along which transports differ. TODO 013 carved out
the first concrete transport (then called group-transport) into its
own draft spec. TE-28 named the 100-year goal and the six load-bearing
constraints C-1 through C-6.

Once those landed, three pressures converged:

1. **TE/DR/TODO/DI files at the top level were starting to mix
   harness-level concerns with per-protocol concerns.** A reader
   couldn't tell at a glance whether a TE applied to the wire-lab
   harness itself or to one specific protocol riding on it.

2. **`group-transport` was being treated as if it occupied a single
   layer** (the transport layer), but its locked DFs T1-T6 (Parents
   header, ack-in-body, flat per-leaf layout) are session/group
   semantics, not L4 framing. PromiseGrid is supposed to ride the
   existing Internet, which means most "transports" we care about
   (UDP, TCP, SCTP, WebSocket, MQTT, file-drop, SMTP, BitTorrent,
   IPFS-pubsub) already exist and have RFCs. We don't re-specify them;
   we specify *how PromiseGrid uses them*.

3. **The freeze ceremony for spec docs** (Alt-A through Alt-F in
   conversation) needs a concrete directory shape to attach to before
   the trade-offs become decidable.

This TE locks the directory shape and the layer decomposition. It
explicitly defers the freeze-ceremony decision to a follow-on TE.

## Locked shape: protocols as simulated repos

Top level of the wire-lab outer repo:

```
wire-lab/
├── README.md                            (byte-identical to origin/main)
├── protocols/                           <- one "simulated repo" per protocol pCID
├── transports/                          <- on-the-wire simulation surface
├── tools/                               (spec checker, reference implementations)
└── manifest.json                        (top-level manifest of all frozen releases)
```

Top-level `DR/`, `TODO/`, `DI/`, and `docs/thought-experiments/` go
away as top-level directories. Their content is absorbed into the
relevant protocol's spec doc as inline sections (Decision Records,
Open Work, Don't-Touch Invariants, Bibliography). Cross-protocol
concerns about the wire-lab harness itself live inside
`protocols/wire-lab.d/specs/harness-spec-draft.md`.

### Protocol layout under `protocols/`

Every protocol gets, at minimum, a live working directory:

```
protocols/<slug>.d/
├── docs/thought-experiments/             (per-protocol TEs)
├── specs/<slug>-draft.md                 (the protocol's draft spec)
└── manifest.json                         (per-protocol release manifest)
```

When a protocol release is frozen, two immutable siblings appear next
to the live `.d/`:

```
protocols/<slug>-<pcid>.md                <- frozen spec doc, pCID = hash of this file
protocols/<slug>-<pcid>.d/                <- frozen design tree (sibling, immutable)
```

Multiple frozen releases coexist as siblings (C-4: forking is normal).
The live `<slug>.d/` continues to evolve toward the next freeze.

### Concrete protocol inventory at the time of this TE

```
protocols/
├── wire-lab.d/                          (the harness itself, treated as a protocol)
├── wire-lab-bafkreih...A1.md            (frozen, when first release ships)
├── wire-lab-bafkreih...A1.d/
│
├── udp-binding.d/                       (L4 binding spec for UDP)
├── tcp-binding.d/                       (L4 binding spec for TCP)
├── websocket-binding.d/                 (L4 binding spec for WebSocket)
├── mqtt-binding.d/                      (L4 binding spec for MQTT)
├── file-drop-binding.d/                 (filesystem-as-transport binding)
├── smtp-binding.d/                      (email-as-transport binding)
│
├── group-session.d/                     (session/group protocol; was "group-transport")
│
└── ppx-dr.d/                            (message protocol: proposals + contests)
```

`group-transport` is renamed to `group-session` because its DFs T1-T6
are session-layer semantics, not L4 framing (see "Layer
decomposition" below). The locked DFs themselves do not change. The
rename is tracked as a follow-on TODO.

`ppx-dr` is a new explicit protocol slug for the proposal/contest
schema that has so far lived as ad-hoc files under `proposals/`. The
name `ppx-dr` matches the existing DR-001 bootstrap convention.

## Layer decomposition under `transports/`

`transports/` is the simulation surface where bytes-as-bytes live.
Every PromiseGrid message file in the simulation lives at a path of
the form:

```
transports/<real-world-transport-name>/<L4-binding-pCID>/<session-pCID>/<message-pCID>/<message-id>.msg
```

Five-level decomposition, each level corresponding to one layer's
identity:

| Level | Name                          | Identity      | Example                                        |
|-------|-------------------------------|---------------|------------------------------------------------|
| 1     | real-world transport          | external slug | `udp`, `tcp`, `websocket`, `mqtt`, `file-drop` |
| 2     | L4-binding pCID               | PromiseGrid   | `udp-binding-bafkrei...U1`                     |
| 3     | session-protocol pCID         | PromiseGrid   | `group-session-bafkreig...B2`                  |
| 4     | message-protocol pCID         | PromiseGrid   | `ppx-dr-bafkrei...D4`                          |
| 5     | message-id (filename)         | content-hash  | `<message-id>.msg`                             |

### Why level 1 is a slug, not a pCID

Level 1 names a real-world transport that already exists and that
PromiseGrid does not get to redefine. UDP is what the IETF assigned
that name to in RFC 768 in 1980; we don't replace that name, we adopt
it. C-1 (no central registry) is satisfied because real-world
transport names are external givens we don't own. The pCID layer
appears at level 2, where we name *our* binding spec for that
transport.

### Why levels 2-4 are separate pCIDs

Each level can be forked independently (C-4):

- Two parties can author competing UDP bindings without forking the
  session protocol above them.
- Two session protocols can run over the same UDP binding.
- Two message protocols (e.g. `ppx-dr` and a future `chat`) can ride
  the same session protocol.

If we collapsed levels into one spec, every change at any layer would
require forking everything below. Per-layer pCIDs make composition
the cheap operation and bundling the explicit, opt-in operation.

### A "stack" is a tuple, not a single pCID

A receiver can only interpret a message if it knows the full
`(level-2, level-3, level-4)` tuple. Two parties saying "I speak
ppx-dr" without naming the session and binding have not actually
agreed on anything they can send each other. Negotiation/handshake
(out of scope for this TE) is the act of agreeing on the tuple.

### Layer count is open

Five levels is the maximum demonstrated by current protocols, not a
floor or a ceiling. Some leaves stop earlier (a binding-only
adversarial probe with no session above it). Future protocols may
need a sixth layer (e.g., explicit framing between binding and
session). The shape generalizes by adding directory levels; we don't
pre-declare layers we don't yet need.

## What the binding layer does

The binding layer owns the round-trip translation between **a
sequence of opaque PromiseGrid messages** (what session and
message-protocol layers produce and consume) and **whatever shape the
underlying real-world transport actually moves** (datagrams, byte
streams, mailbox files, MQTT publishes, etc.).

In scope for a binding spec:

1. **Framing.** How are message boundaries found?
2. **Size limits.** Largest message the binding promises to carry intact.
3. **Fragmentation/reassembly,** if any.
4. **Addressing convention** specific to this transport.
5. **Connection model.** Stateless? Long-lived with reconnect?
6. **Liveness and failure semantics.** Loss, ordering, duplication promises.
7. **Port/path/topic conventions,** including any default and the rule
   that no central allocator exists.
8. **Receive loop responsibilities** for a conformant implementation.
9. **Local I/O details** for the simulation: what bytes go in
   `transports/<wire>/<binding-pCID>/...` files.

Out of scope for a binding spec:

- What messages mean (message-protocol layer's job).
- Who is in the conversation, ordering, ack semantics (session-protocol layer's job).
- Routing across multiple transports (higher-layer concern).
- Cryptographic signing of message *content* (session or message
  layer; the binding may sign per-frame for spam resistance, which is
  a separate use of crypto).

A binding spec is not a re-specification of UDP/TCP/MQTT; those
specs already exist. The binding spec is the much smaller document
that says: given that UDP exists and works as RFC 768 says, here is
exactly how PromiseGrid uses it. UDP-binding v0 is roughly 5-10
pages. TCP-binding more (framing is harder). SMTP-binding more still
(impedance mismatch is bigger).

A reference one-page sketch for UDP-binding v0 is committed alongside
this TE at `protocols/udp-binding.d/specs/udp-binding-draft.md`.

## Send/receive walk-through (canonical example)

Stack: `udp` / `udp-binding-bafkrei...U1` /
`group-session-bafkreig...B2` / `ppx-dr-bafkrei...D4`. Sender wants
to deliver a ppx-dr proposal.

**Send path:**

1. Application produces a `Proposal{...}` struct.
2. ppx-dr layer CBOR-encodes to PAYLOAD bytes (~480 bytes).
3. group-session layer wraps PAYLOAD in a session message with
   Parents/From/Message-ID/Content-Protocol headers, blank line,
   payload. Total ~612 bytes.
4. UDP-binding layer:
   - Checks 612 <= 1232 (UDP-binding v0 size limit). OK.
   - Looks up peer address `198.51.100.7:4646`.
   - `udp_socket.sendto(SESSION_MESSAGE, addr)`.
   - Writes simulation artifact:
     `transports/udp/udp-binding-bafkrei...U1/group-session-bafkreig...B2/ppx-dr-bafkrei...D4/<message-id>.msg`
     with file content equal to the exact 612 bytes.
5. OS prepends UDP+IP headers; bytes traverse the network. Out of
   scope for any PromiseGrid spec.

**Receive path:**

1. UDP-binding recv loop: `recvfrom` returns 612 bytes from peer.
   Write file at the same path on the receive side. Hand 612 bytes up.
2. group-session parses headers, splits at blank line, reads
   `Content-Protocol: ppx-dr-...`, hands payload up.
3. ppx-dr CBOR-decodes to `Proposal{...}`.
4. Application receives the proposal.

**What changes if the binding swaps:**

- TCP-binding: varint length prefix + bytes; long-lived connection;
  no inherent size limit.
- WebSocket-binding: one binary frame per message.
- MQTT-binding: PUBLISH on topic `pgrid/<group-id>` at QoS configured
  by the binding spec.
- File-drop-binding: write `<message-id>.msg` to a shared spool dir;
  receiver watches with inotify.
- SMTP-binding: base64-wrap as email body to a configured address;
  receiver polls IMAP.

What stays identical across all bindings: the session message bytes
themselves. The session and message-protocol layers do not know or
care which binding ran. That is exactly the property the binding
layer preserves.

## Migrations triggered by this TE (follow-on TODOs)

This TE locks shape but does not move files. The following migrations
are tracked as separate TODOs:

1. **TODO 014: protocols-as-simulated-repos migration.** Move
   `specs/harness-spec-draft.md` to
   `protocols/wire-lab.d/specs/harness-spec-draft.md`. Move
   `specs/group-transport-draft.md` (renamed) to
   `protocols/group-session.d/specs/group-session-draft.md`. Move
   `specs/transport-spec-draft.md` content into the binding layer
   (likely deleted; the thin outer rule it captured is now subsumed
   by per-binding specs). Move all top-level `docs/thought-experiments/`
   TEs into the appropriate per-protocol `.d/` (most go under
   `wire-lab.d/`; the group-transport-envelope/transport-axes TEs go
   under `group-session.d/`).

2. **TODO 015: DR/TODO/DI absorption.** Inline current `DR/`,
   `TODO/`, `DI/` content as sections inside the relevant protocol
   spec doc. Per-protocol DRs/TODOs go to that protocol's spec;
   harness-level go to wire-lab harness spec.

3. **TODO 016: proposals as transport messages.** Move
   `proposals/pending/ppx-dr-001-bootstrap/` content to
   `transports/<future-bootstrap-transport>/<udp-binding-pCID>/<group-session-pCID>/ppx-dr-bafkrei.../<message-id>.msg`.
   DI-003 protection follows the bytes, not the path. Renaming
   `.md` to `.msg` is itself a one-time edit that must precede the
   DI-003 anchor moving.

4. **TODO 017: group-transport -> group-session rename.** Pure
   renaming. DFs T1-T6 unchanged.

5. **TODO 018: write UDP-binding v0 reference and test vectors.**
   Once the directory is migrated, flesh out the one-page sketch
   into a v0 spec with test vectors and a Go reference implementation
   stub.

## Open questions

OQ-29.1: **Freeze ceremony.** Two-step commit (Alt-B), Merkle bundle
(Alt-C), release record (Alt-D), git tag (Alt-E), or hybrid (Alt-F)?
Deferred to a follow-on TE after the user answers two clarifying
questions about (a) self-sufficiency of the spec doc taken alone and
(b) doc-centric vs repo-centric simulation.

OQ-29.2: **Filename inside transport leaves.** Content-hash, ULID, or
both (`<ulid>-<hash>.msg`)? Hash gives free cross-transport
deduplication; ULID makes ordering trivial. Lean: hash; ordering
recovered from session-layer Parents header.

OQ-29.3: **Empty-leaf rule.** Is `transports/<wire>/<binding>/`
allowed with no session subdirectory yet? Lean: yes (honest about
provisioning state).

OQ-29.4: **Framing intermediate layer.** Do we need a separate
framing layer between binding and session, or does each binding own
its framing? Lean: each binding owns its framing (no reuse to be had
across UDP/TCP/MQTT/SMTP).

OQ-29.5: **Binding-layer signing.** May a binding sign per-frame for
spam resistance independently of session/message-layer signing? Lean:
yes, but out of scope for v0 specs.

OQ-29.6: **Negotiation/handshake protocol.** Where does
"agree on the (binding, session, message) tuple" live? Probably its
own protocol on a bootstrap transport. Future work.

OQ-29.7: **Real-world transport slug list.** Who decides which slugs
appear at level 1? Lean: any author can introduce a new
`transports/<slug>/` directory; collisions resolve by C-1 (no central
registry). Practical norm: use IANA names where they exist.

OQ-29.8: **Test-harness simulation of network effects.** Does the
wire-lab simulator also emulate loss, reordering, latency, partition?
If yes, where does the emulation live in the directory shape? Lean:
adversarial harness lives in `tools/` and produces test vectors that
violate the binding's promises in controlled ways. ns-3 is the
leading candidate (see OQ-29.9 for empirical sandbox results).

OQ-29.9: **External network simulator integration (ns-3, similar).**
Once at least two binding specs and one session protocol have v0
implementations in `tools/`, can the wire-lab harness drive ns-3 or a
similar packet-level simulator to produce realistic loss/latency
adversarial scenarios?

**Empirical sandbox results (added 2026-05-01, supersedes prior
"start with tc netem" lean):** A separate session in the Perplexity
Computer sandbox (Debian trixie, kernel 6.1.158, sudo available)
tested every candidate from this TE. Findings:

| Candidate          | Status in sandbox          | Notes                                                                       |
|--------------------|----------------------------|-----------------------------------------------------------------------------|
| ns-3 3.43          | **Working**                | Debian package; CMake+Ninja+pkg-config; UDP echo + PCAP smoke-tested.       |
| `tc netem`         | **Unavailable**            | Kernel lacks `sch_netem`, `tbf`, `fq_codel` modules; only `pfifo_fast`.     |
| Mininet 2.3.0      | Working with caveats       | Requires `PYTHONPATH=/usr/lib/python3/dist-packages`; LinuxBridge only.     |
| netns + veth + br  | Working                    | sudo required; clean topology setup; no impairment available without netem. |
| OMNeT++            | Not installed              | No Debian package found; would require source build.                        |
| Shadow             | Not installed              | No usable Debian package; would require source build.                       |

The practical consequence is that **packet-loss / latency / jitter
adversarial scenarios run inside this sandbox have to live in ns-3**,
not in `tc netem`-based emulation. Mininet remains useful for live
topology experiments without impairment. Steve's local dev box (any
normal Linux install) has `tc netem` and can use it; the sandbox does
not. This is a sandbox-specific constraint, not a project-level
recommendation reversal.

Revised lean: **adopt ns-3 as the primary network-realism fixture**
for adversarial scenarios in the wire-lab, with a small Mininet path
for topology-only experiments. The integration shape (ns-3 emulates
the wire; Go reference implementations run as applications above
sockets via tap-bridge or DCE; PCAP output is post-processed into
`transports/<wire>/<binding-pCID>/.../<msg-id>.msg` artifacts) is
unchanged from the original lean. Worth a dedicated TE once the
first Go binding+session implementations exist. Reference smoke-test
artifacts from the sandbox session are recorded outside this repo;
the wire-lab's own ns-3 scenarios will live under
`tools/ns3-harness/` once authored.

## Reference to load-bearing constraints

This TE relies on (TE-28):

- **C-1 no central registry:** Real-world transport slugs are external
  givens, not allocations. Default ports are conventions, not
  allocations.
- **C-3 adversarial-by-default:** Per-layer pCIDs allow adversarial
  variation at any single layer in isolation.
- **C-4 forking is normal:** Per-layer pCIDs make forking the
  cheap-by-default operation; bundling is the explicit one.
- **C-6 signing key is the only structural lock:** Binding-layer
  signing (OQ-29.5) is permitted but does not subsume session/message
  signing.

## Recommendation

Adopt the locked shape and the layer decomposition. Track the five
follow-on TODOs (014-018). Defer the freeze ceremony to its own TE.
Treat ns-3 (or equivalent) as a future external test fixture once
protocols have v0 implementations.
