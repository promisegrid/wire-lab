# Guide-Safe Claim Bundle Pattern

## 1. Publishable layer-scoped conformance manifest

A manifest is a small publishable statement with at least:

- `app_contract_ref`: authoritative protocol-boundary contract or contract-family reference
- `embodiment_ref`: local embodiment label; not authoritative app identity
- `claim_scope`: one of `local-implementation`, `live-transport`, `audit-publication`, `device-effect`
- `status`: `implemented`, `provisional`, or `blocked`
- `wire_artifacts_observable_now`: what Bob can actually inspect
- `dependencies_not_implied_by_protocol`: host/runtime/device dependencies such as browser storage, Node helper, CUPS, libusb, IPP, vendor SDK, or custom relay
- `signing_note`: whether signature carriage is provisional, adapter-local, or frozen elsewhere
- `audit_object_kind`: `none`, `blob`, `snapshot`, `save-blob`, `receipt`, `break-witness`, or similar

## 2. Durable citation rule

- A durable audit claim must cite a concrete CAS object.
- A live session, socket, relay, or transient channel is never itself the durable audit object.
- For immutable blob apps, cite the content hash separately from any availability or retention promise.
- For live CRDT apps, cite a durable milestone object such as a save blob or snapshot, optionally with an op-range digest.
- For physical effects, cite request identity plus a durable receipt or break-witness.
- For key rotation, cite a durable signed continuity artifact rather than a display name or host-local account label.

## 3. Peer-local promise-accounting records

Each peer keeps a local promise-accounting record row with:

- `timestamp`
- `actor`
- `object_ref`
- `promise_class`: `identity`, `availability`, `authorization`, or `effect`
- `promisor`
- `scope_or_retention`
- `result`: `satisfied`, `failed`, `unknown`, `replayed`, `superseded`
- `evidence_ref`: local log, receipt, cited object, or break-witness

## 4. Interpretation notes

- Possession of a hash is not automatically authorization.
- A content hash can identify bytes without promising discovery or availability.
- Shared branding does not prove shared app identity; the contract reference does.
- A published audit claim does not imply that an unfrozen live transport is PromiseGrid-conformant.
- If restart leaves a non-idempotent effect ambiguous, the safe durable output is usually a break-witness, not silent re-execution.

## 5. Intended outcome

The pattern makes partial conformance more checkable, keeps live-versus-audit boundaries honest, and gives Alice, Bob, and Carol enough local evidence to audit stale claims and failures years later.
