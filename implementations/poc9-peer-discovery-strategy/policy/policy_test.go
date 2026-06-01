package policy

import (
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc9-peer-discovery-strategy/token"
)

func TestAliceCountersExpensiveStorageOffer(t *testing.T) {
	// Intent: Alice can ask for different reciprocal promises without commanding
	// Bob to accept them. Source: DI-sipuz
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

func TestDaveAcceptsFirstExpiredTokenForEvidence(t *testing.T) {
	// Intent: Dave's first expired-token acceptance is voluntary evidence
	// gathering about Mallory's presentation, not evidence that Alice broke a
	// promise. Source: DI-sipuz; DI-vujil
	davePolicy := ForNode("dave")
	decision := davePolicy.Decide(ActionContext{
		Action:       ActionAccept,
		Peer:         "mallory",
		SourcePeer:   "mallory",
		Issuer:       "alice",
		ResourceKind: "data",
		ResourceID:   "dataset-expired",
		TransferRule: token.TransferBearer,
	})
	if !decision.Accepted {
		t.Fatalf("dave should accept first expired token for local evidence: %#v", decision)
	}
}

func TestDaveRefusesMalloryAfterExpiredTokenMisrepresentation(t *testing.T) {
	// Intent: Expired-token evidence should affect Dave's trust in Mallory as the
	// circulator, while Alice remains neutral because her signed expiry promise was
	// kept. Source: DI-sipuz; DI-vujil
	davePolicy := ForNode("dave")
	decision := davePolicy.Decide(ActionContext{
		Action:       ActionAccept,
		Peer:         "mallory",
		SourcePeer:   "mallory",
		Issuer:       "alice",
		ResourceKind: "data",
		ResourceID:   "dataset-expired",
		TransferRule: token.TransferBearer,
		SourceTrust:  -2,
		PeerTrust:    -2,
		ExpectedRisk: 3,
	})
	if decision.Accepted {
		t.Fatalf("dave should refuse later mallory expired token after circulator evidence: %#v", decision)
	}
}

func TestCarolAcceptsBearerForNonTransferableCompute(t *testing.T) {
	// Intent: Bearer-for-non-transferable exchange clears Carol's local threshold
	// because Carol values Bob's bearer stake token enough to issue Alice scoped
	// compute access. Source: DI-sipuz
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
