# TE-mupoz: Root and protocol migration scope under simulations

*Thought experiment: if TE-nizor's Pahah-with-contract recommendation is
implemented, how much existing root-level content and `protocols/` content should
move under `simulations/`?*

## TE ID

TE-mupoz

## Status

decided

First drafted 2026-05-10 00:22:21 UTC.

## Minting note

`tools/mint-handle` was attempted first, but the current corpus scan aborts on
the pre-existing duplicate owner pair
`TODO-jivam-turns-149-170-recovery-completion.md` and `DI-jivam`. The fallback
for this TE used the same proquint alphabet and an explicit repository collision
check for `mupoz`; no `TE-mupoz`, `TODO-mupoz`, `DR-mupoz`, `DI-mupoz`, or bare
`mupoz` occurrence existed before this file was created. The duplicate
`jivam`/`DI-jivam` condition should be resolved separately so future handle
minting returns to the normal tool path.

## Decision under test

DUT-mupoz: After TE-nizor, should existing root-level content and existing
`protocols/` content be moved under `simulations/`, copied under
`simulations/`, referenced from `simulations/`, or left entirely outside
`simulations/`?

This is not a cosmetic file-layout question. The answer determines whether
`simulations/` becomes:

- a bounded experiment apparatus that consumes existing protocol drafts and
  specimens as inputs;
- a new home for all wire-lab activity;
- or a confusing duplicate source of truth.

TE-nizor says Pahah can satisfy the recovered turns 149-208 concerns only if
`simulations/` carries a minimal contract: question, concerns matrix, protocol
set, scoped world, events, observations, results, and decision handoff. This TE
tests what that implies for the already-existing corpus.

## Short answer

Move **specimens into simulations** and leave root `protocols/` for the
wire-lab apparatus itself.

Implementing TE-nizor should not mean hoisting the current repo into
`simulations/`. It should mean:

1. Keep top-level apparatus where it is:
   `AGENTS.md`, `README.md`, `DEV-GUIDE-RESOURCES.md`, `docs/`, `DR/`,
   `protocols/wire-lab.d/`, `tools/`, `implementations/`, and most
   `proposals/`.
2. Treat non-wire-lab protocol drafts such as `group-session`, `udp-binding`,
   and `ppx-dr` as candidate PromiseGrid protocol specimens until they
   graduate. They should move under the simulation that is testing them rather
   than remain as root-level siblings of the wire-lab harness protocol.
3. Treat concrete message/state evidence as the class that naturally belongs
   inside a simulation world. The main current example is
   `transports/wire-lab-devs-draft/`.
4. Do not preserve `transports/wire-lab-devs-draft/` in its current location as
   a design constraint. Preserve auditability instead: use `git mv` when moving
   it, record the old path, source commit, and CIDs, and keep any historically
   significant path strings honest as historical evidence.
5. Put **new** concrete world state under `simulations/<sim>/world/` by default
   after Pahah is implemented.

The recommended policy is therefore:

> Existing apparatus stays rooted. Root `protocols/` keeps `wire-lab.d` as the
> harness protocol. Candidate PromiseGrid protocol drafts and concrete
> transport/world specimens move into the first relevant simulation with
> source-path, source-commit, and CID/pCID manifests. New simulation specimens
> and mutable world state live under `simulations/<sim>/`.

## Sources examined

- **TE-pahah** proposes `simulations/` as the primary home for concrete
  experimental worlds, replacing premature root-level `transports/`/`groups/`
  debates with named simulation worlds.
- **TE-nizor** recommends Pahah-with-contract and explicitly leaves protocol,
  artifact, representation, CAS, feed, and promise-economy specimen decisions
  unsettled.
- **TE-vilot** and **TE-hirap** keep apparatus documents ordinary while allowing
  promise-shaped or PromiseGrid-message-shaped artifacts inside named
  simulations.
- **TE-vipir** previously locked `protocols/<slug>.d/` as simulated repos for
  protocol design, with per-protocol specs, TODO queues, manifests, and frozen
  siblings. This TE reopens whether those simulated repos belong at repo root or
  under named simulations.
- **TE-liviv** separates spec-side protocol design from implementation-side
  artifacts under `implementations/`; this TE keeps the split but questions
  whether candidate spec-side trees should be root-level in wire-lab.
- **TE-domat** recommends treating root-level `groups/`/`transports/` as a
  layer-ownership question and warns against silently rewriting existing
  `transports/wire-lab-devs-draft/` history.
