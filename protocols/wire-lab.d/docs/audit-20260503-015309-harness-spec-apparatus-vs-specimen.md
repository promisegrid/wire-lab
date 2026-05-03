# Audit: harness-spec-draft.md — apparatus vs. specimen

**Date:** 2026-05-03 01:53:09 UTC
**Auditor:** stevegt-via-perplexity (bot)
**Target:** `protocols/wire-lab.d/specs/harness-spec-draft.md` (652 lines)
**Status:** Memo only. No edits to the spec yet.
**Companion target (not audited here):** `protocols/wire-lab.d/specs/transport-spec-draft.md`.

## Framing

Per Steve's correction this session, the wire-lab harness is **lab apparatus**. Its spec defines how scenarios run, how actors are named, how transports are simulated, how messages get injected and observed, how trust ledgers are scored, and how outputs are compared. The harness-spec must NOT prescribe wire shape, envelope shape, transport binding, ordering rules, or criticality models. Each of those is a **specimen** — a hypothesis under study — and each lives in its own `protocols/<slug>.d/` directory with its own spec doc and TEs.

`grid([pcid, payload])` is one envelope hypothesis (working but unproven). Promise-stack (`[]Promise` with the `Wrap`/`Peel`/`Project` library API) is another envelope hypothesis. Both belong in their own protocol directories, not in the harness-spec.

**Three-way classification used below:**

- **Apparatus** — describes the lab. Belongs in harness-spec.
- **Specimen** — prescribes a particular wire/envelope/binding/ordering/criticality hypothesis. Must move to a per-protocol spec doc.
- **Ambiguous** — frames a question or taxonomy that is partly apparatus (the harness has to study it) and partly specimen (the spec text leans toward one answer). Needs a TE to disambiguate before editing.

## Section-by-section classification

