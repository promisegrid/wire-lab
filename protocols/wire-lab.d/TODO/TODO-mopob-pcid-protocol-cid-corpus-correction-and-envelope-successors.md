# TODO-mopob: pCID Protocol-CID corpus correction and envelope successors

## Prior aliases

None. This TODO was minted after the proquint-handle migration.

## Status

Open. This TODO owns the corpus-wide correction that `pCID` means
**Protocol CID**: the CID of the protocol spec document, not the CID of a
payload object. The protocol named by a `pCID` can define the full following
slot vector, including payload shape and proof/signature encoding. The whole
message is encoded as CBOR, so arity is carried by the CBOR array header ahead
of the `[pCID, ...]` array contents.
This TODO plans direct source fixes where safe, successor sims where scored
artifacts must remain immutable, and a source-side eradication pass for the
wrong “payload pCID” framing. `SIM-dalor` has one explicit exception because
its score evidence was produced from the bad pCID wording plus a bad layer-local
rubric; the stale run is removed and replaced rather than preserved as current
evidence. Source: `DI-muniz`; `DI-pozom`.

## Decision Intent Log

ID: DI-sizoh
Date: 2026-05-26 11:13:09
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Add the current envelope/proof hierarchy to `DEV-GUIDE-RESOURCES.md`
and expand `DN-jotob` so guide readers see the relationship between the formal
envelope rule, the payload-first example profile, the signed-message example
profile, single-signer varsig-style proofs, and pCID-defined multisig proof
sets or proof chains.
Intent: Make the proof guidance explicit without making a proof slot,
`varsig`, multisig, or any other proof scheme a universal envelope requirement.
The hierarchy should teach the common path while preserving pCID ownership of
slot roles, signable views, proof encodings, and Promise Theory interpretation.
Constraints: Use "example" rather than "default" for profile wording. Keep
`grid([42(pCID), ...protocol-defined-slots])` as the formal rule. Keep
`grid([42(pCID), payload, ...])` and `grid([42(pCID), payload, proof])` as
examples, not inherited requirements.
Affects: `DEV-GUIDE-RESOURCES.md`;
`docs/research/DN-jotob-grid-envelope-tag42-variable-outer-slots.md`;
this TODO.
Supersedes:

ID: DI-punam
Date: 2026-05-26 10:58:44
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Replace "recommended default profile" wording with "recommended
example profile" for `grid([42(pCID), payload, ...])`. The formal envelope rule
remains `grid([42(pCID), ...protocol-defined-slots])`, and
`grid([42(pCID), payload, ...])` remains the preferred teaching and specimen
shape for ordinary payload-first protocols, but it is no longer described as a
default that protocol authors are expected to inherit.
Intent: Avoid accidental normative pressure that could make payload-first slot
layout look mandatory or privileged. PromiseGrid should show simple examples
without implying that the example constrains pCID-defined slot roles.
Constraints: Do not rewrite historical TE bodies for this wording change. Add
refinement notes where needed, update current developer-facing docs, and keep
slot `0` as `42(pCID)`.
Affects: `docs/thought-experiments/TE-lamun-pcid-defined-slot-vector.md`;
`docs/thought-experiments/TE-fikoj-universal-42-pcid-envelope-shape.md`;
`docs/research/DN-jotob-grid-envelope-tag42-variable-outer-slots.md`;
`docs/research/DN-lujad-promisegrid-kernel-role-profile.md`;
`protocols/wire-lab.d/specs/harness-spec-draft.md`; `DEV-GUIDE-RESOURCES.md`;
this TODO.
Supersedes: DI-rojij only for the "recommended default profile" phrase and
its normative force; preserves DI-rojij's formal slot-vector rule.

ID: DI-rojij
Date: 2026-05-26 10:27:23
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Refine the current PromiseGrid outer-envelope rule from hard
`grid([42(pCID), payload, ...])` to formal
`grid([42(pCID), ...protocol-defined-slots])`, where slot `0` remains
`42(pCID)` and the protocol specification named by `pCID` defines every
following slot. Keep `grid([42(pCID), payload, ...])` as the recommended default
profile for ordinary protocols, and require any protocol that gives slot `1` a
different role to justify and specify that deviation explicitly.
Intent: Preserve `TE-fikoj`'s useful universal slot-0 and CBOR-arity result
while removing an unnecessary universal claim that slot `1` must always be the
payload. The pCID-named spec already owns payload shape, proof encoding,
signable view, and validation rules; letting it define the full slot vector
keeps the envelope simpler and more durable without abandoning the payload-first
profile that most protocols should use.
Constraints: Keep `pCID` = Protocol CID. Keep slot `0` as the current
`42(pCID)` standard instance. Do not treat protocol-defined slots as license for
vague layouts: specs must define slot count, slot order, slot meanings,
signable view, validation order, unsupported behavior, and promise semantics.
Carriage remains distinct from semantic acceptance, and trust remains local.
Affects: `docs/thought-experiments/TE-lamun-pcid-defined-slot-vector.md`;
`docs/thought-experiments/TE-fikoj-universal-42-pcid-envelope-shape.md`;
`docs/research/DN-jotob-grid-envelope-tag42-variable-outer-slots.md`;
`docs/research/DN-lujad-promisegrid-kernel-role-profile.md`;
`protocols/wire-lab.d/specs/harness-spec-draft.md`; `DEV-GUIDE-RESOURCES.md`;
this TODO.
Supersedes: DI-sisak for the hard slot-1 payload rule only; preserves
DI-sisak's universal slot-0, tag-42, and CBOR-arity conclusions.