- **Current repo layout** has top-level apparatus (`AGENTS.md`, `README.md`,
  `DEV-GUIDE-RESOURCES.md`, `docs/`, `DR/`, `protocols/wire-lab.d/`), candidate
  protocol trees (`protocols/group-session.d/`, `protocols/udp-binding.d/`,
  `protocols/ppx-dr.d/`), tool and implementation trees (`tools/`,
  `implementations/`), review/proposal records (`proposals/`), and one concrete
  transport specimen (`transports/wire-lab-devs-draft/`).

## Terms

- **Move:** a `git mv` or equivalent path change that makes a file's current
  path live under `simulations/`.
- **Copy:** duplicate bytes under `simulations/` while the original remains.
- **Reference:** name the original path, commit, CID, pCID, or spec path from a
  simulation file without duplicating the bytes.
- **Seed:** input evidence used to initialize or motivate a simulation.
- **Specimen:** message bytes, world state, runtime traces, or other objects
  being tested.
- **Apparatus:** documents and tools used to reason about, govern, run, or
  interpret the tests.

## Current content classification

### Root identity and governance apparatus

Paths:

- `AGENTS.md`
- `AGENTS-codex.md`
- `AGENTS-ppx.md`
- `README.md`
- `LICENSE`
- `DEV-GUIDE-RESOURCES.md`

These files define repo identity, agent rules, guide-writer orientation, and
legal terms. They are not simulation specimens. A simulation may cite them when
explaining its constraints, but moving them under `simulations/` would make the
repo look like it has no top-level operating contract.

Recommended action: **keep rooted; reference only**.

### Coordination and decision apparatus

Paths:

- `docs/thought-experiments/`
- `docs/discussion/`
- `docs/research/`
- `DR/`
- `protocols/wire-lab.d/TODO/TODO.md`

These are the reasoning and decision trail. TE-vilot and TE-hirap both warn that
turning apparatus into message-shaped specimens too early blurs the boundary
Mallory can exploit. TE-nizor uses these files to decide what simulations should
test; it does not say the reasoning trail itself becomes a simulation world.

Recommended action: **keep rooted; simulations cite relevant TE/DR/TODO IDs and
paths; do not move wholesale**.

### Wire-lab apparatus protocol tree

Paths:

- `protocols/wire-lab.d/`

This is the protocol/spec tree for the wire-lab harness itself: the apparatus
that defines how this repo runs experiments, tracks TODOs, records harness
decisions, and produces evidence. It is not a candidate PromiseGrid wire
protocol specimen.

Recommended action: **keep rooted as the only default root `protocols/` member**.

### Candidate PromiseGrid protocol trees

Paths:

- `protocols/group-session.d/`
- `protocols/udp-binding.d/`
- `protocols/ppx-dr.d/`

These are not the wire-lab apparatus itself. They are candidate PromiseGrid
protocols being tested by the lab: session behavior, feed/binding behavior, and
proposal/decision-message behavior. That makes them closer to specimens than to
root governance.

Keeping them as root-level siblings of `wire-lab.d` makes a candidate protocol
look as authoritative as the harness protocol. It also makes `protocols/` do two
jobs at once: "the protocol of this lab" and "protocols currently under test."
Under Pahah, the second job belongs inside named simulations.

Recommended action: **move under the first simulation that uses them, with
source-path/source-commit manifests. Within a simulation, place them under
`simulations/<sim>/protocols/` when they are protocol specimens, and reference
them from `protocol-set.md`. Root `protocols/` should then contain only
`wire-lab.d` unless a later DF creates another apparatus-level protocol**.

### Implementation and tool apparatus

Paths:

- `tools/`
- `implementations/`
- `bin/`

Tools and implementations may be used by simulations, but they are not
simulation state. A simulation can record tool versions and command output in
`observations/` or `results/`; it should not own the tool's source unless the
simulation is explicitly testing a fork of that tool.

Recommended action: **keep rooted; record tool versions and outputs inside
simulation observations when relevant**.

### Proposal and review records

Paths:

- `proposals/pending/`
- `proposals/approved/`

These records are ambiguous. They are apparatus when they document human review
and branch convergence. They can also be specimens if the experiment is "can
ppx-dr proposals travel as PromiseGrid messages?" Existing files should not be
moved by default because they are historical review records with path-based
references and DI/DR provenance.

Recommended action: **keep rooted by default; copy or wrap selected records only
inside simulations that explicitly test proposal/message protocols**.

### Concrete transport/message specimen

Paths:

- `transports/README.md`
- `transports/wire-lab-devs-draft/README.md`
- `transports/wire-lab-devs-draft/*.txt`

