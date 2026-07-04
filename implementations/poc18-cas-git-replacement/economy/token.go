// Package economy models POC18-local promise-token economics.
//
// Intent: POC18 needs spendable storage-payment pressure without inventing a
// global currency, authorization service, or settlement authority. Tokens here
// are issuer promises encoded as signed CWT-style CBOR and locally redeemed by
// the agent deciding whether the reciprocal storage promise is worth making.
// Source: DI-bidum
package economy

import (
	"crypto/rand"
	"fmt"
	"math"
	"time"

	"github.com/fxamacker/cbor/v2"
	cose "github.com/veraison/go-cose"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/graph"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
)

const (
	storagePaymentCapability = "poc18-storage-payment"
	claimCapability          = int64(-90001)
	claimScope               = int64(-90002)
	claimObjectCID           = int64(-90003)
	claimTokenValue          = int64(-90004)
	claimUnit                = int64(-90005)
	claimTransferable        = int64(-90006)
)

// StoragePaymentToken is Alice's signed promise that a bearer can spend a local
// storage credit for a named POC18 retention scope and object.
type StoragePaymentToken struct {
	Issuer       string `json:"issuer"`
	Subject      string `json:"subject"`
	Scope        string `json:"scope"`
	ObjectCID    string `json:"object_cid"`
	Value        int64  `json:"value"`
	Unit         string `json:"unit"`
	ExpiresUnix  int64  `json:"expires_unix"`
	Nonce        string `json:"nonce"`
	Transferable bool   `json:"transferable"`
}

// IssuedToken records exact signed token bytes and their CID after the issuer
// stores them in its sparse CAS.
type IssuedToken struct {
	Token StoragePaymentToken `json:"token"`
	CID   string              `json:"cid"`
	Entry store.Entry         `json:"entry"`
	Bytes []byte              `json:"-"`
}

// ExpectedStoragePayment is Frank's local acceptance policy before a token is
// treated as reciprocal economics for a retention promise.
type ExpectedStoragePayment struct {
	Issuer    string
	Subject   string
	Scope     string
	ObjectCID string
	Value     int64
	Unit      string
}

// RedemptionReport is the fixture-visible result of verifying and consuming one
// storage-payment bearer token.
type RedemptionReport struct {
	Issuer                string `json:"issuer"`
	Subject               string `json:"subject"`
	Redeemer              string `json:"redeemer"`
	Scope                 string `json:"scope"`
	ObjectCID             string `json:"object_cid"`
	TokenCID              string `json:"token_cid"`
	RedemptionMessageCID  string `json:"redemption_message_cid"`
	Value                 int64  `json:"value"`
	Unit                  string `json:"unit"`
	Transferable          bool   `json:"transferable"`
	Issued                bool   `json:"issued"`
	TokenStoredByIssuer   bool   `json:"token_stored_by_issuer"`
	TokenStoredByRedeemer bool   `json:"token_stored_by_redeemer"`
	SignatureVerified     bool   `json:"signature_verified"`
	Redeemed              bool   `json:"redeemed"`
	ReplayRejected        bool   `json:"replay_rejected"`
}

// Ledger is a local spent-token set. It deliberately has no global authority:
// it only records which token CIDs this agent has already accepted.
type Ledger struct {
	spent map[string]bool
}

// NewLedger returns an empty local spent-token ledger.
func NewLedger() *Ledger {
	return &Ledger{spent: map[string]bool{}}
}

