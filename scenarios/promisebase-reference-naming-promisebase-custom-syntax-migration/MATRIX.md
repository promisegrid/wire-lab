# Promisebase custom syntax migration Matrix

## Authority Boundary

This matrix summarizes evidence. It does not declare a winning design by itself.
Source: `DI-faros`; `DI-vabor`; `DI-dimas`; `DI-nanih`.

| Simulation | Scenario | Latest result run | Status | Notes |
|---|---|---|---|---|
| `SIM-ligan-promisebase-reference-naming` | `promisebase-reference-naming-promisebase-custom-syntax-migration` |  | not-run | Source simulation for the mined row. |
| `SIM-kovis-grid-envelope-enc-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload` | `promisebase-reference-naming-promisebase-custom-syntax-migration` | `results/SIM-kovis-grid-envelope-enc-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload/promisebase-reference-naming-promisebase-custom-syntax-migration/openai-gpt-5.5-xhigh/20260519-024746.md` | run | strict baseline: Rejects legacy syntax unless a migration handler is known. |
| `SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload` | `promisebase-reference-naming-promisebase-custom-syntax-migration` | `results/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/promisebase-reference-naming-promisebase-custom-syntax-migration/openai-gpt-5.5-xhigh/20260519-024746.md` | run | risky partial: Best-effort tooling may accidentally preserve known-bad syntax semantics. |
| `SIM-vamaz-grid-envelope-enc-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload` | `promisebase-reference-naming-promisebase-custom-syntax-migration` | `results/SIM-vamaz-grid-envelope-enc-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload/promisebase-reference-naming-promisebase-custom-syntax-migration/openai-gpt-5.5-xhigh/20260519-024746.md` | run | strong: Opaque legacy-byte preservation supports explicit CID-backed migration. |
| `SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload` | `promisebase-reference-naming-promisebase-custom-syntax-migration` | `results/SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload/promisebase-reference-naming-promisebase-custom-syntax-migration/openai-gpt-5.5-xhigh/20260519-024746.md` | run | good strict baseline: CAS-native rejection prevents accidental non-CID syntax acceptance. |
| `SIM-rakir-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload` | `promisebase-reference-naming-promisebase-custom-syntax-migration` | `results/SIM-rakir-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/promisebase-reference-naming-promisebase-custom-syntax-migration/openai-gpt-5.5-xhigh/20260519-024746.md` | run | risky partial: CAS-native salvage helps discovery but risks generic interpretation. |
| `SIM-riliz-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload` | `promisebase-reference-naming-promisebase-custom-syntax-migration` | `results/SIM-riliz-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload/promisebase-reference-naming-promisebase-custom-syntax-migration/openai-gpt-5.5-xhigh/20260519-024746.md` | run | best: CAS-native opaque preservation best supports deliberate CID-backed migration. |
