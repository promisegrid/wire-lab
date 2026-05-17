# Promise Economy Spectrum Scenarios

These scenarios test turn 179's promise-economy mechanism spectrum. They are
simulation pressure, not a final economics mechanism or a base-protocol decision.
Source: `DI-vabij`.

## Scenario Matrix

| Scenario | Setup | What to test | Decision pressure |
|---|---|---|---|
| Peer-local assessment only | Alice records that Bob kept storage promises and Mallory sent corrupt chunks. | Whether pull/keep/advertise decisions can use local observations without token transfer. | The base protocol should not require a token field if peer-local assessment is enough for some groups. |
| Reciprocal barter promise | Alice sends chunk C to Bob if Bob promises onward-restraint or later storage. | Whether reciprocal promises can be represented without fungible units. | Conditional-release and promise accounting need to work before market mechanics are added. |
| Permissioned capability token | Alice issues Bob a non-transferable capability to fetch a chunk family. | Whether permission and redemption can be scoped without implying exchange rates. | Capability semantics may be useful while still avoiding marketplace behavior. |
| Transferable promise token | Bob transfers an Alice-issued promise token to Carol. | Whether provenance, permission, double-spend prevention, and refusal semantics stay local and auditable. | Transferability creates obligations that the base protocol may need to defer. |
| Floating exchange rates | Alice accepts Bob's storage promises at one rate and Carol's at another. | Whether peers can keep local valuation without a central price oracle. | "Everyone is their own central bank" must not become a hidden global currency. |
| Cryptocurrency-toxicity failure | Mallory creates many identities, pumps token value, then defaults on storage. | Which protocol assumptions amplify speculation, Sybil attacks, or central exchange capture. | Any economic mechanism must be tested against failure modes before becoming protocol-visible. |

## Expected Outputs

- A mechanism-neutrality checklist for future PromiseGrid wire formats and
  protocol specs.
- A list of economic mechanism variants that can be simulated independently.
- A warning list of fields or semantics that would prematurely force a token
  marketplace, global price, or central exchange.