// IssueStoragePaymentToken signs token claims as COSE_Sign1 and stores the exact
// token bytes in the issuer's CAS.
//
// Intent: Payment must be a cryptographically verifiable promise object, not a
// text note in reciprocal terms. The token CID names exact signed bytes that can
// be transferred and later redeemed once by a local ledger. Source: DI-bidum
func IssueStoragePaymentToken(cas *store.FileStore, token StoragePaymentToken) (IssuedToken, error) {
	if cas == nil {
		return IssuedToken{}, fmt.Errorf("issuer CAS is required")
	}
	if validateErr := validateToken(token); validateErr != nil {
		return IssuedToken{}, validateErr
	}
	payloadBytes, payloadErr := tokenClaims(token)
	if payloadErr != nil {
		return IssuedToken{}, payloadErr
	}
	signer, signerErr := cose.NewSigner(cose.AlgorithmEdDSA, graph.DeterministicPrivateKey(token.Issuer))
	if signerErr != nil {
		return IssuedToken{}, signerErr
	}
	headers := cose.Headers{
		Protected: cose.ProtectedHeader{
			cose.HeaderLabelAlgorithm: cose.AlgorithmEdDSA,
		},
		Unprotected: cose.UnprotectedHeader{},
	}
	tokenBytes, signErr := cose.Sign1(rand.Reader, signer, headers, payloadBytes, nil)
	if signErr != nil {
		return IssuedToken{}, signErr
	}
	entry, putErr := cas.Put("payment_token", tokenBytes)
	if putErr != nil {
		return IssuedToken{}, putErr
	}
	return IssuedToken{Token: token, CID: entry.CID, Entry: entry, Bytes: tokenBytes}, nil
}

// RedeemStoragePaymentToken verifies a bearer token, stores exact token bytes in
// the redeemer's CAS, and consumes the token CID in the redeemer's local ledger.
func (ledger *Ledger) RedeemStoragePaymentToken(cas *store.FileStore, tokenBytes []byte, expected ExpectedStoragePayment, redeemedBy string, now time.Time) (RedemptionReport, error) {
	if ledger == nil {
		return RedemptionReport{}, fmt.Errorf("payment ledger is required")
	}
	if cas == nil {
		return RedemptionReport{}, fmt.Errorf("redeemer CAS is required")
	}
	if redeemedBy == "" {
		return RedemptionReport{}, fmt.Errorf("redeemer is required")
	}
	token, verifyErr := VerifyStoragePaymentToken(tokenBytes, expected.Issuer, now)
	if verifyErr != nil {
		return RedemptionReport{}, verifyErr
	}
	if termsErr := token.matches(expected); termsErr != nil {
		return RedemptionReport{}, termsErr
	}
	tokenCID := store.CIDText(store.CIDForBytes(tokenBytes))
	if ledger.spent[tokenCID] {
		return RedemptionReport{}, fmt.Errorf("storage payment token %s already redeemed", tokenCID)
	}
	entry, putErr := cas.Put("payment_token", tokenBytes)
	if putErr != nil {
		return RedemptionReport{}, putErr
	}
	ledger.spent[tokenCID] = true
	return RedemptionReport{
		Issuer:                token.Issuer,
		Subject:               token.Subject,
		Redeemer:              redeemedBy,
		Scope:                 token.Scope,
		ObjectCID:             token.ObjectCID,
		TokenCID:              entry.CID,
		Value:                 token.Value,
		Unit:                  token.Unit,
		Transferable:          token.Transferable,
		Issued:                true,
		TokenStoredByRedeemer: true,
		SignatureVerified:     true,
		Redeemed:              true,
	}, nil
}

// StoreRedemptionMessage records Frank's local promise that the token was
// consumed as the reciprocal economics for the related retention promise.
func StoreRedemptionMessage(cas *store.FileStore, report RedemptionReport, redeemedAt string) (graph.StoredMessage, error) {
	tokenCID, tokenErr := store.ParseCIDText(report.TokenCID)
	if tokenErr != nil {
		return graph.StoredMessage{}, tokenErr
	}
	objectCID, objectErr := store.ParseCIDText(report.ObjectCID)
	if objectErr != nil {
		return graph.StoredMessage{}, objectErr
	}
	payload := graph.Payload{
		Promiser:    report.Redeemer,
		Promisee:    report.Issuer,
		PromiseKind: "storage_payment_redemption",
		PromiseBody: graph.StoragePaymentRedemptionBody(
			tokenCID,
			objectCID,
			report.Scope,
			report.Value,
			report.Unit,
			redeemedAt,
			report.Transferable,
		),
		ReciprocalPromise: []any{},
		LocalConstraints: []any{
			"issuer signature verified before local redemption",
			"local spent-token ledger rejects replay",
		},
	}
	return graph.StoreMessage(cas, []graph.Parent{{Role: "redeems_token", CID: tokenCID}}, payload)
}

