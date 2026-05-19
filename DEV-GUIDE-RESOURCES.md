# PromiseGrid Development Guide Resources

This file is a wire-lab source map for people and LLMs writing the
PromiseGrid Development Guide. It is not the guide. The guide is about
PromiseGrid; wire-lab is the experimental simulation space where many
PromiseGrid design choices are derived, tested, and recorded. Source:
`DI-nunut`.

## Authority model

- Before guide prose stabilizes, use this file to find wire-lab evidence,
  open DRs, and current writer notes for each guide audience. Source:
  `DI-nunut`.
- After guide prose stabilizes, treat the PromiseGrid Development Guide as the
  higher-level developer source of truth for the claims it settles. Wire-lab
  remains provenance and design history. Source: `DI-nunut`.
- Frozen protocol specs remain authoritative by pCID wherever they are hosted.
  If the guide cites a frozen spec by pCID, that spec is authoritative for its
  protocol even if the file lives in wire-lab or another repo. Source:
  `DI-nunut`.
- Do not present draft TEs, TODOs, DRs, or transport specimens as final
  PromiseGrid APIs. Use them only as evidence for why a guide claim exists or
  remains unsettled. Source: `DI-nunut`.
- Treat discussion ledgers and session-status inventories as dated provenance,
  not as current API or layout guidance. Re-check current DIs, DRs, TODOs,
  simulations, and specs before turning an old "not yet committed" list or path
  name into guide prose. Source: `DI-pazum`.
- `docs/thought-experiments/TE-gurov-promise-shaped-artifacts.md`,
  `docs/thought-experiments/TE-vilot-promise-shaped-simulation-artifacts.md`,
  `docs/thought-experiments/TE-hirap-artifacts-as-promisegrid-messages.md`, and
  `docs/thought-experiments/TE-nizor-pahah-implementation-sufficiency.md`
  are open design evidence about whether coordination artifacts should be
  written as PT-style promises or full PromiseGrid-message-shaped artifacts. Do
  not teach "all docs are promises" or "all artifacts are PromiseGrid messages"
  as settled guide prose until their DF questions are answered and a DI locks the
  result. Source: `DI-nunut`.
- `docs/thought-experiments/TE-nizor-pahah-implementation-sufficiency.md`
  is also open design evidence that `simulations/` may become the primary
  wire-lab evidence boundary. Do not teach simulation directory shapes as final
  PromiseGrid node layout; treat them as wire-lab apparatus until a frozen spec
  or guide-side decision says otherwise. Source: `DI-nunut`.
- `docs/thought-experiments/TE-mupoz-root-protocol-migration-scope-under-simulations.md`
  is decided design evidence that root `protocols/` remains only the wire-lab
  harness apparatus home (`protocols/wire-lab.d/`), while candidate PromiseGrid
  protocol drafts move into simulations as specimens until they graduate. Do not
  teach simulation-local protocol specimens as authoritative PromiseGrid protocol
  roots unless a later DF/DI, frozen spec, or guide-side decision explicitly
  graduates them. DF-mupoz.3 is locked as Alt 3.A by `DI-pakid`; the remaining
  Mupoz decisions and the first recovery simulation path are locked by
  `DI-fakin`. Do not treat the old root `transports/` or `proposals/` paths as
  PromiseGrid layout commitments. Source: `DI-nunut`; `DI-fakin`.
- `docs/thought-experiments/TE-dojab-simulation-run-model-and-scenario-result-matrix.md`
  is decided, refined design evidence about root `scenarios/`, root
  `results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md`,
  Codex-only first runs, and 2D scenario-by-simulation comparison matrices.
  `model-id` is a specific provider/model/reasoning slug such as
  `openai-GPT-5.5-xhigh`, and result-path timestamps are UTC per `DI-miror`.
  `DF-dojab.1` through `DF-dojab.5` are locked by `DI-faros`. Treat those
  paths as wire-lab comparison apparatus, not PromiseGrid node layout or final
  guide-facing API structure.
- `scenarios/README.md` and `results/README.md` are the current root skeleton
  contracts for shared scenario entries and per-run result evidence. Existing
  sim-local scenario rows have been mined into root entries under `DI-nanih`;
  application-seed entries remain owned by TODO-dadub. Until real run files
  exist, guide writers should treat these as apparatus templates and scenario
  pressure, not as PromiseGrid API examples. Source: `DI-vabor`; `DI-dimas`;
  `DI-nanih`.
