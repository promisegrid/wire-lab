package protocol

// MarshalRelationshipPayloadFields encodes relationship_v1 as a pCID-owned CBOR
// array with an extensible promise-body pair list.
// Intent: Live LLM relationship promises need flexible body details, but fresh
// wire bytes should still be arrays owned by relationship_v1 rather than a
// generic field_* map. Source: DI-dirat
func MarshalRelationshipPayloadFields(fields map[string]string) ([]byte, error) {
	return marshalPairPayload(protocolRelationshipV1, fields)
}

// RelationshipPayloadFields projects relationship_v1 arrays back into runtime
// compatibility fields for existing local trust handlers and analyzer counters.
func RelationshipPayloadFields(payloadBytes []byte) (map[string]string, error) {
	return payloadFieldsFromPairs(protocolRelationshipV1, payloadBytes)
}

// MarshalKernelReceivePayloadFields encodes kernel_receive_v1 app registration
// promises as pCID-owned arrays.
// Intent: App receive-promise registration is still a promise to the local
// kernel, not an out-of-band control map or shared-volume registration file.
// Source: DI-dirat
func MarshalKernelReceivePayloadFields(fields map[string]string) ([]byte, error) {
	if fields["field_promise_about"] == "" {
		fields["field_promise_about"] = "receive_pcid"
	}
	return marshalPairPayload(protocolKernelReceiveV1, fields)
}

// KernelReceivePayloadFields projects kernel_receive_v1 arrays into local
// compatibility fields used by the kernel receive-promise table.
func KernelReceivePayloadFields(payloadBytes []byte) (map[string]string, error) {
	return payloadFieldsFromPairs(protocolKernelReceiveV1, payloadBytes)
}
