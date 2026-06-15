package pcid

import "testing"

func TestRegistryNamesRoundTrip(t *testing.T) {
	registry := NewRegistry()
	for _, name := range []string{
		RelationshipV1,
		PostalScaleV1,
		UPSLabelV1,
		AccountingV1,
		PrinterPortV1,
		KernelReceiveV1,
		CASStorageV1,
		CIDComputeV1,
		IdentityKeyV1,
		RouteV1,
		MessageShapeTransportV1,
		MessageShapeNativeProofV1,
		MessageShapeEnvelopeParentsV1,
		MessageShapePayloadParentsV1,
		MessageShapeCOSEPayloadV1,
		MessageShapeCOSEProofV1,
	} {
		protocolCID, ok := registry.CID(name)
		if !ok {
			t.Fatalf("missing protocol %s", name)
		}
		roundTripName, ok := registry.Name(protocolCID)
		if !ok || roundTripName != name {
			t.Fatalf("round trip name = %q %v, want %q true", roundTripName, ok, name)
		}
	}
}
