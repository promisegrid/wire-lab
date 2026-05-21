# Question

How should PromiseGrid distinguish live-state protocol claims from durable audit
publication claims for real-time apps such as collaborative editors, presence,
multiplayer state, and telemetry? Source: `DI-ragaz`.

Open decision points:

- Is the guide-safe posture "live state stays off-grid for now, audit publishes
  to group-session" or a provisional future live protocol sketch?
- What does an audit message cite: snapshot CID, CRDT save blob, operation log,
  milestone receipt, or another object?
- How should conformance claims avoid conflating live channel behavior with
  group-session audit behavior?