ID: DI-hozal
Date: 2026-05-24 19:20:42
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Update `DEV-GUIDE-RESOURCES.md` so it explicitly points guide writers
to `docs/research/DN-jotob-grid-envelope-tag42-variable-outer-slots.md` as the
plain-English explainer for the currently locked tagged-selector envelope
direction.
Intent: `DEV-GUIDE-RESOURCES.md` already states the locked current wire
direction, but it does not yet route readers to the short plain-English note
that explains why `grid([42(pCID), payload, ...])` is current, why `42(pCID)`
is the present IPLD / IPFS / Bluesky compatibility choice, and what sim lineage
fed the final lock. The guide resource map should expose that explanatory note
close to the top current-state snapshot.
Constraints: Keep this as a source-map clarification only. Do not change the
locked `TE-fikoj` / `DI-sisak` envelope direction. Do not edit `DN-jotob` in
this pass. Make clear that `DN-jotob` is explanatory guidance, not the locking
authority and not a frozen protocol spec.
Affects: `protocols/wire-lab.d/TODO/TODO-mopob-pcid-protocol-cid-corpus-correction-and-envelope-successors.md`;
`DEV-GUIDE-RESOURCES.md`.

ID: DI-zozov
Date: 2026-05-24 19:00:16
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Update `docs/research/DN-jotob-grid-envelope-tag42-variable-outer-slots.md`
to give explicit heritage credit to the envelope-sim lineage that led to the
current tagged-selector direction, with `SIM-dalor` named as the primary
antecedent and `SIM-pobod`, `SIM-jufag`, and `SIM-zukis` named as the nearby
comparison and specimen lineages.
Intent: The note currently explains the chosen direction and the reason for tag
`42`, but it reads as if that direction appeared fully formed at the TE/DI
layer. The heritage should be visible in plain English so readers can see that
the locked direction inherits concrete pressure and ideas from the earlier sim
lineage, especially `SIM-dalor`'s protocol-owned outer-proof slot framing.
Constraints: Keep this as a doc-only lineage clarification. Do not reopen
`TE-fikoj` or `DI-sisak`. Do not claim that any single sim, including
`SIM-dalor`, was itself the final protocol decision. Make clear that the TE/DI
chain remains the actual locking authority, while the sims provided the design
pressure and near-contender evidence.
Affects: `protocols/wire-lab.d/TODO/TODO-mopob-pcid-protocol-cid-corpus-correction-and-envelope-successors.md`;
`docs/research/DN-jotob-grid-envelope-tag42-variable-outer-slots.md`.

ID: DI-bumon
Date: 2026-05-24 18:52:00
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Rename the plain-English tagged-envelope design note to the
proquint-named path
`docs/research/DN-jotob-grid-envelope-tag42-variable-outer-slots.md` and
clarify its slot-0 explanation so it says explicitly that `42` is the current
tagged-selector instance chosen for IPLD / IPFS / Bluesky ecosystem
compatibility. The note must explain what that compatibility buys and what it
does not imply.
Intent: The current note explains the chosen direction but does not yet say
plainly why `42` is preferred or what “ecosystem compatibility” means in
practical developer terms. The note should make that rationale explicit while
keeping the semantic role distinction intact: slot `0` still means Protocol
CID, while `42(pCID)` is the current wire instance chosen for interop with
CID-bearing DAG-CBOR ecosystems.
Constraints: Keep the note in `docs/research/`. Do not open a new TE, DR, or
sim. Preserve the already locked `TE-fikoj` / `DI-sisak` direction. Make clear
that IPLD / IPFS / Bluesky compatibility buys tooling and ecosystem alignment,
not semantic delegation, global trust, or external authority. Use the proquint
path in new references going forward.
Affects: `protocols/wire-lab.d/TODO/TODO-mopob-pcid-protocol-cid-corpus-correction-and-envelope-successors.md`;
`docs/research/DN-jotob-grid-envelope-tag42-variable-outer-slots.md`.

