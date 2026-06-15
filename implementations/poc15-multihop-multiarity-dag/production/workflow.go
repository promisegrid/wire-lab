package production

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// These promise-about names are payload-level meanings used inside pCID-owned
// protocols; they are not top-level wire actions or independent pCIDs.
// Intent: Keep PromiseGrid protocol vocabulary promise-first while moving key
// rotation into identity_key_v1 instead of generic report pCIDs. Source:
// DI-bikit; DI-vipih
const (
	PromiseWeighPackage               = "weigh_package"
	PromiseAddressLookup              = "address_lookup"
	PromisePrintLabel                 = "print_label"
	PromiseShipmentUpdate             = "shipment_update"
	PromiseIssuePrintCapability       = "issue_print_capability"
	PromiseRedeemPrintCapability      = "redeem_print_capability"
	PromiseStoreContent               = "store_content"
	PromiseServeContent               = "serve_content"
	PromiseReplicateContent           = "replicate_content"
	PromiseServeReplicaContent        = "serve_replica_content"
	PromiseReplicaTokenLifecycle      = "replica_token_lifecycle"
	PromisePresentStorageReport       = "present_storage_report"
	PromiseExecuteFunction            = "execute_function"
	PromiseLookupComputeCache         = "lookup_compute_cache"
	PromiseProvideComputeContext      = "provide_compute_context"
	PromiseVerifyComputeResult        = "verify_compute_result"
	PromiseRotateSigningKey           = "rotate_signing_key"
	PromiseLabelFutureMalformedReport = "label_future_malformed_report"
	PromiseUnsupportedVariantProbe    = "unsupported_variant_probe"
	PromiseRouteSetup                 = "route_setup"
	PromiseRouteForward               = "route_forward"
	PromiseRouteReachability          = "route_reachability"
	PrintCapabilityScope              = "print_label"
	PrintCapabilityMaxBytes           = 4096
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
// leaving the fulfillment agent to judge whether those event records are enough.
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
// Intent: The printer promises label event records derived from supplied package
// and address event records; it does not promise shipment success outside its device
// runtime adapter. Source: DI-timah
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
// Intent: The token is an event record that `printer_port` promises bounded future
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
	printEvent := sha256.Sum256([]byte("printer_port|" + token + "|" + tokenID + "|" + scope + "|" + maxBytes))
	token = "pcap1:" + hex.EncodeToString(printEvent[:])
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
	printEvent, issueErr := IssuePrintCapabilityToken(capabilityFields)
	if issueErr != nil {
		return issueErr
	}
	if token != printEvent {
		return fmt.Errorf("print capability token does not match printer_port promise")
	}
	return nil
}

// LabelBytesForShipment returns the bounded bytes that the label-printer app
// asks the printer-port app to write to local printer hardware.
// Intent: The UPS label app can promise label-content generation, while the
// printer-port app separately promises local hardware access event records.
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
// Intent: POC15 avoids real USB dependencies while still making hardware access
// a separate local promise surface with exact print event records. Source: DI-pohaj;
// DI-vutok
func PrintLabelToLocalDevice(fields map[string]string) (string, error) {
	if validateErr := ValidatePrintCapabilityToken(fields); validateErr != nil {
		return "", validateErr
	}
	labelBytes, decodeErr := hex.DecodeString(strings.TrimSpace(fields["field_label_bytes_hex"]))
	if decodeErr != nil {
		return "", fmt.Errorf("label bytes are not valid hex: %w", decodeErr)
	}
	printEvent := sha256.Sum256(labelBytes)
	spoolID := "spool-" + hex.EncodeToString(printEvent[:])[:16]
	return spoolID, nil
}

// ContentCID returns a POC CIDv1 raw sha2-256 identity for stored content,
// function source bytes, inputs, contexts, and compute results.
// Intent: POC15 preserves the distinction between pCID protocol identity and
// payload-level CIDs while making CAS and compute event records exact-byte checkable.
// Source: DI-sinur
func ContentCID(content []byte) string {
	digest := sha256.Sum256(content)
	return "cidv1-raw-sha2-256:" + hex.EncodeToString(digest[:])
}

// VerifyContentCID checks that a byte string matches the payload-level CID an
// agent promised it represented.
// Intent: Receivers independently verify bytes from their local vantage instead
// of accepting sender claims as authority. Source: DI-sinur
func VerifyContentCID(content []byte, cid string) bool {
	return ContentCID(content) == cid
}

