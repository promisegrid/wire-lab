# TE-hirap: Artifacts as PromiseGrid messages

*Thought experiment: in the spirit of TE-pahah and TE-vilot, should durable wire-lab artifacts be shaped like PromiseGrid messages, and how should CBOR and plain-text representations coexist while the canonical format is still unsettled?*

## TE ID

TE-hirap

## Status

needs DF

First drafted 2026-05-09 19:39:04 UTC.

## Decision under test

DUT-hirap: Should all durable artifacts in wire-lab -- TEs, DRs, DIs, TODOs,
specs, guide-resource notes, simulation events, observations, results,
messages, and review artifacts -- be shaped like PromiseGrid messages?

This is a stronger question than TE-gurov's "should every artifact be a
PT-style promise?" A PromiseGrid-message-shaped artifact would usually have:

- an explicit protocol selector, ideally a pCID,
- content-addressed identity over canonical bytes,
- author/promiser/provenance fields,
- parent links or other graph edges,
- a payload/body governed by the selected protocol,
- and possibly evidence such as signatures, channel observations, or review
  attestations.

The representation question is load-bearing:

- **CBOR / structured promise-stack form** is closer to a future machine-first
  wire format.
- **Plain text / `grid <pcid>` form** is closer to the current group-session
  experiment and remains readable in git diffs, chat, and LLM context.

TE-pahah and TE-vilot change the framing. The best question is not "should the
whole repository become a PromiseGrid message store tomorrow?" The better
question is: should wire-lab use named simulations to test message-shaped
artifact protocols, while keeping apparatus documents honest and readable until
evidence supports a broader template?

## Short answer

Do not globally reshape every artifact as a PromiseGrid message yet.

Recommended synthesis:

1. Treat **PromiseGrid-message-shaped artifacts** as a simulation/specimen
   candidate, not a repo-wide apparatus rule.
2. Use **plain-text message shape** first for human/LLM-readable simulations and
   reviewable corpus artifacts:

   ```text
   grid <artifact-protocol-pCID-or-draft-ref>

   Date: 2026-05-09T19:39:04Z
   From: alice@example.net (Alice)
   Parents: <artifact-cid> <artifact-cid>

   I promise ...
   ```

3. Use **CBOR / promise-stack form** as the canonical production-wire candidate
   to be tested by tools and simulations, not as the default authoring format for
   TEs, DRs, TODOs, or guide notes.
4. Do not maintain two equally canonical byte encodings for the same artifact
   unless a later protocol defines an identity rule that prevents CID split. Pick
   one canonical byte representation per protocol and treat other views as
   derived, diagnostic, or quoted.
5. Keep apparatus documents -- TEs, DRs, TODOs, and guide-resource notes -- in
   role-specific Markdown until a simulation result and DI justify changing
   templates.

This extends TE-vilot: promise-shaped simulation artifacts can be even stronger
than `I promise` prose. They can be full message-shaped objects. But the place
to prove that is still `simulations/<sim>/`, not the entire repo.

## Assumptions

- The current PromiseGrid wire format is not frozen.
- The harness spec describes a future direction where messages are promise
  stacks and the on-wire encoding may be a CBOR array of promise frames.
- TE-hogus and the group-session draft lock a plain-text `grid <pcid>` envelope
  for the current group-session-style experiment, not for all future
  PromiseGrid messages.
- A pCID identifies a protocol spec, not a specific payload or artifact.
- A message CID or artifact CID identifies canonical bytes. If the same logical
  artifact has two canonical encodings, it has two byte identities unless a
  higher protocol explicitly defines equivalence.
- Human and LLM readers need direct readability while design work is still in
  progress.
- Machine readers and future kernels need unambiguous canonical bytes,
  structured fields, and evidence semantics.
- Mallory can exploit visual similarity between apparatus files and message
  specimens by making an experimental object look like an adopted decision.

## Threat / trust model

Alice is a human maintainer reviewing ordinary repo changes. Bob is an LLM agent
recovering context from TEs, DRs, TODOs, simulation logs, and message specimens.
Carol publishes conformance claims. Dave runs dogfood simulations. Ellen writes
the PromiseGrid Development Guide. Frank operates large simulation corpora.
Mallory tries to exploit representation ambiguity, forged provenance, generated
views, and visually convincing `grid <pcid>` text to make specimens look like
decisions or alternate encodings look like the same artifact.

