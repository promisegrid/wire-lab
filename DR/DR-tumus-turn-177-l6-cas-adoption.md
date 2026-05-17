# DR-tumus - Turn-177 L6 CAS adoption profile

DR-ID: DR-tumus
Date: 2026-05-17 09:28:02
Asked by: stevegt@t7a.org (Steve Traugott)
State: open
Question: Which concrete L6 CAS profile should TE-43 lock for the first PromiseGrid CAS spec, including CBOR / DAG-CBOR profile, CIDv1 codec object typing, pointer-object shape, chunking algorithm and parameters, and promisebase / pitbase adoption stance?
Why this blocks progress: Turn 177 established that L7 group semantics should resolve through L6 CAS and L5 feeds should move chunks, but the exact CAS object model remains undecided. TODO-kituj / TE-43 cannot draft the first concrete L6 CAS spec, and TODO-pipus cannot design a safe pointer-and-CAS migration, until this profile has a DR/DF/DI path.
Affects: `protocols/wire-lab.d/TODO/TODO-kituj-te-43-promisebase-prior-art-adoption.md`; `simulations/SIM-jomag-cas-object-model/`; `simulations/SIM-zazit-chunk-feed-replication/`; `simulations/SIM-bobud-l6-cas-starting-profile-bakeoff/`; `simulations/SIM-kohad-cas-object-type-binding-bakeoff/`; `simulations/SIM-gobaz-chunking-identity-bakeoff/`; `protocols/wire-lab.d/TODO/TODO-pipus-te-39-wire-lab-devs-migration.md`; `/home/stevegt/lab/promisebase/`.
Unblocks: TODO-kituj / TE-43 drafting; TODO-pipus additive pointer-and-CAS migration planning; TODO-dozak wire-lab / promisebase merge trajectory framing.
Waiting on: stevegt@t7a.org (Steve Traugott)
Decision:
Linked DI: DI-navod; DI-pator; DI-davov; DI-majib; DI-bukoh; DI-molah
Related commits:
Last updated: 2026-05-17 10:09:28

## Event log

- 2026-05-17 09:28:02 — Opened during turn-177 cleanup so concrete L6 CAS decisions are not left only as replay notes, TODO bullets, or simulation scenarios.
- 2026-05-17 09:44:46 — Added unanswered next-DF packet and acceptance criteria under `DI-majib`.
- 2026-05-17 09:55:11 — Routed DF-tumus.1 through DF-tumus.3 through standalone bakeoff simulations after Steve asked for sims instead of direct answers.
- 2026-05-17 10:09:28 — Synthesized the three bakeoff simulations into a final answerable DR packet while leaving the DR open for Steve's decision.

## Evidence

- Turn 177 corrected the layer order: L7 group protocols define semantics, L6 CAS stores/resolves chunks and pointer objects, and L5 feeds advertise/request/replicate chunks.
- `simulations/SIM-jomag-cas-object-model/SCENARIOS.md` records the concrete pressure cases: deterministic CBOR agreement, DAG-CBOR interop, CIDv1 object typing, pointer-object identity, chunker parameter mismatch, promisebase adapter, and small-object degenerate case.
- `simulations/SIM-zazit-chunk-feed-replication/SCENARIOS.md` records feed-side pressure that depends on usable CAS object identity under sparse replication.
- `simulations/SIM-bobud-l6-cas-starting-profile-bakeoff/SCENARIOS.md` tests the starting-profile choice without preferring IPFS / IPLD alignment, promisebase adapter reuse, or a minimal pointer/raw seed.
- `simulations/SIM-kohad-cas-object-type-binding-bakeoff/SCENARIOS.md` tests CID codec-only typing, codec-plus-kind typing, and path suffixes as a negative control.
- `simulations/SIM-gobaz-chunking-identity-bakeoff/SCENARIOS.md` tests pCID-driven chunking, chunking-CID / cCID-style descriptors, profile negotiation, and raw-only deferral as exploratory alternatives.
- `TODO-kituj` owns TE-43 and already lists deterministic CBOR, allowed tags, chunking parameters, CIDv1 object typing, pointer-object shape, and promisebase prior-art adoption as scope.

