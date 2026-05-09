# TE-gurov: Promise-shaped coordination artifacts

*Thought experiment: should durable wire-lab coordination artifacts be written as PT-style promises, and how should author or LLM-model identity appear if they are?*

## TE ID

TE-gurov

## Status

needs DF

First drafted 2026-05-09 19:19:33 UTC.

## Decision under test

Should wire-lab require every durable artifact -- documents, Thought
Experiments (TEs), Decision Requests (DRs), Decision Intents (DIs), TODOs,
specs, review notes, and guide material -- to be expressed as a
Promise-Theory-style promise?

The strong form of the proposal has three linked parts:

- Each artifact clearly identifies the author, responsible promiser, recorder,
  LLM agent, or LLM model that produced or adopted it.
- The body text starts with the phrase `I promise`; Steve's prompt used the
  lowercase phrase `i promise`, so this TE treats casing as a live design
  question rather than a typo to silently normalize.
- The artifact's prose is promise-shaped, not merely about promises.

The decision is not whether promises are useful in the corpus. They already are:
the repo's vocabulary treats specs, claims, and adoption commitments as
promise-adjacent. The decision is whether *every artifact class* should be
forced into a promise form, or whether artifacts should instead contain explicit
promises where the speech act is honest.

## Scope

This TE covers durable coordination and documentation artifacts in wire-lab:

- Protocol specs and draft specs.
- TEs, DRs, DIs, TODOs, review artifacts, and decision logs.
- Development-guide inputs and onboarding prose.
- LLM-assisted outputs, including drafts that may be recorded by a human or
  an agent identity.

This TE does not move files, change templates, alter DR/DI schemas, rewrite
existing TEs, or create the TODO/DI that would be required to adopt any
convention. Its result is `needs DF`.

## Assumptions

- In Promise Theory, a promise is an autonomous speech act by an agent. A
  document can record, carry, or content-address a promise, but the agency
  question remains load-bearing.
- Existing wire-lab process treats DR and DI records as source-of-truth
  coordination artifacts. A DR records an unresolved question and its state. A
  DI records a locked decision and intent provenance.
- A TE is analysis that narrows alternatives before DF. A TE may recommend a
  choice, but it does not itself lock the choice.
- LLM-generated text can be useful evidence, draft material, or a record of an
  agent run. Unless explicitly delegated authority by Steve or another
  accountable operator, an LLM model is not a durable social promiser in the
  same way as a human, organization, repository role, or running service
  identity.
- Readers include human maintainers, LLM agents, external implementers, and
  future guide writers. Tooling includes simple grep, Markdown readers, schema
  checks, and content-addressed artifact stores.
- Mallory may publish a malicious artifact that uses correct promise wording but
  false identity, false authority, or misleading scope.
- The phrase `I promise` can be a natural-language sentence opening, a
  machine-readable token, or both. Those uses have different casing and
  migration costs.

## Alternatives

### Alt-A: Strict universal whole-body promise

Every artifact body starts with the exact required phrase and the whole body is
interpreted as the promise. Under the strictest version, a DR, TE, TODO, README,
guide page, and protocol spec all become direct first-person commitments by the
declared author or model.

What this makes easier:

- One uniform rule: every durable artifact is promise-shaped.
- Simple detection by reading the first words of the body.
- A strong philosophical claim that the repo is promises all the way down.

What this makes harder:

- Analysis artifacts become strained. A TE explores alternatives; it does not
  promise that any one alternative is true or adopted.
- DRs ask questions; forcing a question to be a promise can obscure the
  difference between uncertainty and commitment.
- TODOs are mutable work queues; "I promise" can imply delivery commitment even
  when the artifact only records current priority.
- LLM model identity becomes misleading if the model is presented as a promiser
  without durable agency or accountability.
- Existing DR/DI source-of-truth fields may be demoted behind prose, making
  machine parsing and audit harder.

### Alt-B: Promise wrapper plus ordinary artifact body

