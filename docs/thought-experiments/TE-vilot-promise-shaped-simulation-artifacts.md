# TE-vilot: Promise-shaped artifacts inside simulation-first wire-lab

*Thought experiment: reconsider TE-pahah and TE-gurov together, asking whether the simulation-first structure gives wire-lab the right boundary for testing PT-style promise artifacts without forcing every coordination document to become one.*

## TE ID

TE-vilot

## Status

needs DF

First drafted 2026-05-09 19:28:10 UTC.

## Decision under test

DUT-vilot: If TE-pahah is right that wire-lab should be simulation-first, how
should TE-gurov's "all durable artifacts as PT-style promises" proposal be
reinterpreted?

The combined question is not merely "should every file start with `I promise`?"
It is:

1. which artifacts are part of the **wire-lab apparatus** used to reason and
   decide,
2. which artifacts are **specimens inside a simulation world** that test
   PromiseGrid behavior,
3. which artifacts are **commitments** such as specs, adoption claims,
   conformance records, or review replies,
4. and where human, repo-role, service, agent, and LLM-model identity should be
   represented without confusing provenance with agency.

Steve's positive read of TE-pahah matters. If `simulations/` becomes the primary
home for concrete experimental worlds, then promise-shaped artifacts no longer
need to be answered as a repo-wide formatting mandate. They can be tested inside
explicit simulation worlds, while apparatus documents stay honest about their
roles as questions, analysis, queues, and decision records.

## Short answer

TE-pahah strengthens TE-gurov by providing the missing boundary.

The recommended synthesis is:

- Keep apparatus documents such as TEs, DRs, TODOs, and guide-resource notes in
  their current role-specific forms.
- Permit explicit promise blocks or short role promises in apparatus documents
  only when the speech act is genuinely a promise.
- Use `simulations/<sim>/` to test promise-shaped artifacts as specimens:
  messages, event records, result claims, adoption claims, conformance records,
  and actor commitments can all carry explicit promiser/provenance fields and
  `I promise` bodies when the simulation is designed to test that convention.
- Do not globally rewrite every durable document as a PT promise until a
  simulation has produced evidence and a DI locks a template.

This makes TE-gurov's conservative Alt-C the apparatus default and TE-gurov's
wrapper Alt-B a simulation/specimen candidate. It preserves TE-pahah's core
claim: wire-lab is the experimental simulation space for deriving PromiseGrid
design choices, not the final PromiseGrid node layout and not the dev guide.

## Assumptions

- TE-pahah is directionally correct: concrete experimental worlds should live
  under `simulations/`, and root-level `transports/` versus `groups/` is the
  wrong primary axis.
- TE-gurov is also directionally useful: promise form exposes agency,
  provenance, and scope, but strict universal first-person prose is dishonest
  for questions, analyses, and mutable work queues.
- Existing wire-lab DR/DI rules remain source-of-truth for decisions and open
  questions until a later DI supersedes them.
- An LLM model can be provenance. It is not automatically the promiser unless a
  separate decision defines a durable delegated agent identity.
- The PromiseGrid Development Guide should eventually explain PromiseGrid, not
  wire-lab process mechanics, except where process evidence supports a settled
  PromiseGrid claim.
- Mallory can exploit ambiguous file shapes by making a specimen look like a
  settled repo-wide rule, or by making a generated draft look like an
  accountable promise.

## Alternatives

### Alt A -- Global promise-shaped corpus

Every new durable artifact, including TEs, DRs, TODOs, README files, simulation
records, specs, guide-resource notes, and review artifacts, carries explicit
identity metadata and starts its body with an artifact-specific `I promise`
wrapper.

**Easier:** The repo has one visible philosophical rule. A reader can assume
every artifact is an agent's promise about its role.

**Harder:** The apparatus/specimen boundary becomes blurry. A TE that promises
to analyze is not the same thing as a simulated actor promise. A DR that promises
to record a question is not the same thing as an adopted answer. The uniform
opening risks making process files feel like protocol specimens.

**Obligations:** Define templates for every artifact class, migrate old files or
grandfather them, and explain why questions and TODO queues are "promises"
without misleading guide writers.

