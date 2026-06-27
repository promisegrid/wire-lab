package device

// StorageToken is Bob's compact bearer capability as held by one endpoint.
type StorageToken struct {
	Bytes              []byte
	CID                string
	Issuer             string
	Holder             string
	MaxContentBytes    uint64
	MaxRetainedObjects uint64
	RetentionTerms     string
}
