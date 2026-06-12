# POC15 Route Promises

POC15 routing should be modeled as promises between neighboring agents, not as a
route command, permission check, policy engine, or global path authority. Source:
`DI-pamob`.

## Basic Model

A multi-hop send is a sequence of local exchanges:

1. Alice sends a signed envelope to Bob because Alice locally trusts Bob enough
   for this purpose.
2. Bob decides whether Bob promises one forwarding hop to Carol.
3. Carol decides whether Carol promises one forwarding hop to Dave.
4. Dave receives the envelope and judges the sender, route evidence, pCID, and
   payload locally.

Each hop is a local promise by the current holder of the envelope. A failed hop
is not automatically evidence that the final receiver broke a promise; it may be
local capacity refusal, missing pCID support, weak trust, expired reciprocal
promise, unavailable transport, or a broken forwarding promise by the hop that
actually promised forwarding.

## Route Exclusion

If Alice does not want Alice-owned traffic to transit Mallory, Alice has no
global way to inspect every possible path. The PromiseGrid-compatible approach is
to ask Alice's direct peers for bounded promises such as:

> Bob promises Alice that for this route class, Bob will not intentionally
> forward Alice's envelopes to Mallory or through a route candidate that Bob
> locally knows includes Mallory.

Alice then chooses whether Bob's promise and history are good enough. Bob still
owns Bob's choice. If Bob later breaks that promise, Alice records local
break-history evidence and reduces trust in Bob for future traffic. This is route
selection by local trust, not route enforcement.

## Forwarding Payload Shape

The first executable POC15 slice should keep the top-level action as `promise`
and should avoid inventing workflow verbs. There are two plausible payload
directions:

- Use `relationship_v1` promise payloads for the first route-promise experiment,
  because the semantics are relationship/trust/evidence-heavy and the exact
  forwarding mechanics can remain implementation-local.
- Add a narrow future `route_forwarding_v1` pCID only after a TE/DI settles why
  forwarding needs its own pCID-owned payload shape rather than relationship
  promises plus local transport mechanics.

The analyzer should treat premature generic route actions as drift unless a later
DI explicitly locks them.

## Evidence To Record

- `route_forward_promise_made`: a direct peer promises one forwarding hop.
- `route_forward_promise_kept`: the next hop receives the exact envelope or a
  pCID-defined forwarding wrapper.
- `route_forward_not_promised`: the peer does not currently promise the hop.
- `route_exclusion_promise_made`: a peer promises not to forward through Mallory
  for Alice's named route class.
- `route_exclusion_used_in_choice`: Alice chooses a route because local peer
  promises match Alice's local constraints.
- `route_exclusion_broken_observed`: Alice records local evidence that a chosen
  peer broke a route-exclusion promise.

These names are planning names, not locked wire action kinds. They should remain
events or payload meanings unless POC15 locks a specific pCID-owned payload.
