# POC18 Tangled prior-art review

This design note records the `nahop.18` review of Tangled as adjacent prior art
for POC18. It is not a frozen PromiseGrid protocol spec and it does not adopt
Tangled terminology as PromiseGrid terminology. It answers what POC18 should
learn from, reject from, or explicitly differ from in Tangled before POC18
implementation locks. Source: `DI-dibut`; `DI-zuruj`.

## Source provenance

Primary sources reviewed on 2026-06-27:

- Tangled homepage: <https://tangled.org/>
- Tangled docs single-page edition: <https://docs.tangled.org/single-page>
- Tangled introduction: <https://blog.tangled.org/intro>
- Tangled pull-request lifecycle: <https://blog.tangled.org/pulls>
- Tangled Jujutsu / stacked-PR article: <https://blog.tangled.org/stacking>
- Tangled six-month retrospective: <https://blog.tangled.org/6-months>
- Tangled source monorepo page: <https://tangled.org/tangled.org/core>

Those sources were enough to review Tangled's public architecture and workflow
claims at the level needed for POC18 planning. This note does not claim to audit
Tangled's implementation correctness, security properties, economics, or
long-term governance.

## Tangled summary

Tangled is a decentralized Git hosting and collaboration platform built on the
AT Protocol. Its public materials emphasize three goals: ownership of data, low
barrier to entry, and no compromise on user experience. It remains Git-first:
users create repositories, configure SSH keys, push and pull Git repositories,
migrate existing Git remotes, mirror to multiple Git remotes, and preserve Git
branches and tags.

Tangled's distinctive hosting unit is the **knot**. A knot is a lightweight,
headless Git server that can be single-tenant, multi-tenant, self-hosted on
small infrastructure, or managed by Tangled. Knots host Git repositories and
serve Git operations. Tangled also has an **appview**: the public web view that
aggregates repositories across knots so users can browse, clone, and contribute
without manually caring where each repository lives.

Tangled uses AT Protocol accounts and DIDs for identity. It also exposes
Tangled-owned XRPC APIs and has been moving toward more AT-native data. Its
six-month retrospective says a centralizing knot/appview registration secret was
removed: knots and spindles declare an owner, and appviews can verify ownership.
That is a useful decentralization move, but the public product shape still has a
consolidated appview experience.

For review workflow, Tangled's useful idea is that a pull request has explicit
rounds. A submission begins at round 0; the author chooses when to resubmit and
advance to the next round. Prior rounds remain reviewable, review comments stay
attached to the submission they addressed, and interdiffs can show what changed
between rounds. Tangled also supports stacked pull requests using Jujutsu change
IDs, so a rewritten commit can still be matched to the same logical change.

## Comparison to POC18

POC18 is not trying to be a better Git forge. It is trying to test whether
PromiseGrid can replace the Git/GitHub collaboration substrate with sparse CAS,
parent-linked envelopes, versioned reference-set promises, continuous peer DAG
sync, and local trust. Source: `DI-zuruj`; `DI-dibut`.

The closest Tangled match to POC18 is not a one-to-one component mapping. It is
a set of pressure points:

- Tangled knots prove that self-hosted repository infrastructure and managed
  convenience can coexist in a user-facing code-collaboration system.
- Tangled appview proves that a consolidated social/code view is important for
  user experience, but POC18 must not let an appview become a PromiseGrid
  authority.
- Tangled's SSH/Git compatibility proves that migration paths matter, including
  ordinary Git push and pull. POC18 should support those operations through a
  Git bridge while keeping them separate from native PromiseGrid sync.
- Tangled's round-based review model proves that immutable review submissions,
  interdiffs, and author-chosen resubmission points are useful.
- Tangled's Jujutsu support confirms the need for stable logical change identity
  across rewritten exact commits.

POC18's versioned reference-set design already covers much of that last point.
A logical change in POC18 should be a reference-set promise with its own CID,
parents, targets, promiser, terms, and proof. That is stronger than copying a
change-ID field into a commit-like object because the logical change is itself a
versioned promise object.

## Adopt, reject, or differ

