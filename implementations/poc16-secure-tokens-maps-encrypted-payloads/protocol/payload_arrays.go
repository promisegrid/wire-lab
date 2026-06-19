package protocol

import (
	"fmt"
	"sort"
)

const (
	protocolRelationshipV1      = "relationship_v1"
	protocolPostalScaleV1       = "postal_scale_v1"
	protocolUPSLabelV1          = "ups_label_v1"
	protocolAccountingV1        = "accounting_v1"
	protocolPrinterPortV1       = "printer_port_v1"
	protocolKernelReceiveV1     = "kernel_receive_v1"
	protocolCASStorageV1        = "cas_storage_v1"
	protocolCIDComputeV1        = "cid_compute_v1"
	protocolRouteV1             = "route_v1"
	protocolSecureCapabilityV1  = "secure_capability_v1"
	protocolEncryptedPayloadV1  = "encrypted_payload_v1"
	protocolParserBuilderRoleV1 = "parser_builder_role_v1"
	protocolMapPayloadProfileV1 = "map_payload_profile_v1"
)

type arrayPayloadState struct {
	Outcome     string
	PromiseText string
	Reason      string
}

type arrayPayloadSchema struct {
	promiseAbout string
	bodyFields   []string
}

// MarshalKnownArrayPayload encodes the POC16 pCIDs that have been migrated from
// earlier generic map payloads to pCID-owned CBOR arrays.
// Intent: This helper is a migration adapter, not a PromiseGrid-wide payload
// standard; each pCID still owns its own slot order and local compatibility
// projection. Source: DI-gahuh
func MarshalKnownArrayPayload(protocolName string, fields map[string]string) ([]byte, bool, error) {
	switch protocolName {
	case protocolRelationshipV1:
		payloadBytes, err := MarshalRelationshipPayloadFields(fields)
		return payloadBytes, true, err
	case protocolPostalScaleV1:
		payloadBytes, err := MarshalPostalScalePayloadFields(fields)
		return payloadBytes, true, err
	case protocolUPSLabelV1:
		payloadBytes, err := MarshalUPSLabelPayloadFields(fields)
		return payloadBytes, true, err
	case protocolAccountingV1:
		payloadBytes, err := MarshalAccountingPayloadFields(fields)
		return payloadBytes, true, err
	case protocolPrinterPortV1:
		payloadBytes, err := MarshalPrinterPortPayloadFields(fields)
		return payloadBytes, true, err
	case protocolKernelReceiveV1:
		payloadBytes, err := MarshalKernelReceivePayloadFields(fields)
		return payloadBytes, true, err
	case protocolCASStorageV1:
		payloadBytes, err := MarshalCASStoragePayloadFields(fields)
		return payloadBytes, true, err
	case protocolCIDComputeV1:
		payloadBytes, err := MarshalCIDComputePayloadFields(fields)
		return payloadBytes, true, err
	case protocolRouteV1:
		payloadBytes, err := MarshalRoutePayloadFields(fields)
		return payloadBytes, true, err
	default:
		return nil, false, nil
	}
}

// PayloadFieldsForProtocolName decodes slot-1 bytes with the local pCID name
// that came from slot 0.
// Intent: Several pCID-owned array payloads intentionally share the same outer
// array grammar, so trial-decoding without the pCID can misclassify
// relationship, route, and kernel-receive messages. Source: DI-pusak
func PayloadFieldsForProtocolName(protocolName string, payloadBytes []byte) (map[string]string, error) {
	fields, fieldsErr := UnmarshalStringMap(payloadBytes)
	if fieldsErr == nil {
		normalizeMapFieldsForProtocolName(protocolName, fields)
		return fields, nil
	}
	switch protocolName {
	case protocolRelationshipV1:
		return RelationshipPayloadFields(payloadBytes)
	case protocolPostalScaleV1:
		return PostalScalePayloadFields(payloadBytes)
	case protocolUPSLabelV1:
		return UPSLabelPayloadFields(payloadBytes)
	case protocolAccountingV1:
		return AccountingPayloadFields(payloadBytes)
	case protocolPrinterPortV1:
		return PrinterPortPayloadFields(payloadBytes)
	case protocolKernelReceiveV1:
		return KernelReceivePayloadFields(payloadBytes)
	case protocolCASStorageV1:
		return CASStoragePayloadFields(payloadBytes)
	case protocolCIDComputeV1:
		return CIDComputePayloadFields(payloadBytes)
	case protocolRouteV1:
		return RoutePayloadFields(payloadBytes)
	case protocolIdentityKeyV1:
		return IdentityKeyPayloadFields(payloadBytes)
	default:
		return nil, fieldsErr
	}
}

