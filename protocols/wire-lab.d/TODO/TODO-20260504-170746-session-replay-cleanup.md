# TODO 021 - Session-replay cleanup: walk turns 149-192 chronologically and reconcile dropped threads

Track the work to walk session `ea135ce8` chronologically from turn 149 (2026-05-03 00:05 UTC) through turn 192 (2026-05-04 16:00 UTC), one turn at a time, capturing every dropped thread, open question, and verbal-but-uncommitted decision into the repo BEFORE moving to the next thread. The session compressed context once (around the promisebase pivot at turn 180-187), and at least one earlier rename rule (`draft--<slug>` -> `<slug>-draft`) and one vocabulary correction (`grid envelope` not `carrier line`) were locked verbally but never committed.

The starting state for this TODO is `last-72hr-topic-ledger.md` in the workspace root (not in the repo); that file is the audit input for this work.

The discipline for this TODO: for each turn, (1) read what was discussed, (2) describe it back to the user, (3) ask any outstanding clarifications, (4) write results (or a placeholder note) to the repo, (5) capture all unfinished threads into this file or a successor TODO before moving on, (6) check off the turn here.

## Subtasks

- [x] 021.0 (meta) Create this file, twig `ppx/te-20260504-170746-session-replay-cleanup`, commit the audit input ledger to the repo. **Done 2026-05-04, commit `a6215ae`, pushed.**
- [x] 021.149 Turn 149 (2026-05-03 00:05 UTC) -- "yes" -- (DF-35.1 Cat-1a/Cat-1b split confirmation; led to DI-020-20260502-232651). **Walked 2026-05-04.** Result already on disk: DI-020-20260502-232651 committed `cd82c19`, merged `d8c3e93` to ppx/main, pushed, twig deleted. Bot then presented DF-35.2 Alt-2.a (mandatory DI-ID enumeration as Cat-3 Refinement on TE-34), which is what turn 150 confirmed. No leftover threads.
- [x] 021.150 Turn 150 (2026-05-03 00:23 UTC) -- "yes" -- (DF-35.2 Cat-2 enumeration confirmation). **Walked 2026-05-04.** DF-35.2 Alt-2.a (mandatory enumeration of unchanged DIs by ID in Cat-2 notes) was confirmed and landed as a Cat-3 Refinement on TE-34 (commit `8f8cbba`, merge `394c187`). Bot then presented DF-35.3 Alt-3.a (mandatory cross-TE quotation grep before Cat-2 sweep) which is what turn 151 confirmed. No leftover threads.
- [x] 021.151 Turn 151 (2026-05-03 00:28 UTC) -- "yes" -- (DF-35.3 Cat-2 cross-TE quotation grep confirmation). **Walked 2026-05-04.** DF-35.3 Alt-3.a (mandatory cross-TE quotation grep before any Cat-2 sweep is run, with bot emitting matches for human review) was confirmed and landed as a Cat-3 Refinement on TE-34 (commit `04126ac`, merge `795a846`). Bot then presented DF-35.4 Alt-4.a (top-of-file Status: header field across the corpus) which is what turn 152 confirmed. No leftover threads.
- [x] 021.152 Turn 152 (2026-05-03 00:35 UTC) -- "yes" -- (DF-35.4 top-of-file Status field confirmation). **Walked 2026-05-04.** DF-35.4 Alt-4.a (uniform top-of-file `Status:` header field on every TE) was confirmed and landed as a Cat-3 Refinement on TE-34 (commit `a6295fc`, merge `3b65766`). The 020.10 retrofit immediately followed: all 35 existing TEs got a top-of-file `Status:` line (commit `9eede76`, merge `c2ccbe1`). All four DF-35.* are now answered; the locked TE editing policy is four DIs (DI-020-20260502-213103 / -213104 / -213105 / -232651) plus four Cat-3 Refinements on TE-34. Subtasks 020.5/020.6/020.7/020.10 transitioned from DEFERRED to ready-to-execute. No leftover threads.
- [x] 021.153 Turn 153 (2026-05-03 00:44 UTC) -- "do it" -- (sweep authorization for 020.5/020.6/020.7/020.10). **Walked 2026-05-04.** Bot's turn summary explicitly named only 020.5 (AGENTS rollout, commit `416541a` / merge `4e3c8e6`) and 020.7 (TE-1 Refinements, commit `12c32f9` / merge `41c86f4`) as landed in this turn, but the TODO 020 file shows 020.6 (path-reference sweep, 27 TE files modified, 41 line changes, Cat-1a/Cat-1b classified per DI-020-20260502-232651, with Cat-1b matches enumerated and left untouched per DF-35.3) and 020.10 (Status field retrofit on all 35 TEs) were also completed by 2026-05-03 / 2026-05-02 as part of the same batch. After turn 153, only 020.8 (README reconfirmation) remained open in TODO 020. Spec check OK; both twigs deleted; ppx/main pushed at `4e3c8e6`. No leftover threads.
- [ ] 021.154 Turn 154 (2026-05-03 01:17 UTC) -- "next?" -- (move to TE-36 apparatus vs specimen).
- [ ] 021.155 Turn 155 (2026-05-03 01:38 UTC) -- TE-36 DF-36.1 answered Alt-1.A (strict) and pushed back: "what the heck is projection mode?".
- [ ] 021.156 Turn 156 (2026-05-03 01:39 UTC) -- "wait. this is not harness-specific. this is only one possible wire envelope." -- (envelope locality challenge).
- [ ] 021.157 Turn 157 (2026-05-03 01:42 UTC) -- "you're confused. we're looking for one message envelope that works with all transports..." -- (clarified: one harness-wide envelope, multiple transports).
- [ ] 021.158 Turn 158 (2026-05-03 01:46 UTC) -- "harness-spec is the right home..." -- (envelope stays in harness-spec; carve-out scope shrinks).
- [ ] 021.159 Turn 159 (2026-05-03 01:50 UTC) -- "you got it right. proceed with step 1." -- (carve-out audit committed at `4725b3e`).
- [ ] 021.160 Turn 160 (2026-05-03 02:24 UTC) -- "do it and continue" -- (TE-36 drafted on twig, DF-36.1 through DF-36.7 walking begins).
- [ ] 021.161 Turn 161 (2026-05-03 02:57 UTC) -- "what's the difference between promise-stack and grid-pcid-payload?" -- (raises OQ-36.6).
- [ ] 021.162 Turn 162 (2026-05-03 03:12 UTC) -- "i suspect that promise-stack is too complicated and that you invented it..." -- (recorded as OQ-36.6 on TE-36 twig at `1a75b5b`).
- [ ] 021.163 Turn 163 (2026-05-03 03:34 UTC) -- "that's unreadable. please format it better." -- (formatting fix request; reply not visible in committed corpus).
- [ ] 021.164 Turn 164 (2026-05-03 11:20 UTC) -- "We need to get file based transport working ASAP so JJ and stevegt can use it to collaborate" -- (PIVOT: TE-36 parked at `0230c20`; first TE-37 (git-file-transport) drafted then reverted).
- [ ] 021.165 Turn 165 (2026-05-03 11:21 UTC) -- "JJ is a human collaborator using Claude in a clone of the repo. Don't mention JJ by name in any docs." -- (anonymity rule for committed docs).
- [ ] 021.166 Turn 166 (2026-05-03 11:29 UTC) -- "The group is at least 3 developer agents, and should be named more generically..." -- (slug = `wire-lab-devs`).
- [ ] 021.167 Turn 167 (2026-05-03 11:37 UTC) -- "Use .txt instead of .msg. Members should fetch from all branches but only post to their own author-id branch." -- (committed `463af44`/`ba9f2a6`).
- [ ] 021.168 Turn 168 (2026-05-03 11:45 UTC) -- "Think. do not sequentially number messages. Use the CID of the message as the file name." -- (committed `19e9b37`/`7a903d4`).
- [ ] 021.169 Turn 169 (2026-05-03 11:50 UTC) -- "Do we really need message-id?" -- (committed `31573dc`/`8d3bf04`; current ppx/main HEAD).
- [ ] 021.170 Turn 170 (2026-05-03 16:53 UTC) -- "Shouldn't the group named directory be nested inside a generic git file transport directory?" -- (DISCUSSED, NOT WRITTEN).
- [ ] 021.171 Turn 171 (2026-05-03 16:56 UTC) -- "think. Do we need another layer in there that means git file transfer..." -- (DISCUSSED, NOT WRITTEN).
- [ ] 021.172 Turn 172 (2026-05-03 17:01 UTC) -- "Think from a different angle. What if we wanted to run group sessions over rsync..." -- (substrate-pluralism reframe; DISCUSSED, NOT WRITTEN).
- [ ] 021.173 Turn 173 (2026-05-03 17:06 UTC) -- "Is there precedent for this in practice, in RFCs, in historical networks?" -- (Usenet/Sendmail/email-MTA precedents canvassed; DISCUSSED, NOT WRITTEN).
- [ ] 021.174 Turn 174 (2026-05-03 17:13 UTC) -- "I'm suspicious of the word binding; can't remember using it in Usenet gateways..." -- (drop "binding" for substrates; "feed" candidate; DISCUSSED, NOT WRITTEN).
- [ ] 021.175 Turn 175 (2026-05-03 17:25 UTC) -- **CORRECTION**: "Grid pCID is not the carrier line; it's the grid envelope -- fix that wherever we say it. Feed is okay, but you also used substrate, which I also like. I'm starting to think transports is misnamed. We're not really representing sites very well. We're not representing the decentralized CAS." -- (FOUR locked verbal corrections; NOT YET WRITTEN).
- [ ] 021.176 Turn 176 (2026-05-03 18:05 UTC) -- **CORRECTION**: "i like your layer 0-3 -- perhaps those should be numbered 4-7? i like 'groups' instead of 'forums'. game out 1:1 ephemeral flows like TCP/websocket. consider promisegrid-over-usenet and usenet-over-promisegrid. nested envelopes of different types is fine. **draft--slug.md or draft--slug.d should instead be slug-draft.md or slug-draft.d to align with slug-CID.md.** Symlinks into a CAS directory; rabin chunking of all CAS content; merkle trees. Foundational CID promise: 'I promise that this CID contains a unique and accurate hash of this content.'" -- (SEVEN locked verbal corrections including the slug-state rename; NOT YET WRITTEN).
- [ ] 021.177 Turn 177 (2026-05-03 18:34 UTC) -- "layers 5-7 would be less confusing... TCP/websocket -- the ephemeral version is a session... PromiseGrid-over-Usenet vs Usenet-over-PromiseGrid... a group = a named instance of a group protocol..." -- (layered model details; DISCUSSED, NOT WRITTEN).
- [ ] 021.178 Turn 178 (2026-05-03 19:14 UTC) -- "you need to assume that no site has all CAS objects -- think big... layer 5 promises need to be able to be economic... could a promisegrid app replace BGP?" -- (sparse-CAS, promise-economy framing, BGP question; DISCUSSED, NOT WRITTEN).
- [ ] 021.179 Turn 179 (2026-05-03 19:49 UTC) -- "yes, discovery and transport of chunks are both via promises... 'i will send you the chunk if you promise to not send it to anyone not in group X'. capability tokens MAY be transferable; floating exchange rates; 'everyone is their own central bank'. don't hardcode 'jj' yet." -- (promise-economy details; DISCUSSED, NOT WRITTEN).
- [ ] 021.180 Turn 180 (2026-05-03 20:21 UTC) -- "wait. you read promisebase .md files instead of the code? check git log..." -- (CONTEXT-COMPRESSION POINT: pivot to promisebase code review; everything from turns 170-179 starts to drop out of focus).
- [ ] 021.181 Turn 181 (2026-05-04 05:00 UTC) -- "Read the code, focus on the db/ directory" -- (promisebase code review; randStream regression identified).
- [ ] 021.182 Turn 182 (2026-05-04 05:05 UTC) -- "That used to work" -- (randStream Go-1.20 deprecation traced).
- [ ] 021.183 Turn 183 (2026-05-04 05:27 UTC) -- "Do it" -- (apply randStream fix; commit `d98b5d3` to promisebase main).
- [ ] 021.184 Turn 184 (2026-05-04 05:29 UTC) -- "Do 2" -- (read rest of promisebase code).
- [ ] 021.185 Turn 185 (2026-05-04 14:28 UTC) -- "are you able to hold and use two different PATs, one for wire-lab repo and one for promisebase repo?" -- (yes; PAT separation).
- [ ] 021.186 Turn 186 (2026-05-04 14:29 UTC) -- promisebase PAT supplied -- (saved at `.secrets/gh-pat-promisebase`).
- [ ] 021.187 Turn 187 (2026-05-04 14:32 UTC) -- "i'm solo in promisebase right now. let's hold other promisebase work for now and go back to TE-38" -- (PIVOT BACK; "TE-38" here was bot's number, see 021.192).
- [ ] 021.188 Turn 188 (2026-05-04 15:12 UTC) -- "Wait. Did you push the promise base changes?" -- (yes, randStream fix is on promisebase main).
- [ ] 021.189 Turn 189 (2026-05-04 15:23 UTC) -- "Can you examine the other" -- (truncated; resumed in 021.190).
- [ ] 021.190 Turn 190 (2026-05-04 15:24 UTC) -- "Can you examine the other promisebase branches?" -- (only `main` on remote; 2021 mob-consensus branches were merged and deleted).
- [ ] 021.191 Turn 191 (2026-05-04 15:54 UTC) -- "Be skeptical of what you see in promisebase in terms of design. Consider it a prototype at best. Any conflict between wire-lab and promisebase design choices should be discussed but should prefer wire-lab." -- (canon rule).
- [ ] 021.192 Turn 192 (2026-05-04 16:00 UTC) -- "I do intend to ref, factor, modernize, and use promisebase as one possible layer in promisegrid. If they are all committed to the repo, I'd like to pause on 38 and go back to dogfooding our message transport for collaboration between our group of developers. We need to get that going ASAP so that I'm not working on this solo." -- (NEW PIVOT: dogfood first, big TE later).
- [ ] 021.todo12 Investigate what happened to TODO 12. The master TODO list shows TODO 12 = `protocols/group-session.d/TODO/TODO-20260501-045543-group-transport-envelope.md`. Verify the file exists, read its current state, and determine whether it has been completed, superseded, folded, retired, or merely stalled. If unfinished, capture remaining work and current owner; if finished, mark the file accordingly.
- [ ] 021.pre149 Investigate turns earlier than 149 (turns 1-148, 2026-04-28 through 2026-05-03 00:05 UTC) for additional dropped threads. Walk the turn index from `past_session_contexts/sessions/2026-05-04_2026-05-10/ea135ce8/conversation.md` lines 5-152, classify each turn by status (WRITTEN to repo / PARKED on twig / DROPPED in conversation / REVERTED), and append any newly-discovered dropped threads to this TODO as additional `021.<turn>` subtasks for full-rigor walking. Cross-reference against the existing wire-lab git log for the period through 2026-05-03 00:05 UTC.

## Notes between turns 192 and 021.0

After turn 192 the user (in this fresh session) raised three correction rounds about TE numbering:

- "Wait, are you sure that 38 is not the right number? ... Did we accidentally drop conversation that was supposed to be in 37?" -- bot confirmed TE-37 was never committed; two prior false-start drafts.
- "you have the name wrong" -- bot's first TE-37 attempt used `transports/draft--wire-lab-devs/` instead of the verbally-locked `wire-lab-devs-draft` form.
- "Stop. I thought 37 was about integration of Promisebase." -- numbering confusion: there are at least three different things that have been called "TE-37," and zero of them are committed.
- "You still have the wire lab devs directory name wrong... It's as if as soon as I mentioned promisebase you compressed context and forgot everything we were working on. Please thoroughly review the entire session history in detail between TE 35 and the first mention of promisebase."

This TODO is the response to that final instruction.

## Decision Intent Log

(empty until turns are walked and DIs are filed)

## Outstanding cross-cutting questions tracked by this TODO

These are the questions that span multiple turns. They are listed here so they don't get lost between subtask walks. Each will be resolved or forwarded into a successor TODO/TE before this TODO closes.

1. What goes into TE-37? At least three candidate scopes: (a) promisebase integration / federation layer; (b) layered model + vocabulary fixes; (c) the `wire-lab-devs` git-file-transport bootstrap (already done as committed instance, not a TE).
2. What goes into TE-38? If TE-37 is "promisebase integration," then TE-38 was originally numbered as that and Steve's "go back to TE-38" at turn 187 referred to it. The numbering needs to be settled.
3. Should the verbal Cat-2 sweeps (`grid envelope`/`carrier line`, `draft--slug`/`slug-draft`, `transports`/`groups`) be done as one big sweep or many smaller ones? Each cascades CIDs in committed messages.
4. Pivot priority: dogfood first (Steve's most recent directive) vs. lock the model first.
5. promisebase canonicality: "wire-lab is canon, promisebase is prototype" needs to be written into a wire-lab spec (and ideally promisebase too).
