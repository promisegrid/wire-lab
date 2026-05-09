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

## How writers should use this file

Start in the PromiseGrid Development Guide repo, then use the audience sections
below to find supporting wire-lab evidence. Each guide claim should be treated
as one of three states:

- **Ready for guide prose:** enough provenance exists to write careful
  PromiseGrid-facing prose now.
- **Write as provisional:** useful for orientation, but final wording should
  warn that design work is still in progress.
- **Blocked by DR:** do not write settled guide prose until the cited DR closes.

## Audience Readiness Matrix

This matrix answers the current guide-writer feedback items `FB-gigit`,
`FB-rivot`, `FB-vitih`, `FB-mulaj`, and `FB-rigod` at wire-lab scope. It is
writer guidance, not a final PromiseGrid product/API freeze. Source:
`DI-zalak`.

| Guide audience | Current readiness | What can be written now | What remains provisional or blocked | Likely first normative citations |
|---|---|---|---|---|
| Laypeople | Ready for careful guide prose | PromiseGrid is designed for long-lived decentralized communities of autonomous/free agents; no central registry is a design constraint; protocol forking is normal; multi-generational durability is a first-order requirement. | Specific trust-ledger scoring, final wire format, app APIs, and kernel shape remain out of layperson settled prose. `DR-napum` remains open for final public wording. | The guide itself after stabilization; wire-lab sources stay provenance. |
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
  specific trust-ledger mechanics remain provisional. Source: `DI-zalak`;
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
- `protocols/*/specs/*.md` are the current draft or frozen protocol specs that
  app-facing guidance may eventually cite.

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
- `transports/README.md` describes the current transport simulation surface.

### Writer notes

- Porting guidance should describe PromiseGrid infrastructure boundaries, not
  the mechanics of one wire-lab run.
- Use the apparatus/specimen split: harness mechanics are evidence-gathering
  apparatus; per-protocol specs and implementations are the candidate contracts
  that may become porting targets.
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
- **App Devs:** likely first normative references are frozen protocol specs
  under `protocols/*/specs/` by pCID, followed by B-side implementation
  `CHANGELOG.md` conformance claims. Until freeze, cite draft specs only as
  provisional orientation.
- **Kernel Devs:** likely first normative references are frozen binding,
  session, and message specs by pCID, plus implementation conformance records.
  `protocols/udp-binding.d/specs/udp-binding-draft.md` and
  `protocols/group-session.d/specs/group-session-draft.md` are likely early
  ancestors, but they are not final normative citations until frozen.
- **Non-normative provenance:** TEs, DRs, TODOs, `transports/README.md`,
  `implementations/README.md`, and the harness spec explain how decisions were
  derived; they do not become app or porting APIs merely by being informative.

## Open DRs that block settled guide prose

- `DR-napum` — decides which layperson-facing PromiseGrid claims are stable
  enough for settled Intro and Goals prose.
- `DR-tuhaz` — decides the stable app-developer contract or provisional
  fallback for "How to write a grid app."
- `DR-davod` — decides the stable kernel-developer porting boundary and
  conformance target.

## Maintenance promises

- Keep this file as a source map plus writer notes. Do not move full guide
  prose here. Source: `DI-nunut`.