ID: DI-mabit
Date: 2026-05-24 18:45:00
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Add one direct, breed-eligible specimen simulation,
`SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig`, plus one
plain-English design note at
`docs/research/grid-envelope-tag42-variable-outer-slots-20260524.md` to
implement and explain the already locked `TE-fikoj` / `DI-sisak` direction.
`SIM-zukis` must concretely exercise the family
`grid([42(pCID), payload, ...])` by using one protocol-owned example later outer
slot: slot `2` is a `varsig` proof for this specimen. The broader family rule
remains that later outer-slot roles are defined by the protocol named by
`pCID`, not by the universal envelope.
Intent: The repo now has the locked outer-envelope direction in TE, TODO,
harness-spec, and guide-resource form, but it still lacks one direct specimen
simulation and one developer-facing explanatory note. This task fills both gaps
without reopening the fixed-three-slot versus variable-arity decision. The sim
must make variable arity concrete while staying PT-clean and avoiding
selector-shopping machinery. The design note must explain the current direction
in plain English for kernel and app developers.
Constraints: Keep `pCID` = Protocol CID. Do not create `SIM-META.json`; the new
sim remains breed-eligible. Do not add `sig_pcid`, `env_pcid`, `payload_pcid`,
generic claim headers, statement capsules, or promise-accounting envelope
fields. Do not write a new TE. Keep the design note outside the TE corpus and
outside the formal spec tree. Cite `TE-fikoj`, `DI-sisak`, and the current
`DEV-GUIDE-RESOURCES.md` snapshot from the design note.
Affects: `protocols/wire-lab.d/TODO/TODO-mopob-pcid-protocol-cid-corpus-correction-and-envelope-successors.md`;
`simulations/SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig/QUESTION.md`;
`simulations/SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig/README.md`;
`simulations/SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig/protocols/grid-envelope.d/specs/grid-envelope-draft.md`;
`docs/research/grid-envelope-tag42-variable-outer-slots-20260524.md`.

ID: DI-sisak
Date: 2026-05-24 18:06:40
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Lock the current PromiseGrid outer-envelope direction as the CBOR
array `grid([42(pCID), payload, ...])`. Slot 0 belongs to a tagged CBOR
protocol-selector family; tag `42` is the current standard instance. CBOR array
length carries arity. The protocol named by `pCID` defines the roles and
validity rules of later outer slots. Fixed three-slot
`grid([42(pCID), payload, sig])` remains a considered alternative, but is not
the chosen current direction.
Intent: `TE-fikoj` narrowed the choice to fixed three-slot versus variable
arity once slot 0 was universalized as `42(pCID)`. The user chose the
variable-arity branch and the tagged-selector-family framing. This locks the
current direction so the harness spec and dev-guide resource map stop
presenting the older minimal `[pcid, payload]` wording as the dominant current
shape. Treating tag `42` as the current family instance also leaves room for a
future selector-tag successor without inventing `grid2()` or `grid3()` outer
families.
Constraints: Keep `pCID` = Protocol CID. Keep the outer envelope as CBOR. Do
not rewrite `TE-fikoj`'s body; update only its `## Status` and
`## Decision status` fields. Keep Promise Theory vocabulary: carriage is not
semantic acceptance, and later-slot meaning remains protocol-owned rather than
a generic claim-card map. Update `DEV-GUIDE-RESOURCES.md` in the same pass so
the current design snapshot reflects the locked direction while still
identifying fixed three-slot as a considered but unchosen alternative.
Affects: `protocols/wire-lab.d/TODO/TODO-mopob-pcid-protocol-cid-corpus-correction-and-envelope-successors.md`;
`docs/thought-experiments/TE-fikoj-universal-42-pcid-envelope-shape.md`;
`protocols/wire-lab.d/specs/harness-spec-draft.md`;
`DEV-GUIDE-RESOURCES.md`.
Supersedes: DI-kafat

