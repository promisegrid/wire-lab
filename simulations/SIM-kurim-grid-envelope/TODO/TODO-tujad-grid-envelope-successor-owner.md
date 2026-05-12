# TODO-tujad: Grid-envelope successor owner

## Prior aliases

None. This TODO is created directly as a sim-local successor owner under
`rusis.10`.

## Status

Closed for the turn-157/turn-158 successor-owner scope. This TODO remains the
parent seed and historical owner for TE-40 transferred rows `UT-157.a`,
`UT-157.c`, `UT-158.f`, and `UT-158.h`; the concrete successor path is now
materialized as 24 standalone positional grid-envelope simulations under
`DI-fanah`. Successor-owner routing into this TODO was locked under `DI-mosor`
in `protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`.
Seed anchors were established earlier under `DI-nijon`.

## Decision Intent Log

ID: DI-joroh
Date: 2026-05-12 08:44:53
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Resolve the turn-157 grid-envelope replay cleanup by capturing the
candidate envelope inventory and `grid([pcid, payload])` working hypothesis in
this sim-local successor owner, without locking a canonical envelope or creating
a protocol tree.
Intent: Turn 157 contains load-bearing design material that should not remain
only in replay memory. Capturing the alternatives and hypothesis here gives the
grid-envelope lineage a concrete owner for the open replay residue while
preserving the fact that the hypothesis remains unproven.
Constraints: Preserve `Env-1` through `Env-5` as candidate inventory. Do not
decide the final PromiseGrid envelope. Do not create `protocols/grid-envelope.d/`
or draft a grid-envelope spec in this cleanup pass. Leave turn-158 protocol/spec
work open under `tujad.3`.
Affects: `simulations/SIM-kurim-grid-envelope/TODO/TODO-tujad-grid-envelope-successor-owner.md`;
`protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`;
`protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`.