The design succeeds if cooperative actors can distinguish apparatus from
specimen, canonical bytes from derived views, pCID from message/artifact CID, and
provenance from agency. It fails if a ritual message shape creates false
authority, or if text/CBOR duality splits identity without an explicit protocol
rule.

## Alternatives

### Alt A -- Global CBOR message corpus

Every durable artifact is stored as a CBOR PromiseGrid message or promise stack.
Markdown views are generated for humans but are not canonical.

**Easier:** Strongest machine discipline. Canonical bytes, type tags, nested
promise frames, evidence fields, and pCID dispatch can be exact. The repo would
dogfood the future machine-oriented wire format aggressively.

**Harder:** Humans and LLMs lose direct readability. Git diffs become weak
without specialized tooling. Simple review, grep, merge conflict repair, and
long-horizon archaeological reading all get harder.

**Obligations:** Build encoders, decoders, validators, pretty-printers,
round-trip tests, and recovery procedures before ordinary coordination can
continue. Define how DR/DI authority appears in structured fields without
rewriting the source-of-truth process prematurely.

### Alt B -- Global plain-text `grid <pcid>` corpus

Every durable artifact becomes a canonical UTF-8 text message with a first-line
protocol selector, ordered headers, parent CIDs, and a body. Each artifact class
gets a protocol: TE protocol, DR protocol, TODO protocol, DI protocol, spec
protocol, simulation-result protocol, and so on.

**Easier:** Human-readable, LLM-readable, diffable, and close to existing
group-session specimens. Protocol selection and parent links become visible in
ordinary files.

**Harder:** It still makes apparatus documents look like messages. A DR is an
open question, not necessarily a message in a runtime transport. A TODO is a
mutable queue snapshot, not necessarily a promise stack. Plain text also makes
canonicalization fragile: whitespace, header order, line endings, and body edits
all affect CIDs.

**Obligations:** Define artifact protocols for every file class, migration
rules, header sets, parent-link semantics, and canonicalization policy. Decide
whether the body must start with `I promise`, contain explicit promise blocks, or
only follow its artifact protocol.

### Alt C -- Dual canonical CBOR and plain-text representations

Every artifact has both a plain-text form and a CBOR form. Both are considered
canonical and both may circulate.

**Easier:** Humans get text; machines get structure.

**Harder:** This is the identity trap. If both encodings are canonical bytes,
they produce different CIDs for the same logical artifact. Receipts, parent
links, signatures, and citations can split across representations. Even a
lossless text-to-CBOR mapping needs one representation to be the identity anchor,
or the protocol must define a higher-level equivalence relation.

**Obligations:** Define which CID is authoritative for graph links and
signatures, how derived views are authenticated, and how a receiver detects that
two encodings are the same logical object rather than two conflicting objects.

### Alt D -- Simulation-scoped message-shaped artifacts

Keep repo apparatus ordinary. Inside named simulations, define artifact/message
protocols that can use plain text, CBOR, or both under explicit rules. A
simulation can compare text-authored messages, CBOR-authored messages, and
generated views without implying that every repo file has adopted that shape.

**Easier:** This follows TE-pahah and TE-vilot. The design can be tested in
bounded worlds with actors, sites, groups, CAS, feeds, wires, messages,
observations, and results. Failures are evidence, not repo-wide breakage.

**Harder:** The corpus has multiple styles. Readers must distinguish apparatus
from specimen and know that `simulations/<sim>/` may contain stronger message
objects than the surrounding docs.

**Obligations:** Each simulation that uses message-shaped artifacts must state
its artifact protocol, canonical byte representation, derived-view rules, and
identity/citation policy.

### Alt E -- Commitment and specimen artifacts only

Only artifacts that already behave like commitments or runtime specimens become
PromiseGrid-message-shaped: transport messages, simulation events,
conformance/adoption claims, review replies, and frozen specs. TEs, DRs, TODOs,
and guide notes remain conventional Markdown.

**Easier:** Strong semantic honesty. Runtime-like objects get runtime-like
message form; analysis and decision apparatus keep their existing roles.

**Harder:** It may under-test whether a broader message-shaped corpus would help
PromiseGrid coordination. The boundary between "commitment" and "apparatus" can
be porous.

**Obligations:** Define a classification rule for commitment/specimen artifacts
and a process for promoting an artifact class into message shape.

