# Question

Which L6 CAS object model should PromiseGrid test first for pointer objects,
raw chunks, Merkle nodes, deterministic CBOR / DAG-CBOR encoding, Rabin or
FastCDC chunking, CIDv1 codec-based object typing, and branch-specific
promisebase / pitbase prior-art evidence? Source: `DI-navod`; `DI-tibis`;
`DI-mivap`.

Open decision points:

- What exact CBOR profile produces stable CIDs across independently written
  implementations?
- Which CIDv1 codecs distinguish raw chunks, Merkle nodes, and pointer objects
  without duplicating type information in filenames?
- Which chunking algorithm and parameters are stable enough to become a
  protocol contract?
- Which promisebase / pitbase lessons become PromiseGrid design inputs, and
  which remain prototype history?
- If promisebase is used as prior-art evidence, which tree state is being
  evaluated: `main`, `kv`, a merged state, or no promisebase branch?
