# Simulations

`simulations/` holds bounded experiment worlds for deriving PromiseGrid design
choices from wire-lab evidence. A simulation may contain candidate protocol
specimens, archived historical inputs, concrete world state, observations, and
results, but it is not itself the PromiseGrid development guide or a final
PromiseGrid node layout.

Root wire-lab apparatus stays rooted unless a later DI says otherwise. In the
current Mupoz split, `protocols/wire-lab.d/` remains the harness apparatus home,
while candidate PromiseGrid protocol trees and concrete specimens move into a
named simulation until results graduate through DR, DI, frozen specs, guide
prose, or a future PromiseGrid spec corpus. Source: `DI-pakid`; `DI-fakin`.

Cross-simulation pressure lives at root `scenarios/`, and cross-simulation run
evidence lives at root `results/`. Simulation-local `SCENARIOS.md` files remain
lineage-local pressure sources; `TODO-dadub` mined those rows into root scenario
entries under `DI-nanih` and added root application-seed entries under
`DI-midif`. Root results use
`results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md` and remain
evidence only. Source: `DI-faros`; `DI-miror`; `DI-vabor`; `DI-dimas`;
`DI-nanih`; `DI-midif`.

## Current simulations

| Simulation | Purpose | Status |
|---|---|---|
| `SIM-ludut-wire-lab-devs/` | Concrete wire-lab-devs world evidence and transport specimen state for replay and migration provenance. | Active concrete-world lineage |
| `SIM-rakot-group-session/` | Independent group-session lineage carrying session-envelope protocol drafts and TODO ownership. | Active protocol lineage |
| `SIM-ludaf-udp-feed/` | Independent UDP feed lineage (renamed from legacy `udp-binding` active-tree naming). | Active protocol lineage |
| `SIM-labit-feed-outer/` | Independent thin outer-feed lineage, including extracted feed-outer draft material. | Active protocol lineage |
| `SIM-kurim-grid-envelope/` | Parent seed and successor-owner lineage for positional grid-envelope variants. | Split parent / successor-owner lineage |
| `SIM-hugoj-cas-usenetlike-gitlike/` | Broad design exploration of a CAS + Usenet-like + git-like PromiseGrid lineage, with `group-session` treated as one current specimen rather than the whole subject. Source: `DI-pijun`. | Active design exploration |
| `SIM-jomag-cas-object-model/` | Turn-177 CAS object-model exploration for deterministic CBOR / DAG-CBOR, pointer objects, chunking parameters, and CIDv1 object typing. Source: `DI-navod`; `DI-tibis`. | Active design-point exploration |
| `SIM-ligan-promisebase-reference-naming/` | Turn-178 standalone exploration for promisebase-style human-readable references to CAS roots without custom non-CID syntax or identity confusion. Source: `DI-tibis`. | Active design-point exploration |
| `SIM-bobud-l6-cas-starting-profile-bakeoff/` | DR-tumus bakeoff for the first L6 CAS starting profile: IPFS / IPLD-aligned, promisebase-adapter, minimal pointer/raw, or a refined alternative. Source: `DI-bukoh`. | Active design-point exploration |
| `SIM-kohad-cas-object-type-binding-bakeoff/` | DR-tumus bakeoff for CAS object type binding: CID codec-only, codec plus internal kind, path suffix negative control, or a refined alternative. Source: `DI-bukoh`. | Active design-point exploration |
| `SIM-gobaz-chunking-identity-bakeoff/` | DR-tumus bakeoff for chunking identity: pCID-driven chunking, chunking-CID / cCID-style descriptor, negotiated profile, raw-only deferral, or a refined alternative. Source: `DI-bukoh`. | Active design-point exploration |
| `SIM-zazit-chunk-feed-replication/` | Turn-177 L5 feed exploration where feeds advertise, request, and replicate CAS chunks between sparse sites rather than carrying group messages. Source: `DI-navod`. | Active design-point exploration |
| `SIM-jurar-cas-backed-group-session/` | Turn-177 successor exploration for group-session semantics over CAS roots and pointer objects, without rewriting the existing `.txt` specimen. Source: `DI-navod`. | Active design-point exploration |
| `SIM-rusap-promise-accounting-records/` | Turn-177 promise-economy exploration for peer-local promise accounting records that inform pull, keep, advertise, and refusal decisions. Source: `DI-navod`. | Active design-point exploration |
| `SIM-punaz-bgp-class-routing-app/` | Turn-178 standalone exploration for BGP-class routing-policy applications over peer-local promises without central route authority or global trust registry. Source: `DI-tibis`. | Active design-point exploration |
| `SIM-haros-promise-economy-spectrum/` | Turn-179 standalone exploration for promise-economy mechanism neutrality across peer-local assessment, capability tokens, transferability, floating rates, and cryptocurrency-toxicity failure modes. TODO owner: TODO-rajig. Source: `DI-vabij`; `DI-pidag`. | Active design-point exploration |
| `SIM-zarud-conditional-release-geofencing/` | Question home for conditional-release, onward-restraint, geofencing, and recursive promise-graph ownership from TODO-ralud. Source: `DI-pukap`. | Active design-point exploration |
| `SIM-narok-transport-family-bakeoff/` | Question home for future ring, cluster-of-clusters, gossip, and receipts-at-scale transport-family pressure from TODO-sinuv. Source: `DI-pukap`. | Active design-point exploration |
| `SIM-dihiz-peer-adoption-metadata/` | Question home for peer-level pCID adoption metadata and open-question answer bindings from TODO-nivus. Source: `DI-pukap`. | Active design-point exploration |
| `SIM-ranib-spec-requirement-sections/` | Question home for protocol spec promise-vocabulary, 100-year, and layperson/easy-implementation section requirements from TODO-kulih / DR-robon. Source: `DI-pukap`. | Active design-point exploration |
| `SIM-bohof-group-session-freeze-promise/` | Question home for group-session freeze evidence and `merge-group-transport-spec` promise shape from TODO-bisur. Source: `DI-pukap`. | Active design-point exploration |
| `SIM-kuful-udp-feed-v0-conformance/` | Question home for UDP-feed v0 reference implementation, test-vector, artifact-writer, and ns-3 conformance evidence from TODO-jodon. Source: `DI-pukap`. | Active design-point exploration |

