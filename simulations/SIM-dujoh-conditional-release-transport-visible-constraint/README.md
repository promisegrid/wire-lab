# SIM-dujoh-conditional-release-transport-visible-constraint

This simulation turns the transport/feed-visible constraint alternative from
`SIM-zarud-conditional-release-geofencing` into a concrete candidate specimen.
It tests whether lower layers need enough condition metadata to avoid routing,
replicating, or fetching content in ways that violate release promises. Source:
`DI-fibuv`.

## Design Under Test

Transport and feed messages carry opaque condition references plus small
machine-checkable hints, so peers can decline storage or forwarding when their
local promises would conflict with the condition.

## Boundaries

This simulation does not require transport layers to understand human policy
text. It tests the smallest visible constraint surface that could prevent
obvious violations while keeping payload semantics above transport.
