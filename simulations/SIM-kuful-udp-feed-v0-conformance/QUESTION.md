# Question

Which minimal evidence should prove UDP-feed v0: a Go reference implementation,
test vectors, simulator artifact output, ns-3 end-to-end round trip, and/or a
session-layer composition test? Source: `DI-pukap`; `TODO-jodon`.

Open decision points:

- Should test vectors be authored before or after the reference implementation?
- Is a loopback-only implementation sufficient for v0, or must ns-3 prove the
  binding through an emulated network?
- What simulation-artifact path and metadata prove promise 10 without becoming a
  production API?
- Which conformance claims belong in implementation `CHANGELOG.md` records?
- Does UDP-feed v0 need a group/session layer above it before it counts as a
  useful binding specimen?
