# Question

Which starting profile should the first PromiseGrid L6 CAS spec test: an
IPFS/IPLD-aligned profile, a promisebase-adapter profile, or a minimal
pointer/raw profile? Source: `DI-bukoh`; `DR-tumus`; `DI-mivap`.

Open decision points:

- How much IPFS / IPLD compatibility is useful before PromiseGrid has its own
  operational evidence?
- Which promisebase / pitbase mechanisms are strong substrate candidates rather
  than only prior-art examples?
- If the promisebase-adapter profile survives, which promisebase tree state is
  the adapter source: `main`, `kv`, a merged state, or no branch?
- How small can the first CAS profile be while still supporting pointer objects,
  sparse fetch, and turn-177's L5/L6/L7 inversion?
- Which starting profile best preserves future migration room if TE-43 later
  locks more complete chunked Merkle objects?
