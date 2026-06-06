# CID-named function compute

## Scenario ID

cid-named-function-compute

## Source / Provenance

- Source type: new harness scenario
- Source path: `protocols/wire-lab.d/TODO/TODO-godad-poc13-cas-compute-functions.md`
- Source row/title: `TODO-godad` planned POC13 CID-named compute pressure
- Source DI / TODO / TE: `DI-bibom`; `TODO-godad`; `DI-gopag`

## Purpose

Exercise candidate designs against PromiseGrid compute where protocol pCID names
the compute protocol and payload-level `function_cid` names CAS-stored function
code. The scenario tests whether pure and impure computation can remain
promise-first, replayable, cacheable where safe, and locally auditable.

## Setup

Alice stores a deterministic tax/shipping-cost function in CAS and receives a
`function_cid`. Bob promises to execute that function under a compute pCID for
explicit input objects. Carol promises to retain the function code and common
input objects. Dave promises to cache pure results when exact function, input,
context, ABI/version, and protocol identity match. Mallory presents a result
computed from a different function, missing input object, stale context object,
or hidden timestamp.

## Stimulus

Alice asks Bob to compute a result twice: first as a pure function over exact
input CIDs, then as an apparently impure function that needs a timestamp,
randomness, or sensor observation. The design must decide whether the latter can
be made pure-after-the-fact by turning every ambient input into an explicit
context object that is itself promised, stored, and auditable.

## Expected Pressure

The candidate design must keep pCID as protocol identity, keep `function_cid` as
payload-level CAS identity, separate compute execution from result trust, and
record enough evidence for Alice to decide locally whether Bob kept the compute
promise. Cache hits must be explainable from exact identity, not from a global
compute authority.

## Scenario-Specific Evaluation Questions

- Who promises that the function code bytes named by `function_cid` are retained
  and available?
- Who promises to execute the function, and what input/context objects are part
  of that promise?
- What exactly is the cache key for a pure compute result?
- How are timestamp, randomness, sensor, filesystem, network, or peer-observed
  inputs converted into explicit context objects?
- What local evidence lets Alice distinguish kept compute, malformed result,
  unavailable dependency, broken execution promise, and receiver
  non-commitment?
- How does the design avoid treating a function call as an RPC command rather
  than Bob's voluntary promise to compute under stated terms?
