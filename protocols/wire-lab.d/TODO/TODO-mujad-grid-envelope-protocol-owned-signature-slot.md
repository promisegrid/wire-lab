# TODO-mujad: Grid-envelope protocol-owned signature slot

## Prior aliases

None. This TODO was minted after the proquint-handle migration.

## Status

Open. This TODO owns a focused outer-envelope comparison follow-on: add one new
standalone simulation that tests a three-slot outer envelope
`[pcid, payload, signature]` where `pcid` owns the signature rules, then
compare it against existing outer-shape families. It is separate from
`TODO-tugoz`, which already owns the signable-view and Gordian comparison batch.
Source: `DI-kukuk`.

## Decision Intent Log

ID: DI-kukuk
Date: 2026-05-23 23:27:04
Status: active
Author: jessica@t7a.org (Jessica)
Decision: Track the protocol-owned-signature-slot question as a new harness
TODO that adds one standalone sim,
`SIM-dalor-grid-envelope-protocol-owned-signature-slot`, while reusing existing
comparators for the rest of the outer-shape batch. The new specimen keeps the
outer shape as `[pcid, payload, signature]`, makes the third slot mandatory,
binds the signed bytes to canonical `[pcid, payload]`, and makes `pcid` define
the proof rules, including varsig, multisig, or another proof family. Reuse the
existing six-scenario slice from `TODO-tugoz` for the first scored comparison.
Intent: The repo already tests outer `sig_pcid`, payload-owned proof, wrapper
proof, outer signed summary header, and nested outer attestation. What was
missing was the direct three-slot question: should the outer envelope carry a
proof sibling while letting `pcid` own the proof semantics, instead of adding a
separate outer signature selector? Adding only one new sim keeps the comparison
clean and avoids duplicating already-existing specimens.
Constraints: Keep the `ppx` bot identity rules untouched. Do not rewrite or
mutate historical scored results in place. Keep the current consensus wording in
`DEV-GUIDE-RESOURCES.md` unchanged except for registering the new sim as an
active alternative. The new sim must not silently add extra outer audit fields,
headers, or nested wrapper machinery. The first scored slice should reuse
`portable-signing-key-identity`, `live-crdt-audit-publication`,
`multi-embodiment-app-identity`, `chunk-feed-replication-sparse-advertisement`,
`device-bound-agent-physical-effect`, and `minimal-immutable-blob-app`.
Affects: `protocols/wire-lab.d/TODO/TODO.md`;
`protocols/wire-lab.d/TODO/TODO-mujad-grid-envelope-protocol-owned-signature-slot.md`;
`DR/DR-jalom-grid-envelope-protocol-owned-signature-slot.md`;
`docs/thought-experiments/TE-nahir-grid-envelope-protocol-owned-signature-slot.md`;
`simulations/SIM-dalor-grid-envelope-protocol-owned-signature-slot/`;
`simulations/README.md`; `DEV-GUIDE-RESOURCES.md`.

ID: DI-vatav
Date: 2026-05-24 11:34:06
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Correct `TE-pokul` so the prior Alt-B / Alt-C split is collapsed.
The three-slot shape `[pCID, payload, proof]` already has the intended
current-sender semantics when the protocol named by `pCID` defines the third
slot as the sender's proof over the message. The remaining envelope decision is
therefore minimal `[pCID, payload]` versus three-slot `[pCID, payload, proof]`,
not a separate proof-slot alternative versus a separate outer-attestation
alternative.
Intent: The earlier TE wording created a false distinction and obscured the
actual design question. The sender-proof meaning belongs inside the protocol
definition named by `pCID`; it is not a separate envelope family.
Constraints: Keep `pCID` as Protocol CID. Do not treat the proof slot as
authority or global trust. Do not rewrite scored sim trees or committed result
artifacts in place; supersede them under TODO-mopob where needed.
Affects: `docs/thought-experiments/TE-pokul-cbor-array-extensibility-pcid-owned-proof-and-grid-envelope-landing-shape.md`;
`protocols/wire-lab.d/TODO/TODO-mujad-grid-envelope-protocol-owned-signature-slot.md`;
`protocols/wire-lab.d/specs/harness-spec-draft.md`;
`protocols/wire-lab.d/TODO/TODO-mopob-pcid-protocol-cid-corpus-correction-and-envelope-successors.md`.

## Scope

This TODO covers:

- one new standalone simulation for the three-slot protocol-owned-signature
  outer shape;
- one repo-owned TE/DR packet that explains why this is a new question distinct
  from `TODO-tugoz`;
- simulation/guide-index registration for the new specimen;
- preparing the new specimen for later scored comparison against existing
  envelope families.

This TODO does not itself score the full comparison batch, settle the winning
outer shape, or replace the existing `TODO-tugoz` signable-view/Gordian work.

## Predecessor context

- `TODO-tugoz` already covers signable-view, wrapper-proof, UCAN-like, and
  Gordian comparison sims.
- The positional matrix already includes several mandatory-opaque-bytes
  specimens using `[pcid, payload, signature]`, but none explicitly makes
  `pcid` the named owner of varsig/multisig-style proof semantics as the main
  design move under test.
- `SIM-lotiv` keeps multisig placement modes open, but does not isolate the
  simpler three-slot pCID-owned proof rule as its own candidate specimen.

## Subtasks

- [x] mujad.1 Write a TE that compares a new three-slot pCID-owned proof
  specimen against narrower and broader comparison sets, and lock the broad
  five-family comparison as the right question packet. Source: `DI-kukuk`.
- [x] mujad.2 Create `SIM-dalor-grid-envelope-protocol-owned-signature-slot`
  as a standalone simulation with `README.md`, `QUESTION.md`, and a local
  `grid-envelope-draft.md`. Source: `DI-kukuk`.
- [x] mujad.3 Register the new sim in `simulations/README.md` and list it in
  `DEV-GUIDE-RESOURCES.md` as an active unscored alternative, not consensus.
  Source: `DI-kukuk`.
- [ ] mujad.4 Run the first scored slice using the six existing `TODO-tugoz`
  scenarios once provider-backed scoring is available for this local clone.
  Source: `DI-kukuk`.
- [x] mujad.5 Write a follow-on TE that tests whether CBOR-array extensibility
  plus `pCID`-owned proof semantics change the likely default landing shape
  away from a visible outer future-proofing story and toward a corrected Alt-A
  vs Alt-B decision. Implemented in
  `docs/thought-experiments/TE-pokul-cbor-array-extensibility-pcid-owned-proof-and-grid-envelope-landing-shape.md`.
  Source: `DI-kukuk`; corrected by `DI-vatav`.

## Validation and acceptance criteria

- The new sim clearly states that `[pcid, payload]` are the signed bytes.
- The new sim clearly states that `pcid` owns proof semantics for the third
  slot, including varsig/multisig-style families, without adding a separate
  outer `sig_pcid`.
- `simulations/README.md` and `DEV-GUIDE-RESOURCES.md` place the new sim in an
  alternatives/open-work position rather than presenting it as consensus.
- The TE/DR packet makes the distinction from `TODO-tugoz` explicit so later
  evidence can cite the correct decision owner.
