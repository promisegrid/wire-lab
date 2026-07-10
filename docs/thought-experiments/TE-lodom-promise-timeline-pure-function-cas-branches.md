# Promise timeline, pure-function agents, CAS branches, and token double-spend

TE ID: TE-lodom

## Status

decided, refined

Cat-2 / decision-refinement update per `TODO-nudav` / `DI-mokaz`: "local
ledger" wording is replaced by "derived local projection/index/cache" wording
where it describes the current recommendation. No DIs live in this file.

## Decision under test

Should PromiseGrid add a future POC that treats promises as timeline assertions
stored in decentralized CAS object chains, treats agents as deterministic
pure-function servers over explicit context, and treats capability-token
double-spend as branch evidence that receivers may merge, reject, or remember
locally?

## Assumptions

- PromiseGrid remains Promise Theory first: no agent promises on behalf of
  another agent, all cooperation is voluntary, and trust remains local and
  relationship-relative.
- A promise can be modeled as an assertion by a promiser about some part of the
  universe at a point or interval on a timeline. This includes assertions about
  bytes, computations, device readings, social commitments, economic terms,
  resource retention, and observed peer behavior.
- CAS object identity names exact bytes. A CID does not by itself promise
  availability, correctness, or trustworthiness; those are separate promises.
- An agent's local CAS is its chronological event source and source of truth for
  POC20. Current state is a derived projection from local CAS event objects, not
  a separate authority.
- Local CAS is not automatically public. Some local CAS objects may remain
  private forever, some may be sent only after encryption, and some may be
  plain-shareable. Sending or withholding an object is a local promise decision.
- POC18 already models CAS/VCS parent-linked object chains, reference sets,
  tags, branches, snapshots, logical changes, and continuous peer sync. POC20
  should reuse those lessons rather than inventing a second graph model.
- POC19 remains the production-shaped plumbing path. This TE explores semantic
  model work for a parallel POC20 and should not block POC19 unless a later
  Decision Intent explicitly changes that.

## Alternatives

### Alternative A: Keep token replay ledgers as mutable local side state

Agents continue to reject replayed or already-redeemed tokens with local mutable
spent-token ledgers. Branches may exist for VCS data, but token redemption
history is not itself modeled as branchable CAS timeline state.

### Alternative B: Model all promise history as append-only CAS timelines

Every relevant promise, token issue, token redemption, pure-function result,
observation, correction, and merge decision is a parent-linked CAS object. A
branch is an ordinary local or group timeline. Token double-spend is represented
as conflicting redemption branches rather than hidden mutable state.

### Alternative C: Derive local projections from CAS event timelines

Fast local projections, caches, indexes, JSON files, SQLite tables, or in-memory
summaries may exist for efficiency, but the durable semantic record is the local
CAS event timeline. If a derived projection and the CAS timeline disagree, the
receiver's durable promise graph is the source used for explanation, replay,
merge, and trust review.

## Scenario analysis

### Scenario 1: Promise as timeline assertion

Alice promises that a package weighed 2.40 kg at timestamp T. In a timeline
model, that promise is not just a transient message. It is a CAS object whose
parents identify the local device state, prior calibration promises, and any
earlier related package promises Alice wants receivers to consider.

Bob receives the promise. Bob does not need to accept it as global truth. Bob can
store it as "Alice promised X about package Y at T" and later decide how much to
trust Alice, the scale, the calibration chain, or the transport path. If Alice
later discovers a calibration error, she adds a correction promise that
supersedes or narrows the earlier promise. The old promise remains visible as a
past assertion; the correction changes Alice's later timeline rather than
rewriting history.

Alternative A can store the message, but it does not force the correction and
prior state to become explicit graph structure. Alternative B makes the timeline
model clean and replayable. Alternative C allows indexes for speed while keeping
the durable explanation in CAS.

### Scenario 2: Pure-function agent with explicit context

Bob promises to act as a pure-function server. Given function CID F, input CID I,
and context CID C, Bob promises result CID R. If the computation uses a clock,
randomness, sensor data, peer quotes, exchange rates, or a local model version,
those values must be explicit inputs inside C. For the same F, I, and C, Bob
should always return the same R or later promise that an earlier result was
broken.

This is stronger than "Bob runs code." It says Bob's service promise is about a
stable relation among content-addressed objects. Carol can recompute locally, or
ask Dave to compute the same tuple, without requiring a central compute
authority. If Bob returns different results for the same F/I/C tuple, receivers
can identify the inconsistency as promise evidence.

Alternative A can cache result records but may leave ambient inputs implicit.
Alternative B demands fully explicit context and gives the best replay story.
Alternative C is practical if implementations use runtime caches but write the
CAS timeline as the durable explanation.

### Scenario 3: Local timeline branches

Alice and Bob each keep local branches over a shared topic. Alice's branch says
that token T was issued to Carol and redeemed for storage at time A. Bob's branch
says that he observed token T presented later by Mallory. Neither branch is the
global truth. They are local timelines: Alice's promise history and Bob's
promise/observation history.

