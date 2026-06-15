package relationship

import "sort"

// Outcome names the local interpretation of a peer's promise event record. These
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

const recoveryCautionAfterNegativeEvent = 4
const minTrustScore = -5
const maxTrustScore = 5

// Ledger stores one agent's private trust and direct-link choices.
// Intent: POC15 correlates strong local trust with direct TCP adjacency and
// weak trust with removed adjacency without creating a global peering
// authority. Source: DI-timah
type Ledger struct {
	trustByPeer           map[string]int
	cautionByPeer         map[string]int
	directPeers           map[string]bool
	permanentDistrust     map[string]bool
	transitExcludedByPeer map[string]bool
	strongTrust           int
	weakTrust             int
	decay                 int
}

// State is the durable, local-only relationship snapshot for one agent.
// Intent: Persist only this agent's private trust and current direct-link
// promises, never a global trust authority. Recovery caution is persisted with
// trust so a process restart inside the same run cannot erase recent
// malformed/broken events. Permanent distrust and transit exclusion are also
// local restraint promises, not global reputation or route policy. Source:
// DI-timah; DI-fijov; DI-dubih
type State struct {
	TrustByPeer           map[string]int `json:"trust_by_peer"`
	CautionByPeer         map[string]int `json:"caution_by_peer,omitempty"`
	DirectPeers           []string       `json:"direct_peers"`
	PermanentDistrust     []string       `json:"permanent_distrust,omitempty"`
	TransitExcludedByPeer []string       `json:"transit_excluded_by_peer,omitempty"`
}

// NewLedger initializes one local relationship ledger.
func NewLedger(allPeers []string, initialDirectPeers []string, strongTrust, weakTrust, decay int) *Ledger {
	trustByPeer := make(map[string]int, len(allPeers))
	cautionByPeer := make(map[string]int, len(allPeers))
	for _, peerName := range allPeers {
		trustByPeer[peerName] = 0
		cautionByPeer[peerName] = 0
	}
	directPeers := make(map[string]bool, len(initialDirectPeers))
	for _, peerName := range initialDirectPeers {
		if _, exists := trustByPeer[peerName]; exists {
			directPeers[peerName] = true
		}
	}
	return &Ledger{
		trustByPeer:           trustByPeer,
		cautionByPeer:         cautionByPeer,
		directPeers:           directPeers,
		permanentDistrust:     make(map[string]bool, len(allPeers)),
		transitExcludedByPeer: make(map[string]bool, len(allPeers)),
		strongTrust:           strongTrust,
		weakTrust:             weakTrust,
		decay:                 decay,
	}
}

// Trust returns the local trust score for a peer.
func (ledger *Ledger) Trust(peerName string) int {
	return ledger.trustByPeer[peerName]
}

// Caution returns the remaining positive-event delay for one peer.
// Intent: Runtime and analyzer event records should be able to show that recent
// malformed/broken events delay recovery without exposing mutable ledger
// state or creating a global reputation authority. Source: DI-sihuz
func (ledger *Ledger) Caution(peerName string) int {
	return ledger.cautionByPeer[peerName]
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
	return sortedEnabledPeers(ledger.directPeers)
}

// PermanentDistrustPeers returns peers this local agent has decided not to
// repair or contact without a separate future local decision.
func (ledger *Ledger) PermanentDistrustPeers() []string {
	return sortedEnabledPeers(ledger.permanentDistrust)
}

// TransitExcludedPeers returns peers this local agent will not use as transit
// hops for this agent's own inbound or outbound traffic.
func (ledger *Ledger) TransitExcludedPeers() []string {
	return sortedEnabledPeers(ledger.transitExcludedByPeer)
}

// PermanentlyDistrusted reports whether this local agent has decided not to use
// ordinary future event records from the peer to restore contact automatically.
func (ledger *Ledger) PermanentlyDistrusted(peerName string) bool {
	return ledger.permanentDistrust[peerName]
}

