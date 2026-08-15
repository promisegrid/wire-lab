# PromiseGrid Grid Language Design Handoff

Date: 2026-08-13
Repository: `/home/stevegt/lab/wire-lab`
Related POC: POC21 grid DevOps
Repository milestone: `f234f66 Document POC21 Grid language design`

## Purpose

This retained POC21 discussion document gives another engineer or Codex session
enough context to continue the Grid language design without depending on the
original chat session. It combines committed repository decisions with later
discussion, and maps every unresolved architectural question to an authoritative
DR. Source: `DI-rulub`.

Status labels used below:

- **LOCKED** means an active Decision Intent record commits the current design.
- **RECOMMENDED** means a TE or design discussion has a preferred answer, but no
  final Decision Intent has locked it.
- **OPEN** means a DR or required TE still owns the decision.
- **DISCUSSION-ONLY** means the point emerged after commit `f234f66` and is now
  retained here, but no DI has made it authoritative.

## Executive summary

- **LOCKED:** PromiseGrid has one Grid language family with a common
  effect-free data and declaration layer.
- **LOCKED:** That family has a finite, statically expandable Gridfile profile
  and a Turing-complete `*.grid` program profile.
- **LOCKED:** The installed stage0 `grid` executable reads only the bounded
  effect-free bootstrap data subset.
- **LOCKED:** The first complete Grid program engine is a Go AST interpreter
  fetched and launched as stage1 code.
- **LOCKED:** Inter-agent behavior remains promise-based. The language does not
  turn messages into commands or let one agent make promises for another.
- **RECOMMENDED:** Portable source names the language specification with an
  ordinary CID, not a pCID. Exact runtime selection belongs in a separate
  execution descriptor.
- **OPEN:** Canonical syntax, typing, effects, evaluation, and
  content-addressed definition semantics remain undecided.
- **DISCUSSION-ONLY:** The general Grid language should not block implementation of
  the first stage0 executable. Stage0 can use a small fixed bootstrap-data
  encoding while the full language remains under design. `DR-toras` owns the
  exact relationship between that encoding and canonical Grid syntax.
- **DISCUSSION-ONLY:** A Gridfile is the human-readable intended plan. The executed
  machine event journal is the append-only CAS history of actual attempts and
  results. They are related but are not the same object. `DR-lotir` owns the
  final terminology and object contract.

## Decision ownership audit

| Topic | Status | Authoritative owner |
|---|---|---|
| One Grid language family with finite Gridfile and Turing-complete program profiles | Locked | `DI-rigob` |
| Effect-free stage0 subset and fetched Go AST interpreter | Locked | `DI-bigap` |
| Portable source and exact runtime identity | Recommended, needs DF | `TE-fakof`; `DR-lupiz` |
| Canonical syntax, typing, effects, evaluation, and definitions | Open | `DR-junaz` |
| Gridfile plan versus executed machine-event history | Discussion-only, open | `DR-lotir` |
| Run-once and recurring-entrypoint invocation history | Discussion-only, open | `DR-lotir`; syntax under `DR-junaz` |
| Stage0 bootstrap format relationship to canonical Grid | Discussion-only, open | `DR-toras` |
| Organization-specific command reference sets and composition | Discussion-only, open | `DR-lotur` |
| Stage0 and `gcloud` component analogy | Explanatory only | This handoff; bounded by `DI-bigap` |
| Example Gridfile and `*.grid` syntax | Illustrative only | Future decision under `DR-junaz` |

The full Grid-language program profile remains in POC21 under active
`DI-rigob`. This persistence pass does not reopen that scope decision.

## Architecture context

### Stage0 `grid` executable

The installed `grid` binary is a small bootstrap seed, not the complete
microkernel. Its intended responsibilities are:

- read local bootstrap configuration;
- recognize local owner identity and trust anchors;
- recognize initial peers and a bootstrap Merkle/CAS root CID;
- verify CIDs and a small built-in protocol surface;
- fetch a stage1 descriptor and its executable closure;
- obtain local approval from the owning agent or human;
- materialize and launch stage1;
- replace itself only through a separately reviewed self-update path;
- record minimum bootstrap, launch, and self-update events in local CAS.

Stage0 should not contain the complete Grid interpreter, general-purpose VCS,
WASM runtime, broad protocol handler registry, DevOps planner, or application
stack.

### Stage1 microkernel roles and applications

Stage1 is fetched from CAS and supplies the larger local role set:

