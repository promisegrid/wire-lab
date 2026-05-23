# SIM-nalug-chunking-raw-only-first-profile

This simulation turns the raw-only first-profile alternative from
`SIM-gobaz-chunking-identity-bakeoff` into a concrete candidate specimen. It
tests whether the first L6 CAS profile should avoid chunked Merkle roots and use
raw chunks until the chunking identity rule is ready. Source: `DI-fibuv`.

## Design Under Test

The first profile promises only raw-byte content identity and defers chunking
algorithm promises to a later profile or payload family.

## Boundaries

This simulation does not claim raw-only storage is sufficient forever. It tests
whether deferral reduces early protocol debt or prevents useful large-object
replication evidence from emerging.