This is the one current class that clearly resembles simulation world state:
message bytes, transport instance metadata, participant roster, and bootstrap
round-trip evidence. It is also historical and already path-sensitive:
`transports/wire-lab-devs-draft/README.md` records previous names, bootstrap
state, and freeze gates. Moving it silently would destroy evidence about the
path confusion the current TEs are trying to resolve.

Recommended action: **migrate the existing transport specimen into the first
Pahah simulation if DF locks Pahah/Mupoz, using `git mv` plus a source-path,
source-commit, and CID manifest. Do not preserve the current root
`transports/` location merely for compatibility; preserve the evidence and its
history instead**.

## Alternatives

### Alt A -- Whole-repo simulation move

Move most root-level content into a first simulation, for example:

```
simulations/wire-lab-recovery/
    AGENTS.md
    README.md
    docs/
    DR/
    protocols/
    proposals/
    transports/
    tools/
```

**Easier:** Everything relevant to the recovery is physically co-located. A
developer opening the simulation sees the entire corpus.

**Harder:** The simulation becomes the repo. Top-level identity disappears or
gets duplicated. Apparatus and specimen collapse completely. Future simulations
cannot share the harness rules or candidate protocols without copying again.
This also makes `simulations/` look like an archival bundle instead of an
experimental world.

**Obligations:** Requires broad path rewrites, Cat-1a/Cat-1b classification
across the corpus, updated AGENTS path scopes, broken relative-link repair, and
a new theory of repo identity.

### Alt B -- Move all protocol trees under `simulations/`, keep other apparatus rooted

Keep root docs and tools, but make protocol drafts simulation-local:

```
simulations/wire-lab-recovery/protocols/
    wire-lab.d/
    group-session.d/
    udp-binding.d/
    ppx-dr.d/
```

**Easier:** The simulation has full copies of the protocol drafts it uses. A
replay can be hermetic if those copies never change.

**Harder:** It moves the wire-lab harness protocol itself under a simulation,
which makes the lab apparatus look like a specimen. Multiple simulations then
fork the harness rules by file copy.

**Obligations:** Requires a new rule for identifying the current wire-lab
apparatus after the apparatus has moved into one simulation. This is the same
apparatus/specimen confusion TE-nizor is trying to avoid.

### Alt C -- Move only concrete root-level specimens

Keep apparatus and all current `protocols/` directories rooted. Move existing
concrete world state under the first simulation:

```
simulations/wire-lab-recovery/world/transports/wire-lab-devs-draft/
```

**Easier:** `simulations/` becomes the obvious place for active and historical
world state. Root-level `transports/` stops being a competing active tree.

**Harder:** A literal `git mv` of existing message specimen paths can make
historical references false or misleading. Existing README prose says this
transport lived under `transports/`; moving it may force a large documentation
repair before the simulation even starts. It also risks implying that the old
path never mattered.

**Obligations:** Requires a migration/supersession note at the old path, careful
CID verification after move, and a statement that message content CIDs survive
path migration but path references in surrounding docs remain historical.

### Alt D -- Reference-only simulations

Create `simulations/`, but do not move or copy any existing files. The first
simulation has `protocol-set.md`, `concerns.md`, and `observations/` that point
to existing paths:

```
simulations/wire-lab-recovery/
    protocol-set.md
    concerns.md
    observations/existing-transport.md
```

**Easier:** No duplicate bytes and no path churn. The simulation remains a clean
analysis layer over the current corpus.

**Harder:** It may not feel like a concrete world. New actors still need to jump
between root `transports/`, `protocols/`, and simulation files. If the simulation
needs to mutate or branch world state, reference-only input is insufficient.

**Obligations:** Requires strong source references: commit IDs, paths, CIDs, and
clear rules for when referenced evidence becomes copied seed data.

### Alt E -- Wire-lab-only root protocols plus specimen migration

Keep apparatus and `protocols/wire-lab.d/` rooted. Move candidate PromiseGrid
protocol trees and concrete transport/world specimens into the first simulation
when they are part of the recovery/dogfood world. The first simulation records
each migration as seed evidence and makes protocol specimens explicit:

```
simulations/wire-lab-recovery/
    README.md
    QUESTION.md
    concerns.md
    protocol-set.md
    seed/
        README.md
        protocol-tree-migrations.md
        wire-lab-devs-draft-migration.md
    protocols/
        group-session.d/
        udp-binding.d/
        ppx-dr.d/
    world/
        sites/
        groups/
        cas/
        feeds/
        wires/
        transports/
    events/
    observations/
    results/
    decisions.md
```

