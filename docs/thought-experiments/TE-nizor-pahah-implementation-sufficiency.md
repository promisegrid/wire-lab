# TE-nizor: Can implementing Pahah satisfy the recovered concerns?

*Thought experiment: considering TE-pahah, TE-vilot, TE-hirap, and turns 149-208, can the concerns recovered from the context-loss window be satisfied by implementing Pahah's simulation-first structure?*

## TE ID

TE-nizor

## Status

needs DF

First drafted 2026-05-09 19:57:57 UTC.

## Decision under test

DUT-nizor: Can wire-lab satisfy the concerns recovered from turns 149-208 by
implementing TE-pahah's simulation-first structure, or does the project need a
different immediate organizing decision before it can safely continue?

Steve's prompt said "pahaj"; this TE treats that as `TE-pahah`, because no
`TE-pahaj` exists in the current corpus and the intended reference is the
simulation-first TE just endorsed as "very good."

This is a critical decision because turns 149-208 show a repeated failure mode:
conversation discovered real design structure, but the repo kept trying to
encode the discovery as root-level production-looking paths, premature protocol
specs, or global artifact style rules. TE-pahah proposes a different answer:
make `simulations/` the primary home for concrete worlds and experiments, then
let results drive DR/DI/spec decisions.

## Short answer

Yes, **if "implement Pahah" means more than creating a directory**.

Pahah can satisfy the recovered concerns if implementation includes a minimal
simulation contract:

- each simulation states its question, protocol set, assumptions, actors/sites,
  world state, events, observations, results, and decision handoff;
- simulation worlds can contain scoped `sites/`, `groups/`, `cas/`, `feeds/`,
  `wires/`, and artifact/message-shape experiments without making those names
  repo-root commitments;
- recovered concerns from turns 149-208 are tracked as explicit simulation
  questions or result criteria, not lost in prose;
- apparatus documents remain TEs, DRs, TODOs, and guide-resource notes until
  results graduate through DF/DI;
- the first implementation targets dogfooding and recovery, not full production
  PromiseGrid architecture.

No, **if "implement Pahah" means only "mkdir simulations/"** or "move
`transports/wire-lab-devs-draft/` into a prettier path." A bare directory
skeleton would not satisfy the 149-208 concerns. The key is not the path itself;
the key is the apparatus/specimen boundary plus a repeatable way to run
competing worlds and record what they prove.

## Sources examined

This TE uses:

- **TE-pahah:** recommends top-level `simulations/` as the primary home for
  concrete experimental worlds.
- **TE-vilot:** says promise-shaped artifact protocols should be tested inside
  simulations before changing all apparatus documents.
- **TE-hirap:** extends that boundary to full PromiseGrid-message-shaped
  artifacts and text/CBOR representation bakeoffs.
- **Turns 149-208:** the recovered window that includes TE editing policy
  refinements, apparatus/specimen correction, message-envelope confusion,
  dogfood pressure, directory/layering churn, sites/CAS/feed/group discovery,
  promisebase prototype analysis, and the replay-ledger discipline.
- **Replay ledgers:** `docs/discussion/session-replay-72hr-ledger-20260504.md`
  and `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`
  as repo copies of recovered concern inventories.

## Concern inventory from turns 149-208

### C1 -- TE editing and context preservation

Turns 149-153 locked refinements for TE editing policy: Cat-2 DI enumeration,
cross-TE quotation checks, and top-of-file status markers. Turns 197-208 then
created the recovery discipline: write outstanding questions and unfinished
threads to the repo before changing threads.

Concern: future work must not rely on conversational memory. Every simulation
must have durable question/result/decision handoff files.

### C2 -- Apparatus versus specimen

Turns 155-160 discovered that the harness/spec was being used to prescribe one
candidate wire envelope instead of acting as lab apparatus. Turn 158 is the
load-bearing correction: wire-lab must test multiple hypotheses at all layers,
not name one as the harness answer.

Concern: root-level structures and harness prose must not make a specimen look
like settled architecture.

