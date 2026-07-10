# POC20 timeline pure-function CAS branches design

## Status

Design draft. This is the first implementation-local POC20 artifact. It is not
executable code, not a frozen protocol spec, and not a production API. The
canonical thought experiment remains `TE-lodom`, and the canonical task record
remains `TODO-nudav`. Source: `DI-kakos`; `DI-bibah`; `TE-lodom`; `TODO-nudav`.

## Purpose

POC20 should test the semantic model that POC19 should not accidentally block:
promises as assertions about part of the universe on a timeline, agents as
deterministic pure-function servers over explicit content-addressed context,
decentralized CAS object chains as local or group timelines, and
capability-token double-spend as visible branch evidence rather than hidden
mutable state.

POC20 is parallel to POC19. POC19 owns production-shaped plumbing: daemon/client
shape, transports, VCS/CAS app installation, WASI execution, and runtime
capability promises. POC20 owns the deeper promise timeline model. POC20 should
reuse POC18 and POC19 lessons, but it should not be treated as a POC19
code-generation blocker. Source: `DI-kakos`.

## Core model

Every important POC20 event should be represented as a promise-shaped CAS object
with parent links. Local runtime indexes may exist later, but they are derived
views. The durable explanation is the visible CAS timeline.

The first POC20 model uses these rules:

- A promise is an assertion by a promiser about some part of the universe at a
  point or interval on a timeline.
- A pure-function agent promises that for function CID `F`, input CID `I`, and
  context CID `C`, it returns result CID `R`; timestamps, randomness, sensor
  reads, model versions, exchange rates, and peer observations must be explicit
  context objects.
- A local timeline is an agent's parent-linked promise branch. It is not global
  truth.
- A group timeline is a reference-set-like branch that multiple agents
  voluntarily promise to maintain or interpret together.
- A token issue, transfer, redemption, double-spend, or merge decision is a
  pCID-defined promise object, not hidden ledger mutation.
- A double-spend can appear on parallel branches. Receivers decide locally
  whether to keep, reject, merge, compensate, or leave branches unmerged.

## First executable target

The first executable POC20 slice should use Alice, Bob, Carol, Dave, Ellen, and
Mallory:

- Alice issues a bearer capability token as a promise object on her local
  timeline.
- Bob provides a pure-function service and promises a result for explicit
  function, input, and context CIDs.
- Carol and Dave maintain local branches and optionally agree on a shared group
  branch for selected token and compute history.
- Mallory presents the same bearer token on two parallel branches.
- Ellen receives both branches and records a local merge, non-merge, or
  compensation promise without acting as a global authority.

The run is successful only if the double-spend is visible as branch evidence in
CAS and can be explained without a hidden mutable spent-token table as the source
of truth.

## Non-goals

- No global ledger, global branch authority, global trust authority, or global
  monitor.
- No claim that every production token pCID must use the same double-spend
  semantics.
- No final PromiseGrid token API, app API, or storage profile.
- No requirement that POC19 wait for POC20.
- No code generation until `TE-lodom` is reviewed and the executable slice is
  locked by follow-up DF in `TODO-nudav`.

## Acceptance criteria for future code

The future executable POC20 should not be considered complete until a clean run
shows:

- visible parent-linked CAS timelines for all participating agents;
- at least one voluntarily agreed group timeline branch;
- reproducible pure-function results from explicit function, input, and context
  CIDs;
- token issue, transfer, redemption, double-spend, and merge/non-merge as
  promise-shaped CAS objects;
- a double-spend represented as branch conflict rather than hidden mutable
  ledger state;
- local receiver decisions that can keep, reject, merge, compensate, or leave
  branches unmerged without commanding other agents;
- diagnostics that render the relevant raw CBOR/CAS objects for review.
