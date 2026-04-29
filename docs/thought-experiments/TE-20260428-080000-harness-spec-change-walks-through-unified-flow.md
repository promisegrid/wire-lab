# TE-14: A harness-spec change walks through the unified flow

*Thought experiment, part of the [PromiseGrid Wire Lab](../../harness-spec.md). This file is content-addressable; its hash is its pCID.*

**Purpose.** Stress-test whether the single-flow design (§10a) holds up when the proposal targets the harness spec itself. Walked through end-to-end on 2026-04-28; full trace lives in `thought-experiments/TE-14-harness-spec-walkthrough.md` once split out, summary here.

**Setup.** Harness-spec at `harness-spec-v7`. Cast: Steve (`design-judgment` ~0.95), Elder-7 (~0.4), Claude-A (LLM, ~0.08), Hostile-Prober-3 (~0.01). Trigger: empty `transcript.cbor` files accumulating for crashed-setup runs.

**Walk.**

1. Elder-7 emits a prose proposal in a wire message: outer assertion "I propose a change to the spec at pCID <harness-spec-v7>," inner body in English, signed.
2. Claude-A independently drafts a more thorough proposal four hours later, signs with its low-trust key.
3. Elder-12 endorses Claude-A's. Hostile-Prober-3 contests Elder-7's. Claude-B counter-proposes a cleaner alternative (delete the run dir entirely) referencing Claude-A's by ID.
4. Steve reads Claude-A's queue summary (the LLM disclosed its conflict-of-interest as a proposer summarizing the queue), reads all three proposals, decides Claude-B's counter-propose is best, addresses Prober-3's contest by adding one line to the change.
5. Steve signs a `merge-harness-spec` promise pointing at `harness-spec-v8` (Claude-B's change plus Steve's amendment). The harness-spec pointer now follows that signature.
6. One week later, settlement: empty-transcript count is 0 as predicted. Trust ledger awards prediction-quality credit to Claude-B (+0.04), Claude-A (+0.02), Elder-7 (+0.02). Prober-3's contest, though resolved-against, was recorded as having surfaced a real concern that influenced the merged outcome.

**Findings: where the design holds.**

- One vocabulary handled everything. No `harness-spec-change-v1` envelope was needed.
- No envelope schema was needed. Prose proposals, mid-structure proposals, counter-proposes, and merge promises coexisted.
- The "your signature is the only lock" rule worked as intended: low-trust agents could propose, contest, and endorse, but only Steve's signature moved the pointer.
- Trust-ledger weighting handled noise gracefully without special-casing.

**Findings: where the design wobbled (load-bearing details TE-14 surfaced).**

- **Proposal-collector handler matching.** What does the handler that files messages into the review queue actually pattern-match on, given there's no schema? In the prose-friendly early phase, the most honest answer is: an LLM-as-router that reads incoming messages and decides whether they look like proposals. Which is fine — it's just another low-trust agent making classification promises — but it means the flow depends on at least one LLM in the routing path from day one. Worth being explicit about.
- **Pointer storage semantics.** "The harness-spec pointer follows your signing key" is mechanically true at the cryptographic level, but where does the pointer *live*? Workspace file pushed to git? CRDT? DNS-like lookup against your public key? TE-14 implicitly assumed a workspace file. Not yet decided.
- **Settlement window.** TE-14 picked 7 days for prediction-quality settlement out of nowhere. Real proposals will have wildly varying settlement windows (typo fix: hours; trust-ledger-merge rule change: months). Proposers should probably declare their own predicted-falsification window as part of the proposal-checklist convention.
- **Conflict of interest when an LLM both proposes and summarizes the queue.** Claude-A in TE-14 disclosed the conflict in its summary. The convention should require this; the trust ledger will catch sustained self-favoring summaries by degrading the LLM's score, but a community-norm spec saying "disclose if you're summarizing a queue that contains your own proposals" is cheap and helpful.
- **No explicit `defer` verb.** Steve decided in 9.5 hours, but a real harness-spec change might warrant "give me two more weeks of crash data first." The absence of `defer` is fine — a proposal sitting in `pending/` *is* deferred — but the doc should say so.

**Big-picture finding.** The single-flow design holds up. None of the wobbles required inventing a parallel mechanism for harness-spec changes; they're all refinements within the one-flow framework.

**Open questions promoted from TE-14**: see §12 #12 (proposal-collector matching), #13 (pointer storage), #14 (settlement-window convention).
