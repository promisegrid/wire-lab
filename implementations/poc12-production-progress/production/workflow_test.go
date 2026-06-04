package production

import "testing"

func TestProductionWorkflowDerivesShipmentFacts(t *testing.T) {
	weightOunces, weightErr := WeightForPackage("PKG-1001")
	if weightErr != nil {
		t.Fatalf("weight: %v", weightErr)
	}
	address, addressErr := AddressForOrder("ORDER-1001")
	if addressErr != nil {
		t.Fatalf("address: %v", addressErr)
	}
	trackingNumber, costCents, labelErr := LabelForShipment("PKG-1001", address, weightOunces)
	if labelErr != nil {
		t.Fatalf("label: %v", labelErr)
	}
	if err := ValidateAccountingUpdate("ORDER-1001", trackingNumber, costCents); err != nil {
		t.Fatalf("accounting update: %v", err)
	}
}

func TestProductionWorkflowRejectsBadFacts(t *testing.T) {
	if _, err := WeightForPackage(""); err == nil {
		t.Fatalf("empty package should be rejected")
	}
	if _, err := AddressForOrder("UNKNOWN-1"); err == nil {
		t.Fatalf("unknown order should be rejected")
	}
	if err := ValidateAccountingUpdate("ORDER-1001", "bad", 0); err == nil {
		t.Fatalf("bad accounting update should be rejected")
	}
}
