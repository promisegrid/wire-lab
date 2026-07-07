package sync

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	cidlib "github.com/ipfs/go-cid"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/economy"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/graph"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
)

// AgentState is the local sync-agent checkpoint. It is a performance cache and
// local memory aid, not an authority over peers or branch state.
type AgentState struct {
	Version                 int                        `json:"version"`
	Agent                   string                     `json:"agent"`
	Trust                   map[string]TrustCheckpoint `json:"trust"`
	RedeemedBearerTokenCIDs []string                   `json:"redeemed_bearer_token_cids,omitempty"`
	LastReport              *SchedulerReport           `json:"last_report,omitempty"`
}

// TrustCheckpoint records one local trust estimate derived from retained graph
// objects, optionally falling back to a market signal when local evidence is
// absent.
type TrustCheckpoint struct {
	Peer             string `json:"peer"`
	Score            int    `json:"score"`
	KeptPromises     int    `json:"kept_promises"`
	BrokenPromises   int    `json:"broken_promises"`
	MarketSignal     int    `json:"market_signal,omitempty"`
	MarketSignalUsed bool   `json:"market_signal_used"`
	Reason           string `json:"reason"`
}

// SchedulerConfig controls one local sync-agent scheduling pass.
type SchedulerConfig struct {
	Rounds      int       `json:"rounds"`
	RetainUntil string    `json:"retain_until"`
	Now         time.Time `json:"-"`
}

// CandidatePeer describes one peer the local agent may choose for a scheduled
// sync pass.
type CandidatePeer struct {
	Peer             Peer
	Heads            map[string]cidlib.Cid
	MarketSignal     int
	BearerTokenBytes []byte
	BearerIssuer     string
	PaymentScope     string
	PaymentObjectCID string
	PaymentValue     int64
	PaymentUnit      string
	Capability       string
}

// SchedulerReport records one local scheduler pass.
type SchedulerReport struct {
	Agent                 string                               `json:"agent"`
	ChosenPeer            string                               `json:"chosen_peer,omitempty"`
	Decisions             []PeerDecision                       `json:"decisions"`
	CapabilityRedemptions []economy.CapabilityRedemptionReport `json:"capability_redemptions,omitempty"`
	ContinuousSync        ContinuousSyncReport                 `json:"continuous_sync"`
	IdleReason            string                               `json:"idle_reason,omitempty"`
}

// PeerDecision records why one peer was chosen or skipped.
type PeerDecision struct {
	Peer   string          `json:"peer"`
	Result string          `json:"result"`
	Reason string          `json:"reason"`
	Trust  TrustCheckpoint `json:"trust"`
}

// DefaultAgentState returns an initialized local sync-agent checkpoint.
func DefaultAgentState(agent string) AgentState {
	return AgentState{Version: 1, Agent: agent, Trust: map[string]TrustCheckpoint{}}
}

// LoadAgentState reads local sync-agent checkpoint state, returning a default
// state when the path does not exist.
func LoadAgentState(path string, agent string) (AgentState, error) {
	content, readErr := os.ReadFile(path)
	if os.IsNotExist(readErr) {
		return DefaultAgentState(agent), nil
	}
	if readErr != nil {
		return AgentState{}, readErr
	}
	var state AgentState
	if unmarshalErr := json.Unmarshal(content, &state); unmarshalErr != nil {
		return AgentState{}, unmarshalErr
	}
	if state.Version == 0 {
		state.Version = 1
	}
	if state.Agent == "" {
		state.Agent = agent
	}
	if state.Trust == nil {
		state.Trust = map[string]TrustCheckpoint{}
	}
	if validateErr := validateAgentState(state); validateErr != nil {
		return AgentState{}, validateErr
	}
	return state, nil
}

// SaveAgentState writes local sync-agent checkpoint state atomically.
func SaveAgentState(path string, state AgentState) error {
	if state.Version == 0 {
		state.Version = 1
	}
	if state.Trust == nil {
		state.Trust = map[string]TrustCheckpoint{}
	}
	if validateErr := validateAgentState(state); validateErr != nil {
		return validateErr
	}
	content, marshalErr := json.MarshalIndent(state, "", "  ")
	if marshalErr != nil {
		return marshalErr
	}
	content = append(content, '\n')
	if mkdirErr := os.MkdirAll(filepath.Dir(path), 0o755); mkdirErr != nil {
		return mkdirErr
	}
	tmpPath := path + ".tmp"
	if writeErr := os.WriteFile(tmpPath, content, 0o644); writeErr != nil {
		return writeErr
	}
	return os.Rename(tmpPath, path)
}