- Active guide prose should use **promise accounting records** for peer-local
  relationship accounting. Alice, Bob, Carol, and other
  peers each keep their own records; wire-lab does not define a central or
  harness-owned ledger. Source: `DI-mugar`.
- The live guide-feedback route is this resource map plus the external
  PromiseGrid Development Guide feedback process. Retired `ppx-dr` and archived
  proposal records remain provenance, not a live guide-feedback mechanism.
  Source: `DI-mugar`; `DI-fakin`.

## How writers should use this file

Start in the PromiseGrid Development Guide repo, then use the audience sections
below to find supporting wire-lab evidence. Each guide claim should be treated
as one of three states:

- **Ready for guide prose:** enough provenance exists to write careful
  PromiseGrid-facing prose now.
- **Write as provisional:** useful for orientation, but final wording should
  warn that design work is still in progress.
- **Blocked by DR:** do not write settled guide prose until the cited DR closes.

## Mupoz Relocation Map

The PromiseGrid Development Guide is about PromiseGrid, not wire-lab. After
Mupoz, guide writers should usually focus on one or more named simulations
because simulations show concrete PromiseGrid design questions under test. Root
wire-lab machinery is still useful provenance, but it should not dominate guide
prose unless the claim is about how wire-lab itself derives evidence. Source:
`DI-fakin`.

The relocation table below exists so humans and LLMs can find moved evidence
without treating old root paths as still-current design commitments. The
reasoning is:

- `protocols/wire-lab.d/` stays rooted because it is harness apparatus, not a
  candidate PromiseGrid protocol specimen. Source: `DI-pakid`.
- Candidate protocol trees move with their TODO queues so each simulation can
  test a coherent repo-like specimen rather than detached draft files. Source:
  `DI-fakin`; `docs/thought-experiments/TE-vipir-protocols-as-simulated-repos-and-binding-layer.md`.
- Concrete transport evidence moves into `world/` because it is active
  simulation state with verifiable message CIDs. Source: `DI-fakin`;
  `simulations/SIM-ludut-wire-lab-devs/seed/wire-lab-devs-draft-migration.md`.
- Legacy proposals move into `archive/` because the old proposal mechanism has
  been replaced by this resource map plus the external guide feedback process.
  Source: `DI-fakin`.
- Simulation-local artifacts are evidence. They become authoritative
  PromiseGrid guidance only after DR/DI/spec/dev-guide handoff. Source:
  `DI-fakin`.
- Turn-164 cleanup closed the group-session/feed-outer freeze-boundary split:
  group-session owns fixed configured `<author-id>/main` membership,
  git-binding scope, and legacy `Message-ID:` compatibility in the specimen;
  feed-outer stays thin and spec freeze does not rewrite historical
  transport/feed data. Treat these as provisional simulation evidence, not
  final PromiseGrid API prose. Source: `DI-rurab`; `DI-bomud`.

| Old path | New path | Guide-writer status | Reasoning |
|---|---|---|---|
| `protocols/group-session.d/` | `simulations/SIM-rakot-group-session/protocols/group-session.d/` | Candidate protocol specimen; cite as provisional evidence only. | Group-session remains an important envelope/session draft, but root `protocols/` is now reserved for `wire-lab.d` apparatus until graduation. |
| `protocols/udp-binding.d/` | `simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/` | Candidate binding specimen; cite as provisional evidence only. | The active lineage path is `udp-feed`; old `udp-binding` naming is historical provenance in retained draft filenames and prior records. |
| `protocols/ppx-dr.d/` | `protocols/wire-lab.d/archive/retired/ppx-dr/protocols/ppx-dr.d/` | Retired proposal/review protocol archive. | Proposal-as-message thinking remains evidence, but the old proposal queue is retired and not a live guide-feedback mechanism. |
| `transports/wire-lab-devs-draft/` | `simulations/SIM-ludut-wire-lab-devs/world/transports/wire-lab-devs-draft/` | Active simulation specimen with verified message CIDs. | The bytes are useful for replay and CID examples; the old root `transports/` path is not a PromiseGrid API. |

