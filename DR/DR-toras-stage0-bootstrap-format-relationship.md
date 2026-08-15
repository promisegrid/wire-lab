# DR-toras - Stage0 bootstrap format relationship

DR-ID: DR-toras
Date: 2026-08-13 19:35:34 PDT
Asked by: stevegt@t7a.org (Steve Traugott)
State: open
Question: Should stage0 bootstrap configuration use a frozen strict subset of
canonical Grid syntax, a separately specified fixed bootstrap format, or a
descriptor-selected data language supported by stage0?
Why this blocks progress: `kifok.20` currently places the stage0 data reader
after canonical Grid grammar and source-header decisions. The discussion
identified a credible smaller critical path, but changing that dependency
without a TE and DI could either block stage0 unnecessarily or create a second
configuration language accidentally.
Affects: `protocols/wire-lab.d/TODO/TODO-kifok-poc21-grid-devops.md`;
`implementations/poc21-grid-devops/docs/DESIGN.md`;
`implementations/poc21-grid-devops/docs/discussion/grid-language-handoff-20260813.md`;
`DEV-GUIDE-RESOURCES.md`.
Unblocks: POC21 stage0 bootstrap-data specification, `kifok.20` sequencing, and
the first fetch/verify/launch implementation slice.
Waiting on: stevegt@t7a.org (Steve Traugott)
Decision:
Linked DI:
Related commits:
Last updated: 2026-08-13 19:35:34 PDT

## Event log

- 2026-08-13 19:35:34 PDT — Opened from the Grid-language handoff's finding
  that the full canonical language design need not necessarily block the
  smallest stage0 proof. This DR does not change `DI-bigap` or `kifok.20`.
- 2026-08-13 19:35:34 PDT — The future TE must compare the three alternatives
  under stage0 size, long-term compatibility, independent implementation,
  migration, malformed-input safety, and constrained-device requirements.
