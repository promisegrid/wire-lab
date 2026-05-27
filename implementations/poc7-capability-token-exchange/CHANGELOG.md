# CHANGELOG

```changelog-entry
claim:           implements
spec:            scenario:promise-economy-capability-token-exchange
scope:           executable-proof-of-concept
breaking-change: false
notes:           Adds a bounded five-container POC for capability-token access promises, bearer/non-transferable token exchange, issuer-local redemption/revocation, peer-local exchange offers, and local trust updates. Source: `DI-tugih`.
```

```changelog-entry
claim:           extends
spec:            scenario:promise-economy-capability-token-exchange
scope:           executable-proof-of-concept
breaking-change: false
notes:           Deepens POC7 with signed CBOR grid envelopes, CBOR token bytes, real storage/compute/data redemption payloads, holder-initiated trade, local exchange-state mutation, and a fix for the Carol access-token transaction. Source: `DI-fibok`.
```

```changelog-entry
claim:           extends
spec:            scenario:promise-economy-capability-token-exchange
scope:           executable-proof-of-concept
breaking-change: false
notes:           Replaces HTTP wrapper transport with length-framed TCP and reframes app message kinds around resource promise requests, promise fulfillment presentations, promise receipts, reciprocal exchange promises, local evidence observations, and issuer-local revocation notices. Source: `DI-tanat`.
```

```changelog-entry
claim:           extends
spec:            scenario:promise-economy-capability-token-exchange
scope:           executable-proof-of-concept
breaking-change: false
notes:           Removes the special auditor role and Alice-to-Dave token shortcut; Mallory now voluntarily circulates the stale token to Dave, and Dave redeems it through ordinary trader behavior before updating local trust in Alice and Mallory. Source: `DI-pabot`.
```

```changelog-entry
claim:           extends
spec:            scenario:promise-economy-capability-token-exchange
scope:           executable-proof-of-concept
breaking-change: false
notes:           Adds deterministic per-agent local policy scoring before issue, accept, redeem, transfer, quote, and reciprocal exchange actions; Dave now refuses a later Mallory stale-token transfer after Dave's own broken-redemption evidence changes local trust. Source: `DI-rodog`.
```