## Bakeoff synthesis

The three `DI-bukoh` bakeoff simulations make the next decision packet narrow
enough to answer without treating any simulation as authoritative by itself.
This synthesis is not the decision; it is the recommended packet for Steve to
accept, reject, or refine. Source: DI-molah.

- **Starting profile recommendation:** choose a minimal pointer/raw first
  profile, while carrying DAG-CBOR / CIDv1 bridge constraints as compatibility
  pressure. The minimal profile tests pointer objects, sparse fetch, and
  historical-byte preservation without prematurely adopting an external stack or
  all promisebase / pitbase internals.
- **Type-binding recommendation:** make the CID codec the authoritative object
  type discriminator. Reject filename/path suffixes as type authority because
  they do not travel with content identity. Permit object-internal kind only
  when the chosen codec's payload schema defines it; it must not become a second
  independent type authority.
- **Chunking recommendation:** keep the first profile raw-only. The bakeoff
  shows that same-bytes/different-chunker cases still need a separate design
  decision before chunked Merkle roots can be stable. pCID-driven chunking,
  chunking-CID / cCID-style descriptors, and negotiated profiles remain
  explicit follow-on pressure, not locked first-profile terms.
- **Promisebase recommendation:** use promisebase / pitbase as prior art only
  for the first L6 CAS spec. Its block/tree/stream and Rabin evidence remains
  relevant, but adopting it as the substrate is premature before the first
  PromiseGrid pointer/raw profile and migration path are decided.

## Next DF packet

This is the next user-answerable decision packet for TODO-kituj / TE-43. It is
not answered here. DF-tumus.1 through DF-tumus.3 have been synthesized from the
standalone bakeoff simulations added under `DI-bukoh`. Source: DI-majib;
DI-bukoh; DI-molah.

- **DF-tumus.1 — Starting profile.** Choose Alt-1A minimal pointer/raw first
  profile with DAG-CBOR / CIDv1 bridge constraints (recommended), Alt-1B
  IPFS/IPLD-aligned default, Alt-1C promisebase-adapter default, or Alt-1D a
  refined alternative.
- **DF-tumus.2 — Type binding.** Choose Alt-2A CID codec as authoritative
  object type, path suffixes rejected, and internal kind allowed only as part of
  the codec-defined payload schema (recommended); Alt-2B strict codec-only with
  no internal kind; Alt-2C codec plus object-internal kind as a two-part type
  rule; or Alt-2D a refined alternative.
- **DF-tumus.3 — Chunking lock scope.** Choose Alt-3A raw-only first profile
  with chunked Merkle roots deferred to a follow-on chunking-identity decision
  (recommended), Alt-3B lock a full chunking algorithm and parameters now,
  Alt-3C use pCID-driven chunking now, Alt-3D use a chunking-CID /
  cCID-style descriptor now, Alt-3E use profile-negotiated chunking now, or
  Alt-3F a refined alternative.
- **DF-tumus.4 — Promisebase stance.** Choose Alt-4A use promisebase / pitbase
  as prior art only for the first L6 CAS spec (recommended), Alt-4B
  adopt-as-substrate-with-adapter, Alt-4C defer until TODO-dozak
  merge-trajectory work, or Alt-4D a refined alternative.

## Acceptance criteria

- TODO-kituj / TE-43 can draft without rereading turn 177.
- TODO-pipus can tell whether CAS-backed migration is blocked on a CAS profile
  or can proceed with a minimal placeholder.
- Promisebase / pitbase evidence is classified as adopted substrate, prior art,
  or deferred merge input.
- Steve can answer `DR-tumus` from this packet without rereading the three
  `DI-bukoh` bakeoff simulations.

## Notes

This DR does not choose the L6 profile. It gives TODO-kituj / TE-43 a concrete decision request to answer before drafting or freezing the first L6 CAS spec.
