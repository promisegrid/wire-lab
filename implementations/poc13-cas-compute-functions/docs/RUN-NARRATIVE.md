# POC13 run narrative

POC13 is executable evidence, not a final PromiseGrid storage, compute,
economics, trust, token, or transport API. Source: `DI-lupag`.

## Message sequence

1. Every container records `runtime_readiness_promised`; every local agent records
   `peer_readiness_observed` for the peer containers before application promises
   are sent.
2. Alice records `trust_driven_peer_choice` for Bob as primary storage, Frank as
   replica fallback, Carol as compute peer, and Dave/Grace as verification peers.
3. Alice sends Bob a low-credit `store_request`; Bob records
   `economics_price_refused` because the offer is below his local opportunity
   cost, then returns `price_refusal`.
4. Alice records the refusal and sends Bob a higher-credit `store_request` with
   invoice bytes, content CID, and a request for a serve-once capability token.
5. Bob verifies the content CID, reserves bounded local storage capacity, records
   credits earned, stores the bytes, issues Alice a capability token, and sends
   `store_acceptance` to Alice.
6. Bob sends Frank `replicate_request`; Frank verifies the bytes, reserves local
   replica capacity, records credits earned, stores the replica, and sends Bob
   `replicate_acceptance`.
7. Bob records Frank's replica promise and sends Alice `replica_available`.
8. Alice redeems Bob's token with `serve_request`, but Bob records
   `primary_storage_unavailable` and sends Alice `primary_unavailable_notice`
   instead of serving bytes.
9. Alice receives Bob's notice, chooses Frank because Bob is unavailable and
   Frank is trusted enough, sends Frank `replica_serve_request`, and records
   `replica_recovery_requested`.
10. Frank serves the replica bytes; Alice verifies the content CID and records
    `replica_recovery_succeeded` and `cas_bytes_retrieved`.
11. Alice sends Carol `compute_request` with function bytes, input bytes, context
    CID, and credit offer. Carol records bounded compute capacity and credits
    earned, then asks Ellen for explicit context.
12. Ellen returns the context object; Carol recomputes `fibonacci(n)` from the
    payload-provided function/input/context bytes, records `compute_function_executed`,
    sends the result to Alice, and sends the same result to Dave for cache
    checkpoint evidence.
13. Dave records `compute_cache_checkpointed` for the exact
    protocol/function/input/context/result tuple.
14. Alice verifies Carol's result locally by recomputing the payload-provided
    function and records `compute_result_locally_verified`.
15. Alice asks Dave and Grace to verify the same compute result. Each peer
    recomputes from the signed payload material and returns
    `compute_verification_result`; Alice records `compute_verification_received`.
16. Mallory sends Grace corrupt CAS bytes under a valid-looking content CID;
    Grace rejects the bytes, lowers local trust in Mallory, records corrupt-byte
    evidence, and later records Mallory's voluntary repair promise without
    treating it as proof of future behavior.
17. Each container records `runtime_done_promised` only after local agents finish
    and the runtime observes an idle period with no active TCP handlers.

## Agent incentives

- **Alice** wants durable storage and trustworthy compute. She spends credits
  only after peers make acceptable promises, uses Bob first because of local
  trust, and falls back to Frank when Bob does not currently promise serving.
- **Bob** wants to earn storage credits without overpromising. He refuses low
  price offers, accepts bounded capacity, issues an issuer-local serve token,
  and asks Frank for replication to improve Alice's recovery path.
- **Frank** earns replication value by keeping Bob's replica promise and later
  serving Alice when Bob is unavailable.
- **Carol** earns compute credits by executing only CID-bound function/input/
  context material and by returning result bytes named by result CID.
- **Ellen** earns relationship value by promising explicit context objects so
  impure compute becomes auditable and replayable.
- **Dave** earns relationship value by caching exact compute tuples and verifying
  Alice's result without claiming global truth.
- **Grace** earns relationship value by independently verifying compute and by
  recording corrupt-byte evidence from Mallory.
- **Mallory** applies adversarial pressure. Her corrupt bytes reduce Grace's
  local trust; her later repair promise is only future evidence to watch, not an
  authority override.

## Analyzer report

The analyzer now reports raw event counts plus `score_report`. The score is a
bounded POC fitness summary across transport, storage, compute, economics,
trust, verification, and replica recovery. It is a regression signal for POC13;
it is not a general PromiseGrid quality score.
