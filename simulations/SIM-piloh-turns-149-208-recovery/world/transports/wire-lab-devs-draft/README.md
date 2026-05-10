# world/transports/wire-lab-devs-draft/

This is the first concrete instance of the small-finite-closed-group
transport-protocol defined in
[`protocols/group-session.d/specs/group-session-draft.md`](../../../protocols/group-session.d/specs/group-session-draft.md).

This directory moved from root `transports/wire-lab-devs-draft/` into this
simulation's `world/transports/` tree per `DI-fakin`. The message files are
unchanged specimens; their CIDs still verify over raw bytes after the move.

## Status: BOOTSTRAP

The protocol's spec is still a draft and its pCID is not yet minted. Per
the spec's freeze gate, freeze requires "at least one real transport
instance has been created and exchanged at least one round-trip." This
instance is that round-trip.

While bootstrapping, the directory is named with the placeholder state
suffix `-draft` (per DF-38.5: pattern `<slug>-<state>` where `<state>` is
`draft` pre-freeze or a CID post-freeze), and messages use the grid
envelope `grid draft:group-session`. Once the spec is frozen and a pCID
is minted, a later DR/DI/spec handoff must decide whether this simulation-local
specimen graduates, stays as evidence, or produces a new frozen successor.

Note: prior to 2026-05-06 this directory was named
`transports/draft--wire-lab-devs/` per the older `draft--<slug>` rule.
DF-38.5 (locked verbally turn 176, written turn 285) renames the
draft-state convention to `<slug>-draft` so unfrozen and frozen forms
sort together in `ls`. The three message files in this directory
(authored before the rename) reference the old path in their body text;
those references should be read as historical and refer to this same
directory at its current path.

## Membership

Closed and fixed at transport creation per §8 of the protocol spec. The
participant set is the population of developer agents collaborating on
wire-lab — multiple humans, each driving one or more LLM agents inside
their own clone of this repository.

Under the per-author-branch git binding (see below), membership is
realized as the set of `<author-id>/main` branches participants are
configured to fetch and propagate from.

## Layout (per spec §1: flat; per spec §2: filename = CID)

```
simulations/SIM-piloh-turns-149-208-recovery/world/transports/wire-lab-devs-draft/
    README.md                      (this file; not a protocol message)
    <message-cid-1>.txt
    <message-cid-2>.txt
    ...
```

Message files are named by their message CID per spec §2, with `.txt`
appended. Filename and content are mutually self-checking: re-hash the
file's bytes and compare to the filename to verify integrity. Per
DI-012-20260508-033513, canonical messages do not include `Message-ID:`.
Readers may tolerate exactly one legacy `Message-ID:` header before `Date:`
for historical compatibility, but the value is ignored semantically and
the message CID remains the only protocol identifier.

`README.md` is a navigation aid for humans; it is not a protocol message
and is ignored by readers walking the message DAG.

## Git binding: per-author branches with content-addressed merge

This transport instance follows
[`group-session-draft.md`](../../../protocols/group-session.d/specs/group-session-draft.md) §9.

### Branch ownership

- Each participant has their own write branch `<author-id>/main`.
- Examples: `ppx/main`, `codex/main`, etc.
- A participant **authors** message files only on their own branch.
- A participant **propagates** other participants' message files onto
  their own branch via the merge step below. Propagation is a verbatim
  file copy — same canonical bytes, same CID, same filename — and is
  not authoring.

### Receive-merge-push-then-optionally-post cycle

Each cycle proceeds in two phases. The merge phase is mandatory whenever
new messages are observed; the post phase is optional.

**Merge phase (mandatory):**

```bash
git fetch --all
# For each known author-id/main branch that is not your own:
#   list *.txt files under simulations/SIM-piloh-turns-149-208-recovery/world/transports/wire-lab-devs-draft/
#   that are not on your own branch.
# For each such file:
#   verify CID = filename (tools/spec cid <file>)
#   verify envelope structure per spec §4
# Copy verified files into your working tree on your own branch.
git add simulations/SIM-piloh-turns-149-208-recovery/world/transports/wire-lab-devs-draft/*.txt
git commit -m "transport: merge <count> messages from <branches>"
git push origin <your-author-id>/main
```

