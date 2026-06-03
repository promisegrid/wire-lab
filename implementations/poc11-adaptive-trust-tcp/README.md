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
- Live decisions come from an OpenAI-compatible Responses API call.
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

## Current Limits

POC11 is intentionally nondeterministic at runtime. Config limits bound turns,
per-agent calls, monitor calls, and request timeouts, but specific decisions are
provider outputs. Direct link reconfiguration is a local runtime promise about
dialing or accepting TCP from a peer; the POC does not mutate Docker networks.
This is executable evidence, not a final LLM-agent API, trust API, discovery
protocol, monitor standard, economics protocol, SDK, or provider abstraction.
Source: `DI-hotos`.
