# ups_label_v1

`ups_label_v1` is retained as historical/specimen evidence for the earlier
single-function shipping workflow pCID split. Active POC16 runtime traffic now
uses `production_shipping_v1` for label-printing app promises.

The label printer uses package facts, address facts, and printer-port capability
tokens to promise a shipping label, tracking number, and cost event from its own
local workflow. This retired specimen must not be treated as an active runtime
receive promise after the parser-role correction. Source: `DI-gazin`.
