package poc13

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sort"
	"unicode/utf8"
)

const (
	gridTag        = uint64(0x67726964)
	dagCBORLinkTag = uint64(42)

	CASStorageV1     = "cas_storage_v1"
	CIDComputeV1     = "cid_compute_v1"
	EvidenceReportV1 = "evidence_report_v1"
)

// ProtocolCID is the content-derived identity of a protocol spec document, not
// the identity of one message, payload, function, or content object.
// Intent: POC13 must preserve the pCID distinction while testing payload-level
// content CIDs and function CIDs. Source: DI-notig
type ProtocolCID struct {
	cidBytes []byte
	digest   [32]byte
}

// Registry is the POC13 local table of provisional protocol spec names.
type Registry struct {
	byName map[string]ProtocolCID
	byCID  map[string]string
}

// NewRegistry returns POC-local pCIDs derived from concise stand-in spec text.
func NewRegistry() Registry {
	registry := Registry{byName: make(map[string]ProtocolCID), byCID: make(map[string]string)}
	for _, entry := range []struct {
		name string
		spec string
	}{
		{CASStorageV1, "poc13 decentralized cas storage promise protocol v1"},
		{CIDComputeV1, "poc13 cid named function compute promise protocol v1"},
		{EvidenceReportV1, "poc13 local evidence report promise protocol v1"},
	} {
		registry.register(entry.name, NewProtocolCID([]byte(entry.spec)))
	}
	return registry
}

func (registry Registry) register(name string, protocolCID ProtocolCID) {
	registry.byName[name] = protocolCID
	registry.byCID[protocolCID.String()] = name
}

// MustCID returns a protocol pCID for POC-local programming errors.
func (registry Registry) MustCID(name string) ProtocolCID {
	protocolCID, ok := registry.byName[name]
	if !ok {
		panic(fmt.Sprintf("unknown protocol name %s", name))
	}
	return protocolCID
}

// Name returns the local protocol name for one parsed envelope pCID.
func (registry Registry) Name(protocolCID ProtocolCID) (string, bool) {
	name, ok := registry.byCID[protocolCID.String()]
	return name, ok
}

// NewProtocolCID hashes spec bytes into the POC CIDv1 raw sha2-256 form.
func NewProtocolCID(specBytes []byte) ProtocolCID {
	digest := sha256.Sum256(specBytes)
	cidBytes := make([]byte, 0, 36)
	cidBytes = append(cidBytes, 0x01, 0x55, 0x12, 0x20)
	cidBytes = append(cidBytes, digest[:]...)
	return ProtocolCID{cidBytes: cidBytes, digest: digest}
}

func newProtocolCIDFromBytes(cidBytes []byte) ProtocolCID {
	copiedBytes := append([]byte(nil), cidBytes...)
	var digest [32]byte
	if len(copiedBytes) >= 36 {
		copy(digest[:], copiedBytes[len(copiedBytes)-32:])
	}
	return ProtocolCID{cidBytes: copiedBytes, digest: digest}
}

// Tag42Bytes returns the DAG-CBOR tag-42 byte string payload.
func (protocolCID ProtocolCID) Tag42Bytes() []byte {
	tagBytes := make([]byte, 0, len(protocolCID.cidBytes)+1)
	tagBytes = append(tagBytes, 0x00)
	tagBytes = append(tagBytes, protocolCID.cidBytes...)
	return tagBytes
}

// String renders a stable local evidence string.
func (protocolCID ProtocolCID) String() string {
	return "cidv1-raw-sha2-256:" + hex.EncodeToString(protocolCID.digest[:])
}

// ContentCID returns a POC CIDv1 raw sha2-256 identity for bytes or function
// code carried in payload fields.
func ContentCID(content []byte) string {
	digest := sha256.Sum256(content)
	return "cidv1-raw-sha2-256:" + hex.EncodeToString(digest[:])
}

// VerifyContentCID checks that bytes match a payload-level CID claim.
func VerifyContentCID(content []byte, cid string) bool {
	return ContentCID(content) == cid
}

// Proof is a deterministic POC Ed25519 proof over the pCID-defined signable
// view.
type Proof struct {
	Signer    string
	PublicKey []byte
	Signature []byte
}

// Envelope represents grid([42(pCID), payload, proof]).
// Intent: Go owns exact bytes and signatures; LLM output only influences
// payload meanings after validation. Source: DI-notig
type Envelope struct {
	ProtocolCID ProtocolCID
	Payload     []byte
	Proof       Proof
}

