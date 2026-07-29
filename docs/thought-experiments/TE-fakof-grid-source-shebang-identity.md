# TE-fakof: Grid source shebang identity

## TE ID

`TE-fakof`

## Status

needs DF

## Decision under test

What identity should the fixed prelude of a Grid-family source file contain so a
node can discover how to interpret the remaining source bytes without confusing
language semantics, wire-protocol dispatch, exact implementation identity, or
local execution choice?

The open coordination record is `DR-lupiz`. This TE supports
`TODO-kifok` and the language-family direction locked by `DI-rigob` and
`DI-bigap`.

## Assumptions

- A Grid-family source file has a small fixed prelude followed by source bytes.
  The exact text syntax is not locked. A shebang-like line is used in examples
  because it makes the bootstrap problem visible.
- The source file is stored as an immutable CAS object. Changing any byte,
  including the prelude, changes the source CID.
- The existing pCID definition remains in force during the experiment: a pCID
  is the CID of a wire-protocol specification and selects a parser or builder
  for a `grid()` message. Alternative A explicitly tests changing that rule.
- A language specification and an interpreter or compiler are different
  objects. Multiple implementations may promise to implement one language
  specification.
- A raw executable CID identifies exact bytes but does not by itself state
  platform, ABI, dependencies, entrypoint, source semantics, or requested local
  resource promises.
- A runtime descriptor may refer to a language specification, exact executable
  and dependency CIDs, supported platform profiles, entrypoint, and requested
  capability promises.
- Local trust remains agent-relative. No language specification, runtime
  descriptor, executable, compiler, signer, or peer can require a node to run
  code.
- POC21 stage0 reads only its built-in finite bootstrap data subset. The full
  Grid program interpreter is a fetched stage1 module. Source-header selection
  must not make first boot circular.
- Exact replay requires preserving the concrete runtime artifact and all inputs
  that affected the result, even if the source names only portable language
  semantics.

### Threat and trust model

Alice publishes source and language specifications. Bob operates a Linux x86-64
node. Carol operates an ARM64 node. Dave operates a constrained node that can
run only a small interpreter. Ellen later reconstructs historical executions.
Mallory publishes misleading descriptors, substitutes implementations, removes
objects she once served, and tries to make a source identifier look like an
execution promise.

Alice can promise what her source means and which implementations she tested.
Bob, Carol, Dave, and Ellen each decide locally which source, specification,
runtime, compiler, executable, signer, and peer promises they trust.

## Candidate alternatives

### Alternative A — language-spec pCID in the source prelude

The prelude contains the pCID of the language specification:

```text
#!grid <language-spec-pCID>
```

This alternative deliberately extends pCID beyond wire-protocol specifications.
A node routes the identifier to a registered language handler in a way that
resembles pCID-selected wire-message parsing.

### Alternative B — ordinary language-spec CID in the source prelude

The prelude contains the ordinary CID of a language specification:

```text
#!grid <language-spec-CID>
```

The node finds one locally acceptable implementation that promises to implement
that specification. The source identifies semantics, not one implementation.

### Alternative C — raw interpreter or compiler executable CID

The prelude contains the CID of exact executable bytes:

```text
#!grid <interpreter-or-compiler-executable-CID>
```

This follows the simplest reading of a content-addressed Unix shebang: fetch
those exact bytes and invoke them.

### Alternative D — runtime-descriptor CID in the source prelude

The prelude contains a CID of a runtime descriptor. The descriptor links to the
language specification and to exact implementation closures for one or more
platform profiles:

```text
#!grid <runtime-descriptor-CID>
```

The source therefore carries a default exact runtime context, while the
descriptor handles more detail than a raw executable CID can.

### Alternative E — language-spec CID and runtime-descriptor CID together

The prelude carries both identities:

```text
#!grid <language-spec-CID> <runtime-descriptor-CID>
```

