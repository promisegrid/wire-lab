# postal_scale_v1

`postal_scale_v1` is retained as historical/specimen evidence for the earlier
single-function shipping workflow pCID split. Active POC16 runtime traffic now
uses `production_shipping_v1` for package-weight promises made by a postal scale
agent.

The scale promises only its own measurement event and never promises address
lookup, shipment cost, label printing, or accounting updates. This retired
specimen must not be treated as an active runtime receive promise after the
parser-role correction. Source: `DI-gazin`.
