# Question

Which L6 CAS object model should PromiseGrid test first for pointer objects,
raw chunks, Merkle nodes, deterministic CBOR / DAG-CBOR encoding, Rabin or
FastCDC chunking, and CIDv1 codec-based object typing? Source: `DI-navod`.

Open decision points:

- What exact CBOR profile produces stable CIDs across independently written
  implementations?
- Which CIDv1 codecs distinguish raw chunks, Merkle nodes, and pointer objects
  without duplicating type information in filenames?
- Which chunking algorithm and parameters are stable enough to become a
  protocol contract?
- Which promisebase / pitbase lessons become PromiseGrid design inputs, and
  which remain prototype history?
