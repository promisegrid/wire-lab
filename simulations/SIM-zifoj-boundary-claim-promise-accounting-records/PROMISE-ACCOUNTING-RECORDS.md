# Minimal promise-accounting records

## Row schema

Each local row should record:

- `time`
- `actor`
- `operation_kind`
- `object_or_effect_id`
- `claim_card_ref`
- `promise_made_or_relied_on`
- `counterparty_or_observer`
- `outcome` — success, unavailable, refused, replayed, superseded, broken
- `evidence_refs`
- `notes_on_retention_auth_or_dependency`

## Required operation kinds

- `publish_blob`
- `retrieve_blob`
- `availability_claim`
- `effect_request`
- `effect_receipt`
- `break_witness`
- `snapshot_publish`
- `key_rotation`
- `delegation_record`

## Scenario guidance

### Minimal immutable blob app

Bob's promise-accounting record row for returning a hash should state whether he promised only content identity or also any separate storage, discovery, retention, or authorization service. If none, it must say none.

Carol's failed retrieval row should record:

- requested hash,
- where she looked,
- what availability claim she relied on,
- whether auth was missing versus bytes unavailable.

### Device-bound physical effect

A non-idempotent effect needs a durable `effect_id` or dedupe key that survives restart. Replays should produce either the prior receipt or a break-witness, not a silently duplicated effect.

### Live CRDT audit publication

Keep live-session rows separate from durable audit rows. The durable row cites the snapshot object or manifest; it does not claim the live channel itself was core grid conformance.

### Portable signing-key identity

A key-rotation row should link `old_key_id -> new_key_id` with signed evidence from the old key, the new key, or both, plus any host-specific storage caveat.

## Invariants

- A hash is never enough, by itself, to imply authorization.
- Availability claims must name who made them.
- Host dependency failures are promise-accounting record events, not silent implementation detail.
- Break-witness is a first-class outcome, not an afterthought.
