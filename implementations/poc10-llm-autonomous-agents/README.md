# poc10-llm-autonomous-agents

`poc10-llm-autonomous-agents` is executable POC evidence for replacing scripted
agent strategy with LLM-directed local decisions while preserving PromiseGrid
wire mechanics. It is a sibling successor to POC9, not an in-place rewrite.
Source: `DI-pijan`.

## What This Tests

- Seven agents run in seven containers on a sparse mesh.
- Each agent has a local persona, motivation, neighbor list, trust ledger, and
  bounded action surface.
- Runtime decisions come from a live OpenAI-compatible Responses API call.
- Tests use fake LLM decisions and do not require network or API keys.
- A monitor LLM evaluates the completed logs after the run. The monitor is an
  observer only; it never controls agent decisions.

## Protocol Shape

POC10 keeps one pCID for the whole autonomous-agent protocol. Message kinds are
payload variants under that pCID, not separate pCIDs. Go code, not the LLM,
encodes and signs `grid([42(pCID), payload, proof])` as CBOR. The LLM may choose
semantic decisions, payload field values, or freeform promise text, but it never
writes raw CBOR or signatures. Source: `DI-pijan`.

## Autonomy Profiles

- **Structured action:** the LLM chooses an allowed action and target; Go fills
  protocol-safe fields.
- **Structured payload:** the LLM fills typed payload fields; Go validates them
  before signing.
- **Freeform intent:** the LLM writes promise/refusal/offer text; Go wraps it as
  a typed payload or rejects it.

Default profiles:

- Alice: privacy-sensitive data holder, structured action.
- Bob: storage provider, structured payload.
- Carol: compute provider, structured payload.
- Dave: skeptical evidence accountant and monitor node, structured action.
- Ellen: relationship broker, structured payload.
- Frank: relay helper, structured action.
- Mallory: opportunistic peer, freeform intent.

## Config

Copy the committed template:

```sh
cp config.example.json config.json
```

`config.json` is ignored by Git. It may contain model names, limits, topology,
and the name of the API-key environment variable. It must not contain the API
key itself. The loader rejects fields such as `api_key`, `secret`, `token`,
`bearer_token`, `access_token`, and `auth_token`. Source: `DI-pijan`.

## Run

Set the API key in the environment variable named by `config.json`:

```sh
export OPENAI_API_KEY=...
docker compose up --build --abort-on-container-exit
docker compose down --volumes
```

Runtime logs are written to the shared Docker volume under
`/run/poc10/<run_id>/` and mirrored to stdout as JSONL.

## Current Limits

POC10 is intentionally nondeterministic at runtime. Config limits bound turns,
agent calls, monitor calls, and request timeouts, but the specific decisions are
provider outputs. This is executable evidence, not a final LLM-agent API, trust
API, discovery protocol, monitor standard, economics protocol, SDK, or provider
abstraction. Source: `DI-pijan`.
