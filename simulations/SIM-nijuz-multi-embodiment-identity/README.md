# SIM-nijuz: Multi-embodiment identity

This simulation is a provisional question home for `FB-zazon` and `FB-robif`:
one logical app spanning heterogeneous host embodiments while preserving one
protocol contract and one signing-key identity story. Source: `DI-ragaz`.

## Question

What worked example can show a browser tab and a long-running plugin/helper as
one logical PromiseGrid app without overstating unresolved identity, storage, or
binding decisions? Source: `DI-ragaz`.

## Decision Axes

- **One app, multiple embodiments:** each component claims only the part of the
  shared pCID-selected contract it actually implements.
- **Portable identity:** a user needs continuity across browser and plugin hosts
  without treating display names or local usernames as identity.
- **Key recipe status:** algorithm, rotation, handshake, and constrained-host
  storage guidance may be provisional and pivotable.
- **Host constraints:** browser tabs, Node helpers, plugins, and native daemons
  have different filesystem, process, and key-storage promises.
- **Auditability:** later peers must understand which component made which claim
  under which key and protocol version.

## Related Root Scenarios

- `scenarios/multi-embodiment-app-identity/multi-embodiment-app-identity.md`
- `scenarios/portable-signing-key-identity/portable-signing-key-identity.md`

## Boundaries

This simulation does not bless a permanent cryptographic primitive or storage
mechanism. It tests provisional guide language for conformance and identity
continuity across constrained hosts. Source: `DI-ragaz`.
