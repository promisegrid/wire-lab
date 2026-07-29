# Grid language profiles and runtime descriptors

This design note explains the current POC21 Grid-language direction in plain
English. It is not a frozen language specification, runtime specification, or
wire-protocol specification. Locked decisions come from `DI-rigob` and
`DI-bigap`; open source-header and canonical-language questions remain in
`DR-lupiz` and `DR-junaz`.

## Short version

PromiseGrid should have one language family with a common data/declaration
layer and two executable profiles:

1. A finite **Gridfile journal profile** for human-readable machine-lifecycle
   configuration and ordered change journals.
2. A Turing-complete **Grid program profile** for `*.grid` agents,
   applications, parsers, builders, planners, and pure-function services.

The common data layer also gives ordinary configuration files one shared syntax
without requiring startup to execute arbitrary code. POC21 will implement the
full program profile with a Go AST interpreter fetched as stage1, while stage0
will understand only the small effect-free data subset needed to find and verify
stage1. Source: `DI-rigob`; `DI-bigap`.

## One language family, not one unrestricted mode

The profiles should share:

- scalar, array, map, record, sum-type, and CID values;
- local symbols whose authoritative targets are CIDs;
- CID-pinned imports and reference sets;
- executable and runtime descriptors;
- protocol-spec and pCID values where wire messages require them;
- source locations, diagnostics, formatting, and editor tooling;
- pure functions and deterministic value evaluation where the selected profile
  permits them.

They should not share every control-flow capability. A Gridfile must remain
finite, ordered, and statically expandable. It cannot contain a local pragma
that silently turns its machine-change journal into an unrestricted program.
That boundary preserves the reason for having a Gridfile: a person or agent can
inspect the complete intended order before machine mutation begins.

A `*.grid` program may use unbounded recursion or iteration and is therefore
Turing complete in the abstract. A real node may still promise only bounded CPU,
memory, time, storage, and message resources. Resource limits stop one
invocation; they do not redefine the abstract language as non-Turing-complete.
Source: `DI-rigob`; `DI-bigap`.

## Configuration is the common data subset

The common data layer makes a Grid-family configuration file possible without
making configuration itself executable. It may contain literals, typed values,
CIDs, symbols, records, finite collections, and CID-pinned imports. It may not
loop, recurse, read clocks, sample randomness, access files or devices, send
messages, or request machine mutation.

This distinction is load-bearing for the installed `grid` stage0 binary. At
first boot, stage0 needs enough local information to recognize owner anchors,
bootstrap roots, initial peers, and the stage1 closure. It cannot depend on an
arbitrary interpreter that it has not yet found, verified, and locally chosen to
run. Stage0 therefore contains only the common data reader. Richer configuration
evaluation belongs to fetched stage1 code. Source: `DI-bigap`.

After bootstrap, a node may store selected configuration source and its parsed
objects in local CAS. CAS identity does not make the configuration universally
trusted or shareable. Each agent still decides locally what it retains, sends,
adopts, or keeps private.

## A fixed prelude is unavoidable

A language-selection directive can provide wide future flexibility, but the
first bytes must have a tiny fixed representation. The loader must recognize
those bytes before it knows which language understands the rest of the source.
A conceptual text rendering is:

```text
#!grid <base32 CID>
<source bytes interpreted according to the referenced object>
```

This is illustrative syntax, not a locked header. The CID's meaning is the open
question in `DR-lupiz` and `TE-fakof`.

Once that prelude has been recognized, the remaining bytes could be interpreted
as a canonical Grid program, a Gridfile, another textual language, bytecode, or
another source form described by a locally accepted runtime. The directive
selects semantics or machinery; it does not itself promise that the local node
will fetch, retain, or execute anything.

## Source identity and execution identity differ

One reproducibility model discussed for POC21 is the tuple:

```text
(source CID, language-spec CID, runtime-descriptor CID, input/context CIDs)
```

This exposes a tension that the shebang TE must resolve:

- Putting an exact runtime descriptor in source makes direct execution
  reproducible, but changing the compiler or interpreter also changes the source
  bytes and source CID.
- Putting only language semantics in source keeps the same source portable
  across conforming implementations, but exact replay must record the selected
  implementation elsewhere.
