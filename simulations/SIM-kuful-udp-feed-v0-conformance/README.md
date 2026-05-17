# SIM-kuful: UDP-Feed v0 Conformance

This simulation captures the protocol/specimen questions in `TODO-jodon`: what
minimal reference implementation, test-vector, simulation-artifact, and ns-3
harness surface proves UDP-feed v0 well enough to support later PromiseGrid
design work. It is a standalone design-point simulation, not an implementation
and not a frozen UDP-feed spec. Source: `DI-pukap`.

## Question

What conformance surface should UDP-feed v0 require before the lineage is useful
as a binding specimen: Go reference implementation, deterministic test vectors,
simulator artifact writer, ns-3 round trip, or a smaller/larger combination?
Source: `DI-pukap`; `TODO-jodon`.

## Candidate Shapes

- **Reference-first:** A minimal Go implementation establishes behavior, then
  test vectors and ns-3 scenarios validate it.
- **Vector-first:** Test vectors define the behavior before the implementation
  claims conformance.
- **Harness-first:** An ns-3 round trip proves the binding layer in realistic
  network conditions before polishing API shape.
- **Layer-composition proof:** UDP-feed is not done until a session protocol can
  ride above it and record conformance claims.

## Boundaries

This simulation does not write the Go implementation, test vectors, or ns-3
harness. It captures the conformance-design question so `TODO-jodon` can remain
implementation-focused while still exploring what evidence should count as a
usable v0 binding artifact. Source: `DI-pukap`.