- transport listeners and senders;
- pCID dispatch and pCID-specific parser/builder roles;
- CAS/VCS storage and sparse peer synchronization;
- app reference-set and executable-descriptor resolution;
- local app interface;
- native, WASM, container, and Grid-language runtimes;
- lifecycle supervision and local resource protection;
- operation-scoped secret and signing services;
- machine plan execution and local event recording.

The first stage1 proof is a native/static executable launched as another process
inside the same POC container. Later stage1 runtimes may host WASM, native
programs, containers, or Grid programs.

### Relationship to CAS and VCS

CAS is the source of truth for:

- language specifications;
- source files and libraries;
- runtime and execution descriptors;
- executable artifacts;
- protocol specifications;
- input and output data;
- Gridfiles;
- machine event histories;
- messages, results, and parent links.

Each node has a sparse, partial CAS. No node is assumed to possess every object.
VCS reference sets provide human-usable names, branches, tags, directories,
application roots, runtime roots, command roots, and machine roots over the
underlying CID-addressed objects. Mutable indexes and databases are disposable
views that must be rebuildable from local CAS.

## Locked language-family decisions

### DI-rigob: one family, two executable profiles

`DI-rigob` locks:

- one shared language family;
- a common data/declaration layer;
- a finite Gridfile journal profile;
- a Turing-complete `*.grid` program profile;
- no local pragma that silently promotes a Gridfile into unrestricted program
  mode;
- the general program profile as current POC21 scope.

The intent is to share CID-native values, imports, types, descriptors,
diagnostics, and tooling across configuration, machine planning, agents,
applications, protocol parsers/builders, planners, and pure-function services.

### DI-bigap: stage0 subset and first engine

`DI-bigap` locks:

- stage0 reads only the common effect-free data/declaration subset;
- stage0 configuration cannot loop, recurse, perform I/O, send messages, or
  request local resources;
- the first full Grid engine is a Go AST interpreter fetched as stage1;
- language semantics are specified independently of that implementation;
- later engines may use bytecode, WASM, native code, or constrained-device
  implementations;
- practical runs may be bounded by CPU, memory, time, storage, and messaging
  promises even though the abstract language is Turing complete.

## The three language profiles

### 1. Effect-free data and declaration subset

This is the only part stage0 needs. It should support enough finite structure to
describe:

- scalar values and finite collections;
- CIDs and references;
- owner and trust-anchor identifiers;
- initial peers;
- bootstrap root CIDs;
- runtime descriptors;
- local paths and resource bounds;
- deterministic imports of finite data.

It does not perform effects and does not contain unbounded control flow.

### 2. Gridfile profile

A Gridfile is a finite, ordered, human-readable machine plan. It should support:

- named actions;
- prerequisite relationships;
- finite static expansion before target mutation;
- CID-addressed inputs and executable descriptors;
- explicit local capabilities;
- promises made by the agents that actually perform work;
- validation and trigger relationships;
- stable run-once action identities;
- diagnostics that explain the complete intended order.

A Gridfile must reject recursion, unbounded iteration, dynamically discovered
new plan nodes, or other constructs that prevent the complete plan from being
reviewed before mutation begins.

### 3. Turing-complete `*.grid` program profile

The full program profile may support:

- functions and recursion;
- iteration;
- algebraic or record data types;
- pattern matching;
- modules and CID-pinned imports;
- pure functions;
- explicit capability or effect types;
- agents, applications, parsers, builders, planners, and pure-function servers.

The program may be computationally universal, but actual effects remain bounded
by resources and capabilities promised by local kernel roles.

## Gridfile plan versus machine event journal

**DISCUSSION-ONLY clarification:** Existing POC21 prose sometimes uses “journal” for
two related but different artifacts. Future repository changes should separate
them explicitly after `DR-lotir` is resolved.

### Gridfile: intended ordered plan

The Gridfile says what should be attempted and in what dependency order. It is
human-written or human-reviewable, refactorable, and intended to expose the full
plan before machine mutation.

### Machine event journal: executed history

The machine event journal is the append-only, parent-linked CAS history of what
was actually attempted and what occurred. It records enough exact CIDs and local
outcomes to explain and replay the machine's history.

The two artifacts interact as follows:

```text
Gridfile source CID
        |
        v
finite expanded plan CID
        |
        v
action attempt event -> action result event -> trigger event -> validation event
        |                       |                    |                |
        +-----------------------+--------------------+----------------+
                                |
                                v
                    machine event timeline head
```

