package protocol

import (
	"testing"
	"time"
)

func TestLifecycleTokenCOSEVerification(t *testing.T) {
	protocolCID := MustProtocolCIDText("bafkreidamxalqxl2depjwlzhwdvglpda57fkqy5hvnwiz6jow6tapungeu")
	now := time.Unix(2000, 0)
	terms := LifecycleTokenTerms{
		IssuerRoleID:   "agent:alice",
		AudienceRoleID: "supervisor:alpha",
		RunID:          "test-run",
		RoleKind:       "app",
		ChannelProfile: LifecycleChannelTCP,
		ProtocolCID:    protocolCID,
		GraceMillis:    5000,
		MaxInvocations: 1,
		IssuedAtUnix:   now.Unix(),
		NotBeforeUnix:  now.Add(-time.Second).Unix(),
		ExpiresUnix:    now.Add(time.Minute).Unix(),
		TokenID:        "token-1",
		ShutdownTerms:  []string{"quiesce", "drain", "exit"},
	}
	token, issueErr := IssueLifecycleToken(terms)
	if issueErr != nil {
		t.Fatalf("issue lifecycle token: %v", issueErr)
	}
	verified, verifyErr := VerifyLifecycleToken(token.COSEBytes, terms, now, 0)
	if verifyErr != nil {
		t.Fatalf("verify lifecycle token: %v", verifyErr)
	}
	if verified.IssuerRoleID != terms.IssuerRoleID || verified.ProtocolCID.String() != protocolCID.String() {
		t.Fatalf("verified terms mismatch: %#v", verified)
	}
	if _, replayErr := VerifyLifecycleToken(token.COSEBytes, terms, now, 1); replayErr == nil {
		t.Fatalf("expected replay rejection")
	}
	wrongAudience := terms
	wrongAudience.AudienceRoleID = "supervisor:wrong"
	if _, audienceErr := VerifyLifecycleToken(token.COSEBytes, wrongAudience, now, 0); audienceErr == nil {
		t.Fatalf("expected wrong audience rejection")
	}
	if _, expiryErr := VerifyLifecycleToken(token.COSEBytes, terms, now.Add(time.Hour), 0); expiryErr == nil {
		t.Fatalf("expected expired token rejection")
	}
}

func TestLifecycleMessageUsesSinglePayloadSlot(t *testing.T) {
	protocolCID := MustProtocolCIDText("bafkreidamxalqxl2depjwlzhwdvglpda57fkqy5hvnwiz6jow6tapungeu")
	payload := LifecyclePayload{
		Kind:           LifecycleKindTokenFulfilled,
		Promiser:       "agent:alice",
		Promisee:       "supervisor:alpha",
		RoleID:         "agent:alice",
		RoleKind:       "app",
		ChannelProfile: LifecycleChannelTCP,
		RunID:          "test-run",
		TokenCID:       "bafkreibuvp6v3kqi6wdyrysfdppwr4vvgipbyp672t6eentkpj5swosu3y",
		Outcome:        LifecycleOutcomeKept,
		Detail:         "done",
	}
	frameBytes, encodeErr := EncodeLifecycleMessage(protocolCID, payload)
	if encodeErr != nil {
		t.Fatalf("encode lifecycle message: %v", encodeErr)
	}
	gridMessage, gridErr := ParseGridMessage(frameBytes)
	if gridErr != nil {
		t.Fatalf("parse grid message: %v", gridErr)
	}
	if len(gridMessage.Slots) != 1 {
		t.Fatalf("lifecycle message slots=%d, want one payload slot", len(gridMessage.Slots))
	}
	decoded, decodeErr := DecodeLifecycleMessage(frameBytes, protocolCID)
	if decodeErr != nil {
		t.Fatalf("decode lifecycle message: %v", decodeErr)
	}
	if decoded.Kind != payload.Kind || decoded.Outcome != payload.Outcome {
		t.Fatalf("decoded payload mismatch: %#v", decoded)
	}
}
