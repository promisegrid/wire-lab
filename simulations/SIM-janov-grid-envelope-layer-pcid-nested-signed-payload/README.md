# SIM-janov-grid-envelope-layer-pcid-nested-signed-payload: Grid-envelope nested signed payload probe

This simulation is a standalone grid-envelope specimen. It tests a shared
layer/network pCID in the outer envelope, where that pCID defines a nested
payload structure containing the actual payload pCID, actual payload bytes, and
a signature over the nested payload. It does not claim that this layering is the
canonical PromiseGrid wire format. Source: `DI-joman`.

The local draft spec is
`protocols/grid-envelope.d/specs/grid-envelope-draft.md`.
