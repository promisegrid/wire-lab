# Question

Should PromiseGrid bind chunking through the object protocol pCID, a separate
chunking descriptor such as a possible chunking CID / cCID, a negotiated
profile, or a raw-only first profile? Source: `DI-bukoh`; `DR-tumus`.

Open decision points:

- If pCID defines chunking, does every content object need to know which
  protocol spec produced its chunks?
- If chunking has a separate descriptor, is that descriptor content-addressed,
  protocol-addressed, human-named, or something else?
- How do peers detect that Alice and Bob chunked the same bytes with different
  algorithms or parameters?
- Does the first L6 CAS profile need chunked Merkle roots now, or can raw chunks
  and pointer objects carry the first migration safely?
