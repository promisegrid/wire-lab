# Grid-envelope draft: nested payload with outer attestation Multisig

> **Status: DRAFT.** Not frozen. The pCIDs for this spec are not yet minted.
> Variant: `nested-payload-outer-attestation-multisig`.

## Scope

This spec defines one grid-envelope candidate for wire-lab comparison. It is a specimen inside `SIM-natim-grid-envelope-nested-payload-outer-attestation-multisig`, not a harness rule and not the canonical PromiseGrid envelope. It was promoted from proposal `SIM-natim-child-nested-payload-outer-attestation-multisig` under `DI-fihub`; the scored proposal tree remains raw evidence under `proposals/`.

This specimen combines two parent strengths while making one explicit repair:

- keep a nested payload whose actual payload bytes are signed together with their payload pCID;
- keep exact-byte proof retention and conservative unsupported/quarantine behavior;
- repair the unsigned-outer-layer weakness by making a current-actor outer attestation mandatory.

## Bounded design deltas

1. Replace the parent's unsigned outer conformance promise with a mandatory outer attestation slot.
2. Standardize a minimal outer attestation transcript with role, constraints reference, and freshness reference.
3. Require dual-proof audit outcomes so payload authorship and relay/forwarder authorship are evaluated separately.

## Envelope shape

The outer envelope shape is:

```text
[pcid_a, payload_a, attest_pcid, attest_multisig]
```

Slots are interpreted positionally:

- `pcid_a` identifies the layer profile that defines this outer shape and the attestation transcript rules;
- `payload_a` is the exact bytes of a nested signed payload object;
- `attest_pcid` identifies the outer attestation proof profile;
- `attest_multisig` is an exact Cryptid-style Multisig object used for the outer attestation.

For this variant, `payload_a` is the deterministic CBOR bytes of:

```text
[pcid_b, payload_b, payload_sig_pcid, payload_proof]
```

Nested slots are interpreted as follows:

- `pcid_b` identifies the actual payload protocol;
- `payload_b` is the actual application or evidence payload bytes;
- `payload_sig_pcid` identifies the payload-proof profile;
- `payload_proof` is the exact proof bytes that authenticate `[pcid_b, payload_b]`.

This specimen therefore separates two promises:

- **payload authorship/content promise:** authenticated by `payload_proof` over `[pcid_b, payload_b]`;
- **current actor role promise:** authenticated by `attest_multisig` over the outer attestation transcript.

## Encoding

The outer envelope and nested payload object are deterministic CBOR positional arrays. `pcid_a`, `pcid_b`, `attest_pcid`, and `payload_sig_pcid` are CIDv1 byte strings or links as defined by the local run profile. `payload_a`, `payload_b`, `payload_proof`, and `attest_multisig` are byte strings at the carrier layer.

Canonical outer envelope bytes are the deterministic CBOR bytes of `[pcid_a, payload_a, attest_pcid, attest_multisig]`.

Canonical nested payload bytes are the deterministic CBOR bytes of `[pcid_b, payload_b, payload_sig_pcid, payload_proof]`.

Canonical payload-authorship signed bytes are the deterministic CBOR bytes of `[pcid_b, payload_b]`.

## Outer attestation transcript

The outer Multisig object is detached by default in this specimen. The signed message bytes are the deterministic CBOR bytes of:

```text
["pg-relay-attest-v1", pcid_a, payload_a, role_code, constraints_ref, freshness_ref]
```

Fields:

- `role_code`: one of `origin`, `holder`, `relay`, `forwarder`, or `observer`;
- `constraints_ref`: CID or `null`, pointing to an application-specific object that constrains distribution, export, onward release, availability scope, or similar policy;
- `freshness_ref`: CID or `null`, pointing to an application-specific object that records epoch, sequence, observation time, expiry window, or other freshness evidence.

The outer attestation does **not** replace `payload_b`. It says: a particular actor, in a particular role, attests to carrying or forwarding exactly this `payload_a` under the referenced constraints and freshness context.

## Verification order

A verifier evaluates this specimen in two layers.

### 1. Nested payload verification

If `payload_sig_pcid` is understood, verify `payload_proof` against the canonical bytes of `[pcid_b, payload_b]`.

If `payload_sig_pcid` is not understood, retain the exact nested bytes and mark payload verification unsupported.

### 2. Outer attestation verification

If `attest_pcid` and the Multisig signing codec are understood, verify `attest_multisig` against the canonical outer attestation transcript.