### Alt B -- Simulation-scoped promise specimens, ordinary apparatus

Keep apparatus documents ordinary and role-specific. Use `simulations/<sim>/` as
the place to test promise-shaped documents and messages. A simulation may define
its own artifact protocol: for example, event records may start with `I promise`
and name `Promiser`, `Recorder`, `Drafter`, `Model`, `Promisee`, `Scope`, and
`Evidence`.

**Easier:** Wire-lab can experiment with strong promise conventions without
pretending those conventions are already repo-wide law. Human and LLM readers can
see that a promise-shaped artifact is part of a named simulation and can inspect
that simulation's question, protocol set, observations, and results.

**Harder:** There are two visible styles: apparatus files and specimen files.
Writers need to know which side of the boundary they are editing.

**Obligations:** Each simulation that uses promise-shaped artifacts must state
its artifact protocol. If the experiment succeeds, a later TE/DF/DI can graduate
the convention into a broader template.

### Alt C -- Commitment-only promise form

Use promise form only for specs, adoption claims, conformance records, review
replies, and simulation result claims. TEs, DRs, TODOs, event logs, and ordinary
observations remain conventional Markdown unless they contain explicit promise
blocks.

**Easier:** Promise form stays close to real commitments. This is easiest to
teach and least likely to conflate analysis with adoption.

**Harder:** It does not fully test whether promise-shaped artifacts work as a
general substrate for PromiseGrid coordination. Some useful experiments may be
missed because the convention is too narrow.

**Obligations:** Maintain a clear boundary between commitments and evidence.
Define conformance/adoption/review templates before using them normatively.

### Alt D -- Pure simulation ontology

Everything in `simulations/<sim>/` is promise-shaped, including README files,
questions, events, observations, results, and decisions. Everything outside
`simulations/` remains ordinary.

**Easier:** It is a strong experiment. The repo can test "all artifacts are
promises" in a bounded world without changing the whole corpus.

**Harder:** It may still strain questions and observations. A simulation's
`QUESTION.md` might honestly promise to ask a question, but the useful content is
the question, not the promise. The experiment could produce formatting noise
rather than clearer agency.

**Obligations:** Pick one initial simulation and accept that it is deliberately
testing a strong convention, not defining the repo's default.

### Alt E -- Keep Pahah and Gurov separate

Adopt or reject `simulations/` independently from promise-shaped artifact
templates. Do not let one decision affect the other.

**Easier:** Fewer coupled decisions.

**Harder:** It misses the key design opportunity. Pahah provides a boundary that
makes Gurov safer and more testable. Keeping them separate leaves Gurov as an
abstract corpus-wide style debate.

**Obligations:** Later work will have to rediscover how promise-shaped artifacts
fit inside concrete simulations.

## Scenario analysis

### S1: Alice enters the repo to understand the current design

Alice is a new human developer. She starts at the top-level README, then opens
the TE index, then looks for the experiment that produced a claim about group
messages.

- **Alt A:** Alice sees many files starting with `I promise`. This is memorable,
  but it does not tell her whether she is reading apparatus, specimen, result, or
  adopted rule. The opening phrase becomes too generic.
- **Alt B:** Alice sees TEs and DRs as decision apparatus, then enters a named
  simulation where promise-shaped artifacts are clearly specimens under test.
  The boundary helps her avoid treating the simulation as the whole repo.
- **Alt C:** Alice can identify commitments quickly, but may not see enough
  promise-shaped examples to understand how a PromiseGrid app should behave.
- **Alt D:** Alice gets a bounded "all artifacts are promises" world, but may
  still have to learn why simulation questions and observations have promise
  wrappers.
- **Alt E:** Alice must mentally combine Pahah and Gurov herself.

S1 favors Alt B. It gives Alice a navigable corpus without losing the chance to
study promise-shaped artifacts.

### S2: Bob, an LLM agent, writes a simulation event

Bob observes a test run where Carol sends Dave a message through a sparse CAS
path. The event record is generated by an LLM session and reviewed by Steve.