A `once` action is not complete merely because a mutable marker file exists. A
runtime should find a successful CAS event matching the action definition,
input CIDs, prerequisite result CIDs, relevant context, and action identity.

**DISCUSSION-ONLY run-once model:** The session proposed a stronger replacement
for `touch $@`:

1. Append a start event containing the machine identity, stable entrypoint and
   action identities, exact action/descriptor/input CIDs, and a stable invocation
   CID.
2. Pass that invocation CID to the local module or cooperating peer.
3. Append a completion event that links to the start event only after the
   promised work completes.
4. Skip a later run-once action only when local CAS contains the exact matching
   completion.
5. Treat a start without completion as indeterminate unless the performing agent
   explicitly promised that retrying the same invocation CID is safe.

This model does not claim generic exactly-once effects. A cooperating peer may
promise to deduplicate the invocation CID, but external messages and physical
effects can leave an ambiguous result after interruption. `DR-lotir` owns the
event and invocation identities; exact Gridfile syntax remains under
`DR-junaz`.

**DISCUSSION-ONLY recurring-entrypoint model:** Run-once preparation and recurring
work should remain in one finite prerequisite graph. An external clock,
operator, local scheduler, peer, or device event sends a promise-shaped local
`grid()` message asking the Gridfile engine to consider invoking a named
entrypoint. The engine decides locally whether it promises to do so, pins the
currently adopted Gridfile CID, satisfies missing run-once prerequisites,
executes the entrypoint's per-invocation actions in deterministic order, and
records the invocation in local CAS.

Machine-changing invocations should be serialized. An invocation that begins
under one Gridfile CID finishes under that CID; a newly adopted Gridfile applies
to later invocations. Repeated external requests do not make Gridfile itself
Turing complete because repetition comes from outside the finite plan. Names
such as `once` and `each-invocation` remain illustrative until `DR-junaz` is
decided.

Ordered history matters because each machine change alters the machine that
executes later changes. The same action set in a different order can produce a
different result. Replay supports explanation, duplicate-machine construction,
and disaster recovery, but it does not promise true rollback. External messages,
physical actions, and other side effects cannot generally be undone. Corrections
append new history and proceed by roll-forward.

## Source identity and execution identity

`TE-fakof` compared six alternatives:

- language-spec pCID in source;
- ordinary language-spec CID in source;
- raw interpreter/compiler executable CID in source;
- runtime-descriptor CID in source;
- language CID plus runtime CID in source;
- language-spec CID in source with exact execution described separately.

### TE-fakof recommendation

The strongest survivor is:

```text
portable source
    -> ordinary language-spec CID

exact execution descriptor
    -> source CID
    -> runtime-descriptor CID
    -> selected platform-artifact CID
    -> input and context CIDs
    -> requested local capability promises
```

Reasons:

- source semantics should not change when an interpreter bug is fixed;
- independent implementations should be able to execute the same source;
- exact replay still requires the selected runtime and artifact CIDs;
- pCID should remain the CID of a wire-protocol specification;
- source identifiers and executable identifiers are not trust decisions.

`TE-fakof` remains `needs DF`. Its recommended answer set is:

- `DF-fakof.1`: ordinary language-spec CID;
- `DF-fakof.2`: separate execution descriptor;
- `DF-fakof.3`: runtime descriptor plus selected artifact CID.

## Runtime descriptors

A runtime descriptor bridges portable source and executable bytes. It should be
capable of naming or linking:

- language-spec CID;
- platform and architecture constraints;
- interpreter, compiler, VM, or executable artifact CIDs;
- dependency closure CIDs;
- entrypoint and argument conventions;
- expected pCIDs;
- requested local capabilities;
- input and output conventions;
- lifecycle expectations.

Possessing a descriptor or executable does not require a node to fetch, retain,
or run it. Each node evaluates the descriptor, signer, peer relationship,
requested capabilities, and local resource conditions for itself.

## Promise Theory and effects

The language must preserve the PromiseGrid agent model:

- no agent can make a promise on behalf of another agent;
- cooperation is voluntary;
- trust is local and relationship-relative;
- no identifier, signature, registry entry, or language import creates global
  trust;
- messages remain pCID-defined promises rather than RPC commands;
- importing code does not grant network, filesystem, device, process, secret,
  or peer access;
- effectful operations require explicit local capabilities;
- a kernel role may withdraw a local resource promise when reciprocal terms
  break, but it does not thereby gain authority over independent agents.

