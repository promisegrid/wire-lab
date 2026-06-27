package sim

import (
	"fmt"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/artifact"
	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/protocol"
)

const supervisorRoleID = "supervisor:harness"

type lifecycleSupervisor struct {
	runID       string
	writer      *artifact.Writer
	now         time.Time
	tokens      []protocol.LifecycleToken
	invocations map[string]int64
}

// newLifecycleSupervisor creates the host-local lifecycle/resource role.
//
// Intent: POC17 inherits POC16 lifecycle/resource promises locally while
// keeping lifecycle token bytes off the 200-byte LoRa path. Source: DI-zopub
func newLifecycleSupervisor(runID string, writer *artifact.Writer) *lifecycleSupervisor {
	return &lifecycleSupervisor{
		runID:       runID,
		writer:      writer,
		now:         time.Unix(2000, 0),
		invocations: make(map[string]int64),
	}
}

func (s *lifecycleSupervisor) start(roleIDs ...string) error {
	for _, roleID := range roleIDs {
		token, err := protocol.IssueLifecycleToken(s.termsFor(roleID, fmt.Sprintf("%s-token", roleID)))
		if err != nil {
			return err
		}
		s.tokens = append(s.tokens, token)
		if err := s.recordLifecycle("lifecycle_token_issued", token.Terms.IssuerRoleID, protocol.NewLifecycleTokenIssuedPayload(token), "issued", map[string]any{
			"token_cid":       token.CID,
			"channel_profile": token.Terms.ChannelProfile,
		}); err != nil {
			return err
		}
		if err := s.recordLifecycle("lifecycle_token_ready", token.Terms.IssuerRoleID, protocol.NewLifecycleReadyPayload(token, "work accepted"), "ready", map[string]any{
			"token_cid": token.CID,
		}); err != nil {
			return err
		}
		if err := s.writer.WriteEvent(artifact.Event{
			Type:    "resource_access_promised",
			Actor:   supervisorRoleID,
			Peer:    token.Terms.IssuerRoleID,
			PCID:    protocol.LocalLifecycleV1PCID,
			Outcome: "promised",
			Details: map[string]any{
				"resources": []string{"process", "cpu", "ram", "radio", "flash", "energy"},
				"scope":     "host_local",
				"source":    "DI-zopub",
			},
		}); err != nil {
			return err
		}
	}
	return s.recordTokenRejections()
}

func (s *lifecycleSupervisor) finish() error {
	for _, token := range s.tokens {
		if err := s.invoke(token); err != nil {
			return err
		}
	}
	if len(s.tokens) == 0 {
		return fmt.Errorf("lifecycle supervisor has no tokens")
	}
	// Intent: Resource withdrawal is local promise protection, not command
	// authority over the role or global peer-trust evidence. Source: DI-zopub
	return s.writer.WriteEvent(artifact.Event{
		Type:    "resource_access_withdrawn",
		Actor:   supervisorRoleID,
		Peer:    s.tokens[0].Terms.IssuerRoleID,
		PCID:    protocol.LocalLifecycleV1PCID,
		CID:     s.tokens[0].CID,
		Outcome: "local_resource_withdrawn",
		Details: map[string]any{
			"broken_promise":          "late radio buffer release",
			"mechanism":               "radio_access_withdrawal",
			"not_command_authority":   true,
			"not_peer_trust_evidence": true,
			"scope":                   "host_local",
		},
	})
}

func (s *lifecycleSupervisor) invoke(token protocol.LifecycleToken) error {
	invocation := protocol.NewLifecycleInvocationPayload(supervisorRoleID, token, "run_complete", s.now.Add(5*time.Second))
	if err := s.recordLifecycle("lifecycle_token_invoked", supervisorRoleID, invocation, "invoked", map[string]any{
		"token_cid": token.CID,
		"reason":    "run_complete",
	}); err != nil {
		return err
	}
	tokenBytes, err := invocation.TokenBytes()
	if err != nil {
		return err
	}
	prior := s.invocations[token.CID]
	if _, err := protocol.VerifyLifecycleToken(tokenBytes, token.Terms, s.now, prior); err != nil {
		return s.writer.WriteEvent(artifact.Event{
			Type:    "lifecycle_token_rejected",
			Actor:   token.Terms.IssuerRoleID,
			Peer:    supervisorRoleID,
			PCID:    protocol.LocalLifecycleV1PCID,
			CID:     token.CID,
			Outcome: "rejected",
			Details: map[string]any{"reason": err.Error(), "token_cid": token.CID},
		})
	}
	s.invocations[token.CID] = prior + 1
	return s.recordLifecycle("lifecycle_token_fulfilled", token.Terms.IssuerRoleID, protocol.NewLifecycleFulfilledPayload(token, protocol.LifecycleOutcomeKept, "quiesced, drained, flushed, and exited voluntarily"), "kept", map[string]any{
		"token_cid": token.CID,
	})
}

