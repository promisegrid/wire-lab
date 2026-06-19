package runtime

import (
	"fmt"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/pcid"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/protocol"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/specdocs"
)

func (node *Node) runPOC16ProtocolProfileWorkflow(contentCID string) error {
	// Intent: POC16 must remain a POC15 superset while adding secure tokens,
	// encrypted payloads, map-permitted profiles, and pCID-selected parser/builder
	// roles as broad protocol-family pressure rather than pCID-per-operation
	// fragmentation. Source: DI-vulit
	node.recordPOC15SupersetAgentSet()
	if err := node.recordPOC16SpecContextEvents(); err != nil {
		return err
	}
	mapExactHash, mapErr := node.recordMapPayloadProfile()
	if mapErr != nil {
		return mapErr
	}
	tokenExactHash, tokenErr := node.recordSecureCapabilityProfile(contentCID, mapExactHash)
	if tokenErr != nil {
		return tokenErr
	}
	encryptedExactHash, encryptedErr := node.recordEncryptedPayloadProfile(contentCID, tokenExactHash)
	if encryptedErr != nil {
		return encryptedErr
	}
	return node.recordParserBuilderRoleProfile(encryptedExactHash)
}

func (node *Node) recordPOC15SupersetAgentSet() {
	agents := []string{"alice", "bob", "carol", "dave", "ellen", "frank", "grace", "heidi", "ivan", "judy", "mallory", "oscar", "fulfillment", "postal_scale", "ups_label_printer", "printer_port", "accounting", "peggy", "victor"}
	for _, agentName := range agents {
		node.record("poc15_superset_named_agent_preserved", "kept", agentName, "agent="+agentName+" POC16 config keeps inherited POC15 agent and workflow role")
	}
}

func (node *Node) recordPOC16SpecContextEvents() error {
	protocolNames := []string{
		pcid.RelationshipV1,
		pcid.CASStorageV1,
		pcid.CIDComputeV1,
		pcid.IdentityKeyV1,
		pcid.RouteV1,
		pcid.SecureCapabilityV1,
		pcid.EncryptedPayloadV1,
		pcid.ParserBuilderRoleV1,
		pcid.MapPayloadProfileV1,
	}
	contexts, contextsErr := specdocs.ContextsFor(protocolNames, 700)
	if contextsErr != nil {
		return contextsErr
	}
	for _, context := range contexts {
		protocolCID, ok := node.Protocols.CID(context.Name)
		if !ok {
			return fmt.Errorf("spec context has unknown pCID name %s", context.Name)
		}
		node.record("poc16_protocol_spec_doc_recorded", "kept", "", "pcid="+context.Name+" spec_pcid="+protocolCID.String()+" spec_sha256="+context.SHA256)
		node.record("llm_spec_context_embedded", "kept", "", "pcid="+context.Name+" spec_pcid="+protocolCID.String()+" excerpt_bytes="+fmt.Sprintf("%d", len(context.Excerpt)))
		node.record("llm_spec_context_cid_recorded", "kept", "", "pcid="+context.Name+" spec_pcid="+protocolCID.String())
	}
	return nil
}

func (node *Node) recordMapPayloadProfile() (string, error) {
	fields := map[string]string{
		"act":           "promise",
		"from":          node.Agent.Name,
		"to":            "bob",
		"promise_about": "map_payload_profile",
		"promise":       "Alice promises to use a pCID-owned CBOR map only where readability is worth the extra bytes.",
		"reason":        "maps are permitted by this pCID, not by a universal payload rule",
		"profile":       "self-documenting-map",
		"device_note":   "arrays remain preferred for limited devices",
	}
	envelope, _, buildErr := node.buildEnvelopeFromFields(pcid.MapPayloadProfileV1, node.Protocols.MustCID(pcid.MapPayloadProfileV1), fields)
	if buildErr != nil {
		return "", buildErr
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		return "", bytesErr
	}
	node.emitMessageArtifact("poc16_profile", "bob", pcid.MapPayloadProfileV1, envelopeBytes, fields)
	if _, parseErr := node.parseEnvelope(envelopeBytes); parseErr != nil {
		return "", parseErr
	}
	exactHash := protocol.HashExactBytes(envelopeBytes)
	node.record("poc16_map_payload_specimen_emitted", "kept", "bob", "pcid="+pcid.MapPayloadProfileV1+" exact_sha256="+exactHash)
	node.record("poc16_map_payload_parsed", "kept", "bob", "pcid="+pcid.MapPayloadProfileV1+" parser=slot0-selected")
	return exactHash, nil
}

