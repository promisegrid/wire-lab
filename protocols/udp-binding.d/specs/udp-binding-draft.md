# UDP-binding v0 (draft)

## Status

DRAFT. One-page sketch authored alongside TE-29
(`docs/thought-experiments/TE-20260501-215027-protocols-as-simulated-repos-and-binding-layer.md`).
This is the first concrete L4-binding spec under the layer
decomposition locked in TE-29. Subject to revision before freeze.

## Abstract

This document specifies how PromiseGrid messages are carried over
UDP (RFC 768) in version 0 of the binding. It is intentionally
minimal: one PromiseGrid session message per UDP datagram, no
fragmentation, default port 4646, best-effort delivery promises only.

This binding does not redefine UDP. It defines exactly how
PromiseGrid uses UDP and what a conformant implementation promises.

## Layer position

UDP-binding v0 occupies level 2 in the five-level stack defined by
TE-29:

```
transports/udp/udp-binding-bafkrei...U1/<session-pCID>/<message-pCID>/<message-id>.msg
                ^^^^^^^^^^^^^^^^^^^^^^^^
                this spec
```

Level 1 (`udp`) names the real-world transport (RFC 768), which this
binding does not modify. Levels 3 and above (session, message) are
opaque to this binding.

## Promises (normative)

I promise that:

1. **One message per datagram.** Each PromiseGrid session message
   handed to this binding for sending is transmitted in exactly one
   UDP datagram. The datagram payload equals the session message
   bytes verbatim, with no prefix, suffix, framing, magic number, or
   length field added by this binding.

2. **Maximum size 1232 bytes.** Session messages exceeding 1232
   bytes (the IPv6 minimum MTU minus IPv6 and UDP headers) MUST be
   rejected at the sender with a local error. This binding does not
   fragment. Senders MUST NOT rely on IP fragmentation; if Path MTU
   Discovery reports a path MTU below 1232, the sender MUST surface
   a local error rather than send a fragmented datagram.

3. **No delivery, ordering, or deduplication promises.** This binding
   inherits UDP's best-effort semantics. Datagrams may be lost,
   duplicated, reordered, or delivered to unintended ports. Higher
   layers (session protocol) are responsible for any guarantees
   beyond best-effort.

4. **Default port 4646.** Conformant implementations default to UDP
   port 4646 for both send and receive. Any other port is allowed by
   mutual configuration. This document does not allocate the port; no
   central allocator exists. 4646 is a convention, not a reservation.

5. **Address shape `host:port`.** Peer addresses presented to the
   binding take the form `host:port` where `host` is an IPv4 or IPv6
   literal or a DNS name resolvable at send time, and `port` is a UDP
   port number. Multicast group addresses are permitted; behavior on
   multicast send is identical (one datagram per recipient join).

6. **Receive contract.** A conformant receiver binds a UDP socket to
   the configured port, calls `recvfrom` (or equivalent) with a
   buffer of at least 1232 bytes, and hands each datagram payload up
   to the session layer as one complete message. Datagrams whose
   source address is not in the configured peer set MAY be dropped
   silently or logged; either behavior is conformant.

7. **DSCP default zero.** Senders SHOULD set IP DSCP to 0 (default
   forwarding). Other DSCP values are permitted by configuration but
   are out of scope for this spec.

8. **No checksums beyond UDP's.** This binding does not add a
   checksum. UDP's own checksum (RFC 768) is sufficient. Senders MUST
   NOT disable the UDP checksum (zero-checksum optimization is
   prohibited at this binding).

9. **No connection state.** This binding is stateless. No handshake,
   no keep-alive, no reconnect. Senders may send to any address at
   any time; receivers may receive from any address at any time
   (subject to the peer-set filter at promise 6).

10. **Simulation artifact format.** When this binding runs inside the
    wire-lab simulator, it writes one file per datagram seen on the
    wire to:

    ```
    transports/udp/<this-binding-pCID>/<session-pCID>/<message-pCID>/<message-id>.msg
    ```

    File contents equal the exact datagram payload bytes (the same
    bytes handed up to the session layer, byte-for-byte). The session
    and message-protocol pCIDs in the path are derived by the session
    layer above; this binding does not parse them. The `<message-id>`
    in the filename is supplied by the session layer (typically a
    content hash of the session message bytes).

## Anti-promises (non-normative clarifications)

This binding does not promise:

- **Reliability.** Datagrams may be lost. Use a session protocol that
  handles retransmit if you need reliability.
- **Order.** Datagrams may arrive in any order.
- **Privacy or authenticity.** UDP is unencrypted. This binding does
  not add encryption or signatures. Use session/message-layer
  cryptography (per C-6) for authenticity.
- **Spam resistance.** A receiver bound to UDP/4646 will receive any
  datagram any sender sends to it. Filtering and rate-limiting are
  out of scope for v0.
- **NAT traversal.** Pure UDP without rendezvous works for hosts with
  reachable addresses. NAT traversal is a separate concern (probably
  a future binding or a session-layer feature).
- **Path MTU discovery.** Implementations MAY perform PMTUD; this
  binding only requires that they error out below 1232 rather than
  fragment.

## Test vectors (placeholder)

To be added in TODO 018. At minimum:

- TV-1: a 612-byte session message round-trips byte-for-byte through
  loopback UDP/4646.
- TV-2: a 1232-byte session message round-trips byte-for-byte.
- TV-3: a 1233-byte session message produces a local sender error
  before any datagram leaves the host.
- TV-4: a malformed datagram (e.g., truncated by transport-level
  loss simulation) is handed up to the session layer without
  modification by this binding (it is not this binding's job to
  validate session-layer content).
- TV-5: simulation-artifact file written for TV-1 contains exactly
  the 612 bytes of TV-1's session message.

## Reference implementation

To be authored in TODO 018 under
`tools/udp-binding/` with at minimum:

- `Send(msg []byte, addr Addr) error`
- a recv loop that reads from a configured UDP socket and invokes a
  caller-supplied `Handle(msg []byte, src Addr)` callback.

Language: Go (per Steve's standing preference).

## Forking and versioning

Per C-4 (forking is normal, TE-28), any author may publish a UDP
binding with different choices (e.g., different default port,
different size limit, added framing). Such a fork takes a different
pCID and lives as a sibling under `protocols/`. Two parties wishing
to interoperate must agree on the same UDP-binding pCID; they cannot
silently mix bindings.

This binding may evolve to v1 if a load-bearing change is required
(e.g., explicit fragmentation support, NAT-traversal hooks).
Migration from v0 to v1 follows the rules in TE-28 OQ-100.2.

## Bibliography

- RFC 768 (UDP)
- RFC 8200 (IPv6, for the 1232-byte size derivation)
- TE-29 (this binding's layer position)
- TE-28 (load-bearing constraints C-1 through C-6)

## Open questions

- OQ-UDP-1: Should the default port 4646 be parameterized in some
  registry-free way (e.g., derived from a hash of the binding pCID)?
  Lean: no, fixed convention is simpler.
- OQ-UDP-2: Should the 1232-byte limit be raised for known-IPv4-only
  paths? Lean: no, uniformity is more valuable than the few extra
  bytes.
- OQ-UDP-3: Multicast semantics in the simulation. How does
  `transports/udp/...` represent a single multicast send arriving at N
  receivers? Lean: one file per (sender, receiver) pair; deduplication
  by content hash makes this storage-cheap.
