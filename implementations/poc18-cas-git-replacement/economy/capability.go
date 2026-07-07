package economy

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/fxamacker/cbor/v2"
	cose "github.com/veraison/go-cose"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/graph"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
)

const (
	bearerPaymentCapability  = "poc18-bearer-payment"
	serviceCapabilityPromise = "poc18-service-capability"
	claimServiceCapability   = int64(-90101)
	claimParentBearerCID     = int64(-90102)
)

// BearerToken is a transferable local economic promise. It is not itself a
// permission to use a resource; the accepting resource controller may redeem it
// for a non-transferable capability token.
type BearerToken struct {
	Issuer                 string   `json:"issuer"`
	Scope                  string   `json:"scope"`
	ObjectCID              string   `json:"object_cid"`
	Value                  int64    `json:"value"`
	Unit                   string   `json:"unit"`
	ExpiresUnix            int64    `json:"expires_unix"`
	Nonce                  string   `json:"nonce"`
	Transferable           bool     `json:"transferable"`
	RedeemableCapabilities []string `json:"redeemable_capabilities"`
}

// IssuedBearerToken records exact signed bearer-token bytes and their CID after
// the issuer stores them in its sparse CAS.
type IssuedBearerToken struct {
	Token BearerToken `json:"token"`
	CID   string      `json:"cid"`
	Entry store.Entry `json:"entry"`
	Bytes []byte      `json:"-"`
}

// CapabilityToken is a non-transferable promise by a resource controller to
// provide one specific service capability to one subject.
type CapabilityToken struct {
	Issuer          string `json:"issuer"`
	Subject         string `json:"subject"`
	Capability      string `json:"capability"`
	Scope           string `json:"scope"`
	ObjectCID       string `json:"object_cid"`
	Value           int64  `json:"value"`
	Unit            string `json:"unit"`
	ExpiresUnix     int64  `json:"expires_unix"`
	Nonce           string `json:"nonce"`
	ParentBearerCID string `json:"parent_bearer_cid"`
	Transferable    bool   `json:"transferable"`
}

// ExpectedBearerPayment is the local acceptance policy for redeeming a bearer
// token into a capability token.
type ExpectedBearerPayment struct {
	Issuer     string
	Scope      string
	ObjectCID  string
	Value      int64
	Unit       string
	Capability string
}

// CapabilityRedemptionReport records the fixture-visible result of converting
// a bearer token into a specific non-transferable capability token.
type CapabilityRedemptionReport struct {
	BearerIssuer           string `json:"bearer_issuer"`
	Redeemer               string `json:"redeemer"`
	Subject                string `json:"subject"`
	CapabilityIssuer       string `json:"capability_issuer"`
	Capability             string `json:"capability"`
	Scope                  string `json:"scope"`
	ObjectCID              string `json:"object_cid"`
	Value                  int64  `json:"value"`
	Unit                   string `json:"unit"`
	BearerTokenCID         string `json:"bearer_token_cid"`
	CapabilityTokenCID     string `json:"capability_token_cid"`
	BearerStoredByRedeemer bool   `json:"bearer_stored_by_redeemer"`
	CapabilityStored       bool   `json:"capability_stored"`
	SignatureVerified      bool   `json:"signature_verified"`
	Redeemed               bool   `json:"redeemed"`
	CapabilityBytes        []byte `json:"-"`
}

// IssueBearerToken signs a transferable bearer token and stores the exact token
// bytes in the issuer's CAS.
//
// Intent: Bearer economics should be cryptographic promise objects that can be
// exchanged for specific service capabilities, not informal text credits or a
// global currency. Source: DI-fakop
func IssueBearerToken(cas *store.FileStore, token BearerToken) (IssuedBearerToken, error) {
	if cas == nil {
		return IssuedBearerToken{}, fmt.Errorf("issuer CAS is required")
	}
	if validateErr := validateBearerToken(token); validateErr != nil {
		return IssuedBearerToken{}, validateErr
	}
	payloadBytes, payloadErr := bearerClaims(token)
	if payloadErr != nil {
		return IssuedBearerToken{}, payloadErr
	}
	tokenBytes, signErr := signTokenBytes(token.Issuer, payloadBytes)
	if signErr != nil {
		return IssuedBearerToken{}, signErr
	}
	entry, putErr := cas.Put("payment_token", tokenBytes)
	if putErr != nil {
		return IssuedBearerToken{}, putErr
	}
	return IssuedBearerToken{Token: token, CID: entry.CID, Entry: entry, Bytes: tokenBytes}, nil
}

