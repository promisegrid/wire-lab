# TE-28: The 100-year goal as a load-bearing design constraint

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-20260501-202713

## Status

recommended for immediate adoption

## Decision under test

The wire-lab and PromiseGrid as a whole are designed against what Steve calls **the 100-year goal**: PromiseGrid must survive and remain useful as an open, decentralized community of free agents and humans across roughly a century of evolution, without depending on any central authority that might disappear inside that horizon. Every design choice the wire-lab makes — pCID conventions, transport-spec, group-transport-spec, harness-spec, freeze conventions, branching policy, TE numbering, and everything that comes later — must survive pressure-testing against the 100-year frame. The harness-spec gestures at this in places ("multiple human generations", "centuries-long", "multi-decade scenarios") but does not name the constraint set explicitly, and so does not give either future-Steve, future-LLM-readers, or future-contributors a single citable place where the constraints live.

The decision under test is: **what constraints does the 100-year goal name, and does the current locked design honor each of them?** If any current lock fails a 100-year constraint, the failure must be surfaced as an open question or follow-on TE so a short-horizon assumption does not get smuggled in unnoticed.

## Why this TE now

Two pressures converged during the session that produced this TE:

1. The session was working through how `protocols/` should be shaped and how releases should be expressed (anticipated future TE on simulated-repo releases). Several of the candidate alternatives looked attractive at a one-month horizon and unsafe at a 100-year horizon. The agent caught itself recommending shapes that depended on an assumption (e.g., "the wire-lab repo will exist forever in this form") that the 100-year frame forbids.
2. Steve asked whether the 100-year goal had made it into any documents. It hadn't, beyond softer phrasings. The constraint set has been operating implicitly; making it explicit is overdue.

This TE is the explicit record. Subsequent design TEs (release machinery, cross-protocol citation, pCID registry semantics, freeze ceremony) cite this TE when they pressure-test their alternatives.

## Assumptions

- PromiseGrid is intended to outlive every centralized institution it currently depends on (GitHub, the wire-lab repo, the agent-running-this-session, the cryptographic primitives currently considered safe, every individual contributor, and Steve himself).
- The wire-lab is a research-grade simulation of PromiseGrid; its findings are intended to inform the canonical PromiseGrid spec, not just the wire-lab itself. A lock that fails the 100-year frame inside the wire-lab will fail it outside the wire-lab too.
- Constraints are not goals in themselves. They are filters. A design that satisfies all constraints may still be wrong; a design that violates any constraint is wrong.
- "100 years" is shorthand for "longer than any individual contributor's career or any individual institution's reliable presence." The exact number does not matter; the qualitative regime does.

## The constraint set

The 100-year goal names six load-bearing constraints. Each is stated as a thing the design must do or must not assume. Each was articulated by Steve in earlier sessions and is recorded here as the canonical statement.

### C-1: No central registry

The design must not assume any global registry exists for any naming or coordination function over the 100-year horizon. Any centralized registry that exists today (DNS, ICANN, IANA, GitHub, npm, the Ethereum mainnet, the PromiseGrid project itself, this wire-lab repo) is presumed gone or untrustworthy at some point inside the horizon. Naming, discovery, and trust must work without one.

The pCID design is shaped directly by this constraint: pCIDs are content-addressable and mintable by anyone, so no registry is required for protocol identification. But the constraint extends further. It also forbids assuming a single canonical "list of protocols," a single "pCID-to-prose-name lookup," or a single "deciding signer for the whole network." Each agent maintains its own local pCID-knowledge graph; everything global is convention or emergent property, not infrastructure.

### C-2: Multi-generational durability

The design must support handoff across multiple human generations. Signing keys rotate, retire, get inherited, get lost, get compromised. Reputations accrue, decay, get willed to successors, get repudiated. Communities of interest fragment, merge, fork, fall silent for decades, return. The wire-lab's existing "break-witness" and "trust ledger" machinery is designed for this. So is the "kernel is just a role" framing — no role can outlive the agent currently playing it.

The constraint extends to the wire-lab's own artifacts: TE/DR/DI/TODO files, frozen specs, branching conventions, the harness itself. Any artifact whose meaning depends on Steve being available to interpret it fails this constraint. Artifacts must be self-explanatory to a contributor who arrives 30 years later with no living mentor.

### C-3: Adversarial-by-default

