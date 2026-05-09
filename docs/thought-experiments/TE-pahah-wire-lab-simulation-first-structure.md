# TE-pahah: Wire-lab simulation-first structure for design decisions

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-pahah

## Status

needs DF

## Decision under test

DUT-pahah: What repository structure best accomplishes the wire-lab's simulation and decision-making goals while PromiseGrid is still being derived, before production compatibility or active-use stability matter?

This TE starts from Steve's correction on 2026-05-09: the current trees are not production, not relied on by active users, and not compatibility-bound. They are brainstorming and experiment artifacts. Therefore the structure should optimize for:

1. making experiments easy to set up and compare,
2. making the decision under test explicit,
3. preserving enough evidence for later reasoning without treating every early path as sacred,
4. supporting multiple competing protocol/layer hypotheses side by side,
5. helping humans and LLMs find the right artifact quickly.

This differs from TE-domat's default posture. TE-domat over-weighted backward compatibility and historical preservation. Those concerns still matter for auditability, but they should not dominate the structure while the repo is still a lab.

## Short answer

The best structure for wire-lab's current purpose is **simulation-first**:

```text
simulations/
  README.md
  <sim-handle>-<slug>/
    README.md
    QUESTION.md
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

Keep `protocols/` as the home for candidate protocol/spec families and keep `docs/thought-experiments/`, `DR/`, and `protocols/*/TODO/` as the decision machinery. Put concrete experimental worlds under `simulations/`, not directly at repo-root layer names like `transports/` or `groups/`.

Inside each simulation, use the clearest layer names for the hypothesis being tested. For the current PromiseGrid candidate model, that means `sites/`, `groups/`, `cas/`, `feeds/`, and `wires/` or `transports/` as appropriate. The point is that those layer trees are scoped to a simulation run, so alternative structures can coexist without implying one global production layout.

## Why this TE now

The `transports/` -> `groups/` discussion exposed a deeper issue. We were asking which root-level production-looking tree should own current experimental artifacts. That is the wrong primary question for wire-lab.

The wire-lab is not yet a PromiseGrid node, not yet a production transport, and not yet the final dev guide. It is a decision apparatus. Its filesystem should make experimental questions, candidate structures, observations, and decisions visible.

In that frame:

- A root-level `transports/` tree suggests "this repo is currently simulating transport data as the canonical global state."
- A root-level `groups/` tree suggests "this repo is currently a PromiseGrid group host."
- A root-level `simulations/` tree says the true thing: "this repo contains experimental worlds used to derive PromiseGrid decisions."

The layer names still matter, but they should usually live inside a named simulation until the project intentionally graduates one structure into a reference implementation or dogfood deployment.

## Evidence and prior art inside the repo

- **TE-havib** separates apparatus from specimen. The harness is the apparatus; candidate envelope, transport, group, CAS, and promise-economy designs are specimens under test.
- **TE-vipir** and **TE-sihih** provide useful layer and protocol-family vocabulary, but they still read partly like a target architecture. That is valuable, but it should be instantiated inside simulations when testing concrete consequences.
- **TE-nijab** correctly rejected rewriting historical specimen data merely to satisfy a freeze ceremony. This TE narrows that concern: preserve evidence when it helps decision-making, but do not let historical paths block a clearer experimental layout.
- **TE-domat** correctly noticed the `groups`/`transports` layer conflict but over-corrected toward compatibility preservation. This TE reframes the problem around simulation utility.
- **Turns 170-177** progressively moved from flat `transports/draft--wire-lab-devs/` to protocol grouping, substrate/feed pluralism, rejection of "binding", `groups` vocabulary, sites, CAS, feeds, and L7/L6/L5 layering. That progression is exactly the kind of evolving hypothesis set that benefits from named simulations rather than global root-level commitments.
- **Turns 178-179** added sparse CAS, promise economy, geofencing, conditional release, multi-hop discovery, and promisebase prior art. Those are not just path-naming concerns; they require scenario worlds with actors, sites, promises, chunk possession, policies, observations, and outcomes.

## Assumptions

1. Backward compatibility with current `transports/` paths is not required. Existing paths are not production APIs.
2. Existing experimental artifacts may be moved, renamed, or replaced if doing so makes the lab clearer, provided the change is recorded honestly.
3. Git history is sufficient to recover old layouts if needed. The working tree should optimize for the current best experiment shape, not for preserving every interim path in place.
4. The repo must support side-by-side competing designs. A structure that can host only one current global world is a poor fit.
5. Human and LLM readers should be able to start from a simulation directory and answer: What question is being tested? Which protocols are under test? Who are the actors/sites? What happened? What decision did the evidence support?
6. Production PromiseGrid layout may differ from wire-lab simulation layout. The lab discovers the production layout; it should not prematurely pretend to be it.

## Threat / trust model

Alice is a new human developer. Bob is an LLM agent. Carol is a protocol designer proposing a competing structure. Dave is running a dogfood communication experiment. Ellen is writing the PromiseGrid dev guide. Mallory exploits ambiguity by making a specimen look like settled architecture or by hiding which decision a run was meant to answer.

The structure is successful if Alice, Bob, Carol, Dave, and Ellen can reconstruct experiments and decisions without reading every prior chat turn. It fails if Mallory can point to an arbitrary root-level tree and claim "this is the architecture" when it was only a temporary specimen.

## Alternatives

### Alt A — Root-level production-preview trees

Use root-level trees such as:

```text
groups/
sites/
cas/
feeds/
transports/ or wires/
```

The repo working tree becomes a single current best model of a PromiseGrid-like world.

**Easier:** Very concrete. Humans can browse `groups/` and `sites/` directly. Dogfooding feels real because the repo resembles a live node.

**Harder:** It implies there is one current global world. Side-by-side alternatives need awkward suffixes, branches, or out-of-band notes. It blurs apparatus/specimen boundaries because a root-level `groups/` tree looks canonical even when it is just one candidate.

**Obligations:** Every major structural experiment becomes a root-level migration. The repo history accumulates churn whenever a design candidate changes.

### Alt B — Keep root-level `transports/` as the lab surface

Continue treating `transports/` as the main experimental tree. Put group, feed, CAS, and site evidence underneath it as needed.

**Easier:** Smallest departure from written TE-vipir / TE-sihih / TE-nijab wording. Existing files remain close to their current place.

**Harder:** The name is wrong for the full simulation problem. Sparse CAS, sites, policies, promise economy, geofencing, and group semantics are not all "transports". This structure keeps reintroducing layer confusion.

**Obligations:** Requires repeated explanatory prose telling readers that "transports" means more than transports.

### Alt C — Simulation-first tree with scoped worlds

Create a top-level `simulations/` tree. Each simulation is a named world or run. Its local `world/` directory contains whatever layer trees the candidate model needs.

Example:

```text
simulations/
  README.md
  pahah-wire-lab-devs-bootstrap/
    README.md
    QUESTION.md
    protocol-set.md
    world/
      sites/
        alice/
        bob/
        carol/
      groups/
        group-session/
          wire-lab-devs/
      cas/
        git-cas/
      feeds/
        git-feed/
      wires/
        git/
    events/
      0001-alice-posts-message.md
      0002-bob-fetches-chunk.md
    observations/
      chunk-propagation.md
      confusion-log.md
    results/
      summary.md
    decisions.md
```

**Easier:** Names the repo's true job. Supports competing structures side by side. Makes scenario setup and decision evidence first-class. Prevents premature root-level production claims. Lets one simulation use `groups/` while another tests a different layout.

**Harder:** Adds a level of nesting. A dogfood user must know which simulation is current. Some protocols may need helper docs to map from a protocol spec to the simulations that exercise it.

**Obligations:** Needs a short `simulations/README.md` explaining the contract of a simulation directory. Needs a convention for simulation handles and for connecting results to DR/DI records.

### Alt D — Protocol-first simulations under `protocols/<slug>.d/`

Put simulations under the protocol being tested:

```text
protocols/group-session.d/simulations/
protocols/git-feed.d/simulations/
protocols/cas.d/simulations/
```

**Easier:** Keeps each protocol's examples close to its spec. Simple single-protocol tests are easy to find.

**Harder:** Most important wire-lab questions cross protocols. A realistic scenario includes group semantics, CAS chunks, feeds, sites, promise economy, and often multiple candidate envelope formats. Protocol-first placement makes cross-layer scenarios homeless or duplicated.

**Obligations:** Requires a rule for multi-protocol simulations, which will likely recreate a top-level `simulations/` tree anyway.

### Alt E — Decision-first tree under `DR/` or `TODO/`

Put each simulation under the decision artifact it informs:

```text
DR/DR-.../simulations/
protocols/wire-lab.d/TODO/TODO-.../simulations/
```

**Easier:** Strong traceability. Every simulation directly answers a decision request or TODO.

**Harder:** Simulations can outlive one DR, inform multiple decisions, or become reusable fixtures. Deep paths are hard for humans and tooling. Decision artifacts become bulky and stop being readable logs.

**Obligations:** Requires copying or linking when one simulation informs several decisions.

## Scenario analysis

### S1 — Alice wants to understand what wire-lab is testing this week

Alice opens the repo with no session-log context.

Alt A points her at plausible production-like roots, but she cannot tell whether `groups/` is a settled architecture or an experiment. Alt B gives her an even narrower false cue: everything looks like a transport problem. Alt C gives her `simulations/README.md`, then a list of named simulations and questions. Alt D makes her inspect many protocol directories. Alt E makes her inspect DR/TODO records before seeing the world under test.

Alt C is best because it starts with the lab question rather than a layer name.

### S2 — Bob needs to compare two competing directory structures

Bob wants to test `groups/` at root against `simulations/<run>/world/groups/` and against a pure CAS-frontier design.

Alt A requires branch churn or suffixes at root. Alt B forces all alternatives through `transports/`, biasing the result. Alt C allows:

```text
simulations/abcde-root-groups/
simulations/fghij-scoped-world/
simulations/klmno-cas-frontier/
```

with identical actors and scenario scripts. Alt D scatters the comparison across protocol homes. Alt E ties the comparison to one decision record and makes reuse harder.

Alt C best supports side-by-side comparison.

### S3 — Carol proposes a new feed design

Carol proposes `uucp-feed` and wants to see whether it can carry the same group/CAS scenario as `git-feed`.

Alt A makes feed experiments compete with global state. Alt B is closer, but it still frames the whole world as transport data. Alt C lets Carol clone an existing simulation, swap `world/feeds/git-feed/` for `world/feeds/uucp-feed/`, run the same scenario, and compare `observations/`. Alt D works only if the scenario is feed-local; cross-layer behavior is awkward. Alt E works only if one DR owns the whole comparison.

Alt C best supports controlled substitution.

### S4 — Dave wants dogfood messaging today

Dave wants to use wire-lab-based messaging to coordinate with Alice and Bob.

Alt A is the shortest path to a "real" `groups/wire-lab-devs` directory. Alt C adds nesting, but it also makes the dogfood status honest:

```text
simulations/wire-lab-devs-dogfood/world/groups/group-session/wire-lab-devs/
```

The nested path says this is a dogfood simulation, not yet the production PromiseGrid group layout. If dogfooding becomes active enough that the path is painful, a top-level README or script can point to the current dogfood simulation without changing the underlying structure.

Alt A is faster for typing. Alt C is safer for meaning and still usable.

### S5 — Ellen writes the PromiseGrid dev guide

Ellen needs to know which repo artifacts are final-ish evidence versus discarded experiments.

Alt A gives Ellen a tree that looks authoritative but may just be the latest experiment. Alt B hides too much under transport vocabulary. Alt C gives Ellen `results/` and `decisions.md` inside each simulation, with links to DR/DI records. She can distinguish "this was tried" from "this was chosen". Alt D and Alt E require more cross-navigation.

Alt C best serves dev-guide extraction.

### S6 — Mallory exploits ambiguity

Mallory wants to claim that a half-finished root-level tree is the settled PromiseGrid architecture.

Alt A gives Mallory the strongest opening because root-level production-looking names imply authority. Alt B gives a different ambiguity: lower-layer transport framing can be used to hide group/CAS policy choices. Alt C reduces the attack because every concrete world sits under `simulations/<slug>/` and must state its `QUESTION.md`. Mallory can still lie, but the path itself contradicts the claim of finality.

Alt C is most robust against specimen-as-canon confusion.

### S7 — A simulation needs large content, sparse CAS, and geofencing

Alice has chunks X and Y. Bob has Y and Z. Carol is allowed to fetch X only if she promises not to send it outside group G.

Alt A can represent this, but the root grows into a pseudo-production node. Alt B cannot naturally host sites, policy, groups, and CAS without overloading "transport". Alt C naturally scopes the whole world:

```text
world/sites/alice/cas/
world/sites/bob/cas/
world/groups/G/
world/feeds/
events/carol-requests-X.md
observations/geofence-assessment.md
```

Alt C is best for full-world simulations.

### S8 — We later decide the production layout should be different

Suppose wire-lab simulations show that production PromiseGrid should use a layout unlike any current repo tree.

Alt A creates expensive embarrassment because the repo root looked like the production preview. Alt B similarly binds us to transport vocabulary. Alt C has no problem: simulations are evidence, not production. The final layout can be documented in the dev guide or reference implementation, citing the simulations that led there.

Alt C best protects learning velocity.

## Conclusions

1. **Reject Alt B as the main structure.** `transports/` is too narrow for the simulation and decision-making goals.
2. **Reject Alt E as the main structure.** DR/TODO traceability is necessary, but simulations should not be buried inside decision logs.
3. **Reject Alt D as the main structure.** Protocol-local simulations are useful for unit examples, but realistic wire-lab questions are cross-protocol and cross-layer.
4. **Keep Alt A only as a possible later production-preview or dogfood shortcut.** Root-level `groups/`, `sites/`, `cas/`, and `feeds/` may become appropriate after the lab chooses a stable current-world model. They are premature as the primary structure.
5. **Choose Alt C as the best wire-lab structure now.** A top-level `simulations/` tree with scoped worlds best supports the lab's real job: comparing candidate PromiseGrid structures, running scenario evidence, and driving DR/DI decisions.

## Recommended structure

Top-level:

```text
protocols/                  # Candidate protocol/spec families.
docs/thought-experiments/   # TE analyses.
DR/                         # Open decision requests.
protocols/*/TODO/           # Task queues and DI logs.
implementations/            # Reference/prototype code when it exists.
simulations/                # Concrete experimental worlds and runs.
```

Simulation directory:

```text
simulations/<sim-handle>-<slug>/
  README.md                 # Human orientation: purpose, status, how to inspect.
  QUESTION.md               # Decision under test and alternatives.
  protocol-set.md           # Protocol specs/pCIDs/draft paths used in this simulation.
  world/
    sites/                  # Simulated nodes and local views.
    groups/                 # L7 group state, rosters, frontiers, semantic pointers.
    cas/                    # Shared or modeled CAS objects, if the sim uses a shared fixture.
    feeds/                  # L5 feed configs/claims/protocol-specific fixtures.
    wires/                  # Raw carrier events, packet traces, or substrate observations.
  events/                   # Ordered scenario events; human-readable first.
  observations/             # What happened, including failures/confusion.
  results/                  # Summaries suitable for TE/DR/DI citation.
  decisions.md              # Links to DRs/DIs/TODOs affected by the result.
