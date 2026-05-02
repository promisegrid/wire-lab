# TE-30: TODO numbering and per-protocol TODO shape

## Status

DRAFT. Locks the directory shape and naming convention for TODO
records under the TE-29 protocols-as-simulated-repos layout. Closes
out the TODO-numbering question that TE-29 named but did not resolve.
Renumbers all existing TODOs as part of TODO 014.

## Why this TE

TE-29 locked the shape under which each protocol becomes its own
simulated repo (`protocols/<slug>.d/`). It said decisions, open work,
and don't-touch invariants belong inside protocol spec docs. It did
not say what happens to the top-level `TODO/` directory or to its 19
existing TODO records, beyond a vague "TODO 015: DR/TODO/DI
absorption" that was retired before filing.

That left two real problems:

1. **Mixed scope at the top level.** The current `TODO/` mixes
   harness-level work (the bot/Steve coordination process,
   simulation infrastructure, repo-shape changes) with per-protocol
   work (the group-session envelope, UDP-binding v0, ppx-dr message
   migration). A reader cannot tell which is which without opening
   each file.

2. **Global integer numbering does not scale.** As more protocols
   land, top-level `TODO/` grows linearly forever. By the time the
   wire-lab has 20 protocols, "TODO 137" carries no scope hint at
   all. OQ-100.4 (numbering wrap across centuries, from TE-28) is
   the long-horizon version of this concern; it is already real at
   the small scale.

This TE locks Option D from the conversation of 2026-05-01: per-
protocol TODO subtrees inside each `protocols/<slug>.d/`, no top-
level `TODO/` directory, timestamp-named TODO files matching the TE
convention, **including a full renumber of all existing TODOs**.

## Locked shape

### Per-protocol TODO directories

Every protocol's `.d/` tree contains a `TODO/` subdirectory:

```
protocols/<slug>.d/
├── docs/thought-experiments/
├── specs/<slug>-draft.md
├── TODO/
│   ├── TODO.md                              local queue: this protocol's TODOs
│   ├── TODO-YYYYMMDD-HHMMSS-slug.md         per-TODO record
│   └── ...
└── manifest.json
```

The harness is itself a protocol (`protocols/wire-lab.d/`). Its
TODOs live in `protocols/wire-lab.d/TODO/`. **There is no top-level
`TODO/` directory.** Harness work is not special-cased; it lives in
the wire-lab protocol's `TODO/` like every other protocol's work
lives in its own.

### Filename convention

```
TODO-YYYYMMDD-HHMMSS-<slug>.md
```

Identical to the TE convention (TE-20260501-215027-...). Timestamps
are **first-drafted-time**, anchored at the moment the TODO record
was first authored, never updated. A TODO that moves between
protocols (rare; see "Boundary changes" below) keeps its original
timestamp.

For TODOs that already exist and are being migrated under TODO 014,
the canonical first-drafted-time is recovered from the git commit
that introduced the file. Where a file was introduced via merge from
a long-ago feature branch, the parent-commit date on the originating
branch is used. Where unrecoverable (e.g. squash-merged), the
earliest plausible timestamp from git's history of the file is used.
Migration-time is never used as the canonical timestamp; that would
violate TE-25's drafting-time-anchored invariant.

### Local queue: per-protocol `TODO.md`

Each `protocols/<slug>.d/TODO/TODO.md` lists only that protocol's
TODOs. Format follows the existing top-level TODO.md:

```
- [ ] TODO-YYYYMMDD-HHMMSS-<slug> short description
- [x] TODO-YYYYMMDD-HHMMSS-<slug> short description (done)
- [-] TODO-YYYYMMDD-HHMMSS-<slug> short description (RETIRED/FOLDED)
```

Per-protocol TODO.md is what a contributor working inside that
protocol scans first. It is the "what's open here?" answer.

### Master cross-listed queue: `wire-lab.d/TODO/TODO.md`

The harness protocol's `TODO.md` is the master cross-listed queue.
It lists:

1. All harness-level TODOs (those that live in `wire-lab.d/TODO/`).
2. All per-protocol TODOs from every other protocol's `TODO/`,
   referenced by their full path.

This master queue is what the bot and Steve scan at the start of
each session. It is the "what's open across the whole wire-lab?"
answer.

