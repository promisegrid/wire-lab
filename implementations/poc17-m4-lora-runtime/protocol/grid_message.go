package protocol

import "fmt"

const (
	PCIDDeviceStatus = "device_status_v1"
	PCIDLoRaLink     = "lora_link_v1"
	PCIDPeerStorage  = "peer_storage_v1"
)

// Message is a decoded PromiseGrid envelope for the first POC17 behavior slice.
type Message struct {
	PCID    string
	Payload []byte
	Proof   []byte
}

// Build creates a pCID-selected grid envelope.
func Build(msg Message) ([]byte, error) {
	if msg.PCID == "" {
		return nil, fmt.Errorf("pcid must be set")
	}
	slots := []Item{TagItem(CIDTag, TextItem(msg.PCID)), BytesItem(msg.Payload)}
	if msg.Proof != nil {
		slots = append(slots, BytesItem(msg.Proof))
	}
	return Encode(TagItem(GridTag, ArrayItem(slots...)))
}

// Parse validates the outer grid shape and returns the pCID-owned payload bytes.
func Parse(data []byte) (Message, error) {
	item, err := Decode(data)
	if err != nil {
		return Message{}, err
	}
	if item.Tag == nil || item.Tag.Number != GridTag {
		return Message{}, fmt.Errorf("missing grid tag")
	}
	slots := item.Tag.Value.Array
	if len(slots) != 2 && len(slots) != 3 {
		return Message{}, fmt.Errorf("unsupported grid slot count %d", len(slots))
	}
	pcidTag := slots[0].Tag
	if pcidTag == nil || pcidTag.Number != CIDTag || pcidTag.Value.Text == nil {
		return Message{}, fmt.Errorf("slot 0 must be tag 42 text pCID")
	}
	if slots[1].Bytes == nil {
		return Message{}, fmt.Errorf("slot 1 must be payload bytes")
	}
	msg := Message{
		PCID:    *pcidTag.Value.Text,
		Payload: append([]byte(nil), slots[1].Bytes...),
	}
	if len(slots) == 3 {
		if slots[2].Bytes == nil {
			return Message{}, fmt.Errorf("slot 2 must be proof bytes")
		}
		msg.Proof = append([]byte(nil), slots[2].Bytes...)
	}
	return msg, nil
}

// BuildStatusPayload keeps the small-device payload positional and compact.
func BuildStatusPayload(deviceID, status string, batteryPercent uint64, parents []string) ([]byte, error) {
	parentItems := make([]Item, 0, len(parents))
	for _, parent := range parents {
		parentItems = append(parentItems, TextItem(parent))
	}
	return Encode(ArrayItem(TextItem(deviceID), TextItem(status), UintItem(batteryPercent), ArrayItem(parentItems...)))
}

// ParseStatusPayload decodes the POC17 device_status_v1 positional payload.
func ParseStatusPayload(data []byte) (deviceID string, status string, batteryPercent uint64, parents []string, err error) {
	item, err := Decode(data)
	if err != nil {
		return "", "", 0, nil, err
	}
	if len(item.Array) != 4 || item.Array[0].Text == nil || item.Array[1].Text == nil || item.Array[2].Uint == nil {
		return "", "", 0, nil, fmt.Errorf("invalid device status payload")
	}
	parentItems := item.Array[3].Array
	parents = make([]string, 0, len(parentItems))
	for _, parent := range parentItems {
		if parent.Text == nil {
			return "", "", 0, nil, fmt.Errorf("invalid parent link")
		}
		parents = append(parents, *parent.Text)
	}
	return *item.Array[0].Text, *item.Array[1].Text, *item.Array[2].Uint, parents, nil
}
