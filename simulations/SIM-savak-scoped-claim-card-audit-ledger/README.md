# SIM-savak-scoped-claim-card-audit-ledger: Scoped claim card + audit records

This promoted simulation breeds `SIM-robot-app-semantics-conformance` with `SIM-kugap-live-sync-audit-split` by turning their strongest boundary ideas into concrete evidence shapes.

## Core move

Early apps do not publish one vague conformance statement. They publish:

1. a **scoped claim card** per protocol plane or embodiment;
2. an **audit message** that cites exactly one immutable **audit-anchor object**; and
3. local append-only **promise-accounting records** for failures, retries, retention changes, and break-witnesses.

## Why this should score higher

The parents were strong on honesty and layer separation, but weak on concrete claim format, exact cited audit object, and local evidence shape. This specimen repairs those gaps without freezing an app SDK or live transport.

## Design deltas

### 1. Scoped claim cards
Every app embodiment or service publishes a small claim card with:
- authoritative boundary identity: frozen `pCID` if available, else draft spec path plus revision hash and explicit provisional label;
- scope plane: `live`, `audit`, `storage`, or `device`;
- status: `implemented`, `provisional`, `blocked`, or `local-only`;
- local-only IDs and dependencies;
- wire-visible IDs;
- capability or authorization note;
- availability and retention note;
- optional shared app family ID for multi-embodiment claims.

### 2. Exact audit anchor
Every durable audit publication cites one immutable CAS object. Examples:
- blob CID for a minimal immutable blob app;
- snapshot or manifest CID for CRDT save-time audit;
- effect receipt CID or break-witness CID for device actions;
- key-rotation record CID for signing-key continuity.

The audit message may also include human-readable prose, but the cited object is exact.

### 3. Minimal promise-accounting records
Each peer keeps local append-only records for:
- publish or store attempt;
- retrieval attempt and result;
- availability or retention promise observed;
- retention-policy change;
- operation or effect request key;
- receipt, failure, or break-witness.

For non-idempotent physical effects, ambiguous restart state must default to **break-witness, not re-execute**.

## Scenario fit

- **App semantics partial conformance:** honest, checkable claim structure.
- **Live CRDT audit publication:** exact save-time anchor object and clean live-vs-audit separation.
- **Minimal immutable blob app:** hash is address to content identity, not authorization and not an availability promise unless a claim card says who promises retrieval.
- **Device-bound physical effect:** stable operation key, receipt, and break-witness rule improve replay handling.
- **Multi-embodiment app identity:** same app means shared boundary identity or family plus per-embodiment scoped cards, not branding.
- **Portable signing-key identity:** key rotation can be anchored by a cited rotation record without freezing final crypto forever.

## Boundaries

This simulation still does not define a final live transport, universal capability token format, or frozen witness envelope. It only proposes auditable evidence shapes that keep provisional claims honest.

This specimen was promoted from review-stage child proposal
`SIM-savak-child-scoped-claim-card-audit-ledger` from
`ga-canary-20260521-011601` under `DI-fihub`. The ignored proposal artifacts
remain local raw evidence; this directory is the canonical non-child simulation
home.