Likely effect families include CAS read/write, messaging, process execution,
network use, filesystem access, device access, time, randomness, secret
operations, and lifecycle control. The exact type-and-effect system remains
open under `DR-junaz`.

## Candidate language influences

The committed design note evaluates:

- Go-like syntax for approachable general programming;
- HCL-like syntax for configuration readability;
- Lisp-like syntax for one small code/data grammar;
- Makefile-inspired prerequisites for ordered Gridfiles;
- strict ML/Haskell-like types and explicit effects;
- Unison-like content-addressed definitions;
- Nix-like immutable dependency closure and reproducibility;
- Erlang/Elixir-like agents and supervision;
- Rebol/Red-like code/data blocks;
- Forth-like small interpreters for constrained devices.

The current research recommendation, not yet locked, combines:

- strict ML/Unison-style semantic foundations;
- explicit effects and capabilities;
- CID-addressed definitions and imports;
- Makefile-inspired top-level Gridfile prerequisites;
- simple finite record and collection syntax for configuration;
- a fixed source prelude that can identify other languages indefinitely.

## Illustrative Gridfile

This example is non-normative; canonical syntax is still open.

```grid
machine "alice-web-01"

input nginx_config =
    cid("bafy...nginx-config-bytes")

agent packages =
    runtime(cid("bafy...package-agent-descriptor"))

agent files =
    runtime(cid("bafy...file-agent-descriptor"))

agent services =
    runtime(cid("bafy...service-agent-descriptor"))

capability rootfs =
    local_capability("root-filesystem")

capability service_manager =
    local_capability("service-manager")

once install_nginx:
    ask packages to promise {
        install: "nginx"
        using: rootfs
    }

once configure_nginx after install_nginx:
    ask files to promise {
        write: "/etc/nginx/nginx.conf"
        content: nginx_config
        using: rootfs
    }

once restart_nginx after configure_nginx:
    ask services to promise {
        restart: "nginx"
        using: service_manager
    }

once verify_nginx after restart_nginx:
    ask services to promise {
        check: "http://127.0.0.1/"
        expect_status: 200
    }

after_success verify_nginx:
    adopt machine_root
```

Expanded dependency order:

```text
install_nginx
      |
      v
configure_nginx
      |
      v
restart_nginx
      |
      v
verify_nginx
      |
      v
adopt machine_root
```

## Illustrative Turing-complete Grid program

This example is also non-normative. It follows the `TE-fakof` recommendation by
placing an ordinary language-spec CID immediately after the shebang and keeping
the exact runtime outside the source.

```grid
#!/usr/bin/env grid
language cid("bafy...grid-language-spec")

import Console from cid("bafy...console-library")
import List    from cid("bafy...collections-library")

fn fibonacci(n: Int) -> Int {
    if n <= 1 {
        n
    } else {
        fibonacci(n - 1) + fibonacci(n - 2)
    }
}

fn calculate_all(values: List<Int>) -> List<Int> {
    List.map(values, fibonacci)
}

effect fn main(context: Runtime.Context) -> Result<Unit> {
    let results = calculate_all([5, 10, 20])

    Console.write(
        context.capability("stdout"),
        "Fibonacci results: " + List.format(results)
    )?

    ok()
}
```

The source, language specification, imported libraries, runtime descriptor,
selected interpreter artifact, inputs, and results are separate CAS objects.

## Stage0 implementation dependency tree

**DISCUSSION-ONLY synthesis:** The first stage0 executable may not require
completion of the full language design. `DR-toras` must resolve the bootstrap
format relationship before this recommendation changes `kifok.20`.

```text
write installed grid stage0
|
+-- lock stage0 responsibility
|   +-- bootstrap seed, not whole microkernel
|   +-- fetch and launch stage1
|
+-- lock bootstrap data format
|   +-- identity and owner anchors
|   +-- initial peers and root CID
|   +-- CAS and execution-cache paths
|
+-- lock minimal protocol specs
|   +-- bootstrap/fetch
|   +-- object transfer
|   +-- stage1 launch
|   +-- stage0 self-update
|
+-- lock runtime descriptor
|   +-- platform selection
|   +-- executable closure
|   +-- capability requirements
|
+-- implement shared CID and sparse-CAS support
|
+-- implement exact grid() CBOR over TCP
|
+-- implement fetch, verification, local approval, and materialization
|
+-- implement stage1 launch and lifecycle
|
+-- implement stage0 replacement and restart
|
+-- record bootstrap and self-update CAS events
|
+-- run clean container proof
    |
    +-- then fetch Gridfile parser and full Grid interpreter as stage1 modules
```