### Alt F -- Status quo plus local examples

Do not adopt message-shaped artifact conventions. Continue writing Markdown
docs and plain group-session messages case by case.

**Easier:** No migration and no template churn.

**Harder:** The repo misses a chance to test a central PromiseGrid idea against
its own coordination artifacts. Message-shape questions remain abstract.

**Obligations:** Keep documenting why examples are local and not general.

## Representation model

The key distinction is logical message versus byte representation.

### Logical shape

A logical artifact-as-message could be described as:

```text
ArtifactMessage := {
  protocol: pCID-or-draft-ref
  artifact_class: TE | DR | DI | TODO | spec | event | observation | result | review | message
  cid: hash(canonical-bytes)
  parents: zero-or-more artifact/message CIDs
  actor_fields: promiser, promisee, author, recorder, drafter, model, reviewer
  evidence: signatures, commit provenance, channel observation, review attestation
  payload: bytes | markdown | cbor | nested ArtifactMessage | promise-stack
}
```

This logical shape does not require one universal syntax. Different protocols
can choose different encodings. The pCID names those rules.

### Plain-text representation

Plain text is best when the artifact must be authored, reviewed, quoted, and
searched by humans and LLMs:

```text
grid draft:artifact-message

Date: 2026-05-09T19:39:04Z
From: alice@example.net (Alice)
Parents: bafk...
Artifact-Class: observation
Canonical-Form: text

I promise that this observation records the result of simulation SIM-hirap-demo
under the stated protocol set.
```

Plain text makes the protocol visible and keeps the body legible. Its cost is
canonicalization discipline. If the filename is a CID over canonical bytes, then
header order, blank lines, line endings, and body bytes are load-bearing.

### CBOR representation

CBOR is best when the artifact is intended for machine transport, compact
storage, signatures, binary payloads, or nested promise stacks:

```text
[
  {
    "p": <artifact-protocol-pCID>,
    "from": <agent-id>,
    "parents": [<cid>, <cid>],
    "claim": "payload conforms to p",
    "evidence": <signature-or-channel-proof>,
    "body": <bytes-or-nested-promise-stack>
  }
]
```

The example is illustrative, not a schema. A real CBOR form should be defined by
a spec, likely with deterministic encoding rules. Its cost is opacity in ordinary
git review unless the repo also carries verified text views.

### Derived views

If one representation is canonical, other representations can still be useful:

- A CBOR object can have a generated Markdown rendering for review.
- A plain-text artifact can have a generated CBOR projection for simulator
  ingestion.
- A simulation can record both the canonical object CID and the derived view CID.

But derived views must not silently become alternate identities. If a text view
is derived from CBOR, receipts should cite the CBOR CID unless the protocol says
otherwise. If CBOR is derived from text, signatures should bind the text CID
unless the protocol says otherwise.

## Scenario analysis

### S1: Alice reviews a TE in a normal Git diff

Alice is a human maintainer reviewing a new thought experiment.

- **Alt A:** Alice needs a CBOR viewer before she can review ordinary reasoning.
  This blocks the current low-friction TE workflow.
- **Alt B:** Alice can read the TE directly, but the TE now looks like a runtime
  message. She must understand the TE protocol before she understands the
  analysis.
- **Alt C:** Alice sees a generated text view but must know whether she is
  reviewing the view, the CBOR object, or both.
- **Alt D:** Alice reviews apparatus TEs as Markdown and reviews
  message-shaped specimens inside simulations when the experiment requires it.
- **Alt E:** Same as Alt D for TEs; runtime specimens still get message shape.
- **Alt F:** Easy review, but no broader experiment.

S1 favors Alt D or Alt E. Global CBOR is the wrong default for human-facing
analysis at this stage.

### S2: Bob, an LLM agent, needs to reconstruct evidence

Bob is asked why a design decision was made. He searches TEs, DRs, simulation
results, and message specimens.

- **Alt A:** Bob needs tools to decode CBOR before searching. If tooling is good,
  structured fields help; if tooling is missing, context recovery degrades.
- **Alt B:** Bob can grep and quote directly. Parent CIDs and protocol selectors
  help, but boilerplate may dominate search results.
- **Alt C:** Bob may confuse generated views with canonical artifacts unless
  view status is explicit.
