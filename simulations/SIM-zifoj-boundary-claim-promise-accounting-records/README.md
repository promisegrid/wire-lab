# SIM-zifoj-boundary-claim-promise-accounting-records: Boundary claim package with promise-accounting records

This promoted simulation combines the parent strengths on layer-boundary discipline and honest provisional conformance, then makes them more operational.

## Why this specimen exists

The parents are strong at saying what not to overclaim, but weak at giving a concrete claim package that Bob, Carol, and a future porter can actually inspect. This specimen adds a small auditable structure without pretending unsettled drafts are frozen.

## Core design move

### Delta 1: three claim cards, not one blob of prose
Every implementation publishes up to three separate cards:

1. **Port Boundary Card** — what runtime/dispatcher/library surface is actually implemented now.
2. **App Semantics Card** — what protocol-boundary semantics are claimed now.
3. **Host / Embodiment / Dependency Card** — what remains host-local, live-channel-local, device-local, or embodiment-specific.

This preserves the parents' boundary clarity while reducing ambiguity in kernel-porting, app partial conformance, multi-embodiment identity, live CRDT, and device-agent cases.

### Delta 2: minimal promise-accounting records
Each actor keeps local rows for the promises they make, rely on, fail to meet, or witness failing. The records are lightweight but durable enough for later audit.

### Delta 3: status tags per field
Each claimed item is marked:

- `frozen`
- `provisional`
- `blocked`
- `local-only`

Each non-frozen item must carry a short migration note.

## Expected improvement over parents

- Gives **exact honest partial-conformance wording** instead of only warning against overclaiming.
- Separates **content identity** from **availability, authorization, retention, and discovery**.
- Gives a reusable record shape for **effect dedupe**, **break-witness**, **snapshot citation**, and **key rotation**.
- Keeps harness apparatus, draft live transport, adapter-local signatures, and host dependencies from being mistaken for core PromiseGrid conformance.

## Boundaries

This simulation still does **not** define:

- the final kernel/runtime term of art,
- a frozen app contract,
- a universal witness format,
- a final live transport binding,
- a universal capability-token standard.

It defines only a safer interim packaging for claims and evidence.

This specimen was promoted from review-stage child proposal
`SIM-zifoj-child-boundary-claim-ledger` from `ga-canary-20260521-210902` under
`DI-lanuz`. The ignored proposal artifacts remain local raw evidence; this
directory is the canonical non-child simulation home. Canonical prose uses
peer-local promise-accounting records rather than the proposal's `ledger`
shorthand.
