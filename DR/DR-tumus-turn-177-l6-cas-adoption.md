# DR-tumus - Turn-177 L6 CAS adoption profile

DR-ID: DR-tumus
Date: 2026-05-17 09:28:02
Asked by: stevegt@t7a.org (Steve Traugott)
State: open
Question: Which concrete L6 CAS profile should TE-43 lock for the first PromiseGrid CAS spec, including CBOR / DAG-CBOR profile, CIDv1 codec object typing, pointer-object shape, chunking algorithm and parameters, and promisebase / pitbase adoption stance?
Why this blocks progress: Turn 177 established that L7 group semantics should resolve through L6 CAS and L5 feeds should move chunks, but the exact CAS object model remains undecided. TODO-kituj / TE-43 cannot draft the first concrete L6 CAS spec, and TODO-pipus cannot design a safe pointer-and-CAS migration, until this profile has a DR/DF/DI path.
Affects: `protocols/wire-lab.d/TODO/TODO-kituj-te-43-promisebase-prior-art-adoption.md`; `simulations/SIM-jomag-cas-object-model/`; `simulations/SIM-zazit-chunk-feed-replication/`; `protocols/wire-lab.d/TODO/TODO-pipus-te-39-wire-lab-devs-migration.md`; `/home/stevegt/lab/promisebase/`.
Unblocks: TODO-kituj / TE-43 drafting; TODO-pipus additive pointer-and-CAS migration planning; TODO-dozak wire-lab / promisebase merge trajectory framing.
Waiting on: stevegt@t7a.org (Steve Traugott)
Decision:
Linked DI: DI-navod; DI-pator; DI-davov
Related commits:
Last updated: 2026-05-17 09:28:02

## Event log

- 2026-05-17 09:28:02 — Opened during turn-177 cleanup so concrete L6 CAS decisions are not left only as replay notes, TODO bullets, or simulation scenarios.

## Evidence

- Turn 177 corrected the layer order: L7 group protocols define semantics, L6 CAS stores/resolves chunks and pointer objects, and L5 feeds advertise/request/replicate chunks.
- `simulations/SIM-jomag-cas-object-model/SCENARIOS.md` records the concrete pressure cases: deterministic CBOR agreement, DAG-CBOR interop, CIDv1 object typing, pointer-object identity, chunker parameter mismatch, promisebase adapter, and small-object degenerate case.
- `simulations/SIM-zazit-chunk-feed-replication/SCENARIOS.md` records feed-side pressure that depends on usable CAS object identity under sparse replication.
- `TODO-kituj` owns TE-43 and already lists deterministic CBOR, allowed tags, chunking parameters, CIDv1 object typing, pointer-object shape, and promisebase prior-art adoption as scope.

## Candidate decisions

- **Alt-A: IPFS/IPLD-aligned default.** Prefer DAG-CBOR-compatible object nodes, public multicodec / CIDv1 values, and bridge-friendly defaults unless PromiseGrid has a concrete reason to diverge.
- **Alt-B: promisebase-adapter default.** Treat promisebase / pitbase's existing block/tree/stream model as the implementation substrate, then wrap or translate class headers, stream symlinks, and Rabin parameters into PromiseGrid CID / pointer semantics.
- **Alt-C: minimal first CAS profile.** Lock only raw chunks plus a minimal CBOR pointer object first, defer chunked Merkle objects and promisebase integration until the first additive migration proves the need.

## Notes

This DR does not choose the L6 profile. It gives TODO-kituj / TE-43 a concrete decision request to answer before drafting or freezing the first L6 CAS spec.
