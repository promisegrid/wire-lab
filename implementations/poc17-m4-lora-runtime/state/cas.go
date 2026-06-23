package state

import (
	"crypto/sha256"
	"encoding/hex"
)

// CAS is a tiny sparse store with pressure-driven GC.
type CAS struct {
	limit   int
	order   []string
	objects map[string][]byte
}

// NewCAS creates bounded local state for the simulated M4 agent.
func NewCAS(limit int) *CAS {
	return &CAS{limit: limit, objects: make(map[string][]byte)}
}

// Put stores exact bytes and returns their stable hash.
func (c *CAS) Put(data []byte) (hash string, evicted []string) {
	sum := sha256.Sum256(data)
	hash = hex.EncodeToString(sum[:])
	if _, exists := c.objects[hash]; !exists {
		c.order = append(c.order, hash)
	}
	c.objects[hash] = append([]byte(nil), data...)
	for len(c.order) > c.limit {
		oldest := c.order[0]
		c.order = c.order[1:]
		delete(c.objects, oldest)
		evicted = append(evicted, oldest)
	}
	return hash, evicted
}

// Has reports whether exact bytes are still locally retained.
func (c *CAS) Has(hash string) bool {
	_, ok := c.objects[hash]
	return ok
}

// Count returns the retained object count.
func (c *CAS) Count() int { return len(c.objects) }
