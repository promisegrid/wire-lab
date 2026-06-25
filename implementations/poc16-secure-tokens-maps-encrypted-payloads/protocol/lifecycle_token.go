package protocol

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fxamacker/cbor/v2"
	cose "github.com/veraison/go-cose"
)

const (
	LifecycleKindTokenIssued    = "token_issued"
	LifecycleKindReady          = "ready"
	LifecycleKindTokenInvoked   = "token_invoked"
	LifecycleKindTokenFulfilled = "token_fulfilled"

	LifecycleChannelTCP        = "tcp"
	LifecycleChannelParserPath = "parser_path"
	LifecycleChannelStdio      = "stdio"

	LifecycleOutcomeKept     = "kept"
	LifecycleOutcomeBroken   = "broken"
	LifecycleOutcomeRejected = "rejected"
	LifecycleOutcomeTimedOut = "timed_out"
)

const (
	cwtPrivateProtocolCID      = int64(-70020)
	cwtPrivateRunID            = int64(-70021)
	cwtPrivateRoleKind         = int64(-70022)
	cwtPrivateRoleID           = int64(-70023)
	cwtPrivateChannelProfile   = int64(-70024)
	cwtPrivateShutdownTerms    = int64(-70025)
	cwtPrivateGraceMillis      = int64(-70026)
	cwtPrivateMaxInvocations   = int64(-70027)
	cwtPrivateTokenProfile     = int64(-70028)
	localLifecycleTokenSubject = "local lifecycle shutdown promise"
	localLifecycleTokenProfile = "local_lifecycle_cwt_cose_sign1_eddsa_v1"
)

var lifecycleCBOREncMode = mustLifecycleCBOREncMode()

// LifecycleTokenTerms are the signed promise terms inside a lifecycle
// capability token.
//
// Intent: A local role signs its own shutdown promise before the supervisor ever
// invokes it, so shutdown is token presentation and fulfillment rather than a
// supervisor command. Source: DI-jafoj
type LifecycleTokenTerms struct {
	IssuerRoleID   string
	AudienceRoleID string
	RunID          string
	RoleKind       string
	ChannelProfile string
	ProtocolCID    ProtocolCID
	GraceMillis    int64
	MaxInvocations int64
	IssuedAtUnix   int64
	NotBeforeUnix  int64
	ExpiresUnix    int64
	TokenID        string
	ShutdownTerms  []string
}

// LifecycleToken is the exact COSE/CWT lifecycle capability token plus
// diagnostic identifiers used in local events.
type LifecycleToken struct {
	Terms        LifecycleTokenTerms
	COSEBytes    []byte
	CID          string
	PublicKey    []byte
	COSEBase64   string
	PublicBase64 string
}

// LifecyclePayload is the pCID-owned map payload for local_lifecycle_v1.
type LifecyclePayload struct {
	Kind            string
	Promiser        string
	Promisee        string
	RoleID          string
	RoleKind        string
	ChannelProfile  string
	RunID           string
	TokenCOSEBase64 string
	TokenCID        string
	PublicKeyBase64 string
	Reason          string
	DeadlineUnix    int64
	Outcome         string
	Detail          string
}

// IssueLifecycleToken signs the role's local shutdown promise as a CWT payload
// protected by COSE_Sign1/Ed25519 using the well-known go-cose implementation.
func IssueLifecycleToken(terms LifecycleTokenTerms) (LifecycleToken, error) {
	if err := validateLifecycleTerms(terms); err != nil {
		return LifecycleToken{}, err
	}
	payloadBytes, payloadErr := marshalLifecycleTokenTerms(terms)
	if payloadErr != nil {
		return LifecycleToken{}, payloadErr
	}
	signer, signerErr := cose.NewSigner(cose.AlgorithmEdDSA, DeterministicPrivateKey(terms.IssuerRoleID))
	if signerErr != nil {
		return LifecycleToken{}, signerErr
	}
	headers := cose.Headers{
		Protected: cose.ProtectedHeader{
			cose.HeaderLabelAlgorithm: cose.AlgorithmEdDSA,
		},
		Unprotected: cose.UnprotectedHeader{},
	}
	coseBytes, coseErr := cose.Sign1(rand.Reader, signer, headers, payloadBytes, nil)
	if coseErr != nil {
		return LifecycleToken{}, coseErr
	}
	copiedCOSE := append([]byte(nil), coseBytes...)
	publicKey := DeterministicPublicKey(terms.IssuerRoleID)
	copiedPublic := append([]byte(nil), publicKey...)
	return LifecycleToken{
		Terms:        terms,
		COSEBytes:    copiedCOSE,
		CID:          CIDForExactBytes(copiedCOSE),
		PublicKey:    copiedPublic,
		COSEBase64:   base64.StdEncoding.EncodeToString(copiedCOSE),
		PublicBase64: base64.StdEncoding.EncodeToString(copiedPublic),
	}, nil
}

