# SIM-zazit: Chunk feed replication

This simulation explores the turn-177 inversion that L5 feeds replicate CAS
chunks between sites while L6 CAS stores and resolves those chunks for higher
layers. It treats feeds as meaning-oblivious replication mechanisms rather
than group-message transports. Source: `DI-navod`.

## Question

What does an L5 feed need to advertise, request, and replicate when its unit of
movement is a CAS chunk or CAS object rather than a group-session message?
Source: `DI-navod`.

## Turn 177 pressure

Turn 177 corrected the model from "CAS below feeds that carry messages" to
"feeds below CAS that move chunks." That has practical consequences:

- A feed should not need to understand group-session message semantics.
- Sparse sites should be able to advertise only the chunks they have or want.
- A single feed family may serve group messages, files, and future content that
  share the CAS layer.
- Pull/keep/advertise decisions need inputs from peer-local promise accounting
  records without making L5 itself a central accounting authority.

Turn 178 adds two more pressures to this simulation: every site should be sparse
by default, and some realistic experiments may put each simulated site or large
content corpus in a different repository rather than assuming all message bytes
live in wire-lab. Source: `DI-vaguf`.

This simulation owns the feed-side consequences of that inversion. The
`SIM-jomag` CAS object-model simulation owns object typing and chunking
parameters; this simulation asks how those objects move between sites. Source:
`DI-navod`.

## Decision axes

- **Advertisement unit:** individual chunk CID, Merkle root CID, pointer object
  CID, range/frontier summary, or a mix.
- **Request policy:** how a site decides which advertised chunks to pull under
  sparse-CAS assumptions.
- **Replication promise:** what a feed participant promises about timely
  storage, retransmission, integrity, and refusal.
- **Carrier fit:** how UDP, git, libp2p, IPFS, ATPROTO, or other substrates
  might carry the same feed role without becoming the protocol meaning.
- **Site topology:** whether a simulation site is a directory, a separate repo,
  a remote peer, or an imported fixture when testing sparse large-object flows.
- **Failure behavior:** duplicate advertisements, missing chunks, partial
  Merkle trees, corrupt chunks, and peer-specific refusal.

## Boundaries

This simulation is not the `udp-feed` spec and does not modify the existing
`SIM-ludaf-udp-feed` lineage. It is the turn-177 chunk-replication design-point
workspace that later feed specimens can adopt, reject, or compete against.
Source: `DI-navod`.
