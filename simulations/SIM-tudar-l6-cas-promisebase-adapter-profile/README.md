# SIM-tudar-l6-cas-promisebase-adapter-profile

This simulation turns the promisebase-adapter starting-profile alternative from
`SIM-bobud-l6-cas-starting-profile-bakeoff` into a concrete candidate specimen.
It tests whether PromiseGrid L6 CAS should start from promisebase / pitbase
block, tree, and stream prior art, then map that prior art into PromiseGrid
pointers, pCIDs, and peer-local promises. Source: `DI-fibuv`.

## Design Under Test

The first L6 CAS profile promises a compatibility adapter from existing
promisebase object and stream ideas into a PromiseGrid CAS profile, with explicit
translation rules for object identity, pointer shape, and stream evolution.

## Boundaries

This simulation does not make promisebase the canonical PromiseGrid storage
system. It tests whether adapter-first reuse gives better continuity than either
IPLD alignment or a minimal raw-pointer start.
