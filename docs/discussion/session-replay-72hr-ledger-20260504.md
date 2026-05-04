# Last 72 hours topic ledger (2026-05-01 → 2026-05-04 16:00 UTC)

Source: `past_session_contexts/sessions/2026-05-04_2026-05-10/ea135ce8/conversation.md` turns 87–192. Cross-referenced against `wire-lab` git log (60 commits in window) and `promisebase` git log (1 commit in window). Status legend: **WRITTEN** = committed to a repo; **PARKED** = on a twig, not merged; **DROPPED** = abandoned in conversation; **REVERTED** = drafted then rolled back.

---

## May 1 (turns 87–116)

| Topic | Status | Where it lives |
|---|---|---|
| TE-26 channel/transport types, pCID-keyed paths, DAG message graphs | **WRITTEN** | wire-lab `9d42b0d` → merged `5f4b1cc` (ppx/main). Later **rewritten in place** by TE-27. |
| Two-spec split for codex↔perplexity (thin wire-lab spec + protocol-specific spec) | **WRITTEN** as principle in TE-26 | TE-26 §implications. The codex-perplexity-channel-draft.md filename was discussed but renamed before any spec file landed. |
| `channels/` → `transports/` rename | **WRITTEN** | TE-27 (`a20bbc4` / merge `5115f12`). Migration done in same commit. |
| TE-27 axes-of-transport-differentiation analysis | **WRITTEN** | TE-27 same commit. DF-27.1..27.6 all locked at "use your recommendations." |
| Rewrite TE-26 in place to drop `channel` vocabulary | **WRITTEN** | Same TE-27 commit (`a20bbc4`). |
| `unicast-transport-draft.md` and `group-transport-draft.md` carve-outs | **WRITTEN** as TE-27 anticipated specs; later folded into the protocols-as-simrepos migration (TODO 014) so the actual spec files land under `protocols/group-session.d/specs/group-session-draft.md`. |
| S7 transport migration (closing one transport, opening a new one) | **DROPPED** as a TE-27 lock, **WRITTEN** as a deferred open question (Alt-T6.A). |
| `ack` in body vs headers | **WRITTEN** | TE-27 follow-on TODO. Body. |
| protocols-as-simrepos shape (`protocols/<slug>.d/`) | **WRITTEN** | TE-29 (`d4c82fc`) + TODO 014 migration (`d957ef9`). |
| Protocol-as-simrepo: each `<slug>.d/` is a simulated git repo with its own `specs/`, `implementations/`, `DR/`, `TODO/`, `DI/`, `proposals/` | **WRITTEN** | TE-29 + TODO 014. |
| DR/TODO/DI absorption into harness vs into per-protocol simrepos | **WRITTEN** | TE-29 settled it: per-protocol simrepos own their DRs/TODOs/DIs; harness keeps cross-cutting ones. |
| Five-level transports path `transports/<wire>/<L4-binding-pCID>/<session-pCID>/<message-pCID>/<msg-id>.msg` | **WRITTEN** | TE-29 §architecture. **Locked.** |
| L4 binding-layer concept (UDP example walked through) | **WRITTEN** | TE-29 §binding-layer. |
| TE-28 100-year goal (capture-resistance, generational durability) | **WRITTEN** | `aed5075` / merge `ae8a4b8`. |
| OQ-29.1 group-session freeze ceremony | **WRITTEN** as TE-29 OQ; later **resolved** in TE-31 (spec-doc inversion). |
| OQ-29.9 simulator empirics (ns-3, mininet, tc netem) | **WRITTEN** as OQ then **revised** with sandbox data (`be6dcbf`). |
| TODOs 015/016/017 (RETIRED/BLOCKED/FOLDED stubs) | **WRITTEN** | `65455c5`. |
| TODO 018 UDP-binding v0 reference impl + TODO 019 ns-3 harness scaffold | **WRITTEN** as TODO stubs only | `18cc0f9`. The actual code is **not yet written.** |

---

## May 2 (turns 117–148)

