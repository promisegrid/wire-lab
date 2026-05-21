# Live CRDT Audit Publication

## Scenario ID

live-crdt-audit-publication

## Source / Provenance

- Source type: new harness scenario
- Source path: `/home/stevegt/lab/promisegrid-dev-guide/FEEDBACK.md`
- Source row/title: `FB-hurit` and `FB-nilat`
- Source DI / TODO / TE: `DI-ragaz`; `TODO-rozas`; `DR-tuhaz`; `group-session-draft.md`; `udp-binding-draft.md`

## Purpose

Test real-time app pressure where live state needs reliable, ordered, low-latency
frames, but durable PromiseGrid evidence may be published at milestones.

## Setup

Alice edits a shared document in a browser while Bob edits the same document in
Neovim. Their live CRDT sync needs sub-second in-order delivery. Carol audits
durable snapshots later through group-session messages that cite content-addressed
state. Mallory drops or reorders live frames and delays audit publication.

## Stimulus

The live channel partitions for thirty seconds, then reconnects. Alice and Bob
continue editing. The app emits an audit message at save time with a snapshot
reference and human-readable promise body.

## Expected Pressure

The candidate design must avoid pretending that best-effort datagrams or
git-paced group-session are the live transport, while showing how durable audit
evidence can still survive for 100-year review.

## Scenario-Specific Evaluation Questions

- Should live state be off-grid until a reliable binding exists, or should a
  future live pCID shape be sketched?
- What exact object does the audit message cite?
- How are live-channel conformance claims kept separate from audit-layer claims?
