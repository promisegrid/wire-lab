package state

import (
	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/protocol"
)

// CAS is a tiny sparse store with limit-driven GC.
type CAS struct {
	limit   int
	order   []string
	objects map[string][]byte
}

// NewCAS creates bounded local state for the simulated M4 agent.
func NewCAS(limit int) *CAS {
	return &CAS{limit: limit, objects: make(map[string][]byte)}
}

// Put stores exact bytes and returns their stable CID.
func (c *CAS) Put(data []byte) (contentCID string, evicted []string) {
	contentCID, err := protocol.CIDForBytes(data)
	if err != nil {
		panic(err)
	}
	if _, exists := c.objects[contentCID]; !exists {
		c.order = append(c.order, contentCID)
	}
	c.objects[contentCID] = append([]byte(nil), data...)
	for len(c.order) > c.limit {
		oldest := c.order[0]
		c.order = c.order[1:]
		delete(c.objects, oldest)
		evicted = append(evicted, oldest)
	}
	return contentCID, evicted
}

// Get returns exact retained bytes for a CID.
func (c *CAS) Get(contentCID string) ([]byte, bool) {
	data, ok := c.objects[contentCID]
	if !ok {
		return nil, false
	}
	return append([]byte(nil), data...), true
}

// Has reports whether exact bytes are still locally retained.
func (c *CAS) Has(contentCID string) bool {
	_, ok := c.objects[contentCID]
	return ok
}

// CIDs returns retained object CIDs in retention order.
func (c *CAS) CIDs() []string {
	return append([]string(nil), c.order...)
}

// Count returns the retained object count.
func (c *CAS) Count() int { return len(c.objects) }
