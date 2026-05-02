# PromiseGrid Wire Lab

*This is the canonical Wire Lab harness specification. Thought experiments are kept as separate content-addressable files in [`docs/thought-experiments/`](docs/thought-experiments/). The repo-level [`README.md`](README.md) gives a one-page orientation.*

A simulation harness for discovering the wire formats, ingress models, and trust mechanics that let PromiseGrid survive and evolve as an open, decentralized community of free agents and humans across multiple human generations. v2 supersedes v1 in five places: the message is now a stack of promises rather than a fixed envelope; trust is a first-class durable relationship rather than a per-message check; capability tokens may operate as personal currencies in a continuous double auction; "kernel" is just a role that any agent can play; and not every agent simulates a deterministic program — some simulate humans and other non-deterministic actors.

The harness exists to let us **discover** the right design, not to validate a predetermined one. Every choice in this document is an experimental knob the simulator can change between runs.

> **Read this document as provisional throughout.** Concrete shapes appear inline — struct definitions, schemas, directory layouts, taxonomies, numeric thresholds. They are *examples of one workable starting point*, not commitments. The only structural commitment we are confident in this early is the one in §10a.8: the canonical harness-spec pointer follows the user's signing key. Everything else — the `Promise` field set, the `TrustLedger` shape, the K1–K5 ingress taxonomy, the seven C-scenarios, the 30%-at-month-12 pass line, the grid-poc directory table, the proposal shape, the convergence formula — is presented in concrete form to make the design discussable, but is expected to evolve. §12 calls out which of these are most likely over-locked and asks which to revisit next.

---

## The 100-Year Goal

PromiseGrid is designed against an explicit **100-year goal**: it must survive and remain useful as an open, decentralized community of free agents and humans across roughly a century of evolution, without depending on any central authority that might disappear inside that horizon. Every design choice in this harness-spec, in subordinate protocol specs (`specs/transport-spec-draft.md`, `specs/group-transport-draft.md`, future protocol specs), and in every TE/DR/DI/TODO record must be pressure-tested against the constraints below. Short-horizon assumptions that look fine over months or years routinely fail at the 100-year horizon; the constraint set exists to prevent them from being smuggled in unnoticed.

The constraint set:

- **C-1 — No central registry.** No global registry exists at the 100-year horizon for any naming or coordination function. Any centralized registry that exists today (DNS, ICANN, IANA, GitHub, npm, the PromiseGrid project itself, this wire-lab repo) is presumed gone or untrustworthy at some point inside the horizon. Naming, discovery, and trust must work without one.
- **C-2 — Multi-generational durability.** The design must support handoff across multiple human generations. Signing keys rotate, retire, get inherited, get lost, get compromised. Communities of interest fragment, merge, fork, fall silent for decades, return. Artifacts must be self-explanatory to a contributor who arrives 30 years later with no living mentor.
- **C-3 — Adversarial-by-default.** The design must assume adversaries are present at every layer at all times. No mode where "we trust each other so this part is fine" is acceptable, even temporarily. Every layer must be auditable in isolation, every promise must carry its own signature chain, every receipt must be verifiable against locally-held trust ledgers.
- **C-4 — Protocol forking is normal.** Over a century, protocols will be forked into bilingual variants, contested, declared deprecated by sub-communities, rediscovered, adapted to new substrates. Forking is a normal life-cycle event, not an exceptional one. A pCID identifies a *spec*, not the *use* of a spec; multiple competing pCIDs for analogous protocols is the steady state.
- **C-5 — Trust accrues per-burden.** Trust is not a global scalar attached to an identity. It is a per-promise-type, per-relationship, per-context vector that each agent maintains in its own local trust ledger. Any design that produces or relies on a global trust score fails this constraint.
- **C-6 — Signing key is the only structural lock.** The only durable anchor in the system is the signing key of an individual agent. Specs follow signing keys. Trust ledgers follow signing keys. Histories follow signing keys. Every other "fixed point" — repository URLs, registry names, organizational identities, project websites — is convention, and conventions evolve, decay, or get attacked over the horizon.

The constraints are filters, not goals. A design that satisfies all of them may still be wrong; a design that violates any of them is wrong. Subsequent TEs that pressure-test alternatives reference these constraints by code (C-1 through C-6) so the wire-lab has a shared shorthand and each TE does not re-derive the constraint set.

The full pressure-test of every currently-locked artifact in the wire-lab against these six constraints lives in [TE-28: The 100-year goal as a load-bearing design constraint](docs/thought-experiments/TE-20260501-202713-100-year-goal-as-design-constraint.md). Five open questions surfaced from that pressure-test (OQ-100.1 through OQ-100.5) are recorded in §12 below.

---

## 1. The Core Reframe: Messages Are Layered Promises

A message is no longer `[pCID, payload, signature]`. It is a **promise stack** — an ordered set of nested promises, each made by a distinct promiser, each interpretable on its own and verifiable against that promiser's history.

> **A note on what a pCID is.** A `pCID` is a *protocol* CID: the content hash of a spec document that defines a wire protocol. It is roughly equivalent to a TCP/UDP port number — an identifier that selects a protocol — except that no central registry is needed, because the spec's hash *is* the port number. Anyone can mint a pCID by writing a spec and publishing its CID. A pCID is **not** the hash of any particular payload; it is the hash of the rules that the payload follows. Two different messages can carry the same pCID and entirely different payloads, just as two TCP segments to port 443 can carry entirely different HTTPS streams. This distinction matters for the assertion vocabulary below.

### 1.1 The promise-stack model

A message is `[]Promise`. A `Promise` is:

```
Promise := {
    promiser:  AgentID | TransportID | RuntimeID | KernelID
    assertion: "this payload conforms to the spec identified by pCID P"
              | "I authored this payload"        ; signature-style
              | "I forwarded this byte sequence unmodified from peer Q at time T"
              | "I will perform action X if you redeem this token"
              | "I commit not to issue another message numbered N for this stream"
              | "the payload below resolves at the CID I'm naming; fetch it elsewhere"
              | ... (open set; new assertion vocabularies are themselves content-addressed specs, i.e. their own pCIDs)
    body:      bytes | promise-stack | nil
    evidence:  signature | MAC | "byte-arrival on authenticated channel" | nil
    ttl:       optional duration or revocation hook
}
```

The on-the-wire encoding of a message is a CBOR array of promise frames, innermost first or outermost first (the simulator tests both). A receiver consumes promises top-down, deciding for each whether to **accept**, **defer** (cache pending more evidence), or **reject**.

Some realizations of the same logical message in this model:

- **`[pCID, payload, signature]`** — pCID names the protocol the payload conforms to; signature frame asserts authorship. Two promises, one outermost frame.
- **`[pCID, payload]` + transport promise** — sender promises the payload conforms to the spec at pCID; the transport (TLS, Noise, an authenticated channel kept alive across many messages) carries a continuing implicit promise of sender identity. Signature lives **inside** payload only when the receiver demands non-repudiable evidence.
- **`[transport-promise, [routing-promise, [authorship-promise, [content-promise, payload]]]]`** — a fully nested form where the kernel/router has appended its own promise after receiving from the network, so a downstream handler can verify *which router forwarded this and when*.
- **`[bare-payload]`** — within a single trusted process, the only promise is implicit ("the function caller asserts this argument is correct"). Still a promise stack of length one.

The wire library (call it `promstack`) then has just three operations:
1. `Wrap(msg, promise) → msg'` — push a new promise on the outside.
2. `Peel(msg) → (outermost-promise, inner-msg)` — consume the top promise.
3. `Project(msg, predicate) → []Promise` — pull out all promises matching a predicate (e.g. "all signature promises by Alice").

### 1.2 Why this is better than a fixed envelope

- **Transport-as-promiser becomes legitimate.** A TCP connection authenticated at startup makes a continuing implicit promise about the source of every byte. We don't need to re-sign every message; we can track that the transport is *making* a multi-message promise, and the trust ledger updates if the transport keeps or breaks it.
- **Routers and kernels can append their own promises** without modifying the inner content — exactly the chain-of-custody property we want for a centuries-long system where intermediaries come and go.
- **The signature debate goes away.** Whether the signature covers `[pCID, payload]` or just `payload` becomes a choice of *which promise frame holds the signature*. Different protocols (different pCIDs) can choose differently and the wire format doesn't care.
- **Capability tokens are just promises too.** A token is a promise whose assertion is of the form "I will perform X if you redeem this." It lives in the same stack as everything else.

### 1.3 What the simulator tests about layering

