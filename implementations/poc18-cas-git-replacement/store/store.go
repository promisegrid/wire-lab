// Package store owns POC18's exact-byte CAS and CID handling.
//
// Intent: Keep all object identities as CIDs, with binary CID bytes on the wire
// and CIDv1 base32 when rendered into filenames, logs, or JSON. Source: DI-jifuj;
// DI-harih
package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/fxamacker/cbor/v2"
	cidlib "github.com/ipfs/go-cid"
	mbase "github.com/multiformats/go-multibase"
	mh "github.com/multiformats/go-multihash"
)

const (
	// LinkTagNumber is the DAG-CBOR tag number used for IPLD links.
	LinkTagNumber = uint64(42)
)

var cborMode cbor.EncMode

func init() {
	mode, modeErr := cbor.CanonicalEncOptions().EncMode()
	if modeErr != nil {
		panic("canonical CBOR mode should be constructible: " + modeErr.Error())
	}
	cborMode = mode
}

// Entry describes one locally retained CAS object.
type Entry struct {
	CID  string `json:"cid"`
	Kind string `json:"kind"`
	Path string `json:"path"`
	Size int64  `json:"size"`
}

// FileStore is a sparse local CAS. It is not a global repository and does not
// promise completeness.
type FileStore struct {
	Root string
}

// Open creates or opens a sparse filesystem CAS rooted at root.
func Open(root string) (*FileStore, error) {
	if root == "" {
		return nil, fmt.Errorf("store root is required")
	}
	for _, dir := range []string{
		filepath.Join(root, "objects"),
		filepath.Join(root, "chunks"),
	} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	return &FileStore{Root: root}, nil
}

// CIDForBytes returns the CIDv1 raw sha2-256 identifier for exact object bytes.
func CIDForBytes(content []byte) cidlib.Cid {
	multihash, multihashErr := mh.Sum(content, mh.SHA2_256, -1)
	if multihashErr != nil {
		panic("sha2-256 multihash should not fail: " + multihashErr.Error())
	}
	return cidlib.NewCidV1(cidlib.Raw, multihash)
}

// CIDText renders the canonical CIDv1 base32 text form.
func CIDText(value cidlib.Cid) string {
	text, textErr := value.StringOfBase(mbase.Base32)
	if textErr != nil {
		panic("cid base32 rendering should not fail: " + textErr.Error())
	}
	return text
}

// ParseCIDText validates a printable CIDv1 raw sha2-256 identifier.
func ParseCIDText(cidText string) (cidlib.Cid, error) {
	parsedCID, decodeErr := cidlib.Decode(cidText)
	if decodeErr != nil {
		return cidlib.Undef, decodeErr
	}
	if validateErr := validateSupportedCIDProfile(parsedCID); validateErr != nil {
		return cidlib.Undef, validateErr
	}
	if CIDText(parsedCID) != cidText {
		return cidlib.Undef, fmt.Errorf("cid text must be canonical base32")
	}
	return parsedCID, nil
}

// ParseCIDBytes validates binary CIDv1 raw sha2-256 bytes.
func ParseCIDBytes(cidBytes []byte) (cidlib.Cid, error) {
	parsedCID, castErr := cidlib.Cast(cidBytes)
	if castErr != nil {
		return cidlib.Undef, castErr
	}
	if validateErr := validateSupportedCIDProfile(parsedCID); validateErr != nil {
		return cidlib.Undef, validateErr
	}
	if !bytes.Equal(parsedCID.Bytes(), cidBytes) {
		return cidlib.Undef, fmt.Errorf("cid bytes are not canonical")
	}
	return parsedCID, nil
}

func validateSupportedCIDProfile(value cidlib.Cid) error {
	if !value.Defined() {
		return fmt.Errorf("cid is undefined")
	}
	if value.Version() != 1 {
		return fmt.Errorf("cid must be CIDv1, got v%d", value.Version())
	}
	if value.Type() != cidlib.Raw {
		return fmt.Errorf("cid codec must be raw, got 0x%x", value.Type())
	}
	decodedMultihash, decodeErr := mh.Decode(value.Hash())
	if decodeErr != nil {
		return decodeErr
	}
	if decodedMultihash.Code != mh.SHA2_256 {
		return fmt.Errorf("cid multihash must be sha2-256, got 0x%x", decodedMultihash.Code)
	}
	if decodedMultihash.Length != 32 {
		return fmt.Errorf("cid sha2-256 length must be 32 bytes, got %d", decodedMultihash.Length)
	}
	return nil
}

