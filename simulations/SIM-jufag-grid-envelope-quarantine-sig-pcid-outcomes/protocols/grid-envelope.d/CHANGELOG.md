# CHANGELOG: grid-envelope

A-side CHANGELOG (per TE-liviv) for this simulation-local `grid-envelope` protocol specimen.

## 2026-05-20

Initial promoted draft for
`SIM-jufag-grid-envelope-quarantine-sig-pcid-outcomes`, promoted from
review-stage child proposal
`SIM-jufag-child-grid-envelope-quarantine-sig-pcid-outcomes` under `DI-dipid`
and bred from:

- `SIM-ruzil-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-opaque-bytes`
- `SIM-tohol-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-opaque-bytes`

Bounded deltas:

1. Replace the parent unknown-policy split with quarantined store/forward semantics.
2. Add explicit `sig_pcid` dispatch for signature verification profile discovery.
3. Make receiver outcomes auditable as accepted, quarantined, or rejected.
