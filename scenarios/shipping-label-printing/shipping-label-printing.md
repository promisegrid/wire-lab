# Shipping Label Printing

## Scenario ID

shipping-label-printing

## Source / Provenance

- Source type: application seed
- Source path: `protocols/wire-lab.d/TODO/TODO-dadub-root-scenario-skeleton-and-seed-catalog.md`
- Source row/title: Seed application catalog entry `shipping-label-printing`
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-midif`;
  `TODO-dadub`; `DI-timah`; `DI-pohaj`; `DI-zapab`

## Purpose

Exercise PromiseGrid design candidates against shipping-label-printing application
pressure: carrier selection, label capability promises, address data, payment,
and chain-of-custody promises.

## Setup

Alice depends on an outcome in the Shipping Label Printing domain. Bob makes promises
about carrier selection, label capability, address data, payment, and
chain-of-custody promises. Carol needs enough evidence to rely on Bob's promise
without having complete global state, and Mallory may exploit stale, missing, or
ambiguous evidence.

## Stimulus

A routine application event becomes contested or incomplete. Some relevant objects,
signatures, observations, or relationship records are delayed, partitioned, stale, or
disputed, and each peer must decide what to accept, retry, downgrade, or escalate using
only local evidence.

## Expected Pressure

The candidate simulation must show which promises, CAS objects, feeds, identity claims,
names, and promise accounting records are needed for this application pressure, while
avoiding hidden global state or a central authority that would make the result non-
comparable.

## POC12 Evidence So Far

`implementations/poc12-production-progress/` now exercises this scenario through
one concrete shipment flow. The fulfillment app asks accounting for Alice's
shipping address, asks the postal scale for the package weight, asks the UPS label
printer for label/tracking/cost evidence, and then sends accounting a shipment
update. The UPS label printer cannot promise label evidence until it first asks
the local `printer_port` resource promiser for a scoped future-print
capability-promise token and redeems that token with bounded label bytes.

The kernel remains transport evidence only: it routes exact pCID-selected bytes
between local apps and peer kernels. Trust, workflow decisions, device behavior,
`not_promised` restraint, and duplicate checkpoint evidence remain app-local.
Source: `DI-timah`; `DI-pohaj`; `DI-zapab`.

## Scenario-Specific Evaluation Questions

- Which agent promises each shipping step: address lookup, package weight, label
  bytes, printer-port access, tracking number, cost, and accounting update?
- Which failures are broken peer promises, which are receiver non-commitment, and
  which are local resource or provider unavailability?
- How does duplicate shipment-update evidence remain auditable without repeatedly
  increasing trust for the same order/tracking/cost checkpoint?
- When should fulfillment stop retrying a step after a receiver returns
  `not_promised`, and what new local evidence would justify trying again?
