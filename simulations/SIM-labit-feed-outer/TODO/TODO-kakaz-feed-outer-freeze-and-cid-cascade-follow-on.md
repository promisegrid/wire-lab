# TODO-kakaz: Feed-outer freeze and CID-cascade follow-on

## Prior aliases

None. This TODO is created directly as a sim-local successor owner under
`rusis.10`.

## Status

Closed. The feed-outer specimen slice transferred out of rooted
`TODO-turog` and `TODO-duvuk` by `DI-mosor` is resolved by `DI-bomud`.

## Scope

This TODO owns the feed-outer side of the split follow-on work:

- outer transport-spec freeze boundary cleanup and freeze-gate follow-on
  that is not group-session-local;
- transport-level filename/CID-cascade policy and related provenance
  rules for the feed-outer lineage;
- extraction-handoff prep for `rusis.11` from rooted
  `transport-spec-draft.md`.

Out of scope: group-session-local envelope details (tracked in
`TODO-gapab`).

## Closure summary

`DI-bomud` closes this successor owner. Spec freeze publishes a spec pCID and
does not rename, rehash, or rewrite historical feed/transport specimens.
Draft-era directories remain evidence; any frozen successor or derived mirror
is additive and cites its source evidence. Feed-outer remains a thin outer
convention and does not own group-session message filename/CID,
`Message-ID:`, header, body, or reader/writer rules.

## Subtasks

- [x] kakaz.1 Define feed-outer freeze-boundary wording that supersedes
  stale TE-41 step-5 guidance without mutating historical specimen bytes.
  Closed by `DI-bomud`.
- [x] kakaz.2 Define the feed-outer slice of TE-42 Path-A-vs-Path-B and
  deprecation policy handling. Closed by `DI-bomud`: feed-outer owns no
  message-header compatibility rule.
- [x] kakaz.3 Prepare extraction handoff notes for `rusis.11` from
  rooted `protocols/wire-lab.d/specs/transport-spec-draft.md`. Closed by
  `DI-bomud`; active wording now lives in the simulation-local feed-outer
  draft.
- [x] kakaz.4 Back-link resulting decisions and artifacts to rooted
  `TODO-turog` and `TODO-duvuk` historical records. Closed by `DI-bomud`.

## Decision Intent Log

Successor-owner routing into this TODO was locked under `DI-mosor` in
`protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`.

ID: DI-bomud
Date: 2026-05-13 22:48:21
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Close the feed-outer successor slice for turn-164, TE-41, and TE-42
follow-on work. Spec freeze publishes new spec pCIDs but never rewrites
historical transport/feed data; frozen successors or derived views must be
additive and cite prior evidence. Feed-outer remains a thin outer convention
and does not own group-session message filename/CID, `Message-ID:`, header, or
body parsing rules. A Steve-authored DI is the operative
`merge-transport-spec` promise until cryptographic promise tooling exists.
Intent: Retire the stale rewrite-at-freeze plan while preserving a clean
feed-outer boundary from group-session-local message semantics.
Constraints: Do not mutate historical message bytes, old CIDs, or legacy
transport specimens. Keep group-session-local policy under `DI-rurab`. Do not
implement cryptographic signing or `tools/spec freeze` changes in this pass.
Affects: `simulations/SIM-labit-feed-outer/protocols/feed-outer.d/specs/feed-outer-draft.md`;
`simulations/SIM-ludut-wire-lab-devs/world/transports/wire-lab-devs-draft/README.md`;
`simulations/SIM-labit-feed-outer/TODO/TODO-kakaz-feed-outer-freeze-and-cid-cascade-follow-on.md`;
`protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md`;
`protocols/wire-lab.d/TODO/TODO-duvuk-te-42-filename-cid-cascade-policy.md`;
`protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`.
