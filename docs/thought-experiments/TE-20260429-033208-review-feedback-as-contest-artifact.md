# Thought Experiment - Durable review feedback as contest artifact

TE ID: TE-20260429-033208

## Decision under test

How should Steve's request-changes review of a bot-authored proposal branch be persisted in-repo when he wants the response recorded on `main`?

## Assumptions

- Bot-authored work arrives on `ppx/{twig}` branches and is reviewed locally.
- Steve wants the request-changes response for `ppx/dr-001-bootstrap` preserved on `main`, not only in chat.
- The proposal envelope schema is intentionally not locked yet, so prose artifacts are valid.
- The evolving discourse vocabulary in `harness-spec.md` §10a.3 already distinguishes endorsement, contest, and counter-proposal.
- GitHub PRs are not the desired long-term substrate for this repo's review history.

## Alternatives

1. **Chat only.** Keep request-changes feedback in chat with no durable in-repo artifact.
2. **Repo-root note.** Write a markdown file at the repo root on `main` that contains the review reply.
3. **Queue-local contest artifact.** Record the feedback as a prose `contest-v1` artifact under `proposals/pending/<proposal-id>/` on `main`.

## Scenario analysis

### Scenario 1 - Normal single-branch review

With **chat only**, Steve can reply quickly, but the repo loses the durable record and any future reviewer has to reconstruct intent from external chat logs.

With a **repo-root note**, the record is durable, but the artifact is disconnected from any proposal queue shape and the repo root accumulates process litter.

With a **queue-local contest artifact**, the review reply stays durable and sits beside the proposal it contests in the same namespace.

### Scenario 2 - Several proposal branches pending at once

With **chat only**, multiple review threads are easy to confuse and difficult for agents to correlate with exact branches or commits.

With a **repo-root note**, the files can be named carefully, but the root directory becomes an unsorted queue and humans or agents must infer proposal grouping from filenames alone.

With a **queue-local contest artifact**, each proposal has a dedicated directory and can accumulate endorsements, contests, and follow-ups without flattening everything into one directory.

### Scenario 3 - Long-horizon memory and migration away from GitHub

With **chat only**, the review decision depends on an external system and can disappear, fragment, or become hard to query later.

With a **repo-root note**, the repo preserves the record, but the artifact shape says nothing about how later tooling should distinguish a contest from other prose notes.

With a **queue-local contest artifact**, the durable record already uses the same discourse concept the harness spec expects, so later PromiseGrid-native tooling can ingest it with minimal reinterpretation.

### Scenario 4 - Agent consumption and future automation

With **chat only**, an agent has to scrape or be hand-fed external discussion.

With a **repo-root note**, an agent can read the file, but it still has to guess that the note is a contest and guess which proposal it targets.

With a **queue-local contest artifact**, an agent gets both the semantic role (contest) and the target proposal from path and file content, making later polling or indexing much simpler.

### Scenario 5 - Exact target and replayability

With **chat only**, the review feedback can drift away from the exact commit it reviewed.

With a **repo-root note**, Steve can mention the commit SHA, but the note still lives outside any proposal-specific namespace.

With a **queue-local contest artifact**, the contest can point at both the proposal branch and the exact reviewed commit SHA while remaining colocated with future follow-ups.

## Conclusions

- **Rejected:** chat only. It is too ephemeral and external for a durable review system.
- **Rejected:** repo-root note. It is durable but structurally ad hoc and scales poorly.
- **Surviving alternative:** queue-local prose `contest-v1` artifact under `proposals/pending/<proposal-id>/`.

## Implications for the repo's open TODOs and pending DIs

- Lock the decision in `DI-002-20260429-033208` under `TODO/002-review-feedback-as-contest-artifact.md`.
- Record the decision request and outcome in `DR/DR-004-20260429-033208-review-feedback-as-contest-artifact.md`.
- Publish the first durable contest artifact for `ppx/dr-001-bootstrap` on `main`.
