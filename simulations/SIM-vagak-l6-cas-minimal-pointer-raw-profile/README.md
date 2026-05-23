# SIM-vagak-l6-cas-minimal-pointer-raw-profile

This simulation turns the minimal pointer/raw starting-profile alternative from
`SIM-bobud-l6-cas-starting-profile-bakeoff` into a concrete candidate specimen.
It tests whether PromiseGrid should begin L6 CAS with raw chunks plus a small
CBOR pointer object while deferring Merkle DAG structure and substrate-specific
commitments. Source: `DI-fibuv`.

## Design Under Test

The first L6 CAS profile promises that raw bytes and a minimal deterministic
pointer object are enough to test fetch, keep, advertise, and reference
promises before adding chunked Merkle roots or richer object families.

## Boundaries

This simulation does not reject future DAG-CBOR, IPLD, or chunked Merkle
profiles. It tests whether the smallest useful CAS profile reduces migration
debt and keeps early implementations durable.
