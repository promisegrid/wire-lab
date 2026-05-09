# TE-david: PromiseGrid dev-guide resources from wire-lab evidence

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-david

## Status

decided

## Decision under test

DUT-david: Decide what the wire-lab repo must provide so humans and LLMs can write the PromiseGrid Development Guide without confusing the guide's subject, authority model, or intended audience.

The PromiseGrid Development Guide is about PromiseGrid. Wire-lab is not the product guide and is not the stable product source of truth. Wire-lab is the experimental simulation space where PromiseGrid design choices are explored, tested, derived, and recorded before they stabilize.

This TE tests four coupled decisions:

- where guide writers should start inside wire-lab;
- how much synthesis wire-lab should provide versus forcing writers to read the whole corpus;
- how to prevent writers from treating experimental wire-lab mechanics as PromiseGrid commitments;
- how authority should shift after the PromiseGrid guide stabilizes.

## Assumptions

- The development guide's current audience structure is already locked in the `promisegrid-dev-guide` repo: Laypeople, App Devs, and Kernel Devs.
- Wire-lab contains the current design provenance for many PromiseGrid choices, but much of that provenance is experimental, provisional, or explicitly superseded.
- Human writers and LLM writers need different amounts of orientation, but both need the same authority rules.
- A stabilized PromiseGrid Development Guide should become a higher-level source of truth than exploratory wire-lab notes.
- Hashed, frozen protocol spec docs remain authoritative by pCID wherever they are hosted. A frozen spec may live in wire-lab, another PromiseGrid repo, or a future archival location; the pCID, not the repo path, carries protocol-spec identity.
- The first implementation task should not edit the guide repo. It should create wire-lab-side planning artifacts that make the later guide-writing work safer.

## Alternatives

### Alt A — Top-level resource map plus writer notes

Create a top-level `DEV-GUIDE-RESOURCES.md` in wire-lab later, tracked by a harness TODO. The file maps guide audiences to wire-lab evidence and adds short PromiseGrid-facing writer notes. It does not draft the guide itself.

This alternative gives writers a clear starting point and enough synthesis to avoid obvious misreadings. Its obligation is maintenance: when TEs, DIs, DRs, protocol specs, or guide prose stabilize, the resource map must be updated.

### Alt B — README-only orientation

Put the guide-writer orientation directly in wire-lab's top-level `README.md`.

This keeps the path obvious but overloads the repo README. The README already has to explain wire-lab as a simulation harness; embedding the guide-writing source map there risks making wire-lab look like the guide's subject.

### Alt C — Edit the guide repo directly now

Use wire-lab evidence to start writing the PromiseGrid Development Guide bodies immediately.

This moves fast but skips the core safety problem: the writer still needs a stable map from experimental evidence to PromiseGrid-facing claims. It also creates cross-repo mutation before wire-lab has recorded what it is promising to provide.

### Alt D — Source map only

Create `DEV-GUIDE-RESOURCES.md`, but include only source links and gap labels.

This has the lowest drift risk, but it leaves too much synthesis to each writer. LLMs in particular are likely to overgeneralize from a single TE or mistake harness mechanics for PromiseGrid commitments.

### Alt E — Source map plus short excerpts

Create a source map with writer notes and selected quotes from wire-lab sources.

This is immediately useful but creates the highest maintenance burden. Quoted text can become stale or historically contextual under the TE editing policy, and future guide writers may cargo-cult old wording instead of following the current authority chain.

## Scenario analysis

### Alice — human writer drafting layperson intro and goals

Alice starts in the PromiseGrid Development Guide repo. She sees the locked Laypeople sections but does not know which wire-lab documents explain the current PromiseGrid story. She needs motivation, goals, and stable vocabulary, not the mechanics of a simulation harness.

Under Alt B, Alice reads the wire-lab README and may conclude the guide should explain wire-lab. That is wrong: the guide should explain PromiseGrid, while wire-lab supplies provenance for design choices.

Under Alt D, Alice gets links but little help deciding which claims are stable enough for lay readers. She may quote an old experiment or bury the guide in internal process.

Under Alt E, Alice gets quotable phrases but may copy historical or provisional wording into the guide.

Alt A gives Alice the right path: a PromiseGrid-facing resource section that says which sources support layperson motivation, which claims remain unsettled, and which wire-lab details should stay out of the guide unless they explain why a PromiseGrid choice exists.