func (s *lifecycleSupervisor) recordTokenRejections() error {
	if len(s.tokens) == 0 {
		return nil
	}
	token := s.tokens[0]
	checks := []struct {
		reason string
		fn     func() error
	}{
		{reason: "malformed_cose", fn: func() error {
			_, err := protocol.VerifyLifecycleToken([]byte("not cose"), token.Terms, s.now, 0)
			return err
		}},
		{reason: "wrong_run", fn: func() error {
			expected := token.Terms
			expected.RunID = "wrong-run"
			_, err := protocol.VerifyLifecycleToken(token.COSEBytes, expected, s.now, 0)
			return err
		}},
		{reason: "wrong_pcid", fn: func() error {
			expected := token.Terms
			expected.ProtocolPCID = protocol.MustPCIDForName(protocol.ProtocolDeviceStatus)
			_, err := protocol.VerifyLifecycleToken(token.COSEBytes, expected, s.now, 0)
			return err
		}},
		{reason: "expired", fn: func() error {
			_, err := protocol.VerifyLifecycleToken(token.COSEBytes, token.Terms, s.now.Add(time.Hour), 0)
			return err
		}},
		{reason: "replayed", fn: func() error {
			_, err := protocol.VerifyLifecycleToken(token.COSEBytes, token.Terms, s.now, 1)
			return err
		}},
	}
	for _, check := range checks {
		err := check.fn()
		if err == nil {
			return fmt.Errorf("expected lifecycle rejection for %s", check.reason)
		}
		if writeErr := s.writer.WriteEvent(artifact.Event{
			Type:    "lifecycle_token_rejected",
			Actor:   token.Terms.IssuerRoleID,
			Peer:    supervisorRoleID,
			PCID:    protocol.LocalLifecycleV1PCID,
			CID:     token.CID,
			Outcome: "rejected",
			Details: map[string]any{
				"reason":    check.reason,
				"error":     err.Error(),
				"token_cid": token.CID,
			},
		}); writeErr != nil {
			return writeErr
		}
	}
	return nil
}

func (s *lifecycleSupervisor) recordLifecycle(eventType, actor string, payload protocol.LifecyclePayload, outcome string, details map[string]any) error {
	frame, err := protocol.EncodeLifecycleMessage(payload)
	if err != nil {
		return err
	}
	frameCID, path, err := s.writer.RecordLifecycle(frame)
	if err != nil {
		return err
	}
	if details == nil {
		details = make(map[string]any)
	}
	details["frame_cid"] = frameCID
	details["host_local_only"] = true
	return s.writer.WriteEvent(artifact.Event{
		Type:    eventType,
		Actor:   actor,
		Peer:    payload.Promisee,
		PCID:    protocol.LocalLifecycleV1PCID,
		CID:     frameCID,
		Path:    path,
		Outcome: outcome,
		Details: details,
	})
}

func (s *lifecycleSupervisor) termsFor(roleID, tokenID string) protocol.LifecycleTokenTerms {
	return protocol.LifecycleTokenTerms{
		IssuerRoleID:   roleID,
		AudienceRoleID: supervisorRoleID,
		RunID:          s.runID,
		RoleKind:       "app",
		ChannelProfile: protocol.LifecycleChannelStdio,
		ProtocolPCID:   protocol.LocalLifecycleV1PCID,
		GraceMillis:    5000,
		MaxInvocations: 1,
		IssuedAtUnix:   s.now.Unix(),
		NotBeforeUnix:  s.now.Add(-time.Second).Unix(),
		ExpiresUnix:    s.now.Add(time.Minute).Unix(),
		TokenID:        tokenID,
		ShutdownTerms: []string{
			"quiesce",
			"drain accepted work",
			"close local sessions",
			"flush local events",
			"emit role summary",
			"exit voluntarily",
		},
	}
}