ID: DI-fanah
Date: 2026-05-12 09:22:37
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Split the turn-158 grid-envelope protocol/spec successor work into
24 standalone positional grid-envelope variant simulations instead of choosing
one draft or keeping the work only in this parent owner.
Intent: Turn 158 requires candidate envelopes to behave as specimens that can
evolve and compete independently. The grid-envelope successor therefore needs a
variant matrix that exposes the open design axes without making the harness or
the parent simulation prefer one answer.
Constraints: All envelope shapes are positional. Include the encoding axis
`cbor` versus `dag-cbor`; the unknown-pCID axis `opaque`, `hard-reject`, and
`best-effort`; and the signature axis `wrapper-pcid`, `unsigned-v0`,
`mandatory-opaque-bytes`, and `mandatory-sig-pcid-payload`. Do not create
map-shaped envelope variants, choose a winning variant, edit raw replay logs,
edit historical TE bodies, or make grid-envelope canonical.
Affects: `simulations/README.md`;
`simulations/SIM-kurim-grid-envelope/README.md`;
`simulations/SIM-kurim-grid-envelope/QUESTION.md`;
`simulations/SIM-kurim-grid-envelope/TODO/TODO-tujad-grid-envelope-successor-owner.md`;
`protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`;
`protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`simulations/SIM-*-grid-envelope-enc-<cbor|dag-cbor>-unknown-<opaque|hard-reject|best-effort>-sig-<wrapper-pcid|unsigned-v0|mandatory-opaque-bytes|mandatory-sig-pcid-payload>/`.

## Scope

This TODO owns grid-envelope follow-on that was previously parked under
`TODO-kugod` as "until grid-envelope successor exists":

- candidate envelope inventory ownership;
- `grid([pcid, payload])` working-hypothesis prose ownership;
- turn-158 parallel grid-hypothesis TODO ownership;
- concrete successor planning for grid-envelope protocol directory/spec
  work in this lineage.

Anchor seed note:
`simulations/SIM-kurim-grid-envelope/seed/extraction-sources.md`.

## Candidate Envelope Inventory

Turn 157 named these candidate envelope alternatives as inputs to later
grid-envelope and envelope-bakeoff work. This inventory records the candidates;
it does not select a winner. Source: `DI-joroh`.

- `Env-1`: `grid([pcid, payload])`. A two-element CBOR array where the first
  element is a pCID identifying which protocol, handler, or assertion type
  interprets the second element. This is the current working hypothesis. A
  payload may itself be another `grid([pcid, payload])` value if recursion is
  needed.
- `Env-2`: Promise stack of grid frames. A CBOR sequence of
  `grid([pcid, payload])` frames where stack semantics apply at the sequence
  level and the grid shape applies at each frame. This candidate reconciles the
  grid hypothesis with the earlier TE-famar promise-stack work.
- `Env-3`: Bare CBOR with no shared envelope. Each protocol chooses its own
  message shape. This is maximally permissive, but it may leave the harness
  without a shared parser to exercise across candidate transports.
- `Env-4`: Capability-port triplet. A direct `(promiser, assertion, body)`
  structure with no grid indirection. This is closest to the older
  harness-spec `Promise` shape.
- `Env-5`: Tagged union over `Env-1` and `Env-2`. Single-frame messages use
  `grid([pcid, payload])`; multi-frame messages use a stack of grid frames; a
  top-level tag distinguishes the two cases.

## Grid Envelope Working Hypothesis

`grid([pcid, payload])` is the current working hypothesis for a
transport-agnostic message envelope, but turn 157 explicitly says it has not
been proven. This simulation owns that hypothesis as a candidate specimen to be
tested against alternatives, not as a settled harness rule. Source: `DI-joroh`.

The harness may use this inventory to construct later bakeoffs, but any final
canonical-envelope decision still needs its own TE/DF/DI path. Turn 158's
apparatus-vs-specimen correction remains controlling: the harness compares
candidate envelopes rather than declaring this one canonical in advance.

## Successor Variant Simulations

`DI-fanah` closes the concrete successor path by creating 24 standalone
positional grid-envelope simulations. Each variant carries its own local draft
under `protocols/grid-envelope.d/` so future evolution can compare specimens
without relying on this parent lineage as a shared protocol bundle.

| Simulation | Encoding | Unknown-pCID policy | Signature policy |
|---|---|---|---|
| `../SIM-mahih-grid-envelope-enc-cbor-unknown-opaque-sig-wrapper-pcid/` | CBOR | Opaque store/forward | Wrapper pCID |
| `../SIM-gasus-grid-envelope-enc-cbor-unknown-opaque-sig-unsigned-v0/` | CBOR | Opaque store/forward | Unsigned v0 |
| `../SIM-vutar-grid-envelope-enc-cbor-unknown-opaque-sig-mandatory-opaque-bytes/` | CBOR | Opaque store/forward | Mandatory opaque bytes |
| `../SIM-vamaz-grid-envelope-enc-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload/` | CBOR | Opaque store/forward | Mandatory signature pCID + payload |
| `../SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid/` | CBOR | Hard reject | Wrapper pCID |
| `../SIM-gazan-grid-envelope-enc-cbor-unknown-hard-reject-sig-unsigned-v0/` | CBOR | Hard reject | Unsigned v0 |
| `../SIM-hupir-grid-envelope-enc-cbor-unknown-hard-reject-sig-mandatory-opaque-bytes/` | CBOR | Hard reject | Mandatory opaque bytes |
| `../SIM-kovis-grid-envelope-enc-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload/` | CBOR | Hard reject | Mandatory signature pCID + payload |
| `../SIM-vivus-grid-envelope-enc-cbor-unknown-best-effort-sig-wrapper-pcid/` | CBOR | Best-effort inspection | Wrapper pCID |
| `../SIM-fonig-grid-envelope-enc-cbor-unknown-best-effort-sig-unsigned-v0/` | CBOR | Best-effort inspection | Unsigned v0 |
| `../SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/` | CBOR | Best-effort inspection | Mandatory opaque bytes |
| `../SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/` | CBOR | Best-effort inspection | Mandatory signature pCID + payload |
| `../SIM-gojot-grid-envelope-enc-dag-cbor-unknown-opaque-sig-wrapper-pcid/` | DAG-CBOR | Opaque store/forward | Wrapper pCID |
| `../SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0/` | DAG-CBOR | Opaque store/forward | Unsigned v0 |
| `../SIM-ruzil-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-opaque-bytes/` | DAG-CBOR | Opaque store/forward | Mandatory opaque bytes |
| `../SIM-riliz-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload/` | DAG-CBOR | Opaque store/forward | Mandatory signature pCID + payload |
| `../SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/` | DAG-CBOR | Hard reject | Wrapper pCID |
| `../SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0/` | DAG-CBOR | Hard reject | Unsigned v0 |
| `../SIM-tohol-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-opaque-bytes/` | DAG-CBOR | Hard reject | Mandatory opaque bytes |
| `../SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload/` | DAG-CBOR | Hard reject | Mandatory signature pCID + payload |
| `../SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/` | DAG-CBOR | Best-effort inspection | Wrapper pCID |
| `../SIM-zifik-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-unsigned-v0/` | DAG-CBOR | Best-effort inspection | Unsigned v0 |
| `../SIM-fonol-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/` | DAG-CBOR | Best-effort inspection | Mandatory opaque bytes |
| `../SIM-rakir-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/` | DAG-CBOR | Best-effort inspection | Mandatory signature pCID + payload |

## Subtasks

- [x] tujad.1 Materialize the candidate envelope inventory owner record
  for `UT-157.a`.
- [x] tujad.2 Materialize the `grid([pcid, payload])`
  working-hypothesis owner record for `UT-157.c`.
- [x] tujad.3 Define and track the concrete successor path for
  grid-envelope protocol directory/spec work for `UT-158.f`. Closed by
  `DI-fanah` via the 24 standalone positional successor simulations.
- [x] tujad.4 Back-link resulting decisions and artifacts to
  `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`.
- [x] tujad.5 Record that this TODO satisfies `UT-158.h` as the
  turn-158 parallel grid-hypothesis TODO; `tujad.3` is now closed by
  the `DI-fanah` successor split.
