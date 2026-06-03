package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sort"
	"unicode/utf8"
)

// cborWriter writes the deterministic CBOR subset POC10 needs.
// Intent: Exercise real byte-level CBOR grid envelopes without making this POC
// depend on a broad codec stack. Source: DI-pijan
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

// cborReader reads the deterministic CBOR subset POC10 writes.
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
	value := make([]byte, int(length))
	copy(value, reader.data[reader.offset:reader.offset+int(length)])
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

// UnmarshalStringMap decodes a deterministic CBOR map of string fields.
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
