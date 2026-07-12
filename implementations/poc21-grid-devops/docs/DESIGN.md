# POC21 grid DevOps design

## Status

Planned pre-code design. POC21 applies PromiseGrid to pull-based machine
administration: a target machine runs `grid`, pulls CAS/VCS data and executables
from peers, locally decides whether to adopt them, runs approved executables
in-place, and records the complete ordered change history for replay. The
canonical task record is `TODO-kifok`. Source: `DI-zosol`; `DI-tanov`;
`DI-dahaj`; `DI-moson`; `TODO-kifok`.

## Locked direction

POC21 is a separate DevOps POC, not additional POC20 scope. It inherits:

- POC18's sparse CAS/VCS, reference-set, chunking, CAR transfer, continuous sync,
  and in-band large-object lessons;
- POC19's stable `grid` stage0 bootstrap, fetched stage1 microkernel modules,
  local operator root adoption, and host-capability promises;
- POC20's CAS event-source timeline, replayable projections, root-decision
  records, and corrective-history model.

The first executable target is a disposable container that runs `grid` as root.
This is intentionally not host-root mutation. The container target is enough to
test filesystem mutation, process restart, service-trigger promises, executable
replacement, and replay without making the developer's host the experimental
machine. Source: `DI-tanov`.

## First POC21 slice order

1. **Actual stage0 self-update.** Bob's target container starts with a minimal
   `grid` stage0 binary. Alice offers a candidate stage0 replacement as a
   CID-named CAS object. Bob pulls it over exact PromiseGrid `grid()` CBOR over
   TCP, verifies CID/proof/local trust criteria, records a local approval or
   rejection promise, writes the approved binary beside the current binary,
   switches atomically where the container filesystem supports it, and restarts
   into the approved bytes. The local CAS-backed machine-change journal records
   only the minimum replay facts needed for this proof: current running stage0 CID
   if known, candidate stage0 CID, local decision state, and post-restart running
   CID/outcome. Local reason, impact summary CID, and recovery notes are optional
   when useful. Source: `DI-nafat`.
2. **Ordered configuration replay.** Bob pulls an ordered machine-change journal:
   prerequisites, file changes, executable changes, trigger promises, validation
   promises, parent links, and resulting events. Bob applies changes in tested
   order, runs selected triggers, validates locally, and records each action and
   result as parent-linked CAS events. Carol starts as a duplicate target and
   replays the same journal to prove the sequence can reproduce the materialized
   state when local inputs and promises match.
3. **Package/container-image artifact distribution.** Dave offers package-like
   blobs or OCI-layer-shaped artifacts as in-band CAS/VCS objects. Bob pulls
   missing chunks on demand and applies them through the same ordered journal
   mechanism rather than through an external package store.

Source: `DI-dahaj`.

## Runtime and safety model

Future POC21 runtime state is bounded under:

```text
/tmp/wire-lab-poc21-run/<run_id>/
```

Expected generated subtrees include:

```text
/tmp/wire-lab-poc21-run/<run_id>/targets/<agent>/rootfs/
/tmp/wire-lab-poc21-run/<run_id>/targets/<agent>/cas/
/tmp/wire-lab-poc21-run/<run_id>/targets/<agent>/events/
/tmp/wire-lab-poc21-run/<run_id>/diagnostics/
```

All cross-agent behavior must use promise-shaped `grid()` CBOR over TCP. Docker
volumes may mount fixture inputs, target root filesystems, and diagnostics, but
must not become the communication path between agents. Any observer or analyzer
is test machinery, not a production monitor or trust source.

POC21 must keep the infrastructures.org/isconf/Turing-equivalence lesson
explicit: updating a machine changes the machine that runs later updates. The
tested sequence and production sequence therefore need ordered replay, not
unordered desired-state convergence alone. Recovery is corrective roll-forward
with preserved history; prior binaries, prior Merkle/CAS roots, and prior
materialized trees are recovery inputs, not promises that the prior universe can
be restored. Source: `DI-moson`.

## Future acceptance gates

The future `poc21-analyze` must fail unless it proves:

- no host-root path was mutated;
- root execution happened only inside disposable target containers;
- every inter-agent message was exact promise-shaped `grid()` CBOR over TCP;
- no inter-agent payload was transferred through a shared Docker volume;
- stage0 self-update executed the replacement binary inside the target container
  and recorded the minimum replay facts in target-local CAS;
- ordered change journals, triggers, validations, and results are replayable from
  target-local CAS events;
- a duplicate target can replay the same journal to matching materialized state
  when local inputs and promises match;
- package/image-style artifacts are stored and transferred as in-band CAS/VCS
  objects;
- derived indexes, diagnostics, and summaries can be deleted and rebuilt from
  CAS-backed source events.

## Non-goals

- No host-root mutation.
- No bare-metal target in the first slice.
- No VM requirement before the container proof works.
- No package manager, container registry, or update decision maker that overrides
  local adoption promises.
- No rollback guarantee after newer code has touched data, binaries, peers,
  devices, or outputs.
