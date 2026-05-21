# Grid Envelope Variant Spec (DRAFT)

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `enc-dag-cbor` / `unknown-quarantine` / `sig-mandatory-sig-pcid-payload`.
> Bred from two simulation-local parent specimens. Source lineage: `DI-fanah` plus GA child synthesis promoted under `DI-dipid`.

## Purpose

This spec defines one full positional grid-envelope candidate for wire-lab
comparison. It is a specimen inside `SIM-bimos-grid-envelope-quarantine-sig-pcid-audit-tuple`, not a harness rule
and not the canonical PromiseGrid wire format.

The design goal is to keep the stricter parent's explicit signature dispatch and
hard validity boundary, while importing the other parent's better mixed-version
evolution behavior by preserving rejected unknown envelopes as exact opaque
evidence.

## Positional Envelope Shape

The envelope shape for this variant is:

```text
[pcid, payload, sig_pcid, sig_payload]
```

Slots are interpreted positionally:

- `pcid` identifies the protocol/spec/handler that interprets `payload`.
- `payload` is opaque bytes until interpreted by the handler named by `pcid`.
- `sig_pcid` identifies the signature or proof protocol.
- `sig_payload` is opaque bytes for the handler named by `sig_pcid`.

A `payload` may itself be the canonical bytes of another grid envelope when the
protocol named by `pcid` specifies recursive nesting. The outer grid envelope
does not prescribe the payload's internal organization beyond the bytes boundary.

## Encoding

This variant encodes the envelope as DAG-CBOR-compatible positional arrays.
`pcid` and `sig_pcid` are DAG-CBOR Link values. `payload` and `sig_payload` are
byte strings. The envelope remains positional: no map/object envelope fields are
introduced.

The canonical bytes for hashing, archival, signing input, and evidence relay are
the DAG-CBOR bytes of the exact positional array under this spec.

## Unknown Protocol Policy: `unknown-quarantine`

If a receiver lacks a handler for `pcid` or `sig_pcid`, the envelope enters
**quarantine**.

Quarantine means all of the following:

- The receiver MUST NOT accept the envelope as a valid interpreted message under
  this variant.
- The receiver MUST NOT parse `payload` speculatively without the handler named
  by `pcid`.
- The receiver MUST NOT treat `sig_payload` as verified without the handler named
  by `sig_pcid`.
- The receiver MUST preserve the exact canonical bytes if it stores the object.
- The receiver MAY archive the bytes together with local rejection metadata.
- The receiver MAY relay the exact bytes onward only as **quarantined opaque
  evidence**, and MUST preserve the fact that local interpretation failed.

Relay of quarantined evidence is not acceptance. A relaying peer is forwarding
raw evidence for later interpretation by a different peer, version, or archive.

This policy is intended to avoid the main failure modes seen in the parents:

- unlike `unknown-hard-reject`, useful future evidence is not forced to vanish;
- unlike bare `unknown-opaque`, unknown content is not mistaken for locally valid
  understood traffic.

## Signature and Authorship Policy

The third and fourth positional slots are mandatory.

- `sig_pcid` identifies the signature or proof protocol.
- `sig_payload` is opaque bytes interpreted by that signature protocol.

Unless `sig_pcid` publishes stricter rules, the signature or proof covers the
canonical unsigned prefix:

```text
[pcid, payload]
```

Signing `pcid` together with `payload` reduces type-confusion risk and makes the
dispatch rule part of the auditable signed evidence.

This envelope layer enforces presence, position, and byte-shape. Signer
identity, authority, rotation, revocation, delegation, and freshness policy are
still defined by the protocol ecosystems named by `pcid` and `sig_pcid`.

## Local Audit Tuple

To improve long-term auditability without adding more wire-level fields,
implementations SHOULD retain a local audit tuple for each received envelope:

- exact canonical envelope bytes;
- observed receive time and source context;
- extracted `pcid` and `sig_pcid` bytes;
- local interpretation result: accepted, rejected, or quarantined;
- local verification result and verifier version, if verification was possible;
- any locally known content-addressed retrieval path for the specs named by
  `pcid` and `sig_pcid`.

This tuple is local evidence, not a global registry requirement and not a new
consensus rule.

## Layering-Test Behavior

This variant answers the layering pressure as follows:

- Ordering disagreements are handled by the protocol named by `pcid`; the grid
  envelope preserves only dispatch identity and exact bytes.
- Forwarding or relay evidence can be represented by wrapper envelopes, by the
  payload protocol, or by quarantine relay of exact rejected bytes.
- External or content-addressed body references live inside `payload` under the
  protocol named by `pcid`.
- Incompatible interpretation rules fail visibly at the `pcid` or `sig_pcid`
  boundary, but evidence can still survive for later audit or upgraded peers.

## Scenario Pressure Notes

### IoT fleet maintenance

This variant does not define device identity, maintenance history, firmware
approval, telemetry, or access-control objects by itself. It does improve the
envelope substrate for long-lived mixed fleets by allowing unknown future update
or telemetry envelopes to be retained and escalated instead of dropped, while
still preventing accidental acceptance.

### BGP routing

This variant does not define route, withdrawal, freshness, or leak-policy
semantics by itself. It does improve survivability of unfamiliar route evidence
across sparse peers: newer route objects can be quarantined and relayed for
later inspection rather than black-holed.

### CAS application object families

The first binding rule remains extensible because `pcid` identifies the payload
family without requiring reinterpretation of old bytes. Unknown future families
can be preserved as quarantined evidence, which softens rollout brittleness while
keeping local validity boundaries explicit.

## Non-Goals

This draft does not:

- declare a winning envelope;
- define a central pCID registry;
- freeze a final PromiseGrid signing scheme;
- define application payloads for IoT, routing, or CAS families;
- claim that quarantine relay settles trust or freshness.

## Freeze Gate

This draft can freeze only after at least one simulation run compares it against
sibling positional grid-envelope variants and a maintainer signs a merge/freeze
promise for this specific specimen.
