# poc11-adaptive-trust-tcp

`poc11-adaptive-trust-tcp` is executable POC evidence for live-agent autonomy,
multi-round relationship decay/repair, reciprocal economics, adversarial
malformed-decision tests, and adaptive direct TCP links. It builds on POC10 but
removes POC10's multi-action surface. Source: `DI-hotos`.

## What This Tests

- Twelve agents run across eight containers.
- Each container runs one supervisor process and one or more independent agent
  kernel processes.
- Each agent keeps its own local relationship ledger, budget, capacity, recent
  evidence, and direct TCP link choices.
- Strong local trust can create or preserve a direct TCP peer relationship.
  Broken or malformed promise evidence can remove that direct link.
- Live decisions come from an OpenAI-compatible Responses API call with a
  strict structured-output schema, followed by Go-owned validation and one
  bounded repair attempt for common missing-shape errors.
- Tests use fake decisions and do not require network or API keys.
- A monitor LLM evaluates completed logs after the run. The monitor is an
  observer only; it never controls agent decisions.

## Protocol Shape

POC11 keeps one pCID for the whole adaptive-trust protocol. Every message uses
one top-level semantic act:

```text
promise
```

Repair, refusal, observation, route preference, economics, and link preference
are promise payload semantics under that pCID. They are not separate action
kinds. Go code, not the LLM, encodes and signs
`grid([42(pCID), payload, proof])` as CBOR. The LLM chooses promise text and
pCID-owned payload fields, but it never writes raw CBOR or proof bytes. Source:
`DI-hotos`.

## Agent Graph

The example run uses this initial sparse graph:

```text
alice -- bob -- carol -- dave -- ellen -- frank -- mallory -- oscar -- ivan
  |       |       |        |       |        |          |        |       |
ellen   frank   grace    heidi   alice    bob        judy     mallory grace

grace -- judy -- heidi -- alice
```

Container placement:

- `alice`: Alice.
- `bob`: Bob.
- `carol`: Carol.
- `dave`: Dave and the observer-only monitor role.
- `ellen-frank`: Ellen and Frank.
- `grace-heidi`: Grace and Heidi.
- `ivan-judy`: Ivan and Judy.
- `mallory-oscar`: Mallory and Oscar.

## Config

Copy the committed template:

```sh
cp config.example.json config.json
```

`config.json` is ignored by Git. It may contain model names, limits, topology,
and the name of the API-key environment variable. It must not contain the API
key itself. The loader rejects fields such as `api_key`, `secret`, `token`,
`bearer_token`, `access_token`, and `auth_token`. Source: `DI-hotos`.

## Run

Set the API key in the environment variable named by `config.json`:

```sh
printf '%s' "$OPENAI_API_KEY" > openai_api_key.txt
chmod 600 openai_api_key.txt
docker compose up --build --abort-on-container-exit
docker compose down --volumes
```

`openai_api_key.txt` is ignored by Git and mounted as a Compose secret. The live
client reads `OPENAI_API_KEY_FILE` inside the container so `docker compose
config` does not render the key value. Source: `DI-mudar`.

Runtime logs are written to the shared Docker volume under
`/run/poc11/<run_id>/` and mirrored to stdout as JSONL.

Summarize one completed run from inside the Compose volume:

```sh
docker compose run --rm --entrypoint /usr/local/bin/poc11-analyze dave /run/poc11/poc11-demo
```

## First Live Run

The first live Docker run completed successfully: all twelve agents wrote done
markers and Dave wrote the observer-only monitor report. The monitor scored the
run as Promise Theory fit `4/5`, autonomy `5/5`, protocol validity `3/5`, local
trust correctness `5/5`, and imposition avoidance `5/5`. Source: `DI-horuh`.

The useful transactions were differentiated and self-scoped: Alice/Bob used
low-risk storage/trust-building promises; Bob/Carol exchanged storage and
compute willingness; Dave made evidence-gathering promises; Ellen made
observation/introduction promises limited to what she observed; Frank made
bounded relay-like promises without claiming delivery control; Grace, Judy,
Heidi, and Ivan exercised reciprocal-value, repair-oriented, and local
observation promises; Mallory and Oscar created adversarial pressure and then
sent self-scoped promise-only messages. Source: `DI-horuh`.

