# poc9-peer-discovery-strategy

`poc9-peer-discovery-strategy` is executable POC evidence for
`scenarios/promise-economy-capability-token-exchange/` and TODO-vopus. It is a
sibling successor to POC8, not an in-place rewrite. Source: `DI-sipuz`.

The demo uses seven containers on a deterministic sparse mesh:

```text
Alice:   Bob, Ellen
Bob:     Alice, Carol, Frank
Carol:   Bob, Dave, Mallory
Dave:    Carol, Ellen
Ellen:   Alice, Dave, Frank
Frank:   Bob, Ellen, Mallory
Mallory: Carol, Frank
```

The approved mesh edges are `Alice-Bob`, `Alice-Ellen`, `Bob-Carol`,
`Bob-Frank`, `Carol-Dave`, `Carol-Mallory`, `Dave-Ellen`, `Ellen-Frank`, and
`Frank-Mallory`. Direct TCP is allowed only to mesh neighbors; non-neighbor
communication requires locally learned route promises.

Each container runs one local kernel boundary plus local promise-economy app
roles inside a single bounded process. Peers use length-framed TCP. Each app
message is a signed CBOR `grid(...)` tag (`0x67726964`, decimal `1735551332`)
wrapping `[42(pCID), payload, proof]`.

POC9 uses **one pCID for the whole discovery/economy protocol**. The payload
field `kind` selects protocol-defined variants under that one pCID:
`need_advertisement`, `offer_promise`, `counter_promise`,
`acceptance_promise`, `token_issue_promise`, `token_redemption_promise`,
`outcome_observation`, `referral_promise`, `introduction_promise`, and
`route_promise`. These are payload kinds, not separate pCIDs. There is no
special `probe_promise` kind and no `exchange_rate_quote` kind. Low-risk probes
are ordinary low-value promises selected by local strategy; exchange terms live
inside ordinary offer/counter/acceptance promises. A `route_promise` carries its
promised peer path in `promised_route`; the envelope's `route` field remains the
current transport path for the signed message itself. Source: `DI-sipuz`.

## What Takes Place

- **Neighbors publish introductions and route promises.** Each node promises only
  its own local route offer. Recipients record those route promises as evidence,
  not as a global routing table.
- **Bob refers Alice to Carol.** Bob's `referral_promise` is Bob's promise about
  Bob's local evidence. Alice does not treat it as trust transfer to Carol.
- **Alice starts with ordinary low-risk promises.** Alice asks Bob for public
  storage and Carol for known public compute before sending private data or
  higher-value work.
- **Alice escalates after kept evidence.** After Bob keeps public storage and
  Carol keeps public compute, Alice sends private storage to Bob and higher-value
  compute to Carol.
- **Carol observes malformed transport bytes.** Mallory sends length-framed bytes
  that are not a CBOR grid envelope. Carol records local transport-break evidence
  and later refuses a route promise from Mallory.
- **Mallory routes around local refusal.** Mallory observes Carol's refusal and
  later reaches Dave through Frank and Ellen, proving that refusal is local and
  scoped rather than global banishment.
- **Dave judges expired-token misuse locally.** Dave accepts Mallory's first
  expired Alice token for evidence, redeems it with Alice, observes that Alice's
  signed expiry promise was kept, keeps Alice trust neutral, lowers local trust
  in Mallory for presenting expired bytes as useful, and refuses Mallory's later
  expired-token offer.

## Promise Theory Fit

- Agents only promise their own behavior or the behavior of resources they
  control. No agent makes a promise on behalf of another agent.
- Trust is local and scoped. Bob's referral about Carol is evidence about Bob;
  Carol still earns Alice's trust only through Carol's own kept promises.
- Tokens are issuer promises, not global permission objects. Bearer tokens can be
  presented by holders, but redemption outcomes remain issuer-local and
  observer-local evidence.
- Expiry is part of the issuer's signed promise. If Alice promises a token is
  valid only until a stated time, redemption after that time is expiry evidence,
  not a broken Alice promise. Source: `DI-vujil`.
- TCP is evidence, not authority. Successful frames, failed dials, malformed
  bytes, and refused relay paths influence local relationship accounting without
  making an open TCP connection equivalent to trust.
- POC9 avoids a central directory, service registry, central exchange, global
  trust score, global price oracle, shared token-status ledger, permission
  authority, authorization authority, and conformance authority. Source:
  `DI-sipuz`.

## Current Limits

POC9 is still a bounded deterministic POC. The sparse mesh, scenario timing, and
local strategy thresholds are evidence scaffolding, not a final PromiseGrid
discovery protocol, routing protocol, trust API, token standard, economics model,
kernel API, storage API, compute API, transport standard, or SDK. Source:
`DI-sipuz`.

Run:

```sh
go test ./...
POC9_RUN_ID="$(date -u +%Y%m%d%H%M%S)" docker compose up --build --abort-on-container-exit
docker compose down --volumes
```
