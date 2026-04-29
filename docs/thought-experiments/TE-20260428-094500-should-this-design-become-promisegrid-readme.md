# TE-15: Should this design become `promisegrid/promisegrid/README.md`?

*Thought experiment, part of the [PromiseGrid Wire Lab](../../harness-spec.md). This file is content-addressable; its hash is its pCID.*

**Purpose.** Decide whether the Wire Lab harness specification should eventually replace, augment, or be linked from the canonical PromiseGrid org-level README at [github.com/promisegrid/promisegrid](https://github.com/promisegrid/promisegrid). Run on 2026-04-28 with the actual current README in front of us.

## What the existing README is, today

- Last meaningful push: February 2026.
- Public, GPL-3.0, fairly readable, ~365 lines.
- Frames PromiseGrid as a consensus-based decentralized computing system, designed to address tragedy-of-the-commons, owned by its users rather than any single entity.
- Established core ideas already in the README: capability-as-promise, content-addressable code, merge-as-consensus, nested kernels, Multihash addressing, Promise Theory grounding, ToC as the root motivation.

## Where the existing README and the Wire Lab design line up

- **Capability-as-promise.** Already there ("a capability token represents a promise that the issuer will either break or fulfill at a later time"). The Wire Lab's layered-promise framing is a deepening of this, not a replacement.
- **ToC as the central problem.** Already there. The Wire Lab's C1–C7 scenarios and Ostrom audit operationalize it.
- **Content-addressable everything.** Already there. The Wire Lab's pCID = "hash of the spec document" extends it.
- **Merge-as-consensus.** Already there. The Wire Lab's proposal/endorse/contest/counter-propose flow is a more concrete instance.
- **Nested kernels.** Already there. The Wire Lab's K1–K5 ingress models unpack the spectrum.

## Where the Wire Lab design has moved past the existing README

These are the places the README is no longer accurate, not just incomplete.

1. **§"PromiseGrid Universal Protocol"** says "A message consists of a capability token followed by a payload" and conflates the capability token with the protocol-version identifier. The Wire Lab separates these: a capability token is a promise of action; a pCID is a separate hash naming the spec the payload conforms to. The README's framing is now wrong.
2. **§"Capability-as-Promise"** treats trust implicitly. The Wire Lab introduces a per-peer durable trust ledger with per-assertion-type scoring, relationship age, decay, and break-witnesses. The README is silent on this.
3. **§"Architecture"** describes the system as primarily a decentralized kernel with WASI as the portability target. The Wire Lab's K1–K5 explicitly de-thrones the kernel and treats WASM/WASI as one runtime among several (browser, native, MCU, container). The README is stale here.
4. **§Milestones** has framings ("Write WASI target") that no longer match the multi-runtime stance.
5. **The README has no methodology section.** It presents a finished-sounding system. The Wire Lab is fundamentally about how the design itself evolves: prose specs graduating to structured ones, proposals reviewed via the unified flow, the user's signing key as the single load-bearing lock. This is the most important missing piece and the one most worth eventually integrating.

## Where the Wire Lab harness should live

A separate question from "should the README change," but the answer feeds back. Two options:

**Option A — Wire Lab inside `promisegrid/promisegrid` under `wire-lab/`.** One repo. Atomic updates. But changes the existing repo's nature into a monorepo; every Wire Lab commit churns the canonical project's pushed-at timestamp.

**Option B — Wire Lab in its own `promisegrid/wire-lab` repo.** Matches the precedent already set by `promisegrid/grid-poc` (also its own repo). README cadence and Wire Lab cadence stay separate. Two scoped PATs are easier to reason about than one.

**Recommendation: Option B**, because the precedent already exists, the cadences are different, and limiting bot write access to a single experimental repo is the safer default.

## When and how the canonical README should change

A four-phase plan:

**Phase 1 (now):** Create `promisegrid/wire-lab` as its own repo. Push the harness-spec, thought experiments, and a small repo-level README that orients readers to *the harness*. The org-level README at `promisegrid/promisegrid` does **not** change yet.

**Phase 2 (after one or two real harness runs have produced findings):** Minor org-README update. Add two paragraphs: "The protocol details are being prototyped in [`promisegrid/wire-lab`](https://github.com/promisegrid/wire-lab). Findings to date are updating our understanding of …" Additive, low-risk, doesn't touch the architecture sections.

**Phase 3 (months out, once the Wire Lab framing has stabilized):** Substantive org-README rewrite. Replace §"PromiseGrid Universal Protocol" with a layered-promises framing. Add §"Trust as Durable Relationship." De-canonicalize the kernel framing in §"Architecture." Add a methodology section. Update milestones.

**Phase 4 (eventually):** The org-level README becomes something close to what the Wire Lab harness-spec is now, plus quick-start material, plus the narrative of how the design got there. The Wire Lab repo continues to exist for the running harness, but the README is the public face.

## The reflexive validation milestone

The first non-trivial use of the unified proposal flow on a document that *truly matters to Steve* should be the Phase 3 rewrite of the canonical README. That is: the methodology gets validated by being applied to update the README that describes the methodology. If the unified flow can produce a README rewrite that Steve is willing to sign and merge, the methodology has earned its standing. If it cannot, the methodology is not yet ready and the harness has done its job by surfacing that.

## Findings

- The existing README is older but not wrong-in-spirit. It is wrong-in-detail in the places enumerated above, and silent on several big topics (trust ledger, methodology, multi-runtime stance).
- Two repos, two PATs, two cadences. Wire Lab gets its own repo with its own write access.
- The canonical PromiseGrid README stays under Steve's hands directly until the harness produces findings worth integrating.
- The eventual Phase 3 README rewrite is the natural validation milestone for the unified flow.

## Open questions promoted from TE-15

- See harness-spec §12 #17 (when does Phase 2 trigger?).
- See harness-spec §12 #18 (does the canonical README get its own pCID, signed by Steve, distinct from the Wire Lab harness-spec pCID?).