### C3 -- One envelope across transports, but not proven

Turns 157 and 206 clarified that wire-lab is looking for one message envelope
that can ride over all transports, with `grid([pcid, payload])` as a working
hypothesis, not a proven conclusion. Turn 207 then questioned whether "promise
stack" itself was a mistranslation of nested promise envelopes.

Concern: envelope candidates need a bakeoff surface. Pahah must not hard-code
promise-stack, text `grid <pcid>`, CBOR, or any other envelope as final.

### C4 -- Dogfood urgency

Turns 164-169 pushed hard toward file-based group communication so multiple
developer agents and humans could collaborate through wire-lab rather than
leaving Steve solo. Those turns also locked or nearly locked several practical
message rules: `.txt`, CID-as-filename, fetch/merge/push/post cycle, and
removing `Message-ID`.

Concern: the next structure cannot be so theoretical that it delays dogfooding.
It needs a small path for current messages while preserving experiment clarity.

### C5 -- Root path and vocabulary churn

Turns 170-176 iterated through `transports/`, protocol-slug nesting, feed versus
binding, `groups` versus `forums`, `wire-lab-devs-draft`, `grid envelope`,
headers, nested envelopes, and `draft--slug` versus `slug-draft`.

Concern: root-level names keep implying settled architecture. The repo needs a
place where competing naming and layer hypotheses can coexist without global
renames on every correction.

### C6 -- Sites, sparse CAS, feeds, groups, and chunking

Turns 175-178 discovered the deeper model: users are not sites; sites hold local
CAS replicas, run feeds, participate in groups, and make/assess promises. No
site has all CAS objects. Feeds may carry chunks, not merely messages. Rabin or
FastCDC chunking, pointer files, Merkle roots, and CIDv1 codec choices matter
but are not all needed for immediate dogfood.

Concern: a flat message directory cannot test federation. The first useful
simulation must represent sites, local CAS state, feeds, groups, and partial
replication even if some are stubbed.

### C7 -- Promise economy and trust relationships

Turns 177-179 made the social motivation explicit: PromiseGrid tries to avoid
centralized capture by making promises and peer-relative assessments first-class
at every layer. Conditional release, geofencing, multi-hop discovery/fetch, and
possibly capability-token economies must be testable without assuming one
economic model.

Concern: the structure must allow promise/trust/economic behavior to be modeled
as simulation behavior, not baked into a single premature protocol.

### C8 -- Promisebase as prototype, not authority

Turns 180-192 corrected a bad pivot: promisebase docs were partly refactoring
plans, while `db/` contains useful implemented CAS/Merkle/Rabin primitives.
Steve then locked the rule: promisebase is a prototype; when wire-lab and
promisebase conflict, discuss the conflict and prefer wire-lab.

Concern: simulations can use promisebase evidence or code as a specimen, but
wire-lab decisions remain canonical for PromiseGrid design until a later
handoff changes that.

### C9 -- Artifact promise/message shape

TE-vilot and TE-hirap generalize turns 176-179 and 207-208. Some artifacts may
need to be promise-shaped or full message-shaped objects, and CBOR/plain-text
representations need bakeoffs. But forcing every TE, DR, TODO, or guide note
into that shape would erase the apparatus/specimen distinction.

Concern: artifact-shape experiments belong inside simulations or commitment
templates first, not as a global repo rewrite.

### C10 -- Dev-guide source discipline

The guide is about PromiseGrid, not wire-lab. Wire-lab is evidence until the
guide and frozen specs stabilize. Turns 149-208 showed why future guide writers
need a durable source map and clear "this is settled / this is simulation
evidence / this is open" boundaries.

Concern: implementing Pahah must make results easy to cite without teaching
wire-lab internals as final PromiseGrid API.

## Alternatives

### Alt A -- Minimal Pahah

Create `simulations/` with Pahah's directory skeleton, then move current
bootstrap data under the first simulation world.

**Easier:** Fast. Low conceptual overhead. Establishes the path name.

