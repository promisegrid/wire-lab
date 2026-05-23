# SIM-hopiv-peer-adoption-spec-side-answer-vocabulary

This simulation turns the spec-side answer-vocabulary alternative from
`SIM-dihiz-peer-adoption-metadata` into a concrete candidate specimen. It tests
whether the spec itself should define the allowed open-question answer keys and
values while peers publish compact bindings. Source: `DI-fibuv`.

## Design Under Test

The spec pCID promises the answer vocabulary. Each peer publishes only the
selected answer values, relying on the spec content to define meaning and
migration rules.

## Boundaries

This simulation does not decide the final answer syntax. It tests whether
spec-local vocabularies reduce peer metadata size or make old adoption promises
fragile when specs evolve.