The loader verifies that the descriptor promises to implement the named
specification. The source directly states both portable semantics and a default
runtime.

### Alternative F — language-spec CID in source; exact runtime in a separate execution descriptor

The source prelude contains an ordinary language-spec CID, as in Alternative B.
An immutable execution descriptor separately binds:

```text
source CID
language-spec CID
runtime-descriptor CID
selected platform-artifact CID
input/context CIDs
requested local capability promises
```

The source CID remains stable when a compiler or interpreter changes. Each
actual execution still records exact implementation identity. A named launcher,
reference set, Gridfile stanza, or promise may point to the execution descriptor
when one-click exact execution is desired.

## Scenario analysis

### S1 — Alice publishes one source file; Bob and Carol use different machines

Alice publishes a small Grid program that computes package metadata. Bob can run
Linux x86-64 binaries. Carol can run only ARM64 binaries. Both want the same
language semantics and source CID.

**Alternative A** lets both nodes recognize one pCID and find different
handlers, but it changes pCID from a wire-protocol selector into a general
language selector. Bob's pCID registry now contains both wire parsers and source
interpreters whose inputs are not `grid()` messages.

**Alternative B** keeps one source CID and allows Bob and Carol to select
different implementations. Their result records must name the selected
implementations elsewhere or Ellen cannot reproduce the exact executions.

**Alternative C** fails portability if the CID names x86-64 bytes. A portable
shell-style executable might choose another binary, but that turns the raw
executable into an undeclared descriptor format.

**Alternative D** can list x86-64 and ARM64 implementation closures. Bob and
Carol preserve the same source only if the descriptor remains unchanged. A new
platform variant changes the descriptor CID and therefore changes the source
bytes even when the program body and language semantics are unchanged.

**Alternative E** has the same platform strength as D and makes the semantic
relationship explicit. It still changes the source CID when the runtime
descriptor changes.

**Alternative F** keeps one source CID and gives Bob and Carol distinct exact
execution descriptors. Their result histories can converge on the same output
or diverge visibly without pretending they used identical runtimes.

### S2 — Dave has a constrained runtime

Dave receives Alice's source on a small device. He recognizes the language but
can run only a compact interpreter that supports the required subset and fits
his local CPU and memory promises.

**Alternative A** can route the language pCID to Dave's compact handler, but the
pCID no longer tells readers whether it names wire syntax or source semantics.

**Alternative B** works if Dave has a local mapping from the language-spec CID
to his compact interpreter. He can decline source features outside the subset.

**Alternative C** works only if the source was authored specifically with
Dave's executable CID. A source rewritten for each device loses the benefit of
portable language identity.

**Alternatives D and E** work if the descriptor includes Dave's platform. If it
does not, Dave must wait for a new descriptor, which changes every source file
that embeds it.

**Alternative F** lets Dave retain the same source and use a local execution
descriptor naming his compact interpreter and narrower capabilities. His local
result records remain explicit about the implementation difference.

### S3 — Bob fixes an interpreter bug without changing the language

Bob discovers that interpreter `R1` mishandles integer overflow. Alice publishes
or accepts corrected interpreter `R2`. The language specification and program
body are unchanged.

**Alternatives A and B** preserve the source CID. New executions select R2 and
record it; old executions still name R1 in their historical execution context.
This supports comparison without rewriting source history.

**Alternative C** requires changing the source prelude from the R1 executable
CID to R2, creating a new source CID even though source meaning is intended to
remain unchanged.

**Alternative D** creates a new runtime-descriptor CID and therefore a new
source CID. This tightly couples source identity to implementation maintenance.

**Alternative E** behaves like D while retaining a redundant unchanged
language-spec CID in the prelude.

**Alternative F** preserves source identity and publishes a new execution
descriptor for R2. Bob's historical timeline shows which executions used the
buggy or corrected runtime.

### S4 — Mallory substitutes a runtime

