package protocol

import (
	"crypto/ed25519"
	"fmt"
)

const (
	coseSign1Tag       = uint64(18)
	coseHeaderAlg      = int64(1)
	coseAlgorithmEdDSA = int64(-8)
)

// COSESign1 is the parsed subset of COSE_Sign1 that POC15 needs for envelope
// shape pressure tests.
// Intent: POC15 should test COSE as a protocol-owned payload/proof option while
// keeping COSE out of the universal envelope invariant. Source: DI-mosat
type COSESign1 struct {
	ProtectedHeader []byte
	Payload         []byte
	Signature       []byte
	Detached        bool
}

// EncodeCOSESign1 writes a deterministic COSE_Sign1 object using Ed25519 over
// the standard COSE Sig_structure for either embedded or detached payloads.
func EncodeCOSESign1(payload []byte, signer string, detached bool) ([]byte, error) {
	protectedHeader, protectedErr := encodeCOSEProtectedHeader()
	if protectedErr != nil {
		return nil, protectedErr
	}
	toSign, signErr := encodeCOSESignatureStructure(protectedHeader, payload)
	if signErr != nil {
		return nil, signErr
	}
	signature := ed25519.Sign(DeterministicPrivateKey(signer), toSign)
	writer := &cborWriter{}
	if err := writer.writeTag(coseSign1Tag); err != nil {
		return nil, err
	}
	if err := writer.writeArrayHeader(4); err != nil {
		return nil, err
	}
	if err := writer.writeBytes(protectedHeader); err != nil {
		return nil, err
	}
	if err := writer.writeMapHeader(0); err != nil {
		return nil, err
	}
	if detached {
		if err := writer.writeNull(); err != nil {
			return nil, err
		}
	} else if err := writer.writeBytes(payload); err != nil {
		return nil, err
	}
	if err := writer.writeBytes(signature); err != nil {
		return nil, err
	}
	return writer.buffer.Bytes(), nil
}

// VerifyCOSESign1 checks the POC15 COSE_Sign1 subset and requires protected
// header alg=-8 so algorithm choice is explicit in signed bytes.
func VerifyCOSESign1(coseBytes []byte, detachedPayload []byte, signer string) error {
	coseSign1, parseErr := parseCOSESign1(coseBytes)
	if parseErr != nil {
		return parseErr
	}
	if headerErr := verifyCOSEProtectedHeader(coseSign1.ProtectedHeader); headerErr != nil {
		return headerErr
	}
	payload := coseSign1.Payload
	if coseSign1.Detached {
		if len(detachedPayload) == 0 {
			return fmt.Errorf("detached COSE_Sign1 needs payload bytes")
		}
		payload = detachedPayload
	}
	toSign, signErr := encodeCOSESignatureStructure(coseSign1.ProtectedHeader, payload)
	if signErr != nil {
		return signErr
	}
	if !ed25519.Verify(DeterministicPublicKey(signer), toSign, coseSign1.Signature) {
		return fmt.Errorf("COSE_Sign1 signature failed")
	}
	return nil
}

func encodeCOSEProtectedHeader() ([]byte, error) {
	writer := &cborWriter{}
	if err := writer.writeMapHeader(1); err != nil {
		return nil, err
	}
	if err := writer.writeSignedInt(coseHeaderAlg); err != nil {
		return nil, err
	}
	if err := writer.writeSignedInt(coseAlgorithmEdDSA); err != nil {
		return nil, err
	}
	return writer.buffer.Bytes(), nil
}

func verifyCOSEProtectedHeader(protectedHeader []byte) error {
	reader := &cborReader{data: protectedHeader}
	length, lengthErr := reader.readTypeAndLength(5)
	if lengthErr != nil {
		return lengthErr
	}
	if length != 1 {
		return fmt.Errorf("COSE protected header must have one entry, got %d", length)
	}
	key, keyErr := reader.readSignedInt()
	if keyErr != nil {
		return keyErr
	}
	value, valueErr := reader.readSignedInt()
	if valueErr != nil {
		return valueErr
	}
	if key != coseHeaderAlg || value != coseAlgorithmEdDSA {
		return fmt.Errorf("COSE protected header alg=%d:%d, want 1:-8", key, value)
	}
	if reader.offset != len(reader.data) {
		return fmt.Errorf("trailing cbor bytes in COSE protected header: %d", len(reader.data)-reader.offset)
	}
	return nil
}

func encodeCOSESignatureStructure(protectedHeader, payload []byte) ([]byte, error) {
	writer := &cborWriter{}
	if err := writer.writeArrayHeader(4); err != nil {
		return nil, err
	}
	if err := writer.writeString("Signature1"); err != nil {
		return nil, err
	}
	if err := writer.writeBytes(protectedHeader); err != nil {
		return nil, err
	}
	if err := writer.writeBytes(nil); err != nil {
		return nil, err
	}
	if err := writer.writeBytes(payload); err != nil {
		return nil, err
	}
	return writer.buffer.Bytes(), nil
}

func parseCOSESign1(coseBytes []byte) (COSESign1, error) {
	reader := &cborReader{data: coseBytes}
	tag, tagErr := reader.readTypeAndLength(6)
	if tagErr != nil {
		return COSESign1{}, tagErr
	}
	if tag != coseSign1Tag {
		return COSESign1{}, fmt.Errorf("COSE_Sign1 must use tag 18, got %d", tag)
	}
	arrayLength, arrayErr := reader.readTypeAndLength(4)
	if arrayErr != nil {
		return COSESign1{}, arrayErr
	}
	if arrayLength != 4 {
		return COSESign1{}, fmt.Errorf("COSE_Sign1 array must have 4 slots, got %d", arrayLength)
	}
	protectedHeader, protectedErr := reader.readBytes()
	if protectedErr != nil {
		return COSESign1{}, protectedErr
	}
	unprotectedLength, unprotectedErr := reader.readTypeAndLength(5)
	if unprotectedErr != nil {
		return COSESign1{}, unprotectedErr
	}
	if unprotectedLength != 0 {
		return COSESign1{}, fmt.Errorf("POC15 COSE_Sign1 unprotected header must be empty")
	}
	payload, detached, payloadErr := readCOSEPayloadSlot(reader)
	if payloadErr != nil {
		return COSESign1{}, payloadErr
	}
	signature, signatureErr := reader.readBytes()
	if signatureErr != nil {
		return COSESign1{}, signatureErr
	}
	if reader.offset != len(reader.data) {
		return COSESign1{}, fmt.Errorf("trailing cbor bytes in COSE_Sign1: %d", len(reader.data)-reader.offset)
	}
	return COSESign1{
		ProtectedHeader: protectedHeader,
		Payload:         payload,
		Signature:       signature,
		Detached:        detached,
	}, nil
}

func readCOSEPayloadSlot(reader *cborReader) ([]byte, bool, error) {
	if reader.offset >= len(reader.data) {
		return nil, false, fmt.Errorf("missing COSE payload slot")
	}
	if reader.data[reader.offset] == 0xf6 {
		reader.offset++
		return nil, true, nil
	}
	payload, payloadErr := reader.readBytes()
	return payload, false, payloadErr
}
