# TE-2: Trust-ledger merge after partition

*Thought experiment, part of the [PromiseGrid Wire Lab](../../harness-spec.md). This file is content-addressable; its hash is its pCID.*

Two clusters trade independently for simulated days, then reconnect. Each side has trust scores for peers from before the partition that have evolved differently. How do they merge? Does the lower score win, the higher, the more-witnessed? Outcome: a merge function with documented properties.
