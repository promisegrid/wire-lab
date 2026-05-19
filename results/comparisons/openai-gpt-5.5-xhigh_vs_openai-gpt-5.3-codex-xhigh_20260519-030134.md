# Cross-Model Drift Report: openai-gpt-5.5-xhigh vs openai-gpt-5.3-codex-xhigh

- Baseline model: `openai-gpt-5.5-xhigh`
- Comparison model: `openai-gpt-5.3-codex-xhigh`
- Cells compared: `28`
- Verdict text changes: `2`
- Score/rank changes: `1`

## Simulation Ranking Shift

| Simulation | Old avg score | Old rank | New avg score | New rank | Rank delta |
|---|---:|---:|---:|---:|---:|
| `SIM-riliz-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload` | 5.67 | 1 | 5.67 | 1 | +0 |
| `SIM-vamaz-grid-envelope-enc-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload` | 5.67 | 2 | 5.33 | 2 | +0 |
| `SIM-gazan-grid-envelope-enc-cbor-unknown-hard-reject-sig-unsigned-v0` | 3.00 | 3 | 3.00 | 3 | +0 |
| `SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload` | 3.00 | 4 | 3.00 | 4 | +0 |
| `SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid` | 3.00 | 5 | 3.00 | 5 | +0 |
| `SIM-rakir-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload` | 3.00 | 6 | 3.00 | 6 | +0 |
| `SIM-tohol-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-opaque-bytes` | 3.00 | 7 | 3.00 | 7 | +0 |
| `SIM-kovis-grid-envelope-enc-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload` | 2.67 | 8 | 2.67 | 8 | +0 |
| `SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload` | 2.67 | 9 | 2.67 | 9 | +0 |
| `SIM-kuful-udp-feed-v0-conformance` | 0.33 | 10 | 0.33 | 10 | +0 |

## Per-Cell Drift

| Simulation | Scenario | Old verdict | New verdict | Score delta |
|---|---|---|---|---:|
| `SIM-gazan-grid-envelope-enc-cbor-unknown-hard-reject-sig-unsigned-v0` | `shipping-label-printing` | partial fit, likely negative-control value | partial fit, negative-control baseline | +0 |
| `SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload` | `aerospace-project-funding` | good partial fit | good partial fit | +0 |
| `SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload` | `journalism-source-provenance` | partial with privacy risk | partial with privacy risk | +0 |
| `SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload` | `promisebase-reference-naming-promisebase-custom-syntax-migration` | risky partial fit | risky partial fit | +0 |
| `SIM-kovis-grid-envelope-enc-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload` | `aerospace-project-funding` | partial but brittle | partial but brittle | +0 |
| `SIM-kovis-grid-envelope-enc-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload` | `journalism-source-provenance` | partial but brittle | partial but brittle | +0 |
| `SIM-kovis-grid-envelope-enc-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload` | `promisebase-reference-naming-promisebase-custom-syntax-migration` | good strict baseline | good strict baseline | +0 |
| `SIM-kuful-udp-feed-v0-conformance` | `chunking-identity-bakeoff-profile-negotiation-mismatch` | weak fit, transport-only evidence | weak fit, transport-only evidence | +0 |
| `SIM-kuful-udp-feed-v0-conformance` | `municipal-governance` | poor standalone fit | poor standalone fit | +0 |
| `SIM-kuful-udp-feed-v0-conformance` | `promise-economy-spectrum-cryptocurrency-toxicity-failure` | poor fit, useful negative control | poor fit, useful negative control | +0 |
| `SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid` | `chunking-identity-bakeoff-profile-negotiation-mismatch` | partial fit | partial fit | +0 |
| `SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid` | `municipal-governance` | weak-to-partial fit | weak-to-partial fit | +0 |
| `SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid` | `promise-economy-spectrum-cryptocurrency-toxicity-failure` | negative-control partial fit | negative-control partial fit | +0 |
| `SIM-rakir-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload` | `aerospace-project-funding` | good partial fit | good partial fit | +0 |
| `SIM-rakir-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload` | `journalism-source-provenance` | partial with privacy risk | partial with privacy risk | +0 |
| `SIM-rakir-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload` | `promisebase-reference-naming-promisebase-custom-syntax-migration` | risky partial fit | risky partial fit | +0 |
| `SIM-riliz-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload` | `aerospace-project-funding` | strong fit | strong fit | +0 |
| `SIM-riliz-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload` | `journalism-source-provenance` | strongest fit in this slice | strongest fit in this slice | +0 |
| `SIM-riliz-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload` | `promisebase-reference-naming-promisebase-custom-syntax-migration` | best migration fit | best migration fit | +0 |
| `SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload` | `aerospace-project-funding` | partial but brittle | partial but brittle | +0 |
| `SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload` | `journalism-source-provenance` | partial but brittle | partial but brittle | +0 |
| `SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload` | `promisebase-reference-naming-promisebase-custom-syntax-migration` | good strict baseline | good strict baseline | +0 |
| `SIM-tohol-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-opaque-bytes` | `chunking-identity-bakeoff-profile-negotiation-mismatch` | partial fit, stronger audit baseline | partial fit, stronger audit baseline | +0 |
| `SIM-tohol-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-opaque-bytes` | `municipal-governance` | partial fit | partial fit | +0 |
| `SIM-tohol-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-opaque-bytes` | `promise-economy-spectrum-cryptocurrency-toxicity-failure` | partial guardrail, not an economy solution | partial guardrail, not an economy solution | +0 |
| `SIM-vamaz-grid-envelope-enc-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload` | `aerospace-project-funding` | strong fit | strong fit | +0 |
| `SIM-vamaz-grid-envelope-enc-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload` | `journalism-source-provenance` | strongest fit in this slice | strong fit | -1 |
| `SIM-vamaz-grid-envelope-enc-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload` | `promisebase-reference-naming-promisebase-custom-syntax-migration` | best migration fit | best migration fit | +0 |

## Scenario Aggregates

| Scenario | Old avg score | New avg score | Delta |
|---|---:|---:|---:|
| `aerospace-project-funding` | 3.67 | 3.67 | +0.00 |
| `chunking-identity-bakeoff-profile-negotiation-mismatch` | 2.33 | 2.33 | +0.00 |
| `journalism-source-provenance` | 3.33 | 3.17 | -0.17 |
| `municipal-governance` | 2.00 | 2.00 | +0.00 |
| `promise-economy-spectrum-cryptocurrency-toxicity-failure` | 2.00 | 2.00 | +0.00 |
| `promisebase-reference-naming-promisebase-custom-syntax-migration` | 4.33 | 4.33 | +0.00 |
| `shipping-label-printing` | 3.00 | 3.00 | +0.00 |

## Notes

- This report compares verdict lines from result artifacts; it does not execute protocol harness code.
- Paths and line numbers for each verdict are in the source result files under `results/<sim>/<scenario>/<model>/<timestamp>.md`.
