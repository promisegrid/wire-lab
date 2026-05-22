# GA Child Promotion Procedure

This procedure is the operator contract for promoting generated GA child
simulations into canonical simulation and result homes. Use it when Steve says:

```text
promote <child-proquint> [<child-proquint> ...]
```

For example, `promote natim maraz` means "find the generated child simulations
whose IDs start with `SIM-natim-child-` and `SIM-maraz-child-`, review their
evidence, and promote the selected designs using this procedure." Source:
`DI-dikoh`; `DI-zadik`; `DI-higot`.

## Authority Boundary

Promotion is a review decision. A generated child is not canonical merely because
it exists under `proposals/<run-group-id>/`. `tools/ga-runner accept` records the
reviewed proposal evidence in GA state, but it does not by itself create final
canonical simulation/result artifacts or stage files for commit. Under
`DI-higot`, scored simulation files and scored result JSON bytes are
append-only. A promotion procedure may not rewrite them in place to add cleanup
metadata, canonical IDs, or canonical-home notes. Source: `DI-lirat`;
`DI-dikoh`; `DI-zadik`; `DI-higot`.

The Jufag/Bimos promotion from `ga-canary-20260520-221953` is the historical
copy-style precedent. `DI-zadik` superseded that operational style with
move-cleanup intent, but `DI-higot` adds a harder constraint: scored artifact
bytes are append-only. A future promotion may move or copy accepted scored
artifacts unchanged into canonical homes, but this procedure must not instruct
post-score rewrites of those artifact bytes.

## 1. Resolve Requested Children

For each requested proquint:

1. Search `results/state/ga-canary-*.json` for child IDs matching
   `SIM-<proquint>-child-*`.
2. Require exactly one usable match unless Steve explicitly resolves the
   ambiguity.
3. The child status must normally be `generated`. `accepted` is allowed only when
   resuming an interrupted promotion that already passed the accept step.
4. Record the run group, child ID, proposal sim path, parent IDs, tree hash, and
   child result paths.
5. Ignore old queued, skipped, failed, or stale `running` children unless Steve
   explicitly asks to recover that run first.

If there are multiple matches for a proquint, stop before editing and report the
candidate run groups and child IDs.

## 2. Validate Candidate Evidence

Before canonical edits, verify each selected child:

1. The proposal sim path from GA state exists under
   `proposals/<run-group-id>/simulations/<child-id>/`.
2. The current proposal tree hash matches the state `tree_hash`.
3. Every selected child result JSON path exists under
   `proposals/<run-group-id>/results/<child-id>/<scenario>/<model>/<timestamp>.json`.
4. `tools/ga-runner validate` accepts the selected result JSON files.
5. Each result JSON's `source.sim_path` points at the proposal sim path and
   `source.simulation_tree_hash` matches the proposal tree.
6. The child has enough standing simulation material to be reviewable as scored
   evidence. At minimum this normally means `README.md` and `QUESTION.md`; if
   comparable sims in the same family carry local protocol files such as
   `protocols/<slug>.d/CHANGELOG.md`, `manifest.json`, or draft specs, record
   any missing canonicalization work as follow-up rather than editing the scored
   child in place.

Do not repair the proposal tree in place before `accept`; `accept` verifies the
raw scored child. After `accept`, do not repair the scored child in place
either. If canonicalization is required, route it through surrounding docs or a
future byte-identical move/copy path. The original proposal path in `source.*`
remains historical scored-source metadata.

## 3. Choose Final Names

Default final simulation ID:

```text
SIM-<same-proquint>-<child-slug-with-child-removed>
```

Examples:

- `SIM-natim-child-nested-payload-outer-attestation-multisig` becomes
  `SIM-natim-nested-payload-outer-attestation-multisig`, unless a clearer
  domain-qualified final name is needed.
- If adding a family prefix improves index consistency, use the smallest clear
  final name, such as
  `SIM-natim-grid-envelope-nested-payload-outer-attestation-multisig`.

Stop and ask Steve before editing when:

- the final name would collide with an existing `simulations/SIM-*` or
  `results/SIM-*` path;
- the generated slug is generic, misleading, or still contains `ga-child`,
  `pending`, or another process artifact;
- two promoted children should be merged instead of kept as separate competing
  sims.

## 4. Record the Promotion DI

Before canonical promotion edits, append a new DI to
`protocols/wire-lab.d/TODO/TODO-tapur-ga-runner-json-fitness-and-child-sim-search.md`.
The DI must name:

- selected run group and child IDs;
- final canonical simulation IDs;
- whether each child is promoted, deferred, or rejected;
- the move-cleanup policy for proposal artifacts;
- the scored-artifact handling policy;
- the provenance rule that `source.*` remains the historical scored-source path
  and hashes even after the proposal root is moved away;
- all expected canonical sim, result, state, and index paths.

Use that DI ID in promoted sim docs, result `promotion` metadata, and final
handoff evidence.

## 5. Record Acceptance in GA State

Run `accept` before moving canonical artifacts:

```bash
cd tools/ga-runner
go run . accept \
  -repo-root ../.. \
  -run-group-id <run-group-id> \
  -child <child-id> \
  -result proposals/<run-group-id>/results/<child-id>/<scenario-1>/<model>/<timestamp>.json \
  -result proposals/<run-group-id>/results/<child-id>/<scenario-2>/<model>/<timestamp>.json \
  -reviewer-note '<why this child is accepted; cite DI-<handle>>'
```

Repeat `-child` and `-result` as needed when promoting multiple children from the
same run. If children come from different run groups, run one `accept` command
per run group. The command should print review paths and remind the operator
that promotion is still required before `git add`.

## 6. Record Canonical Intent Without Rewriting Scored Bytes

For each accepted child:

1. Record the intended final canonical simulation ID in the promotion DI and
   reviewer note.
2. Keep useful provenance language in surrounding docs: parent IDs, run group,
   original child ID, and promotion DI.
3. Do not rewrite scored simulation files or scored result JSONs to replace
   child IDs, add cleanup notes, or add canonical result metadata.
4. If Steve explicitly wants canonical homes before a later byte-identical path
   is locked, stop and resolve whether to move or copy the scored bytes
   unchanged.

This keeps the scored artifact bytes historical while still letting TODOs,
indexes, and review notes say which child was accepted and what canonical name
is intended.

## 7. Update Indexes and Cross-References

Update the public navigation surfaces that make the new canonical sim
discoverable:

- the relevant TODO/DR/DI text when the promotion closes or routes open work;
- `simulations/README.md` only if it can point at scored evidence without
  pretending a rewritten canonical sim tree exists;
- `DEV-GUIDE-RESOURCES.md` when the accepted child materially changes the guide
  evidence picture and the reference can be made without claiming a rewritten
  canonical artifact already exists.

For family-specific sims, preserve local family grouping. For example,
grid-envelope promotions belong with the grid-envelope table, while guide-level
claim/conformance promotions should be indexed as their own simulation family
instead of being forced into grid-envelope language.

## 8. Validate Before Handoff

Run targeted validation on every scored JSON file whose surrounding docs or
state you touched:

```bash
cd tools/ga-runner
go run . validate -repo-root ../.. -result ../../proposals/<run-group-id>/results/<child-id>/<scenario-id>/<model-id>/<timestamp>.json
```

If a future byte-identical promotion path creates canonical result files, also
validate those copied/moved files explicitly. For larger batches, validate the
shared timestamp/model after the explicit file checks:

```bash
cd tools/ga-runner
go run . validate -repo-root ../.. -model <model-id> -timestamp <timestamp>
```

Then run repository-level documentation checks that are cheap and relevant:

```bash
git diff --check
```

If code changed, also run the code-specific checks required by `AGENTS.md`,
including comment-delta audit and `errcheck ./...` for Go changes. A normal
promotion should not require Go code changes.

## 9. Handoff

The final response must include:

- accepted child IDs and intended final sim IDs;
- accepted run group(s) and original child IDs;
- where the scored evidence currently lives;
- validation commands run;
- explicit statement that scored artifact bytes were left unchanged, or that any
  move/copy was byte-identical;
- Decision Compliance, Decision Matrix, Inline diff annotations, Runtime Path
  Touch Matrix, Comment audit, Intent provenance audit, and Exceptions.

Do not commit unless Steve says `commit`.

## 10. Cull Rejected Children Separately

Promotion moves accepted children and therefore does not use `cull` for those
children. Rejected or unpromoted children remain under `proposals/` until Steve
chooses to cull them. For rejected children, use `tools/ga-runner cull` with
explicit child IDs and reasons so destructive cleanup is state-bound and
auditable. Source: `DI-kofil`; `DI-dikoh`; `DI-zadik`; `DI-higot`.
