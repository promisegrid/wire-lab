# SIM-punaz: BGP-class routing app

This simulation explores turn 178's question about whether a PromiseGrid
application could model BGP-class inter-domain routing policy without inheriting
BGP's central trust assumptions, route-leak failure modes, or free hijack
dynamics. It is a standalone application-pressure simulation, not a base
PromiseGrid routing spec and not a general promise-accounting record store.
Source: `DI-tibis`.

## Question

Can a PromiseGrid application model BGP-class reachability, path selection, and
policy trust using peer-local promises and relationship-specific accounting,
without requiring a central route authority or global trust registry?
Source: `DI-tibis`.

## Turn 178 pressure

Turn 178 raised BGP replacement as an application-class question: could a
PromiseGrid app get away from BGP vulnerabilities by making route announcements,
transit promises, policy refusals, and observed failures accountable in the same
peer-relative promise economy used elsewhere in PromiseGrid?

This simulation is split out from `SIM-rusap-promise-accounting-records` because
BGP-class routing is more than one accounting example. It combines graph
selection, path policy, adversarial announcements, convergence, partial
knowledge, and long-lived relationship records. Source: `DI-tibis`.

## Decision Axes

- **Announcement shape:** whether route advertisements are promises, claims,
  signed observations, or a composition of those forms.
- **Policy expression:** how peers express reachability, cost, preference,
  refusal, jurisdiction, and transit constraints.
- **Local decision inputs:** what Alice records locally about Bob, Carol, and
  Mallory before preferring one path over another.
- **Failure handling:** how route leaks, hijacks, partitions, stale paths, and
  false reachability affect future decisions.
- **Layer boundary:** what belongs in an L7 routing app versus generic L5/L6
  PromiseGrid transport and CAS behavior.

## Boundaries

This simulation does not choose a replacement for Internet BGP, does not define
a base PromiseGrid routing protocol, and does not require the core PromiseGrid
protocol to adopt any global route table. It exists so BGP-class pressure can
evolve independently and inform later app-level design decisions.
Source: `DI-tibis`.