## Stage0 and the `gcloud` component analogy

**DISCUSSION-ONLY analogy:** The installed stage0 `grid` executable resembles the
`gcloud` front end in one limited respect: one durable command can discover,
download, and manage subcommand components rather than embedding all behavior in
its original binary.

PromiseGrid differs in important ways:

- components and dependency closures are identified by CID;
- components are exchanged through promise-shaped messages;
- installation is local adoption of CAS/VCS reference sets;
- the owner evaluates trust and requested capabilities locally;
- adoption and execution are recorded in local CAS history;
- organizations publish candidate component roots but do not impose them.

## Organization-specific Grid subcommands

**DISCUSSION-ONLY design direction:** An organization might publish a versioned
command reference set. `DR-lotur` owns whether this becomes the POC21 design and
how command sets compose:

```text
organization command root: bafy...acme-command-root
|
+-- commands/deploy     -> bafy...deploy-descriptor
+-- commands/payroll    -> bafy...payroll-descriptor
+-- commands/inventory  -> bafy...inventory-descriptor
+-- commands/compliance -> bafy...compliance-descriptor
```

Each command descriptor can name:

- executable, WASM, container, or Grid-program CIDs;
- dependency and library CIDs;
- argument and output conventions;
- pCIDs sent or received;
- requested capabilities;
- lifecycle and resource expectations;
- help text and provenance.

Running `grid deploy` would cause stage1 to resolve the locally adopted command
root, fetch missing objects, verify CIDs and signatures, evaluate requested
capabilities, and decide whether it promises to run the component.

The organization publishes a root and promises what it contains. Each node's
owner decides whether to adopt it. Updates create new root CIDs and preserve old
roots.

The mechanism and composition rule remain open: personal, team, organization, and upstream
command sets might use explicit namespaces, ordered overlays, or a merged
reference set with collision detection. `DR-lotur` requires a future TE, DF,
and DI before the command registry is implemented.

## Current planning defect

**DISCUSSION-ONLY finding:** `kifok.20` currently says the stage0 common-data reader
must wait for canonical grammar and source-header identity. This unnecessarily
may put `DR-junaz` and `DR-lupiz` on the first stage0 critical path.

Recommended correction:

- lock a deliberately small bootstrap data format for stage0;
- implement the first stage0 fetch/verify/launch proof without the full Grid
  grammar;
- let stage1 fetch the Gridfile parser and full Grid language engine;
- later decide whether the bootstrap subset is a strict subset of canonical Grid
  syntax or merely a stable bootstrap format translated by stage1 tooling.

The last bullet is itself an architectural choice. `DR-toras` owns the required
TE, DF, and DI before the repo TODO is changed.

## Open decisions

### DR-lupiz: source and runtime identity

Still requires Steve to answer `DF-fakof.1` through `DF-fakof.3`. The TE
recommends ordinary language-spec CID, separate execution descriptor, and
runtime descriptor plus selected artifact CID.

### DR-junaz: canonical language design

Still requires a dedicated TE comparing:

- surface syntax;
- static or dynamic typing;
- algebraic types and pattern matching;
- capability/effect typing;
- evaluation order and determinism;
- module and import semantics;
- content-addressed definitions;
- Gridfile and program shared grammar;
- diagnostics and source positions;
- standard library boundaries;
- constrained-device implementation feasibility.

### DR-lotir: Gridfile plan and machine event history

The repo should decide whether to lock distinct terms and identities for the
intended plan and executed event history. Discussion recommends `Gridfile` for
the plan and `machine event journal` or `machine timeline` for actual CAS
events. `DR-lotir` remains open pending its TE and DF.

The same DR owns start/completion invocation identity, indeterminate interrupted
runs, and recurring external entrypoint requests over a locally adopted finite
graph. `DR-junaz` owns the eventual source syntax for those semantics.

### DR-toras: Bootstrap subset relationship

Decide whether the stage0 data format must be:

- a frozen strict subset of canonical Grid syntax;
- a separately specified bootstrap format;
- or a descriptor-selected data language supported by stage0.

The third option risks making stage0 larger and should not be selected without
strong evidence. `DR-toras` remains open pending its TE and DF.

### DR-lotur: Command reference-set composition

Decide namespace and precedence rules for personal, team, organization, and
upstream command roots. The DR also owns whether command reference sets become
the POC21 extension mechanism at all.

