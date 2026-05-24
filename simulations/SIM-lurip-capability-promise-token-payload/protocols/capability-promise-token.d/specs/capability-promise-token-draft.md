# Capability Promise Token Draft

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Source: `DI-kafiz`.

## Purpose

This higher-layer payload protocol represents a capability-like token as a
promiser's signed promise. The token is evidence that a receiver may assess; it
is not a universal permission object.

## Candidate Payload Shape

```cddl
capability-promise-token = [
  promiser: agent-ref,
  promisee: agent-ref / nil,
  action_pcid: cid,
  scope_bytes: bstr,
  transfer_rule: "nontransferable" / "transferable-with-history" / "bearer-like-promise",
  expiry_or_freeze_ref: cid / nil,
  promiser_proof: bstr
]

capability-transfer-record = [
  transfer_promiser: agent-ref,
  token_ref: cid,
  recipient: agent-ref,
  transfer_observation: bstr,
  transfer_proof: bstr
]
```

## Promise-Theory Rules

- Alice's token is Alice's promise; it does not force Alice or any other agent
  to act.
- Bob cannot promise Alice's future behavior when transferring a token; Bob can
  only promise or attest to Bob's own handling.
- Carol assesses Alice's promise, Bob's transfer history, expiry/freeze evidence,
  and Carol's own relationship history before relying on the token.
- The pCID-named protocol defines action scope, transfer evidence, freeze or
  expiry references, and proof encoding.

## Non-Goals

This draft does not define generic assertion machinery, a universal
authorization ledger, or a base-envelope capability selector.
