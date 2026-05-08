# TE-nijab: Transport layering and freeze boundaries

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-nijab

## Status

decided

## Decision under test

DUT-nijab: Resolve the category error exposed by DR-suhod B/C/D: whether `transports/` directories are protocol specs that freeze and rewrite when a higher-layer pCID is minted, or whether they are append-only simulation data for lower-layer network/feed behavior that can carry many higher-layer protocols over time.

The decision matters because the current corpus contains two incompatible shapes:

- Older outer transport text says `transports/<pcid>--<slug>/` names a single transport-protocol pCID and that the pCID defines the directory interior.
- TE-vipir and TE-sihih move toward layered paths in which `transports/` models the real network/feed substrate and higher-layer protocols ride above it, much as TCP can carry HTTP, SSH, SMTP, or any other application protocol.

This TE does not lock the answer. It tests the alternatives and frames the DF questions that should be answered before DR-suhod B/C/D, TODO-turog, TODO-pipus, or TODO-duvuk mutate transport data.

## Assumptions

- `transports/` is committed repository data. Existing message files under `transports/wire-lab-devs-draft/` are historical specimens and must not be edited merely because a spec later freezes.
- A spec doc can freeze by minting a pCID and writing a frozen snapshot. A historical transport message cannot be "refrozen" by rewriting its bytes without changing its message CID.
- TE-sihih's L5/L6/L7 model is the current best prior art: feeds move chunks over substrates, CAS stores chunks, and group/message protocols define meaning above those lower layers.
- TE-numan's migration invariants remain in force: migration must be reconstructible, must not silently rewrite pre-migration message bytes, and must preserve enough old-frontier evidence for future readers.
- DR-suhod intake A is already settled by DI-012-20260508-033513. This TE concerns only the remaining B/C/D cluster: manifest scope, transport directory naming/layering, and stale spec references.

## Scope and systems affected

In scope:

- The meaning of `transports/` as a repo tree.
- The freeze checklist in TODO-turog and the group-session freeze gate.
- DR-suhod intakes B/C/D.
- Stale path and lifecycle claims in `transport-spec-draft.md`, `group-session-draft.md`, `transports/README.md`, and `transports/wire-lab-devs-draft/README.md`.

Out of scope:

- Actually migrating `transports/wire-lab-devs-draft/` to the TE-sihih CAS-pointer shape.
- Changing existing `.txt` transport messages.
- Changing `tools/spec` behavior before the DF questions below are answered.
- Designing the full feed/CAS migration workflow owned by TODO-pipus, TODO-duvuk, and later CAS/feed TEs.

## Threat / trust model

- Cooperative actors are honest but may be running mixed versions of the docs. A safe rule must let them distinguish "old specimen data" from "current target architecture" without rewriting either.
- Mallory can exploit ambiguous freeze language to smuggle in a history rewrite: if a mechanical rewrite is described as "freezing," she can replace old message bytes while claiming to preserve the same conversation.
- Future readers may have only git history and on-disk files. They must be able to reconstruct whether a file was originally authored under draft semantics or generated later as a derived/migrated artifact.
- The filesystem and git object store are assumed reliable. This TE is not about corruption recovery; it is about category boundaries and auditability.

## Alternatives

### Alt A — One higher-layer protocol per transport directory

Keep the older `transports/<pcid>--<slug>/` rule. The pCID in the directory name is the protocol that owns the directory interior. Freezing a draft higher-layer spec can require renaming the directory and rewriting message envelopes to use the minted pCID.

This is close to the current `transport-spec-draft.md` and early group-session language. It is simple when a directory carries exactly one protocol forever. It becomes a layer inversion when the same lower-layer transport should carry multiple higher-layer protocols.

### Alt B — Layered transport path, future data

Treat `transports/` as the lower-layer network/feed simulation surface. Future transport data uses a layered path whose components identify the wire/feed/group/message stack, following TE-vipir/TE-sihih rather than the older single-pCID directory rule.

The exact final path may still be refined by a later DF, but the important property is fixed for this alternative: a lower-layer transport/feed is not owned by one higher-layer group/message protocol. A TCP-like feed can carry multiple L7 protocols.