Under this policy:

- `seed/protocol-tree-migrations.md` records original paths, source commits,
  destination paths, and draft/frozen pCID status for candidate protocol trees.
- `seed/wire-lab-devs-draft-migration.md` records original path, source commit,
  message CIDs, and the `git mv` destination.
- `protocol-set.md` references simulation-local protocol specimens by path and
  commit, and later by pCID when frozen.
- migrated and new world state live under `world/`;
- old root `transports/` does not remain active; it may become a temporary
  redirect note during migration or disappear once references are updated.
- root `protocols/` retains `wire-lab.d` as the harness apparatus protocol.

**Easier:** Preserves evidence, avoids treating candidate PromiseGrid protocols
as root apparatus, and makes the simulation world concrete. It also handles the
real recovery need: make the concerns visible and testable without pretending
the corpus was always arranged this way.

**Harder:** The migration must be explicit. Historical path references need to
remain historically true, while current pointers must move to the simulation.
The first simulation README must distinguish harness apparatus, protocol
specimens, migrated seed evidence, and new world state.

**Obligations:** Define migration manifest formats for protocol trees and
transport specimens, verify CIDs before and after the move, classify
path-reference updates as Cat-1a/Cat-1b, and define a graduation rule for moving
results into DR/DI/spec work or an external PromiseGrid spec corpus.

### Alt F -- Future-first only

Keep every existing file where it is. Declare that `simulations/` applies only
to new work after the DF.

**Easier:** Zero migration risk. No historical path repair.

**Harder:** It fails the recovery purpose of TE-nizor. Turns 149-208 concerns
would remain scattered across TEs, TODOs, DRs, root transport files, and session
logs, while the new simulation contract starts empty. It would let the same
context-loss gap persist.

**Obligations:** Requires a separate recovery index anyway, which would likely
recreate `seed/` and `concerns.md` under another name.

## Scenario analysis

### S1 -- Alice writes the PromiseGrid dev guide

Alice needs to explain PromiseGrid, not wire-lab. She needs stable design
outputs, not every experimental byway.

- **Alt A:** Alice sees the whole repo under one simulation and may mistake one
  experiment bundle for the PromiseGrid design.
- **Alt B:** Alice may cite simulation-local `wire-lab.d` copies as if they are
  the harness rules.
- **Alt C:** Alice sees simulation-local world evidence, but may lose the
  historical reason old transport paths existed.
- **Alt D:** Alice sees clean references but may still have to chase too many
  scattered paths.
- **Alt E:** Alice gets the best source map: root `wire-lab.d` for harness
  method, simulation-local protocol specimens for evidence, and simulation
  `results/` for what the experiment showed.
- **Alt F:** Alice gets no useful synthesis from the simulation.

S1 favors Alt E. The dev guide needs simulations to point at evidence, not
become the source of protocol truth.

### S2 -- Bob implements the first Pahah simulation

Bob needs to create a working directory that can answer "does this structure
recover turns 149-208 and unblock dogfood?"

- **Alt A:** Bob spends most of his time repairing paths.
- **Alt B:** Bob becomes responsible for maintaining duplicate harness and
  protocol trees.
- **Alt C:** Bob can start from concrete message bytes but must handle path
  history carefully.
- **Alt D:** Bob can write the contract quickly, but lacks local world bytes if
  he wants replay.
- **Alt E:** Bob can migrate candidate protocol trees into `protocols/`, migrate
  the concrete transport specimen into `world/`, keep the harness rooted, and
  start replay from migration manifests.
- **Alt F:** Bob has no existing evidence inside the simulation, so the result
  does not prove recovery.

S2 favors Alt E, with Alt D as an acceptable first sub-step if time is tight.

### S3 -- Carol maintains `group-session.d`

Carol owns candidate protocol design and freeze work. She needs to know whether
she is editing wire-lab apparatus or a protocol specimen under test.

- **Alt A:** Carol's protocol source moves with the whole repo and loses clear
  top-level discoverability.
- **Alt B:** Carol must reconcile simulation-local `wire-lab.d` copies with the
  root harness.
- **Alt C:** Carol keeps `protocols/group-session.d/` authoritative, but that
  makes the candidate protocol look like root apparatus.
- **Alt D:** Same as Alt C.
- **Alt E:** Carol edits candidate protocol specimens inside the simulation that
  tests them; if they graduate, their result points to a DR/DI/spec handoff.
- **Alt F:** Same as Alt C, but simulations provide less feedback.

