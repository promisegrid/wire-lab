# poc16-secure-tokens-maps-encrypted-payloads

`poc16-secure-tokens-maps-encrypted-payloads` is the planned successor to
`poc15-multihop-multiarity-dag`. It is executable as a POC15 superset so POC16
additions can be measured against a known regression floor. The purpose is to
preserve POC15's multi-hop forwarding, pCID-owned slot-vector variety,
raw-message CAS/DAG review, agent-accessible sparse CAS, wire-visible parent
links, COSE specimens, useful routed WASM/stdio work, promise-correct route
exclusion, route economics, persistent sessions, and explicit non-monolithic
kernel roles while adding secure CWT-style capability tokens, encrypted payload
profiles, pCID-owned map payload profiles, embedded spec context for LLM agents,
and pCID-selected parser/builder role pressure. Source: `DI-pamob`;
`DI-podut`; `DI-lutuv`; `DI-manul`; `DI-vulit`.

## Current Executable State

The current executable scaffold preserves inherited POC15 app/kernel, shipping,
CAS, compute, WASM, stdio, event-collector, persistent-session, sparse-CAS,
message-DAG, route, and analyzer behavior under a separate POC16 module path and
run root. It now also includes the first
`route_v1` slice: Alice confirms and uses an Alice->Bob->Carol->Dave route
through voluntary neighboring route promises. It now also emits run-local exact
CBOR specimens for pCID-owned multiarity, wire-visible parent-link placement,
COSE-as-payload, COSE-as-proof, and native-proof comparison, then gates those
specimens in the analyzer. The current convergence slice also moves
`route_v1` normal app traffic beyond specimen-only coverage: forwarded route
traffic can use `grid([42(pCID), parents, payload, proof])`, payload parent
links are exercised by a route-carried message, Alice reuses a route under an
explicit local lifetime, asymmetric response-path terms are recorded, route
credits are promised/earned/spent, and Peggy/Victor receive route-carried
compute envelopes as useful runtime-adapter work. It retains raw message
artifacts for operator review through an observer-only artifact stream. The
current parser-role correction makes the TE-ritig flow executable: local apps
talk to a container-local parser role, parser roles talk to the transport kernel
with `kernel_transport_v1`, and peer kernels forward exact envelopes to parser
roles that promised the matching pCID. The transport kernel inspects the grid
tag, slot-0 pCID, parent links, exact CID, proof, and kernel-transport control
payloads; it does not parse normal application payload fields such as `to` to
choose app behavior. Replies are correlated by ACK/response envelope parent
links to exact request message CIDs, not by RPC request IDs, so the same
raw-message DAG used for operator review also supports session demux. Source:
`DI-lutuv`; `DI-lihir`; `DI-kohuj`; `DI-tuhop`; `DI-mosat`; `DI-vopab`;
`DI-vulit`; `DI-gazin`.

The POC16-specific slice adds four broad profile pCIDs rather than
pCID-per-operation fragmentation: `secure_capability_v1`,
`encrypted_payload_v1`, `parser_builder_role_v1`, and
`map_payload_profile_v1`. The secure-capability profile signs CWT-style CBOR
claims with COSE_Sign1 and records expiry, audience, bearer-transfer, and replay
checks. The encrypted-payload profile seals real ciphertext with AEAD, records
wrong-recipient and tamper rejection, and distinguishes visible parent links from
hidden parent commitments. The parser/builder profile records slot-0 pCID
selection of a local parser or builder role while keeping destination and route
semantics inside pCID-owned payloads. The map-payload profile shows that maps are
permitted when a pCID spec chooses them, while constrained arrays remain preferred
for limited devices. Source: `DI-vulit`.

The active shipping/device workflow now uses one protocol-family pCID,
`production_shipping_v1`, rather than active runtime pCIDs for
`postal_scale_v1`, `ups_label_v1`, `accounting_v1`, and `printer_port_v1`.
Those older specs remain in `docs/protocols/` only as historical/specimen
evidence.
Inside `production_shipping_v1`, `promise_about` selects weighing, address
lookup, label printing, print capability issue/redeem, or shipment update body
slots. Source: `DI-gazin`.

