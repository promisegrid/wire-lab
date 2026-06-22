package specdocs

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// docs embeds the implementation-local POC16 protocol specs used for pCID
// derivation and live-agent prompt context.
//
// Intent: LLM-backed agents should see pCID spec prose for the protocols they
// promise to use, while runtime code still treats slot 0 as bytes and never does
// network routing by prose lookup. The source corpus stays under this POC's
// docs/protocols path so root docs do not become a stale competing authority.
// Source: DI-vulit; DI-magug
//
//go:embed *.md
var docs embed.FS

var protocolFiles = map[string]string{
	"accounting_v1":                     "accounting-v1.md",
	"cas_storage_v1":                    "cas-storage-v1.md",
	"cid_compute_v1":                    "cid-compute-v1.md",
	"encrypted_payload_v1":              "encrypted-payload-v1.md",
	"identity_key_v1":                   "identity-key-v1.md",
	"kernel_receive_v1":                 "kernel-receive-v1.md",
	"kernel_transport_v1":               "kernel-transport-v1.md",
	"map_payload_profile_v1":            "map-payload-profile-v1.md",
	"message_shape_cose_payload_v1":     "message-shape-cose-payload-v1.md",
	"message_shape_cose_proof_v1":       "message-shape-cose-proof-v1.md",
	"message_shape_envelope_parents_v1": "message-shape-envelope-parents-v1.md",
	"message_shape_native_proof_v1":     "message-shape-native-proof-v1.md",
	"message_shape_payload_parents_v1":  "message-shape-payload-parents-v1.md",
	"message_shape_transport_v1":        "message-shape-transport-v1.md",
	"parser_builder_role_v1":            "parser-builder-role-v1.md",
	"postal_scale_v1":                   "postal-scale-v1.md",
	"printer_port_v1":                   "printer-port-v1.md",
	"production_shipping_v1":            "production-shipping-v1.md",
	"relationship_v1":                   "relationship-v1.md",
	"route_v1":                          "route-v1.md",
	"secure_capability_v1":              "secure-capability-v1.md",
	"ups_label_v1":                      "ups-label-v1.md",
}

// BytesFor returns the exact embedded spec bytes for one protocol name.
func BytesFor(protocolName string) ([]byte, error) {
	fileName, ok := protocolFiles[protocolName]
	if !ok {
		return nil, fmt.Errorf("no embedded POC16 spec for %s", protocolName)
	}
	specBytes, readErr := docs.ReadFile(fileName)
	if readErr != nil {
		return nil, readErr
	}
	copiedBytes := make([]byte, len(specBytes))
	copy(copiedBytes, specBytes)
	return copiedBytes, nil
}

// SpecContext is one embedded spec excerpt plus a content hash of the embedded
// bytes supplied to an agent.
type SpecContext struct {
	Name    string
	SHA256  string
	Excerpt string
}

// ContextsFor returns deterministic spec snippets for locally supported pCIDs.
func ContextsFor(protocolNames []string, maxExcerptBytes int) ([]SpecContext, error) {
	names := append([]string{}, protocolNames...)
	sort.Strings(names)
	contexts := make([]SpecContext, 0, len(names))
	for _, name := range names {
		context, contextErr := ContextFor(name, maxExcerptBytes)
		if contextErr != nil {
			return nil, contextErr
		}
		contexts = append(contexts, context)
	}
	return contexts, nil
}

// ContextFor returns the embedded spec snippet for one protocol name.
func ContextFor(protocolName string, maxExcerptBytes int) (SpecContext, error) {
	specBytes, bytesErr := BytesFor(protocolName)
	if bytesErr != nil {
		return SpecContext{}, bytesErr
	}
	digest := sha256.Sum256(specBytes)
	excerpt := strings.TrimSpace(string(specBytes))
	if maxExcerptBytes > 0 && len(excerpt) > maxExcerptBytes {
		excerpt = excerpt[:maxExcerptBytes] + "\n..."
	}
	return SpecContext{
		Name:    protocolName,
		SHA256:  hex.EncodeToString(digest[:]),
		Excerpt: excerpt,
	}, nil
}
