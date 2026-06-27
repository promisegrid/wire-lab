package protocol

import (
	"testing"
	"time"
)

func TestLifecycleTokenCOSEVerification(t *testing.T) {
	now := time.Unix(2000, 0)
	terms := testLifecycleTerms(now)
	token, err := IssueLifecycleToken(terms)
	if err != nil {
		t.Fatalf("issue lifecycle token: %v", err)
	}
	if err := ValidateCIDText(token.CID); err != nil {
		t.Fatalf("token CID is not canonical: %v", err)
	}
	verified, err := VerifyLifecycleToken(token.COSEBytes, terms, now, 0)
	if err != nil {
		t.Fatalf("verify lifecycle token: %v", err)
	}
	if verified.IssuerRoleID != terms.IssuerRoleID || verified.ProtocolPCID != LocalLifecycleV1PCID {
		t.Fatalf("verified terms mismatch: %#v", verified)
	}
	if _, err := VerifyLifecycleToken(token.COSEBytes, terms, now, 1); err == nil {
		t.Fatalf("expected replay rejection")
	}
	wrongRun := terms
	wrongRun.RunID = "wrong-run"
	if _, err := VerifyLifecycleToken(token.COSEBytes, wrongRun, now, 0); err == nil {
		t.Fatalf("expected wrong-run rejection")
	}
	wrongPCID := terms
	wrongPCID.ProtocolPCID = MustPCIDForName(ProtocolDeviceStatus)
	if _, err := VerifyLifecycleToken(token.COSEBytes, wrongPCID, now, 0); err == nil {
		t.Fatalf("expected wrong-pCID rejection")
	}
	if _, err := VerifyLifecycleToken(token.COSEBytes, terms, now.Add(time.Hour), 0); err == nil {
		t.Fatalf("expected expired token rejection")
	}
	if _, err := VerifyLifecycleToken([]byte("not cose"), terms, now, 0); err == nil {
		t.Fatalf("expected malformed token rejection")
	}
}

func TestLifecycleMessageUsesMapPayload(t *testing.T) {
	now := time.Unix(2000, 0)
	token, err := IssueLifecycleToken(testLifecycleTerms(now))
	if err != nil {
		t.Fatalf("issue lifecycle token: %v", err)
	}
	frame, err := EncodeLifecycleMessage(NewLifecycleTokenIssuedPayload(token))
	if err != nil {
		t.Fatalf("encode lifecycle frame: %v", err)
	}
	decoded, err := DecodeLifecycleMessage(frame)
	if err != nil {
		t.Fatalf("decode lifecycle frame: %v", err)
	}
	if decoded.Kind != LifecycleKindTokenIssued || decoded.TokenCID != token.CID {
		t.Fatalf("decoded payload mismatch: %#v", decoded)
	}
	tokenBytes, err := decoded.TokenBytes()
	if err != nil {
		t.Fatalf("decode token bytes: %v", err)
	}
	if _, err := VerifyLifecycleToken(tokenBytes, token.Terms, now, 0); err != nil {
		t.Fatalf("verify decoded token: %v", err)
	}
}

func testLifecycleTerms(now time.Time) LifecycleTokenTerms {
	return LifecycleTokenTerms{
		IssuerRoleID:   "agent:m4-ivan",
		AudienceRoleID: "supervisor:harness",
		RunID:          "test-run",
		RoleKind:       "app",
		ChannelProfile: LifecycleChannelStdio,
		ProtocolPCID:   LocalLifecycleV1PCID,
		GraceMillis:    5000,
		MaxInvocations: 1,
		IssuedAtUnix:   now.Unix(),
		NotBeforeUnix:  now.Add(-time.Second).Unix(),
		ExpiresUnix:    now.Add(time.Minute).Unix(),
		TokenID:        "token-1",
		ShutdownTerms:  []string{"quiesce", "drain accepted work", "flush local events", "exit voluntarily"},
	}
}