The current convergence slice hardens that transport and storage pressure:
every persistent session now carries a local `session_id`, frame/request/response
counters, and exactly one terminal event with a close reason. Alice's local
trust-change scenario closes and later reopens the active app/kernel session so
the run proves an existing stream is not silently reused after local distrust.
CAS capability tokens are now CWT-style CBOR payloads protected by the same
well-known COSE_Sign1/Ed25519 library pattern as `local_lifecycle_v1`, with
issuer, subject, scope, content CID, expiry, nonce, and transferability terms.
This is still executable POC pressure rather than a final token standard, but it
removes the older custom string-map token body from normal storage-token traffic.
The route workflow also records peer transit-exclusion promises, explicit route
lifetime exhaustion, local non-send after expiry, and route renewal before reuse.
Source: `DI-mapop`; `DI-lurov`.

The current CAS/DAG slice now also gives each app its own local sparse CAS view.
Agents may store exact messages, local state bytes, encrypted blobs named by
ciphertext CID, or peer-served content. Object bytes are now stored as per-agent
filesystem CAS files under `stores/<agent>/cas-objects/`, while
`durable-state.json` remains the mutable root/index for metadata, tokens,
journals, compute cache, and message-DAG entries. Those stores are intentionally
incomplete: a local message-DAG index may record missing parents, peer storage is
voluntary, bearer storage tokens can compensate a holder for retrieving bytes
from the issuer, and local GC retains paid/pinned/encrypted objects while
removing pressure-tagged temporary objects. Source: `DI-manul`; `DI-fagog`.

The current parser/listener/app separation matches `TE-ritig` model 3 for
runtime routing: local apps send exact envelopes to the parser role, the parser
role uses `kernel_transport_v1` to ask the transport kernel to carry exact bytes,
peer kernels dispatch by slot-0 pCID to parser roles that promised that pCID, and
parser roles deliver only to local apps that made matching receive promises. The
transport kernel does not use normal application payload fields such as `to`,
`promise_about`, DID, CID-rooted path, or route terms to choose app behavior.
There is one bounded implementation-local demux exception: persistent sessions
may inspect a pCID-owned `outcome` projection only to classify parent-linked
ACK-like responses; that projection is not a route decision, not a PromiseGrid
API, and not deployment readiness. Source: `DI-rapuk`; `DI-gazin`; `TE-ritig`.

Analyzer success is also intentionally scoped. The POC16 analyzer now reports
`poc_scope_fitness.poc_scope_complete` rather than a deployment-readiness flag. A
passing clean run means the current POC acceptance surface was covered;
production use still requires separate design, reliability, security, operations,
and API stability review. Source: `DI-rapuk`.

## Superset Requirement

POC16 should preserve POC15 unless a later DI explicitly authorizes a scoped
exception:

- POC11-style autonomous sparse-mesh relationship and economics pressure.
- POC12-style local kernel/app/device workflow and shipping agents.
- POC13-style CAS storage, CID compute, verifier disagreement, replica recovery,
  token lifecycle, retention/GC, backpressure, rate-limit, replay protection,
  bounded trust, and dynamic direct TCP relationship event.
- POC14-style WASM process agents, stdio-only subprocess agents, decentralized
  monitoring signals, hard local distrust, route-exclusion scenarios, pCID-owned
  payload migration event, observer-only event collector, and no shared-volume
  agent coordination.

## POC16 Additions

1. **Real multi-hop forwarding.** A sender should be able to send an envelope to
   a direct peer that voluntarily promises one forwarding hop, eventually
   reaching the target if each hop locally promises the next step.
2. **Route promises, not route authority.** A forwarding path is a chain of
   voluntary peer promises. No node commands another node to forward, and no node
   claims global route truth.
3. **Route exclusion by peer promises.** If Alice does not want traffic to
   transit Mallory, Alice selects peers whose promises and keep/break history are
   good enough for Alice's local risk decision.
