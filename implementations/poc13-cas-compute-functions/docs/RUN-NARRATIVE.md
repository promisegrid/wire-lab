# POC13 run narrative

POC13 is executable evidence, not a final PromiseGrid storage, compute,
economics, trust, token, evidence-report, or transport API. Source: `DI-lupag`;
`DI-nisaz`.

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
6. Bob sends Frank `replicate_request`; Frank verifies the bytes, records
   credits earned, stores the replica, issues Alice a Frank-local replica token,
   sends Bob `replicate_acceptance`, and sends Alice `replica_available`.
7. Bob records Frank's replica promise without promising on Frank's behalf.
8. Alice tries to redeem Bob's token with `serve_request`, but this scenario
   routes that one send to a closed TCP port. Alice records `tcp_message_send_failed`
   and `primary_storage_unavailable`; Bob does not send a cooperative
   unavailable reply.
9. Alice receives Frank's `replica_available`, stores Frank's token locally,
   chooses Frank because Bob is unreachable and Frank is trusted enough, sends
   Frank `replica_serve_request` with the token, and records
   `replica_recovery_requested`.
10. Frank redeems his own replica token, serves the replica bytes, and Alice
    verifies the content CID before recording `replica_recovery_succeeded` and
    `cas_bytes_retrieved`.
11. Alice sends Carol `compute_request` with function bytes, input bytes, context
    CID, and credit offer. Carol records bounded compute capacity and credits
    earned, then asks Ellen for explicit context.
12. Ellen returns the context object; Carol recomputes `fibonacci(n)` from the
    payload-provided function/input/context bytes, records `compute_function_executed`,
    sends Alice one hash-valid but semantically wrong `compute_result`, then
    sends Alice the correct result and sends the correct result to Dave for cache
    checkpoint evidence.
13. Alice rejects the wrong result by local recomputation, withholds payment,
    and asks Dave and Grace to verify the same signed payload material.
14. Dave and Grace reject the wrong result by recomputation and return
    `evidence_report` promises under `evidence_report_v1`.
15. Alice verifies Carol's later correct result locally, spends the agreed
    compute credits, and asks Dave and Grace to verify it.
16. Dave records `compute_cache_checkpointed` for the exact
    protocol/function/input/context/result tuple. Dave and Grace verify the
    correct result and return `evidence_report` promises; Alice records
    `compute_verification_received`.
17. Mallory sends Grace corrupt CAS bytes under a valid-looking content CID;
    Grace rejects the bytes, lowers local trust in Mallory, records corrupt-byte
    evidence, and later records Mallory's voluntary repair promise without
    treating it as proof of future behavior.
18. Mallory also sends an unknown-pCID promise and an unsupported storage
    variant; Grace records `unknown_pcid_not_promised` and
    `promise_variant_not_promised` rather than crashing or treating either as an
    authority failure.
19. Mallory sends Carol a competing compute request marked as capacity pressure;
    Carol records `economics_capacity_refused` as a local non-commitment.
20. Each container records `runtime_done_promised` only after local agents finish
    and the runtime observes an idle period with no active TCP handlers.

## Agent incentives

- **Alice** wants durable storage and trustworthy compute. She spends credits
  only after peers make acceptable promises, uses Bob first because of local
  trust, and falls back to Frank when Bob does not currently promise serving.
- **Bob** wants to earn storage credits without overpromising. He refuses low
  price offers, accepts bounded capacity, issues an issuer-local serve token,
  and asks Frank for replication to improve Alice's recovery path.
- **Frank** earns replication value by keeping Bob's replica promise and later
  serving Alice when Bob is unreachable, but only after Alice presents Frank's
  own replica token.
- **Carol** earns compute credits by executing only CID-bound function/input/
  context material. When she sends one bad result, Alice withholds payment and
  verification peers record broken evidence; she earns only after the correct
  result verifies.
- **Ellen** earns relationship value by promising explicit context objects so
  impure compute becomes auditable and replayable.
- **Dave** earns relationship value by caching exact compute tuples and sending
  evidence-report promises about Alice's results without claiming global truth.
- **Grace** earns relationship value by independently verifying compute, recording
  corrupt-byte evidence from Mallory, and declining unknown or unsupported
  promises locally.
- **Mallory** applies adversarial pressure. Her corrupt bytes reduce Grace's
  local trust; her later repair promise is only future evidence to watch, not an
  authority override.

## Analyzer report

The analyzer now reports raw event counts plus `score_report`. The score is a
bounded POC fitness summary across transport, storage, compute, economics,
trust, verification, and replica recovery. It also gates on `evidence_report_v1`,
unknown-pCID non-commitment, unsupported-variant non-commitment, TCP send
failure, replica token redemption, bad-result rejection, and capacity refusal.
It is a regression signal for POC13; it is not a general PromiseGrid quality
score.
