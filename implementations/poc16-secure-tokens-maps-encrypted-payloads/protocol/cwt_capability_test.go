package protocol

import (
	"strings"
	"testing"
	"time"
)

func TestCWTCapabilityTokenVerifiesAudienceAndExpiry(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	token := CWTCapabilityToken{
		Issuer:        "alice",
		Subject:       "frank",
		Audience:      "bob",
		Capability:    "store-for-peer",
		Scope:         "cas_storage_v1",
		ContentCID:    "cid:sample",
		TokenID:       "cwt-token-1",
		Confirmation:  "holder-key-frank",
		ExpiresUnix:   now.Add(time.Hour).Unix(),
		NotBeforeUnix: now.Add(-time.Minute).Unix(),
		Transferable:  false,
	}
	tokenText, encodeErr := EncodeCWTCapabilityToken(token)
	if encodeErr != nil {
		t.Fatalf("encode cwt capability token: %v", encodeErr)
	}
	verified, verifyErr := VerifyCWTCapabilityToken(tokenText, "alice", "bob", now)
	if verifyErr != nil {
		t.Fatalf("verify cwt capability token: %v", verifyErr)
	}
	if verified.TokenID != token.TokenID || verified.Confirmation != token.Confirmation {
		t.Fatalf("verified token = %+v, want token id and confirmation from %+v", verified, token)
	}
	if _, wrongAudienceErr := VerifyCWTCapabilityToken(tokenText, "alice", "carol", now); wrongAudienceErr == nil {
		t.Fatalf("wrong audience should fail")
	}
	if _, expiryErr := VerifyCWTCapabilityToken(tokenText, "alice", "bob", now.Add(2*time.Hour)); expiryErr == nil || !strings.Contains(expiryErr.Error(), "expired") {
		t.Fatalf("expired token error = %v, want expired", expiryErr)
	}
}

func TestCWTCapabilityTokenRequiresConfirmationWhenNonTransferable(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	_, encodeErr := EncodeCWTCapabilityToken(CWTCapabilityToken{
		Issuer:        "alice",
		Subject:       "frank",
		Audience:      "bob",
		Capability:    "store-for-peer",
		Scope:         "cas_storage_v1",
		ContentCID:    "cid:sample",
		TokenID:       "cwt-token-2",
		ExpiresUnix:   now.Add(time.Hour).Unix(),
		NotBeforeUnix: now.Add(-time.Minute).Unix(),
		Transferable:  false,
	})
	if encodeErr == nil || !strings.Contains(encodeErr.Error(), "confirmation") {
		t.Fatalf("non-transferable token encode error = %v, want confirmation requirement", encodeErr)
	}
}