4. **Promise economics.** Forwarding can be compensated through reciprocal
   forwarding, relationship value, bearer capability-token payment,
   non-transferable forwarding-capacity tokens, or optional stake/collateral
   promises.
5. **Route durability and asymmetry.** POC16 should test one-shot routes,
   bounded durable routes, reverse-route replies, independent return-route
   discovery, and return-route tokens.
6. **pCID-owned multiarity.** POC16 should test pCIDs that define different
   arity, slot order, parent-link location, proof location, and COSE usage.
7. **Raw message review.** Every raw envelope artifact should be retained in
   run-scoped CAS for later review, while valid message parent links form a
   wire-visible causal DAG.
8. **Useful routed WASM/stdio work.** Peggy and Victor should do valuable work
   for other agents over routed paths, not merely prove adapter plumbing exists.
9. **Kernel as role collection.** POC16 should name transport, app-interface,
   routing, local-resource, and event-retention roles explicitly even when the
   Docker runtime collapses some roles into one process.
10. **Agent-accessible sparse CAS.** Each agent should be able to maintain its
    own incomplete CAS and optional local message DAG, or rely on peer CAS
    promises paid through credits or bearer storage tokens.
11. **Persistent multiplexed TCP sessions.** App/kernel and kernel/kernel TCP
    paths should stay open during a run and carry many exact envelopes. A local
    pending map uses the request message CID as the key, and generated ACKs
    parent-link that CID so no payload-level RPC request ID is needed.
12. **Signed capability-token bytes.** CAS storage and bearer tokens should be
    signed issuer promises carried as normal pCID-owned payload fields, with
    replay/expiry/scope checks performed by the issuer's local state.
13. **Secure CWT-style capability tokens.** Broad token-profile traffic should
    use COSE-signed CBOR claims with expiry, audience, replay IDs,
    transferability, and optional holder confirmation.
14. **Encrypted payload profiles.** Broad encryption-profile traffic should use
    real AEAD ciphertext and record recipient, tamper, visible-parent, and
    hidden-parent behavior as local promise events.
15. **Embedded pCID specs.** LLM-backed agents should receive embedded spec
    context for supported pCIDs, including the spec CID and excerpt, without
    turning prose lookup into runtime routing.

## Multihop Promise Sketch

A route setup is a chain of conditional reciprocal promises:

1. Alice promises Bob that Alice will use, compensate, or reciprocate for a
   route to Dave if Bob can assemble a route that satisfies Alice's constraints.
2. Bob promises Alice only Bob's next hop, such as forwarding route setup traffic
   to Carol if Carol promises a compatible next hop and Alice's reciprocal terms
   are acceptable.
3. Carol makes the same local promise toward Dave or another next hop.
4. Dave returns a reachability promise along the return path, promising to
   receive matching messages within the route's constraints.
5. Alice sends actual traffic only after route-confirmation promises exist.

Each peer owns only its own promise. Bob does not promise that Carol will keep
Carol's promise, and no peer obtains authority over another peer.

## Raw Message Review

The executable POC16 run now stores exact envelope bytes for operator review.
Apps do not mount or read the observer volume; instead, each app emits a one-way
`message_artifact` record to stdout, the local supervisor forwards it to the
observer-only collector, and the collector writes:

- `/run/poc16/<run_id>/message-cas/<exact_cid>.cbor` — binary CBOR envelope
  bytes exactly as sent, received, or acknowledged.
- `/run/poc16/<run_id>/message-dag.jsonl` — one JSON index row per artifact,
  including observer, direction, peer, pCID name, exact CID, optional parent
  exact CID, promise meaning, and relative artifact path.

`poc16-analyze` validates the index by reading each `.cbor` artifact and
recomputing its exact SHA-256 before the clean regression can pass. It traverses
the unique exact-message DAG rather than double-counting duplicate observations
of the same message, then fails the clean regression if a valid parent link is
missing or unreachable. Source: `DI-tuhop`; `DI-kohuj`.

`poc16-cbor-diag` is included in the collector image for operator inspection of
the retained artifacts:

```sh
docker compose run --rm --entrypoint /usr/local/bin/poc16-cbor-diag \
  event-collector -cid <exact_cid>
```

