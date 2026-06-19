# Protocol Spec Documents

`docs/protocols/` stores protocol specification documents whose bytes can be
hashed into PromiseGrid protocol CIDs. These docs are not a central authority or
registry; each pCID remains the content-derived name of one specific protocol
spec, and each agent locally decides which pCIDs it promises to use.

Editable specs use human-readable slug filenames while they are being drafted,
for example `route-v1.md`. When a spec is frozen for an executable POC, add a
symlink whose filename is the POC's CIDv1 text form for the exact spec bytes:

```text
docs/protocols/<cidv1-text>.md -> <slug>.md
```

Intent: POC16 needs real protocol specs for pCID provenance, while keeping
runtime routing free of prose-spec lookups during ordinary message handling.
Source: DI-mubul; DI-nogij

The POC16 transport listener should treat slot 0 as a protocol-family selector,
not as an app address, message kind, RPC method, or service-registry entry. The
pCID-selected parser or builder role owns pCID-specific payload decoding,
including whether the payload is a CBOR array, CBOR map, COSE object, encrypted
bytes, nested selector, app-local address, route promise, or another shape
defined by the spec. TE-ritig keeps the exact POC16 parser/builder split in
needs-DF state until implementation decisions are locked.

POC16 currently uses `cidv1-raw-sha2-256:<hex>.md` symlink names because its
minimal CID implementation renders CIDv1 raw/sha2-256 bytes in that stable text
form. LLM-backed agents should receive the exact relevant spec prose through
`go:embed` prompt-context construction for every pCID they promise to send,
receive, redeem, verify, store, compute, or route. Runs should record which spec
CIDs were supplied to each LLM agent so later review can detect prompt-context
drift.
