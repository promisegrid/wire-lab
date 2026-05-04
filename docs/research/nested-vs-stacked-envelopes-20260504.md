# Nested vs Stacked Envelopes: How Long-Lived Decentralized Protocols Resolved This Design Question

**File:** `wire-lab/docs/research/nested-vs-stacked-envelopes-20260504.md`  
**Date:** 2026-05-04  
**Context:** PromiseGrid wire-lab — candidate envelope shape research  
**Status:** Reference material for design decision

---

## 1. Executive Synthesis

The record across twenty-five years of long-lived decentralized protocol work is unambiguous: **recursive nesting (each layer is itself a complete, self-describing message containing another message) has aged dramatically better than ordered frame sequences in the category of multi-hop, content-addressed, capability-bearing envelopes.** The frame-sequence model ages well only in narrower, session-bound, crypto-handshake contexts where a fixed state machine controls both ends.

Three concrete anchors illustrate why:

**First, IPFS/IPLD's CID+codec model** — the design that most directly parallels `grid([pcid, payload])` — has been one of the most durable parts of that ecosystem. A CIDv1 is literally `[version | codec-id | multihash]`: a self-describing recursive identifier where each block carries the selector for how to interpret it. The lesson from IPLD's evolution is that when they *baked* codec semantics *into* a single monolithic codec (dag-pb), the result was brittle and rigid; when they separated the codec selector (multicodec) from the hash and from the payload, the design became composable and long-lived. See the [IPLD post-mortem on dag-pb](https://ipld.io/docs/synthesis/how-ipfs-web-gateways-work/) for the explicit retrospective.

**Second, JOSE/JWE nested JWTs** represent the cautionary tale about nesting done *wrong*. JOSE supports nesting syntactically (a JWT whose payload is another JWT, signaled by `cty: JWT`), but the nesting is implicit, depth-unbounded, and relies on recursive re-entry in the validation path. This has produced: algorithm confusion attacks, DoS vulnerabilities via unbounded recursion (CVE-2025-53864 in Nimbus JOSE+JWT), and signing-order ambiguity. The failure was not nesting per se — it was nesting without a codec/type selector attached to the outer envelope before you begin unwrapping. PASETO's explicit rejection of JOSE's model, and COSE's redesign of JOSE for CBOR (with protected headers as bstr-wrapped byte strings rather than inline JSON), both confirm the same diagnosis.

**Third, Nostr's deliberately flat event format** shows that when the design space is *single-hop, single-semantic* (post, sign, broadcast), flatness wins on implementation simplicity. But when Nostr needs layered semantics (encrypted DMs via NIP-17, gift wrapping via NIP-59), it reinvents a manual nesting structure anyway — confirming that nesting is the right primitive for multi-semantic contexts and flat is a special case, not a general principle.

**For a 100-year-lived envelope, the correct answer is: recursive nesting with an explicit, content-addressed codec/handler identifier attached at every layer, such that any layer can be peeled, identified, and validated independently.** The `grid([pcid, payload])` hypothesis aligns with what has aged well. The "promise stack" as an ordered-frame protocol is a reasonable *transport* primitive but is the wrong abstraction at the *message-semantic* layer of a content-addressed, capability-bearing system.

The key insight from studying all these systems: **what makes nesting brittle is not recursion itself but the absence of a self-describing type token at every layer.** JOSE failed because it made type tokens optional and implicit. IPLD succeeded because CIDs carry their codec identifier unconditionally. PromiseGrid's `pcid` — a content-addressed identifier of which protocol/handler interprets the payload — is precisely that mandatory type token. This is the load-bearing design.

---

## 2. The Two Models, Defined

### Model A: Recursive Nesting (Onion)

A message is a complete, self-describing unit that optionally contains another complete message as its payload. The structure is:

```
Message := [type-selector, payload]
payload  := bytes | Message   // recursion
```

Each layer is:
- **Self-describing**: contains its own type/codec identifier
- **Independently valid**: can be parsed, verified, and interpreted without context from enclosing layers
- **Peelable by semantics**: the receiver uses the type-selector to dispatch to the right handler, which may produce a new inner Message

Key invariant: **you know what you have before you open it.** The type-selector (codec ID, MIME type, CID, etc.) is unconditionally present at every layer.

Examples in this report: CID+block in IPLD, COSE_Sign1/COSE_Encrypt0, DIDComm encrypted envelope, Cap'n Proto capability references, CapTP object references.

### Model B: Ordered Frame Sequence (Stack)

A message is a flat list of typed frames, processed in order. Each frame peels one semantic layer:

```
Message := [Frame_0, Frame_1, Frame_2, ...]
Frame_i  := [frame-type, frame-data]
```

Processing is a state machine: `Frame_0` is processed first and modifies state that affects how `Frame_1` is interpreted. The sequence is not a tree; it is a tape.

Key invariant: **the meaning of a frame depends on what precedes it in the sequence.** Frames are not independently meaningful.

Examples: Noise Protocol Framework handshake message tokens (`e`, `ee`, `s`, `es` processed sequentially), TLS record layer, QUIC frames.

### The Critical Difference

| Property | Recursive Nesting | Ordered Frame Sequence |
|---|---|---|
| Each layer self-describing | Yes — type at every level | No — meaning depends on position |
| Can peel any single layer independently | Yes | No — must process from start |
| Supports content-addressing per layer | Yes — hash each complete nested message | No — frame i is not addressable independently |
| Suitable for multi-hop routing | Yes — each hop sees a complete message | Only if hop count is fixed at design time |
| Suitable for promise pipelining | Yes — any layer can carry a capability reference | Awkward — references must be flattened |
| Scales to unbounded depth | Yes, with depth guards | Yes, trivially |
| Implementation complexity | Moderate (recursive dispatch) | Low (linear scan) |

---

## 3. Per-System Findings

### 3.1 E Language and CapTP

**Choice made:** Neither pure nesting nor pure sequence — CapTP uses a **table-indexed reference model** with flat messages referencing object IDs from the "Four Tables" (Imports, Exports, Questions, Answers). A message on the wire is roughly: `[message-type, target-reference-id, args...]`. Arguments can themselves be capability references (table indices). This is closer to a **flat record with typed slots** than either deep nesting or a sequence of frames.

**Primary source:** [E Language CapTP Four Tables](http://erights.org/elib/distrib/captp/4tables.html) — Mark S. Miller's original design document. The four-table scheme (Imports/Exports for object references, Questions/Answers for promise resolution) gives each capability an integer ID over a connection. Nested object graphs are represented as multiple messages, not as a deeply nested single message.

**Rejected alternative:** Deep arbitrary nesting of capability references. CapTP explicitly avoids this because it creates distributed object graphs with no clear lifecycle; the table approach gives both sides explicit control over reference counting and garbage collection.

**Cap'n Proto** (Kenton Varda's implementation of CapTP concepts) documents that messages use **offset-based pointers within a flat arena**, where *outer objects appear before inner objects* in the wire encoding — the inverse of conventional nested serialization. Capability references are table IDs passed as typed fields. See [Cap'n Proto RPC](https://capnproto.org/rpc.html) and [Cap'n Proto Introduction](https://capnproto.org) which notes: *"outer objects appear entirely before inner objects (as opposed to most encodings, where outer objects encompass inner objects)"* — an explicit design choice for incremental reads without full buffer arrival.

**Post-hoc assessment:** Cap'n Proto's design has held up well. Kenton Varda has noted in multiple venues that the promise-pipelining design (sending `bar(foo())` as two messages in a single round trip, without waiting for `foo()` to resolve) is the most important architectural feature, and it works precisely because capability references are flat table IDs, not nested objects. There is no formal "this aged badly" retrospective; the system is actively maintained as of 2025.

**Spritely Goblins / OCapN:** As documented in the [Spritely Core whitepaper](https://spritely.institute/static/papers/spritely-core.html) and [OCapN GitHub repository](https://github.com/ocapn/ocapn), OCapN is still pre-standardization (as of mid-2026), representing a convergence of Agoric/Spritely/Cap'n Proto implementations. The CapTP abstract protocol uses an abstract "netlayer" interface and carries object references as integer IDs. The [Spritely Goblins CapTP documentation](https://files.spritely.institute/docs/guile-goblins/0.13.0/CapTP-The-Capability-Transport-Protocol.html) emphasizes that the protocol eliminates the need for explicit message-structure coordination in most applications. The serialization format for OCapN messages is expected to be Syrup (a simple, canonicalizable encoding), but the spec is not finalized.

**Implication for PromiseGrid:** CapTP's approach suggests that `pcid` should be a content-addressed reference (like a CID) rather than a nested capability, and that capability references in `payload` should be table-IDs or CID-like references, not deeply nested capability objects. The analogy between CapTP's exports table and IPFS's content-address space is direct.

---

### 3.2 Tahoe-LAFS

**Choice made:** **Flat capability strings with hierarchical internal structure.** A Tahoe capability (e.g., `URI:CHK:key:hash:k:n:size`) is a plain string URI where each colon-separated field has a specific type role. The capability is not a nested message; it is a typed, structured identifier.

**Primary source:** [Tahoe-LAFS URI Specification](https://tahoe-lafs.readthedocs.io/en/latest/specifications/uri.html). The architecture layers capabilities through a strict abstraction stack:
- Key-value store: capability → bytes (no metadata)
- Filestore: capability → named directory entry → bytes

**Rejected alternative:** Self-describing nested message format. Tahoe chose not to embed structural metadata in the capability itself (beyond the type prefix and essential parameters) because self-description would bloat the short capability strings users must copy-paste. The separation of "what the capability is" (the URI string) from "what it contains" (the referenced bytes) is deliberate.

**Post-hoc assessment:** The design has been stable for 15+ years. The flat URI format aged well for the intended use case (human-copyable, bookmarkable capabilities). The known limitation is that the capability string encodes encoding parameters (erasure coding k/n), meaning changing erasure parameters creates a new capability — this conflates storage policy with identity, a regret noted in the codebase. A richer outer envelope would have helped. This is an example where the flat capability string was correct for the URI layer but insufficient at the envelope layer above it.

---

### 3.3 IPFS / IPLD and Multiformats

**Choice made:** **Recursive nesting via CID prefix chain.** The [CIDv1 specification](https://github.com/multiformats/cid) is:

```
CIDv1 := [0x01 (CID version)] [multicodec (content-type)] [multihash]
multihash := [hash-function-code] [digest-length] [digest-bytes]
```

This is a chain of self-describing prefixes. Every block carries its codec selector. The IPLD data model then defines Nodes as typed values (maps, lists, strings, bytes, links), where a **Link** is always a CID — a content-addressed pointer that carries its own codec and hash function.

**Primary sources:** [Multiformats multicodec](https://github.com/multiformats/multicodec), [IPLD codecs documentation](https://ipld.io/docs/codecs/), [IPLD brief primer](https://ipld.io/docs/intro/primer/).

**Post-hoc assessment — dag-pb failure:** The [IPLD synthesis document on web gateways](https://ipld.io/docs/synthesis/how-ipfs-web-gateways-work/) provides an explicit retrospective:

> *"dag-pb included many responsibilities that we now handle separately: it included the concept of directories (and sharding them) in the codec itself; similarly, it included the concept of files (and sharding their component bytes) in the codec itself. By baking these things into the codec itself, there were several major disadvantages: it was not possible to reuse sharding mechanisms with other codecs; there was no way to switch between high-level and raw views of the data; dag-pb had very (very) rigid ideas about the topology of data. The huge behavioral gaps and asymmetry in representability between dag-pb and other early IPLD codecs produced a great deal of confusion."*

The lesson is sharp: **when you bake semantic interpretation into the codec itself, you lose composability.** The successful design separates: (1) the content-address/hash, (2) the codec selector, (3) the codec implementation. This maps directly to `grid([pcid, payload])` where `pcid` is (1)+(2) and `payload` is the bytes the codec processes.

**Multiformats prefix chain:** The multicodec table approach — varint-prefixed codec codes — has proven generally robust, though there is a known practical issue: the registry is a flat global namespace that requires central coordination. Post-hoc critique from the ecosystem (captured in [IPFS forum discussions](https://discuss.ipfs.tech/t/deprecated-multicodec-implementation-in-go/6992)) is that Go library churn around multicodec registries created confusion, but the conceptual model was sound.

**libp2p protocol negotiation:** Uses multistream-select — a flat text-based negotiation: each side sends the protocol ID as a newline-terminated string, and the other echoes it to confirm. This is an ordered-frame-sequence for the *negotiation handshake*, but after negotiation, the chosen protocol takes over entirely. The design clearly separates transport negotiation (flat sequence) from application protocol framing (protocol-defined). [libp2p protocols documentation](https://libp2p.io/docs/protocols/).

**Iroh** uses QUIC streams with length-prefix framing for its custom protocols, as documented in [iroh message framing tutorial](https://www.iroh.computer/blog/message-framing-tutorial). Iroh's ALPN string (e.g., `iroh/smol/0`) is the protocol identifier for a connection — a flat string per connection, not per message. This is appropriate: QUIC already provides stream multiplexing, so per-message codec dispatch is done at the application layer (where iroh-blobs, for example, uses content-addressed blocks).

---

### 3.4 Secure Scuttlebutt (SSB)

**Choice made:** **Flat JSON message with nested `content` field.** The canonical SSB message structure ([SSB protocol guide](https://ssbc.github.io/scuttlebutt-protocol-guide/)) is:

```json
{
  "previous": "%hash.sha256",
  "author": "@pubkey.ed25519",
  "sequence": 2,
  "timestamp": 1514517078157,
  "hash": "sha256",
  "content": { "type": "post", "text": "message text" },
  "signature": "base64sig.sig.ed25519"
}
```

The outer envelope is flat; the `content` object is where application types live. Private messages replace `content` with an encrypted base64 blob (ending in `.box`), requiring the receiver to detect encryption by the `.box` suffix — a workaround, not a design.

**Post-hoc assessment — critical:** SSB's signing approach — signing the canonical JSON serialization of the message — has been widely criticized. A 2020 analysis ([Secure Scuttlebutt is a cool idea whose realization has fatal flaws](https://derctuo.github.io/notes/secure-scuttlebutt.html)) documents:

> *"The `[Message format, signature]` section finally explains the JSON canonical serialization used for signatures, which is kind of terrible; it refers to the ECMA-262 6th-edition spec for `JSON.stringify`! ...it entirely fails to mention the order of dictionary keys... This is the same design error as XML canonicalization and ASN.1 DER, only botched."*

A 2020 Hacker News thread ([Scuttlebutt is a neat concept, burdened by a bad protocol](https://news.ycombinator.com/item?id=22910752)) adds: *"Signing a message in a format that wasn't designed to be signed"* as the headline critique.

**Metafeeds (Bendy Butt format):** The SSB team's own upgrade uses Bencode (a binary, canonically serializable format) and replaces the flat JSON with a list: `[author, sequence, previous, timestamp, contentSection]`. This migration away from JSON canonical signing toward a proper binary format is itself a post-hoc verdict that the original approach was wrong.

**Implication for PromiseGrid:** SSB's failure mode is precisely what `grid([pcid, payload])` avoids: SSB tried to sign a JSON structure (not bytes), which required a fragile canonical-JSON hack. Using CBOR with `payload` as bytes means you always sign and address bytes, not a higher-level representation.

---

### 3.5 Noise Protocol Framework

**Choice made:** **Ordered flat byte sequence.** The [Noise specification](https://noiseprotocol.org/noise.html) (Trevor Perrin, rev. 2018) defines handshake messages as flat byte buffers assembled by processing a sequence of tokens:

```
tokens: ["e", "ee", "s", "es"]
→ [e.public_key bytes] [s.encrypted bytes] [payload.encrypted bytes]
```

No length fields, no type fields, no nested structures within a Noise message. The spec explicitly states: *"All Noise messages can be processed without parsing"*.

**Rationale documented in spec:**
> *"Noise messages are not self-describing. No type or length fields exist within Noise messages themselves... This means simpler testing (easy to test maximum sizes), memory safety (reduces integer overflow/memory handling errors), and streaming support."*

**Why this was correct for Noise:** Noise is a *session-establishment protocol*, not a *content-addressed message format*. Both endpoints run identical state machines. The ordered frame sequence is the right abstraction when: (a) the message sequence is fixed by protocol design, (b) both endpoints are synchronized, and (c) there are no multi-hop routing requirements. None of these hold for PromiseGrid.

**Post-hoc assessment:** Noise has aged extremely well for its intended purpose. It is used in WireGuard, WhatsApp, Signal, and many others. The lesson is: frame sequences are correct for *handshakes*, not for *application-layer envelopes*. The Noise spec itself notes that *higher protocols use a standard 16-bit length field* — i.e., Noise delegates encapsulation to a layer above.

---

### 3.6 JOSE (JWS / JWE)

**Choice made:** **Implicit recursive nesting with optional type-token.** JOSE supports nested JWTs: a JWT whose payload is another signed or encrypted JWT, indicated by `cty: JWT` in the JOSE header ([RFC 7519, Section 5.2](https://www.rfc-editor.org/rfc/rfc7519)). The structure is:

```
JWE([header: {cty: "JWT"}, enc(JWS([header: {...}, payload, signature]))])
```

Critically, `cty: JWT` is **OPTIONAL** in the inner JWT and **MUST** be present in the outer envelope only when nesting occurs. This means the outer structure signals "there's nesting" but the inner structures do not self-announce their depth.

**Failure modes documented:**

1. **DoS via unbounded recursion:** CVE-2025-53864 (Nimbus JOSE+JWT, [GitHub advisory](https://github.com/advisories/GHSA-xwmg-2g98-w7v9)): *"Connect2id Nimbus JOSE+JWT before 10.0.2 allows a remote attacker to cause a denial of service via a deeply nested JSON object supplied in a JWT claim set, because of uncontrolled recursion."*

2. **Algorithm confusion attacks:** The `alg: none` vulnerability (allowing unsigned tokens to be accepted), and RSA-vs-HMAC key confusion, both stem from JOSE's algorithm agility — each envelope independently negotiates its own algorithm, with no guarantee of consistency.

3. **Signing-order ambiguity:** RFC 7519 Section 11.2 notes that *"syntactically the signing and encryption operations may be applied in any order"* — a normative ambiguity that has caused real-world vulnerabilities.

**PASETO's explicit rejection ([Paragon Initiative, 2018](https://paragonie.com/blog/2018/03/paseto-platform-agnostic-security-tokens-is-secure-alternative-jose-standards-jwt-etc)):**

> *"The JOSE standards give developers enough rope to hang themselves... That is why we made PASETO, which exposes JWT-like validation claims without any of the runtime protocol negotiation and cryptography protocol joinery that caused so many critical JWT security failures."*

PASETO replaces JOSE's algorithm-per-token negotiation with versioned protocols (*"each version has One True Ciphersuite"*). PASETO is intentionally not designed for deep nesting.

**Post-hoc assessment:** JOSE nesting is widely considered a mistake in hindsight. The [IETF guidance document for COSE/JOSE designers](https://www.ietf.org/archive/id/draft-tschofenig-jose-cose-guidance-00.html) acknowledges that COSE re-examined JOSE's design decisions and changed them. JOSE nesting failed because it provided nesting without mandatory, unconditional type-tokens at every layer — the exact problem that `pcid` in `grid([pcid, payload])` solves.

---

### 3.7 COSE (CBOR Object Signing and Encryption)

**Choice made:** **Explicit structure with unconditionally present algorithm identifiers.** [RFC 9052](https://datatracker.ietf.org/doc/rfc9052/) defines COSE structures as CBOR arrays with protected headers wrapped as bstr (byte strings):

```
COSE_Sign1 := [
  protected: bstr .cbor header-map,  // always a serialized CBOR map
  unprotected: header-map,
  payload: bstr / nil,
  signature: bstr
]
```

The key design improvement over JOSE: **protected headers are always a bstr-wrapped CBOR byte string**, not inline JSON. This means the signed bytes are always exactly the bytes of the header encoding — no canonical-JSON problem.

**Design changes from JOSE (RFC 9052, Section 1.3):**

> *"A single overall message structure has been defined so that encrypted, signed, and MACed messages can easily be identified and still have a consistent view. Signed messages distinguish between the protected and unprotected header parameters that relate to the content and those that relate to the signature."*

COSE explicitly supports multi-layer structures (COSE_Encrypt wrapping COSE_Sign1), but each layer is a complete, tagged CBOR structure — not implicit by convention.

**Post-hoc assessment:** COSE is newer than JOSE (RFC 9052 published 2022, updating RFC 8152 from 2017) and is the actively recommended format for constrained device environments (IoT, hardware tokens). It has avoided JOSE's worst failure modes by making structure explicit and unconditional. The [VC JOSE-COSE specification](https://www.w3.org/TR/vc-jose-cose/) now uses both in parallel, with COSE preferred for CBOR contexts.

---

### 3.8 DIDComm v1 vs v2

**Choice made (v1):** Ad-hoc JSON envelope format (Aries RFC 0019/0587) with custom pack/unpack. [Aries RFC 0019](https://identity.foundation/aries-rfcs/latest/features/0019-encryption-envelope/) defined a JWM-like structure with Authcrypt/Anoncrypt, using libsodium primitives.

**Choice made (v2):** JOSE-aligned JWM (JSON Web Messages) with defined nesting rules. [DIDComm Messaging v2.x specification](https://identity.foundation/didcomm-messaging/spec/) defines three nested envelope types:
- `application/didcomm-plain+json` — plaintext (innermost)
- `application/didcomm-signed+json` — sign(plaintext)
- `application/didcomm-encrypted+json` — encrypt(plaintext) or encrypt(sign(plaintext))

The valid nesting rules are explicitly enumerated (not implicit), with `typ` in every JWE making the outer structure self-describing. The v2 specification explicitly states:

> *"`typ` to make JOSE structure formats self-describing. This is particularly helpful in the outermost envelope of any DIDComm message, before unwrapping begins."*

**The shift from v1 to v2** was precisely about making the nesting explicit and typed — adding `application/didcomm-encrypted+json` as an IANA media type so that every envelope is self-announcing. The v1 approach relied on application context to know what envelope format was in use; v2 does not.

**Multi-hop routing design:** DIDComm v2 uses layered nesting for onion routing:
```
Alice → Mediator1 → Mediator2 → Bob:
encrypt(forward(encrypt(forward(encrypt(plaintext, Bob), M2), M1))
```

Each mediator decrypts one layer and sees a `forward` message containing the next encrypted blob. This is explicit recursive nesting — not a frame sequence. The mediator does not need to know the depth of nesting; it just peels one layer.

**Post-hoc assessment:** DIDComm v2 is widely regarded as a significant improvement over v1, specifically because it made type tokens unconditional. The [DIDComm v2.1 specification](https://identity.foundation/didcomm-messaging/spec/v2.1/) is actively maintained as of 2025.

---

### 3.9 Matrix Protocol

**Choice made:** **Flat event format with DAG references.** A Matrix [Persistent Data Unit (PDU)](https://spec.matrix.org/v1.1/rooms/v3/) is a flat JSON object with fields including `type`, `content`, `auth_events` (list of event IDs), and `prev_events` (list of event IDs). The `content` field contains the event-type-specific data:

```json
{
  "type": "m.room.message",
  "content": { "msgtype": "m.text", "body": "Hello" },
  "auth_events": ["$eventId1", "$eventId2"],
  "prev_events": ["$eventId3"],
  "hashes": { "sha256": "..." },
  "signatures": { "example.com": { ... } }
}
```

Events do not nest other events. The DAG structure is through `prev_events` and `auth_events` references — content-addressed pointers to prior events, not inline nesting. Starting from [Room Version 3](https://spec.matrix.org/v1.1/rooms/v3/), event IDs are the base64url encoding of the reference hash of the event — making them content-addressed.

**Post-hoc assessment:** The flat event format with DAG references has proven robust. The change to content-addressed event IDs in Room Version 3+ was a design improvement that aligned with the general principle that content-addressing is correct for long-lived, distributed data. The `type`+`content` structure (`type` as a dispatch key, `content` as opaque bytes/object) is a direct parallel to `grid([pcid, payload])`.

---

### 3.10 ActivityPub

**Choice made:** **Nested JSON-LD objects with remote context loading.** ActivityPub ([W3C specification](https://w3c.github.io/activitypub/)) uses Activity Streams 2.0 with JSON-LD, resulting in structures like:

```json
{
  "@context": "https://www.w3.org/ns/activitystreams",
  "type": "Create",
  "actor": { "type": "Person", ... },
  "object": {
    "type": "Note",
    "content": "Hello, world!"
  }
}
```

The nesting can be multiple levels deep: an Activity contains an Object, which may contain another Object, and so on. Crucially, `@context` is a remote URL that must be fetched and resolved to understand the message.

**Post-hoc assessment — remote context loading:** The 2018 analysis [A proposal for standardising a subset of JSON-LD](https://stephank.nl/p/2018-10-20-a-proposal-for-standardising-a-subset-of-json-ld.html) documents the core problem:

> *"Your app trying to interpret a JSON-LD document from an unknown producer needs all context metadata, including any remote references, or it cannot proceed. All the other complexity can be lumped under algorithmic complexity, but remote context loading adds a slew of potential network and security issues. This is no good in ActivityPub, which is about lots of people exchanging lots of content at potentially very high rates."*

A 2025 security disclosure ([SocialHub, July 2025](https://socialhub.activitypub.rocks/t/potential-security-vulnerability-remote-json-ld-contexts-may-be-used-to-bypass-restrictions-when-arbitrary-objects-are-allowed-to-be-created/5439)) revealed that remote JSON-LD context loading creates an active attack vector: a malicious actor could publish a context document that, when expanded, redefines fields to impersonate other actors.

**The specific failure mode:** ActivityPub's nesting failed not because it uses nesting, but because the type-selector (`@context`) is a *remote URL* rather than a *content-addressed identifier*. If `@context` were a CID (content-addressed), the attack surface would vanish: the context cannot change under you if you hold its hash.

---

### 3.11 Nostr

**Choice made:** **Flat event structure, no envelope nesting.** [NIP-01](https://github.com/nostr-protocol/nips/blob/master/01.md) defines:

```json
{
  "id": "<sha256 of serialized event>",
  "pubkey": "<32-byte hex public key>",
  "created_at": <unix timestamp>,
  "kind": <integer 0-65535>,
  "tags": [["e", "<event-id>"], ["p", "<pubkey>"]],
  "content": "<arbitrary string>",
  "sig": "<64-byte signature of id>"
}
```

No nesting of events within events. `content` is a string; `tags` are flat arrays of strings. References to other events are via `e` tags (event IDs as hex strings), not inline objects.

**Rationale:** Nostr's design philosophy, attributed to fiatjaf, centers on extreme implementability: a relay should be buildable in a weekend, and clients should be able to process events with minimal parsing logic. The [nostr-protocol/nips README](https://github.com/nostr-protocol/nips) makes the simplicity value explicit through the protocol architecture.

**Where nesting reappears:** When Nostr needs multi-layer semantics, it defines new event kinds on top of the flat base. NIP-59 ("Gift Wrap") defines event kind 1059 whose `content` is a JSON-encrypted event — effectively manual nesting using the content string. This is nesting by convention, without a type-dispatch mechanism. It requires client-side knowledge of which `kind` values use encryption. The flat format forced this workaround.

**Post-hoc assessment:** Nostr's flatness is correct for its use case (single-hop broadcast, verifiable posts). Where it reaches for multi-layer semantics, it unsatisfyingly re-invents manual nesting. The design is a good example of "flat is a special case of nested" — the degenerate one-level case.

---

### 3.12 Bluesky / AT Protocol

**Choice made:** **Typed records in CBOR (DRISL), with CID-based linking.** The [AT Protocol data model](https://atproto.com/specs/data-model) uses DRISL-CBOR (a normalized subset of DAG-CBOR) for all signed and addressed data. Records are typed CBOR maps with a `$type` field pointing to a Lexicon schema. Links between records are CIDs (content-addressed hashes). The Lexicon schema language defines the expected shape of each record type.

**Key design choices:**
- Signing and hashing always operate on CBOR bytes, not JSON strings — avoiding SSB's canonical-JSON failure
- `$type` is always present in records — avoiding JOSE's optional type-token failure
- Links are CIDs, not URLs — avoiding ActivityPub's remote-context-loading vulnerability
- Records are flat within their type (no arbitrary recursive nesting of arbitrary records) but can reference other records by CID

**Post-hoc assessment:** AT Protocol (launched 2023, actively developed as of 2026) appears to have learned from prior ecosystem mistakes. The [Lexicon style guide](https://atproto.com/guides/lexicon-style-guide) explicitly warns against overly complex schemas. The use of CIDs for record links is a direct adoption of the IPLD model.

---

### 3.13 Promise Theory (Mark Burgess)

**What the theory says about nesting vs stacking:**

Mark Burgess's [Promise Theory](https://markburgess.org/promises.html) (Burgess & Bergstra, 2004 onward) does not directly specify wire formats. However, it provides structural guidance that maps onto the nesting-vs-stacking question:

**Against sequential chains:** From [Some Notes About Promise Theory](http://markburgess.org/PromiseMethod.pdf):

> *"Avoid serial (intermediate) agent chains, as these represent hard dependencies that increase both the number of agents and the number of promises. The probability of a chain of agents not keeping its promise is like the sum of probabilities for each agent in the chain (≃Np → 1)."*

This is a direct argument against the "ordered frame sequence" model as a general pattern: a promise stack where each frame's meaning depends on the prior frame is exactly the serial chain that promise theory warns against.

**The bottom-up aggregation rule:** Promise theory recommends:

> *"Put all promisable capabilities and properties into the smallest agents that can represent them, and inherit the promises upwards by aggregation."*

This maps to recursive nesting: each layer is a complete, independently-meaningful agent-level unit (the `grid([pcid, payload])` structure), and layers are composed by containment rather than by sequence.

**The atomicity rule:** *"Represent every property, which can be determined independently, as a separate agent or component."* This strongly implies that the codec selector (`pcid`) and the payload should be independently addressable — which is true in `grid([pcid, payload])` and false in an ordered frame sequence where meaning is positional.

**No wire-format operationalization found:** No primary source was found documenting an attempt to directly operationalize promise theory into a wire protocol. The Promise Theory perspective on data networks paper ([arXiv 1405.2627](https://arxiv.org/pdf/1405.2627)) uses promises to model IP addressing (multiplet/container model with n-tuple addressing) but does not prescribe a specific wire format. The "containers in containers" model in that paper — where each container level has its own forwarding promises — is architecturally similar to recursive nesting.

---

## 4. Cross-Cutting Failure Modes

### 4.1 The Optional Type-Token Failure (JOSE, ActivityPub, SSB)

The most common failure mode across these systems is **making the type-selector at each layer optional or implicit.** In JOSE, `cty: JWT` is optional in nested JWTs; in ActivityPub, `@context` is a remote URL that may change; in SSB's classic format, the content `type` is present by convention but there is no outer envelope type-token. The consequence in each case is: receivers cannot know what they have without examining the content — and when content is encrypted, they cannot examine it. This forces implementations to make unsafe assumptions about what layer they're at, leading to algorithm confusion, type confusion, and signature-stripping attacks.

The fix is always the same: **a mandatory, unconditional, content-addressed type-token at every layer.** COSE's `protected: bstr .cbor header-map` (always present, always a serialized byte string), Matrix's `type` field (always present in every PDU), and Nostr's `kind` field (always an integer) are all correct implementations of this principle.

### 4.2 Remote-URL Type Selectors (ActivityPub, early IPLD/dag-pb)

ActivityPub's `@context: "https://www.w3.org/ns/activitystreams"` and dag-pb's implicit codec assumptions both demonstrate the failure mode of **type selectors that resolve through network I/O.** The problems compound: the resource at the URL can change, network fetches introduce latency and attack surface, and caches can serve stale or malicious data. The [2025 ActivityPub security disclosure](https://socialhub.activitypub.rocks/t/potential-security-vulnerability-remote-json-ld-contexts-may-be-used-to-bypass-restrictions-when-arbitrary-objects-are-allowed-to-be-created/5439) is the most recent concrete manifestation. Content-addressed type selectors (CIDs) eliminate this entirely: the hash of the protocol definition is the selector, and a different definition is simply a different CID.

### 4.3 Canonical Serialization for Signing (SSB, early JWS)

SSB's requirement to sign canonical JSON created an ongoing interoperability hazard: key ordering is not defined by JSON, whitespace rules differ, and ECMA-262's `JSON.stringify` behavior is implementation-defined. The [critique by Kragen Javier Sitaker](https://derctuo.github.io/notes/secure-scuttlebutt.html) and the HN discussion ([2020](https://news.ycombinator.com/item?id=29672518)) confirm this is not theoretical. The SSB team's own migration to Bendy Butt (Bencode-based) validates the diagnosis. **The fix: sign bytes, not structured objects.** COSE's protected headers as `bstr` and AT Protocol's DRISL-CBOR (CBOR that has canonical normalization rules) both do this correctly. `grid([pcid, payload])` where `payload` is always `bytes` avoids this problem entirely.

### 4.4 Algorithm Agility in Nested Structures (JOSE)

JOSE's per-token algorithm negotiation — where each envelope independently specifies its cryptographic algorithm — creates an unbounded combinatorial attack surface. PASETO's [2018 analysis](https://paragonie.com/blog/2018/03/paseto-platform-agnostic-security-tokens-is-secure-alternative-jose-standards-jwt-etc) documents: *"most vulnerabilities are in the 'joinery between components' rather than direct cipher breaks."* The attack is not breaking AES or HMAC individually; it is confusing which algorithm to apply to which layer. Frame-sequence models are particularly vulnerable here because the algorithm for frame `i` may be specified in frame `i-1`.

### 4.5 Depth Attacks in Recursive Nesting (JOSE, JSON-LD)

The Nimbus JOSE+JWT CVE-2025-53864 and related JSON-LD issues demonstrate that **recursive structures without depth limits enable resource exhaustion attacks.** A legitimate recursive design can be exploited by a malicious actor who sends a 1000-level-deep nested structure. The fix is **always a configurable depth limit**, which every content-addressed recursive system should implement. This is not an argument against nesting; it is an argument for an explicit depth budget in the parser/dispatcher. COSE does not have this problem because nested COSE structures are uncommon by design and depth is bounded by application schema.

### 4.6 Codec Rigidity Baked Into Data (dag-pb, early IPLD)

The dag-pb case — where sharding mechanisms and file/directory semantics were baked into the codec — is the textbook example of a **codec doing too much.** When the codec owns too many responsibilities, changing any one of them requires a new codec, breaking all existing data. The lesson is that **a codec/protocol identifier should be narrow and stable** — it should specify how to interpret the payload bytes, not what the payload should contain. `pcid` as a content-addressed identifier of the *handler* is the correct minimal interface: the pcid is stable (it's a hash), the handler can evolve by publishing a new pcid.

---

## 5. Implications for PromiseGrid Wire-Lab

Based on the cross-system evidence, the following recommendations apply directly to the PromiseGrid envelope design:

### 5.1 The `grid([pcid, payload])` model is the right shape

The hypothesis `grid([pcid, payload])` — a CBOR two-element array with a content-addressed protocol identifier and opaque payload bytes — is structurally isomorphic to the patterns that aged well across multiple systems:

- It parallels **IPLD's CIDv1**: a codec-selector prefix followed by addressed content
- It parallels **COSE_Sign1**: a protected header (bstr-wrapped, always present) followed by opaque payload bytes
- It parallels **Matrix PDU**: a `type` field (always present dispatch key) followed by `content` (opaque to the router)
- It is **structurally consistent with CapTP's object references**: `pcid` is the capability/protocol reference, `payload` is the argument

The recursive form — `payload` being another `grid([pcid', payload'])` — is exactly how the above systems compose layers, and it is the right mechanism. It is not a "misreading" of promise theory; it is the operationalization of the bottom-up aggregation rule: each layer is a complete promise-making agent that contains another such agent.

### 5.2 `pcid` must be unconditional and content-addressed

Do not make `pcid` optional. Do not allow `pcid` to be a URL that resolves over the network. Both are failure modes confirmed by ActivityPub's `@context` and JOSE's optional `cty` token. `pcid` must be:

1. **Present at every layer** — even if the "protocol" is just "raw bytes, pass through"
2. **Content-addressed** — a hash or CID-like value, not a mutable registry entry or URL
3. **Deterministic** — given a `pcid`, any compliant implementation knows the dispatch without network I/O

This is the key design invariant that distinguishes the `grid([pcid, payload])` approach from JOSE-style nesting. JOSE failed because `cty: JWT` was optional and the algorithm was a string negotiated per-token. `grid([pcid, payload])` succeeds if `pcid` is always a content hash.

### 5.3 The "promise stack" framing is a transport concern, not an envelope concern

The Noise Protocol Framework, QUIC framing, and libp2p multistream-select all demonstrate that **ordered frame sequences are correct at the transport/handshake layer, not the application-message layer.** A "promise stack" of frames to peel is appropriate for describing the sequence of cryptographic operations in a session handshake, or for QUIC stream multiplexing. It is the wrong abstraction for the content-addressed, capability-bearing message envelope of PromiseGrid.

If PromiseGrid has a session-establishment phase (e.g., negotiating which pcids are in scope, establishing shared secrets), an ordered frame sequence is appropriate for that phase. Once in the data plane, recursive nested `grid([pcid, payload])` is appropriate.

### 5.4 Payload must always be bytes (not a re-parsed structure)

SSB's canonical-JSON failure, JOSE's algorithm-confusion attacks, and ActivityPub's JSON-LD expansion problems all stem from the same root: the payload is parsed *before* the signature is verified, at the wrong layer. In `grid([pcid, payload])`, `payload` should always be `bytes` (a CBOR bstr). The bytes are what get signed and hashed. The handler identified by `pcid` is responsible for parsing those bytes. This separation is not optional — it is what makes content-addressing composable.

### 5.5 Implement a depth limit

Whatever the final design, implement a configurable recursion depth limit (e.g., `max_depth=32` as a default, configurable to 0 for "trust only top-level"). CVE-2025-53864 and the general class of recursive-descent DoS attacks are real and have hit production systems. The depth limit should be enforced by the dispatcher, not left to each handler.

### 5.6 Explicit recommendation on `grid([pcid, payload])` vs flat frame sequence

| Criterion | `grid([pcid, payload])` recursive | Flat frame sequence |
|---|---|---|
| Content-addressability per layer | ✅ hash each `grid(...)` | ❌ frames are not independently addressable |
| Handler dispatch by content | ✅ `pcid` is the dispatch key | ❌ position is the dispatch key |
| Multi-hop routing (pass-through without inspection) | ✅ router sees one layer | ✅ router sees all frames |
| Promise theory bottom-up aggregation | ✅ each layer is an independent agent | ❌ frames depend on predecessors |
| Long-term evolvability | ✅ new pcid = new protocol version | ❌ adding a frame type breaks all parsers |
| Match for aligned systems (IPLD, COSE, Matrix, DIDComm v2) | ✅ | ❌ |

**Recommendation: Use `grid([pcid, payload])` with payload recursion. Do not use a flat frame sequence as the primary message model.** The frame-sequence model can be used for session handshakes and transport negotiation, but not for the grid message envelope itself.

---

## 6. Open Questions for the User

**Q1. What is the expected depth budget for nested `grid` messages?**
Most real deployments will have 3-5 layers (e.g., network encryption / routing / protocol / application-content). But the theoretical maximum needs to be defined. Is there a use case that requires unbounded depth? If not, a fixed max_depth in the spec is the right call now, before implementations diverge.

**Q2. Is `pcid` a pure content hash, or does it carry version/routing metadata?**
CIDv1 carries both the codec identifier and the multihash. If `pcid` is purely a content hash of the protocol spec, you need a bootstrap registry or well-known pcids for the `grid` protocol itself. If `pcid` carries a version varint prefix (like CIDv1's multicodec prefix), you gain forward compatibility at the cost of slightly more complex parsing. Which tradeoff is right for PromiseGrid?

**Q3. How are capability references (not just content references) represented in payload?**
CapTP uses integer table IDs for capability references. IPFS uses CIDs for content references. PromiseGrid claims to be capability-based *and* content-addressed. Are capability references always content-addressed (i.e., the capability is a CID of a protocol/actor definition)? Or are there ephemeral capabilities that cannot be content-addressed? The answer changes whether the session-level table-of-capabilities model (like CapTP's four tables) is needed alongside `grid([pcid, payload])`.

**Q4. What is the canonical serialization for `grid` messages used for signing and hashing?**
`grid([pcid, payload])` as a CBOR two-element array has a canonical CBOR serialization (deterministic map key ordering, length-delimited arrays). But does this apply before or after any encryption layer? The lesson from SSB and JOSE is that this must be specified unambiguously in the spec before any implementation exists.

**Q5. Does PromiseGrid need onion routing (multi-hop privacy-preserving routing), and if so, does each hop's `grid` layer need to be distinguishable by the router without decryption?**
DIDComm v2 and Tor both solve this by making the outermost layer's `pcid` a well-known "forward" protocol, so the router knows to peel exactly one layer without seeing the inner content. If PromiseGrid has routing mediators that should not see the inner message, this implies a "routing pcid" vs "application pcid" distinction. Alternatively, if PromiseGrid is a pure content-addressed lookup (no routing privacy), this question may not apply.

---

## Sources Consulted

| System | Primary Sources | Notes |
|---|---|---|
| E language / CapTP | [erights.org/elib/distrib/captp/4tables.html](http://erights.org/elib/distrib/captp/4tables.html) | Only structural overview retrieved; wire format detail limited |
| Spritely / OCapN | [files.spritely.institute/.../CapTP](https://files.spritely.institute/docs/guile-goblins/0.13.0/CapTP-The-Capability-Transport-Protocol.html), [spritely-core.html](https://spritely.institute/static/papers/spritely-core.html), [github.com/ocapn/ocapn](https://github.com/ocapn/ocapn) | Spec pre-standardization; wire format not finalized |
| Cap'n Proto | [capnproto.org/rpc.html](https://capnproto.org/rpc.html), [capnproto.org](https://capnproto.org) | ✅ |
| Tahoe-LAFS | [tahoe-lafs.readthedocs.io/…/uri.html](https://tahoe-lafs.readthedocs.io/en/latest/specifications/uri.html), [tahoe-lafs.org/trac/…/Capabilities](https://tahoe-lafs.org/trac/tahoe-lafs/wiki/Capabilities) | ✅ |
| IPFS / IPLD | [ipld.io/docs/codecs](https://ipld.io/docs/codecs/), [ipld.io/docs/intro/primer](https://ipld.io/docs/intro/primer/), [IPLD synthesis](https://ipld.io/docs/synthesis/how-ipfs-web-gateways-work/) | ✅ |
| Multiformats / CID | [github.com/multiformats/cid](https://github.com/multiformats/cid), [github.com/multiformats/multicodec](https://github.com/multiformats/multicodec) | ✅ |
| libp2p / Iroh | [libp2p.io/docs/protocols](https://libp2p.io/docs/protocols/), [iroh.computer/blog/message-framing-tutorial](https://www.iroh.computer/blog/message-framing-tutorial) | ✅ |
| SSB | [ssbc.github.io/scuttlebutt-protocol-guide](https://ssbc.github.io/scuttlebutt-protocol-guide/) | ✅ |
| SSB critique | [derctuo.github.io/notes/secure-scuttlebutt.html](https://derctuo.github.io/notes/secure-scuttlebutt.html), [HN 22910752](https://news.ycombinator.com/item?id=22910752) | ✅ |
| Noise Protocol | [noiseprotocol.org/noise.html](https://noiseprotocol.org/noise.html) | ✅ |
| JOSE / JWT | [RFC 7519](https://www.rfc-editor.org/rfc/rfc7519), [CVE-2025-53864](https://github.com/advisories/GHSA-xwmg-2g98-w7v9) | ✅ |
| COSE | [RFC 9052 datatracker](https://datatracker.ietf.org/doc/rfc9052/), [IETF guidance draft](https://www.ietf.org/archive/id/draft-tschofenig-jose-cose-guidance-00.html) | ✅ |
| PASETO | [paragonie.com/blog/2018/03/paseto…](https://paragonie.com/blog/2018/03/paseto-platform-agnostic-security-tokens-is-secure-alternative-jose-standards-jwt-etc) | ✅ |
| DIDComm v1 | [Aries RFC 0019](https://identity.foundation/aries-rfcs/latest/features/0019-encryption-envelope/) | ✅ |
| DIDComm v2 | [identity.foundation/didcomm-messaging/spec](https://identity.foundation/didcomm-messaging/spec/), [v2.1](https://identity.foundation/didcomm-messaging/spec/v2.1/) | ✅ |
| Matrix | [spec.matrix.org/v1.1/rooms/v3](https://spec.matrix.org/v1.1/rooms/v3/) | ✅ |
| ActivityPub | [w3c.github.io/activitypub](https://w3c.github.io/activitypub/), [JSON-LD subset proposal](https://stephank.nl/p/2018-10-20-a-proposal-for-standardising-a-subset-of-json-ld.html), [SocialHub security issue 2025](https://socialhub.activitypub.rocks/t/potential-security-vulnerability-remote-json-ld-contexts-may-be-used-to-bypass-restrictions-when-arbitrary-objects-are-allowed-to-be-created/5439) | ✅ |
| Nostr | [NIP-01](https://github.com/nostr-protocol/nips/blob/master/01.md), [nostr-protocol/nips](https://github.com/nostr-protocol/nips) | ✅; fiatjaf's own blog post on design rationale not found |
| Bluesky / AT Protocol | [atproto.com/specs/data-model](https://atproto.com/specs/data-model), [atproto.com/guides/overview](https://atproto.com/guides/overview), [lexicon style guide](https://atproto.com/guides/lexicon-style-guide) | ✅ |
| Promise Theory | [markburgess.org/PromiseMethod.pdf](http://markburgess.org/PromiseMethod.pdf), [arXiv 1405.2627](https://arxiv.org/pdf/1405.2627) | ✅ |
| Promise Theory wire format | (none found) | No documented operationalization of promise theory into a wire format found in any primary source |

### Sources tried but inaccessible or incomplete
- `ocapn.org/spec/` — 404; OCapN spec is not yet published at a stable URL
- `erights.org/elib/distrib/captp/4tables.html` — retrieved but truncated; wire-level message format not in the accessible portion
- Fiatjaf's personal blog posts on Nostr design rationale — no specific post on nested-vs-flat design decision was found; this claim is marked **unsourced — flag for follow-up**
- The CapTP message format specification for OCapN (Syrup encoding) — not yet published in a normative form as of mid-2026
