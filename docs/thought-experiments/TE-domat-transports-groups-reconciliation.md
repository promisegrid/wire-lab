# TE-domat: `transports/` / `groups/` reconciliation after TE-vipir, TE-sihih, and turn 176

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-domat

## Status

needs DF

## Decision under test

DUT-domat: Determine how to reconcile the conversation-level `transports/` -> `groups/` direction from turn 176 with the written TE-vipir / TE-sihih / TE-nijab corpus, the current live `transports/wire-lab-devs-draft/` specimen, and the emerging L5/L6/L7 PromiseGrid model.

The narrow question is not just "what should the directory be called?" The real question is which layer each top-level tree represents:

- `transports/` has been used for the wire-lab's observable transport/feed simulation data.
- `groups/` is the verbal L7 vocabulary Steve chose for named group instances.
- TE-sihih keeps `transports/` in the locked path scheme while renaming path-level-3 vocabulary from "forum" to "group".
- TE-nijab later decides that `transports/` should be lower-layer simulation data, not a namespace owned by one higher-layer protocol.

This TE therefore tests whether the turn-176 direction should land as a literal top-level rename, a split between lower-layer `transports/` and higher-layer `groups/`, or some other path discipline.

## Why this TE now

DR-nugog currently asks whether the outer `transports/` tree should group instances under a protocol slug, keep the flat shape, or defer. That framing is incomplete after rereading the turns 149-208 recovery material. The missing issue is that Steve had already chosen "groups" over "forums" in turn 176, while the committed TEs never cleanly wrote a `transports/` -> `groups/` decision.

The repo now contains three conflicting classes of evidence:

1. **Written TE locks that preserve `transports/`.** TE-vipir says `transports/` is the layer-decomposition surface for bytes-as-bytes simulation and gives the path `transports/<real-world-transport-name>/<L4-binding-pCID>/<session-pCID>/<message-pCID>/<message-id>.txt`. TE-sihih later says TE-vipir's path scheme is unchanged except for vocabulary (`binding` -> `feed`, `forum` -> `group`) and an L6 CAS subtree.
2. **Conversation locks that choose "groups".** Turn 175 says "`transports/` currently holds instances of message-exchange protocols" and proposes `forums/` as a better word. Turn 176 explicitly changes that to `groups/`, says `forums/` -> `groups/` in the directory tree, and maps `transports/draft--wire-lab-devs/` to `groups/group-session/wire-lab-devs-draft/`.
3. **Later recovery and layering decisions that make a direct rename suspicious.** Turn 177 corrects the layer order: feeds move chunks, CAS stores/resolves chunks, groups define message semantics. TE-nijab then decides `transports/` is lower-layer network/feed simulation data and current `transports/wire-lab-devs-draft/` remains historical pre-layered specimen data until an additive migration exists.

The TE must preserve all three kinds of evidence without silently pretending that any one of them erased the others.

## Evidence map

### Written corpus

- **TE-zalut** locked the early thin outer rule: a transport directory is keyed by pCID plus slug, messages do not declare their transport, and the handler for a pCID owns the directory interior.
- **TE-junil** renamed `channels/` to `transports/` because the directory was understood as the wire-lab's simulation surface for observable message flow. It reserved "channel" as a possible future above-transport concept but did not introduce that layer.
- **TE-vipir** moved the model toward a layered path under `transports/` and introduced protocol simulated repositories. Its path model treats the path itself as the instance/stack declaration.
- **TE-sihih** adopted L5/L6/L7 vocabulary and L6 CAS pointer indirection, but explicitly kept TE-vipir's path shape. Its vocabulary table changes "forum" to "group" and `<forum-pCID>` to `<group-pCID>`, not top-level `transports/` to `groups/`.
- **TE-nijab** decided that `transports/` is lower-layer simulation data and that current `transports/wire-lab-devs-draft/` remains historical pre-layered specimen data. It rejects rewriting historical transport messages at freeze time.
- **`transport-spec-draft.md`, `group-session-draft.md`, and `transports/README.md`** still contain older flat `transports/<pcid>--<slug>/` / `transports/wire-lab-devs-draft/` wording. Those files are stale relative to TE-vipir, TE-sihih, and TE-nijab, but they remain part of the current contradiction.

