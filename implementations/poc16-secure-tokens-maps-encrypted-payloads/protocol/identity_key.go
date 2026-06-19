package protocol

import "fmt"

const protocolIdentityKeyV1 = "identity_key_v1"

const identityKeyRotateSigningKey = "rotate_signing_key"
const identityKeyRotateSigningKeyAck = "rotate_signing_key_ack"

// IdentityKeyRotationPayload is the POC16 identity_key_v1 request payload.
// Intent: `identity_key_v1` is the first POC16 pCID migrated away from the
// earlier generic string-map scaffold into a pCID-owned CBOR array. Source:
// DI-vipih; DI-pusak
type IdentityKeyRotationPayload struct {
	Promiser      string
	Promisee      string
	NewKeyLabel   string
	RotationScope string
}

// IdentityKeyRotationAckPayload is the POC16 identity_key_v1 ACK payload.
// Intent: ACKs for migrated pCIDs should also remain pCID-owned arrays instead
// of reintroducing generic map payloads through ACK construction. Source:
// DI-vipih
type IdentityKeyRotationAckPayload struct {
	Promiser      string
	Promisee      string
	Outcome       string
	PromiseText   string
	NewKeyLabel   string
	RotationScope string
}

// MarshalIdentityKeyRotationPayload encodes:
// [promiser, promisee, "rotate_signing_key", [new_key_label, rotation_scope]].
func MarshalIdentityKeyRotationPayload(payload IdentityKeyRotationPayload) ([]byte, error) {
	writer := &cborWriter{}
	if err := writer.writeArrayHeader(4); err != nil {
		return nil, err
	}
	for _, value := range []string{payload.Promiser, payload.Promisee, identityKeyRotateSigningKey} {
		if err := writer.writeString(value); err != nil {
			return nil, err
		}
	}
	if err := writer.writeArrayHeader(2); err != nil {
		return nil, err
	}
	for _, value := range []string{payload.NewKeyLabel, payload.RotationScope} {
		if err := writer.writeString(value); err != nil {
			return nil, err
		}
	}
	return writer.buffer.Bytes(), nil
}

// MarshalIdentityKeyRotationAckPayload encodes:
// [promiser, promisee, "rotate_signing_key_ack",
// [outcome, promise_text, new_key_label, rotation_scope]].
func MarshalIdentityKeyRotationAckPayload(payload IdentityKeyRotationAckPayload) ([]byte, error) {
	writer := &cborWriter{}
	if err := writer.writeArrayHeader(4); err != nil {
		return nil, err
	}
	for _, value := range []string{payload.Promiser, payload.Promisee, identityKeyRotateSigningKeyAck} {
		if err := writer.writeString(value); err != nil {
			return nil, err
		}
	}
	if err := writer.writeArrayHeader(4); err != nil {
		return nil, err
	}
	for _, value := range []string{payload.Outcome, payload.PromiseText, payload.NewKeyLabel, payload.RotationScope} {
		if err := writer.writeString(value); err != nil {
			return nil, err
		}
	}
	return writer.buffer.Bytes(), nil
}

// IdentityKeyPayloadFields exposes routing and ACK fields from identity_key_v1
// array payloads while the rest of POC16 is still migrating from string maps.
// Intent: These fields are compatibility projections for the current runtime;
// the wire payload remains the pCID-owned array above. Source: DI-vipih
func IdentityKeyPayloadFields(payloadBytes []byte) (map[string]string, error) {
	reader := &cborReader{data: payloadBytes}
	arrayLength, err := reader.readTypeAndLength(4)
	if err != nil {
		return nil, err
	}
	if arrayLength != 4 {
		return nil, fmt.Errorf("identity_key_v1 payload must have 4 slots, got %d", arrayLength)
	}
	promiser, err := reader.readString()
	if err != nil {
		return nil, err
	}
	promisee, err := reader.readString()
	if err != nil {
		return nil, err
	}
	promiseKind, err := reader.readString()
	if err != nil {
		return nil, err
	}
	bodyLength, err := reader.readTypeAndLength(4)
	if err != nil {
		return nil, err
	}
	fields := map[string]string{
		"act":           "promise",
		"from":          promiser,
		"to":            promisee,
		"promise_about": identityKeyRotateSigningKey,
	}
	switch promiseKind {
	case identityKeyRotateSigningKey:
		if bodyLength != 2 {
			return nil, fmt.Errorf("rotate_signing_key body must have 2 slots, got %d", bodyLength)
		}
		newKeyLabel, readErr := reader.readString()
		if readErr != nil {
			return nil, readErr
		}
		rotationScope, readErr := reader.readString()
		if readErr != nil {
			return nil, readErr
		}
		fields["new_key_label"] = newKeyLabel
		fields["rotation_scope"] = rotationScope
	case identityKeyRotateSigningKeyAck:
		if bodyLength != 4 {
			return nil, fmt.Errorf("rotate_signing_key_ack body must have 4 slots, got %d", bodyLength)
		}
		outcome, readErr := reader.readString()
		if readErr != nil {
			return nil, readErr
		}
		promiseText, readErr := reader.readString()
		if readErr != nil {
			return nil, readErr
		}
		newKeyLabel, readErr := reader.readString()
		if readErr != nil {
			return nil, readErr
		}
		rotationScope, readErr := reader.readString()
		if readErr != nil {
			return nil, readErr
		}
		fields["promise_about"] = identityKeyRotateSigningKey
		fields["outcome"] = outcome
		fields["promise"] = promiseText
		fields["new_key_label"] = newKeyLabel
		fields["rotation_scope"] = rotationScope
	default:
		return nil, fmt.Errorf("unsupported identity_key_v1 promise kind %q", promiseKind)
	}
	if reader.offset != len(reader.data) {
		return nil, fmt.Errorf("trailing cbor bytes in identity_key_v1 payload: %d", len(reader.data)-reader.offset)
	}
	return fields, nil
}
