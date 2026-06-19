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
			"compute_status",
			"exchange_id",
			"function_cid",
			"function_b64",
			"input_cid",
			"input_b64",
			"context_cid",
			"context_b64",
			"result_cid",
			"result_b64",
			"bad_result_cid",
			"bad_result_b64",
			"credit_offer",
			"units",
			"capacity_probe",
		},
	},
	{
		promiseAbout: computePromiseLookupComputeCache,
		bodyFields: []string{
			"cache_key",
			"exchange_id",
			"cache_status",
			"function_cid",
			"function_b64",
			"input_cid",
			"input_b64",
			"context_cid",
			"context_b64",
			"result_cid",
			"result_b64",
		},
	},
	{
		promiseAbout: computePromiseProvideContext,
		bodyFields: []string{
			"function_cid",
			"exchange_id",
			"input_cid",
			"context_cid",
			"context_b64",
		},
	},
	{
		promiseAbout: computePromiseVerifyResult,
		bodyFields: []string{
			"verdict",
			"exchange_id",
			"subject_peer",
			"subject_result_cid",
			"result_promiser",
			"function_cid",
			"function_b64",
			"input_cid",
			"input_b64",
			"context_cid",
			"context_b64",
			"result_cid",
			"result_b64",
			"disagreement_probe",
		},
	},
}

// MarshalCIDComputePayloadFields encodes cid_compute_v1 as a pCID-owned CBOR
// array. The pCID owns these slots; the compatibility fields exist only so the
// current runtime handlers can keep their local promise vocabulary.
// Intent: Compute promises should be positional protocol payloads about exact
// function/input/context/result CIDs, not universal maps. Source:
// DI-gahuh
func MarshalCIDComputePayloadFields(fields map[string]string) ([]byte, error) {
	return marshalArrayPayload(fields, cidComputeSchemas)
}

// CIDComputePayloadFields projects cid_compute_v1 arrays back into runtime
// fields while the POC keeps existing compute handlers and analyzer event records.
func CIDComputePayloadFields(payloadBytes []byte) (map[string]string, error) {
	return payloadFieldsFromArray(protocolCIDComputeV1, payloadBytes, cidComputeSchemas)
}