### Session-log evidence

- **Turn 170** asks whether the group-named directory should be nested under a generic git-file transport directory, especially if a second named group appears. The answer frames protocol-slug nesting under `transports/` as the recommended near-term shape.
- **Turn 172** reframes the problem around running group sessions over git, rsync, unison, UUCP, UDP, SVN, or CVS. That exposes the feed/substrate axis and makes group-session substrate-agnostic.
- **Turn 173** surveys precedent: email, Usenet, FidoNet, CORBA, SOAP, gRPC, libp2p, Matrix, and git all separate message identity/meaning from delivery substrate.
- **Turn 174** rejects "binding" vocabulary, prefers "feed", rejects a `bindings/` subdirectory inside message data, and moves toward an instance manifest or external feed declarations.
- **Turn 175** says `transports/` is misnamed for the current content, proposes `forums/`, introduces first-class sites and CAS, and sketches a layered TE with "forums sit on top of feeds sit on top of CAS".
- **Turn 176** chooses `groups` over `forums`, validates 1:1 ephemeral TCP/WebSocket as a size-2 group over an appropriate feed, allows nested grid envelopes, chooses `<slug>-draft` / `<slug>-<cid>`, raises symlink/pointer-to-CAS options, and maps `transports/draft--wire-lab-devs/` to `groups/group-session/wire-lab-devs-draft/`.
- **Turn 177** corrects the layer order: L7 groups define semantics, L6 CAS stores/resolves chunks, and L5 feeds advertise/request/replicate chunks. This makes groups oblivious to feeds and makes a single lower-layer feed usable by many higher-layer protocols.
- **Turns 178-179** add sparse CAS, promise-economy, conditional release, geofencing, multi-hop discovery, and promisebase prior-art pressure. These do not settle the directory name, but they strongly support a split between semantic group views and lower-layer chunk movement.
- **Turn 196** records `transports/` -> `groups/` plus protocol-slug layer as a verbal rule not yet committed.

### Current tracking evidence

- **TODO-lilar / dropped-thread disposition** treat `groups` vs `transports` as a turn-176 correction that was not cleanly committed.
- **TODO-jivam** keeps the turn-149-through-208 recovery open until items like this are either recovered, transferred, or explicitly deferred.
- **TODO-kugod / DR-nugog** currently own the residual transport-tree structure question, but DR-nugog's alternatives are too narrow because they ask only "nested under `transports/` or flat?" and do not account for a possible top-level `groups/` tree.

## Assumptions

1. Existing message bytes under `transports/wire-lab-devs-draft/` are historical specimen data. They may be cited, superseded, or used as migration input, but they must not be silently rewritten or moved in a way that pretends their original path never existed.
2. A pCID names a protocol/spec, not a particular human-readable directory name. Human-readable slugs are navigation aids unless a spec explicitly makes them semantic.
3. The current best layer model is the turn-177 / TE-sihih model: L7 groups, L6 CAS, L5 feeds. Groups talk in CIDs and semantics; CAS stores/resolves CIDs; feeds move chunks between sites.
4. A top-level tree should have one dominant responsibility. If a tree is both a lower-layer feed trace and a higher-layer group registry, future readers will misread it.
5. The wire-lab is an experimental simulation space. It must preserve wrong turns and pre-layered specimens as evidence when those artifacts help future developers understand why the design changed.
6. This TE explores structure only. It does not move `transports/wire-lab-devs-draft/`, create `groups/`, rewrite specs, or choose final names.

## Threat / trust model

Alice, Bob, and Carol are cooperative future contributors. They may have different versions of the repo and may lack the session logs. They need the filesystem and docs to tell them which artifacts are historical specimens, which trees represent current target architecture, and where to write new messages.

