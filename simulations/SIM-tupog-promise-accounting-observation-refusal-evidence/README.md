# SIM-tupog-promise-accounting-observation-refusal-evidence

This higher-layer simulation preserves the useful pressure from the rejected
`hadit`/`jogoh` children without moving generic assertion machinery into the
base envelope. It tests pCID-owned promise-accounting payloads
for signed refusal, silence/timeout, exact-byte observation evidence, and local
trust updates. Source: `DI-kafiz`.

## Design Under Test

Alice asks Bob to store or relay exact bytes under a pCID-named protocol. Bob may
promise to accept the work, may sign a refusal, or may remain silent until Alice
records a timeout. Each record is evidence for a local trust assessment, not a
global judgment:

- Bob signs only Bob's own promise or refusal.
- Alice records exact bytes she observed and her local interpretation of Bob's
  response, timeout, or later keep/break evidence.
- Carol may receive Alice's evidence, but Carol decides locally whether Alice's
  observations matter to Carol's relationship with Bob.
- The base grid envelope remains pCID-selected; this simulation lives in the
  payload protocol that defines the observation/refusal record shape.

## Local Draft Spec

The local draft spec in
`protocols/promise-accounting-evidence.d/specs/promise-accounting-evidence-draft.md`
defines a candidate payload protocol for observation/refusal evidence. The draft
does not freeze a pCID and does not define a universal assertion artifact.

## Boundaries

This simulation does not promote `hadit` or `jogoh`, does not add rejected
multi-selector envelope stacks, and does not make refusal evidence
authoritative. It only asks whether higher-layer pCID-owned evidence records are
enough for PromiseGrid peers to distinguish signed refusal from silence or
timeout.