### Alt C — Preserve current bootstrap data as a historical specimen

Keep `transports/wire-lab-devs-draft/` readable and append-only as pre-layered bootstrap data. Do not rewrite existing `.txt` messages. When a layered target exists, add a migration/supersession artifact and create new data under the new layered shape.

This is not a full future architecture; it is a compatibility bridge that lets the corpus keep moving without falsifying existing CIDs or message bodies.

### Alt D — Rewrite-at-freeze

At freeze time, rename `transports/wire-lab-devs-draft/` and rewrite every message's `grid draft:group-session` line to a minted pCID. Recompute CIDs and update filenames in a single mechanical commit.

This gives a tidy final tree but violates the expectation that transport data is append-only evidence. It also makes old message CIDs stop naming the bytes that were actually exchanged.

## Scenario analysis

### S1 — Normal operation: one carrier, multiple higher-layer protocols

Alice and Bob have a TCP feed. At 10:00 they exchange a group-session message. At 10:01 the same feed carries an SSH-like maintenance protocol. At 10:02 it carries a ppx-dr proposal message. The lower-layer transport is the same, while the higher-layer message semantics differ.

Alt A makes this awkward because one directory pCID appears to own the whole interior. Either all higher-layer protocols are forced behind one aggregate spec, or each protocol gets a separate directory that no longer represents the same lower-layer feed.

Alt B models the case directly: the feed path identifies the lower-layer carrier, and higher-layer pCIDs identify the protocols carried over it. This matches the OSI intuition that TCP can carry HTTP, SSH, SMTP, and many other protocols without becoming those protocols.

Alt C keeps today's bootstrap data usable while acknowledging it predates the layered shape. It does not solve the future case alone, but it avoids corrupting the evidence while Alt B is built.

Alt D provides no benefit here; rewriting old files to a single group-session pCID makes the multi-protocol future harder to explain.

### S2 — Spec freeze

The group-session draft freezes and receives a pCID. The transport/feed stack that carried the draft messages does not change: the bytes already exchanged remain the bytes already exchanged.

Alt A encourages a misleading action: rename the data directory and rewrite envelopes so old data looks as if it had originally used the frozen pCID. That conflates spec publication with historical data mutation.

Alt B treats freeze as a spec-side event. New messages can cite the frozen pCID through the layered stack; old messages remain draft-era artifacts.

Alt C adds a bridge: a migration or supersession record can say "these draft-era messages led to this frozen spec and this future path," without changing the messages.

Alt D fails the CID discipline. Even if the rewrite is mechanical, the rewritten messages are new messages. Treating them as the same history is false.

### S3 — Reader from the future

Ellen clones the repo in 2030. She sees both frozen specs and old transport data. She needs to reconstruct what happened without access to chat logs.

Alt A leaves Ellen with ambiguous evidence if old data was rewritten: did the messages originally use the frozen pCID, or were they rewritten after the fact?

Alt B is reconstructible if the docs say `transports/` is data and the layered path encodes stack choices.

Alt C is the strongest bridge for Ellen. A recorded migration/supersession artifact can preserve old directory, old frontier CID, new destination, and reason for the change.

Alt D is actively hostile to Ellen because it destroys the distinction between "authored under draft semantics" and "rewritten after freeze."

### S4 — Mixed-version participants

Alice has pulled the layered docs. Bob still has the pre-layered `wire-lab-devs-draft` instructions. Both act honestly.

Alt A creates a race: if freeze/rename/rewrite is treated as the canonical operation, Bob may author against old paths while Alice rewrites them.

Alt B lets participants choose stack paths explicitly. A mixed-version participant can keep writing to the old specimen until a migration record tells them it is closed.

Alt C gives the minimum safety rule: old data remains present, and new data starts elsewhere only after an explicit additive record.

Alt D makes mixed-version operation dangerous because it changes the ground under Bob's local history.

### S5 — CAS pointer migration

TE-sihih says future messages under `transports/` become CBOR pointers into `cas/`. Today's `wire-lab-devs-draft` messages are inline text files.