S3 rejects Alt B and weakens Alt C. The root should not mix harness apparatus
with candidate PromiseGrid protocols. Candidate protocols can still be shared by
copying, referencing, or freezing simulation-local specimens by pCID.

### S4 -- Dave replays the context-loss recovery window

Dave wants to verify that turns 149-208 concerns did not disappear.

- **Alt A:** Dave can find everything, but cannot tell which files are apparatus
  and which are specimens.
- **Alt B:** Dave may treat simulation-local harness copies as evidence that a
  harness decision was locked when it was only simulation-local.
- **Alt C:** Dave sees concrete messages under the simulation if they were moved,
  but needs a migration audit to trust the move.
- **Alt D:** Dave can trace concerns through references, but still has to read
  root files manually.
- **Alt E:** Dave sees each recovered concern mapped to migrated protocol
  specimens, migrated transport evidence, world state, observations, and result
  status.
- **Alt F:** Dave cannot use the simulation to complete recovery.

S4 favors Alt E. The critical object is not a moved directory; it is the concern
matrix tied to seed evidence and result criteria.

### S5 -- Ellen opens the repo in 2030 after another context gap

Ellen has no chat memory. She trusts file layout and explicit source maps.

- **Alt A:** Ellen thinks `simulations/wire-lab-recovery/` is the real repo.
- **Alt B:** Ellen finds multiple harness copies and cannot tell which one won.
- **Alt C:** Ellen sees a cleaner active world tree but may not understand
  historical root paths.
- **Alt D:** Ellen sees no duplication, but only if the references are complete.
- **Alt E:** Ellen sees root wire-lab apparatus, simulation-local protocol
  specimens, migration manifests, and results that point to DR/DI/spec outcomes.
- **Alt F:** Ellen rediscovers the same scattered state.

S5 strongly favors Alt E.

### S6 -- Frank runs two simulations with different protocol variants

Frank wants one simulation using current `group-session-draft.md` and another
using a forked envelope candidate.

- **Alt A:** Both simulations share the whole repo bundle awkwardly.
- **Alt B:** Each simulation may fork both harness and candidate protocols, so
  the comparison point is unclear.
- **Alt C:** Frank can reference the same top-level candidate protocol draft from
  both and record variant deltas separately, but the draft still looks
  apparatus-level.
- **Alt D:** Same as Alt C, but with no local snapshots.
- **Alt E:** Frank compares protocol specimens inside each simulation and names
  the exact local candidate or pCID in `protocol-set.md`.
- **Alt F:** Future simulations can work, but old recovery evidence is missing.

S6 favors Alt E. Simulations should test protocol variants as specimens without
turning every candidate into a root protocol.

### S7 -- Mallory exploits duplicate authority

Mallory wants a weak reader to cite the wrong file as authoritative.

- **Alt A:** Mallory points at a simulation copy of `AGENTS.md` or a protocol
  spec and claims it is repo policy.
- **Alt B:** Mallory points at a simulation-local harness draft and claims it is
  the canonical one.
- **Alt C:** Mallory has less room, but can still exploit moved historical
  transport paths if migration notes are weak.
- **Alt D:** Mallory has little duplicate authority to exploit, but can exploit
  missing references.
- **Alt E:** Mallory has the least leverage if migration manifests clearly say
  "protocol specimen, not harness apparatus" and `protocol-set.md` identifies
  exactly which candidate each simulation tests.
- **Alt F:** Mallory can exploit scattered unintegrated context.

S7 favors Alt E and requires strict labels: root `wire-lab.d`, simulation
`protocols/`, `seed`, `world`, `observations`, `results`, and `protocol-set`
are different artifact classes.

### S8 -- A write fails halfway through migration

An agent starts moving root files into `simulations/` and crashes.

- **Alt A:** Catastrophic partial migration risk.
- **Alt B:** Duplicate harness/protocol sources may be half-created.
- **Alt C:** Risk is limited to concrete specimens, but still affects historical
  path integrity.
- **Alt D:** No move risk.
- **Alt E:** The move is limited to candidate protocol specimens and concrete
  transport/world specimens. It can be staged: prepare the simulation contract
  and migration manifests first, verify CIDs/pCIDs, then `git mv` specimen
  trees. The root harness apparatus never enters the partial-move blast radius.
- **Alt F:** No migration risk but also no recovery progress.

S8 favors Alt E, implemented in phases: contract and manifest first, CID checks
second, physical `git mv` third.

### S9 -- Scale increases

Hundreds of simulations and many protocol drafts exist.