Turn 193 is the urgency source for the `wire-lab-devs` dogfood lineage: the
developer group needed message-transport dogfooding ASAP so Steve was not
working solo. Guide writers should point to `SIM-ludut` as dogfood evidence, not
as a final PromiseGrid group layout or CAS migration answer. Source: `DI-vuzot`.
Turn 194 is the naming/context caution for that same lineage: preserve the
current `wire-lab-devs-draft` specimen name when citing dogfood evidence, and
do not let adjacent promisebase framing erase the current wire-lab-devs path.
Source: `DI-fugod`.
| `transports/README.md` | `protocols/wire-lab.d/archive/migrations/SIM-piloh-turns-149-208-recovery/archive/transports/README.md` | Historical wire-lab design note. | Preserve as provenance for why the transport surface existed, not as current guide layout. |
| `proposals/approved/` | `protocols/wire-lab.d/archive/migrations/SIM-piloh-turns-149-208-recovery/archive/proposals/approved/` | Historical review evidence. | Approved proposal records are useful history; current guide feedback belongs outside wire-lab in the guide process. |
| `proposals/pending/` | `protocols/wire-lab.d/archive/migrations/SIM-piloh-turns-149-208-recovery/archive/proposals/pending/` | Historical pending/contested review evidence. | Pending proposal records are not a live queue after the dev-guide resource/feedback process replaced proposals. |
| `protocols/wire-lab.d/` | unchanged | Harness apparatus. | Use for wire-lab process/provenance claims, not as a PromiseGrid app/kernel API. |

## Turn-177 CAS / Feed / Group Simulations

`DI-navod` created four standalone simulation charters so turn 177's CAS,
feed, group-session, and promise-economy conclusions are visible as simulation
specimens instead of living only in replay cleanup notes. These simulations are
design-point workspaces, not final PromiseGrid APIs and not shared protocol
homes. Guide writers may cite them as evidence that the L5/L6/L7 implications
of turn 177 are under active exploration. Do not infer a canonical
grid-envelope, CAS object model, feed wire format, group message shape, or
promise-accounting scheme from their existence. Source: `DI-navod`.

| Simulation | What it is for | Guide-writer status |
|---|---|---|
| `simulations/SIM-jomag-cas-object-model/` | Tests L6 CAS object-model questions: deterministic CBOR / DAG-CBOR, pointer objects, chunking parameters, CIDv1 codec object typing, and promisebase / pitbase prior-art pressure. | Provisional design-point evidence. |
| `simulations/SIM-ligan-promisebase-reference-naming/` | Tests the turn-178 promisebase human-readable reference-symbol / hash-name-resolution question as its own reference-naming lineage, separate from CAS byte-shape decisions. | Provisional design-point evidence. |
| `simulations/SIM-bobud-l6-cas-starting-profile-bakeoff/` | Tests which first L6 CAS starting profile should be evaluated; its synthesis now feeds the answerable `DR-tumus` packet. | Provisional design-point evidence. |
| `simulations/SIM-kohad-cas-object-type-binding-bakeoff/` | Tests CAS object type-binding choices, including path suffixes as a negative control; its synthesis now feeds the answerable `DR-tumus` packet. | Provisional design-point evidence. |
| `simulations/SIM-gobaz-chunking-identity-bakeoff/` | Tests chunking identity choices, including pCID-driven and chunking-CID / cCID-style alternatives as exploratory terms; its synthesis now feeds the answerable `DR-tumus` packet. | Provisional design-point evidence. |
| `simulations/SIM-zazit-chunk-feed-replication/` | Tests the turn-177 inversion where L5 feeds advertise, request, and replicate CAS chunks between sparse sites rather than carrying group messages. | Provisional design-point evidence. |
| `simulations/SIM-jurar-cas-backed-group-session/` | Tests a successor group-session shape where L7 group semantics point at L6 CAS roots and pointer objects without rewriting historical `.txt` specimens. | Provisional design-point evidence. |
| `simulations/SIM-rusap-promise-accounting-records/` | Tests peer-local promise accounting records for pull, keep, advertise, refusal, and 100-year mental-model pressure without defining a central accounting authority. | Provisional design-point evidence. |
| `simulations/SIM-punaz-bgp-class-routing-app/` | Tests turn-178 BGP-class routing-policy application pressure as a standalone L7 app question rather than a generic promise-accounting example. | Provisional design-point evidence. |
| `simulations/SIM-haros-promise-economy-spectrum/` | Tests turn-179 promise-economy mechanism neutrality across peer-local assessment, capability tokens, transferability, floating exchange rates, cryptocurrency-toxicity failure modes, and the turn-184 RFC-1005 test-driven-fabric prior-art seed. TODO owner: TODO-rajig. | Provisional design-point evidence. |

