# Question

What is the smallest CAS-backed group-session specimen that can preserve L7
group semantics while using L6 pointer objects, CAS roots, chunks, and L5 feed
replication for storage and movement? Source: `DI-navod`.

Open decision points:

- What is the group-visible identifier: message root CID, pointer-object CID,
  or another typed object?
- How do `Parents:`-style relationships translate when parents may be pointer
  objects or Merkle roots instead of inline text files?
- Which parts of current `group-session` remain semantic obligations, and
  which are artifacts of the `.txt` specimen?
- What additive migration artifact should appear first without rewriting
  historical message bytes?
