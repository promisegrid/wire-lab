package protocol

import "fmt"

const (
	PeerStorageGrant      = "grant"
	PeerStoragePut        = "put"
	PeerStoragePutResult  = "put_result"
	PeerStorageGet        = "get"
	PeerStorageGetResult  = "get_result"
	PeerStorageGetRefusal = "get_refusal"
)

// PeerStoragePayload is the decoded shape of the pCID-defined peer_storage
// tuple. Kind is inferred from tuple shape; it is not sent as a redundant slot.
type PeerStoragePayload struct {
	Kind               string
	Issuer             string
	Holder             string
	Token              []byte
	TokenCID           string
	AllowedKinds       []string
	MaxContentBytes    uint64
	MaxRetainedObjects uint64
	RetentionTerms     string
	Accepted           bool
	RelatedMessageCID  string
	ContentCID         string
	Content            []byte
	Reason             string
}

// BuildPeerStorageGrant creates Bob's compact storage capability grant.
func BuildPeerStorageGrant(payload PeerStoragePayload) ([]byte, error) {
	return Encode(ArrayItem(
		TextItem(payload.Issuer),
		TextItem(payload.Holder),
		BytesItem(payload.Token),
		textArray(payload.AllowedKinds),
		UintItem(payload.MaxContentBytes),
		UintItem(payload.MaxRetainedObjects),
		TextItem(payload.RetentionTerms),
	))
}

// BuildPeerStoragePut creates Ivan's promise to put exact bytes with Bob.
func BuildPeerStoragePut(payload PeerStoragePayload) ([]byte, error) {
	contentCID, err := cidItem(payload.ContentCID)
	if err != nil {
		return nil, err
	}
	return Encode(ArrayItem(
		TextItem(payload.Holder),
		TextItem(payload.Issuer),
		BytesItem(payload.Token),
		contentCID,
		BytesItem(payload.Content),
		TextItem(payload.Reason),
	))
}

// BuildPeerStoragePutResult creates Bob's acceptance or refusal of a put.
func BuildPeerStoragePutResult(payload PeerStoragePayload) ([]byte, error) {
	tokenCID, err := cidItem(payload.TokenCID)
	if err != nil {
		return nil, err
	}
	contentCID, err := cidItem(payload.ContentCID)
	if err != nil {
		return nil, err
	}
	accepted := uint64(0)
	if payload.Accepted {
		accepted = 1
	}
	return Encode(ArrayItem(
		TextItem(payload.Issuer),
		TextItem(payload.Holder),
		tokenCID,
		contentCID,
		UintItem(accepted),
		TextItem(payload.Reason),
	))
}

// BuildPeerStorageGet creates Ivan's promise to get exact bytes by CID.
func BuildPeerStorageGet(payload PeerStoragePayload) ([]byte, error) {
	contentCID, err := cidItem(payload.ContentCID)
	if err != nil {
		return nil, err
	}
	return Encode(ArrayItem(
		TextItem(payload.Holder),
		TextItem(payload.Issuer),
		BytesItem(payload.Token),
		contentCID,
		TextItem(payload.Reason),
	))
}

// BuildPeerStorageGetResult creates Bob's fulfillment with exact bytes.
func BuildPeerStorageGetResult(payload PeerStoragePayload) ([]byte, error) {
	tokenCID, err := cidItem(payload.TokenCID)
	if err != nil {
		return nil, err
	}
	contentCID, err := cidItem(payload.ContentCID)
	if err != nil {
		return nil, err
	}
	return Encode(ArrayItem(
		TextItem(payload.Issuer),
		TextItem(payload.Holder),
		tokenCID,
		contentCID,
		BytesItem(payload.Content),
	))
}

// BuildPeerStorageGetRefusal creates Bob's refusal to return bytes.
func BuildPeerStorageGetRefusal(payload PeerStoragePayload) ([]byte, error) {
	tokenCID, err := cidItem(payload.TokenCID)
	if err != nil {
		return nil, err
	}
	contentCID, err := cidItem(payload.ContentCID)
	if err != nil {
		return nil, err
	}
	return Encode(ArrayItem(
		TextItem(payload.Issuer),
		TextItem(payload.Holder),
		tokenCID,
		contentCID,
		TextItem(payload.Reason),
	))
}

// ParsePeerStoragePayload decodes the pCID-defined tuple and infers its kind
// from slot count and slot types instead of consuming a redundant kind slot.
func ParsePeerStoragePayload(data []byte) (PeerStoragePayload, error) {
	item, err := Decode(data)
	if err != nil {
		return PeerStoragePayload{}, err
	}
	slots := item.Array
	switch len(slots) {
	case 5:
		return parsePeerStorageFiveSlot(slots)
	case 6:
		return parsePeerStorageSixSlot(slots)
	case 7:
		return parsePeerStorageSevenSlot(slots)
	default:
		return PeerStoragePayload{}, fmt.Errorf("invalid peer_storage slot count %d", len(slots))
	}
}

