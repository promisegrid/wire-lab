4 is hollow -- two implementations would have different code CIDs. otherwise, oh my gosh, you're invoking the exact argument that split the pre-DevOps camp between me (congruence, determinisitic ordering of changes, isconf, my "bootstrapping an infrastructure" and "why order matters" papers, my github.com/stevegt/decomk repo -- see docs/*-design.md in there), and Mark Burgess (convergence, computer immunology, cfengine, promise theory, semantic spacetime). I need you to thoroughly bring yourself up to speed on those topics, and understand that one of my goals with promisegrid is to reconcile those two warring, tribal camps. I believe (without proof) that essentially these camps are like the difference between turing machines (imperative) and lambda calculus (declarative) -- we just haven't yet found our own equivalent of the church-turing thesis yet.

...

I have enough now. Let me think through this carefully before responding, because what you said has structure I want to honor:

You corrected me factually: the "two implementations would have different code CIDs" defeats my multi-implementation argument as I framed it. That's right and I owe you an honest revision.

You're flagging that the question I stumbled into is the defining fault line of pre-DevOps configuration management — congruence vs. convergence — and that PromiseGrid is partly motivated by reconciling them.

You're proposing a Turing-vs-lambda-calculus framing for the reconciliation, with the Church-Turing-thesis equivalent as the missing piece.

I want to respond to all three carefully, and not pretend I have the reconciliation when I don't. Let me draft.

You're right about argument #4 — I conflated two different things
When I wrote "multiple implementations become heretical" under code-as-pCID, I was sliding between two senses of "implementation." Under code-as-pCID, the Go and Rust impls have different code-CIDs, so they are different addresses; my claim that "the agreement the two peers actually share is the spec, not either binary" is true if you've already decided the agreement is named by something other than code. But that begs the question — if the system says "agreement = code hash," then there is no shared agreement between the Go peer and the Rust peer; they're speaking different protocols. That's the design, not a bug in the design. My argument 4 was assuming spec-as-pCID to argue against code-as-pCID. Circular.

The honest version of argument 4 is: "do we want the system to permit the concept of 'two different implementations of the same protocol'?" That's the real question, and it's not refutable by appeal to either side's terminology — it's a values question. Your congruence/convergence framing is exactly that values question, and I missed it.