Each artifact has explicit identity metadata and a short opening promise that
states the artifact's function. The rest of the body remains appropriate to the
artifact type. For example, a TE might open with "I promise to analyze these
alternatives under the stated assumptions"; a DR might open with "I promise to
record this unresolved question and its current state"; a DI might open with "I
promise this record captures the decision intent authorized by the named
decision-maker."

What this makes easier:

- The artifact is promise-shaped without pretending that every sentence is a
  commitment.
- Existing structured sections can remain source-of-truth.
- LLM provenance can be separated from accountable authorship.
- Readers can distinguish the promise wrapper from contained analysis,
  questions, evidence, and logs.

What this makes harder:

- Tooling must understand at least two layers: identity/front matter and body.
- The convention is more verbose than status quo.
- Authors must learn artifact-specific promise openings.

### Alt-C: Ordinary artifacts containing explicit promise blocks

Artifacts keep their current shapes. When a real promise appears, it is placed
in an explicit block with promiser, promisee, scope, and provenance. The body
does not need to start with `I promise`.

What this makes easier:

- Existing DR/DI/TE structures remain intact.
- Only genuine promises get promise machinery.
- It is honest for analysis artifacts: they may contain promises without being
  reducible to promises.

What this makes harder:

- The corpus does not get a universal promise-shaped opening.
- Promise discovery depends on section conventions or tooling.
- The philosophical "all artifacts are promises" claim is weaker.

### Alt-D: Specs and adoption claims only

Only protocol specs, conformance statements, review replies, and adoption
claims are promise-shaped. TEs, DRs, TODOs, and guide prose remain ordinary
coordination artifacts unless they contain a specific promise.

What this makes easier:

- Promise form stays close to runtime behavior and conformance.
- The existing DR/DI source-of-truth protocol remains untouched.
- Guide writers have less vocabulary to explain.

What this makes harder:

- The boundary between spec-like and non-spec artifacts must be maintained.
- Review and governance artifacts may still need promise language, so the
  boundary will be porous.
- Some coordination artifacts that already behave like commitments may remain
  under-specified.

### Alt-E: Status quo

Do not add a promise-body requirement. Continue using existing TE/DR/DI/TODO
forms and add promise language case by case.

What this makes easier:

- No migration cost.
- No risk of false agency from LLM-generated first-person prose.
- No disruption to source-of-truth schemas.

What this makes harder:

- The corpus continues to mix promise-theoretic ideas with conventional
  documentation forms.
- Promise commitments remain discoverable only by human reading.
- The dev guide may need to explain why specs are promise-like but nearby
  process artifacts are not.

## Scenario analysis

### S1: Normal TE authoring

Alice writes a TE comparing two transport designs. Bob reads it before a DF
round. Carol later uses the TE to understand why one option was rejected.

- **Alt-A:** Alice's TE starts with `I promise`, but it is unclear what she
  promises. If she promises the final answer, the TE is dishonest because DF has
  not happened. If she promises only to analyze, then the whole-body rule is not
  actually whole-body; it has become a wrapper.
- **Alt-B:** Alice promises to analyze the alternatives under named assumptions.
  Bob can trust the scope of the work without mistaking it for a decision.
  Carol can still read ordinary scenario prose.
- **Alt-C:** Alice keeps the TE format and may include a promise block such as
  "Alice promises that the rejected alternatives were evaluated under the same
  assumptions." This is honest but less uniform.
- **Alt-D:** The TE is ordinary analysis. Promise form waits until a spec or
  adoption claim is written.
- **Alt-E:** Nothing changes.

S1 favors Alt-B or Alt-C. Alt-A collapses unless its "whole-body" promise is
weakened into a wrapper.

### S2: DR with an unresolved question

Bob opens a DR asking whether a protocol should require deterministic
serialization. The correct state is "open"; no answer has been chosen.

- **Alt-A:** "I promise deterministic serialization is required" would be false.
  "I promise to ask whether deterministic serialization is required" is true,
  but then the promise is about record-keeping, not the whole body.
- **Alt-B:** The DR can promise to record the question, why it blocks progress,
  who is waiting, and how the eventual decision will be linked. The DR remains a
  question artifact.
- **Alt-C:** The DR can remain ordinary and later contain the promise made by
  the decision-maker when the question is resolved.
