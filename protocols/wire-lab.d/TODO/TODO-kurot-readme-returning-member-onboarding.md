# TODO-kurot: README returning-member onboarding

## Status

Planned and first pass implemented. Owns the README rewrite that helps a
returning PromiseGrid team member come up to speed on wire-lab's current POCs,
wire protocol direction, and production-shape trajectory. Source: `DI-minol`.

This TODO also records the writing guidance to reduce vague use of the word
"pressure" in future work when a clearer term fits. Source: `DI-gapav`.

## Decision Intent Log

ID: DI-minol
Date: 2026-07-08 17:56:01 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Rewrite the root `README.md` as a long-form onboarding primer for a
returning team member, using a brief overview, pre-POC design-search context,
grouped POC journey, current wire/CAS/token model, and POC19 production-shape
direction; refresh `DEV-GUIDE-RESOURCES.md` so it points to the improved README
as the first stop for returning contributors.
Intent: The existing README was too short for someone who knows the basic
PromiseGrid philosophy but has missed months of simulations, GA design search,
POC evolution, wire-protocol convergence, CAS/VCS work, and POC19 planning. The
root README should orient that person quickly without duplicating the detailed
source map in `DEV-GUIDE-RESOURCES.md`.
Constraints: Keep the README readable as a clean narrative rather than a
postmortem; center the POC journey and end with the current envelope, pCID, CBOR,
CID, CAS, parent-link, token, local-trust, and POC19 `grid daemon` / `grid run`
model; distinguish current direction from frozen API only moderately; do not
change implementation code or generated artifacts.
Affects: `README.md`; `DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO.md`;
`protocols/wire-lab.d/TODO/TODO-kurot-readme-returning-member-onboarding.md`.

ID: DI-gapav
Date: 2026-07-08 17:56:01 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Reduce future use of the word "pressure" in new or touched
documentation, plans, prompts, and summaries when a clearer term captures the
meaning.
Intent: "Pressure" has been used as a catch-all for test cases, constraints,
risks, design comparisons, scenarios, implementation requirements, and evidence.
That makes onboarding prose less clear and can blur the difference between
design search, executable evidence, and locked decisions. Future work should
prefer precise terms such as "design search," "comparison," "calibration,"
"candidate design," "scenario," "constraint," "requirement," "risk,"
"tradeoff," or "evidence gathering," depending on meaning.
Constraints: This DI applies to new and touched prose going forward; it does not
require a repo-wide historical sweep and does not forbid "pressure" when it is
the best technical word in a specific context.
Affects: `README.md`; `DEV-GUIDE-RESOURCES.md`; future docs, TODOs, prompts,
plans, and summaries.

## Tasks

- [x] kurot.1 Lock README onboarding decisions in `DI-minol`.
- [x] kurot.2 Lock future wording guidance for "pressure" in `DI-gapav`.
- [x] kurot.3 Rewrite `README.md` as the returning-member onboarding primer.
- [x] kurot.4 Refresh `DEV-GUIDE-RESOURCES.md` to point to the README as the
  first stop for returning contributors.
- [x] kurot.5 Cross-list this TODO in `protocols/wire-lab.d/TODO/TODO.md`.
- [ ] kurot.6 After the returning member reviews the README, capture any missing
  context or follow-up questions as TODO updates.