The design must assume adversaries are present at every layer at all times. Slander attacks, Sybil-by-witnessing, typo-squatting popular pCIDs, fork-and-rename-as-original, governance capture, "tragedy of the commons" failure modes (bandwidth flood, storage parasite, compute leech, reputation harvester, naming squatter, governance capture, liability dodge — all seven from the harness-spec). No mode where "we trust each other so this part is fine" is acceptable, even temporarily.

This constraint is the reason every layer must be auditable in isolation, every promise must carry its own signature chain, every receipt must be verifiable against locally-held trust ledgers, and every spec must define its failure modes as carefully as its success modes.

### C-4: Protocol forking is normal

Over a century, protocols will be forked into bilingual variants, contested, declared deprecated by sub-communities, rediscovered, adapted to new substrates. The design must treat forking as a normal life-cycle event, not an exceptional one. Forks must be representable, distinguishable, and individually subject to trust scoring. A pCID identifies a *spec*, not the *use* of a spec; multiple competing pCIDs for analogous protocols is the steady state, not the failure case.

The wire-lab's existing per-axis meta-rule from TE-27 (visibility and routing topology warrant distinct pCIDs; cardinality is a parameter except at extremes; etc.) is partially shaped by this constraint: it's the discipline that prevents one protocol from quietly becoming several without a fork-event being visible. The release machinery TE that follows this one will need to address how forks get expressed in `protocols/`.

### C-5: Trust accrues per-burden

Trust is not a global scalar attached to an identity. It is a per-promise-type, per-relationship, per-context vector that each agent maintains in its own local trust ledger. Alice may trust Bob's promises about routing latency and not his promises about identity verification. Communities may converge on rough trust patterns through gossip and witnessing, but no global trust score exists, and any design that produces or relies on one fails this constraint.

This is why the harness-spec emphasizes "trust ledger" rather than "reputation system," why receipts are scoped to the promise they acknowledge rather than to "good behavior" generically, and why the "kernel is just a role" framing matters: a kernel-role agent's trust must be earned per-burden, not by virtue of its role.

### C-6: Signing key is the only structural lock

The only durable anchor in the system is the signing key of an individual agent. Specs follow signing keys. Trust ledgers follow signing keys. Histories follow signing keys. Every other "fixed point" — repository URLs, registry names, organizational identities, project websites — is convention, and conventions evolve, decay, or get attacked over the horizon. The signing key is the one thing the design grants permanent structural status to.

This constraint is what makes pCIDs the right shape: they are content-addressable, derived from bytes, and the bytes can be re-signed by any current key-holder without changing the pCID. The signature chain is what attests current authority; the pCID is what attests document identity; neither requires a registry.

## Pressure-test of the current locked design

This section walks each currently-locked artifact in the wire-lab through the six constraints. Where the artifact honors the constraint, that is recorded. Where it fails or is silent on the constraint, that is recorded as an open question or as a candidate for follow-on TE work.

### `protocols/wire-lab.d/specs/harness-spec-draft.md` and the frozen `harness-spec-bafkrei...md`

- **C-1:** Honored. The harness-spec explicitly says "there is no global pCID registry in production PromiseGrid (because it has to last centuries with no central authority), and the harness must not have one either." The harness-spec is consistent with C-1.
- **C-2:** Honored in spirit (the "multi-generational handoff" subsection in §6, the "multi-generational durability" property under break-witnesses). Could be strengthened by explicitly naming the 100-year goal at the top of §1 or §2 so contributors see it before reading any prose. **Action: add a short subsection.**
- **C-3:** Honored. The seven tragedy-of-the-commons scenarios, the Ostrom-report scoring, and the adversarial scenarios in §6 enumerate the constraint's implications.
- **C-4:** Partially honored. The harness-spec acknowledges that competing pCIDs for analogous protocols may coexist, but does not lock a discipline for how forks should be represented or how trust gets distributed across them. **Open question for follow-on TE: how does protocol-fork representation interact with `protocols/<slug>.d/` shape?**
- **C-5:** Honored. Per-burden trust ledgers are first-class throughout the harness-spec.
- **C-6:** Honored. The harness-spec is explicit that "the only durable anchor in the system is the signing key" and that "specs follow signing keys."

### `protocols/wire-lab.d/specs/transport-spec-draft.md` (TE-26 / TE-27 outer rule)

