# SIM-gadol-peer-adoption-promise-message

This simulation turns the peer-adoption promise-message alternative from
`SIM-dihiz-peer-adoption-metadata` into a concrete candidate specimen. It tests
whether adoption should be expressed as an ordinary PromiseGrid promise message
that can be stored, relayed, superseded, and accounted for. Source: `DI-fibuv`.

## Design Under Test

Alice publishes a promise message: she promises to behave according to pCID X
with a listed set of answer bindings. Later messages can supersede that promise
under normal promise-accounting rules.

## Boundaries

This simulation does not choose a final promise-message envelope. It tests
whether adoption should reuse the ordinary promise path instead of a special
metadata object.
