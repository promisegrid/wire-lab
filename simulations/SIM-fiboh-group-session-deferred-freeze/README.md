# SIM-fiboh-group-session-deferred-freeze

This simulation turns the deferred-freeze alternative from
`SIM-bohof-group-session-freeze-promise` into a concrete candidate specimen. It
tests whether group-session should remain provisional until cryptographic
promise tooling and frozen pCID references exist. Source: `DI-fibuv`.

## Design Under Test

No freeze promise is accepted until the exact artifact set can be named by
frozen pCID references and verified by the intended promise tooling.

## Boundaries

This simulation does not block all group-session work. It tests whether deferral
prevents false certainty or delays useful convergence evidence too long.
