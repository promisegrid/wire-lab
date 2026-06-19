# secure_capability_v1

Agents use `secure_capability_v1` for cryptographically signed capability-token
promises. POC16 uses COSE_Sign1 over a CWT-style CBOR claim map containing
issuer, subject, audience, expiry, not-before time, token ID, capability, scope,
content CID, transferability, and optional holder confirmation.
