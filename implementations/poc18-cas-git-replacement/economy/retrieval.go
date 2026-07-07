package economy

import (
	"fmt"
	"time"

	"github.com/fxamacker/cbor/v2"
	cose "github.com/veraison/go-cose"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/graph"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
)

const (
	objectRetrievalCapability = "poc18-object-retrieval"
	claimObjectCIDs           = int64(-90201)
)

// RetrievalCapabilityToken is a CWT/COSE promise by the issuer that the named
// subject may redeem this exact token for object bytes in the stated scope.
type RetrievalCapabilityToken struct {
	Issuer       string   `json:"issuer"`
	Subject      string   `json:"subject"`
	Scope        string   `json:"scope"`
	ObjectCIDs   []string `json:"object_cids"`
	ExpiresUnix  int64    `json:"expires_unix"`
	Nonce        string   `json:"nonce"`
	Transferable bool     `json:"transferable"`
}

// IssuedRetrievalCapability records exact signed token bytes after local CAS
// storage by the issuer.
type IssuedRetrievalCapability struct {
	Token RetrievalCapabilityToken `json:"token"`
	CID   string                   `json:"cid"`
	Entry store.Entry              `json:"entry"`
	Bytes []byte                   `json:"-"`
}

// ExpectedRetrievalCapability is the local verifier policy for accepting a
// retrieval token.
type ExpectedRetrievalCapability struct {
	Issuer     string
	Subject    string
	Scope      string
	ObjectCIDs []string
}

// IssueRetrievalCapability signs an object-retrieval capability as COSE_Sign1
// and stores the exact token bytes in the issuer's local sparse CAS.
//
// Intent: POC18 TCP retrieval must use cryptographically checkable capability
// promises, not unstructured text tokens, while still leaving acceptance as a
// local verifier decision. Source: DI-koriz
func IssueRetrievalCapability(cas *store.FileStore, token RetrievalCapabilityToken) (IssuedRetrievalCapability, error) {
	if cas == nil {
		return IssuedRetrievalCapability{}, fmt.Errorf("issuer CAS is required")
	}
	if validateErr := validateRetrievalCapability(token); validateErr != nil {
		return IssuedRetrievalCapability{}, validateErr
	}
	payloadBytes, payloadErr := retrievalClaims(token)
	if payloadErr != nil {
		return IssuedRetrievalCapability{}, payloadErr
	}
	tokenBytes, signErr := signTokenBytes(token.Issuer, payloadBytes)
	if signErr != nil {
		return IssuedRetrievalCapability{}, signErr
	}
	entry, putErr := cas.Put("capability_token", tokenBytes)
	if putErr != nil {
		return IssuedRetrievalCapability{}, putErr
	}
	return IssuedRetrievalCapability{Token: token, CID: entry.CID, Entry: entry, Bytes: tokenBytes}, nil
}

// VerifyRetrievalCapability verifies a CWT/COSE retrieval promise against one
// agent's local acceptance policy.
func VerifyRetrievalCapability(tokenBytes []byte, expected ExpectedRetrievalCapability, now time.Time) (RetrievalCapabilityToken, error) {
	var message cose.Sign1Message
	if unmarshalErr := message.UnmarshalCBOR(tokenBytes); unmarshalErr != nil {
		return RetrievalCapabilityToken{}, unmarshalErr
	}
	verifier, verifierErr := cose.NewVerifier(cose.AlgorithmEdDSA, graph.DeterministicPublicKey(expected.Issuer))
	if verifierErr != nil {
		return RetrievalCapabilityToken{}, verifierErr
	}
	if verifyErr := message.Verify(nil, verifier); verifyErr != nil {
		return RetrievalCapabilityToken{}, verifyErr
	}
	token, tokenErr := retrievalFromClaims(message.Payload)
	if tokenErr != nil {
		return RetrievalCapabilityToken{}, tokenErr
	}
	if token.Issuer != expected.Issuer {
		return RetrievalCapabilityToken{}, fmt.Errorf("retrieval token issuer=%s, want %s", token.Issuer, expected.Issuer)
	}
	if token.Subject != expected.Subject {
		return RetrievalCapabilityToken{}, fmt.Errorf("retrieval token subject=%s, want %s", token.Subject, expected.Subject)
	}
	if token.Scope != expected.Scope {
		return RetrievalCapabilityToken{}, fmt.Errorf("retrieval token scope=%s, want %s", token.Scope, expected.Scope)
	}
	if token.ExpiresUnix <= now.Unix() {
		return RetrievalCapabilityToken{}, fmt.Errorf("retrieval token expired at %d", token.ExpiresUnix)
	}
	if !sameStringSet(token.ObjectCIDs, expected.ObjectCIDs) {
		return RetrievalCapabilityToken{}, fmt.Errorf("retrieval token object set mismatch")
	}
	return token, nil
}

