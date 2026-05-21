# Grid-envelope draft: signed summary header with nested schema payload

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `signed-summary-header-nested-schema`.

## Scope

This spec defines one grid-envelope candidate for wire-lab comparison. It is a specimen inside `SIM-maraz-grid-envelope-signed-summary-header-nested-schema`, not a harness rule and not the canonical PromiseGrid envelope. It was promoted from proposal `SIM-maraz-child-signed-summary-header-nested-schema` under `DI-fihub`; the scored proposal tree remains raw evidence under `proposals/`.

This specimen combines two parent moves:

- From `SIM-janov-grid-envelope-layer-pcid-nested-signed-payload`: keep a nested signed payload so the actual payload remains explicitly bound to its own `pcid`.
- From `SIM-sajar-grid-envelope-variable-arity-pcid-defined-fields`: allow the inner payload schema to remain pCID-defined and flexible, including variable-arity designs, instead of forcing one universal application field layout.

The main repair is a small universal outer header plus an outer signature, so unknown inner schemas do not collapse into unauditable opaque bytes.

## Envelope Shape

The outer envelope shape for this variant is:

```text
[pcid_a, header_a, payload_a, sig_a]
```

Slots are interpreted as follows:

- `pcid_a`: identifies this shared outer envelope variant.
- `header_a`: deterministic CBOR map with universal audit-visible fields.
- `payload_a`: canonical bytes of a nested signed payload.
- `sig_a`: outer signature over the canonical bytes of `[pcid_a, header_a, payload_a]`.

For this candidate, `pcid_a` defines `payload_a` as the canonical bytes of:

```text
[pcid_b, payload_b, sig_b]
```

Nested slots are interpreted as follows:

- `pcid_b`: identifies the actual payload schema.
- `payload_b`: payload bytes interpreted only by the `pcid_b` schema.
- `sig_b`: nested signature over the canonical bytes of `[pcid_b, payload_b]`.

`pcid_b` MUST equal `header_a.content_pcid`.

## Universal Header Fields

`header_a` is a deterministic CBOR map with these required keys:

- `content_pcid`: CIDv1 for the inner payload schema. Must match `pcid_b`.
- `msg_class`: one of `advertisement`, `conditional-release`, `route-claim`, `withdrawal`, or `other`.
- `summary_refs`: array of zero to sixteen CIDv1 references to audit-visible CAS objects.
- `signer_ref`: reference to the outer signer identity or key material used to verify `sig_a`.

Optional keys:

- `related_ref`: CIDv1 of related, prior, superseded, disputed, or parent evidence.
- `freshness`: map containing either or both of `issuer_seq` and `not_after_unix_s`.
- `schema_ref`: CIDv1 for an immutable schema bundle, handler bundle, or documentation object that helps future peers recover `content_pcid` semantics.

Generic tooling may index only these universal fields. It MUST NOT infer full application semantics from them.

## Why `summary_refs` Exists

The parent variable-arity design made too much evidence invisible to generic peers. This specimen requires every sender to project a minimal stable summary into `summary_refs` so routing, replication, and later audit can still refer to durable objects even when `pcid_b` is unknown.

Examples:

- Sparse replication advertisement: `summary_refs` may include a Merkle root CID, a frontier summary CID, and a pointer-summary CID.
- Conditional release: `summary_refs` may include a restraint-graph root CID and a recipient-acceptance CID.
- BGP-like routing: `summary_refs` may include a route-set root CID, path-evidence CID, or withdrawal object CID.

The header does not claim that these references are sufficient to fully understand the message. It only makes contested evidence discoverable and indexable without hidden global state.

## Encoding and Canonicalization

- The outer envelope is a deterministic CBOR positional array of arity 4.
- `header_a` is a deterministic CBOR map.
- `payload_a` is a byte string containing the deterministic CBOR bytes of `[pcid_b, payload_b, sig_b]`.
- `pcid_a` and `pcid_b` are CIDv1 byte strings.
- `sig_a` and `sig_b` are signature objects encoded as byte strings or deterministic CBOR values defined by the signing suite.
- The canonical bytes for `sig_a` are the deterministic CBOR bytes of `[pcid_a, header_a, payload_a]`.
- The canonical bytes for `sig_b` are the deterministic CBOR bytes of `[pcid_b, payload_b]`.

