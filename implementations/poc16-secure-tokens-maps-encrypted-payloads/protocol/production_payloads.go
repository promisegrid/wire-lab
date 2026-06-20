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
			"package_id",
			"exchange_id",
			"weight_ounces",
		},
	},
}

var upsLabelSchemas = []arrayPayloadSchema{
	{
		promiseAbout: productionPromisePrintLabel,
		bodyFields: []string{
			"package_id",
			"exchange_id",
			"shipping_address",
			"weight_ounces",
			"tracking_number",
			"cost_cents",
			"printer_spool_id",
		},
	},
}

var accountingSchemas = []arrayPayloadSchema{
	{
		promiseAbout: productionPromiseAddressLookup,
		bodyFields: []string{
			"order_id",
			"exchange_id",
			"shipping_address",
		},
	},
	{
		promiseAbout: productionPromiseShipmentUpdate,
		bodyFields: []string{
			"order_id",
			"exchange_id",
			"tracking_number",
			"cost_cents",
			"duplicate_shipment_update",
		},
	},
}

var printerPortSchemas = []arrayPayloadSchema{
	{
		promiseAbout: productionPromiseIssuePrintCapability,
		bodyFields: []string{
			"issuee",
			"exchange_id",
			"print_capability_issuee",
			"print_capability_token",
			"print_capability_token_id",
			"print_capability_scope",
			"print_capability_max_bytes",
		},
	},
	{
		promiseAbout: productionPromiseRedeemPrintCapability,
		bodyFields: []string{
			"print_capability_issuee",
			"exchange_id",
			"print_capability_token",
			"print_capability_token_id",
			"print_capability_scope",
			"print_capability_max_bytes",
			"label_bytes_hex",
			"printer_spool_id",
		},
	},
}

var productionShippingSchemas = append(append(append(append([]arrayPayloadSchema{},
	postalScaleSchemas...),
	upsLabelSchemas...),
	accountingSchemas...),
	printerPortSchemas...)

// MarshalPostalScalePayloadFields encodes postal_scale_v1 as a pCID-owned CBOR
// array. Intent: Device payloads should be protocol-owned slot values, not
// generic map messages. Source: DI-dirat
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

// MarshalProductionShippingPayloadFields encodes the active
// production_shipping_v1 protocol-family payload.
// Intent: POC16 now treats weighing, address lookup, label printing,
// printer-port token issue/redeem, and shipment update as operations inside one
// production-shipping protocol family instead of active pCID-per-operation
// fragments. Source: DI-gazin
func MarshalProductionShippingPayloadFields(fields map[string]string) ([]byte, error) {
	return marshalArrayPayload(fields, productionShippingSchemas)
}

// ProductionShippingPayloadFields projects production_shipping_v1 arrays into
// local compatibility fields for the existing shipping workflow handlers.
func ProductionShippingPayloadFields(payloadBytes []byte) (map[string]string, error) {
	return payloadFieldsFromArray(protocolProductionShippingV1, payloadBytes, productionShippingSchemas)
}
