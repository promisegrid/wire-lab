package protocol

import "fmt"

const (
	protocolCASStorageV1 = "cas_storage_v1"
	protocolCIDComputeV1 = "cid_compute_v1"
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

// MarshalKnownArrayPayload encodes the POC14 pCIDs that have been migrated from
// legacy field maps to pCID-owned CBOR arrays.
// Intent: This helper is a migration boundary, not a PromiseGrid-wide payload
// standard; each pCID still owns its own slot order and local compatibility
// projection. Source: DI-gahuh
func MarshalKnownArrayPayload(protocolName string, fields map[string]string) ([]byte, bool, error) {
	switch protocolName {
	case protocolCASStorageV1:
		payloadBytes, err := MarshalCASStoragePayloadFields(fields)
		return payloadBytes, true, err
	case protocolCIDComputeV1:
		payloadBytes, err := MarshalCIDComputePayloadFields(fields)
		return payloadBytes, true, err
	default:
		return nil, false, nil
	}
}

func marshalArrayPayload(fields map[string]string, schemas []arrayPayloadSchema) ([]byte, error) {
	schema, schemaFound := findArrayPayloadSchema(fields["field_promise_about"], schemas)
	if !schemaFound {
		return nil, fmt.Errorf("unsupported pCID-owned array payload promise_about %q", fields["field_promise_about"])
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
		"act":                    "promise",
		"from":                   promiser,
		"to":                     promisee,
		"field_promise_about":    promiseAbout,
		"field_payload_protocol": protocolName,
		"field_payload_shape":    "cbor_array",
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