// VerifyLifecycleToken verifies the presented COSE/CWT token and checks local
// invocation terms before a role treats it as its own shutdown promise.
func VerifyLifecycleToken(coseBytes []byte, expected LifecycleTokenTerms, now time.Time, priorInvocations int64) (LifecycleTokenTerms, error) {
	var message cose.Sign1Message
	if unmarshalErr := message.UnmarshalCBOR(coseBytes); unmarshalErr != nil {
		return LifecycleTokenTerms{}, unmarshalErr
	}
	verifier, verifierErr := cose.NewVerifier(cose.AlgorithmEdDSA, DeterministicPublicKey(expected.IssuerRoleID))
	if verifierErr != nil {
		return LifecycleTokenTerms{}, verifierErr
	}
	if verifyErr := message.Verify(nil, verifier); verifyErr != nil {
		return LifecycleTokenTerms{}, verifyErr
	}
	terms, termsErr := unmarshalLifecycleTokenTerms(message.Payload)
	if termsErr != nil {
		return LifecycleTokenTerms{}, termsErr
	}
	if err := compareLifecycleTerms(terms, expected, now, priorInvocations); err != nil {
		return LifecycleTokenTerms{}, err
	}
	return terms, nil
}

// NewLifecycleTokenIssuedPayload describes the startup promise token a role
// voluntarily gives to its local supervisor.
func NewLifecycleTokenIssuedPayload(token LifecycleToken) LifecyclePayload {
	return LifecyclePayload{
		Kind:            LifecycleKindTokenIssued,
		Promiser:        token.Terms.IssuerRoleID,
		Promisee:        token.Terms.AudienceRoleID,
		RoleID:          token.Terms.IssuerRoleID,
		RoleKind:        token.Terms.RoleKind,
		ChannelProfile:  token.Terms.ChannelProfile,
		RunID:           token.Terms.RunID,
		TokenCOSEBase64: token.COSEBase64,
		TokenCID:        token.CID,
		PublicKeyBase64: token.PublicBase64,
	}
}

// NewLifecycleReadyPayload tells the supervisor that the role is holding its
// token promise ready for invocation.
func NewLifecycleReadyPayload(token LifecycleToken, detail string) LifecyclePayload {
	payload := NewLifecycleTokenIssuedPayload(token)
	payload.Kind = LifecycleKindReady
	payload.TokenCOSEBase64 = ""
	payload.PublicKeyBase64 = ""
	payload.Detail = detail
	return payload
}

// NewLifecycleInvocationPayload presents an exact signed token back to its
// issuer under the token's local shutdown terms.
func NewLifecycleInvocationPayload(supervisorRoleID string, token LifecycleToken, reason string, deadline time.Time) LifecyclePayload {
	return LifecyclePayload{
		Kind:            LifecycleKindTokenInvoked,
		Promiser:        supervisorRoleID,
		Promisee:        token.Terms.IssuerRoleID,
		RoleID:          token.Terms.IssuerRoleID,
		RoleKind:        token.Terms.RoleKind,
		ChannelProfile:  token.Terms.ChannelProfile,
		RunID:           token.Terms.RunID,
		TokenCOSEBase64: token.COSEBase64,
		TokenCID:        token.CID,
		Reason:          reason,
		DeadlineUnix:    deadline.Unix(),
	}
}

// NewLifecycleFulfilledPayload records the issuer's local outcome after token
// invocation.
func NewLifecycleFulfilledPayload(token LifecycleToken, outcome, detail string) LifecyclePayload {
	return LifecyclePayload{
		Kind:           LifecycleKindTokenFulfilled,
		Promiser:       token.Terms.IssuerRoleID,
		Promisee:       token.Terms.AudienceRoleID,
		RoleID:         token.Terms.IssuerRoleID,
		RoleKind:       token.Terms.RoleKind,
		ChannelProfile: token.Terms.ChannelProfile,
		RunID:          token.Terms.RunID,
		TokenCID:       token.CID,
		Outcome:        outcome,
		Detail:         detail,
	}
}

