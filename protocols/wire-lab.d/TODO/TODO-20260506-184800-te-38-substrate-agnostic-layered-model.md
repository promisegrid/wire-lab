# TODO 22: TE-38 substrate-agnostic layered model

Parent TODO for TE-38 work on the wire-lab branch
`ppx/te-20260506-184800-substrate-agnostic-layered-model`. TE-38 is a
small vocabulary refactor of TE-29 plus a 100-year-goal citation, not
a redesign (per the 2026-05-06 scope confirmation; see Q-22.3 below).

## Status

In progress. Twig: `ppx/te-20260506-184800-substrate-agnostic-layered-model`.

## Foundational invariants (referenced, not redecided)

These are settled in memory and in prior TEs. TE-38 cites them; it
does not reopen them.

- **Layer assignments** (memory `projects.promisegrid.protocol_layers`,
  2026-05-03): L7 = group/forum protocols (top, application-most);
  L6 = CAS protocols (Rabin chunking, Merkle trees, CIDv1 codec
  distinguishing leaf vs dag-cbor; foundational pCID promise "I
  promise the bytes I serve under this CID hash to this CID");
  L5 = feed protocols (replicating CAS chunks between sites).

- **100-year goal invariants** (TE-28; memory
  `projects.promisegrid.design_principles`, 2026-05-01): no central
  registry; multi-generational durability; adversarial-by-default;
  protocol forking is normal; trust accrues per-burden; signing key
  is the only structural lock.

- **Existing architecture** (TE-29, 2026-05-01): each protocol is
  a simulated repo at `protocols/<slug>.d/` plus frozen siblings
  `protocols/<slug>-<pcid>.{md,d/}`; transport simulation uses 5-
  level paths
  `transports/<wire>/<binding-pCID>/<session-pCID>/<message-pCID>/<msg-id>.txt`
  where the path itself declares the instance/stack.

- **Vocabulary** (memory `projects.promisegrid.vocabulary.wire`,
  2026-05-03): "feed" not "binding"; "groups" not "forums" (or
  "transports", "channels"); "grid envelope" not "carrier line";
  slug-state is `<slug>-draft` then `<slug>-<cid>`, not
  `draft--<slug>`.

- **Sparse CAS** (memory `projects.promisegrid.architecture.sparse_cas`,
  2026-05-03): no site has all content; every site is sparse.

## Question log

Per the AGENTS-ppx Question-logging discipline (added 2026-05-06 at
commit 10ecc77), every question asked of Steve in TE-38 work is
logged here at the moment of asking and checked off only after the
resulting product is committed and pushed.

- [~] **Q-22.1** DF-38.A: layered-model foundation (full model + numbered layers / substrate-axis only, defer / full model + layer numbers as display nicknames)
    opened: 2026-05-06 18:48 UTC
    asked of: stevegt@t7a.org
    blocks: TE-38 scope and DF order
    alternatives: Alt-A.1 / Alt-A.2 / Alt-A.3
    recommendation: Alt-A.3
    retracted: 2026-05-06 -- L5/L6/L7 layer assignments and 100-year-goal invariants are already settled in memory and in TE-28; the question was the wrong question. Bot retraction at Steve's prompt.

- [x] **Q-22.2** DF-38.1: are L5 feeds a first-class protocol family parallel to L7 group protocols?
    opened: 2026-05-06 19:55 UTC
    asked of: stevegt@t7a.org
    blocks: TE-38 tree-shape DFs (where L5 protocol directories live)
    alternatives: Alt-1.A (top-level under protocols/) / Alt-1.B (under protocols/feeds/) / Alt-1.C (no separate L5 family)
    recommendation: Alt-1.A
    answered: 2026-05-06 21:25 UTC -- Alt-1.A
    resolved: 2026-05-06 @ 10ecc77
    product: this TODO file's "Foundational invariants" section cites the lock; further codification lands in the TE-38 doc itself once drafted (TE doc product will move resolution forward to a later commit if needed)
    note: per discipline §5, "resolved" requires committed-and-pushed product. The Alt-1.A lock is captured in this TODO file (committed) and in memory (`projects.promisegrid.protocol_layers` updated 2026-05-06). The TE-38 doc will further codify it. If a strict reading of §5 requires the TE doc itself before checking off, downgrade to [ ] with `answered:` only and re-resolve at TE-38 doc landing.

- [~] **Q-22.3** DF-38.2: how does an instance declare its bindings to L5 feeds?
    opened: 2026-05-06 22:34 UTC
    asked of: stevegt@t7a.org
    blocks: TE-38 instance shape
    alternatives: Alt-2.A (feeds/ subdir) / Alt-2.B (INSTANCE.md field) / Alt-2.C (top-level protocols/bindings.d/)
    recommendation: Alt-2.A
    retracted: 2026-05-06 -- TE-29 already locked the answer: the directory path under transports/ IS the instance/stack declaration; no separate per-instance feeds/ subdirectory or INSTANCE.md field is needed. Bot retraction at Steve's prompt to "dig deeper -- don't we have this in TE already?".

- [x] **Q-22.4** TE-38 scope: vocabulary refactor + 100-year-goal citation, OR larger?
    opened: 2026-05-06 22:36 UTC
    asked of: stevegt@t7a.org
    blocks: TE-38 DF count and structure
    alternatives: (a) vocabulary refactor + citation / (b) larger / (c) reframe scope
    recommendation: (a) -- but framed as a check rather than a recommendation
    answered: 2026-05-06 22:56 UTC -- (a) vocabulary refactor + citation
    resolved: 2026-05-06 @ 10ecc77
    product: this TODO file (Foundational invariants section + question log records the scope lock); TE-38 doc will be the further codification

- [x] **Q-22.5** TE-29 update: Cat-3 entry / supersede / leave alone?
    opened: 2026-05-06 22:36 UTC
    asked of: stevegt@t7a.org
    blocks: TE-38 -- TE-29 reconciliation mechanism
    alternatives: (a) Cat-3 entry on TE-29 / (b) supersede TE-29 / (c) leave TE-29 alone
    recommendation: (a) -- but framed as a check rather than a recommendation
    answered: 2026-05-06 22:56 UTC -- (a) Cat-3 entry on TE-29
    resolved: 2026-05-06 @ 10ecc77
    product: this TODO file records the lock; TE-29 Cat-3 entry will be drafted in the TE-38 commit cycle

- [ ] **Q-22.6** DF-38.M: where does L6 CAS fit in TE-29's path scheme?
    opened: 2026-05-06 22:58 UTC
    asked of: stevegt@t7a.org
    blocks: TE-38 -- whether the path scheme needs an L6 level or stays 5-deep
    alternatives: Alt-M.1 (L6 implicit; messages are leaf chunks) / Alt-M.2 (explicit L6 path level, 6-deep paths) / Alt-M.3 (cas/ subtree + pointer files, with Alt-M.1 as migration path)
    recommendation: Alt-M.3 long-term, Alt-M.1 today
    note: 2026-05-06 -- explanation of all three alternatives delivered to Steve at 15:08 PT on his request; question stands open pending answer.

- [x] **Q-22.7** OPEN-THREADS.md migration scope: which deprecation shape?
    opened: 2026-05-06 23:08 UTC
    asked of: stevegt@t7a.org
    blocks: deprecating OPEN-THREADS.md and consolidating thread tracking into TODO files
    alternatives:
      Alt-OT.A -- one-row-per-thread: each open T-* thread becomes its own
                  TODO file under protocols/<owner>.d/TODO/ (17 new files);
                  closed threads recorded in a single OPEN-THREADS-archive.md
      Alt-OT.B -- one TODO file per topical bundle: group related threads
                  (e.g. T-FILENAME-CID-CASCADE goes inside the TE-42 parent
                  TODO; T-MIG-OPS inside the TE-39 parent TODO; etc.) using
                  the same parent-TODO pattern as TODO 22 for TE-38
      Alt-OT.C -- single new TODO 23 'thread index' that absorbs OPEN-THREADS
                  verbatim as a question/thread log; existing closed entries
                  carry over with [x] preserved
    recommendation: Alt-OT.B -- it matches the structure already established
                    by TODO 22 (TE-38) and the dependency-sorted TE roster
                    in dropped-thread-disposition-20260506.md. Each TE that
                    will be drafted (TE-39 through TE-45 plus TE-36-followon)
                    gets its own parent TODO file with the relevant T-*
                    threads merged in as content. Anticipated-future-TE
                    threads (T-RING-TRANSPORT, T-CLUSTER-OF-CLUSTERS-TRANSPORT,
                    T-GOSSIP-TRANSPORT, T-RECEIPTS-AT-SCALE) get a single
                    'anticipated TEs' TODO since they're not yet scheduled.
                    T-021-CC-Q3/Q4/Q5 are cross-cutting and stay in TODO 21
                    (already there). The deprecation step renames
                    OPEN-THREADS.md -> OPEN-THREADS-DEPRECATED.md with a
                    pointer header explaining the migration.
    answered: 2026-05-07 00:22 UTC -- Alt-OT.B
    resolved: 2026-05-07 @ 5670eca (wire-lab) and 3bed63a (session-logs)
    product: protocols/wire-lab.d/TODO/TODO-20260507-002306-te-{36-followon,39,40,41,42,43,44,45}-*.md, TODO-20260507-002306-anticipated-future-tes-transport-family.md, AGENTS-ppx.md (deprecation note)

- [x] **Q-22.8** OPEN-THREADS.md deprecation mechanics: rename or delete?
    opened: 2026-05-06 23:08 UTC
    asked of: stevegt@t7a.org
    blocks: completion of OPEN-THREADS.md deprecation
    alternatives:
      Alt-D.A -- rename to OPEN-THREADS-DEPRECATED.md with a header pointing
                 readers to the new TODO-based scheme; preserves git history
                 and avoids breaking any reference that cites the old name
      Alt-D.B -- delete OPEN-THREADS.md outright; rely on git history for
                 the audit trail; smaller tree, less clutter
      Alt-D.C -- replace contents with a tombstone pointer (1-page file
                 explaining the deprecation, listing where each thread
                 migrated to); keeps the filename live as a redirect
    recommendation: Alt-D.C -- a tombstone is the most discoverable for
                    anyone (human or bot) who searches for OPEN-THREADS.md
                    and finds the new file, while keeping the audit trail
                    in git. Alt-D.A keeps the renamed file but anyone
                    citing 'OPEN-THREADS.md' will get a 404; Alt-D.B loses
                    the redirect benefit entirely.
    answered: 2026-05-07 00:22 UTC -- Alt-D.B (delete outright)
    resolved: 2026-05-07 @ 3bed63a (session-logs wire-lab branch)
    product: deletion commit on stevegt/session-logs wire-lab branch; AGENTS-ppx.md (commit 5670eca on wire-lab ppx/main) deprecation note explains where each thread went so readers searching git history can find them

## Open questions retired by this TODO

(None yet; backfilled DF questions Q-22.1 through Q-22.5 are recorded
above with their actual asked-and-answered timeline.)

## Decision Intent Log

(Will be populated as DFs lock and product lands.)
