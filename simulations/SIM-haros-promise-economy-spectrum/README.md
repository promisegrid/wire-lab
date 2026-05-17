# SIM-haros: Promise economy spectrum

This simulation explores turn 179's warning that "promise economy" might range
from peer-local social assessment to transferable capability-token marketplaces
with floating exchange rates, and that the cryptocurrency-like end of that
spectrum could go badly wrong. It is a standalone design-point simulation, not a
final economics mechanism and not a central trust ledger. Source: `DI-vabij`.
TODO owner: `protocols/wire-lab.d/TODO/TODO-rajig-promise-economy-spectrum.md`.
Source: `DI-pidag`.

## Question

How can PromiseGrid test different promise-economy mechanisms without making
the wire format, group semantics, feed behavior, or CAS policy prematurely
commit to one economic model? Source: `DI-vabij`.

## Turn 179 pressure

Turn 179 made the spectrum explicit: promises may be simple peer-relative
records, permissioned capability tokens, transferable commitments, or something
with floating exchange rates where each peer effectively issues its own units.
Steve also warned that this is only one possible model and may degenerate into
cryptocurrency toxicity if the protocol bakes in the wrong assumptions too
early.

This simulation exists so that mechanism neutrality is tested directly instead
of being hidden inside generic promise accounting records. Source: `DI-vabij`.

## Decision Axes

- **Mechanism neutrality:** what fields or object shapes would accidentally
  force fungibility, transferability, exchange rates, or pure-social accounting.
- **Local versus transferable:** which promises are only relationship records and
  which, if any, can move between peers.
- **Permissioning:** how a promise token or record expresses who may hold,
  transfer, redeem, or assess it.
- **Pathology resistance:** what scenarios reveal hoarding, speculation, Sybil
  amplification, laundering, or central exchange capture.
- **Layer boundary:** which economics choices belong in L7 applications, which
  are visible to L5/L6 pull and retention decisions, and which must stay outside
  the base protocol.

## Boundaries

This simulation does not choose a currency model, define a token standard,
settle reputation scoring, or make economic settlement part of the base
PromiseGrid protocol. It exists to keep the protocol-neutrality question
discoverable while later simulations and TODO owners test concrete variants.
Source: `DI-vabij`.
