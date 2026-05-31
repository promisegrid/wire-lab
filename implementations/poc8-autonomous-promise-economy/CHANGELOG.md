# CHANGELOG

```changelog-entry
claim:           implements
spec:            scenario:promise-economy-capability-token-exchange
scope:           executable-proof-of-concept
breaking-change: false
notes:           Adds POC8 as a new successor to POC7. POC8 preserves signed CBOR grid envelopes, one pCID-selected promise-economy protocol, issuer-local token promises, framed TCP, and real storage/compute redemption while replacing Alice's centralized transaction script with autonomous local need advertisements, offer promises, counter promises, accept/refuse decisions, bearer-for-non-transferable exchange, collateral/stake promises, peer-local exchange quotes, and stale-token trust decay. Source: `DI-sirus`.
```

```changelog-entry
claim:           extends
spec:            scenario:promise-economy-capability-token-exchange
scope:           prior-evidence-lineage
breaking-change: false
notes:           Inherits the useful POC7 evidence lineage: capability tokens as issuer promises, bearer/non-transferable token distinction, issuer-local redemption/revocation, local trust updates, exact signed CBOR grid bytes, and Mallory-to-Dave stale-token pressure. POC7 remains preserved as prior evidence. Source: `DI-tugih`; `DI-fibok`; `DI-tanat`; `DI-pabot`; `DI-rodog`; `DI-hanih`; `DI-sirus`.
```
