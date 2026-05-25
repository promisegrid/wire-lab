# SIM-fovip: Kernel promise boundary port contract

This simulation is the active evidence surface for the `DR-davod` kernel design
packet. It follows `TE-mazop`, which found that `TE-jimar` was not enough to
close the kernel-developer porting boundary, and now incorporates the follow-on
`TE-pudiv`, `TE-dunas`, and `TE-gakoh` questions. Source: `DI-funaf`;
`DI-fidot`.

## Question

What minimum PromiseGrid kernel implementation promises can be claimed across
rich native nodes, browser/WASM hosts, mobile sandboxes, MCU/header-only ports,
and split local service graphs without pretending that one process shape,
namespace, API, or prior-art pattern is universal?

## Candidate specimen

The specimen under test is a kernel implementation promise record: a local
implementation promises its own behavior by publishing:

- profile name and runtime class;
- supported pCIDs and unsupported-pCID behavior;
- app-facing promises for storage, compute, network send/receive, key use,
  device access, lifecycle, pCID dispatch, refusal, receipt, evidence,
  namespace, reference, and checkpoint behavior;
- host/runtime assumptions that the port depends on but does not itself promise;
- explicitly unsupported features;
- exact evidence records for kept, refused, unavailable, and broken promises;
- adapter promises when local APIs wrap pCID-selected grid messages;
- voluntary namespace promises when the port projects group namespaces;
- CID-rooted promise-bound reference behavior for cross-agent path sharing.

## Runtime classes

- **Native node:** Bob runs a local service with storage, network, keys, feed,
  CAS, dispatch, lifecycle, and evidence roles.
- **Browser/WASM:** Alice runs inside a host that owns storage, network, clocks,
  keys, and lifecycle.
- **Mobile sandbox:** Dave can run only while the OS permits background work and
  network access.
- **MCU/header-only:** Carol supports one or two pCIDs, bounded evidence, and a
  hardware/device promise.
- **Split local services:** Ellen separates dispatch, storage, keys, networking,
  and evidence into local services with separate promises.

## Basic principles under test

- Kernel is a role/profile set, not a ruler.
- Everything useful is a promise: app/kernel operations, resources, namespaces,
  references, and kernel implementation promise records all help agents make or
  evaluate promises.
- The app/kernel boundary is a promise boundary; exposed operations are
  pCID-selected `grid([42(pCID), payload, ...])` messages, even when local APIs
  provide ergonomic adapters.
- A kernel implementation promise record is not a global certificate. Alice,
  Bob, Carol, and later agents evaluate the record and make/break history
  locally.
- Host assumptions are not implementation promises unless the host is also an
  explicit promiser.
- Voluntary group namespaces may exist inside trust relationships, but imposed
  universal namespaces are rejected.
- File-like resources are promise-log projections or checkpoints, not evidence
  that PromiseGrid is filesystem-first.

## Evidence axes

The simulation should let reviewers ask whether the candidate:

- names the local promiser for each storage, compute, network, key, device,
  lifecycle, dispatch, namespace, reference, and evidence promise;
- maps every exposed app/kernel operation to a pCID-selected message or explains
  why the operation is outside the PromiseGrid boundary;
- records exact bytes when needed for proof, replay, unsupported-pCID carriage,
  or broken-promise evidence;
- states host/runtime assumptions separately from the port's own promises;
- names unsupported pCIDs and unsupported roles directly;
- keeps trust local and avoids global permission, namespace, conformance, or
  policy authority;
- treats V, Amoeba, Plan 9, and Hurd as pattern pressure, not imported design
  authority;
- supports voluntary group namespaces and CID-rooted promise-bound references
  without treating Alice's local path as global truth;
- represents file/resource state as checkpoints or projections over selected
  promise/event log frontiers.

## Boundaries

This simulation does not close `DR-davod` and does not define a final
PromiseGrid kernel API. It tests whether kernel implementation promises give
guide writers enough evidence to discuss kernel developers without promising a
daemon, microkernel, browser host, mobile runtime, MCU library, namespace
protocol, or SDK as the single correct implementation shape.

The current envelope shape `grid([42(pCID), payload, ...])` is input evidence,
not a reopened decision.

## Related evidence

- `docs/research/DN-lujad-promisegrid-kernel-role-profile.md`
- `docs/thought-experiments/TE-jimar-kernel-runtime-portability-boundary.md`
- `docs/thought-experiments/TE-mazop-kernel-promise-boundary-and-minimum-port-contract.md`
- `docs/thought-experiments/TE-pudiv-app-kernel-grid-message-boundary.md`
- `docs/thought-experiments/TE-dunas-prior-art-influence-on-promisegrid-kernel.md`
- `docs/thought-experiments/TE-gakoh-local-views-over-promise-event-hypergraph.md`
- `DR/DR-davod-promisegrid-kernel-dev-porting-boundary.md`
- `scenarios/kernel-porting-boundary/kernel-porting-boundary.md`
- `simulations/SIM-funas-kernel-porting-boundary/`