This preserves the Janov strength that the actual payload bytes are signed together with their payload `pcid`.

## Unknown pCID Policy

If a receiver lacks a handler for `pcid_a`, it may preserve or forward the full outer envelope as opaque evidence, but it MUST NOT claim to parse `header_a`, verify either signature, or expose application meaning.

If a receiver understands `pcid_a` but lacks a handler for `content_pcid` / `pcid_b`, it:

- MAY parse `header_a`.
- MAY verify `sig_a` using `signer_ref` and the signature suite.
- MAY retain `payload_a` and `schema_ref` for later recovery.
- MAY verify `sig_b` only if the signature object format is known without needing the `pcid_b` payload semantics.
- MUST NOT claim to understand `payload_b` semantics.
- MUST NOT invent missing fields beyond `header_a`.

This is stricter than generic variable-arity opacity but more useful than total opacity because summary references, signer, related evidence, and freshness remain auditable.

## Signature and Authorship Policy

This specimen separates two claims:

- `sig_b` authenticates the actual inner payload claim.
- `sig_a` authenticates the outer sender or relay endorsement of carriage, header projection, and stated relationship to other evidence.

The same actor may produce both signatures, but the design does not require that.

This repairs the main Janov weakness: a valid inner signature alone no longer leaves the outer conformance promise dependent on transport identity.

## Generic Audit Behavior

A generic peer that knows only `pcid_a` can still do the following:

- verify outer endorsement with `sig_a`;
- index by `msg_class`;
- track which immutable objects are referenced in `summary_refs`;
- compare `related_ref` chains for supersession, dispute, or parent-child linkage;
- apply freshness checks using `issuer_seq` or `not_after_unix_s`;
- preserve `schema_ref` so later auditors can recover the missing inner schema.

This is intended to score better under sparse advertisement, onward-restraint chaining, and contested routing evidence than either parent alone.

## Limits and Denial-of-Service Controls

This variant fixes several limits that the parent variable-arity draft left open:

- outer envelope arity is exactly 4;
- nested payload arity is exactly 3;
- total encoded outer envelope size MUST NOT exceed 65536 bytes;
- encoded `header_a` size MUST NOT exceed 2048 bytes;
- `summary_refs` length MUST NOT exceed 16 entries;
- duplicate map keys and non-canonical CBOR MUST be rejected.

Inner `payload_b` may still have variable structure under `pcid_b`, but those inner rules operate inside bounded outer transport and audit scaffolding.

## Scenario Pressure Notes

### Sparse advertisement

This variant does not define a chunk replication protocol, but it requires generic advertisement-visible references. A peer can at least see which root or summary objects are being advertised, who endorsed the advertisement, and whether it is fresh, even when it cannot parse the full replication schema.

### Onward-restraint chain

This variant does not declare where all restraint semantics must live, but it allows a sender to expose a restraint-graph root and acceptance records as `summary_refs`, while the full recursive graph remains inside `payload_b`. `related_ref` can link release evidence to prior promises.

### BGP-like routing

This variant does not define routing policy, but it makes it easier to distinguish origin payload authorship from propagation endorsement. Route claims, withdrawals, and related contested evidence can be chained through `msg_class`, `related_ref`, `summary_refs`, and `freshness`.

## Open Questions

- Is the fixed `msg_class` set small enough to stay durable but large enough to be useful?
- Should `schema_ref` become mandatory for long-term audit of old envelopes?
- Is header projection sufficient to prevent audit collapse for unknown inner schemas, or do some workloads need one more universal field?

## Non-Canonical Status

This draft does not declare a winning envelope, does not define a central pCID registry, and does not constrain sibling simulations. It exists to compare a signed summary header plus nested schema payload against both unsigned nested payload and fully schema-local outer-variable-arity alternatives.
