# Profile negotiation mismatch Matrix

## Authority Boundary

This matrix summarizes evidence. It does not declare a winning design by itself.
Source: `DI-faros`; `DI-vabor`; `DI-dimas`; `DI-nanih`.

| Simulation | Scenario | Latest result run | Status | Notes |
|---|---|---|---|---|
| `SIM-gobaz-chunking-identity-bakeoff` | `chunking-identity-bakeoff-profile-negotiation-mismatch` |  | not-run | Source simulation for the mined row. |
| `SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid` | `chunking-identity-bakeoff-profile-negotiation-mismatch` | `results/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/chunking-identity-bakeoff-profile-negotiation-mismatch/openai-gpt-5.3-codex-xhigh/20260519-025634.md` | run | partial: Fails closed on unknown profile pCIDs; no chunker negotiation semantics. |
| `SIM-tohol-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-opaque-bytes` | `chunking-identity-bakeoff-profile-negotiation-mismatch` | `results/SIM-tohol-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-opaque-bytes/chunking-identity-bakeoff-profile-negotiation-mismatch/openai-gpt-5.3-codex-xhigh/20260519-025634.md` | run | partial: Signature bytes help audit profile advertisement; mismatch recovery remains external. |
| `SIM-kuful-udp-feed-v0-conformance` | `chunking-identity-bakeoff-profile-negotiation-mismatch` | `results/SIM-kuful-udp-feed-v0-conformance/chunking-identity-bakeoff-profile-negotiation-mismatch/openai-gpt-5.3-codex-xhigh/20260519-025634.md` | run | weak: Can prove byte preservation, not chunk-profile identity behavior. |
