# SIM-nonib-udp-feed-harness-first-conformance

This simulation turns the harness-first UDP-feed conformance alternative from
`SIM-kuful-udp-feed-v0-conformance` into a concrete candidate specimen. It tests
whether an ns-3 or equivalent network harness should prove the binding layer
under realistic packet loss, timing, and topology pressure before polishing API
shape. Source: `DI-fibuv`.

## Design Under Test

The v0 conformance surface begins with simulator artifacts that reproduce
network behavior and require implementations to record packet-level evidence for
delivery, loss, and recovery promises.

## Boundaries

This simulation does not build the ns-3 harness. It tests whether harness-first
evidence is necessary for a transport binding or too expensive for v0.
