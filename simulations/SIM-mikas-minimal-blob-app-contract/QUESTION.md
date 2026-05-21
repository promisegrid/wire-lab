# Question

What minimal CAS-facing app contract can PromiseGrid teach for immutable blob
write/read behavior without overpromising replication, discovery, or
hash-as-capability semantics while `DR-tuhaz` and `DR-tumus` remain open? Source:
`DI-ragaz`.

Open decision points:

- Can the guide publish a provisional `blob in -> hash out; hash in -> blob out`
  pattern tied to a named draft CAS profile?
- What must be stated as host-defined: retention, replication, ingress,
  discovery, and authorization?
- Should "holding the hash means you can read" be rejected as base PromiseGrid
  semantics unless an app/spec explicitly promises that convention?