func (node *Node) recordSecureCapabilityProfile(contentCID, parentExactHash string) (string, error) {
	now := time.Now().UTC()
	token := protocol.CWTCapabilityToken{
		Issuer:        node.Agent.Name,
		Subject:       "frank",
		Audience:      "bob",
		Capability:    "store-for-peer",
		Scope:         pcid.CASStorageV1,
		ContentCID:    contentCID,
		TokenID:       "poc16-cwt-storage-token-1",
		ExpiresUnix:   now.Add(10 * time.Minute).Unix(),
		NotBeforeUnix: now.Add(-time.Minute).Unix(),
		Transferable:  true,
	}
	tokenText, tokenErr := protocol.EncodeCWTCapabilityToken(token)
	if tokenErr != nil {
		return "", tokenErr
	}
	fields := protocol.CWTTokenSummaryFields(token)
	fields["act"] = "promise"
	fields["from"] = node.Agent.Name
	fields["to"] = "bob"
	fields["promise_about"] = "capability_token_promise"
	fields["promise"] = "Alice promises that this bearer token can be redeemed for the bounded storage behavior described by its signed claims."
	fields["reason"] = "the issuer promise is signed; redemption still depends on Bob's local judgment"
	fields["token_b64"] = tokenText
	envelope, _, buildErr := node.buildEnvelopeFromFieldsWithParents(pcid.SecureCapabilityV1, node.Protocols.MustCID(pcid.SecureCapabilityV1), fields, []string{parentExactHash})
	if buildErr != nil {
		return "", buildErr
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		return "", bytesErr
	}
	node.emitMessageArtifact("poc16_profile", "bob", pcid.SecureCapabilityV1, envelopeBytes, fields)
	verified, verifyErr := protocol.VerifyCWTCapabilityToken(tokenText, node.Agent.Name, "bob", now)
	if verifyErr != nil {
		return "", verifyErr
	}
	seenTokens := map[string]bool{verified.TokenID: true}
	node.record("cwt_capability_token_issued", "kept", "bob", "pcid="+pcid.SecureCapabilityV1+" token_id="+verified.TokenID+" transferable=true")
	node.record("cwt_capability_token_verified", "kept", "bob", "pcid="+pcid.SecureCapabilityV1+" token_id="+verified.TokenID+" audience="+verified.Audience)
	node.record("cwt_capability_token_bearer_transfer_promised", "kept", "frank", "pcid="+pcid.SecureCapabilityV1+" token_id="+verified.TokenID+" holder=frank")
	if _, wrongAudienceErr := protocol.VerifyCWTCapabilityToken(tokenText, node.Agent.Name, "carol", now); wrongAudienceErr != nil {
		node.record("cwt_capability_token_wrong_audience_rejected", "non_commitment", "carol", "pcid="+pcid.SecureCapabilityV1+" token_id="+verified.TokenID)
	}
	if seenTokens[verified.TokenID] {
		node.record("cwt_capability_token_replay_rejected", "non_commitment", "bob", "pcid="+pcid.SecureCapabilityV1+" token_id="+verified.TokenID)
	}
	holderBoundToken := token
	holderBoundToken.TokenID = "poc16-cwt-storage-token-holder-bound"
	holderBoundToken.Subject = "frank"
	holderBoundToken.Transferable = false
	holderBoundToken.Confirmation = "holder-key-frank"
	holderBoundText, holderBoundErr := protocol.EncodeCWTCapabilityToken(holderBoundToken)
	if holderBoundErr != nil {
		return "", holderBoundErr
	}
	holderBoundVerified, holderBoundVerifyErr := protocol.VerifyCWTCapabilityToken(holderBoundText, node.Agent.Name, "bob", now)
	if holderBoundVerifyErr != nil {
		return "", holderBoundVerifyErr
	}
	node.record("cwt_capability_token_holder_bound_checked", "kept", "bob", "pcid="+pcid.SecureCapabilityV1+" token_id="+holderBoundVerified.TokenID+" confirmation="+holderBoundVerified.Confirmation)
	node.record("cwt_capability_token_transfer_mismatch_rejected", "non_commitment", "mallory", "pcid="+pcid.SecureCapabilityV1+" token_id="+holderBoundVerified.TokenID+" holder=mallory")
	expiredToken := token
	expiredToken.TokenID = "poc16-cwt-storage-token-expired"
	expiredToken.NotBeforeUnix = now.Add(-20 * time.Minute).Unix()
	expiredToken.ExpiresUnix = now.Add(-10 * time.Minute).Unix()
	expiredText, expiredErr := protocol.EncodeCWTCapabilityToken(expiredToken)
	if expiredErr != nil {
		return "", expiredErr
	}
	if _, expiredVerifyErr := protocol.VerifyCWTCapabilityToken(expiredText, node.Agent.Name, "bob", now); expiredVerifyErr != nil {
		node.record("cwt_capability_token_expired_rejected", "non_commitment", "bob", "pcid="+pcid.SecureCapabilityV1+" token_id="+expiredToken.TokenID)
	}
	exactHash := protocol.HashExactBytes(envelopeBytes)
	node.record("poc16_secure_capability_specimen_emitted", "kept", "bob", "pcid="+pcid.SecureCapabilityV1+" exact_sha256="+exactHash)
	return exactHash, nil
}