```

The `world/` subtrees are not necessarily production PromiseGrid paths. They are the simulation's model of the world. If a simulation is specifically testing a proposed production filesystem layout, `QUESTION.md` must say so.

## Implications for `transports/` and `groups/`

Under this TE's recommendation:

- Do not choose between top-level `transports/` and top-level `groups/` as the primary experiment home.
- Move future concrete dogfood/simulation worlds under `simulations/<sim>/world/`.
- Use `world/groups/` when modeling L7 group semantics.
- Use `world/feeds/` for L5 feed behavior and `world/wires/` or `world/transports/` only for raw substrate/transport evidence.
- Existing `transports/wire-lab-devs-draft/` can be migrated into a simulation if that makes the repo clearer. Since no production compatibility is required, the question is not "can we preserve the old path?" but "what migration leaves future readers with the clearest experiment record?"

## DF questions exposed

### DF-pahah.1 — Adopt `simulations/` as the primary home for concrete worlds?

Recommended answer: yes, Alt C.

Surviving alternatives:

- **1.A — Yes, simulation-first (recommended).** Create `simulations/` and put concrete experimental worlds under `simulations/<sim-handle>-<slug>/`.
- **1.B — No, root-level production-preview.** Use top-level `groups/`, `sites/`, `cas/`, and `feeds/` as the current world.
- **1.C — Hybrid.** Use `simulations/` for comparisons, but keep one top-level current dogfood world as a shortcut.

### DF-pahah.2 — What should happen to current `transports/wire-lab-devs-draft/`?

Recommended answer: migrate it into the first dogfood simulation when implementation begins.

Surviving alternatives:

- **2.A — Migrate into `simulations/<sim>/world/` (recommended).** Treat current data as seed evidence for the new simulation structure.
- **2.B — Leave it temporarily and mark it legacy.** Add README pointers only.
- **2.C — Delete/recreate if cleaner.** Accept that current data is disposable brainstorming output; preserve only git history and any still-useful CIDs in a note.

### DF-pahah.3 — Should simulation IDs use the proquint handle namespace?

Recommended answer: yes, if simulations become citable decision evidence.

Surviving alternatives:

- **3.A — `SIM-<handle>` IDs.** Mint handles with the same global handle discipline when a simulation needs stable citations.
- **3.B — Slug-only simulation directories.** Use human names until collisions or citations force stronger IDs.

## Implications for open work

- **TE-domat** should be treated as useful but too compatibility-weighted. Its key question should be re-read through this TE: root-level `groups/` versus `transports/` is less important than where simulations live.
- **DR-nugog** should be reframed again. The best next DR question is whether to adopt `simulations/` and migrate the current bootstrap data into the first simulation world, not whether `transports/` should be flat or nested.
- **TODO-kugod / UT-159.b** should remain open until the transport-spec companion audit distinguishes apparatus-level simulation structure from candidate production layout.
- **TODO-turog / TODO-pipus / TODO-duvuk** should not plan migrations directly from old root paths to new root paths. They should first decide whether the target is a simulation world, a protocol spec example, or a future production/reference layout.
- **DEV-GUIDE-RESOURCES.md** should eventually point guide writers at `simulations/*/results/` once such directories exist, because those results will be more useful than raw path history.

## Decision status

`needs DF`. This TE answers the design question by recommending a simulation-first structure, but it does not lock the decision. The next decision is DF-pahah.1: adopt `simulations/` as the primary home for concrete wire-lab worlds, or keep a root-level current-world layout.
