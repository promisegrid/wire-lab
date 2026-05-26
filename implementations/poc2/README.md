# poc2

`poc2` is PromiseGrid proof of concept 2. It is separate from the external
`grid-poc` repo and is not a final PromiseGrid API. Source: `DI-ratij`;
`DI-tijat`.

The demo proves one narrow claim: app/kernel and kernel/kernel boundaries can use
pCID-selected `grid([42(pCID), payload, ...])` messages without making the kernel
an RPC authority.

## Demo

```sh
cd implementations/poc2
docker compose up --build --abort-on-container-exit
```

Expected path:

```text
Alice app -> Alice kernel -> Bob kernel -> Bob app
```

Both kernels write newline-delimited evidence records. Bob's app prints the hello
text after receiving it through Bob's local kernel.

## Commands

The image installs the same compiled Go binary under two names:

- `poc2-kernel`
- `poc2-hello`

`poc2-kernel` flags:

```text
--node alice|bob
--app-listen 127.0.0.1:<port>
--peer-listen 0.0.0.0:<port>
--peer bob:<port>
--evidence ./run/evidence.jsonl
```

`poc2-hello` flags:

```text
--node alice|bob
--kernel 127.0.0.1:<port>
--mode send|receive
--to bob
--text "hello from Alice"
```

## Non-promises

`poc2` does not promise production identity, production signatures, storage,
consensus, namespace stability, or final SDK shape. It records exact message bytes
as evidence but does not claim secure peer identity.
