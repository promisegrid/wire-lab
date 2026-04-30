# Congruence, Convergence, and the Grid

*An essay on why PromiseGrid's choice of spec-as-pCID is formalism-neutral, why congruence and convergence are dual rather than opposed, and on a thread to pull on for the missing equivalence theorem.*

— Drafted 2026-04-29 in conversation between Steve Traugott and stevegt-via-perplexity, ahead of TE-23. This is a framing essay, not a normative spec. Nothing here is yet locked as a Decision Intent.

---

## 1. Why this essay exists

A small architectural question — *should the protocol CID be the hash of the spec document or the hash of the code that implements it?* — turned out, on the third bounce, to be the question that splits the pre-DevOps configuration management world in two. On one side, the congruence camp: Steve Traugott, isconf, decomk, "Bootstrapping an Infrastructure" (LISA '98), "Why Order Matters" (LISA '02). On the other, the convergence camp: Mark Burgess, cfengine, "Computer Immunology" (LISA '98), the Promise Theory program, semantic spacetime. Both camps published their foundational papers within months of each other in the same venue. Neither has ever fully accepted the other's framing.

PromiseGrid's stated ambition, articulated in conversation with Steve, is to host both camps in a single substrate without privileging either. This essay is the part of that ambition we can write down now: a faithful summary of each camp's worldview on its own terms, an argument for why the grid's top-level addressing primitive (pCID = content hash of a spec document) is formalism-neutral rather than partisan, a speculative sketch of what a Church-Turing-equivalent reconciliation would have to look like, and a small but real bridge concept — *promise-about-trajectory* — that might be where the equivalence theorem could anchor.

The essay is written for two audiences who do not usually share a vocabulary: the trajectory-believer and the attractor-believer. It tries to use words both camps can hear without flinching. Where it has to pick a vocabulary, it says so.

---

## 2. The two camps, on their own terms

### 2.1 The congruence camp

The starting observation, due to Traugott and Brown, is that *disk state is fully describable; behavior is not*. Two hosts that pass the same behavioral tests today may diverge tomorrow because of bits we did not look at. The only auditable, replayable, low-cost-to-validate object is therefore the disk itself, and more precisely the *ordered sequence of bit-changes* that produced the disk. Two hosts are "the same" if and only if the same change-journal was replayed in the same order on both. Anything weaker is a hope.

This is not a stylistic preference. The "Why Order Matters" paper grounds it in a Turing-equivalence claim: a self-administered host running an automated systems administration tool (an ASAT, in the paper's vocabulary) is equivalent to a Universal Turing Machine with a rewritable program tape. The host's tools fetch and execute new programs that modify the host's own runtime environment, including the very tools doing the fetching. The order in which instructions are loaded changes the machine itself: a Turing virtual machine with instruction set AB is *a different machine* than one with instruction set BA, not just a machine that happened to take a different trajectory. Order is not a constraint on outcomes; order is a property of identity.

From this, the cost argument follows mechanically. If a host's behavior is determined by the ordered sequence of bit-changes applied to it, then the cheapest way to ensure two hosts behave identically is to apply identical ordered sequences. Validating, inspecting, testing, and deploying a single ordered sequence (Ctest in the paper's notation) is the lowest cost. Adequately testing partially-ordered sequences (Cpartial) requires more work, because the orthogonality of partial orders is itself undecidable in general. Predicting orthogonality (Cpredict) is more expensive still. And testing arbitrary unordered sequences (Crandom) is dominated by N-factorial combinatorial explosion, which is intractable beyond a handful of changes. The paper concludes — and this is the load-bearing claim — that *the least-cost method of administering more than one host is to apply the same ordered changes to each*. Every weaker discipline is more expensive.

Three concrete artifacts realize this view. **Isconf** treats `hosts.conf` as an inventory layer that composes ordered tuples by macro expansion (DEFAULT first, then host-specific overrides), then hands the resulting tuple stream to GNU make as an explicit ordered goal list. The make execution graph is the change-journal; stamps in `$stampdir` are the idempotency record. **Decomk** generalizes this into a three-concern model: identity selection (which keys apply), tuple resolution (expand keys to NAME=value pairs), action selection (positional args become make targets). The canonical environment contract is that one tuple stream feeds both the env-export file and the make argv, in deterministic order: incoming DECOMK_* passthrough, then resolved config tuples, then decomk-computed tuples, last-wins. Make is again the execution graph; idempotency is again observable through stamps. **The wire-lab harness-spec** carries the same DNA — DR records ordered, DI records ordered, settled statements traceable to the changes that settled them — even though it operates on prose rather than disk bits.

The unifying thesis: *the ground truth of an infrastructure is its trajectory*. A correct infrastructure is one whose trajectory has been observed-and-tested-once and replayed-faithfully many times. Convergence approaches this ideal asymptotically and never reaches it, because convergence operates on *subsets* of disk state, leaving the unmanaged complement undecidable. The paper states, flatly, "We know of no previously divergent infrastructure that, through convergence alone, has reached congruence." That sentence is the hinge.

### 2.2 The convergence camp

The starting observation, due to Burgess, is the opposite: *behavior is what matters; bits are the substrate behavior emerges from*. A computer system in production is a complex, stochastic, partially-observable artifact embedded in an adversarial environment. It is more like an organism than a program. The right metaphor is biology, not engineering: cells die and are replaced; immune systems detect and neutralize foreign material; homeostasis is maintained by feedback loops, not by faithful execution of a master plan. Burgess's "Computer Immunology" paper makes this explicit. The computer's environment is hostile; perfect prevention is impossible; the system must be designed to *recover* from departures from the desired state, not to prevent them.

Cfengine operationalizes this with **converging semantics**: the administrator describes what a system should look like (file present, permissions correct, daemon running, configuration line set), and cfengine, running periodically on each host, samples the relevant aspects of state and applies whatever corrective action is needed to move the system toward the description. When the system matches the description, cfengine becomes inert. There is no master journal; there is no commander; there is only a host that knows what it is supposed to look like and corrects itself toward that picture, repeatedly, indefinitely.

Two intellectual moves follow from this design and become foundational in their own right.

**Promise Theory**, developed by Burgess starting around 2004 and refined with Bergstra, formalizes why convergence requires autonomy. Burgess found that obligation-based logics — "host X *must* be in state S" — were "wishful thinking," because no agent in a real distributed system can in fact compel the behavior of another. Networks partition; daemons crash; humans patch things by hand. The only commitments that are not wishful thinking are commitments an agent makes about *its own* behavior. Hence: agents are autonomous (causally independent); agents can only promise about themselves; cooperation is voluntary; an obligation imposed externally is not real. Command-and-control is reproduced *inside* this framework as the special case where one agent voluntarily promises to follow another's instructions — and the autonomy framing is preserved because that promise can be withdrawn. Promises are mathematically more primitive than graph adjacencies, because a graph edge requires the mutual consent of both endpoints to be a real edge.

**Semantic spacetime**, Burgess's later work, extends this into a unifying picture: space *is* the structure of promises between agents. Locality emerges from which agents promise to interact with which others. Causality emerges from the temporal order of promises being made and kept. Time-as-sequence is one observer's interpretation of a promise graph, not a given backdrop. Functional spaces (a parking lot, a database, a hotel) are spaces with semantics — and the semantics is *exactly the joint promise structure of the agents in that space*. This is not metaphor. Burgess's program is to show that the same formalism works at every scale, from a cfengine cluster to a human organization to (speculatively) physical spacetime itself.

The unifying thesis: *the ground truth of an infrastructure is its attractor*. A correct infrastructure is one whose joint promise structure has the desired state as a stable fixed point, such that arbitrary perturbations decay back toward it. The cost of validation is not the cost of testing every ordered sequence; it is the cost of describing the attractor and verifying that each agent's promises individually contribute to it. Multi-agent systems become tractable not because we control them but because we accept that we don't and design for the relaxation dynamics anyway.

### 2.3 What each camp considers load-bearing, and what each considers a bug in the other

A careful reading of both camps reveals something subtle that took this conversation to surface: *autonomy is load-bearing in both worldviews*. The split is not autonomy-vs-no-autonomy. The split is over what an autonomous agent's promise is *about*.

In the congruence camp, autonomy is a *desired and implemented feature*, not a bug. This is not always obvious from a quick read of the cfengine-vs-isconf debate, but it is the load-bearing point of the push-vs-pull argument that runs through "Bootstrapping an Infrastructure," through isconf, and through decomk. *Hosts pull.* Pushing commands at a host is an imposition and is fragile in exactly the way Promise Theory says it is — networks partition, daemons crash, trust is unilateral, and "host X *must* run command Y" is wishful thinking. So the congruence camp's actual design is: a host *promises* to fetch the journal, *promises* to execute it in order, and *promises* to record the result. A host might not keep those promises — but in that case, the host (and its human sysadmin by extension) is breaking an earlier promise to be a cooperative member of the infrastructure. The journal's authority is not a command issued at the host; it is a contract the host has voluntarily entered. Replay is real because the host is autonomous *and* keeps its promises, not in spite of autonomy.

This is structurally important. It means the disagreement between the camps is not "autonomy yes vs autonomy no." Both camps treat autonomy as primitive. The disagreement is over what an autonomous agent finds it natural to promise *about*: an ordered trajectory of changes (congruence) versus a fixed-point state attractor (convergence). Both kinds of promise are coherent Promise-Theory promises; they differ in what they assert.

With autonomy off the table as a point of disagreement, the actual genuine disagreements come into sharper focus.

For congruence, **convergence-without-trajectory is a bug**. If each host is autonomously promising only about its current-and-future state, then no two hosts have necessarily traversed the same trajectory, and the trajectory was the thing whose identity congruence cared about. The infrastructure ceases to have a *history* you can replay; it has only a *current condition* you can sample. Even with autonomy preserved, you have lost the auditable change-journal that made multi-host correctness cheap.

For convergence, **mandatory ordering is a bug**. If every host is committed to replaying a master journal in a specific order, the design over-specifies: it requires the journal author to anticipate every contingency, embeds a single point of failure in the journal-authoring process, and cannot scale to environments where partial communication and unanticipated drift are the norm. Even with autonomy preserved (hosts *do* promise to follow the journal), the design buys determinism at the cost of resilience.

Each camp can simulate the other. You can build convergence on top of a journal: just make the journal contain "converge to state S" as one of its steps. Traugott actually did exactly this, at NASA Ames Research Center — an ordered series of mixed cfengine and shell script runs that ran every hour in order to constrain changes on supercomputers, making reasoning about them easier. You can build congruence on top of cfengine: just extend the cfengine language syntax to allow the sysadmin to specify Makefile-like prereq chains for any subset of changes. Traugott's proposal for this, and the ensuing backlash, was the origin event for his "Why Order Matters" paper. But each camp considers the *other's* primitive a mere consequence of its own, while considering its own primitive load-bearing. This is the same shape as the Turing-machines-vs-lambda-calculus argument in computability theory: each formalism can encode the other, but each formalism makes different things easy to *think about*. Turing machines make state-and-time first-class; lambda calculus makes substitution-and-equality first-class. The Church-Turing thesis says they compute the same class of functions, but the two formalisms remain in active use because they reveal different structure.

Congruence and convergence might be in exactly that relationship. With autonomy clarified as common ground rather than disputed territory, the remaining disagreement — *what an autonomous agent should promise about* — looks much more like a formalism-choice question than a values-clash question. We do not yet know whether the two formalisms admit the same class of administrable infrastructures, because no one has stated and proved the equivalence theorem.

---

## 3. The pCID question, viewed through the camps

The pCID question we started from — *is the protocol CID the hash of the spec document or the hash of the code that implements it?* — looks innocent until you notice that it is a vote.

If pCID is the hash of the *code*, then the agreement between two peers is "we run byte-identical code." Two peers either match bit-for-bit or they speak different protocols. There is no daylight between protocol identity and implementation identity. Multiple implementations of "the same protocol" cannot exist as a first-class concept; what exists is multiple protocols that happen to interoperate in practice. This is the *congruence* camp's worldview made into a primitive: the trajectory (the code) *is* the agreement.

If pCID is the hash of the *spec document*, then the agreement between two peers is "we both promise to behave according to this prose." Multiple implementations are first-class — Go, Rust, hand-rolled C, all promising to keep the same protocol. Conformance is a relationship between an implementation and the spec, not a property of the implementation in isolation. This is the *convergence* camp's worldview made into a primitive: the attractor (the spec) is the agreement, and each peer autonomously promises to be in its basin.

Stated this way, the pCID question is not architectural; it is tribal. Whichever way the grid decides, one of the two communities is told their primitive is second-class.

There is a third option, and it is the one PromiseGrid has, in fact, already chosen: pCID hashes the **spec document**, and *specific protocols* — defined inside the grid, addressed by their own pCIDs at the top — are free to declare that their payloads carry inner code-hashes (call them hCIDs, fCIDs, whatever the protocol's spec wants to call them). The grid's top-level primitive is the wider, formalism-neutral one (the agreement). The narrower, formalism-specific one (the implementation) lives one level down, scoped to the protocols that need it.

This is not a compromise. It is structurally analogous to a Church-Turing-thesis-shaped move: name the equivalence class at the top; let the formalisms live inside. A grid whose top-level addressing is spec-as-pCID is a grid that *can host* a function-call protocol whose spec says "treat me congruently — my payloads are addressed to inner code-hashes and my receivers run that exact code." It is also a grid that can host a desired-state protocol whose spec says "treat me convergently — my payloads are state predicates and my receivers run whatever local logic gets them there." Both are valid pCIDs at the top. Both inherit the autonomy framing from spec-as-promise. Neither tribe is told their primitive is second-class.

The reason this works is subtle and worth stating directly: **a spec doc can natively host either worldview, because a spec doc is a thing autonomous agents promise about.** A spec that says "implementations may vary; here is the desired behavior" hosts convergence directly. A spec that says "implementations must be byte-identical to the code at hash X" hosts congruence directly. The pCID is the same kind of thing in both cases; only the contents of the spec differ. That is the move.

This also dissolves the "two implementations have different code CIDs, therefore multiple implementations can't share an agreement" objection. They share an agreement *one level up* — at the pCID of the spec — even though they have different code-CIDs *one level down*. The grid sees the pCID; the protocol-specific layer sees the inner CIDs. Naming them differently is not bookkeeping; it is the structural acknowledgement that they live at different levels and do different work.

### 3.1 Nesting

A natural follow-on question: can these CIDs nest? Yes, and the four cases are worth being explicit about.

**pCID inside pCID** is routine: one protocol whose payloads are messages of another protocol. This is the layered-protocol case, like TCP inside IP inside Ethernet. Each outer pCID names the calling convention; each inner pCID names the next layer in. There is no upper bound on nesting depth that the grid imposes; specific protocols may impose their own.

**hCID inside pCID** is the canonical congruent-protocol-inside-the-grid case: a protocol whose spec says "my payloads name an inner code-hash and my receivers run that exact code." This is the function-call protocol of your original instinct. The outer pCID is the calling convention (a spec); the inner hCID is the realization (code).

**hCID inside hCID** is rarer but coherent: code that calls code, where both are content-addressed. It is the content-addressed-RPC case where one function invokes another by hash. Not philosophically problematic; the grid does not need to forbid it.

**pCID inside hCID** is the meta-circular case and the one that takes a moment to see clearly. It says: "this exact code is the receiver, and the messages it receives name a spec to interpret them." This is the shape of a Lisp `eval`, a Smalltalk image, an interpreter generally — congruent code whose job is to look up a spec and apply it.

A practical and increasingly common instance of pCID-inside-hCID is *deterministic code calling an LLM*. Most LLM client programs fit this pattern: the client code itself is referenceable by hCID (it runs byte-identically on every host that uses it), and the system message it sends to the model is naturally a pCID — a spec the model is asked to behave according to. The user's prompt is then function input, the payload addressed to that pCID. This pattern shows up in production today wherever "a tool wraps an LLM call." The hCID names the deterministic harness; the pCID names the behavioral spec the model is being asked to honor; the prompt is parameter data. The hCID and pCID together constitute the interpreter; the convergent receiver (the model) does its work inside that frame. It is a clean example because it shows that pCID-inside-hCID is not exotic — it is the natural decomposition of any "deterministic code that orchestrates a non-deterministic responder." The grid does not need to forbid it; it shows up whenever congruent code orchestrates a convergent receiver.

The grid does not need to know about any of these patterns at the top level. Each protocol's spec is free to constrain what nests inside its own payloads. This is the same discipline that lets TCP not need to know about HTTP: the layering is a per-protocol concern, not a substrate concern.

---

## 4. The reconciliation question, honestly

I want to be careful here: I do not have the equivalence theorem. I have the *shape* the theorem would need to take, and a thread to pull on, and a reason to think the thread is real. Stating those clearly without overclaiming is the goal of this section.

### 4.1 The shape

The Church-Turing thesis says: any function that is *effectively computable* in any reasonable formalism is computable by a Turing machine, and equivalently by lambda calculus, and equivalently by general recursive functions, and equivalently by several other formalisms discovered later. The thesis is informal — "effectively computable" is a pretheoretical notion that the thesis pins down post hoc — but it is load-bearing because it lets us talk about *the class of computable functions* without privileging a formalism. Different formalisms then become tools for thinking, chosen for which structure they reveal, not for which class of functions they admit.

A congruence-convergence equivalent would have to say something like:

> *Any effectively-administrable infrastructure-state-trajectory is realizable as both an ordered journal of bit-changes (the congruent realization) and as a fixed point of a system of autonomous-agent promises with appropriate convergence dynamics (the convergent realization). The two formalisms admit the same class of administrable infrastructures.*

If something like that were true, then code-as-pCID and spec-as-pCID would stop being a values choice and become a formalism choice — like choosing Turing machines for an OS course versus lambda calculus for a type-theory course. Each formalism would be picked for the reasoning it makes easy, knowing the other is available when better suited.

I want to flag the words *effectively-administrable* and *appropriate convergence dynamics* in that statement. Both are doing the same kind of pretheoretical work that *effectively computable* did in 1936. Pinning them down rigorously is most of the work, and I do not have a candidate definition. But the *shape* is recognizable, and the shape is what tells us this might be a theorem rather than a slogan.

### 4.2 Why this is not obviously vacuous

A skeptic could say: "Of course you can encode either formalism in the other; both are Turing-complete substrates. The 'equivalence' is trivially true and tells us nothing." This deserves an answer.

The Church-Turing thesis is also "trivially true" in the same sense — Turing machines and lambda calculus are both Turing-complete, and any Turing-complete formalism can encode any other. What makes the thesis non-trivial is the claim that the *informal notion* of effective computability is captured by the formal definitions, *and* that the formalisms turn out to suggest the same boundary on what is and is not computable. The halting problem is undecidable in any of them. The same functions are recursive in all of them. The non-triviality is in the alignment between intuition and multiple independent formalisms, not in the inter-encodability.

The congruence-convergence question would be non-trivial in the same way if (a) "effectively administrable" turned out to be a meaningful pretheoretical notion (it does — sysadmins know it when they see it), (b) the two formalisms turned out to admit the same boundary on what is and is not administrable (this is the open question), and (c) some infrastructures turned out to be administrable by neither, with the same infrastructures excluded by both formalisms. (c) is the tell. If you find an infrastructure that is administrable congruently but not convergently, or vice versa, the equivalence is false. If every counterexample turns out to be administrable by both or neither, the equivalence is on its way to being a theorem.

I do not know whether (c) holds. I think it is the right experimental question.

### 4.3 The thread to pull on: promise-about-trajectory

Here is the small but real thing the wire-lab harness-spec already says that neither original camp said. From the locked vocabulary: *promises are assertions of state in the past, present, or future, often conditional*.

That phrase is doing more work than it looks like.

Burgess's original Promise Theory describes promises as commitments by an autonomous agent about its *current and future behavior*. The classical examples are "I will deliver this byte stream," "I will keep my CPU usage under 80%," "I will respond to requests on port 443." These are promises about state-as-of-now and state-going-forward. They are the natural primitive for convergence: an attractor is described as "everyone keeps their promises," and the attractor is reached when those promises are jointly kept.

Traugott's "Why Order Matters" describes a host's identity as the *ordered history* of bit-changes that produced it. Trajectories, not states. A host is what its journal says it is. And — with §2.3's correction in hand — a host is what its journal says it is *because the host has autonomously promised to fetch and apply that journal in order*. This is the natural primitive for congruence, and it is already a Promise-Theory-shaped commitment: a host promises that it pulled changes C1, C2, C3 in that order, and an administrator who relies on the journal is relying on that promise being kept.

The wire-lab vocabulary's "past, present, or future" extension makes the bridge explicit. A promise-about-the-past is, structurally, *a claim about a trajectory*. "I promise that I, peer P, applied changes C1, C2, C3 in that order" is a perfectly well-formed promise in the extended Promise-Theory sense, and it is also exactly the kind of statement a congruence-camp host's journal records. The journal entry and the promise are the same object, viewed from two sides. The journal-author side sees "the canonical record of what should have happened on this host"; the host side sees "the promise I am keeping by replaying this record." Both views are simultaneously true. Neither view alone is the whole picture.

If that is right — and I want to emphasize *if*, because I am speculating — then the equivalence theorem might land here. A congruent infrastructure is one where every host's identity is defined by the joint set of *promises-about-trajectory* it has made and is keeping about itself. A convergent infrastructure is one where every host's identity is defined by the joint set of *promises-about-state* it is keeping about itself. Both are sets of self-promises by autonomous agents in the extended sense; both are realizable in either substrate. The equivalence would be: any infrastructure expressible as a coherent set of promises-about-trajectory has a corresponding expression as a coherent set of promises-about-state-with-convergence-dynamics, and vice versa, modulo the kind of pretheoretical caveats Church-Turing always carries.

The empirical content of this claim, if it is a claim, is something like: *autonomous agents who promise about their own trajectories produce the same observable behavior, in the limit, as autonomous agents who promise about their own attractor states under appropriate dynamics*. I think that is checkable in small cases (model both kinds of agent in a simulator; observe; compare). I think it is an open question in the general case.

This is the thread. It is a real thread because the wire-lab vocabulary already extended Promise Theory in exactly the direction that makes the bridge possible. Burgess never wrote "promise about the past" because his application — convergence — did not need it. Traugott wrote about pull-not-push and host autonomy as load-bearing operational features, but he did not frame those features in Promise-Theory vocabulary, because Promise Theory had not yet been written down when "Bootstrapping an Infrastructure" appeared. The wire-lab gets to write in the vocabulary both camps now share, after the fact, and so it has — almost by accident — written down the primitive that might unify them.

I want to be very clear: writing down a primitive is not the same as proving an equivalence. The thread is gold, in the sense that it is the right thing to pull on. But pulling on it is real work and may take a long time, and may end at "the equivalence does not hold and here is the obstruction," which would itself be a result.

---

## 5. An observation about how the original camps wrote, and what the wire-lab has the chance to do differently

Reading "Why Order Matters" and "Computer Immunology" back-to-back, the most striking thing — and I mean this as an observation, not a criticism — is that *neither paper engages the other's primitive on its own terms*. Traugott and Brown argue order is load-bearing, and treat convergence as a special case that fails to deliver the determinism their framework needs. Burgess argues autonomy is load-bearing, and treats imperative ordering as a category error of the obligation-logic kind. Neither paper says "the other person's primitive is real and here's how mine accommodates it." The conversation across the wall has been, at best, polite mutual incomprehension.

This is not unusual. It is what foundational camps do during their formative period: each defends its primitive because giving ground feels like giving up the program. Turing did not write much about lambda calculus. Church did not write much about Turing machines. The equivalence between them was proved by other people, later, after both formalisms had matured enough that proving them equivalent did not threaten either.

The wire-lab is in a different position. It is being designed *after* both camps have matured. It does not have to pick a side to get started; the foundational work both sides need is already done. This is a rare and probably temporary opportunity. The grid's top-level vocabulary can be built to *host* the equivalence rather than to *embody* one side of it. That is what the spec-as-pCID choice does, and that is why it is worth being deliberate about — not because it solves the reconciliation problem, but because it is the choice that does not foreclose on it.

If we make that choice carefully, the grid becomes a place where someone could, eventually, write the equivalence theorem (or its refutation). The promise-about-trajectory thread is the most concrete reason to think the theorem is even reachable. The harness-spec's "past, present, or future" wording is the reason it is reachable from inside the grid's existing vocabulary, without retrofitting.

---

## 6. What I am claiming, and what I am not

What I am claiming, in order of confidence, highest first:

1. **The pCID-hashes-the-spec choice is formalism-neutral.** It does not privilege congruence or convergence. Both worldviews can be expressed as protocols nested inside the grid, addressed by their own pCIDs. This is the right choice for the grid's top-level addressing and was correct on its merits before this conversation; the conversation only made the reasoning explicit.

2. **Autonomy is common ground.** Both camps treat autonomy as load-bearing. The congruence camp's pull-not-push discipline is a Promise-Theory-shaped design that predates Promise Theory's vocabulary: hosts autonomously promise to fetch and replay the journal. The disagreement between the camps is not autonomy-vs-no-autonomy; it is over what an autonomous agent's promise is *about* (an ordered trajectory vs a fixed-point state).

3. **The two camps are dual rather than opposed.** They make different things easy to think about; they each consider the other's primitive a mere consequence of their own; they have not engaged each other's framing on its own terms in twenty-five years; the structural shape of the disagreement matches the Turing-vs-lambda-calculus shape in computability theory.

4. **The promise-about-trajectory thread is a candidate bridge.** The wire-lab harness-spec's locked phrasing — *promises are assertions of state in the past, present, or future* — extends Promise Theory in exactly the direction that lets a journal entry and a host's pull-and-replay promise be two views of the same object. This is a real bridge concept, not an analogy.

5. **A Church-Turing-equivalent might be reachable from there.** "Effectively administrable" is a pretheoretical notion in the same sense that "effectively computable" was in 1936. The two formalisms (ordered journals replayed under host promise, autonomous-agent promise systems with convergence dynamics) might admit the same class of administrable infrastructures. The empirical question — whether some infrastructure is administrable by exactly one of the two — is the experiment.

What I am *not* claiming:

- I am not claiming the equivalence theorem holds. I do not have a proof, and I do not have a refutation. I have a shape and a thread.
- I am not claiming the wire-lab needs to commit to the reconciliation as a deliverable. The grid does not need the theorem to be useful; it needs only the choice that does not foreclose on the theorem. We are making that choice.
- I am not claiming that hCID, fCID, or any inner-CID name should be promoted to the harness-spec's top-level vocabulary. Per Steve's instruction in this conversation, those names live inside the specific protocols that use them. The grid's top-level vocabulary stays small.
- I am not claiming Burgess and Traugott are wrong about anything. Both camps' foundational papers are correct on their own terms. The reconciliation, if it exists, is *additive*: it does not invalidate either side's results; it reveals a wider frame both sides occupy.

---

## 7. A note on tone, for both audiences

To the trajectory-believer: I have not asked you to give up determinism. The grid hosts your worldview directly. A congruent protocol whose spec demands inner hCIDs is a perfectly normal protocol on this grid, and your replay-from-a-journal discipline is preserved inside it. The only thing the grid asks is that the protocol *itself* be named by the hash of its spec — which is just an ordered prose document, exactly the kind of artifact your discipline already produces and tracks.

To the attractor-believer: I have not asked you to give up autonomy. Every pCID names an agreement that autonomous agents promise to keep, on their own terms, with no commander. A convergent protocol whose spec describes desired state and leaves implementation to the implementor is the natural shape of pCID itself; nothing in the grid forces you to nest code-hashes inside if your protocol does not need them. Promise Theory is not just compatible with this design; it is the design.

To both: the disagreement between your camps is real and deserves more careful work than either side has yet put into engaging the other's primitive on its own terms. This essay is a small step toward that work, written in the hope that the grid we build together can be the substrate where the next generation of that work happens.

---

## 8. Appendix: where this essay's claims came from

This is a working essay drafted in conversation, not a polished paper, so I want to be transparent about what informs each section.

- §2.1 (congruence camp) draws on `doc/decomk-design.md` and `doc/isconf-design.md` from `github.com/stevegt/decomk` (read in full), and on Traugott & Brown, "Why Order Matters: Turing Equivalence in Automated Systems Administration," LISA '02 (read in full via [USENIX](https://www.usenix.org/legacy/events/lisa2002/tech/full_papers/traugott/traugott.pdf)). The Turing-equivalence quote and the cost-ordering inequality come directly from that paper.

- §2.2 (convergence camp) draws on Burgess, "Computer Immunology," LISA '98 (read via [USENIX](https://www.usenix.org/event/lisa98/full_papers/burgess/burgess.pdf)), the [Wikipedia entry on Promise Theory](https://en.wikipedia.org/wiki/Promise_theory) (which is high-quality and reflects Burgess and Bergstra's own framing), and Burgess's [Semantic Spacetime project introduction](https://markburgess.org/blog_spacetime3.html). The "wishful thinking" characterization of obligation-based logics is Burgess's own phrasing.

- §2.3 (load-bearing vs bug-for-the-other-camp) was substantially revised after Steve corrected an earlier draft: autonomy is *not* a bug for congruence — it is a desired and implemented feature, manifest in the pull-not-push design that runs through "Bootstrapping an Infrastructure," isconf, and decomk. The corrected framing locates the disagreement at "what does an autonomous agent promise about?" rather than "is autonomy admitted at all?" The historical examples of each camp simulating the other — Traugott's hourly mixed-cfengine-and-shell ordered runs at NASA Ames Research Center, and his proposal to extend cfengine's syntax with Makefile-like prereq chains (whose backlash was the origin event for the "Why Order Matters" paper) — are Steve's first-person contribution from this conversation.

- §3 and §4 are original synthesis from this conversation. The pCID-as-formalism-neutral framing emerged from Steve's instinct that hCID/fCID belongs as a per-protocol concept rather than a top-level grid concept; the Church-Turing-shaped framing was Steve's explicit suggestion in his message of 2026-04-29 21:19 PDT. The pCID-inside-hCID example using a deterministic LLM client (hCID = client code, pCID = system message, prompt = payload) is Steve's contribution from the same conversation.

- §4.3 (promise-about-trajectory) is original to this essay. The observation that the wire-lab's "past, present, or future" phrasing is the bridge primitive is mine; the phrasing itself is Steve's, locked earlier in the harness-spec's vocabulary. The strengthened version of the bridge — that a host's pull-and-replay promise *is* the journal entry, viewed from the other side — emerged from Steve's §2.3 correction.

- §5 (how the original camps wrote) is observational, based on having read both foundational papers in the same session. The Turing/Church historical analogy is a standard one in computability theory pedagogy.

- §6 (claims and non-claims) is meant to be re-readable as a standalone summary.

This essay is a draft and will probably change. It is being saved as `docs/essays/congruence-convergence-and-the-grid.md` rather than as a TE because it is framing rather than a decision request. If we want to lock any of its claims as DIs, those will be in TE-23 and possibly TE-24, with this essay as a referenced source.