// VerifyStoragePaymentToken verifies the COSE signature and returns the signed
// CWT-style storage-payment claims.
func VerifyStoragePaymentToken(tokenBytes []byte, expectedIssuer string, now time.Time) (StoragePaymentToken, error) {
	var message cose.Sign1Message
	if unmarshalErr := message.UnmarshalCBOR(tokenBytes); unmarshalErr != nil {
		return StoragePaymentToken{}, unmarshalErr
	}
	verifier, verifierErr := cose.NewVerifier(cose.AlgorithmEdDSA, graph.DeterministicPublicKey(expectedIssuer))
	if verifierErr != nil {
		return StoragePaymentToken{}, verifierErr
	}
	if verifyErr := message.Verify(nil, verifier); verifyErr != nil {
		return StoragePaymentToken{}, verifyErr
	}
	token, tokenErr := tokenFromClaims(message.Payload)
	if tokenErr != nil {
		return StoragePaymentToken{}, tokenErr
	}
	if token.Issuer != expectedIssuer {
		return StoragePaymentToken{}, fmt.Errorf("storage payment token issuer=%s, want %s", token.Issuer, expectedIssuer)
	}
	if token.ExpiresUnix <= now.Unix() {
		return StoragePaymentToken{}, fmt.Errorf("storage payment token expired at %d", token.ExpiresUnix)
	}
	return token, nil
}

func (token StoragePaymentToken) matches(expected ExpectedStoragePayment) error {
	if token.Subject != expected.Subject {
		return fmt.Errorf("storage payment token subject=%s, want %s", token.Subject, expected.Subject)
	}
	if token.Scope != expected.Scope {
		return fmt.Errorf("storage payment token scope=%s, want %s", token.Scope, expected.Scope)
	}
	if token.ObjectCID != expected.ObjectCID {
		return fmt.Errorf("storage payment token object_cid=%s, want %s", token.ObjectCID, expected.ObjectCID)
	}
	if token.Value != expected.Value {
		return fmt.Errorf("storage payment token value=%d, want %d", token.Value, expected.Value)
	}
	if token.Unit != expected.Unit {
		return fmt.Errorf("storage payment token unit=%s, want %s", token.Unit, expected.Unit)
	}
	return nil
}

func validateToken(token StoragePaymentToken) error {
	if token.Issuer == "" || token.Subject == "" || token.Scope == "" || token.ObjectCID == "" || token.Unit == "" || token.Nonce == "" {
		return fmt.Errorf("storage payment token needs issuer, subject, scope, object CID, unit, and nonce")
	}
	if token.Value <= 0 {
		return fmt.Errorf("storage payment token value must be positive")
	}
	if token.ExpiresUnix == 0 {
		return fmt.Errorf("storage payment token needs expiration")
	}
	if _, parseErr := store.ParseCIDText(token.ObjectCID); parseErr != nil {
		return parseErr
	}
	return nil
}

func tokenClaims(token StoragePaymentToken) ([]byte, error) {
	claims := map[any]any{
		cose.CWTClaimIssuer:         token.Issuer,
		cose.CWTClaimSubject:        token.Subject,
		cose.CWTClaimExpirationTime: token.ExpiresUnix,
		cose.CWTClaimCWTID:          []byte(token.Nonce),
		claimCapability:             storagePaymentCapability,
		claimScope:                  token.Scope,
		claimObjectCID:              token.ObjectCID,
		claimTokenValue:             token.Value,
		claimUnit:                   token.Unit,
		claimTransferable:           token.Transferable,
	}
	return store.MarshalCBOR(claims)
}

