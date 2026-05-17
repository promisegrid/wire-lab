# Question

What should count as a valid freeze or merge promise for a group-session
specimen, and how should that promise distinguish group-session semantics from
outer/feed rules? Source: `DI-pukap`; `TODO-bisur`.

Open decision points:

- Which artifacts must be named by pCID or path before a group-session specimen
  can be called frozen?
- Does freeze require implementation conformance tests, a verified message DAG,
  human-readable evidence, or all of those?
- What exact scope does a `merge-group-transport-spec` promise cover?
- How does the promise avoid treating a provisional group-session specimen as
  the canonical PromiseGrid wire format?
- What happens if outer/feed rules change after the group-session semantics are
  otherwise ready to freeze?
