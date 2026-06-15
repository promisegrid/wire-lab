# POC15 Message Shapes, Parent Links, COSE, and Raw CAS DAG

POC15 should test that `pCID` really owns arity, slot meaning, signable view,
proof location, and payload interpretation. The outer invariant remains a CBOR
`grid(...)` array whose slot 0 is `42(pCID)`, but slots 1..N are defined by the
protocol spec named by that pCID. Source: `DI-podut`.

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

POC15 should add pCIDs for these specimens:

1. **Transport/session-auth-only:** `grid([42(pCID), payload])`.
2. **Common signed message:** `grid([42(pCID), payload, proof])`.
3. **Envelope parents before body:** `grid([42(pCID), parents, payload, proof])`.
4. **Envelope body before parents:** `grid([42(pCID), payload, parents, proof])`.
5. **Payload-owned parents:** `grid([42(pCID), payload, proof])`, where the
   payload includes parent links in a pCID-defined position.
6. **COSE as payload:** `grid([42(pCID), cose_sign1])`.
7. **COSE as proof:** `grid([42(pCID), payload, cose_sign1_detached])`.

Root messages use an empty parent array for specimens with explicit parent
slots. Parent links use CBOR tag-42 IPLD links where present.

## COSE Specimens

POC15 should test COSE in two roles:

- **COSE-as-payload:** slot 1 is a COSE object, such as `COSE_Sign1`, whose
  payload is the promise body under the pCID's rules.
- **COSE-as-proof:** a later slot is a detached COSE proof over the pCID-defined
  signable view.

The first executable slice should use EdDSA/Ed25519 and require COSE `alg` in the
protected header. Analyzer gates should reject wrong algorithms, unprotected
algorithm-only claims, tampered payloads, and mismatched detached signable views.

## Transport Proof Versus Message Proof

Transport/session signatures are hop-local promises about a connection or direct
send. Envelope/message proofs are durable object-level promises about the message
bytes that can survive CAS storage, forwarding, replay, and offline review.
POC15 should test transport-auth-only messages, but it must not report them as
self-contained Alice-authored promise objects unless the pCID explicitly defines
that semantics and the corresponding transport/session events are present.

## Analyzer Expectations

The analyzer should report:

- counts by pCID specimen and arity,
- counts by proof style: native proof, COSE payload, COSE proof,
  transport-auth-only,
- parent-link counts and DAG reconstruction success,
- raw CAS object counts by artifact kind,
- orphaned parent links,
- malformed artifacts retained before parse failure,
- events that reference missing raw artifacts,
- any accidental return to one universal payload shape.
