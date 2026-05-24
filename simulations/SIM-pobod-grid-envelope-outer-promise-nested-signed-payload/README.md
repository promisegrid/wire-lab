# SIM-pobod-grid-envelope-outer-promise-nested-signed-payload

This simulation is the Promise-Theory-first successor to
`SIM-janov-grid-envelope-layer-pcid-nested-signed-payload`. It keeps the nested
signed-payload idea but makes the outer-layer promise explicit: an agent
promises that the carried nested bytes match the pCID-defined structure and that
the nested signature is the evidence object for the inner payload claim. Source:
`DI-tavaz`.

## Design Under Test

The outer envelope remains small, but its semantics are restated in Promise
Theory terms:

- the sender promises that the outer pCID defines how to parse the nested body;
- the nested body carries the actual payload pCID, payload bytes, and signature;
- receivers decide locally whether they trust the sender enough to believe the
  outer promise and whether the nested signature is enough evidence for the
  inner payload claim;
- no unsigned outer field is treated as self-authenticating authority.

## Why this differs from `janov`

`janov` described the shape well but left the outer-layer promise implicit,
which made authorship and trust too dependent on unstated transport context.
`pobod` keeps the same structural pressure while making the promise boundary
explicit.

## Boundaries

This simulation does not declare the winning grid-envelope format. It only asks
whether a nested signed payload becomes more PT-clean once the outer-layer
promise is made explicit.
