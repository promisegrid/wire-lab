# GA Backfill Comparison Report

- Run group: `ga-backfill-20260522-215638-medium`
- State: `results/state/ga-backfill-20260522-215638-medium.json`
- Comparison basis: latest exact-match canonical `promisegrid.ga.result.v1` record for the same `sim_id` + `scenario_id`, preferring the same `runner.api_model` when available.
- Compared cells: `22`
- Unmatched v2 cells: `0`
- Ambiguous matched pairs: `7`

## Sim Rank Drift

| V2 rank | V1 rank | Δ rank | Sim | Family | Cells | V1 avg | V2 avg | Δ avg | Vocab | Source |
|---:|---:|---:|---|---|---:|---:|---:|---:|---|---|
| 1 | 1 | +0 | `SIM-robot-app-semantics-conformance` | conformance-family | 7 | 66.14 | 70.57 | +4.43 | hard_hit | historical |
| 2 | 2 | +0 | `SIM-kurim-grid-envelope` | grid-envelope | 15 | 45.54 | 60.27 | +14.73 | clean | historical |

## Largest Cell Deltas

| Δ score | Sim | Scenario | V1 | V2 | V1 result | V2 result |
|---:|---|---|---:|---:|---|---|
| +28.00 | `SIM-kurim-grid-envelope` | `bgp-class-routing-app-conflicting-policies` | 34.00 | 62.00 | `results/SIM-kurim-grid-envelope/bgp-class-routing-app-conflicting-policies/openai-gpt-5.4-xhigh/20260522-012332.json` | `results/SIM-kurim-grid-envelope/bgp-class-routing-app-conflicting-policies/openai-gpt-5.4-medium/20260523-045646.json` |
| +23.00 | `SIM-kurim-grid-envelope` | `udp-feed-v0-conformance-malformed-datagram` | 37.00 | 60.00 | `results/SIM-kurim-grid-envelope/udp-feed-v0-conformance-malformed-datagram/openai-gpt-5.4-xhigh/20260522-012332.json` | `results/SIM-kurim-grid-envelope/udp-feed-v0-conformance-malformed-datagram/openai-gpt-5.4-medium/20260523-045646.json` |
| +21.00 | `SIM-kurim-grid-envelope` | `bgp-class-routing-app-route-leak` | 37.00 | 58.00 | `results/SIM-kurim-grid-envelope/bgp-class-routing-app-route-leak/openai-gpt-5.4-xhigh/20260522-012332.json` | `results/SIM-kurim-grid-envelope/bgp-class-routing-app-route-leak/openai-gpt-5.4-medium/20260523-045646.json` |
| +19.50 | `SIM-kurim-grid-envelope` | `chunk-feed-replication-carrier-independence` | 42.50 | 62.00 | `results/SIM-kurim-grid-envelope/chunk-feed-replication-carrier-independence/openai-gpt-5.4-xhigh/20260522-012332.json` | `results/SIM-kurim-grid-envelope/chunk-feed-replication-carrier-independence/openai-gpt-5.4-medium/20260523-045646.json` |
| +18.00 | `SIM-kurim-grid-envelope` | `bgp-class-routing-app-honest-reachability-promise` | 40.00 | 58.00 | `results/SIM-kurim-grid-envelope/bgp-class-routing-app-honest-reachability-promise/openai-gpt-5.4-xhigh/20260522-012332.json` | `results/SIM-kurim-grid-envelope/bgp-class-routing-app-honest-reachability-promise/openai-gpt-5.4-medium/20260523-045646.json` |
| +18.00 | `SIM-kurim-grid-envelope` | `bgp-class-routing-app-sparse-knowledge` | 40.00 | 58.00 | `results/SIM-kurim-grid-envelope/bgp-class-routing-app-sparse-knowledge/openai-gpt-5.4-xhigh/20260522-012332.json` | `results/SIM-kurim-grid-envelope/bgp-class-routing-app-sparse-knowledge/openai-gpt-5.4-medium/20260523-045646.json` |
| +17.00 | `SIM-kurim-grid-envelope` | `chunk-feed-replication-corrupt-chunk` | 37.00 | 54.00 | `results/SIM-kurim-grid-envelope/chunk-feed-replication-corrupt-chunk/openai-gpt-5.4-xhigh/20260522-012332.json` | `results/SIM-kurim-grid-envelope/chunk-feed-replication-corrupt-chunk/openai-gpt-5.4-medium/20260523-045646.json` |
| +15.00 | `SIM-kurim-grid-envelope` | `cas-object-type-binding-bakeoff-unknown-typed-object` | 57.00 | 72.00 | `results/SIM-kurim-grid-envelope/cas-object-type-binding-bakeoff-unknown-typed-object/openai-gpt-5.4-xhigh/20260522-012332.json` | `results/SIM-kurim-grid-envelope/cas-object-type-binding-bakeoff-unknown-typed-object/openai-gpt-5.4-medium/20260523-045646.json` |
| +14.00 | `SIM-kurim-grid-envelope` | `cas-object-type-binding-bakeoff-raw-chunk-versus-pointer-bytes` | 50.00 | 64.00 | `results/SIM-kurim-grid-envelope/cas-object-type-binding-bakeoff-raw-chunk-versus-pointer-bytes/openai-gpt-5.4-xhigh/20260522-012332.json` | `results/SIM-kurim-grid-envelope/cas-object-type-binding-bakeoff-raw-chunk-versus-pointer-bytes/openai-gpt-5.4-medium/20260523-045646.json` |
| +12.00 | `SIM-kurim-grid-envelope` | `udp-feed-v0-conformance-session-layer-composition` | 60.00 | 72.00 | `results/SIM-kurim-grid-envelope/udp-feed-v0-conformance-session-layer-composition/openai-gpt-5.4-xhigh/20260522-012332.json` | `results/SIM-kurim-grid-envelope/udp-feed-v0-conformance-session-layer-composition/openai-gpt-5.4-medium/20260523-045646.json` |
| +11.00 | `SIM-kurim-grid-envelope` | `conditional-release-geofencing-opaque-lower-layer-carriage` | 51.00 | 62.00 | `results/SIM-kurim-grid-envelope/conditional-release-geofencing-opaque-lower-layer-carriage/openai-gpt-5.4-xhigh/20260522-012332.json` | `results/SIM-kurim-grid-envelope/conditional-release-geofencing-opaque-lower-layer-carriage/openai-gpt-5.4-medium/20260523-045646.json` |
| +11.00 | `SIM-robot-app-semantics-conformance` | `live-crdt-audit-publication` | 69.00 | 80.00 | `results/SIM-robot-app-semantics-conformance/live-crdt-audit-publication/openai-gpt-5.4-xhigh/20260521-234037.json` | `results/SIM-robot-app-semantics-conformance/live-crdt-audit-publication/openai-gpt-5.4-medium/20260523-045646.json` |

