# Question

If L5 feeds move CAS chunks instead of group messages, what must a feed
advertise, request, replicate, and promise so sparse sites can converge without
assuming every site stores all CAS objects? Source: `DI-navod`.

Open decision points:

- Does a feed advertise leaves, roots, pointer objects, or compact frontiers?
- Which layer decides that an advertised chunk is worth pulling or retaining?
- What feed behavior remains substrate-neutral across UDP, git, libp2p, IPFS,
  ATPROTO, and future carriers?
- When do sparse sites need separate repos or external content corpora rather
  than directories inside one wire-lab checkout?
- How are corrupt, missing, duplicate, or refused chunks represented without
  leaking group-session semantics into L5?