Mallory is not assumed to break git or hashes. Mallory exploits ambiguity. If the docs say "rename transport data at freeze" or "transports means groups now" without preserving history, Mallory can present a rewritten artifact as if it were the original exchange. The defense is an explicit additive migration record and clear layer ownership for every tree.

## Alternatives

### Alt A — Literal top-level rename: `transports/` -> `groups/`

Move the current and future message-instance tree from `transports/` to `groups/`. The turn-176 mapping becomes literal: `transports/wire-lab-devs-draft/` eventually becomes `groups/group-session/wire-lab-devs-draft/`.

**Easier:** Matches the verbal turn-176 correction directly. Human readers looking for group conversations see a plainly named `groups/` tree. It eliminates the misleading OSI-loaded word `transport` for group messages.

**Harder:** Conflicts with TE-junil's reason for `transports/`, TE-vipir's layered path, TE-sihih's unchanged path scheme, and TE-nijab's decision that `transports/` is lower-layer simulation data. It risks collapsing L5 feed traces and L7 group semantics into one tree under a L7 name.

**Obligations:** Requires a superseding TE/DI for any older written TEs whose settled path claims become false. Requires a migration/supersession record for current `transports/wire-lab-devs-draft/` rather than an unexplained `git mv`.

### Alt B — Keep `transports/` as the only operational tree

Do not create `groups/`. Future data follows TE-vipir / TE-sihih / TE-nijab through a lower-layer `transports/` path or metadata shape. "Group" is a vocabulary term for the L7 protocol/component inside that path, not a top-level directory.

**Easier:** Maximally preserves the written corpus. It respects TE-nijab's lower-layer simulation-surface decision and avoids a broad path churn.

**Harder:** It leaves the turn-176 "groups" direction partly unlanded and keeps humans looking under a word that Steve already flagged as misleading for named message spaces. It also forces every group-level reader to understand a lower-layer tree before finding the group.

**Obligations:** Requires explicit docs saying the turn-176 `groups/` wording was narrowed to vocabulary only. Without that, recovery readers will continue to see a missing implementation.

### Alt C — Split responsibilities: `groups/` for L7 group views, `transports/` for lower-layer feed/specimen data

Keep `transports/` as the lower-layer simulation/evidence tree described by TE-nijab. Introduce or plan `groups/` as the L7 group registry/view tree: group identity, manifests, membership/roster pointers, current frontier pointers, and references to CAS roots. Message bytes live in CAS; lower-layer movement lives in `transports/` or future feed-specific trees; `groups/` records group semantics and navigable group state.

Under this alternative, turn 176 was right that the current group-message directory was conceptually a group, but later turns refined the implementation: the correct target is not a destructive whole-tree rename. It is a layer split.

**Easier:** Reconciles all evidence. `groups/` satisfies the human/L7 vocabulary lock; `transports/` remains the lower-layer specimen/log surface from TE-junil, TE-vipir, TE-sihih, and TE-nijab. It also matches turn 175's realization that a group directory should become a manifest/roster/pointer set while messages live in CAS and sites hold sparse views.

**Harder:** Introduces two trees and therefore a discovery rule between them. A developer must know when to look at `groups/` versus `transports/` versus `cas/` versus `sites/`.

**Obligations:** Requires a spec paragraph or README diagram that defines the relationship among groups, sites, CAS, and transports/feeds. Requires a future migration plan that preserves `transports/wire-lab-devs-draft/` as historical specimen data and creates new group-view artifacts additively.

### Alt D — Rename semantic content to `groups/` and move lower-layer evidence to a new `feeds/` or `wires/` tree

Use `groups/` for group-level data and abandon `transports/` as the lower-layer tree name. Replace it with a clearer lower-layer name such as `feeds/`, `wires/`, `links/`, or `substrates/`.

**Easier:** Removes the overloaded and OSI-confusing word `transport` entirely. A future reader sees `groups/` for L7 and `feeds/` for L5.

