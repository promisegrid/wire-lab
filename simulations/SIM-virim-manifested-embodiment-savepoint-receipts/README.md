# SIM-virim-manifested-embodiment-savepoint-receipts: Manifested embodiment savepoint receipts

This promoted simulation breeds `SIM-nijuz-multi-embodiment-identity` and `SIM-kugap-live-sync-audit-split` into a more concrete worked profile.

The core move is simple:

1. One logical app is bound by a content-addressed App Manifest.
2. Each embodiment publishes a scoped Embodiment Claim under a durable signing-identity continuity chain.
3. Live collaboration stays off-grid for now, while durable PromiseGrid publication cites a concrete Savepoint Audit Envelope.

## Why this should score higher

- It repairs the parent ambiguity about what makes browser and helper embodiments one app: shared manifest CID, not branding or host coincidence.
- It repairs the parent ambiguity about authoritative identity and rotation evidence: continuity receipts link embodiment keys and rotated keys; display names and local handles are non-authoritative hints only.
- It repairs the parent ambiguity about what an audit message cites: every durable publication cites a Savepoint Audit Envelope that names the exact snapshot, receipt, or break-witness object.
- It keeps the parent strength of strict layer separation: live transport is not silently upgraded into PromiseGrid conformance.

## Artifact set

### 1. App Manifest or AM

A content-addressed object that defines the logical app boundary.

Required fields:
- manifest CID
- logical app label
- contract-family CID or draft reference set
- expected embodiment classes such as browser, plugin helper, daemon
- shared semantics required across embodiments
- explicitly blocked or excluded semantics
- migration or successor-manifest reference if the app evolves

Rule: two embodiments count as one app only when their claims reference the same AM CID or an auditable successor chain.

### 2. Embodiment Claim or EC

Each runtime component signs its own honest partial-conformance claim.

Required fields:
- AM CID
- embodiment identifier and runtime class
- implemented contract subset
- runtime and storage limits
- live-channel role, if any
- excluded guarantees
- supported audit-envelope variants
- current signing key reference

Rule: the authoritative protocol-boundary identity is the signing key named in the continuity chain, not a display name, local username, browser profile name, or adapter-local ID.

### 3. Identity Continuity Receipt or ICR

A signed linkage object for browser-helper binding and key rotation.

Required fields:
- prior key reference, or null for first publication
- next key reference
- scope: user, embodiment, or both
- reason: rotation, embodiment-link, replacement, revocation, loss
- time or event context
- signatures by both keys when possible
- if both signatures are impossible, a one-key statement plus later witness or break-witness

Rule: if continuity cannot be shown, publish a break-witness instead of implying continuity.

### 4. Savepoint Audit Envelope or SAE

The only durable on-grid citation surface for live-state apps in this profile.

Required fields:
- AM CID
- EC reference set
- active ICR reference set
- cited durable object: snapshot CID, save blob CID, physical-effect receipt CID, or break-witness CID
- human-readable promise body
- publication channel information
- explicit statement that live transport was off-grid and outside this conformance claim

Rule: group-session may publish the SAE, but group-session is not the live transport.

## Local promise accounting

Each actor keeps small local promise-accounting records with:
- AM CID used
- EC and ICR references observed
- request, savepoint, or effect ID
- expected SAE or receipt reference
- observed failure class such as missing object, replay suspicion, retention loss, stale claim, or broken continuity

This gives later auditors durable explanation even when storage, transport, or host recovery fails.

This specimen was promoted from review-stage child proposal
`SIM-virim-child-manifested-embodiment-savepoint-receipts` from
`ga-canary-20260521-210902` under `DI-lanuz`. The ignored proposal artifacts
remain local raw evidence; this directory is the canonical non-child simulation
home.

## Scenario pressure fit

- `app-semantics-partial-conformance`: EC gives honest wording and explicit excluded guarantees.
- `multi-embodiment-app-identity`: AM makes one app a contract fact rather than a UX story.
- `portable-signing-key-identity`: ICR gives an auditable rotation and cross-host linkage object.
- `live-crdt-audit-publication`: SAE fixes the cited durable object while keeping live sync off-grid.
- `device-bound-agent-physical-effect`: the same envelope can cite a physical-effect receipt or break-witness rather than a snapshot.
- `minimal-immutable-blob-app`: content hash names content identity only; availability, authorization, replication, and retention need separate receipts or claims.

## Boundaries

This simulation does not define a reliable live binding, a final cryptographic suite, or a permanent storage recipe. It standardizes only a guide-safe v0 audit and claim surface: manifest, embodiment claim, continuity receipt, and savepoint envelope.