Open turn-177 decision records: `DR-tumus` owns the concrete L6 CAS adoption
question, `DR-gabif` owns additive CAS-backed group-session migration, and
`DR-robon` owns the turn-177 spec-shape requirement question. Guide prose should
treat those areas as unsettled until the DRs close and their linked DIs land.
Each DR now has an unanswered next-DF packet and acceptance criteria for its
owner TODO. `DR-tumus` DF-tumus.1 through DF-tumus.3 were additionally routed
through the three `DI-bukoh` bakeoff sims, then synthesized by `DI-molah` into
the current answerable packet. Source: `DI-davov`; `DI-majib`; `DI-bukoh`;
`DI-molah`.

Turn 178 adds guide-facing pressure but does not make new guide prose settled.
`DI-vaguf` / `DI-lusum` / `DI-tibis` route sparse-CAS, pull-decision,
BGP-class app, group-identity, promisebase reference-naming, multi-repo
site-topology, and capture-resistance narrative points into sims, TODO owners,
and `DR-napum`. The BGP-class app question and promisebase reference-naming
question now have their own standalone sims rather than living under
promise-accounting or CAS-object-model surfaces. Guide writers should treat
those as exploratory evidence until the relevant DRs and guide decisions close.
Source: `DI-vaguf`; `DI-lusum`; `DI-tibis`.

Turn 179 adds a separate economics-model warning. `DI-vabij` routes the
promise-economy spectrum into its own simulation, and `DI-pidag` gives that
simulation TODO-rajig as a discoverable owner. Conditional-release / geofencing
stays with TODO-ralud and `SIM-zarud`; the promisebase wholesale-adoption pivot
is treated as invalidated by later code-first review. Guide writers should not
describe transferable promise tokens, floating rates, or promisebase adoption as
settled PromiseGrid design. Source: `DI-vabij`; `DI-pidag`.

Turn 184 adds promisebase code-audit evidence but still does not settle guide
prose. `DI-nulak` routes RFC-1005's test tree CID + executable tree CID +
cache-on-pass pattern to TODO-rajig / `SIM-haros`, and routes promisebase
`db/`/`kv/fs`/Docker/FUSE partial-rot evidence to TODO-kituj / `DR-tumus`.
Turn 190 adds the corrective `kv` branch target issue: guide writers must not
talk about promisebase as a single-branch prior-art source until `DR-tumus`
decides whether TE-43 evaluates `main`, `kv`, a merged state, or no
promisebase branch. Guide writers should cite that material as prior-art
pressure only until the relevant DR/DI path lands. Turn 191 adds the
prototype-not-canon rule: guide writers must not present promisebase design as
PromiseGrid canon, should document material promisebase-vs-wire-lab conflicts,
should prefer wire-lab unless a later locked decision says otherwise, and should
cite exact 2021 RFC/prototype dates and 2025 PromiseGrid message-format dates
rather than loose rounded-age claims when chronology matters. Turn 192 adds that
promisebase may become one possible PromiseGrid layer after reference,
factoring, and modernization, but guide writers should not describe that as a
settled merge or settled L6 substrate until `DR-tumus` and TODO-dozak close.
Source: `DI-nulak`; `DI-mivap`; `DI-sapiv`; `DI-rupuh`.

## Protocol/Specimen TODO Question Simulations

`DI-pukap` created six standalone question-home simulations for
protocol/specimen TODO questions that were still visible only as TODO prose.
These are unsettled design workspaces, not final PromiseGrid APIs, not active
protocol homes, and not guide-ready normative text. Guide writers may cite them
as evidence that the questions are being explored independently before a TE, DR,
DI, frozen spec, or guide section settles the answer. Source: `DI-pukap`.

| Simulation | Source question | Guide-writer status |
|---|---|---|
| `simulations/SIM-zarud-conditional-release-geofencing/` | Conditional-release, onward-restraint, geofencing, and recursive promise-graph ownership from TODO-ralud. | Provisional design-point evidence. |
| `simulations/SIM-narok-transport-family-bakeoff/` | Future ring, cluster-of-clusters, gossip, and receipts-at-scale transport-family pressure from TODO-sinuv. | Provisional design-point evidence. |
| `simulations/SIM-dihiz-peer-adoption-metadata/` | Peer-level pCID adoption metadata and open-question answer bindings from TODO-nivus. | Provisional design-point evidence. |
| `simulations/SIM-ranib-spec-requirement-sections/` | Promise-vocabulary, 100-year, and layperson/easy-implementation spec-section requirements from TODO-kulih / DR-robon. | Provisional design-point evidence. |
| `simulations/SIM-bohof-group-session-freeze-promise/` | Group-session freeze evidence and `merge-group-transport-spec` promise shape from TODO-bisur. | Provisional design-point evidence. |
| `simulations/SIM-kuful-udp-feed-v0-conformance/` | UDP-feed v0 reference implementation, test-vector, artifact-writer, and ns-3 conformance evidence from TODO-jodon. | Provisional design-point evidence. |

