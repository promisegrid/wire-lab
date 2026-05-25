# SIM-fovip: Kernel promise boundary port contract

This simulation is the concrete evidence surface for the next `DR-davod` packet.
It follows `TE-mazop`, which found that `TE-jimar` is not enough to close the
kernel-developer porting boundary because the app/kernel/host promise surface
and minimum first-port contract are still under-specified. Source: `DI-funaf`.

## Question

What minimum PromiseGrid port contract can be claimed across rich native nodes,
browser/WASM hosts, mobile sandboxes, MCU/header-only ports, and split local
service graphs without pretending that one process shape is universal?

## Candidate specimen

The specimen under test is a declared profile record. A port promises its own
behavior by publishing:

- profile name and runtime class;
- supported pCIDs and unsupported-pCID behavior;
- app-facing promises for storage, compute, network, key use, device access,
  lifecycle, pCID dispatch, and evidence recording;
- host/runtime assumptions that the port depends on but does not itself promise;
- explicitly unsupported features;
- exact evidence records for kept, refused, unavailable, and broken promises.

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

## Boundaries

This simulation does not close `DR-davod` and does not define a final
PromiseGrid kernel API. It tests whether a profile-declared port contract gives
guide writers enough evidence to discuss kernel developers without promising a
daemon, microkernel, browser host, mobile runtime, or MCU library as the single
correct implementation shape.

The current envelope shape `grid([42(pCID), payload, ...])` is input evidence,
not a reopened decision.

## Related evidence

- `docs/thought-experiments/TE-jimar-kernel-runtime-portability-boundary.md`
- `docs/thought-experiments/TE-mazop-kernel-promise-boundary-and-minimum-port-contract.md`
- `DR/DR-davod-promisegrid-kernel-dev-porting-boundary.md`
- `scenarios/kernel-porting-boundary/kernel-porting-boundary.md`
- `simulations/SIM-funas-kernel-porting-boundary/`
