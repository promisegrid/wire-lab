package protocol

const (
	productionPromiseWeighPackage          = "weigh_package"
	productionPromiseAddressLookup         = "address_lookup"
	productionPromisePrintLabel            = "print_label"
	productionPromiseShipmentUpdate        = "shipment_update"
	productionPromiseIssuePrintCapability  = "issue_print_capability"
	productionPromiseRedeemPrintCapability = "redeem_print_capability"
)

var postalScaleSchemas = []arrayPayloadSchema{
	{
		promiseAbout: productionPromiseWeighPackage,
		bodyFields: []string{
			"field_package_id",
			"field_exchange_id",
			"field_weight_ounces",
		},
	},
}

var upsLabelSchemas = []arrayPayloadSchema{
	{
		promiseAbout: productionPromisePrintLabel,
		bodyFields: []string{
			"field_package_id",
			"field_exchange_id",
			"field_shipping_address",
			"field_weight_ounces",
			"field_tracking_number",
			"field_cost_cents",
			"field_printer_spool_id",
		},
	},
}

var accountingSchemas = []arrayPayloadSchema{
	{
		promiseAbout: productionPromiseAddressLookup,
		bodyFields: []string{
			"field_order_id",
			"field_exchange_id",
			"field_shipping_address",
		},
	},
	{
		promiseAbout: productionPromiseShipmentUpdate,
		bodyFields: []string{
			"field_order_id",
			"field_exchange_id",
			"field_tracking_number",
			"field_cost_cents",
			"field_duplicate_shipment_update",
		},
	},
}

var printerPortSchemas = []arrayPayloadSchema{
	{
		promiseAbout: productionPromiseIssuePrintCapability,
		bodyFields: []string{
			"field_issuee",
			"field_exchange_id",
			"field_print_capability_issuee",
			"field_print_capability_token",
			"field_print_capability_token_id",
			"field_print_capability_scope",
			"field_print_capability_max_bytes",
		},
	},
	{
		promiseAbout: productionPromiseRedeemPrintCapability,
		bodyFields: []string{
			"field_print_capability_issuee",
			"field_exchange_id",
			"field_print_capability_token",
			"field_print_capability_token_id",
			"field_print_capability_scope",
			"field_print_capability_max_bytes",
			"field_label_bytes_hex",
			"field_printer_spool_id",
		},
	},
}

// MarshalPostalScalePayloadFields encodes postal_scale_v1 as a pCID-owned CBOR
// array. Intent: Device payloads should be protocol-owned slot values, not
// generic field-map messages. Source: DI-dirat
func MarshalPostalScalePayloadFields(fields map[string]string) ([]byte, error) {
	return marshalArrayPayload(fields, postalScaleSchemas)
}

func PostalScalePayloadFields(payloadBytes []byte) (map[string]string, error) {
	return payloadFieldsFromArray(protocolPostalScaleV1, payloadBytes, postalScaleSchemas)
}

// MarshalUPSLabelPayloadFields encodes ups_label_v1 as a pCID-owned CBOR array.
// Intent: Label generation and local print results stay under the UPS-label
// protocol body instead of a universal payload map. Source: DI-dirat
func MarshalUPSLabelPayloadFields(fields map[string]string) ([]byte, error) {
	return marshalArrayPayload(fields, upsLabelSchemas)
}

func UPSLabelPayloadFields(payloadBytes []byte) (map[string]string, error) {
	return payloadFieldsFromArray(protocolUPSLabelV1, payloadBytes, upsLabelSchemas)
}

// MarshalAccountingPayloadFields encodes accounting_v1 as a pCID-owned CBOR
// array. Intent: Address and shipment-update records are accounting protocol
// body slots, not app/kernel control fields. Source: DI-dirat
func MarshalAccountingPayloadFields(fields map[string]string) ([]byte, error) {
	return marshalArrayPayload(fields, accountingSchemas)
}

func AccountingPayloadFields(payloadBytes []byte) (map[string]string, error) {
	return payloadFieldsFromArray(protocolAccountingV1, payloadBytes, accountingSchemas)
}

// MarshalPrinterPortPayloadFields encodes printer_port_v1 as a pCID-owned CBOR
// array. Intent: Capability-token issue and redemption are promises by the
// printer-port resource owner, not permission records from a kernel authority.
// Source: DI-dirat
func MarshalPrinterPortPayloadFields(fields map[string]string) ([]byte, error) {
	return marshalArrayPayload(fields, printerPortSchemas)
}

func PrinterPortPayloadFields(payloadBytes []byte) (map[string]string, error) {
	return payloadFieldsFromArray(protocolPrinterPortV1, payloadBytes, printerPortSchemas)
}
