# BGP-Class Routing App Scenarios

These scenarios test turn 178's BGP-class application pressure independently of
the generic promise-accounting simulation. They are not a base protocol decision
and not a final route-selection algorithm. Source: `DI-tibis`.

## Scenario Matrix

| Scenario | Setup | What to test | Decision pressure |
|---|---|---|---|
| Honest reachability promise | Alice advertises reachability to Carol through Bob, and Bob later forwards as promised. | What Alice, Bob, and Carol each record locally after the path works. | Route-like promises need observable kept/broken outcomes without a global route authority. |
| Route hijack | Mallory advertises a short attractive path to Carol but cannot deliver traffic or chunks. | How peers detect failed promises and locally downgrade future route choices. | PromiseGrid routing apps need hijack costs that are local and relationship-specific. |
| Route leak | Bob repeats Alice's restricted transit offer outside the intended policy context. | Whether policy scope and onward-restraint promises can be represented and audited. | Routing policy needs more than reachability; it needs promise scope and violation evidence. |
| Conflicting policies | Alice prefers paths that avoid a jurisdiction; Carol prefers cheapest path; Bob has both offers. | Whether route choice can be policy-relative instead of globally best. | A PromiseGrid routing app should support peer-specific preference and refusal semantics. |
| Partition and stale path | A formerly good path becomes unavailable during intermittent connectivity. | How stale promises, timeouts, and withdrawal notices affect local decisions. | Long-lived routing records need aging and repair without central convergence machinery. |
| Sparse knowledge | Alice knows only Bob and Carol; Bob knows Dave and Ellen; no peer has the whole graph. | Whether multi-hop discovery can find acceptable paths without requiring full topology replication. | BGP-class pressure must compose with sparse-CAS and sparse relationship knowledge. |

## Expected Outputs

- A list of route-app promise concepts: reachability, transit, withdrawal,
  refusal, policy scope, onward-restraint, stale path, and violation evidence.
- A boundary note for future specs: BGP-class routing is an L7 application
  pressure case, not a reason to centralize base PromiseGrid accounting.
- Pressure cases for later promise-accounting and conditional-release sims.
