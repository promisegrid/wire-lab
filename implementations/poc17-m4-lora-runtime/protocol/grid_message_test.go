package protocol

import "testing"

func TestBuildParseGridMessage(t *testing.T) {
	payload, err := BuildStatusPayload("m4-ivan", "ready", 87, []string{"missing-parent"})
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
	if device != "m4-ivan" || status != "ready" || battery != 87 || len(parents) != 1 || parents[0] != "missing-parent" {
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
