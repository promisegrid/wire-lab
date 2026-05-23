# SIM-tuhas-group-session-two-surface-freeze-gate

This simulation turns the two-surface freeze-gate alternative from
`SIM-bohof-group-session-freeze-promise` into a concrete candidate specimen. It
tests whether outer/feed rules and group-session semantics must freeze
separately, with any merge promise naming both scopes. Source: `DI-fibuv`.

## Design Under Test

A freeze record separates envelope/feed promises from group-session semantic
promises. A merge promise is valid only when it states which surface is frozen
and which remains provisional.

## Boundaries

This simulation does not choose final outer/feed or group-session specs. It
tests whether separate freeze surfaces reduce ambiguity or create excessive
coordination overhead.
