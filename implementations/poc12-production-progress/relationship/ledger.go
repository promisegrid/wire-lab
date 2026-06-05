package relationship

import "sort"

// Outcome names the local interpretation of a peer's promise evidence. These
// are not protocol actions; they are local ledger conclusions about observed
// promise keep/break history.
type Outcome string
type Transition string

const (
	OutcomeKept          Outcome = "kept"
	OutcomeBroken        Outcome = "broken"
	OutcomeMalformed     Outcome = "malformed"
	OutcomeNonCommitment Outcome = "non_commitment"
	OutcomeRepairKept    Outcome = "repair_kept"
	OutcomeDiscoveryKept Outcome = "discovery_kept"

	TransitionAdded     Transition = "direct_peer_added"
	TransitionRemoved   Transition = "direct_peer_removed"
	TransitionUnchanged Transition = "direct_peer_unchanged"
)

// Ledger stores one agent's private trust and direct-link choices.
// Intent: POC12 correlates strong local trust with direct TCP adjacency and
// weak trust with removed adjacency without creating a global peering
// authority. Source: DI-timah
type Ledger struct {
	trustByPeer map[string]int
	directPeers map[string]bool
	strongTrust int
	weakTrust   int
	decay       int
}

// State is the durable, local-only relationship snapshot for one agent.
// Intent: Persist only this agent's private trust and current direct-link
// promises, never a global trust authority. Source: DI-timah
type State struct {
	TrustByPeer map[string]int `json:"trust_by_peer"`
	DirectPeers []string       `json:"direct_peers"`
}

// NewLedger initializes one local relationship ledger.
func NewLedger(allPeers []string, initialDirectPeers []string, strongTrust, weakTrust, decay int) *Ledger {
	trustByPeer := make(map[string]int, len(allPeers))
	for _, peerName := range allPeers {
		trustByPeer[peerName] = 0
	}
	directPeers := make(map[string]bool, len(initialDirectPeers))
	for _, peerName := range initialDirectPeers {
		if _, exists := trustByPeer[peerName]; exists {
			directPeers[peerName] = true
		}
	}
	return &Ledger{
		trustByPeer: trustByPeer,
		directPeers: directPeers,
		strongTrust: strongTrust,
		weakTrust:   weakTrust,
		decay:       decay,
	}
}

// Trust returns the local trust score for a peer.
func (ledger *Ledger) Trust(peerName string) int {
	return ledger.trustByPeer[peerName]
}

// Snapshot returns a copy of all local trust scores.
func (ledger *Ledger) Snapshot() map[string]int {
	snapshot := make(map[string]int, len(ledger.trustByPeer))
	for peerName, trustScore := range ledger.trustByPeer {
		snapshot[peerName] = trustScore
	}
	return snapshot
}

// DirectPeers returns currently promised direct TCP peers.
func (ledger *Ledger) DirectPeers() []string {
	peers := make([]string, 0, len(ledger.directPeers))
	for peerName, direct := range ledger.directPeers {
		if direct {
			peers = append(peers, peerName)
		}
	}
	sort.Strings(peers)
	return peers
}

// CanDial reports whether the local agent currently promises to dial a peer.
func (ledger *Ledger) CanDial(peerName string) bool {
	return ledger.directPeers[peerName]
}

// CanAccept reports whether the local agent currently promises to accept a peer.
func (ledger *Ledger) CanAccept(peerName string) bool {
	return ledger.directPeers[peerName]
}

// ObserveOutcome updates local trust from one observed promise outcome and
// reports how that evidence changed this agent's direct TCP relationship.
func (ledger *Ledger) ObserveOutcome(peerName string, outcome Outcome) Transition {
	if _, exists := ledger.trustByPeer[peerName]; !exists {
		return TransitionUnchanged
	}
	wasDirect := ledger.directPeers[peerName]
	switch outcome {
	case OutcomeKept:
		ledger.trustByPeer[peerName]++
	case OutcomeRepairKept:
		ledger.trustByPeer[peerName] += 2
	case OutcomeDiscoveryKept:
		// Intent: One kept low-risk discovery promise can form a direct peer
		// because the strong-trust threshold is deliberately small in POC12.
		// Source: DI-timah
		ledger.trustByPeer[peerName] += 2
	case OutcomeBroken, OutcomeMalformed:
		ledger.trustByPeer[peerName] -= 3
	case OutcomeNonCommitment:
		// Intent: A local non-commitment means the peer did not promise the
		// requested exchange; it is not evidence that the peer broke an explicit
		// promise. Source: DI-jinoz
	}
	ledger.reconfigure(peerName)
	isDirect := ledger.directPeers[peerName]
	if !wasDirect && isDirect {
		return TransitionAdded
	}
	if wasDirect && !isDirect {
		return TransitionRemoved
	}
	return TransitionUnchanged
}

// Export returns a durable copy of the local relationship state.
// Intent: Make restart tests inspectable without exposing mutable ledger maps.
// Source: DI-timah
func (ledger *Ledger) Export() State {
	return State{
		TrustByPeer: ledger.Snapshot(),
		DirectPeers: ledger.DirectPeers(),
	}
}

// ApplyState restores locally known peer state without accepting unknown peers.
// Intent: A persisted file may restore only relationships this configuration
// already recognizes, preventing stale files from inventing peers.
// Source: DI-timah
func (ledger *Ledger) ApplyState(state State) {
	for peerName, trustScore := range state.TrustByPeer {
		if _, exists := ledger.trustByPeer[peerName]; exists {
			ledger.trustByPeer[peerName] = trustScore
		}
	}
	ledger.directPeers = make(map[string]bool, len(state.DirectPeers))
	for _, peerName := range state.DirectPeers {
		if _, exists := ledger.trustByPeer[peerName]; exists {
			ledger.directPeers[peerName] = true
		}
	}
	for peerName := range ledger.trustByPeer {
		ledger.reconfigure(peerName)
	}
}

// DecayRound slowly moves idle relationships toward lower confidence.
func (ledger *Ledger) DecayRound() {
	if ledger.decay <= 0 {
		return
	}
	for peerName, trustScore := range ledger.trustByPeer {
		if trustScore > 0 {
			ledger.trustByPeer[peerName] -= ledger.decay
			if ledger.trustByPeer[peerName] < 0 {
				ledger.trustByPeer[peerName] = 0
			}
		}
		ledger.reconfigure(peerName)
	}
}

func (ledger *Ledger) reconfigure(peerName string) {
	trustScore := ledger.trustByPeer[peerName]
	if trustScore >= ledger.strongTrust {
		ledger.directPeers[peerName] = true
		return
	}
	if trustScore <= ledger.weakTrust {
		delete(ledger.directPeers, peerName)
	}
}