// TransitExcluded reports whether this local agent refuses to use the peer as a
// transit hop for this agent's own traffic.
func (ledger *Ledger) TransitExcluded(peerName string) bool {
	return ledger.transitExcludedByPeer[peerName]
}

// CanDial reports whether the local agent currently promises to dial a peer.
func (ledger *Ledger) CanDial(peerName string) bool {
	if ledger.permanentDistrust[peerName] {
		return false
	}
	return ledger.directPeers[peerName]
}

// CanAccept reports whether the local agent currently promises to accept a peer.
func (ledger *Ledger) CanAccept(peerName string) bool {
	if ledger.permanentDistrust[peerName] {
		return false
	}
	return ledger.directPeers[peerName]
}

// PermanentlyDistrust records a local durable decision not to treat future
// ordinary contact or repair from this peer as enough to restore relationship.
// Intent: Permanent distrust is Alice's local restraint promise, not punishment
// imposed on Mallory or a global reputation fact. Source: DI-dubih
func (ledger *Ledger) PermanentlyDistrust(peerName string) {
	if _, exists := ledger.trustByPeer[peerName]; !exists {
		return
	}
	ledger.permanentDistrust[peerName] = true
	ledger.trustByPeer[peerName] = minTrustScore
	delete(ledger.directPeers, peerName)
}

// ExcludeTransit records a local route promise: this peer must not appear as a
// transit hop in this agent's own traffic candidates.
// Intent: Transit exclusion is local path selection, not network-wide route
// enforcement or permission. Source: DI-dubih
func (ledger *Ledger) ExcludeTransit(peerName string) {
	if _, exists := ledger.trustByPeer[peerName]; exists {
		ledger.transitExcludedByPeer[peerName] = true
	}
}

// RouteAllowed reports whether this local agent can use a candidate path without
// crossing a locally excluded transit peer.
func (ledger *Ledger) RouteAllowed(route []string) bool {
	for routeIndex, peerName := range route {
		if routeIndex == 0 || routeIndex == len(route)-1 {
			continue
		}
		if ledger.transitExcludedByPeer[peerName] || ledger.permanentDistrust[peerName] {
			return false
		}
	}
	return true
}

