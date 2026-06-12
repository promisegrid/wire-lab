# POC14 Implementation Notes

POC14 was scaffolded from POC13 and then renamed locally to
`implementations/poc14-wasm/`. The copy is intentional: future POCs should be a
superset of prior POCs unless a scoped DI explicitly says otherwise. Source:
`DI-sihuz`; `DI-linof`.

## What Changed From POC13

- Module path and command names now use `promisegrid.dev/wire-lab/implementations/poc14-wasm` and `poc14-*` binaries.
- Docker run state moved to `/run/poc14` and the compose network/volume are
  named `poc14` / `poc14-run`.
- Peggy and Victor were added to `config.example.json` in a new `wasm-stdio`
  container.
- `poc14-wasm-agent` records WASM module-boundary evidence without adding a WASM
  host-call RPC surface.
- `poc14-stdio-adapter` and `poc14-stdio-worker` record stdio-only worker
  messaging while preserving exact PromiseGrid envelopes.
- `poc14-analyze` now reports `heterogeneous_boundary_counts`,
  `decentralized_monitor_counts`, `migration_counts`, `restart_counts`, and
  score dimensions for `boundary`, `monitoring`, `migration`, and `restart`.
- Scripted `cas_storage_v1` and `cid_compute_v1` messages now use pCID-owned
  CBOR array payloads on the wire, with `field_*` compatibility projections only
  inside local handlers and analyzer evidence. Source: `DI-gahuh`.
- Peggy and Victor now each send one useful routed `relationship_v1` promise to
  Dave so the heterogeneous-boundary agents do more than record boundary
  existence. Source: `DI-pamob`.
- POC15 planning lives under `implementations/poc15-multihop/` and should make
  real multi-hop forwarding and route-exclusion-by-peer-promise executable.
  Source: `DI-pamob`.

## Non-Goals

- POC14 does not add global monitoring, global trust, global exchange rates,
  authorization services, permission checks, conformance authority, or RPC verbs.
- POC14 does not persist POC state across clean runs.
- POC14 does not replace app-local trust judgment with kernel or analyzer
  judgment.
- POC14 does not implement real multi-hop forwarding; it records route-candidate
  and route-exclusion evidence so POC15 can implement forwarding without
  pretending to own a global route view.
