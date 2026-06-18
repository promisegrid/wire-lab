package protocol

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

// SignedCapabilityToken is POC15's signed-CBOR pressure model for capability
// promise tokens.
// Intent: Tokens are issuer promises encoded as signed bytes, not authority
// objects. POC15 uses the existing COSE_Sign1/Ed25519 subset to make replay,
// expiry, scope, and transferability concrete while deferring full CWT adoption
// to a later pCID decision. Source: DI-mapop
type SignedCapabilityToken struct {
	Issuer       string
	Subject      string
	Scope        string
	ContentCID   string
	ExpiresUnix  int64
	Nonce        string
	Transferable bool
}

// EncodeSignedCapabilityToken encodes token claims as deterministic CBOR and
// signs them with the issuer's deterministic Ed25519 key.
func EncodeSignedCapabilityToken(token SignedCapabilityToken) (string, error) {
	if token.Issuer == "" || token.Subject == "" || token.Scope == "" || token.ContentCID == "" || token.Nonce == "" {
		return "", fmt.Errorf("signed capability token needs issuer, subject, scope, content CID, and nonce")
	}
	payloadBytes, payloadErr := MarshalStringMap(signedCapabilityTokenClaims(token))
	if payloadErr != nil {
		return "", payloadErr
	}
	coseBytes, coseErr := EncodeCOSESign1(payloadBytes, token.Issuer, false)
	if coseErr != nil {
		return "", coseErr
	}
	return base64.StdEncoding.EncodeToString(coseBytes), nil
}

// VerifySignedCapabilityToken verifies the issuer signature and returns the
// signed claims without consulting any local redemption table.
func VerifySignedCapabilityToken(tokenText, expectedIssuer string, now time.Time) (SignedCapabilityToken, error) {
	coseBytes, decodeErr := base64.StdEncoding.DecodeString(tokenText)
	if decodeErr != nil {
		return SignedCapabilityToken{}, decodeErr
	}
	if verifyErr := VerifyCOSESign1(coseBytes, nil, expectedIssuer); verifyErr != nil {
		return SignedCapabilityToken{}, verifyErr
	}
	coseSign1, parseErr := parseCOSESign1(coseBytes)
	if parseErr != nil {
		return SignedCapabilityToken{}, parseErr
	}
	fields, fieldsErr := UnmarshalStringMap(coseSign1.Payload)
	if fieldsErr != nil {
		return SignedCapabilityToken{}, fieldsErr
	}
	token, tokenErr := signedCapabilityTokenFromClaims(fields)
	if tokenErr != nil {
		return SignedCapabilityToken{}, tokenErr
	}
	if token.Issuer != expectedIssuer {
		return SignedCapabilityToken{}, fmt.Errorf("signed capability token issuer=%s, want %s", token.Issuer, expectedIssuer)
	}
	if token.ExpiresUnix <= now.Unix() {
		return SignedCapabilityToken{}, fmt.Errorf("signed capability token expired at %d", token.ExpiresUnix)
	}
	return token, nil
}

// RedeemSignedCapabilityToken verifies the token and checks the expected local
// promise terms before the issuer consumes it from local state.
func RedeemSignedCapabilityToken(tokenText, expectedIssuer, expectedSubject, expectedScope, expectedContentCID string, now time.Time) (SignedCapabilityToken, error) {
	token, verifyErr := VerifySignedCapabilityToken(tokenText, expectedIssuer, now)
	if verifyErr != nil {
		return SignedCapabilityToken{}, verifyErr
	}
	if token.Subject != expectedSubject {
		return SignedCapabilityToken{}, fmt.Errorf("signed capability token subject=%s, want %s", token.Subject, expectedSubject)
	}
	if token.Scope != expectedScope {
		return SignedCapabilityToken{}, fmt.Errorf("signed capability token scope=%s, want %s", token.Scope, expectedScope)
	}
	if token.ContentCID != expectedContentCID {
		return SignedCapabilityToken{}, fmt.Errorf("signed capability token content_cid=%s, want %s", token.ContentCID, expectedContentCID)
	}
	return token, nil
}

func signedCapabilityTokenClaims(token SignedCapabilityToken) map[string]string {
	return map[string]string{
		"content_cid":  token.ContentCID,
		"expires_unix": strconv.FormatInt(token.ExpiresUnix, 10),
		"issuer":       token.Issuer,
		"nonce":        token.Nonce,
		"scope":        token.Scope,
		"subject":      token.Subject,
		"transferable": strconv.FormatBool(token.Transferable),
		"type":         "signed_capability_token_v1",
	}
}

func signedCapabilityTokenFromClaims(fields map[string]string) (SignedCapabilityToken, error) {
	expiresUnix, expiresErr := strconv.ParseInt(fields["expires_unix"], 10, 64)
	if expiresErr != nil {
		return SignedCapabilityToken{}, expiresErr
	}
	transferable, transferableErr := strconv.ParseBool(fields["transferable"])
	if transferableErr != nil {
		return SignedCapabilityToken{}, transferableErr
	}
	if fields["type"] != "signed_capability_token_v1" {
		return SignedCapabilityToken{}, fmt.Errorf("signed capability token type=%q", fields["type"])
	}
	return SignedCapabilityToken{
		Issuer:       fields["issuer"],
		Subject:      fields["subject"],
		Scope:        fields["scope"],
		ContentCID:   fields["content_cid"],
		ExpiresUnix:  expiresUnix,
		Nonce:        fields["nonce"],
		Transferable: transferable,
	}, nil
}
