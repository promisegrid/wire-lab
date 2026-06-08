# POC13 clean-run analyzer output, 2026-06-08

This file archives the exact analyzer JSON captured after a clean Docker run of
`implementations/poc13-cas-compute-functions/scripts/run-clean.sh`. The analyzer
was rerun with:

```sh
docker compose -f implementations/poc13-cas-compute-functions/compose.yaml run --rm --entrypoint /usr/local/bin/poc13-analyze alice-bob /run/poc13/poc13-demo
```

Source: `DI-nisaz`.

```json
{
  "run_dir": "/run/poc13/poc13-demo/run",
  "total_events": 264,
  "event_counts": {
    "agent_waiting_for_peer_promises": 6,
    "app_receive_promise_registered": 8,
    "capability_token_issued": 1,
    "capability_token_received": 1,
    "cas_bytes_retrieved": 1,
    "cas_bytes_stored": 1,
    "cas_corrupt_bytes_rejected": 1,
    "cas_corrupt_evidence_recorded": 1,
    "cas_replica_serve_promised": 1,
    "cas_replica_stored": 1,
    "cas_replication_confirmed": 1,
    "cas_replication_promised": 1,
    "cas_retention_promised": 1,
    "cas_serve_promised": 1,
    "cas_storage_promised": 1,
    "cas_verification_promised": 1,
    "cid_compute_promised": 1,
    "compute_bad_result_promised": 1,
    "compute_cache_checkpointed": 1,
    "compute_context_promised": 1,
    "compute_function_executed": 1,
    "compute_result_locally_rejected": 1,
    "compute_result_locally_verified": 3,
    "compute_result_peer_rejected": 2,
    "compute_result_peer_verified": 2,
    "compute_result_promised": 1,
    "compute_result_received": 2,
    "compute_verification_received": 4,
    "economics_capacity_refused": 1,
    "economics_capacity_reserved": 2,
    "economics_credit_accepted": 3,
    "economics_credit_offered": 2,
    "economics_credits_earned": 3,
    "economics_credits_spent": 2,
    "economics_payment_withheld": 1,
    "economics_price_probe": 1,
    "economics_price_refused": 2,
    "evidence_report_received": 4,
    "llm_decision_live": 8,
    "peer_readiness_observed": 32,
    "primary_storage_unavailable": 1,
    "promise_envelope_validated": 56,
    "promise_variant_not_promised": 1,
    "replica_capability_token_issued": 1,
    "replica_capability_token_received": 1,
    "replica_capability_token_redeemed": 1,
    "replica_recovery_requested": 1,
    "replica_recovery_succeeded": 1,
    "runtime_done_promised": 8,
    "runtime_readiness_promised": 8,
    "tcp_message_received": 28,
    "tcp_message_send_failed": 1,
    "tcp_message_sent": 28,
    "trust_driven_peer_choice": 3,
    "trust_repair_promise_recorded": 1,
    "trust_updated": 14,
    "unknown_pcid_not_promised": 1
  },
  "outcome_counts": {
    "broken": 4,
    "kept": 243,
    "malformed": 9,
    "non_commitment": 8
  },
  "agent_counts": {
    "alice": 83,
    "bob": 33,
    "carol": 31,
    "dave": 24,
    "ellen": 14,
    "frank": 27,
    "grace": 35,
    "mallory": 17
  },
  "protocol_counts": {
    "cas_storage_v1": 84,
    "cid_compute_v1": 68,
    "cidv1-raw-sha2-256:dd8aed1efb771d22b92347b51891053c9c1333fc975e203e2c175c2baa3f8f4d": 4,
    "evidence_report_v1": 24
  },
  "trust_driven_choice_counts": {
    "alice": 3
  },
  "economics_counts": {
    "alice": 7,
    "bob": 4,
    "carol": 4,
    "frank": 2
  },
  "verification_counts": {
    "alice": 5,
    "dave": 2,
    "grace": 3
  },
  "replica_recovery_counts": {
    "alice": 2,
    "frank": 1
  },
  "rpc_drift_counts": {},
  "placeholder_live_decision_counts": {},
  "score_report": {
    "overall": 5,
    "transport": 5,
    "storage": 5,
    "compute": 5,
    "economics": 5,
    "trust": 5,
    "verification": 5,
    "replica": 5
  }
}
```
