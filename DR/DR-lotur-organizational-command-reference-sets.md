# DR-lotur - Organizational command reference sets

DR-ID: DR-lotur
Date: 2026-08-13 19:35:34 PDT
Asked by: stevegt@t7a.org (Steve Traugott)
State: open
Question: Should organizations customize `grid` commands by publishing
versioned CID-addressed command reference sets, and if so how should personal,
team, organization, and upstream command sets compose, handle collisions, and
remain subject to local adoption?
Why this blocks progress: The discussion produced a promising `gcloud`-like
component model for organization-specific commands, but neither the command-root
mechanism nor namespace and precedence behavior is locked. Implementing it now
would turn illustrative examples into accidental API.
Affects: `protocols/wire-lab.d/TODO/TODO-kifok-poc21-grid-devops.md`;
`implementations/poc21-grid-devops/docs/DESIGN.md`;
`implementations/poc21-grid-devops/docs/discussion/grid-language-handoff-20260813.md`;
`DEV-GUIDE-RESOURCES.md`.
Unblocks: Organization-specific `grid` subcommands, command descriptor design,
local command-root adoption, help discovery, and namespace collision behavior.
Waiting on: stevegt@t7a.org (Steve Traugott)
Decision:
Linked DI:
Related commits:
Last updated: 2026-08-13 19:35:34 PDT

## Event log

- 2026-08-13 19:35:34 PDT — Opened while retaining the session's stage0 and
  `gcloud` analogy. The handoff illustrates an organization command root but
  labels it discussion-only rather than a settled PromiseGrid interface.
- 2026-08-13 19:35:34 PDT — A future TE must compare explicit namespaces,
  ordered overlays, and merged reference sets with collision detection while
  preserving local owner adoption and old-root history.
