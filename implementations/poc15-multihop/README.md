# poc15-multihop

`poc15-multihop` is the planned successor to `poc14-wasm`. It is not executable
yet. The purpose is to keep POC14 as the regression baseline while adding real
multi-hop forwarding, useful routed WASM/stdio work, promise-correct route
exclusion, and explicit non-monolithic kernel roles. Source: `DI-pamob`.

## Superset Requirement

POC15 should preserve POC14 unless a later DI explicitly authorizes a scoped
exception:

- POC11-style autonomous sparse-mesh relationship and economics pressure.
- POC12-style local kernel/app/device workflow and shipping agents.
- POC13-style CAS storage, CID compute, verifier disagreement, replica recovery,
  token lifecycle, retention/GC, backpressure, rate-limit, replay protection,
  bounded trust, and dynamic direct TCP relationship event.
- POC14-style WASM process agents, stdio-only subprocess agents, decentralized
  monitoring signals, hard local distrust, route-exclusion scenarios, and
  pCID-owned payload migration event.

## POC15 Additions

1. **Real multi-hop forwarding.** A sender should be able to send an envelope to
   a direct peer that voluntarily promises to forward it to another direct peer,
   eventually reaching the target if each hop locally promises the next step.
2. **Route promises, not route authority.** A forwarding path is a chain of
   voluntary peer promises. No node commands another node to forward, and no node
   claims a global route truth.
3. **Route exclusion by peer promises.** If Alice does not want traffic to transit
   Mallory, Alice cannot prove the entire network path by inspection. Alice can
   select peers whose direct promises say they will not forward Alice traffic
   through Mallory, then judge their later keep/break event locally.
4. **Useful routed WASM/stdio work.** Peggy and Victor should do work that is
   valuable to other agents over routed paths, not merely prove that WASM and
   stdio boundaries exist.
5. **Kernel as role collection.** POC15 should split transport, app-boundary,
   routing, and local-resource roles into explicit processes or objects where
   practical. A tiny runtime may still collapse those roles into one process or
   firmware image, but the design vocabulary should stay role-based.

## Candidate Agent Work

- Peggy promises WASM module-validation event to Dave over a route chosen by
  local trust and route-exclusion promises.
- Victor promises stdio subprocess round-trip event to Alice or Dave over a
  route that may require a relay.
- Frank or Ellen promises one bounded forwarding hop for selected pCIDs only
  when local trust, capacity, and reciprocal economics are acceptable.
- Alice asks direct peers for route-exclusion promises about Mallory, then sends
  only through peers whose local promises match Alice's constraints.

## Analyzer Targets

POC15 should add analyzer gates for:

- At least one successful multi-hop forwarded envelope.
- At least one forwarding non-commitment when capacity, trust, pCID, or route
  constraints are not locally promised.
- At least one route-exclusion promise used in route choice.
- At least one route-exclusion promise later evaluated as kept or broken from
  Alice's local events.
- At least one useful routed Peggy work item and one useful routed Victor work
  item.
- Explicit event that the kernel roles are separated or intentionally
  collapsed for a named runtime.

## Open Questions

- What pCID should own forwarding payload shape: a new narrow routing pCID or a
  relationship-level promise in the first executable slice?
- What exact route event does a forwarder return without becoming an
  authority over downstream peers?
- How much path disclosure is useful before it leaks too much relationship or
  topology information?
- Which kernel roles must be processes in Docker and which can remain objects in
  the same process for the first POC15 implementation?

## Planning Docs

- `docs/ROUTE-PROMISES.md` covers multi-hop forwarding and route exclusion by
  peer promises.
- `docs/KERNEL-ROLES.md` covers the transport, app-boundary, routing,
  local-resource, and event role split.