ID: DI-kafat
Date: 2026-05-24 17:06:23
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Add a new superseding TE, `TE-fikoj`, under `TODO-mopob` instead of
editing `TE-tilir` in place. `TE-fikoj` must test the stronger premise that the
outer PromiseGrid envelope standard itself may freeze slot 0 as `42(pCID)`,
then compare fixed three-slot `grid([42(pCID), payload, sig])` versus
variable-arity `grid([42(pCID), payload, ...])` under that premise. If both
branches remain plausible after scenario analysis, the TE should lean variable
arity.
Intent: The new user position is stronger than `TE-tilir`: it treats the
tag-42 wrapper as acceptable durable boilerplate, values native interop with
the CID / DAG-CBOR ecosystem highly, and no longer finds "carrier profile" a
helpful distinction for this decision. That is a substantive change under the
TE editing policy, so it must land as a new superseding TE. The new TE should
test whether CBOR-array arity plus pCID-defined outer slot roles is a better
long-horizon story than freezing an exactly-three-slot outer contract too early.
Constraints: Do not rewrite `TE-tilir`'s body; update only its `## Status` and
`## Decision status` fields. Keep `pCID` = Protocol CID. Do not regress into
payload-CID wording. Keep `TE-vujaj` alive as the wording/history facet. Do not
erase `TE-pokul`'s older anti-unnecessary-outer-slot pressure, but test whether
the stronger universal-`42(pCID)` premise changes the landing shape. Do not
update `DEV-GUIDE-RESOURCES.md` or simulation lineages in this task.
Affects: `protocols/wire-lab.d/TODO/TODO-mopob-pcid-protocol-cid-corpus-correction-and-envelope-successors.md`;
`docs/thought-experiments/TE-fikoj-universal-42-pcid-envelope-shape.md`;
`docs/thought-experiments/TE-tilir-slot0-ipld-interop-and-bootstrap-supersedence.md`;
`docs/thought-experiments/README.md`;
`protocols/wire-lab.d/specs/harness-spec-draft.md`.

ID: DI-lasah
Date: 2026-05-24 16:17:29
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Add a new superseding TE, `TE-tilir`, under `TODO-mopob` instead of
editing `TE-titaj` in place. `TE-tilir` must add an explicit ATProto / Bluesky
ecosystem-interop scenario and a byte-level bootstrap comparison between raw
CID bytes and DAG-CBOR `42(pCID)` Link form, then rebalance the conclusion so
that PromiseGrid keeps `pCID` as the semantic role of slot 0 while treating
`42(pCID)` as the strongly preferred representation when the carrier profile is
DAG-CBOR / IPLD.
Intent: The requested change is substantive under the TE editing policy, so it
must land as a new TE with supersedence markers rather than as an in-place
rewrite of `TE-titaj`. The earlier TE was directionally right about separating
role from representation, but it understated the practical interop value of
aligning with the wider CID / DAG-CBOR ecosystem and overstated the byte-level
burden of `42(pCID)` on constrained receivers.
Constraints: Do not rewrite `TE-titaj`'s body; update only its `## Status` and
`## Decision status` fields. Do not reopen `TE-pokul`'s two-slot-versus-three-
slot family question. Do not collapse `42(pCID)` into the timeless semantic
meaning of slot 0. Keep `TE-vujaj` alive as the wording-focused facet. Do not
update `DEV-GUIDE-RESOURCES.md` or simulation lineages in this task.
Affects: `protocols/wire-lab.d/TODO/TODO-mopob-pcid-protocol-cid-corpus-correction-and-envelope-successors.md`;
`docs/thought-experiments/TE-tilir-slot0-ipld-interop-and-bootstrap-supersedence.md`;
`docs/thought-experiments/TE-titaj-pcid-slot0-bootstrap-across-decades.md`;
`docs/thought-experiments/README.md`;
`protocols/wire-lab.d/specs/harness-spec-draft.md`.

ID: DI-pikos
Date: 2026-05-24 15:11:21
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Add a new standalone TE, `TE-titaj`, under `TODO-mopob` to test
slot-0 standardization from the receiver's point of view across decades,
devices, runtimes, and identifier-ecosystem drift. The TE compares
`grid([42(pCID), payload, sig])`, `grid([pCID, payload, sig])`, and a corrected
bootstrap-oriented `grid([slot0_bytes, payload, sig])` alternative in which the
PromiseGrid envelope standard, not the protocol spec, defines how a receiver
recovers `pCID` from slot-0 bytes.
Intent: The earlier wording TE (`TE-vujaj`) established that `42(...)` is a
representation question, not a role question, but it did not fully pressure the
receiver bootstrap problem over a 100-year horizon. This TE is needed because a
future receiver may know only CBOR plus local archives, while IPLD, CID, or
multihash conventions may survive, fragment, or fade. The repo needs a
receiver-centric record of what the lowest common denominator actually is before
locking a slot-0 standard.
Constraints: Keep `pCID` = Protocol CID throughout. Do not let the protocol spec
bootstrap slot 0; that is circular. Keep `TE-pokul` as the owner of the
two-slot-versus-three-slot family question. Keep the TE status at `needs DF`.
Do not update `DEV-GUIDE-RESOURCES.md` or simulation lineages in this task.
Affects: `protocols/wire-lab.d/TODO/TODO-mopob-pcid-protocol-cid-corpus-correction-and-envelope-successors.md`;
`docs/thought-experiments/TE-titaj-pcid-slot0-bootstrap-across-decades.md`;
`docs/thought-experiments/README.md`;
`protocols/wire-lab.d/specs/harness-spec-draft.md`.