// NewEnvelope signs one pCID-owned payload field map.
func NewEnvelope(protocolCID ProtocolCID, fields map[string]string, signer string) (Envelope, error) {
	payloadBytes, marshalErr := MarshalStringMap(fields)
	if marshalErr != nil {
		return Envelope{}, marshalErr
	}
	envelope := Envelope{ProtocolCID: protocolCID, Payload: payloadBytes}
	proof, proofErr := SignEnvelope(envelope, signer)
	if proofErr != nil {
		return Envelope{}, proofErr
	}
	envelope.Proof = proof
	return envelope, nil
}

// SignableBytes serializes grid([42(pCID), payload]) for proof generation.
func (envelope Envelope) SignableBytes() ([]byte, error) {
	writer := &cborWriter{}
	if err := writer.writeTag(gridTag); err != nil {
		return nil, err
	}
	if err := writer.writeArrayHeader(2); err != nil {
		return nil, err
	}
	if err := writer.writeTag(dagCBORLinkTag); err != nil {
		return nil, err
	}
	if err := writer.writeBytes(envelope.ProtocolCID.Tag42Bytes()); err != nil {
		return nil, err
	}
	if err := writer.writeBytes(envelope.Payload); err != nil {
		return nil, err
	}
	return writer.buffer.Bytes(), nil
}

// Bytes serializes the full signed envelope.
func (envelope Envelope) Bytes() ([]byte, error) {
	proofBytes, proofErr := MarshalStringMap(map[string]string{
		"signer":     envelope.Proof.Signer,
		"public_key": hex.EncodeToString(envelope.Proof.PublicKey),
		"signature":  hex.EncodeToString(envelope.Proof.Signature),
	})
	if proofErr != nil {
		return nil, proofErr
	}
	writer := &cborWriter{}
	if err := writer.writeTag(gridTag); err != nil {
		return nil, err
	}
	if err := writer.writeArrayHeader(3); err != nil {
		return nil, err
	}
	if err := writer.writeTag(dagCBORLinkTag); err != nil {
		return nil, err
	}
	if err := writer.writeBytes(envelope.ProtocolCID.Tag42Bytes()); err != nil {
		return nil, err
	}
	if err := writer.writeBytes(envelope.Payload); err != nil {
		return nil, err
	}
	if err := writer.writeBytes(proofBytes); err != nil {
		return nil, err
	}
	return writer.buffer.Bytes(), nil
}

// ParseEnvelope parses grid([42(pCID), payload, proof]).
func ParseEnvelope(envelopeBytes []byte) (Envelope, error) {
	reader := &cborReader{data: envelopeBytes}
	outerTag, outerTagErr := reader.readTypeAndLength(6)
	if outerTagErr != nil {
		return Envelope{}, outerTagErr
	}
	if outerTag != gridTag {
		return Envelope{}, fmt.Errorf("outer envelope must be grid tag %d, got %d", gridTag, outerTag)
	}
	arrayLength, arrayErr := reader.readTypeAndLength(4)
	if arrayErr != nil {
		return Envelope{}, arrayErr
	}
	if arrayLength != 3 {
		return Envelope{}, fmt.Errorf("poc13 envelope must have three slots, got %d", arrayLength)
	}
	tagNumber, tagErr := reader.readTypeAndLength(6)
	if tagErr != nil {
		return Envelope{}, tagErr
	}
	if tagNumber != dagCBORLinkTag {
		return Envelope{}, fmt.Errorf("slot 0 must be tag 42, got %d", tagNumber)
	}
	tagBytes, tagBytesErr := reader.readBytes()
	if tagBytesErr != nil {
		return Envelope{}, tagBytesErr
	}
	if len(tagBytes) < 2 || tagBytes[0] != 0x00 {
		return Envelope{}, fmt.Errorf("tag 42 payload must start with DAG-CBOR CID sentinel")
	}
	payloadBytes, payloadErr := reader.readBytes()
	if payloadErr != nil {
		return Envelope{}, payloadErr
	}
	proofBytes, proofErr := reader.readBytes()
	if proofErr != nil {
		return Envelope{}, proofErr
	}
	if reader.offset != len(reader.data) {
		return Envelope{}, fmt.Errorf("trailing cbor bytes in envelope: %d", len(reader.data)-reader.offset)
	}
	proofFields, fieldsErr := UnmarshalStringMap(proofBytes)
	if fieldsErr != nil {
		return Envelope{}, fieldsErr
	}
	publicKey, publicErr := hex.DecodeString(proofFields["public_key"])
	if publicErr != nil {
		return Envelope{}, publicErr
	}
	signature, signatureErr := hex.DecodeString(proofFields["signature"])
	if signatureErr != nil {
		return Envelope{}, signatureErr
	}
	return Envelope{
		ProtocolCID: newProtocolCIDFromBytes(tagBytes[1:]),
		Payload:     payloadBytes,
		Proof: Proof{
			Signer:    proofFields["signer"],
			PublicKey: publicKey,
			Signature: signature,
		},
	}, nil
}

