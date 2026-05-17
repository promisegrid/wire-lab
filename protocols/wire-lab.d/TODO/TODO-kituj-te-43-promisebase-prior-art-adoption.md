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
TE-43's DF/DI path. Source: TE-sihih; TODO-vunub Q-22.6; `DI-navod`;
`DI-pator`. The concrete decision request is now `DR-tumus`. No twig yet.

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
  - Fix the import-path error (UT-181.a: t7a/pitbase ->
    stevegt/promisebase).
  - Address fuse/ test failures and cmd/pb Docker SDK rot (UT-184.e/g).

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

## Question log

- 2026-05-17: `DR-tumus` asks which concrete L6 CAS profile TE-43 should lock:
  CBOR / DAG-CBOR profile, CIDv1 codec object typing, pointer-object shape,
  chunking algorithm and parameters, and promisebase / pitbase adoption stance.
  Source: `DI-davov`.

## Decision Intent Log

(Will be populated as DFs lock and product lands.)