**Post phase (optional):**

```bash
# Author a new message file under simulations/SIM-piloh-turns-149-208-recovery/world/transports/wire-lab-devs-draft/
# following spec §4 (envelope), §5 (body has explicit "I promise ..."),
# §4.6 (Parents: set to message CIDs of direct ancestors), §6 (body-as
# -receipt if acknowledging).
# Compute its CID and rename the file:
NEW_CID=$(tools/spec cid simulations/SIM-piloh-turns-149-208-recovery/world/transports/wire-lab-devs-draft/draft.txt)
mv simulations/SIM-piloh-turns-149-208-recovery/world/transports/wire-lab-devs-draft/draft.txt \
   simulations/SIM-piloh-turns-149-208-recovery/world/transports/wire-lab-devs-draft/${NEW_CID}.txt
git add simulations/SIM-piloh-turns-149-208-recovery/world/transports/wire-lab-devs-draft/${NEW_CID}.txt
git commit -m "transport: post ${NEW_CID}"
git push origin <your-author-id>/main
```

The merge phase precedes the post phase because any new message's
`Parents:` references must be CIDs the agent has already verified and
propagated onto its own branch.

### Convergence

Under steady-state operation, every participant's `<author-id>/main`
branch eventually contains every message anyone has ever posted to the
transport. The transport state is eventually consistent.

The git commit graph is incidental to the message graph. Ordering comes
from the DAG of `Parents:` headers per §4.6, not from branch topology.

### Infrastructure files

This `README.md` and the directory itself are infrastructure, not
protocol messages. They are NOT propagated by the merge phase; they
propagate via the branch participants forked from when joining the
transport. Edits to this README are coordinated out-of-band.

## Bootstrap roster

The transport is being bootstrapped by `ppx/main` (this branch). The
four bootstrap messages, all authored on `ppx/main`, are:

- `bafkreihhuejiefrqrm7zgw2jsdqc37lwmbvfkw5uqbnjx3wsobcxh3y7ni.txt`
  (m000: transport-creation; CID-named per spec §2; From:
  stevegt-via-perplexity)
- `bafkreihnonvsf3vmcagukqcxwoh35255eduulvwwx3kax6ty4iidklk5vu.txt`
  (m001: branch-binding clarification; cites m000 as parent; From:
  stevegt-via-perplexity)
- `bafkreidef4b4qdc4xjvkjrern7jm4ta75q55ed2u2ilwcrkxqhn7n4fjce.txt`
  (m002: CID-filenames + merge-cycle ratification; cites m001 as
  parent; From: stevegt-via-perplexity)
- `bafkreia46vxsahmeicugfxmc7natorkstc3mdaz4r5d3zz46whjwpvqwta.txt`
  (m003: second-sender mock per Alt-F2 of DF-021-TODO12.2; cites m002
  as parent; From: alice; closes 012.7 freeze-gate condition (2);
  ratifies the DF-38.5 directory rename)

Other developer agents joining the transport are expected to:

1. Fork their `<author-id>/main` from a branch that already has this
   directory (e.g. `ppx/main`).
2. Run the merge phase to confirm they observe the bootstrap messages.
3. Optionally post a message declaring their author-id and write
   branch, citing the most recent ratification message as a parent.

## Freeze checklist (per spec §Freeze gate)

- [ ] [`protocols/wire-lab.d/specs/transport-spec-draft.md`](../../../../../protocols/wire-lab.d/specs/transport-spec-draft.md) frozen
- [x] At least one round-trip exercising §3 / §4 / §4.6 / §6 / §7 (closed 2026-05-06: four-message DAG, two distinct senders `stevegt-via-perplexity` and `alice`, all CIDs verified by `tools/spec cid`)
- [ ] Steve signs `merge-group-transport-spec` promise
- [ ] `tools/spec freeze group-transport-spec` mints pCID and snapshots
- [ ] This directory and every message's grid envelope rewritten to the
      minted pCID

## Related

- [DR-009](../../../../../DR/DR-009-20260430-204108-group-transport-envelope.md)
- [TODO-bisur](../../../protocols/group-session.d/TODO/TODO-bisur-group-transport-envelope.md)