Alt A and Alt D tempt an all-at-once rewrite. That rewrite may be appropriate for a derived view, but not for the historical specimen itself.

Alt B makes the target architecture explicit: future messages follow the layered path and pointer shape.

Alt C keeps the migration honest: the old inline files remain as historical input; the new pointer files are new artifacts that can cite the old frontier.

## Conclusions

- Reject Alt D for transport data. Rewriting historical transport messages at freeze time violates CID-based identity and append-only evidence.
- Reject Alt A as the long-term model for `transports/`. It collapses lower-layer carrier/feed identity into higher-layer protocol identity and does not model "one transport carries many protocols."
- Keep Alt B as the surviving architecture for future data: `transports/` models the lower-layer network/feed simulation surface, and higher-layer protocol pCIDs ride above it in the path/metadata shape locked or refined by TE-vipir, TE-sihih, and follow-on DFs.
- Keep Alt C as the surviving compatibility rule for current data: `transports/wire-lab-devs-draft/` remains a historical pre-layered specimen until a future migration explicitly supersedes it.

## DF questions exposed

### DF-nijab.1 — What is `transports/`?

Recommended answer: `transports/` is lower-layer simulation data for wires/feeds/substrates. It is not a namespace of frozen higher-layer protocol specs.

Surviving alternatives:

- **1.A — Lower-layer simulation surface.** Future transport data uses a layered path/metadata shape so one carrier/feed can carry many higher-layer protocols.
- **1.B — Single-protocol directory.** Keep `transports/<pcid>--<slug>/` as one pCID owning each directory's interior.

### DF-nijab.2 — What happens to `wire-lab-devs-draft`?

Recommended answer: preserve it as historical pre-layered specimen data and migrate additively later.

Surviving alternatives:

- **2.A — Preserve and supersede.** Do not rewrite old messages; add a migration/supersession record and write future data elsewhere.
- **2.B — Derived rewrite only.** Permit a derived mirror or generated view that rewrites old messages, but keep the original specimen authoritative.

### DF-nijab.3 — How should freeze docs be corrected?

Recommended answer: specs freeze; transport data does not freeze. Replace "rename and rewrite every message" checklist language with "record migration/supersession when the new layered target exists."

Surviving alternatives:

- **3.A — Correct now after DF.** Update TODO-turog, group-session freeze gate, transport README, and DR-suhod after this TE's DFs lock.
- **3.B — Park behind TODO-pipus/TE-43.** Leave current freeze checklist marked suspect and do not change spec docs until the first CAS/feed migration is designed.

## Implications for open work

- **DR-suhod B** should no longer be framed only as "`tools/spec` is hard-coded to `wire-lab.d`." The deeper issue is that spec freeze tooling should freeze specs, not transport data.
- **DR-suhod C** should be reframed from "which frozen directory name?" to "what path/metadata shape represents a lower-layer transport carrying many higher-layer protocols?"
- **DR-suhod D** should rewrite stale references only after DF-nijab.1 through DF-nijab.3 lock.
- **TODO-turog** must not execute its current step 5 as written until this TE is resolved. The phrase "rewrite every message's grid envelope" is under direct challenge here.
- **TODO-pipus** remains the likely home for the eventual operational migration of `wire-lab-devs-draft` to the TE-sihih-aligned shape.
- **TODO-duvuk** remains the likely home for CID-cascade policy once the project distinguishes original transport specimens from generated/derived rewrites.

## Decision status

`decided` — Steve answered the DFs on 2026-05-08. DF-nijab.1 -> 1.A, recorded in DI-026-20260508-054722: `transports/` is a lower-layer network/feed simulation surface. DF-nijab.2 -> 2.A, recorded in DI-026-20260508-054723: `transports/wire-lab-devs-draft/` remains historical specimen data and future migration is additive. DF-nijab.3 -> 3.B, recorded in DI-026-20260508-054724: freeze-doc cleanup is parked behind TODO-pipus/TE-43, so this TE's 3.A recommendation is not executed in this pass.
