# SIM-kabiv-chunking-pcid-driven-identity

This simulation turns the pCID-driven chunking identity alternative from
`SIM-gobaz-chunking-identity-bakeoff` into a concrete candidate specimen. It
tests whether the protocol specification identified by pCID should define the
chunking algorithm and parameters for the objects it governs. Source: `DI-fibuv`.

## Design Under Test

The pCID promises both payload interpretation and the chunking recipe that turns
a byte stream into leaves and roots. Receivers use the pCID to reconstruct the
same chunk boundaries and reject mismatched roots.

## Boundaries

This simulation does not choose a final algorithm such as Rabin or FastCDC. It
tests whether pCID ownership keeps chunk identity deterministic or couples too
much storage policy to protocol semantics.
