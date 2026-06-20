# kernel_transport_v1

`kernel_transport_v1` is the local transport-control protocol between a
container parser role and the container transport kernel. It is not an
application protocol and is not a registry.

The payload is a pCID-owned CBOR array using the pair-payload profile. The
promiser promises only local parser/kernel behavior:

- `receive_pcid`: the parser role promises that it can receive exact envelopes
  for a named application pCID and forward those bytes to local apps that have
  made matching receive promises.
- `carry_exact_envelope`: the parser role asks the transport kernel to carry
  exact envelope bytes toward a named target agent. The target name is
  transport-control data supplied by the parser role; the transport kernel does
  not decode the embedded application payload to discover it.

The transport kernel may inspect the outer control payload because this pCID is
explicitly kernel-handled. For embedded application envelopes, the kernel may
inspect only the grid tag, slot-0 pCID, parent links, exact message hash, and
proof validity. Application destination, operation, route, promise body, and
business semantics remain owned by the parser role and app protocols.
