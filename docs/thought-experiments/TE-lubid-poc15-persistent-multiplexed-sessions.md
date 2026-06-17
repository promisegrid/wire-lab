# TE-lubid: POC15 persistent multiplexed TCP sessions

TE ID: TE-lubid
## Status
decided

## Decision under test

Should POC15 replace per-message TCP dials with persistent multiplexed sessions across app/kernel and kernel/kernel paths, using exact message CIDs and parent links rather than RPC request IDs for correlation?

## Assumptions

- POC15 remains the implementation target.
- All TCP paths should be persistent where TCP exists: app-to-local-kernel, app receive registration/delivery, and kernel-to-kernel peer forwarding.
- Message identity is the exact SHA-256 / raw CIDv1 of the complete CBOR grid envelope bytes.
- ACKs and response envelopes should parent-link the request message CID when POC15 code generates them.
- Promise Theory constraints remain primary: sessions are local transport promises, not authority, permission, or RPC command channels.
- Clean-run state remains resettable POC state; persistent means persistent during a process/run, not between clean runs.

## Alternatives

1. Keep one-shot dials for outbound messages and only document the current behavior.
2. Add persistent but sequential sessions: one TCP connection per peer/path with at most one in-flight request.
3. Add persistent multiplexed sessions with in-flight requests keyed by exact request message CID and ACK/response parent links.
4. Add an RPC-like request ID or method/session layer.

## Scenario analysis

### Normal operation

Alice sends multiple relationship, storage, compute, and route promises to Bob and Carol during one run. With one-shot dials, every message pays connection setup cost and the run cannot test long-lived relationship between transport promises and app promises. Sequential persistent sessions reduce dial churn but still force head-of-line blocking. Multiplexed sessions let multiple promise envelopes be in flight while still preserving message identity: the response names the request by parent link, and the local pending map uses the request message CID. RPC-style request IDs would work technically but would create a second identity system parallel to the message DAG.

### Concurrent actors and out-of-order replies

Alice may ask Bob for storage and Carol for compute while Dave verifies a result. A sequential session either serializes these or requires several ad-hoc connections. A multiplexed session can accept ACKs in any order because the parent link identifies the request. This matches the existing message-DAG direction better than RPC sequencing.

### App/kernel boundary

A local app should not dial peer apps directly. It should maintain a persistent app/kernel session that carries receive promises, outbound envelopes, inbound deliveries, and ACKs. The kernel should demux ACKs by parent CID and route fresh app envelopes by pCID and target. This keeps the app/kernel boundary as ordinary PromiseGrid envelopes rather than a command API.

### Kernel/kernel boundary

A kernel-to-kernel TCP connection can be a long-lived transport promise between two local kernels, but it does not make either kernel authoritative. A peer kernel may send fresh envelopes or ACKs on the same stream. The receiving kernel decides locally whether a frame is a pending response or a fresh inbound promise by checking parent links against its pending map.

### Failure and reconnect

If a session closes with pending requests, the local side cannot know whether the peer broke a promise unless the peer had made a concrete promise and the application semantics say the promise is broken. Most session failures are local transport uncertainty. The implementation should mark pending waits as local failure/non-commitment, record session close/reconnect events, and allow later sends to open a fresh persistent session.

### Backpressure and rate limiting

A persistent session needs bounded queues and timeouts. These are local promises about what the process will accept or send. Queue-full behavior is a local non-commitment, not a peer failure. This fits existing POC pressure events.

### Long-horizon evolution

CID/parent correlation keeps session multiplexing compatible with the raw-message CAS/DAG work. Future pCIDs can define richer response shapes, but POC15 does not need to add a universal request ID. If later transport profiles add TLS/QUIC/Noise, the same message-level correlation can remain unchanged.

## Conclusions

Reject alternatives 1, 2, and 4 for this task. Alternative 3 is the chosen design: persistent multiplexed sessions keyed by exact message CIDs, with ACK/response parent links required for POC15-generated responses.

## Implications for open TODOs and pending DIs

- TODO-gogug should gain a DI locking persistent multiplexed POC15 sessions.
- POC15 transport code should gain a reusable persistent-session abstraction.
- POC15 runtime and kernel code should use persistent app/kernel and kernel/kernel sessions.
- Analyzer gates should check session reuse and ACK parent-link coverage.
- README and DEV-GUIDE-RESOURCES should describe persistent sessions as transport promises, not RPC.

## Decision status

locked by DI-vopab
