# TE-nahir: Grid-envelope protocol-owned signature slot

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-nahir

## Status

decided

## Decision under test

What is the right next comparison packet for the unresolved outer-envelope
question raised by the three-slot shape `[pcid, payload, signature]` where
`pcid` owns the signature rules?

## Assumptions

- The existing `TODO-tugoz` six-sim batch remains intact and should not be
  repurposed.
- Existing grid-envelope specimens should be reused where they already cover a
  comparison family.
- The user wants to answer “what should the preferred default outer shape be,”
  not merely “can the three-slot shape work at all?”
- For this new three-slot specimen, the third slot signs canonical
  `[pcid, payload]` bytes, and `pcid` defines the proof family.

## Alternatives

- **Alt-A — narrow 3-way.** Add the three-slot sim and compare it only against
  minimal `[pcid, payload]` and explicit outer `sig_pcid`.
- **Alt-B — focused 4-way.** Add the three-slot sim and compare it against
  minimal `[pcid, payload]`, explicit outer `sig_pcid`, and payload-owned proof.
- **Alt-C — wide 5-way.** Add the three-slot sim and compare it against minimal
  `[pcid, payload]`, explicit outer `sig_pcid`, payload-owned proof, and
  wrapper proof.

## Scenario analysis

### Scenario 1 — Picking a preferred default outer shape

Alice wants evidence strong enough to recommend a default outer shape. Alt-A
leaves too much ambiguity because it omits the strongest “keep proof out of the
outer envelope entirely” counterexample. Alt-B improves that comparison but
still omits wrapper-proof designs that may win on audit clarity. Alt-C answers
the broader default-choice question directly.

### Scenario 2 — Keeping implementation cost bounded

Bob wants the comparison to stay bounded. Alt-C is still workable because only
one new sim is needed; the other four families already have usable
representatives. This weakens the main objection to the wider comparison set.

### Scenario 3 — Avoiding comparison noise

Carol worries that too many families will blur the result. Alt-A is clearest
but too narrow. Alt-B balances clarity and coverage. Alt-C adds one more
counterexample family, but it is materially different enough to justify its
place in the packet.

### Scenario 4 — Long-term audit and generic tooling pressure

Dave wants to know whether the three-slot design is better than designs that
either name proof rules explicitly or move proof semantics fully into the
payload/wrapper layer. Alt-C is the only packet that directly exposes all of
those tradeoffs together.

## Conclusions

- Alt-A is rejected because it is too narrow to answer the preferred-default
  outer-shape question.
- Alt-B survives as a reasonable fallback if implementation cost grows.
- Alt-C survives and is preferred because it answers the real design question
  without requiring a batch of new sims; only one new specimen is missing.

## Recommendation

Adopt Alt-C. Add one new standalone simulation for the three-slot
protocol-owned-signature shape and compare it against existing representatives
for minimal `[pcid, payload]`, explicit outer `sig_pcid`, payload-owned proof,
and wrapper proof. Reuse the six-scenario slice already used for focused
envelope comparison work.

## Implications for open work

- Create a new TODO/DR thread separate from `TODO-tugoz` because this question
  is about whether the outer envelope needs a separate proof selector at all.
- Add one new standalone simulation and register it as an active alternative,
  not consensus.

## Decision status

`decided` — locked by `DI-kukuk` in
`protocols/wire-lab.d/TODO/TODO-mujad-grid-envelope-protocol-owned-signature-slot.md`.

