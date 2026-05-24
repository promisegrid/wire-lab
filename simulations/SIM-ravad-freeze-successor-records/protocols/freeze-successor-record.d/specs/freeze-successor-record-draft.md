# Freeze Successor Record Draft

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Source: `DI-kafiz`.

## Purpose

This higher-layer payload protocol records an agent's own freeze or successor
promise for exact protocol/specimen bytes. It does not make the base envelope
understand freeze semantics.

## Candidate Payload Shape

```cddl
freeze-successor-record = [
  promiser: agent-ref,
  record_kind: "freeze" / "successor" / "withdrawal",
  predecessor_ref: cid / nil,
  successor_ref: cid / nil,
  exact_bytes_ref: cid,
  promise_text: bstr,
  promiser_proof: bstr
]
```

## Promise-Theory Rules

- Alice may promise that she will treat `exact_bytes_ref` as frozen for her
  future work, but Alice cannot freeze Bob's or Carol's behavior.
- A successor record links evidence; it does not rewrite the predecessor.
- Bob may adopt, ignore, contest, or supersede Alice's successor locally.
- Exact-byte references let later peers audit the promise without trusting a
  mutable registry.

## Non-Goals

This draft does not define a global freeze authority, a universal statement
capsule, or base-envelope fields for freeze metadata.
