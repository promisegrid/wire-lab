package policy

import (
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc8-autonomous-promise-economy/token"
)

func TestAliceCountersExpensiveStorageOffer(t *testing.T) {
	// Intent: Alice can ask for different reciprocal promises without commanding
	// Bob to accept them. Source: DI-sirus
	alicePolicy := ForNode("alice")
	decision := alicePolicy.Decide(ActionContext{
		Action:       ActionCounter,
		Peer:         "bob",
		ResourceKind: "storage",
		ResourceID:   "alice-report",
		Price:        8,
		MaxPrice:     8,
		Stake:        2,
	})
	if !decision.Accepted {
		t.Fatalf("alice should counter storage terms she can locally accept: %#v", decision)
	}
}

func TestDaveAcceptsFirstStaleTokenForEvidence(t *testing.T) {
	// Intent: Dave's first stale-token acceptance is voluntary evidence gathering,
	// not permission from Alice or any central authority. Source: DI-sirus
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
	// Source: DI-sirus
	davePolicy := ForNode("dave")
	decision := davePolicy.Decide(ActionContext{
		Action:       ActionAccept,
		Peer:         "mallory",
		SourcePeer:   "mallory",
		Issuer:       "alice",
		ResourceKind: "data",
		ResourceID:   "dataset-stale",
		TransferRule: token.TransferBearer,
		IssuerTrust:  -1,
		SourceTrust:  -1,
		PeerTrust:    -1,
	})
	if decision.Accepted {
		t.Fatalf("dave should refuse later mallory stale token after broken evidence: %#v", decision)
	}
}

func TestCarolAcceptsBearerForNonTransferableCompute(t *testing.T) {
	// Intent: Bearer-for-non-transferable exchange clears Carol's local threshold
	// because Carol values Bob's bearer stake token enough to issue Alice scoped
	// compute access. Source: DI-sirus
	carolPolicy := ForNode("carol")
	decision := carolPolicy.Decide(ActionContext{
		Action:       ActionAccept,
		Peer:         "alice",
		Issuer:       "bob",
		ResourceKind: "storage",
		ResourceID:   "bob-storage-stake",
		TransferRule: token.TransferBearer,
	})
	if !decision.Accepted {
		t.Fatalf("carol should accept Bob bearer stake token for non-transferable compute access: %#v", decision)
	}
}
