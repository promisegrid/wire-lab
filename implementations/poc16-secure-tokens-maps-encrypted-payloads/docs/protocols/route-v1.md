# route_v1

## Status

Active POC16 sparse-network routing protocol. The exact embedded Markdown bytes
are hashed to derive this pCID. Source: `DI-bitug`; `DI-lihir`.

## Abstract

`route_v1` carries voluntary promises about bounded forwarding through sparse
peer relationships. Each hop promises only its own forwarding behavior for a
route, route lifetime, parent-link handling, and reciprocal compensation. It is
not OSPF, not a global route table, and not an obligation imposed by the sender.

## pCID and envelope

The active shape is:

```text
grid([42(pCID), payload, proof])
```

Payload addressing belongs to this pCID. Slot 0 selects the route parser; it does
not identify the route destination.

## Promise Theory model

Every hop is an autonomous promiser. Alice may promise to use and compensate a
route if Bob promises to forward; Bob may promise to forward if Carol promises to
forward; each promise is local and revocable according to its stated terms. Route
failure lowers trust only where a local agent judges that a promise was broken.

## Payload grammar

The payload is the pCID-owned map-body profile:

```text
payload = [
  promiser: text,
  promisee: text,
  promise_about: text,
  state: [outcome: text, promise_text: text, reason: text, turn: text],
  body: {detail_key: text => detail_value: text, ...}
]
```

All core slots are REQUIRED. `body` contains protocol-owned text key/value
details in a nested CBOR map namespace. A parser MUST reject non-arrays, wrong
array lengths, non-text core fields, non-map bodies, duplicate body keys,
reserved/core body keys, non-text body keys or values, or trailing CBOR bytes.

Common `promise_about` values are `route_probe`, `forward_if_next_promises`,
`route_confirmation`, `route_use`, `route_renewal`, `route_failure`, and
`route_exclusion`. Common details are `route_id`, `next_peer`, `final_peer`,
`ttl_turns`, `max_hops`, `credit_offer`, `credit_token`, `parent_cid`,
`excluded_peer`, and `failure_reason`.

## Sender behavior

A sender SHOULD probe before sending valuable payloads over a route. A forwarding
promise SHOULD state lifetime, next peer, compensation, and whether parent links
are preserved. A route exclusion is a promise preference by the sender; it is
credible only to the extent that direct peers voluntarily promise not to forward
through the excluded peer.

## Receiver and parser behavior

A receiver MAY promise to forward, decline silently, send a non-commitment, or
counter with different terms. A receiver MUST NOT forward if doing so would break
its own local promises about route exclusion, capacity, privacy, or compensation.

## Protocol state machine

```text
[no route]
    | route_probe
    v
[hop negotiating] --next hop promises--> [route confirmed]
      | no promise / bad terms                  |
      v                                        | route_use
[route not promised]                           v
                                         [forwarding active]
                                               |
                             expiry/failure/withdrawal
                                               v
                                         [route inactive]
```

Asymmetric routes are represented by separate route promises with distinct
`route_id` values and parent links.

## State, CAS, DAG, and retention

Route probes, confirmations, and use messages SHOULD parent-link to prior route
messages so agents can reconstruct why a route was locally trusted. No global
route DAG is assumed.

## Security considerations

Routes can leak traffic patterns. A sender SHOULD avoid putting sensitive payloads
on a route unless it has sufficient local trust in every promised hop or unless
the payload is encrypted so transit peers cannot read it. Compensation tokens
must be verified before being treated as payment promises.

## Interoperability notes

Production implementations may use this protocol for manual sparse routes,
lightning-style route discovery, or relay-market experiments. It deliberately
avoids central route authorities.

## Examples

```text
grid([42(pCID),
  ["bob", "alice", "forward_if_next_promises",
    ["kept", "I promise to forward route r7 traffic to Carol for three turns if Carol promises the next hop.",
     "Alice offered relay credit", "turn-12"],
    {"route_id": "r7", "next_peer": "carol", "final_peer": "grace",
     "ttl_turns": "3", "credit_offer": "alice-relay-credit:5"}
  ], proof
])
```
