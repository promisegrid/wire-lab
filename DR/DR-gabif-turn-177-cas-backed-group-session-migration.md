# DR-gabif - Turn-177 CAS-backed group-session migration

DR-ID: DR-gabif
Date: 2026-05-17 09:28:02
Asked by: stevegt@t7a.org (Steve Traugott)
State: open
Question: What additive migration contract should TODO-pipus use to move from pre-CAS inline `.txt` group-session evidence to a CAS-backed group-session specimen with pointer objects, CAS roots, chunk feeds, and preserved historical bytes?
Why this blocks progress: Turn 177 made pointer files and CAS roots structurally important, but the current group-session specimen still has inline `.txt` files named by CID over canonical text bytes. Rewriting those files is disallowed, while doing nothing leaves the successor architecture untested. TODO-pipus needs an explicit migration decision request before it can design overlap, back-reference, group-identity continuity, authorizing promise, and seal mechanics.
Affects: `protocols/wire-lab.d/TODO/TODO-pipus-te-39-wire-lab-devs-migration.md`; `simulations/SIM-jurar-cas-backed-group-session/`; `simulations/SIM-ludut-wire-lab-devs/world/transports/wire-lab-devs-draft/`; `simulations/SIM-rakot-group-session/protocols/group-session.d/specs/group-session-draft.md`; `simulations/SIM-zazit-chunk-feed-replication/`.
Unblocks: TODO-pipus T-MIG-OPS; additive CAS-backed group-session specimen design; any later migration from historical transport evidence into pointer-and-CAS form.
Waiting on: stevegt@t7a.org (Steve Traugott)
Decision:
Linked DI: DI-navod; DI-pator; DI-davov; DI-majib; DI-vuzot
Related commits:
Last updated: 2026-05-17 22:15:24

## Event log

- 2026-05-17 09:28:02 — Opened during turn-177 cleanup so TODO-pipus has a DR for the concrete additive migration contract instead of relying on replay notes and scenario docs alone.
- 2026-05-17 09:44:46 — Added unanswered next-DF packet and acceptance criteria under `DI-majib`.
- 2026-05-17 22:15:24 — Added turn-193 dogfood urgency under `DI-vuzot`; DF-gabif.4 must account for the need to dogfood developer messaging ASAP without rewriting historical `.txt` evidence or blocking on full CAS migration.

## Evidence

- TE-sihih says historical pre-CAS message files remain readable as historical state and should not be rewritten merely because the pointer-and-CAS shape lands later.
- `simulations/SIM-jurar-cas-backed-group-session/SCENARIOS.md` records the migration pressure cases: additive successor specimen, group-visible identity, parent links through CAS, arbitrary body shape, missing pointee, envelope independence, and historical compatibility.
- `simulations/SIM-zazit-chunk-feed-replication/SCENARIOS.md` records the feed pressure that any migration must respect: feeds move chunks, not group messages, and sparse-CAS behavior is normal.
- `TODO-pipus` owns the operational migration scope and still lists close-old-vs-overlap-vs-atomic-swap, back-reference format, message disposition, authorizing promise, seal mechanics, group-identity continuity, and trigger discipline as open DFs.
- Turn 193 makes developer-message dogfooding urgent again after the promisebase pause. That urgency is evidence for trigger discipline, not authorization to rewrite historical inline `.txt` messages or skip the additive migration decision. Source: DI-vuzot.

## Candidate decisions

- **Alt-A: additive overlap.** Keep historical `.txt` evidence immutable, add a new CAS-backed specimen beside it, and publish explicit provenance/back-reference links from successor records to the historical evidence.
- **Alt-B: sealed successor cutover.** Keep old evidence immutable, seal it with a final migration promise, then start a new CAS-backed group-session instance with continuity metadata.
- **Alt-C: defer migration until L6 profile locks.** Keep only scenario docs for now and do not create any CAS-backed group-session specimen until DR-tumus / TE-43 answers the concrete CAS profile.

## Next DF packet

This is the next user-answerable decision packet for TODO-pipus. It is not
answered here. Source: DI-majib.

- **DF-gabif.1 — Migration shape.** Choose Alt-A additive overlap
  (recommended once a minimal L6 profile exists), Alt-B sealed successor
  cutover, or Alt-C defer until DR-tumus closes.
- **DF-gabif.2 — Historical reference.** Choose explicit back-reference from
  new CAS-backed records to historical `.txt` CIDs, a separate migration
  manifest, or both.
- **DF-gabif.3 — Group-visible identity.** Choose pointer-object CID,
  message-root CID, or envelope/wrapper CID as the identifier that group tools
  and parent links cite.
- **DF-gabif.4 — Trigger discipline.** Choose manual migration after DF/DI,
  migration only after first frozen L6 pCID, or migration only after dogfood
  tooling exists. Any answer should say how near-term dogfood proceeds while the
  CAS-backed successor is still undecided. Source: DI-vuzot.

## Acceptance criteria

- Historical message bytes remain immutable.
- Old and new identities have a documented relationship.
- The first CAS-backed group-session specimen can be created without implying
  that grid-envelope variants or existing `.txt` specimens have been selected
  or superseded.

## Notes

This DR does not authorize rewriting historical message files. Any resulting migration must preserve existing CID evidence and explain how old and new identities relate.