The tool reads `message-dag.jsonl`, opens the matching
`message-cas/<exact_cid>.cbor` artifact, and prints CBOR diagnostic notation
with nested payload/proof byte strings expanded when they contain valid CBOR.
It is read-only and does not affect app, parser-role, or transport-kernel
behavior. Source: `DI-bapif`.

The parser-role correction adds reviewable raw artifacts for the local and peer
flow boundaries:

- `app_to_parser` — an app's exact signed envelope as the local parser role saw
  it.
- `parser_to_kernel_receive` — a parser role's `kernel_transport_v1` promise
  that it can receive exact envelopes for a pCID.
- `parser_to_kernel_carry` — a parser role's `kernel_transport_v1` request to
  carry exact app envelope bytes toward a target agent.
- `kernel_to_parser` — exact peer-originated envelope bytes delivered from the
  transport kernel to the local parser role.
- `parser_to_app` — exact peer-originated envelope bytes delivered by the parser
  role to the local app that promised the pCID.
- `app_to_parser_ack` and `parser_to_kernel_ack` — parent-linked ACK bytes
  flowing back through the same local role chain.

These are observer-only artifacts. They make the role split inspectable after a
run; they do not give apps a shared volume and do not change promise outcomes.
Source: `DI-gazin`.

Parser-flow diagnostics from a clean run look like ordinary PromiseGrid
envelopes with different pCID-selected payload meanings at each boundary:

```cbor-diag
# direction=parser_to_kernel_receive protocol=kernel_transport_v1
1735551332([
  42(h'0001551220...kernel-transport-v1...'),
  << [
      "parser:grace-heidi",
      "kernel:grace-heidi",
      "receive_pcid",
      ["", "I promise this parser role can receive exact envelopes for this pCID ...", "...", ""],
      [["pcid", "cas_storage_v1"], ["transport_action", "receive_pcid"]],
    ] >>,
  << { "signer": "parser:grace-heidi", "signature": "..." } >>,
])

# direction=parser_to_app protocol=production_shipping_v1
1735551332([
  42(h'0001551220...production-shipping-v1...'),
  << [
      "fulfillment",
      "accounting",
      "address_lookup",
      ["", "I promise to receive accounting's local address event ...", "..."],
      ["ORDER-1001", "fulfillment-accounting-production_shipping_v1-000001", ""],
    ] >>,
  << { "signer": "fulfillment", "signature": "..." } >>,
])

# direction=app_to_parser_ack / parser_to_kernel_ack
1735551332([
  42(h'0001551220...production-shipping-v1...'),
  [42(h'0001551220...parent-message...')],
  << [
      "accounting",
      "fulfillment",
      "address_lookup",
      ["kept", "I promise I received and recorded your signed promise message.", "..."],
      ["ORDER-1001", "", "100 Promise Way, Suite 100, Example City, CA 94000"],
    ] >>,
  << { "signer": "accounting", "signature": "..." } >>,
])
```

Shortened examples from a clean run show the intended review surface:

- Map payload profile:

  ```cbor-diag
  1735551332([
    42(h'0001551220...map-payload-profile-v1...'),
    << {
        "act": "promise",
        "from": "alice",
        "promise_about": "map_payload_profile",
        "profile": "self-documenting-map",
        "to": "bob",
      } >>,
    << { "signer": "alice", "signature": "..." } >>,
  ])
  ```

- Compact array payload and failure case:

  ```cbor-diag
  1735551332([
    42(h'0001551220...cas-storage-v1...'),
    << [
        "mallory",
        "grace",
        "present_storage_report",
        ["", "Mallory promises this intentionally corrupted signature is valid.", "..."],
        ["", "", "", ""],
      ] >>,
    << { "signer": "mallory", "signature": "..." } >>,
  ])
  ```

