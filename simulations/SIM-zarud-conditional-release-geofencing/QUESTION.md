# Question

When Alice sends content to Bob with conditions such as "only forward if the
next recipient promises onward-restraint" or "only dispatch inside this
geofence," what PromiseGrid protocol layer owns the recursive promise graph and
what evidence must travel with the content? Source: `DI-pukap`; `TODO-ralud`.

Open decision points:

- Does conditional release belong inside group-session semantics, in a separate
  conditional-release protocol family, in transport/feed metadata, or in a
  hybrid split?
- What counts as sufficient evidence that Bob accepted Alice's onward-restraint
  condition before Bob forwards content to Carol?
- Does geofencing constrain membership, message dispatch, storage, fetch,
  replication, or all of those operations?
- Can lower layers safely store or route encrypted/opaque content when they do
  not understand the condition vocabulary?
- How does Mallory's attempted replay outside the allowed group or geography get
  detected, refused, or made auditable?