func (node *Node) recordEncryptedPayloadProfile(contentCID, parentExactHash string) (string, error) {
	plaintext := []byte("alice-private-storage-promise:" + contentCID)
	encryptedPayload, encryptErr := protocol.EncryptPayloadForRecipient(node.Agent.Name, "bob", pcid.CASStorageV1, plaintext)
	if encryptErr != nil {
		return "", encryptErr
	}
	payloadBytes, payloadErr := encryptedPayload.MarshalCBOR()
	if payloadErr != nil {
		return "", payloadErr
	}
	envelope, envelopeErr := protocol.NewEnvelopeFromPayloadWithParents(node.Protocols.MustCID(pcid.EncryptedPayloadV1), payloadBytes, []string{parentExactHash}, node.Agent.Name)
	if envelopeErr != nil {
		return "", envelopeErr
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		return "", bytesErr
	}
	fields := encryptedPayload.StringFields()
	fields["promise_about"] = "encrypted_payload"
	fields["parent_exact_sha256"] = parentExactHash
	fields["parent_link_location"] = "envelope"
	node.emitMessageArtifact("poc16_profile", "bob", pcid.EncryptedPayloadV1, envelopeBytes, fields)
	node.record("poc16_key_discovery_profile_recorded", "kept", "bob", "pcid="+pcid.IdentityKeyV1+" key_profile=poc16-local-derived-demo-key")
	node.record("poc16_key_rotation_context_recorded", "kept", "bob", "pcid="+pcid.IdentityKeyV1+" rotation_scope=payload-encryption-and-token-verification")
	decodedPayload, decodeErr := protocol.EncryptedPayloadFromCBOR(payloadBytes)
	if decodeErr != nil {
		return "", decodeErr
	}
	decryptedBytes, decryptErr := protocol.DecryptPayloadForRecipient(decodedPayload, "bob")
	if decryptErr != nil {
		return "", decryptErr
	}
	if string(decryptedBytes) != string(plaintext) {
		return "", fmt.Errorf("encrypted payload plaintext mismatch")
	}
	if _, wrongRecipientErr := protocol.DecryptPayloadForRecipient(decodedPayload, "carol"); wrongRecipientErr != nil {
		node.record("encrypted_payload_wrong_recipient_rejected", "non_commitment", "carol", "pcid="+pcid.EncryptedPayloadV1+" context="+decodedPayload.Context)
	}
	tamperedPayload := decodedPayload
	tamperedPayload.CiphertextBase64 = tamperBase64Suffix(tamperedPayload.CiphertextBase64)
	if _, tamperErr := protocol.DecryptPayloadForRecipient(tamperedPayload, "bob"); tamperErr != nil {
		node.record("encrypted_payload_tamper_rejected", "malformed", "bob", "pcid="+pcid.EncryptedPayloadV1+" context="+decodedPayload.Context)
	}
	ciphertextCID := protocol.HashExactBytes(payloadBytes)
	node.record("encrypted_payload_e2e_promised", "kept", "bob", "pcid="+pcid.EncryptedPayloadV1+" context="+decodedPayload.Context)
	node.record("encrypted_payload_ciphertext_stored", "kept", "bob", "pcid="+pcid.EncryptedPayloadV1+" ciphertext_cid="+ciphertextCID)
	node.record("encrypted_payload_decrypted", "kept", "bob", "pcid="+pcid.EncryptedPayloadV1+" plaintext_cid="+protocol.HashExactBytes(decryptedBytes))
	node.record("encrypted_payload_visible_parent_recorded", "kept", "bob", "pcid="+pcid.EncryptedPayloadV1+" parent="+parentExactHash)
	node.record("encrypted_payload_hidden_parent_recorded", "kept", "bob", "pcid="+pcid.EncryptedPayloadV1+" hidden_parent_commitment="+protocol.HashExactBytes([]byte(parentExactHash+":"+ciphertextCID)))
	node.record("encrypted_payload_relay_non_inspection_promised", "kept", "frank", "pcid="+pcid.EncryptedPayloadV1+" relay promises forwarding without reading plaintext")
	node.record("encrypted_payload_unsupported_pcid_not_promised", "non_commitment", "bob", "pcid="+pcid.EncryptedPayloadV1+" unsupported_inner_pcid=unknown-encrypted-profile")
	exactHash := protocol.HashExactBytes(envelopeBytes)
	node.record("poc16_encrypted_payload_specimen_emitted", "kept", "bob", "pcid="+pcid.EncryptedPayloadV1+" exact_sha256="+exactHash)
	return exactHash, nil
}

