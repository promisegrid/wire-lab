# PromiseGrid Wire Lab

A simulation harness for discovering the wire formats, ingress models, and trust mechanics that let [PromiseGrid](https://github.com/promisegrid/promisegrid) survive and evolve as an open, decentralized community of free agents and humans across multiple human generations.

The harness exists to **discover** the right design, not to validate a predetermined one. Every choice is an experimental knob the simulator can change between runs.

## Where to read

- **[`specs/harness-spec-draft.md`](specs/harness-spec-draft.md)** — the canonical spec. This is the single document the simulator, the production code generators, and human reviewers all key off. Its content hash is its pCID.
- **[`docs/thought-experiments/`](docs/thought-experiments/)** — one file per thought experiment, each content-addressable and standing alone. The harness-spec links into them; the [index](docs/thought-experiments/README.md) lists them in chronological order.

## Branch layout

At the moment, the `ppx/main` branch is the active development branch, and the `main` branch is a review branch. The `ppx/main` branch contains the latest code and documentation.

## Status

Provisional. Almost everything in the harness-spec is an experimental knob, not a commitment. The single load-bearing structural decision is that the canonical pointer to the harness-spec is whatever pCID Steve has most recently signed a `merge-harness-spec` promise for — the lock is the key, not a document.

## How to propose a change

1. Branch off `main`.
2. Edit `specs/harness-spec-draft.md` and/or add a new TE file under `docs/thought-experiments/` named `TE-YYYYMMDD-HHMMSS-some-phrase.md` (the timestamp is when the TE was first written; the slug is a kebab-case rendering of its title).
3. Open a PR. The PR body is the proposal; reviewers comment in-thread.
4. Steve is the only signer who can merge to `main`. Branch protection enforces this.

## License

GPL-3.0, matching the rest of [PromiseGrid](https://github.com/promisegrid/promisegrid).