The run also exposed production-readiness gaps: rejected missing-target,
non-direct-target, and self-target decisions; one invalid JSON decision; repeated
near-duplicate promises; some receipts after `node_done`; and normal listener
shutdown logged as `accept_failed`. These are protocol hygiene and lifecycle
issues, not reasons to reintroduce top-level RPC-like action kinds. Source:
`DI-horuh`; `DI-mosoj`.

## Hardening After First Live Run

The next hardening pass addresses those protocol-hygiene gaps without changing
the one-action Promise Theory boundary. Normal listener close is now clean
`listener_closed` evidence; each node stops accepting new frames, drains active
receive handlers for a bounded interval, persists its local relationship
snapshot, and then writes idempotent `node_done` evidence. The monitor also
treats `monitor.done` as idempotent and drains in-flight receipts before reading
logs. Source: `DI-duhub`.

Live provider calls now request structured JSON: decisions return `act`,
`target`, `promise`, `reason`, and a strict key/value `fields` list; monitor
reports return the five bounded scores plus summary and concerns. The runtime
still validates every decision locally, rejects authority/RPC/prompt-injection
wording, and allows only a single bounded repair for missing `act`, missing
single direct-peer `target`, or missing promise text. Source: `DI-duhub`.

Resource promises now have executable POC checks. Storage or compute fulfillment
promises must declare positive units that fit local budget and capacity before
they are sent. Inbound resource promises with malformed or extreme unit claims
are rejected as local evidence. Broken promises with stake or collateral fields
spend local budget units; no central penalty authority is introduced. Source:
`DI-duhub`.

## Targeting And Discovery Follow-up

The hifud.16-hifud.24 pass tightens live targeting and adds a real candidate-link
path while preserving the same single top-level `promise` act. Provider-side
structured output now constrains `target` to locally visible direct or candidate
peer names. Runtime validation still rejects ordinary candidate traffic: a
candidate peer is valid only when the promise payload says
`promise_about=link_discovery`. Bundled targets such as `bob,ellen` are repaired
by narrowing to the first locally valid target; unresolved ambiguous targets are
still rejected. Source: `DI-nanud`.

Shutdown now has a configurable `shutdown_grace_millis` interval. During that
interval, a node writes `turns_done` evidence and keeps its listener open until
all agents have finished active turns, or until the grace interval expires, so
lagging peers can finish already-planned sends before listeners close. The
monitor no longer emits duplicate `inflight_drained` evidence. Source:
`DI-nanud`.

The follow-up live run completed with all containers exiting `0`. The analyzer
reported `143` events, `12` `turns_done` markers, `12` `node_done` markers,
`12` clean `shutdown_grace_elapsed` events, `12` `inflight_drained` events, no
`decision_rejected` events, no refused-connection transport failures, `2`
`send_failed` events caused by explicit `not_promised` acknowledgements, `2`
`message_not_promised` events, and `2` `promise_withheld` non-commitments. The
monitor scored Promise Theory fit `4/5`, autonomy `5/5`, protocol validity
`4/5`, local trust correctness `4/5`, and imposition avoidance `4/5`. Source:
`DI-nanud`.

## Current Limits

POC11 is intentionally nondeterministic at runtime. Config limits bound turns,
per-agent calls, monitor calls, and request timeouts, but specific decisions are
provider outputs. Direct link reconfiguration is a local runtime promise about
dialing or accepting TCP from a peer; the POC does not mutate Docker networks.
This is executable evidence, not a final LLM-agent API, trust API, discovery
protocol, monitor standard, economics protocol, SDK, or provider abstraction.
The structured-output field list is a provider-boundary accommodation, not a
new wire format; the CBOR payload remains the pCID-owned string map. Source:
`DI-hotos`; `DI-duhub`; `DI-nanud`.
