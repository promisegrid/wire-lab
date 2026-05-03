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
(`m000-transport-creation-20260503T113024Z.msg`) and may be tightened
in a follow-up before freeze.

## Layout (per spec §1: flat)

```
transports/draft--wire-lab-devs/
    README.md                              (this file)
    m000-transport-creation-20260503T113024Z.msg
    <message-id-1>.msg
    <message-id-2>.msg
    ...
```

`README.md` is a navigation aid for humans; it is not a protocol message
and is ignored by readers walking the message DAG.

## How to send

1. Pull the latest from `ppx/main` (or whatever branch the transport
   currently lives on; during bootstrap that is the twig
   `ppx/te-20260503-112348-git-file-transport`).
2. Author a message file under this directory following
   `protocols/group-session.d/specs/group-session-draft.md` §4 (envelope),
   §5 (body must contain at least one explicit `I promise ...` clause),
   and §4.6 (`Parents:` set to the message CIDs of the messages your
   message acknowledges as direct ancestors).
3. Compute the message CID per §3 (CIDv1, base32, sha2-256, raw, over
   full canonical bytes).
4. Commit. Do not edit the file again after it appears in a commit on a
   branch any other participant has pulled (§7 append-only); editing
   would change the canonical bytes and break every `Parents:`
   reference.
5. Push.

## How to receive

1. Pull.
2. List `*.msg` files under this directory.
3. Verify each new message's carrier line is `grid draft:group-session`
   (until freeze) and re-compute its CID to verify integrity.
4. Walk the `Parents:` DAG to reconstruct ordering.

## Freeze checklist (per spec §Freeze gate)

- [ ] `protocols/wire-lab.d/specs/transport-spec-draft.md` frozen
- [ ] At least one round-trip exercising §3 / §4 / §4.6 / §6 / §7
- [ ] Steve signs `merge-group-transport-spec` promise
- [ ] `tools/spec freeze group-transport-spec` mints pCID and snapshots
- [ ] This directory and every message's carrier line rewritten to the
      minted pCID

## Related

- [DR-009](../../DR/DR-009-20260430-204108-group-transport-envelope.md)
- [TODO 012](../../protocols/group-session.d/TODO/TODO-20260501-045543-group-transport-envelope.md)
