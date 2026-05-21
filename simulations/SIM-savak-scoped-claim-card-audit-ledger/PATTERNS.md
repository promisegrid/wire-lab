# Patterns

## 1. Scoped claim card

Minimal fields:
- `claim_id`
- `subject_kind`: app | embodiment | storage-service | device-agent
- `authoritative_boundary_id`
- `scope_plane`: live | audit | storage | device
- `status`: implemented | provisional | blocked | local-only
- `wire_visible_ids`
- `local_only_ids`
- `dependency_disclosure`
- `capability_auth_note`
- `availability_retention_note`
- `cited_audit_anchor` (optional until a durable publication exists)
- `shared_app_family_id` (required when multiple embodiments claim to be one app)
- `supersedes_or_rotation_ref` (optional)

Rules:
- local IDs must never be the authoritative boundary identity;
- a live-plane claim never implies audit-plane conformance;
- a storage-plane hash return does not imply retention, replication, discovery, or read authority unless the card says who promises those things;
- two embodiments are the same app only if they say so with shared boundary identity or family linkage, not by UX branding alone.

## 2. Audit-anchor object

A durable audit message cites exactly one immutable object:
- blob CID;
- snapshot or manifest CID;
- effect receipt CID;
- break-witness CID; or
- key-rotation record CID.

Human-readable promise text may explain the claim, but the anchor object is the precise replay target for later audit.

## 3. Minimal promise-accounting records

Recommended local entry types:
- `store_attempt`
- `publish_claim`
- `retrieve_attempt`
- `availability_observation`
- `retention_change_observation`
- `operation_request`
- `effect_receipt`
- `break_witness`
- `key_rotation_observation`

Recommended common fields:
- local timestamp
- actor
- counterparty or source
- object CID or operation/effect key
- observed result
- referenced claim card IDs
- referenced audit-anchor CID

## 4. Replay rule for non-idempotent effects

If a device agent restarts and cannot prove whether operation key `K` already completed, it must:
1. not execute the effect again;
2. emit a break-witness that names `K` and the evidence gap; and
3. let later reconciliation or operator action resolve the ambiguity.

That keeps at-most-once posture auditable even when host-local state is incomplete.
