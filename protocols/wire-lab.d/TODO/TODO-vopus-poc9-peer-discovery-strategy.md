# TODO-vopus: POC9 peer discovery strategy

## Status

Planned. Owns the next proof-of-concept pressure after POC8: agents should
discover peers, offers, routes, and trial promises through local relationships
and signed promise evidence rather than through fixed role scripts, static peer
knowledge, or a central directory. Source: `DI-vorus`.

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
- Use low-risk probe promises to build relationship evidence before higher-risk
  storage, compute, relay, or bearer-token exchange.
- Treat transport behavior as local evidence. Promise relationships influence
  whether an agent opens, keeps, retries, or uses TCP; TCP delivery, malformed
  bytes, disconnects, timeouts, retries, and latency feed local
  promise-accounting records.
- Do not add a central directory, service registry, exchange authority, global
  trust score, global price oracle, shared token-status ledger, permission
  authority, authorization authority, or conformance authority.

## Subtasks

- [ ] vopus.1 Run a DF round for POC9 implementation paths, package/function
  names, runtime paths, and payload variant names before code edits.
- [ ] vopus.2 Define the provisional one-pCID discovery payload vocabulary:
  `need_advertisement`, `offer_promise`, `referral_promise`,
  `introduction_promise`, `route_promise`, `probe_promise`, and
  `outcome_observation`.
- [ ] vopus.3 Add peer bootstrap behavior from existing contacts without
  introducing a global directory or hidden service registry.
- [ ] vopus.4 Implement scoped local trust thresholds for public gossip,
  storage, compute, relay, exchange, and introduction promises.
- [ ] vopus.5 Add low-risk probe promises that let Alice build trust before
  sending sensitive data or trusting returned computation results.
- [ ] vopus.6 Add referrals and introductions where the receiving agent records
  the referrer's promise separately from any signed offer by the referred agent.
- [ ] vopus.7 Add dynamic route promises so agents can reach non-neighbors
  through locally trusted relay promises without making the kernel a router or
  authority.
- [ ] vopus.8 Add repeated rounds where agents revise local strategy from
  promise keep/break history instead of following a fixed scenario script.
- [ ] vopus.9 Record TCP connection outcomes as transport-promise evidence
  without treating an open TCP connection as trust.
- [ ] vopus.10 Validate the POC with deterministic tests, Docker output, and
  guide-resource updates that label the result as evidence only.

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
