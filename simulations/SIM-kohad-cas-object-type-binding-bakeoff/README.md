# SIM-kohad: CAS object type binding bakeoff

This simulation explores how PromiseGrid should bind object type into CAS
identity. It exists because `DR-tumus` DF-tumus.2 was not ready for a direct
answer: Steve asked for simulations to test the type-binding choice first. It
is a standalone design-point simulation, not a frozen codec registry and not a
shared protocol bundle for other simulations. Source: `DI-bukoh`.

## Question

How should PromiseGrid distinguish raw chunks, Merkle nodes, pointer objects,
and future application-shaped CAS objects without creating multiple sources of
truth for object identity? Source: `DI-bukoh`; `DR-tumus`.

## Candidate Bindings

- **CID codec only:** The CIDv1 codec / multicodec value is the object-type
  discriminator.
- **CID codec plus internal kind:** The CID codec identifies a broad object
  family, and the object bytes carry a kind field for finer distinctions.
- **Path suffix negative control:** Filename or path suffixes are tested as a
  deliberate negative control because they do not travel with content identity.

## Boundaries

This simulation does not allocate multicodec values or choose a final type
system. It gives `TODO-kituj` / TE-43 pressure cases for answering `DR-tumus`
while keeping path names, local storage layout, and object identity separate.
Source: `DI-bukoh`.