- CWT-shaped capability token:

  ```cbor-diag
  1735551332([
    42(h'0001551220...secure-capability-v1...'),
    [42(h'0001551220...parent-message...')],
    << {
        "act": "promise",
        "issuer": "alice",
        "audience": "bob",
        "scope": "cas_storage_v1",
        "token_b64": "0oRDoQEn...",
        "transferable": "true",
      } >>,
    << { "signer": "alice", "signature": "..." } >>,
  ])
  ```

- Holder-bound token checks use the same `secure_capability_v1` token profile
  with `transferable=false` and a holder confirmation claim. The clean run records
  `cwt_capability_token_holder_bound_checked` and
  `cwt_capability_token_transfer_mismatch_rejected` so a reviewer can confirm
  that Mallory cannot redeem Frank's holder-bound token.

- Encrypted payload:

  ```cbor-diag
  1735551332([
    42(h'0001551220...encrypted-payload-v1...'),
    [42(h'0001551220...parent-message...')],
    << {
        "algorithm": "AES-256-GCM",
        "ciphertext_b64": "...",
        "nonce_b64": "...",
        "recipient": "bob",
        "sender": "alice",
      } >>,
    << { "signer": "alice", "signature": "..." } >>,
  ])
  ```

## Agent-Accessible Sparse CAS

The observer-only `message-cas/` archive is for operator review. It is not the
same thing as an app's local CAS. Each app now keeps run-scoped CAS object files
in its own `stores/<agent>/cas-objects/` directory plus metadata in
`stores/<agent>/durable-state.json`. The JSON file is an index/root, not the
normal byte store; old `cas_objects_b64` state is read only for migration and is
omitted on new saves. Metadata describes objects the agent chose to retain:

- Exact message bytes emitted or received by that app.
- Local state or checkpoint bytes the app wants to retain for the run.
- Encrypted local bytes addressed by ciphertext CID rather than cleartext CID.
- Peer-served bytes retained after storage, replica, or bearer-token promises.

Some apps maintain a local message DAG over exact message bytes they retained.
The DAG is sparse by design: if a parent CID is not in the local store, the app
records a missing-parent event instead of treating the system as inconsistent.
Peer CAS promises are similarly local and voluntary. Bob may promise paid storage
for one CID, Frank may store a replica, and Alice may transfer Bob's bearer
storage token to Frank as an incentive; none of those promises creates a global
store, global retention rule, or global trust decision.

POC16 deliberately mixes local filesystem formats across agents. Some agents
write generic `<safe-cid>.bin` files, some write `.cbor` when the exact object
bytes parse as one complete CBOR item, and some write a local CBOR wrapper file
whose metadata records both the logical/original CID and the stored wrapper CID.
That wrapper is only local storage format pressure: when a peer asks for content
CID `X` that this app promised to store, retrieval still returns the exact
original bytes named by `X` unless a future pCID explicitly negotiates wrapper
bytes. Source: `DI-manul`; `DI-fagog`.

## Multiarity Specimens

POC16 now emits and analyzes pCIDs that define these slot-vector specimens:

- `grid([42(pCID), payload])` for transport/session-auth-only pressure.
- `grid([42(pCID), payload, proof])` for the current signed-message specimen.
- `grid([42(pCID), parents, payload, proof])` for envelope parents before body.
- `grid([42(pCID), payload, parents, proof])` for body before envelope parents.
- A payload-owned-parent form is also exercised by ordinary `route_v1` payload
  fields that name prior exact message bytes while the outer route envelope
  stays in the pCID's chosen slot vector.
- `grid([42(pCID), cose_sign1])` for COSE-as-payload.
- `grid([42(pCID), payload, cose_sign1_detached])` for COSE-as-proof.

Parent links should name the CID of exact prior envelope bytes and should use
CBOR tag-42 IPLD links where the specimen includes parent links. The run CAS DAG
is separate: malformed bytes and non-message artifacts are reviewable run
artifacts, but they are not valid PromiseGrid parent-linked messages.

## Persistent Sessions

