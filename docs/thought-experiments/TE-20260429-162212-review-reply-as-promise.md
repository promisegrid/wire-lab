# Thought Experiment - Review reply as promise

TE ID: TE-20260429-162212

## Decision under test

How should Steve's durable request-changes reply to a bot-authored proposal branch be expressed so that it fits the repo's promise-stack model, remains easy to read, and preserves a durable in-repo review record?

## Assumptions

- Bot-authored work arrives on `ppx/{twig}` branches and is reviewed locally.
- Steve wants the review reply preserved on `main`, not only in chat.
- The review queue remains under `proposals/pending/` while a proposal is neither merged nor finally rejected.
- The proposal envelope schema is intentionally not locked yet, so prose messages are acceptable.
- Steve's request-changes reply is usually conditional rather than final: the branch may become acceptable after revision.

## Alternatives

1. **`contest-v1`.** Record the reply as a peer-style contest against the proposal.
2. **Final approve/reject decision.** Record the reply as a terminal decision on the proposal as currently written.
3. **Conditional review promise.** Record the reply as a prose review message whose core sentence is: if the proposer publishes a revised branch that satisfies named conditions, Steve promises to review it again; until then, Steve does not promise to merge the current revision.
4. **Counter-proposal.** Replace the review reply with a new proposal that supersedes the current one.

## Scenario analysis

### Scenario 1 - Normal request-changes review on one pending branch

With **`contest-v1`**, the repo gets a durable record, but Steve's message is framed as if he were merely another peer contesting someone else's judgment rather than stating his own conditional commitment.

With a **final approve/reject decision**, the semantics are too strong. The current branch is not mergeable, but Steve is also not saying the work is dead forever.

With a **conditional review promise**, the durable record says exactly what Steve means: fix these things and I will review again; until then I am not promising to merge.

With a **counter-proposal**, the review reply becomes heavier than the situation requires.

### Scenario 2 - Several proposal branches pending at once

With **`contest-v1`**, each review file is durable, but later readers still need to mentally translate a contest into Steve's actual review posture.

With a **final decision**, the queue becomes artificially binary and loses the useful middle state of pending-but-needs-revision.

With a **conditional review promise**, each pending proposal can carry a clear set of conditions for the next review pass without pretending that Steve already made a terminal decision.

With a **counter-proposal**, the queue gains unnecessary extra proposals for what are really ordinary review notes.

### Scenario 3 - Long-horizon auditability of Steve's reasoning

With **`contest-v1`**, future readers see the objections, but not Steve's own promised next step.

With a **final decision**, future readers can tell acceptance state, but not the narrow path by which the proposal could have become acceptable.

With a **conditional review promise**, future readers see both the blockers and Steve's explicit commitment about what happens after those blockers are addressed.

With a **counter-proposal**, Steve's review intent is scattered across two artifacts.

### Scenario 4 - Revised branch arrives later

With **`contest-v1`**, the follow-up branch can answer the objections, but the original artifact still names the wrong speech-act.

With a **final decision**, the proposer may need an entirely new branch even when the work was merely incomplete.

With a **conditional review promise**, the proposer knows exactly what conditions unlock the next review pass.

With a **counter-proposal**, the branch history becomes harder to track because the review reply itself looks like a competing design proposal.

### Scenario 5 - Proposal remains pending while requested conditions are outstanding

With **`contest-v1`**, the queue state and the artifact semantics are slightly at odds.

With a **final decision**, the queue state is wrong because the proposal is not actually finished being considered.

With a **conditional review promise**, the proposal cleanly remains in `pending/` until the proposer satisfies the named conditions and Steve fulfills the promise to review again.

With a **counter-proposal**, the pending state of the original proposal becomes ambiguous.

## Conclusions

- **Rejected:** `contest-v1`. Durable, but semantically off for Steve's role.
- **Rejected:** final approve/reject decision. Too strong for ordinary request-changes review.
- **Rejected:** counter-proposal. Too heavy for routine review feedback.
- **Surviving alternative:** a simple prose review message whose core content is a conditional promise.

## Implications for the repo's open TODOs and pending DIs

- Lock the semantic correction in `DI-003-20260429-162212` under `protocols/wire-lab.d/TODO/TODO-20260429-162412-review-reply-as-promise.md`.
- Record the question and decided answer in `DR/DR-005-20260429-162212-review-reply-as-promise.md`.
- Publish a new promise-shaped review reply for `ppx/te-20260428-202400-promise-stack-ordering`.
- Publish a superseding promise-shaped review reply for `ppx/dr-001-bootstrap` without deleting the earlier `contest` artifact.
