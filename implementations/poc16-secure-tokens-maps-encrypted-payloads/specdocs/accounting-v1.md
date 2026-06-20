# accounting_v1

`accounting_v1` is retained as historical/specimen evidence for the earlier
single-function shipping workflow pCID split. Active POC16 runtime traffic now
uses `production_shipping_v1` for address lookup, shipment update, duplicate
update recognition, and confirmation.

Accounting still promises only its own records and local update outcomes; this
retired specimen must not be treated as an active runtime receive promise after
the parser-role correction. Source: `DI-gazin`.