// ObserveOutcome updates local trust from one observed promise outcome and
// reports how that event record changed this agent's direct TCP relationship.
func (ledger *Ledger) ObserveOutcome(peerName string, outcome Outcome) Transition {
	if _, exists := ledger.trustByPeer[peerName]; !exists {
		return TransitionUnchanged
	}
	if ledger.permanentDistrust[peerName] {
		delete(ledger.directPeers, peerName)
		return TransitionUnchanged
	}
	wasDirect := ledger.directPeers[peerName]
	switch outcome {
	case OutcomeKept:
		if !ledger.consumeRecoveryCaution(peerName) {
			ledger.trustByPeer[peerName]++
		}
	case OutcomeRepairKept:
		// Intent: Explicitly kept repair promises may rebuild local confidence,
		// but the separate caution counter still records that recent negative
		// event existed. Source: DI-fijov
		ledger.consumeRecoveryCaution(peerName)
		ledger.trustByPeer[peerName] += 2
	case OutcomeDiscoveryKept:
		// Intent: One kept low-risk discovery promise can form a direct peer
		// because the strong-trust threshold is deliberately small in POC15.
		// If this peer has recent malformed/broken events, discovery first
		// works off local caution instead of immediately forming adjacency.
		// Source: DI-timah; DI-fijov
		if !ledger.consumeRecoveryCaution(peerName) {
			ledger.trustByPeer[peerName] += 2
		}
	case OutcomeBroken, OutcomeMalformed:
		ledger.trustByPeer[peerName] -= 3
		ledger.addRecoveryCaution(peerName)
	case OutcomeNonCommitment:
		// Intent: A local non-commitment means the peer did not promise the
		// requested exchange; it is not an event showing that the peer broke an explicit
		// promise. Source: DI-jinoz
	}
	ledger.saturateTrust(peerName)
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

// consumeRecoveryCaution makes ordinary positive events pay down recent
// malformed/broken events before they can raise trust again.
// Intent: Local trust should recover through a sequence of observed kept
// promises, not by immediately accepting several neat-looking messages after a
// malformed or broken promise. Source: DI-fijov
func (ledger *Ledger) consumeRecoveryCaution(peerName string) bool {
	if ledger.cautionByPeer[peerName] <= 0 {
		return false
	}
	ledger.cautionByPeer[peerName]--
	return true
}

// addRecoveryCaution records that this peer now needs extra locally observed
// kept behavior before positive trust can climb again.
// Intent: The caution counter is local relationship memory, not a punishment
// imposed by another agent or a global reputation system. Source: DI-fijov
func (ledger *Ledger) addRecoveryCaution(peerName string) {
	if ledger.cautionByPeer[peerName] < recoveryCautionAfterNegativeEvent {
		ledger.cautionByPeer[peerName] = recoveryCautionAfterNegativeEvent
	}
}

// saturateTrust keeps the local trust score in a small comparable range.
// Intent: POC15 monitor output should not show unbounded local trust values
// that look like absolute reputation; trust remains only this agent's bounded
// relationship event records for this run. Source: DI-sihuz
func (ledger *Ledger) saturateTrust(peerName string) {
	if ledger.trustByPeer[peerName] > maxTrustScore {
		ledger.trustByPeer[peerName] = maxTrustScore
	}
	if ledger.trustByPeer[peerName] < minTrustScore {
		ledger.trustByPeer[peerName] = minTrustScore
	}
}

// Export returns a durable copy of the local relationship state.
// Intent: Make restart tests inspectable without exposing mutable ledger maps.
// Source: DI-timah
func (ledger *Ledger) Export() State {
	cautionSnapshot := make(map[string]int, len(ledger.cautionByPeer))
	for peerName, cautionCount := range ledger.cautionByPeer {
		if cautionCount > 0 {
			cautionSnapshot[peerName] = cautionCount
		}
	}
	return State{
		TrustByPeer:           ledger.Snapshot(),
		CautionByPeer:         cautionSnapshot,
		DirectPeers:           ledger.DirectPeers(),
		PermanentDistrust:     ledger.PermanentDistrustPeers(),
		TransitExcludedByPeer: ledger.TransitExcludedPeers(),
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
			ledger.saturateTrust(peerName)
		}
	}
	for peerName, cautionCount := range state.CautionByPeer {
		if _, exists := ledger.cautionByPeer[peerName]; exists && cautionCount > 0 {
			ledger.cautionByPeer[peerName] = cautionCount
		}
	}
	ledger.directPeers = make(map[string]bool, len(state.DirectPeers))
	for _, peerName := range state.DirectPeers {
		if _, exists := ledger.trustByPeer[peerName]; exists {
			ledger.directPeers[peerName] = true
		}
	}
	ledger.permanentDistrust = make(map[string]bool, len(state.PermanentDistrust))
	for _, peerName := range state.PermanentDistrust {
		if _, exists := ledger.trustByPeer[peerName]; exists {
			ledger.permanentDistrust[peerName] = true
			delete(ledger.directPeers, peerName)
		}
	}
	ledger.transitExcludedByPeer = make(map[string]bool, len(state.TransitExcludedByPeer))
	for _, peerName := range state.TransitExcludedByPeer {
		if _, exists := ledger.trustByPeer[peerName]; exists {
			ledger.transitExcludedByPeer[peerName] = true
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
	if ledger.permanentDistrust[peerName] {
		delete(ledger.directPeers, peerName)
		return
	}
	trustScore := ledger.trustByPeer[peerName]
	if trustScore >= ledger.strongTrust {
		ledger.directPeers[peerName] = true
		return
	}
	if trustScore <= ledger.weakTrust {
		delete(ledger.directPeers, peerName)
	}
}

func sortedEnabledPeers(peerMap map[string]bool) []string {
	peers := make([]string, 0, len(peerMap))
	for peerName, enabled := range peerMap {
		if enabled {
			peers = append(peers, peerName)
		}
	}
	sort.Strings(peers)
	return peers
}
