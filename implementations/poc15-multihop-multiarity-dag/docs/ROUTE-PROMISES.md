# POC15 Route Promises

POC15 routing should be modeled as promises between neighboring agents, not as a
route command, permission check, policy engine, or global path authority. Source:
`DI-pamob`; `DI-podut`.

## Basic Model

A multi-hop send is a sequence of local exchanges:

1. Alice sends a route-setup promise to Bob because Alice locally trusts Bob
   enough for this purpose.
2. Bob decides whether Bob promises one forwarding hop to Carol.
3. Carol decides whether Carol promises one forwarding hop to Dave.
4. Dave receives the setup message and judges the sender, route terms, pCID, and
   payload locally.
5. Dave returns a reachability promise along a promised return path.
6. Alice sends actual traffic only after route-confirmation promises exist.

Each hop is a local promise by the current holder of the envelope. A failed hop
is not automatically an event that the final receiver broke a promise; it may be
local capacity refusal, missing pCID support, weak trust, expired reciprocal
promise, unavailable transport, or a broken forwarding promise by the hop that
actually promised forwarding.

## Lightning-Leader Route Setup

The POC15 route setup should test a lightning-leader pattern without importing
command/control vocabulary:

- Alice promises Bob that Alice will use, compensate, or reciprocate for a route
  to Dave if Bob can assemble a route satisfying Alice's route constraints.
- Bob may decline, counter, or promise only Bob's next hop.
- Bob forwards a compatible route-setup promise to Carol only if Bob locally
  judges the terms acceptable.
- Carol repeats the pattern until Dave receives the setup message.
- Dave returns a reachability promise that identifies the route, constraints,
  parent message CID, expiration, and receive willingness.

The setup message is not an order to route. It is an offer to form a chain of
voluntary forwarding promises.

## Incentives

Forwarding is not free. POC15 should test incentives as ordinary reciprocal
promises:

- **Reciprocal forwarding:** Bob forwards for Alice because Alice promises to
  forward for Bob later under bounded terms.
- **Relationship value:** Bob forwards because kept promises improve Bob's local
  relationship with Alice or Carol.
- **Bearer capability-token payment:** Alice pays Bob with a transferable token
  issued by an agent Bob locally values.
- **Non-transferable forwarding-capacity token:** Bob issues Alice a scoped token
  promising future forwarding capacity for a route class.
- **Stake or collateral:** Alice or Bob stakes promise tokens that can be locally
  judged broken or forfeited if the promised forwarding behavior fails.

Bearer tokens help when peers locally value them, but there is no central
exchange rate. Each peer decides whether the offered token is worth its own
capacity, risk, and opportunity cost.

## Route Durability

A route is durable only if each hop promises durability. POC15 should test:

- **One-shot route:** valid for one carried message or one request/response pair.
- **Bounded durable route:** valid until an expiration, byte limit, message count,
  pCID set, local trust threshold, or capacity limit is reached.
- **Revoked or not-renewed route:** a peer locally promises no further forwarding
  after expiration, trust decay, capacity pressure, or changed reciprocal terms.

Alice may cache a route, but the cache is a record of promises and constraints,
not a guaranteed circuit.

## Asymmetric Routes

The same promise process supports asymmetric routes:

- Dave can reply over the reverse route if every reverse hop promised that use.
- Dave can run independent route discovery back to Alice.
- Alice can include a return-route token or reply-route promise in the setup
  message.

Parent links should connect request and response DAGs even when the forward path
and return path differ.

## Route Exclusion

If Alice does not want Alice-owned traffic to transit Mallory, Alice has no
global way to inspect every possible path. The PromiseGrid-compatible approach is
to ask Alice's direct peers for bounded promises such as:

> Bob promises Alice that for this route class, Bob will not intentionally
> forward Alice's envelopes to Mallory or through a route candidate that Bob
> locally knows includes Mallory.

Alice then chooses whether Bob's promise and history are good enough. Bob still
owns Bob's choice. If Bob later breaks that promise, Alice records local
break-history event and reduces trust in Bob for future traffic. This is route
selection by local trust, not route enforcement.

## Forwarding Payload Shape

POC15 should add a narrow route pCID because route setup and forwarding have real
protocol structure: route ID, final target, next-hop promise, pCID constraints,
size/TTL limits, expiration, parent links, route-exclusion terms, compensation,
reciprocal promises, carried-message CID, and carried-message bytes.

The top-level semantic action remains `promise`. Route setup, route confirmation,
forwarding, refusal, compensation, and repair are pCID-owned payload semantics or
local event interpretations, not reusable top-level action kinds.

## Failure Semantics

Analyzer and runtime logs must distinguish:

- **No promise:** a peer declines, lacks capacity, lacks pCID support, or does not
  trust the terms. No trust penalty by default.
- **Could not send:** a local transport or resource failure. Do not blame the peer
  automatically.
- **Promise broken:** a peer promised a hop and then failed under the promised
  conditions after the deadline or outcome rule applies.
- **Malformed or tampered:** reject exact bytes and reduce trust in the peer that
  sent those bytes, not automatically the original author.
- **Unknown outcome:** timeout or missing downstream information remains unresolved
  until the local promise's deadline passes.

## Events To Record

- `route_setup_promise_made`: a direct peer receives or makes a conditional route
  setup promise.
- `route_forward_promise_made`: a direct peer promises one forwarding hop.
- `route_reachability_promised`: the distant peer promises to receive matching
  route traffic under bounded terms.
- `route_forward_promise_kept`: the next hop receives the exact envelope or a
  pCID-defined forwarding wrapper.
- `route_forward_not_promised`: the peer does not currently promise the hop.
- `route_payment_promised`: a peer accepts compensation or reciprocal terms.
- `route_exclusion_promise_made`: a peer promises not to forward through Mallory
  for Alice's named route class.
- `route_transit_exclusion_peer_promised`: each forwarding peer promises its own
  hop will not intentionally use Mallory as transit for Alice's named route.
- `route_exclusion_used_in_choice`: Alice chooses a route because local peer
  promises match Alice's local constraints.
- `route_exclusion_broken_observed`: Alice records local events that a chosen
  peer broke a route-exclusion promise.
- `route_lifetime_exhausted`: Alice records that the bounded route lifetime has
  been consumed.
- `route_expired_message_not_sent`: Alice locally declines to send after the
  route's promised lifetime is exhausted.
- `route_renewal_requested` / `route_renewal_confirmed`: Alice asks for and
  receives a fresh neighboring route promise before relying on the path again.

These names are planning names, not locked wire action kinds. They should remain
events or payload meanings unless POC15 locks a specific pCID-owned payload.