- **Alt D:** Bob can read apparatus normally, then inspect simulation specimens
  whose protocols explain whether text or CBOR is canonical.
- **Alt E:** Bob has fewer message-shaped examples but a clearer boundary.
- **Alt F:** Bob gets ordinary docs but less protocol discipline.

S2 favors Alt D because it supports both human/LLM reconstruction and structured
experiments.

### S3: Carol signs a conformance claim

Carol maintains an implementation and publishes a claim that it conforms to a
frozen protocol pCID.

- **Alt A:** CBOR is attractive: Carol signs a deterministic object with exact
  fields.
- **Alt B:** Plain text is also attractive if humans need to review the claim and
  the canonical bytes are stable.
- **Alt C:** Dual canonical claims risk signature split.
- **Alt D:** A simulation can compare text-signed and CBOR-signed claims before
  freezing a conformance-claim protocol.
- **Alt E:** Conformance claims are exactly the kind of artifact that should
  become message-shaped early.
- **Alt F:** Status quo leaves claim shape under-specified.

S3 favors Alt E for commitment artifacts, with Alt D as the way to choose text
versus CBOR before freezing.

### S4: Dave runs a dogfood communication simulation

Dave wants real message files, receipts, missing-message detection, and
cross-representation tests.

- **Alt A:** CBOR is close to production, but difficult for manual repair and
  chat-based collaboration.
- **Alt B:** Plain text is easy to author and review. It exercises pCID,
  message-CID, parent links, and promise bodies directly.
- **Alt C:** Useful as a deliberate test, but only if the simulation explicitly
  models representation equivalence and CID choice.
- **Alt D:** Best fit. Dave can run a text-message world, a CBOR-message world,
  and a dual-view world side by side under `simulations/`.
- **Alt E:** Works if dogfood message files are classified as specimens.
- **Alt F:** Does not test enough.

S4 strongly favors Alt D.

### S5: Ellen writes the PromiseGrid Development Guide

Ellen needs to explain what developers should implement.

- **Alt A:** Ellen might prematurely tell developers that PromiseGrid artifacts
  are CBOR promise stacks, even though that is not frozen.
- **Alt B:** Ellen might prematurely teach `grid <pcid>` text as the canonical
  PromiseGrid wire format, even though TE-hogus scoped it to the first
  group-session experiment.
- **Alt C:** Ellen faces the hardest explanation: two canonical encodings with
  identity rules not yet settled.
- **Alt D:** Ellen can say wire-lab is testing representations in named
  simulations; only frozen specs by pCID become normative.
- **Alt E:** Ellen can discuss message-shaped conformance/adoption artifacts
  once those protocols freeze.
- **Alt F:** Ellen has less evidence to cite.

S5 favors Alt D and explicitly warns against guide prose that treats any
representation as final before a frozen spec exists.

### S6: Mallory exploits representation ambiguity

Mallory publishes a text object and a CBOR object that claim to be the same
artifact but differ in a load-bearing field.

- **Alt A:** There is only CBOR if the rule is followed, so ambiguity is lower;
  but humans may rely on an unauthenticated text rendering.
- **Alt B:** There is only text if the rule is followed, so ambiguity is lower;
  but parsers must reject non-canonical text reliably.
- **Alt C:** This is dangerous unless equivalence, canonical identity, and
  signature rules are exact.
- **Alt D:** A simulation can test this attack explicitly before adopting
  dual-view tooling.
- **Alt E:** Commitment artifacts can choose one canonical representation per
  protocol.
- **Alt F:** Ambiguity remains informal.

S6 rejects dual canonical representation as a default. It is valid only as an
explicit simulation or a carefully specified protocol.

### S7: Frank operates at scale

Frank runs thousands of simulations and stores large artifacts in CAS.

- **Alt A:** CBOR and CAS work well at scale. Chunking, binary payloads, and
  deterministic parsing are natural.
- **Alt B:** Plain text scales operationally less well for large binary bodies,
  but it remains excellent for small coordination messages and summaries.
- **Alt C:** Dual canonical storage doubles indexing and validation burden.
- **Alt D:** Simulations can use CBOR for large/machine-heavy worlds and text for
  coordination-heavy worlds.
- **Alt E:** Commitment/specimen artifacts can pick representations by need.
- **Alt F:** Scale questions stay deferred.

S7 favors representation choice by protocol, not one global artifact syntax.