// RedeemBearerForCapability verifies one bearer token, consumes it in the local
// redeemer ledger, and issues a non-transferable capability token.
func (ledger *Ledger) RedeemBearerForCapability(cas *store.FileStore, bearerBytes []byte, expected ExpectedBearerPayment, capabilityIssuer string, subject string, now time.Time) (CapabilityRedemptionReport, error) {
	if ledger == nil {
		return CapabilityRedemptionReport{}, fmt.Errorf("payment ledger is required")
	}
	if cas == nil {
		return CapabilityRedemptionReport{}, fmt.Errorf("redeemer CAS is required")
	}
	if capabilityIssuer == "" || subject == "" {
		return CapabilityRedemptionReport{}, fmt.Errorf("capability issuer and subject are required")
	}
	bearer, verifyErr := VerifyBearerToken(bearerBytes, expected.Issuer, now)
	if verifyErr != nil {
		return CapabilityRedemptionReport{}, verifyErr
	}
	if matchErr := bearer.matches(expected); matchErr != nil {
		return CapabilityRedemptionReport{}, matchErr
	}
	bearerCID := store.CIDText(store.CIDForBytes(bearerBytes))
	if ledger.spent[bearerCID] {
		return CapabilityRedemptionReport{}, fmt.Errorf("bearer token %s already redeemed", bearerCID)
	}
	bearerEntry, putErr := cas.Put("payment_token", bearerBytes)
	if putErr != nil {
		return CapabilityRedemptionReport{}, putErr
	}
	ledger.spent[bearerCID] = true
	capabilityToken := CapabilityToken{
		Issuer:          capabilityIssuer,
		Subject:         subject,
		Capability:      expected.Capability,
		Scope:           bearer.Scope,
		ObjectCID:       bearer.ObjectCID,
		Value:           bearer.Value,
		Unit:            bearer.Unit,
		ExpiresUnix:     bearer.ExpiresUnix,
		Nonce:           bearer.Nonce + ":" + capabilityIssuer + ":" + subject + ":" + expected.Capability,
		ParentBearerCID: bearerEntry.CID,
		Transferable:    false,
	}
	capabilityBytes, issueErr := issueCapabilityBytes(capabilityToken)
	if issueErr != nil {
		return CapabilityRedemptionReport{}, issueErr
	}
	capabilityEntry, capabilityPutErr := cas.Put("capability_token", capabilityBytes)
	if capabilityPutErr != nil {
		return CapabilityRedemptionReport{}, capabilityPutErr
	}
	return CapabilityRedemptionReport{
		BearerIssuer:           bearer.Issuer,
		Redeemer:               capabilityIssuer,
		Subject:                subject,
		CapabilityIssuer:       capabilityIssuer,
		Capability:             expected.Capability,
		Scope:                  bearer.Scope,
		ObjectCID:              bearer.ObjectCID,
		Value:                  bearer.Value,
		Unit:                   bearer.Unit,
		BearerTokenCID:         bearerEntry.CID,
		CapabilityTokenCID:     capabilityEntry.CID,
		BearerStoredByRedeemer: true,
		CapabilityStored:       true,
		SignatureVerified:      true,
		Redeemed:               true,
		CapabilityBytes:        capabilityBytes,
	}, nil
}

