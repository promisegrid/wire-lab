# AT Protocol repo/CAR interop

## Scenario ID

atproto-repo-car-interop

## Source / Provenance

- Source type: new harness scenario
- Source path: `docs/research/grid-envelope-signature-prior-art-20260522.md`;
  `simulations/SIM-pamap-grid-envelope-signable-view-atproto-like/`
- Source simulation: `SIM-pamap-grid-envelope-signable-view-atproto-like/`
- Source row/title: AT Protocol / Bluesky ecosystem repo, CAR, DAG-CBOR, and
  signable-view interop pressure
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-kizal`;
  `DI-kafot`; `DI-nohir`; `DI-nidop`

## Purpose

Test whether a candidate PromiseGrid design can interoperate with AT Protocol /
Bluesky ecosystem object movement without importing that ecosystem as a
PromiseGrid authority. The pressure is repo/CAR/DAG-CBOR and CID-linked object
compatibility, plus exact-byte and local promise-evidence preservation across
the boundary.

## Setup

Alice has PromiseGrid CAS/feed objects and local evidence records for promises
she observed. Bob promises to mirror or inspect an atproto-style repo/CAR export
that may contain DAG-CBOR objects, CID links, and signed-object patterns. Carol
later audits what crossed the boundary using only the bytes, links, and evidence
available to her sparse local view. Mallory may provide a stale, partial, or
misleading repo snapshot.

No actor treats a Bluesky PDS, relay, appview, DID registry, CAR file, or
DAG-CBOR codec as a PromiseGrid trust authority. Those artifacts may carry useful
bytes or claims, but Alice, Bob, and Carol still make local trust judgments from
their own promises, observations, and make/break history.

## Stimulus

Run the candidate simulation against this source test: Alice exports or
references PromiseGrid CAS/feed material through an atproto-style repo/CAR
boundary, Bob imports or mirrors it, and Carol later attempts to verify what
promise-relevant bytes, links, and evidence survived the crossing.

## Expected Pressure

The simulation should explain which semantics stay in PromiseGrid promises and
which are merely atproto-compatible carriage mechanics. It should preserve exact
bytes, CID/link identity, signable-view or proof evidence when present, and
local promise-accounting records without assuming global completeness or
external identity/trust authority.

## Scenario-Specific Evaluation Questions

- Which PromiseGrid objects can be represented directly as atproto-compatible
  DAG-CBOR/CID-linked objects, and which require a PromiseGrid-specific wrapper
  or pointer object?
- What exact bytes does Bob promise to preserve when importing, mirroring, or
  re-exporting the repo/CAR material?
- How does Carol distinguish "atproto tooling can parse this object" from
  "Carol locally trusts the promise evidence this object carries"?
- What happens when Mallory provides a stale CAR export, omits a linked object,
  rewrites metadata, or presents a valid atproto object that is irrelevant to
  Alice's PromiseGrid promises?
