# SIM-kifas: App/kernel surface promises and local observations

This simulation is a Hozif-derived successor. It keeps `SIM-hozif-app-pcid-promises-local-observations` as the base design and mines only two useful ideas from the culled Fozid proposal: explicit per-surface promiser entries and a sharper kept/refused/not-promised/unavailable/broken/timed-out observation vocabulary. Source: `DI-pavub`.

The correction is to keep those ideas inside one stable app/kernel promise-accounting protocol pCID. Related promise and observation records evolve together, so message-level variation belongs in the payload `kind`, not in separate tightly coupled pCIDs. The outer carrier remains the current PromiseGrid envelope direction: CBOR `grid([42(pCID), payload, ...])` with slot roles owned by the protocol named by the Protocol CID.

## Candidate specimen

The specimen under test has two parts:

- the current envelope direction: CBOR `grid([42(pCID), payload, ...])`;
- one `app-kernel-surface-promises` payload protocol whose payload kinds are `surface_promise_v1` and `promise_observation_v1`.

A `surface_promise_v1` record lets one agent promise its own behavior for one local surface: receive, handle, store, compute, send, key-use, lifecycle, namespace-view, reference-view, evidence, or device-effect. A `promise_observation_v1` record lets one observer remember what it locally saw about a promise attempt.

## What this intentionally does not do

- It does not create external authority over agents, global namespace truth, or global trust scores.
- It does not say that a kernel knows what an app can do. The kernel may remember what an app promised this kernel and what the kernel later observed.
- It does not move promise-accounting semantics into the universal envelope layer.
- It does not use Fozid's untagged slot-0 carriage shape; the tag-42 pCID link remains the current envelope direction.
- It does not split closely coupled promise and observation records into separate pCIDs unless a later design demonstrates an independent deployment or layer boundary.

## Good parts mined from Fozid

Fozid showed that first-port work benefits from naming the concrete local surface and promiser instead of hiding everything behind one broad kernel entry. It also showed that kept, refused, not-promised, unavailable, broken, timed-out, unreadable, and corrupt-observation outcomes should not collapse into one failure bucket.

This successor keeps those insights but rewrites them in Promise Theory terms:

- each listed surface has a specific promiser;
- omitted surfaces mean no promise;
- assumptions are not silently upgraded into promises;
- local next-step effects are observer-local notes, not portable rules;
- evidence references are exact bytes, CIDs, bounded local notes, or witness records, not authority facts.

## Scenario pressure

This sim should be scored against:

- `kernel-porting-boundary`: an app promises a local kernel which pCID-selected messages it receives or handles, while the kernel promises only its own delivery and observation behavior.
- `promise-accounting-records-kept-storage-promise`: Alice sends data to Bob only after Bob has made bounded storage promises and Alice's local trust threshold is met.
- `promise-accounting-records-refused-service`: Carol asks Alice for data while making reciprocal promises; Alice records honest refusal or non-promise distinctly from silence or breakage.
- `cas-object-model-dag-cbor-interop`: exact-byte and CID-rooted evidence remains compatible with CBOR/DAG-CBOR without requiring every small device to run a full IPLD stack.
- `l6-cas-starting-profile-bakeoff-long-horizon-reprofile`: older promise and observation records remain meaningful as pCID specs and storage profiles evolve.
