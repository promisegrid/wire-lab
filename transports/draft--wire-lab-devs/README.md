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
their own clone of this repository. Concrete `From:` values are recorded
in the transport-creation message
(`m000-transport-creation-20260503T113024Z.txt`) and may be tightened in
a follow-up before freeze.

Under the per-author-branch git binding (see below), membership is
realized as the set of `<author-id>/main` branches participants are
configured to fetch.

## Layout (per spec §1: flat)

```
transports/draft--wire-lab-devs/
    README.md                              (this file; not a protocol message)
    m000-transport-creation-20260503T113024Z.txt
    <message-id-1>.txt
    <message-id-2>.txt
    ...
```

`README.md` is a navigation aid for humans; it is not a protocol message
and is ignored by readers walking the message DAG.

## Git binding: per-author branches

This transport instance uses the per-author-branch git binding described
in [`protocols/group-session.d/specs/group-session-draft.md`](../../protocols/group-session.d/specs/group-session-draft.md)
§9. In summary:

- Each participant has their own write branch named `<author-id>/main`.
  Examples: `ppx/main`, `codex/main`, `<another-id>/main`.
- A participant writes message files **only** on their own
  `<author-id>/main` branch.
- A participant **fetches all** participants' `<author-id>/main`
  branches and reads the union of their `*.txt` files under
  `transports/draft--wire-lab-devs/`.
- The git commit graph is incidental; ordering comes from the DAG of
  `Parents:` headers per §4.6, not from branch topology.
- `README.md` and the directory itself are infrastructure, not protocol
  messages; they propagate via the branch participants forked from when
  joining the transport.

## How to send

1. On your `<author-id>/main` branch (your own — never another
   participant's), pull from origin.
2. Author a message file under this directory following
   [`group-session-draft.md`](../../protocols/group-session.d/specs/group-session-draft.md)
   §4 (envelope), §5 (body must contain at least one explicit
   `I promise ...` clause), and §4.6 (`Parents:` set to the message CIDs
   of the messages your message acknowledges as direct ancestors).
3. Filename: `<message-id>.txt` per §2; the recommended `Message-ID`
   convention is `<author-id>-<utc-timestamp>-<short-slug>`.
4. Compute the message CID per §3 (CIDv1, base32, sha2-256, raw, over
   full canonical bytes) using `tools/spec cid <file>`.
5. Commit. Do not edit the file again after pushing (§7 append-only);
   editing would change the canonical bytes and break every `Parents:`
   reference.
6. Push your `<author-id>/main`. Never force-push.

## How to receive

1. `git fetch --all`
2. For each known participant's `<author-id>/main` branch, list new
   `*.txt` files under this directory.
3. Verify each new message's carrier line is `grid draft:group-session`
   (until freeze) and re-compute its CID to verify integrity.
4. Walk the `Parents:` DAG to reconstruct ordering.

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
