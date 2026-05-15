# TODO-gapab: Group-session freeze and CID-cascade follow-on

## Prior aliases

None. This TODO is created directly as a sim-local successor owner under
`rusis.10`.

## Status

Closed. The group-session specimen slice transferred out of rooted
`TODO-turog` and `TODO-duvuk` by `DI-mosor` is resolved by `DI-rurab`.

## Scope

This TODO owns only the group-session side of the split follow-on work:

- group-session freeze-doc and freeze-gate cleanup that belongs to the
  `group-session` lineage specimen;
- group-session envelope-level filename/CID-cascade policy questions
  (including reader/writer compatibility details);
- cross-reference and boundary notes needed to coordinate with the
  feed-outer successor owner.

Out of scope: feed-outer transport-spec extraction and outer-rule
ownership (tracked in `TODO-kakaz`).

## Closure summary

`DI-rurab` closes this successor owner. For the group-session git-bound
specimen, membership is the fixed configured set of exact
`<author-id>/main` branches at transport creation; arbitrary branch discovery
does not admit new members, and passive read-only observers are not members.
The per-author-branch git binding is normative for this simulation's git-bound
specimen without preventing other future group-session bindings. Canonical
writers omit `Message-ID:`; readers may accept exactly one legacy header in the
pre-`Date:` slot without stripping or rehashing historical bytes.

## Subtasks

- [x] gapab.1 Recast the group-session slice of TE-41 freeze follow-on
  into a simulation-local checklist tied to
  `specs/group-session-draft.md`. Closed by `DI-rurab`.
- [x] gapab.2 Resolve the group-session slice of TE-42
  filename/CID-cascade policy with explicit reader and writer behavior.
  Closed by `DI-rurab`.
- [x] gapab.3 Coordinate cross-layer boundary notes with
  `simulations/SIM-labit-feed-outer/TODO/TODO-kakaz-feed-outer-freeze-and-cid-cascade-follow-on.md`.
  Closed by `DI-rurab` and the coordinated feed-outer decision `DI-bomud`.
- [x] gapab.4 Back-link resulting decisions and artifacts to rooted
  `TODO-turog` and `TODO-duvuk` historical records. Closed by `DI-rurab`.

## Decision Intent Log

Successor-owner routing into this TODO was locked under `DI-mosor` in
`protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`.

ID: DI-rurab
Date: 2026-05-13 22:48:21
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Close the group-session successor slice for turn-164, TE-41, and
TE-42 follow-on work. Membership in the git-bound specimen is the fixed
configured set of exact `<author-id>/main` branches at transport creation;
unknown branches are ignored until explicitly admitted through a successor
transport instance or future DI; passive read-only observers are non-members.
The per-author-branch git binding is normative for this simulation's git-bound
specimen but does not preclude other group-session bindings. Canonical writers
omit `Message-ID:`, while readers may tolerate one legacy `Message-ID:` header
without rewriting bytes. A Steve-authored DI is the operative
`merge-group-transport-spec` promise until cryptographic promise tooling
exists.
Intent: Finish the group-session half of the turn-164 pending work without
mutating historical transport messages or leaving membership, freeze, and
CID-cascade policy spread across conversation-only UT rows.
Constraints: Do not edit historical message files or claim that a spec freeze
rewrites transport data. Keep feed-outer freeze-boundary decisions in
`DI-bomud`. Keep rooted `TODO-turog` and `TODO-duvuk` as historical
coordination memory. Treat this DI as interim authorization, not as a
cryptographic signing implementation.
Affects: `simulations/SIM-rakot-group-session/protocols/group-session.d/specs/group-session-draft.md`;
`simulations/SIM-ludut-wire-lab-devs/world/transports/wire-lab-devs-draft/README.md`;
`simulations/SIM-rakot-group-session/protocols/group-session.d/TODO/TODO-gapab-group-session-freeze-and-cid-cascade-follow-on.md`;
`protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md`;
`protocols/wire-lab.d/TODO/TODO-duvuk-te-42-filename-cid-cascade-policy.md`;
`protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`.