**Harder:** Does not by itself encode the recovered concerns, artifact/message
protocols, representation policy, replay links, or graduation rules. It risks
creating a new empty apparatus that future agents fill inconsistently.

**Obligations:** Still requires later TEs/TODOs for almost every concern above.

### Alt B -- Pahah plus simulation contract

Implement `simulations/` with a required minimal contract: `README.md`,
`QUESTION.md`, `protocol-set.md`, `concerns.md`, `world/`, `events/`,
`observations/`, `results/`, and `decisions.md`. The first simulation imports
the current dogfood data as a specimen, explicitly tracks turns 149-208 concerns
as test criteria, and states which parts are stubs.

**Easier:** Preserves Pahah's insight while preventing context loss. Gives Vilot
and Hirap a place to run without changing global templates. Lets dogfood start
with a text-first specimen while still recording open CBOR/chunking/promise
economy questions.

**Harder:** Slightly more work than a bare directory. Requires discipline to
write small result summaries and avoid making the first simulation too large.

**Obligations:** Define the simulation contract before moving current data.
Treat simulation results as evidence that feed DRs/DIs, not as decisions by
themselves.

### Alt C -- Root model first

Implement top-level `groups/`, `sites/`, `cas/`, `feeds/`, and perhaps
`wires/` now, using the current best model from turns 175-178.

**Easier:** Concrete and satisfying. Dogfood paths look like the target model.

**Harder:** Repeats the mistake Pahah identified: a root-level tree looks like
settled architecture even when it is still a specimen. It makes every new
correction a repo-wide migration.

**Obligations:** Immediately decide unresolved layer, naming, chunking,
representation, and artifact-shape questions.

### Alt D -- TE/DR/TODO only; no simulations yet

Keep writing TEs and DRs until the model is clear enough to implement.

**Easier:** No path churn. Maximum analytical caution.

**Harder:** Fails the dogfood urgency from turns 164 and 193. Keeps many
questions abstract. Encourages more context compression because there is no
concrete world to inspect.

**Obligations:** Eventually still needs Pahah or an equivalent experiment home.

### Alt E -- Hybrid current dogfood root plus simulations

Keep one top-level current dogfood world for convenience and use `simulations/`
for bakeoffs.

**Easier:** Immediate path access for active collaborators.

**Harder:** Creates two authorities: the current root world and the simulation
worlds. Readers may treat the root dogfood tree as more real than the simulations
even when the root tree is just an implementation shortcut.

**Obligations:** Define which path is authoritative for evidence and how it
mirrors into simulations.

## Scenario analysis

### S1 -- Recovering from context loss

Alice starts after a context-loss event. She needs to know what turns 149-208
discovered and what remains open.

- **Alt A:** Alice sees `simulations/` but not why it exists. She still has to
  read TEs, TODOs, and session logs.
- **Alt B:** Alice enters the first recovery/dogfood simulation and sees the
  question, concern list, world state, events, observations, results, and
  decision handoff. This directly addresses the failure mode from turns 197-208.
- **Alt C:** Alice sees plausible production trees and may mistake them for
  settled architecture.
- **Alt D:** Alice reads many docs but no concrete world.
- **Alt E:** Alice must learn which of two worlds is authoritative.

S1 favors Alt B.

### S2 -- Dogfooding without blocking on full architecture

Bob wants the developer group to start using wire-lab-based messaging today.

- **Alt A:** Bob can move files quickly, but the path alone does not say which
  unresolved questions are being deferred.
- **Alt B:** Bob can start with a text-first dogfood simulation that explicitly
  stubs CBOR, chunking, promise economy, and multi-feed behavior. The simulation
  can record when those stubs become blockers.
- **Alt C:** Bob must pick root-level `groups/`, `sites/`, `cas/`, and `feeds/`
  semantics before using them.
- **Alt D:** Bob does not get a dogfood path.
- **Alt E:** Bob gets a shortcut but risks future merge/confusion.

S2 favors Alt B, with an explicit constraint: do not implement all Pahah
subdirectories in production depth before dogfooding.