func normalizeMapFieldsForProtocolName(protocolName string, fields map[string]string) {
	// Intent: Map payloads remain pCID-owned, but the current POC runtime still
	// expects local compatibility fields named from/to for delivery. This adapter is
	// per-pCID projection, not a universal payload schema. Source: DI-vulit
	if protocolName == protocolEncryptedPayloadV1 {
		if fields["from"] == "" {
			fields["from"] = fields["sender"]
		}
		if fields["to"] == "" {
			fields["to"] = fields["recipient"]
		}
	}
	fields["payload_protocol"] = protocolName
}

type pairPayloadField struct {
	key   string
	value string
}

// marshalPairPayload encodes flexible pCID-owned payloads whose protocol specs
// intentionally allow an extensible body of local key/value promise details.
// Intent: Relationship and kernel-receive payloads need extension room for live
// LLM details and app registration details, but those details are still body
// slots owned by their pCID rather than a universal named map on the wire.
// Source: DI-dirat; DI-pusak
func marshalPairPayload(protocolName string, fields map[string]string) ([]byte, error) {
	promiseAbout := fields["promise_about"]
	if promiseAbout == "" {
		promiseAbout = "local_observation"
	}
	writer := &cborWriter{}
	if err := writer.writeArrayHeader(5); err != nil {
		return nil, err
	}
	for _, value := range []string{fields["from"], fields["to"], promiseAbout} {
		if err := writer.writeString(value); err != nil {
			return nil, err
		}
	}
	if err := writer.writeArrayHeader(4); err != nil {
		return nil, err
	}
	for _, value := range []string{fields["outcome"], fields["promise"], fields["reason"], fields["turn"]} {
		if err := writer.writeString(value); err != nil {
			return nil, err
		}
	}
	pairs := pairPayloadFields(fields)
	if err := writer.writeArrayHeader(len(pairs)); err != nil {
		return nil, err
	}
	for _, pair := range pairs {
		if err := writer.writeArrayHeader(2); err != nil {
			return nil, err
		}
		if err := writer.writeString(pair.key); err != nil {
			return nil, err
		}
		if err := writer.writeString(pair.value); err != nil {
			return nil, err
		}
	}
	_ = protocolName
	return writer.buffer.Bytes(), nil
}

func pairPayloadFields(fields map[string]string) []pairPayloadField {
	pairs := make([]pairPayloadField, 0, len(fields))
	for key, value := range fields {
		if value == "" || pairPayloadCoreField(key) {
			continue
		}
		wireKey := key
		if wireKey == "payload_protocol" || wireKey == "promise_about" {
			continue
		}
		pairs = append(pairs, pairPayloadField{key: wireKey, value: value})
	}
	sort.Slice(pairs, func(leftIndex, rightIndex int) bool {
		return pairs[leftIndex].key < pairs[rightIndex].key
	})
	return pairs
}

func pairPayloadCoreField(key string) bool {
	switch key {
	case "act", "from", "to", "outcome", "promise", "reason", "turn", "protocol":
		return true
	default:
		return false
	}
}

func payloadFieldsFromPairs(protocolName string, payloadBytes []byte) (map[string]string, error) {
	reader := &cborReader{data: payloadBytes}
	arrayLength, err := reader.readTypeAndLength(4)
	if err != nil {
		return nil, err
	}
	if arrayLength != 5 {
		return nil, fmt.Errorf("%s pair payload must have 5 slots, got %d", protocolName, arrayLength)
	}
	promiser, err := reader.readString()
	if err != nil {
		return nil, err
	}
	promisee, err := reader.readString()
	if err != nil {
		return nil, err
	}
	promiseAbout, err := reader.readString()
	if err != nil {
		return nil, err
	}
	stateLength, err := reader.readTypeAndLength(4)
	if err != nil {
		return nil, err
	}
	if stateLength != 4 {
		return nil, fmt.Errorf("%s pair-payload state must have 4 slots, got %d", protocolName, stateLength)
	}
	outcome, err := reader.readString()
	if err != nil {
		return nil, err
	}
	promiseText, err := reader.readString()
	if err != nil {
		return nil, err
	}
	reason, err := reader.readString()
	if err != nil {
		return nil, err
	}
	turn, err := reader.readString()
	if err != nil {
		return nil, err
	}
	fields := map[string]string{
		"act":              "promise",
		"from":             promiser,
		"to":               promisee,
		"promise_about":    promiseAbout,
		"payload_protocol": protocolName,
	}
	for _, optional := range []struct {
		key   string
		value string
	}{
		{key: "outcome", value: outcome},
		{key: "promise", value: promiseText},
		{key: "reason", value: reason},
		{key: "turn", value: turn},
	} {
		if optional.value != "" {
			fields[optional.key] = optional.value
		}
	}
	pairCount, err := reader.readTypeAndLength(4)
	if err != nil {
		return nil, err
	}
	for pairIndex := uint64(0); pairIndex < pairCount; pairIndex++ {
		pairLength, pairLengthErr := reader.readTypeAndLength(4)
		if pairLengthErr != nil {
			return nil, pairLengthErr
		}
		if pairLength != 2 {
			return nil, fmt.Errorf("%s pair payload entry must have 2 slots, got %d", protocolName, pairLength)
		}
		key, keyErr := reader.readString()
		if keyErr != nil {
			return nil, keyErr
		}
		value, valueErr := reader.readString()
		if valueErr != nil {
			return nil, valueErr
		}
		if key != "" && value != "" {
			fields[""+key] = value
		}
	}
	if reader.offset != len(reader.data) {
		return nil, fmt.Errorf("trailing cbor bytes in %s pair payload: %d", protocolName, len(reader.data)-reader.offset)
	}
	return fields, nil
}

