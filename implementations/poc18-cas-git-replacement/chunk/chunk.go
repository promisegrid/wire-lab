// Package chunk stores regular file bytes as Rabin content-defined chunks plus
// PromiseGrid manifest objects.
//
// Intent: POC18 must handle large files in-band instead of adding a Git LFS-like
// side channel. Source: DI-dofoj; DI-jifuj
package chunk

import (
	"bytes"
	"fmt"
	"io"
	"os"

	cidlib "github.com/ipfs/go-cid"
	"github.com/restic/chunker"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
)

const (
	// DefaultPolynomial is restic's documented example irreducible Rabin
	// polynomial. Keeping it constant makes POC18 test output deterministic.
	DefaultPolynomial = 0x3DA3358B4DC173
	DefaultMinSize    = 2 * 1024
	DefaultMaxSize    = 32 * 1024
	DefaultAvgBits    = 13
)

// Parameters records the chunker contract used to produce a manifest.
type Parameters struct {
	Polynomial string `json:"polynomial"`
	MinSize    uint   `json:"min_size"`
	MaxSize    uint   `json:"max_size"`
	Average    int    `json:"average_bits"`
}

// ChunkRef names one retained chunk.
type ChunkRef struct {
	Offset int64  `json:"offset"`
	Length uint   `json:"length"`
	CID    string `json:"cid"`
}

// Manifest describes the ordered chunks for one file.
type Manifest struct {
	FileSize int64      `json:"file_size"`
	Chunker  string     `json:"chunker"`
	Params   Parameters `json:"params"`
	Chunks   []ChunkRef `json:"chunks"`
}

// StoredManifest is the result of chunking and retaining a regular file.
type StoredManifest struct {
	Manifest    Manifest
	ManifestCID cidlib.Cid
	ChunkCIDs   []cidlib.Cid
}

// DefaultParameters returns the deterministic first-slice Rabin configuration.
func DefaultParameters() Parameters {
	return Parameters{
		Polynomial: fmt.Sprintf("0x%x", uint64(DefaultPolynomial)),
		MinSize:    DefaultMinSize,
		MaxSize:    DefaultMaxSize,
		Average:    DefaultAvgBits,
	}
}

// StoreFile reads path and stores its chunks and manifest in cas.
func StoreFile(cas *store.FileStore, path string) (StoredManifest, error) {
	content, readErr := os.ReadFile(path)
	if readErr != nil {
		return StoredManifest{}, readErr
	}
	return StoreBytes(cas, content)
}

// StoreBytes chunks content, stores each chunk, and stores a CBOR manifest.
func StoreBytes(cas *store.FileStore, content []byte) (StoredManifest, error) {
	params := DefaultParameters()
	rabin := chunker.New(
		bytes.NewReader(content),
		chunker.Pol(DefaultPolynomial),
		chunker.WithBoundaries(params.MinSize, params.MaxSize),
		chunker.WithAverageBits(params.Average),
	)
	buffer := make([]byte, params.MaxSize)
	manifest := Manifest{FileSize: int64(len(content)), Chunker: "rabin", Params: params}
	chunkCIDs := []cidlib.Cid{}
	for {
		nextChunk, nextErr := rabin.Next(buffer)
		if nextErr == io.EOF {
			break
		}
		if nextErr != nil {
			return StoredManifest{}, nextErr
		}
		entry, putErr := cas.Put("chunk", nextChunk.Data)
		if putErr != nil {
			return StoredManifest{}, putErr
		}
		chunkCID, parseErr := store.ParseCIDText(entry.CID)
		if parseErr != nil {
			return StoredManifest{}, parseErr
		}
		manifest.Chunks = append(manifest.Chunks, ChunkRef{
			Offset: int64(nextChunk.Start),
			Length: nextChunk.Length,
			CID:    entry.CID,
		})
		chunkCIDs = append(chunkCIDs, chunkCID)
	}
	manifestEntry, _, manifestErr := cas.PutCBOR("chunk_manifest", manifest.CBOR())
	if manifestErr != nil {
		return StoredManifest{}, manifestErr
	}
	manifestCID, parseErr := store.ParseCIDText(manifestEntry.CID)
	if parseErr != nil {
		return StoredManifest{}, parseErr
	}
	return StoredManifest{Manifest: manifest, ManifestCID: manifestCID, ChunkCIDs: chunkCIDs}, nil
}

