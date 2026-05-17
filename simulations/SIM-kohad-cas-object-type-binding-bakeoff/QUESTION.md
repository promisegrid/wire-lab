# Question

Should PromiseGrid bind CAS object type through CIDv1 codec values alone,
through CID codec plus object-internal kind fields, or through some other
mechanism that avoids filename and path suffix dependence? Source: `DI-bukoh`;
`DR-tumus`.

Open decision points:

- Can public or PromiseGrid-owned codec values distinguish raw chunks, Merkle
  nodes, pointer objects, and future object families without ambiguity?
- Does an internal kind field add useful evolvability, or does it duplicate the
  CID codec and create contradictory type sources?
- What breaks when type meaning lives in a path suffix instead of the CAS
  object identity?
- How should unknown object types be stored, rejected, or forwarded by peers
  that only understand the first L6 CAS profile?
