# SIM-vopit-transport-receipts-at-scale

This simulation turns the receipts-at-scale transport alternative from
`SIM-narok-transport-family-bakeoff` into a concrete candidate specimen. It
tests whether large-N acknowledgement evidence should use vectors, summaries,
compact proofs, or other receipt aggregation structures. Source: `DI-fibuv`.

## Design Under Test

Transport messages produce local receipt promises that can be summarized without
requiring every peer to store every individual acknowledgement forever.

## Boundaries

This simulation does not settle a receipt-vector format. It tests whether
receipt aggregation improves auditability at scale or creates unverifiable
summary claims that weaken peer-local evidence.