ID: DI-falap
Date: 2026-05-24 14:09:36
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Add a standalone TE, `TE-vujaj`, under `TODO-mopob` to test the
wording boundary between conceptual PromiseGrid envelope shape and
DAG-CBOR-specific link encoding. The TE will ask whether PromiseGrid should say
the envelope shape is `grid([42(pCID), payload, sig])`, or instead say the
conceptual shape is `grid([pCID, payload, sig])` and reserve `42(pCID)` for
DAG-CBOR encoding examples and profile-specific specimens.
Intent: The current corpus has already locked that `pCID` means Protocol CID and
that `42(...)` is a DAG-CBOR link representation, but it has not yet written a
repo-owned TE focused on whether using `42(pCID)` in the canonical envelope
phrase would collapse encoding-level detail into the semantic shape. A separate
TE is safer than rewriting `TE-pokul` because `TE-pokul` already owns the
two-slot-versus-three-slot envelope-family question.
Constraints: Keep `pCID` as Protocol CID throughout. Do not pre-lock the TE's
conclusion in the DI: deciding whether the preferred phrase is
`grid([pCID, payload, sig])`, `grid([42(pCID), payload, sig])`, or a stricter
conceptual-versus-wire distinction is the purpose of `TE-vujaj`. Do not update
`DEV-GUIDE-RESOURCES.md` or active simulation lineages in this task. Keep the TE
status at `needs DF` unless the TE itself truly closes the question under
existing evidence.
Affects: `protocols/wire-lab.d/TODO/TODO-mopob-pcid-protocol-cid-corpus-correction-and-envelope-successors.md`;
`docs/thought-experiments/TE-vujaj-grid-envelope-dag-cbor-link-wording.md`;
`docs/thought-experiments/README.md`;
`protocols/wire-lab.d/specs/harness-spec-draft.md`.

ID: DI-muniz
Date: 2026-05-24 23:59:00
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Track the `pCID` meaning correction as one harness TODO. The locked
meaning is: `pCID` means Protocol CID, i.e. the CID of the protocol spec
document. A `pCID` is not now and never was the CID of the payload. The
protocol named by a `pCID` may define payload shape, signable view, proof
encoding, verification rules, and related semantics. The whole message uses
CBOR; array arity is carried by the CBOR array header, not by extra outer
future-proofing fields. Safe non-evidence docs may be corrected in place.
Scored sims and committed result artifacts must not be rewritten in place;
instead, they are superseded by successor sims or clarified by later source
artifacts.
Intent: The current corpus still contains wrong “payload pCID” phrasing in
active sim docs and already-scored sim lineages. That mistake distorts the
envelope debate by making `pCID` sound like a payload-object hash instead of a
protocol selector. The correction must be explicit, repo-owned, and applied
consistently so future sims, scoring, and dev-guide prose stop inheriting the
error.
Constraints: Do not rewrite committed `results/` evidence. Do not rewrite
already-scored canonical sim trees in place; use successor sims instead. Keep
historical evidence readable, but route future breeding and comparison work to
the corrected successor lineages. Preserve the distinction between direct safe
doc fixes and superseded scored specimens.
Affects: `protocols/wire-lab.d/TODO/TODO-mopob-pcid-protocol-cid-corpus-correction-and-envelope-successors.md`;
`protocols/wire-lab.d/TODO/TODO.md`; successor `simulations/SIM-*/` trees for
the corrected envelope lineages; future `simulations/README.md`; future
`DEV-GUIDE-RESOURCES.md`; related source docs that still say or imply “payload
pCID”.

