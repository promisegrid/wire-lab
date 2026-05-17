# SIM-jurar: CAS-backed group session

This simulation explores a CAS-backed successor shape for group-session
messages. It asks how L7 group semantics should point at L6 CAS roots and
pointer objects after turn 177, without rewriting the existing `.txt`
group-session specimen or declaring one final PromiseGrid message format.
Source: `DI-navod`.

## Question

How should a group-session specimen behave when group messages are resolved
through CAS pointer objects and roots rather than stored directly as inline
`<cid>.txt` message files? Source: `DI-navod`.

## Turn 177 pressure

Turn 177 preserved the idea that groups define message semantics while moving
storage and replication beneath them:

- L7 group protocols define the meaning of messages, parents, acknowledgements,
  membership, and application payloads.
- L6 CAS stores and resolves pointer objects, chunks, and Merkle roots.
- L5 feeds advertise/request/replicate chunks between sites.
- Pointer files are content-addressed objects in their own right, not symlinks.
- CBOR message or pointer bytes must be canonical enough for CID agreement.

This simulation keeps those group-session consequences visible while
`SIM-rakot-group-session` remains the historical/current `.txt` draft lineage
and `TODO-pipus` owns operational migration from pre-CAS inline specimens.
Source: `DI-navod`.

## Decision axes

- **Group root:** what object a group member publishes as the current root or
  frontier of a group view.
- **Parent links:** whether parent references are direct message-root CIDs,
  pointer-object CIDs, Merkle roots, or a typed combination.
- **Body shape:** whether message bodies are CBOR text strings, CBOR maps,
  encrypted blobs, signed payloads, or arbitrary CAS roots.
- **Migration:** how existing `.txt` group-session evidence can remain
  historical while an additive CAS-backed specimen appears beside it.
- **Envelope relation:** how a candidate grid envelope wraps or points to the
  CAS-backed message object without making any grid-envelope variant canonical.

## Boundaries

This simulation does not supersede `SIM-rakot-group-session` by itself. It is a
standalone successor specimen charter: later work can instantiate files,
fixtures, or protocol drafts here if the CAS-backed group-session line proves
useful. Source: `DI-navod`.