// SampleContentBytes supplies deterministic private bytes for the executable
// CAS path without committing real user data.
func SampleContentBytes() []byte {
	return []byte("poc15 sample invoice bytes")
}

// SampleSecondContentBytes creates independent storage pressure so the CAS path
// cannot pass by handling only one hard-coded object.
func SampleSecondContentBytes() []byte {
	return []byte("poc15 second sample receipt bytes")
}

// CorruptContentBytes is Mallory's malformed event probe; its claimed CID is
// intentionally different from the byte content.
func CorruptContentBytes() []byte {
	return []byte("poc15 sample invoice bytes corrupted by mallory")
}

func SampleFunctionBytes() []byte {
	return []byte("poc15 function: fibonacci(n) v1")
}

func SampleSumFunctionBytes() []byte {
	return []byte("poc15 function: sum(values) v1")
}

func SampleInputBytes() []byte {
	return []byte("n=9")
}

func SampleSumInputBytes() []byte {
	return []byte("values=2,3,5")
}

func SampleContextBytes() []byte {
	return []byte("timestamp=2026-06-06T00:00:00Z;randomness=explicit-none")
}

// ComputeInputs carries the exact bytes named by a cid_compute_v1 promise.
// Intent: Runtime adapters should share the same decode and CID-check path as
// ordinary Go compute peers, so WASM and stdio execution differ only by local
// runtime mechanics rather than by protocol semantics. Source: DI-sivis
type ComputeInputs struct {
	FunctionBytes []byte
	InputBytes    []byte
	ContextBytes  []byte
}

// ComputeCacheKey names one exact local compute checkpoint over all bytes that
// make the result meaningful.
// Intent: Cache reuse is an event record over a pCID-defined tuple, not a global
// execution result authority. Source: DI-sinur
func ComputeCacheKey(protocolName, functionCID, inputCID, contextCID, resultCID string) string {
	return ContentCID([]byte(protocolName + "|" + functionCID + "|" + inputCID + "|" + contextCID + "|" + resultCID))
}

// ExecuteFunction runs only the two bounded function forms used by this POC.
// Intent: POC15 tests CID-named arbitrary function bytes without embedding a
// general-purpose remote execution engine or RPC command surface. Source: DI-sinur
func ExecuteFunction(functionBytes, inputBytes, contextBytes []byte) ([]byte, error) {
	functionText := strings.TrimSpace(string(functionBytes))
	inputText := strings.TrimSpace(string(inputBytes))
	switch {
	case strings.Contains(functionText, "fibonacci"):
		n, err := ParseFibonacciInput(inputBytes)
		if err != nil {
			return nil, err
		}
		return FibonacciResultBytes(n, uint64(fibonacci(n)), contextBytes), nil
	case strings.Contains(functionText, "sum"):
		if !strings.HasPrefix(inputText, "values=") {
			return nil, fmt.Errorf("unsupported sum input %q", inputText)
		}
		rawValues := strings.TrimPrefix(inputText, "values=")
		total := 0
		for _, rawValue := range strings.Split(rawValues, ",") {
			value, err := strconv.Atoi(strings.TrimSpace(rawValue))
			if err != nil {
				return nil, err
			}
			total += value
		}
		return []byte(fmt.Sprintf("sum(%s)=%d;context_cid=%s", rawValues, total, ContentCID(contextBytes))), nil
	default:
		return nil, fmt.Errorf("unsupported function source %q", functionText)
	}
}

// DecodeComputeInputs decodes the compatibility field projection of a
// cid_compute_v1 payload back into exact function/input/context bytes.
// Intent: Runtime adapters and ordinary compute peers should verify the same
// exact bytes before promising compute results. Source: DI-sivis
func DecodeComputeInputs(fields map[string]string) (ComputeInputs, error) {
	functionBytes, functionErr := base64.StdEncoding.DecodeString(fields["field_function_b64"])
	if functionErr != nil {
		return ComputeInputs{}, functionErr
	}
	inputBytes, inputErr := base64.StdEncoding.DecodeString(fields["field_input_b64"])
	if inputErr != nil {
		return ComputeInputs{}, inputErr
	}
	contextBytes, contextErr := base64.StdEncoding.DecodeString(fields["field_context_b64"])
	if contextErr != nil {
		return ComputeInputs{}, contextErr
	}
	return ComputeInputs{FunctionBytes: functionBytes, InputBytes: inputBytes, ContextBytes: contextBytes}, nil
}