- **Alt A:** Each simulation is too large if it bundles the repo.
- **Alt B:** Harness and protocol copies multiply and drift.
- **Alt C:** Simulation worlds remain manageable, but active and historical
  specimen placement must be clear.
- **Alt D:** Smallest storage footprint but can become reference-heavy.
- **Alt E:** Good storage properties if migrated protocol specimens can be
  content-addressed and shared by pCID when needed.
- **Alt F:** Future scale is okay, but legacy evidence remains unsynthesized.

S9 favors Alt E with optional content-addressed snapshots only when the
simulation must be reproducible independent of future draft edits.

## Evaluation matrix

| Content class | Existing paths | Move under `simulations/`? | Recommended treatment |
|---|---|---:|---|
| Repo identity/rules | `AGENTS*.md`, `README.md`, `LICENSE` | No | Keep rooted; simulations cite constraints if relevant. |
| Guide-writer source map | `DEV-GUIDE-RESOURCES.md` | No | Keep rooted; add pointers to simulation results when they exist. |
| TEs, discussion, research | `docs/` | No | Keep rooted as apparatus; simulations cite IDs and paths. |
| DRs | `DR/` | No | Keep rooted; simulations cite open/settled decisions. |
| Harness protocol/TODO/DI logs | `protocols/wire-lab.d/` | No | Keep rooted as the lab apparatus protocol. |
| Candidate protocol drafts/TODOs | `protocols/group-session.d/`, `protocols/udp-binding.d/`, `protocols/ppx-dr.d/` | Yes, after DF | Move into the first simulation as protocol specimens with source-path/source-commit/pCID migration manifests. |
| Candidate protocol manifests/changelogs | non-wire-lab `protocols/*/manifest.json`, `CHANGELOG.md` | Yes, with their protocol tree | Move with the specimen tree; mark specimen status and graduation state. |
| Tools | `tools/`, `bin/` | No | Keep rooted; record tool versions/results in observations. |
| Implementations | `implementations/` | No | Keep rooted; simulations may invoke or compare implementations. |
| Proposal/review records | `proposals/` | Usually no | Keep rooted; copy/wrap only for simulations testing proposal-as-message protocols. |
| Existing transport specimen | `transports/wire-lab-devs-draft/` | Yes, after DF | Move into the first simulation with a source-path/source-commit/CID migration manifest. |
| Future concrete world state | none yet | Yes | Create under `simulations/<sim>/world/`. |
| Future simulation results | none yet | Yes | Create under `simulations/<sim>/results/` and point to DR/DI/spec updates. |

## Consequences for `protocols/`

Root `protocols/` should probably contain only `wire-lab.d`.

The stronger statement is:

> A simulation consumes and tests candidate PromiseGrid protocols as specimens.
> The root `protocols/` tree names the wire-lab harness protocol unless a later
> DF creates another apparatus-level protocol.

That means the first simulation's `protocol-set.md` should name migrated
protocol inputs in a way that is stable enough for replay:

```
Protocol: group-session
Role in simulation: session protocol under test
Original source path: protocols/group-session.d/specs/group-session-draft.md
Simulation path: simulations/<first-sim>/protocols/group-session.d/specs/group-session-draft.md
Source commit: <git commit>
Spec pCID: draft, not frozen
Local variant: none
Notes: used as current draft; simulation results may propose DR/DI/spec edits
```

If a simulation tests a fork, the fork should be visibly local:

```
Protocol: group-session
Role in simulation: session protocol candidate B
Base specimen path: simulations/<first-sim>/protocols/group-session.d/specs/group-session-draft.md
Source commit: <git commit>
Local variant file: protocol-variants/group-session-alt-frontier-ack.md
Spec pCID: specimen only, not canonical
```

Only after simulation results survive review should changes graduate out of the
simulation. The destination is not necessarily root `protocols/<slug>.d/`; it may
be a future PromiseGrid spec corpus, the PromiseGrid dev guide, or a frozen
specimen/result under the simulation. Root `protocols/` should not regain
non-wire-lab candidate protocols unless a later DF says wire-lab itself should
host a broader spec registry.

## Consequences for root-level `transports/`

Root-level `transports/` should not be preserved in its current location merely
because it exists today. The evidence matters; the root path does not.

The next implementation should migrate `transports/wire-lab-devs-draft/` into
the first simulation after the DF locks the target simulation path. The
migration should create an audit record such as:

```
simulations/<first-sim>/
    seed/
        wire-lab-devs-draft-migration.md
    world/
        transports/
            wire-lab-devs-draft/
```

The migration manifest should list:

