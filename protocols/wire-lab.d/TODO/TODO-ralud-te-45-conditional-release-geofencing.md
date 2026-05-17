# TODO-ralud: TE-45 conditional-release / geofencing / recursive promise-graph

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TODO-30` (integer alias)
- `TODO-20260507-002306` (timestamp alias and pre-migration filename)

## Status

Open. Orthogonal to TE-sihih layered model; can land independently.
No twig yet.
Simulation question home: `simulations/SIM-zarud-conditional-release-geofencing/`.
Source: `DI-pukap`.

## Threads absorbed from OPEN-THREADS.md

### T-CONDITIONAL-RELEASE (formerly OPEN-THREADS, opened 2026-05-06)

1 UT but raised in turn 179 as orthogonal architectural axis.

Scope: model conditional-release promises ("send only if recipient
promises onward-restraint") as recursive promise-graph that the
protocol must track (UT-179.c); add the geofencing constraint
dimension as a first-class property of group membership and message
dispatch (UT-179.d); establish whether this lives in transport-spec,
group-session, or a new conditional-release.d/ protocol family.

Blocking: orthogonal to TE-sihih layered model; can land independently.

Anchor: turn 179 framing.
Disposition-file pointer: `dropped-thread-disposition-20260506.md`
§ TE-45 cluster (1 UT).

## Question log

- 2026-05-17: The conditional-release, onward-restraint, geofencing, and
  recursive promise-graph ownership question is now captured as a standalone
  simulation question in
  `simulations/SIM-zarud-conditional-release-geofencing/`. This does not answer
  the owner-layer question or close the TODO. Source: `DI-pukap`.

## Decision Intent Log

(Will be populated as DFs lock and product lands.)