func tokenFromClaims(payloadBytes []byte) (StoragePaymentToken, error) {
	var claims map[any]any
	if unmarshalErr := cbor.Unmarshal(payloadBytes, &claims); unmarshalErr != nil {
		return StoragePaymentToken{}, unmarshalErr
	}
	nonce, nonceErr := bytesClaim(claims, cose.CWTClaimCWTID)
	if nonceErr != nil {
		return StoragePaymentToken{}, nonceErr
	}
	transferable, transferableErr := boolClaim(claims, claimTransferable)
	if transferableErr != nil {
		return StoragePaymentToken{}, transferableErr
	}
	if capability, capabilityErr := stringClaim(claims, claimCapability); capabilityErr != nil {
		return StoragePaymentToken{}, capabilityErr
	} else if capability != storagePaymentCapability {
		return StoragePaymentToken{}, fmt.Errorf("storage payment token capability mismatch")
	}
	token := StoragePaymentToken{
		Issuer:       mustStringClaim(claims, cose.CWTClaimIssuer),
		Subject:      mustStringClaim(claims, cose.CWTClaimSubject),
		Scope:        mustStringClaim(claims, claimScope),
		ObjectCID:    mustStringClaim(claims, claimObjectCID),
		Value:        mustInt64Claim(claims, claimTokenValue),
		Unit:         mustStringClaim(claims, claimUnit),
		ExpiresUnix:  mustInt64Claim(claims, cose.CWTClaimExpirationTime),
		Nonce:        string(nonce),
		Transferable: transferable,
	}
	if validateErr := validateToken(token); validateErr != nil {
		return StoragePaymentToken{}, validateErr
	}
	return token, nil
}

func mustStringClaim(claims map[any]any, key any) string {
	value, err := stringClaim(claims, key)
	if err != nil {
		return ""
	}
	return value
}

func mustInt64Claim(claims map[any]any, key any) int64 {
	value, err := int64Claim(claims, key)
	if err != nil {
		return 0
	}
	return value
}

func stringClaim(claims map[any]any, key any) (string, error) {
	value, ok := claimValue(claims, key)
	if !ok {
		return "", fmt.Errorf("missing string claim %v", key)
	}
	text, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("claim %v must be text", key)
	}
	return text, nil
}

func bytesClaim(claims map[any]any, key any) ([]byte, error) {
	value, ok := claimValue(claims, key)
	if !ok {
		return nil, fmt.Errorf("missing bytes claim %v", key)
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil, fmt.Errorf("claim %v must be bytes", key)
	}
	return bytes, nil
}

func int64Claim(claims map[any]any, key any) (int64, error) {
	value, ok := claimValue(claims, key)
	if !ok {
		return 0, fmt.Errorf("missing integer claim %v", key)
	}
	switch typed := value.(type) {
	case int64:
		return typed, nil
	case uint64:
		if typed > math.MaxInt64 {
			return 0, fmt.Errorf("claim %v integer is too large", key)
		}
		return int64(typed), nil
	case int:
		return int64(typed), nil
	default:
		return 0, fmt.Errorf("claim %v must be integer", key)
	}
}

func boolClaim(claims map[any]any, key any) (bool, error) {
	value, ok := claimValue(claims, key)
	if !ok {
		return false, fmt.Errorf("missing bool claim %v", key)
	}
	typed, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("claim %v must be bool", key)
	}
	return typed, nil
}

func claimValue(claims map[any]any, key any) (any, bool) {
	for claimKey, value := range claims {
		if claimKeysEqual(claimKey, key) {
			return value, true
		}
	}
	return nil, false
}

func claimKeysEqual(left any, right any) bool {
	leftInt, leftOK := claimKeyInt64(left)
	rightInt, rightOK := claimKeyInt64(right)
	if leftOK && rightOK {
		return leftInt == rightInt
	}
	return left == right
}

func claimKeyInt64(value any) (int64, bool) {
	switch typed := value.(type) {
	case int64:
		return typed, true
	case uint64:
		if typed > math.MaxInt64 {
			return 0, false
		}
		return int64(typed), true
	case int:
		return int64(typed), true
	default:
		return 0, false
	}
}