| Section | Lines | Class | Brief rationale |
|---|---|---|---|
| Preamble + provisional note | 1–10 | Apparatus | Frames the harness as discovery vehicle; explicitly flags concrete shapes as provisional. Aligns with apparatus framing. |
| The 100-Year Goal (C-1..C-6) | 13–28 | Apparatus | Constraint set is meta-level; pressure-tests both apparatus and specimens. Stays. |
| §1 Core Reframe header + pCID note | 32–36 | Ambiguous | The pCID-as-spec-hash note is meta (apparatus); but "Messages Are Layered Promises" framing already prescribes one envelope. Header itself is a specimen claim. |
| §1.1 Promise-stack model (`Promise` struct, `[]Promise` shape, CBOR array, top-down accept/defer/reject, `promstack` Wrap/Peel/Project library API) | 38–70 | **Specimen** | Prescribes one specific envelope hypothesis (nested promise stack), one specific receiver-side library API (`promstack` with three operations), and one specific encoding family (CBOR array of promise frames). This is exactly the kind of prescription the harness-spec must NOT carry. Belongs in a `protocols/promise-stack.d/specs/` doc (directory does not yet exist). |
| §1.2 Why this is better than a fixed envelope | 72–77 | **Specimen** | Argues for the promise-stack hypothesis vs. an alternative. Argument-for-a-specimen, not apparatus. Belongs alongside §1.1 in the promise-stack protocol's spec/TE corpus. |
| §1.3 What the simulator tests about layering | 79–84 | Ambiguous | "Receiver can handle out-of-order promise stacks" / "forwarding node can strip its own promise" / "promise about a missing inner body" — these read as apparatus tests, but they presume the promise-stack envelope. If we keep them, they have to be reframed at the apparatus-level: "the harness exercises any envelope under out-of-order, forwarding-strip, and merkle-reference scenarios, and reports on each candidate envelope's behavior." Otherwise they migrate with §1.1. |
| §2 Trust as Durable Relationship (header + introduction) | 88–90 | Apparatus | Framing — trust as first-class subject of experimentation. |
| §2.1 Per-agent trust ledger (`TrustLedger` struct shape) | 92–110 | **Specimen** (struct shape) / Apparatus (existence of a trust ledger) | The existence of per-agent trust ledgers and the need to record interactions/keep/break events is apparatus. The specific field set (`first_seen_ns`, `interactions`, `kept`, `broken`, `evidence_chain`, `open_promises`, `score`, `score_components`, `reputation_imports`, `relationship_age_ns`, `last_drift_ns`) is one shape proposal. Should be lifted out into a candidate-trust-ledger-shape protocol or TE; the harness-spec should say "the harness records, per peer pair, evidence of kept/broken/open promises sufficient to score under any of the families in §2.1's table." |
| §2.1 Three scoring families (Beta-Bayesian / Worldline-weighted / EigenTrust-like) | 112–120 | Apparatus | Comparative taxonomy of scoring rules — the harness studies all three. Stays. |
| §2.2 "Broken promise" event is itself a promise | 122–126 | **Specimen** | "B emits a witnessed-break promise" presumes the promise-stack envelope. Apparatus-level statement would be "the harness studies how break-witness messages get carried by candidate envelopes and how recipients weigh them." Belongs in promise-stack protocol spec (or wherever break-witness lives). |
| §2.3 Durable relationship features | 128–134 | Apparatus | First-encounter rituals, reputation portability, decay, defection cost, multi-generational handoff — all knobs/scenarios the harness must support. Stays (with light wording fixes for "promise" usage). |
| §2.4 Trust is per-assertion-type | 136–138 | Apparatus | Asserts a constraint (trust must be vector, not scalar) that any candidate ledger must honor. Aligned with C-5. Stays. |
| §3 header | 142–144 | Apparatus | Frames CDA as one of several economic-coordination experiments. Stays. |
| §3.1 The model (token issuance, portfolio, CDA) | 146–152 | Apparatus | Describes the CDA experimental setup the harness will run. Stays. |
| §3.2 Exchange rates as trust signals | 154–158 | Apparatus | Hypothesis the harness will test (rate paths as trust signal). Stays. |
| §3.3 What the simulator must verify | 160–166 | **Specimen-leaning** | Bullet 1 explicitly says "the wire-level encoding of orders, fills, redemptions composes from the same `Promise` primitive (no special market protocol envelope)." That presumes the promise-stack envelope. Other bullets are apparatus. Reframe bullet 1 as "the harness studies whether the candidate envelope can carry orders/fills/redemptions without a separate market envelope." |
| §3.4 Alternative economic models (gift, mutual-credit, stake-bonded, quadratic, kernel-role auctions) | 168–178 | Apparatus | Comparative taxonomy. Closing line ("the wire format and trust ledger must support all of them") is apparatus-level requirement on whichever envelope wins, not a prescription of one. Stays. |
| §4 Ingress models K1–K5 | 182–203 | Apparatus | Taxonomy of ingress models the harness will study. Stays. |
| §4.3 Hybrid scenarios | 205–211 | Apparatus | Acceptance criterion ("yes, identically at the wire level") that the harness will enforce on any candidate envelope. Apparatus-level. Stays. |
| §5 Agents (deterministic + stochastic; profiles; scenario generators) | 215–257 | Apparatus | Agent typology and stochastic-agent profiles are scenario infrastructure. Stays. |
| §5.4 "Stochastic agents talk the same wire as deterministic ones" | 255–257 | Apparatus | An apparatus-level invariant. Stays. |
| §6 Tragedy of the Commons (C1–C7 scenarios + metrics + pass line) | 261–294 | Apparatus | Acceptance test scenarios for the whole system. Stays. (The 30%-at-month-12 pass line is one knob but already flagged as provisional in §12 OQ-11.) |
| §7.1 Edges carry promises, not bytes (transport-promise as outermost frame) | 302–304 | **Specimen** | "This promise becomes the outermost frame on the message handed to Layer C" presumes the promise-stack envelope and prescribes a specific frame placement. Apparatus-level statement would be "every edge annotates inbound traffic with observed metadata (source, time, integrity class, liveness) and hands that metadata to Layer C in whatever shape the candidate envelope demands." Belongs in promise-stack protocol spec; the apparatus version belongs in harness-spec. |
| §7.2 pCID registry is itself a commons | 306–312 | Apparatus | Acceptance scenarios for any naming scheme. Stays. |
| §7.3 Time is multi-clock | 314–322 | Apparatus | Per-agent clock infrastructure. Stays. |
| §7.4 Snapshot/fork/replay/diff | 324–331 | Apparatus | Harness facilities. Stays. |
| §7.5 Real-world cross-pollination | 333–342 | Apparatus | Realism upgrades to the harness. Stays. |
| §7.6 Adversarial-by-default review | 344–346 | Apparatus | Shadow-adversary acceptance gate. Stays. |
| §8 Thought Experiments index | 350–387 | Apparatus | Cross-references to docs/thought-experiments/. Stays (entries themselves may need link updates as protocols carve out, but the index is apparatus-level). |
| §9 Realism Suggestions | 391–421 | Apparatus | Cross-cutting realism advice. Stays. |
| §10 grid-poc directory table | 425–440 | Ambiguous | Names `x/wire-v2/promstack/` as the home of the promise-stack library and `x/wire-v2/trust/` as the home of trust primitives. After the protocol carve-out, those become **implementations** of specimens, not harness components. Apparatus-level reframe: "harness fixtures live under `x/sims/...`; specimen implementations live alongside `implementations/<impl-name>/` per TE-32." Stays in some form, but rows for `promstack` and `trust` migrate. |
| §10a.1 Specs may be prose, structured, or both | 454–463 | Apparatus | Methodology. Stays. |
| §10a.2 Proposals are messages | 465–476 | **Specimen-leaning** | "A proposal is a promise stack like any other message" presumes promise-stack envelope. Reframe at apparatus-level: "a proposal is a message in whatever envelope the host transport carries; the proposal-checklist convention is content-level, not envelope-level." Move envelope-specific text out. |
| §10a.3 Endorsement/contest/counter-propose | 478–488 | Specimen-leaning | "endorse-v1 / contest-v1 / counter-propose-v1" pCIDs; these are themselves protocols and belong in their own `protocols/...d/` directories (or in `protocols/ppx-dr.d/`). The discourse-vocabulary as a *requirement on the apparatus* (the harness must support proposal/endorse/contest workflows) stays; specific pCIDs and named handlers move out. |
| §10a.4 Cross-run aggregation + design-knob registry (`runs/index.parquet`, `knobs.cbor`) | 490–508 | Apparatus | Harness-side data infrastructure. Stays. |
| §10a.5 Outcome attribution (sweeps, counterfactual replay) | 510–517 | Apparatus | Harness study tools. Stays. |
| §10a.6 Counterfactual evidence for in-sim agents (`hypothesis` and `hypothesis-result-v1` assertions) | 519–525 | Specimen-leaning | The hypothesis-emitter mechanism is apparatus, but the named assertion types presume one envelope vocabulary. Reframe at apparatus-level. |
| §10a.7 Durable cross-run agent memory | 527–538 | Apparatus | Per-agent-id `personal-archive/` directory. Stays. |
| §10a.8 What's load-bearing: signing key | 540–550 | Apparatus | Constitutional rule for the harness-spec pointer itself. Already flagged in §10 as the only `[Production-bound]` claim. Stays. |
| §10a.9 Convergence indicator | 552–561 | Apparatus | Harness logging discipline. Stays. |
| §10a.10 Review queue + your role as agent | 563–576 | Apparatus | Workflow. Stays. |
| §10a.11 Canonical iteration | 578–590 | Apparatus | Loop description. Stays. |
| §10a.12 What this enables and does not | 592–603 | Apparatus | Closing scope statement. Stays. |
| §11 The Decision the Harness Is Trying to Make | 607–613 | Apparatus | Single-sentence framing of what the simulator is for. Stays. |
| §12 Open Questions (33 entries) | 617–651 | Mostly apparatus, mixed | Each OQ classifies independently; most are apparatus questions about the harness. OQ-1 (encoding), OQ-19 (group-transport graduation), OQ-2 (currency vs mutual-credit), OQ-6 (one canonical break-witness pCID), OQ-7 (vector vs scalar projection) are specimen-side questions about candidate protocols, but their placement in §12 is acceptable as "harness-spec asks the human these questions about the specimens it is studying." Stays as-is at this audit; some questions may relocate to per-protocol TE/DR queues during step 4 of the 6-step plan. |

