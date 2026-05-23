# SIM-falun-conditional-release-separate-protocol-family

This simulation turns the separate conditional-release family alternative from
`SIM-zarud-conditional-release-geofencing` into a concrete candidate specimen.
It tests whether release conditions should be a dedicated PromiseGrid protocol
that group/session, feed, and application protocols can cite. Source: `DI-fibuv`.

## Design Under Test

Alice publishes a conditional-release promise object under its own pCID. Other
protocols reference that object instead of embedding release policy directly in
their own message formats.

## Boundaries

This simulation does not define a final policy language. It tests whether a
separate protocol family improves reuse and auditability or adds too much
indirection for ordinary content transfer.
