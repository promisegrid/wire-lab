# DR-davod - PromiseGrid kernel developer porting boundary

DR-ID: DR-davod
Date: 2026-05-08 17:25:12
Asked by: stevegt@t7a.org (Steve Traugott)
State: open
Question: What stable PromiseGrid infrastructure boundary should the Development Guide present to kernel developers who are porting PromiseGrid to a new platform or language?
Why this blocks progress: Wire-lab has a strong apparatus/specimen split and layered substrate model, but the guide still needs an explicit developer-facing porting target: which frozen specs, runtime expectations, conformance claims, and implementation obligations a port must satisfy.
Affects: `DEV-GUIDE-RESOURCES.md`; `/home/stevegt/lab/promisegrid-dev-guide/README.md`; `/home/stevegt/lab/promisegrid-dev-guide/TODO/TODO-binap-readme-outline-lock.md`; `protocols/wire-lab.d/specs/harness-spec-draft.md`; `transports/README.md`; `implementations/README.md`
Unblocks: `TODO-rozas.7`; `TODO-binap.7`
Waiting on: stevegt@t7a.org (Steve Traugott)
Decision:
Linked DI: DI-zalak (wire-lab guide-resource guidance; DR remains open); DI-funaf (TE-mazop/SIM-fovip evidence packet; DR remains open); DI-gumum (TE-pudiv app/kernel grid-message evidence; DR remains open); DI-somok (TE-dunas prior-art influence guardrails; DR remains open); DI-kuvum (TE-gakoh local-view promise/event hypergraph synthesis; DR remains open)
Related commits:
Last updated: 2026-05-25 10:01:19

## Event log

- 2026-05-08 17:25:12 — Opened from `TODO-rozas.9` so kernel-developer guide resources can cite an explicit DR for unsettled porting boundaries.
- 2026-05-08 20:32:35 — Narrowed current wire-lab answer for `FB-vitih`: the porting target is not the wire-lab harness. A PromiseGrid port should be framed as a runtime/implementation that hosts pCID-selected protocol handlers, implements the frozen binding/session/message specs it claims, and records those claims in implementation conformance records. This does not close the DR; the first required frozen spec set, runtime expectations, and implementation obligations remain open.
- 2026-05-24 14:03:42 — Added `TE-jimar` as current needs-DF evidence for the kernel/runtime portability boundary. Its provisional recommendation is to define "kernel" as a portable local infrastructure role/profile around pCID-selected messages, exact bytes, local promises, peer-local trust assessment, and evidence records, with daemon, microkernel, host-runtime, header/library-only, and split-object shapes as profiles rather than a universal process shape. This does not close the DR.
- 2026-05-25 08:26:41 — Added `TE-mazop` as current needs-DF evidence that `TE-jimar` is not enough to close this DR. The missing packet is the app/kernel/host promise surface and the minimum credible first-port contract. `SIM-fovip` is the intended concrete simulation to answer that packet. This does not close the DR.
- 2026-05-25 08:43:37 — Added `TE-pudiv` as follow-up needs-DF evidence that the app/kernel boundary may use the same `grid([42(pCID), payload, ...])` message format as the wire boundary. This narrows the next simulation question without closing the DR.
- 2026-05-25 09:12:02 — Added `TE-dunas` as needs-DF evidence for which lower-patent-risk distributed-OS prior art may influence PromiseGrid kernel work. V, Amoeba, Plan 9 / 9P, and GNU Hurd translators are retained as pattern or simulation pressure; Spring and modern multikernel details are excluded from current influence. This does not close the DR.
- 2026-05-25 10:01:19 — Added `TE-gakoh` as needs-DF evidence for reconciling decentralized-mainframe UX, local trust-filtered views, promise-bound path references, Burgess-style non-human agents, IPLD-compatible hypergraph storage, and event/command sourcing. This does not close the DR.

## Next DF packet

`TE-mazop` narrows the next unanswered packet without deciding it. Before this
DR closes, wire-lab needs evidence for:

- the smallest profile declaration that makes a first PromiseGrid port credible;
- required app-facing promises for storage, compute, networking, key use, device
  access, lifecycle, pCID dispatch, and evidence recording;
- host/runtime assumptions that must stay separate from PromiseGrid promises;
- evidence records for kept, refused, unavailable, or broken app-facing promises;
- the minimum pCID-selected spec coverage required for a first port to claim
  support, including explicit unsupported-pCID and unsupported-role handling.
- whether app/kernel operations are themselves pCID-selected grid messages, and
  whether local APIs are only adapters over those messages and evidence records.
- how much V-style same-message IPC, Amoeba-style capability/resource services,
  Plan 9-style uniform resource views, and Hurd-style replaceable/deferred
  services should influence the sim packet without importing transparency,
  global permission, or filesystem-authority assumptions.
- whether the kernel/app model should expose local file-like views over
  CID-rooted promise-bound references into a shared promise/event hypergraph,
  whether trusted groups should be able to maintain voluntary namespace
  frontiers by reciprocal promises, and whether IPLD-compatible objects should
  carry that graph without owning PromiseGrid semantics.

`SIM-fovip-kernel-promise-boundary-port-contract` is the current simulation home
for answering this packet, but `TE-pudiv`, `TE-dunas`, and `TE-gakoh` leave open
whether `SIM-fovip` should be extended or succeeded by a narrower same-grid
app/kernel boundary, prior-art-influence, and local-view hypergraph sim.
`DR-davod` remains open until that evidence is reviewed and a later DI locks the
guide-facing porting boundary.
