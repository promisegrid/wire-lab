# Chunk Feed Replication Scenarios

These scenarios make the turn-177 L5 feed inversion concrete: feeds move CAS
chunks or CAS objects between sparse sites, while L7 group semantics remain
above CAS. They are simulation inputs for TODO-kituj / TE-43 and TODO-pipus,
not a frozen feed wire format. Source: `DI-pator`.

## Scenario Matrix

| Scenario | Setup | What to test | Decision pressure |
|---|---|---|---|
| Sparse advertisement | Alice has a subset of chunks for a Merkle root; Bob has a different subset. | Whether the feed advertises leaves, roots, pointer objects, frontiers, or compact summaries without assuming full replication. | Feed specs must work when no site has all CAS objects. |
| Pull decision | Bob receives an advertisement for chunk C and has peer-local promise accounting records about Alice. | Which inputs decide whether Bob pulls, delays, refuses, or asks another peer. | The "decides" step needs an explicit cross-layer interface instead of a hand wave. |
| Multi-repo sparse site | Alice's site state, Bob's site state, and a large shared corpus live in separate repos or fixtures. | Whether feed promises and CAS object references remain meaningful when the harness orchestrates multiple storage roots. | Turn 178's multi-repo question should be explored without assuming one repo contains every site's content. |
| Partial Merkle fetch | Bob wants root R but only some children are locally available. | Whether the feed can request missing children without understanding group-session message semantics. | L5 should remain meaning-oblivious while still serving L6 CAS repair. |
| Corrupt chunk | Mallory advertises or sends bytes whose hash does not match CID C. | Whether rejection, accounting, and retry behavior are local enough to avoid central enforcement. | Feed behavior must compose with CAS hash verification and peer-local accounting records. |
| Duplicate advertisement | Alice and Carol both advertise chunk C. | Whether duplicate offers are harmless and how Bob chooses between peers. | Promise accounting can influence peer choice without making the feed a central reputation service. |
| Refusal or non-response | Alice refuses to send C or never answers. | Whether refusal is a normal observed outcome that can feed future local decisions. | The feed protocol needs space for refusal and timeout outcomes without treating every miss as corruption. |
| Carrier independence | The same chunk exchange is attempted over UDP, git, libp2p, IPFS, or ATPROTO-adjacent carriers. | Which semantics belong to the feed role and which are carrier mechanics. | The simulation should preserve turn-177's claim that feeds move chunks independent of substrate. |

## Expected Outputs

- A candidate list of L5 feed observations that can be recorded without group
  semantics leaking into the feed layer.
- A TODO-pipus migration constraint: successor group-session specimens should
  depend on chunk replication, not message-file transport.
- A TE-43 constraint: CAS object and chunking decisions must be usable by
  sparse feed replication.
