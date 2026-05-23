# SIM-bilam-udp-feed-layer-composition-conformance

This simulation turns the layer-composition UDP-feed conformance alternative
from `SIM-kuful-udp-feed-v0-conformance` into a concrete candidate specimen. It
tests whether UDP-feed is only useful once a session protocol can ride above it
and record its own promises. Source: `DI-fibuv`.

## Design Under Test

UDP-feed conformance requires at least one higher-layer session or feed use case
to pass through the binding and produce auditable composition evidence.

## Boundaries

This simulation does not choose the higher-layer protocol. It tests whether
composition proof prevents isolated transport success from being overvalued or
sets the v0 bar too high.
