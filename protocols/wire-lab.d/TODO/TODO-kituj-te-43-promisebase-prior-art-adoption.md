# TODO-kituj: TE-43 promisebase prior-art adoption

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TODO-28` (integer alias)
- `TODO-20260507-002306` (timestamp alias and pre-migration filename)

## Status

Open. TE-sihih's L5/L6/L7 layer model and L6 CAS subtree have landed, so the
old TE-sihih prerequisite is satisfied. This TODO still owns the concrete
promisebase / pitbase prior-art adoption question for the first L6 CAS spec and
remains undrafted. The simulation-facing CAS object-model pressure from turn
177 now lives in `simulations/SIM-jomag-cas-object-model/`, with feed-side
chunk movement pressure in `simulations/SIM-zazit-chunk-feed-replication/`;
their `SCENARIOS.md` files make the object-model and chunk-replication pressure
concrete. Those simulations are inputs to this TODO, not replacements for
TE-43's DF/DI path. Three narrower `DR-tumus` bakeoff sims now test the
starting-profile, object-type-binding, and chunking-identity choices before the
DF is asked again, and their synthesis now lives in the final answerable
`DR-tumus` packet. Source: TE-sihih; TODO-vunub Q-22.6; `DI-navod`;
`DI-pator`; `DI-bukoh`; `DI-molah`. The concrete decision request is now
`DR-tumus`. No twig yet.

## Threads absorbed from OPEN-THREADS.md

### T-PROMISEBASE-ADOPTION (formerly OPEN-THREADS, opened 2026-05-06)

Largest cluster after TE-sihih (25 UTs).

Scope:
  - Ratify the prototype-not-canon stance (Steve turn-191: "prototype
    at best; prefer wire-lab in conflict") and document concrete
    wire-lab-vs-promisebase conflicts before applying the wire-lab default.
    Source: `DI-sapiv`.
  - Document pitbase as L6 substrate prior art with 17/17 chunker-
    merkle tests green (UT-183.c).
  - Specify the concrete CBOR profile for L6 pointer and CAS-object bytes:
    deterministic/canonical encoding, allowed tags, and text-string versus
    byte-string boundaries (UT-177.c).
  - Lock CAS object typing through CIDv1 codec / multicodec choices for raw
    chunks, Merkle nodes, and pointer objects rather than filename suffixes
    (UT-177.h, refined by UT-178.l).
  - Decide chunking algorithm and full parameter set, including the turn-177
    FastCDC/Rabin-size proposal and the later pitbase mismatch.
  - Decide promisebase's human-readable reference-symbol / hash-name-resolution
    lesson separately from CBOR profile, chunking, and CIDv1 object typing:
    adopt it as L6 reference objects, route it to L7 group/session metadata,
    or reject/defer custom syntax in favor of CID-backed pointers. The
    simulation question lives in
    `simulations/SIM-ligan-promisebase-reference-naming/`. Source: `DI-lusum`;
    `DI-tibis`.
  - Account for turn 184's `kv/fs` refactor pressure before choosing a concrete
    promisebase / pitbase dependency target: pinned `db/`, future `db` on
    `kv/fs`, or `kv/fs` plus wire-lab-owned CAS/tree logic. Source: `DI-nulak`.
  - Reconcile the kv branch on the remote (UT-190.d -- undiscovered
    through turn 192) by deciding whether TE-43 / `DR-tumus` evaluates
    promisebase `main`, `kv`, a merged state, or no promisebase branch as the
    prior-art / adoption target. Source: `DI-mivap`.
  - Preserve turn 192's future-facing modernization pressure without treating
    current promisebase design as canon: Steve intends to reference, factor,
    modernize, and use promisebase as one possible PromiseGrid layer, so TE-43 /
    `DR-tumus` must distinguish today's L6 CAS profile from any later
    promisebase modernization or adoption path. Source: `DI-rupuh`.
  - When citing promisebase RFC age or provenance, use exact git-history dates:
    RFC-1003 / RFC-1004 / RFC-1005 trace to 2021-04-28; RFC-1003 and
    RFC-1005 were edited on 2021-07-08; RFC-1006 traces from 2021-06-23
    through 2021-07-06; RFC-1007 image artifacts trace to 2021-07-10 and
    2021-08-13; promisebase `x/message-format.md` introduces the PromiseGrid
    message-format draft on 2025-09-24. Source: `DI-sapiv`.
  - Resolve the Rabin-vs-FastCDC chunking parameter mismatch (UT-181.b:
    pitbase 512 KiB min / 8 MiB max vs turn-177 ~16 KiB average).
  - Use `SIM-jomag-cas-object-model` and
    `SIM-zazit-chunk-feed-replication` as simulation-facing pressure tests for
    the object-model and chunk-movement consequences of TE-43 decisions
    (DI-navod). Their `SCENARIOS.md` files add concrete cases for deterministic
    CBOR agreement, DAG-CBOR interop, CIDv1 object typing, pointer-object
    identity, chunker parameter mismatch, sparse advertisement, partial Merkle
    fetch, corrupt chunks, and carrier independence (DI-pator).
  - Use `SIM-bobud-l6-cas-starting-profile-bakeoff`,
    `SIM-kohad-cas-object-type-binding-bakeoff`, and
    `SIM-gobaz-chunking-identity-bakeoff` before answering `DR-tumus`
    DF-tumus.1 through DF-tumus.3. Their synthesis now recommends minimal
    pointer/raw first profile, CID codec authority with path suffixes rejected,
    raw-only first profile, and promisebase / pitbase as prior art only, while
    leaving the actual decision to Steve (DI-bukoh; DI-molah).
  - Fix the import-path error (UT-181.a: t7a/pitbase ->
    stevegt/promisebase).
  - Address fuse/ test failures, cmd/pb Docker SDK rot, server/daemon
    uncertainty, and the resulting partial-rot / test-threshold question
    (UT-184.e/g). Source: `DI-nulak`.

Blocking: TE-sihih (substrate-agnostic layered model) L6 substrate
definition must land first. The remaining decision gate is `DR-tumus`.

Anchor: turn-191 canon rule; pitbase main + kv branch on
stevegt/promisebase.
Turn 192 adds future-layer pressure, but not a settled adoption mandate.
Disposition-file pointer: `dropped-thread-disposition-20260506.md`
§ TE-43 cluster (25 UTs).

### T-021-CC-Q5 (carried forward, open)

promisebase canonicality: "wire-lab is canon, promisebase is
prototype" -- needs to be either ratified by a TE or formalized as a
canon rule somewhere readers will find. TE-43 is the natural home
since it locks the prototype-not-canon stance.

## Subtasks

- [ ] kituj.1 Answer `DR-tumus` / DF-tumus.1 through DF-tumus.4 before drafting
  TE-43's concrete L6 CAS profile. The current answerable packet is in
  `DR-tumus` after `DI-molah` synthesis. Source: `DI-majib`; `DI-bukoh`;
  `DI-molah`.
- [ ] kituj.2 After `DR-tumus` is decided, write the TE-43 DI(s) for CBOR /
  DAG-CBOR profile, CIDv1 object typing, pointer-object shape, chunking scope,
  and promisebase / pitbase stance.
- [ ] kituj.3 Update `SIM-jomag-cas-object-model` and
  `SIM-zazit-chunk-feed-replication` with forward pointers to the TE-43 result
  once the DR/DI path lands.
- [x] kituj.4 Review `SIM-bobud-l6-cas-starting-profile-bakeoff`,
  `SIM-kohad-cas-object-type-binding-bakeoff`, and
  `SIM-gobaz-chunking-identity-bakeoff`, then revise `DR-tumus` DF-tumus.1
  through DF-tumus.3 into final answerable choices. Done via `DI-molah`.
- [ ] kituj.5 Decide the promisebase human-readable reference-symbol /
  hash-name-resolution lesson separately from CBOR, chunking, CIDv1 object
  typing, and pointer-object-shape decisions, using
  `simulations/SIM-ligan-promisebase-reference-naming/` as the standalone
  question home. Source: `DI-lusum`; `DI-tibis`.
- [ ] kituj.6 Before TE-43 uses promisebase evidence for the L6 CAS stance,
  decide which promisebase tree state is the evidence/adoption target:
  `main`, `kv`, a merged state, or no branch. Source: `DI-mivap`.
- [ ] kituj.7 When drafting TE-43 or answering `DR-tumus`, document each
  material conflict between wire-lab design choices and promisebase prior art,
  state the wire-lab-preferred default, and cite exact RFC/prototype dates
  instead of loose rounded-age claims. Source: `DI-sapiv`.
- [ ] kituj.8 Before drafting TE-43, reconcile the historical promisebase DF
  lists from turns 187, 191, and 192 into the current `DR-tumus` packet; state
  which list supersedes which, clarify or quote Steve's ambiguous `ref` wording,
  and pick a fresh concise twig only when the actual TE-43 work starts. Source:
  `DI-rupuh`.
- [x] kituj.9 Add `implementations/poc6-dag-cbor-interop/` as executable POC
  evidence for the `cas-object-model-dag-cbor-interop` scenario, using real
  IPLD / DAG-CBOR libraries to test CID links, byte strings, tag-42 link
  encoding, stable CIDs, and local evidence without requiring an IPFS daemon.
  Source: `DI-sagos`.

## Question log

- 2026-05-17: `DR-tumus` asks which concrete L6 CAS profile TE-43 should lock:
  CBOR / DAG-CBOR profile, CIDv1 codec object typing, pointer-object shape,
  chunking algorithm and parameters, and promisebase / pitbase adoption stance.
  Source: `DI-davov`.
- 2026-05-17: `DR-tumus` now has an unanswered next-DF packet and acceptance
  criteria for the first TE-43 decision pass. Source: `DI-majib`.
- 2026-05-17: DF-tumus.1 through DF-tumus.3 are routed through three standalone
  bakeoff simulations before final answers are requested. Source: `DI-bukoh`.
- 2026-05-17: The three bakeoff simulations were synthesized into final
  answerable `DR-tumus` choices, but `DR-tumus` remains open for Steve's actual
  decision. Source: `DI-molah`.
- 2026-05-17: Turn 178's promisebase reference-symbol / hash-name-resolution
  loose end is routed here as `kituj.5` and into
  `SIM-ligan-promisebase-reference-naming/QUESTION.md`, separate from CBOR,
  chunking, and CID object-typing work. Source: `DI-lusum`; `DI-tibis`.
- 2026-05-17: Turn 184's full promisebase audit routes `kv/fs` refactor
  pressure, `cmd/pb` Docker SDK rot, FUSE/server/daemon uncertainty, and the
  resulting "what test status is enough?" question here and to `DR-tumus`.
  Source: `DI-nulak`.
- 2026-05-17: Turn 190's branch-enumeration correction routes the promisebase
  `kv` branch target question here, to `DR-tumus`, and to the CAS simulation
  questions. Source: `DI-mivap`.
- 2026-05-17: Turn 191's canon rule routes the promisebase prototype-not-canon
  stance here: TE-43 / `DR-tumus` should discuss material conflicts, prefer
  wire-lab unless Steve locks an exception, and use exact RFC/prototype dates
  when chronology supports the prior-art argument. Source: `DI-sapiv`.
- 2026-05-17: Turn 192's active-prototype framing is routed here: promisebase may
  become one possible PromiseGrid layer after reference/factoring/modernization,
  but TE-43 must keep current promisebase design non-authoritative, account for
  the `kv` branch as potential modernization evidence, and preserve the
  `ref`/DF-list/twig loose ends until the actual TE-43 pass. Source: `DI-rupuh`.

## Decision Intent Log

ID: DI-sagos
Date: 2026-05-26 11:42:44
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Create `implementations/poc6-dag-cbor-interop/` as the first POC
using the new `pocN-{slug}` naming rule. The POC implements the
`scenarios/cas-object-model-dag-cbor-interop/` pressure as standalone Go tests
using real IPLD / DAG-CBOR libraries, not a hand-rolled CBOR encoder.
Intent: Produce cheap executable evidence for whether CID links, byte strings,
and tag-42 link encoding stay compatible with IPLD / IPFS-style tooling without
requiring an IPFS daemon or deciding the final L6 CAS profile.
Constraints: Keep this POC standalone and test-driven. Do not add a Docker demo
or app/kernel flow. Use Promise Theory vocabulary: Alice and Bob make and judge
local promises; the POC records evidence rather than global truth. Update
`implementations/README.md` with the `pocN-{slug}` naming rule and cite this POC
as provisional implementation evidence only.
Affects: `implementations/poc6-dag-cbor-interop/**`;
`implementations/README.md`;
`protocols/wire-lab.d/TODO/TODO-kituj-te-43-promisebase-prior-art-adoption.md`.
Supersedes:
