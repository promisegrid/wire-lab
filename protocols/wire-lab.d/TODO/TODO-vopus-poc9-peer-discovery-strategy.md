# TODO-vopus: POC9 peer discovery strategy

## Status

Implemented. Owns `implementations/poc9-peer-discovery-strategy/`, the
proof-of-concept pressure after POC8: agents discover peers, offers, routes, and
ordinary low-risk promises through local relationships and signed promise
evidence rather than through fixed role scripts, static peer knowledge, or a
central directory. Source: `DI-vorus`; `DI-sipuz`; `DI-vujil`.

## Scope

- Treat POC9 as executable design evidence, not a final PromiseGrid discovery,
  trust, routing, transport, token, storage, compute, or economics API.
- Preserve the POC8 one-pCID lesson: message kinds are payload variants under
  one protocol pCID, not separate pCIDs per message type.
- Use signed CBOR `grid([42(pCID), payload, proof])` envelopes over framed TCP
  unless a later DF explicitly chooses a different implementation shape.
- Start discovery from existing contacts: configured seed peers, prior trusted
  peers, local or physical pairing, imported references, and introductions from
  peers already known locally.
- Model referrals as evidence, not authority. Bob's statement that Carol offers
  storage is Bob's promise unless Bob carries Carol's signed offer; Alice still
  makes her own local trust judgment either way.
- Keep peering scoped. Alice may trust Bob to gossip public offers without
  trusting Bob to store private data, return computation results, relay sensitive
  messages, or price bearer tokens.
- Use ordinary low-risk storage, compute, relay, or exchange promises to build
  relationship evidence before higher-risk storage, compute, relay, or
  bearer-token exchange. Do not add a special `probe_promise` payload kind.
- Treat transport behavior as local evidence. Promise relationships influence
  whether an agent opens, keeps, retries, or uses TCP; TCP delivery, malformed
  bytes, disconnects, timeouts, retries, and latency feed local
  promise-accounting records.
- Treat signed token expiry as part of the issuer's promise. If Alice clearly
  promises a token is valid only until a stated time, redemption after that time
  is neutral expiry evidence about Alice and may be negative evidence about a
  peer who presented the expired token as useful.
- Do not add a central directory, service registry, exchange authority, global
  trust score, global price oracle, shared token-status ledger, permission
  authority, authorization authority, or conformance authority.

## Subtasks

- [x] vopus.1 Run a DF round for POC9 implementation paths, package/function
  names, runtime paths, payload variant names, sparse mesh shape, ordinary-probe
  semantics, and completion markers before code edits. Done under `DI-sipuz`.
- [x] vopus.2 Define the one-pCID POC9 payload vocabulary:
  `need_advertisement`, `offer_promise`, `counter_promise`,
  `acceptance_promise`, `token_issue_promise`, `token_redemption_promise`,
  `outcome_observation`, `referral_promise`, `introduction_promise`, and
  `route_promise`. Done under `DI-sipuz`.
- [x] vopus.3 Add peer bootstrap behavior from existing contacts without
  introducing a global directory or hidden service registry. Done under
  `DI-sipuz`.
- [x] vopus.4 Implement scoped local trust thresholds for public gossip, storage,
  compute, relay, exchange, introduction promises, and transport-path refusal.
  Done under `DI-sipuz`.
- [x] vopus.5 Add ordinary low-risk promises that let Alice build trust before
  sending sensitive data or trusting returned computation results. Done under
  `DI-sipuz`.
- [x] vopus.6 Add referrals and introductions where the receiving agent records
  the referrer's promise separately from any promise by the referred agent. Done
  under `DI-sipuz`.
- [x] vopus.7 Add dynamic route promises so agents can reach non-neighbors
  through locally trusted relay promises without making the kernel a router or
  authority. Done under `DI-sipuz`.
- [x] vopus.8 Add repeated public/private and expired-token misuse rounds where
  agents revise local strategy from promise keep/break history instead of
  following a single fixed transaction script. Done under `DI-sipuz`;
  corrected under `DI-vujil`.
- [x] vopus.9 Record TCP connection outcomes as transport-promise evidence
  without treating an open TCP connection as trust. Done under `DI-sipuz`.
- [x] vopus.10 Validate the POC with deterministic tests, Docker output, and
  guide-resource updates that label the result as evidence only. Done under
  `DI-sipuz`.
- [x] vopus.11 Correct POC9 token semantics so explicitly expired Alice tokens
  do not lower Dave's local trust in Alice; only Mallory's local trust decreases
  when Mallory presents expired token bytes as useful. Done under `DI-vujil`.

## Decision Intent Log

