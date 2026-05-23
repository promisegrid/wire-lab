# SIM-lofij-peer-adoption-hybrid-pointer

This simulation turns the hybrid peer-adoption pointer alternative from
`SIM-dihiz-peer-adoption-metadata` into a concrete candidate specimen. It tests
whether a compact peer promise should point at a richer adoption object or
answer-profile object when more detail is needed. Source: `DI-fibuv`.

## Design Under Test

The peer publishes a small adoption promise naming the spec pCID and a pointer
to an optional richer object containing answer bindings, scope, freshness, and
signature evidence.

## Boundaries

This simulation does not decide which fields are mandatory. It tests whether a
hybrid pointer avoids both oversized inline promises and under-specified compact
adoption records.