func marshalArrayPayload(fields map[string]string, schemas []arrayPayloadSchema) ([]byte, error) {
	schema, schemaFound := findArrayPayloadSchema(fields["promise_about"], schemas)
	if !schemaFound {
		return nil, fmt.Errorf("unsupported pCID-owned array payload promise_about %q", fields["promise_about"])
	}
	state := arrayPayloadState{
		Outcome:     fields["outcome"],
		PromiseText: fields["promise"],
		Reason:      fields["reason"],
	}
	writer := &cborWriter{}
	if err := writer.writeArrayHeader(5); err != nil {
		return nil, err
	}
	for _, value := range []string{fields["from"], fields["to"], schema.promiseAbout} {
		if err := writer.writeString(value); err != nil {
			return nil, err
		}
	}
	if err := writer.writeArrayHeader(3); err != nil {
		return nil, err
	}
	for _, value := range []string{state.Outcome, state.PromiseText, state.Reason} {
		if err := writer.writeString(value); err != nil {
			return nil, err
		}
	}
	if err := writer.writeArrayHeader(len(schema.bodyFields)); err != nil {
		return nil, err
	}
	for _, fieldName := range schema.bodyFields {
		if err := writer.writeString(fields[fieldName]); err != nil {
			return nil, err
		}
	}
	return writer.buffer.Bytes(), nil
}

func payloadFieldsFromArray(protocolName string, payloadBytes []byte, schemas []arrayPayloadSchema) (map[string]string, error) {
	reader := &cborReader{data: payloadBytes}
	arrayLength, err := reader.readTypeAndLength(4)
	if err != nil {
		return nil, err
	}
	if arrayLength != 5 {
		return nil, fmt.Errorf("%s payload must have 5 slots, got %d", protocolName, arrayLength)
	}
	promiser, err := reader.readString()
	if err != nil {
		return nil, err
	}
	promisee, err := reader.readString()
	if err != nil {
		return nil, err
	}
	promiseAbout, err := reader.readString()
	if err != nil {
		return nil, err
	}
	schema, schemaFound := findArrayPayloadSchema(promiseAbout, schemas)
	if !schemaFound {
		return nil, fmt.Errorf("unsupported %s promise_about %q", protocolName, promiseAbout)
	}
	stateLength, err := reader.readTypeAndLength(4)
	if err != nil {
		return nil, err
	}
	if stateLength != 3 {
		return nil, fmt.Errorf("%s state must have 3 slots, got %d", protocolName, stateLength)
	}
	outcome, err := reader.readString()
	if err != nil {
		return nil, err
	}
	promiseText, err := reader.readString()
	if err != nil {
		return nil, err
	}
	reason, err := reader.readString()
	if err != nil {
		return nil, err
	}
	bodyLength, err := reader.readTypeAndLength(4)
	if err != nil {
		return nil, err
	}
	if bodyLength != uint64(len(schema.bodyFields)) {
		return nil, fmt.Errorf("%s body for %s must have %d slots, got %d", protocolName, promiseAbout, len(schema.bodyFields), bodyLength)
	}
	fields := map[string]string{
		"act":              "promise",
		"from":             promiser,
		"to":               promisee,
		"promise_about":    promiseAbout,
		"payload_protocol": protocolName,
	}
	for _, optional := range []struct {
		key   string
		value string
	}{
		{key: "outcome", value: outcome},
		{key: "promise", value: promiseText},
		{key: "reason", value: reason},
	} {
		if optional.value != "" {
			fields[optional.key] = optional.value
		}
	}
	for _, fieldName := range schema.bodyFields {
		value, readErr := reader.readString()
		if readErr != nil {
			return nil, readErr
		}
		if value != "" {
			fields[fieldName] = value
		}
	}
	if reader.offset != len(reader.data) {
		return nil, fmt.Errorf("trailing cbor bytes in %s payload: %d", protocolName, len(reader.data)-reader.offset)
	}
	return fields, nil
}

func findArrayPayloadSchema(promiseAbout string, schemas []arrayPayloadSchema) (arrayPayloadSchema, bool) {
	for _, schema := range schemas {
		if schema.promiseAbout == promiseAbout {
			return schema, true
		}
	}
	return arrayPayloadSchema{}, false
}
