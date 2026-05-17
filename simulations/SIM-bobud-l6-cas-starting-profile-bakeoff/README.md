# SIM-bobud: L6 CAS starting profile bakeoff

This simulation explores which starting profile should guide the first
PromiseGrid L6 CAS spec. It exists because `DR-tumus` DF-tumus.1 was not ready
for a direct answer: Steve asked for simulations to test the starting-profile
choice first. It is a standalone design-point simulation, not a frozen L6 CAS
spec and not a shared protocol bundle for other simulations. Source:
`DI-bukoh`.

## Question

Which first L6 CAS profile gives PromiseGrid enough concrete structure to test
turn-177's CAS model without prematurely locking the whole long-term storage
architecture? Source: `DI-bukoh`; `DR-tumus`.

## Candidate Profiles

- **IPFS / IPLD-aligned:** Start near DAG-CBOR, CIDv1, and public multicodec
  conventions so bridgeability is easy to evaluate early.
- **Promisebase adapter:** Start from promisebase / pitbase block, tree, and
  stream prior art, then map it into PromiseGrid pointer and CID semantics.
- **Minimal pointer/raw:** Start with raw chunks plus a minimal CBOR pointer
  object, deferring chunked Merkle objects and implementation substrate
  commitments until the migration path proves the need.

## Boundaries

This simulation does not choose the starting profile. It documents the pressure
that `TODO-kituj` / TE-43 must consider before answering `DR-tumus`. It should
not imply that any profile is a canonical PromiseGrid API, a preferred harness
home, or a dependency for another simulation. Source: `DI-bukoh`.