- **Alt-D:** The DR stays outside promise form.
- **Alt-E:** Existing DR practice continues.

S2 rejects strict whole-body interpretation for DRs. A DR can honestly contain
or wrap a promise to maintain state; it cannot honestly be the decision it does
not yet contain.

### S3: DI source-of-truth and decision authority

Carol records a DI after Steve decides a naming rule. Dave, an LLM agent, drafts
the patch. Ellen audits the file later.

- **Alt-A:** If Dave's generated body starts "I promise", Ellen may infer that
  Dave or the model is the decision-maker. That conflicts with the rule that DI
  authorship belongs to the decision-maker unless authority is explicitly
  delegated.
- **Alt-B:** The artifact can separate fields: `Decision-maker`, `Recorder`,
  `Drafter`, `Model`, and `Promiser`. The promise says the DI record captures
  the decision intent authorized by the named decision-maker. Dave's model can
  be provenance, not agency.
- **Alt-C:** The DI can preserve its existing required fields and add an
  explicit promise block only if the decision-maker wants one.
- **Alt-D:** DI remains source-of-truth but not promise-shaped.
- **Alt-E:** Existing practice remains least ambiguous for audits.

S3 favors Alt-B when promise form is desired. It also exposes a hard rule:
model identity and promiser identity must not be conflated.

### S4: LLM-generated draft with partial accountability

Ellen asks an LLM to draft a guide page. The model name is available, but the
model cannot receive future complaints, repair misunderstandings, or accept
social consequences. Frank reviews and commits the page.

- **Alt-A:** "I promise" by the model is agency theater unless the repo defines
  the deployed agent as an accountable promiser. If Frank is the promiser, the
  model should be listed as a drafter or tool, not as the "I".
- **Alt-B:** Front matter can state: responsible author/promiser = Frank;
  drafter = LLM agent session; model = named model if known; promise = Frank
  promises the page is a reviewed guide artifact. The body can still credit the
  LLM.
- **Alt-C:** A provenance block can identify the LLM without making it the
  promiser.
- **Alt-D:** No promise form is needed unless the guide page asserts a contract.
- **Alt-E:** Existing authorship conventions still need improvement for LLM
  provenance.

S4 rejects "LLM model as default promiser." The safe framing is accountable
author plus LLM provenance.

### S5: Mixed migration and incomplete writes

Alice's branch uses promise-shaped TEs. Bob's branch still uses legacy TEs.
Carol's parser expects `## TE ID` and `## Status`. A write is interrupted after
the opening `I promise` line but before identity metadata.

- **Alt-A:** The first line is present, so a shallow checker may accept a file
  whose promiser and scope are missing. Legacy files fail the new rule even
  though they may be valid historical records.
- **Alt-B:** A valid file requires both identity metadata and the artifact-type
  sections. The opening promise is not enough. Legacy files can be grandfathered
  while new files use the wrapper.
- **Alt-C:** Existing parsers continue to work. Promise-block validation is
  local to blocks.
- **Alt-D:** Only spec/adoption tooling changes.
- **Alt-E:** No migration risk.

S5 favors Alt-B or Alt-C if the repo migrates. It also argues for validators
that check complete structure, not just the prefix phrase.

### S6: Trust boundary and Mallory

Mallory publishes a malicious artifact that starts "I promise", claims to be
from Alice, and includes plausible-looking LLM provenance.

- **Alt-A:** Uniform promise wording may create false confidence. The phrase
  does not authenticate identity or authority.
- **Alt-B:** The wrapper can require identity fields, but those fields still
  need signatures, repository provenance, branch policy, or review history to be
  trustworthy. The wrapper improves audit shape but does not solve trust.
- **Alt-C:** Promise blocks have the same authentication need.
- **Alt-D:** Spec/adoption promises still require verification.
- **Alt-E:** Status quo also requires verification, but has less risk of
  ritualized promise language being mistaken for proof.

S6 says promise form is not a security boundary. Any adopted convention must say
that "I promise" is a claim to be assessed, not proof.

### S7: Readability, guide cost, and scale

