package protocol

import "testing"

func TestBuildParseGridMessage(t *testing.T) {
	parentCID, err := CIDForBytes([]byte("poc17 missing parent fixture"))
	if err != nil {
		t.Fatal(err)
	}
	payload, err := BuildStatusPayload("m4-ivan", "ready", 87, []string{parentCID})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := Build(Message{ProtocolName: ProtocolDeviceStatus, Payload: payload, Proof: []byte("proof")})
	if err != nil {
		t.Fatal(err)
	}
	msg, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if msg.ProtocolName != ProtocolDeviceStatus || msg.PCID != MustPCIDForName(ProtocolDeviceStatus) || string(msg.Proof) != "proof" {
		t.Fatalf("unexpected message: %+v", msg)
	}
	device, status, battery, parents, err := ParseStatusPayload(msg.Payload)
	if err != nil {
		t.Fatal(err)
	}
	if device != "m4-ivan" || status != "ready" || battery != 87 || len(parents) != 1 || parents[0] != parentCID {
		t.Fatalf("unexpected payload %q %q %d %#v", device, status, battery, parents)
	}
}

func TestBuildUsesCIDBytesInSlotZero(t *testing.T) {
	raw, err := Build(Message{ProtocolName: ProtocolOrderStatus, Payload: []byte("payload")})
	if err != nil {
		t.Fatal(err)
	}
	item, err := Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	slotZero := item.Tag.Value.Array[0]
	if slotZero.Tag == nil || slotZero.Tag.Number != CIDTag || slotZero.Tag.Value.Bytes == nil {
		t.Fatalf("slot 0 must be tag-42 CID bytes: %+v", slotZero)
	}
	if len(slotZero.Tag.Value.Bytes) != 37 || slotZero.Tag.Value.Bytes[0] != 0x00 {
		t.Fatalf("slot 0 must contain DAG-CBOR CID data, got %x", slotZero.Tag.Value.Bytes)
	}
}

func TestParseRejectsBrokenGrid(t *testing.T) {
	if _, err := Parse([]byte{0xff, 0x01}); err == nil {
		t.Fatal("expected malformed CBOR rejection")
	}
}

func TestValidateCIDTextRejectsHexDigest(t *testing.T) {
	if err := ValidateCIDText("aeb977a9a7cf9076107dc7cb3a901d7e0c4e37f6e785a73e6658e702da275068"); err == nil {
		t.Fatal("expected bare hex digest rejection")
	}
	if err := ValidateCIDText(MustPCIDForName(ProtocolOrderStatus)); err != nil {
		t.Fatal(err)
	}
}

func TestBuildParseOrderStatusPayload(t *testing.T) {
	payload, err := BuildOrderStatusPayload(OrderStatusPayload{
		Type:        "MSG",
		Source:      "gateway-bob",
		Dest:        "m4-ivan",
		Counter:     7,
		OrderNumber: "BT-1042",
		Status:      "cut",
	})
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := ParseOrderStatusPayload(payload)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.Type != "MSG" || decoded.Source != "gateway-bob" || decoded.Dest != "m4-ivan" || decoded.Counter != 7 || decoded.OrderNumber != "BT-1042" || decoded.Status != "cut" {
		t.Fatalf("unexpected order payload: %+v", decoded)
	}
}

func TestBuildParsePeerStoragePayloadsUseShape(t *testing.T) {
	content := []byte("parent fixture")
	contentCID, err := CIDForBytes(content)
	if err != nil {
		t.Fatal(err)
	}
	requestCID, err := CIDForBytes([]byte("peer_storage request"))
	if err != nil {
		t.Fatal(err)
	}
	tokenCID, err := CIDForBytes([]byte("compact-token"))
	if err != nil {
		t.Fatal(err)
	}
	token := []byte("compact-token")
	put, err := BuildPeerStoragePut(PeerStoragePayload{
		Holder:     "m4-ivan",
		Issuer:     "gateway-bob",
		Token:      token,
		ContentCID: contentCID,
		Content:    content,
		Reason:     "cas_retention_limit",
	})
	if err != nil {
		t.Fatal(err)
	}
	item, err := Decode(put)
	if err != nil {
		t.Fatal(err)
	}
	if len(item.Array) != 6 {
		t.Fatalf("put payload should not include a redundant kind slot: %#v", item.Array)
	}
	if item.Array[3].Tag == nil || item.Array[3].Tag.Number != CIDTag {
		t.Fatalf("content CID must use tag 42, got %#v", item.Array[3])
	}
	decoded, err := ParsePeerStoragePayload(put)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.Kind != PeerStoragePut || decoded.ContentCID != contentCID || string(decoded.Content) != string(content) {
		t.Fatalf("unexpected peer_storage put payload: %+v", decoded)
	}
	putResult, err := BuildPeerStoragePutResult(PeerStoragePayload{
		Issuer:            "gateway-bob",
		Holder:            "m4-ivan",
		TokenCID:          tokenCID,
		RelatedMessageCID: requestCID,
		ContentCID:        contentCID,
		Accepted:          true,
		Reason:            "retain_until_replaced",
	})
	if err != nil {
		t.Fatal(err)
	}
	decodedPutResult, err := ParsePeerStoragePayload(putResult)
	if err != nil {
		t.Fatal(err)
	}
	if decodedPutResult.Kind != PeerStoragePutResult || decodedPutResult.RelatedMessageCID != requestCID || !decodedPutResult.Accepted {
		t.Fatalf("unexpected peer_storage put result payload: %+v", decodedPutResult)
	}
	getResult, err := BuildPeerStorageGetResult(PeerStoragePayload{
		Issuer:            "gateway-bob",
		Holder:            "m4-ivan",
		TokenCID:          tokenCID,
		RelatedMessageCID: requestCID,
		ContentCID:        contentCID,
		Content:           content,
	})
	if err != nil {
		t.Fatal(err)
	}
	decodedGetResult, err := ParsePeerStoragePayload(getResult)
	if err != nil {
		t.Fatal(err)
	}
	if decodedGetResult.Kind != PeerStorageGetResult || decodedGetResult.RelatedMessageCID != requestCID || string(decodedGetResult.Content) != string(content) {
		t.Fatalf("unexpected peer_storage get result payload: %+v", decodedGetResult)
	}
}