## Rejected alternatives

- **Reject Alt A as the global default now.** CBOR may be the right production
  wire candidate, but using it for all repo coordination artifacts would make the
  current lab less inspectable by humans and LLMs before the required tooling and
  canonical schemas exist.
- **Reject Alt B as the global default now.** Plain-text `grid <pcid>` is strong
  prior art for group-session-style specimens, but making every TE, DR, TODO, and
  guide note a text message would blur apparatus/specimen boundaries.
- **Reject Alt C as an unqualified default.** Dual canonical text and CBOR forms
  create CID and signature split unless one representation is the identity anchor
  or an explicit equivalence protocol exists.
- **Reject Alt F as sufficient.** Status quo avoids churn but does not test
  artifact-as-message enough to inform PromiseGrid design.

## Surviving alternatives

- **Alt D survives as the strongest candidate.** Simulation-scoped
  artifact-message protocols let wire-lab test text, CBOR, and derived views
  without turning every apparatus file into a message.
- **Alt E survives as the likely first graduation path.** Commitment and specimen
  artifacts are the safest classes to promote if simulations show value.
- **Alt C survives only inside explicit simulations or frozen protocols.** It is
  worth testing text/CBOR duality, but only with a stated canonical-identity
  rule.

## Reconsideration of TE-pahah

TE-pahah becomes more important when the artifact shape becomes message-like.
Without a simulation boundary, "all artifacts are PromiseGrid messages" turns
the repo itself into a half-built PromiseGrid node. That prematurely mixes:

- apparatus used to reason,
- specimens used to test,
- results derived from a test,
- decisions locked by DR/DI,
- and future production/reference layout.

With `simulations/`, the repo can create bounded worlds that answer concrete
questions:

```text
simulations/
  SIM-<handle>-artifact-message-shape/
    README.md
    QUESTION.md
    protocol-set.md
    artifact-message-protocol.md
    representation-policy.md
    world/
      actors/
      sites/
      groups/
      cas/
      feeds/
      wires/
      artifacts/
        text/
        cbor/
    events/
    observations/
    results/
    decisions.md
```

That structure lets wire-lab test text-only, CBOR-only, and derived-view models
side by side without turning top-level TEs and DRs into runtime messages.

## Reconsideration of TE-vilot

TE-vilot said promise-shaped artifacts should be tested inside simulations and
that apparatus documents should stay ordinary unless they contain genuine
promise blocks. TE-hirap strengthens that recommendation.

Message shape is stronger than promise prose. It adds pCID dispatch, canonical
bytes, parent links, artifact/message CIDs, evidence, and representation rules.
Therefore it is even less appropriate as a blanket rule for all apparatus
documents before a simulation proves it.

TE-vilot's "simulation boundary" should be read as the right boundary not only
for `I promise` wrappers, but also for:

- text `grid <pcid>` artifact-message envelopes,
- CBOR promise-stack artifact objects,
- dual-view tooling,
- representation-equivalence tests,
- and message-shaped conformance/review/adoption claims.

## Reconsideration of TE-hogus and group-session

TE-hogus selected `grid <pcid>` plain text for the wire-lab's first
group-session-style transport protocol because the project needed readable
protocol-selection-by-pCID, parent-CID references, and explicit promise bodies.

That decision should not be generalized too far. It proves that plain text is
excellent for current human/LLM-centered experiments. It does not prove that:

- all repo artifacts should be text messages,
- the final PromiseGrid wire format should be text,
- CBOR should be rejected,
- or every document should start with a carrier line.

The better lesson is narrower and more useful: plain text is a good specimen
format for early simulations; CBOR remains a strong candidate for production
wire/storage protocols; both should be tested under explicit pCID-governed
protocols.

## Recommended conclusion

Adopt Alt D as the active experimental posture, with Alt E as the likely first
graduation path:

1. Do not reshape all current artifacts as PromiseGrid messages.
2. Create message-shaped artifact protocols inside named simulations when the
   decision under test requires them.
3. Prefer plain text for early human/LLM-authored specimens and CBOR for
   machine-heavy or production-wire simulations.
4. Never treat both text and CBOR as equally canonical unless a protocol defines
   an explicit identity/equivalence rule.
5. Promote commitment/specimen classes first: conformance claims, adoption
   claims, review replies, simulation results, and runtime messages.