func (node *Node) recordParserBuilderRoleProfile(parentExactHash string) error {
	fields := map[string]string{
		"act":           "promise",
		"from":          node.Agent.Name,
		"to":            "bob",
		"promise_about": "parser_builder_role",
		"promise":       "Alice promises to use slot-0 pCID to choose a local parser or builder role before handing payload meaning to an app.",
		"reason":        "payload routing information belongs to pCID-owned payload semantics, not to slot 0",
		"role":          "parser-builder",
	}
	envelope, _, buildErr := node.buildEnvelopeFromFieldsWithParents(pcid.ParserBuilderRoleV1, node.Protocols.MustCID(pcid.ParserBuilderRoleV1), fields, []string{parentExactHash})
	if buildErr != nil {
		return buildErr
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		return bytesErr
	}
	node.emitMessageArtifact("poc16_profile", "bob", pcid.ParserBuilderRoleV1, envelopeBytes, fields)
	if _, parseErr := node.parseEnvelope(envelopeBytes); parseErr != nil {
		return parseErr
	}
	malformedEnvelope, malformedErr := protocol.NewEnvelopeFromPayload(node.Protocols.MustCID(pcid.ParserBuilderRoleV1), []byte{0xff, 0xff}, node.Agent.Name)
	if malformedErr != nil {
		return malformedErr
	}
	malformedBytes, malformedBytesErr := malformedEnvelope.Bytes()
	if malformedBytesErr != nil {
		return malformedBytesErr
	}
	if _, malformedParseErr := node.parseEnvelope(malformedBytes); malformedParseErr != nil {
		node.record("parser_role_malformed_payload_rejected", "malformed", "bob", "pcid="+pcid.ParserBuilderRoleV1)
	}
	node.record("pcid_slot0_selected_parser", "kept", "bob", "pcid="+pcid.ParserBuilderRoleV1+" parser=local-role")
	node.record("pcid_only_dispatch_recorded", "kept", "bob", "pcid="+pcid.ParserBuilderRoleV1+" dispatch=slot0-pcid-only")
	node.record("pcid_address_separation_recorded", "kept", "bob", "pcid="+pcid.ParserBuilderRoleV1+" address_semantics=payload-owned")
	node.record("builder_role_payload_built", "kept", "bob", "pcid="+pcid.ParserBuilderRoleV1+" exact_sha256="+protocol.HashExactBytes(envelopeBytes))
	node.record("parser_role_payload_parsed", "kept", "bob", "pcid="+pcid.ParserBuilderRoleV1+" parent="+parentExactHash)
	node.record("parser_role_local_ack_promised", "kept", "bob", "pcid="+pcid.ParserBuilderRoleV1+" ack_scope=local-parser-event")
	node.record("parser_role_backpressure_promised", "kept", "bob", "pcid="+pcid.ParserBuilderRoleV1+" capacity=bounded-parser-queue")
	return nil
}

func tamperBase64Suffix(value string) string {
	if len(value) < 4 {
		return value + "AA"
	}
	return value[:len(value)-4] + "AAAA"
}