- original path;
- source commit;
- message filenames/CIDs;
- verification command and result;
- destination path;
- whether any root `transports/` redirect note remains temporarily;
- historical caveat that old path names are evidence, not current design
  commitments.

Future message/world state should stay under the simulation:

```
simulations/<first-sim>/world/
    sites/
    groups/
    cas/
    feeds/
    wires/
    transports/
```

After that, root-level `transports/` should either disappear or remain only as a
short redirect/deprecation note for current readers. It should not remain as an
active specimen location.

## What this TE rejects

- Moving `protocols/wire-lab.d/` under `simulations/`.
- Moving `docs/`, `DR/`, or TODO/DI apparatus wholesale under `simulations/`.
- Treating a simulation-local candidate protocol as harness apparatus or as a
  graduated PromiseGrid spec before DR/DI/spec handoff.
- Silently `git mv`-ing `transports/wire-lab-devs-draft/` without a migration
  manifest, source-path note, and DF.
- Starting Pahah with future-only empty simulations that do not map recovered
  turns 149-208 concerns into seed evidence and result criteria.

## What this TE preserves

- TE-vipir's insight that protocol work benefits from repo-like trees, while
  relocating candidate protocol trees under simulations until graduation.
- TE-liviv's spec-side versus implementation-side split.
- TE-vilot and TE-hirap's apparatus/specimen boundary.
- TE-nizor's need for a concrete simulation contract.
- TE-domat's caution that current transport evidence is historical specimen data
  and should not be silently rewritten. This preserves the evidence, not the old
  root path.

## DF questions exposed

### DF-mupoz.1 -- What is the default migration scope when Pahah is implemented?

Locked answer: 1.A per `DI-fakin`.

Surviving alternatives:

- **1.A -- Wire-lab-only root protocols plus specimen migration.** Keep
  apparatus and `protocols/wire-lab.d/` rooted; move candidate protocol trees
  and concrete transport/world specimens into the first simulation with
  migration notes; write new world state under simulations. **Locked by
  `DI-fakin` on 2026-05-10.**
- **1.B -- Move only concrete specimens.** Physically move existing root
  transport specimens into the first simulation, but without the broader
  `protocol-set.md` and graduation rules.
- **1.C -- Reference-only first.** Create simulations that reference all existing
  evidence but copy/move nothing until a later implementation pass.
- **1.D -- Move all `protocols/` under simulations.** Rejected because it moves
  `wire-lab.d` apparatus under a specimen world.

### DF-mupoz.2 -- What happens to `transports/wire-lab-devs-draft/`?

Locked answer: 2.A per `DI-fakin`.

Surviving alternatives:

- **2.A -- Migrate with source manifest.** `git mv` the directory into the
  first simulation, create a seed/migration manifest, verify message CIDs, and
  update current pointers while leaving historical references honest. **Locked
  by `DI-fakin` on 2026-05-10.**
- **2.B -- Read-only seed copy.** Keep the original and copy message bytes into
  the first simulation for replay, with CID verification and source-path notes.
- **2.C -- Preserve in place.** Keep `transports/wire-lab-devs-draft/` rooted
  and only reference it from simulations. Rejected as over-preserving the old
  path now that compatibility is not a constraint.

### DF-mupoz.3 -- What remains in root `protocols/`?

Locked answer: 3.A per `DI-pakid`.

Surviving alternatives:

- **3.A -- Root `protocols/` contains only `wire-lab.d` (recommended).**
  Candidate PromiseGrid protocols move under simulations as specimens. **Locked
  by `DI-pakid` on 2026-05-10.**
- **3.B -- Root `protocols/` remains a lab-wide candidate protocol registry.**
  Simulations reference root protocol drafts. Rejected as likely to confuse
  apparatus with specimens.
- **3.C -- All protocol trees move under simulations, including `wire-lab.d`.**
  Rejected because the harness apparatus needs a root home.

### DF-mupoz.4 -- What happens to `proposals/`?

Locked answer: 4.C per `DI-fakin`.

Surviving alternatives:

- **4.A -- Keep proposal records rooted.** Rejected after Steve clarified that
  the legacy proposal mechanism has been replaced by `DEV-GUIDE-RESOURCES.md`
  and the PromiseGrid dev-guide feedback process.
- **4.B -- Copy selected proposals into a ppx-dr simulation.** Acceptable for a
  proposal-as-message test.
- **4.C -- Move all proposals under simulations.** Move the legacy proposal
  records into the first simulation archive as historical evidence, not active
  world state. **Locked by `DI-fakin` on 2026-05-10.**

