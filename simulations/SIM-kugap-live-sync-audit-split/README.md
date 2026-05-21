# SIM-kugap: Live sync and audit split

This simulation is a provisional question home for `FB-hurit` and `FB-nilat`:
apps that need low-latency reliable live state but still want durable
PromiseGrid audit publication. Source: `DI-ragaz`.

## Question

Should the guide describe a future reliable live-state protocol shape, an
off-grid live channel plus PromiseGrid audit publication pattern, or both as
explicitly provisional? Source: `DI-ragaz`.

## Decision Axes

- **Live transport needs:** reliable, in-order, frame-preserving delivery with
  sub-second latency for CRDT sync or presence.
- **Audit protocol needs:** durable, replayable publication of snapshots,
  milestones, receipts, or break-witnesses.
- **Layer separation:** live state and durable audit may use different pCIDs,
  timings, durability promises, and conformance claims.
- **Group-session role:** group-session may be an audit layer without being the
  live transport.
- **Failure behavior:** dropped live frames, stale audit snapshots, and partition
  recovery must be explicit.

## Related Root Scenario

- `scenarios/live-crdt-audit-publication/live-crdt-audit-publication.md`

## Boundaries

This simulation does not define a reliable low-latency binding. It tests whether
the guide can name the split honestly without implying that `udp-feed` or
`group-session` already solves live CRDT transport. Source: `DI-ragaz`.
