package carbundle

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"sort"

	cidlib "github.com/ipfs/go-cid"
	carv2 "github.com/ipld/go-car/v2"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
)

// Block records one exact CAR block and its original local CAS kind.
type Block struct {
	CID  string `json:"cid"`
	Kind string `json:"kind"`
	Size int64  `json:"size"`
}

// Encode writes a CARv1 bundle containing the exact bytes for the supplied CIDs.
//
// Intent: object_bytes promises carry a standard CARv1 package, not ad-hoc byte
// arrays. The provider still promises only what it has in its local sparse CAS,
// and the receiver verifies each CID before storing. Source: DI-koriz
func Encode(cas *store.FileStore, roots []cidlib.Cid) ([]byte, []Block, error) {
	if cas == nil {
		return nil, nil, fmt.Errorf("CAS is required")
	}
	if len(roots) == 0 {
		return nil, nil, fmt.Errorf("CAR must have at least one root")
	}
	var buffer bytes.Buffer
	if err := writeHeader(&buffer, roots); err != nil {
		return nil, nil, err
	}
	blocks := make([]Block, 0, len(roots))
	for _, root := range roots {
		content, entry, getErr := cas.Get(root)
		if getErr != nil {
			return nil, nil, getErr
		}
		if err := writeBlock(&buffer, root, content); err != nil {
			return nil, nil, err
		}
		blocks = append(blocks, Block{CID: store.CIDText(root), Kind: entry.Kind, Size: entry.Size})
	}
	return buffer.Bytes(), blocks, nil
}

// DecodeAndStore verifies every CAR block by CID and stores exact bytes locally.
func DecodeAndStore(cas *store.FileStore, carBytes []byte, kinds map[string]string) ([]Block, error) {
	if cas == nil {
		return nil, fmt.Errorf("CAS is required")
	}
	// Intent: POC18 receivers keep their local exact-CID check, but also prove
	// that the transferred package is a standard-readable CAR stream instead of
	// an accidental private framing format. Source: DI-biruf
	if verifyErr := VerifyStandard(carBytes, nil, mapKeys(kinds)); verifyErr != nil {
		return nil, verifyErr
	}
	reader := bytes.NewReader(carBytes)
	if err := readHeader(reader); err != nil {
		return nil, err
	}
	blocks := []Block{}
	for {
		blockCID, rawData, nextErr := readBlock(reader)
		if nextErr == io.EOF {
			break
		}
		if nextErr != nil {
			return nil, nextErr
		}
		if !store.CIDForBytes(rawData).Equals(blockCID) {
			return nil, fmt.Errorf("CAR block bytes do not match CID %s", store.CIDText(blockCID))
		}
		cidText := store.CIDText(blockCID)
		kind := kinds[cidText]
		if kind == "" {
			kind = "message"
		}
		entry, putErr := cas.Put(kind, rawData)
		if putErr != nil {
			return nil, putErr
		}
		if entry.CID != cidText {
			return nil, fmt.Errorf("CAR block CID mismatch: stored %s wanted %s", entry.CID, cidText)
		}
		blocks = append(blocks, Block{CID: entry.CID, Kind: entry.Kind, Size: entry.Size})
	}
	return blocks, nil
}

// VerifyStandard asks the Go CAR library to parse and hash-check the supplied
// bytes, then compares optional expected roots and block CIDs.
//
// Intent: The observer/analyzer should be able to prove collected object_bytes
// payloads are real CARv1 packages that ordinary IPLD tooling can inspect, while
// PromiseGrid still treats the enclosing grid message as the actual promise.
// Source: DI-biruf
func VerifyStandard(carBytes []byte, expectedRoots []string, expectedBlocks []string) error {
	reader, readerErr := carv2.NewReader(bytes.NewReader(carBytes))
	if readerErr != nil {
		return readerErr
	}
	if reader.Version != 1 {
		return fmt.Errorf("CAR version=%d, want 1", reader.Version)
	}
	roots, rootsErr := reader.Roots()
	if rootsErr != nil {
		return rootsErr
	}
	if len(expectedRoots) > 0 {
		actualRoots := make([]string, 0, len(roots))
		for _, root := range roots {
			actualRoots = append(actualRoots, store.CIDText(root))
		}
		if !sameStringSet(actualRoots, expectedRoots) {
			return fmt.Errorf("CAR roots %v do not match expected %v", actualRoots, expectedRoots)
		}
	}
	if _, inspectErr := reader.Inspect(true); inspectErr != nil {
		return inspectErr
	}
	if len(expectedBlocks) == 0 {
		return nil
	}
	actualBlocks, actualErr := blockCIDs(carBytes)
	if actualErr != nil {
		return actualErr
	}
	if !sameStringSet(actualBlocks, expectedBlocks) {
		return fmt.Errorf("CAR blocks %v do not match expected %v", actualBlocks, expectedBlocks)
	}
	return nil
}

