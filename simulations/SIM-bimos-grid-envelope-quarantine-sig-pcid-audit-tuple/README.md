# SIM-bimos-grid-envelope-quarantine-sig-pcid-audit-tuple: Grid-envelope variant

This simulation is a standalone positional grid-envelope specimen bred from exactly two parent variants.

It combines the parent strengths into one candidate:

- explicit signature-protocol dispatch via a 4-slot envelope;
- strict no-speculative-parse rejection for unknown protocols;
- exact-byte preservation of rejected unknown envelopes as quarantined evidence.

Variant under test: `enc-dag-cbor` / `unknown-quarantine` / `sig-mandatory-sig-pcid-payload`.

This promoted specimen is expected to improve on the parents under mixed-version
IoT fleets, route-evidence exchange, and future CAS object-family rollout because
unknown envelopes are not silently accepted, but they also are not needlessly
discarded.

This specimen was promoted from review-stage child proposal
`SIM-bimos-child-grid-envelope-quarantine-sig-pcid-audit-tuple` from
`ga-canary-20260520-221953` under `DI-dipid`. The ignored proposal artifacts
remain local raw evidence; this directory is the canonical non-child simulation
home.

The local draft spec is `protocols/grid-envelope.d/specs/grid-envelope-draft.md`.