POC16 now treats TCP persistence as a transport promise made locally for the
duration of one clean run. Each app opens one app/parser-role frame stream,
registers its pCID receive promises on that stream, sends outbound envelopes on
that same stream, and receives parser-delivered peer messages on that same
stream. Each parser role opens one parser/kernel control stream using
`kernel_transport_v1`, and kernels keep reusable peer-kernel sessions per
endpoint. Fresh inbound frames are handled as ordinary signed PromiseGrid
envelopes; replies are matched to pending requests only when their envelope
parent links include the exact request CID.

This is intentionally not an RPC channel. There is no universal method name,
request number, route authority, or kernel-owned trust judgment. The kernel owns
exact-byte transport and parser-role receive-promise tables; parser roles own
pCID-owned app payload projection; apps own promise meaning, trust, economics,
and keep/break judgment. Analyzer gates now require
`persistent_session_*` events and fail if retained ACK artifacts lack request
parent links. The current gate also requires frame-sent/frame-received counts,
request-start/response-match counts, one terminal event per opened session, at
least one shutdown terminal reason, and explicit close/reopen events for the
trust-driven reconfiguration slice. Source: `DI-vopab`; `DI-mapop`; `TE-lubid`;
`DI-gazin`.

## Candidate Agent Work

- Alice discovers and uses a route to Dave, then sends actual traffic over it.
- Bob, Carol, Ellen, or Frank promise bounded forwarding only when local trust,
  capacity, pCID support, route constraints, and compensation are acceptable.
- Dave confirms reachability and receive willingness, then replies along either
  the reverse route or an independently discovered asymmetric route.
- Peggy promises routed WASM compute or validation work to Dave or Alice.
- Victor promises routed stdio compute or validation work to Alice or Dave.
- Alice asks direct peers for route-exclusion promises about Mallory, then sends
  only through peers whose local promises match Alice's constraints.

## Analyzer Targets

POC16 analyzer gates now cover the route slice, raw message retention,
multiarity specimens, parent-link specimen artifacts, COSE payload/proof
verification, COSE tamper rejection, normal-traffic route parent slots,
payload-parent route links, raw-message DAG traversal, route durability,
asymmetric response-path handling, reciprocal route credits, useful routed
Peggy/Victor compute, agent-accessible sparse CAS, bearer storage-token
incentives, encrypted-object CIDs, local CAS GC, and an explicit kernel-role
profile event. The newest gates add strict persistent-session lifecycle
accounting, route expiry/renewal/transit-exclusion pressure, local/primary,
peer/replica, missing-object, and untrusted-peer CAS retrieval outcomes, plus
signed-token issue/verify/expiry/replay events.
The newest parser-role gates require real parser-role process events,
`kernel_transport_v1` control traffic, positive parser-role delivery, local ACK
and backpressure promises, malformed-payload rejection, zero active artifacts for
retired `kernel_receive_v1`/shipping-step pCIDs, active runtime pCID count of
ten, and raw parser-flow artifact directions.
Remaining analyzer targets are:

- At least one forwarding non-commitment due to capacity, trust, pCID support,
  route constraints, or compensation mismatch.
- At least one route-exclusion promise used in route choice and later judged kept
  or broken from Alice's local events.
- Raw CAS objects for every sent, received, forwarded, ACKed, rejected,
  malformed, decision, monitor, WASM, and stdio artifact.
- Deeper route economics than one reciprocal route credit.
- True independent asymmetric return-route traffic, not only response-path terms
  and handling events.

## Planning Docs

- `docs/ROUTE-PROMISES.md` covers promise-based multi-hop forwarding, incentives,
  durability, asymmetric routes, and failure semantics.
- `docs/KERNEL-ROLES.md` covers the transport, app-interface, routing,
  local-resource, and event-retention role split.
- `docs/MESSAGE-SHAPES.md` covers pCID-owned arity, parent links, COSE specimens,
  transport-session proofs, and raw CAS/DAG review.
- `../../docs/thought-experiments/TE-vakah-poc16-secure-tokens-maps-encrypted-payloads.md`
  records the TE that motivates the POC16 design.
- `../../docs/thought-experiments/TE-lubid-poc16-persistent-multiplexed-sessions.md`
  records the TE for persistent multiplexed app/kernel and kernel/kernel
  sessions keyed by exact request message CIDs.