// VerifyBearerToken verifies the COSE signature and returns the signed
// CWT-style bearer token claims.
func VerifyBearerToken(tokenBytes []byte, expectedIssuer string, now time.Time) (BearerToken, error) {
	var message cose.Sign1Message
	if unmarshalErr := message.UnmarshalCBOR(tokenBytes); unmarshalErr != nil {
		return BearerToken{}, unmarshalErr
	}
	verifier, verifierErr := cose.NewVerifier(cose.AlgorithmEdDSA, graph.DeterministicPublicKey(expectedIssuer))
	if verifierErr != nil {
		return BearerToken{}, verifierErr
	}
	if verifyErr := message.Verify(nil, verifier); verifyErr != nil {
		return BearerToken{}, verifyErr
	}
	token, tokenErr := bearerFromClaims(message.Payload)
	if tokenErr != nil {
		return BearerToken{}, tokenErr
	}
	if token.Issuer != expectedIssuer {
		return BearerToken{}, fmt.Errorf("bearer token issuer=%s, want %s", token.Issuer, expectedIssuer)
	}
	if token.ExpiresUnix <= now.Unix() {
		return BearerToken{}, fmt.Errorf("bearer token expired at %d", token.ExpiresUnix)
	}
	return token, nil
}

func signTokenBytes(issuer string, payloadBytes []byte) ([]byte, error) {
	signer, signerErr := cose.NewSigner(cose.AlgorithmEdDSA, graph.DeterministicPrivateKey(issuer))
	if signerErr != nil {
		return nil, signerErr
	}
	headers := cose.Headers{
		Protected:   cose.ProtectedHeader{cose.HeaderLabelAlgorithm: cose.AlgorithmEdDSA},
		Unprotected: cose.UnprotectedHeader{},
	}
	return cose.Sign1(rand.Reader, signer, headers, payloadBytes, nil)
}

func issueCapabilityBytes(token CapabilityToken) ([]byte, error) {
	if validateErr := validateCapabilityToken(token); validateErr != nil {
		return nil, validateErr
	}
	payloadBytes, payloadErr := capabilityClaims(token)
	if payloadErr != nil {
		return nil, payloadErr
	}
	return signTokenBytes(token.Issuer, payloadBytes)
}

func bearerClaims(token BearerToken) ([]byte, error) {
	claims := map[any]any{
		cose.CWTClaimIssuer:         token.Issuer,
		cose.CWTClaimExpirationTime: token.ExpiresUnix,
		cose.CWTClaimCWTID:          []byte(token.Nonce),
		claimCapability:             bearerPaymentCapability,
		claimScope:                  token.Scope,
		claimObjectCID:              token.ObjectCID,
		claimTokenValue:             token.Value,
		claimUnit:                   token.Unit,
		claimTransferable:           token.Transferable,
		claimServiceCapability:      token.RedeemableCapabilities,
	}
	return store.MarshalCBOR(claims)
}

func capabilityClaims(token CapabilityToken) ([]byte, error) {
	claims := map[any]any{
		cose.CWTClaimIssuer:         token.Issuer,
		cose.CWTClaimSubject:        token.Subject,
		cose.CWTClaimExpirationTime: token.ExpiresUnix,
		cose.CWTClaimCWTID:          []byte(token.Nonce),
		claimCapability:             serviceCapabilityPromise,
		claimScope:                  token.Scope,
		claimObjectCID:              token.ObjectCID,
		claimTokenValue:             token.Value,
		claimUnit:                   token.Unit,
		claimTransferable:           token.Transferable,
		claimServiceCapability:      token.Capability,
		claimParentBearerCID:        token.ParentBearerCID,
	}
	return store.MarshalCBOR(claims)
}

