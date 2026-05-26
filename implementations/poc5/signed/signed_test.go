package signed

import "testing"

func TestSignedResultVerifiesAsLocalEvidence(t *testing.T) {
	request, err := NewSignRequest("Alice", "alice-hello-app", "dave-signed-app", "hello")
	if err != nil {
		t.Fatalf("NewSignRequest returned error: %v", err)
	}
	_, fields, err := envelopeFields(request)
	if err != nil {
		t.Fatalf("envelopeFields returned error: %v", err)
	}
	result, err := SignedApp{NodeName: "Dave", AppName: "dave-signed-app"}.resultEnvelope(fields, "request-hash")
	if err != nil {
		t.Fatalf("resultEnvelope returned error: %v", err)
	}
	resultFields, err := VerifyEnvelope(result)
	if err != nil {
		t.Fatalf("VerifyEnvelope returned error: %v", err)
	}
	if resultFields["request_hash"] != "request-hash" {
		t.Fatalf("request_hash = %q", resultFields["request_hash"])
	}
	if resultFields["text"] != "hello" {
		t.Fatalf("text = %q", resultFields["text"])
	}
}

func envelopeFields(envelope interface {
	PayloadFields() (map[string]string, error)
}) (string, map[string]string, error) {
	fields, err := envelope.PayloadFields()
	if err != nil {
		return "", nil, err
	}
	return fields["kind"], fields, nil
}
