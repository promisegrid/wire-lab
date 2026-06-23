package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	// GridTag is the explicit PromiseGrid outer CBOR tag.
	GridTag uint64 = 1735551332
	// CIDTag is the tag used for the slot-0 protocol selector.
	CIDTag uint64 = 42
)

// Item is the small CBOR value subset needed by the POC17 behavior simulator.
type Item struct {
	Uint  *uint64
	Text  *string
	Bytes []byte
	Array []Item
	Tag   *Tag
}

// Tag represents a tagged CBOR value.
type Tag struct {
	Number uint64
	Value  Item
}

// UintItem creates an unsigned integer CBOR value.
func UintItem(v uint64) Item { return Item{Uint: &v} }

// TextItem creates a text-string CBOR value.
func TextItem(v string) Item { return Item{Text: &v} }

// BytesItem creates a byte-string CBOR value.
func BytesItem(v []byte) Item {
	cp := append([]byte(nil), v...)
	return Item{Bytes: cp}
}

// ArrayItem creates an array CBOR value.
func ArrayItem(v ...Item) Item { return Item{Array: append([]Item(nil), v...)} }

// TagItem creates a tagged CBOR value.
func TagItem(number uint64, value Item) Item { return Item{Tag: &Tag{Number: number, Value: value}} }

// Encode serializes the supported CBOR subset.
func Encode(item Item) ([]byte, error) {
	var buf bytes.Buffer
	if err := writeItem(&buf, item); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode parses the supported CBOR subset and rejects trailing bytes.
func Decode(data []byte) (Item, error) {
	r := bytes.NewReader(data)
	item, err := readItem(r)
	if err != nil {
		return Item{}, err
	}
	if r.Len() != 0 {
		return Item{}, fmt.Errorf("trailing CBOR bytes: %d", r.Len())
	}
	return item, nil
}

func writeItem(w io.Writer, item Item) error {
	switch {
	case item.Uint != nil:
		return writeTypeValue(w, 0, *item.Uint)
	case item.Text != nil:
		data := []byte(*item.Text)
		if err := writeTypeValue(w, 3, uint64(len(data))); err != nil {
			return err
		}
		_, err := w.Write(data)
		return err
	case item.Bytes != nil:
		if err := writeTypeValue(w, 2, uint64(len(item.Bytes))); err != nil {
			return err
		}
		_, err := w.Write(item.Bytes)
		return err
	case item.Array != nil:
		if err := writeTypeValue(w, 4, uint64(len(item.Array))); err != nil {
			return err
		}
		for _, child := range item.Array {
			if err := writeItem(w, child); err != nil {
				return err
			}
		}
		return nil
	case item.Tag != nil:
		if err := writeTypeValue(w, 6, item.Tag.Number); err != nil {
			return err
		}
		return writeItem(w, item.Tag.Value)
	default:
		return fmt.Errorf("empty CBOR item")
	}
}

func writeTypeValue(w io.Writer, major byte, value uint64) error {
	prefix := major << 5
	switch {
	case value < 24:
		_, err := w.Write([]byte{prefix | byte(value)})
		return err
	case value <= 0xff:
		_, err := w.Write([]byte{prefix | 24, byte(value)})
		return err
	case value <= 0xffff:
		var tmp [3]byte
		tmp[0] = prefix | 25
		binary.BigEndian.PutUint16(tmp[1:], uint16(value))
		_, err := w.Write(tmp[:])
		return err
	case value <= 0xffffffff:
		var tmp [5]byte
		tmp[0] = prefix | 26
		binary.BigEndian.PutUint32(tmp[1:], uint32(value))
		_, err := w.Write(tmp[:])
		return err
	default:
		var tmp [9]byte
		tmp[0] = prefix | 27
		binary.BigEndian.PutUint64(tmp[1:], value)
		_, err := w.Write(tmp[:])
		return err
	}
}

func readItem(r *bytes.Reader) (Item, error) {
	b, err := r.ReadByte()
	if err != nil {
		return Item{}, fmt.Errorf("read CBOR prefix: %w", err)
	}
	major := b >> 5
	additional := b & 0x1f
	value, err := readValue(r, additional)
	if err != nil {
		return Item{}, err
	}
	switch major {
	case 0:
		return UintItem(value), nil
	case 2:
		data, err := readBytes(r, value)
		if err != nil {
			return Item{}, err
		}
		return BytesItem(data), nil
	case 3:
		data, err := readBytes(r, value)
		if err != nil {
			return Item{}, err
		}
		text := string(data)
		return TextItem(text), nil
	case 4:
		items := make([]Item, 0, value)
		for i := uint64(0); i < value; i++ {
			child, err := readItem(r)
			if err != nil {
				return Item{}, err
			}
			items = append(items, child)
		}
		return ArrayItem(items...), nil
	case 6:
		child, err := readItem(r)
		if err != nil {
			return Item{}, err
		}
		return TagItem(value, child), nil
	default:
		return Item{}, fmt.Errorf("unsupported CBOR major type %d", major)
	}
}

func readValue(r *bytes.Reader, additional byte) (uint64, error) {
	switch {
	case additional < 24:
		return uint64(additional), nil
	case additional == 24:
		b, err := r.ReadByte()
		return uint64(b), err
	case additional == 25:
		data, err := readBytes(r, 2)
		if err != nil {
			return 0, err
		}
		return uint64(binary.BigEndian.Uint16(data)), nil
	case additional == 26:
		data, err := readBytes(r, 4)
		if err != nil {
			return 0, err
		}
		return uint64(binary.BigEndian.Uint32(data)), nil
	case additional == 27:
		data, err := readBytes(r, 8)
		if err != nil {
			return 0, err
		}
		return binary.BigEndian.Uint64(data), nil
	default:
		return 0, fmt.Errorf("unsupported CBOR additional info %d", additional)
	}
}

func readBytes(r *bytes.Reader, n uint64) ([]byte, error) {
	if n > uint64(r.Len()) {
		return nil, fmt.Errorf("short CBOR read: need %d have %d", n, r.Len())
	}
	data := make([]byte, n)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}
	return data, nil
}
