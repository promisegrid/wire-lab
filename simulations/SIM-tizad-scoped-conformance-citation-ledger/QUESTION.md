# Question

How should PromiseGrid guide authors require scoped conformance claims and durable citation records so that early apps can be honest, auditable, and evolvable before final contracts freeze?

Open decision points:

- What minimum fields must a publishable conformance manifest contain so Bob can verify a partial claim from observed artifacts?
- Which contract or contract-family reference is authoritative at the protocol boundary when one logical app has multiple embodiments?
- When Bob returns a hash, what promise class is being made: content identity only, availability, authorization, or some app-specific combination?
- For live apps, what exact durable milestone object should audit publication cite: snapshot, save blob, op-log checkpoint, receipt, or break-witness?
- For replay-prone physical effects, when must a restarted agent emit a durable receipt versus a break-witness rather than re-executing?
- What minimal local promise-accounting records should Alice, Bob, and Carol retain so later auditors can interpret failures after hosts, drivers, or storage policies change?
