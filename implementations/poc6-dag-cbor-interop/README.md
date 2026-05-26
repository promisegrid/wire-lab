# poc6-dag-cbor-interop

`poc6-dag-cbor-interop` is executable POC evidence for
`scenarios/cas-object-model-dag-cbor-interop/`.

The POC tests whether a PromiseGrid CAS-facing object can use real
DAG-CBOR/IPLD tooling for CID links, byte strings, tag-42 link encoding, stable
CID calculation, and local promise evidence without requiring an IPFS daemon.

It is not a final PromiseGrid CAS API.

Source: `DI-sagos`.

## Run

```sh
go test ./...
```
