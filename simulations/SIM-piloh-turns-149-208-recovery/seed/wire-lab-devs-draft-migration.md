# wire-lab-devs-draft Migration

The `wire-lab-devs-draft` transport evidence moved into this simulation's world
so it can be replayed and evaluated as specimen data without preserving root
`transports/` as an active layout commitment. Source: `DI-fakin`.

| Field | Value |
|---|---|
| Original path | `transports/wire-lab-devs-draft/` |
| New path | `simulations/SIM-piloh-turns-149-208-recovery/world/transports/wire-lab-devs-draft/` |
| Method | `git mv` |
| Source commit | `780f56525a8d528d3d5caf58ab18f9a7f41da892` |
| CID parameters | CIDv1, raw codec, sha2-256 multihash, base32 multibase |
| Verification result | PASS on 2026-05-10: all four `bafk*.txt` filenames matched raw CIDv1 over file bytes after migration. |

## Verified message CIDs

| Message file | Verification |
|---|---|
| `bafkreia46vxsahmeicugfxmc7natorkstc3mdaz4r5d3zz46whjwpvqwta.txt` | PASS |
| `bafkreidef4b4qdc4xjvkjrern7jm4ta75q55ed2u2ilwcrkxqhn7n4fjce.txt` | PASS |
| `bafkreihhuejiefrqrm7zgw2jsdqc37lwmbvfkw5uqbnjx3wsobcxh3y7ni.txt` | PASS |
| `bafkreihnonvsf3vmcagukqcxwoh35255eduulvwwx3kax6ty4iidklk5vu.txt` | PASS |

The message files are not edited by this migration. Their body text may mention
old paths such as `transports/draft--wire-lab-devs/`; those references are
historical evidence and are preserved to keep CIDs stable.
