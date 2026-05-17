# Question

How can PromiseGrid test different promise-economy mechanisms without making
the wire format, group semantics, feed behavior, or CAS policy prematurely
commit to one economic model? Source: `DI-vabij`.

Open decision points:

- What is the minimum protocol surface that supports peer-local promise
  assessment, barter-like reciprocal promises, capability tokens, transferable
  commitments, and floating-rate models without choosing among them?
- Which fields would accidentally bake in fungibility, token balances,
  universal pricing, or a central exchange assumption?
- How can a peer advertise, refuse, pull, keep, or forward content when its
  local economics model differs from a peer's model?
- Which experiments are needed before PromiseGrid allows transferable or
  permissioned promise tokens in any first-class protocol?
- What must remain out of scope until the simulations show that a proposed
  mechanism avoids cryptocurrency-like failure modes?
