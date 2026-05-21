# Claim cards

## Status vocabulary

Every field in every card must be tagged with one of:

- `frozen`: claimed against a frozen obligation
- `provisional`: useful now, expected to evolve
- `blocked`: intentionally not claimed yet
- `local-only`: true only inside one adapter, host, or embodiment

Non-frozen fields should include `migration_note`.

---

## 1. Port Boundary Card

Use for infrastructure ports.

### Required fields

- `claim_subject`
- `role_term_used` — e.g. runtime, dispatcher, library
- `spec_refs_frozen`
- `spec_refs_orientation_only`
- `implemented_surfaces`
- `explicit_non_goals`
- `harness_apparatus_excluded`
- `conformance_scope`
- `status_tags`
- `evidence_refs`

### Honest wording pattern

> This implementation claims conformance only for the listed frozen surfaces. Draft and harness-era materials are cited as orientation only. Unlisted ingress, feed, CAS, session, app, or handler behaviors are not claimed.

---

## 2. App Semantics Card

Use for app-side partial conformance.

### Required fields

- `app_identity_at_protocol_boundary`
- `contract_ref_or_family_ref`
- `supported_operations`
- `unsupported_or_blocked_operations`
- `local_ids_are_internal_only`
- `signature_carriage_status`
- `capability_or_witness_language_status`
- `status_tags`
- `evidence_refs`

### Honest wording pattern

> This app implementation supports only the listed protocol-boundary operations. Local IDs, local storage handles, and adapter-local message shapes are internal and are not claimed as interoperable contract identity. Signatures below are evidence for this implementation's current behavior and are not, by themselves, a claim that a final wire-level signature contract is frozen.

### CAS-safe wording for minimal blob apps

> Returning a content hash asserts byte-identity of the referenced content object only. It does not, by itself, assert continued availability, replication, retention duration, location, or read authorization.

---

## 3. Host / Embodiment / Dependency Card

Use when behavior depends on browser storage, Node helpers, CUPS, libusb, live channels, vendor SDKs, or device daemons.

### Required fields

- `embodiment_id`
- `shared_app_identity_ref` or `same_app_basis`
- `host_capabilities_and_limits`
- `external_dependencies`
- `failure_surfaces`
- `storage_and_key_handling_notes`
- `replay_restart_posture`
- `live_channel_claims_separate_from_audit_claims`
- `status_tags`
- `evidence_refs`

### Honest wording patterns

For live CRDT:

> Real-time sync uses a host-specific live channel and is not claimed here as core PromiseGrid transport conformance. Durable audit publication claims apply only to the cited snapshot objects and audit messages.

For device-bound agents:

> Physical-effect execution depends on the listed host drivers and daemon state. This card declares implementation dependencies and replay posture; it does not turn those host dependencies into PromiseGrid wire obligations.

For browser/plugin identity:

> Shared app identity is claimed via the cited contract reference and signed continuity records, not by display name, local username, or host-specific storage location alone.