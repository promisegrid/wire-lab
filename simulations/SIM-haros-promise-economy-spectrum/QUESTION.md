# Question

How can PromiseGrid test different promise-economy mechanisms without making
the wire format, group semantics, feed behavior, or CAS policy prematurely
commit to one economic model? Source: `DI-vabij`.

TODO owner: `protocols/wire-lab.d/TODO/TODO-rajig-promise-economy-spectrum.md`.
Source: `DI-pidag`.

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
- How should RFC-1005's content-addressable test-driven fabric -- test tree CID,
  executable tree CID, arguments, and cache-on-pass semantics -- map into or stay
  separate from PromiseGrid promise-economy vocabulary? Source: `DI-nulak`.