func retrievalClaims(token RetrievalCapabilityToken) ([]byte, error) {
	claims := map[any]any{
		cose.CWTClaimIssuer:         token.Issuer,
		cose.CWTClaimSubject:        token.Subject,
		cose.CWTClaimExpirationTime: token.ExpiresUnix,
		cose.CWTClaimCWTID:          []byte(token.Nonce),
		claimCapability:             objectRetrievalCapability,
		claimScope:                  token.Scope,
		claimObjectCIDs:             token.ObjectCIDs,
		claimTransferable:           token.Transferable,
	}
	return store.MarshalCBOR(claims)
}

func retrievalFromClaims(payloadBytes []byte) (RetrievalCapabilityToken, error) {
	var claims map[any]any
	if unmarshalErr := cbor.Unmarshal(payloadBytes, &claims); unmarshalErr != nil {
		return RetrievalCapabilityToken{}, unmarshalErr
	}
	capability, capabilityErr := stringClaim(claims, claimCapability)
	if capabilityErr != nil {
		return RetrievalCapabilityToken{}, capabilityErr
	}
	if capability != objectRetrievalCapability {
		return RetrievalCapabilityToken{}, fmt.Errorf("retrieval token capability mismatch")
	}
	nonce, nonceErr := bytesClaim(claims, cose.CWTClaimCWTID)
	if nonceErr != nil {
		return RetrievalCapabilityToken{}, nonceErr
	}
	transferable, transferableErr := boolClaim(claims, claimTransferable)
	if transferableErr != nil {
		return RetrievalCapabilityToken{}, transferableErr
	}
	objectCIDs, objectCIDsErr := stringSliceClaim(claims, claimObjectCIDs)
	if objectCIDsErr != nil {
		return RetrievalCapabilityToken{}, objectCIDsErr
	}
	token := RetrievalCapabilityToken{
		Issuer:       mustStringClaim(claims, cose.CWTClaimIssuer),
		Subject:      mustStringClaim(claims, cose.CWTClaimSubject),
		Scope:        mustStringClaim(claims, claimScope),
		ObjectCIDs:   objectCIDs,
		ExpiresUnix:  mustInt64Claim(claims, cose.CWTClaimExpirationTime),
		Nonce:        string(nonce),
		Transferable: transferable,
	}
	if validateErr := validateRetrievalCapability(token); validateErr != nil {
		return RetrievalCapabilityToken{}, validateErr
	}
	return token, nil
}

func validateRetrievalCapability(token RetrievalCapabilityToken) error {
	if token.Issuer == "" || token.Subject == "" || token.Scope == "" || token.Nonce == "" {
		return fmt.Errorf("retrieval token needs issuer, subject, scope, and nonce")
	}
	if token.Transferable {
		return fmt.Errorf("retrieval token must be non-transferable")
	}
	if token.ExpiresUnix == 0 {
		return fmt.Errorf("retrieval token needs expiration")
	}
	if len(token.ObjectCIDs) == 0 {
		return fmt.Errorf("retrieval token needs at least one object CID")
	}
	for _, objectCID := range token.ObjectCIDs {
		if _, parseErr := store.ParseCIDText(objectCID); parseErr != nil {
			return parseErr
		}
	}
	return nil
}

func sameStringSet(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	counts := map[string]int{}
	for _, value := range left {
		counts[value]++
	}
	for _, value := range right {
		counts[value]--
	}
	for _, count := range counts {
		if count != 0 {
			return false
		}
	}
	return true
}