// MarshalCBOR writes canonical CBOR so the same promise object has one stable
// CID across process runs.
func MarshalCBOR(value any) ([]byte, error) {
	return cborMode.Marshal(value)
}

// UnmarshalCBOR decodes a CBOR object into the caller's target.
func UnmarshalCBOR(data []byte, target any) error {
	return cbor.Unmarshal(data, target)
}

// LinkTag wraps a CID as a DAG-CBOR tag 42 link payload.
func LinkTag(value cidlib.Cid) cbor.Tag {
	tagBytes := make([]byte, 0, len(value.Bytes())+1)
	tagBytes = append(tagBytes, 0x00)
	tagBytes = append(tagBytes, value.Bytes()...)
	return cbor.Tag{Number: LinkTagNumber, Content: tagBytes}
}

// CIDFromLinkTag validates a decoded DAG-CBOR tag 42 link.
func CIDFromLinkTag(value any) (cidlib.Cid, error) {
	tag, ok := value.(cbor.Tag)
	if !ok {
		return cidlib.Undef, fmt.Errorf("value is not cbor tag")
	}
	if tag.Number != LinkTagNumber {
		return cidlib.Undef, fmt.Errorf("expected tag 42, got tag %d", tag.Number)
	}
	tagBytes, ok := tag.Content.([]byte)
	if !ok {
		return cidlib.Undef, fmt.Errorf("tag 42 content must be bytes")
	}
	if len(tagBytes) < 2 || tagBytes[0] != 0x00 {
		return cidlib.Undef, fmt.Errorf("tag 42 content must start with DAG-CBOR CID sentinel")
	}
	return ParseCIDBytes(tagBytes[1:])
}

// Put stores exact bytes and records a local sparse-CAS index entry.
func (fileStore *FileStore) Put(kind string, content []byte) (Entry, error) {
	if kind == "" {
		return Entry{}, fmt.Errorf("object kind is required")
	}
	objectCID := CIDForBytes(content)
	cidText := CIDText(objectCID)
	path := fileStore.pathFor(kind, cidText)
	if existing, readErr := os.ReadFile(path); readErr == nil {
		if !bytes.Equal(existing, content) {
			return Entry{}, fmt.Errorf("existing CAS bytes for %s do not match CID content", cidText)
		}
		return fileStore.recordEntry(cidText, kind, path, int64(len(content)))
	} else if !os.IsNotExist(readErr) {
		return Entry{}, readErr
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return Entry{}, err
	}
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, content, 0o644); err != nil {
		return Entry{}, err
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return Entry{}, err
	}
	return fileStore.recordEntry(cidText, kind, path, int64(len(content)))
}

// PutCBOR stores a canonical CBOR object under its CID.
func (fileStore *FileStore) PutCBOR(kind string, value any) (Entry, []byte, error) {
	content, marshalErr := MarshalCBOR(value)
	if marshalErr != nil {
		return Entry{}, nil, marshalErr
	}
	entry, putErr := fileStore.Put(kind, content)
	if putErr != nil {
		return Entry{}, nil, putErr
	}
	return entry, content, nil
}

// Get returns exact bytes for a locally retained CID.
func (fileStore *FileStore) Get(objectCID cidlib.Cid) ([]byte, Entry, error) {
	cidText := CIDText(objectCID)
	index, indexErr := fileStore.loadIndex()
	if indexErr != nil {
		return nil, Entry{}, indexErr
	}
	if entry, ok := index[cidText]; ok {
		content, readErr := os.ReadFile(entry.Path)
		if readErr != nil {
			return nil, Entry{}, readErr
		}
		if !CIDForBytes(content).Equals(objectCID) {
			return nil, Entry{}, fmt.Errorf("CAS path %s does not match CID %s", entry.Path, cidText)
		}
		return content, entry, nil
	}
	for _, path := range fileStore.probePaths(cidText) {
		content, readErr := os.ReadFile(path)
		if readErr == nil {
			if !CIDForBytes(content).Equals(objectCID) {
				return nil, Entry{}, fmt.Errorf("CAS path %s does not match CID %s", path, cidText)
			}
			entry := Entry{CID: cidText, Kind: kindFromPath(path), Path: path, Size: int64(len(content))}
			return content, entry, nil
		}
		if !os.IsNotExist(readErr) {
			return nil, Entry{}, readErr
		}
	}
	return nil, Entry{}, fmt.Errorf("missing local CAS object %s", cidText)
}

