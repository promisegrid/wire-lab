# Question

Which positional grid-envelope variant or variants, if any, survive comparison
across encoding, unknown-pCID behavior, and signature placement?

Secondary questions from the nested-vs-stacked turn-208 research that any
surviving grid-envelope variant must answer:

- What recursion depth budget should nested `grid([pcid, payload])` messages
  enforce?
- Is `pcid` a pure content hash, or does it carry version / routing metadata?
- How are capability references represented when the payload needs more than
  content references?
- What canonical serialization is used for signing and hashing grid messages?
- Does PromiseGrid need onion-routing layers, and if so what must the outer
  `pcid` reveal to routers?

Source: `DI-limom`, `DI-rugig`, `DI-fanah`, `DI-kabuk`.