- **Alt A:** If every artifact says `I promise`, Bob's generated first-person
  text risks implying that the model is the accountable promiser.
- **Alt B:** The simulation artifact protocol can separate `Promiser: Carol`,
  `Recorder: bob-agent`, `Model: <model-name-if-known>`, and `Reviewer: Steve`.
  The event body can contain Carol's actual promise while the file metadata
  records Bob's provenance.
- **Alt C:** The event stays conventional unless Carol's message contains a
  promise block. This is safe but less useful for testing artifact protocols.
- **Alt D:** Every simulation file tests the strong convention, which may be
  useful if the run's question is specifically "can all simulation artifacts be
  promises?"
- **Alt E:** Bob has no integrated guidance.

S2 favors Alt B for ordinary simulations and Alt D only for a simulation whose
explicit question is the strong convention itself.

### S3: Carol opens a DR after a simulation result conflicts with a TE

Carol's simulation says promise wrappers improve auditability. A prior TE says
strict promise wrappers hurt readability. Carol opens a DR asking which template
should be adopted.

- **Alt A:** The DR itself has a promise wrapper. That is tolerable only if the
  wrapper promises to record the unresolved question, not to settle it.
- **Alt B:** The DR remains an ordinary source-of-truth question, cites the
  simulation result, and may include an explicit promise block if someone commits
  to running another test.
- **Alt C:** Same as Alt B for the DR, but the simulation evidence may be thinner.
- **Alt D:** The simulation may overfit to its own convention; the DR still needs
  ordinary decision semantics outside the simulation.
- **Alt E:** The source-of-truth relationship remains unclear.

S3 says DR/DI apparatus should stay role-specific even when they discuss
promise-shaped simulation artifacts.

### S4: Dave dogfoods a real coordination channel

Dave wants a dogfood world where humans and LLMs exchange promise-shaped
messages. The world needs sites, groups, CAS, feeds, wires, receipts, and
message-level pCIDs.

- **Alt A:** The dogfood records and repo apparatus look too similar. Dave may
  accidentally treat a TODO as a message or a message as a TODO.
- **Alt B:** Dave creates or uses a simulation whose `world/` contains the
  dogfood artifacts. Message files can be strongly promise-shaped because that is
  the specimen under test. Apparatus files describe the experiment.
- **Alt C:** Dave can make messages promise-shaped, but the surrounding event and
  result artifacts may not capture enough agency/provenance for later replay.
- **Alt D:** This is viable if Dave explicitly wants the whole dogfood world,
  including observations and results, to test the strong convention.
- **Alt E:** The dogfood world risks reusing ambiguous root-level paths.

S4 favors Alt B with optional Alt D experiments.

### S5: Ellen writes the PromiseGrid Development Guide

Ellen wants to write guidance for app developers: choose a protocol, use pCID,
write handlers, publish conformance claims, and understand promise semantics.

- **Alt A:** Ellen may overteach wire-lab process artifacts as PromiseGrid API
  rules. "Every doc starts with `I promise`" becomes guide noise unless the
  protocol actually requires that.
- **Alt B:** Ellen can cite simulation results as evidence and cite frozen specs
  as normative. She can say promise-shaped artifact templates were tested in
  named simulations before being adopted, if a later DI adopts them.
- **Alt C:** Ellen has fewer experimental artifacts to draw from, but less risk
  of overclaiming.
- **Alt D:** Ellen can cite the bounded strong experiment as evidence for or
  against broad promise-shaped documentation.
- **Alt E:** Ellen must infer the relationship from scattered TEs.

S5 favors Alt B. It supports guide writing without making wire-lab artifact
style into a public PromiseGrid rule prematurely.

### S6: Mallory exploits ambiguity

Mallory publishes a file that starts with `I promise`, claims to be from Alice,
and sits near real decision records. The goal is to make readers treat it as an
authorized decision or conformance claim.

- **Alt A:** Uniform promise language increases the chance that a ritual phrase
  is mistaken for authority.
