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
- Agents exchange signed envelopes over length-framed TCP between Docker
  containers; the analyzer requires both sent and received TCP evidence.
- Containers write readiness and done evidence under the Docker run volume so
  startup waits for peer readiness and shutdown waits for local runtime
  quiescence instead of fixed sleeps.
- Agents record local evidence for real CAS storage, token-scoped serving,
  retrieval, replication, corrupt bytes, explicit context objects, dynamic
  function execution, compute results, and cache keys.
- Capability tokens, credit-style economics, trust updates, and voluntary repair
  promises are modeled as local pCID-owned promise payloads and evidence, not as
  global authority.
- LLM decisions are attempted only when `live_decisions` is true and the named
  API-key environment variable is present; config files must not store keys.
- Live provider runs must emit meaningful provider text. The analyzer rejects
  placeholder live-decision evidence such as `provider returned no output_text`
  because that proves reachability, not useful local promise judgment.

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

The repo-local clean regression runner resets the Docker volume before running
and analyzing:

```sh
implementations/poc13-cas-compute-functions/scripts/run-clean.sh
```

If no provider key is available, agents use explicit local fallback decisions
and still emit deterministic protocol evidence.