### S3 -- Testing one envelope across all transports

Carol wants to test whether `grid([pcid, payload])`, nested grid envelopes,
plain text `grid <pcid>`, CBOR promise stacks, or another shape should become
the shared envelope.

- **Alt A:** The skeleton can host a bakeoff only if someone invents conventions
  later.
- **Alt B:** A simulation can define an envelope-bakeoff protocol set, run the
  same scenario over git-feed, websocket-feed, and uucp-feed specimens, and
  store results without changing the harness spec.
- **Alt C:** Root-level current-world paths imply one current envelope.
- **Alt D:** Envelope debate remains abstract.
- **Alt E:** The root dogfood envelope may bias the simulation.

S3 favors Alt B.

### S4 -- Modeling sites, sparse CAS, feeds, and promise economy

Dave wants to model no-site-has-all-content, multi-hop discovery, conditional
release, geofencing, chunk-storage promises, and peer-relative trust.

- **Alt A:** Possible, but underspecified.
- **Alt B:** The simulation world can contain `world/sites`, `world/cas`,
  `world/feeds`, `world/groups`, promise/assessment logs, and explicit
  observations. The model can start with small fake CIDs and later swap in
  promisebase or real chunking.
- **Alt C:** Forces immediate root-level semantics for each layer.
- **Alt D:** Cannot observe distributed behavior.
- **Alt E:** Can work, but split authority complicates trust observations.

S4 favors Alt B. This is the strongest argument that Pahah is not just tidy
filesystem design; it is the right simulation apparatus.

### S5 -- Promisebase integration

Ellen wants to use promisebase as evidence or code without letting it override
wire-lab decisions.

- **Alt A:** No explicit slot for promisebase-as-specimen.
- **Alt B:** A simulation can include `protocol-set.md` and `world/cas` notes
  saying "promisebase db/ is the candidate CAS implementation under test" while
  preserving wire-lab authority in `decisions.md`.
- **Alt C:** Root-level CAS paths may imply promisebase or another CAS model is
  already the architecture.
- **Alt D:** Keeps promisebase in prose only.
- **Alt E:** Ambiguous if the root dogfood path uses promisebase differently
  than the simulation.

S5 favors Alt B.

### S6 -- Artifact promises and message-shaped artifacts

Frank wants to test whether simulation events, observations, conformance
claims, or review replies should be PT promises or PromiseGrid messages.

- **Alt A:** No standard place to state the artifact protocol.
- **Alt B:** A simulation can include `artifact-protocol.md` or a section in
  `protocol-set.md`, define text/CBOR canonicality, and keep ordinary TEs/DRs
  outside the specimen.
- **Alt C:** Root-level artifacts may become message-shaped by accident.
- **Alt D:** No concrete artifact protocol gets tested.
- **Alt E:** Two places can disagree about whether artifact-message shape is
  required.

S6 favors Alt B and reinforces Vilot/Hirap.

### S7 -- Guide writing and external contributor onboarding

Ellen writes the PromiseGrid Development Guide; Mallory tries to cite a
specimen as final API.

- **Alt A:** The guide still needs to infer which simulation outputs matter.
- **Alt B:** Results and decisions are explicit. Guide writers can cite a
  simulation result as evidence and a DI/frozen spec as authority.
- **Alt C:** Mallory can point at root-level `groups/` or `cas/` and claim
  finality.
- **Alt D:** Guide writers get evidence scattered across TEs.
- **Alt E:** Mallory can exploit the hybrid shortcut.

S7 favors Alt B.

### S8 -- Long-horizon evolution

Years later, Gina wants to understand why root-level `groups/` was not adopted,
why CBOR was deferred, or why chunking was stubbed.

- **Alt A:** Skeleton history is weak evidence.
- **Alt B:** Each simulation captures the question, assumptions, alternatives,
  events, results, and decision handoff. This is the best long-horizon record.
- **Alt C:** Root history shows churn, not necessarily rationale.
- **Alt D:** TEs explain rationale but may lack observed evidence.
- **Alt E:** Hybrid roots age badly unless mirrored carefully.

