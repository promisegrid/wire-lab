# SIM-rusap: Promise accounting records

This simulation explores peer-local promise accounting records for the
turn-177 promise-economy requirement. Alice, Bob, Carol, and other peers each
keep their own records for relationships and observations; this simulation
does not define a central record store, harness-owned accounting authority, or
global trust database. Source: `DI-navod`.

## Question

What should a peer record locally so it can decide which chunks to pull, keep,
advertise, or refuse while preserving PromiseGrid's 100-year, peer-relative
trust model? Source: `DI-navod`.

## Turn 177 pressure

Turn 177 made the promise economy load-bearing at every layer:

- Sites need promises about storing, sharing, and accurately advertising CAS
  chunks.
- Feed behavior needs incentives and accountability so chunk replication does
  not collapse into spam, capture, or blind flooding.
- Group-level semantics need human-readable promise vocabulary and easy mental
  models.
- Long-lived deployment needs local, repairable records rather than central
  enforcement.

This simulation uses the phrase **promise accounting records** deliberately.
The intended model is peer-local accounting maintained by each actor for its
own relationships, not a central or shared record controlled by the harness.
Source: `DI-navod`.

## Decision axes

- **Record scope:** per peer, per group, per site, per feed, per chunk, or a
  composition of those scopes.
- **Observation shape:** promises made, promises kept, failures, refusals,
  corrupt data, latency, storage duration, and context.
- **Layer interaction:** how L7 group choices guide L6 CAS retention and L5
  feed pull/advertise decisions without centralizing all authority.
- **Human model:** how laypeople understand "sites make promises and keep
  them or do not" without needing protocol internals.
- **Longevity:** what records remain useful under key rotation, migration,
  sparse storage, intermittent connectivity, and changing social norms.

## Boundaries

This simulation does not decide final scoring, reputation, consensus, or
economic settlement mechanics. It exists to test what must be locally recorded
before specs require promise-vocabulary, 100-year pressure-test, or easy
mental-model sections. `TODO-kulih` remains the spec-doc-shape owner for those
documentation requirements. Source: `DI-navod`.
