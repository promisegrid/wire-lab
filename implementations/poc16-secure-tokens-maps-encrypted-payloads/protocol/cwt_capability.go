package protocol

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

const (
	cwtClaimIssuer       = int64(1)
	cwtClaimSubject      = int64(2)
	cwtClaimAudience     = int64(3)
	cwtClaimExpiration   = int64(4)
	cwtClaimNotBefore    = int64(5)
	cwtClaimTokenID      = int64(7)
	cwtClaimCapability   = int64(-70000)
	cwtClaimScope        = int64(-70001)
	cwtClaimContentCID   = int64(-70002)
	cwtClaimTransferable = int64(-70003)
	cwtClaimConfirmation = int64(-70004)
)

// CWTCapabilityToken is POC16's compact token profile for issuer promises that
// can be redeemed later by a holder.
//
// Intent: Capability tokens are promises by the issuer, not grants from an
// authority. The token bytes are cryptographically signed COSE_Sign1 over a
// narrow CWT-style CBOR claim map so expiry, audience, transferability, and
// replay identifiers are concrete in POC16. Source: DI-vulit
type CWTCapabilityToken struct {
	Issuer        string
	Subject       string
	Audience      string
	Capability    string
	Scope         string
	ContentCID    string
	TokenID       string
	Confirmation  string
	ExpiresUnix   int64
	NotBeforeUnix int64
	Transferable  bool
}

// EncodeCWTCapabilityToken signs a CWT-style claim map with the issuer's
// deterministic POC16 Ed25519 key and returns base64 COSE_Sign1 bytes.
func EncodeCWTCapabilityToken(token CWTCapabilityToken) (string, error) {
	if err := validateCWTCapabilityToken(token); err != nil {
		return "", err
	}
	claimsBytes, claimsErr := marshalCWTCapabilityClaims(token)
	if claimsErr != nil {
		return "", claimsErr
	}
	coseBytes, coseErr := EncodeCOSESign1(claimsBytes, token.Issuer, false)
	if coseErr != nil {
		return "", coseErr
	}
	return base64.StdEncoding.EncodeToString(coseBytes), nil
}

// VerifyCWTCapabilityToken verifies COSE bytes, checks the issuer/audience/time
// promises, and returns the signed claims for local redemption logic.
func VerifyCWTCapabilityToken(tokenText, expectedIssuer, expectedAudience string, now time.Time) (CWTCapabilityToken, error) {
	coseBytes, decodeErr := base64.StdEncoding.DecodeString(tokenText)
	if decodeErr != nil {
		return CWTCapabilityToken{}, decodeErr
	}
	if verifyErr := VerifyCOSESign1(coseBytes, nil, expectedIssuer); verifyErr != nil {
		return CWTCapabilityToken{}, verifyErr
	}
	coseSign1, parseErr := parseCOSESign1(coseBytes)
	if parseErr != nil {
		return CWTCapabilityToken{}, parseErr
	}
	token, tokenErr := unmarshalCWTCapabilityClaims(coseSign1.Payload)
	if tokenErr != nil {
		return CWTCapabilityToken{}, tokenErr
	}
	if token.Issuer != expectedIssuer {
		return CWTCapabilityToken{}, fmt.Errorf("cwt capability issuer=%s, want %s", token.Issuer, expectedIssuer)
	}
	if token.Audience != expectedAudience {
		return CWTCapabilityToken{}, fmt.Errorf("cwt capability audience=%s, want %s", token.Audience, expectedAudience)
	}
	if now.Unix() < token.NotBeforeUnix {
		return CWTCapabilityToken{}, fmt.Errorf("cwt capability not valid before %d", token.NotBeforeUnix)
	}
	if now.Unix() >= token.ExpiresUnix {
		return CWTCapabilityToken{}, fmt.Errorf("cwt capability expired at %d", token.ExpiresUnix)
	}
	return token, nil
}

func validateCWTCapabilityToken(token CWTCapabilityToken) error {
	if token.Issuer == "" || token.Subject == "" || token.Audience == "" || token.Capability == "" || token.Scope == "" || token.ContentCID == "" || token.TokenID == "" {
		return fmt.Errorf("cwt capability token needs issuer, subject, audience, capability, scope, content CID, and token ID")
	}
	if token.NotBeforeUnix == 0 || token.ExpiresUnix == 0 {
		return fmt.Errorf("cwt capability token needs nbf and exp")
	}
	if token.ExpiresUnix <= token.NotBeforeUnix {
		return fmt.Errorf("cwt capability token exp must be after nbf")
	}
	if !token.Transferable && token.Confirmation == "" {
		return fmt.Errorf("non-transferable cwt capability token needs holder confirmation")
	}
	return nil
}

