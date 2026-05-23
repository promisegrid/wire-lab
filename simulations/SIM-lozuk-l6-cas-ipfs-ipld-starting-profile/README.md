# SIM-lozuk-l6-cas-ipfs-ipld-starting-profile

This simulation turns the IPFS / IPLD-aligned starting-profile alternative from
`SIM-bobud-l6-cas-starting-profile-bakeoff` into a concrete candidate specimen.
It tests whether PromiseGrid L6 CAS should begin near CIDv1, multicodec,
DAG-CBOR, and IPLD traversal conventions so bridgeability and prior tooling can
be evaluated early. Source: `DI-fibuv`.

## Design Under Test

The first L6 CAS profile promises that stored objects use CIDv1-compatible
identity, deterministic DAG-CBOR where structured objects are needed, and
explicit pCID-owned payload semantics above the CAS object layer.

## Boundaries

This simulation does not freeze the CAS profile or require IPFS as the runtime
substrate. It exists to test whether reuse of IPLD-shaped conventions improves
interoperability without importing unwanted complexity or centralized registry
assumptions.