ID: DI-pozom
Date: 2026-05-24 11:54:14
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Apply a one-time `SIM-dalor` exception to the usual scored-artifact
immutability rule. Edit `SIM-dalor-grid-envelope-protocol-owned-signature-slot`
in place so `pCID` is consistently Protocol CID, the envelope-level signed
promise is explicit, and promise-accounting remains a higher-layer payload
concern. Remove the stale `20260524-070634` canary run evidence and the
uncommitted `20260524-174215` dalor-breed evidence, then replace them with a
fresh forced dalor-breed run under the corrected source and rubric.
Intent: The low dalor scores were not reliable design evidence. They were
created from two coupled mistakes: source text that said “payload pCID” and a
scorer rubric that let higher-layer promise-accounting scenarios penalize an
envelope-layer design for correctly leaving accounting semantics to the payload
protocol. A dalor envelope is itself a signed promise that the payload bytes are
shaped according to the protocol specification named by the `pCID`; peers still
make local trust judgments from that evidence. Replacing this narrow run keeps
future GA breeding from suppressing a valid envelope candidate because of known
bad inputs.
Constraints: This exception applies only to `SIM-dalor` and the named stale
runs. It does not authorize broad in-place rewrites of other scored sims or
historical evidence. The scorer must evaluate candidates at their claimed
layer: envelope sims may score highly for PT cleanliness when they express the
envelope-level signed promise, even if higher-layer promise accounting belongs
inside the payload protocol. Scenario fit may still be lower when a scenario
asks for higher-layer semantics the envelope intentionally does not provide.
Affects: `protocols/wire-lab.d/TODO/TODO-mopob-pcid-protocol-cid-corpus-correction-and-envelope-successors.md`;
`simulations/SIM-dalor-grid-envelope-protocol-owned-signature-slot/`;
`tools/ga-runner/score.go`; `tools/ga-runner/result.go`;
`tools/ga-runner/run-canary.sh`; stale `results/**/*20260524-070634.json`;
stale `results/jobs/ga-canary-20260524-pt-gate-v3-successors-plus-dalor/`;
stale `results/state/ga-canary-20260524-pt-gate-v3-successors-plus-dalor.json`;
uncommitted `results/**/*20260524-174215.json`;
uncommitted `results/jobs/ga-canary-20260524-dalor-breed/`;
uncommitted `results/state/ga-canary-20260524-dalor-breed.json`.
Supersedes: DI-muniz for this narrow `SIM-dalor` stale-run replacement only.

## Scope

This TODO covers:

- locating and classifying every remaining source-side use of the wrong
  “payload pCID” framing;
- fixing safe non-evidence docs in place;
- superseding scored or otherwise lineaged sim docs with corrected successor
  sims, except for the explicit `SIM-dalor` stale-run replacement authorized by
  `DI-pozom`;
- updating sim-index and guide prose to route readers to the corrected
  successor lineages;
- keeping committed results immutable while documenting how to interpret older
  wording safely.

This TODO does not normally rewrite committed `results/` artifacts, does not
normally mutate already-scored canonical sims in place, and does not by itself
lock the final winning envelope shape. The only current exception is the
`SIM-dalor` stale-run replacement authorized by `DI-pozom`.

## Predecessor context

- `AGENTS.md` already locks the repo-wide definition: `pCID` means Protocol CID,
  not a message, payload, or promise-object hash.
- `TODO-golad` already recorded the same correction explicitly for the ppx-main
  integration branch.
- `TE-pokul` narrows the envelope question to minimal `[pCID, payload]` versus
  `[pCID, payload, proof]` where the protocol named by `pCID` defines the proof
  encoding and current-sender proof semantics. It rejects visible outer
  future-proofing as the reason to add fields and rejects the earlier false
  Alt-B / Alt-C split. Source: `DI-vatav`.
- Current source-side error occurrences remain in active sim docs such as
  `SIM-dalor`, `SIM-pamap`, `SIM-pobod`, `SIM-maraz`, `SIM-janov`,
  `SIM-natim`, and `SIM-fitin`.

## Error model

The wrong framing has three recurring forms:

1. text that literally says or implies “payload pCID” as if the `pCID` were the
   CID of the payload object;
2. text that says proof semantics are “owned by the payload pCID” instead of by
   the protocol spec named by the outer or nested `pCID`;
3. text that treats variable arity or future-proofing as a reason to add outer
   tuple machinery even though the whole message is already CBOR and the array
   header carries arity.

## Correction rule

All corrected successor artifacts should say the same thing plainly:

- `pCID` = Protocol CID = CID of the protocol spec document;
- the protocol named by a `pCID` defines slots `1..N`, including payload shape,
  signable projection, proof/signature encoding, and verification rules;
- the whole message is CBOR;
- CBOR array arity is already in the array header;
- extra outer fields need a semantic reason, not an extensibility reason.
- do not create a separate “outer attestation” successor family from the false
  Alt-B / Alt-C split; if `[pCID, payload, proof]` is specified so the protocol
  named by `pCID` defines `proof` as the sender's evidence over the message,
  that is the corrected three-slot family. Source: `DI-vatav`.

## Current known source-side hits

The following active source artifacts currently need direct correction or
supersession:

