package production

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const (
	PromiseWeighPackage   = "weigh_package"
	PromiseAddressLookup  = "address_lookup"
	PromisePrintLabel     = "print_label"
	PromiseShipmentUpdate = "shipment_update"
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
