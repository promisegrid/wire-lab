# Historical networks: substrate precedent for federated message protocols

**Date:** 2026-05-03
**Status:** Research note. Source material for TE-sihih (binding /
substrate / feed protocols as a first-class family). The naming
question itself ("binding" vs. "substrate" vs. "feed" vs. other) is
unresolved as of this writing; this document uses "substrate" as the
neutral working term for the carrier mechanism (git, rsync, NNTP,
UUCP, ...) and notes where specific precedents use specific words.

Current reading after TE-sihih: forward-going wire-lab vocabulary is
`feed` for the L5 protocol role and `substrate` for explanatory prose.
The old `binding` wording remains historical precedent vocabulary, not
the active wire-lab term. Source: `DI-pijun`; TE-sihih.

## Question being grounded

Wire-lab's group-session transport currently spells out a git-shaped
delivery cycle in its §9 (per-author-branch, fetch-merge-push). Steve's
turn-172 prompt enumerated rsync, unison, UUCP, UDP, SVN, CVS as peer
substrates -- forcing the question of whether the same group-session
protocol can run over multiple delivery substrates simultaneously, and
if so, whether the substrate axis deserves its own first-class
abstraction parallel to the existing transport-protocol axis.

The question this document grounds: **is there practical precedent in
RFCs and historical networks for separating an abstract federated-
message-protocol from its concrete delivery substrate?**

Short answer: yes, deep and unambiguous, going back to the late 1970s.

## Six concrete precedent systems

### 1. Email (SMTP / UUCP / X.400) -- the canonical example

