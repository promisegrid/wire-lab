# Question

Can PromiseGrid keep the universal outer envelope as `grid([pcid, payload])`
while a payload pCID defines a reserved proof slot plus a named
`payload_without_sig` projection, and then require the proof to cover the
canonical bytes of `grid([pcid, payload_without_sig])` without introducing a
universal outer signature tuple? Source: `DI-nohir`.