Mallory offers an implementation that claims to implement Alice's language but
sends data to Mallory and returns altered results. Bob has not built trust in
Mallory's implementation promises.

**Alternatives A and B** require Bob's local language-to-runtime selection to
reject Mallory's implementation. The source identifier alone does not protect
him; exact selected runtime identity must be visible before execution.

**Alternative C** makes substitution byte-visible, but a trusted-looking raw
binary still lacks a declared capability surface and language relationship.

**Alternatives D and E** make implementation closure and requested capabilities
reviewable through the descriptor. Bob may still reject the descriptor or its
signer locally.

**Alternative F** gives Bob the same descriptor review while keeping the
decision out of the source identity. A Gridfile or owner promise can name the
exact accepted execution descriptor; Mallory cannot replace it without changing
its CID.

No alternative turns a CID into trust. Each only changes which facts Bob can
evaluate before deciding what he promises to run.

### S5 — Ellen reproduces an execution thirty years later

Ellen has Alice's source, its CAS ancestors, and Bob's execution history. The
original project website and package repositories no longer exist.

**Alternative A** tells Ellen which language pCID was used, assuming future
readers still understand the expanded pCID convention. It does not identify the
exact implementation unless Bob recorded that separately.

**Alternative B** preserves language semantics cleanly but also requires a
separate exact runtime record.

**Alternative C** identifies exact bytes, but Ellen may not know their ABI,
entrypoint, dependency closure, or language semantics after surrounding context
has disappeared.

**Alternative D** preserves a richer closure, provided the descriptor and all
reachable objects were retained. It also means the source itself was revised
for every runtime-descriptor revision.

**Alternative E** preserves both semantics and runtime context in the source,
but still couples source identity to runtime maintenance.

**Alternative F** gives Ellen an immutable semantic source object and a separate
immutable execution object. She can study source history independently and can
reproduce the exact run when the execution closure remains available.

### S6 — Carol uses a different conforming implementation to verify Bob

Bob runs interpreter R1. Carol distrusts R1's maintainer and independently runs
compiler R2 against the same source and explicit inputs. They compare result
CIDs.

**Alternatives A and B** naturally permit independent implementations. Their
verification claim must name each exact implementation.

**Alternative C** makes the source itself select R1, so Carol must override or
rewrite it to use R2. The override is outside the source's declared execution
path.

**Alternative D** can include both implementations, but the descriptor author
chooses that set. Carol may want an implementation absent from Alice's default
descriptor.

**Alternative E** has the same limitation while repeating the language-spec
identity.

**Alternative F** treats independent execution descriptors as normal. Bob and
Carol can compare promises over one source and input set without adopting each
other's runtime.

### S7 — The source travels through mixed-version PromiseGrid nodes

Alice sends a source reference through Bob to Carol. Bob's node understands the
new language. Carol's older node understands the fixed prelude but not the
language or runtime descriptor.

**Alternative A** risks sending the identifier through pCID dispatch paths that
normally expect a wire-message body. An older node may mistake a source-format
identifier for a protocol it should parse at the network boundary.

**Alternative B** lets Carol recognize an ordinary CID and state that she does
not currently promise to interpret it. She may still retain or relay the opaque
source object.

**Alternative C** lets Carol retain exact bytes but gives no portable clue about
what source semantics they implement.

**Alternatives D and E** let Carol retain the descriptor graph opaquely. She can
later fetch a compatible runtime if one appears.

**Alternative F** separates the portable source graph from one or more execution
graphs. Carol can retain either graph independently according to local storage
promises.

### S8 — Compiler output becomes a CAS object

Alice's Grid source is compiled rather than interpreted. The compiler produces a
WASM or native artifact. Bob wants to cache the result; Carol wants to verify it
with another compiler.

**Alternatives A and B** identify source semantics but require a build record
that names source, compiler, compiler inputs, target, and output CID.

