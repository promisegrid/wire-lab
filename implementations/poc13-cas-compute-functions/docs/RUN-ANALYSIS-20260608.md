# POC13 clean-run analyzer output, 2026-06-08

This file archives the exact analyzer JSON captured after a clean Docker run of
`implementations/poc13-cas-compute-functions/scripts/run-clean.sh`. Source:
`DI-nisaz`; `DI-kikoj`.

```json
{
  "run_dir": "/run/poc13/poc13-demo/run",
  "total_events": 511,
  "event_counts": {
    "agent_waiting_for_peer_promises": 6,
    "app_receive_promise_registered": 8,
    "bad_proof_rejected": 1,
    "bad_proof_sent": 1,
    "capability_token_expired": 2,
    "capability_token_issued": 2,
    "capability_token_received": 2,
    "capability_token_renewal_requested": 2,
    "capability_token_renewed": 4,
    "capability_token_revoked": 4,
    "capability_token_ttl_observed": 2,
    "capability_token_ttl_promised": 2,
    "cas_bytes_retrieved": 2,
    "cas_bytes_stored": 2,
    "cas_corrupt_bytes_rejected": 1,
    "cas_corrupt_evidence_recorded": 1,
    "cas_multi_object_pressure": 1,
    "cas_replica_serve_promised": 2,
    "cas_replica_stored": 2,
    "cas_replication_confirmed": 2,
    "cas_replication_promised": 2,
    "cas_retention_promised": 2,
    "cas_serve_promised": 2,
    "cas_storage_promised": 2,
    "cas_verification_promised": 1,
    "cid_compute_promised": 2,
    "compute_alternate_function_executed": 1,
    "compute_bad_result_promised": 2,
    "compute_cache_checkpointed": 2,
    "compute_cache_hit": 2,
    "compute_cache_miss": 1,
    "compute_cache_miss_observed": 1,
    "compute_cache_reused": 2,
    "compute_context_promised": 2,
    "compute_disagreement_resolved_locally": 2,
    "compute_followup_function_requested": 1,
    "compute_function_executed": 2,
    "compute_result_locally_rejected": 2,
    "compute_result_locally_verified": 6,
    "compute_result_peer_rejected": 4,
    "compute_result_peer_verified": 2,
    "compute_result_promised": 2,
    "compute_result_received": 6,
    "compute_verification_received": 6,
    "compute_verifier_disagreement": 4,
    "dynamic_peer_choice_from_persisted_trust": 2,
    "economics_capacity_refused": 1,
    "economics_capacity_reserved": 4,
    "economics_credit_accepted": 6,
    "economics_credit_offered": 2,
    "economics_credits_earned": 6,
    "economics_credits_spent": 4,
    "economics_payment_withheld": 2,
    "economics_price_probe": 1,
    "economics_price_refused": 2,
    "evidence_report_received": 8,
    "key_rotation_promise_recorded": 1,
    "llm_decision_live": 8,
    "network_outage_variant_selected": 4,
    "peer_readiness_observed": 32,
    "persisted_trust_history_loaded": 8,
    "primary_storage_unavailable": 4,
    "promise_envelope_validated": 123,
    "promise_variant_not_promised": 1,
    "replica_capability_token_issued": 2,
    "replica_capability_token_received": 2,
    "replica_capability_token_redeemed": 2,
    "replica_recovery_requested": 2,
    "replica_recovery_succeeded": 2,
    "runtime_done_promised": 8,
    "runtime_readiness_promised": 8,
    "tcp_message_received": 60,
    "tcp_message_send_failed": 4,
    "tcp_message_sent": 60,
    "trust_driven_peer_choice": 4,
    "trust_repair_promise_recorded": 1,
    "trust_updated": 26,
    "unknown_pcid_not_promised": 1
  },
  "outcome_counts": {
    "broken": 8,
    "kept": 455,
    "malformed": 19,
    "non_commitment": 29
  },
  "agent_counts": {
    "alice": 196,
    "bob": 53,
    "carol": 52,
    "dave": 51,
    "ellen": 20,
    "frank": 66,
    "grace": 52,
    "mallory": 21
  },
  "protocol_counts": {
    "cas_storage_v1": 192,
    "cid_compute_v1": 156,
    "cidv1-raw-sha2-256:dd8aed1efb771d22b92347b51891053c9c1333fc975e203e2c175c2baa3f8f4d": 4,
    "evidence_report_v1": 55
  },
  "trust_driven_choice_counts": {
    "alice": 4
  },
  "economics_counts": {
    "alice": 10,
    "bob": 7,
    "carol": 7,
    "frank": 4
  },
  "verification_counts": {
    "alice": 8,
    "dave": 4,
    "grace": 3
  },
  "replica_recovery_counts": {
    "alice": 4,
    "frank": 2
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