// RunScheduledSync chooses the best peer from local evidence, redeems any
// offered bearer token into a peer-specific capability token, and runs one
// continuous-sync pass with the chosen peer.
//
// Intent: This is the gpg-agent-like scheduler core for `grid sync ...`: local
// state is created only when sync commands ask for it, peer choice comes from
// local retained promises, and token economics stay capability-specific. Source:
// DI-fakop
func RunScheduledSync(local Peer, candidates []CandidatePeer, state AgentState, config SchedulerConfig) (AgentState, SchedulerReport, error) {
	if local.Agent == "" || local.CAS == nil {
		return AgentState{}, SchedulerReport{}, fmt.Errorf("local peer name and CAS are required")
	}
	if state.Version == 0 {
		state = DefaultAgentState(local.Agent)
	}
	if state.Agent == "" {
		state.Agent = local.Agent
	}
	if state.Trust == nil {
		state.Trust = map[string]TrustCheckpoint{}
	}
	now := config.Now
	if now.IsZero() {
		now = time.Date(2026, 7, 4, 12, 0, 0, 0, time.UTC)
	}
	report := SchedulerReport{Agent: local.Agent}
	if len(candidates) == 0 {
		report.IdleReason = "no candidate peers"
		state.LastReport = &report
		return state, report, nil
	}
	scored, scoreErr := scoreCandidates(local.CAS, candidates)
	if scoreErr != nil {
		return AgentState{}, SchedulerReport{}, scoreErr
	}
	for _, scoredCandidate := range scored {
		state.Trust[scoredCandidate.candidate.Peer.Agent] = scoredCandidate.trust
	}
	sort.Slice(scored, func(left, right int) bool {
		if scored[left].trust.Score == scored[right].trust.Score {
			return scored[left].candidate.Peer.Agent < scored[right].candidate.Peer.Agent
		}
		return scored[left].trust.Score > scored[right].trust.Score
	})
	selected := scored[0]
	if selected.trust.Score <= 0 || len(selected.candidate.Heads) == 0 {
		for _, scoredCandidate := range scored {
			report.Decisions = append(report.Decisions, PeerDecision{
				Peer:   scoredCandidate.candidate.Peer.Agent,
				Result: "skipped",
				Reason: "no positive local trust score or no advertised heads",
				Trust:  scoredCandidate.trust,
			})
		}
		report.IdleReason = "no trusted candidate with advertised heads"
		state.LastReport = &report
		return state, report, nil
	}
	for _, scoredCandidate := range scored {
		result := "skipped"
		reason := "lower local trust score"
		if scoredCandidate.candidate.Peer.Agent == selected.candidate.Peer.Agent {
			result = "chosen"
			reason = selected.trust.Reason
		}
		report.Decisions = append(report.Decisions, PeerDecision{Peer: scoredCandidate.candidate.Peer.Agent, Result: result, Reason: reason, Trust: scoredCandidate.trust})
	}
	ledger, ledgerErr := economy.NewLedgerWithSpent(state.RedeemedBearerTokenCIDs)
	if ledgerErr != nil {
		return AgentState{}, SchedulerReport{}, ledgerErr
	}
	if len(selected.candidate.BearerTokenBytes) > 0 {
		redemption, redeemErr := redeemCandidateBearer(ledger, local, selected.candidate, now)
		if redeemErr != nil {
			return AgentState{}, SchedulerReport{}, redeemErr
		}
		report.CapabilityRedemptions = append(report.CapabilityRedemptions, redemption)
		state.RedeemedBearerTokenCIDs = appendUniqueCID(state.RedeemedBearerTokenCIDs, redemption.BearerTokenCID)
	}
	syncConfig := ContinuousSyncConfig{Rounds: config.Rounds, RetainUntil: config.RetainUntil, Offer: "scheduled_capability_sync"}
	continuous, syncErr := RunContinuousDAGSync(selected.candidate.Peer, local, selected.candidate.Heads, nil, syncConfig)
	if syncErr != nil {
		return AgentState{}, SchedulerReport{}, syncErr
	}
	report.ChosenPeer = selected.candidate.Peer.Agent
	report.ContinuousSync = continuous
	state.LastReport = &report
	return state, report, nil
}

type scoredCandidate struct {
	candidate CandidatePeer
	trust     TrustCheckpoint
}

func scoreCandidates(localCAS *store.FileStore, candidates []CandidatePeer) ([]scoredCandidate, error) {
	scored := make([]scoredCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate.Peer.Agent == "" || candidate.Peer.CAS == nil {
			return nil, fmt.Errorf("candidate peer name and CAS are required")
		}
		trust, trustErr := TrustFromLocalGraph(localCAS, candidate.Peer.Agent, candidate.MarketSignal)
		if trustErr != nil {
			return nil, trustErr
		}
		scored = append(scored, scoredCandidate{candidate: candidate, trust: trust})
	}
	return scored, nil
}

