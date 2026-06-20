# production_shipping_v1

`production_shipping_v1` is the active POC16 protocol-family pCID for the
shipping/device workflow. It consolidates package weighing, accounting address
lookup, UPS label generation, printer-port capability issue/redeem, and shipment
update promises under one protocol family instead of treating each operation as a
separate active runtime pCID.

The payload is a pCID-owned CBOR array selected by `promise_about`:

- `weigh_package`
- `address_lookup`
- `print_label`
- `shipment_update`
- `issue_print_capability`
- `redeem_print_capability`

The named `promise_about` value chooses the remaining body slots. Local parser
roles and apps may project those slots into compatibility fields for this POC,
but pCID slot 0 still names the protocol-family spec, not the target app,
operation authority, or a service registry entry.
