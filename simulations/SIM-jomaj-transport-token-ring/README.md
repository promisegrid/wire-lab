# SIM-jomaj-transport-token-ring

This simulation turns the token-ring transport alternative from
`SIM-narok-transport-family-bakeoff` into a concrete candidate specimen. It
tests whether a PromiseGrid transport family should use explicit turn-taking,
per-hop authorization, and stronger ordering than gossip. Source: `DI-fibuv`.

## Design Under Test

Peers pass a logical token that grants the next sender permission to emit
ordered transport messages. Each hop records the promise it made to forward,
wait, or refuse.

## Boundaries

This simulation does not replace UDP-feed or group-session. It tests whether a
ring transport gives useful ordering and accountability or creates fragile
membership and liveness constraints.