- **C-1:** Honored. The pCID-keyed transport directory convention does not assume any registry. Each transport-protocol's pCID is mintable from its spec's bytes; no central authority blesses it.
- **C-2:** Mostly honored. The four locked principles do not depend on any short-horizon institution. The slug component of `transports/<pcid>--<slug>/` is a human-readable convenience that may drift in meaning over decades; this is acknowledged in the spec as presentational rather than load-bearing.
- **C-3:** Honored. The "code-as-handler" principle means the handler is the spec's executable interpretation; an adversarial reader cannot trick the system by claiming a different interpretation, because the spec's pCID covers the canonical bytes from which the handler is derived.
- **C-4:** Partially honored. The per-axis meta-rule (visibility, routing topology) gives a discipline for when distinct pCIDs are warranted, but does not address forks of a single protocol over time. **Open question: when an existing transport-protocol is forked (incompatible v2 of group-transport, or a community-specific variant), how does the new pCID get expressed in `transports/`?**
- **C-5:** Out of scope for the outer transport-spec. Trust over message receipts is a per-protocol concern; the transport-spec is silent.
- **C-6:** Honored. The pCID is the only anchor; the slug is convention.

### `protocols/group-session.d/specs/group-session-draft.md` (TODO 013 carve-out)

- **C-1:** Honored. The protocol's identity is its pCID, computed from the spec's bytes. Membership of any individual transport instance is closed and named by the slug, not by a registry.
- **C-2:** Mostly honored. The append-only persistence property and the "membership change creates a new transport instance" rule mean the transport's history survives any individual member's departure or death. The `Date:` header in messages is informational; the protocol does not depend on synchronized clocks across the horizon.
- **C-3:** Partially honored. The v0 contract leaves identity (`From:`) as a free-form string, with cryptographic signing of `From:`, message bodies, and message CIDs explicitly deferred. **Open question recorded in OQ-G1:** when the protocol gains cryptographic signing, how does the v0-to-v1 migration honor C-2 (existing v0 transports survive) and C-4 (v0 and v1 are siblings, not parent-and-child)?
- **C-4:** Acknowledged. The spec recognizes that "different transport-protocol classes (ring, gossip, hub-mediated, large-N, ephemeral, etc.) will produce their own envelope decisions in their own spec docs." But what about forks of the *same* class — incompatible v2 of group-transport-protocol? The spec is silent. **Open question for follow-on TE.**
- **C-5:** Partially honored. The body-level acknowledgement scheme records *what* was acknowledged but not the trust score behind the acknowledgement. Trust is a reader-side concern. **Open question recorded as OQ-G2:** cumulative-prefix or frontier-style ack semantics under the DAG model. C-5 implies these semantics must be local-trust-vector-aware, not global.
- **C-6:** Honored. The protocol's identity is its pCID; signing of individual messages is a v1 concern.

### TE-24 (group-transport envelope) and DR-009 / TODO 012

- These are decision records for the group-transport contract. The constraints apply transitively: anything the group-transport spec gets right or gets wrong about the constraints is reflected here.
- **Specific honor of C-2:** the rename TODO 013 just performed (channel→transport, channel-carrier→group-transport-envelope) demonstrates that the wire-lab is willing to refactor decision-record vocabulary when a clearer framing emerges. That willingness is a 100-year discipline: contributors arriving in 2046 should not have to decode 2026 vocabulary that turned out to be wrong.

### TE-25 (numbering collision), TE-26 (transport-protocol types), TE-27 (transports rename)

- These are repo-mechanics and outer-spec TEs. C-2 dominates: future-readers must be able to follow the numbering and rename mechanics without living institutional memory.
- **TE-25** locks the drafting-time invariant for TE numbering. This works under C-2 because the timestamp slug is content-addressable in the same sense as a pCID — any reader can recover the intended numbering by reading the timestamps.
- **TE-26** locks four principles for `transports/`. Honored under C-1, C-3, C-6. Partial under C-4 (silent on forks).
- **TE-27** locks the transports rename and the per-axis meta-rule. Honored under C-1, C-3, C-6. The per-axis rule indirectly addresses C-4 by giving a discipline for *when* a fork is warranted vs. when a parameter suffices.

### Branching policy and ppx/main convention

- All bot work happens on `ppx/<twig>` branches that merge into `ppx/main`. ppx/main is intended as a long-lived integration branch.
- **C-2:** This convention assumes the existence of `ppx/main` and the wire-lab repo on GitHub. Both of these are short-horizon institutions. The convention is appropriate for current work but is presumed transient over the 100-year horizon. No action needed; the convention is honestly scoped.
- **C-6:** The branching policy is convention, not a structural lock. Steve's signing key is the structural lock; the branches are bookkeeping. This is consistent.

