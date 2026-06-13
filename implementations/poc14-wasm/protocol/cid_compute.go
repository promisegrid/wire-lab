package protocol

const (
	computePromiseExecuteFunction    = "execute_function"
	computePromiseLookupComputeCache = "lookup_compute_cache"
	computePromiseProvideContext     = "provide_compute_context"
	computePromiseVerifyResult       = "verify_compute_result"
)

var cidComputeSchemas = []arrayPayloadSchema{
	{
		promiseAbout: computePromiseExecuteFunction,
		bodyFields: []string{
			"field_compute_status",
			"field_function_cid",
			"field_function_b64",
			"field_input_cid",
			"field_input_b64",
			"field_context_cid",
			"field_context_b64",
			"field_result_cid",
			"field_result_b64",
			"field_bad_result_cid",
			"field_bad_result_b64",
			"field_credit_offer",
			"field_units",
			"field_capacity_probe",
		},
	},
	{
		promiseAbout: computePromiseLookupComputeCache,
		bodyFields: []string{
			"field_cache_key",
			"field_cache_status",
			"field_function_cid",
			"field_function_b64",
			"field_input_cid",
			"field_input_b64",
			"field_context_cid",
			"field_context_b64",
			"field_result_cid",
			"field_result_b64",
		},
	},
	{
		promiseAbout: computePromiseProvideContext,
		bodyFields: []string{
			"field_function_cid",
			"field_input_cid",
			"field_context_cid",
			"field_context_b64",
		},
	},
	{
		promiseAbout: computePromiseVerifyResult,
		bodyFields: []string{
			"field_verdict",
			"field_subject_peer",
			"field_subject_result_cid",
			"field_result_promiser",
			"field_function_cid",
			"field_function_b64",
			"field_input_cid",
			"field_input_b64",
			"field_context_cid",
			"field_context_b64",
			"field_result_cid",
			"field_result_b64",
			"field_disagreement_probe",
		},
	},
}

// MarshalCIDComputePayloadFields encodes cid_compute_v1 as a pCID-owned CBOR
// array. The pCID owns these slots; the compatibility fields exist only so the
// current runtime handlers can keep their local promise vocabulary.
// Intent: Compute promises should be positional protocol payloads about exact
// function/input/context/result CIDs, not universal field maps. Source:
// DI-gahuh
func MarshalCIDComputePayloadFields(fields map[string]string) ([]byte, error) {
	return marshalArrayPayload(fields, cidComputeSchemas)
}

// CIDComputePayloadFields projects cid_compute_v1 arrays back into runtime
// fields while the POC keeps existing compute handlers and analyzer event records.
func CIDComputePayloadFields(payloadBytes []byte) (map[string]string, error) {
	return payloadFieldsFromArray(protocolCIDComputeV1, payloadBytes, cidComputeSchemas)
}