- `simulations/SIM-dalor-grid-envelope-protocol-owned-signature-slot/protocols/grid-envelope.d/specs/grid-envelope-draft.md`
- `simulations/SIM-pamap-grid-envelope-signable-view-atproto-like/QUESTION.md`
- `simulations/SIM-pobod-grid-envelope-outer-promise-nested-signed-payload/README.md`
- `simulations/SIM-maraz-grid-envelope-signed-summary-header-nested-schema/QUESTION.md`
- `simulations/SIM-janov-grid-envelope-layer-pcid-nested-signed-payload/README.md`
- `simulations/SIM-natim-grid-envelope-nested-payload-outer-attestation-multisig/protocols/grid-envelope.d/specs/grid-envelope-draft.md`
- `simulations/SIM-fitin-gordian-payload-wrapper/protocols/grid-envelope.d/specs/grid-envelope-draft.md`

## Subtasks

- [ ] mopob.1 Run a source-side grep/classification pass for every remaining
  “payload pCID” / equivalent phrase in active docs, sims, and guide prose, and
  classify each hit as direct-fix, successor-needed, or historical-evidence-only.
  Source: `DI-muniz`.
- [ ] mopob.2 Record a short repo-owned clarification note, if needed, stating
  how to read older scored artifacts that still use the wrong shorthand without
  treating them as current terminology. Source: `DI-muniz`.
- [x] mopob.3 Correct `SIM-dalor` in place under the one-time stale-run
  exception. Keep the three-slot probe but state explicitly that the outer
  `pCID` names the protocol spec, that the protocol may define proof encoding,
  that CBOR already carries arity, and that the envelope-level signature is
  evidence of the sender's promise that the payload is shaped according to the
  protocol named by the `pCID`. Source: `DI-muniz`; `DI-pozom`.
- [ ] mopob.4 Create a corrected successor for `SIM-pamap` that states the
  signable-view/proof-slot rule as part of the protocol named by the `pCID`,
  not as something owned by a “payload pCID”. Source: `DI-muniz`.
- [ ] mopob.5 Create a corrected successor for `SIM-pobod` that states the
  nested body carries an inner protocol `pCID`, payload bytes, and proof, not
  an “actual payload pCID” in the object-hash sense. Source: `DI-muniz`.
- [ ] mopob.6 Create a corrected successor for `SIM-maraz` that frames the
  signed summary header as a projection over protocol-defined nested structure,
  not “inner payload pCID-defined” structure. Source: `DI-muniz`.
- [ ] mopob.7 Leave `SIM-janov` historical and route future work to a corrected
  successor lineage instead of rewriting `janov` in place. Source: `DI-muniz`.
- [ ] mopob.8 Create a corrected successor for `SIM-natim` that clearly
  separates the outer protocol `pCID`, nested protocol `pCID`, and any proof
  profile selectors without calling any of them payload-object identifiers.
  Source: `DI-muniz`.
- [ ] mopob.9 Create a corrected successor for `SIM-fitin` that replaces
  “support this payload pCID” with “support the protocol named by this `pCID`”
  and keeps the Gordian comparison lineage intact. Source: `DI-muniz`.
- [ ] mopob.10 Update `simulations/README.md` and any relevant guide prose so
  readers are pointed at the corrected successors and the old mistaken lineages
  are treated as historical evidence or non-breeding question homes. Source:
  `DI-muniz`.
- [ ] mopob.11 Decide, per lineage, whether the corrected predecessor should be
  demoted to `question-home`, `bakeoff`, or other non-breeding status after its
  successor exists. Source: `DI-muniz`.
- [ ] mopob.12 Run a focused rescoring slice for the corrected successors
  against the same envelope scenarios so predecessor-vs-successor deltas are
  comparable. Source: `DI-muniz`.
- [ ] mopob.13 After successor creation, rerun the source grep and eradicate or
  supersede every remaining active source-side mention of the error. Leave only
  historical evidence that is intentionally frozen. Source: `DI-muniz`.
- [x] mopob.14 Add `TE-vujaj` as a standalone thought experiment for the wording
  boundary between conceptual envelope shape and DAG-CBOR link encoding. The TE
  belongs under `TODO-mopob` because it clarifies how Protocol-CID semantics
  should be described without silently recasting `42(...)` as the semantic
  envelope shape. Source: `DI-falap`.
- [x] mopob.15 Add `TE-titaj` as a receiver-centric thought experiment for
  slot-0 bootstrap across decades, devices, runtimes, and identifier-ecosystem
  drift. The TE should test what future receivers can still parse locally if
  IPLD survives, if IPLD fades, or if CID/multihash conventions themselves
  become archaeology, and it should keep the bootstrap rule outside the
  protocol spec itself. Source: `DI-pikos`.
