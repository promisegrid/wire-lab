# identity_key_v1

Agents use `identity_key_v1` to promise key-rotation and key-continuity facts.
Payloads are pCID-owned CBOR arrays, not universal maps. A key record is a local
promise about the promiser's future signing identity, not a global identity
authority.