Dave writes the PromiseGrid development guide. Ellen onboards as a new
contributor. Frank operates tooling over thousands of artifacts.

- **Alt-A:** Every page begins similarly. That is memorable, but it becomes
  monotonous and can make analytical prose feel cultish or less technical.
  Guide writers must explain why questions, TODOs, and failed alternatives are
  all "promises."
- **Alt-B:** The guide can teach a consistent pattern: every artifact promises
  its function, while sections preserve normal meanings. The cost is real but
  teachable.
- **Alt-C:** The guide teaches promises as explicit blocks. Lower conceptual
  cost, weaker unifying story.
- **Alt-D:** The guide focuses promise language on specs and adoption.
- **Alt-E:** The guide spends less time on process philosophy but may leave the
  repo's promise vocabulary feeling inconsistent.

S7 rejects Alt-A for broad readability. Alt-B is viable only if the opening
promise is short, artifact-specific, and visibly separate from evidence and
analysis.

## Casing analysis: `i promise` vs `I promise`

Steve's prompt used lowercase `i promise`. There are three plausible readings:

1. **Natural-language reading:** The intended phrase is `I promise`, with the
   normal uppercase first-person pronoun. This is easiest for humans and least
   surprising in Markdown prose.
2. **Protocol-token reading:** The exact lowercase bytes `i promise` are a
   token. This is easy for machines but awkward in human-facing English and
   visually signals either informality or error.
3. **Case-insensitive transition reading:** Validators accept either spelling,
   but canonical writers emit `I promise`.

The recommended survivor is canonical `I promise` for visible prose, with
case-insensitive detection only as a migration or linting aid if a tooling
decision later requires it. If the phrase becomes a machine token rather than a
sentence, the token should probably move to front matter instead of forcing
lowercase prose.

## Compatibility with the DR/DI source-of-truth protocol

Strict universal promise form conflicts with the existing DR/DI protocol if it
turns questions into commitments or model drafts into decision authority.
Compatibility is possible if promise form is layered:

- A **TE** may promise fair analysis under stated assumptions; it must not
  promise that its recommendation is adopted.
- A **DR** may promise to record an unresolved question, blockers, waiting
  parties, and eventual linkage; it must not promise an answer before decision.
- A **DI** may promise that the record captures an authorized decision intent;
  it must keep decision-maker, recorder, and LLM drafter distinct.
- A **TODO** may promise to represent current work state and dependency order;
  it must not imply unconditional delivery unless a human explicitly makes that
  commitment.
- A **spec** may be promise-shaped at the protocol/adoption layer, but open
  questions and assumptions remain part of the honest scope of the promise.

Under Alt-B or Alt-C, DR and DI fields remain source-of-truth. Promise text
summarizes the artifact's role; it does not replace required fields, event
logs, decision authority, or append-only provenance.

## Rejected alternatives

- **Reject Alt-A as a strict whole-body rule.** It is not honest for DRs, TEs,
  and TODOs unless it is weakened into an artifact-function wrapper.
- **Reject lowercase-only `i promise` as canonical prose.** It may be acceptable
  as a temporary parser match, but human-visible canonical prose should use
  `I promise`.
- **Reject LLM model as default promiser.** A model can be named as drafter,
  model, tool, or provenance source. It should not be the accountable "I"
  unless a separate DI defines an agent identity that can bear promises.

## Surviving alternatives

- **Alt-B survives as the strongest candidate:** identity/front matter plus a
  short artifact-specific `I promise` wrapper, with existing body structure
  preserved.
- **Alt-C survives as the conservative candidate:** ordinary artifacts with
  explicit promise blocks only where a genuine promise exists.
- **Alt-D survives as the narrow candidate:** promise form only for specs,
  adoption claims, review replies, and other artifacts that already behave like
  commitments.

## Recommended conclusion

Adopt Alt-B only if Steve wants the corpus to make a visible philosophical
commitment that every durable artifact has a promise-shaped function. Otherwise
adopt Alt-C as the safer engineering default.

