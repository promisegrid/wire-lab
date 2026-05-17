# Question

Can a PromiseGrid application model BGP-class reachability, path selection, and
policy trust using peer-local promises and relationship-specific accounting,
without requiring a central route authority or global trust registry?
Source: `DI-tibis`.

Open decision points:

- What is the promise shape for route-like announcements, transit offers,
  refusals, withdrawals, and observed failures?
- Which parts of route selection belong in an L7 routing app, and which must be
  exposed by generic PromiseGrid feed, CAS, or promise-accounting surfaces?
- How should peers respond locally to route leaks, hijacks, stale paths, and
  conflicting policy claims?
- Can the model converge under sparse knowledge and intermittent connectivity
  without assuming a global table or central arbiter?
- What must remain out of scope until a real routing-app specimen exists?