func marshalCWTCapabilityClaims(token CWTCapabilityToken) ([]byte, error) {
	writer := &cborWriter{}
	if err := writer.writeMapHeader(11); err != nil {
		return nil, err
	}
	stringClaims := []struct {
		label int64
		value string
	}{
		{label: cwtClaimIssuer, value: token.Issuer},
		{label: cwtClaimSubject, value: token.Subject},
		{label: cwtClaimAudience, value: token.Audience},
		{label: cwtClaimCapability, value: token.Capability},
		{label: cwtClaimScope, value: token.Scope},
		{label: cwtClaimContentCID, value: token.ContentCID},
		{label: cwtClaimConfirmation, value: token.Confirmation},
	}
	for _, claim := range stringClaims {
		if err := writer.writeSignedInt(claim.label); err != nil {
			return nil, err
		}
		if err := writer.writeString(claim.value); err != nil {
			return nil, err
		}
	}
	for _, claim := range []struct {
		label int64
		value int64
	}{
		{label: cwtClaimExpiration, value: token.ExpiresUnix},
		{label: cwtClaimNotBefore, value: token.NotBeforeUnix},
	} {
		if err := writer.writeSignedInt(claim.label); err != nil {
			return nil, err
		}
		if err := writer.writeSignedInt(claim.value); err != nil {
			return nil, err
		}
	}
	if err := writer.writeSignedInt(cwtClaimTokenID); err != nil {
		return nil, err
	}
	if err := writer.writeBytes([]byte(token.TokenID)); err != nil {
		return nil, err
	}
	if err := writer.writeSignedInt(cwtClaimTransferable); err != nil {
		return nil, err
	}
	if err := writer.writeBool(token.Transferable); err != nil {
		return nil, err
	}
	return writer.buffer.Bytes(), nil
}

func unmarshalCWTCapabilityClaims(claimsBytes []byte) (CWTCapabilityToken, error) {
	reader := &cborReader{data: claimsBytes}
	claimCount, claimErr := reader.readTypeAndLength(5)
	if claimErr != nil {
		return CWTCapabilityToken{}, claimErr
	}
	token := CWTCapabilityToken{}
	for index := uint64(0); index < claimCount; index++ {
		label, labelErr := reader.readSignedInt()
		if labelErr != nil {
			return CWTCapabilityToken{}, labelErr
		}
		switch label {
		case cwtClaimIssuer:
			value, err := reader.readString()
			if err != nil {
				return CWTCapabilityToken{}, err
			}
			token.Issuer = value
		case cwtClaimSubject:
			value, err := reader.readString()
			if err != nil {
				return CWTCapabilityToken{}, err
			}
			token.Subject = value
		case cwtClaimAudience:
			value, err := reader.readString()
			if err != nil {
				return CWTCapabilityToken{}, err
			}
			token.Audience = value
		case cwtClaimExpiration:
			value, err := reader.readSignedInt()
			if err != nil {
				return CWTCapabilityToken{}, err
			}
			token.ExpiresUnix = value
		case cwtClaimNotBefore:
			value, err := reader.readSignedInt()
			if err != nil {
				return CWTCapabilityToken{}, err
			}
			token.NotBeforeUnix = value
		case cwtClaimTokenID:
			value, err := reader.readBytes()
			if err != nil {
				return CWTCapabilityToken{}, err
			}
			token.TokenID = string(value)
		case cwtClaimCapability:
			value, err := reader.readString()
			if err != nil {
				return CWTCapabilityToken{}, err
			}
			token.Capability = value
		case cwtClaimScope:
			value, err := reader.readString()
			if err != nil {
				return CWTCapabilityToken{}, err
			}
			token.Scope = value
		case cwtClaimContentCID:
			value, err := reader.readString()
			if err != nil {
				return CWTCapabilityToken{}, err
			}
			token.ContentCID = value
		case cwtClaimTransferable:
			value, err := reader.readBool()
			if err != nil {
				return CWTCapabilityToken{}, err
			}
			token.Transferable = value
		case cwtClaimConfirmation:
			value, err := reader.readString()
			if err != nil {
				return CWTCapabilityToken{}, err
			}
			token.Confirmation = value
		default:
			if err := reader.skipItem(); err != nil {
				return CWTCapabilityToken{}, err
			}
		}
	}
	if reader.offset != len(reader.data) {
		return CWTCapabilityToken{}, fmt.Errorf("trailing cbor bytes in cwt capability claims: %d", len(reader.data)-reader.offset)
	}
	if err := validateDecodedCWTCapabilityToken(token); err != nil {
		return CWTCapabilityToken{}, err
	}
	return token, nil
}

func validateDecodedCWTCapabilityToken(token CWTCapabilityToken) error {
	if token.Issuer == "" || token.Subject == "" || token.Audience == "" || token.Capability == "" || token.Scope == "" || token.ContentCID == "" || token.TokenID == "" {
		return fmt.Errorf("decoded cwt capability token is missing required claims")
	}
	if token.ExpiresUnix == 0 || token.NotBeforeUnix == 0 {
		return fmt.Errorf("decoded cwt capability token is missing exp or nbf")
	}
	if !token.Transferable && token.Confirmation == "" {
		return fmt.Errorf("decoded non-transferable cwt capability token is missing confirmation")
	}
	return nil
}

// CWTTokenSummaryFields renders non-secret token facts for event records and
// analyzer-friendly payload specimens.
func CWTTokenSummaryFields(token CWTCapabilityToken) map[string]string {
	return map[string]string{
		"audience":     token.Audience,
		"capability":   token.Capability,
		"content_cid":  token.ContentCID,
		"expires_unix": strconv.FormatInt(token.ExpiresUnix, 10),
		"issuer":       token.Issuer,
		"not_before":   strconv.FormatInt(token.NotBeforeUnix, 10),
		"scope":        token.Scope,
		"subject":      token.Subject,
		"token_id":     token.TokenID,
		"transferable": strconv.FormatBool(token.Transferable),
		"confirmation": token.Confirmation,
	}
}