[RFC 822](https://www.rfc-editor.org/rfc/rfc822) (1982) and its
successor [RFC 5322](https://www.rfc-editor.org/rfc/rfc5322) (2008)
define the *message format*: headers (From, To, Subject, Date,
Message-ID), addressing, body conventions. The *delivery substrate*
has been, in production, all of:

- SMTP over TCP ([RFC 5321](https://www.rfc-editor.org/rfc/rfc5321))
- UUCP over dialup modems, with bang-path addressing
  (`host1!host2!user`) -- same RFC 822 messages, completely different
  substrate
- BSMTP (batched SMTP) wrapped inside UUCP transfers -- SMTP-the-
  protocol carried over UUCP-the-substrate
- X.400 over OSI stacks, gatewayed bidirectionally to RFC 822 mail
  via [RFC 1327](https://www.rfc-editor.org/rfc/rfc1327) and
  [RFC 2156/2157](https://www.rfc-editor.org/rfc/rfc2156)
- More recently: SMTP over TLS, SMTP over Tor hidden services, SMTP
  over I2P

Same message identity (Message-ID, headers, body), same routing
semantics, totally different wires. Sendmail's configuration language
has explicit *relay* declarations per substrate (`UUCP_RELAY`,
`SMART_HOST`, etc.) -- the "instance declares its substrates" pattern
in production form.

### 2. Usenet (NNTP / UUCP) -- the closest structural match

Usenet is structurally near-identical to wire-lab group-session: a
federated, store-and-forward, append-only message space with content-
addressed-ish identity (Message-ID) and DAG-shaped reply threading
(References:). It ran over multiple substrates from day one:

- Originally pure UUCP store-and-forward over modems
- [NNTP over TCP](https://www.rfc-editor.org/rfc/rfc977) added in the
  mid-1980s; [RFC 3977](https://www.rfc-editor.org/rfc/rfc3977) is the
  current revision
- For decades, both ran simultaneously, with sites using UUCP for some
  peers and NNTP for others
- The articles themselves were byte-identical regardless of which
  substrate carried them -- same Message-ID, same headers, same body
- A site running NNTP could receive an article via UUCP from one peer
  and serve it via NNTP to another, with no change to the article

Each Usenet site's `newsfeeds` file declares which peers it talks to
over which substrate (UUCP / NNTP / ...). This is the per-instance
substrate-declaration pattern.

Usenet does NOT use the word "binding" for these substrate mappings.
Operational vocabulary used "feed," "newsfeed," "peer," "transport."
This is one source of the naming-question tension: W3C-style formal
vocabulary (binding) vs. store-and-forward operational vocabulary
(feed).

The same precedent also supports a broader PromiseGrid design reading: current
`group-session` can be read as one early specimen on a
content-addressed-Usenet successor line, especially when paired with git's own
object/transport separation as a modern substrate precedent. This remains an
exploratory framing rather than a frozen protocol identity claim; the point of
the analogy is to surface which structural invariants aged well, not to assert
that PromiseGrid must import every Usenet mechanism unchanged. Source:
`DI-pijun`.

### 3. FidoNet -- same pattern, different substrate set

FidoNet messages had a fixed format
([FTS-0001](http://ftsc.org/docs/fts-0001.016)) and were carried over
dialup modems originally, then later over TCP via binkp
([FTS-1026](http://ftsc.org/docs/fts-1026.001)) and ifcico, and
gatewayed to/from Usenet and email. Same separation: message format
is one spec, substrate transports are separate specs.

FidoNet vocabulary used "transport," "mailer," and "tosser" for
different layers; again, no "binding."

### 4. CORBA GIOP / IIOP -- the formal "protocol with substrates" model

GIOP (General Inter-ORB Protocol) is the abstract message format;
IIOP is GIOP-over-TCP; bindings also existed for HTTP, raw IPC, and
shared memory. The
[OMG CORBA specs](https://www.omg.org/spec/CORBA/) explicitly call
these "transports" with GIOP as the substrate-agnostic core.

CORBA's word: "transport." (Conflicts with wire-lab's existing
"transport protocol" usage, which is one reason the wire-lab vocabulary
is contested.)

### 5. SOAP / WSDL -- explicit "binding" terminology

[WSDL 1.1](https://www.w3.org/TR/wsdl) and
[WSDL 2.0](https://www.w3.org/TR/wsdl20/) have a formal `<binding>`
element that maps an abstract operation onto a concrete wire (SOAP-
over-HTTP, SOAP-over-SMTP, SOAP-over-JMS, SOAP-over-MQ). The W3C uses
the word *binding* for exactly the substrate-mapping concept, and the
artifact is a separate document from the protocol's abstract
definition.

This is the closest formal precedent for the spec-structure we are
considering: an abstract protocol document plus separate per-substrate
binding documents.

W3C vocabulary: "binding."

### 6. Modern systems with the same pattern

- **gRPC** -- protobuf service definitions are substrate-agnostic;
  gRPC-over-HTTP/2 is the default but
  [gRPC-Web](https://github.com/grpc/grpc-web), gRPC-over-QUIC, and
  in-process bindings exist as peers
- **NATS / Kafka / AMQP** -- message-format specs and transport-binding
  specs are separate documents
  ([AMQP 1.0 spec](https://www.amqp.org/sites/amqp.org/files/amqp.pdf))
- **libp2p** -- explicitly designed around "transports" (TCP, QUIC,
  WebRTC, WebSockets, Bluetooth) as pluggable peers under a substrate-
  agnostic protocol stack
  ([libp2p protocols](https://libp2p.io/docs/protocols/))
- **Matrix** -- the Matrix federation protocol is HTTPS-bound today
  but the [spec](https://spec.matrix.org/) is careful to call out
  that the message format and event graph are independent of HTTP
- **Git itself** -- git distinguishes the *object format* (content-
  addressed objects: blobs, trees, commits, tags) from the *transport*
  (`git://`, `ssh://`, `https://`, `file://`, dumb-HTTP, smart-HTTP,
  bundle files). Same objects, same identity (SHA), multiple
  substrates. This is the most relevant precedent of all because it
  is literally the substrate wire-lab is using today, and wire-lab's
  CID layer stacks on top of git's blob-SHA layer when the substrate
  is git.

Modern vocabulary varies: "transport" (libp2p, gRPC), "wire format"
(protobuf), "broker" (Kafka, AMQP). Again, no consensus on a single
word.

## Three recurring patterns

Across all six systems, three patterns recur:

### Pattern A -- Substrate specs as separate documents, parallel to the protocol

- W3C: WSDL `<binding>` documents
- IETF: separate RFCs for "X over Y" -- e.g.
  [RFC 1149 IP-over-avian-carriers](https://www.rfc-editor.org/rfc/rfc1149)
  is the joke version;
  [RFC 3261 SIP-over-TCP/UDP/SCTP/TLS](https://www.rfc-editor.org/rfc/rfc3261)
  is the serious one
- OMG: GIOP spec + IIOP spec as separate documents

### Pattern B -- Substrates declared per-instance / per-deployment, not per-protocol

- Sendmail: each site's `sendmail.cf` declares which relays/substrates
  it uses for which destinations
- Usenet: each site's `newsfeeds` file declares which peers it talks
  to over which substrate (UUCP / NNTP / ...)
- gRPC: each service deployment declares its listening transports

### Pattern C -- Substrate-agnostic message identity

- Email Message-ID (RFC 822 §4.6.1)
- Usenet Message-ID (RFC 3977 §3.6)
- Git SHA (content-addressed object hash)
- All computed from message content; *invariant* under substrate
  change. A message with the same ID is the same message regardless
  of which wire carried it.

## How the patterns map onto wire-lab

If wire-lab adopts the same architecture:

- Pattern A -> `protocols/<substrate-slug>-<word>.d/` -- separate spec
  documents per substrate (where `<word>` is the as-yet-unsettled term:
  `binding`, `feed`, `substrate`, ...)
- Pattern B -> a per-instance declaration of which substrates this
  instance is mirrored over (a `bindings/` subdirectory, or
  equivalent, with one entry per active substrate; see TE-sihih DF-38.2)
- Pattern C -> message CIDs are already substrate-agnostic by
  construction. The canonical bytes (group-session §3, §4) do not
  include any substrate-specific addressing -- so the CID is invariant
  whether the file arrived over git, rsync, UUCP, etc.

Current reading after TE-sihih: Pattern A maps to L5 feed specs, using
`feed` rather than `binding`. Pattern B does not require the abandoned
turn-172/174 `bindings/` directory or `INSTANCE.md` feed-manifest
sketch; TODO-vunub Q-22.3 retracts that mechanism because TE-vipir's
path-as-declaration already identifies the wire/feed/group/message stack.
Pattern C remains unchanged. Source: `DI-pijun`; TODO-vunub Q-22.3;
TE-sihih.

## Counterfactual: systems that did NOT separate protocol from substrate

The systems that conflated protocol with substrate from the start
tend to acquire painful retrofits when a second substrate becomes
necessary:

- **HTTP** was designed as a single TCP-based protocol; HTTP/2,
  HTTP/3, HTTP-over-QUIC, WebSockets, and Server-Sent Events are
  essentially "HTTP-the-message-format with new substrates under it"
  decades after the fact. Each retrofit has been disruptive.
- **IRC** ([RFC 1459](https://www.rfc-editor.org/rfc/rfc1459) /
  [RFC 2810](https://www.rfc-editor.org/rfc/rfc2810)) chose monolithic
  TCP and is criticized for it; federation is fragile and
  substrate-locked.

The implicit argument: a system designed for multiple substrates from
day one is structurally cheaper than one that retrofits additional
substrates later under load.

## Negative-precedent check

Turn-173 recovery revisited the obvious counter-evidence candidates.
The result is mixed cautionary precedent, not a direct disproof of the
substrate-as-first-class shape. Source: `DI-pijun`.

- **XMPP** is not a negative case. [RFC 6120](https://www.rfc-editor.org/rfc/rfc6120)
  defines the TCP binding, [BOSH](https://xmpp.org/extensions/xep-0124.html)
  and [RFC 7395](https://www.rfc-editor.org/rfc/rfc7395) add HTTP/WebSocket
  options, and the protocol did not abandon the idea that one message protocol
  can have multiple connection methods.
- **JMS / Jakarta Messaging** is a caution, not a match. The
  [Jakarta Messaging specification](https://jakarta.ee/specifications/messaging/3.1/jakarta-messaging-spec-3.1.html)
  standardizes an API and explicitly does not define a wire protocol, so it
  shows the portability limits of an API-only abstraction across providers
  rather than failure of a federated wire protocol with explicit substrate
  specs.
- **BITNET, early P2P systems, and other store-and-forward systems**
  remain useful cautionary comparisons, but this note does not record a
  known system that cleanly tried the wire-lab target shape -- stable
  message identity plus first-class substrate specs plus per-site
  substrate declarations -- and then abandoned that shape as unworkable.

The conclusion is deliberately bounded: no convincing negative precedent
is recorded in this note, but future TEs should still treat the analogy as
evidence to test, not as proof that every Usenet, XMPP, or JMS mechanism
belongs in PromiseGrid. Source: `DI-pijun`.

## Naming question (DF-38.5 input)

The precedent enumeration above does NOT settle the naming question.
Each precedent system has its own term; the wire-lab choice will
depend on which precedent set is weighted most heavily:

| Precedent set | Word | Notes |
|---|---|---|
| W3C / SOAP / WSDL | binding | Closest formal precedent for spec structure |
| CORBA / libp2p / gRPC | transport | Conflicts with existing wire-lab "transport protocol" usage |
| Usenet / FidoNet operational | feed | Store-and-forward heritage; per-site `newsfeeds` file |
| Generic / informal | substrate | Architecture-neutral; describes the carrier mechanism |
| PromiseGrid existing | grid envelope | Already reserved for the `grid <pcid>` selector/envelope, not for feed naming |

Bot's turn-173 recommendation was "binding" (W3C-precedent-weighted).
Steve's turn-174 challenge will be that Usenet-gateway operational
vocabulary did not use "binding"; turn 175 will endorse "feed" and
"substrate" jointly. The disagreement is about precedent-weight, not
about either set being wrong. (Logged as UT-173.b in the session-
replay ledger.)

This document deliberately uses "substrate" throughout as the neutral
working term, and notes the contested naming explicitly so a future
reader can substitute the eventually-chosen word without losing the
substantive content.

## Source

This research was assembled in turn 173 of the 2026-05-03 wire-lab
session. The conversation also surfaced a non-obvious implication:
under git, wire-lab's CID layer stacks on top of git's blob-SHA layer,
giving two CAS layers; under a non-CAS substrate (rsync, UUCP), only
the wire-lab CID layer remains. The substrate-binding spec for each
substrate must account for the cardinality of CAS layers. (Logged as
UT-173.e; connects to turn 175's "We're not really representing the
decentralized CAS" correction.)

## Citing this doc

When TE-sihih is drafted, its background section can cite:

> See `docs/research/historical-networks-20260503.md` for the
> precedent enumeration grounding the feed-protocol-as-first-class
> architectural claim.

This keeps the precedent material in one durable location, separate
from any specific TE's argument structure, so it can be re-cited from
later TEs (substrate-naming, CAS-cardinality, gateway design) without
duplication.