| Topic | Status | Where it lives |
|---|---|---|
| TE-30 TODO numbering rules and per-protocol TODO shape | **WRITTEN** | `4c8462b` / merge `9d037fc`. |
| TE-31 spec-doc inversion (spec-doc upstream, simrepo as implementation, closes OQ-29.1) | **WRITTEN** | `c95c112` / merge `a8f85e1`. |
| TE-32 spec-side vs implementation-side split, `implementations/` top-level | **WRITTEN** | `1522e5e` / merge `f23b89c`. |
| TE-33 spec-doc Informative References to its workshop, RFC-shaped | **WRITTEN** | `414d706` / merge `e5f3c6a`. |
| TODO 014 protocols-as-simrepos migration | **WRITTEN** | `d957ef9` / merge `ecc75fc`. |
| TE-34 TE editing policy and the TE corpus as one document with facets (Cat-1/Cat-2/Cat-3/Cat-5) | **WRITTEN** | `ca92000` / merge `d841526`. |
| TODO 020 — lock TE-34 DIs (1.C/2.A/3.C); defer AGENTS rollout behind 020.9 tabletop TE | **WRITTEN** | `ff87f31` / merge `910b58e`. |
| TE-35 tabletop simulation of the TE editing policy (Alice/Bob/Carol/Dave/Ellen/Frank/Mallory) | **WRITTEN** | `8f7143e` / merge `6be8ea0`. |
| Cat-1a (current-pointer paths) vs Cat-1b (historical-quotation paths) split | **WRITTEN** as new DI | `cd82c19` / merge `d8c3e93`. |
| DF-35.2 Cat-2 enumeration of unchanged DIs (mandatory) | **WRITTEN** | `8f8cbba` / merge `394c187`. |
| DF-35.3 Cat-2 cross-TE quotation grep (mandatory) | **WRITTEN** | `04126ac` / merge `795a846`. |
| DF-35.4 top-of-file Status field on every TE | **WRITTEN** | `a6295fc` / merge `3b65766`; followed by 020.10 retrofit across all 35 TEs (`9eede76` / merge `c2ccbe1`). |
| 020.6 Cat-1a path-reference sweep across 27 TEs | **WRITTEN** | `bc99a23` / merge `a77f678`. |
| 020.7 TE-1 Refinements section (Cat-3/Cat-4 navigational) | **WRITTEN** | `12c32f9` / merge `41c86f4`. |
| 020.5 codify TE editing policy in AGENTS.md / AGENTS-codex.md / AGENTS-ppx.md | **WRITTEN** | `416541a` / merge `4e3c8e6`. |
| 020.8 reconfirm `docs/thought-experiments/README.md` editing policy | **WRITTEN** | `52b8d19` / merge `a5e431d`. |

---

## May 3 morning (turns 149–164)

| Topic | Status | Where it lives |
|---|---|---|
| TE-36 apparatus vs specimen carve-out audit | **WRITTEN** as audit report | `4725b3e` (step 1 of 6-step plan). |
| TE-36 draft itself (DF-36.1..36.7) | **PARKED** on twig `ppx/te-20260503-022446-apparatus-vs-specimen` at `0230c20`. DF-36.5/36.6 locked; 36.7 still open. |
| TE-36 PT vocabulary recast (drop "contest"; use imposition + assessment) | **WRITTEN** to twig | `3cc1cb7`. |
| TE-36 replace 'imposition' with 'conditional-promise' (PT canon correction) | **WRITTEN** to twig | `d2278a1`. |
| TE-36 DF-36.5 Alt-5.C "Both" (apparatus summary + specimen detail) | **WRITTEN** to twig | `cbf7f41`. |
| TE-36 OQ-36.6 (promise-stack may be redundant with grid-pcid-payload) | **WRITTEN** to twig | `1a75b5b`. |
| TE-36 DF-36.6 Alt-6.D lazy carve-out (no trust-ledger.d yet) | **WRITTEN** to twig | `0230c20`. |
| Carve-out plan steps 2–6 (per-DF moves of envelope/library-API/ledger material out of harness-spec) | **NOT STARTED** — gated on TE-36 DF-36.7 lock. |