## Grid-Envelope Variant Simulations

`DI-fanah` split the grid-envelope successor path into 24 standalone positional
variant simulations. These are candidate specimens for comparison, not a
preferred envelope family and not final PromiseGrid guide prose. Guide writers
may cite them as evidence that encoding, unknown-pCID handling, and signature
placement are still under active bakeoff. The parent owner lineage is
`simulations/SIM-kurim-grid-envelope/`; the 24 child simulations below are the
actual independently evolvable specimens. Source: `DI-fanah`;
`simulations/SIM-kurim-grid-envelope/TODO/TODO-tujad-grid-envelope-successor-owner.md`.

| Simulation | What it is for | Encoding | Unknown-pCID policy | Signature policy |
|---|---|---|---|---|
| `simulations/SIM-mahih-grid-envelope-enc-cbor-unknown-opaque-sig-wrapper-pcid/` | Tests a positional CBOR envelope that stores unknown protocol payloads opaquely and identifies signature semantics through the wrapper pCID. | CBOR | Opaque store/forward | Wrapper pCID |
| `simulations/SIM-gasus-grid-envelope-enc-cbor-unknown-opaque-sig-unsigned-v0/` | Tests a positional CBOR envelope that stores unknown protocol payloads opaquely and allows an unsigned v0 baseline. | CBOR | Opaque store/forward | Unsigned v0 |
| `simulations/SIM-vutar-grid-envelope-enc-cbor-unknown-opaque-sig-mandatory-opaque-bytes/` | Tests a positional CBOR envelope that stores unknown protocol payloads opaquely and requires opaque signature bytes. | CBOR | Opaque store/forward | Mandatory opaque bytes |
| `simulations/SIM-vamaz-grid-envelope-enc-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload/` | Tests a positional CBOR envelope that stores unknown protocol payloads opaquely and requires a signature pCID plus signature payload. | CBOR | Opaque store/forward | Mandatory signature pCID + payload |
| `simulations/SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid/` | Tests a positional CBOR envelope that rejects unknown protocol payloads and identifies signature semantics through the wrapper pCID. | CBOR | Hard reject | Wrapper pCID |
| `simulations/SIM-gazan-grid-envelope-enc-cbor-unknown-hard-reject-sig-unsigned-v0/` | Tests a positional CBOR envelope that rejects unknown protocol payloads and allows an unsigned v0 baseline. | CBOR | Hard reject | Unsigned v0 |
| `simulations/SIM-hupir-grid-envelope-enc-cbor-unknown-hard-reject-sig-mandatory-opaque-bytes/` | Tests a positional CBOR envelope that rejects unknown protocol payloads and requires opaque signature bytes. | CBOR | Hard reject | Mandatory opaque bytes |
| `simulations/SIM-kovis-grid-envelope-enc-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload/` | Tests a positional CBOR envelope that rejects unknown protocol payloads and requires a signature pCID plus signature payload. | CBOR | Hard reject | Mandatory signature pCID + payload |
| `simulations/SIM-vivus-grid-envelope-enc-cbor-unknown-best-effort-sig-wrapper-pcid/` | Tests a positional CBOR envelope that attempts best-effort inspection of unknown protocol payloads and identifies signature semantics through the wrapper pCID. | CBOR | Best-effort inspection | Wrapper pCID |
| `simulations/SIM-fonig-grid-envelope-enc-cbor-unknown-best-effort-sig-unsigned-v0/` | Tests a positional CBOR envelope that attempts best-effort inspection of unknown protocol payloads and allows an unsigned v0 baseline. | CBOR | Best-effort inspection | Unsigned v0 |
| `simulations/SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/` | Tests a positional CBOR envelope that attempts best-effort inspection of unknown protocol payloads and requires opaque signature bytes. | CBOR | Best-effort inspection | Mandatory opaque bytes |
| `simulations/SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/` | Tests a positional CBOR envelope that attempts best-effort inspection of unknown protocol payloads and requires a signature pCID plus signature payload. | CBOR | Best-effort inspection | Mandatory signature pCID + payload |
| `simulations/SIM-gojot-grid-envelope-enc-dag-cbor-unknown-opaque-sig-wrapper-pcid/` | Tests a positional DAG-CBOR envelope that stores unknown protocol payloads opaquely and identifies signature semantics through the wrapper pCID. | DAG-CBOR | Opaque store/forward | Wrapper pCID |
| `simulations/SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0/` | Tests a positional DAG-CBOR envelope that stores unknown protocol payloads opaquely and allows an unsigned v0 baseline. | DAG-CBOR | Opaque store/forward | Unsigned v0 |
| `simulations/SIM-ruzil-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-opaque-bytes/` | Tests a positional DAG-CBOR envelope that stores unknown protocol payloads opaquely and requires opaque signature bytes. | DAG-CBOR | Opaque store/forward | Mandatory opaque bytes |
| `simulations/SIM-riliz-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload/` | Tests a positional DAG-CBOR envelope that stores unknown protocol payloads opaquely and requires a signature pCID plus signature payload. | DAG-CBOR | Opaque store/forward | Mandatory signature pCID + payload |
| `simulations/SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/` | Tests a positional DAG-CBOR envelope that rejects unknown protocol payloads and identifies signature semantics through the wrapper pCID. | DAG-CBOR | Hard reject | Wrapper pCID |
| `simulations/SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0/` | Tests a positional DAG-CBOR envelope that rejects unknown protocol payloads and allows an unsigned v0 baseline. | DAG-CBOR | Hard reject | Unsigned v0 |
| `simulations/SIM-tohol-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-opaque-bytes/` | Tests a positional DAG-CBOR envelope that rejects unknown protocol payloads and requires opaque signature bytes. | DAG-CBOR | Hard reject | Mandatory opaque bytes |
| `simulations/SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload/` | Tests a positional DAG-CBOR envelope that rejects unknown protocol payloads and requires a signature pCID plus signature payload. | DAG-CBOR | Hard reject | Mandatory signature pCID + payload |
| `simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/` | Tests a positional DAG-CBOR envelope that attempts best-effort inspection of unknown protocol payloads and identifies signature semantics through the wrapper pCID. | DAG-CBOR | Best-effort inspection | Wrapper pCID |
| `simulations/SIM-zifik-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-unsigned-v0/` | Tests a positional DAG-CBOR envelope that attempts best-effort inspection of unknown protocol payloads and allows an unsigned v0 baseline. | DAG-CBOR | Best-effort inspection | Unsigned v0 |
| `simulations/SIM-fonol-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/` | Tests a positional DAG-CBOR envelope that attempts best-effort inspection of unknown protocol payloads and requires opaque signature bytes. | DAG-CBOR | Best-effort inspection | Mandatory opaque bytes |
| `simulations/SIM-rakir-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/` | Tests a positional DAG-CBOR envelope that attempts best-effort inspection of unknown protocol payloads and requires a signature pCID plus signature payload. | DAG-CBOR | Best-effort inspection | Mandatory signature pCID + payload |