### Bob — LLM writer drafting app-developer guidance

Bob is asked to write "How to write a grid app." He can read the repo quickly, but he is vulnerable to source blending. If he sees pCIDs, transports, apparatus/specimen language, and experimental TODOs without orientation, he may hallucinate a stable app API or present wire-lab protocol specimens as required PromiseGrid behavior.

Alt B increases that risk because README prominence does not equal design stability.

Alt D reduces navigation cost but still leaves Bob to infer which sources are provenance, which are normative, and which are obsolete.

Alt E gives Bob more text to copy, which may improve fluency while worsening authority mistakes.

Alt A is strongest for Bob because writer notes can say what each source does and does not authorize. For example, a note can say that pCID identity matters for frozen protocol specs, but that a draft transport experiment is not a PromiseGrid app API. That prevents the most likely LLM failure: turning lab evidence into product commitments.

### Carol — kernel developer drafting porting guidance

Carol needs to explain how to port PromiseGrid infrastructure to a new platform or language. She must understand the apparatus/specimen split, the layered transport/feed model, protocol specs, pCID identity, and DR/DI provenance. She also needs to know when wire-lab stops being the authority.

Alt B mixes Carol's deep porting needs into a general README and makes maintenance awkward.

Alt D makes Carol read the whole corpus and reconstruct the authority chain herself.

Alt E may tempt Carol to quote old TE passages whose main value is historical reasoning, not current porting instruction.

Alt A lets wire-lab provide a precise bridge: the future resource map can point Carol to the apparatus/specimen TE, transport-layering TE, protocol specs, and open DRs while stating that stabilized guide prose supersedes exploratory notes. Frozen specs remain exceptions: when the guide points to a hashed protocol spec by pCID, that frozen spec is the protocol authority whether it is hosted in wire-lab or elsewhere.

## Conclusions

- Reject Alt B. A README pointer is useful, but the full resource map belongs in a dedicated top-level file so the wire-lab README can remain about wire-lab.
- Reject Alt C for this task. Direct guide prose should wait until wire-lab records the writer-resource promise and source map obligations.
- Reject Alt D as too thin for mixed human/LLM authorship.
- Reject Alt E for the first pass because excerpts create unnecessary drift and quotation-maintenance risk.
- Choose Alt A: create a future `DEV-GUIDE-RESOURCES.md` as a source map plus short PromiseGrid-facing writer notes.

## DF questions and locked answers

### DF-david.1 — What is the guide about?

Locked answer: the PromiseGrid Development Guide is about PromiseGrid, not wire-lab. Wire-lab is design provenance and experimental simulation evidence.

### DF-david.2 — Where should guide writers start in wire-lab?

Locked answer: create a future top-level `DEV-GUIDE-RESOURCES.md`, with a concise pointer from top-level `README.md`.

### DF-david.3 — What should the resource file contain?

Locked answer: source map plus short writer notes. It should not contain full guide prose and should not copy large excerpts from source docs.

### DF-david.4 — What is the authority model after guide stabilization?

Locked answer: once the PromiseGrid Development Guide stabilizes, stable guide prose becomes the higher-level source of truth for developers. Exploratory wire-lab notes remain provenance. Hashed, frozen protocol specs are the exception: when the guide cites a frozen spec by pCID, that spec is authoritative for its protocol regardless of hosting repo.

## Implications for open work

- File a harness TODO that tracks the future `DEV-GUIDE-RESOURCES.md` work and keeps it visibly PromiseGrid-facing.
- The resource map should be organized by the guide's locked audiences: Laypeople, App Devs, and Kernel Devs.
- Each audience section should include current wire-lab sources, writer notes, open gaps or DRs, and DI/DR provenance.
- The wire-lab top-level README should eventually link to `DEV-GUIDE-RESOURCES.md`, but should not absorb the whole resource map.
- The resource map must include an authority-transition note so future writers know when to prefer the stabilized guide over exploratory wire-lab material.

## Decision status

`decided` — Steve selected the dedicated `DEV-GUIDE-RESOURCES.md` entrypoint, source map plus writer notes, one harness TODO, and the authority-transition framing in planning discussion on 2026-05-08. The locked implementation decision is recorded in `DI-nunut` in `protocols/wire-lab.d/TODO/TODO-rozas-promisegrid-dev-guide-resources.md`.