// EncodeLifecycleMessage writes grid([42(local_lifecycle_v1_pCID), payload])
// without a generic proof slot because COSE/CWT owns token proof for this pCID.
func EncodeLifecycleMessage(protocolCID ProtocolCID, payload LifecyclePayload) ([]byte, error) {
	payloadBytes, marshalErr := MarshalStringMap(payload.lifecycleFields())
	if marshalErr != nil {
		return nil, marshalErr
	}
	return EncodeGridMessage(protocolCID, RawCBORGridSlot(payloadBytes))
}

// DecodeLifecycleMessage parses the pCID-owned lifecycle map payload.
func DecodeLifecycleMessage(frameBytes []byte, expectedProtocolCID ProtocolCID) (LifecyclePayload, error) {
	gridMessage, gridErr := ParseGridMessage(frameBytes)
	if gridErr != nil {
		return LifecyclePayload{}, gridErr
	}
	if !gridMessage.ProtocolCID.Equal(expectedProtocolCID) {
		return LifecyclePayload{}, fmt.Errorf("lifecycle pCID=%s, want %s", gridMessage.ProtocolCID.String(), expectedProtocolCID.String())
	}
	if len(gridMessage.Slots) != 1 {
		return LifecyclePayload{}, fmt.Errorf("local_lifecycle_v1 expects one payload slot, got %d", len(gridMessage.Slots))
	}
	fields, fieldsErr := UnmarshalStringMap(gridMessage.Slots[0].RawCBOR)
	if fieldsErr != nil {
		return LifecyclePayload{}, fieldsErr
	}
	return lifecyclePayloadFromFields(fields)
}

// TokenBytes decodes the exact COSE token bytes carried in a lifecycle payload.
func (payload LifecyclePayload) TokenBytes() ([]byte, error) {
	if strings.TrimSpace(payload.TokenCOSEBase64) == "" {
		return nil, fmt.Errorf("lifecycle payload has no token bytes")
	}
	tokenBytes, decodeErr := base64.StdEncoding.DecodeString(payload.TokenCOSEBase64)
	if decodeErr != nil {
		return nil, decodeErr
	}
	if payload.TokenCID != "" && CIDForExactBytes(tokenBytes) != payload.TokenCID {
		return nil, fmt.Errorf("lifecycle token CID mismatch")
	}
	return tokenBytes, nil
}

func mustLifecycleCBOREncMode() cbor.EncMode {
	mode, err := cbor.CanonicalEncOptions().EncMode()
	if err != nil {
		panic(err)
	}
	return mode
}

func marshalLifecycleTokenTerms(terms LifecycleTokenTerms) ([]byte, error) {
	cwtTerms := map[any]any{
		cose.CWTClaimIssuer:         terms.IssuerRoleID,
		cose.CWTClaimSubject:        localLifecycleTokenSubject,
		cose.CWTClaimAudience:       terms.AudienceRoleID,
		cose.CWTClaimIssuedAt:       terms.IssuedAtUnix,
		cose.CWTClaimNotBefore:      terms.NotBeforeUnix,
		cose.CWTClaimExpirationTime: terms.ExpiresUnix,
		cose.CWTClaimCWTID:          []byte(terms.TokenID),
		cwtPrivateProtocolCID:       terms.ProtocolCID.Bytes(),
		cwtPrivateRunID:             terms.RunID,
		cwtPrivateRoleKind:          terms.RoleKind,
		cwtPrivateRoleID:            terms.IssuerRoleID,
		cwtPrivateChannelProfile:    terms.ChannelProfile,
		cwtPrivateShutdownTerms:     terms.ShutdownTerms,
		cwtPrivateGraceMillis:       terms.GraceMillis,
		cwtPrivateMaxInvocations:    terms.MaxInvocations,
		cwtPrivateTokenProfile:      localLifecycleTokenProfile,
	}
	return lifecycleCBOREncMode.Marshal(cwtTerms)
}

