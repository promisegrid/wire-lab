# TODO-mopob: pCID Protocol-CID corpus correction and envelope successors

## Prior aliases

None. This TODO was minted after the proquint-handle migration.

## Status

Open. This TODO owns the corpus-wide correction that `pCID` means
**Protocol CID**: the CID of the protocol spec document, not the CID of a
payload object. The protocol named by a `pCID` can define both payload shape
and proof/signature encoding. The whole message is encoded as CBOR, so arity is
carried by the CBOR array header ahead of the `[pCID, ...]` array contents.
This TODO plans direct source fixes where safe, successor sims where scored
artifacts must remain immutable, and a source-side eradication pass for the
wrong “payload pCID” framing. `SIM-dalor` has one explicit exception because
its score evidence was produced from the bad pCID wording plus a bad layer-local
rubric; the stale run is removed and replaced rather than preserved as current
evidence. Source: `DI-muniz`; `DI-pozom`.

## Decision Intent Log

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
- the protocol named by a `pCID` defines payload shape, signable projection,
  proof/signature encoding, and verification rules;
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
