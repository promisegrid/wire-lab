# CHANGELOG

```changelog-entry
claim:           implements
spec:            scenario:promise-economy-capability-token-exchange
scope:           executable-proof-of-concept
breaking-change: false
notes:           Adds POC9 as a sibling successor to POC8. POC9 preserves signed CBOR grid envelopes, one pCID-selected discovery/economy protocol, issuer-local token promises, framed TCP, and real storage/compute redemption while replacing POC8's known-peer/static-route pressure with deterministic sparse-mesh discovery, route promises, referral promises, ordinary low-risk promises before private escalation, malformed TCP evidence, local route refusal, and Mallory-to-Dave expired-token misuse evidence through an alternate route. Source: `DI-sipuz`; `DI-vujil`.
```

```changelog-entry
claim:           extends
spec:            scenario:promise-economy-capability-token-exchange
scope:           prior-evidence-lineage
breaking-change: false
notes:           Inherits the useful POC7 and POC8 evidence lineage: capability tokens as issuer promises, bearer/non-transferable token distinction, issuer-local redemption/revocation, local trust updates, exact signed CBOR grid bytes, autonomous need/offer/acceptance vocabulary, and Mallory-to-Dave token-pressure evidence. POC9 corrects the POC8-style stale-token wording by treating signed expiry as neutral for Alice and negative only for Mallory when Mallory presents expired bytes as useful. POC7 and POC8 remain preserved as prior evidence. Source: `DI-tugih`; `DI-fibok`; `DI-tanat`; `DI-pabot`; `DI-rodog`; `DI-hanih`; `DI-sirus`; `DI-sipuz`; `DI-vujil`.
```