func unmarshalLifecycleTokenTerms(payloadBytes []byte) (LifecycleTokenTerms, error) {
	var cwtTerms map[any]any
	if err := cbor.Unmarshal(payloadBytes, &cwtTerms); err != nil {
		return LifecycleTokenTerms{}, err
	}
	protocolCID, cidErr := protocolCIDTerm(cwtValue(cwtTerms, cwtPrivateProtocolCID))
	if cidErr != nil {
		return LifecycleTokenTerms{}, cidErr
	}
	shutdownTerms, shutdownErr := stringSliceTerm(cwtValue(cwtTerms, cwtPrivateShutdownTerms))
	if shutdownErr != nil {
		return LifecycleTokenTerms{}, shutdownErr
	}
	tokenID, tokenErr := bytesTextTerm(cwtValue(cwtTerms, cose.CWTClaimCWTID))
	if tokenErr != nil {
		return LifecycleTokenTerms{}, tokenErr
	}
	terms := LifecycleTokenTerms{
		IssuerRoleID:   stringTerm(cwtValue(cwtTerms, cose.CWTClaimIssuer)),
		AudienceRoleID: stringTerm(cwtValue(cwtTerms, cose.CWTClaimAudience)),
		RunID:          stringTerm(cwtValue(cwtTerms, cwtPrivateRunID)),
		RoleKind:       stringTerm(cwtValue(cwtTerms, cwtPrivateRoleKind)),
		ChannelProfile: stringTerm(cwtValue(cwtTerms, cwtPrivateChannelProfile)),
		ProtocolCID:    protocolCID,
		GraceMillis:    int64Term(cwtValue(cwtTerms, cwtPrivateGraceMillis)),
		MaxInvocations: int64Term(cwtValue(cwtTerms, cwtPrivateMaxInvocations)),
		IssuedAtUnix:   int64Term(cwtValue(cwtTerms, cose.CWTClaimIssuedAt)),
		NotBeforeUnix:  int64Term(cwtValue(cwtTerms, cose.CWTClaimNotBefore)),
		ExpiresUnix:    int64Term(cwtValue(cwtTerms, cose.CWTClaimExpirationTime)),
		TokenID:        tokenID,
		ShutdownTerms:  shutdownTerms,
	}
	if stringTerm(cwtValue(cwtTerms, cose.CWTClaimSubject)) != localLifecycleTokenSubject {
		return LifecycleTokenTerms{}, fmt.Errorf("lifecycle token subject mismatch")
	}
	if stringTerm(cwtValue(cwtTerms, cwtPrivateTokenProfile)) != localLifecycleTokenProfile {
		return LifecycleTokenTerms{}, fmt.Errorf("lifecycle token profile mismatch")
	}
	if err := validateLifecycleTerms(terms); err != nil {
		return LifecycleTokenTerms{}, err
	}
	return terms, nil
}

func compareLifecycleTerms(actual, expected LifecycleTokenTerms, now time.Time, priorInvocations int64) error {
	if actual.IssuerRoleID != expected.IssuerRoleID {
		return fmt.Errorf("lifecycle token issuer=%s, want %s", actual.IssuerRoleID, expected.IssuerRoleID)
	}
	if actual.AudienceRoleID != expected.AudienceRoleID {
		return fmt.Errorf("lifecycle token audience=%s, want %s", actual.AudienceRoleID, expected.AudienceRoleID)
	}
	if actual.RunID != expected.RunID {
		return fmt.Errorf("lifecycle token run_id=%s, want %s", actual.RunID, expected.RunID)
	}
	if actual.RoleKind != expected.RoleKind {
		return fmt.Errorf("lifecycle token role_kind=%s, want %s", actual.RoleKind, expected.RoleKind)
	}
	if actual.ChannelProfile != expected.ChannelProfile {
		return fmt.Errorf("lifecycle token channel_profile=%s, want %s", actual.ChannelProfile, expected.ChannelProfile)
	}
	if !actual.ProtocolCID.Equal(expected.ProtocolCID) {
		return fmt.Errorf("lifecycle token pCID=%s, want %s", actual.ProtocolCID.String(), expected.ProtocolCID.String())
	}
	if now.Unix() < actual.NotBeforeUnix {
		return fmt.Errorf("lifecycle token not valid before %d", actual.NotBeforeUnix)
	}
	if now.Unix() >= actual.ExpiresUnix {
		return fmt.Errorf("lifecycle token expired at %d", actual.ExpiresUnix)
	}
	if priorInvocations >= actual.MaxInvocations {
		return fmt.Errorf("lifecycle token invocation limit reached")
	}
	return nil
}

func validateLifecycleTerms(terms LifecycleTokenTerms) error {
	if terms.IssuerRoleID == "" || terms.AudienceRoleID == "" || terms.RunID == "" || terms.RoleKind == "" || terms.ChannelProfile == "" || terms.TokenID == "" {
		return fmt.Errorf("lifecycle token needs issuer, audience, run id, role kind, channel profile, and token id")
	}
	if len(terms.ProtocolCID.Bytes()) == 0 {
		return fmt.Errorf("lifecycle token needs protocol CID")
	}
	if terms.GraceMillis <= 0 || terms.MaxInvocations <= 0 {
		return fmt.Errorf("lifecycle token needs positive grace and invocation count")
	}
	if terms.NotBeforeUnix == 0 || terms.ExpiresUnix == 0 || terms.IssuedAtUnix == 0 {
		return fmt.Errorf("lifecycle token needs iat, nbf, and exp")
	}
	if terms.ExpiresUnix <= terms.NotBeforeUnix {
		return fmt.Errorf("lifecycle token exp must be after nbf")
	}
	if len(terms.ShutdownTerms) == 0 {
		return fmt.Errorf("lifecycle token needs shutdown terms")
	}
	return nil
}

