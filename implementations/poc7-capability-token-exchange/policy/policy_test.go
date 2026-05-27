package policy

import (
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/token"
)

func TestDaveAcceptsFirstStaleTokenForEvidence(t *testing.T) {
	// Intent: Dave's first stale-token acceptance is voluntary evidence gathering,
	// not permission from Alice or any central authority. Source: DI-rodog
	davePolicy := ForNode("dave")
	decision := davePolicy.Decide(ActionContext{
		Action:       ActionAccept,
		Peer:         "mallory",
		SourcePeer:   "mallory",
		Issuer:       "alice",
		ResourceKind: "data",
		ResourceID:   "dataset-stale",
		TransferRule: token.TransferBearer,
	})
	if !decision.Accepted {
		t.Fatalf("dave should accept first stale token for local evidence: %#v", decision)
	}
}

func TestDaveRefusesMalloryAfterBrokenEvidence(t *testing.T) {
	// Intent: Broken redemption evidence must change Dave's later local choice so
	// Mallory cannot keep circulating stale tokens as if trust were unchanged.
	// Source: DI-rodog
	davePolicy := ForNode("dave")
	decision := davePolicy.Decide(ActionContext{
		Action:          ActionAccept,
		Peer:            "mallory",
		SourcePeer:      "mallory",
		Issuer:          "alice",
		ResourceKind:    "data",
		ResourceID:      "dataset-stale",
		TransferRule:    token.TransferBearer,
		IssuerTrust:     -1,
		SourcePeerTrust: -1,
		PeerTrust:       -1,
	})
	if decision.Accepted {
		t.Fatalf("dave should refuse later mallory stale token after broken evidence: %#v", decision)
	}
}

func TestMalloryChoosesToCirculateBearerToken(t *testing.T) {
	// Intent: Mallory still acts from Mallory-local incentives; the scenario does
	// not force the transfer even when the outcome is bad for Dave. Source: DI-rodog
	malloryPolicy := ForNode("mallory")
	decision := malloryPolicy.Decide(ActionContext{
		Action:         ActionTransfer,
		Recipient:      "dave",
		RecipientTrust: 0,
		Issuer:         "alice",
		ResourceKind:   "data",
		ResourceID:     "dataset-stale",
		TransferRule:   token.TransferBearer,
	})
	if !decision.Accepted {
		t.Fatalf("mallory should choose short-term stale-token circulation: %#v", decision)
	}
}

func TestReciprocalExchangeClearsLocalThreshold(t *testing.T) {
	// Intent: Reciprocal exchange should clear Alice's own local threshold when
	// Carol offers a bearer token Alice values enough. Source: DI-rodog
	alicePolicy := ForNode("alice")
	decision := alicePolicy.Decide(ActionContext{
		Action:       ActionTrade,
		Peer:         "carol",
		Issuer:       "bob",
		ResourceKind: "storage",
		ResourceID:   "storage-slot-trade",
		DesiredKind:  "data",
		DesiredID:    "dataset-private",
		TransferRule: token.TransferBearer,
	})
	if !decision.Accepted {
		t.Fatalf("alice should accept reciprocal exchange for private data token: %#v", decision)
	}
}
