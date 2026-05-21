# CHANGELOG: grid-envelope

A-side CHANGELOG (per TE-liviv) for this simulation-local `grid-envelope` protocol specimen.

## 2026-05-21

- Initial promoted draft created from `SIM-janov-grid-envelope-layer-pcid-nested-signed-payload` and `SIM-sajar-grid-envelope-variable-arity-pcid-defined-fields`.
- Added a fixed signed outer summary header to expose generic audit-visible fields.
- Retained a nested signed inner payload so payload semantics can still evolve per `content_pcid`.
- Added universal arity and size limits plus `schema_ref` guidance for durable recovery of old schemas.
- Promoted from proposal `SIM-maraz-child-signed-summary-header-nested-schema`
  from `ga-canary-20260521-003110` under `DI-fihub`.
