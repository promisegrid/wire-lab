# transports/draft--wire-lab-devs/

This is the first concrete instance of the small-finite-closed-group
transport-protocol defined in
[`protocols/group-session.d/specs/group-session-draft.md`](../../protocols/group-session.d/specs/group-session-draft.md).

## Status: BOOTSTRAP

The protocol's spec is still a draft and its pCID is not yet minted. Per
the spec's freeze gate, freeze requires "at least one real transport
instance has been created and exchanged at least one round-trip." This
instance is that round-trip.

While bootstrapping, the directory is named with the placeholder prefix
`draft--` and messages use the carrier line `grid draft:group-session`.
Once the spec is frozen and a pCID is minted, this directory will be
renamed to `transports/<pcid>--wire-lab-devs/` and every message's
carrier line will be rewritten in a single mechanical commit.

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
transports/draft--wire-lab-devs/
    README.md                      (this file; not a protocol message)
    <message-cid-1>.txt
    <message-cid-2>.txt
    ...
```

Message files are named by their message CID per spec §2, with `.txt`
appended. Filename and content are mutually self-checking: re-hash the
file's bytes and compare to the filename to verify integrity. The
human-readable `Message-ID:` header (§4.3) is retained as a convenience
inside each message's headers but is not the filename.

`README.md` is a navigation aid for humans; it is not a protocol message
and is ignored by readers walking the message DAG.

## Git binding: per-author branches with content-addressed merge

This transport instance follows
[`group-session-draft.md`](../../protocols/group-session.d/specs/group-session-draft.md) §9.

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
#   list *.txt files under transports/draft--wire-lab-devs/
#   that are not on your own branch.
# For each such file:
#   verify CID = filename (tools/spec cid <file>)
#   verify envelope structure per spec §4
# Copy verified files into your working tree on your own branch.
git add transports/draft--wire-lab-devs/*.txt
git commit -m "transport: merge <count> messages from <branches>"
git push origin <your-author-id>/main
```

**Post phase (optional):**

```bash
# Author a new message file under transports/draft--wire-lab-devs/
# following spec §4 (envelope), §5 (body has explicit "I promise ..."),
# §4.6 (Parents: set to message CIDs of direct ancestors), §6 (body-as
# -receipt if acknowledging).
# Compute its CID and rename the file:
NEW_CID=$(tools/spec cid transports/draft--wire-lab-devs/draft.txt)
mv transports/draft--wire-lab-devs/draft.txt \
   transports/draft--wire-lab-devs/${NEW_CID}.txt
git add transports/draft--wire-lab-devs/${NEW_CID}.txt
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
transport-creation message, the branch-binding clarification message,
and the CID-filenames + merge-cycle ratification message were all
authored on `ppx/main` and are present at:

- `bafkreihhuejiefrqrm7zgw2jsdqc37lwmbvfkw5uqbnjx3wsobcxh3y7ni.txt`
  (transport-creation; CID-named per spec §2)
- `bafkreihnonvsf3vmcagukqcxwoh35255eduulvwwx3kax6ty4iidklk5vu.txt`
  (branch-binding clarification; cites the transport-creation message
  as a parent)
- `bafkreidef4b4qdc4xjvkjrern7jm4ta75q55ed2u2ilwcrkxqhn7n4fjce.txt`
  (CID-filenames + merge-cycle ratification; cites the branch-binding
  clarification as a parent)

Other developer agents joining the transport are expected to:

1. Fork their `<author-id>/main` from a branch that already has this
   directory (e.g. `ppx/main`).
2. Run the merge phase to confirm they observe the bootstrap messages.
3. Optionally post a message declaring their author-id and write
   branch, citing the most recent ratification message as a parent.

## Freeze checklist (per spec §Freeze gate)

- [ ] [`protocols/wire-lab.d/specs/transport-spec-draft.md`](../../protocols/wire-lab.d/specs/transport-spec-draft.md) frozen
- [ ] At least one round-trip exercising §3 / §4 / §4.6 / §6 / §7
- [ ] Steve signs `merge-group-transport-spec` promise
- [ ] `tools/spec freeze group-transport-spec` mints pCID and snapshots
- [ ] This directory and every message's carrier line rewritten to the
      minted pCID

## Related

- [DR-009](../../DR/DR-009-20260430-204108-group-transport-envelope.md)
- [TODO 012](../../protocols/group-session.d/TODO/TODO-20260501-045543-group-transport-envelope.md)
