# poc15-multihop-multiarity-dag

`poc15-multihop-multiarity-dag` is the planned successor to `poc14-wasm`.
It is not executable yet. The purpose is to keep POC14 as the regression
baseline while adding real multi-hop forwarding, pCID-owned slot-vector variety,
raw-message CAS/DAG review, wire-visible parent links, COSE specimens, useful
routed WASM/stdio work, promise-correct route exclusion, route economics, and
explicit non-monolithic kernel roles. Source: `DI-pamob`; `DI-podut`.

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
7. **Raw message review.** Every raw artifact should be retained in run-scoped
   CAS for later review, while valid message parent links form a wire-visible
   causal DAG.
8. **Useful routed WASM/stdio work.** Peggy and Victor should do valuable work
   for other agents over routed paths, not merely prove adapter plumbing exists.
9. **Kernel as role collection.** POC15 should name transport, app-boundary,
   routing, local-resource, and event-retention roles explicitly even when the
   Docker runtime collapses some roles into one process.

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

## Multiarity Specimens

POC15 should add pCIDs that define these slot-vector specimens:

- `grid([42(pCID), payload])` for transport/session-auth-only pressure.
- `grid([42(pCID), payload, proof])` for the current signed-message specimen.
- `grid([42(pCID), parents, payload, proof])` for envelope parents before body.
- `grid([42(pCID), payload, parents, proof])` for body before envelope parents.
- `grid([42(pCID), payload, proof])` where parent links live inside the
  pCID-owned payload.
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

POC15 should add analyzer gates for:

- At least one successful multi-hop route setup and forwarded actual message.
- At least one forwarding non-commitment due to capacity, trust, pCID support,
  route constraints, or compensation mismatch.
- At least one one-shot route and one bounded durable route.
- At least one asymmetric reply route.
- At least one incentive, payment, or reciprocal-token exchange.
- At least one route-exclusion promise used in route choice and later judged kept
  or broken from Alice's local events.
- At least one useful routed Peggy work item and one useful routed Victor work
  item.
- Raw CAS objects for every sent, received, forwarded, ACKed, rejected,
  malformed, decision, monitor, WASM, and stdio artifact.
- Valid parent DAG reconstruction for envelope-parent and payload-parent
  specimens.
- COSE payload and COSE proof validation plus tamper rejection.
- Parser dispatch by pCID-owned arity and slot semantics.
- Explicit event that kernel roles are separated or intentionally collapsed for
  a named runtime.

## Planning Docs

- `docs/ROUTE-PROMISES.md` covers promise-based multi-hop forwarding, incentives,
  durability, asymmetric routes, and failure semantics.
- `docs/KERNEL-ROLES.md` covers the transport, app-boundary, routing,
  local-resource, and event-retention role split.
- `docs/MESSAGE-SHAPES.md` covers pCID-owned arity, parent links, COSE specimens,
  transport-session proofs, and raw CAS/DAG review.
- `../../docs/thought-experiments/TE-vakah-poc15-multihop-multiarity-dag.md`
  records the TE that motivates the POC15 design.
