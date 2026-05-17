# UDP-Feed v0 Conformance Scenarios

These scenarios are evidence for `TODO-jodon`. They are not an implementation
and not a frozen UDP-feed spec. Source: `DI-pukap`.

## Scenario Matrix

| Scenario | Setup | What to test | Decision pressure |
|---|---|---|---|
| Loopback round trip | Alice sends a 612-byte payload to Bob over local UDP. | Whether the reference implementation preserves bytes and exposes the expected send/receive API. | A minimal reference may be enough for first v0 evidence if vectors lock the bytes. |
| Boundary-size payload | Alice sends exactly 1232 bytes, then 1233 bytes. | Whether the implementation honors the size promise and errors locally before oversize send. | Size behavior should be in vectors before wider conformance claims. |
| Malformed datagram | Bob receives arbitrary bytes that do not parse at higher layers. | Whether UDP-feed passes bytes upward unchanged rather than inventing message semantics. | Binding conformance must stay below session semantics. |
| Simulation artifact output | A simulator-mode send writes an artifact file for the transmitted bytes. | Whether artifact output proves promise 10 without becoming production behavior. | The artifact contract needs to be testable and explicitly scoped. |
| ns-3 two-node path | Alice and Bob communicate through an ns-3-emulated UDP network. | Whether the v0 reference survives non-loopback timing, interface, and packet-capture conditions. | ns-3 may be the evidence that separates a useful specimen from a local toy. |
| Session-layer composition | A minimal group/session message rides above UDP-feed v0. | Whether UDP-feed's API is sufficient for the next layer without leaking binding details. | If composition is required, TODO-jodon's done criteria must include more than UDP round trip. |

## Expected Outputs

- Evidence for TODO-jodon's done criteria and for whether TODO-bihon's ns-3
  scaffold is required before UDP-feed v0 is considered usable.
- A conformance checklist covering implementation API, vectors, artifact writer,
  ns-3 proof, and implementation conformance claims.