S8 favors Alt B.

## Rejected alternatives

- **Reject Alt A as insufficient.** Pahah's directory name is not the design.
  The simulation contract is the design.
- **Reject Alt C as premature.** Root-level layer trees are exactly what caused
  `transports/`/`groups`/`feeds` churn to look like architecture instead of
  specimens.
- **Reject Alt D as too slow.** It ignores the dogfood pressure and keeps
  critical distributed-system questions abstract.
- **Reject Alt E as a default.** A top-level dogfood shortcut may become useful
  later, but only after the simulation authority model is settled.

## Surviving alternatives

- **Alt B survives as the recommended interpretation of "implement Pahah."**
  Implement the simulation-first tree with a minimal contract and use it as the
  boundary for recovered concerns, promise/message artifact experiments, and
  dogfood specimens.
- **Alt E survives as a possible later convenience.** If a current dogfood root
  is needed for ergonomics, it should mirror or be generated from a named
  simulation, not supersede it.

## What "implement Pahah" must include

Implementing Pahah sufficiently means the first pass creates at least:

```text
simulations/
  README.md
  SIM-<handle>-<slug>/
    README.md
    QUESTION.md
    concerns.md
    protocol-set.md
    world/
      sites/
      groups/
      cas/
      feeds/
      wires/
    events/
    observations/
    results/
    decisions.md
```

The files do not need to be elaborate. The point is that every concrete world
has:

- a question under test;
- scoped layer names that do not claim repo-root authority;
- a list of recovered concerns or test criteria;
- an explicit protocol/representation/artifact-shape policy;
- event and observation records;
- result summaries;
- and a handoff to DR/DI/spec work.

## What Pahah does not decide

Even with Alt B, implementing Pahah does **not** decide:

- whether `grid([pcid, payload])` is the final envelope;
- whether current group-session plain text graduates to canonical PromiseGrid
  wire format;
- whether all artifacts are promises;
- whether all artifacts are PromiseGrid messages;
- whether CBOR or plain text is canonical for any artifact class;
- whether Rabin/FastCDC chunking is required in the first CAS protocol;
- whether promisebase becomes a dependency, specimen, implementation substrate,
  or separate evolving repo;
- what the production PromiseGrid node filesystem looks like.

That non-decision is a feature. Pahah is an apparatus decision: it creates the
place where those specimen decisions can be tested.

## Recommended conclusion

Adopt Alt B. Treat `TE-pahah` as the immediate structural decision to implement,
but define "implement" as **simulation-first with contract**, not "make a
directory."

This satisfies the concerns from turns 149-208 better than any root-level model:

- It preserves apparatus/specimen separation.
- It gives dogfooding a concrete but explicitly experimental home.
- It gives sites/CAS/feeds/groups a scoped world where they can be modeled
  before root-level adoption.
- It gives Vilot/Hirap a place to test promise-shaped and message-shaped
  artifacts without global rewrites.
- It gives recovery and guide writers durable, citable result records.

The remaining risk is scope creep. The first implementation should not try to
build full CAS chunking, CBOR tooling, promise economy, and multi-feed behavior
at once. It should create the simulation contract, import or reference the
current dogfood evidence, and mark missing capabilities as stubs or future
simulation questions.

## DF questions exposed

### DF-nizor.1 -- Is Pahah-with-contract sufficient as the immediate umbrella decision?

Recommended answer: yes.

Surviving alternatives:

- **1.A -- Yes, Pahah-with-contract (recommended).** Implement `simulations/`
  with a minimal contract and use it as the boundary for recovered concerns.
- **1.B -- No, implement root-level model first.** Adopt root `groups/`, `sites/`,
  `cas/`, and `feeds/` before simulations.
- **1.C -- No, write more TEs first.** Continue analysis before creating any
  simulation structure.

### DF-nizor.2 -- What is the minimum simulation contract?

Recommended answer: question, concerns, protocol set, world, observations,
results, and decision handoff.

