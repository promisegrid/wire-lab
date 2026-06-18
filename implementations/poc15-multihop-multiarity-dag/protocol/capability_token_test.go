package protocol

import (
	"encoding/base64"
	"testing"
	"time"
)

func TestSignedCapabilityTokenRedeemsExpectedClaims(t *testing.T) {
	// Intent: POC15 capability tokens should be real signed CBOR bytes whose
	// issuer, subject, scope, content CID, expiry, nonce, and transferability are
	// checked before local redemption. Source: DI-mapop
	token := SignedCapabilityToken{
		Issuer:       "bob",
		Subject:      "alice",
		Scope:        "serve-once",
		ContentCID:   "cidv1-test-content",
		ExpiresUnix:  4_102_444_800,
		Nonce:        "token-000001",
		Transferable: false,
	}
	tokenText, encodeErr := EncodeSignedCapabilityToken(token)
	if encodeErr != nil {
		t.Fatalf("encode token: %v", encodeErr)
	}
	redeemed, redeemErr := RedeemSignedCapabilityToken(tokenText, "bob", "alice", "serve-once", "cidv1-test-content", time.Unix(1_800_000_000, 0))
	if redeemErr != nil {
		t.Fatalf("redeem token: %v", redeemErr)
	}
	if redeemed.Nonce != token.Nonce || redeemed.Transferable {
		t.Fatalf("redeemed claims = %#v", redeemed)
	}
}

func TestSignedCapabilityTokenRejectsTamperExpiryAndScopeMismatch(t *testing.T) {
	token := SignedCapabilityToken{
		Issuer:       "bob",
		Subject:      "bearer",
		Scope:        "bearer-storage",
		ContentCID:   "cidv1-test-content",
		ExpiresUnix:  4_102_444_800,
		Nonce:        "token-000002",
		Transferable: true,
	}
	tokenText, encodeErr := EncodeSignedCapabilityToken(token)
	if encodeErr != nil {
		t.Fatalf("encode token: %v", encodeErr)
	}
	if _, redeemErr := RedeemSignedCapabilityToken(tokenText, "bob", "bearer", "serve-once", "cidv1-test-content", time.Unix(1_800_000_000, 0)); redeemErr == nil {
		t.Fatalf("scope mismatch should fail")
	}
	if _, redeemErr := RedeemSignedCapabilityToken(tokenText, "bob", "bearer", "bearer-storage", "cidv1-test-content", time.Unix(4_102_444_801, 0)); redeemErr == nil {
		t.Fatalf("expired token should fail")
	}
	tamperedBytes, decodeErr := base64.StdEncoding.DecodeString(tokenText)
	if decodeErr != nil {
		t.Fatalf("decode token: %v", decodeErr)
	}
	tamperedBytes[len(tamperedBytes)-1] ^= 0xff
	tamperedText := base64.StdEncoding.EncodeToString(tamperedBytes)
	if _, redeemErr := RedeemSignedCapabilityToken(tamperedText, "bob", "bearer", "bearer-storage", "cidv1-test-content", time.Unix(1_800_000_000, 0)); redeemErr == nil {
		t.Fatalf("tampered token should fail")
	}
}