## Audience Readiness Matrix

This matrix answers the current guide-writer feedback items `FB-gigit`,
`FB-rivot`, `FB-vitih`, `FB-mulaj`, and `FB-rigod` at wire-lab scope. It is
writer guidance, not a final PromiseGrid product/API freeze. Source:
`DI-zalak`.

| Guide audience | Current readiness | What can be written now | What remains provisional or blocked | Likely first normative citations |
|---|---|---|---|---|
| Laypeople | Ready for careful guide prose | PromiseGrid is designed for long-lived decentralized communities of autonomous/free agents; no central registry is a design constraint; protocol forking is normal; multi-generational durability is a first-order requirement. | Specific promise-accounting scoring, final wire format, app APIs, and kernel shape remain out of layperson settled prose. `DR-napum` remains open for final public wording. | The guide itself after stabilization; wire-lab sources stay provenance. |
| App Devs | Provisional | The minimum current contract is: choose an explicit protocol spec, use its pCID when frozen, let that spec define payload/handler semantics, and publish implementation conformance claims rather than relying on branch paths. | No stable SDK, handler ABI, universal app message API, or app protocol subset is frozen yet. `DR-tuhaz` remains open. | Future frozen `protocols/*/specs/*.md` docs by pCID, plus B-side `CHANGELOG.md` conformance claims. |
| Kernel Devs | Provisional / blocked for final porting instructions | The porting target is not wire-lab. A porter should expect to implement pCID-selected protocol handlers, substrate/binding/session/message layers claimed by the port, and conformance records for those claims. | The first required frozen spec set, runtime expectations, and implementation obligations are not locked yet. `DR-davod` remains open. | Future frozen binding/session/message specs by pCID, implementation `CHANGELOG.md` conformance claims, and guide prose once stabilized. |

