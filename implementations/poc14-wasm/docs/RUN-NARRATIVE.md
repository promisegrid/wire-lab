# POC14 Run Narrative

POC14 is the current executable superset target after POC13. It preserves the
POC11/POC12/POC13 behavior set and adds heterogeneous WASM/stdio runtime
adapters plus decentralized monitoring event records. Source: `DI-sihuz`;
`DI-linof`; `DI-lulof`; `DI-kimim`; `DI-sivis`.

## Expected Clean-Run Shape

1. Each Docker service starts one local `poc14-kernel` plus its configured app
   processes.
2. Relationship, storage, compute, shipping, device, and accounting agents keep
   the inherited POC13 startup and turn event.
3. Peggy starts as `poc14-wasm-agent`, compiles, instantiates, and calls an
   embedded no-import WASM Fibonacci module with wazero, sends WASM-adapter
   event as an ordinary `relationship_v1` promise, promises Dave reusable
   module-execution event, and later keeps Alice's `cid_compute_v1` promise by
   executing that WASM export.
4. Victor starts as `poc14-stdio-adapter`, launches `poc14-stdio-worker`, receives
   an exact signed envelope as a CBOR byte string inside a length-prefixed CBOR
   frame over stdout, forwards it through the local kernel, returns the exact
   peer ACK as a CBOR byte string over stdin, promises Dave reusable subprocess
   round-trip event, and later delegates Alice's exact inbound `cid_compute_v1`
   envelope to the worker so the worker returns a signed compute ACK.
5. Alice records decentralized monitoring candidates as local events: local
   event summaries, peer-carried attestations, bearer-token exchange-rate
   signals, topology signals, and voluntary gossip.
6. Alice records hard local trust and routing events: permanent local distrust of
   Mallory and Alice's promise that Alice's inbound/outbound traffic should not
   transit Mallory.
7. Alice records mixed-version pCID migration and same-run restart recovery
   event as local promises.
8. `poc14-analyze` counts inherited POC13 events plus POC14 runtime-adapter,
   decentralized-monitoring, migration, and restart events.
9. Scripted `cas_storage_v1` and `cid_compute_v1` promises travel as pCID-owned
   CBOR array payloads and appear in local handlers only as compatibility
   projections.

## Runtime Adapter Events

- `wasm_process_agent_started`
- `wasm_module_instantiated`
- `wasm_export_called`
- `wasm_export_result_observed`
- `wasm_adapter_promise_sent`
- `wasm_adapter_ack_received`
- `wasm_compute_request_promised`
- `wasm_compute_request_received`
- `wasm_compute_function_executed`
- `wasm_compute_result_promised`
- `wasm_compute_result_verified`
- `wasm_useful_work_promised`
- `wasm_useful_work_ack_received`
- `stdio_worker_started`
- `stdio_cbor_request_sent`
- `stdio_worker_envelope_received`
- `stdio_cbor_envelope_received`
- `stdio_adapter_kernel_forwarded`
- `stdio_cbor_ack_sent`
- `stdio_worker_ack_event`
- `stdio_cbor_ack_event`
- `stdio_compute_request_promised`
- `stdio_compute_worker_started`
- `stdio_compute_request_forwarded`
- `stdio_compute_worker_executed`
- `stdio_compute_ack_received`
- `stdio_compute_result_verified`
- `stdio_useful_work_promised`
- `stdio_useful_work_ack_received`

## Migrated Payload Event

- `pcid_owned_array_payload_sent`
- `pcid_owned_array_payload_received`
- `pcid_owned_array_ack_sent`
- `pcid_owned_array_ack_received`

## Decentralized Monitoring Event

- `decentralized_monitoring_model_recorded`
- `local_event_summary_promised`
- `peer_carried_attestation_promised`
- `bearer_token_exchange_rate_observed`
- `relationship_topology_signal_observed`
- `voluntary_gossip_promised`

## Hard Trust And Routing Events

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

## Migration And Restart Event

- `mixed_version_pcid_migration_promised`
- `mixed_version_legacy_pcid_observed`
- `mixed_version_successor_pcid_selected`
- `run_internal_restart_orchestration_promised`
- `run_internal_restart_checkpoint_promised`
- `run_internal_restart_recovery_observed`

## Current Status

The latest clean Docker baseline after `DI-sivis` is the 2026-06-13
`poc14-demo` run after Peggy and Victor became useful `cid_compute_v1` peers.
The run produced 2269 total events; all analyzer score dimensions were `5`;
`runtime_adapter_event_counts` included one each of the Peggy
`wasm_compute_request_promised`, `wasm_compute_request_received`,
`wasm_compute_function_executed`, `wasm_compute_result_promised`, and
`wasm_compute_result_verified` records, and one each of the Victor
`stdio_compute_request_promised`, `stdio_compute_worker_started`,
`stdio_compute_request_forwarded`, `stdio_compute_worker_executed`,
`stdio_compute_ack_received`, and `stdio_compute_result_verified` records.
`protocol_counts` included `identity_key_v1=42`, `cas_storage_v1=269`,
`cid_compute_v1=260`, `relationship_v1=677`, `accounting_v1=41`,
`postal_scale_v1=19`, `printer_port_v1=32`, and `ups_label_v1=19`; and
`rpc_drift_counts` plus `resource_trust_coupling_counts` were empty.

The run remains not production-ready: the production-fitness report kept
`ready_for_production=false` because the monitor scored `protocol_validity` at
`4/5`. Source: `DI-linof`; `DI-lulof`; `DI-kinaf`; `DI-dubih`; `DI-gahuh`;
`DI-kimim`; `DI-sivis`.