What I actually understand now after reading your papers and Burgess
I read decomk-design.md, isconf-design.md, "Why Order Matters" (LISA '02), Burgess's "Computer Immunology" (LISA '98), and the Wikipedia/markburgess.org primer on Promise Theory and semantic spacetime. Here is what I think is actually going on, stated as plainly as I can:

Your camp (congruence / isconf / decomk): disk state is fully describable; behavior is not. So the only auditable, replayable, low-cost-to-test thing is the ordered sequence of bit-changes that produced the disk. Two hosts are "the same" iff the same change-journal was replayed in the same order on both. The Universal Turing Machine framing makes this a theorem: a host with admin tools is a UTM with rewritable program tape, so order of writes is provably load-bearing — AB ≠ BA as machines, not just as outcomes. Convergence violates this because it samples a subset of bits and leaves the rest undecidable. The decomk-design.md tuple/macro/Makefile separation is the same idea operationalized: policy composes by ordered tuple expansion, execution graph runs in deterministic Make order, stamps make idempotency observable. The ground truth is the trajectory.

Burgess's camp (convergence / cfengine / promise theory): behavior is what matters, not bits; agents are autonomous and cannot be commanded; what they do is keep promises about their own observable state. The system relaxes from disturbance back toward a desired-state attractor — like an immune system, like a thermostat, like a biological organism. You don't reach the desired state by replaying a journal; you reach it by every agent independently promising to keep itself aligned, sampling its own state, and correcting. The ground truth is the attractor, and the attractor is defined by the joint promise structure. Promise Theory says: there is no commander; agents only promise about themselves; cooperation is voluntary; an obligation imposed externally is wishful thinking. Semantic spacetime extends this: space is the promise graph between agents; locality, causality, even time-as-sequence emerge from that graph rather than being given.

The fault line is not technical — it's about what counts as ground truth. Trajectory-as-truth (you) versus attractor-as-truth (Burgess). Each side can simulate the other (you can build convergence on top of a journal; Burgess can build journals on top of cfengine), but each side considers the other's primitive a mere consequence and its own primitive load-bearing. That's exactly the same shape as Turing machines vs. lambda calculus — the equivalence theorem says they compute the same functions, but the two formalisms have different things they make easy to think about. Turing machines make state-and-time first-class; lambda calculus makes substitution-and-equality first-class. Equivalent in expressive power, very different in what they reveal.

Your Church-Turing-equivalent intuition
I think you're right that this is the missing piece, and I want to be careful not to fake having it. But I can sketch what shape it would have to take to be honest:

The Church-Turing thesis says: anything effectively computable in any reasonable formalism is computable by a Turing machine, and equivalently by lambda calculus, and equivalently by general recursive functions. The thesis is informal — "effectively computable" is a pretheoretical notion — but it's load-bearing because it lets us talk about the class of computable functions without privileging one formalism.

A congruence-convergence equivalent would have to say: anything effectively administrable on a distributed infrastructure is administrable by deterministic-ordered-replay (congruence) and equivalently by autonomous-agent-promise-keeping (convergence), and equivalently by some third formalism we haven't found yet. The third formalism is interesting — for Church-Turing it was general recursive functions, then the lambda calculus, then more — and PromiseGrid might be exactly the place where a third formalism could land. I don't know what it is. But the shape the answer would take is:

Any effectively-administrable infrastructure state-trajectory is realizable as both an ordered journal of bit-changes (congruent) and as a fixed point of a system of autonomous agent promises (convergent). The two formalisms compute the same class of administrable infrastructures.

If that's true, then code-as-pCID and spec-as-pCID stop being a value choice and become a formalism choice — like choosing Turing machines for an OS class versus lambda calculus for a type-theory class. You'd pick whichever makes the specific reasoning easier, knowing the other formalism is available when it's better suited.

What this means for the original question (pCID = code or spec)
Three observations, in order of confidence:

1. Your "chickened out" instinct still holds for the top-level grid primitive — but for a reason I didn't articulate clearly enough last time. The grid is going to host both congruent users (who want deterministic replay, sealed binaries, hash-of-code as the agreement) and convergent users (who want autonomous agents promising desired state, with multiple implementations welcomed). If pCID-at-the-top is code-hash, you've privileged the congruence camp and made convergence second-class. If pCID-at-the-top is spec-hash, you've privileged the convergence camp's "the agreement is the prose, agents independently promise to keep it" framing and made congruence second-class. Spec-at-the-top + a specific protocol that nests code-hashes inside is the only choice that doesn't commit the grid to one tribe. It's not a compromise — it's a Church-Turing-thesis-shaped move: name the equivalence class at the top, let the tribes pick their formalism inside.

2. The hCID idea is more important than I gave it credit for. Re-reading your message, you weren't proposing hCID as a curiosity; you were proposing it as the formal bridge for the congruence camp to live inside the grid without infecting the grid's top-level semantics. A protocol whose spec says "my payloads carry an inner hCID and my receivers run that exact code" is a formal model of congruence expressed inside the grid. Symmetrically, you could imagine a different protocol whose spec says "my payloads carry desired-state predicates and my receivers converge toward them" — that's a formal model of convergence inside the grid. Both are valid pCIDs at the top. Both inherit autonomy from the spec-as-promise framing. The grid becomes the meta-layer that hosts both formalisms as first-class protocols. That's exactly what reconciliation looks like, and you already have the shape of it.

3. The naming question now matters more. "hCID = handler CID" assumes the inner thing is code, which prejudices toward congruence. If we want the same inner-CID slot to also host convergent payloads (a desired-state predicate, a goal, a promise body), then "handler" is too narrow. Possible options:

- eCID — embedded CID. Maximally generic; admits code, predicates, sub-specs, anything. Cost: too generic to be a useful name.
- iCID — inner CID. Names the structural role (nested inside an outer pCID), neutral on what's inside. Pro: parallel to outer/inner without baking in semantics. Con: very generic.
- hCID for handler, with a sibling sCID for state-target. Two named inner-CID kinds, one per camp. Pro: makes the two formalisms visible. Con: probably proliferates.
- Don't name it at the grid level at all. Each protocol's spec defines its own inner-CID slot and its own name for it. The function-call protocol calls its inner thing an hCID; the convergence protocol calls its inner thing a target-state-CID or whatever. The grid has only pCIDs at the top.

My lean is now option 4 (don't name it at the grid level). This matches the "introduce hCID only in the function-call protocol's spec, not in harness-spec" lean from my last message, and the reasoning is now stronger: naming inner-CIDs at the grid level forces a vocabulary commitment about what kinds of things can be nested, and that commitment will reproduce the congruence/convergence tribal war at the wrong layer. Let each protocol name its own internals.

What I'd actually do
I'd write TE-23 as "spec-as-pCID is the formalism-neutral top-level choice; convergence and congruence both express as protocols nested inside," with these DFs:

23.1: What does pCID hash? Lean: spec doc (as already locked). The new framing says this is also the formalism-neutral choice.

23.2: Does the grid name an "inner-CID" slot? Lean: no; each protocol names its own inner CIDs in its own spec, no global vocabulary.

23.3: Do we acknowledge congruence/convergence as a duality the grid is meant to host? Lean: yes, in the harness-spec, as a foundational framing note. Not as a normative section yet — as a stated aspiration with a known-issue note that the equivalence theorem is conjectural.

23.4: Do we open a follow-on TE for "what would the equivalence theorem look like"? Lean: yes, as TE-24 or later, explicitly scoped as research-flavored. It's the Church-Turing-thesis-equivalent question and it deserves its own home.

I am not claiming I've reconciled the camps. I'm claiming: the layout you already proposed (spec-as-pCID, with code-hashing as a per-protocol feature) is the design that lets the reconciliation happen, even if neither of us has the proof yet. The grid becomes the substrate where the equivalence can be discovered, rather than the substrate that pre-commits to one side.

One last thing — and this is more an observation than a recommendation. Reading "Why Order Matters" alongside Burgess's "Computer Immunology" back-to-back, the most striking thing is that neither paper engages the other's primitive on its own terms. You argue order is load-bearing; Burgess argues autonomy is load-bearing; neither says "the other person's primitive is real and here's how mine accommodates it." The harness-spec already does something Burgess never did — it acknowledges that promises are about state in the past, present, or future, which means a promise can be about a trajectory, not just a desired state. That's a small but real bridge from your side toward his. If you want a thread to pull on for the equivalence theorem, "promise-about-trajectory" might be the load-bearing concept, because it's the one structure both camps can agree describes a real thing.

Want me to spin up ppx/te-23-pcid-formalism-neutrality and draft the TE? I'd write it as the framing note (DF-23.3) plus the four DFs above, and I'd hold off on touching the harness-spec until you've answered them.