ID: DI-vorus
Date: 2026-05-31 18:46:26
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Fold POC8 findings into the capability-token exchange scenario and
plan POC9 as a peer-discovery and durable-strategy successor. POC9 should test
how agents discover peers, offers, route promises, and trial promises through
existing relationships, referrals, signed promise evidence, and low-risk probes.
Promise relationships should influence TCP connection choices, and TCP outcomes
should feed local promise-accounting evidence, but TCP connectivity must not be
treated as trust.
Intent: POC8 demonstrates a better autonomous promise-economy vocabulary than
POC7, but it still relies on a bounded deterministic world with fixed role
plans, known peers, static topology, seeded stale-token history, and a harness
completion marker. The next evidence step should show agents discovering useful
relationships and escalating trust scope from local observations instead of
depending on a central registry, global trust authority, hidden script, or
service-discovery shortcut.
Constraints: Keep POC9 planned until a later DF locks implementation paths,
runtime paths, names, and payload details. Preserve POC8 as historical evidence.
Do not implement POC9 in this step. Do not add a central directory, service
registry, exchange authority, global trust score, global price oracle, shared
token-status ledger, permission authority, authorization authority, or
conformance authority. Keep all trust and valuation local to the observing
agent.
Affects:
`scenarios/promise-economy-capability-token-exchange/promise-economy-capability-token-exchange.md`;
`protocols/wire-lab.d/TODO/TODO-rajig-promise-economy-spectrum.md`;
`protocols/wire-lab.d/TODO/TODO-vopus-poc9-peer-discovery-strategy.md`;
`protocols/wire-lab.d/TODO/TODO.md`; `DEV-GUIDE-RESOURCES.md`.

ID: DI-sipuz
Date: 2026-06-01 05:11:42
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Implement `implementations/poc9-peer-discovery-strategy/` as a
sibling successor to POC8. POC9 uses seven actors on a deterministic sparse
mesh, one pCID-selected signed CBOR `grid([42(pCID), payload, proof])`
discovery/economy protocol, ordinary low-risk promises instead of a special
probe kind, offer/counter/acceptance promises instead of a separate
exchange-rate quote kind, route/referral/introduction promises, malformed TCP
evidence, local transport-path refusal, and all-node done markers under
`/run/poc9/<run_id>/<node>.done`.
Intent: POC8 proved autonomous promise-economy vocabulary but still assumed
known peers, static route tables, and a bounded harness completion marker.
POC9 should test the next pressure case: Alice discovers useful peers through
signed local promises, weighs referrals as evidence rather than authority, uses
ordinary public/known promises before private or higher-value work, records TCP
behavior as local evidence, and lets local trust changes affect later route,
storage, compute, and stale-token choices without creating a central directory,
global trust score, service registry, exchange authority, permission authority,
authorization authority, or conformance authority.
Constraints: POC9 is executable evidence only, not a final PromiseGrid discovery
protocol, routing protocol, trust API, token standard, economics model, kernel
API, storage API, compute API, transport standard, or SDK. Preserve POC7 and
POC8 unchanged. Keep all trust and valuation local to the observing agent. Use
in-run memory only; do not add restart persistence. Direct TCP is allowed only
to approved sparse-mesh neighbors.
Affects: `implementations/poc9-peer-discovery-strategy/**`;
`implementations/README.md`; `DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-vopus-poc9-peer-discovery-strategy.md`;
`protocols/wire-lab.d/TODO/TODO.md`.

ID: DI-vujil
Date: 2026-06-01 15:40:39
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Correct POC9's Mallory-to-Dave token pressure from stale/revoked-token
semantics to signed-expiry semantics. Alice's token carries an explicit
`expires_at_unix` promise. Redemption after that expiry returns `expired`, keeps
Dave's local trust in Alice neutral, and lowers only Dave's local trust in
Mallory when Mallory presents the expired token bytes as useful.
Intent: A token whose issuer clearly promised a TTL has not been broken when it
expires. Treating that outcome as a broken Alice promise misapplies Promise
Theory by blaming the issuer for keeping the stated expiry promise. The useful
POC9 pressure is Dave's local judgment about Mallory's circulation behavior, not
a false Alice trust penalty.
Constraints: Preserve explicit revocation as a separate broken-promise case in
token tests. Do not add wall-clock dependence; use deterministic logical time for
the bounded POC. Preserve one pCID and signed CBOR grid envelopes. Keep POC9
evidence-only, not a final token expiry standard.
Affects: `implementations/poc9-peer-discovery-strategy/**`;
`implementations/README.md`; `DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-vopus-poc9-peer-discovery-strategy.md`;
`protocols/wire-lab.d/TODO/TODO.md`.
Supersedes: `DI-sipuz` for POC9 stale/revoked token semantics only.
