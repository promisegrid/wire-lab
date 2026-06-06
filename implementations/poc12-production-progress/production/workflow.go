package production

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const (
	PromiseWeighPackage          = "weigh_package"
	PromiseAddressLookup         = "address_lookup"
	PromisePrintLabel            = "print_label"
	PromiseShipmentUpdate        = "shipment_update"
	PromiseIssuePrintCapability  = "issue_print_capability"
	PromiseRedeemPrintCapability = "redeem_print_capability"
	PrintCapabilityScope         = "print_label"
	PrintCapabilityMaxBytes      = 4096
)

// WeightForPackage simulates a deterministic postal scale reading.
// Intent: Device agents make promises about local device state; they do not
// invent weights or accept commands from workflow agents. Source: DI-timah
func WeightForPackage(packageID string) (int, error) {
	packageID = strings.TrimSpace(packageID)
	if packageID == "" {
		return 0, fmt.Errorf("package_id is required")
	}
	if strings.Contains(strings.ToLower(packageID), "jammed") {
		return 0, fmt.Errorf("scale cannot currently read package %s", packageID)
	}
	return 16 + len(packageID)*3, nil
}

// AddressForOrder simulates an accounting-system address lookup.
// Intent: The accounting agent promises only address records it locally has,
// leaving the fulfillment agent to judge whether that evidence is enough.
// Source: DI-timah
func AddressForOrder(orderID string) (string, error) {
	orderID = strings.TrimSpace(orderID)
	if orderID == "" {
		return "", fmt.Errorf("order_id is required")
	}
	if strings.Contains(strings.ToLower(orderID), "unknown") {
		return "", fmt.Errorf("order %s is unknown to accounting", orderID)
	}
	return "100 Promise Way, Suite " + strconv.Itoa(len(orderID)*10) + ", Example City, CA 94000", nil
}

// LabelForShipment simulates a UPS label printer response from shipment facts.
// Intent: The printer promises label evidence derived from supplied package and
// address evidence; it does not promise shipment success outside its device
// boundary. Source: DI-timah
func LabelForShipment(packageID, address string, weightOunces int) (string, int, error) {
	packageID = strings.TrimSpace(packageID)
	address = strings.TrimSpace(address)
	if packageID == "" {
		return "", 0, fmt.Errorf("package_id is required")
	}
	if address == "" {
		return "", 0, fmt.Errorf("shipping address is required")
	}
	if weightOunces <= 0 {
		return "", 0, fmt.Errorf("positive weight_ounces is required")
	}
	digest := sha256.Sum256([]byte(packageID + "|" + address + "|" + strconv.Itoa(weightOunces)))
	trackingNumber := "1Z" + strings.ToUpper(hex.EncodeToString(digest[:]))[:14]
	costCents := 700 + weightOunces*4
	return trackingNumber, costCents, nil
}

// IssuePrintCapabilityToken creates a deterministic capability-promise token
// from the printer-port issuer's local promise fields.
// Intent: The token is evidence that `printer_port` promises bounded future
// label printing for one issuee and scope; it is not permission from an
// authority. Source: DI-pohaj; DI-vutok
func IssuePrintCapabilityToken(fields map[string]string) (string, error) {
	tokenID := strings.TrimSpace(fields["field_print_capability_token_id"])
	if tokenID == "" {
		tokenID = "printcap-" + strings.TrimSpace(fields["field_print_capability_issuee"])
	}
	scope := strings.TrimSpace(fields["field_print_capability_scope"])
	if scope == "" {
		scope = PrintCapabilityScope
	}
	maxBytes := fields["field_print_capability_max_bytes"]
	if maxBytes == "" {
		maxBytes = strconv.Itoa(PrintCapabilityMaxBytes)
	}
	token := strings.TrimSpace(fields["field_print_capability_issuee"])
	if token == "" {
		token = strings.TrimSpace(fields["to"])
	}
	if token == "" {
		return "", fmt.Errorf("print capability issuee is required")
	}
	if scope != PrintCapabilityScope {
		return "", fmt.Errorf("print capability scope %q is not supported", scope)
	}
	printEvidence := sha256.Sum256([]byte("printer_port|" + token + "|" + tokenID + "|" + scope + "|" + maxBytes))
	token = "pcap1:" + hex.EncodeToString(printEvidence[:])
	return token, nil
}