- A separate execution descriptor can bind one source CID to an exact runtime,
  platform artifact, inputs, capabilities, and entrypoint without rewriting the
  source whenever the implementation changes.

The current recommendation is to keep source semantics and exact execution
context separable, but that recommendation is not locked. `DR-lupiz` remains
open pending `TE-fakof`.

## pCID is not automatically a language identifier

The current PromiseGrid definition of pCID is the CID of a wire-protocol
specification used to select a parser or builder for a `grid()` message. A
programming-language specification is not automatically a wire protocol.

`TE-fakof` deliberately tests both possibilities instead of assuming the answer:

- extend pCID vocabulary so a language specification may be used in the source
  prelude; or
- keep pCID specific to wire protocols and use an ordinary CID for a language
  specification.

The first alternative may reuse familiar pCID-selected dispatch machinery. It
also risks conflating a source file with a wire message and making pCID mean two
different kinds of specification. No production-facing document should present
that extension as settled while `DR-lupiz` is open.

## Runtime descriptors

A raw executable CID identifies exact bytes but does not by itself explain how
to run those bytes. A runtime descriptor can instead name or link:

- the language specification CID;
- an interpreter, compiler, bytecode engine, or runtime implementation;
- exact executable and dependency CIDs for supported platform profiles;
- the source media type and entrypoint convention;
- required local resource and capability promises;
- expected input and output shapes;
- deterministic-build or result-verification information where available.

The descriptor is a local decision input, not a command. A node may recognize
the descriptor, fetch its closure, decline it, or offer a narrower resource
promise than requested. Exact execution records should name the concrete
artifact actually selected for the local platform.

This continues earlier experiments in `promisebase`, where a `lang1` source
file began with the content address of its interpreter and `pb exec` fetched the
interpreter from CAS. The older `rfc-lang` sketch also used hash-pinned imports.
POC21 replaces raw hashes with CIDs and replaces an implicit host environment
with explicit runtime and capability descriptions.

## Candidate language influences

Canonical syntax and semantics remain open under `DR-junaz`. The following
families contribute different useful ideas.

### Go-like

Go-like blocks, functions, explicit control flow, and familiar record syntax are
approachable for general programming. They are less natural for Makefile-style
ordered prerequisite journals unless the language adds special top-level
declarations.

```grid
func factorial(n Int) Int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}
```

### HCL-like

HCL-like blocks are strong for readable configuration and structured
declarations. They become awkward when general functions, recursion, and agent
event loops require embedding another expression language.

```grid
peer "alice" {
    endpoint = "tcp://alice:7447"
}
```

### Lisp-like

Lisp gives one small uniform grammar for code and data and makes
metaprogramming inexpensive. It is less familiar to many operators reading
machine configuration and journals.

```grid
(define (factorial n)
  (if (<= n 1) 1 (* n (factorial (- n 1)))))
```

### Makefile-inspired

Makefile prerequisite and recipe structure remains the clearest known model for
the ordered Gridfile profile because isconf and decomk proved that people can
read, review, and repair a machine-lifetime journal in that form.

```grid
once install-http: base-os
    invoke tools.install packages.http

entry boot: install-http
    invoke tools.start services.hello
```

This need not reproduce Make's implicit variables, string expansion phases,
tab hazards, shell dependence, or implicit-rule machinery.

### Strict ML or Haskell-like

Immutable values, algebraic data types, pattern matching, pure functions, and
explicit effects fit PromiseGrid's deterministic-computation and pCID-selected
payload goals. Strict evaluation appears easier to meter and reason about than
Haskell-style laziness.

```grid
factorial : Int -> Int
factorial n =
    if n <= 1 then 1
    else n * factorial (n - 1)
```

### Unison-like

Unison's most relevant idea is content-addressed definition identity. A human
name may change without changing a function's identity; a body change creates a
new identity. That aligns with CAS-native imports, local symbol tables, and
reproducible execution.

### Nix-like

Nix-like pure evaluation, content-addressed inputs, and derivation-oriented
outputs fit configuration and build planning. Lazy evaluation and derivation
semantics are less naturally suited to long-running autonomous agents unless
substantially extended.

### Erlang or Elixir-like

Pattern matching, isolated processes, supervision, and receive loops fit
autonomous agents. A Grid language must still avoid treating ordinary actor
messages as commands: inter-agent messages remain voluntary, pCID-defined
promises.

