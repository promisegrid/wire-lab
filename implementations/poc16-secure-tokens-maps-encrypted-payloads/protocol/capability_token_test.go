package protocol

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/fxamacker/cbor/v2"
	cose "github.com/veraison/go-cose"
)

func TestSignedCapabilityTokenRedeemsExpectedClaims(t *testing.T) {
	// Intent: POC16 capability tokens should be real signed CBOR bytes whose
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

func TestSignedCapabilityTokenUsesCWTStyleCOSEPayload(t *testing.T) {
	// Intent: CAS/storage tokens should now match the local_lifecycle_v1
	// implementation pattern: library-verified COSE_Sign1 over a numeric
	// CWT-style CBOR term map, not the older custom string-map token body.
	// Source: DI-lurov
	token := SignedCapabilityToken{
		Issuer:       "bob",
		Subject:      "frank",
		Scope:        "bearer-storage",
		ContentCID:   "bafkrei-test-content",
		ExpiresUnix:  4_102_444_800,
		Nonce:        "token-000003",
		Transferable: true,
	}
	tokenText, encodeErr := EncodeSignedCapabilityToken(token)
	if encodeErr != nil {
		t.Fatalf("encode token: %v", encodeErr)
	}
	coseBytes, decodeErr := base64.StdEncoding.DecodeString(tokenText)
	if decodeErr != nil {
		t.Fatalf("decode token: %v", decodeErr)
	}
	var message cose.Sign1Message
	if unmarshalErr := message.UnmarshalCBOR(coseBytes); unmarshalErr != nil {
		t.Fatalf("unmarshal COSE: %v", unmarshalErr)
	}
	verifier, verifierErr := cose.NewVerifier(cose.AlgorithmEdDSA, DeterministicPublicKey("bob"))
	if verifierErr != nil {
		t.Fatalf("new verifier: %v", verifierErr)
	}
	if verifyErr := message.Verify(nil, verifier); verifyErr != nil {
		t.Fatalf("verify COSE: %v", verifyErr)
	}
	var terms map[any]any
	if unmarshalErr := cbor.Unmarshal(message.Payload, &terms); unmarshalErr != nil {
		t.Fatalf("unmarshal CWT-style terms: %v", unmarshalErr)
	}
	tokenID, tokenIDErr := bytesTextTerm(cwtValue(terms, cose.CWTClaimCWTID))
	if tokenIDErr != nil {
		t.Fatalf("token id: %v", tokenIDErr)
	}
	transferable, transferableErr := boolTerm(cwtValue(terms, cwtClaimTransferable))
	if transferableErr != nil {
		t.Fatalf("transferable: %v", transferableErr)
	}
	if stringTerm(cwtValue(terms, cose.CWTClaimIssuer)) != token.Issuer ||
		stringTerm(cwtValue(terms, cose.CWTClaimSubject)) != token.Subject ||
		stringTerm(cwtValue(terms, cwtClaimScope)) != token.Scope ||
		stringTerm(cwtValue(terms, cwtClaimContentCID)) != token.ContentCID ||
		int64Term(cwtValue(terms, cose.CWTClaimExpirationTime)) != token.ExpiresUnix ||
		tokenID != token.Nonce ||
		!transferable {
		t.Fatalf("CWT-style terms = %#v", terms)
	}
	if _, hasLegacyStringKey := terms["issuer"]; hasLegacyStringKey {
		t.Fatalf("legacy string-map issuer key should not be present: %#v", terms)
	}
}
