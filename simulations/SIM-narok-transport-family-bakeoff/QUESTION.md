# Question

How should PromiseGrid model future transport families beyond the current
lineages: token-ring ordering, cluster-of-clusters routing, gossip convergence,
and receipts-at-scale? Source: `DI-pukap`; `TODO-sinuv`.

Open decision points:

- Which semantics distinguish a ring transport from a centralized or
  hub-mediated transport?
- What does a cluster promise when the cluster, rather than a single peer, is the
  addressable transport unit?
- How do IHave/IWant-style claims generalize under gossip without implying
  delivery, ordering, or storage guarantees the transport cannot keep?
- Do receipts at multi-writer or large-N scale become vectors, summaries,
  per-peer promises, or a separate proof object?
- Which questions are independent enough to become separate sims or TEs later?
