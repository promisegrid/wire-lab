# SIM-hugoj: CAS / Usenet-like / Git-like lineage

This simulation explores whether part of PromiseGrid should be
understood as a content-addressed-Usenet successor line with git-like
hashed and hash-chained objects and transport separation. It is a
broad design exploration, not a frozen protocol claim and not a claim
that PromiseGrid simply *is* Usenet or git. Source: `DI-pijun`.

`group-session` is the current worked specimen inside this simulation, not the
identity of the simulation itself. The point of the simulation is to ask which
historical invariants age well enough to carry forward into PromiseGrid, which
ones should be adapted, and which ones should be rejected. Source: `DI-pijun`.

## Question

Which parts of the Usenet / FidoNet / email lineage, when combined with
content-addressed storage and git-like transport separation, describe a useful
PromiseGrid design line rather than a misleading analogy? Source: `DI-pijun`.

## Current specimen

The current specimen is `group-session`, because it already exhibits several of
the relevant traits: append-only message growth, DAG/thread structure,
content-addressed message identity, and multiple possible delivery substrates.
That makes it a useful worked example, but not the whole design space. A later
PromiseGrid protocol may preserve some of these traits while diverging sharply
from `group-session` in governance, payload shape, replication, or feed
mechanics. Source: `DI-pijun`.

## What maps cleanly from precedent

- Message identity that stays stable across delivery substrates.
- Store-and-forward replication among named peers or sites.
- Per-instance declaration of how a site exchanges messages with peers.
- Thread/DAG semantics that survive relay across heterogeneous wires.
- Separation between protocol meaning and the substrate that carries bytes.

These are the parts of the turn-173 precedent survey that look structurally
durable across email, Usenet, FidoNet, and modern git-like or gRPC/libp2p-like
systems. See `docs/research/historical-networks-20260503.md`. Source:
`DI-pijun`.

## What does not map cleanly

- Usenet control-message conventions do not automatically become PromiseGrid
  governance rules.
- Git's object model and transport behavior are precedent, not identity; using
  git as a current substrate does not imply PromiseGrid should clone git's
  object types or workflow.
- A content-addressed-Usenet framing does not by itself settle feed naming,
  site manifests, sparse-CAS shape, moderation, or freeze behavior.
- The simulation does not assume that every future PromiseGrid protocol is a
  `group-session` variant.

These are precisely the places where the analogy is useful only if it remains
explicitly exploratory. Source: `DI-pijun`.

## Provenance and authority

The historical grounding for this simulation comes from
`docs/research/historical-networks-20260503.md`, especially its Usenet,
FidoNet, email, and git-adjacent precedent analysis. The immediate trigger for
filing this simulation is the turn-173 replay recovery in
`protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`
and the unresolved follow-on notes in
`protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`.

Nothing in this directory overrides later DRs, DIs, frozen specs, or guide
prose. This simulation is a bounded design workspace for testing one broad
lineage claim and keeping the claim visible until it either graduates or is
rejected. Source: `DI-pijun`.
