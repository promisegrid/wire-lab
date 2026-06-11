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
6. Alice records mixed-version pCID migration and same-run restart recovery
   evidence as local promises.
7. `poc14-analyze` counts inherited POC13 evidence plus POC14 boundary,
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

## Migration And Restart Evidence

- `mixed_version_pcid_migration_promised`
- `mixed_version_legacy_pcid_observed`
- `mixed_version_successor_pcid_selected`
- `run_internal_restart_orchestration_promised`
- `run_internal_restart_checkpoint_promised`
- `run_internal_restart_recovery_observed`

## Current Status

No fresh POC14 Docker result is recorded in this document yet. The first clean
run should be executed with `scripts/run-clean.sh`, then this narrative should be
updated with the actual run ID, event totals, analyzer score report, and any
production-fitness gaps. Source: `DI-linof`.