func bearerFromClaims(payloadBytes []byte) (BearerToken, error) {
	var claims map[any]any
	if unmarshalErr := cbor.Unmarshal(payloadBytes, &claims); unmarshalErr != nil {
		return BearerToken{}, unmarshalErr
	}
	capability, capabilityErr := stringClaim(claims, claimCapability)
	if capabilityErr != nil {
		return BearerToken{}, capabilityErr
	}
	if capability != bearerPaymentCapability {
		return BearerToken{}, fmt.Errorf("bearer token capability mismatch")
	}
	nonce, nonceErr := bytesClaim(claims, cose.CWTClaimCWTID)
	if nonceErr != nil {
		return BearerToken{}, nonceErr
	}
	transferable, transferableErr := boolClaim(claims, claimTransferable)
	if transferableErr != nil {
		return BearerToken{}, transferableErr
	}
	capabilities, capabilitiesErr := stringSliceClaim(claims, claimServiceCapability)
	if capabilitiesErr != nil {
		return BearerToken{}, capabilitiesErr
	}
	token := BearerToken{
		Issuer:                 mustStringClaim(claims, cose.CWTClaimIssuer),
		Scope:                  mustStringClaim(claims, claimScope),
		ObjectCID:              mustStringClaim(claims, claimObjectCID),
		Value:                  mustInt64Claim(claims, claimTokenValue),
		Unit:                   mustStringClaim(claims, claimUnit),
		ExpiresUnix:            mustInt64Claim(claims, cose.CWTClaimExpirationTime),
		Nonce:                  string(nonce),
		Transferable:           transferable,
		RedeemableCapabilities: capabilities,
	}
	if validateErr := validateBearerToken(token); validateErr != nil {
		return BearerToken{}, validateErr
	}
	return token, nil
}

func (token BearerToken) matches(expected ExpectedBearerPayment) error {
	if token.Scope != expected.Scope {
		return fmt.Errorf("bearer token scope=%s, want %s", token.Scope, expected.Scope)
	}
	if token.ObjectCID != expected.ObjectCID {
		return fmt.Errorf("bearer token object_cid=%s, want %s", token.ObjectCID, expected.ObjectCID)
	}
	if token.Value != expected.Value {
		return fmt.Errorf("bearer token value=%d, want %d", token.Value, expected.Value)
	}
	if token.Unit != expected.Unit {
		return fmt.Errorf("bearer token unit=%s, want %s", token.Unit, expected.Unit)
	}
	if !containsString(token.RedeemableCapabilities, expected.Capability) {
		return fmt.Errorf("bearer token does not promise capability %s", expected.Capability)
	}
	return nil
}

func validateBearerToken(token BearerToken) error {
	if token.Issuer == "" || token.Scope == "" || token.ObjectCID == "" || token.Unit == "" || token.Nonce == "" {
		return fmt.Errorf("bearer token needs issuer, scope, object CID, unit, and nonce")
	}
	if token.Value <= 0 {
		return fmt.Errorf("bearer token value must be positive")
	}
	if token.ExpiresUnix == 0 {
		return fmt.Errorf("bearer token needs expiration")
	}
	if !token.Transferable {
		return fmt.Errorf("bearer token must be transferable")
	}
	if len(token.RedeemableCapabilities) == 0 {
		return fmt.Errorf("bearer token needs at least one redeemable capability")
	}
	if _, parseErr := store.ParseCIDText(token.ObjectCID); parseErr != nil {
		return parseErr
	}
	return nil
}

func validateCapabilityToken(token CapabilityToken) error {
	if token.Issuer == "" || token.Subject == "" || token.Capability == "" || token.Scope == "" || token.ObjectCID == "" || token.Unit == "" || token.Nonce == "" || token.ParentBearerCID == "" {
		return fmt.Errorf("capability token needs issuer, subject, capability, scope, object CID, unit, nonce, and parent bearer CID")
	}
	if token.Transferable {
		return fmt.Errorf("capability token must be non-transferable")
	}
	if token.Value <= 0 {
		return fmt.Errorf("capability token value must be positive")
	}
	if token.ExpiresUnix == 0 {
		return fmt.Errorf("capability token needs expiration")
	}
	if _, parseErr := store.ParseCIDText(token.ObjectCID); parseErr != nil {
		return parseErr
	}
	if _, parseErr := store.ParseCIDText(token.ParentBearerCID); parseErr != nil {
		return parseErr
	}
	return nil
}

func stringSliceClaim(claims map[any]any, key any) ([]string, error) {
	value, ok := claimValue(claims, key)
	if !ok {
		return nil, fmt.Errorf("missing string list claim %v", key)
	}
	items, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("claim %v must be array", key)
	}
	result := make([]string, 0, len(items))
	for _, item := range items {
		text, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("claim %v item must be text", key)
		}
		result = append(result, text)
	}
	return result, nil
}

func containsString(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}