Surviving alternatives:

- **2.A -- Full minimal contract (recommended).** Require `QUESTION.md`,
  `concerns.md`, `protocol-set.md`, `world/`, `events/`, `observations/`,
  `results/`, and `decisions.md`.
- **2.B -- Light contract.** Require only `README.md`, `world/`, and `results/`.
- **2.C -- No fixed contract.** Let each simulation choose its own shape.

### DF-nizor.3 -- What should the first simulation be?

Recommended answer: a recovery/dogfood simulation seeded by current
`wire-lab-devs` evidence.

Surviving alternatives:

- **3.A -- Recovery/dogfood simulation (recommended).** Seed from the current
  dogfood message evidence and turns 149-208 concerns.
- **3.B -- Envelope bakeoff simulation.** Start with `grid([pcid,payload])`,
  promise-stack, text, and CBOR envelope variants.
- **3.C -- Sites/CAS/feed simulation.** Start with sparse CAS and multi-hop feed
  behavior.
- **3.D -- Artifact-shape simulation.** Start with promise/message-shaped
  artifacts from Vilot/Hirap.

### DF-nizor.4 -- How should turn-recovery concerns enter simulations?

Recommended answer: copy concern IDs into `concerns.md` and map them to results.

Surviving alternatives:

- **4.A -- Concerns matrix (recommended).** Each simulation has `concerns.md`
  mapping recovered concerns to world elements, events, observations, and
  results.
- **4.B -- Free prose only.** Mention concerns in `README.md`.
- **4.C -- TODO-only.** Keep concerns only in TODO/DR files and do not repeat
  them inside simulations.

### DF-nizor.5 -- How much architecture should the first implementation include?

Recommended answer: stub broad layers, implement only what dogfood needs.

Surviving alternatives:

- **5.A -- Stub broad layers, implement narrow dogfood (recommended).** Include
  `sites/`, `groups/`, `cas/`, `feeds/`, and `wires/` directories, but use small
  stub files and current text messages until evidence demands more.
- **5.B -- Implement full layered model now.** Add real CAS/chunking/feed logic
  immediately.
- **5.C -- Dogfood only.** Ignore sites/CAS/feeds until later.

### DF-nizor.6 -- What graduates from a simulation?

Recommended answer: results feed DR/DI/spec work; simulations do not decide.

Surviving alternatives:

- **6.A -- Results feed DR/DI/specs (recommended).** Simulation `results/` and
  `decisions.md` point to DRs and DIs; they do not lock decisions alone.
- **6.B -- Simulation results can directly lock policy.** A result file can be
  authoritative if reviewed.
- **6.C -- Simulations are illustrative only.** They cannot be cited as evidence
  in decisions.

## Implications for open TODOs, DRs, and DIs

- **TE-pahah:** should be treated as the current structural umbrella, but only
  under the richer "with contract" interpretation described here.
- **TE-vilot / TE-hirap:** become downstream uses of the simulation boundary,
  not separate reasons to delay implementing Pahah.
- **DR-nugog:** should not decide root `transports/`/`groups` layout until
  DF-nizor.1 and DF-nizor.2 settle whether the current evidence moves into a
  named simulation.
- **TODO-kugod:** should keep residual recovery open until at least one
  simulation concern matrix demonstrates where turns 149-208 concerns now live.
- **DEV-GUIDE-RESOURCES.md:** should continue warning guide writers that
  simulation results are evidence, not final PromiseGrid API, until results point
  to frozen specs or locked DIs.
- **Future implementation TODO:** if DF-nizor.1/2 lock the recommendation, file
  or update a TODO for creating `simulations/README.md` and the first
  recovery/dogfood simulation skeleton.

## Decision status

`needs DF`. This TE recommends implementing Pahah as a simulation-first
structure with a minimal simulation contract. It concludes that Pahah can
satisfy the recovered concerns from turns 149-208 as an apparatus decision, but
cannot by itself settle the protocol, artifact, representation, CAS, feed, or
promise-economy specimen decisions.