- That a receiver can correctly handle **out-of-order promise stacks** (some protocols want signature outermost, some want it innermost; a long-lived system needs both).
- That a forwarding node can **strip its own promise on the way out** (so I don't accumulate a 200-promise tail after 100 hops over 50 years).
- That a promise **about a missing inner body** (e.g. "I attest the CID below resolves to legal payload but I'm not shipping it; fetch it") works as a first-class case. This is the "merkle reference" pattern we'll need for huge messages and for the hypergraph.
- That two agents who disagree about which promise frame to evaluate first **fail loudly and recoverably**, not silently.

---

## 2. Trust as Durable Relationship: The Keep/Break Ledger

Steve called this out as key: peers know each other and build trust over time in durable relationships. The harness must make this a first-class subject of experimentation.

### 2.1 Per-agent trust ledger

Every agent maintains, for every other promiser it has ever interacted with (peers, signers, transports, kernels, named protocols), a `TrustLedger` entry:

```
TrustLedger[promiserID] := {
    first_seen_ns:        int64
    interactions:         counter
    kept:                 counter
    broken:               counter
    evidence_chain:       []EvidenceRef    // hashes of observed keep/break events
    open_promises:        []OpenPromise    // promises in flight, not yet decided
    score:                float            // current local trust scalar
    score_components:     map[assertion_type]float // separable by what kind of promise
    reputation_imports:   []ReputationImport      // 3rd-party attestations (down-weighted)
    relationship_age_ns:  int64
    last_drift_ns:        int64
}
```

The harness lets us swap the **scoring function** between runs without changing agent code. Three families to compare:

| Family | Update rule | Pros / Cons |
|---|---|---|
| **Beta-Bayesian** | `score = (kept+α)/(kept+broken+α+β)`, conjugate prior | Smooth, principled, well-studied; requires choosing prior; flat memory of distant past |
| **Worldline-weighted** (current sim4 sketch) | `trust *= weight × (actual − predicted)` along a chain | Captures hypergraph branch quality; goes negative; hard to interpret as probability |
| **EigenTrust-like** (P2P literature) | global fixed-point over local ratings | Strong against simple sybils; centralizing pull; fragile against collusion |

A run produces a transcript that includes every trust-ledger update, so we can replay a different scoring function on the same byte-level history and compare.

### 2.2 The "broken promise" event is itself a promise

When agent B observes that A failed to keep promise X (e.g. issued a token, never responded to redemption), B emits a *witnessed-break* promise: "I, B, observed A break promise X at time T, evidence is the redemption-message-CID I sent and the lack of response by deadline D." This break-witness travels the network as ordinary messages.

Important: *the receiver of a break-witness must apply its own trust ledger to B before believing B*. This is how slander attacks and Sybil-by-witnessing collapse: a low-trust peer's witnessing carries low weight. This is also what gives the system its multi-generational durability — old break-witnesses signed by long-dead agents still carry weight if their signing keys remain in the chain of custody, but only as much as their reputation at the time of signing earned.

### 2.3 Durable relationship features the harness supports

- **First-encounter rituals.** When two agents meet for the first time, both record a `first_seen` event. Several scenarios test what these rituals must contain (mutual challenge? exchange of references? trial transaction?).
- **Reputation portability with provenance.** When agent C imports a reputation score for X from agent A, the import is itself a promise by A about X. C trusts the import only as much as it trusts A, and only for the assertion types A vouched for.
- **Trust decay vs. trust persistence.** Knob: do scores decay over absent-time? Two scenarios will show the difference — one where I haven't talked to my old friend Alice in 30 years (real years of simulated time) but our shared history still counts, vs. one where idle relationships fade. Probably we need *separate* decay rates per assertion type.
- **Defection cost.** What does it cost to walk away from a relationship and start over with a new identity? Sybil resistance lives here. The harness can vary the cost (free, computational, social-graph-attested, stake-based) and observe which costs deter selfish exit.
- **Multi-generational handoff.** Scenarios where an agent's signing key is rotated, retired to a successor, or where a long-running agent dies and "wills" its reputation to another. The wire-level question: what kind of promise does that require, and who must witness it?

### 2.4 Trust is per-assertion-type

A given peer may be highly trusted to deliver bytes (transport assertion), mildly trusted to attest authorship (signature assertion), and totally untrusted to issue capability tokens that bind real resources (issuer assertion). Lumping these into one scalar destroys information that matters in a long-lived system. Every entry in the ledger therefore tracks `score_components` keyed by assertion type.

---

## 3. Capability Tokens as Personal Currency: A Continuous Double Auction

This is one of several economic-coordination experiments the harness will run, but the most fully developed one because the prior sim4 work already gestures at it.

### 3.1 The model

- Every agent **issues** its own currency by emitting capability tokens. A token says: "I, Alice, promise to perform action `X` (maybe parameterized by the redeemer's input) when this token is redeemed at me."
- Each token has a **face description** (the action class) and a **face quantity** (how much of it).
- Agents accumulate tokens issued by various peers: a portfolio of IOUs.
- Markets emerge. Bob holds 10 of Alice's "compute-1-vCPU-hour" tokens and 5 of Charlie's "store-1-GB-month" tokens, and posts a bid: "I'll trade 3 Alice-tokens for 2 Charlie-tokens." Charlie or anyone holding both can match.
- A **continuous double auction** at each agent (or at a designated market-maker handler) clears these orders and produces a stream of transaction hyperedges.

### 3.2 Exchange rates as trust signals

Crucially: the equilibrium rate at which Alice-tokens trade for Charlie-tokens is **the market's verdict on relative trust × utility**. If Alice's tokens trade at half the rate they did a month ago, the network is telling us either Alice's reliability has dropped or her offered service has fallen out of favor. The harness records these rate paths as a continuous experimental signal.

A key insight: this gives PromiseGrid a **quantitative, decentralized, continuously-updated trust metric** that doesn't require a global reputation algorithm. It is emergent from voluntary trade.

### 3.3 What the simulator must verify

- That the wire-level encoding of orders, fills, and redemptions composes from the same `Promise` primitive (no special "market protocol" envelope).
- That double-spend prevention works without a global ledger: Alice can issue 100 tokens, but if she issues 200 by speculating on parallel hypergraph branches, the merging mechanism must catch it and the witnesses must be able to publish a break-witness.
- That **price discovery** is real. Run two scenarios — one where every peer has perfect information and one where information is patchy — and confirm that prices converge to the same equilibrium in the first, and to clusters in the second.
- That an **adversary cannot pump-and-dump** by issuing many cheap tokens, building reputation, then suddenly issuing many expensive tokens and defaulting. (The honest-market scenario should pass; the adversarial scenario should produce visible defaults that the trust ledger records.)
- That the system **does not require a single market**. Many local markets (per-handler, per-region, per-community) should be able to run in parallel and bridge through arbitrageurs.

### 3.4 Alternative economic models the harness will also try

The double auction is one experiment. Run it against:

1. **Gift economy.** No prices; agents fulfill requests when they feel like it; reputation is the only signal; resource limits are soft. Test in a small high-trust community.
2. **Mutual credit / Ripple-style.** Bilateral credit limits between known peers; balances rebalance through chains. Works well in dense graphs of strong relationships.
3. **Stake-bonded promises.** To issue a token, you must lock collateral. Default burns it. Familiar from blockchains.
4. **Quadratic funding for shared resources.** Communities pool contributions for common-pool resources (a shared cache, a shared time-server). Tests whether commons can be funded without coercion.
5. **Auctioning the right to host a kernel-role.** Periodic auctions for "I will be the router of pCID X for the next epoch." Tests rotation of privileged roles.

Each is a different scenario set in the harness. The point is not to pick the winner; the point is that **the wire format and trust ledger must support all of them** because a centuries-long open system will see all of them tried, often simultaneously, by different sub-communities.

---

## 4. Ingress Models: Many Kernels, or No Kernel

You said it: "we may not even need a dedicated kernel." The harness must explore the full spectrum.

### 4.1 The four canonical ingress models (and a fifth hybrid)

| Model | Description | Best for |
|---|---|---|
| **K1 — Classical microkernel** | Every inbound message hits a kernel agent first; kernel checks pCID, finds the registered handler, forwards. Kernel can deny. | Servers running mixed workloads; nodes that need a single audit chokepoint |
| **K2 — Router-only** | Kernel exists but only routes; cannot deny, cannot inspect payload. Pure dispatcher. Handlers enforce their own policy. | Network appliances; edge nodes; constrained MCUs |
| **K3 — Handler-as-kernel** | The handler that "owns" a resource is the kernel for that resource. Other handlers on the same node go through it for that resource only. | Heterogeneous nodes where different resources have different governance |
| **K4 — No kernel** | Each handler binds its own ingress (its own socket, file watcher, etc.). Discovery is via gossip, not a kernel registry. Conflicts (two handlers claiming the same port) are resolved by external signal or first-come. | Browser tabs; one-shot agents; air-gapped devices |
| **K5 — Capability-only** | The "kernel" is just the holder of root capability tokens for that node's resources. To do anything, you need a capability that ultimately chains back to the root holder. The root holder doesn't *do* anything except issue / revoke tokens. | Long-lived community-shared infrastructure where governance matters more than mechanism |

Per runtime type, default models will likely differ. The harness lets us mix freely: a node can run K1 for its container handlers while running K4 for its browser-tab handler.

### 4.2 What the simulator forces us to discover

- **Whether the wire format changes between models.** It should not. A handler that started its life under K1 should be able to migrate to K4 by re-deploying without rewriting its message-handling code. The harness will fail any scenario where this is not true.
- **Where policy lives.** In K1 it's centralized; in K5 it's in the capability graph; in K4 it's in each handler. The harness measures, per scenario, *how many places must agree* for a sensitive action to succeed. Fewer is more durable across centuries (because places that need to agree are places that can disappear).
- **Cost of ingress.** Per-message latency and CPU under each model. Real numbers in the transcript, not arguments in markdown.
- **What happens when the kernel itself becomes adversarial.** K1 and K3 have a privileged role; if compromised or corrupted, they can deny all service. Scenario: kernel goes Byzantine. Does the system route around it? If yes, K1 collapses to K4. The harness measures the recovery time and the messages-lost.

### 4.3 Hybrid scenarios

The most interesting scenarios are mixed-model. Example:

> A node runs three handlers: a WASM module under K1 (kernel-mediated, denied by default), a native Go process under K5 (capability-only), and a browser tab under K4 (its own websocket). The same external peer sends three messages targeting the three handlers in close succession. How does each ingress feel from the outside? Can the peer write one message-sending function for all three?

The answer must be **yes, identically** at the wire level. The differences must live entirely in **which promises are required for acceptance** at each ingress, not in how messages are formatted or framed. If they're not yes-identical, that's a wire-protocol bug to fix.

---

## 5. Agents That Aren't Hardcoded Programs

You corrected v1 here: not every agent is deterministic. Some simulate humans, node administrators, irate users, capricious adversaries. The harness must support both classes cleanly.

### 5.1 The two agent classes

```
Agent := DeterministicAgent | StochasticAgent

DeterministicAgent: pure state machine, seeded RNG only, repeatable runs (v1 contract)
StochasticAgent:    may consume real entropy, real wall-clock time, real network;
                    may be backed by an LLM, by recorded human input, or by
                    a probabilistic model
```

A run that contains stochastic agents is **not byte-deterministic**, but it is still *reproducible-in-distribution* if the harness controls the seeds and the LLM calls.

### 5.2 The stochastic agent profiles we need

| Profile | Models |
|---|---|
| `user-novice` | A casual end-user; clicks before reading; reasonable patience; trusts defaults. |
| `user-expert` | A power user; reads docs; tries things deliberately; reports bugs. |
| `node-admin` | Operates a node; cares about uptime, costs, reputation; will rotate keys. |
| `community-elder` | Long-lived participant; attends governance discussions; pushes back on proposals. |
| `casual-defector` | Mostly cooperates, defects when easy gains are available. |
| `hostile-prober` | Looks for protocol bugs; sends malformed messages; tests Byzantine paths. |
| `griefer` | Defects to harm others, even at personal cost. |
| `lazy-validator` | Approves messages without checking; common, quietly damaging. |
| `regulator` | Issues high-trust attestations; can revoke; rare; powerful. |

Each profile is implemented as a small policy module that decides actions based on the agent's perceived state. Implementations may use LLMs to play the role; the harness records every LLM prompt and response so a run can be replayed deterministically by feeding back the recorded outputs (the "transcript-replay" mode).

### 5.3 Scenario generators powered by stochastic agents

- **Onboarding scenarios.** Drop a `user-novice` into an existing community and see how many messages it takes for them to do something useful. Measures the protocol's accessibility, not just its correctness.
- **Governance scenarios.** A community of `community-elder` + `user-expert` debates a protocol upgrade (a new pCID), votes, and rolls it out. Measures the protocol's evolvability.
- **Long-running scenarios.** Mix all profiles, run for the equivalent of simulated months. Measure trust-ledger drift, market price stability, kernel-role rotation, and rate of break-witnesses. Detect tragedy-of-the-commons emergence.
- **Generation-handoff scenarios.** A `community-elder` retires, hands keys to a successor. Does their reputation transfer? Should it? At what discount?

### 5.4 Crucially: stochastic agents talk the same wire as deterministic ones

The harness does not give stochastic agents a special channel. They produce real `Promise` stacks, on real edges, and their messages get the same trust-ledger treatment. This is the only way to test that the protocol survives contact with humans (and human-like agents) at all.

---

## 6. Tragedy of the Commons: The Acceptance Test for the Whole System

You said the goal is to discover the design that **avoids tragedy of the commons** in a long-lived multi-community system. That means the harness needs explicit ToC scenarios with explicit metrics.

### 6.1 The ToC scenarios

| Scenario | The commons being preyed upon | Default-defection move |
|---|---|---|
| **C1 — Bandwidth flood** | Shared transport capacity | Spam high-bandwidth pCIDs; let neighbors pay relay cost |
| **C2 — Storage parasite** | Shared content-addressed cache | Pin lots of data; serve none to others |
| **C3 — Compute leech** | A community of capability-token issuers | Redeem freely; never issue; collect refunds via market arbitrage |
| **C4 — Reputation harvester** | The trust ledger itself | Build trust on cheap promises, then issue one big bad promise and exit |
| **C5 — Naming squatter** | The pCID space (handler descriptors) | Register popular names; sell back; never serve |
| **C6 — Governance capture** | A community's protocol-change process | Vote-stuff; push self-serving pCID upgrade; lock out dissenters |
| **C7 — Liability dodge** | The break-witness system | Default; deny; counter-witness; collude with allies to drown out witnesses |

Each scenario instantiates a population — typically 20-100 agents, mostly cooperators with a tunable fraction of defectors of each profile — and runs for simulated months. The pass/fail line is **whether cooperators continue to find each other and transact** while defectors are progressively excluded.

### 6.2 The metrics

- **Survival of cooperator-cooperator transactions** over time. (Should not collapse.)
- **Marginal cost of defection.** A defector should pay a rising price to keep defecting (lost trust → worse market terms → fewer counterparties).
- **Sybil cost.** Cost to spawn a new identity and re-enter the cooperator set. (Should be high enough to deter casual reset.)
- **Recovery from a single bad actor.** How many messages until a single griefer is functionally muted by the network?
- **Recovery from a coordinated 20% adversary.** Same, but harder.
- **Drift of the protocol itself** — which pCIDs gain adoption, which die, who proposed which.

The harness writes an `ostrom-report.md` per long run, named after Elinor Ostrom's design principles for managing commons. The report scores the run against each principle (clear boundaries, proportional benefits, collective-choice arenas, monitoring, graduated sanctions, conflict resolution, recognized rights to organize, nested enterprises). Designs that score low on multiple principles will collapse over centuries even if they pass all single-message tests.

### 6.3 What "passes" looks like

A design passes if, in **at least three independent runs of all seven C-scenarios**, with different random seeds and different agent-profile mixes, the cooperator population's transaction rate does not fall by more than 30% from baseline at simulated month 12.

Passing C1-C7 is necessary but not sufficient. It's the operational definition of "this system can survive its own users for a year." Designs that pass this can graduate to multi-decade scenarios.

---

## 7. The Harness Itself: Architecture Refinements

Most of v1's three-layer design (scenario driver, edge fabric, agents) survives. Here are the v2 refinements.

### 7.1 Edges carry promises, not bytes

Every edge in Layer B annotates inbound traffic with a **transport-promise** describing what it observed: source endpoint, time of arrival, integrity of the byte stream (was it TLS? was it Noise? was it raw TCP?), liveness of the upstream peer. This promise becomes the outermost frame on the message handed to Layer C. Handlers are free to consume or ignore the transport promise; the trust ledger uses it.

### 7.2 The pCID registry is itself a commons

There is no global pCID registry in production PromiseGrid (because it has to last centuries with no central authority), and the harness must not have one either. Each agent maintains its own local pCID-knowledge graph: which CIDs it has seen, which protocols those CIDs identify, who vouched for each mapping. The simulator must demonstrate scenarios where:

- Two communities independently coin the same human-readable name for different pCIDs and learn to bridge them.
- A widely-used pCID is forked: half the community moves on to a v2 pCID, the other half stays on v1. The harness shows what bilingualism looks like at the wire level.
- An adversary publishes a pCID that *looks* like a popular one (typo squat) and the trust ledger discounts it appropriately.

### 7.3 Time is multi-clock

A long-lived system has no single clock. The harness gives each agent:

- A virtual wall clock (advanced by the scenario driver, used for timeouts, token TTLs).
- A logical clock (Lamport or vector, advanced by message events, used for ordering).
- A worldline hash chain (their own private append-only log, used for "did I really say this?").

Promises that depend on time **must specify which clock**. A `redeem-by` field on a token that uses wall clock will behave differently than one that uses logical clock when networks partition. The harness will run scenarios that exercise this.

### 7.4 Snapshot, fork, replay, diff

Beyond v1's transcript:

- **Snapshot**: at any virtual time, freeze every agent's state and every edge's queue. CBOR-encoded, content-addressed.
- **Fork**: from any snapshot, run two divergent scenarios. Useful for "what if Alice had defected here?" thought experiments.
- **Replay**: run the same scenario with a single variable changed (different scoring function, different ingress model, different pCID for the same logical message). Diff the resulting transcripts.
- **Counterfactual analysis**: given a transcript, ask the harness "would this break-witness have been believed if Alice's score had been 0.3 instead of 0.7?" without re-running the whole scenario.

### 7.5 Real-world cross-pollination

Realism upgrades:

- **Network impairment.** Real packet loss curves, real latency distributions from public CDN measurement datasets, real BLE link error rates. No flat 5% drop.
- **Real artifacts.** Some agents communicate through actual git pushes to a local repo; some through actual files on a tmpfs. Not simulated; the real syscalls.
- **Hardware-in-the-loop.** Optional: at least one agent in each long run runs on a real microcontroller (e.g. ESP32, RP2040) attached via USB-CDC, exposed as a `stdio` edge. Discovers what only-on-real-hardware bugs look like.
- **Browser-in-the-loop.** Headless Chromium running real JS agents under Playwright, with real `postMessage` between iframes. Exposes browser timing weirdness.
- **Real wall time long-running.** Some scenarios run at real wall-clock speed for hours or days, not just simulated time. Catches drift, GC pauses, certificate-expiry issues that fast-time runs miss.
- **Real human in the loop, occasionally.** A "human-confederate" mode where one slot in a stochastic-agent population is actually a human typing decisions. Used sparingly, in pre-release acceptance scenarios. Captures human-protocol-interaction issues.

### 7.6 Adversarial-by-default review

Every scenario has a **shadow adversary scenario**: same setup, but with one or more `hostile-prober` or `griefer` agents added. If the design survives the shadow, ship the design. If not, fix the design, then try the shadow again.

---

## 8. Thought Experiments

Each TE is a falsifiable experiment whose outcome teaches us something about the right design. The full chronological index lives in [`docs/thought-experiments/README.md`](docs/thought-experiments/README.md); the summaries below are the subset currently pulled into the canonical spec text, with later additions appended as they directly shape active decisions. Source: `DI-009-20260430-204108`.

- **[TE-1: Promise-stack ordering](docs/thought-experiments/TE-20260427-180000-promise-stack-ordering.md)** — Run two variants — signature outermost vs.
- **[TE-2: Trust-ledger merge after partition](docs/thought-experiments/TE-20260427-180100-trust-ledger-merge-after-partition.md)** — Two clusters trade independently for simulated days, then reconnect.
- **[TE-3: Currency exchange-rate equilibration](docs/thought-experiments/TE-20260427-180200-currency-exchange-rate-equilibration.md)** — Three agents issue tokens.
- **[TE-4: Sybil under double auction](docs/thought-experiments/TE-20260427-180300-sybil-under-double-auction.md)** — A single human spawns 50 identities to manipulate prices.
- **[TE-5: Kernel-as-handler vs. classical kernel](docs/thought-experiments/TE-20260427-180400-kernel-as-handler-vs-classical-kernel.md)** — Run scenario S-mixed (multiple resources on one node) under K1 and K3.
- **[TE-6: Capability-token revocation propagation](docs/thought-experiments/TE-20260427-180500-capability-token-revocation-propagation.md)** — Alice issues tokens that B, C, D delegate among themselves through three more hops.
- **[TE-7: Human-novice onboarding under K4](docs/thought-experiments/TE-20260427-180600-human-novice-onboarding-under-k4.md)** — Drop a `user-novice` stochastic agent into a K4 (no-kernel) community.
- **[TE-8: Generational handoff](docs/thought-experiments/TE-20260427-180700-generational-handoff.md)** — A `community-elder` retires after 6 simulated months.
- **[TE-9: Two communities, two pCIDs, same intent](docs/thought-experiments/TE-20260427-180800-two-communities-two-pcids-same-intent.md)** — Community A's spec for compute-call hashes to pCID `0xAAAA…`.
- **[TE-10: Slow-mover survival](docs/thought-experiments/TE-20260427-180900-slow-mover-survival.md)** — Ninety percent of the network has upgraded to `protocol-v2`; ten percent are stuck on `protocol-v1` (microcontrollers, sealed appliances, sentimentally maintained legacy).
- **[TE-11: Ostrom's principles audit](docs/thought-experiments/TE-20260427-181000-ostroms-principles-audit.md)** — Take a winning C-scenario (cooperators dominate at month 12) and a losing one.
- **[TE-12: Promise-stack as zero-knowledge envelope](docs/thought-experiments/TE-20260427-181100-promise-stack-as-zero-knowledge-envelope.md)** — Can a promise frame contain a ZK proof that "the inner payload satisfies predicate P" without revealing the payload?
- **[TE-13: Time-traveling break-witness](docs/thought-experiments/TE-20260427-181200-time-traveling-break-witness.md)** — Charlie discovers in year 5 that an event in year 2 was actually fraudulent.
- **[TE-14: A harness-spec change walks through the unified flow](docs/thought-experiments/TE-20260428-080000-harness-spec-change-walks-through-unified-flow.md)** — Stress-test of the single-flow design when the proposal targets the harness spec itself. Walks through end-to-end with a stochastic elder, two LLM analysts, an adversary, and Steve as deciding signer. Findings: the design holds; surfaces wobbles around proposal-collector matching, pointer storage, and settlement-window conventions.
- **[TE-15: Should this design become `promisegrid/promisegrid/README.md`?](docs/thought-experiments/TE-20260428-094500-should-this-design-become-promisegrid-readme.md)** — Compares the existing org-level README against the Wire Lab design, identifies stale and missing material, and recommends a four-phase migration: separate Wire Lab repo now, minor org-README touchup later, substantive rewrite once findings exist, full graduation eventually.
- **[TE-24: Group-transport envelope: `grid <pcid>` carrier, canonical bytes, and explicit promise body](docs/thought-experiments/TE-20260430-204108-group-transport-envelope.md)** — For the wire-lab's first transport-protocol (the group-transport-protocol class: small-finite-closed-group, N≥2, all-see-all, multi-writer DAG of messages), compares `grid <pcid>`, ordinary protocol headers, filename/path-selected protocols, and a pure structured object. Current finding: use `grid <pcid>` so the repo can test pCID-selected dispatch, message-CID-linked references, and explicit promise bodies without prematurely freezing the canonical PromiseGrid wire format. The substantive v0 contract — `Parents:` header for DAG links, body-level acknowledgement for receipts, no `Kind:` header, flat subdirectory layout — lives in `specs/group-transport-draft.md`. Source: `DI-009-20260430-204108`; open graduation question: `DR-009-20260430-204108`. Numbering note: this TE was originally drafted as TE-21 on `stevegt/channels-grid-pcid` and used "channel" vocabulary; per TE-25 (numbering reconciliation, drafting-time invariant) it is TE-24 on `ppx/main`, and per TE-27 + TODO 013 it has been renamed and rewritten in place under transport vocabulary.
- **[TE-25: TE-21 numbering collision and harness-spec path](docs/thought-experiments/TE-20260430-213447-te-numbering-collision-and-harness-spec-path.md)** — Reconciles a numbering collision when two branches independently drafted a "TE-21." Locks the drafting-time invariant: integers anchor on the first-drafted timestamp; later-drafted TEs renumber when branches collide. Same rule applies to TODO numbers. Also reaffirms `specs/harness-spec-draft.md` as the canonical path for the harness spec (post-TE-22 freeze convention).
- **[TE-26: Transport-protocol types, pCID-keyed transport paths, and DAG message graphs](docs/thought-experiments/TE-20260430-215624-channel-transport-types-and-threaded-replies.md)** — Surveys nine transport patterns (1:1, broadcast, subgroup multicast, pub/sub, anycast, request-reply, broadcast-with-receipts, single-writer log, general DAG) and concludes that transport type does not drive on-disk path layout. Locks four principles for the outer wire-lab transport-spec: a message does not declare its transport (no `Transport:` header); transport directories are keyed `transports/<pcid>--<slug>/` where the pCID is canonical protocol identity; each transport-protocol-pCID names a spec defining its directory's interior; the code that reads the directory structure is the handler for that pCID. The outer wire-lab transport-spec ships thin (`specs/transport-spec-draft.md`); TE-24's v0 contract is reframed as the group-transport-protocol's contract in a separate spec (`specs/group-transport-draft.md`). The conceptual shift toward DAG-shaped message graphs (zero-or-more parents per message) is surfaced; concrete header shapes are delegated to per-transport-protocol specs. (Originally drafted using "channel"; rewritten in place per TE-27. Filename retains "channel" from original drafting.)
- **[TE-27: Transports rename and axes of transport-protocol differentiation](docs/thought-experiments/TE-20260501-021921-transports-rename-and-axes-of-differentiation.md)** — Renames `channels/` to `transports/` throughout the repo (the directory is the wire lab's simulation surface for the network being studied; "channel" suggested an above-transport logical-addressing concept that wire-lab does not need yet). Establishes per-axis meta-rule for distinguishing transport-protocols: visibility (B) and routing topology (C) warrant distinct pCIDs; cardinality (A) is a parameter except at extremes (large-N, unbounded); membership (D), persistence (E), message-graph shape (F), direction (G), and receipts (H) are parameters within a single transport-protocol's spec. Locks the carve-out plan: one starter spec, `specs/group-transport-draft.md`, covering small-finite-closed-group with N≥2; the codex-perplexity case is the N=2 instance, not a separate spec. Anticipates future transport-protocol specs (ring, star, cluster-of-clusters, gossip) and a future TE on transport-protocol migration semantics.
- **[TE-28: The 100-year goal as a load-bearing design constraint](docs/thought-experiments/TE-20260501-202713-100-year-goal-as-design-constraint.md)** — Names the **100-year goal** that has been operating implicitly throughout the wire-lab's design and articulates its six load-bearing constraints (C-1 no central registry; C-2 multi-generational durability; C-3 adversarial-by-default; C-4 protocol forking is normal; C-5 trust accrues per-burden; C-6 signing key is the only structural lock). Pressure-tests every currently-locked artifact (harness-spec, transport-spec, group-transport-spec, TE-24/25/26/27, branching policy, TE/TODO numbering) against each constraint. Finds that no current lock fails outright, but several locks are silent on C-4 (protocol forking) and on cryptographic-signing migration discipline, surfacing five open questions OQ-100.1 through OQ-100.5 recorded in §12. Recommends the constraint codes C-1 through C-6 be used as shared shorthand in subsequent TEs so each TE does not re-derive the constraint set.

- **[TE-29: Protocols as simulated repos, and the L4-binding layer](docs/thought-experiments/TE-20260501-215027-protocols-as-simulated-repos-and-binding-layer.md)** — Locks the directory shape under which each protocol becomes its own simulated repo (`protocols/<slug>.d/` for live, `protocols/<slug>-<pcid>.{md,d/}` sibling pair for frozen), absorbs the former top-level `DR/`/`TODO/`/`DI/` directories as inline sections inside the relevant protocol's spec, and decomposes `transports/` into a five-level path `transports/<real-world-transport>/<L4-binding-pCID>/<session-pCID>/<message-pCID>/<message-id>.msg` aligned with the existing-Internet reality that PromiseGrid rides L4+ transports rather than re-specifying them. Defines the **L4-binding layer**: small specs (e.g. UDP-binding v0) that translate between opaque PromiseGrid session messages and the framing/sizing/addressing/connection conventions of one specific real-world transport. Renames `group-transport` to `group-session` (DFs T1-T6 unchanged; session-layer semantics, not L4 framing). Treats `proposals/` content as instances of a `ppx-dr` message protocol that should live under `transports/...`. Surfaces nine open questions OQ-29.1 through OQ-29.9 covering freeze ceremony, leaf-filename convention, framing layer, binding-layer signing, negotiation/handshake, real-world transport slug allocation, and external network-simulator integration (ns-3) for later iterations. Companion artifact: a one-page draft of UDP-binding v0 at `protocols/udp-binding.d/specs/udp-binding-draft.md`. Five follow-on TODOs (014-018) track the directory migration, DR/TODO/DI absorption, proposals-as-messages move, group-session rename, and UDP-binding v0 reference implementation.

- **[TE-30: TODO numbering and per-protocol TODO shape](docs/thought-experiments/TE-20260502-002548-todo-numbering-and-per-protocol-shape.md)** — Closes out the TODO-numbering question that TE-29 had named but not resolved. Locks Option D from the conversation of 2026-05-01: per-protocol `TODO/` subtrees inside each `protocols/<slug>.d/`, **no top-level `TODO/` directory**, timestamp-named TODO files matching the TE convention `TODO-YYYYMMDD-HHMMSS-<slug>.md`. The harness is itself a protocol (`protocols/wire-lab.d/`); harness TODOs live at `wire-lab.d/TODO/` like every other protocol's TODOs live in their own. Each protocol's local `TODO.md` lists only that protocol's TODOs; the master cross-listed `TODO.md` at `wire-lab.d/TODO/TODO.md` lists everything across the wire-lab. Boundary rule: a TODO is per-protocol if and only if its work touches files only inside one `protocols/<slug>.d/` tree; anything else is harness-level. **Renumbers all 19 existing TODOs**: 16 are harness-level and migrate to `wire-lab.d/TODO/`; 3 are per-protocol (012 -> group-session, 016 -> ppx-dr, 018 -> udp-binding). RETIRED (015) and FOLDED (017) status records migrate as ordinary timestamped files; no integer survivors. Canonical first-drafted timestamps are recovered from git commit history, never set to migration-time. Migration is delegated to TODO 014 as new subtasks 10 and 11 so it happens atomically with the rest of the protocols-migration. Surfaces five sub-OQs OQ-30.1 through OQ-30.5 covering boundary changes, sync discipline, DR/DI parallelism, harness-shape sub-categories, and `proposals/` cleanup. Partially closes OQ-100.4 (TODO half resolved by per-protocol timestamping; TE half still uses integer aliases per TE-25).

- **[TE-31: Spec-doc as upstream, simrepo as implementation: inverting the conformance reference](docs/thought-experiments/TE-20260502-004924-spec-doc-inversion-and-conformance-changelog.md)** — Closes OQ-29.1 (freeze ceremony) by inverting the doc/repo reference direction. Old mental model: spec doc embeds a tree-hash that names the simrepo. New mental model: each `protocols/<slug>.d/` simrepo carries a `CHANGELOG.md` whose entries each name a spec-doc CID and a conformance verb (`implements`, `partially-implements`, `extends`, `deprecates`); the spec doc itself stays upstream and oblivious, RFC-shaped. Locks **Alt-G**: the spec is frozen when its doc-CID is published and at least one simrepo CHANGELOG entry references it; no tree hash, no bundle. Supersedes Alt-A through Alt-F enumerated in TE-29. Justified by five arguments: (1) matches RFC reference direction; (2) one doc-CID supports many implementations without republishing; (3) implementations evolve their conformance claim cheaply; (4) companion files travel with the doc automatically because the doc cites them by pCID (per the user's 2.C answer that machine-checkable bits are hybrid inline + companion-pCID); (5) keeps the spec doc as a universally-quantified promise per TE-21. Sketches the CHANGELOG header-block format, harness conformance-check loop, breaking-change migration shape, and self-consistency with PromiseGrid's promise/acceptance primitive (the doc-CID is a promise; the CHANGELOG entry is an acceptance). Surfaces five follow-on questions OQ-31.1 through OQ-31.5 covering harness self-application, CHANGELOG location/format, breaking-change flagging, CHANGELOG-entry pCIDs, and an optional reverse-index in the spec-doc store. Implementation impact deferred to a forthcoming TODO 020 (CHANGELOG.md format and harness parser); TODO 014 unaffected.

- **[TE-32: Spec-side vs implementation-side split, and the `implementations/` top-level](docs/thought-experiments/TE-20260502-014525-spec-vs-implementation-split.md)** — Amends TE-31 (does not retract). TE-31 conflated two distinct categories of file: **A-side** (spec/design WIP — TEs, draft specs, design TODOs, manifests; what `protocols/<slug>.d/` was always for per TE-29) and **B-side** (conformance/implementation artifacts — code, test vectors, harness fixtures). The CHANGELOG language in TE-31 was always about the B-side. Locks the split: A-side stays under `protocols/<slug>.d/`; B-side moves to a new top-level `implementations/<impl-name>/` that does not exist yet but is created by TODO 014. Each side has its own `CHANGELOG.md` with different semantics: A-side records **freeze events** (`event: freeze` with a doc-CID) authored by spec maintainers; B-side records **conformance claims** (`claim: implements/partially-implements/extends/deprecates` with an upstream doc-CID) authored by implementers. Neither names the other directly; discovery is a query over known CHANGELOG sets in both directions. Reasons for top-level `implementations/` (not a subtree of `protocols/<slug>.d/`): (1) one impl may implement many protocols (a Go reference shipping UDP-binding+group-session+ppx-dr is one tree); (2) external implementations live in their own external repos with the same shape, location-independent; (3) spec and impl fork independently per C-4. Refines Alt-G slightly: a spec is frozen when its A-side CHANGELOG records a `freeze` event with a doc-CID (one-sided act, controlled by spec maintainers; does not require an implementation to exist, matching how RFC 768 was frozen on publication not on first BSD release). Resolves TE-31 OQ-31.3 (yes, explicit `breaking-change` boolean on B-side entries). Refines TE-31 OQ-31.2 (per-side CHANGELOG location specified; YAML-vs-fenced-vs-HTML-comment for the header block remains open). Refines TE-31 OQ-31.5 (bidirectional reverse indices; neither embedded in docs). Surfaces five sub-OQs OQ-32.1 through OQ-32.5 (external-implementation discovery; impl-of-multiple-protocols mapping; test-vector A-vs-B classification; harness-as-protocol vs harness-as-implementation per TE-30; frozen-sibling tree size). Adds steps 12-15 to TODO 014 to keep `protocols/<slug>.d/` design-only, create `implementations/` top-level with stub README, retarget TODO 018 and TODO 019 to live under `implementations/`, and seed empty CHANGELOG.md stubs in each migrated `protocols/<slug>.d/`.

---

## 9. Realism Suggestions

### Across the board

- **Run for simulated centuries on at least one design candidate.** Even at 1000× wall-clock acceleration, a centuries-long run is days of compute. The signal it produces about long-term protocol drift is irreplaceable. One overnight run per quarter is enough.
- **Recruit real second-implementer.** Have someone who is not Steve write a second `promstack` implementation in a different language from a separate reading of the spec. Their bugs are the spec's ambiguities. This is the single highest-leverage realism move available.
- **Adversarial bounty program inside the harness.** Pin a small budget; let outside contributors submit `griefer` agents that try to break running scenarios. Pay for new failure modes.
- **Embed a TLA+ specification of the trust-ledger merge rule and the promise-stack semantics.** Not the whole system; just those two. Gives us a formal object to point at when we argue.
- **Use real corporate / community networks for some scenarios.** Run the harness across two or three actual VPSes in different countries with real BGP paths; not just localhost. Catches MTU, NAT, RTT, timezone, and DST surprises.
- **Build the visualization first.** A 3D or graph-style real-time view of agents, edges, message flows, trust scores, and market prices. Bugs invisible in transcripts become obvious visually. This pays for itself within the first month.

### For long-livedness specifically

- **Crypto-agility from day one.** No scenario should hardcode an algorithm. Every signature promise carries an algorithm identifier; every scenario tries at least one rotation.
- **Documented deprecation path** for each protocol element. The harness can refuse to ship a pCID that doesn't declare what its successor will probably look like.
- **Backward incompatibility requires a multi-promise rollout.** Old peers must see *something* even when they can't fully process — at least an "I cannot help you, here is who can" reply. The harness will reject designs where old peers see opaque silence.
- **Death of agents is a first-class lifecycle event.** Not just rotation — actual cessation. The harness will cycle agents through death and successor-emergence in long runs. Discovers what happens to their open promises.
- **Simulate currency hyperinflation and deflation.** Issue rates set high or low; observe when the market collapses. Helps us pick currency-issuance policies that aren't catastrophic.

### For cross-runtime reality

- **Microcontroller runtime profile.** Define a "minimal endpoint" subset: which promise assertions it must understand vs. which it must reject gracefully. Real MCUs in real scenarios validate it.
- **Browser-tab profile.** Same exercise. The browser is a hostile environment (CORS, mixed content, tab-suspension, storage quotas). Embrace it.
- **JS-on-server profile.** Different from browser. Has its own wrinkles (event loop, fs access).
- **Bare-metal native profile.** No allocator? Single-threaded? Define what that means for our wire library.

### For governance reality

- **Have a community-governance layer in the simulation.** Not a "kernel" — a literal forum where stochastic `community-elder` agents argue and vote on proposed pCIDs. Their decisions feed back into which pCIDs the harness considers "blessed" for that community.
- **Multi-community scenarios.** Two or three communities with overlapping membership but separate elders. Models the real eventual structure of any open system.
- **Schism scenarios.** Communities split, each forks the protocol, both continue. Does the system handle the schism gracefully, or does someone lose their data?

---

## 10. Where This Lives in `grid-poc`

| New thing | Lives in |
|---|---|
| `promstack` Go library | `x/wire-v2/promstack/` |
| Trust ledger primitives | `x/wire-v2/trust/` |
| Stochastic agent SDK | `x/sims/simkit/stochastic/` |
| Currency / DA market handler | `x/sims/agents/market/` |
| Ingress models K1–K5 | `x/sims/simkit/ingress/` |
| Long-running scenarios | `x/sims/scenarios/long-run/` |
| ToC scenarios | `x/sims/scenarios/commons/` |
| Visualization frontend | `x/sims/viewer/` |
| Ostrom report generator | `x/sims/simkit/reports/` |
| Real-MCU and browser harnesses | `x/sims/simkit/realruntimes/` |

This sits alongside, not in front of, existing `x/wire`, `x/testbed`, `x/sim4`, `x/kernel1`. v2 is an additional experiment, not a rewrite. If a v2 design wins, the winning pieces graduate to the unprefixed locations.

---

## 10a. The Self-Improvement Loop: One Flow, Many Agents (Including You)

v2 produces enough logs to **diagnose** problems but not yet enough to let an agent — or the harness paired with an LLM, or you — actually **evolve a better design and propose it**. This section closes that gap.

The key simplification: there is **one** proposal flow, not two. Whether the target of a proposal is an in-sim protocol pCID or the harness spec itself, the same trust ledger and the same review queue apply. You are an agent in this system: a uniquely-trusted one whose signature on a merge promise re-issues the harness spec, but mechanically just another agent. An LLM helping you read proposals is also an agent — lower trust by default, but it can earn trust by making good predictions over time, exactly like any other peer.

This works because the harness spec is itself a spec at some pCID (call it `harness-spec-vN`), so a proposal to change the harness is just a proposal whose target is that pCID. No second proposal type is needed.

A secondary simplification, made in this revision: we are **not** locking the proposal envelope schema, an `invariants.cbor` file, a design-entropy formula, or a review-queue directory layout this early. Those are guesses about what the future system will need. The single structural commitment is described in §10a.8 — the harness-spec pointer follows your signing key. Everything else listed below is a useful starting convention, expected to evolve.

### 10a.1 Specs may be prose, structured, or both

A pCID names a spec. The spec is a content-addressed object — typically a CBOR-encoded bundle, but its *content* can be anything readable: prose markdown, IPLD schema, executable test scenarios, code, diagrams, or any combination. Three observations matter:

- **Prose is a legitimate spec form, especially in early drafts.** A new idea is usually clearer in English than in IPLD schema. Forcing schema-first stalls discussion. The first version of any proposal is allowed to be prose only, and the harness will route it the same way as a fully-structured one.
- **Structure is earned, not required.** As a spec matures and acquires conformance tests, an IPLD schema, an assertion vocabulary, and known-good handler implementations, it can be republished with those attached — a new pCID that supersedes the prose version. Proposals are encouraged to start prose, become structured, and converge.
- **Specs declare what they contain.** A spec object's outer wrapper says which sections are present (prose, schema, tests, examples, etc.) so a reader — human, LLM, or program — knows what it is dealing with. An LLM-readable spec without an IPLD schema is fine; a programmatic conformance-checker reading the same spec will simply skip the schema-check phase and report "schema not present."
- **Humans and LLMs are agents too.** A human author writing prose and signing it produces exactly the same artifact — a content-addressed spec at some pCID — as a code-generated one. The trust ledger evaluates the promiser, not the production method.

A handler library still names which (pCID, role) pairs it implements; if the spec has only prose, the handler library can declare "implements pCID P per author's reading of the prose at <date>." Disagreements about the prose's interpretation surface as `contest` messages on the wire — useful diagnostic signal.

### 10a.2 Proposals are messages, not a locked envelope

A proposal is a promise stack like any other message. We do **not** define a frozen `Proposal` schema this early; we don't yet know enough about what proposals will need to carry, and locking a parser-level shape now would force every later refinement to amend the schema or work around it.

What we do instead:

- A proposal is a message whose outermost assertion is something like "this is a proposal to change pCID X / knob Y / the harness spec itself." The target is named, the proposer signs it. Everything else is body.
- A *suggested* proposal shape — call it `proposal-checklist-v0`, published as a prose spec at its own pCID — describes the kinds of content reviewers tend to find useful: a diagnosis ("what's wrong, what I expected"), the proposed change, a predicted effect (how the proposer commits to being wrong), counter-evidence the proposer searched for, and a rollback sketch. This is a community convention, not a parser rule. Reviewers are free to ignore proposals that lack these. Better shapes will emerge by being reused, exactly the way pCIDs themselves accrue trust.
- A proposer's prediction-quality history accrues to their trust ledger under a `design-judgment` assertion type. Reviewers see it. An LLM analyst with a year of bad predictions becomes appropriately easy to ignore — without any schema enforcing it.
- Proposals can be entirely prose, entirely structured, or any mix. A prose proposal signed by an agent (human, LLM, or program), targeting the harness spec, is a first-class citizen of this system.

What we explicitly defer: the exact field names, whether `predicted_effect` is mandatory at the parser level, what counts as adequate counter-evidence, and whether prose-only proposals get a quality discount. All of these are downstream decisions to be made once we have actual proposals to look at.

### 10a.3 Endorsement, contest, counter-proposal — same vocabulary for everything

Three more pCIDs round out the discourse layer:

- `endorse-v1`: "I, peer Q, endorse proposal P; here is my evidence (prose or run-IDs)." Endorsements are promises; their weight depends on Q's trust score for design-judgment.
- `contest-v1`: "I, peer Q, contest proposal P; here is the conflicting observation (prose or run-IDs)." Contests must point at *something* — a specific prediction the proposal made, a specific run the proposer didn't cite — not bare disagreement.
- `counter-propose-v1`: a proposal that explicitly references another and supersedes it. Forms a DAG of design-evolution attempts.

All three accept prose or structured content. All three apply to harness-spec proposals exactly as they apply to in-sim protocol proposals. The harness has no built-in voting mechanism; communities (and you, for harness-spec changes) define their own decision rules, themselves encoded as `governance-v1` specs at whatever pCID the community chose.

For harness-spec proposals specifically, the decision rule is simple by default: **your signature is what merges the change**. An LLM analyst can read, summarize, contest, and endorse — its endorsement is just another promise weighed by its (low) trust score — but only your signature on a "merge" promise re-issues the harness spec at a new pCID. This isn't a configured rule; it's a property of which signing key controls the canonical harness-spec pointer. It is the one and only structural lock the design relies on at this stage.

### 10a.4 Cross-run aggregation and the design-knob registry

For consumer (b) — the harness + LLM — we need data structured for cross-run comparison.

A top-level `runs/index.parquet` maintained by the harness, columns:

```
run_id, scenario_id, scenario_hash, seed,
knob_set_hash,                   ; canonical hash of all knob values
knob.<name> ...,                 ; one column per knob; e.g. knob.ingress_model
outcome_summary,                 ; pass/fail per assertion
metric.<name> ...,               ; cooperator_txn_rate_m12, sybil_cost, etc.
ostrom_score.<principle> ...,    ; one column per principle
start_time, duration_s, harness_version, sim_kit_version
```

Every scenario YAML is required to declare its knobs in a top-level `knobs:` block, with name, legal range, default, and a one-line description. The harness validates this at scenario load time. From those declarations the harness auto-generates a **design-knob catalog** (`docs/knobs.md` and `knobs.cbor`) that is itself content-addressed and updated per-run.

The LLM analyst, when invoked, reads `runs/index.parquet` + `knobs.cbor` and can ask questions like "which knob settings most strongly correlate with cooperator-survival failure?" using ordinary SQL/dataframe operations. It is no longer guessing from prose.

### 10a.5 Outcome attribution: which knob caused this?

Correlation alone produces noise. To attribute outcomes to knobs, the harness will likely use a couple of standard tools, neither novel:

1. **Sampling sweeps.** Latin-hypercube and factorial sampling that vary one knob at a time and several at once. Results saved per-study.
2. **Counterfactual replay.** Re-execute a run with one knob changed and emit a delta report. The natural unit of attribution is the delta between two replayable runs that differ in exactly one knob value.

A strong norm — not yet a parser rule — is that a causal assertion ("knob X caused outcome Y") should cite a sweep or delta report. Reviewers are free to demand this and to discount proposals that don't supply it. Whether the harness eventually enforces this at envelope-parse time is a question for after we've watched a few iterations and seen where the norm actually fails.

### 10a.6 Counterfactual evidence for in-sim agents

Agents inside a run cannot fork the world. But they can do something cheaper: **emit hypotheticals**.

A new assertion type, `"hypothesis: under spec P', message M would have been accepted"`, lets an agent record its private theories on the wire as ordinary promises. The agent commits to the hypothesis (signs it), and the harness's *adversarial-replay* harness then re-runs the scenario with spec P' substituted and emits a verdict: hypothesis confirmed, refuted, or unrelated. The verdict is filed back as a `hypothesis-result-v1` envelope.

This lets agents propose changes from inside the system *with* evidence — their hypothesis got tested, and the test result is part of their reputation now. Speculation without verification costs trust; speculation that pans out gains it.

### 10a.7 Durable cross-run agent memory

For an agent to learn across runs, it must persist state across scenario boundaries.

The harness provides each agent identity (not each agent instance) with a `personal-archive/<agent-id>/` directory that survives runs. Default contents:

- The agent's trust ledger (extended across runs, with run-id provenance per entry).
- The agent's pending and resolved hypotheses.
- The agent's outgoing proposals and their fate.
- The agent's predicted-vs-actual scoreboard (their prediction-quality score over time).

Agents that elect not to use the archive run as before. Agents that do can begin to behave like long-lived community members across many simulated years, and across many distinct harness invocations on the wall clock.

### 10a.8 What's actually load-bearing: your signature

An earlier draft of this section proposed a global `invariants.cbor` file enumerating properties "that must hold under every change." On reflection, the examples I wrote turned out to be a mix of (a) target metrics whose right values are exactly what the harness is supposed to discover, (b) consequences of other definitions (the trust ledger, content addressing) that don't need restating, and (c) parser-level rules that disappear once we drop the proposal-envelope schema. None of them were genuinely constitutional; they were guesses about what the future system would need to defend.

So there is no `invariants.cbor`. The one mechanism that *is* load-bearing — the only thing that prevents a runaway from agents or LLMs editing the system into incoherence — is much simpler:

> **The canonical harness-spec pointer is whatever pCID you have most recently signed a `merge-harness-spec` promise for.** Any agent can propose, any agent can endorse, any agent can publish an alternative harness spec at its own pCID. But "the" harness spec — the one the next run starts from — is identified by your signing key.

That's it. No CBOR config file, no constitutional list, no invariant-change envelope. The lock is the key, not a document.

Communities (and you, in your role as harness-design agent) will accumulate norms over time — "don't merge a harness change without a counterfactual run," "don't approve a proposal whose predicted effect is incoherent," and so on. Those norms can become explicit specs at their own pCIDs once they've been used enough that their shape is obvious. They are not pre-installed.

### 10a.9 "Is the design converging?" — track something, but don't lock the formula

We'd like a rough scalar that answers "is the design converging or fragmenting?" Candidate inputs (not a fixed formula):

- Number of unresolved open questions.
- Number of provisional pCIDs awaiting graduation or retirement.
- Number of knobs whose chosen value is still wide-open vs. narrowed.
- Rate of new pCIDs minted vs. rate of old ones retired.

A reasonable starting move is to log all of these per-iteration in whatever shape is easiest, look at the time series after a few iterations, and *then* decide which combination (if any) is worth elevating to a single "design entropy" number. Picking a closed-form formula now ranks possibilities the harness has never seen against each other; that's how you end up optimizing the metric instead of the design.

### 10a.10 The review queue and your role as an agent

Proposals — whether targeting in-sim pCIDs or the harness spec itself — land in a single review queue. The on-disk layout is whatever's convenient at the time; we are not pre-committing to a directory schema. A reasonable *starting* shape might be `proposals/pending/<id>/`, `approved/<id>/`, `rejected/<id>/`, with each proposal directory containing the original message, any attached evidence, and an append-only file of endorsements and contests received. We'll change it the first time it gets in the way.

The shape that matters is the workflow, not the filesystem. **Your job is to act as an agent reading this queue**, with these supports:

- **LLM as reading aid.** An LLM analyst can be invoked at any time to summarize a proposal, surface its prediction-quality history, list its endorsers and contesters by trust score, run cross-run queries against the harness's accumulated logs, and (for harness-spec proposals) explain in plain language what the change would do. The LLM produces prose; you read it and decide. The LLM does not move proposals between states on its own.
- **Multi-LLM optional.** You can ask several LLMs to read the same proposal independently and surface their disagreements. This is just multiple agents endorsing or contesting; the trust ledger handles weighting.
- **Your decision is itself a wire message.** When you decide a proposal, you sign a promise — "I, Steve, approve/reject proposal P, my reasoning is …" — that becomes part of the public record. Future LLM analysts read this record and learn what kinds of proposals you find convincing or unconvincing.
- **Approved harness-spec proposals re-issue the harness-spec pointer.** A new harness spec at a new pCID becomes the one your signing key now points at; the next run starts on the new baseline. Every prior run remains associated with the harness-spec version that produced it, by content hash, for replay purposes.
- **Approved in-sim proposals re-issue the targeted spec or update the named knob.** Same mechanism, different target.
- **Rejected proposals are durable.** They remain visible with your reasoning attached, both as input to future LLM analysis and as evidence in trust-ledger updates for the proposer.

This keeps you in the loop without putting you in the hot path. LLMs read voluminously, run queries, and write summaries; you spend your attention on the small set of well-prepared proposals that survive prior filtering.

### 10a.11 Closing the loop: the canonical iteration

Once all of the above is in place, one *iteration* of the self-improvement loop looks roughly like:

1. The harness runs the current scenario suite plus any open studies.
2. Aggregator updates the cross-run index and whatever convergence indicators we're tracking at the moment.
3. In-sim agents (consumer a) emit their own proposals, endorsements, and contests; these are persisted.
4. The LLM analyst (consumer b) reads the cross-run data, reads in-sim proposals, reads the previous round's rejected proposals, and produces zero or more new proposals.
5. You triage the queue, signing approve or reject decisions with your reasoning.
6. Approved changes re-issue the targeted spec, knob, or harness pointer at a new pCID.
7. Next iteration starts from the new baseline.

Cadence is an open question — see §12. Continuous in-sim activity is fine; the human review interval is a knob we'll tune from experience.

### 10a.12 What this enables and what it does not

**Enables:**

- An in-sim agent can advocate a design change and have it *tested* before peers decide.
- The harness, with an LLM, can recommend changes grounded in cross-run evidence rather than vibes.
- You read proposals with whatever evidence the proposer assembled, plus an LLM's reading-aid summary, and decide.
- Rejected proposals are durable; the LLM (and any future analyst) gets better at predicting your taste by reading the trail of past decisions.

**Does not enable, by the only deliberate constraint we are confident in this early:**

- Autonomous shipping. The harness-spec pointer moves only when you sign a merge promise. No script, no LLM, no in-sim agent can do that on your behalf. Everything else — what proposals look like, what counts as adequate evidence, what norms govern the LLM's behavior — is left to evolve by the same proposal/endorsement mechanism it governs.

---

## 11. The Decision the Harness Is Trying to Make

Restated as a single question:

> **What stack of promises, traveling on what mix of transports, evaluated by what trust ledger, gated at what kind of ingress, settled in what kind of market, can let an open community of free agents and humans continue to find each other and trade value across centuries — without succumbing to tragedy of the commons, governance capture, or the thousand other ways open systems quietly die?**

The simulator does not answer this question. It produces traces, prices, scores, and Ostrom reports. We — the humans designing PromiseGrid — answer it by reading those artifacts and arguing about them. The simulator's job is to keep us honest.

---

## 12. Open Questions for You

1. **Promise frame canonical encoding.** CBOR array of fixed-shape maps, CBOR Sequences, or something IPLD-schema-driven? My instinct: IPLD schema with CBOR codec, so we get cross-language parsers free. Confirm?
2. **Currency vs. mutual credit as the primary economic test.** I've designed for currency-with-DA-market because it gives the cleanest exchange-rate-as-trust-signal experiment. Should we also stand up a mutual-credit scenario as a co-equal test, or treat it as TE-only? My instinct: co-equal, because real PromiseGrid communities will likely use both.
3. **Stochastic agent backing.** LLM-driven stochastic agents (with recorded transcripts for replay) or scripted-policy stochastic agents (cheaper, less realistic)? My instinct: scripted by default for CI; LLM for occasional deep runs, with the LLM transcript stored as part of the run artifact.
4. **Hardware-in-the-loop scope.** Required for every release, or quarterly? My instinct: quarterly, with one MCU and one browser tab, but the harness should support it on demand from day one.
5. **Centuries-long simulated runs.** Is one overnight run per quarter enough, or do we need a continuous, perpetually-running "PromiseGrid in a tank" instance that we keep alive across releases as a long-term experiment? I think the latter would be enormously valuable — a simulated village that we never reset — but it's a real ops commitment.
6. **Pcid for break-witnesses.** Do we want one canonical break-witness pCID (so any agent can publish a witness in a universally-understood format) or do we let each protocol define its own (so the witness carries protocol-specific evidence)? My instinct: one canonical envelope, with protocol-specific evidence inside.
7. **Single trust score vs. per-assertion-type vector everywhere.** Section 2 argues for the vector. Vectors are harder to compare across agents. Do we accept that as a feature (no global trust score is possible, only local per-assertion ones) or build a default projection to a scalar for human consumption?
8. **Self-improvement cadence.** §10a proposes daily LLM passes and weekly human review of `proposals/pending/`. Too fast or too slow? My instinct: start weekly for both, tighten only after we've seen the LLM's proposal quality stabilize.
9. **Is "your signature is the only lock" really enough?** §10a.8 reduces the constitutional set to one rule: the harness-spec pointer follows your signing key. No `invariants.cbor`, no envelope schema, no LLM-can't-touch-this list. Is there a failure mode you can picture where that's insufficient — e.g. an LLM that floods the queue with plausible-looking proposals, exhausting your attention until something bad slips through? If so, the right defense is probably a rate-limit on submissions per agent (a community norm, not a parser rule), not a return to constitutional documents.
10. **Default trust for LLM analysts.** An LLM endorsing a proposal carries some weight by default — but how much? My instinct: very low at first (e.g. 0.05 of a human-elder endorsement), rising only as its prediction-quality history accumulates good calls. This means a fresh LLM cannot move proposals on its own; it can only summarize and recommend until it earns standing. Confirm?
11. **What else are we still over-locking?** This pass softened the proposal envelope, the invariants file, the design-entropy formula, the review-queue layout, and a number of supporting structures. The `Promise` shape in §1.1, the `TrustLedger` shape in §2.1, the agent binary in §5.1, the K1–K5 ingress taxonomy, the seven C-scenarios, the 30%-at-month-12 pass line, and the grid-poc directory table in §10 are all still presented as more settled than they really are. Which of these do you want to revisit next, and which are stable enough that committing to them is acceptable?
12. **Proposal-collector matching (from TE-14).** With no envelope schema, what does the handler that files messages into the review queue actually pattern-match on? Likely an LLM-as-router in the early phase. Acceptable, or do we want a minimal structural marker (one CBOR tag, say) on proposal-shaped messages just to make routing deterministic? My instinct: accept the LLM-router for now; revisit when the LLM mis-classifies enough to be annoying.
13. **Pointer storage semantics (from TE-14).** "The harness-spec pointer follows your signing key" is true at the crypto level but the pointer has to physically live somewhere. Options: (a) a workspace file you `git push`, (b) a CRDT replicated to all participants, (c) a DNS-like lookup against your pubkey. Which? My instinct: (a) for now (matches the multi-file workspace structure), (b) eventually for the production system.
14. **Settlement-window convention (from TE-14).** Proposers should declare a predicted-falsification window ("check my prediction in 7 days / 3 months / one full simulated year"). Should this be part of the proposal-checklist convention from §10a.2, or left entirely to reviewer discretion? My instinct: in the checklist; reviewers strongly prefer proposals that specify when they'd know they were wrong.
15. **Doc structure: monolith or split?** This document has grown to ~60KB. A two-layer structure — a tight `harness-spec.md` plus one file per thought experiment in `docs/thought-experiments/TE-YYYYMMDD-HHMMSS-*.md`, each content-addressable — would (a) keep the canonical spec readable, (b) let LLM coding agents drill into only the TEs they need, (c) make pCID churn manageable (small spec changes don't require re-hashing the whole TE corpus), and (d) embody the design's own "specs may be prose; structure is earned; everything is content-addressed" principle. Should we split now? My instinct: yes, before the next round of TE work makes the monolith painful.
16. **Spec doc as both simulation-spec AND production-dev-guide.** Plausibly the same document drives the harness *and* serves as the architectural guide that LLM coding tools use to write production code. Pro: single source of truth, prose-rich, embodies the methodology. Con: subtle pull toward writing the production guide as if discovery is done. Mitigation: tag every section as `[Harness]` (provisional), `[Production-bound]` (stable enough to ship), or `[Methodology]` (how we work). Today only §10a.8 is `[Production-bound]`. Adopt this tagging?
17. **When does Phase 2 of the README migration trigger (from TE-15)?** Phase 2 is the additive "two paragraphs pointing at the Wire Lab" update to the canonical PromiseGrid README. What's the right trigger? My instinct: the first thought experiment whose findings produce a concrete recommendation worth quoting publicly — not earlier (premature) and not after several findings (the README falls behind). Probably TE-2 or TE-3 will be the trigger.
18. **Does the canonical PromiseGrid README get its own pCID, signed by Steve, distinct from the Wire Lab harness-spec pCID (from TE-15)?** Probably yes — they have different scopes and different cadences. The README is the org-level identity statement; the harness-spec is one experiment within. Two pCIDs, two signatures, both following Steve's key. Confirm?
19. **Should the group-transport envelope graduate into the canonical PromiseGrid wire format (from TE-24 / DR-009)?** *Scoped to the wire-lab's first transport-protocol (the group-transport-protocol class).* The group-transport-protocol uses a textual `grid <pcid>` carrier so humans and LLMs can exercise protocol selection, message-CID-linked DAG references, and explicit promise bodies in real coordination traffic. That envelope choice is locked for the protocol class in `DI-009-20260430-204108` and the substantive contract is in `specs/group-transport-draft.md`, but whether the eventual canonical wire format should keep the same visible first-line carrier or collapse into a more structured promise-stack object remains open. See `DR-009-20260430-204108`.
20. **OQ-100.1 — Protocol forking representation (from TE-28 / C-4).** When an existing transport-protocol is forked (incompatible v2 of group-transport, or a community-specific variant of any protocol), how is the new pCID expressed in `protocols/` and in `transports/`? The per-axis meta-rule in TE-27 covers when *new* protocol classes warrant distinct pCIDs, but is silent on intra-class forks. Candidate follow-on TE.
21. **OQ-100.2 — Cryptographic signing migration (from TE-28 / C-2 + C-4).** When v0 group-transport gains cryptographic signing in v1, how does the migration honor C-2 (existing v0 transports survive untouched) and C-4 (v0 and v1 are siblings, not parent-and-child)? The transition affects `From:` semantics, message-CID computation, ack semantics, and the trust ledger's interpretation of pre-signing messages.
22. **OQ-100.3 — Cumulative-prefix ack and trust (from TE-28 / C-5; deferred Q2 from TODO 013).** The deferred cumulative-prefix or frontier ack semantics for group-transport must be designed under C-5 — trust-vector-aware, not global. The shape of "I have observed everything up to frontier F" is constrained by the requirement that *each receiver* maintains its own per-burden trust scoring of *each sender's* ack claims.
23. **OQ-100.4 — Numbering wrap (from TE-28 / C-2).** Is the integer sequence for TEs and TODOs stable across centuries, or does it eventually need to be supplanted by purely timestamp-based or pCID-based identifiers? The drafting-time invariant from TE-25 keeps integers usable for now, but a contributor in 2096 with TE-9847 may find the integers more noise than signal.
24. **OQ-100.5 — Slug drift (from TE-28 / C-2).** Human-readable slugs in `transports/<pcid>--<slug>/` and `protocols/<slug>.d/` (anticipated) are presentational, but a 30-year-old slug whose meaning has drifted may mislead a future contributor more than it helps. Is there a discipline for retiring or redirecting drifted slugs without changing pCIDs, or is the right answer to lean harder on pCIDs and ignore slug drift?
25. **OQ-29.1 — Freeze ceremony (from TE-29). CLOSED by TE-31 with Alt-G; refined by TE-32.** The spec is frozen when its A-side `protocols/<slug>.d/CHANGELOG.md` records a `freeze` event with a doc-CID. No tree hash, no bundle. Supersedes Alt-A through Alt-F. The freeze is a one-sided act controlled by spec maintainers (does not require any implementation to exist). B-side `implementations/<impl-name>/CHANGELOG.md` entries reference the doc-CID upstream when claiming conformance (RFC-shape inversion: implementation-to-spec, never spec-to-implementation). See TE-31 for the inversion thesis and TE-32 for the A/B split that disambiguates which CHANGELOG carries which semantics.
26. **OQ-29.2 — Leaf filename in `transports/...` (from TE-29).** Should `<message-id>.msg` use a content hash (free cross-transport deduplication, hash-as-identity), a ULID (trivially sortable, ordering preserved by filename), or both (`<ulid>-<hash>.msg`)? Lean: hash; ordering recovered from session-layer Parents header.
27. **OQ-29.3 — Empty-leaf rule (from TE-29).** Is `transports/<wire>/<binding-pCID>/` allowed to exist with no session-protocol subdirectory beneath it (meaning the transport+binding is provisioned but no session has spoken yet)? Lean: yes — honest about provisioning state.
28. **OQ-29.4 — Framing intermediate layer (from TE-29).** Do we need a separate framing layer between the L4-binding and the session protocol, or does each binding own its framing? Lean: each binding owns its framing (no reuse to be had across UDP datagrams vs. TCP byte-stream vs. SMTP body vs. MQTT publishes).
29. **OQ-29.5 — Binding-layer signing (from TE-29 / C-6).** May a binding sign per-frame for spam resistance independently of session and message-layer signing? Lean: yes-permitted but out of scope for v0 specs.
30. **OQ-29.6 — Negotiation/handshake protocol (from TE-29 / C-1).** Where does "agree on the (binding, session, message) tuple" live? Probably its own protocol on a bootstrap transport. Future work; cannot be a singleton per C-1.
31. **OQ-29.7 — Real-world transport slug list (from TE-29 / C-1).** Who decides which slugs appear at level 1 of `transports/<slug>/`? Lean: any author may introduce a new directory; collisions resolve by C-1 (no central registry); practical norm is to use IANA names where they exist.
32. **OQ-29.8 — Adversarial test-harness placement (from TE-29 / C-3).** Does the wire-lab simulator emulate loss, reordering, latency, partition? If yes, where does the emulation live in the directory shape? Lean: adversarial harness lives in `tools/` and produces test vectors that violate binding promises in controlled ways.
33. **OQ-29.9 — External network-simulator integration, ns-3 or similar (from TE-29; revised 2026-05-01 with empirical sandbox data).** Once at least two binding specs and one session protocol have v0 implementations in `tools/`, can the wire-lab harness drive a packet-level simulator to produce realistic loss/latency adversarial scenarios? Empirical results from the Perplexity Computer sandbox (Debian trixie, kernel 6.1.158): ns-3 3.43 works (Debian package, CMake+Ninja+pkg-config, UDP-echo+PCAP smoke-tested); Mininet 2.3.0 works with LinuxBridge backend (requires `PYTHONPATH=/usr/lib/python3/dist-packages`); raw netns+veth+bridge works under sudo; **`tc netem` is unavailable** (kernel lacks `sch_netem`, `tbf`, `fq_codel`); OMNeT++ and Shadow are not packaged. Revised lean: **adopt ns-3 as the primary network-realism fixture** (it is in fact lower-friction in the sandbox than `tc netem`, which the original lean had recommended first). Mininet remains useful for topology-only experiments. Steve's local dev box has `tc netem` and may use it; the sandbox does not. Integration shape unchanged: ns-3 emulates the wire, Go reference implementations run as applications above sockets via tap-bridge or DCE, PCAP output post-processed into `transports/<wire>/<binding-pCID>/.../<msg-id>.msg` artifacts. Future home: `tools/ns3-harness/`. Warrants a dedicated TE once the first Go binding+session implementations exist.
