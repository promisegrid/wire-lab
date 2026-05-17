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
    at best; prefer wire-lab in conflict").
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
    through turn 192).
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

## Decision Intent Log

(Will be populated as DFs lock and product lands.)