- [x] mopob.16 Add `TE-tilir` as the substantive superseding TE for `TE-titaj`
  so the corpus can incorporate ATProto / Bluesky ecosystem interop and the
  small-byte-delta bootstrap packet without violating the TE editing policy.
  Keep `pCID` as the semantic role, but rebalance the recommendation so
  `42(pCID)` becomes the strongly preferred representation when the carrier
  profile is DAG-CBOR / IPLD. Source: `DI-lasah`.
- [x] mopob.17 Add `TE-fikoj` as the substantive superseding TE for `TE-tilir`
  so the corpus can test the stronger universal-`42(pCID)` premise directly.
  Compare fixed three-slot versus variable-arity outer arrays, and if both
  survive, lean toward variable arity with pCID-defined slot roles. Source:
  `DI-kafat`.
- [x] mopob.18 Lock `TE-fikoj`'s recommended current direction as
  `grid([42(pCID), payload, ...])`, with slot 0 in the tagged CBOR
  protocol-selector family, tag `42` as the current standard instance, CBOR
  array length carrying arity, and later outer-slot roles defined by the
  protocol named by `pCID`. Update the TE status, harness evidence pointer, and
  `DEV-GUIDE-RESOURCES.md` in the same pass. Source: `DI-sisak`.
- [x] mopob.19 Add direct specimen
  `SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig` to implement the
  locked `TE-fikoj` direction as one concrete
  `grid([42(pCID), payload, varsig])` member of the broader
  `grid([42(pCID), payload, ...])` family. Keep slot `2` protocol-owned, keep
  the sim breed-eligible, and compare it against `SIM-dalor`, `SIM-pobod`, and
  `SIM-jufag`. Source: `DI-mabit`.
- [x] mopob.20 Add plain-English design note
  `docs/research/grid-envelope-tag42-variable-outer-slots-20260524.md` to
  explain the locked `TE-fikoj` / `DI-sisak` direction to developers without
  TE scaffolding or formal-spec tone. Source: `DI-mabit`.
- [x] mopob.21 Rename the plain-English design note to
  `docs/research/DN-jotob-grid-envelope-tag42-variable-outer-slots.md` and make
  its slot-0 explanation state explicitly that `42` is the current
  IPLD / IPFS / Bluesky ecosystem-compatibility choice, including what that
  compatibility buys and what it does not imply. Source: `DI-bumon`.
- [x] mopob.22 Update `DN-jotob` to credit its design heritage explicitly,
  especially `SIM-dalor` as the primary ancestor, while also naming `SIM-pobod`,
  `SIM-jufag`, and `SIM-zukis` as the surrounding lineage. Keep `TE-fikoj` /
  `DI-sisak` as the locking authority rather than attributing the final decision
  to any one sim. Source: `DI-zozov`.
- [x] mopob.23 Update `DEV-GUIDE-RESOURCES.md` so the top current-state envelope
  snapshot explicitly points guide writers to `DN-jotob` as the plain-English
  explainer for the locked `grid([42(pCID), payload, ...])` direction, while
  keeping `TE-fikoj` / `DI-sisak` as the locking authority. Source: `DI-hozal`.
- [x] mopob.24 Add `TE-lamun` and lock `DI-rojij` so the formal outer-envelope
  rule becomes `grid([42(pCID), ...protocol-defined-slots])`, with
  `grid([42(pCID), payload, ...])` retained as the recommended example profile
  rather than a universal slot-1 law. Source: `DI-rojij`.
- [x] mopob.25 Lock `DI-punam` so current guidance says "recommended example
  profile" rather than "recommended default profile" for
  `grid([42(pCID), payload, ...])`. Source: `DI-punam`.
- [x] mopob.26 Add the envelope/proof hierarchy to
  `DEV-GUIDE-RESOURCES.md` and the fuller `DN-jotob` explanation, using example
  profile language for payload and signed-message shapes. Source: `DI-sizoh`.

## Validation and acceptance criteria

- Every new corrected successor says explicitly that `pCID` means Protocol CID,
  not payload CID.
- Every corrected successor states that the protocol named by a `pCID` may
  define both payload shape and proof/signature encoding.
- Every corrected successor states that the message is CBOR and that array
  arity is carried by the CBOR header.
- No new source artifact reintroduces the phrase “payload pCID” or equivalent
  wording in active prose.
- Already-scored sims and committed `results/` artifacts remain immutable; any
  correction path for them is successor-based or clarified by later notes.
- `simulations/README.md` and related guide prose route readers to the
  corrected lineage once successors exist.
