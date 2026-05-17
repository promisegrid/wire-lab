# SIM-zarud: Conditional Release / Geofencing

This simulation captures the protocol/specimen question from `TODO-ralud`:
conditional-release promises such as "send only if the recipient promises
onward-restraint," geofencing constraints, and the recursive promise graph that
results when those conditions propagate beyond the first recipient. It is a
standalone design-point simulation, not a frozen PromiseGrid protocol and not a
shared protocol bundle for other simulations. Source: `DI-pukap`.

## Question

Which layer or protocol family should own conditional release, onward-restraint,
geofencing, and recursive promise-graph tracking when Alice sends content to Bob
under conditions that Bob must preserve when forwarding to Carol? Source:
`DI-pukap`; `TODO-ralud`.

## Candidate Shapes

- **Group-session-local:** Treat conditional release as a group-membership and
  dispatch rule enforced by the session semantics that deliver the content.
- **Separate conditional-release family:** Model conditional release as its own
  promise protocol that group/session protocols can invoke or cite.
- **Transport/feed-visible constraint:** Push enough condition metadata down so
  routing, replication, or fetch decisions can avoid violating geofence or
  onward-restraint promises.
- **Hybrid graph:** Keep human-facing semantics at the group layer while lower
  layers carry opaque condition references and proof hooks.

## Boundaries

This simulation does not choose the owner layer, define the final promise graph
object model, or modify the existing group-session specimen. It exists so
`TODO-ralud` can evaluate the architectural pressure in simulation space before
opening or answering a later TE/DR. Source: `DI-pukap`.