When Alice and Bob sync, each may choose whether to retain the other's branch,
request missing parents, merge selected objects, or ignore the branch. If Bob
trusts Alice's issuer history but distrusts Mallory, Bob may mark Mallory's
presentation as suspicious without treating Alice as broken. If Alice later
promises that T was bearer-transferable and not expired, Bob may update his local
branch differently.

Alternative A cannot explain this cleanly if token state lives only in mutable
spent-token tables. Alternative B makes branches first-class. Alternative C
permits local projections while keeping branch facts explainable from CAS.

### Scenario 4: Agreed group timeline

Alice, Bob, and Carol may voluntarily promise to maintain a shared reference set
for a subset of their timelines. For example, they may all promise to track the
same project release branch or the same shipment audit branch. The group branch
is not globally authoritative. It is a shared artifact backed by the agents'
reciprocal promises to maintain, sync, and interpret it similarly.

If Dave joins later, Dave may accept the group branch because he trusts Alice and
Carol, or he may fork his own branch and reconcile only some objects. If Carol
leaves or breaks promises, the group may continue without her, fork around her,
or create a new group branch with different membership promises.

Alternative B gives the clearest expression of "agreed portion of the timeline."
Alternative C is acceptable if implementations maintain local materialized views
of group branches as derived indexes over CAS objects.

### Scenario 5: Double-spend as branch conflict

Mallory receives a bearer token issued by Alice. Mallory presents it to Bob on
one branch and to Carol on another branch. In a hidden mutable-index model, Bob
or Carol tries to decide whether the token is already spent by consulting local
state or the issuer. In a branch model, both redemptions can exist as promises on
parallel branches.

The important question is not whether the universe allows the bytes to appear
twice; it does. The question is how each agent interprets the branches. Bob may
accept his branch because he saw the token first. Carol may accept hers because
her branch has stronger parent evidence. Alice may promise that the token's
pCID-defined semantics allow only one successful redemption per branch family.
Dave may receive both branches and decide that the double-spend lowers trust in
Mallory, not necessarily in Alice.

This model avoids pretending there is a global clock or global ledger. It also
lets a later merge explicitly represent the conflict: both redemptions are
visible, but the merged branch may mark one redemption as accepted, both as
compensated, or the branch family as unmergeable under that token protocol.

Alternative A is simpler but hides branch structure and may recreate a central
ledger instinct. Alternative B tests the strongest PromiseGrid-native model.
Alternative C is likely implementation-practical: fast local replay projections
plus durable branch conflict objects.

### Scenario 6: Merge or non-merge after double-spend

Carol receives Bob's accepted-redemption branch and Ellen's conflicting
accepted-redemption branch. Carol asks whether the token pCID permits both,
whether the issue promise named a single-use family, whether either branch
contains expiry or transfer evidence, and whether Mallory promised exclusivity.

Carol may choose to:

- keep Bob's branch and reject Ellen's;
- keep Ellen's branch and reject Bob's;
- merge both branches with a conflict promise;
- require compensation before merging;
- keep both branches unmerged and lower trust in Mallory;
- ask Alice for an issuer promise about the intended token family semantics.

None of these choices requires Carol to command Bob, Ellen, Alice, or Mallory.
They are Carol's local promises and local trust updates.

## Conclusions

Alternative C is the recommended next path. It preserves the semantic clarity of
Alternative B while admitting that production implementations will need indexes,
caches, and fast local projections. The durable model should be: promises live in
parent-linked local CAS timelines; mutable projections are derived views unless a
pCID explicitly defines an additional promise object.

POC20 should therefore test:

- promise objects as timeline assertions;
- pure-function service promises over explicit context objects;
- local branches and group-agreed branches over CAS objects;
- private, encrypted-shareable, and plain-shareable local CAS objects, with
  sharing treated as a separate local promise decision;
- token issue, transfer, redemption, double-spend, and merge/non-merge behavior
  as branchable promise history;
- local trust decisions based on branch evidence without a global ledger or
  global monitor.

POC19 should not be blocked, but POC19 should avoid choices that would prevent
POC20's model. In particular, POC19 should keep token state explainable by CAS
objects, keep parent links visible, keep pCID-defined token semantics explicit,
and avoid hiding important state in unreviewable daemon internals or
non-rebuildable local projections.

## Implications for open TODOs and pending DIs

- `TODO-vumas` / POC19: keep the production-shaped plumbing path moving, but add
  a note that POC20 owns the deeper promise-timeline semantic model. POC19
  storage and token choices should not preclude branch-based promise histories.
- `TODO-nahop` / POC18: POC18 reference sets, snapshots, parent links, and
  continuous peer DAG sync are the closest existing implementation substrate.
- Future POC20 TODO: define an executable scenario with Alice, Bob, Carol, Dave,
  Ellen, and Mallory exercising pure-function results, local/group timelines,
  and double-spend branches.
- Future protocol docs: token pCIDs should define whether double-redemption is
  forbidden, compensable, branch-local, or pCID-specific rather than assuming a
  universal answer.

## Decision status

decided by `DI-mokaz`. POC20 should use local CAS timelines as the durable
semantic model while permitting derived local projections, indexes, caches, JSON
files, SQLite tables, or in-memory summaries only when they are rebuildable from
local CAS. The follow-up DF still needs to lock the first executable POC20 slice
before code generation.
