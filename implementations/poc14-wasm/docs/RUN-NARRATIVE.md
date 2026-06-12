# POC14 Run Narrative

POC14 is the current executable superset target after POC13. It preserves the
POC11/POC12/POC13 behavior set and adds heterogeneous WASM/stdio process
boundaries plus decentralized monitoring evidence. Source: `DI-sihuz`;
`DI-linof`; `DI-lulof`.

## Expected Clean-Run Shape

1. Each Docker service starts one local `poc14-kernel` plus its configured app
   processes.
2. Relationship, storage, compute, shipping, device, and accounting agents keep
   the inherited POC13 startup and turn evidence.
3. Peggy starts as `poc14-wasm-agent`, validates the minimal WASM module fixture,
   and sends WASM-boundary evidence as an ordinary `relationship_v1` promise.
4. Victor starts as `poc14-stdio-adapter`, launches `poc14-stdio-worker`, receives
   an exact signed envelope over stdout, forwards it through the local kernel,
   and returns the exact peer ACK over stdin.
5. Alice records decentralized monitoring candidates as local evidence: local
   evidence summaries, peer-carried attestations, bearer-token exchange-rate
   signals, topology signals, and voluntary gossip.
6. Alice records hard local trust-boundary evidence: permanent local distrust of
   Mallory and Alice's promise that Alice's inbound/outbound traffic should not
   transit Mallory.
7. Alice records mixed-version pCID migration and same-run restart recovery
   evidence as local promises.
8. `poc14-analyze` counts inherited POC13 evidence plus POC14 boundary,
   decentralized-monitoring, migration, and restart events.

## Boundary Evidence

- `wasm_process_agent_started`
- `wasm_module_header_validated`
- `wasm_boundary_promise_sent`
- `wasm_boundary_ack_received`
- `stdio_worker_started`
- `stdio_worker_envelope_received`
- `stdio_adapter_kernel_forwarded`
- `stdio_worker_ack_observed`

## Decentralized Monitoring Evidence

- `production_monitor_boundary_recorded`
- `local_evidence_summary_promised`
- `peer_carried_attestation_promised`
- `bearer_token_exchange_rate_observed`
- `relationship_topology_signal_observed`
- `voluntary_gossip_promised`

## Hard Trust-Boundary Evidence

- `permanent_distrust_decided`
- `permanent_distrust_future_repair_not_promised`
- `permanent_distrust_direct_peer_removed`
- `permanent_distrust_send_blocked`
- `transit_exclusion_promised`
- `input_transit_exclusion_recorded`
- `output_transit_exclusion_recorded`
- `transit_candidate_rejected`
- `transit_route_candidate_blocked`
- `transit_safe_route_selected`

## Migration And Restart Evidence

- `mixed_version_pcid_migration_promised`
- `mixed_version_legacy_pcid_observed`
- `mixed_version_successor_pcid_selected`
- `run_internal_restart_orchestration_promised`
- `run_internal_restart_checkpoint_promised`
- `run_internal_restart_recovery_observed`

## Current Status

The 2026-06-11 `poc14-demo` clean Docker run passed `poc14-analyze` after adding
the `DI-dubih` behavioral hard-boundary scenarios. The run produced 2079 total
events, all analyzer score dimensions were `5`, and the run included one event
each for `permanent_distrust_send_blocked` and
`transit_route_candidate_blocked` in addition to the earlier hard-boundary
evidence events. `rpc_drift_counts` and `resource_trust_coupling_counts` were
empty.

The run remains not production-ready: the production-fitness report kept
`ready_for_production=false` because the monitor scored `protocol_validity`,
`local_trust_correctness`, and `imposition_avoidance` at `4/5`. Source:
`DI-linof`; `DI-lulof`; `DI-kinaf`; `DI-dubih`.