## Concrete specimen-bearing material that must move out

The clearest specimen-level prescriptions, in priority order for migration:

1. **§1.1 Promise struct + `[]Promise` message shape + `promstack` Wrap/Peel/Project library API.** Belongs in a (not-yet-existing) `protocols/promise-stack.d/specs/promise-stack-draft.md`. TE-1 (promise-stack ordering) already exists and would relocate to that protocol's TE corpus.

2. **§1.2 Why this is better than a fixed envelope.** Argument for the promise-stack hypothesis. Belongs in the same promise-stack spec doc or TE.

3. **§1.3 simulator tests about layering.** Either (a) reframe at apparatus-level ("harness exercises any candidate envelope under out-of-order / forwarding-strip / merkle-reference scenarios") or (b) migrate with §1.1.

4. **§2.1 `TrustLedger` struct shape.** The struct as one candidate trust-ledger shape. May warrant its own protocol slug if multiple candidate shapes will be studied; otherwise lifts to apparatus-level invariants ("any trust ledger must record kept/broken/open evidence per peer pair, support per-assertion-type score components, support reputation imports as down-weighted").

5. **§2.2 "broken promise event is itself a promise" mechanism.** Belongs with whichever envelope owns the break-witness vocabulary (likely promise-stack).

6. **§3.3 bullet 1** ("composes from the same `Promise` primitive"). Reframe at apparatus-level.