**Harder:** It is more disruptive than Alt C because it renames both sides of the split. It also reopens TE-junil and TE-nijab vocabulary at the same time, turning one recovery problem into two.

**Obligations:** Requires its own naming TE/DF for the lower-layer tree. Requires a much larger Cat-2/Cat-3 sweep and likely supersedence notes on several written TEs.

### Alt E — Defer and keep DR-nugog as-is

Do nothing now. Leave `transports/` flat/stale in specs and keep DR-nugog's current alternatives.

**Easier:** Zero immediate file moves and no naming churn.

**Harder:** It preserves a known contradiction. Future contributors will rediscover the same problem, and the recovery work remains incomplete. It also makes DR-nugog misleading because it omits the top-level `groups/` evidence from turn 176.

**Obligations:** If chosen, the repo must explicitly say this is an indefinite deferral, not an unresolved accidental omission.

## Scenario analysis

### S1 — Normal dogfooding today: Alice posts to `wire-lab-devs`

Alice writes a group-session message and Bob reads it. The team needs a path that is easy for humans and LLMs to find while the system is still being dogfooded.

Alt A is easiest for a human who thinks "I want the group conversation." Alt B is easiest for tooling that already follows current `transports/` paths. Alt C gives both: new human-facing group state can live under `groups/`, while the existing transport specimen stays where old tooling and history expect it. Alt D is clear after migration but expensive before dogfooding. Alt E leaves the confusion active.

### S2 — A second named group appears

Alice and Bob use `wire-lab-devs`; Carol creates `wire-lab-ops`.

Alt A handles this by `groups/group-session/wire-lab-devs-draft/` and `groups/group-session/wire-lab-ops-draft/`. Alt B can represent it under a layered `transports/` path, but the group list is not obvious without understanding lower-layer components. Alt C lets `groups/group-session/` become the human registry for both groups while lower-layer observations remain under `transports/`. Alt D also handles it, after a lower-layer rename. Alt E repeats turn 170's original failure mode.

### S3 — The same group is carried over git, rsync, and UUCP

The group identity and messages are invariant; delivery substrates differ. A message may arrive through git on Alice's site and UUCP on Bob's site.

Alt A is weak if `groups/` contains the only files, because feed-specific evidence has nowhere layer-correct to live. Alt B is strong for lower-layer evidence but weak for group-level navigation. Alt C is strongest: `groups/` describes the L7 identity/frontier, CAS stores bytes, and `transports/` / feed traces describe how chunks moved. Alt D can be equally strong if its replacement lower-layer name is settled. Alt E leaves each new feed to invent placement ad hoc.

### S4 — A 1:1 ephemeral WebSocket flow

Alice and Bob have a private, possibly non-durable exchange over WebSocket. The turn-176 reasoning says this is still a size-2 group over a websocket-feed, possibly with an ephemeral CAS.

Alt A can name it as a group but risks over-promising durability if `groups/` is assumed append-only. Alt B can model the lower-layer WebSocket flow but makes it less obvious that the same L7 group semantics apply. Alt C can represent the L7 group as ephemeral or bounded-retention while `transports/` records whatever lower-layer evidence the simulation preserves. Alt D also works after the lower-layer rename. Alt E gives no guidance for ephemeral cases.

### S5 — PromiseGrid-over-Usenet and Usenet-over-PromiseGrid

PromiseGrid-over-Usenet uses Usenet/NNTP as a feed/substrate carrying PromiseGrid messages. Usenet-over-PromiseGrid uses a PromiseGrid group protocol whose payloads are Usenet articles.

Alt A blurs these unless it also has a separate lower-layer tree. Alt B cleanly models PromiseGrid-over-Usenet but gives no top-level place for the Usenet-over-PromiseGrid L7 group. Alt C makes the direction explicit: if Usenet is the carrier, look in lower-layer feed/transport evidence; if Usenet articles are payloads inside a group, look in `groups/`. Alt D can do the same if `feeds/` replaces `transports/`. Alt E preserves ambiguity.

### S6 — Sparse CAS and large chunked content