// PayloadFields decodes the pCID-owned payload.
func (envelope Envelope) PayloadFields() (map[string]string, error) {
	return UnmarshalStringMap(envelope.Payload)
}

// SignEnvelope signs the pCID-defined signable view.
func SignEnvelope(envelope Envelope, signer string) (Proof, error) {
	signableBytes, signableErr := envelope.SignableBytes()
	if signableErr != nil {
		return Proof{}, signableErr
	}
	privateKey := deterministicPrivateKey(signer)
	return Proof{
		Signer:    signer,
		PublicKey: deterministicPublicKey(signer),
		Signature: ed25519.Sign(privateKey, signableBytes),
	}, nil
}

// VerifyEnvelope checks the Ed25519 proof.
func VerifyEnvelope(envelope Envelope) error {
	if len(envelope.Proof.PublicKey) != ed25519.PublicKeySize {
		return fmt.Errorf("unexpected public key length %d", len(envelope.Proof.PublicKey))
	}
	if len(envelope.Proof.Signature) != ed25519.SignatureSize {
		return fmt.Errorf("unexpected signature length %d", len(envelope.Proof.Signature))
	}
	signableBytes, signableErr := envelope.SignableBytes()
	if signableErr != nil {
		return signableErr
	}
	if !ed25519.Verify(envelope.Proof.PublicKey, signableBytes, envelope.Proof.Signature) {
		return fmt.Errorf("envelope signature failed")
	}
	return nil
}

func deterministicPrivateKey(seedText string) ed25519.PrivateKey {
	seed := sha256.Sum256([]byte("poc13 protocol signer: " + seedText))
	return ed25519.NewKeyFromSeed(seed[:])
}

func deterministicPublicKey(seedText string) ed25519.PublicKey {
	return deterministicPrivateKey(seedText).Public().(ed25519.PublicKey)
}

// HashExactBytes returns a sha256 hex digest for local evidence correlation.
func HashExactBytes(exactBytes []byte) string {
	digest := sha256.Sum256(exactBytes)
	return hex.EncodeToString(digest[:])
}

type cborWriter struct {
	buffer bytes.Buffer
}

func (writer *cborWriter) writeTypeAndLength(major byte, length uint64) error {
	prefix := major << 5
	switch {
	case length < 24:
		return writer.buffer.WriteByte(prefix | byte(length))
	case length <= 0xff:
		if err := writer.buffer.WriteByte(prefix | 24); err != nil {
			return err
		}
		return writer.buffer.WriteByte(byte(length))
	case length <= 0xffff:
		if err := writer.buffer.WriteByte(prefix | 25); err != nil {
			return err
		}
		return binary.Write(&writer.buffer, binary.BigEndian, uint16(length))
	case length <= 0xffffffff:
		if err := writer.buffer.WriteByte(prefix | 26); err != nil {
			return err
		}
		return binary.Write(&writer.buffer, binary.BigEndian, uint32(length))
	default:
		if err := writer.buffer.WriteByte(prefix | 27); err != nil {
			return err
		}
		return binary.Write(&writer.buffer, binary.BigEndian, length)
	}
}

func (writer *cborWriter) writeArrayHeader(length int) error {
	return writer.writeTypeAndLength(4, uint64(length))
}

func (writer *cborWriter) writeMapHeader(length int) error {
	return writer.writeTypeAndLength(5, uint64(length))
}

func (writer *cborWriter) writeTag(tag uint64) error {
	return writer.writeTypeAndLength(6, tag)
}

func (writer *cborWriter) writeBytes(value []byte) error {
	if err := writer.writeTypeAndLength(2, uint64(len(value))); err != nil {
		return err
	}
	_, err := writer.buffer.Write(value)
	return err
}

