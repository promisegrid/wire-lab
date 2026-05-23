# SIM-votoj-peer-adoption-structured-object

This simulation turns the structured adoption object alternative from
`SIM-dihiz-peer-adoption-metadata` into a concrete candidate specimen. It tests
whether a peer's promise to follow a spec pCID plus open-question answers should
be recorded in a content-addressed object. Source: `DI-fibuv`.

## Design Under Test

The adoption object records peer identity, spec pCID, answer set, time or scope,
and signature evidence. Other peers fetch and verify that object before
interpreting the peer's protocol behavior.

## Boundaries

This simulation does not define the final adoption object schema. It tests
whether a structured object gives durable auditability or becomes too heavy for
common peer declarations.
