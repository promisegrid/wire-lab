# poc3

`poc3` is PromiseGrid proof of concept 3. It follows `poc2` and is not a final
PromiseGrid API. Source: `DI-vosof`; `DI-fubir`; `DI-horak`.

The demo tests one kernel plus three apps per node. Each app uses the same local
kernel boundary but its own app pCID. The kernel accepts local receive promises,
routes exact envelope bytes, refuses locally unsupported pCIDs, and writes its
own evidence. Apps make their own local promise judgments.

## Demo

```sh
cd implementations/poc3
docker compose up --build --abort-on-container-exit
```

Expected path:

```text
Alice hello app -> Alice kernel -> Bob kernel -> Bob hello app
Alice echo app -> Alice kernel -> Bob kernel -> Bob echo app -> Bob kernel -> Alice kernel -> Alice echo app
Alice signed app -> Alice kernel -> Bob kernel -> Bob signed app
```

## Commands

- `poc3-kernel`
- `poc3-hello`
- `poc3-echo`
- `poc3-signed`

Each app has an importable package plus a small command wrapper. The apps use
per-app pCIDs; receive promises use the kernel-local pCID.

## Non-promises

`poc3` does not promise production identity, production signatures, storage,
consensus, namespace stability, final SDK shape, or global trust. The signed app
only verifies exact bytes under a POC Ed25519 key.