// ValidatePrintCapabilityToken verifies that a redemption presents the exact
// future-print promise token the printer-port app would have issued locally.
// Intent: Token redemption is local promise recognition by the issuer, not a
// global authorization check. Source: DI-pohaj; DI-vutok
func ValidatePrintCapabilityToken(fields map[string]string) error {
	token := strings.TrimSpace(fields["field_print_capability_token"])
	if token == "" {
		return fmt.Errorf("print capability token is required")
	}
	tokenID := strings.TrimSpace(fields["field_print_capability_token_id"])
	if tokenID == "" {
		return fmt.Errorf("print capability token_id is required")
	}
	scope := strings.TrimSpace(fields["field_print_capability_scope"])
	if scope != PrintCapabilityScope {
		return fmt.Errorf("print capability scope %q is not supported", scope)
	}
	maxBytes := intField(fields, "field_print_capability_max_bytes")
	if maxBytes <= 0 || maxBytes > PrintCapabilityMaxBytes {
		return fmt.Errorf("print capability max_bytes %d is outside local bounds", maxBytes)
	}
	labelBytes, decodeErr := hex.DecodeString(strings.TrimSpace(fields["field_label_bytes_hex"]))
	if decodeErr != nil {
		return fmt.Errorf("label bytes are not valid hex: %w", decodeErr)
	}
	if len(labelBytes) == 0 {
		return fmt.Errorf("label bytes are required")
	}
	if len(labelBytes) > maxBytes {
		return fmt.Errorf("label bytes length %d exceeds capability max_bytes %d", len(labelBytes), maxBytes)
	}
	capabilityFields := map[string]string{
		"to":                               strings.TrimSpace(fields["from"]),
		"field_print_capability_issuee":    strings.TrimSpace(fields["from"]),
		"field_print_capability_token_id":  tokenID,
		"field_print_capability_scope":     scope,
		"field_print_capability_max_bytes": strconv.Itoa(maxBytes),
	}
	printEvidence, issueErr := IssuePrintCapabilityToken(capabilityFields)
	if issueErr != nil {
		return issueErr
	}
	if token != printEvidence {
		return fmt.Errorf("print capability token does not match printer_port promise")
	}
	return nil
}

// LabelBytesForShipment returns the bounded bytes that the label-printer app
// asks the printer-port app to write to local printer hardware.
// Intent: The UPS label app can promise label-content generation, while the
// printer-port app separately promises local hardware access evidence.
// Source: DI-pohaj; DI-vutok
func LabelBytesForShipment(fields map[string]string) ([]byte, error) {
	if strings.TrimSpace(fields["field_package_id"]) == "" {
		return nil, fmt.Errorf("package_id is required")
	}
	if strings.TrimSpace(fields["field_tracking_number"]) == "" {
		return nil, fmt.Errorf("tracking_number is required")
	}
	if strings.TrimSpace(fields["field_cost_cents"]) == "" {
		return nil, fmt.Errorf("cost_cents is required")
	}
	labelBytes := []byte("UPS-LABEL\npackage=" + fields["field_package_id"] + "\ntracking=" + fields["field_tracking_number"] + "\ncost_cents=" + fields["field_cost_cents"] + "\n")
	if len(labelBytes) > PrintCapabilityMaxBytes {
		return nil, fmt.Errorf("label bytes length %d exceeds local max %d", len(labelBytes), PrintCapabilityMaxBytes)
	}
	return labelBytes, nil
}

// PrintLabelToLocalDevice simulates writing bounded label bytes to the local
// printer device owned by the printer-port kernel role.
// Intent: POC12 avoids real USB dependencies while still making hardware access
// a separate local promise surface with exact print evidence. Source: DI-pohaj;
// DI-vutok
func PrintLabelToLocalDevice(fields map[string]string) (string, error) {
	if validateErr := ValidatePrintCapabilityToken(fields); validateErr != nil {
		return "", validateErr
	}
	labelBytes, decodeErr := hex.DecodeString(strings.TrimSpace(fields["field_label_bytes_hex"]))
	if decodeErr != nil {
		return "", fmt.Errorf("label bytes are not valid hex: %w", decodeErr)
	}
	printEvidence := sha256.Sum256(labelBytes)
	spoolID := "spool-" + hex.EncodeToString(printEvidence[:])[:16]
	return spoolID, nil
}

// ValidateAccountingUpdate checks that a shipment update carries the minimum
// evidence the accounting agent needs before it promises to record the update.
func ValidateAccountingUpdate(orderID, trackingNumber string, costCents int) error {
	if strings.TrimSpace(orderID) == "" {
		return fmt.Errorf("order_id is required")
	}
	if !strings.HasPrefix(strings.TrimSpace(trackingNumber), "1Z") {
		return fmt.Errorf("tracking_number must start with 1Z")
	}
	if costCents <= 0 {
		return fmt.Errorf("positive shipping cost_cents is required")
	}
	return nil
}

func intField(fields map[string]string, keys ...string) int {
	for _, key := range keys {
		value := fields[key]
		if value == "" {
			continue
		}
		parsedValue, parseErr := strconv.Atoi(value)
		if parseErr == nil {
			return parsedValue
		}
	}
	return 0
}