## Grid-envelope positional variant matrix

`DI-fanah` split the pending grid-envelope successor work from
`SIM-kurim-grid-envelope/` into standalone positional specimens. These variants
compete independently across encoding, unknown-pCID handling, and signature
placement; none is preferred by the harness or by this index.

| Simulation | Encoding | Unknown-pCID policy | Signature policy |
|---|---|---|---|
| `SIM-mahih-grid-envelope-enc-cbor-unknown-opaque-sig-wrapper-pcid/` | CBOR | Opaque store/forward | Wrapper pCID |
| `SIM-gasus-grid-envelope-enc-cbor-unknown-opaque-sig-unsigned-v0/` | CBOR | Opaque store/forward | Unsigned v0 |
| `SIM-vutar-grid-envelope-enc-cbor-unknown-opaque-sig-mandatory-opaque-bytes/` | CBOR | Opaque store/forward | Mandatory opaque bytes |
| `SIM-vamaz-grid-envelope-enc-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload/` | CBOR | Opaque store/forward | Mandatory signature pCID + payload |
| `SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid/` | CBOR | Hard reject | Wrapper pCID |
| `SIM-gazan-grid-envelope-enc-cbor-unknown-hard-reject-sig-unsigned-v0/` | CBOR | Hard reject | Unsigned v0 |
| `SIM-hupir-grid-envelope-enc-cbor-unknown-hard-reject-sig-mandatory-opaque-bytes/` | CBOR | Hard reject | Mandatory opaque bytes |
| `SIM-kovis-grid-envelope-enc-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload/` | CBOR | Hard reject | Mandatory signature pCID + payload |
| `SIM-vivus-grid-envelope-enc-cbor-unknown-best-effort-sig-wrapper-pcid/` | CBOR | Best-effort inspection | Wrapper pCID |
| `SIM-fonig-grid-envelope-enc-cbor-unknown-best-effort-sig-unsigned-v0/` | CBOR | Best-effort inspection | Unsigned v0 |
| `SIM-guhor-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/` | CBOR | Best-effort inspection | Mandatory opaque bytes |
| `SIM-jokak-grid-envelope-enc-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/` | CBOR | Best-effort inspection | Mandatory signature pCID + payload |
| `SIM-gojot-grid-envelope-enc-dag-cbor-unknown-opaque-sig-wrapper-pcid/` | DAG-CBOR | Opaque store/forward | Wrapper pCID |
| `SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0/` | DAG-CBOR | Opaque store/forward | Unsigned v0 |
| `SIM-ruzil-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-opaque-bytes/` | DAG-CBOR | Opaque store/forward | Mandatory opaque bytes |
| `SIM-riliz-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload/` | DAG-CBOR | Opaque store/forward | Mandatory signature pCID + payload |
| `SIM-nipoh-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-wrapper-pcid/` | DAG-CBOR | Hard reject | Wrapper pCID |
| `SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0/` | DAG-CBOR | Hard reject | Unsigned v0 |
| `SIM-tohol-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-opaque-bytes/` | DAG-CBOR | Hard reject | Mandatory opaque bytes |
| `SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload/` | DAG-CBOR | Hard reject | Mandatory signature pCID + payload |
| `SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/` | DAG-CBOR | Best-effort inspection | Wrapper pCID |
| `SIM-zifik-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-unsigned-v0/` | DAG-CBOR | Best-effort inspection | Unsigned v0 |
| `SIM-fonol-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-mandatory-opaque-bytes/` | DAG-CBOR | Best-effort inspection | Mandatory opaque bytes |
| `SIM-rakir-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-mandatory-sig-pcid-payload/` | DAG-CBOR | Best-effort inspection | Mandatory signature pCID + payload |