- **Alt B:** The path and simulation protocol help scope the claim. A promise
  artifact inside a simulation is evidence from that simulation, not a repo-wide
  DI. Authority still depends on signatures, review, branch policy, and DR/DI
  linkage.
- **Alt C:** Fewer promise-shaped files reduce attack surface, but commitments
  still need authentication.
- **Alt D:** Bounded experiments keep the risk contained if readers understand
  the boundary.
- **Alt E:** Ambiguity remains wherever promise language appears ad hoc.

S6 favors Alt B plus explicit warnings: promise form is not authentication, and
simulation specimens are not decisions.

### S7: Thousands of artifacts accumulate

Frank runs many simulations. Each produces events, observations, results, and
candidate messages. LLMs later search the corpus for evidence.

- **Alt A:** Uniform wrappers create large amounts of boilerplate. Search results
  may be dominated by identical promise openings.
- **Alt B:** Simulations can choose the artifact protocol that fits the question.
  Only runs testing promise-shaped records pay the full verbosity cost.
- **Alt C:** Storage and search stay simpler, but less evidence accumulates about
  promise-shaped coordination.
- **Alt D:** Strong worlds are useful as bakeoffs but expensive as the default
  simulation style.
- **Alt E:** Tooling remains inconsistent.

S7 favors Alt B with per-simulation protocols and result summaries.

## Reconsideration of TE-pahah

TE-pahah's simulation-first recommendation becomes stronger after TE-gurov.
Without `simulations/`, the promise-shaped artifact question has no safe
experimental boundary. It becomes either a global style rule or an ad hoc
template discussion.

With `simulations/`, wire-lab can run explicit artifact-protocol experiments:

```text
simulations/
  SIM-<handle>-promise-artifacts/
    README.md
    QUESTION.md
    protocol-set.md
    artifact-protocol.md
    world/
      actors/
      sites/
      groups/
      cas/
      feeds/
      wires/
    events/
    observations/
    results/
    decisions.md
```

This does not mean the repo should create that directory now. It means the
Pahah structure is expressive enough to test Gurov's proposal without making it
repo-wide first.

TE-pahah should therefore be read as more than a `transports/` versus `groups/`
answer. It is also the cleanest way to preserve a distinction between:

- decision apparatus,
- simulation specimen,
- simulation result,
- and eventual PromiseGrid source of truth.

## Reconsideration of TE-gurov

TE-gurov correctly rejects strict whole-body promises for TEs, DRs, and TODOs.
TE-pahah explains where the rejected strong idea can still be tested.

The combined reading is:

- **Strict global promise corpus:** still rejected.
- **Promise wrapper for every new durable apparatus artifact:** no longer the
  strongest default, because `simulations/` gives a safer test boundary.
- **Explicit promise blocks in apparatus:** recommended default now.
- **Promise-shaped simulation artifacts:** recommended experiment path.
- **Commitment artifacts:** strong candidates for promise wrappers after their
  templates are tested or otherwise locked.

This narrows TE-gurov's DF-gurov.1. The practical choice is not simply "Alt-B or
Alt-C everywhere." The practical choice is "Alt-C for apparatus, Alt-B or Alt-D
inside simulations that explicitly test promise-shaped artifacts, and
commitment-specific templates when a later DI adopts them."

## Recommended conclusion

Adopt the synthesis represented by Alt B:

1. Treat `simulations/` as the boundary where wire-lab can test
   promise-shaped artifact protocols.
2. Keep repo apparatus documents role-specific and ordinary unless they contain
   an explicit promise block.
3. Do not require every new TE, DR, TODO, README, or guide-resource note to open
   with `I promise`.
4. For simulation artifacts that do use promise form, separate at least:
   `Promiser`, `Promisee`, `Recorder`, `Drafter`, `Model`, `Reviewer`, `Scope`,
   `Evidence`, and `Source protocol/pCID` where applicable.
5. Graduate any successful template through the normal TE -> DF -> DI path
   before applying it outside its simulation.

This answer is conservative about the corpus and ambitious about experiments.
It preserves Pahah's momentum while giving Gurov a concrete place to become
evidence instead of philosophy.

## DF questions exposed

