# poc15-multihop-multiarity-dag

`poc15-multihop-multiarity-dag` is the planned successor to `poc14-wasm`.
It is executable as a mechanically renamed POC14 baseline so POC15 additions can
be measured against a known regression floor. The purpose is to preserve POC14
while adding real multi-hop forwarding, pCID-owned slot-vector variety,
raw-message CAS/DAG review, agent-accessible sparse CAS, wire-visible parent
links, COSE specimens, useful routed WASM/stdio work, promise-correct route
exclusion, route economics, and explicit non-monolithic kernel roles. Source:
`DI-pamob`; `DI-podut`; `DI-lutuv`; `DI-manul`.

## Current Executable State

The current executable scaffold preserves inherited POC14 app/kernel, shipping,
CAS, compute, WASM, stdio, event-collector, and analyzer behavior under a
separate POC15 module path and run root. It now also includes the first
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
artifacts for operator review through an observer-only artifact stream. A
validated clean run for this update reconstructs a reachable 222-node
exact-message DAG from 413
artifact observations, with parent-link coverage in both envelope slots and
pCID-defined payload fields. Source: `DI-lutuv`; `DI-lihir`; `DI-kohuj`;
`DI-tuhop`; `DI-mosat`.

The current CAS/DAG slice now also gives each app its own local sparse CAS view.
Agents may store exact messages, local state bytes, encrypted blobs named by
ciphertext CID, or peer-served content. Those stores are intentionally
incomplete: a local message-DAG index may record missing parents, peer storage is
voluntary, bearer storage tokens can compensate a holder for retrieving bytes
from the issuer, and local GC retains paid/pinned/encrypted objects while
removing pressure-tagged temporary objects. Source: `DI-manul`.

## Superset Requirement

POC15 should preserve POC14 unless a later DI explicitly authorizes a scoped
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

## POC15 Additions

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
5. **Route durability and asymmetry.** POC15 should test one-shot routes,
   bounded durable routes, reverse-route replies, independent return-route
   discovery, and return-route tokens.
6. **pCID-owned multiarity.** POC15 should test pCIDs that define different
   arity, slot order, parent-link location, proof location, and COSE usage.
7. **Raw message review.** Every raw envelope artifact should be retained in
   run-scoped CAS for later review, while valid message parent links form a
   wire-visible causal DAG.
8. **Useful routed WASM/stdio work.** Peggy and Victor should do valuable work
   for other agents over routed paths, not merely prove adapter plumbing exists.
9. **Kernel as role collection.** POC15 should name transport, app-interface,
   routing, local-resource, and event-retention roles explicitly even when the
   Docker runtime collapses some roles into one process.
10. **Agent-accessible sparse CAS.** Each agent should be able to maintain its
    own incomplete CAS and optional local message DAG, or rely on peer CAS
    promises paid through credits or bearer storage tokens.

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

The executable POC15 run now stores exact envelope bytes for operator review.
Apps do not mount or read the observer volume; instead, each app emits a one-way
`message_artifact` record to stdout, the local supervisor forwards it to the
observer-only collector, and the collector writes:

- `/run/poc15/<run_id>/message-cas/<exact_sha256>.cbor` — binary CBOR envelope
  bytes exactly as sent, received, or acknowledged.
- `/run/poc15/<run_id>/message-dag.jsonl` — one JSON index row per artifact,
  including observer, direction, peer, pCID name, exact hash, optional parent
  exact hash, promise meaning, and relative artifact path.

`poc15-analyze` validates the index by reading each `.cbor` artifact and
recomputing its exact SHA-256 before the clean regression can pass. It traverses
the unique exact-message DAG rather than double-counting duplicate observations
of the same message, then fails the clean regression if a valid parent link is
missing or unreachable. Source: `DI-tuhop`; `DI-kohuj`.

`poc15-cbor-diag` is included in the collector image for operator inspection of
the retained artifacts:

```sh
docker compose run --rm --entrypoint /usr/local/bin/poc15-cbor-diag \
  event-collector -hash <exact_sha256>
```

The tool reads `message-dag.jsonl`, opens the matching
`message-cas/<exact_sha256>.cbor` artifact, and prints CBOR diagnostic notation
with nested payload/proof byte strings expanded when they contain valid CBOR.
It is read-only and does not affect app/kernel behavior. Source: `DI-bapif`.

## Agent-Accessible Sparse CAS

The observer-only `message-cas/` archive is for operator review. It is not the
same thing as an app's local CAS. Each app now keeps run-scoped CAS metadata in
its own `stores/<agent>/durable-state.json` file, and that metadata describes
objects the agent chose to retain:

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
store, global retention rule, or global trust decision. Source: `DI-manul`.

## Multiarity Specimens

POC15 now emits and analyzes pCIDs that define these slot-vector specimens:

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

POC15 analyzer gates now cover the route slice, raw message retention,
multiarity specimens, parent-link specimen artifacts, COSE payload/proof
verification, COSE tamper rejection, normal-traffic route parent slots,
payload-parent route links, raw-message DAG traversal, route durability,
asymmetric response-path handling, reciprocal route credits, useful routed
Peggy/Victor compute, agent-accessible sparse CAS, bearer storage-token
incentives, encrypted-object CIDs, local CAS GC, and an explicit kernel-role
profile event.
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
- `../../docs/thought-experiments/TE-vakah-poc15-multihop-multiarity-dag.md`
  records the TE that motivates the POC15 design.
