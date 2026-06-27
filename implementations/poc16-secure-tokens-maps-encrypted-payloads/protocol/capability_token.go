package protocol

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/fxamacker/cbor/v2"
	cose "github.com/veraison/go-cose"
)

const signedCapabilityTokenCapability = "cas-storage-token"

// SignedCapabilityToken is POC16's signed-CBOR pressure model for capability
// promise tokens.
// Intent: Tokens are issuer promises encoded as signed bytes, not authority
// objects. CAS/storage tokens now use the same CWT-style numeric term map and
// well-known COSE_Sign1/Ed25519 library pattern as local lifecycle tokens, while
// preserving the issuer-local serve-once and bearer-storage semantics already
// exercised by the storage workflows. Source: DI-lurov
type SignedCapabilityToken struct {
	Issuer       string
	Subject      string
	Scope        string
	ContentCID   string
	ExpiresUnix  int64
	Nonce        string
	Transferable bool
}

// EncodeSignedCapabilityToken encodes token terms as canonical CWT-style CBOR
// and signs them with the issuer's deterministic POC16 Ed25519 key through the
// same COSE library used by local_lifecycle_v1.
func EncodeSignedCapabilityToken(token SignedCapabilityToken) (string, error) {
	if err := validateSignedCapabilityToken(token); err != nil {
		return "", err
	}
	payloadBytes, payloadErr := signedCapabilityTokenClaims(token)
	if payloadErr != nil {
		return "", payloadErr
	}
	signer, signerErr := cose.NewSigner(cose.AlgorithmEdDSA, DeterministicPrivateKey(token.Issuer))
	if signerErr != nil {
		return "", signerErr
	}
	headers := cose.Headers{
		Protected: cose.ProtectedHeader{
			cose.HeaderLabelAlgorithm: cose.AlgorithmEdDSA,
		},
		Unprotected: cose.UnprotectedHeader{},
	}
	coseBytes, coseErr := cose.Sign1(rand.Reader, signer, headers, payloadBytes, nil)
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
	var message cose.Sign1Message
	if unmarshalErr := message.UnmarshalCBOR(coseBytes); unmarshalErr != nil {
		return SignedCapabilityToken{}, unmarshalErr
	}
	verifier, verifierErr := cose.NewVerifier(cose.AlgorithmEdDSA, DeterministicPublicKey(expectedIssuer))
	if verifierErr != nil {
		return SignedCapabilityToken{}, verifierErr
	}
	if verifyErr := message.Verify(nil, verifier); verifyErr != nil {
		return SignedCapabilityToken{}, verifyErr
	}
	token, tokenErr := signedCapabilityTokenFromClaims(message.Payload)
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

// validateSignedCapabilityToken keeps the local promise-token terms explicit so
// malformed issuer promises fail before signing or after decoding.
func validateSignedCapabilityToken(token SignedCapabilityToken) error {
	if token.Issuer == "" || token.Subject == "" || token.Scope == "" || token.ContentCID == "" || token.Nonce == "" {
		return fmt.Errorf("signed capability token needs issuer, subject, scope, content CID, and nonce")
	}
	if token.ExpiresUnix == 0 {
		return fmt.Errorf("signed capability token needs exp")
	}
	return nil
}

// signedCapabilityTokenClaims is the CAS/storage token's CWT-style term map.
// Intent: Reuse the same numeric CWT labels as the broader capability-token
// specimen wherever they fit, while keeping this runtime token's promise terms
// narrow: issuer, subject, scope, content CID, expiry, nonce, and transferability.
// Source: DI-lurov
func signedCapabilityTokenClaims(token SignedCapabilityToken) ([]byte, error) {
	terms := map[any]any{
		cose.CWTClaimIssuer:         token.Issuer,
		cose.CWTClaimSubject:        token.Subject,
		cose.CWTClaimExpirationTime: token.ExpiresUnix,
		cose.CWTClaimCWTID:          []byte(token.Nonce),
		cwtClaimCapability:          signedCapabilityTokenCapability,
		cwtClaimScope:               token.Scope,
		cwtClaimContentCID:          token.ContentCID,
		cwtClaimTransferable:        token.Transferable,
	}
	return lifecycleCBOREncMode.Marshal(terms)
}

func signedCapabilityTokenFromClaims(payloadBytes []byte) (SignedCapabilityToken, error) {
	var terms map[any]any
	if err := cbor.Unmarshal(payloadBytes, &terms); err != nil {
		return SignedCapabilityToken{}, err
	}
	nonce, nonceErr := bytesTextTerm(cwtValue(terms, cose.CWTClaimCWTID))
	if nonceErr != nil {
		return SignedCapabilityToken{}, nonceErr
	}
	transferable, transferableErr := boolTerm(cwtValue(terms, cwtClaimTransferable))
	if transferableErr != nil {
		return SignedCapabilityToken{}, transferableErr
	}
	if stringTerm(cwtValue(terms, cwtClaimCapability)) != signedCapabilityTokenCapability {
		return SignedCapabilityToken{}, fmt.Errorf("signed capability token capability mismatch")
	}
	token := SignedCapabilityToken{
		Issuer:       stringTerm(cwtValue(terms, cose.CWTClaimIssuer)),
		Subject:      stringTerm(cwtValue(terms, cose.CWTClaimSubject)),
		Scope:        stringTerm(cwtValue(terms, cwtClaimScope)),
		ContentCID:   stringTerm(cwtValue(terms, cwtClaimContentCID)),
		ExpiresUnix:  int64Term(cwtValue(terms, cose.CWTClaimExpirationTime)),
		Nonce:        nonce,
		Transferable: transferable,
	}
	if err := validateSignedCapabilityToken(token); err != nil {
		return SignedCapabilityToken{}, err
	}
	return token, nil
}

func boolTerm(value any) (bool, error) {
	boolValue, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("expected bool")
	}
	return boolValue, nil
}