// VerifyComputeInputCIDs checks that the exact bytes carried in a compute
// promise match the CIDs that make the promise meaningful.
// Intent: Heterogeneous compute runtimes must not bypass pCID-owned byte/CID
// checks while proving useful work. Source: DI-sivis
func VerifyComputeInputCIDs(fields map[string]string, inputs ComputeInputs) error {
	if !VerifyContentCID(inputs.FunctionBytes, fields["field_function_cid"]) {
		return fmt.Errorf("function bytes do not match function CID")
	}
	if !VerifyContentCID(inputs.InputBytes, fields["field_input_cid"]) {
		return fmt.Errorf("input bytes do not match input CID")
	}
	if !VerifyContentCID(inputs.ContextBytes, fields["field_context_cid"]) {
		return fmt.Errorf("context bytes do not match context CID")
	}
	return nil
}

// ExecuteComputePromiseFields runs a cid_compute_v1 execute_function promise and
// returns the ACK fields common to Go, WASM, and stdio compute peers.
// Intent: Peggy and Victor should differ only in local execution mechanics; the
// ACK payload remains the existing compute promise shape. Source: DI-sivis
func ExecuteComputePromiseFields(fields map[string]string, resultBytes []byte) map[string]string {
	return map[string]string{
		"field_promise_about": PromiseExecuteFunction,
		"field_function_cid":  fields["field_function_cid"],
		"field_function_b64":  fields["field_function_b64"],
		"field_input_cid":     fields["field_input_cid"],
		"field_input_b64":     fields["field_input_b64"],
		"field_context_cid":   fields["field_context_cid"],
		"field_context_b64":   fields["field_context_b64"],
		"field_result_cid":    ContentCID(resultBytes),
		"field_result_b64":    base64.StdEncoding.EncodeToString(resultBytes),
	}
}

// ParseFibonacciInput returns the bounded Fibonacci input from existing POC15
// compute payload bytes.
// Intent: Peggy's embedded WASM module accepts a numeric input but the protocol
// still carries the existing exact input bytes. Source: DI-sivis
func ParseFibonacciInput(inputBytes []byte) (int, error) {
	inputText := strings.TrimSpace(string(inputBytes))
	if !strings.HasPrefix(inputText, "n=") {
		return 0, fmt.Errorf("unsupported fibonacci input %q", inputText)
	}
	var n int
	if _, err := fmt.Sscanf(inputText, "n=%d", &n); err != nil {
		return 0, err
	}
	if n < 0 || n > 40 {
		return 0, fmt.Errorf("fibonacci input out of bounded POC range: %d", n)
	}
	return n, nil
}

// FibonacciResultBytes formats a result from a local runtime-specific Fibonacci
// executor into the same bytes expected by cid_compute_v1 verification.
// Intent: WASM execution should produce the same promise-visible result bytes as
// ordinary compute verification expects. Source: DI-sivis
func FibonacciResultBytes(n int, value uint64, contextBytes []byte) []byte {
	return []byte(fmt.Sprintf("fibonacci(%d)=%d;context_cid=%s", n, value, ContentCID(contextBytes)))
}

func fibonacci(n int) int {
	if n < 2 {
		return n
	}
	previous, current := 0, 1
	for index := 2; index <= n; index++ {
		previous, current = current, previous+current
	}
	return current
}

// BadComputeResultBytes produces bytes that have their own valid CID but fail
// the promised function/input/context semantics.
// Intent: Verifiers must detect semantic breakage, not just malformed hashes.
// Source: DI-sinur
func BadComputeResultBytes(correctResultBytes []byte) []byte {
	return append(append([]byte(nil), correctResultBytes...), []byte(";bad-poc15-result")...)
}

// FunctionKind gives logs a compact label without changing the pCID-owned
// payload contract.
func FunctionKind(functionBytes []byte) string {
	functionText := strings.TrimSpace(string(functionBytes))
	if strings.Contains(functionText, "fibonacci") {
		return "fibonacci"
	}
	if strings.Contains(functionText, "sum") {
		return "sum"
	}
	return "unknown"
}

// ValidateAccountingUpdate checks that a shipment update carries the minimum
// event records the accounting agent needs before it promises to record the update.
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