No site has every CAS object. Alice's site has chunks X and Y; Bob's site has Y and Z; Carol sees only a group frontier CID and must decide which chunks to fetch based on promises and permissions.

Alt A is insufficient by itself because named group directories are not where sparse chunk possession lives. Alt B captures lower-layer movement but not the semantic group frontier. Alt C maps cleanly: `groups/` references frontiers and group promises, `sites/` records local views, `cas/` stores chunks or pointers, and `transports/` records feed observations. Alt D can map cleanly too, but with larger vocabulary churn. Alt E blocks the sparse-CAS design from becoming readable.

### S7 — Spec freeze and CID preservation

The group-session draft freezes. Existing pre-freeze messages under `transports/wire-lab-devs-draft/` have old CIDs and old paths.

Alt A is dangerous if implemented as a direct move/rewrite because it can make old data look as if it was authored in the new tree. Alt B and Alt C both preserve old data. Alt C additionally gives a place to create new group-level manifest/frontier records without mutating old transport specimens. Alt D is safe only if it uses additive migration, not rewrite. Alt E leaves freeze tooling prone to the same rewrite confusion TE-nijab rejected.

### S8 — Future reader after context loss

Ellen clones the repo in 2030 and lacks the chat logs. She sees TE-vipir, TE-sihih, TE-nijab, `transports/wire-lab-devs-draft/`, and maybe new `groups/` files.

Alt A requires careful supersedence notes or Ellen will believe TE-vipir/sihih were simply wrong. Alt B requires a note explaining why turn-176's `groups` choice did not create a path. Alt C is most reconstructible if it ships with a migration/supersession record: old transport specimen here, new group view there, CAS and site roles explicit. Alt D is reconstructible only with a larger migration note. Alt E fails the recovery objective.

### S9 — Mallory exploits path ambiguity

Mallory wants to pass off rewritten pre-freeze messages as original messages, or wants to smuggle a higher-layer protocol through a lower-layer path whose owner is unclear.

Alt A gives Mallory room if the rename is treated as a move rather than a migration. Alt B and Alt C resist this because they preserve old specimens and use explicit layer boundaries. Alt C is stronger for human review because it separates "what the group says" from "how chunks moved." Alt D depends on disciplined migration. Alt E leaves ambiguity exploitable.

## Conclusions

1. **Reject Alt E.** Deferral without a corrected DR keeps the recovery hole open.
2. **Reject a naive Alt A direct rename.** Turn 176 did choose "groups" for the semantic concept, but later evidence makes a literal whole-tree `transports/` -> `groups/` rename too broad. It would conflict with TE-nijab's lower-layer meaning for `transports/` and with TE-sihih's preserved path scheme.
3. **Keep Alt B as a valid minimal path if the project wants no new top-level `groups/` tree yet.** It is coherent with the written corpus, but it under-serves the human/group-navigation goal Steve identified.
4. **Prefer Alt C as the surviving recommendation.** `groups/` should be the L7 group registry/view if and when that tree is introduced. `transports/` should remain lower-layer transport/feed simulation evidence unless a separate naming TE chooses a better lower-layer name. Current `transports/wire-lab-devs-draft/` remains historical pre-layered specimen data.
5. **Keep Alt D as a later vocabulary cleanup, not this decision.** Renaming the lower-layer tree from `transports/` to `feeds/` or another term may be attractive, but it is a separate larger decision.

The practical interpretation is:

> Turn 176 locked the word **group** and exposed that the current `transports/wire-lab-devs-draft/` directory was carrying L7 group semantics under the wrong-looking name. It did not, after turn 177 and TE-nijab, justify erasing `transports/` as the lower-layer specimen tree. The stable architecture is a split: group semantics in `groups/`, lower-layer feed/transport evidence in `transports/`, content in `cas/`, and local partial views under `sites/`.

## DF questions exposed

### DF-domat.1 — What should `groups/` mean?

Recommended answer: Alt C.

Surviving alternatives:

- **1.A — L7 group registry/view beside `transports/` (recommended).** `groups/` names group-level manifests, frontiers, rosters, and semantic pointers. It does not replace lower-layer transport/feed evidence.
- **1.B — No top-level `groups/` for now.** Keep "group" as vocabulary only and use `transports/` / CAS / sites until a concrete group-view artifact is needed.
- **1.C — Literal `transports/` -> `groups/` rename.** Rejected by this TE unless Steve explicitly chooses the churn and supersedence obligations.

### DF-domat.2 — What happens to `transports/wire-lab-devs-draft/`?

Recommended answer: Preserve it as historical pre-layered specimen data and migrate additively later, matching TE-nijab.

Surviving alternatives:

- **2.A — Preserve and supersede (recommended).** Do not move or rewrite existing messages. Add a migration/supersession record when new group/CAS/feed artifacts exist.
- **2.B — Derived mirror only.** Permit generated or derived views under `groups/`, but keep original `transports/wire-lab-devs-draft/` authoritative for historical bytes.

### DF-domat.3 — How should DR-nugog be reframed?

Recommended answer: Rewrite the open question as a split-tree / layer-ownership decision, not merely protocol-slug nesting under `transports/`.

Surviving alternatives:

- **3.A — Reframe DR-nugog around Alt C (recommended).** DR asks whether to introduce `groups/` as L7 group view while preserving `transports/` as lower-layer specimen data.
- **3.B — Close DR-nugog and open a new DR.** Use a fresh DR if the existing question is too stale to amend cleanly.

## Implications for open work

- **DR-nugog** should stop asking only "protocol-slug nesting under `transports/` vs flat." It should reference this TE and ask the layer-ownership question exposed above.
- **TODO-kugod / UT-159.b** remains the residual owner until DR-nugog is answered. The transport-spec companion audit cannot be finished while the top-level tree semantics remain open.
- **TODO-turog** must not perform any freeze-time move/rewrite of `transports/wire-lab-devs-draft/`; TE-nijab already rejected that for historical specimen data.
- **TODO-pipus / TODO-duvuk** remain likely owners for the additive migration and CID-cascade policy once the target group/CAS/feed shape is locked.
- **TE-vipir and TE-sihih** do not need body edits now. If DF-domat later locks Alt C, add Cat-3 refinements to point readers from their `transports/` path language to the split-tree interpretation.
- **`transport-spec-draft.md`, `group-session-draft.md`, `transports/README.md`, and `transports/wire-lab-devs-draft/README.md`** need later cleanup after DF-domat and DR-nugog settle. This TE intentionally does not edit those specs.

## Decision status

`needs DF`. This TE recommends Alt C / DF-domat.1.A, DF-domat.2.A, and DF-domat.3.A, but does not lock them. Steve still needs to choose whether `groups/` becomes a real L7 tree now, remains vocabulary-only for now, or revives the literal top-level rename despite the conflicts described above.

## Refinements

### 2026-05-09 — Simulation-first reframing

TE-pahah (`docs/thought-experiments/TE-pahah-wire-lab-simulation-first-structure.md`) reframes this TE's root-level `transports/` / `groups/` question under Steve's clarification that current wire-lab artifacts are brainstorming experiments, not production or active-use compatibility commitments. TE-pahah should be read before answering DF-domat: if concrete worlds move under `simulations/<sim>/world/`, then root-level `transports/` versus root-level `groups/` becomes a secondary question rather than the main structure decision.

### 2026-05-10 — Preserve evidence, not the root `transports/` path

TE-mupoz (`docs/thought-experiments/TE-mupoz-root-protocol-migration-scope-under-simulations.md`) further narrows this TE's preservation language. Steve clarified that `transports/` does not need to remain in its current location. The current reading is therefore: preserve the historical evidence, source path, source commit, and message CIDs, but allow the concrete `transports/wire-lab-devs-draft/` specimen to move into the first Pahah simulation with an explicit migration manifest. The root path itself is not a compatibility constraint.