The two queues serve different audiences and are kept in sync by
discipline: when a TODO is added to a protocol's local queue, the
same line is added to the master queue with the protocol's path
prefix.

### Boundary rule

A TODO is **per-protocol** if and only if its work, when executed,
touches files only inside one `protocols/<slug>.d/` tree.

Anything else is **harness-level** and lives in `wire-lab.d/TODO/`.
Specifically:

- TODOs that touch multiple protocols are harness-level.
- TODOs that touch top-level files (`README.md`, top-level
  `manifest.json`, `tools/`) are harness-level.
- TODOs that change the repo shape itself are harness-level.
- TODOs about the bot/Steve workflow, branching policy, commit
  message rules, etc. are harness-level.

The reasoning: the harness's job is precisely cross-protocol
coordination. Anything coordinating across protocols is the
harness's work, not any one protocol's.

### No integer survivors

The full renumber means every existing TODO 001 through 019 migrates
to a timestamp filename inside the appropriate `protocols/<slug>.d/TODO/`.
No integer-named TODOs survive at the top level. RETIRED and FOLDED
status records (currently 015 and 017) migrate as timestamped files
under `wire-lab.d/TODO/` with the same status content; they do not
get a special exemption.

The integer system is retired for TODOs from this point forward.
TEs continue to use small integer aliases (TE-1, TE-2, ...) per
TE-25, but TODO records do not. The asymmetry is deliberate:
- TEs are a relatively small, slow-growing canonical sequence
  worth aliasing with short integers.
- TODOs are a much larger, faster-growing per-protocol stream where
  short integer aliases create more confusion than they save.

## Existing TODO triage

Going through the 19 existing top-level TODOs and assigning each to
its destination under the locked shape:

| Current | Scope            | Destination                                        |
|---------|------------------|----------------------------------------------------|
| 001     | harness          | `wire-lab.d/TODO/`                                 |
| 002     | harness          | `wire-lab.d/TODO/`                                 |
| 003     | harness          | `wire-lab.d/TODO/`                                 |
| 004     | harness          | `wire-lab.d/TODO/`                                 |
| 005     | harness          | `wire-lab.d/TODO/`                                 |
| 006     | harness          | `wire-lab.d/TODO/`                                 |
| 007     | harness          | `wire-lab.d/TODO/`                                 |
| 008     | harness          | `wire-lab.d/TODO/`                                 |
| 009     | harness          | `wire-lab.d/TODO/`                                 |
| 010     | harness          | `wire-lab.d/TODO/`                                 |
| 011     | harness          | `wire-lab.d/TODO/`                                 |
| 012     | per-protocol     | `group-session.d/TODO/`                            |
| 013     | harness          | `wire-lab.d/TODO/` (multi-protocol shape change)   |
| 014     | harness          | `wire-lab.d/TODO/` (the migration itself)          |
| 015     | harness          | `wire-lab.d/TODO/` (RETIRED status preserved)      |
| 016     | per-protocol     | `ppx-dr.d/TODO/` (BLOCKED status preserved)        |
| 017     | harness          | `wire-lab.d/TODO/` (FOLDED status preserved)       |
| 018     | per-protocol     | `udp-binding.d/TODO/`                              |
| 019     | harness          | `wire-lab.d/TODO/` (tools/ infrastructure)         |

Sixteen of nineteen are harness-level (most of the existing TODOs
are about the bot/Steve workflow and repo-shape evolution; only
three are bounded inside one protocol).

Three move into per-protocol homes:
- 012 (group-transport-envelope) -> `group-session.d/TODO/`
- 016 (proposals-as-transport-messages, BLOCKED) -> `ppx-dr.d/TODO/`
- 018 (udp-binding-v0-reference-implementation) -> `udp-binding.d/TODO/`

TODO 019 (ns-3 harness scaffold) is *harness-level*, not per-binding,
because the scaffold lives under `tools/ns3-harness/` (top-level
infrastructure that serves all bindings) and its first scenario
proves UDP-binding only as a bring-up gate. Future binding-specific
test scenarios may spawn protocol-local TODOs; the scaffold itself
is harness.

## Migration mechanics (delegated to TODO 014)

