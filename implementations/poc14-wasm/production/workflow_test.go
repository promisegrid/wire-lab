package production

import (
	"encoding/hex"
	"strconv"
	"testing"
)

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

func TestPrintCapabilityTokenAndLocalPrintEvent(t *testing.T) {
	capabilityFields := map[string]string{
		"to":                               "ups_label_printer",
		"field_print_capability_issuee":    "ups_label_printer",
		"field_print_capability_token_id":  "printcap-ups_label_printer",
		"field_print_capability_scope":     PrintCapabilityScope,
		"field_print_capability_max_bytes": strconv.Itoa(PrintCapabilityMaxBytes),
	}
	token, tokenErr := IssuePrintCapabilityToken(capabilityFields)
	if tokenErr != nil {
		t.Fatalf("issue token: %v", tokenErr)
	}
	labelBytes, labelErr := LabelBytesForShipment(map[string]string{
		"field_package_id":      "PKG-1001",
		"field_tracking_number": "1Z71051733616616",
		"field_cost_cents":      "860",
	})
	if labelErr != nil {
		t.Fatalf("label bytes: %v", labelErr)
	}
	redemptionFields := map[string]string{
		"from":                             "ups_label_printer",
		"field_print_capability_issuee":    "ups_label_printer",
		"field_print_capability_token":     token,
		"field_print_capability_token_id":  "printcap-ups_label_printer",
		"field_print_capability_scope":     PrintCapabilityScope,
		"field_print_capability_max_bytes": strconv.Itoa(PrintCapabilityMaxBytes),
		"field_label_bytes_hex":            hex.EncodeToString(labelBytes),
	}
	if err := ValidatePrintCapabilityToken(redemptionFields); err != nil {
		t.Fatalf("validate token: %v", err)
	}
	spoolID, printErr := PrintLabelToLocalDevice(redemptionFields)
	if printErr != nil {
		t.Fatalf("print local device: %v", printErr)
	}
	if spoolID == "" {
		t.Fatalf("spool id should be recorded")
	}
}

func TestPrintCapabilityRejectsWrongToken(t *testing.T) {
	redemptionFields := map[string]string{
		"from":                             "ups_label_printer",
		"field_print_capability_issuee":    "ups_label_printer",
		"field_print_capability_token":     "pcap1:bad",
		"field_print_capability_token_id":  "printcap-ups_label_printer",
		"field_print_capability_scope":     PrintCapabilityScope,
		"field_print_capability_max_bytes": strconv.Itoa(PrintCapabilityMaxBytes),
		"field_label_bytes_hex":            hex.EncodeToString([]byte("label")),
	}
	if err := ValidatePrintCapabilityToken(redemptionFields); err == nil {
		t.Fatalf("wrong token should be rejected")
	}
}
