# POC13 CAS storage and CID-named compute functions

POC13 is executable evidence for decentralized CAS storage promises and
CID-named compute promises. It is not a stable PromiseGrid storage, compute,
cache, app, or kernel API. Source: `DI-bibom`; `DI-notig`.

## What it tests

- `cas_storage_v1` is the storage protocol pCID; content CIDs live in payloads.
- `cid_compute_v1` is the compute protocol pCID; `function_cid`, input CIDs,
  context CIDs, and result CIDs live in payloads.
- Every message keeps one top-level act, `promise`, inside
  `grid([42(pCID), payload, proof])`.
- Agents record local evidence for storage, retention, serving, replication,
  corrupt bytes, explicit context objects, compute results, and cache keys.
- LLM decisions are attempted only when `live_decisions` is true and the named
  API-key environment variable is present; config files must not store keys.

## Run

```sh
cp config.example.json config.json
docker compose up --build
docker compose run --rm --entrypoint /usr/local/bin/poc13-analyze alice-bob /run/poc13/poc13-demo
```

To run without creating local config, point Compose at the committed example:

```sh
POC13_CONFIG=./config.example.json docker compose up --build
POC13_CONFIG=./config.example.json docker compose run --rm --entrypoint /usr/local/bin/poc13-analyze alice-bob /run/poc13/poc13-demo
```

If no provider key is available, agents use explicit local fallback decisions
and still emit deterministic protocol evidence.
