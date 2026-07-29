# DR-junaz - Canonical Grid language design

DR-ID: DR-junaz
Date: 2026-07-29 11:39:35 PDT
Asked by: stevegt@t7a.org (Steve Traugott)
State: open
Question: Which syntax, type system, effect model, evaluation strategy, and
content-addressed definition model should the canonical Turing-complete Grid
programming language use?
Why this blocks progress: POC21 has locked a shared language family, a finite
Gridfile journal profile, a Turing-complete `*.grid` program profile, and an AST
interpreter as the first engine. It still cannot implement the canonical program
profile until it compares ML/Haskell, Unison, Nix, Erlang, Go, HCL, Lisp,
Makefile, Rebol, Forth, and related influences without confusing a recommendation
with a settled language specification.
Affects: `protocols/wire-lab.d/TODO/TODO-kifok-poc21-grid-devops.md`;
`implementations/poc21-grid-devops/docs/DESIGN.md`; the requested Grid-language
design note; future Grid language and runtime specifications.
Unblocks: POC21 parser, type checker, interpreter, standard library, Gridfile
grammar, and canonical `*.grid` source examples.
Waiting on: stevegt@t7a.org (Steve Traugott)
Decision:
Linked DI:
Related commits:
Last updated: 2026-07-29 11:39:35 PDT

## Event log

- 2026-07-29 11:39:35 PDT — Opened after the language-planning discussion
  selected a full POC21 language implementation but deliberately left canonical
  syntax and typing unresolved pending a dedicated thought experiment.
- 2026-07-29 11:40:39 PDT — `DN-gagog` recorded the current candidate influences
  and the provisional strict-ML/Unison, explicit-effects, Makefile-journal
  synthesis without presenting it as a settled language specification.