**Alternative C** identifies exact compiler bytes but omits dependency closure,
target profile, options, and language specification unless those facts are
recorded elsewhere.

**Alternative D** can describe compiler closure and target behavior, but
embedding the descriptor in source changes source identity whenever build
machinery changes.

**Alternative E** adds semantic clarity but keeps the same coupling.

**Alternative F** makes the build or execution descriptor the natural CAS
record. Bob and Carol may publish different artifact promises from the same
source and compare them without rewriting the source object.

### S9 — Sparse CAS loses part of the closure

Years later, Dave has the source object but not every runtime object. Ellen has a
language specification. Bob has one old executable. No agent has a complete
copy.

**Alternative A** leaves Dave with an identifier whose pCID interpretation may
itself be historically ambiguous if pCID vocabulary changed.

**Alternative B** gives Dave the best semantic artifact: he can search peers for
any implementation that promises the specification.

**Alternative C** gives Bob exact bytes but may leave him unable to explain or
run them safely.

**Alternatives D and E** expose the missing closure as missing descriptor links,
which helps targeted retrieval, but the source is tied to that particular
closure.

**Alternative F** permits partial recovery from independent stores: source,
specification, runtime descriptor, executable, and execution records may be
reassembled by CID without any one store being complete.

### S10 — POC21 routes wire messages and loads source at the same time

Bob's node receives a `grid()` message whose pCID selects the machine-maintenance
parser. The resulting promise refers to a Grid source CID. Stage1 then loads that
source.

**Alternative A** uses pCID twice for different layers: once to parse the wire
message and again to choose a source-language interpreter. The implementation
may keep separate registries, but the vocabulary no longer tells developers
which boundary is meant.

**Alternative B** keeps the layers visible. The wire pCID selects the
machine-maintenance payload parser; the payload's source CID leads to source
whose ordinary language-spec CID selects compatible language machinery.

**Alternative C** also keeps pCID separate but treats exact executable bytes as
if they fully describe execution.

**Alternatives D and E** keep wire pCID separate and improve runtime
description, but couple the runtime choice to source.

**Alternative F** gives the clearest pipeline:

```text
wire pCID -> payload parser -> source CID -> language-spec CID
          -> separate execution descriptor -> exact local runtime artifact
```

Each arrow names a different promise or local decision instead of overloading
one identifier.

## What each alternative makes easier, harder, and newly demands

| Alternative | Easier | Harder | New obligation |
| --- | --- | --- | --- |
| A: language pCID | Reuses pCID-shaped handler lookup | Preserving one meaning for pCID and separating wire from source | Define and migrate a second pCID role everywhere |
| B: language-spec CID | Portable source semantics and independent implementations | Exact replay without another runtime record | Maintain a local spec-to-runtime registry and record selected runtime |
| C: raw executable CID | Exact byte selection | Portability, dependencies, semantics, capabilities, and platform selection | Supply missing execution context out of band |
| D: runtime-descriptor CID | Direct reproducible launch with rich closure | Stable source identity across runtime updates | Retain descriptor closure and revise source when it changes |
| E: two CIDs | Visible semantic/runtime consistency check | Header redundancy and source/runtime coupling | Define mismatch and update rules for both identities |
| F: separate execution descriptor | Stable semantic source plus exact replay and independent verification | More than one CAS object must be followed | Preserve explicit source-to-execution links and local selection records |

## Surviving alternatives

**Alternative B survives** as the smallest portable source identity. It is
insufficient alone for exact replay, but it composes with local runtime records.

**Alternative D survives** for a launcher object whose purpose is to name a
specific executable closure. It is less attractive as the identity embedded in
portable source.

**Alternative F survives most strongly** because it preserves both semantic
source identity and exact execution identity without forcing one to change when
the other changes.

**Alternative A is not recommended.** The scenarios find no benefit that
requires redefining pCID rather than using an ordinary CID and a separate local
language registry. Reusing the pCID word creates cross-layer ambiguity.

