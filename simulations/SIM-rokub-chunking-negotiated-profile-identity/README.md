# SIM-rokub-chunking-negotiated-profile-identity

This simulation turns the negotiated-profile chunking identity alternative from
`SIM-gobaz-chunking-identity-bakeoff` into a concrete candidate specimen. It
tests whether peers should advertise or negotiate chunking profiles and then
record enough identity evidence to detect mismatches. Source: `DI-fibuv`.

## Design Under Test

Alice and Bob promise supported chunking profiles before exchanging objects.
Each object root records enough profile identity to let Carol verify that a
received object used the promised boundary rules.

## Boundaries

This simulation does not define a final negotiation protocol. It tests whether
profile negotiation is useful flexibility or a source of durable ambiguity.