### Rebol or Red-like

Blocks that may be data or code make configuration and programming share one
compact syntax. Dynamic interpretation and unfamiliar tooling make behavior
harder to review statically.

### Forth-like

Forth offers a tiny interpreter and excellent constrained-device potential. Its
stack-oriented source is a poor default for rich typed payloads and
human-readable machine journals.

## Current semantic recommendation

The current recommendation, still open under `DR-junaz`, combines:

- strict ML-like pure functions and algebraic data;
- Unison-like content-addressed definitions and local names;
- explicit capability or effect types for CAS, messaging, process, network,
  clock, randomness, and device interactions;
- Makefile-inspired top-level Gridfile prerequisite rules;
- simple finite record and collection expressions for configuration;
- a fixed source prelude that permits other language descriptors indefinitely.

This combination joins POC20 pure-function context, POC18 content-addressed
history, POC19 runtime descriptors, and POC21 ordered machine journals. It is a
research recommendation, not yet a language specification or locked DI.

## First execution engine and cost

POC21 will first implement the full program profile as a Go AST interpreter.
Language semantics must remain independent of that engine so later bytecode,
WASM, native, or constrained-device engines can implement the same language.
Source: `DI-bigap`.

Planning estimates for a POC-quality implementation with source positions,
useful errors, config/journal/program profiles, CID-pinned imports, resource
metering, tests, and documentation are:

| First engine | Approximate code | Approximate model tokens | Relative cost |
| --- | ---: | ---: | ---: |
| AST interpreter | 4–7 KLOC | 180k–300k | 1.0x |
| Bytecode VM | 7–11 KLOC | 300k–500k | 1.7x |
| WASM compiler | 10–16 KLOC | 450k–800k | 2.5x |

These are design estimates, not billing promises. A dynamic functional
interpreter with pattern matching may rise toward 220k–350k model tokens.
Static ML-style inference, algebraic data, and checked effects may raise the
language portion toward 400k–700k tokens. The remaining POC21 DevOps work is
separate.

## Effects and PromiseGrid behavior

Pure computation should be the default. Filesystem, network, clock, randomness,
process, key, storage, and device interactions should occur through explicit
local kernel-role messages and capability promises. A function or agent should
not acquire ambient access merely because its source is locally present.

Inter-agent behavior remains exact promise-shaped `grid()` CBOR over TCP. A Grid
program may implement pCID-specific parser or builder roles, but its source
language CID is not a destination, operation code, peer address, or pCID.

The word `promise` should retain its Promise Theory meaning. The language should
not reuse it as the common asynchronous future/promise abstraction and thereby
make program text ambiguous.

## POC21 consequences

POC21 now includes more than a DevOps journal parser. Its implementation plan
must provide:

- a stage0 reader for the finite bootstrap data subset;
- a fetched stage1 AST interpreter for the program profile;
- a Gridfile parser and expander sharing the same data and diagnostic layer;
- CID-addressed source, runtime descriptors, imports, inputs, and outputs;
- explicit local resource limits and capability promises;
- at least one real Grid program used by the POC21 DevOps scenario;
- tests proving Gridfile rejects unrestricted program constructs;
- tests proving the program profile can perform unbounded computation in the
  abstract while practical runs remain resource-bounded;
- diagnostic output that exposes the source CID, language identity, selected
  runtime artifact, inputs, and resulting local events.

## Open decisions

- `DR-lupiz` asks what the fixed source prelude identifies. `TE-fakof` provides
  the scenario comparison before DF locks the answer.
- `DR-junaz` asks which canonical syntax, typing, effect, evaluation, and
  content-addressed-definition model POC21 implements.
- Exact source-header spelling, file extensions, standard-library names,
  function names, variable names, and runtime paths remain future DF work.

## Source trail

- `DI-rigob` locks the shared language family, finite Gridfile profile,
  Turing-complete program profile, and full POC21 scope.
- `DI-bigap` locks the stage0 data subset and first Go AST interpreter.
- `DR-lupiz` owns source-header identity.
- `DR-junaz` owns canonical language design.
- `TE-fakof` compares shebang identity alternatives.
- `TODO-kifok` owns POC21 planning and implementation.
