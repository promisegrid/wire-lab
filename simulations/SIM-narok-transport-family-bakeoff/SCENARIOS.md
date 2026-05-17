# Transport Family Bakeoff Scenarios

These scenarios are evidence for `TODO-sinuv`. They are not transport-family
decisions and not frozen PromiseGrid transport specs. Source: `DI-pukap`.

## Scenario Matrix

| Scenario | Setup | What to test | Decision pressure |
|---|---|---|---|
| Token rotates among writers | Alice, Bob, and Carol can publish only while holding a ring token. | Whether token possession is a transport-level ordering promise, a group/session policy, or both. | If token state is load-bearing for delivery semantics, a ring transport deserves its own specimen. |
| Per-hop authorization failure | Bob receives a ring message but is not authorized to forward it to Carol. | Whether authorization failure breaks the ring, skips a hop, records refusal, or reconfigures membership. | Ring semantics need a failure model before they can be compared with gossip. |
| Cluster boundary crossing | Alice sends from cluster A to Dave in cluster B through cluster representatives. | Whether clusters are addressable principals, routing hints, policy scopes, or aggregation boundaries. | Cluster-of-clusters transport needs explicit promises at each boundary. |
| Gossip convergence | Alice emits a message; Bob, Carol, and Dave learn it through epidemic propagation with duplicate paths. | Whether convergence evidence is enough without total ordering or delivery promises. | Gossip is attractive only if duplicate suppression and missing-object repair stay explainable. |
| Large-N receipts | One thousand peers may acknowledge possession or processing of the same object. | Whether receipts become vectors, compact summaries, per-peer promise records, or separate content-addressed proofs. | Receipt scale may be a cross-family extension rather than a transport family by itself. |

## Expected Outputs

- A decision map showing which future transport questions are independent and
  which should stay coupled.
- Evidence for whether `TODO-sinuv` should split into separate TE/DR owners for
  ring, cluster, gossip, and receipts-at-scale.