This TE locks shape but does not move files. The renumber is added
as new subtasks 10 and 11 of TODO 014 (the protocols-as-simulated-
repos migration). Doing it inside TODO 014 means the move happens
atomically with the rest of the protocols/ migration; the wire-lab
is never in a half-renumbered state.

TODO 014 step 10: For each existing TODO, recover the canonical
first-drafted-time from git history. Output: a mapping of
`<old-integer-filename>` -> `<new-timestamp-filename>`.

TODO 014 step 11: Move each TODO into its destination
`protocols/<slug>.d/TODO/` directory under its new timestamp
filename. Write each protocol's local `TODO.md`. Write the master
cross-listed `TODO.md` at `protocols/wire-lab.d/TODO/TODO.md`.
Delete the top-level `TODO/` directory entirely. Update all in-tree
references (DR-009, DI-009, TE-24, TE-29, TE-30 itself, harness-spec
section 8, TE indexes, prior-TODO references) to point at new paths.

## Open questions

OQ-30.1: **Boundary changes after first-drafted-time.** A TODO may
start as per-protocol and grow harness-level scope as work proceeds
(or vice versa). When this happens, does the file move directories
(option A: filename-and-timestamp pinned at first-drafted-time, only
the location changes) or stay put with a forwarding stub (option B:
file stays put, redirect record at new location)? Lean: A. Files
move freely; timestamps are pinned at origin and have no semantics
about location.

OQ-30.2: **Sync discipline for the master TODO.md.** Each per-
protocol TODO.md is the source of truth for that protocol's TODOs;
the master TODO.md cross-lists them. Two options: (a) the master is
manually kept in sync by whoever adds a TODO; (b) a tool generates
the master from the per-protocol TODO.md files. Lean: (a) for now,
upgrade to (b) when manual drift becomes painful. Probably bundled
into a future enhancement of `tools/spec`.

OQ-30.3: **Do DRs and DIs follow the same shape?** TE-29 said
DR/DI/TODO would be absorbed into spec docs (as inline sections).
TODO 015 retired the directory-absorption framing for DR/DI because
those directories never existed. But this TE introduces a real
per-protocol `TODO/` subdirectory shape; should there be parallel
`DR/` and `DI/` subdirectories per protocol, or do DRs/DIs stay as
inline spec-doc sections? Lean: inline in the spec doc. DRs and DIs
are a much smaller stream than TODOs and benefit from being read in
spec context. If that proves wrong later, file a new TE.

OQ-30.4: **Harness TODOs that are really about the simulated-repo
shape itself.** TODOs like 013, 014, this-very-renumber are
harness-level by the boundary rule, but they are specifically about
the wire-lab's *shape*, not about any particular protocol. Should
there be a sub-category (`wire-lab.d/TODO/shape/`) or do they all
sit alongside ordinary harness TODOs? Lean: alongside. The
distinction is rarely useful; one less directory.

OQ-30.5: **Top-level cleanup.** Once `TODO/` is gone from the top
level, the only top-level directories are `protocols/`, `transports/`,
`tools/`, `proposals/` (until TODO 016 unblocks), and the
top-level files. This is much cleaner. Should `proposals/` retain
its top-level position until TODO 016 unblocks, or should it move to
`transports/<bootstrap>/...` immediately as a one-off rename without
the full freeze ceremony? Lean: keep it where it is. TODO 016 is
where that work belongs.

## Reference to load-bearing constraints

This TE relies on (TE-28):

- **C-1 no central registry:** TODO numbering across protocols was
  a soft form of central registry (a single integer sequence
  allocating names across the whole repo). Per-protocol timestamping
  removes that registry; each protocol generates its own TODO names
  independently.
- **C-4 forking is normal:** A protocol fork (under a different
  pCID) brings its own `TODO/` with it. There is no risk of TODO-
  number collision between forks because there are no TODO numbers.
- **OQ-100.4 numbering wrap (from TE-28):** Per-protocol timestamping
  fully closes the TODO half of this question. The TE half remains
  open (TEs still use integer aliases per TE-25); future TE may
  decide to extend the timestamp-only convention to TEs as well.

## Recommendation

Adopt the locked shape. Bake the full renumber into TODO 014 as
new subtasks 10 and 11. Update TODO/TODO.md and all cross-references
during that migration. Top-level `TODO/` ceases to exist after TODO
014 lands.