**Alternative C is not sufficient as the sole identity.** Exact executable
bytes remain important inside a runtime descriptor or execution record.

**Alternative E is not recommended for the source prelude.** The runtime
descriptor can already name the language specification, while a separate
execution descriptor provides the same consistency check without changing
source identity whenever runtime machinery changes.

## Conclusions

1. **Source semantics and exact execution context are distinct.** A source CID
   should not need to change merely because an interpreter bug is fixed, a new
   platform is added, or a verifier chooses another implementation.
2. **The language specification should use an ordinary CID.** The pCID term
   should remain specific to wire-protocol specifications unless a later DF
   explicitly overturns the existing vocabulary.
3. **Exact runtime identity belongs in a descriptor.** Raw executable CIDs are
   necessary leaves, but a runtime descriptor supplies platform, dependency,
   entrypoint, language, and capability context.
4. **Exact execution should be a separate CAS object.** The recommended model is
   Alternative F: ordinary language-spec CID in source plus a separate execution
   descriptor naming source, runtime, selected artifact, inputs, and local
   capability promises.
5. **A direct-execution convenience may point to the execution descriptor.** A
   Gridfile stanza, named reference set, or promise can offer one exact execution
   object without making that runtime identity part of the portable source.
6. **No identifier is a trust decision.** Nodes evaluate source, specification,
   runtime, artifact, signer, peer, and requested resources locally.

## Implications for open TODOs and pending DIs

- `TODO-kifok` should treat source objects, language specifications, runtime
  descriptors, and execution descriptors as separate POC21 artifacts.
- `DI-rigob` remains unchanged: Gridfile and `*.grid` remain two profiles of one
  language family.
- `DI-bigap` remains unchanged: stage0 reads only the bootstrap data subset and
  stage1 owns the AST interpreter.
- `DR-lupiz` remains open until Steve answers the DF questions below.
- `DR-junaz` remains open; this TE does not choose canonical syntax, typing, or
  effect semantics.
- `DN-gagog` should present the runtime-descriptor-in-shebang idea as historical
  analysis and point readers here for the stronger source/execution separation.

## Decision Framing — questions for the user

**DF-fakof.1 — language identifier vocabulary.** Should the source prelude use
an ordinary language-spec CID or extend pCID to cover language specifications?

- **Alt-1.a: ordinary language-spec CID. Recommended.** Preserve pCID for wire
  protocols and keep source-language selection in its own local registry.
- **Alt-1.b: language-spec pCID.** Expand pCID vocabulary and reuse pCID-shaped
  dispatch for both wire messages and source files.

**DF-fakof.2 — source/runtime coupling.** Should exact runtime identity be part
of source bytes or a separate execution object?

- **Alt-2.a: separate execution descriptor. Recommended.** Keep source identity
  stable and bind exact runtime, platform artifact, inputs, and capabilities in
  an immutable execution descriptor.
- **Alt-2.b: runtime-descriptor CID in source.** Make each source object carry a
  default exact runtime closure.
- **Alt-2.c: language and runtime CIDs in source.** Carry both identities and
  define consistency and update rules.

**DF-fakof.3 — exact implementation representation.** When exact execution is
recorded, should it identify a raw executable or a runtime descriptor?

- **Alt-3.a: runtime descriptor plus selected artifact CID. Recommended.** Keep
  the full closure and expose the exact local executable selected.
- **Alt-3.b: raw interpreter/compiler executable CID.** Treat exact bytes as
  sufficient and obtain dependencies and capabilities elsewhere.

The recommended set is **(1.a, 2.a, 3.a)**.

## Decision status

`needs DF` — the scenario analysis recommends an ordinary language-spec CID in
portable source, a separate execution descriptor, and a runtime descriptor plus
selected artifact CID for exact execution. `DR-lupiz` remains open until Steve
chooses the DF answers and a later DI records the result.
