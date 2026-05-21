# SIM-funas: Kernel porting boundary

This simulation is a provisional question home for `FB-vitih`, `FB-mulum`, and
`FB-potin`: what "kernel developer" or "porter" should mean while `DR-davod`
remains open. Source: `DI-ragaz`.

## Question

What can the guide safely teach about the PromiseGrid porting target before the
stable kernel/runtime boundary is decided? Source: `DI-ragaz`.

## Decision Axes

- **Term of art:** kernel, runtime, dispatcher, handler host, or library surface.
- **Minimum viable port:** which frozen binding/session/message specs and
  conformance claims a port must implement first.
- **Runtime obligations:** handler dispatch, pCID routing, storage, key handling,
  ingress, feeds, CAS subtree, and implementation changelog evidence.
- **Layer boundaries:** what belongs to substrate/feed/group/CAS/app layers and
  what must not be taught as a monolithic harness clone.
- **Provisional versus blocked:** which K1-K5 ingress and runtime details are
  teachable orientation versus blocked settled instructions.

## Related Root Scenario

- `scenarios/kernel-porting-boundary/kernel-porting-boundary.md`

## Boundaries

This simulation does not define the final PromiseGrid porting API. It keeps the
kernel/runtime/dispatcher framing testable until `DR-davod` closes. Source:
`DI-ragaz`.
