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
			"field_storage_status",
			"field_exchange_id",
			"field_content_cid",
			"field_content_b64",
			"field_credit_offer",
			"field_units",
			"field_capability_token",
			"field_token_style",
			"field_bearer_token",
			"field_replica_peer",
			"field_replica_token",
		},
	},
	{
		promiseAbout: casPromiseServeContent,
		bodyFields: []string{
			"field_token_status",
			"field_exchange_id",
			"field_content_cid",
			"field_content_b64",
			"field_token",
			"field_token_style",
			"field_missing_object_probe",
		},
	},
	{
		promiseAbout: casPromiseReplicateContent,
		bodyFields: []string{
			"field_content_cid",
			"field_exchange_id",
			"field_content_b64",
			"field_issuee",
			"field_units",
			"field_replica_token",
		},
	},
	{
		promiseAbout: casPromiseServeReplicaContent,
		bodyFields: []string{
			"field_token_status",
			"field_exchange_id",
			"field_content_cid",
			"field_content_b64",
			"field_token",
			"field_missing_object_probe",
		},
	},
	{
		promiseAbout: casPromiseReplicaTokenLifecycle,
		bodyFields: []string{
			"field_token_status",
			"field_exchange_id",
			"field_content_cid",
			"field_token",
			"field_bearer_token",
			"field_token_style",
			"field_issuer_peer",
			"field_redeem_peer",
		},
	},
	{
		promiseAbout: casPromisePresentStorageReport,
		bodyFields: []string{
			"field_verdict",
			"field_exchange_id",
			"field_content_cid",
			"field_content_b64",
		},
	},
	{
		promiseAbout: casPromiseLabelFutureMalformedReport,
		bodyFields: []string{
			"field_repair_status",
			"field_exchange_id",
		},
	},
	{
		promiseAbout: casPromiseUnsupportedVariantProbe,
		bodyFields: []string{
			"field_variant_status",
			"field_exchange_id",
		},
	},
}

// MarshalCASStoragePayloadFields encodes cas_storage_v1 as a pCID-owned CBOR
// array. The field names remain local compatibility names only; they are not
// serialized on the wire.
// Intent: CAS storage promises should exercise pCID-owned positional payloads so
// the POC no longer advertises field_* maps as the protocol target. Bearer-token
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
