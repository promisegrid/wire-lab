package token

import "testing"

func TestNonTransferableTokenCannotMove(t *testing.T) {
	alice := NewIssuer("alice")
	issued, issueErr := alice.Issue("alice-data-bob-1", "bob", "data", "dataset-a", TransferNonTransferable)
	if issueErr != nil {
		t.Fatalf("issue: %v", issueErr)
	}
	event := alice.Redeem("carol", issued)
	if event.Outcome != OutcomeRefused {
		t.Fatalf("non-transferable token redeemed by wrong holder: %#v", event)
	}
}

func TestBearerTokenTransferAndRedeem(t *testing.T) {
	alice := NewIssuer("alice")
	bob := NewWallet("bob")
	carol := NewWallet("carol")
	issued, issueErr := alice.Issue("alice-bearer-1", "bob", "data", "dataset-b", TransferBearer)
	if issueErr != nil {
		t.Fatalf("issue: %v", issueErr)
	}
	bob.Add(issued, "alice issued bearer token to bob")
	if _, transferErr := bob.Transfer(issued.ID, carol); transferErr != nil {
		t.Fatalf("transfer: %v", transferErr)
	}
	event := alice.Redeem("carol", issued)
	if event.Outcome != OutcomeKept {
		t.Fatalf("bearer token was not redeemed by new holder: %#v", event)
	}
}

func TestRevokedTokenRecordsBrokenPromiseEvidence(t *testing.T) {
	alice := NewIssuer("alice")
	mallory := NewWallet("mallory")
	issued, issueErr := alice.Issue("alice-bearer-revoked", "mallory", "data", "dataset-c", TransferBearer)
	if issueErr != nil {
		t.Fatalf("issue: %v", issueErr)
	}
	mallory.Add(issued, "mallory received bearer token")
	if revokeErr := alice.Revoke(issued.ID, "issuer-local revocation before redemption"); revokeErr != nil {
		t.Fatalf("revoke: %v", revokeErr)
	}
	event := alice.Redeem("mallory", issued)
	if event.Outcome != OutcomeBroken {
		t.Fatalf("revoked token should create broken-promise evidence: %#v", event)
	}
	mallory.ApplyRedemption(event)
	if mallory.Trust("alice") >= 0 {
		t.Fatalf("mallory trust in alice should decrease after broken redemption")
	}
}

func TestExpiredTokenDoesNotLowerIssuerTrust(t *testing.T) {
	alice := NewIssuer("alice")
	alice.SetNowUnix(200)
	dave := NewWallet("dave")
	issued, issueErr := alice.IssueExpiring("alice-bearer-expired", "mallory", "data", "dataset-expired", TransferBearer, 100)
	if issueErr != nil {
		t.Fatalf("issue expiring: %v", issueErr)
	}
	dave.Add(issued, "dave accepted expired-token bytes for evidence")
	event := alice.Redeem("dave", issued)
	if event.Outcome != OutcomeExpired {
		t.Fatalf("expired token should be neutral issuer-expiry evidence: %#v", event)
	}
	dave.ApplyRedemption(event)
	if dave.Trust("alice") != 0 {
		t.Fatalf("dave trust in alice should stay neutral after signed expiry, got %d", dave.Trust("alice"))
	}
	dave.ApplyPeerObservation("mallory", event)
	if dave.Trust("mallory") >= 0 {
		t.Fatalf("dave trust in mallory should decrease after mallory circulates expired-token evidence")
	}
}

func TestPeerObservationUpdatesCirculatorTrust(t *testing.T) {
	dave := NewWallet("dave")
	dave.ApplyPeerObservation("mallory", Event{Observer: "dave", Event: "held_redemption_observed", Outcome: OutcomeBroken, TokenID: "alice-revoked"})
	if dave.Trust("mallory") >= 0 {
		t.Fatalf("dave trust in mallory should decrease after mallory circulates broken token evidence")
	}
}

func TestTokenSignatureDetectsMutation(t *testing.T) {
	alice := NewIssuer("alice")
	issued, issueErr := alice.Issue("alice-bearer-proof", "bob", "data", "dataset-d", TransferBearer)
	if issueErr != nil {
		t.Fatalf("issue: %v", issueErr)
	}
	issued.ResourceID = "rewritten"
	if verifyErr := VerifyToken(issued); verifyErr == nil {
		t.Fatalf("mutated token proof unexpectedly verified")
	}
}

func TestTokenCBORRoundTrip(t *testing.T) {
	alice := NewIssuer("alice")
	issued, issueErr := alice.Issue("alice-cbor-token", "bob", "data", "dataset-e", TransferBearer)
	if issueErr != nil {
		t.Fatalf("issue: %v", issueErr)
	}
	tokenBytes, encodeErr := Encode(issued)
	if encodeErr != nil {
		t.Fatalf("encode: %v", encodeErr)
	}
	decoded, decodeErr := Decode(tokenBytes)
	if decodeErr != nil {
		t.Fatalf("decode: %v", decodeErr)
	}
	if decoded.ID != issued.ID || decoded.PublicKeyHex == "" || decoded.SignatureHex == "" {
		t.Fatalf("decoded token lost proof fields: %#v", decoded)
	}
	if verifyErr := VerifyToken(decoded); verifyErr != nil {
		t.Fatalf("verify decoded: %v", verifyErr)
	}
}
