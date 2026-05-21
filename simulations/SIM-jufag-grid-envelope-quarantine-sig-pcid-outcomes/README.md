# SIM-jufag-grid-envelope-quarantine-sig-pcid-outcomes: Grid-envelope variant

This simulation is a standalone positional grid-envelope specimen bred from exactly two parents:

- `SIM-ruzil-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-opaque-bytes`
- `SIM-tohol-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-opaque-bytes`

It keeps the shared parent strengths of canonical `enc-dag-cbor` positional encoding and mandatory signatures, but makes one bounded design move on unknown handling and one bounded design move on signature audibility:

- `unknown-quarantine` instead of either unconditional opaque validity or hard reject
- explicit `sig_pcid` instead of leaving signature semantics entirely off-envelope

This specimen was promoted from review-stage child proposal
`SIM-jufag-child-grid-envelope-quarantine-sig-pcid-outcomes` from
`ga-canary-20260520-221953` under `DI-dipid`. The ignored proposal artifacts
remain local raw evidence; this directory is the canonical non-child simulation
home.

The local draft spec is `protocols/grid-envelope.d/specs/grid-envelope-draft.md`.