## Family Highlights

### grid-envelope

| Sim | V1 avg | V2 avg | Δ avg | Vocab | Source |
|---|---:|---:|---:|---|---|
| `SIM-kurim-grid-envelope` | 45.54 | 60.27 | +14.73 | clean | historical |

### conformance-family

| Sim | V1 avg | V2 avg | Δ avg | Vocab | Source |
|---|---:|---:|---:|---|---|
| `SIM-robot-app-semantics-conformance` | 66.14 | 70.57 | +4.43 | hard_hit | historical |

## Cell Detail

| Sim | Scenario | V1 model | V2 model | V1 | V2 | Δ | Vocab | Source |
|---|---|---|---|---:|---:|---:|---|---|
| `SIM-kurim-grid-envelope` | `bgp-class-routing-app-conflicting-policies` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 34.00 | 62.00 | +28.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `bgp-class-routing-app-honest-reachability-promise` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 40.00 | 58.00 | +18.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `bgp-class-routing-app-route-leak` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 37.00 | 58.00 | +21.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `bgp-class-routing-app-sparse-knowledge` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 40.00 | 58.00 | +18.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `cas-backed-group-session-envelope-independence` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 63.00 | 64.00 | +1.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `cas-object-type-binding-bakeoff-raw-chunk-versus-pointer-bytes` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 50.00 | 64.00 | +14.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `cas-object-type-binding-bakeoff-unknown-typed-object` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 57.00 | 72.00 | +15.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `chunk-feed-replication-carrier-independence` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 42.50 | 62.00 | +19.50 | clean | historical |
| `SIM-kurim-grid-envelope` | `chunk-feed-replication-corrupt-chunk` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 37.00 | 54.00 | +17.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `chunk-feed-replication-duplicate-advertisement` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 43.00 | 52.00 | +9.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `conditional-release-geofencing-opaque-lower-layer-carriage` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 51.00 | 62.00 | +11.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `conditional-release-geofencing-replay-outside-conditions` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 43.00 | 52.00 | +9.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `transport-family-bakeoff-per-hop-authorization-failure` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 48.60 | 54.00 | +5.40 | clean | historical |
| `SIM-kurim-grid-envelope` | `udp-feed-v0-conformance-malformed-datagram` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 37.00 | 60.00 | +23.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `udp-feed-v0-conformance-session-layer-composition` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 60.00 | 72.00 | +12.00 | clean | historical |
| `SIM-robot-app-semantics-conformance` | `app-semantics-partial-conformance` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 77.00 | 76.00 | -1.00 | hard_hit | historical |
| `SIM-robot-app-semantics-conformance` | `device-bound-agent-physical-effect` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 57.50 | 60.00 | +2.50 | hard_hit | historical |
| `SIM-robot-app-semantics-conformance` | `kernel-porting-boundary` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 43.00 | 52.00 | +9.00 | hard_hit | historical |
| `SIM-robot-app-semantics-conformance` | `live-crdt-audit-publication` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 69.00 | 80.00 | +11.00 | hard_hit | historical |
| `SIM-robot-app-semantics-conformance` | `minimal-immutable-blob-app` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 69.00 | 72.00 | +3.00 | hard_hit | historical |
| `SIM-robot-app-semantics-conformance` | `multi-embodiment-app-identity` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 75.00 | 82.00 | +7.00 | hard_hit | historical |
| `SIM-robot-app-semantics-conformance` | `portable-signing-key-identity` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 72.50 | 72.00 | -0.50 | hard_hit | historical |
