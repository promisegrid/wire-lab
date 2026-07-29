# DR-lupiz - Grid source shebang identity

DR-ID: DR-lupiz
Date: 2026-07-29 11:39:35 PDT
Asked by: stevegt@t7a.org (Steve Traugott)
State: open
Question: Should a Grid source shebang identify a language specification by
pCID, identify it by an ordinary CID, identify an exact interpreter or compiler,
identify a runtime descriptor, carry both semantic and runtime identities, or
leave exact runtime selection to a separate execution descriptor?
Why this blocks progress: POC21 cannot freeze Grid source identity, parser
selection, reproducible execution context, or the boundary between source and
runtime objects until these alternatives are compared under PromiseGrid's pCID,
CAS, portability, and local capability rules.
Affects: `protocols/wire-lab.d/TODO/TODO-kifok-poc21-grid-devops.md`;
`implementations/poc21-grid-devops/docs/DESIGN.md`; the requested Grid-language
design note and shebang thought experiment.
Unblocks: POC21 Grid-language source format, runtime-descriptor integration, and
the later implementation of Gridfile and `*.grid` execution.
Waiting on: stevegt@t7a.org (Steve Traugott)
Decision:
Linked DI:
Related commits:
Last updated: 2026-07-29 11:39:35 PDT

## Event log

- 2026-07-29 11:39:35 PDT — Opened from the POC21 language-planning discussion.
  The requested TE must test both the existing pCID vocabulary and an ordinary
  language-spec CID rather than assuming the terminology outcome in advance.
- 2026-07-29 11:40:39 PDT — `TE-fakof` filed at
  `docs/thought-experiments/TE-fakof-grid-source-shebang-identity.md` with status
  `needs DF`. Its recommended set is an ordinary language-spec CID in portable
  source, a separate execution descriptor, and a runtime descriptor plus
  selected artifact CID for exact execution.
