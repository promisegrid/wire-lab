# SIM-mikas: Minimal blob app contract

This simulation is a provisional question home for the App Dev feedback item
`FB-vopik`: what a minimal immutable blob app may safely assume about
PromiseGrid-facing CAS behavior. It is not a final CAS profile, SDK, or app API.
Source: `DI-ragaz`.

## Question

What can the guide safely tell an app developer whose whole first app is
`blob in -> hash out` and `hash in -> blob out`, while `DR-tuhaz` and
`DR-tumus` remain open? Source: `DI-ragaz`.

## Decision Axes

- **Write/read contract:** whether a draft profile can promise content-addressed
  identity without promising universal availability.
- **Replication boundary:** what caching, pinning, retention, or pull behavior is
  host/app policy rather than base PromiseGrid semantics.
- **Discovery and ingress:** how much node discovery, rendezvous, and ingress
  remains outside the blob app contract.
- **Hash-as-capability:** whether possession of a hash is a PromiseGrid read
  capability, an app convention, or only an addressing fact.
- **100-year survival:** what survives after hosts, storage policies, keys, and
  CAS implementations change.

## Related Root Scenario

- `scenarios/minimal-immutable-blob-app/minimal-immutable-blob-app.md`

## Boundaries

This simulation should compare candidate CAS profiles and app-contract prose
without treating any one CAS implementation, promisebase branch, or hash-access
policy as canonical. Source: `DI-ragaz`.
