# POC13 run narrative

POC13 is executable evidence, not a final PromiseGrid storage, compute,
economics, trust, token, evidence-report, key, or transport API. Source:
`DI-lupag`; `DI-nisaz`; `DI-kikoj`.

## Message sequence

1. Every container records `runtime_readiness_promised`; every local agent records
   `peer_readiness_observed` for the peer containers before application promises
   are sent.
2. Each agent loads container-local persisted trust history. Alice records
   `trust_driven_peer_choice` and `dynamic_peer_choice_from_persisted_trust`
   for Bob as primary storage, Frank as replica fallback, Carol as compute peer,
   and Dave/Grace as verification peers.
3. Alice sends Bob a low-credit `store_request`; Bob records
   `economics_price_refused` because the offer is below his local opportunity
   cost, then returns `price_refusal`.
4. Alice records the refusal and sends Bob two higher-credit `store_request`
   promises: one invoice object and one receipt object. The second object records
   `cas_multi_object_pressure`.
5. Bob verifies both content CIDs, reserves bounded local storage capacity for
   both, records credits earned, stores both byte strings, issues Alice
   capability tokens, and sends `store_acceptance` for each object.
6. Bob sends Frank `replicate_request` for both objects; Frank verifies bytes,
   records credits earned, stores both replicas, promises one-use TTL for each
   Frank-local replica token, sends Bob `replicate_acceptance`, and sends Alice
   `replica_available` for each object.
7. Bob records Frank's replica promise without promising on Frank's behalf.
8. Alice tries to redeem Bob's token with named transport outage variants:
   `container_stopped` and `network_partition`. Both route to closed local TCP
   ports, producing `tcp_message_send_failed` and `primary_storage_unavailable`
   evidence; Bob does not send a cooperative unavailable reply.
9. Alice receives Frank's `replica_available`, stores Frank's token locally,
   chooses Frank because Bob is unreachable and Frank is trusted enough, sends
   Frank `replica_serve_request` with the token, and records
   `replica_recovery_requested`.
10. Frank redeems his own replica token, serves the replica bytes, then records
    TTL expiry and revocation. Alice asks for renewal and records the renewed
    token. Alice verifies each content CID before recording
    `replica_recovery_succeeded` and `cas_bytes_retrieved`.
11. Alice asks Dave for a cache lookup before compute; Dave records
    `compute_cache_miss` and Alice records `compute_cache_miss_observed`.
12. Alice sends Carol `compute_request` with function bytes, input bytes, context
    CID, and credit offer. Carol records bounded compute capacity and credits
    earned, then asks Ellen for explicit context.
13. Ellen returns the context object; Carol recomputes `fibonacci(n)` from the
    payload-provided function/input/context bytes, records `compute_function_executed`,
    sends Alice one hash-valid but semantically wrong `compute_result`, then
    sends Alice the correct result and sends the correct result to Dave for cache
    checkpoint evidence.
14. Alice rejects the wrong result by local recomputation, withholds payment,
    and asks Dave and Grace to verify the same signed payload material.
15. Dave and Grace reject the wrong result by recomputation and return
    `evidence_report` promises under `evidence_report_v1`.
16. Alice verifies Carol's later correct result locally, spends the agreed
    compute credits, requests a second payload-provided `sum(values)` function,
    and asks Dave and Grace to verify the correct Fibonacci result.
17. Dave records `compute_cache_checkpointed` and `compute_cache_hit` for exact
    result tuples, then sends cache responses that Alice records as
    `compute_cache_reused`. Carol later executes the alternate sum function and
    records `compute_alternate_function_executed`.
18. Dave verifies correct results. Grace deliberately sends disagreement reports
    for correct results; Alice records `compute_verifier_disagreement` and
    `compute_disagreement_resolved_locally` using her own recompute plus Dave's
    report rather than treating Grace as an authority.
19. Mallory sends Grace corrupt CAS bytes under a valid-looking content CID;
    Grace rejects the bytes, lowers local trust in Mallory, records corrupt-byte
    evidence, and later records Mallory's voluntary repair promise without
    treating it as proof of future behavior.
20. Mallory also sends an unknown-pCID promise and an unsupported storage
    variant; Grace records `unknown_pcid_not_promised` and
    `promise_variant_not_promised` rather than crashing or treating either as an
    authority failure.
21. Mallory sends a bad-proof envelope; Grace records `bad_proof_rejected`.
    Mallory also sends a `key_rotation_promise` that Grace records as a future
    promise, not proof that Mallory's future key will be trustworthy.
22. Mallory sends Carol a competing compute request marked as capacity pressure;
    Carol records `economics_capacity_refused` as a local non-commitment.
23. Each container records `runtime_done_promised` only after local agents finish
    and the runtime observes an idle period with no active TCP handlers.

## Agent incentives

- **Alice** wants durable storage and trustworthy compute. She spends credits
  only after peers make acceptable promises, uses Bob first because persisted
  local trust favors him, falls back to Frank when Bob's TCP path fails, resolves
  verifier disagreement locally, and benefits from Dave's cache reuse.
- **Bob** wants to earn storage credits without overpromising. He refuses low
  price offers, accepts bounded capacity, issues an issuer-local serve token,
  and asks Frank for replication to improve Alice's recovery path.
- **Frank** earns replication value by keeping Bob's replica promise and later
  serving Alice when Bob is unreachable, but only after Alice presents Frank's
  own replica token. Frank limits exposure by expiring, revoking, and renewing
  tokens explicitly.
- **Carol** earns compute credits by executing only CID-bound function/input/
  context material. When she sends one bad result, Alice withholds payment and
  verification peers record broken evidence; she earns only after the correct
  result verifies. Carol also proves she is not hard-coded to Fibonacci by
  executing a separate sum function.
- **Ellen** earns relationship value by promising explicit context objects so
  impure compute becomes auditable and replayable.
- **Dave** earns relationship value by caching exact compute tuples and sending
  evidence-report promises about Alice's results without claiming global truth.
  Cache hits let Alice reuse prior exact results.
- **Grace** earns relationship value by independently verifying compute, recording
  corrupt-byte evidence from Mallory, declining unknown or unsupported promises
  locally, rejecting bad proofs, and showing that verifier disagreement is just
  local evidence Alice must resolve.
- **Mallory** applies adversarial pressure. Her corrupt bytes reduce Grace's
  local trust; her later repair promise is only future evidence to watch, not an
  authority override.

## Analyzer report

The analyzer now reports raw event counts plus `score_report`. The score is a
bounded POC fitness summary across transport, storage, compute, economics,
trust, verification, and replica recovery. It also gates on `evidence_report_v1`,
unknown-pCID non-commitment, unsupported-variant non-commitment, TCP send
failure, named outage variants, token TTL/expiry/revocation/renewal, multi-object
CAS pressure, cache miss/hit/reuse, alternate compute functions, verifier
disagreement resolution, persisted trust, bad-proof rejection, key rotation
promises, replica token redemption, bad-result rejection, and capacity refusal.
It is a regression signal for POC13; it is not a general PromiseGrid quality
score.
