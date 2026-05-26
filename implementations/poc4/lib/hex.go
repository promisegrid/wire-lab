package lib

import "encoding/hex"

// HexBytes renders bytes as a string field for relay wrapper payloads.
func HexBytes(value []byte) string {
	return hex.EncodeToString(value)
}

// ParseHexBytes decodes bytes carried in a relay wrapper payload field.
func ParseHexBytes(value string) ([]byte, error) {
	return hex.DecodeString(value)
}
