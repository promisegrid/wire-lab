# SIM-numop-transport-gossip

This simulation turns the gossip transport alternative from
`SIM-narok-transport-family-bakeoff` into a concrete candidate specimen. It
tests whether epidemic propagation, IHave/IWant-style exchange, and eventual
convergence should be a PromiseGrid transport family. Source: `DI-fibuv`.

## Design Under Test

Peers promise to advertise, request, relay, and suppress duplicate messages
using compact knowledge summaries instead of strict global ordering.

## Boundaries

This simulation does not choose a final gossip protocol. It tests whether gossip
improves resilience and sparse knowledge handling or makes promise accounting
and delivery evidence too probabilistic.
