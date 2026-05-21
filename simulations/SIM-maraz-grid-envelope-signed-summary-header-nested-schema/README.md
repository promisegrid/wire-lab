# SIM-maraz-grid-envelope-signed-summary-header-nested-schema: Grid-envelope signed summary header probe

This promoted simulation breeds `SIM-janov-grid-envelope-layer-pcid-nested-signed-payload` and `SIM-sajar-grid-envelope-variable-arity-pcid-defined-fields`.

It keeps Janov's nested signed payload idea, but adds a fixed signed outer summary header so generic peers can audit who endorsed carriage, what class of message it is, which CAS objects summarize the claim, and how to recover the inner schema later. The inner payload remains pCID-defined and may still use Sajar-style variable arity.

Goal: improve sparse-knowledge routing, conditional-release audit, and contested evidence handling without requiring every peer to fully understand every inner payload schema.

This specimen was promoted from review-stage child proposal
`SIM-maraz-child-signed-summary-header-nested-schema` from
`ga-canary-20260521-003110` under `DI-fihub`. The ignored proposal artifacts
remain local raw evidence; this directory is the canonical non-child simulation
home.

The local draft spec is `protocols/grid-envelope.d/specs/grid-envelope-draft.md`.