### DF-vilot.1 -- Should `simulations/` be the test boundary for promise-shaped artifacts?

Recommended answer: yes.

Surviving alternatives:

- **1.A -- Yes, simulation boundary (recommended).** Promise-shaped artifact
  conventions should be tested inside named simulations before repo-wide
  adoption.
- **1.B -- No, global corpus rule.** Apply promise wrappers to all new durable
  artifacts immediately after a DI.
- **1.C -- No, commitment-only rule.** Apply promise form only to commitments;
  do not run broader artifact simulations.

### DF-vilot.2 -- What is the default for apparatus documents?

Recommended answer: explicit promise blocks only.

Surviving alternatives:

- **2.A -- Ordinary apparatus plus explicit promise blocks (recommended).** TEs,
  DRs, TODOs, and guide-resource notes keep their current structures unless a
  genuine promise appears.
- **2.B -- Short role wrapper.** New apparatus docs begin with an
  artifact-specific `I promise` wrapper while preserving required sections.
- **2.C -- Status quo.** No promise-specific apparatus convention.

### DF-vilot.3 -- What identity fields should promise-shaped simulation artifacts use?

Recommended answer: accountable promiser plus provenance fields.

Surviving alternatives:

- **3.A -- Separate accountable agency and provenance (recommended).** Record
  promiser/promisee separately from recorder, drafter, model, and reviewer.
- **3.B -- Agent identity may promise when delegated.** A configured agent can be
  promiser only within authority explicitly delegated by DI.
- **3.C -- Model name as promiser.** The model name is the `I`; not recommended
  without a stronger theory of model agency.

### DF-vilot.4 -- Should current TE/DR/TODO templates change now?

Recommended answer: no.

Surviving alternatives:

- **4.A -- No template change now (recommended).** First use a simulation or
  commitment-specific template experiment.
- **4.B -- Add optional promise-block section.** Add a non-required section to
  templates where genuine promises can be recorded.
- **4.C -- Add mandatory promise wrapper.** Change templates after a DI locks
  that every new apparatus artifact needs a role promise.

### DF-vilot.5 -- What should the dev guide say before DF is resolved?

Recommended answer: provisional evidence only.

Surviving alternatives:

- **5.A -- Treat as open wire-lab evidence (recommended).** Guide writers may
  cite TE-gurov and TE-vilot as evidence of an unsettled artifact-design
  question, not as PromiseGrid API prose.
- **5.B -- Teach commitment artifacts only.** The guide may describe specs,
  adoption claims, and conformance records as promises, but not process docs.
- **5.C -- Teach all-docs-as-promises.** Not recommended until a DI locks it and
  at least one simulation has shown the convention works.

## Implications for open TODOs, DRs, and DIs

- **TE-pahah:** strengthened. Its `simulations/` tree is the right place to test
  promise-shaped artifact protocols.
- **TE-gurov:** narrowed. Its global Alt-B should not be the default until a
  simulation or commitment-specific experiment produces evidence.
- **DR-nugog:** should stay blocked on DF-pahah/DF-vilot rather than deciding a
  root-level `transports/` or `groups/` layout in isolation.
- **TODO-kugod:** remains the recovery owner for the apparatus/specimen
  distinction; the residual checklist should treat promise-shaped artifacts as
  simulation/template work, not as a reason to rewrite every recovered note.
- **TODO-rozas / DEV-GUIDE-RESOURCES.md:** guide writers should treat this as
  open evidence. The guide should not teach "all docs are promises" as settled
  until the DF/DI path closes.
- **Future simulation TODO:** if DF-vilot.1 is accepted, file a TODO to create a
  first artifact-protocol simulation rather than changing all templates directly.
- **Future DIs:** any DI that adopts promise-shaped artifacts outside a
  simulation should state which artifact classes are apparatus, specimen,
  commitment, or guide-facing prose.

## Decision status

`needs DF`. This TE recommends simulation-scoped promise-artifact experiments
and ordinary apparatus documents with explicit promise blocks, but it does not
lock the policy. The next decisions are DF-vilot.1 and DF-vilot.2.
