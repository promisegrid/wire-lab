package protocol

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fxamacker/cbor/v2"
	cidlib "github.com/ipfs/go-cid"
	cose "github.com/veraison/go-cose"
)

const (
	LocalLifecycleV1PCID = "bafkreidamxalqxl2depjwlzhwdvglpda57fkqy5hvnwiz6jow6tapungeu"

	LifecycleKindTokenIssued    = "token_issued"
	LifecycleKindReady          = "ready"
	LifecycleKindTokenInvoked   = "token_invoked"
	LifecycleKindTokenFulfilled = "token_fulfilled"

	LifecycleChannelStdio = "stdio"

	LifecycleOutcomeKept     = "kept"
	LifecycleOutcomeRejected = "rejected"
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

// LifecycleTokenTerms are the CWT claims a local role signs for shutdown.
//
// Intent: POC17 inherits POC16 lifecycle/resource promises locally without
// sending large CWT/COSE tokens over the constrained LoRa path. Source: DI-zopub
type LifecycleTokenTerms struct {
	IssuerRoleID   string
	AudienceRoleID string
	RunID          string
	RoleKind       string
	ChannelProfile string
	ProtocolPCID   string
	GraceMillis    int64
	MaxInvocations int64
	IssuedAtUnix   int64
	NotBeforeUnix  int64
	ExpiresUnix    int64
	TokenID        string
	ShutdownTerms  []string
}

// LifecycleToken carries exact COSE/CWT bytes plus printable diagnostics.
type LifecycleToken struct {
	Terms        LifecycleTokenTerms
	COSEBytes    []byte
	CID          string
	PublicKey    []byte
	COSEBase64   string
	PublicBase64 string
}

// LifecyclePayload is the host-local local_lifecycle_v1 map payload.
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

// IssueLifecycleToken signs a local role's lifecycle promise using COSE_Sign1.
func IssueLifecycleToken(terms LifecycleTokenTerms) (LifecycleToken, error) {
	if err := validateLifecycleTerms(terms); err != nil {
		return LifecycleToken{}, err
	}
	payloadBytes, err := marshalLifecycleTokenTerms(terms)
	if err != nil {
		return LifecycleToken{}, err
	}
	signer, err := cose.NewSigner(cose.AlgorithmEdDSA, DeterministicPrivateKey(terms.IssuerRoleID))
	if err != nil {
		return LifecycleToken{}, err
	}
	headers := cose.Headers{
		Protected: cose.ProtectedHeader{cose.HeaderLabelAlgorithm: cose.AlgorithmEdDSA},
	}
	coseBytes, err := cose.Sign1(rand.Reader, signer, headers, payloadBytes, nil)
	if err != nil {
		return LifecycleToken{}, err
	}
	cid, err := CIDForBytes(coseBytes)
	if err != nil {
		return LifecycleToken{}, err
	}
	publicKey := DeterministicPublicKey(terms.IssuerRoleID)
	copiedCOSE := append([]byte(nil), coseBytes...)
	copiedPublic := append([]byte(nil), publicKey...)
	return LifecycleToken{
		Terms:        terms,
		COSEBytes:    copiedCOSE,
		CID:          cid,
		PublicKey:    copiedPublic,
		COSEBase64:   base64.StdEncoding.EncodeToString(copiedCOSE),
		PublicBase64: base64.StdEncoding.EncodeToString(copiedPublic),
	}, nil
}

// VerifyLifecycleToken checks the signed CWT terms before local invocation.
func VerifyLifecycleToken(coseBytes []byte, expected LifecycleTokenTerms, now time.Time, priorInvocations int64) (LifecycleTokenTerms, error) {
	var message cose.Sign1Message
	if err := message.UnmarshalCBOR(coseBytes); err != nil {
		return LifecycleTokenTerms{}, err
	}
	verifier, err := cose.NewVerifier(cose.AlgorithmEdDSA, DeterministicPublicKey(expected.IssuerRoleID))
	if err != nil {
		return LifecycleTokenTerms{}, err
	}
	if err := message.Verify(nil, verifier); err != nil {
		return LifecycleTokenTerms{}, err
	}
	terms, err := unmarshalLifecycleTokenTerms(message.Payload)
	if err != nil {
		return LifecycleTokenTerms{}, err
	}
	if err := compareLifecycleTerms(terms, expected, now, priorInvocations); err != nil {
		return LifecycleTokenTerms{}, err
	}
	return terms, nil
}

// NewLifecycleTokenIssuedPayload records a role's startup lifecycle promise.
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

// NewLifecycleReadyPayload records that a role is ready for token invocation.
func NewLifecycleReadyPayload(token LifecycleToken, detail string) LifecyclePayload {
	payload := NewLifecycleTokenIssuedPayload(token)
	payload.Kind = LifecycleKindReady
	payload.TokenCOSEBase64 = ""
	payload.PublicKeyBase64 = ""
	payload.Detail = detail
	return payload
}

// NewLifecycleInvocationPayload presents the exact token back to its issuer.
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

// NewLifecycleFulfilledPayload records the issuer's lifecycle outcome.
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

// EncodeLifecycleMessage writes grid([42(local_lifecycle_v1_pCID), map]).
func EncodeLifecycleMessage(payload LifecyclePayload) ([]byte, error) {
	tag42Data, err := Tag42DataForCIDText(LocalLifecycleV1PCID)
	if err != nil {
		return nil, err
	}
	outer := cbor.Tag{
		Number:  GridTag,
		Content: []any{cbor.Tag{Number: CIDTag, Content: tag42Data}, payload.lifecycleFields()},
	}
	return lifecycleCBOREncMode.Marshal(outer)
}

// DecodeLifecycleMessage parses a host-local local_lifecycle_v1 frame.
func DecodeLifecycleMessage(frameBytes []byte) (LifecyclePayload, error) {
	var outer cbor.Tag
	if err := cbor.Unmarshal(frameBytes, &outer); err != nil {
		return LifecyclePayload{}, err
	}
	if outer.Number != GridTag {
		return LifecyclePayload{}, fmt.Errorf("missing grid tag")
	}
	slots, ok := outer.Content.([]any)
	if !ok || len(slots) != 2 {
		return LifecyclePayload{}, fmt.Errorf("local_lifecycle_v1 expects two grid slots")
	}
	pcidTag, ok := slots[0].(cbor.Tag)
	if !ok || pcidTag.Number != CIDTag {
		return LifecyclePayload{}, fmt.Errorf("slot 0 must be tag 42 CID bytes")
	}
	tag42Data, ok := pcidTag.Content.([]byte)
	if !ok {
		return LifecyclePayload{}, fmt.Errorf("slot 0 CID content must be bytes")
	}
	pcidText, err := CIDTextFromTag42Data(tag42Data)
	if err != nil {
		return LifecyclePayload{}, err
	}
	if pcidText != LocalLifecycleV1PCID {
		return LifecyclePayload{}, fmt.Errorf("lifecycle pCID=%s, want %s", pcidText, LocalLifecycleV1PCID)
	}
	fields, ok := slots[1].(map[any]any)
	if !ok {
		return LifecyclePayload{}, fmt.Errorf("slot 1 must be lifecycle map")
	}
	return lifecyclePayloadFromFields(fields)
}

// TokenBytes decodes exact COSE token bytes carried by a lifecycle payload.
func (payload LifecyclePayload) TokenBytes() ([]byte, error) {
	if strings.TrimSpace(payload.TokenCOSEBase64) == "" {
		return nil, fmt.Errorf("lifecycle payload has no token bytes")
	}
	tokenBytes, err := base64.StdEncoding.DecodeString(payload.TokenCOSEBase64)
	if err != nil {
		return nil, err
	}
	if payload.TokenCID != "" {
		cid, err := CIDForBytes(tokenBytes)
		if err != nil {
			return nil, err
		}
		if cid != payload.TokenCID {
			return nil, fmt.Errorf("lifecycle token CID mismatch")
		}
	}
	return tokenBytes, nil
}

// DeterministicPrivateKey is a POC-only key source for reproducible signatures.
func DeterministicPrivateKey(seedText string) ed25519.PrivateKey {
	seed := sha256.Sum256([]byte("poc17 lifecycle signer: " + seedText))
	return ed25519.NewKeyFromSeed(seed[:])
}

// DeterministicPublicKey returns the POC public key for a named local role.
func DeterministicPublicKey(seedText string) ed25519.PublicKey {
	return DeterministicPrivateKey(seedText).Public().(ed25519.PublicKey)
}

func mustLifecycleCBOREncMode() cbor.EncMode {
	mode, err := cbor.CanonicalEncOptions().EncMode()
	if err != nil {
		panic(err)
	}
	return mode
}

func marshalLifecycleTokenTerms(terms LifecycleTokenTerms) ([]byte, error) {
	pcidBytes, err := cidBytesForText(terms.ProtocolPCID)
	if err != nil {
		return nil, err
	}
	cwtTerms := map[any]any{
		cose.CWTClaimIssuer:         terms.IssuerRoleID,
		cose.CWTClaimSubject:        localLifecycleTokenSubject,
		cose.CWTClaimAudience:       terms.AudienceRoleID,
		cose.CWTClaimIssuedAt:       terms.IssuedAtUnix,
		cose.CWTClaimNotBefore:      terms.NotBeforeUnix,
		cose.CWTClaimExpirationTime: terms.ExpiresUnix,
		cose.CWTClaimCWTID:          []byte(terms.TokenID),
		cwtPrivateProtocolCID:       pcidBytes,
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
	protocolPCID, err := pcidTextTerm(cwtValue(cwtTerms, cwtPrivateProtocolCID))
	if err != nil {
		return LifecycleTokenTerms{}, err
	}
	shutdownTerms, err := stringSliceTerm(cwtValue(cwtTerms, cwtPrivateShutdownTerms))
	if err != nil {
		return LifecycleTokenTerms{}, err
	}
	tokenID, err := bytesTextTerm(cwtValue(cwtTerms, cose.CWTClaimCWTID))
	if err != nil {
		return LifecycleTokenTerms{}, err
	}
	terms := LifecycleTokenTerms{
		IssuerRoleID:   stringTerm(cwtValue(cwtTerms, cose.CWTClaimIssuer)),
		AudienceRoleID: stringTerm(cwtValue(cwtTerms, cose.CWTClaimAudience)),
		RunID:          stringTerm(cwtValue(cwtTerms, cwtPrivateRunID)),
		RoleKind:       stringTerm(cwtValue(cwtTerms, cwtPrivateRoleKind)),
		ChannelProfile: stringTerm(cwtValue(cwtTerms, cwtPrivateChannelProfile)),
		ProtocolPCID:   protocolPCID,
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
	if actual.ProtocolPCID != expected.ProtocolPCID {
		return fmt.Errorf("lifecycle token pCID=%s, want %s", actual.ProtocolPCID, expected.ProtocolPCID)
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
	if terms.ProtocolPCID == "" {
		return fmt.Errorf("lifecycle token needs protocol pCID")
	}
	if err := ValidateCIDText(terms.ProtocolPCID); err != nil {
		return fmt.Errorf("invalid lifecycle pCID: %w", err)
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

func (payload LifecyclePayload) lifecycleFields() map[string]any {
	fields := map[string]any{
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
		fields["deadline_unix"] = payload.DeadlineUnix
	}
	if payload.Outcome != "" {
		fields["outcome"] = payload.Outcome
	}
	if payload.Detail != "" {
		fields["detail"] = payload.Detail
	}
	return fields
}

func lifecyclePayloadFromFields(fields map[any]any) (LifecyclePayload, error) {
	payload := LifecyclePayload{
		Kind:            stringMapTerm(fields, "kind"),
		Promiser:        stringMapTerm(fields, "promiser"),
		Promisee:        stringMapTerm(fields, "promisee"),
		RoleID:          stringMapTerm(fields, "role_id"),
		RoleKind:        stringMapTerm(fields, "role_kind"),
		ChannelProfile:  stringMapTerm(fields, "channel_profile"),
		RunID:           stringMapTerm(fields, "run_id"),
		TokenCOSEBase64: stringMapTerm(fields, "token_cose_b64"),
		TokenCID:        stringMapTerm(fields, "token_cid"),
		PublicKeyBase64: stringMapTerm(fields, "public_key_b64"),
		Reason:          stringMapTerm(fields, "reason"),
		DeadlineUnix:    int64MapTerm(fields, "deadline_unix"),
		Outcome:         stringMapTerm(fields, "outcome"),
		Detail:          stringMapTerm(fields, "detail"),
	}
	if payload.Kind == "" || payload.Promiser == "" || payload.Promisee == "" || payload.RoleID == "" || payload.RoleKind == "" || payload.ChannelProfile == "" || payload.RunID == "" {
		return LifecyclePayload{}, fmt.Errorf("lifecycle payload missing common fields")
	}
	return payload, nil
}

func cidBytesForText(cidText string) ([]byte, error) {
	parsed, err := cidlib.Decode(cidText)
	if err != nil {
		return nil, err
	}
	return parsed.Bytes(), nil
}

func pcidTextTerm(value any) (string, error) {
	bytesValue, ok := value.([]byte)
	if !ok {
		return "", fmt.Errorf("expected protocol CID bytes")
	}
	parsed, err := cidlib.Cast(bytesValue)
	if err != nil {
		return "", err
	}
	text := parsed.String()
	if err := ValidateCIDText(text); err != nil {
		return "", err
	}
	return text, nil
}

func cwtValue(fields map[any]any, label any) any {
	value, ok := fields[label]
	if !ok {
		for candidate, candidateValue := range fields {
			if numericLabelsEqual(candidate, label) {
				return candidateValue
			}
		}
		return nil
	}
	return value
}

func numericLabelsEqual(left, right any) bool {
	leftValue, leftOK := labelInt64(left)
	rightValue, rightOK := labelInt64(right)
	return leftOK && rightOK && leftValue == rightValue
}

func labelInt64(value any) (int64, bool) {
	switch typed := value.(type) {
	case int:
		return int64(typed), true
	case int64:
		return typed, true
	case int32:
		return int64(typed), true
	case uint:
		return int64(typed), true
	case uint64:
		if typed > uint64(^uint64(0)>>1) {
			return 0, false
		}
		return int64(typed), true
	case uint32:
		return int64(typed), true
	default:
		return 0, false
	}
}

func stringTerm(value any) string {
	if value == nil {
		return ""
	}
	if text, ok := value.(string); ok {
		return text
	}
	return ""
}

func stringSliceTerm(value any) ([]string, error) {
	if value == nil {
		return nil, fmt.Errorf("expected string slice")
	}
	values, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("expected string slice")
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		text, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("expected string item")
		}
		result = append(result, text)
	}
	return result, nil
}

func bytesTextTerm(value any) (string, error) {
	switch typed := value.(type) {
	case []byte:
		return string(typed), nil
	case string:
		return typed, nil
	default:
		return "", fmt.Errorf("expected bytes text, got %T", value)
	}
}

func int64Term(value any) int64 {
	switch v := value.(type) {
	case int64:
		return v
	case uint64:
		return int64(v)
	case int:
		return int64(v)
	case uint:
		return int64(v)
	default:
		return 0
	}
}

func stringMapTerm(fields map[any]any, key string) string {
	value, ok := fields[key]
	if !ok {
		return ""
	}
	return stringTerm(value)
}

func int64MapTerm(fields map[any]any, key string) int64 {
	value, ok := fields[key]
	if !ok {
		return 0
	}
	if text, ok := value.(string); ok {
		parsed, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return 0
		}
		return parsed
	}
	return int64Term(value)
}