---

## May 3 mid-day pivot (turns 164–169) — git-as-file-transport for collaboration

| Topic | Status | Where it lives |
|---|---|---|
| Initial TE-37 idea: design a new git-file-transport protocol | **REVERTED** | Twig `ppx/te-20260503-112348-git-file-transport` reset back to ppx/main HEAD `4725b3e`. **Realization:** group-session already covers this. |
| `transports/draft--wire-lab-devs/` instance bootstrap (group-session N≥3) | **WRITTEN** | `a1c85f3` / merge `a1ecc72` on ppx/main. |
| Slug choice "wire-lab-devs" (no people-names; JJ stays anonymous in committed docs) | **WRITTEN** | Same commit. |
| Filename = message CID rule (drop sequential numbering) | **WRITTEN** | `19e9b37` / merge `7a903d4`. Existing two messages renamed in place; CIDs preserved. |
| Receive-merge-push-then-optionally-post cycle (S2 + S9) | **WRITTEN** | Same commit. group-session §9 split into 9.1–9.8. |
| Per-author-branch binding (each author writes to `<author-id>/main`; ppx/main is ppx's own write branch); fetch from all branches | **WRITTEN** | `463af44` + `ba9f2a6` / merge `ea3ce84`. group-session §9 includes S9. |
| `.msg` → `.txt` leaf filename rename | **WRITTEN** | Same commit. |
| `Message-ID` header — needed? | **DROPPED** as required field | `31573dc` / merge `8d3bf04`: header removed; bootstrap messages rewritten/rehashed. **Current ppx/main HEAD.** |
| New on-wire message ratifying CID-filenames + merge cycle | **WRITTEN** | message `bafkreihnvl6o5wqczwn3apzhgqfmeytvcbovo5jlv2khqc6enbtzkceyim.txt` in `transports/draft--wire-lab-devs/`. |
| Bootstrap `README.md` inside the instance with cycle as shell commands | **WRITTEN** | Same commit. |

---

## May 3 afternoon/evening (turns 170–180) — directory structure + layered model + promisebase intro

This thread proposed many things; **none of it was committed.** All discussion-only.

| Topic | Status |
|---|---|
| Add a layer between `transports/` and instance dirs, e.g. `transports/group-session/draft--wire-lab-devs/` | **DISCUSSED, NOT WRITTEN** |
| Run group-session over rsync, NNTP, UUCP — substrate pluralism | **DISCUSSED, NOT WRITTEN** |
| Drop "binding"; use "feed" instead (Usenet lineage); spec naming `protocols/git-feed.d/` etc. | **DISCUSSED, NOT WRITTEN** |
| `INSTANCE.md` manifest per instance (sites + feeds in one file, like Usenet `newsfeeds`) | **DISCUSSED, NOT WRITTEN** |
| Vocabulary fix: "grid envelope" replaces "carrier line" everywhere | **DISCUSSED, NOT WRITTEN** — Steve's correction at turn 175 is the lock. m000 has "carrier line" in its body bytes; rewriting cascades CIDs. |
| "headers" replaces inner "envelope" usage | **DISCUSSED, NOT WRITTEN** |
| Rename `transports/` → `forums/` or `groups/` | **DISCUSSED**: Steve picked **`groups/`** at turn 176 ("i like 'groups' instead of 'forums'"). **NOT WRITTEN.** |
| Layer numbering: layers 0–3 → 4–7 (or 5–7) | **DISCUSSED, NOT WRITTEN.** Steve at turn 176: "perhaps those should be numbered 4–7" then turn 177: "layers 5–7 would be less confusing." |
| Slug/state naming convention: `draft--<slug>` → `<slug>-draft`, `<pcid>--<slug>` → `<slug>-<pcid>` (sort-friendly; aligns with `slug-CID.md`) | **DISCUSSED, NOT WRITTEN.** Steve's lock at turn 176. Affects every existing path: `transports/draft--wire-lab-devs/` should become `groups/group-session/wire-lab-devs-draft/` (combining the rename + the new layer + the new state-naming). |
| Nested grid envelopes as a use case (PromiseGrid-over-Usenet, Usenet-over-PromiseGrid) | **DISCUSSED, NOT WRITTEN** |
| First-class `sites/` tree representing each simulated participant (each site has a sparse CAS, pulls/keeps/advertises chunks) | **DISCUSSED, NOT WRITTEN** — Steve's "we're not really representing sites very well" at turn 175. |
| Decentralized CAS as a core part of promisegrid | **DISCUSSED, NOT WRITTEN** — Steve's "we're not representing the decentralized CAS that is supposed to be a core part" at turn 175. |
| Sparse-CAS assumption ("no site has all CAS objects") | **DISCUSSED, NOT WRITTEN** — Steve's "think big" at turn 178. |
| Promise economy: layer-5 promises must be economic; possible capability tokens with floating exchange rates; "everyone is their own central bank"; risk of cryptocurrency toxicity | **DISCUSSED, NOT WRITTEN** |
| Switch messages to CBOR (`{CID}.cbor` instead of `{CID}.txt`) | **DISCUSSED, NOT WRITTEN** |
| Rabin chunking for all CAS content; symlinks pointing into a CAS directory | **DISCUSSED, NOT WRITTEN** |
| BGP-replacement as a possible PromiseGrid app | **DISCUSSED, NOT WRITTEN** |
| TE-38 (originally numbered) — adopt promisebase architecture; define wire-lab's federation layer; 11 DFs sketched | **DROPPED at first sketch**, then re-sketched after Steve flagged that I was reading promisebase `.md` files instead of code. **NEVER COMMITTED.** |
| Promisebase referenced as prior art (chunking with restic, merkle trees, streaming bytes into CAS; stalled on FUSE + container hosting; didn't use CIDs but should have) | **DISCUSSED, NOT WRITTEN** |
| pgmsg tool (thin CLI wrapper over promisebase libraries) | **DISCUSSED, NOT WRITTEN** |
| Compatibility/coexistence with libp2p, IPFS, ATPROTO | **DISCUSSED, NOT WRITTEN** |
| Multi-repo simulation: each simulated site is its own repo (or wire-lab-related messages stay in wire-lab, others in other repos) | **DISCUSSED, NOT WRITTEN** |
| "Don't hardcode 'jj'" rule | **WRITTEN** (as a standing rule applied to the group slug; "wire-lab-devs" was chosen instead). |

---

## May 4 (turns 181–192)

| Topic | Status |
|---|---|
| Read promisebase `db/` directory code (instead of `.md` files) | **DONE** — read db/, kv/fs/, cmd/pb/, fuse/, server/, RFCs 1003–1007. |
| randStream Go-1.20 `rand.Seed` deprecation regression | **WRITTEN** to **promisebase** `main` | commit `d98b5d3` ("db: fix randStream to use local rand.Rand source"). All db/ tests now pass. |
| `cmd/pb` Docker SDK rename drift (6 build errors), `fuse/` 3 test failures, `server/` 1 test failure | **NOT FIXED** — Steve said "hold other promisebase work for now and go back to TE-38." |
| Holding two PATs (wire-lab + promisebase) | **DONE** — promisebase PAT saved at `/home/user/workspace/.secrets/gh-pat-promisebase` (mode 600). Wire-lab PAT still in original session prompt. |
| Examined other promisebase branches | **DONE** — only `main` exists on remote. The 2021 mob-consensus branches (`fskv`, `angela/fskv`, etc.) were merged and deleted. |
| Steve's design-skepticism rule: wire-lab is canon; promisebase is a prototype | **VERBAL RULE** — not yet committed to any spec. |
| TE-37 numbering investigation | **DONE** — confirmed TE-37 has never been committed; there were two prior false-start drafts (the May 3 git-file-transport one, reverted; and an earlier in-this-session draft, removed before commit). |
| Current twig `ppx/te-20260504-090842-wire-lab-canon-pb-modernization` | **EMPTY** — has no commits beyond ppx/main HEAD `8d3bf04`. The off-base TE-37 draft was deleted before any commit. |
| Steve's final pivot: pause promisebase; pivot to message-transport dogfooding ASAP so wire-lab devs can collaborate without one person carrying solo | **VERBAL DIRECTIVE** at turn 187 — not yet acted on. |

---

## Net effect on the repo

**ppx/main (current HEAD `8d3bf04`)** contains TE-26 through TE-35 plus the supporting DI/TODO work, plus the bootstrap `transports/draft--wire-lab-devs/` instance with three messages exchanging the per-author-branch binding + filename=CID rules.

**Twig `ppx/te-20260503-022446-apparatus-vs-specimen` at `0230c20`** holds in-progress TE-36 (DF-36.5/36.6 locked; DF-36.7 open).

**Twig `ppx/te-20260504-090842-wire-lab-canon-pb-modernization`** is currently empty.

**Promisebase `main`** has the randStream fix (`d98b5d3`) and nothing else from this session.

---

## Verbal-rules-not-yet-committed (in priority order)

These are decisions Steve has stated but which have **not yet been swept into the corpus.** A Cat-2 vocabulary sweep would be needed for the corpus-wide ones; the structural ones are bigger.

1. **`draft--<slug>` → `<slug>-draft`** and **`<pcid>--<slug>` → `<slug>-<pcid>`** naming convention (Cat-2 sweep across `transports/draft--wire-lab-devs/`, group-session spec §2, harness-spec, TE-26..TE-29).
2. **"grid envelope"** replaces **"carrier line"** everywhere (Cat-2 sweep; cascades through m000/m001/m002 CIDs as a side effect — same shape as the Message-ID removal).
3. **"headers"** replaces inner-block uses of **"envelope"**.
4. **`transports/` → `groups/`** (rename of the simrepo top level).
5. **Add a protocol-slug layer**: `groups/group-session/wire-lab-devs-draft/` (combines rename + state-naming + new layer).
6. **`feed` (or `substrate`) is the right word** for git/rsync/NNTP/UUCP delivery mechanisms; it is *not* a "binding" (binding belongs to the L4-binding layer per TE-29 — UDP/TCP/QUIC). udp-binding remains correctly named.
7. **First-class `sites/` tree** representing each simulated participant; each site has a sparse CAS.
8. **Decentralized CAS** is core to promisegrid and must be represented in wire-lab's simulation.
9. **Layer numbering 4–7 (or 5–7)** above the L4 binding layer; nothing committed yet that uses these numbers.
10. **promisegrid uses CBOR + CIDv1**; messages eventually `{CID}.cbor` not `{CID}.txt`.
11. **Wire-lab is design canon; promisebase is a prototype** — when they conflict, wire-lab wins; promisebase will be refactored/modernized to align.
12. **Pivot now to dogfooding** so the dev group is exchanging messages over the live `wire-lab-devs` instance (Steve's most recent verbal directive).

---

## Where TE-37 stands today

There are **at least three different things** that have been called "TE-37" in this session, and zero of them have been committed:

- **TE-37a (May 3 11:23 UTC)**: a planned "git-file-transport" protocol TE. Reverted at turn 165 because group-session already covers it.
- **TE-37b (May 3 18:00 UTC, conversation-only)**: planned as the layered-model + promisebase-adoption + federation-layer founding TE. Got 11 DFs sketched in conversation; never written to a file.
- **TE-37c (May 4, this session)**: an on-disk draft titled "Wire-lab as design canon for the CAS storage substrate; promisebase as the prototype." Steve flagged it as off-base; deleted before commit.

Steve's question "I thought 37 was about integration of Promisebase" matches TE-37b. The session prompt's TODO list ("Git-file-transport — unblock collaboration") matches TE-37a (already done as the `transports/draft--wire-lab-devs/` instance, no TE needed).
