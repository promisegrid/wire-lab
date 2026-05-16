# TODO-vunub: TE-sihih substrate-agnostic layered model

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TODO-22` (integer alias)
- `TODO-20260506-184800` (timestamp alias and pre-migration filename)

Parent TODO for TE-38 work on the wire-lab branch
`ppx/te-20260506-184800-substrate-agnostic-layered-model`. TE-38 is a
small vocabulary refactor of TE-29 plus a 100-year-goal citation, not
a redesign (per the 2026-05-06 scope confirmation; see Q-22.3 below).

## Status

Closed. TE-sihih is decided and the product landed at
`docs/thought-experiments/TE-sihih-substrate-agnostic-layered-model.md`;
TE-vipir received the companion Cat-3 refinement. The historical working twig
was `ppx/te-20260506-184800-substrate-agnostic-layered-model`. Later migration
or root-layout questions are owned by successor TEs, DRs, and simulation-local
TODOs rather than by this parent TODO.

## Foundational invariants (referenced, not redecided)

These are settled in memory and in prior TEs. TE-sihih cites them; it
does not reopen them.

- **Layer assignments** (memory `projects.promisegrid.protocol_layers`,
  2026-05-03): L7 = group/forum protocols (top, application-most);
  L6 = CAS protocols (Rabin chunking, Merkle trees, CIDv1 codec
  distinguishing leaf vs dag-cbor; foundational pCID promise "I
  promise the bytes I serve under this CID hash to this CID");
  L5 = feed protocols (replicating CAS chunks between sites).

- **100-year goal invariants** (TE-dajot; memory
  `projects.promisegrid.design_principles`, 2026-05-01): no central
  registry; multi-generational durability; adversarial-by-default;
  protocol forking is normal; trust accrues per-burden; signing key
  is the only structural lock.

- **Existing architecture** (TE-vipir, 2026-05-01): each protocol is
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
commit 10ecc77), every question asked of Steve in TE-sihih work is
logged here at the moment of asking and checked off only after the
resulting product is committed and pushed.

- [~] **Q-22.1** DF-38.A: layered-model foundation (full model + numbered layers / substrate-axis only, defer / full model + layer numbers as display nicknames)
    opened: 2026-05-06 18:48 UTC
    asked of: stevegt@t7a.org
    blocks: TE-sihih scope and DF order
    alternatives: Alt-A.1 / Alt-A.2 / Alt-A.3
    recommendation: Alt-A.3
    retracted: 2026-05-06 -- L5/L6/L7 layer assignments and 100-year-goal invariants are already settled in memory and in TE-dajot; the question was the wrong question. Bot retraction at Steve's prompt.

- [x] **Q-22.2** DF-38.1: are L5 feeds a first-class protocol family parallel to L7 group protocols?
    opened: 2026-05-06 19:55 UTC
    asked of: stevegt@t7a.org
    blocks: TE-sihih tree-shape DFs (where L5 protocol directories live)
    alternatives: Alt-1.A (top-level under protocols/) / Alt-1.B (under protocols/feeds/) / Alt-1.C (no separate L5 family)
    recommendation: Alt-1.A
    answered: 2026-05-06 21:25 UTC -- Alt-1.A
    resolved: 2026-05-06 @ 10ecc77
    product: this TODO file's "Foundational invariants" section cites the lock; further codification lands in the TE-sihih doc itself once drafted (TE doc product will move resolution forward to a later commit if needed)
    note: per discipline §5, "resolved" requires committed-and-pushed product. The Alt-1.A lock is captured in this TODO file (committed) and in memory (`projects.promisegrid.protocol_layers` updated 2026-05-06). The TE-sihih doc will further codify it. If a strict reading of §5 requires the TE doc itself before checking off, downgrade to [ ] with `answered:` only and re-resolve at TE-sihih doc landing.

- [~] **Q-22.3** DF-38.2: how does an instance declare its bindings to L5 feeds?
    opened: 2026-05-06 22:34 UTC
    asked of: stevegt@t7a.org
    blocks: TE-sihih instance shape
    alternatives: Alt-2.A (feeds/ subdir) / Alt-2.B (INSTANCE.md field) / Alt-2.C (top-level protocols/bindings.d/)
    recommendation: Alt-2.A
    retracted: 2026-05-06 -- TE-vipir already locked the answer: the directory path under transports/ IS the instance/stack declaration; no separate per-instance feeds/ subdirectory or INSTANCE.md field is needed. Bot retraction at Steve's prompt to "dig deeper -- don't we have this in TE already?".

- [x] **Q-22.4** TE-sihih scope: vocabulary refactor + 100-year-goal citation, OR larger?
    opened: 2026-05-06 22:36 UTC
    asked of: stevegt@t7a.org
    blocks: TE-sihih DF count and structure
    alternatives: (a) vocabulary refactor + citation / (b) larger / (c) reframe scope
    recommendation: (a) -- but framed as a check rather than a recommendation
    answered: 2026-05-06 22:56 UTC -- (a) vocabulary refactor + citation
    resolved: 2026-05-06 @ 10ecc77
    product: this TODO file (Foundational invariants section + question log records the scope lock); TE-sihih doc will be the further codification

- [x] **Q-22.5** TE-vipir update: Cat-3 entry / supersede / leave alone?
    opened: 2026-05-06 22:36 UTC
    asked of: stevegt@t7a.org
    blocks: TE-sihih -- TE-vipir reconciliation mechanism
    alternatives: (a) Cat-3 entry on TE-vipir / (b) supersede TE-vipir / (c) leave TE-vipir alone
    recommendation: (a) -- but framed as a check rather than a recommendation
    answered: 2026-05-06 22:56 UTC -- (a) Cat-3 entry on TE-vipir
    resolved: 2026-05-06 @ 10ecc77
    product: this TODO file records the lock; TE-vipir Cat-3 entry will be drafted in the TE-sihih commit cycle

- [x] **Q-22.6** DF-38.M: where does L6 CAS fit in TE-vipir's path scheme?
    opened: 2026-05-06 22:58 UTC
    asked of: stevegt@t7a.org
    blocks: TE-sihih -- whether the path scheme needs an L6 level or stays 5-deep
    alternatives: Alt-M.1 (L6 implicit; messages are leaf chunks) / Alt-M.2 (explicit L6 path level, 6-deep paths) / Alt-M.3 (cas/ subtree + pointer files, with Alt-M.1 as migration path) / Alt-M.4 (cas/ subtree, all messages are CBOR pointers, no exceptions -- bot recommendation revised)
    recommendation: Alt-M.4 -- per Steve's 2026-05-06 17:54 PT note 'all
                    chunks should go into CAS; exceptions make it
                    complicated', the small-message-is-leaf exception in
                    Alt-M.1 / Alt-M.3 is dropped. Every message file in
                    transports/ is unconditionally a CBOR pointer file
                    {cas: <cas-protocol-pCID>, root: <chunk-cid>}; every
                    chunk lives in cas/<cas-protocol-pCID>/<chunk-cid>.
                    Today's inline-message files (transports/wire-lab-
                    devs-draft/ and similar) are migrated to pointer
                    form when the first L6 CAS spec lands (likely
                    TE-43 promisebase adoption).
    note: 2026-05-06 -- Alt-M explanation delivered at 15:08 PT;
          Steve flagged the small/large asymmetry at 17:54 PT;
          Alt-M.4 added with revised recommendation pending Steve's
          confirmation.
    answered: 2026-05-07 01:00 UTC -- Alt-M.4
    resolved: 2026-05-07 -- TE-sihih drafted with locked Alt-M.4 shape;
              TE-vipir Cat-3 Refinement filed citing TE-sihih; TODO-vunub
              Q-22.6 checked off.
    product: docs/thought-experiments/TE-sihih-substrate-agnostic-layered-model.md (new TE-sihih doc), docs/thought-experiments/TE-vipir-protocols-as-simulated-repos-and-binding-layer.md (Cat-3 Refinement section appended)

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
                  TODO; T-MIG-OPS inside the TE-mumuv parent TODO; etc.) using
                  the same parent-TODO pattern as TODO-vunub for TE-sihih
      Alt-OT.C -- single new TODO-lilok 'thread index' that absorbs OPEN-THREADS
                  verbatim as a question/thread log; existing closed entries
                  carry over with [x] preserved
    recommendation: Alt-OT.B -- it matches the structure already established
                    by TODO-vunub (TE-sihih) and the dependency-sorted TE roster
                    in dropped-thread-disposition-20260506.md. Each TE that
                    will be drafted (TE-mumuv through TE-45 plus TE-havib-followon)
                    gets its own parent TODO file with the relevant T-*
                    threads merged in as content. Anticipated-future-TE
                    threads (T-RING-TRANSPORT, T-CLUSTER-OF-CLUSTERS-TRANSPORT,
                    T-GOSSIP-TRANSPORT, T-RECEIPTS-AT-SCALE) get a single
                    'anticipated TEs' TODO since they're not yet scheduled.
                    T-021-CC-Q3/Q4/Q5 are cross-cutting and stay in TODO-lilar
                    (already there). The deprecation step renames
                    OPEN-THREADS.md -> OPEN-THREADS-DEPRECATED.md with a
                    pointer header explaining the migration.
    answered: 2026-05-07 00:22 UTC -- Alt-OT.B
    resolved: 2026-05-07 @ 5670eca (wire-lab) and 3bed63a (session-logs)
    product: protocols/wire-lab.d/TODO/TODO-sinuv-te-{36-followon,39,40,41,42,43,44,45}-*.md, TODO-sinuv-anticipated-future-tes-transport-family.md, AGENTS-ppx.md (deprecation note)

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

- The turn-170 `DF-37.1` flat-versus-nested root `transports/` framing is
  superseded for the current specimen by TE-sihih's L5/L6/L7 model, DR-nugog,
  and `DI-fakin`; future root/reference layout questions graduate through later
  DR/DI/spec work rather than through this TODO.
- The turn-171 substrate-path-layer and instance-manifest hooks are retired by
  Q-22.2 and Q-22.3: L5 feeds are first-class protocol families, while
  per-instance feed declarations remain path-as-declaration rather than a new
  manifest field for this TODO's scope. The current wire-lab-devs git binding
  remains inline in group-session §9 under `DI-rurab` until a later feed-spec
  graduation chooses otherwise.
- The turn-172 substrate-pluralism / `binding`-family proposal is retired by
  the same landed TE-sihih products: DF-38.1 / Q-22.2 locks L5 feed protocols
  as first-class, DF-38.2 / Q-22.3 retracts the per-instance feed-manifest
  mechanism in favor of path-as-declaration, and Q-22.6 / Alt-M.4 moves CAS
  concerns into L6 rather than leaving them as an unmodeled fourth family.
  Future extraction of the current git feed prose from group-session §9 is
  per-feed successor work, not open turn-172 recovery work.
- The turn-173 historical-precedent questions are retired for this TODO by
  `DI-pijun`, `docs/research/historical-networks-20260503.md`, and
  `simulations/SIM-hugoj-cas-usenetlike-gitlike/README.md`: the Usenet/git
  lineage is captured as exploratory design evidence, the rejected `binding`
  vocabulary is replaced by `feed`, the negative-precedent check is recorded,
  and CAS-cardinality is carried by the L5/L6/L7 model plus successor
  simulation work.
- Backfilled DF questions Q-22.1 through Q-22.8 are recorded above with their
  actual asked-and-answered or retracted timeline.

## Decision Intent Log

(No local DI entries were added before this file closed; the locks are recorded in the question log above and in the landed TE-sihih / TE-vipir products.)
