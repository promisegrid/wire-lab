# SIM-dutam-grid-envelope-fixed-header-variable-body

This simulation is the Promise-Theory-first successor to
`SIM-sajar-grid-envelope-variable-arity-pcid-defined-fields`. It keeps the
desire for flexible body structure but narrows the universal outer contract to a
small fixed header plus a pCID-defined variable body. Source: `DI-tavaz`.

## Design Under Test

The outer envelope keeps a stable, cross-peer-auditable header:

- a fixed slot for the layer or body pCID;
- a fixed slot for the body bytes or body reference;
- an optional fixed evidence/proof slot when the variant needs one.

Only the body shape varies under the pCID-defined protocol. That preserves room
for flexible payload families without making the entire outer envelope opaque to
generic later audit.

## Why this differs from `sajar`

`sajar` let the first-slot pCID define the entire outer arity, which pushed too
much stable meaning behind handler-specific interpretation. `dutam` keeps the
body flexible while preserving a small common outer shape that later peers can
reason about generically.

## Boundaries

This simulation does not pick the final fixed header. It tests whether a small
stable outer contract is a better PromiseGrid fit than fully variable outer
arity.
