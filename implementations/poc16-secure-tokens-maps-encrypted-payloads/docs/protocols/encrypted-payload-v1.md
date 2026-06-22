# encrypted_payload_v1

## Status

Active POC16 encrypted-payload profile. The embedded Markdown bytes derive this
pCID. Source: `DI-bitug`; `DI-vulit`.

## Abstract

`encrypted_payload_v1` tests real authenticated encryption of a pCID-owned slot-1
payload. It covers recipient checks, tamper rejection, visible parent links, and
hidden parent-link commitments while explicitly leaving production key agreement
to later protocols.

## pCID and envelope

The active shape is:

```text
grid([42(pCID), payload, proof])
```

The payload is a CBOR string map chosen by this pCID for readability.
Intermediate routers may inspect the envelope shell and visible parent links but
MUST NOT need plaintext to forward exact bytes.

## Promise Theory model

The sender promises that it encrypted plaintext for the named recipient under the
stated context. The recipient locally decides whether it can decrypt and trust the
plaintext. Transit peers promise only carriage of ciphertext bytes.

## Payload grammar

```text
payload = {
  "sender": text,
  "recipient": text,
  "context": text,
  "algorithm": "AES-256-GCM",
  "key_profile": "poc16-local-derived-demo-key",
  "nonce_b64": base64 bytes,
  "ciphertext_b64": base64 bytes
}
```

The AEAD additional authenticated data is the byte string
`poc16-encrypted-payload-v1\0sender\0recipient\0context`. POC16 derives a local
demo key as SHA-256 over the same tuple; production MUST replace that with real
relationship-specific key promises.

## Sender behavior

The sender MUST set non-empty `sender`, `recipient`, and `context`; generate a
fresh nonce; and encrypt the plaintext with AES-256-GCM. If parent links are
sensitive, the sender should put them inside encrypted plaintext or commit to
them with a hidden digest rather than visible envelope parent slots.

## Receiver and parser behavior

A parser MUST decode the CBOR map and project `sender` to local `from` and
`recipient` to local `to` only for routing compatibility. The recipient MUST
reject wrong algorithm, wrong key profile, wrong recipient, invalid base64, or
AEAD authentication failure. Rejection is local non-commitment unless the sender
promised decryptability under terms the receiver accepted.

## Protocol state machine

```text
[plaintext local]
    | encrypt for recipient/context
    v
[ciphertext sent] --recipient decrypts and authenticates--> [plaintext accepted]
       | wrong recipient / tamper / key unavailable
       v
[local non-commitment]
```

## State, CAS, DAG, and retention

Ciphertext envelopes may be stored in CAS without revealing plaintext. A receiver
may store decrypted plaintext only according to local promises. Visible parent
links are public to any envelope reader; hidden links require pCID-specific
plaintext or commitment handling.

## Security considerations

The POC16 key profile is not production security. Real deployments need key
exchange, rotation, compromise recovery, nonce discipline, and relationship-level
promises. Encrypting payloads does not hide metadata in the outer envelope.

## Interoperability notes

CBOR maps are acceptable here because the payload is not intended for constrained
IoT devices. A constrained encrypted profile may choose positional arrays under a
new pCID.

## Examples

```text
grid([42(pCID),
  {"sender":"alice", "recipient":"bob", "context":"shipment:42",
   "algorithm":"AES-256-GCM", "key_profile":"poc16-local-derived-demo-key",
   "nonce_b64":"...", "ciphertext_b64":"..."},
  proof
])
```
