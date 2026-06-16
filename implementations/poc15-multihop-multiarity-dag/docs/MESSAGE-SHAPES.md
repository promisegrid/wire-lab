# POC15 Message Shapes, Parent Links, COSE, and Raw CAS DAG

POC15 tests that `pCID` really owns arity, slot meaning, signable view,
proof location, and payload interpretation. The outer invariant remains a CBOR
`grid(...)` array whose slot 0 is `42(pCID)`, but slots 1..N are defined by the
protocol spec named by that pCID. Source: `DI-podut`; `DI-mosat`.

## Raw Artifact CAS

Every raw artifact observed during a POC15 run should be retained under the
run-scoped output tree as CAS-addressed bytes or text:

- valid CBOR envelope frames,
- malformed frames before parse failure,
- app/kernel frames,
- peer/kernel frames,
- forwarded carried-message bytes,
- ACK and non-commitment messages,
- raw live-agent decision JSON,
- monitor prompts and responses,
- WASM adapter inputs and outputs,
- stdio worker inputs and outputs,
- analyzer summaries that link back to raw artifacts.

The raw CAS is a local review archive, not a production global store. Clean/reset
removes it unless the user explicitly exports a run.

## Agent Sparse CAS

The raw artifact CAS is also separate from each app's own sparse CAS. In the
current executable slice, apps mirror exact sent/received message bytes into
their local run-scoped CAS metadata, may keep arbitrary local bytes, may retain
encrypted blobs by ciphertext CID, and may keep peer-served content after
storage or bearer-token promises. These app stores are intentionally incomplete:
an app-local message DAG can record that a parent CID is missing without making
the run-level artifact DAG invalid. Source: `DI-manul`.

## Run DAG Versus Wire DAG

POC15 should keep two DAG concepts separate:

- **Run CAS DAG:** local review links among raw artifacts, events, decisions,
  rejects, forwards, ACKs, monitor reports, and adapter I/O.
- **Wire message DAG:** valid PromiseGrid messages that include pCID-defined
  parent links to exact prior envelope CIDs.

Malformed bytes and non-message artifacts can be linked in the run CAS DAG, but
they are not valid PromiseGrid parent-linked messages.

## Exact Message Identity

The parent-link identity for a valid message is the CID of the exact signed
envelope bytes observed on the wire. Payload CIDs and signable-view CIDs may be
recorded as secondary indexes, but parent links for POC15 should point at exact
envelope CIDs so replay and review are byte-faithful.

## Slot-Vector Specimens

POC15 now emits exact raw-message specimens for these pCID-owned slot vectors:

1. **Transport/session-auth-only:** `grid([42(pCID), payload])`.
2. **Common signed message:** `grid([42(pCID), payload, proof])`.
3. **Envelope parents before body:** `grid([42(pCID), parents, payload, proof])`.
4. **Envelope body before parents:** `grid([42(pCID), payload, parents, proof])`.
5. **Payload-owned parents:** `grid([42(pCID), payload, proof])`, where the
   payload includes parent links in a pCID-defined position.
6. **COSE as payload:** `grid([42(pCID), cose_sign1])`.
7. **COSE as proof:** `grid([42(pCID), payload, cose_sign1_detached])`.

The current executable slice covers items 1, 2, 3, 4, 6, and 7 as run-local
operator-review specimens, and it also exercises item 5 in ordinary `route_v1`
traffic by carrying a payload parent field that points at an exact prior message
hash. Normal `route_v1` forwarding can now use envelope-parent slots, while the
route payload can separately carry payload-owned parent links. Parent links use
CBOR tag-42 IPLD links where they appear as envelope slots; payload-owned parent
links remain pCID-defined payload fields in this POC slice. Source: `DI-kohuj`.

## COSE Specimens

POC15 tests COSE in two roles:

- **COSE-as-payload:** slot 1 is a COSE object, such as `COSE_Sign1`, whose
  payload is the promise body under the pCID's rules.
- **COSE-as-proof:** a later slot is a detached COSE proof over the pCID-defined
  signable view.

The current executable slice uses EdDSA/Ed25519, requires COSE `alg` in the
protected header, verifies COSE-as-payload and detached COSE-as-proof specimens,
records a tamper rejection, and has unit tests that reject wrong algorithms,
unprotected-header variants, and mismatched detached payloads. Source:
`DI-kohuj`.

## Transport Proof Versus Message Proof

Transport/session signatures are hop-local promises about a connection or direct
send. Envelope/message proofs are durable object-level promises about the message
bytes that can survive CAS storage, forwarding, replay, and offline review.
POC15 should test transport-auth-only messages, but it must not report them as
self-contained Alice-authored promise objects unless the pCID explicitly defines
that semantics and the corresponding transport/session events are present.
The current POC records this as `transport_proof_comparison_recorded`: a
transport/session signature is useful hop-local context, while the retained
envelope proof remains the durable object-level promise record for offline CAS
review unless a future pCID explicitly chooses otherwise. Source: `DI-kohuj`.

## Analyzer Expectations

The analyzer now reports and gates the first specimen layer:

- counts by pCID specimen and arity,
- counts by proof style: native proof, COSE payload, COSE proof,
  transport-auth-only,
- parent-link artifact counts,
- raw CAS object counts for retained exact message bytes,
- agent-local sparse CAS counts for retained app bytes, encrypted-object CIDs,
  peer storage/retrieval promises, bearer storage-token flow, local GC, and
  sparse message-DAG missing-parent records,
- any accidental return to one universal payload shape.

The analyzer now reconstructs the retained raw-message DAG enough to report root
count, reachable count, maximum depth, missing-parent count, parent-link count,
and parent-link location counts. Missing parents, non-reachable records, and
loss of either envelope or payload parent-link coverage fail the clean
regression. Duplicate observations of the same exact envelope are retained as
artifact rows, but DAG reachability is judged over unique exact-message hashes
so retransmits and sender/receiver observations do not look like missing graph
nodes. Raw artifact counts by every POC artifact kind remain a planned
follow-up gate. Source: `DI-kohuj`.