func parsePeerStorageFiveSlot(slots []Item) (PeerStoragePayload, error) {
	if slots[0].Text == nil || slots[1].Text == nil {
		return PeerStoragePayload{}, fmt.Errorf("invalid five-slot peer_storage identity slots")
	}
	if slots[2].Bytes != nil && slots[4].Text != nil {
		contentCID, err := cidText(slots[3])
		if err != nil {
			return PeerStoragePayload{}, err
		}
		return PeerStoragePayload{Kind: PeerStorageGet, Holder: *slots[0].Text, Issuer: *slots[1].Text, Token: append([]byte(nil), slots[2].Bytes...), ContentCID: contentCID, Reason: *slots[4].Text}, nil
	}
	if slots[2].Tag != nil {
		tokenCID, err := cidText(slots[2])
		if err != nil {
			return PeerStoragePayload{}, err
		}
		contentCID, err := cidText(slots[3])
		if err != nil {
			return PeerStoragePayload{}, err
		}
		payload := PeerStoragePayload{Issuer: *slots[0].Text, Holder: *slots[1].Text, TokenCID: tokenCID, ContentCID: contentCID}
		if slots[4].Bytes != nil {
			payload.Kind = PeerStorageGetResult
			payload.Content = append([]byte(nil), slots[4].Bytes...)
			return payload, nil
		}
		if slots[4].Text != nil {
			payload.Kind = PeerStorageGetRefusal
			payload.Reason = *slots[4].Text
			return payload, nil
		}
	}
	return PeerStoragePayload{}, fmt.Errorf("invalid five-slot peer_storage payload")
}

func parsePeerStorageSixSlot(slots []Item) (PeerStoragePayload, error) {
	if slots[0].Text == nil || slots[1].Text == nil {
		return PeerStoragePayload{}, fmt.Errorf("invalid six-slot peer_storage identity slots")
	}
	if slots[2].Bytes != nil && slots[4].Bytes != nil && slots[5].Text != nil {
		contentCID, err := cidText(slots[3])
		if err != nil {
			return PeerStoragePayload{}, err
		}
		return PeerStoragePayload{
			Kind:       PeerStoragePut,
			Holder:     *slots[0].Text,
			Issuer:     *slots[1].Text,
			Token:      append([]byte(nil), slots[2].Bytes...),
			ContentCID: contentCID,
			Content:    append([]byte(nil), slots[4].Bytes...),
			Reason:     *slots[5].Text,
		}, nil
	}
	if slots[2].Tag != nil && slots[3].Tag != nil && slots[4].Uint != nil && slots[5].Text != nil {
		tokenCID, err := cidText(slots[2])
		if err != nil {
			return PeerStoragePayload{}, err
		}
		contentCID, err := cidText(slots[3])
		if err != nil {
			return PeerStoragePayload{}, err
		}
		return PeerStoragePayload{Kind: PeerStoragePutResult, Issuer: *slots[0].Text, Holder: *slots[1].Text, TokenCID: tokenCID, ContentCID: contentCID, Accepted: *slots[4].Uint == 1, Reason: *slots[5].Text}, nil
	}
	return PeerStoragePayload{}, fmt.Errorf("invalid six-slot peer_storage payload")
}

func parsePeerStorageSevenSlot(slots []Item) (PeerStoragePayload, error) {
	if slots[0].Text == nil || slots[1].Text == nil || slots[2].Bytes == nil || slots[4].Uint == nil || slots[5].Uint == nil || slots[6].Text == nil {
		return PeerStoragePayload{}, fmt.Errorf("invalid seven-slot peer_storage payload")
	}
	allowed, err := parseTextArray(slots[3])
	if err != nil {
		return PeerStoragePayload{}, err
	}
	return PeerStoragePayload{
		Kind:               PeerStorageGrant,
		Issuer:             *slots[0].Text,
		Holder:             *slots[1].Text,
		Token:              append([]byte(nil), slots[2].Bytes...),
		AllowedKinds:       allowed,
		MaxContentBytes:    *slots[4].Uint,
		MaxRetainedObjects: *slots[5].Uint,
		RetentionTerms:     *slots[6].Text,
	}, nil
}

func textArray(values []string) Item {
	items := make([]Item, 0, len(values))
	for _, value := range values {
		items = append(items, TextItem(value))
	}
	return ArrayItem(items...)
}

func cidItem(cid string) (Item, error) {
	data, err := Tag42DataForCIDText(cid)
	if err != nil {
		return Item{}, err
	}
	return TagItem(CIDTag, BytesItem(data)), nil
}

func cidText(item Item) (string, error) {
	if item.Tag == nil || item.Tag.Number != CIDTag || item.Tag.Value.Bytes == nil {
		return "", fmt.Errorf("expected tag-42 CID")
	}
	return CIDTextFromTag42Data(item.Tag.Value.Bytes)
}

func parseTextArray(item Item) ([]string, error) {
	values := make([]string, 0, len(item.Array))
	for _, child := range item.Array {
		if child.Text == nil {
			return nil, fmt.Errorf("invalid text array")
		}
		values = append(values, *child.Text)
	}
	return values, nil
}