7. **§7.1 transport-promise as outermost frame.** Reframe at apparatus-level (edge annotates with observed metadata; placement in the candidate envelope is the specimen's concern).

8. **§10 grid-poc directory rows for `promstack` and `trust`.** Migrate to per-implementation table once `implementations/<impl-name>/` directories exist (per TE-32).

9. **§10a.2/§10a.3/§10a.6 named pCIDs and assertion types** (`endorse-v1`, `contest-v1`, `counter-propose-v1`, `hypothesis-result-v1`, etc.). Move to whichever protocol owns the discourse vocabulary; likely `protocols/ppx-dr.d/`.

## Material that is unambiguously apparatus and stays

100-year goal, 100-year constraints, scoring family taxonomy, ingress model taxonomy K1–K5, agent profile taxonomy, ToC scenarios C1–C7, Ostrom report, snapshot/fork/replay, real-world cross-pollination, adversarial-by-default review, TE index, realism suggestions (most), §10a.4 cross-run aggregation, §10a.5 attribution, §10a.7 personal archive, §10a.8 signing-key lock, §10a.9 convergence indicator, §10a.10 review queue workflow, §10a.11 canonical iteration, §11 framing.

## Sections needing TE before edits (Ambiguous)

- §1.3 "what the simulator tests about layering" — reframe at apparatus-level vs. migrate.
- §10 grid-poc directory table — which rows are apparatus fixtures and which are specimen implementations.
- §10a.2 / §10a.3 / §10a.6 — apparatus-level vocabulary for proposals/endorse/contest/hypothesis without naming a specific envelope.

## Recommended next move

Per the 6-step corrected plan, step 2 is to file a harness-level TE on the apparatus-vs-specimen split itself (alternatives: strict carve-out, mixed/hybrid, registry-of-knobs). That TE should reference this audit by tree-hash once committed and use the section-by-section table above as input. The TE's tabletop scenarios should walk Alice, Bob, Carol, Dave, Ellen, Frank, Mallory through (1) reading §1.1 today, (2) reading §1.1 after carve-out, (3) attempting to add a second envelope hypothesis that competes with the promise-stack, (4) attempting to fork the promise-stack hypothesis, (5) Mallory injecting a malicious envelope claim, (6) a 30-years-later contributor finding the harness-spec but not the carved-out protocol specs.

After step 2 locks the carve-out shape, step 3 creates the protocol directory (`protocols/promise-stack.d/`), step 4 sweeps harness-spec under Cat-1a/Cat-2 of the editing policy citing the new DI, step 5 reframes TODO 5 under the promise-stack protocol's TODO/, step 6 files a parallel TODO for the `grid([pcid, payload])` envelope hypothesis.

## Companion file pending audit

`protocols/wire-lab.d/specs/transport-spec-draft.md` (≈9 KB) likely contains analogous specimen-vs-apparatus crossings (transport-spec is meant to be apparatus-level per TE-26/TE-27/TE-29, but probably names specific binding shapes inline). Not audited in this memo; flag for follow-on audit before step 4.

## What is NOT recommended

This memo is read-only. No edits to harness-spec, transport-spec, or any TE/DI/DR were made or are recommended on this pass. The 6-step plan calls for a TE first, then a DI, then sweeps citing the locked DI. The audit is the input to step 2, not a license to skip it.
