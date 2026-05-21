# CHANGELOG: grid-envelope

A-side CHANGELOG (per TE-liviv) for this simulation-local `grid-envelope`
protocol specimen.

## 2026-05-20

- Promoted specimen
  `SIM-bimos-grid-envelope-quarantine-sig-pcid-audit-tuple` from review-stage
  child proposal `SIM-bimos-child-grid-envelope-quarantine-sig-pcid-audit-tuple`
  under `DI-dipid`.
- Bred from:
  - `SIM-ruzil-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-opaque-bytes`
  - `SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload`
- Introduced `unknown-quarantine` policy: unknown protocol instances are rejected for validity but may be archived and relayed as explicitly rejected opaque evidence.
- Standardized on explicit `sig_pcid` + `sig_payload` signature dispatch.