6. Keep TEs, DRs, DIs, TODOs, and guide notes as apparatus Markdown until a later
   DI changes their templates.

This preserves the insight behind "artifacts as messages" while keeping the
repo usable as a lab.

## DF questions exposed

### DF-hirap.1 -- Should artifact-as-message be tested inside `simulations/` first?

Recommended answer: yes.

Surviving alternatives:

- **1.A -- Simulation-first message shape (recommended).** Test artifact-message
  protocols inside named simulations before repo-wide adoption.
- **1.B -- Global text message corpus.** Convert all new durable artifacts to
  plain-text `grid <pcid>` style after a DI.
- **1.C -- Global CBOR message corpus.** Store all new durable artifacts as CBOR
  promise-stack objects after a DI.
- **1.D -- Commitment/specimen only.** Apply message shape only to runtime-like
  and commitment artifacts.

### DF-hirap.2 -- Which representation should early simulations prefer?

Recommended answer: plain text first, CBOR as a parallel candidate.

Surviving alternatives:

- **2.A -- Plain text first (recommended).** Use text for early authoring,
  review, LLM inspection, and small coordination specimens.
- **2.B -- CBOR first.** Use structured promise-stack objects from the start and
  require generated human views.
- **2.C -- Bakeoff.** Run paired text and CBOR simulations with the same logical
  artifact protocol.

### DF-hirap.3 -- How should dual representations handle identity?

Recommended answer: one canonical representation per protocol.

Surviving alternatives:

- **3.A -- One canonical CID per protocol (recommended).** Pick text or CBOR as
  canonical; other views are derived and cite the canonical CID.
- **3.B -- Two canonical CIDs with explicit equivalence.** Allow both only if a
  protocol defines equivalence, signatures, receipts, and conflict handling.
- **3.C -- No dual representation.** Forbid derived views in early experiments
  to reduce ambiguity.

### DF-hirap.4 -- Should apparatus docs get carrier lines now?

Recommended answer: no.

Surviving alternatives:

- **4.A -- No carrier lines for apparatus (recommended).** TEs, DRs, DIs, TODOs,
  and guide notes stay Markdown unless a later template decision changes them.
- **4.B -- Optional carrier lines for new apparatus.** Allow experimental
  `grid <pcid>` headers in new docs after a protocol exists.
- **4.C -- Mandatory carrier lines.** Require all new apparatus docs to become
  message-shaped after a DI.

### DF-hirap.5 -- Which artifact classes should graduate first if experiments succeed?

Recommended answer: commitments and specimens.

Surviving alternatives:

- **5.A -- Commitment/specimen classes first (recommended).** Runtime messages,
  simulation events/results, conformance claims, adoption claims, and review
  replies.
- **5.B -- All simulation artifacts first.** Everything inside a named simulation
  uses message shape.
- **5.C -- All repo artifacts.** Every durable file class gets an artifact
  protocol.

## Implications for open TODOs, DRs, and DIs

- **TE-pahah:** strengthened. `simulations/` is the right place to test
  artifact-message protocols and representation policies.
- **TE-vilot:** strengthened. Its simulation boundary applies to full
  PromiseGrid-message shape, not only to `I promise` wrappers.
- **TE-gurov:** narrowed further. Promise prose is only one layer; message shape
  also needs pCID, CID, parent links, evidence, and canonical representation
  decisions.
- **TE-hogus / DR-009:** remain scoped to group-session-style text messages.
  Their `grid <pcid>` decision is useful prior art, not a repo-wide artifact
  template.
- **DEV-GUIDE-RESOURCES.md:** guide writers should treat CBOR, text envelopes,
  and artifact-message protocols as unsettled wire-lab evidence until frozen
  specs by pCID exist.
- **Future simulation TODO:** if DF-hirap.1 is accepted, file a TODO for an
  artifact-message-shape simulation that compares text, CBOR, and derived-view
  representations.
- **Future DIs:** any DI adopting message-shaped artifacts must name the
  artifact classes, canonical representation, identity rule, and whether
  generated views are authoritative or diagnostic.

## Decision status

`needs DF`. This TE recommends simulation-scoped artifact-message experiments,
plain text first for human/LLM-centered specimens, CBOR as a structured
production-wire candidate, and one canonical representation per protocol. It
does not lock any repo-wide artifact template.