func (writer *cborWriter) writeString(value string) error {
	if !utf8.ValidString(value) {
		return fmt.Errorf("invalid utf8 string")
	}
	if err := writer.writeTypeAndLength(3, uint64(len(value))); err != nil {
		return err
	}
	_, err := writer.buffer.WriteString(value)
	return err
}

// MarshalStringMap encodes payload fields as a deterministic CBOR map.
func MarshalStringMap(fields map[string]string) ([]byte, error) {
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	writer := &cborWriter{}
	if err := writer.writeMapHeader(len(keys)); err != nil {
		return nil, err
	}
	for _, key := range keys {
		if err := writer.writeString(key); err != nil {
			return nil, err
		}
		if err := writer.writeString(fields[key]); err != nil {
			return nil, err
		}
	}
	return writer.buffer.Bytes(), nil
}

type cborReader struct {
	data   []byte
	offset int
}

func (reader *cborReader) readByte() (byte, error) {
	if reader.offset >= len(reader.data) {
		return 0, fmt.Errorf("unexpected end of cbor data")
	}
	value := reader.data[reader.offset]
	reader.offset++
	return value, nil
}

func (reader *cborReader) readTypeAndLength(expectedMajor byte) (uint64, error) {
	initial, err := reader.readByte()
	if err != nil {
		return 0, err
	}
	major := initial >> 5
	if major != expectedMajor {
		return 0, fmt.Errorf("expected cbor major %d, got %d", expectedMajor, major)
	}
	additional := initial & 0x1f
	switch additional {
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23:
		return uint64(additional), nil
	case 24:
		value, readErr := reader.readByte()
		return uint64(value), readErr
	case 25:
		if reader.offset+2 > len(reader.data) {
			return 0, fmt.Errorf("truncated uint16")
		}
		value := binary.BigEndian.Uint16(reader.data[reader.offset : reader.offset+2])
		reader.offset += 2
		return uint64(value), nil
	case 26:
		if reader.offset+4 > len(reader.data) {
			return 0, fmt.Errorf("truncated uint32")
		}
		value := binary.BigEndian.Uint32(reader.data[reader.offset : reader.offset+4])
		reader.offset += 4
		return uint64(value), nil
	case 27:
		if reader.offset+8 > len(reader.data) {
			return 0, fmt.Errorf("truncated uint64")
		}
		value := binary.BigEndian.Uint64(reader.data[reader.offset : reader.offset+8])
		reader.offset += 8
		return value, nil
	default:
		return 0, fmt.Errorf("unsupported cbor additional information %d", additional)
	}
}

func (reader *cborReader) readBytes() ([]byte, error) {
	length, err := reader.readTypeAndLength(2)
	if err != nil {
		return nil, err
	}
	if reader.offset+int(length) > len(reader.data) {
		return nil, fmt.Errorf("truncated byte string")
	}
	value := append([]byte(nil), reader.data[reader.offset:reader.offset+int(length)]...)
	reader.offset += int(length)
	return value, nil
}

func (reader *cborReader) readString() (string, error) {
	length, err := reader.readTypeAndLength(3)
	if err != nil {
		return "", err
	}
	if reader.offset+int(length) > len(reader.data) {
		return "", fmt.Errorf("truncated string")
	}
	value := string(reader.data[reader.offset : reader.offset+int(length)])
	reader.offset += int(length)
	if !utf8.ValidString(value) {
		return "", fmt.Errorf("invalid utf8 string")
	}
	return value, nil
}

// UnmarshalStringMap decodes a deterministic CBOR string map.
func UnmarshalStringMap(payloadBytes []byte) (map[string]string, error) {
	reader := &cborReader{data: payloadBytes}
	length, err := reader.readTypeAndLength(5)
	if err != nil {
		return nil, err
	}
	fields := make(map[string]string, int(length))
	for index := uint64(0); index < length; index++ {
		key, keyErr := reader.readString()
		if keyErr != nil {
			return nil, keyErr
		}
		value, valueErr := reader.readString()
		if valueErr != nil {
			return nil, valueErr
		}
		fields[key] = value
	}
	if reader.offset != len(reader.data) {
		return nil, fmt.Errorf("trailing cbor bytes in map: %d", len(reader.data)-reader.offset)
	}
	return fields, nil
}
