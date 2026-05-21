# Question

How should PromiseGrid require scoped claim cards, exact audit-anchor citations, and minimal local promise-accounting records so early apps can make honest partial-conformance claims without conflating:
- local IDs with protocol-boundary identity,
- hash/address with authorization,
- live transport with durable audit publication, or
- request receipt with physical-effect completion?

## Open decision points

- When no frozen `pCID` exists, is `draft-path + revision hash + provisional label` sufficient as the authoritative boundary identity?
- Which claim-card fields are mandatory for every embodiment, and which are plane-specific?
- Should the default rule be that possession of a hash is never authorization unless a separate capability artifact is named?
- For non-idempotent device actions, must ambiguous post-restart state always emit a break-witness instead of retrying execution?
- What is the minimum promise-accounting record set needed for 100-year audit without over-prescribing implementation?
