# GA Backfill Comparison Report

- Run group: `ga-backfill-20260523-060223-dedupe`
- State: `results/state/ga-backfill-20260523-060223-dedupe.json`
- Comparison basis: latest exact-match canonical `promisegrid.ga.result.v1` record for the same `sim_id` + `scenario_id`, preferring the same `runner.api_model` when available and collapsing same-lineage historical reruns to the latest record.
- Compared cells: `22`
- Unmatched v2 cells: `0`
- Historical rerun groups collapsed: `7`
- Ambiguous matched pairs: `0`

## Sim Rank Drift

| V2 rank | V1 rank | Δ rank | Sim | Family | Cells | V1 avg | V2 avg | Δ avg | Vocab | Source |
|---:|---:|---:|---|---|---:|---:|---:|---:|---|---|
| 1 | 1 | +0 | `SIM-robot-app-semantics-conformance` | conformance-family | 7 | 66.14 | 66.57 | +0.43 | hard_hit | historical |
| 2 | 2 | +0 | `SIM-kurim-grid-envelope` | grid-envelope | 15 | 45.54 | 58.00 | +12.46 | clean | historical |

## Largest Cell Deltas

| Δ score | Sim | Scenario | V1 | V2 | V1 result | V2 result |
|---:|---|---|---:|---:|---|---|
| +24.00 | `SIM-kurim-grid-envelope` | `cas-object-type-binding-bakeoff-raw-chunk-versus-pointer-bytes` | 50.00 | 74.00 | `results/SIM-kurim-grid-envelope/cas-object-type-binding-bakeoff-raw-chunk-versus-pointer-bytes/openai-gpt-5.4-xhigh/20260522-012332.json` | `results/SIM-kurim-grid-envelope/cas-object-type-binding-bakeoff-raw-chunk-versus-pointer-bytes/openai-gpt-5.4-medium/20260523-060228.json` |
| +22.00 | `SIM-kurim-grid-envelope` | `bgp-class-routing-app-conflicting-policies` | 34.00 | 56.00 | `results/SIM-kurim-grid-envelope/bgp-class-routing-app-conflicting-policies/openai-gpt-5.4-xhigh/20260522-012332.json` | `results/SIM-kurim-grid-envelope/bgp-class-routing-app-conflicting-policies/openai-gpt-5.4-medium/20260523-060228.json` |
| +19.50 | `SIM-kurim-grid-envelope` | `chunk-feed-replication-carrier-independence` | 42.50 | 62.00 | `results/SIM-kurim-grid-envelope/chunk-feed-replication-carrier-independence/openai-gpt-5.4-xhigh/20260522-012332.json` | `results/SIM-kurim-grid-envelope/chunk-feed-replication-carrier-independence/openai-gpt-5.4-medium/20260523-060228.json` |
| +19.00 | `SIM-kurim-grid-envelope` | `bgp-class-routing-app-route-leak` | 37.00 | 56.00 | `results/SIM-kurim-grid-envelope/bgp-class-routing-app-route-leak/openai-gpt-5.4-xhigh/20260522-012332.json` | `results/SIM-kurim-grid-envelope/bgp-class-routing-app-route-leak/openai-gpt-5.4-medium/20260523-060228.json` |
| +19.00 | `SIM-kurim-grid-envelope` | `chunk-feed-replication-corrupt-chunk` | 37.00 | 56.00 | `results/SIM-kurim-grid-envelope/chunk-feed-replication-corrupt-chunk/openai-gpt-5.4-xhigh/20260522-012332.json` | `results/SIM-kurim-grid-envelope/chunk-feed-replication-corrupt-chunk/openai-gpt-5.4-medium/20260523-060228.json` |
| +18.00 | `SIM-kurim-grid-envelope` | `bgp-class-routing-app-honest-reachability-promise` | 40.00 | 58.00 | `results/SIM-kurim-grid-envelope/bgp-class-routing-app-honest-reachability-promise/openai-gpt-5.4-xhigh/20260522-012332.json` | `results/SIM-kurim-grid-envelope/bgp-class-routing-app-honest-reachability-promise/openai-gpt-5.4-medium/20260523-060228.json` |
| +15.00 | `SIM-kurim-grid-envelope` | `udp-feed-v0-conformance-malformed-datagram` | 37.00 | 52.00 | `results/SIM-kurim-grid-envelope/udp-feed-v0-conformance-malformed-datagram/openai-gpt-5.4-xhigh/20260522-012332.json` | `results/SIM-kurim-grid-envelope/udp-feed-v0-conformance-malformed-datagram/openai-gpt-5.4-medium/20260523-060228.json` |
| +15.00 | `SIM-kurim-grid-envelope` | `conditional-release-geofencing-replay-outside-conditions` | 43.00 | 58.00 | `results/SIM-kurim-grid-envelope/conditional-release-geofencing-replay-outside-conditions/openai-gpt-5.4-xhigh/20260522-012332.json` | `results/SIM-kurim-grid-envelope/conditional-release-geofencing-replay-outside-conditions/openai-gpt-5.4-medium/20260523-060228.json` |
| +13.00 | `SIM-kurim-grid-envelope` | `conditional-release-geofencing-opaque-lower-layer-carriage` | 51.00 | 64.00 | `results/SIM-kurim-grid-envelope/conditional-release-geofencing-opaque-lower-layer-carriage/openai-gpt-5.4-xhigh/20260522-012332.json` | `results/SIM-kurim-grid-envelope/conditional-release-geofencing-opaque-lower-layer-carriage/openai-gpt-5.4-medium/20260523-060228.json` |
| +12.00 | `SIM-kurim-grid-envelope` | `bgp-class-routing-app-sparse-knowledge` | 40.00 | 52.00 | `results/SIM-kurim-grid-envelope/bgp-class-routing-app-sparse-knowledge/openai-gpt-5.4-xhigh/20260522-012332.json` | `results/SIM-kurim-grid-envelope/bgp-class-routing-app-sparse-knowledge/openai-gpt-5.4-medium/20260523-060228.json` |
| +11.00 | `SIM-kurim-grid-envelope` | `chunk-feed-replication-duplicate-advertisement` | 43.00 | 54.00 | `results/SIM-kurim-grid-envelope/chunk-feed-replication-duplicate-advertisement/openai-gpt-5.4-xhigh/20260522-012332.json` | `results/SIM-kurim-grid-envelope/chunk-feed-replication-duplicate-advertisement/openai-gpt-5.4-medium/20260523-060228.json` |
| -10.00 | `SIM-kurim-grid-envelope` | `udp-feed-v0-conformance-session-layer-composition` | 60.00 | 50.00 | `results/SIM-kurim-grid-envelope/udp-feed-v0-conformance-session-layer-composition/openai-gpt-5.4-xhigh/20260522-012332.json` | `results/SIM-kurim-grid-envelope/udp-feed-v0-conformance-session-layer-composition/openai-gpt-5.4-medium/20260523-060228.json` |

