# Journalism Source Provenance Matrix

## Authority Boundary

This matrix summarizes evidence. It does not declare a winning design by itself.
Source: `DI-faros`; `DI-vabor`; `DI-dimas`; `DI-botup`; `DI-midif`.

| Simulation | Scenario | Latest result run | Status | Notes |
|---|---|---|---|---|
| `SIM-kovis-grid-envelope-enc-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload` | `journalism-source-provenance` | `results/SIM-kovis-grid-envelope-enc-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload/journalism-source-provenance/openai-gpt-5.3-codex-xhigh/20260519-030134.md` | run | partial brittle: Fail-closed source evidence; weak custody across handler gaps. |
| `SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload` | `journalism-source-provenance` | `results/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/journalism-source-provenance/openai-gpt-5.3-codex-xhigh/20260519-025634.md` | run | privacy-risk partial: Best-effort salvage can leak or overinterpret source evidence. |
| `SIM-vamaz-grid-envelope-enc-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload` | `journalism-source-provenance` | `results/SIM-vamaz-grid-envelope-enc-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload/journalism-source-provenance/openai-gpt-5.3-codex-xhigh/20260519-030134.md` | run | strong: Opaque preservation fits source protection and later authorized interpretation. |
| `SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload` | `journalism-source-provenance` | `results/SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload/journalism-source-provenance/openai-gpt-5.3-codex-xhigh/20260519-025634.md` | run | strict partial: CAS-native fail-closed baseline; brittle for mixed-version custody. |
| `SIM-rakir-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload` | `journalism-source-provenance` | `results/SIM-rakir-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/journalism-source-provenance/openai-gpt-5.3-codex-xhigh/20260519-025634.md` | run | salvage partial: CAS-native inspection helps corrections but risks source exposure. |
| `SIM-riliz-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload` | `journalism-source-provenance` | `results/SIM-riliz-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload/journalism-source-provenance/openai-gpt-5.3-codex-xhigh/20260519-025634.md` | run | strongest: CAS-native opaque preservation best matches protected evidence chains. |