func writeHeader(writer *bytes.Buffer, roots []cidlib.Cid) error {
	rootValues := make([]any, 0, len(roots))
	for _, root := range roots {
		rootValues = append(rootValues, store.LinkTag(root))
	}
	headerBytes, marshalErr := store.MarshalCBOR(map[string]any{
		"roots":   rootValues,
		"version": uint64(1),
	})
	if marshalErr != nil {
		return marshalErr
	}
	return writeSection(writer, headerBytes)
}

func readHeader(reader *bytes.Reader) error {
	headerBytes, readErr := readSection(reader)
	if readErr != nil {
		return readErr
	}
	var header map[string]any
	if unmarshalErr := store.UnmarshalCBOR(headerBytes, &header); unmarshalErr != nil {
		return unmarshalErr
	}
	version, ok := header["version"].(uint64)
	if !ok || version != 1 {
		return fmt.Errorf("CAR header version must be 1")
	}
	roots, ok := header["roots"].([]any)
	if !ok || len(roots) == 0 {
		return fmt.Errorf("CAR header must include roots")
	}
	for _, rootValue := range roots {
		if _, cidErr := store.CIDFromLinkTag(rootValue); cidErr != nil {
			return cidErr
		}
	}
	return nil
}

func writeBlock(writer *bytes.Buffer, objectCID cidlib.Cid, rawData []byte) error {
	section := make([]byte, 0, len(objectCID.Bytes())+len(rawData))
	section = append(section, objectCID.Bytes()...)
	section = append(section, rawData...)
	return writeSection(writer, section)
}

func readBlock(reader *bytes.Reader) (cidlib.Cid, []byte, error) {
	section, readErr := readSection(reader)
	if readErr != nil {
		return cidlib.Undef, nil, readErr
	}
	cidLength, objectCID, cidErr := cidlib.CidFromBytes(section)
	if cidErr != nil {
		return cidlib.Undef, nil, cidErr
	}
	if cidLength >= len(section) {
		return cidlib.Undef, nil, fmt.Errorf("CAR block is missing data for CID %s", store.CIDText(objectCID))
	}
	rawData := append([]byte(nil), section[cidLength:]...)
	return objectCID, rawData, nil
}

func writeSection(writer *bytes.Buffer, content []byte) error {
	var lengthPrefix [binary.MaxVarintLen64]byte
	prefixLength := binary.PutUvarint(lengthPrefix[:], uint64(len(content)))
	if _, writeErr := writer.Write(lengthPrefix[:prefixLength]); writeErr != nil {
		return writeErr
	}
	_, writeErr := writer.Write(content)
	return writeErr
}

func readSection(reader *bytes.Reader) ([]byte, error) {
	length, lengthErr := binary.ReadUvarint(reader)
	if lengthErr != nil {
		return nil, lengthErr
	}
	if length == 0 {
		return nil, fmt.Errorf("CAR section length must be positive")
	}
	if length > uint64(reader.Len()) {
		return nil, fmt.Errorf("CAR section length %d exceeds remaining bytes %d", length, reader.Len())
	}
	section := make([]byte, int(length))
	if _, readErr := io.ReadFull(reader, section); readErr != nil {
		return nil, readErr
	}
	return section, nil
}

func blockCIDs(carBytes []byte) ([]string, error) {
	reader := bytes.NewReader(carBytes)
	if err := readHeader(reader); err != nil {
		return nil, err
	}
	cids := []string{}
	for {
		blockCID, _, readErr := readBlock(reader)
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return nil, readErr
		}
		cids = append(cids, store.CIDText(blockCID))
	}
	return cids, nil
}

func mapKeys(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sameStringSet(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	leftCopy := append([]string(nil), left...)
	rightCopy := append([]string(nil), right...)
	sort.Strings(leftCopy)
	sort.Strings(rightCopy)
	for index := range leftCopy {
		if leftCopy[index] != rightCopy[index] {
			return false
		}
	}
	return true
}
