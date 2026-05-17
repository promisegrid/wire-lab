# SIM-gobaz: Chunking identity bakeoff

This simulation explores how chunking choices should become part of PromiseGrid
object identity. It exists because `DR-tumus` DF-tumus.3 was not ready for a
direct answer: Steve asked for simulations and raised whether chunking should
be pCID-driven, have a separate chunking CID, or use a similar mechanism such
as a possible cCID. This simulation treats those names as exploratory, not
settled protocol terms. Source: `DI-bukoh`.

## Question

Where should PromiseGrid bind the chunking algorithm and parameters that turn a
byte stream into CAS leaves and Merkle roots? Source: `DI-bukoh`; `DR-tumus`.

## Candidate Bindings

- **pCID-driven chunking:** The protocol spec identified by pCID defines the
  chunking algorithm and parameters for objects under that protocol.
- **Chunking CID / cCID-style identity:** A separate content-addressed or
  protocol-addressed chunking descriptor identifies algorithm and parameters.
- **Profile-negotiated chunking:** Peers negotiate or advertise a named profile,
  and object identity records enough information to detect mismatches.
- **Raw-only first profile:** The first L6 profile avoids chunked Merkle roots
  and uses raw chunks until the chunking identity rule is ready.

## Boundaries

This simulation does not coin a final `cCID` term and does not choose Rabin,
FastCDC, or any parameter set. It exists so `TODO-kituj` / TE-43 can decide
whether chunking identity belongs in a pCID, a distinct descriptor, a negotiated
profile, or a deferred later profile. Source: `DI-bukoh`.
