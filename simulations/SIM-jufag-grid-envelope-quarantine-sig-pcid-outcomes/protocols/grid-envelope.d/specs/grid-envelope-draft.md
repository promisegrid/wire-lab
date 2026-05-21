# Grid Envelope Variant Spec (DRAFT)

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `enc-dag-cbor` / `unknown-quarantine` / `sig-mandatory-profiled-opaque-bytes`.
> Source: `DI-fanah`; promoted under `DI-dipid` from a GA child bred from the two parent specimens named in `README.md`.

## Purpose

This spec defines one full positional grid-envelope candidate for wire-lab comparison. It preserves the parent envelope simplicity while repairing two observed weaknesses from sampled runs:

- mixed-version peers need a way to preserve unknown future evidence without treating it as semantically valid;
- mandatory signature bytes should name their verification profile explicitly enough to improve auditability and long-term recovery.

This is a specimen inside `SIM-jufag-grid-envelope-quarantine-sig-pcid-outcomes`, not a harness rule and not the canonical PromiseGrid wire format.

## Positional Envelope Shape

The envelope shape for this variant is:

```text
[pcid, payload, sig_pcid, signature]
```

Slots are interpreted positionally:

- `pcid` identifies the protocol/spec/handler that interprets `payload`.
- `payload` is opaque bytes until interpreted by the handler named by `pcid`.
- `sig_pcid` identifies the protocol/spec/handler that interprets `signature` and verifies authorship semantics.
- `signature` is mandatory opaque bytes interpreted only by the handler named by `sig_pcid`.

A `payload` may itself be the canonical bytes of another grid envelope when the protocol named by `pcid` specifies recursive nesting. The outer grid envelope does not prescribe the payload's internal organization beyond the bytes boundary.

## Encoding

This variant encodes the envelope as a DAG-CBOR-compatible positional array of length 4.

- `pcid` and `sig_pcid` are DAG-CBOR Link values.
- `payload` and `signature` are byte strings.

The canonical bytes for hashing, storage, relay, and signing are the DAG-CBOR bytes of the exact positional array under this spec.

The signature input is the canonical DAG-CBOR bytes of the unsigned prefix:

```text
[pcid, payload, sig_pcid]
```

This binds the payload-dispatch rule and the signature-profile rule together so a relay cannot swap verification semantics without changing the signed bytes.

## Receiver Outcomes

A receiver MUST classify each received byte sequence into one of these audit outcomes:

1. `accepted`
2. `quarantined`
3. `rejected`

### Accepted

An envelope is `accepted` only if all of the following hold:

- the envelope is structurally well-formed under this spec;
- the receiver has a handler for `pcid`;
- the receiver has a handler for `sig_pcid`;
- the `sig_pcid` handler verifies `signature` over the canonical unsigned prefix;
- local policy allows the result to be relied on.

### Quarantined

An envelope is `quarantined` if it is well-formed but cannot yet be semantically accepted because dispatch knowledge is incomplete. At minimum:

- unknown `pcid` => `quarantined-unknown-pcid`
- unknown `sig_pcid` => `quarantined-unknown-sig-pcid`

For a quarantined envelope, a receiver:

- MAY store the exact canonical bytes;
- MAY relay the exact canonical bytes unchanged;
- MAY attach local diagnostics out of band;
- MUST NOT speculatively parse `payload` without the handler named by `pcid`;
- MUST NOT count the envelope as authenticated or semantically valid without the handler named by `sig_pcid`.

Quarantine is therefore fail-closed for meaning, but not destructive to evidence retention.

### Rejected

An envelope is `rejected` if it is malformed or known-invalid, including:

- wrong array length;
- wrong slot types;
- duplicate or non-canonical encoding under the receiver's DAG-CBOR implementation;
- known `sig_pcid` with failed signature verification.

A rejected envelope is not a valid grid-envelope message under this variant. A receiver MAY retain local forensic notes, but the byte sequence is not to be forwarded as a valid instance of this envelope format.

## Unknown pCID Policy

This variant replaces the parent split between `unknown-opaque` and `unknown-hard-reject` with a bounded middle rule:

- unknown dispatch knowledge does **not** imply semantic acceptance;
- unknown dispatch knowledge does **not** force destruction of potentially important evidence.

This is intended to fit heterogeneous fleets, sparse routing peers, and future application object families better than either parent alone.

## Signature and Authorship Policy

The envelope layer enforces:

- presence of `sig_pcid`;
- presence of `signature`;
- canonical byte coverage of `[pcid, payload, sig_pcid]`.

The handler named by `sig_pcid` defines:

- signature algorithm;
- signer/key binding;
- delegation rules;
- revocation rules;
- time or freshness hooks, if any;
- multi-signature or threshold behavior, if any.

This keeps the outer envelope small while making signature interpretation recoverable from named dispatch rather than from undocumented ecosystem convention.

## Layering-Test Behavior

This variant answers the harness §1.3 layering scenarios as follows:

- Ordering disagreements are handled by the protocol named by `pcid`; the grid envelope preserves bytes plus the dispatch names needed to make disagreements explicit.
- Forwarding, relay, or hop-local evidence may be represented by wrapper envelopes or by payload protocols; quarantined forwarding allows mixed-version peers to preserve evidence without over-accepting it.
- External or content-addressed body references live inside `payload` under the protocol named by `pcid`; the outer envelope does not reinterpret them.
- Incompatible interpretation rules fail visibly at either the `pcid` or `sig_pcid` dispatch boundary.

## Scenario Fit Hypothesis

### IoT fleet maintenance

- Future firmware, maintenance, or telemetry object families can cross older gateways as quarantined evidence instead of being silently trusted or destroyed.
- Explicit `sig_pcid` gives a named place to define signer binding, key rotation, and revocation behavior per profile.
- The envelope still does not define device identity or access-control semantics itself; those remain payload/profile work.

### BGP routing

- Unknown future route object families or signature profiles can be preserved for incident forensics and delayed analysis.
- Fail-closed acceptance prevents unknown route semantics from being treated as valid reachability claims.
- Freshness, withdrawals, and policy semantics still belong to routing payloads or signature profiles, not the outer envelope.

### Application object family

- Future application families can mint new `pcid` values without reinterpreting old envelope bytes.
- Old peers can keep exact bytes as quarantined CAS artifacts, improving migration safety over hard reject.
- Adding `sig_pcid` avoids forcing all future families into one implicit signature convention.

## Non-Goals

This draft does not:

- declare a winning envelope for PromiseGrid;
- define a central pCID registry;
- solve payload-level identity, freshness, routing, or authorization semantics;
- require that every quarantined envelope eventually become interpretable;
- make the parent variants obsolete.

## Freeze Gate

This draft can freeze only after at least one simulation run compares it against sibling positional grid-envelope variants and Steve signs a merge/freeze promise for this specific specimen.