## Laypeople

Use this section for the guide's Laypeople / Intro and Laypeople / Goals
sections.

### Current sources

- `README.md` explains wire-lab's role as a simulation harness rather than the
  PromiseGrid guide itself.
- `docs/thought-experiments/TE-dajot-100-year-goal-as-design-constraint.md`
  records the 100-year design constraint that should shape the public story.
- `docs/thought-experiments/TE-dodaf-should-this-design-become-promisegrid-readme.md`
  compares PromiseGrid public README needs against wire-lab findings.
- `docs/thought-experiments/TE-rotim-ostroms-principles-audit.md` is useful
  provenance for governance and commons-management framing.
- `docs/thought-experiments/TE-sigan-generational-handoff.md` is useful
  provenance for multi-generational continuity.
- `docs/essays/congruence-convergence-and-the-grid.md` is useful background
  for explaining why PromiseGrid can host multiple operational traditions.

### Writer notes

- Lead with PromiseGrid's purpose, not wire-lab's method. Wire-lab details
  belong in footnotes, provenance notes, or contributor-facing material.
- Prefer durable human outcomes: long-lived communities, autonomy, repairable
  trust, forkability, and survival across generations.
- Avoid promising a final wire format, app API, or kernel shape in layperson
  prose. Those belong to App Dev and Kernel Dev sections after the relevant
  DRs close.
- Safe settled claims for current guide prose: the 100-year goal, autonomous
  agents, no central registry, protocol forking as normal lifecycle, and
  multi-generational durability. Repairable trust is safe as a design goal;
  specific promise-accounting mechanics remain provisional. Source: `DI-zalak`;
  final public wording remains tracked by `DR-napum`.

## App Devs

Use this section for the guide's App Devs / How to write a grid app section.

### Current sources

- `docs/thought-experiments/TE-nibar-spec-doc-as-promise.md` explains the
  spec-doc-as-promise model.
- `docs/thought-experiments/TE-lozip-congruence-convergence-duality-and-pcid-framing.md`
  explains pCID-selected protocol framing and payload recursion.
- `docs/thought-experiments/TE-zukug-spec-doc-inversion-and-conformance-changelog.md`
  explains spec freezing and conformance reference direction.
- `docs/thought-experiments/TE-liviv-spec-vs-implementation-split.md` explains
  the A-side spec/design versus B-side implementation split.
- `implementations/README.md` records the current local shape for reference
  implementations and conformance claims.
- `simulations/SIM-rakot-group-session/protocols/group-session.d/specs/*.md`
  and `simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/specs/*.md` are
  current draft protocol specimens after the rusis split. App-facing guidance
  should cite frozen pCIDs or guide-side stable prose, not transient simulation
  paths.

### Writer notes

- Explain that app code should target PromiseGrid protocols by spec identity,
  not by whatever wire-lab experiment happened to produce the current evidence.
- Keep pCID language: a pCID identifies the protocol/spec being spoken, while
  payload content and capability promises live at the layer that spec defines.
- The minimum stable app-developer contract is currently a discipline, not an
  SDK: select a protocol spec, cite/use its pCID after freeze, follow that
  spec's payload and handler rules, and publish conformance claims for code.
  Source: `DI-zalak`.
- Do not invent a stable app SDK, handler ABI, or application message API from
  wire-lab notes. Settled app-dev instructions remain blocked by `DR-tuhaz`.
- If the guide cites a frozen protocol spec by pCID, that citation can be
  normative for app developers. Draft paths without pCIDs are provenance or
  provisional orientation only.

## Kernel Devs

Use this section for the guide's Kernel Devs / How to port the infrastructure
section.

### Current sources

- `docs/thought-experiments/TE-havib-apparatus-vs-specimen-carve-out.md`
  separates wire-lab apparatus from candidate protocol specimens.
- `docs/thought-experiments/TE-jikaf-kernel-as-handler-vs-classical-kernel.md`
  compares kernel-as-handler and classical-kernel framing.
