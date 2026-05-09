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

## How writers should use this file

Start in the PromiseGrid Development Guide repo, then use the audience sections
below to find supporting wire-lab evidence. Each guide claim should be treated
as one of three states:

- **Ready for guide prose:** enough provenance exists to write careful
  PromiseGrid-facing prose now.
- **Write as provisional:** useful for orientation, but final wording should
  warn that design work is still in progress.
- **Blocked by DR:** do not write settled guide prose until the cited DR closes.

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
- Treat final public narrative wording as blocked by `DR-napum` if a claim must
  be presented as settled rather than provisional.

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
- Do not invent a stable app SDK, handler ABI, or application message API from
  wire-lab notes. The guide may describe the direction, but settled app-dev
  instructions are blocked by `DR-tuhaz`.
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
- Settled porting guidance is blocked by `DR-davod` until the stable kernel
  boundary, runtime expectations, and conformance obligations are explicit.

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