// Has reports whether exact bytes for objectCID are present locally.
func (fileStore *FileStore) Has(objectCID cidlib.Cid) bool {
	_, _, getErr := fileStore.Get(objectCID)
	return getErr == nil
}

// List returns the local sparse-CAS index sorted by CID text.
func (fileStore *FileStore) List() ([]Entry, error) {
	index, indexErr := fileStore.loadIndex()
	if indexErr != nil {
		return nil, indexErr
	}
	entries := make([]Entry, 0, len(index))
	for _, entry := range index {
		entries = append(entries, entry)
	}
	sort.Slice(entries, func(left, right int) bool {
		return entries[left].CID < entries[right].CID
	})
	return entries, nil
}

// Delete removes one locally retained CAS object by CID.
//
// Intent: POC18 retention and GC are local resource decisions. Deleting from one
// sparse file store never asserts that the object should disappear from any peer
// store or from the global DAG. Source: DI-mivur
func (fileStore *FileStore) Delete(objectCID cidlib.Cid) (Entry, error) {
	cidText := CIDText(objectCID)
	index, indexErr := fileStore.loadIndex()
	if indexErr != nil {
		return Entry{}, indexErr
	}
	entry, found := index[cidText]
	if !found {
		for _, path := range fileStore.probePaths(cidText) {
			if _, statErr := os.Stat(path); statErr == nil {
				entry = Entry{CID: cidText, Kind: kindFromPath(path), Path: path}
				found = true
				break
			} else if !os.IsNotExist(statErr) {
				return Entry{}, statErr
			}
		}
	}
	if !found {
		return Entry{}, fmt.Errorf("missing local CAS object %s", cidText)
	}
	if readContent, readErr := os.ReadFile(entry.Path); readErr == nil {
		if !CIDForBytes(readContent).Equals(objectCID) {
			return Entry{}, fmt.Errorf("CAS path %s does not match CID %s", entry.Path, cidText)
		}
		entry.Size = int64(len(readContent))
	} else if !os.IsNotExist(readErr) {
		return Entry{}, readErr
	}
	if removeErr := os.Remove(entry.Path); removeErr != nil && !os.IsNotExist(removeErr) {
		return Entry{}, removeErr
	}
	delete(index, cidText)
	if saveErr := fileStore.saveIndex(index); saveErr != nil {
		return Entry{}, saveErr
	}
	return entry, nil
}

func (fileStore *FileStore) recordEntry(cidText, kind, path string, size int64) (Entry, error) {
	index, indexErr := fileStore.loadIndex()
	if indexErr != nil {
		return Entry{}, indexErr
	}
	entry := Entry{CID: cidText, Kind: kind, Path: path, Size: size}
	index[cidText] = entry
	if saveErr := fileStore.saveIndex(index); saveErr != nil {
		return Entry{}, saveErr
	}
	return entry, nil
}

func (fileStore *FileStore) indexPath() string {
	return filepath.Join(fileStore.Root, "index.json")
}

func (fileStore *FileStore) loadIndex() (map[string]Entry, error) {
	index := map[string]Entry{}
	content, readErr := os.ReadFile(fileStore.indexPath())
	if os.IsNotExist(readErr) {
		return index, nil
	}
	if readErr != nil {
		return nil, readErr
	}
	if len(content) == 0 {
		return index, nil
	}
	if err := json.Unmarshal(content, &index); err != nil {
		return nil, err
	}
	return index, nil
}

func (fileStore *FileStore) saveIndex(index map[string]Entry) error {
	content, marshalErr := json.MarshalIndent(index, "", "  ")
	if marshalErr != nil {
		return marshalErr
	}
	content = append(content, '\n')
	tmpPath := fileStore.indexPath() + ".tmp"
	if err := os.WriteFile(tmpPath, content, 0o644); err != nil {
		return err
	}
	return os.Rename(tmpPath, fileStore.indexPath())
}

func (fileStore *FileStore) pathFor(kind, cidText string) string {
	if kind == "chunk" {
		return filepath.Join(fileStore.Root, "chunks", cidText+".bin")
	}
	return filepath.Join(fileStore.Root, "objects", cidText+".cbor")
}

func (fileStore *FileStore) probePaths(cidText string) []string {
	return []string{
		filepath.Join(fileStore.Root, "objects", cidText+".cbor"),
		filepath.Join(fileStore.Root, "chunks", cidText+".bin"),
	}
}

func kindFromPath(path string) string {
	if filepath.Base(filepath.Dir(path)) == "chunks" {
		return "chunk"
	}
	return "object"
}