If `attest_pcid` or the signing codec is not understood, retain the exact proof bytes and mark outer attestation unsupported.

A receiver must not collapse these two checks into one result.

## Decision classes

This specimen introduces explicit audit outcomes:

- `accepted-dual-bound`: both payload proof and outer attestation verify, and local policy accepts the referenced role, constraints, and freshness;
- `payload-authenticated-only`: payload proof verifies but outer attestation is missing, unsupported, stale, or invalid;
- `relay-attested-opaque`: outer attestation verifies but payload proof or payload semantics are unsupported;
- `unsupported`: exact bytes retained, but required codecs, pCIDs, or attributes are not understood;
- `quarantine`: structurally parseable evidence exists but local policy detects mismatch, stale freshness, forbidden role, or critical unknown attributes;
- `reject`: bytes or signatures are malformed, inconsistent, or cryptographically invalid.

These decision classes are part of the specimen and must be locally recordable.

## Unknown pCID and unknown codec policy

This specimen adopts a conservative rule set:

- if `pcid_a` is unknown, preserve exact outer bytes as opaque evidence and do not claim the attestation transcript shape;
- if `pcid_b` is unknown, preserve exact nested bytes and do not claim payload interpretation;
- if `payload_sig_pcid` is unknown, do not claim payload verification even if the payload bytes are stored;
- if `attest_pcid` or the Multisig signing codec is unknown, do not claim current-actor attestation even if the proof bytes are stored;
- unknown critical attributes in either proof profile force `unsupported` or `quarantine`, never success.

Structural skippability is not validity.

## Scenario pressure notes

### Sparse chunk-feed replication

Use `payload_b` for the sparse advertisement object itself. Use `role_code = holder` when the current actor claims partial possession or service scope, and `role_code = relay` when merely forwarding another actor's advertisement. `constraints_ref` may point to an availability-scope object, and `freshness_ref` may point to a frontier or observation-window object.

This does not define the chunk schema, but it does define how to distinguish signed content from current availability/relay promises.

### Conditional-release onward-restraint chain

Use `payload_b` for the release object or restraint graph. A forwarding actor uses `role_code = forwarder` and binds `constraints_ref` to the onward-restraint policy object or recipient-acceptance object. A receiver can therefore audit whether Bob merely relayed content or attested to forwarding under a specific restraint reference.

This does not settle whether the full graph belongs at session, transport, or CAS-object level, but it provides an explicit binding point.

### BGP-style routing pressure

Use `payload_b` for route advertisement, withdrawal, or observation objects. `role_code = origin` and `role_code = relay` distinguish origin claims from propagation claims. `constraints_ref` may point to export-policy or path-policy objects; `freshness_ref` may point to observation or expiry objects. This helps prevent a valid payload proof from being mistaken for a current propagation endorsement.

## Audit retention requirements

A local audit record for this specimen must retain:

- exact outer envelope bytes;
- exact `payload_a` bytes;
- exact payload-authorship signed bytes `[pcid_b, payload_b]`;
- exact `payload_proof` bytes;
- exact outer attestation transcript bytes;
- exact `attest_multisig` bytes;
- verification profile identifiers used for both layers;
- decision class and reason.

## Benefits under parent weaknesses

This specimen directly repairs the shared parent weakness where a valid inner proof could exist but the current relay or forwarder promise remained tied to transport identity or informal context. By forcing a separately auditable outer attestation, the design better distinguishes origin, holder, relay, and forwarder claims while keeping nested payload evolution and exact-byte auditability.

## Non-goals

This draft does not:

- declare a final PromiseGrid envelope;
- freeze any pCID;
- define final chunk-feed, onward-restraint, or routing payload schemas;
- define global key discovery or revocation;
- require a central authority;
- require Cryptid Multisig for inner payload proofs.

## Open questions

- Should `constraints_ref` and `freshness_ref` remain generic CID hooks, or should a future profile require more structured outer metadata?
- Is `role_code` sufficient for the sampled scenarios, or do some domains require finer distinctions such as `aggregator` or `withdrawer`?
- Should the outer attestation always be detached, or may a future frozen profile allow a combined-message Multisig form if it preserves the same exact transcript binding?

## Freeze gate

This draft can freeze only after comparison runs show that mandatory outer attestation improves auditability and failure handling over unsigned-outer nested-payload specimens, and after at least one profile fully specifies acceptable role codes, codec handling, and stale-evidence policy.
