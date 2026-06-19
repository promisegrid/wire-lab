package protocol

const (
	casPromiseStoreContent               = "store_content"
	casPromiseServeContent               = "serve_content"
	casPromiseReplicateContent           = "replicate_content"
	casPromiseServeReplicaContent        = "serve_replica_content"
	casPromiseReplicaTokenLifecycle      = "replica_token_lifecycle"
	casPromisePresentStorageReport       = "present_storage_report"
	casPromiseLabelFutureMalformedReport = "label_future_malformed_report"
	casPromiseUnsupportedVariantProbe    = "unsupported_variant_probe"
)

var casStorageSchemas = []arrayPayloadSchema{
	{
		promiseAbout: casPromiseStoreContent,
		bodyFields: []string{
			"storage_status",
			"exchange_id",
			"content_cid",
			"content_b64",
			"credit_offer",
			"units",
			"capability_token",
			"token_style",
			"bearer_token",
			"replica_peer",
			"replica_token",
		},
	},
	{
		promiseAbout: casPromiseServeContent,
		bodyFields: []string{
			"token_status",
			"exchange_id",
			"content_cid",
			"content_b64",
			"token",
			"token_style",
			"missing_object_probe",
		},
	},
	{
		promiseAbout: casPromiseReplicateContent,
		bodyFields: []string{
			"content_cid",
			"exchange_id",
			"content_b64",
			"issuee",
			"units",
			"replica_token",
		},
	},
	{
		promiseAbout: casPromiseServeReplicaContent,
		bodyFields: []string{
			"token_status",
			"exchange_id",
			"content_cid",
			"content_b64",
			"token",
			"missing_object_probe",
		},
	},
	{
		promiseAbout: casPromiseReplicaTokenLifecycle,
		bodyFields: []string{
			"token_status",
			"exchange_id",
			"content_cid",
			"token",
			"bearer_token",
			"token_style",
			"issuer_peer",
			"redeem_peer",
		},
	},
	{
		promiseAbout: casPromisePresentStorageReport,
		bodyFields: []string{
			"verdict",
			"exchange_id",
			"content_cid",
			"content_b64",
		},
	},
	{
		promiseAbout: casPromiseLabelFutureMalformedReport,
		bodyFields: []string{
			"repair_status",
			"exchange_id",
		},
	},
	{
		promiseAbout: casPromiseUnsupportedVariantProbe,
		bodyFields: []string{
			"variant_status",
			"exchange_id",
		},
	},
}

// MarshalCASStoragePayloadFields encodes cas_storage_v1 as a pCID-owned CBOR
// array. The field names remain local compatibility names only; they are not
// serialized on the wire.
// Intent: CAS storage promises should exercise pCID-owned positional payloads so
// the POC no longer advertises generic maps as the protocol target. Bearer-token
// and sparse-probe fields are still pCID-owned slots inside cas_storage_v1 rather
// than separate message kinds. Source: DI-gahuh; DI-manul
func MarshalCASStoragePayloadFields(fields map[string]string) ([]byte, error) {
	return marshalArrayPayload(fields, casStorageSchemas)
}

// CASStoragePayloadFields projects cas_storage_v1 arrays back into runtime
// fields while handlers are still being migrated incrementally.
func CASStoragePayloadFields(payloadBytes []byte) (map[string]string, error) {
	return payloadFieldsFromArray(protocolCASStorageV1, payloadBytes, casStorageSchemas)
}
