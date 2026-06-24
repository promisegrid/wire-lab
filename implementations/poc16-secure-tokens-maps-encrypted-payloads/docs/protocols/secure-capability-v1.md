# secure_capability_v1

## Status

Active POC16 secure capability-token profile. The embedded Markdown bytes derive
this pCID. Source: `DI-bitug`; `DI-vulit`; `DI-mapop`.

## Abstract

`secure_capability_v1` specifies cryptographically signed capability-token
promises. A token is a signed promise by the issuer that a holder may later
redeem under stated terms. It is not an authorization object from a central
authority.

## pCID and envelope

The active specimen shape is:

```text
grid([42(pCID), payload, proof])
```

The payload MAY be a CBOR map because this pCID chooses a self-documenting token
profile. Capability tokens embedded in other protocols, such as `cas_storage_v1`
or `production_shipping_v1`, follow the token byte rules here even when carried
as text fields there.

## Promise Theory model

The issuer is the promiser. The subject, audience, or holder is the promisee or
redeemer. A bearer token may be transferred if the token says it is transferable;
a holder-bound token requires confirmation material. Redemption is a local event
where the issuer decides whether the token terms match a still-live promise.

## Payload grammar

POC16 has two compatible token encodings:

```text
signed_capability_token_v1 = base64(COSE_Sign1(CBOR string map claims))
claims = {
  "type": "signed_capability_token_v1",
  "issuer": text,
  "subject": text,
  "scope": text,
  "content_cid": text,
  "expires_unix": text integer,
  "nonce": text,
  "transferable": "true" / "false"
}

cwt_capability_token_v1 = base64(COSE_Sign1(CWT-style CBOR claim map))
CWT labels: 1 issuer, 2 subject, 3 audience, 4 exp, 5 nbf, 7 token_id,
-70000 capability, -70001 scope, -70002 content_cid,
-70003 transferable, -70004 confirmation.
```

COSE_Sign1 MUST use tag 18, protected header `{1: -8}` for EdDSA, empty
unprotected header, embedded payload, and an Ed25519 signature in POC16.

## Sender behavior

An issuer MUST sign the token and set expiry, nonce/token ID, scope, content CID,
and transferability. A non-transferable CWT token MUST include confirmation. A
sender offering a token as reciprocal payment SHOULD state what promise the token
is intended to satisfy.

## Receiver and parser behavior

A receiver MUST base64-decode, verify COSE, verify issuer and audience where
applicable, reject tokens before `nbf`, reject expired tokens at or after `exp`,
check scope/content CID on redemption, and replay-check nonce or token ID against
local issuer state. A failed token is local non-commitment unless a prior token
promise was broken.

## Protocol state machine

```text
[no token]
    | issuer signs token promise
    v
[token issued] --transfer if bearer--> [token held by new holder]
      | redeem with matching terms
      v
[token redeemed / consumed]
      | replay / expired / wrong audience / revoked
      v
[locally rejected]
```

## State, CAS, DAG, and retention

Issuers SHOULD remember live, redeemed, revoked, and expired token IDs for the
run. Holders MAY store token bytes in local CAS. Parent links SHOULD connect token
issue, transfer, and redemption messages when carried in envelopes.

## Security considerations

Tokens are bearer-value when transferable and must be protected like money.
Replay, confused deputy, wrong audience, expiry, and stolen holder-bound tokens
are the primary threats. POC16 deterministic keys are not production keys.

## Interoperability notes

CWT is recommended for production successors because it has compact numeric
claims and established COSE composition. The string-map signed token remains a
POC16 pressure model for readability and transition.

## Examples

```text
grid([42(pCID),
  {"issuer":"bob", "subject":"alice", "audience":"bob", "scope":"serve-once",
   "content_cid":"bafkrei...", "token_id":"tok-7",
   "transferable":"false", "confirmation":"alice-key", "exp":"1800000000"},
  proof
])
```
