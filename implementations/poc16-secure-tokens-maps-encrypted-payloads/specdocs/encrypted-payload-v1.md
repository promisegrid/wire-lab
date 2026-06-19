# encrypted_payload_v1

Agents use `encrypted_payload_v1` when the pCID-owned slot-1 payload is
ciphertext plus metadata. POC16 uses AEAD encryption to test recipient checks,
tamper rejection, visible parent links, and hidden parent-link commitments while
leaving production key agreement to a later protocol.