## Recommended continuation sequence

1. Keep the full Grid-language work inside POC21 under active `DI-rigob`; do not
   reopen that scope without a superseding TE and DI.
2. Run the `DR-toras` TE separating the minimal stage0 bootstrap format from
   canonical Grid language syntax; correct `kifok.20` only after DF and DI.
3. Answer `DF-fakof.1` through `DF-fakof.3`; record a DI and close `DR-lupiz`.
4. Run the canonical-language TE required by `DR-junaz`.
5. Lock syntax, types, effects, evaluation order, import semantics, and
   Gridfile/program grammar through DF and DI.
6. Write complete language and Gridfile specifications as CID-addressed spec
   documents.
7. Implement the stage0 finite bootstrap-data reader independently of the full
   interpreter.
8. Implement source positions, lexer/parser, common value model, diagnostics,
   and CID import resolution for stage1.
9. Implement the Go AST interpreter for the full program profile.
10. Implement Gridfile parsing and finite static expansion over the shared value
    and diagnostic layer.
11. Implement execution descriptors and runtime selection with exact CID
    recording.
12. Implement explicit local capabilities/effects and PromiseGrid message
    integration.
13. Add one real POC21 DevOps Gridfile and one real Turing-complete Grid program.
14. Add analyzer gates for finite expansion, import closure, resource bounds,
    source/runtime separation, and Promise Theory correctness.
15. Run all generated POC code in containers and preserve exact raw messages and
    CAS objects for review.

## Acceptance scenarios for the language work

- Stage0 can read bootstrap configuration without loading the full interpreter.
- A Gridfile expands to a complete finite dependency graph before mutation.
- A Gridfile containing recursion or dynamic plan growth is rejected.
- A run-once action is recognized from CAS event identity, not a mutable marker
  alone.
- An interrupted start without a matching completion is indeterminate unless
  the performing agent promised safe same-invocation retry.
- A recurring external request enters through a promise-shaped local `grid()`
  message, reuses the pinned finite graph, and does not create an implicit loop
  inside Gridfile.
- A Turing-complete Grid program demonstrates recursion or unbounded iteration.
- The same portable source can execute under two conforming runtime descriptors.
- Exact replay records the selected runtime descriptor and executable artifact.
- Missing imported CAS objects produce an explicit missing-object result.
- A library import grants no implicit network, filesystem, or device access.
- Local capabilities bound all effects.
- A protocol parser implemented in Grid receives messages selected by wire pCID;
  its language CID is never used as a destination, operation, or wire pCID.
- Malicious or malformed descriptors are rejected without crashing stage0 or
  stage1.
- Organization command-root updates preserve previous roots and require local
  adoption.

## Authoritative repository sources

- `implementations/poc21-grid-devops/docs/DESIGN.md`
- `protocols/wire-lab.d/TODO/TODO-kifok-poc21-grid-devops.md`
- `docs/research/DN-gagog-grid-language-profiles-and-runtime-descriptors.md`
- `docs/thought-experiments/TE-fakof-grid-source-shebang-identity.md`
- `DR/DR-junaz-canonical-grid-language-design.md`
- `DR/DR-lupiz-grid-source-shebang-identity.md`
- `DR/DR-lotir-gridfile-plan-and-machine-event-history.md`
- `DR/DR-toras-stage0-bootstrap-format-relationship.md`
- `DR/DR-lotur-organizational-command-reference-sets.md`
- `implementations/poc21-grid-devops/docs/discussion/grid-language-handoff-20260813.md`
- `implementations/poc19-production-shape/docs/DESIGN.md`
- `implementations/poc20-timeline-pure-function-cas-branches/docs/DESIGN.md`
- `DEV-GUIDE-RESOURCES.md`

## Handoff warnings

- Do not treat example syntax in this document as canonical.
- Do not call a language-spec CID a pCID unless a later DI explicitly changes
  the pCID vocabulary.
- Do not put the full interpreter or unrestricted configuration execution into
  stage0.
- Do not conflate the Gridfile plan with the executed machine event history.
- Do not let the full language design block the smallest stage0 fetch-and-launch
  proof without an explicit new decision.
- Do not grant effects merely because code or descriptors are present in CAS.
- Do not model inter-agent behavior as RPC commands or external authority.
- Do not assume that any node has the complete CAS.
- Do not implement organization command roots until namespace composition and
  collision handling are decided under `DR-lotur`.
