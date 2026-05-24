# SIM-ravad-freeze-successor-records

This higher-layer simulation tests freeze successor records as promises about
lineage, exact bytes, and future interpretation. It preserves the freeze pressure
from the rejected `hadit`/`jogoh` review without treating freeze metadata as a
base-envelope assertion artifact. Source: `DI-kafiz`.

## Design Under Test

Alice freezes a protocol specimen by publishing exact bytes and a promise about
how she will refer to them later. When Alice later proposes a successor, she
publishes a separate successor record instead of mutating the frozen artifact:

- Alice promises only Alice's own freeze or successor intent.
- Bob and Carol independently decide whether to adopt Alice's successor path.
- Exact frozen bytes remain available for audit and rollback.
- Successor records live under a pCID-selected higher-layer protocol, not in the
  universal grid-envelope header.

## Local Draft Spec

The local draft spec in
`protocols/freeze-successor-record.d/specs/freeze-successor-record-draft.md`
defines a candidate payload protocol for freeze and successor records. The draft
does not freeze a pCID and does not create an external authority.

## Boundaries

This simulation does not freeze any current PromiseGrid protocol and does not
settle the group-session freeze procedure. It only asks whether autonomous
agents can publish durable freeze/successor evidence without generic envelope
assertion artifacts.
