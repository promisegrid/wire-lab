# POC14 Implementation Notes

POC14 was scaffolded from POC13 and then renamed locally to
`implementations/poc14-wasm/`. The copy is intentional: future POCs should be a
superset of prior POCs unless a scoped DI explicitly says otherwise. Source:
`DI-sihuz`; `DI-linof`.

## What Changed From POC13

- Module path and command names now use `promisegrid.dev/wire-lab/implementations/poc14-wasm` and `poc14-*` binaries.
- Docker run state visible to the observer-only collector lives at `/run/poc14`.
  Agent containers do not mount that observer run volume or coordinate through
  marker files. Source: `DI-dirat`.
- Peggy and Victor were added to committed `config.json` in a new `wasm-stdio`
  container.
- `poc14-wasm-agent` records real wazero compile/instantiate/call events for an
  embedded no-import Fibonacci module and keeps Alice's `cid_compute_v1` promise
  through that export without adding a WASM host-call RPC surface. Source:
  `DI-kimim`; `DI-sivis`.
- `poc14-stdio-adapter` and `poc14-stdio-worker` record stdio-only worker
  messaging with length-prefixed CBOR frames while preserving exact PromiseGrid
  envelopes as CBOR byte strings. Victor now delegates Alice's exact inbound
  `cid_compute_v1` envelope to the worker, and the worker returns an exact signed
  compute ACK over stdout. Source: `DI-kimim`; `DI-sivis`.
- `poc14-analyze` now reports `runtime_adapter_event_counts`,
  `decentralized_monitor_counts`, `migration_counts`, `restart_counts`, and
  score dimensions for `runtime_adapter`, `monitoring`, `migration`, and
  `restart`.
- Known POC14 protocol messages now use pCID-owned CBOR array payloads on the
  wire, with `field_*` compatibility projections only inside local handlers and
  analyzer event records. Source: `DI-gahuh`; `DI-dirat`.
- Peggy and Victor now each send one useful routed `relationship_v1` promise to
  Dave and each keep one useful `cid_compute_v1` promise for Alice so the
  heterogeneous-runtime-adapter agents do more than record runtime-adapter existence.
  Source: `DI-pamob`; `DI-sivis`.
- POC15 planning lives under `implementations/poc15-multihop-multiarity-dag/`
  and should make real multi-hop forwarding, route-exclusion-by-peer-promise,
  pCID-owned multiarity, raw-message DAG review, and parent-link specimens
  executable. Source: `DI-pamob`; `DI-podut`.

## Non-Goals

- POC14 does not add global monitoring, global trust, global exchange rates,
  authorization services, permission checks, conformance authority, or RPC verbs.
- POC14 does not persist POC state across clean runs.
- POC14 does not replace app-local trust judgment with kernel or analyzer
  judgment.
- POC14 does not implement real multi-hop forwarding; it records route-candidate
  and route-exclusion event so POC15 can implement forwarding without
  pretending to own a global route view.