If Alt-B is chosen, the promise should be a wrapper around the artifact's role,
not a claim that every sentence is itself a promise. The canonical phrase should
be `I promise`, not lowercase `i promise`, and identity fields should separate
responsible promiser, decision-maker, recorder, drafter, and LLM model.

## DF questions

### DF-gurov.1: Artifact promise scope

Which scope should wire-lab adopt for promise-shaped artifacts?

- **B: Promise wrapper for new durable artifacts (recommended if adopting the
  philosophy).** New artifacts identify the responsible promiser/provenance and
  open with an artifact-specific `I promise` wrapper while preserving normal
  TE/DR/DI/TODO sections.
- **C: Promise blocks only (recommended conservative option).** Existing
  artifact formats remain ordinary; explicit promise blocks are added only when
  a real promise is being made.
- **D: Specs/adoption only.** Promise form applies to specs, conformance,
  adoption, review replies, and similar commitment artifacts, not to all docs.
- **E: Status quo.** No new promise-form convention.

### DF-gurov.2: Casing

If a visible phrase is required, what spelling is canonical?

- **`I promise` canonical, case-insensitive lint during migration
  (recommended).** Human prose stays grammatical; tooling may detect lowercase
  as a normalization issue.
- **Exact lowercase `i promise`.** Treat the phrase as a machine token even in
  Markdown body prose.
- **No body phrase.** Put any machine token in front matter and leave prose
  natural.

### DF-gurov.3: LLM identity and agency

How should LLM-generated artifacts identify agency?

- **Accountable promiser plus LLM provenance (recommended).** The human,
  organization, repo role, or delegated agent is the promiser; the LLM model is
  recorded as drafter/tool/provenance.
- **LLM agent as promiser when named.** A configured agent identity may promise,
  but only within authority explicitly delegated by a DI.
- **Model as promiser.** The model name is the "I" in `I promise`; this is not
  recommended without a stronger theory of model agency.

### DF-gurov.4: DR/DI compatibility rule

If promise form is adopted, what protects source-of-truth semantics?

- **Promise text never replaces required DR/DI fields (recommended).** DR/DI
  schemas, event logs, decision authority, and append-only provenance remain
  normative.
- **Promise text becomes normative summary.** Required fields remain, but the
  opening promise controls interpretation if there is conflict.
- **Promise text is decorative.** It is allowed but has no source-of-truth
  effect.

### DF-gurov.5: Migration posture

How should existing artifacts be treated?

- **New artifacts only until a later TE/DI.** No retroactive sweep without a
  separate migration decision.
- **Opportunistic touch.** Add wrappers when files are otherwise edited.
- **Corpus-wide migration.** Rewrite existing artifacts after a dedicated
  source-of-truth and TE-editing-policy review.

## Implications for open TODOs, DRs, and DIs

- **TODO-kulih / TE-nibar:** This TE generalizes "spec doc as promise" from
  specs to all coordination artifacts. It should not supersede TE-nibar; it
  should feed TE-nibar's DF by distinguishing doc-level promises from
  artifact-function wrappers.
- **TODO-kavug / DR-005:** Review replies are strong candidates for promise
  form because they are direct social commitments. The wrapper model may give
  them clearer promiser/promisee/scope fields.
- **TODO-fofas:** Durable review feedback as a contest artifact may benefit
  from contained promise blocks, especially when reviewer authority differs from
  recorder or LLM drafter.
- **TODO-misul and TODO-diliz:** DI provenance and DR backfills should remain
  source-of-truth work. A promise wrapper must not replace DI IDs or DR records.
- **TODO-fonuz / DI-nisam:** Any promise-shaped artifact template should reuse
  the proquint handle namespace rather than introducing a second identity
  scheme.
- **TODO-milum / DR-034:** Agent-instruction consolidation should decide where
  LLM model provenance belongs if Alt-B or Alt-C is adopted.
- **TODO-rozas / DR-tuhaz / DR-napum:** Dev-guide resources and layperson-guide
  claims need the readability result from this TE before teaching "all docs are
  promises" as a settled rule.
- **Pending future DIs:** No DI should state a universal promise-body
  requirement until DF-gurov.1 through DF-gurov.5 are answered and recorded.

## Decision status

needs DF