// TrustFromLocalGraph computes a local trust checkpoint from retained message
// promises. Market signal is used only when no local promise history exists.
func TrustFromLocalGraph(localCAS *store.FileStore, peer string, marketSignal int) (TrustCheckpoint, error) {
	if localCAS == nil {
		return TrustCheckpoint{}, fmt.Errorf("local CAS is required")
	}
	entries, listErr := localCAS.List()
	if listErr != nil {
		return TrustCheckpoint{}, listErr
	}
	checkpoint := TrustCheckpoint{Peer: peer, MarketSignal: marketSignal}
	for _, entry := range entries {
		if entry.Kind != "message" {
			continue
		}
		objectCID, parseErr := store.ParseCIDText(entry.CID)
		if parseErr != nil {
			return TrustCheckpoint{}, parseErr
		}
		content, _, getErr := localCAS.Get(objectCID)
		if getErr != nil {
			return TrustCheckpoint{}, getErr
		}
		view, parseEnvelopeErr := graph.ParseEnvelope(content)
		if parseEnvelopeErr != nil {
			continue
		}
		promiser := fmt.Sprint(view.Payload[0])
		if promiser != peer {
			continue
		}
		kind, kindErr := view.PayloadKind()
		if kindErr != nil {
			return TrustCheckpoint{}, kindErr
		}
		if promiseLooksBroken(kind, view.Payload[3]) {
			checkpoint.BrokenPromises++
			continue
		}
		checkpoint.KeptPromises++
		checkpoint.Score += trustWeight(kind)
	}
	if checkpoint.KeptPromises == 0 && checkpoint.BrokenPromises == 0 && marketSignal > 0 {
		checkpoint.Score = marketSignal
		checkpoint.MarketSignalUsed = true
		checkpoint.Reason = "used market signal because local graph had no promise history"
		return checkpoint, nil
	}
	checkpoint.Score -= checkpoint.BrokenPromises * 4
	checkpoint.Reason = "computed from retained local promise graph"
	return checkpoint, nil
}

func redeemCandidateBearer(ledger *economy.Ledger, local Peer, candidate CandidatePeer, now time.Time) (economy.CapabilityRedemptionReport, error) {
	if candidate.BearerIssuer == "" || candidate.PaymentScope == "" || candidate.PaymentObjectCID == "" || candidate.PaymentValue == 0 || candidate.PaymentUnit == "" || candidate.Capability == "" {
		return economy.CapabilityRedemptionReport{}, fmt.Errorf("candidate bearer redemption needs issuer, scope, object CID, value, unit, and capability")
	}
	redemption, redeemErr := ledger.RedeemBearerForCapability(candidate.Peer.CAS, candidate.BearerTokenBytes, economy.ExpectedBearerPayment{
		Issuer:     candidate.BearerIssuer,
		Scope:      candidate.PaymentScope,
		ObjectCID:  candidate.PaymentObjectCID,
		Value:      candidate.PaymentValue,
		Unit:       candidate.PaymentUnit,
		Capability: candidate.Capability,
	}, candidate.Peer.Agent, local.Agent, now)
	if redeemErr != nil {
		return economy.CapabilityRedemptionReport{}, redeemErr
	}
	if len(redemption.CapabilityBytes) > 0 {
		if _, putErr := local.CAS.Put("capability_token", redemption.CapabilityBytes); putErr != nil {
			return economy.CapabilityRedemptionReport{}, putErr
		}
	}
	return redemption, nil
}

func validateAgentState(state AgentState) error {
	if state.Version != 1 {
		return fmt.Errorf("unsupported sync-agent state version %d", state.Version)
	}
	if state.Agent == "" {
		return fmt.Errorf("sync-agent state agent is required")
	}
	for peer, trust := range state.Trust {
		if peer == "" || trust.Peer == "" {
			return fmt.Errorf("sync-agent trust peer is required")
		}
	}
	for _, tokenCID := range state.RedeemedBearerTokenCIDs {
		if _, parseErr := store.ParseCIDText(tokenCID); parseErr != nil {
			return parseErr
		}
	}
	return nil
}

func promiseLooksBroken(kind string, body any) bool {
	if strings.Contains(kind, "broken") {
		return true
	}
	return strings.Contains(fmt.Sprint(body), "broken")
}

func trustWeight(kind string) int {
	switch kind {
	case "object_availability":
		return 3
	case "object_retention":
		return 2
	case "storage_payment_redemption":
		return 2
	case "reference_set", "snapshot", "review_statement":
		return 1
	default:
		return 1
	}
}

func appendUniqueCID(values []string, next string) []string {
	for _, value := range values {
		if value == next {
			return values
		}
	}
	return append(values, next)
}