// CBOR returns the pCID-owned manifest object shape.
func (manifest Manifest) CBOR() []any {
	chunks := make([]any, 0, len(manifest.Chunks))
	for _, chunkRef := range manifest.Chunks {
		chunkCID, parseErr := store.ParseCIDText(chunkRef.CID)
		if parseErr != nil {
			panic("stored chunk CID should parse: " + parseErr.Error())
		}
		chunks = append(chunks, []any{chunkRef.Offset, chunkRef.Length, store.LinkTag(chunkCID)})
	}
	return []any{
		"poc18_chunk_manifest",
		manifest.FileSize,
		manifest.Chunker,
		[]any{manifest.Params.Polynomial, manifest.Params.MinSize, manifest.Params.MaxSize, manifest.Params.Average},
		chunks,
	}
}

// DecodeManifest reads a manifest object stored by StoreBytes.
func DecodeManifest(content []byte) (Manifest, error) {
	var raw []any
	if err := store.UnmarshalCBOR(content, &raw); err != nil {
		return Manifest{}, err
	}
	if len(raw) != 5 {
		return Manifest{}, fmt.Errorf("manifest must have 5 slots, got %d", len(raw))
	}
	label, ok := raw[0].(string)
	if !ok || label != "poc18_chunk_manifest" {
		return Manifest{}, fmt.Errorf("manifest label mismatch")
	}
	fileSize, ok := raw[1].(uint64)
	if !ok {
		return Manifest{}, fmt.Errorf("manifest file size must be uint")
	}
	chunkerName, ok := raw[2].(string)
	if !ok {
		return Manifest{}, fmt.Errorf("manifest chunker must be text")
	}
	paramsRaw, ok := raw[3].([]any)
	if !ok || len(paramsRaw) != 4 {
		return Manifest{}, fmt.Errorf("manifest params must have 4 slots")
	}
	polynomial, ok := paramsRaw[0].(string)
	if !ok {
		return Manifest{}, fmt.Errorf("manifest polynomial must be text")
	}
	minSize, minOK := paramsRaw[1].(uint64)
	maxSize, maxOK := paramsRaw[2].(uint64)
	average, avgOK := paramsRaw[3].(uint64)
	if !minOK || !maxOK || !avgOK {
		return Manifest{}, fmt.Errorf("manifest size parameters must be uint")
	}
	chunkRows, ok := raw[4].([]any)
	if !ok {
		return Manifest{}, fmt.Errorf("manifest chunks must be array")
	}
	manifest := Manifest{
		FileSize: int64(fileSize),
		Chunker:  chunkerName,
		Params: Parameters{
			Polynomial: polynomial,
			MinSize:    uint(minSize),
			MaxSize:    uint(maxSize),
			Average:    int(average),
		},
	}
	for _, rowValue := range chunkRows {
		row, rowOK := rowValue.([]any)
		if !rowOK || len(row) != 3 {
			return Manifest{}, fmt.Errorf("manifest chunk row must have 3 slots")
		}
		offset, offsetOK := row[0].(uint64)
		length, lengthOK := row[1].(uint64)
		if !offsetOK || !lengthOK {
			return Manifest{}, fmt.Errorf("manifest chunk offset and length must be uint")
		}
		chunkCID, cidErr := store.CIDFromLinkTag(row[2])
		if cidErr != nil {
			return Manifest{}, cidErr
		}
		manifest.Chunks = append(manifest.Chunks, ChunkRef{
			Offset: int64(offset),
			Length: uint(length),
			CID:    store.CIDText(chunkCID),
		})
	}
	return manifest, nil
}

// Reassemble reconstructs content from locally retained chunks.
func Reassemble(cas *store.FileStore, manifest Manifest) ([]byte, error) {
	var output bytes.Buffer
	for _, chunkRef := range manifest.Chunks {
		chunkCID, parseErr := store.ParseCIDText(chunkRef.CID)
		if parseErr != nil {
			return nil, parseErr
		}
		content, _, getErr := cas.Get(chunkCID)
		if getErr != nil {
			return nil, getErr
		}
		if uint(len(content)) != chunkRef.Length {
			return nil, fmt.Errorf("chunk %s length mismatch", chunkRef.CID)
		}
		if _, writeErr := output.Write(content); writeErr != nil {
			return nil, writeErr
		}
	}
	if int64(output.Len()) != manifest.FileSize {
		return nil, fmt.Errorf("manifest file size mismatch: got %d want %d", output.Len(), manifest.FileSize)
	}
	return output.Bytes(), nil
}