| Tangled element | POC18 stance | Reason |
|---|---|---|
| Self-hosted knots | Adopt the lesson, not the Git-server shape | POC18 should make self-hosting and small-community hosting normal, but the storage unit is sparse CAS plus promises, not a Git repository as authority. |
| Managed convenience | Adopt as optional product pressure | A hosted PromiseGrid peer can exist, but it is only another promiser under local trust, not the source of truth. |
| Consolidated appview | Explicitly differ | A PromiseGrid appview may be a read model or discovery peer, but it must not decide global truth, global access, global merge status, or global repository identity. |
| AT Protocol identity | Adopt as interop pressure | POC18 can learn from DID/handle usability and future AT/Bluesky interop, but identity remains input to local trust, not a global trust authority. |
| Role-based access control | Reject as PromiseGrid vocabulary | POC18 should express collaboration as promises, capability tokens, local retention, local sharing, and local trust updates, not permission/conformance language. |
| SSH/Git push-pull | Adopt as required Git bridge compatibility | POC18 should push to and pull from ordinary Git repositories through shared bridge code, while native POC18 sync remains continuous peer DAG reconciliation: agents advertise, request, verify, retain, forward, or ignore reference-set and CAS promises. |
| Git branches/tags | Differ | POC18 uses versioned reference sets as the common abstraction for branches, tags, releases, directories, review threads, logical changes, and workspaces. |
| Hidden tracking refs for fork PRs | Learn from the problem, not the mechanism | POC18 still needs comparable-target materialization, but it should model those targets as reference-set promises rather than hidden Git refs. |
| Round-based reviews | Adopt the idea | POC18 review threads should preserve immutable submissions, explicit resubmissions, interdiffs, and review comments attached to exact versions. |
| Jujutsu change IDs | Adopt the requirement, differ in mechanism | POC18 needs stable logical change identity, but should represent it as a logical-change reference set rather than a raw change-ID field. |
| Spindles / CI | Defer but keep as compute-pressure input | POC18 may later attach test/build promises to review or release reference sets, but POC18 should not expand into CI runner design unless needed. |

## PromiseGrid-specific conclusions

Tangled is useful because it is Git-compatible and user-friendly while moving
some hosting and identity concerns out of a traditional forge. Its limits are
equally useful: a Git-first product retains Git remotes, refs, SSH, branch
operations, repository-local hidden refs, and role-based access vocabulary.
PromiseGrid should bridge to those Git concepts for interoperability but should
not inherit them as native concepts.

For POC18, the native collaboration loop should stay promise-first:

```text
Alice promises a versioned reference set.
Bob locally judges whether he trusts Alice's promise.
Bob asks trusted peers for missing target CIDs and parent-linked CAS objects.
Peers promise what they are willing to store, forward, verify, or retain.
Alice, Bob, and peers continuously update their own local DAG views.
```

This makes Tangled's "appview" concept a possible local or hosted view over
promises, not a global app authority. It makes a "knot" concept a possible
storage/transport peer, not a repository authority. It makes review rounds
reference-set versions, not mutable PR state owned by a forge.

## POC18 implications

POC18 should include these requirements before implementation locks:

- A review-thread reference-set role that can point at immutable submission
  versions, review comments, tests, acceptance promises, and interdiff material.
- A logical-change reference-set role that gives evolving work stable identity
  without copying Jujutsu's exact change-ID field as a primitive.
- A Git bridge path that preserves content and DAG semantics for import, export,
  push, and pull, while treating Git remotes, branches, tags, and pushes as
  compatibility artifacts.
- A discovery/read-model story that can look appview-like for users but is still
  local, sparse, and promise-relative.
- Analyzer gates that reject authority drift: no global forge, appview, role
  permission service, merge authority, branch authority, tag authority, or Git
  remote as the PromiseGrid source of truth.
- Scenario coverage for Alice, Bob, Carol, and Dave collaborating through sparse
  CAS peers where they do not all share the same complete repository view.

## Decision status

`nahop.18` is complete. Tangled should influence POC18's self-hosting,
migration, social-code UX, review rounds, stable logical-change identity, and
conventional Git push/pull bridge requirements. POC18 should explicitly differ
from Tangled by keeping Git/SSH push-pull out of the native sync model and by
keeping appview aggregation, role-based access control, hidden Git refs, and raw
Jujutsu change IDs out of the native PromiseGrid model unless a later TE/DI
narrows one of those choices.

The existing POC18 DIs remain sufficient for the original Tangled conclusion:
`DI-zuruj` locks versioned reference sets and the POC16 baseline, and `DI-dibut`
locks the Tangled prior-art review plus continuous peer DAG sync direction.
`DI-dofoj` refines the Git compatibility result by requiring a shared Git bridge
for import, export, push, and pull.