## Family Highlights

### grid-envelope

| Sim | V1 avg | V2 avg | Δ avg | Vocab | Source |
|---|---:|---:|---:|---|---|
| `SIM-kurim-grid-envelope` | 45.54 | 58.00 | +12.46 | clean | historical |

### conformance-family

| Sim | V1 avg | V2 avg | Δ avg | Vocab | Source |
|---|---:|---:|---:|---|---|
| `SIM-robot-app-semantics-conformance` | 66.14 | 66.57 | +0.43 | hard_hit | historical |

## Cell Detail

| Sim | Scenario | V1 model | V2 model | V1 | V2 | Δ | Vocab | Source |
|---|---|---|---|---:|---:|---:|---|---|
| `SIM-kurim-grid-envelope` | `bgp-class-routing-app-conflicting-policies` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 34.00 | 56.00 | +22.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `bgp-class-routing-app-honest-reachability-promise` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 40.00 | 58.00 | +18.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `bgp-class-routing-app-route-leak` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 37.00 | 56.00 | +19.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `bgp-class-routing-app-sparse-knowledge` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 40.00 | 52.00 | +12.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `cas-backed-group-session-envelope-independence` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 63.00 | 62.00 | -1.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `cas-object-type-binding-bakeoff-raw-chunk-versus-pointer-bytes` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 50.00 | 74.00 | +24.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `cas-object-type-binding-bakeoff-unknown-typed-object` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 57.00 | 60.00 | +3.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `chunk-feed-replication-carrier-independence` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 42.50 | 62.00 | +19.50 | clean | historical |
| `SIM-kurim-grid-envelope` | `chunk-feed-replication-corrupt-chunk` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 37.00 | 56.00 | +19.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `chunk-feed-replication-duplicate-advertisement` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 43.00 | 54.00 | +11.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `conditional-release-geofencing-opaque-lower-layer-carriage` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 51.00 | 64.00 | +13.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `conditional-release-geofencing-replay-outside-conditions` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 43.00 | 58.00 | +15.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `transport-family-bakeoff-per-hop-authorization-failure` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 48.60 | 56.00 | +7.40 | clean | historical |
| `SIM-kurim-grid-envelope` | `udp-feed-v0-conformance-malformed-datagram` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 37.00 | 52.00 | +15.00 | clean | historical |
| `SIM-kurim-grid-envelope` | `udp-feed-v0-conformance-session-layer-composition` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 60.00 | 50.00 | -10.00 | clean | historical |
| `SIM-robot-app-semantics-conformance` | `app-semantics-partial-conformance` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 77.00 | 80.00 | +3.00 | hard_hit | historical |
| `SIM-robot-app-semantics-conformance` | `device-bound-agent-physical-effect` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 57.50 | 64.00 | +6.50 | hard_hit | historical |
| `SIM-robot-app-semantics-conformance` | `kernel-porting-boundary` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 43.00 | 48.00 | +5.00 | hard_hit | historical |
| `SIM-robot-app-semantics-conformance` | `live-crdt-audit-publication` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 69.00 | 74.00 | +5.00 | hard_hit | historical |
| `SIM-robot-app-semantics-conformance` | `minimal-immutable-blob-app` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 69.00 | 68.00 | -1.00 | hard_hit | historical |
| `SIM-robot-app-semantics-conformance` | `multi-embodiment-app-identity` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 75.00 | 68.00 | -7.00 | hard_hit | historical |
| `SIM-robot-app-semantics-conformance` | `portable-signing-key-identity` | `openai-gpt-5.4-xhigh` | `openai-gpt-5.4-medium` | 72.50 | 64.00 | -8.50 | hard_hit | historical |