- `docs/thought-experiments/TE-sihih-substrate-agnostic-layered-model.md`
  records the current layered model for substrate, CAS, session, and message
  concerns.
- `docs/thought-experiments/TE-nijab-transport-layering-and-freeze-boundaries.md`
  records the lower-layer transport/feed interpretation and freeze-boundary
  hazards.
- `docs/thought-experiments/TE-vipir-protocols-as-simulated-repos-and-binding-layer.md`
  records the protocols-as-simulated-repos and binding-layer model.
- `protocols/wire-lab.d/specs/harness-spec-draft.md` describes the harness
  apparatus used to study candidates.
- `simulations/SIM-ludut-wire-lab-devs/world/transports/wire-lab-devs-draft/`
  contains current transport-message specimen evidence with verified CIDs.
- `protocols/wire-lab.d/archive/migrations/SIM-piloh-turns-149-208-recovery/archive/transports/README.md`
  preserves the old root transport design surface as historical evidence.

### Writer notes

- Porting guidance should describe PromiseGrid infrastructure boundaries, not
  the mechanics of one wire-lab run.
- Use the apparatus/specimen split: harness mechanics are evidence-gathering
  apparatus; per-protocol specs and implementations are the candidate contracts
  that may become porting targets.
- Use promise-accounting wording for peer-local relationship records. Do not
  describe a global or harness-owned accounting ledger.
- Do not tell kernel developers to implement every wire-lab draft artifact.
  Frozen pCID specs, conformance claims, and implementation CHANGELOGs are the
  relevant authority path once they exist.
- The current porting target is a pCID-selected protocol stack plus explicit
  conformance claims, not the wire-lab harness. A port implements the frozen
  binding/session/message specs it claims, and records those claims in an
  implementation CHANGELOG. Source: `DI-zalak`.
- Settled porting guidance is blocked by `DR-davod` until the stable kernel
  boundary, runtime expectations, and conformance obligations are explicit.

## Likely Normative Citation Path

These are early pointers for `FB-rigod`; they are expected citation paths, not
final normative references. Source: `DI-zalak`.

- **Laypeople:** once the PromiseGrid Development Guide stabilizes, cite the
  guide itself for public narrative claims. Use wire-lab TEs and essays only as
  provenance unless the guide explicitly points readers there.
- **App Devs:** likely first normative references are frozen protocol specs by
  pCID, followed by B-side implementation `CHANGELOG.md` conformance claims.
  Until freeze, cite draft specs only as provisional orientation and include
  their current repo path or simulation-local path.
- **Kernel Devs:** likely first normative references are frozen binding,
  session, and message specs by pCID, plus implementation conformance records.
  `simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/specs/udp-binding-draft.md`
  and
  `simulations/SIM-rakot-group-session/protocols/group-session.d/specs/group-session-draft.md`
  are likely early ancestors, but they are not final normative citations until
  frozen.
- **Non-normative provenance:** TEs, DRs, TODOs, simulation archives,
  `implementations/README.md`, and the harness spec explain how decisions were
  derived; they do not become app or porting APIs merely by being informative.

## Open DRs that block settled guide prose

- `DR-napum` — decides which layperson-facing PromiseGrid claims are stable
  enough for settled Intro and Goals prose.
- `DR-tuhaz` — decides the stable app-developer contract or provisional
  fallback for "How to write a grid app."
- `DR-davod` — decides the stable kernel-developer porting boundary and
  conformance target.
- `DR-tumus` — decides the concrete L6 CAS adoption profile exposed by turn
  177; DF-tumus.1 through DF-tumus.3 now include the `DI-molah` synthesis of
  the `DI-bukoh` bakeoff simulations, and DF-tumus.4 now includes the
  promisebase `main`/`kv`/merged/no-branch evidence-target choice plus the
  turn-191 prototype-not-canon / documented-conflict constraint and turn-192
  active-prototype modernization pressure.
- `DR-gabif` — decides additive migration from historical inline group-session
  evidence to CAS-backed group-session specimens; turn 193 adds dogfood urgency
  that the trigger discipline must preserve without rewriting historical bytes.
- `DR-robon` — decides whether turn-177 promise-vocabulary, 100-year, and
  mental-model requirements become required protocol-spec sections.

## Maintenance promises

- Keep this file as a source map plus writer notes. Do not move full guide
  prose here. Source: `DI-nunut`.