### TE numbering and TODO numbering

- Both are integer sequences anchored on first-drafted timestamps (TE-25). Under C-1, no registry is needed; the integers are derived from sortable timestamps.
- Under C-2, the sequences will eventually wrap or become unwieldy over the horizon. **Open question for very-long-term TE: is the integer sequence stable across centuries, or does it eventually need to be supplanted by purely timestamp-based or pCID-based identifiers?**

## Conclusions

1. **The 100-year goal is real and is already shaping the locked design** — pCIDs, content-addressed specs, per-burden trust ledgers, multi-generational handoff scenarios. The harness-spec's existing "multiple human generations" and "centuries-long" prose is a soft surface on a hard underlying constraint set.
2. **Naming the constraint set explicitly costs almost nothing and prevents short-horizon assumptions** from being smuggled in by future-LLM-readers, future-contributors, or future-Steve under time pressure. The current implicit framing is sufficient when Steve is reading every PR; it is insufficient over the horizon the design is for.
3. **Several open questions surface from the pressure-test** and warrant follow-on TEs:
   - **OQ-100.1 — Protocol forking:** when an existing transport-protocol is forked (incompatible v2, or a community-specific variant), how is the new pCID expressed in `protocols/` and `transports/`? Cross-references to TE-27 per-axis meta-rule.
   - **OQ-100.2 — Cryptographic signing migration:** when v0 group-transport gains cryptographic signing in v1, how does the migration honor C-2 (existing v0 transports survive) and C-4 (v0 and v1 are siblings, not parent-and-child)?
   - **OQ-100.3 — Cumulative-prefix ack and trust:** the deferred Q2 from TODO 013 (cumulative-prefix or frontier ack semantics) must be designed under C-5 — trust-vector-aware, not global.
   - **OQ-100.4 — Numbering wrap:** is the integer sequence for TEs and TODOs stable across centuries, or does it eventually need to be supplanted by purely timestamp-based or pCID-based identifiers?
   - **OQ-100.5 — Slug drift:** human-readable slugs in `transports/<pcid>--<slug>/` and `protocols/<slug>.d/` are presentational, but a 30-year-old slug whose meaning has drifted may mislead a future contributor more than it helps. Is there a discipline for retiring or redirecting drifted slugs without changing pCIDs?
4. **No current lock fails a 100-year constraint outright.** Several locks are partial-honor or silent-on a constraint; those become open questions to address in subsequent TEs. The current design is consistent with the 100-year goal as far as it goes; the goal does not invalidate any prior decision.

## Recommendation

Adopt the 100-year goal as a named, citable constraint set. Specifically:

- **Add a `## The 100-Year Goal` subsection to `protocols/wire-lab.d/specs/harness-spec-draft.md`** (placed between the abstract and §1 to make it visible before any concrete spec content) that states the goal in one paragraph, lists the six constraints C-1 through C-6 by name, and points readers at this TE for the pressure-test.
- **Cite this TE** from the §8 TE bibliography in the harness-spec.
- **Reference the constraint codes (C-1 through C-6) from future TEs** when those TEs are pressure-testing alternatives. This gives the wire-lab a shared shorthand and prevents each TE from re-deriving the constraint set.
- **Open the five open questions** (OQ-100.1 through OQ-100.5) into the harness-spec's open-questions list (currently §12), each as its own DR or each as a candidate for a follow-on TE.

## Implications for follow-on work

- The anticipated TE on simulated-repo release machinery (release-as-tree-hash, two-step ceremony, Merkle bundle, release record file, etc.) must explicitly pressure-test each alternative against C-1 through C-6 before recommending. Several alternatives that look attractive at a one-month horizon fail at the 100-year horizon and need to be dismissed accordingly.
- The anticipated TEs on ring, gossip, hub-mediated, and large-N transport-protocols must pressure-test their proposed envelopes against C-3 (adversarial-by-default) and C-5 (per-burden trust) more carefully than v0 group-transport had to, because each new protocol adds an attack surface and a per-burden trust dimension.
- Future spec rewrites (e.g., a v1 group-transport that adds cryptographic signing) inherit the OQ-100.2 migration discipline.

## Decision status

`recommended for immediate adoption` — the constraint set is articulated; the harness-spec edit and the open-question additions are mechanical. No DR/DI is required because no decision is locked beyond what was already implicit; this TE makes the implicit explicit.