### DF-mupoz.5 -- What is the graduation rule?

Locked answer: 5.A per `DI-fakin`.

Surviving alternatives:

- **5.A -- Results feed DR/DI/spec handoff.** Simulation results point to DRs,
  DIs, frozen specs, dev-guide prose, or a future PromiseGrid spec corpus;
  rooted apparatus changes only through the normal decision process. **Locked by
  `DI-fakin` on 2026-05-10.**
- **5.B -- Results can directly become root protocols.** Rejected because it lets
  simulation-local evidence bypass decision provenance and recreates the
  apparatus/specimen ambiguity.
- **5.C -- Results are illustrative only.** Rejected because simulations need to
  produce actionable design evidence.

## Locked implementation sequence

1. Create `simulations/README.md` explaining that simulations are experiment
   worlds and may contain candidate protocol specimens; root `protocols/` is for
   the wire-lab apparatus protocol.
2. Create `simulations/SIM-piloh-turns-149-208-recovery/` with TE-nizor's
   minimal contract.
3. Add `seed/protocol-tree-migrations.md` that records old paths, source
   commits, destinations, and pCID status for `group-session`, `udp-binding`,
   and `ppx-dr`.
4. `git mv` candidate protocol trees into
   `simulations/<sim>/protocols/<slug>.d/`.
5. Add `protocol-set.md` that references simulation-local protocol specimens by
   path and commit.
6. Add `seed/wire-lab-devs-draft-migration.md` that records the old
   `transports/wire-lab-devs-draft/` path, source commit, destination, and CIDs.
7. Verify current message CIDs, then `git mv` the transport specimen into the
   first simulation's `world/` tree.
8. Add `concerns.md` mapping turns 149-208 concerns to world elements,
   observations, results, and unresolved DFs.
9. Put any new dogfood world state under `simulations/<sim>/world/`.
10. Move all tracked legacy `proposals/` records into the first simulation's
    archive; leave no tracked root `proposals/` path.
11. Leave top-level `protocols/wire-lab.d/`, `docs/`, `DR/`, `tools/`, and
   `implementations/` in place.
12. File or update DR/DI/spec records only after simulation results justify an
   apparatus, dev-guide, or external spec-corpus change.

## Implications for open TODOs, DRs, and DIs

- **TE-nizor:** should be read as "create a simulation boundary and contract,"
  not "move the repo under simulations."
- **TE-pahah:** remains the structural umbrella, but needs this migration-scope
  DF before any broad path movement.
- **TE-vipir:** remains valid in its repo-like protocol-tree insight, but
  `DI-pakid` and `DI-fakin` supersede root-level placement for non-wire-lab
  candidate protocols with simulation-local placement.
- **TE-domat / DR-nugog:** root `transports` preservation becomes less
  important once the transport specimen migrates into a simulation; the remaining
  question is how future group/feed/CAS semantics graduate from simulation
  worlds into any root-level reference layout.
- **TODO-kugod:** should keep residual recovery open until the first simulation
  maps recovered concerns to seed evidence and results.
- **TODO-vuhuj:** should not be interpreted as requiring non-wire-lab protocol
  candidates to stay rooted forever. It remains the prior protocolization
  migration record and needs follow-up under `DI-pakid`.
- **DEV-GUIDE-RESOURCES.md:** should warn guide writers that simulation-local
  protocol specimens produce evidence; they are not themselves final PromiseGrid
  API layout.
- **Future implementation TODO:** after this first simulation migration, file or
  update successor TODOs only for results that need further DR/DI/spec or
  guide-resource handoff.

## Decision status

`decided`. DF-mupoz.3 is locked as Alt 3.A by `DI-pakid`; DF-mupoz.1,
DF-mupoz.2, DF-mupoz.4, DF-mupoz.5, and the first simulation path are locked by
`DI-fakin`. The decided policy is the wire-lab-only root protocols plus
specimen-migration policy:

- keep apparatus and `protocols/wire-lab.d/` rooted;
- move candidate protocol trees such as `group-session`, `udp-binding`, and
  `ppx-dr` into the first simulation that tests them;
- move existing `transports/wire-lab-devs-draft/` into the first simulation
  with CID verification and migration evidence;
- move legacy `proposals/` records into the first simulation archive because
  the live guide-feedback mechanism has moved elsewhere;
- preserve historical evidence with source-path, source-commit, and CID
  migration records rather than preserving the old root location;
- write new concrete world state under `simulations/<sim>/world/`;
- let simulation results graduate through DR/DI/spec work instead of directly
  becoming root protocol sources.