func (payload LifecyclePayload) lifecycleFields() map[string]string {
	fields := map[string]string{
		"channel_profile": payload.ChannelProfile,
		"kind":            payload.Kind,
		"promiser":        payload.Promiser,
		"promisee":        payload.Promisee,
		"role_id":         payload.RoleID,
		"role_kind":       payload.RoleKind,
		"run_id":          payload.RunID,
	}
	if payload.TokenCOSEBase64 != "" {
		fields["token_cose_b64"] = payload.TokenCOSEBase64
	}
	if payload.TokenCID != "" {
		fields["token_cid"] = payload.TokenCID
	}
	if payload.PublicKeyBase64 != "" {
		fields["public_key_b64"] = payload.PublicKeyBase64
	}
	if payload.Reason != "" {
		fields["reason"] = payload.Reason
	}
	if payload.DeadlineUnix != 0 {
		fields["deadline_unix"] = strconv.FormatInt(payload.DeadlineUnix, 10)
	}
	if payload.Outcome != "" {
		fields["outcome"] = payload.Outcome
	}
	if payload.Detail != "" {
		fields["detail"] = payload.Detail
	}
	return fields
}

func lifecyclePayloadFromFields(fields map[string]string) (LifecyclePayload, error) {
	deadlineUnix := int64(0)
	if fields["deadline_unix"] != "" {
		parsedDeadline, parseErr := strconv.ParseInt(fields["deadline_unix"], 10, 64)
		if parseErr != nil {
			return LifecyclePayload{}, parseErr
		}
		deadlineUnix = parsedDeadline
	}
	payload := LifecyclePayload{
		Kind:            fields["kind"],
		Promiser:        fields["promiser"],
		Promisee:        fields["promisee"],
		RoleID:          fields["role_id"],
		RoleKind:        fields["role_kind"],
		ChannelProfile:  fields["channel_profile"],
		RunID:           fields["run_id"],
		TokenCOSEBase64: fields["token_cose_b64"],
		TokenCID:        fields["token_cid"],
		PublicKeyBase64: fields["public_key_b64"],
		Reason:          fields["reason"],
		DeadlineUnix:    deadlineUnix,
		Outcome:         fields["outcome"],
		Detail:          fields["detail"],
	}
	if payload.Kind == "" || payload.Promiser == "" || payload.Promisee == "" || payload.RoleID == "" || payload.RoleKind == "" || payload.ChannelProfile == "" || payload.RunID == "" {
		return LifecyclePayload{}, fmt.Errorf("lifecycle payload missing common fields")
	}
	return payload, nil
}

func stringTerm(value any) string {
	if text, ok := value.(string); ok {
		return text
	}
	return ""
}

func cwtValue(values map[any]any, label int64) any {
	if value, ok := values[label]; ok {
		return value
	}
	if label >= 0 {
		if value, ok := values[uint64(label)]; ok {
			return value
		}
		if value, ok := values[int(label)]; ok {
			return value
		}
	}
	return nil
}

func int64Term(value any) int64 {
	switch typed := value.(type) {
	case uint64:
		if typed <= uint64(^uint64(0)>>1) {
			return int64(typed)
		}
	case int64:
		return typed
	case int:
		return int64(typed)
	}
	return 0
}

func bytesTextTerm(value any) (string, error) {
	bytesValue, ok := value.([]byte)
	if !ok {
		return "", fmt.Errorf("expected token id bytes")
	}
	return string(bytesValue), nil
}

func protocolCIDTerm(value any) (ProtocolCID, error) {
	bytesValue, ok := value.([]byte)
	if !ok {
		return ProtocolCID{}, fmt.Errorf("expected protocol CID bytes")
	}
	return ParseProtocolCIDFromBytes(bytesValue)
}

func stringSliceTerm(value any) ([]string, error) {
	rawTerms, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("expected shutdown terms array")
	}
	terms := make([]string, 0, len(rawTerms))
	for _, rawTerm := range rawTerms {
		term, ok := rawTerm.(string)
		if !ok {
			return nil, fmt.Errorf("shutdown term is not text")
		}
		terms = append(terms, term)
	}
	return terms, nil
}
