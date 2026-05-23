# SIM-gujiv-conditional-release-hybrid-reference-graph

This simulation turns the hybrid conditional-release graph alternative from
`SIM-zarud-conditional-release-geofencing` into a concrete candidate specimen.
It tests whether human-facing release semantics can stay at the group layer
while lower layers carry opaque references and proof hooks. Source: `DI-fibuv`.

## Design Under Test

Release promises form a reference graph. Group/session payloads explain the
human promise, while lower layers carry stable condition CIDs and local evidence
that the referenced promise was preserved or refused.

## Boundaries

This simulation does not settle the graph object model. It tests whether a
hybrid reference graph preserves layer boundaries better than putting all
conditions in either session messages or transport metadata.
