# SIM-narok: Transport Family Bakeoff

This simulation captures the protocol/specimen questions parked in
`TODO-sinuv`: ring transport, cluster-of-clusters transport, gossip transport,
and receipts-at-scale. It is a standalone design-point simulation for future
transport-family pressure, not a final transport spec and not a shared protocol
bundle for other simulations. Source: `DI-pukap`.

## Question

Which future transport-family shapes should PromiseGrid explore after the
current UDP-feed and group-session specimens: token-ring, cluster-of-clusters,
gossip, receipt-vector scaling, or some combination that splits those concerns?
Source: `DI-pukap`; `TODO-sinuv`.

## Candidate Shapes

- **Token-ring transport:** Stronger ordering than gossip, explicit per-hop
  authorization, and bounded turn-taking.
- **Cluster-of-clusters transport:** Hierarchical addressing where a cluster is
  itself a first-class routing and policy unit.
- **Gossip transport:** Epidemic propagation with convergence and IHave/IWant
  generalization rather than strong ordering.
- **Receipts-at-scale extension:** Multi-writer or large-N acknowledgement
  shapes such as vectors, summaries, or compact proofs.

## Boundaries

This simulation does not schedule the future TEs or pick a transport family. It
keeps the parked protocol/specimen questions visible so later work can compare
them independently instead of treating the current UDP-feed or group-session
lineages as the only possible transport shapes. Source: `DI-pukap`.
