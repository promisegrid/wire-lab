package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sort"
)

// cborWriter implements the tiny deterministic CBOR subset needed by poc2:
// arrays, tag 42, byte strings, text strings, and string maps.
//
// Intent: Avoid an external dependency while still putting real CBOR bytes on
// every app/kernel and kernel/kernel boundary. Source: DI-ratij; DI-tijat.
type cborWriter struct {
	buffer bytes.Buffer
}

func (writer *cborWriter) writeArrayHeader(length int) error {
	return writer.writeTypeAndLength(4, uint64(length))
}

func (writer *cborWriter) writeMapHeader(length int) error {
	return writer.writeTypeAndLength(5, uint64(length))
}

func (writer *cborWriter) writeTag(tagNumber uint64) error {
	return writer.writeTypeAndLength(6, tagNumber)
}

func (writer *cborWriter) writeBytes(byteString []byte) error {
	if err := writer.writeTypeAndLength(2, uint64(len(byteString))); err != nil {
		return err
	}
	_, writeErr := writer.buffer.Write(byteString)
	return writeErr
}

func (writer *cborWriter) writeText(text string) error {
	if err := writer.writeTypeAndLength(3, uint64(len([]byte(text)))); err != nil {
		return err
	}
	_, writeErr := writer.buffer.WriteString(text)
	return writeErr
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
		var lengthBytes [2]byte
		binary.BigEndian.PutUint16(lengthBytes[:], uint16(length))
		_, writeErr := writer.buffer.Write(lengthBytes[:])
		return writeErr
	case length <= 0xffffffff:
		if err := writer.buffer.WriteByte(prefix | 26); err != nil {
			return err
		}
		var lengthBytes [4]byte
		binary.BigEndian.PutUint32(lengthBytes[:], uint32(length))
		_, writeErr := writer.buffer.Write(lengthBytes[:])
		return writeErr
	default:
		return fmt.Errorf("poc2 cbor length too large: %d", length)
	}
}

func marshalStringMap(fields map[string]string) ([]byte, error) {
	writer := &cborWriter{}
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	if err := writer.writeMapHeader(len(keys)); err != nil {
		return nil, err
	}
	for _, key := range keys {
		if err := writer.writeText(key); err != nil {
			return nil, err
		}
		if err := writer.writeText(fields[key]); err != nil {
			return nil, err
		}
	}
	return writer.buffer.Bytes(), nil
}

type cborReader struct {
	data   []byte
	offset int
}

func (reader *cborReader) readTypeAndLength(expectedMajor byte) (uint64, error) {
	if reader.offset >= len(reader.data) {
		return 0, fmt.Errorf("unexpected end of cbor")
	}
	initial := reader.data[reader.offset]
	reader.offset++
	major := initial >> 5
	additional := initial & 0x1f
	if major != expectedMajor {
		return 0, fmt.Errorf("unexpected cbor major type %d, wanted %d", major, expectedMajor)
	}
	switch additional {
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23:
		return uint64(additional), nil
	case 24:
		if reader.offset+1 > len(reader.data) {
			return 0, fmt.Errorf("truncated cbor uint8 length")
		}
		length := uint64(reader.data[reader.offset])
		reader.offset++
		return length, nil
	case 25:
		if reader.offset+2 > len(reader.data) {
			return 0, fmt.Errorf("truncated cbor uint16 length")
		}
		length := uint64(binary.BigEndian.Uint16(reader.data[reader.offset:]))
		reader.offset += 2
		return length, nil
	case 26:
		if reader.offset+4 > len(reader.data) {
			return 0, fmt.Errorf("truncated cbor uint32 length")
		}
		length := uint64(binary.BigEndian.Uint32(reader.data[reader.offset:]))
		reader.offset += 4
		return length, nil
	default:
		return 0, fmt.Errorf("unsupported cbor additional info %d", additional)
	}
}

func (reader *cborReader) readBytes() ([]byte, error) {
	length, err := reader.readTypeAndLength(2)
	if err != nil {
		return nil, err
	}
	if reader.offset+int(length) > len(reader.data) {
		return nil, fmt.Errorf("truncated cbor byte string")
	}
	byteString := reader.data[reader.offset : reader.offset+int(length)]
	reader.offset += int(length)
	copiedBytes := make([]byte, len(byteString))
	copy(copiedBytes, byteString)
	return copiedBytes, nil
}

func (reader *cborReader) readText() (string, error) {
	length, err := reader.readTypeAndLength(3)
	if err != nil {
		return "", err
	}
	if reader.offset+int(length) > len(reader.data) {
		return "", fmt.Errorf("truncated cbor text string")
	}
	text := string(reader.data[reader.offset : reader.offset+int(length)])
	reader.offset += int(length)
	return text, nil
}

func unmarshalStringMap(payloadBytes []byte) (map[string]string, error) {
	reader := &cborReader{data: payloadBytes}
	length, err := reader.readTypeAndLength(5)
	if err != nil {
		return nil, err
	}
	fields := make(map[string]string, length)
	for entryIndex := uint64(0); entryIndex < length; entryIndex++ {
		key, keyErr := reader.readText()
		if keyErr != nil {
			return nil, keyErr
		}
		value, valueErr := reader.readText()
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
