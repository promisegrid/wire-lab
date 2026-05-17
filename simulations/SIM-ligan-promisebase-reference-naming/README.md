# SIM-ligan: Promisebase reference naming

This simulation explores promisebase's human-readable symbol pointing at a hash
problem as its own reference-resolution design question. It is split out from
the CAS object-model simulation so reference names can evolve independently from
CBOR profile, chunking, CIDv1 object typing, and pointer-object byte shape.
Source: `DI-tibis`.

## Question

How should PromiseGrid handle human-readable references to CAS roots learned
from promisebase without creating competing identities, mutable-name confusion,
or custom non-CID syntax?
Source: `DI-tibis`.

## Turn 178 pressure

Turn 178 surfaced promisebase as prior art and specifically called out a prior
reference problem: how should a human-readable symbol point at a hash, especially
when the earlier prototype did not use CIDs and invented custom syntax?

This simulation treats that as a separate question from the core CAS object
model. A CAS object model decides bytes, codecs, chunks, and pointer-object
identity. Reference naming decides whether and how humans get stable or mutable
names for roots without confusing the name with the CID it names. Source:
`DI-tibis`.

## Decision Axes

- **Layer home:** L6 reference objects, L7 group/session metadata, a separate
  reference protocol, or explicit rejection/deferment.
- **Mutability:** immutable labels, mutable refs with history, signed updates,
  or time-scoped aliases.
- **Identity boundary:** how to keep a reference name distinct from the pointer
  object CID and the root CID.
- **Conflict handling:** collisions, squatting, divergent local refs, and
  malicious updates.
- **Interop:** whether the shape should align with git refs, IPNS-like naming,
  AT-style handles, or avoid importing any existing naming model prematurely.

## Boundaries

This simulation does not choose the CBOR profile, chunking algorithm, CID codec
set, or pointer-object shape for L6 CAS. Those remain with
`SIM-jomag-cas-object-model` and TODO-kituj / `DR-tumus`. This simulation only
keeps the reference-symbol / hash-name-resolution question visible as its own
design lineage. Source: `DI-tibis`.
